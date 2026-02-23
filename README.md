# Go Boilerplate для pet-проектов

Этот репозиторий - минимальный шаблон для запуска новых Go pet-проектов с акцентом на чистую архитектуру и готовую инфраструктуру.

## Что уже есть

- Базовый каркас приложения: `cmd/service/main.go`
- Конфигурация и инфраструктурные зависимости:
  - логгер (`internal/configure/telemetry.go`)
  - подключение к БД (`internal/configure/configure.go`)
  - graceful shutdown по сигналам (`internal/util/signal/signal.go`)
- Docker-окружение:
  - `docker-compose.yaml` (Postgres, PgBouncer, app, migrator)
  - `Dockerfile` (dev/prod стадии)
- Миграции БД:
  - пример `create table example` в `database/migrations/0001_init.up.sql`
  - откат в `database/migrations/0001_init.down.sql`
- Примеры infra-адаптеров:
  - `internal/repository/example_pg.go`
  - `internal/gateway/example_http.go`

## Архитектурная структура

Проект следует адаптированной Clean Architecture:

- `internal/domain` - внутренний слой (DDD-сущности и value objects)
- `internal/usecase` + `internal/model` - слой application
- `internal/repository` + `internal/gateway` + `internal/presenter` - слой infra

Отдельного слоя `ports/adapters` нет.

Подробные правила для AI-агентов и проектирования находятся в `agents.md`.

## Быстрый старт

1. Скопируй и обнови `.env` под свой проект.
2. Подними окружение:

```bash
make run
```

3. Убедись, что код собирается:

```bash
go build ./...
```

## Как использовать этот шаблон

- Заменить примерный `Product context` в `agents.md` на описание конкретного проекта.
- Заполнить `internal/domain`, `internal/usecase`, `internal/model` своей предметной областью.
- Расширить infra-адаптеры в `internal/repository`, `internal/gateway`, `internal/presenter`.
- Добавить свои миграции поверх примерной таблицы `example`.
