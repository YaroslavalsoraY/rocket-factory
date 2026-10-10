#!/usr/bin/env python3
"""Run the full local service chain with a fake Telegram API."""

import argparse
import json
import os
from pathlib import Path
import re
import shutil
import signal
import socket
import subprocess
import sys
import tempfile
import threading
import time
import urllib.error
import urllib.request
import uuid
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

ROOT = Path(__file__).resolve().parents[2]
COMPOSE = ROOT / "deploy" / "compose"
SERVICES = ("iam", "payment", "inventory", "order", "assembly", "notification")
CONTAINERS = {
    "core": ("kafka",),
    "iam": ("postgres-iam", "redis-iam"),
    "order": ("order-postgres",),
    "inventory": ("mongo-inventory",),
}
PORTS = (50053, 50052, 50051, 8080)


def run(command, *, input_text=None, check=True, timeout=120):
    result = subprocess.run(
        command, cwd=ROOT, input=input_text, text=True, capture_output=True, timeout=timeout
    )
    if check and result.returncode:
        detail = (result.stderr or result.stdout).strip()
        detail = re.sub(r"://[^/@\s]+:[^/@\s]+@", "://***:***@", detail)
        raise RuntimeError(f"{' '.join(map(str, command[:3]))} failed: {detail[-800:]}")
    return result


def compose(group, *args):
    folder = COMPOSE / group
    return run(
        ["docker", "compose", "--env-file", str(folder / ".env"),
         "-f", str(folder / "docker-compose.yml"), *args],
        timeout=300,
    )


def container_running(name):
    result = run(["docker", "inspect", "-f", "{{.State.Running}}", name], check=False)
    return result.returncode == 0 and result.stdout.strip() == "true"


def wait_until(label, predicate, timeout=90):
    end = time.monotonic() + timeout
    while time.monotonic() < end:
        if predicate():
            return
        time.sleep(0.5)
    raise RuntimeError(f"Timeout waiting for {label}")


def tcp_open(port):
    try:
        with socket.create_connection(("127.0.0.1", port), timeout=0.4):
            return True
    except OSError:
        return False


def docker_exec(name, script, *, input_text=None):
    return run(["docker", "exec", "-i", name, "sh", "-c", script], input_text=input_text)


class TelegramMock(ThreadingHTTPServer):
    daemon_threads = True

    def __init__(self):
        super().__init__(("127.0.0.1", 0), TelegramHandler)
        self.messages = []
        self.lock = threading.Lock()


class TelegramHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.headers.get("Transfer-Encoding", "").lower() == "chunked":
            chunks = []
            while True:
                size = int(self.rfile.readline().split(b";", 1)[0], 16)
                if size == 0:
                    while self.rfile.readline() not in (b"\r\n", b"\n", b""):
                        pass
                    break
                chunks.append(self.rfile.read(size))
                self.rfile.read(2)
            body = b"".join(chunks)
        else:
            body = self.rfile.read(int(self.headers.get("Content-Length", "0")))
        method = self.path.rsplit("/", 1)[-1]
        if method == "getMe":
            result = {"id": 1, "is_bot": True, "first_name": "Rocket Test"}
        elif method == "getUpdates":
            time.sleep(0.5)
            result = []
        elif method == "sendMessage":
            with self.server.lock:
                self.server.messages.append(body.decode("utf-8", "replace"))
            result = {"message_id": 1, "date": 0,
                      "chat": {"id": 723424205, "type": "private"}, "text": "ok"}
        else:
            result = True
        payload = json.dumps({"ok": True, "result": result}).encode()
        try:
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)
        except (BrokenPipeError, ConnectionResetError):
            pass

    def log_message(self, *_args):
        pass


def grpc(method, data, *, session=None):
    command = [str(ROOT / "bin" / "grpcurl"), "-plaintext"]
    if session:
        command += ["-H", "session-uuid: " + session]
    command += ["-d", json.dumps(data), "localhost:50053" if method.startswith(("user.", "iam.")) else "localhost:50051", method]
    return json.loads(run(command).stdout)


def http(method, path, *, body=None, session=None):
    headers = {"Content-Type": "application/json"}
    if session:
        headers["X-Session-Uuid"] = session
    request = urllib.request.Request(
        "http://localhost:8080/api/v1/orders" + path,
        data=json.dumps(body).encode() if body is not None else None,
        headers=headers,
        method=method,
    )
    try:
        with urllib.request.urlopen(request, timeout=10) as response:
            raw = response.read()
            return response.status, json.loads(raw) if raw else None
    except urllib.error.HTTPError as exc:
        raw = exc.read()
        return exc.code, json.loads(raw) if raw else None


