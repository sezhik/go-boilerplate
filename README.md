# Go Boilerplate для pet-проектов

Этот репозиторий — шаблон для запуска новых Go pet-проектов с акцентом на чистую архитектуру и готовую инфраструктуру: HTTP-сервис на Gin, Postgres, метрики, reverse proxy и деплой на VPS.

## Что уже есть

- HTTP-сервис на Gin: `cmd/service/main.go` (роутер, health, метрики, wiring зависимостей)
- Пример полного вертикального слайса `example`:
  - `internal/domain/example` — rich-сущность
  - `internal/model/example` — input/output модели usecase
  - `internal/usecase/example` — сценарий с `contract.go`, gomock-моками и table-driven тестом
  - `internal/presenter/example` — transport DTO в `struct.go`
  - `internal/handlers/http/example` — Gin-handler
  - `internal/repository/example` — Postgres-репозиторий на goqu + pgx (`struct.go` с db-тегами)
  - `internal/gateway/example` — HTTP-гейтвей во внешний API
- Observability:
  - Prometheus-метрики HTTP (`internal/metrics`), эндпоинт `/metrics`
  - Готовые Prometheus + Grafana (provisioning + дашборд) + node-exporter в Compose
- Health check: `GET /_health` (`internal/handlers/http/health`)
- Конфигурация и инфраструктурные зависимости:
  - логгер (`internal/configure/telemetry.go`)
  - подключение к БД с ожиданием миграций (`internal/configure/configure.go`)
- Docker-окружение:
  - `docker-compose.yaml` (Postgres, PgBouncer, app, migrator, Caddy, Prometheus, Grafana, node-exporter; профили dev/prod)
  - `Dockerfile` (dev/prod стадии, hot-reload через air)
  - Caddy как reverse proxy и раздача статики из `public/` (`etc/caddy/`)
- Деплой на VPS через Ansible (`etc/ansible/`): prepare + deploy playbooks, шаблон production Compose
- Миграции БД (`database/migrations/`, migrate/migrate)
- Тесты:
  - unit: testify + go.uber.org/mock, table-driven (`internal/usecase/example/example_test.go`)
  - smoke: `test/smoke/` под build-тегом `smoke`, ходят по HTTP в живой dev-стек
- Публичный API-контракт: `public/openapi.yaml`

## Архитектурная структура

Проект следует адаптированной Clean Architecture:

- `internal/domain` — внутренний слой (DDD-сущности и value objects)
- `internal/usecase` + `internal/model` — слой application
- `internal/repository` + `internal/gateway` + `internal/presenter` + `internal/handlers` — слой infra

Отдельного слоя `ports/adapters` нет.

Подробные правила для AI-агентов и проектирования находятся в `agents.md`.

## Быстрый старт

1. Создай `.env` из шаблона и обнови под свой проект:

```bash
cp .env.example .env
```

2. Подними окружение:

```bash
make run
```

3. Проверь, что стек живой:

```bash
curl http://localhost:8080/_health
curl http://localhost:8080/api/examples
make smoke
```

4. Убедись, что код собирается и тесты проходят:

```bash
go build ./...
go test ./...
```

Grafana: http://127.0.0.1:3000 (admin/admin), Prometheus: http://127.0.0.1:9090.

## Деплой

1. Заполни `etc/ansible/inventory.ini` (адрес VPS) и `etc/ansible/group_vars/all.yml` (`project_name`, `app_domain`, `app_dir`, `deploy_image`).
2. Создай `.env.prod` в корне (не коммитится).
3. Подготовь VPS и задеплой:

```bash
ansible-playbook -i etc/ansible/inventory.ini etc/ansible/playbooks/prepare.yml
make deploy
```

## Как использовать этот шаблон

- Заменить module name `go-boilerplate` в `go.mod` и импортах на имя проекта.
- Заменить примерный `Product context` в `agents.md` на описание конкретного проекта.
- Заполнить `internal/domain`, `internal/usecase`, `internal/model` своей предметной областью по образцу слайса `example`.
- Расширить infra-адаптеры в `internal/repository`, `internal/gateway`, `internal/presenter`, `internal/handlers`.
- Добавить свои миграции поверх примерной таблицы `example`.
- Актуализировать `public/openapi.yaml` при изменениях публичного API.
