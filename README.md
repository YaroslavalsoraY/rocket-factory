# rocket factory project

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/YaroslavalsoraY/0ef8ab849dc6604bff79f92c46b9e7d5/raw/coverage.json)

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
sudo snap install task --classic
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