def assert_response(actual, code):
    status, body = actual
    if status != code:
        raise AssertionError(f"Expected HTTP {code}, got {status}: {body}")
    return body


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--keep-infra", action="store_true", help="leave started Docker infrastructure running")
    args = parser.parse_args()
    for name in ("docker", "go"):
        if not shutil.which(name):
            raise RuntimeError(f"{name} is required")
    if not (ROOT / "bin" / "grpcurl").is_file():
        raise RuntimeError("Run task grpcurl:install first")
    for group in CONTAINERS:
        if not (COMPOSE / group / ".env").is_file():
            raise RuntimeError(f"Missing deploy/compose/{group}/.env; run task env:generate")
    for group in ("assembly", "notification", "payment"):
        if not (COMPOSE / group / ".env").is_file():
            raise RuntimeError(f"Missing deploy/compose/{group}/.env; run task env:generate")
    occupied = [port for port in PORTS if tcp_open(port)]
    if occupied:
        raise RuntimeError(f"Application ports already in use: {occupied}")

    preexisting = {group: any(container_running(name) for name in names)
                   for group, names in CONTAINERS.items()}
    started = []
    processes = []
    logs = {}
    state = {}
    mock = None
    temp = tempfile.TemporaryDirectory(prefix="rocket-e2e-")
    try:
        for group, names in CONTAINERS.items():
            print(f"Starting {group} infrastructure...", flush=True)
            compose(group, "up", "-d", *names)
            started.append(group)
        wait_until("Kafka", lambda: run(
            ["docker", "inspect", "-f", "{{.State.Health.Status}}", "kafka"],
            check=False).stdout.strip() == "healthy", 120)
        wait_until("PostgreSQL/Redis/MongoDB ports",
                   lambda: all(tcp_open(p) for p in (5435, 5444, 6333, 27017)), 90)
        wait_until("PostgreSQL and Redis health", lambda: all(
            run(["docker", "inspect", "-f", "{{.State.Health.Status}}", name],
                check=False).stdout.strip() == "healthy"
            for name in ("postgres-iam", "redis-iam", "order-postgres")
        ), 90)
        mongo_ping = ('mongosh --quiet -u "$MONGO_INITDB_ROOT_USERNAME" '
                      '-p "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase admin '
                      '--eval "db.adminCommand({ping:1}).ok"')
        wait_until("MongoDB ping", lambda: run(
            ["docker", "exec", "mongo-inventory", "sh", "-c", mongo_ping],
            check=False, timeout=8).stdout.strip() == "1", 60)
        wait_until("MongoDB healthcheck", lambda: run(
            ["docker", "inspect", "-f", "{{.State.Health.Status}}", "mongo-inventory"],
            check=False).stdout.strip() == "healthy", 90)

        for service in SERVICES:
            run(["go", "build", "-o", str(Path(temp.name) / service),
                 f"./{service}/cmd/main.go"], timeout=180)
        mock = TelegramMock()
        threading.Thread(target=mock.serve_forever, daemon=True).start()
        for service in SERVICES:
            log = open(Path(temp.name) / f"{service}.log", "w+")
            logs[service] = log
            env = os.environ.copy()
            cwd = ROOT
            if service == "inventory":
                cwd = COMPOSE / "inventory"
                env.update(MONGO_HOST="localhost", IAM_GRPC_HOST="localhost",
                           GRPC_HOST="localhost")
            if service == "notification":
                env.update(TELEGRAM_BOT_TOKEN="000000000:LOCAL_TEST",
                           TELEGRAM_API_URL=f"http://127.0.0.1:{mock.server_port}")
            proc = subprocess.Popen(
                [str(Path(temp.name) / service)], cwd=cwd, env=env,
                stdout=log, stderr=subprocess.STDOUT, start_new_session=True,
            )
            processes.append((service, proc))
            if service in ("iam", "payment", "inventory", "order"):
                port = {"iam": 50053, "payment": 50052, "inventory": 50051, "order": 8080}[service]
                wait_until(service, lambda: tcp_open(port) or proc.poll() is not None, 45)
                if proc.poll() is not None:
                    raise RuntimeError(f"{service} exited during startup")
        print("All six Go services started.", flush=True)
        stamp = uuid.uuid4().hex[:12]
        login = f"rocket-e2e-{stamp}@example.test"
        password = "RocketE2E-12345!"
        user = grpc("user.v1.UserService/Register",
                    {"user": {"login": login, "email": login, "password": password}})
        state["user"] = str(uuid.UUID(user["userUuid"]))
        session = grpc("iam.v1.IamService/Login",
                       {"login": login, "password": password})["sessionUuid"]
        state["session"] = session
        identity = grpc("iam.v1.IamService/Whoami", {"session_uuid": session})
        assert identity["session"]["userId"] == state["user"]
        print("IAM register/login/whoami: OK", flush=True)

        part = str(uuid.uuid4())
        state["part"] = part
        mongo_js = (f'db.parts.insertOne({{_id:"{part}",name:"E2E engine",'
                    'desription:"test part",price:12.5,stock_quantity:10,'
                    'category:1,tags:["e2e"],metadata:{},created_at:new Date()})')
        docker_exec("mongo-inventory",
                    'mongosh --quiet -u "$MONGO_INITDB_ROOT_USERNAME" '
                    '-p "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase admin '
                    '"$MONGO_INITDB_DATABASE"', input_text=mongo_js)
        listed = grpc("inventory.v1.InventoryService/ListParts",
                      {"filter": {"uuids": [part]}}, session=session)
        assert listed["parts"][0]["uuid"] == part
        got_part = grpc("inventory.v1.InventoryService/GetPart",
                        {"uuid": part}, session=session)
        assert got_part["part"]["uuid"] == part
        print("Inventory ListParts/GetPart: OK", flush=True)

        assert_response(http("POST", "", body={"part_uuids": [part]}), 401)
        first = assert_response(http("POST", "", body={"part_uuids": [part]},
                                     session=session), 200)
        state["orders"] = [str(uuid.UUID(first["order_uuid"]))]
        order = assert_response(http("GET", "/" + state["orders"][0],
                                     session=session), 200)
        assert order["user_uuid"] == state["user"] and order["status"] == "PENDING_PAYMENT"
        print("Order auth/Create/Get: OK", flush=True)

        second = assert_response(http("POST", "",
                            body={"user_uuid": str(uuid.uuid4()), "part_uuids": [part]},
                            session=session), 200)
        state["orders"].append(str(uuid.UUID(second["order_uuid"])))
        second_order = assert_response(http("GET", "/" + state["orders"][1],
                                            session=session), 200)
        assert second_order["user_uuid"] == state["user"]
        assert_response(http("POST", "/" + state["orders"][1] + "/cancel",
                             session=session), 204)
        cancelled = assert_response(http("GET", "/" + state["orders"][1],
                                           session=session), 200)
        assert cancelled["status"] == "CANCELLED"
        print("Order ignores client user_uuid and cancels: OK", flush=True)

        paid = assert_response(http("POST", "/" + state["orders"][0] + "/pay",
                                    body={"payment_method": "PAYMENT_METHOD_CARD"},
                                    session=session), 200)
        assert uuid.UUID(paid["transaction_uuid"])
        def completed():
            status, result = http("GET", "/" + state["orders"][0], session=session)
            return status == 200 and result["status"] == "COMPLETED"
        wait_until("OrderPaid -> Assembly -> COMPLETED", completed, 25)
        def notified():
            with mock.lock:
                return sum(state["orders"][0] in msg for msg in mock.messages) >= 2
        wait_until("paid and assembled Telegram notifications", notified, 15)
        print("Payment/Kafka/Assembly/Notification mock: OK", flush=True)
        print("Full end-to-end smoke test passed.", flush=True)
    finally:
        if mock:
            mock.shutdown()
            mock.server_close()
        for _, proc in reversed(processes):
            if proc.poll() is None:
                os.killpg(proc.pid, signal.SIGTERM)
        for _, proc in reversed(processes):
            try:
                proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                os.killpg(proc.pid, signal.SIGKILL)
                proc.wait()
        for log in logs.values():
            log.close()
        if "part" in state:
            part = state["part"]
            docker_exec("mongo-inventory",
                'mongosh --quiet -u "$MONGO_INITDB_ROOT_USERNAME" '
                '-p "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase admin '
                '"$MONGO_INITDB_DATABASE"',
                input_text=f'db.parts.deleteOne({{_id:"{part}"}})',)
        if state.get("orders"):
            ids = ",".join(f"'{item}'" for item in state["orders"])
            docker_exec("order-postgres",
                'psql -q -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"',
                input_text=f"DELETE FROM orders WHERE id IN ({ids});")
        if state.get("session"):
            run(["docker", "exec", "redis-iam", "redis-cli", "DEL",
                 "sessions:" + state["session"],
                 "user_sessions:" + state["user"]], check=False)
        if state.get("user"):
            docker_exec("postgres-iam",
                'psql -q -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"',
                input_text=f"DELETE FROM users WHERE id = '{state['user']}';")
        if not args.keep_infra:
            for group in reversed(started):
                if not preexisting[group]:
                    compose(group, "down")
        temp.cleanup()


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(f"E2E failed: {exc}", file=sys.stderr)
        sys.exit(1)
