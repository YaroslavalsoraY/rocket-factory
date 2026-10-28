#!/bin/bash

# Директория со скриптом
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$SCRIPT_DIR"
COMPOSE_DIR="$SCRIPT_DIR/../compose"

# Используем переданную переменную ENV_SUBST или системный envsubst
if [ -z "$ENV_SUBST" ]; then
  if ! command -v envsubst &> /dev/null; then
    echo "❌ Ошибка: envsubst не найден в системе и не передан через ENV_SUBST!"
    echo "Запустите скрипт через task env:generate"
    exit 1
  fi
  ENV_SUBST=envsubst
fi

# Загружаем основной .env файл
if [ ! -f "$SCRIPT_DIR/.env" ]; then
  echo "Ошибка: Файл $SCRIPT_DIR/.env не найден!"
  exit 1
fi

# Экспортируем все переменные из .env для использования в envsubst
set -a
source "$SCRIPT_DIR/.env"
set +a

# Функция для обработки шаблона и создания .env файла
process_template() {
  local service=$1
  local template="$TEMPLATE_DIR/${service}.env.template"
  local output="$COMPOSE_DIR/${service}/.env"

  if [ "$service" = "core" ]; then
    local required_key
    for required_key in \
      CORE_ELASTICSEARCH_HOST CORE_ELASTICSEARCH_PORT CORE_ELASTICSEARCH_HEAP_SIZE \
      CORE_ELASTICSEARCH_LOGS_INDEX CORE_KIBANA_HOST CORE_KIBANA_PORT \
      CORE_PROMETHEUS_PORT CORE_PROMETHEUS_SCRAPE_INTERVAL \
      CORE_PROMETHEUS_EVALUATION_INTERVAL ORDER_HTTP_PORT ASSEMBLY_METRICS_PORT \
      CORE_GRAFANA_PORT CORE_GRAFANA_ADMIN_USER CORE_GRAFANA_ADMIN_PASSWORD; do
      if [ -z "${!required_key:-}" ]; then
        echo "❌ Укажите $required_key в $SCRIPT_DIR/.env"
        return 1
      fi
    done
  fi
  
  echo "Обработка шаблона для сервиса $service..."
  
  if [ ! -f "$template" ]; then
    echo "⚠️ Шаблон $template не найден, пропускаем..."
    return 0
  fi
  
  # Создаем директорию, если она еще не существует
  mkdir -p "$(dirname "$output")"
  
  # Используем envsubst для замены переменных в шаблоне
  $ENV_SUBST < "$template" > "$output"
  
  echo "✅ Создан файл $output"

  if [ "$service" = "core" ]; then
    local otel_template="$COMPOSE_DIR/core/otel/otel-collector-config.template.yaml"
    local otel_output="$COMPOSE_DIR/core/otel/otel-collector-config.yaml"

    if [ -f "$otel_template" ]; then
      if ! (
        set -a
        source "$output"
        set +a
        "$ENV_SUBST" < "$otel_template" > "$otel_output"
      ); then
        echo "❌ Не удалось создать $otel_output"
        return 1
      fi
      echo "✅ Создан файл $otel_output"
    fi

    local prometheus_template="$COMPOSE_DIR/core/prometheus/prometheus.template.yml"
    local prometheus_output="$COMPOSE_DIR/core/prometheus/prometheus.yml"

    if [ -f "$prometheus_template" ]; then
      if ! (
        set -a
        source "$output"
        set +a
        "$ENV_SUBST" < "$prometheus_template" > "$prometheus_output"
      ); then
        echo "❌ Не удалось создать $prometheus_output"
        return 1
      fi
      echo "✅ Создан файл $prometheus_output"
    fi
  fi
}

# Определяем список сервисов из переменной окружения
if [ -z "$SERVICES" ]; then
  echo "⚠️ Переменная SERVICES не задана. Нет сервисов для обработки."
  exit 0
fi

# Разделяем список сервисов по запятой
IFS=',' read -ra services <<< "$SERVICES"
echo "🔍 Обрабатываем сервисы: ${services[*]}"

# Обрабатываем шаблоны для всех указанных сервисов
success_count=0
skip_count=0
for service in "${services[@]}"; do
  process_template "$service" || exit 1
  if [ -f "$TEMPLATE_DIR/${service}.env.template" ]; then
    ((success_count++))
  else
    ((skip_count++))
  fi
done

if [ $success_count -eq 0 ]; then
  echo "⚠️ Ни один .env файл не создан. Проверьте список сервисов и наличие шаблонов."
else
  echo "🎉 Генерация завершена: $success_count файлов создано, $skip_count шаблонов пропущено"
fi
