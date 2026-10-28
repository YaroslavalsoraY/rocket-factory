# Rocket Factory

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/YaroslavalsoraY/0ef8ab849dc6604bff79f92c46b9e7d5/raw/coverage.json)

Много-сервисная платформа для управления заказами, оплатой, уведомлениями, инвентаризацией и авторизацией.

## О проекте

`rocket-factory` реализует архитектуру микросервисов на Go с использованием:

- `go.work` для монорепозитория и модульной разработки
- gRPC и Protocol Buffers для межсервисного API
- PostgreSQL для долгосрочного хранения
- Redis для кеширования и сессий
- Kafka для событийной интеграции
- Docker Compose для локального запуска
- `buf` для проверки и генерации Protobuf

## Структура репозитория

- `assembly/` — сервис агрегации данных и оркестрации
- `iam/` — сервис авторизации и управления пользователями
- `inventory/` — сервис управления складом
- `notification/` — сервис отправки уведомлений
- `order/` — сервис управления заказами
- `payment/` — сервис управления оплатой
- `platform/` — общие утилиты, middleware и библиотеки
- `shared/` — общие Protobuf, OpenAPI и пакеты
- `deploy/` — конфигурация Docker Compose и шаблоны окружения
- `Taskfile.yml` — задачи разработки и автоматизации
- `go.work` — объединяет все модули Go для единого рабочего пространства

## Быстрый старт

### 1. Установка зависимостей

Проект использует локальные бинарники в `./bin`. Перед работой убедитесь, что они установлены:

```bash
task install-formatters
 task install-golangci-lint
 task install-buf
 task proto:install-plugins
```

### 2. Генерация protobuf

```bash
task proto:gen
```

### 3. Локальный запуск сервисов

В директории `deploy/compose/` находятся `docker-compose.yml` для базовых сервисов и каждого модуля. Для старта используйте:

```bash
docker compose -f deploy/compose/core/docker-compose.yml up -d
```

Для запуска конкретного модуля используйте `docker compose -f deploy/compose/<module>/docker-compose.yml up -d`.

> Внимание: текущий репозиторий содержит несколько модулей, поэтому перед запуском убедитесь, что переменные окружения определены через файлы `deploy/env/*.env.template`.

## Разработка

### Структура модулей

Каждый сервис содержит собственный модуль Go с общими слоями:

- `cmd/` — точка входа
- `internal/app/` — инициализация приложения и DI-контейнер
- `internal/api/` — gRPC-обработчики
- `internal/service/` — бизнес-логика
- `internal/repository/` — работа с БД и хранилищем
- `internal/config/` — конфигурация и загрузка окружения

### Команды для разработки

- Форматирование кода:
  ```bash
  task format
  ```
- Линтинг:
  ```bash
  task lint
  ```
- Генерация Protobuf:
  ```bash
  task proto:gen
  ```
- Проверка Protobuf:
  ```bash
  task proto:lint
  ```

### Сборка и запуск модуля

Для запуска сервиса из корня проекта используйте `go run` через модуль:

```bash
go run ./iam/cmd/main.go
```

Для сборки конкретного сервиса:

```bash
go build -o bin/iam ./iam/cmd/main.go
```

## Тестирование

### Запуск тестов для всех модулей

```bash
for mod in assembly iam inventory notification order payment platform shared; do
  echo "=== $mod ==="
  go test ./$mod/...
 done
```

### Запуск тестов для отдельного модуля

```bash
go test ./iam/...
```

### Сквозной тест всех сервисов

Нужны Docker, Go и файлы `deploy/compose/*/.env`. Если файлов окружения нет, создайте их через `task env:generate`. `grpcurl` задача установит автоматически. Затем запустите:

```bash
task test-e2e
```

Тест поднимает Kafka, PostgreSQL, Redis и MongoDB, собирает и запускает все шесть Go-сервисов, проверяет регистрацию, склад, создание, оплату, сборку и отмену заказа. Уведомления отправляются только локальному Telegram mock. Тест удаляет созданные записи и останавливает контейнеры, которые запустил сам; уже работавшую инфраструктуру оставляет запущенной.

## Конфигурация

### Переменные окружения

Каждый сервис загружает конфигурацию через собственные файлы окружения и поддерживает шаблоны в `deploy/env/*.env.template`.

Основные параметры:

- `DATABASE_URL` — подключение к PostgreSQL
- `REDIS_URL` — подключение к Redis
- `KAFKA_BROKERS` — адреса Kafka
- `GRPC_PORT` — порт gRPC сервера
- `LOG_LEVEL` — уровень логирования

### Docker Compose и env

- `deploy/env/core.env.template` — базовые переменные окружения
- `deploy/env/iam.env.template` — конфигурация IAM сервиса
- `deploy/env/inventory.env.template` — конфигурация Inventory сервиса
- `deploy/env/notification.env.template` — конфигурация Notification сервиса
- `deploy/env/order.env.template` — конфигурация Order сервиса
- `deploy/env/payment.env.template` — конфигурация Payment сервиса
- `deploy/env/assembly.env.template` — конфигурация Assembly сервиса

Для генерации окружения используйте утилиту `deploy/env/generate-env.sh`.

После `task up-core` Grafana доступна по адресу http://localhost:3000. Логин и пароль задаются переменными `CORE_GRAFANA_ADMIN_USER` и `CORE_GRAFANA_ADMIN_PASSWORD` в `deploy/env/.env`. Prometheus добавляется как источник данных автоматически.

Задачи `up-*` и `down-*` используют Docker Desktop (`desktop-linux`). Для другого Engine задайте `ROCKET_DOCKER_CONTEXT`, например `ROCKET_DOCKER_CONTEXT=default task up-all`.

## Архитектурные особенности

- Микросервисная архитектура на Go с модульным `go.work`.
- Сервисы общаются по gRPC и через события Kafka.
- `platform/` содержит общие вспомогательные пакеты: логирование, кеш, gRPC middleware, миграции.
- `shared/` содержит общие protobuf-схемы и API-обертки.
- Каждая бизнес-логика реализуется через `service` + `repository` + `api` слои.

## Поддержка качества

Рекомендуется соблюдать следующие практики:

- Код должен быть отформатирован `gofumpt` и импортирован `gci`
- Выполнять `golangci-lint` для всех модулей перед коммитом
- Обновлять protobuf схемы только через `shared/proto` и генерировать код с `buf`
- Избегать глобальных `panic` в инициализации, вместо этого возвращать ошибки и обрабатывать их на уровне приложения
- Использовать `context.Context` для управления временем жизни запросов и операций

## Полезные команды

- Установить инструменты:
  ```bash
  task install-formatters
  task install-golangci-lint
  task install-buf
  ```
- Проверить общий статус линтинга:
  ```bash
  task lint
  ```
- Создать OpenAPI бандл для order:
  ```bash
  task redocly-cli:order-v1-bundle
  ```

## Совет по локальной разработке

Для комфортной разработки:

- Используйте `go.work` для работы с модулями без `replace`
- Запускайте сервисы по одному, проверяя зависимости через `deploy/compose`
- Храните чувствительные данные вне репозитория и загружайте через `.env`
- Документируйте контракт protobuf и OpenAPI вместе с изменением логики

## Контакты и поддержка

Если нужно добавить новый модуль или расширить API, следуйте текущему шаблону:

1. Создайте отдельный модуль в корне
2. Добавьте его в `go.work`
3. Вынесите общие модели в `shared/`
4. Добавьте тесты и линтинг
