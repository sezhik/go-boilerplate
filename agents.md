# AGENTS GUIDE

Этот файл задает единые правила для AI-агентов в моих pet-проектах.

## Стек

- **HTTP**: Gin
- **DB**: PostgreSQL 17 через pgx/v5, query builder goqu
- **Metrics**: Prometheus + Grafana
- **Connection pooling**: PGBouncer
- **Reverse proxy**: Caddy
- **Go**: 1.25, module name `go-boilerplate` (заменить под проект)

## Команды

- `make run` — поднять полный dev-стек через Docker Compose (с hot-reload).
- `make reset` — снести dev-окружение и volumes.
- `make build` — собрать локальный Docker-образ.
- `make build-production` — собрать production-образ (amd64).
- `make smoke` — smoke-тесты против живого dev-стека.
- `make migration-down` — откатить все миграции.
- `make deploy` — Ansible-деплой на VPS.
- `go test ./...` — все unit-тесты.
- `go generate ./...` — перегенерировать все моки.

## Архитектурный стиль

Используется Clean Architecture в адаптированном виде, без отдельных каталогов `ports`/`adapters`:

1. **Domain (внутренний слой)**
   - Папка: `internal/domain/`.
   - Rich-модели сущностей (Entity) и value objects.
   - Логика домена не должна зависеть от infra-деталей.

2. **Application**
   - Папки: `internal/usecase/`, `internal/model/`.
   - `usecase` — сценарии приложения.
   - `model` — input/output модели для usecase.

3. **Infra**
   - Папки: `internal/repository/`, `internal/gateway/`, `internal/presenter/`, `internal/handlers/`.
   - Реализации внешних интеграций и адаптеров.

### Request flow

```
HTTP Handler → validates input → builds model → Usecase → Repository/Gateway → domain logic
    ↓
Presenter formats usecase output → JSON response
```

### Wiring

Вся dependency injection происходит в `cmd/service/main.go`: DB pool → repositories → usecases → presenters → handlers → Gin router.

## Правила проектирования

- Соблюдать элементы DDD:
  - использовать value objects, чтобы избегать primitive obsession;
  - использовать rich-model подход для сущностей.
- Применять Dependency Inversion.
- Интерфейсы для инверсии зависимостей размещать рядом с местом использования:
  - в том же пакете;
  - в файле `contract.go`.
- Не выносить интерфейсы в общий пакет «на будущее» без необходимости.
- Новые usecase должны зависеть от абстракций, а не от конкретных infra-реализаций.

## Типовой поток выполнения запроса

- Запрос приходит в `handler`.
- В `handler` выполняется валидация входных данных:
  - через отдельный объект-валидатор, подключённый как зависимость;
  - в простом случае — прямо в `handler`.
- В `handler` собирается Input-модель (`internal/model`) и передаётся в `usecase`.
- `usecase` выполняет оркестрацию:
  - обращается к `repository` и/или `gateway`;
  - работает с доменными сущностями;
  - содержит необходимую условную логику и правила сценария.
- `usecase` возвращает Output-модель.
- `handler` формирует ответ:
  - в простом случае презентует Output-модель самостоятельно;
  - либо передаёт Output-модель в `presenter`, который формирует transport-ответ.

## Базовые соглашения по коду

- Имена пакетов: короткие, предметные, в lower-case.
- Публичные типы и функции оставлять только там, где это действительно нужно.
- Ошибки оборачивать с контекстом: `fmt.Errorf("...: %w", err)`.
- В возвращаемых ошибках не использовать динамические данные (`id`, `email`, `token` и т.д.); для них использовать структурированное логирование.
- Избегать заикания (stuttering) в ошибках: если оборачивается sentinel error (`NotFound` и т.п.), не дублировать тот же смысл в тексте контекста.
- Избегать заикания в именах пакетов, типов и функций; ориентироваться на Effective Go: https://go.dev/doc/effective_go#package-names.
- Для имен внутри пакета не дублировать имя пакета в идентификаторе:
  - в пакете `auth`: `SetClaims`, `Claims`, `ErrClaimsMissing`;
  - в пакете `session`: `Session.ID()`.
- Перед добавлением новых имён проверять, не повторяют ли они уже заданный контекст пакета или типа.
- Конструкторы именовать в стиле `NewX(...)`.
- Функции, которые могут вызвать `panic`, именовать с префиксом `Must`.
- Не создавать приватные helper-функции/методы, если они вызываются только один раз; в таком случае предпочитать inline-логику и выделять helper только при реальном переиспользовании или заметном упрощении чтения.
- Внешние DTO не смешивать с доменными сущностями.
- Output-модели usecase не отдавать напрямую наружу как transport-ответ; в handler/presenter маппить их в отдельные infra DTO.
- `json` и `db` теги могут быть только в структурах infra-слоя.
- Инициализация объектов, таких как gateway или repository, происходит в пакете `configure`.
- Каждый gateway размещать в отдельном пакете.
- Каждый repository размещать в отдельном пакете, но семантически связанные таблицы (например, образующие один агрегат) допустимо объединять в один repository в одном пакете — с одним типом, реализующим методы для всех таблиц.
- Каждый presenter размещать в отдельном пакете (например, `internal/presenter/item/item.go`) и использовать presenter как отдельную зависимость handler.
- В каждом presenter transport DTO выносить в отдельный файл `struct.go` рядом с основным файлом presenter.
- При изменениях публичного HTTP API обязательно актуализировать `public/openapi.yaml` в рамках той же задачи; в спецификацию включать только публичные API-роуты, без admin-only маршрутов.
- Если в пакете только один файл, называть файл по имени пакета (например, `session/session.go`, а не `entity.go`).
- Если хотя бы один метод типа использует pointer receiver, то все методы этого типа должны использовать pointer receiver.
- Не использовать более двух возвращаемых значений в функциях и методах.

## Соглашения по Postgres и repository

- Названия таблиц в Postgres использовать в единственном числе (например, `exercise`, `exercise_step`).
- В схемах БД не использовать `ON DELETE CASCADE`.
- В схемах БД не использовать foreign key (`REFERENCES`) ограничения.
- Для scan из Postgres в repository использовать возможности pgx сканирования в структуру (`RowToStructByName`, `CollectRows` и т.д.).
- Структура для db-сканирования должна быть с `db`-тегами и находиться в infra-слое рядом с repository в файле `struct.go`.
- В запросах repository использовать именованные параметры (например, `@user_id`) вместо позиционных (`$1`, `$2`, ...), если число параметров больше 1.

## Тесты и моки

- Для unit-тестов usecase использовать table-driven формат:
  - `tests := []struct{ ... }` с полями `name`, `input`, `prepare`, `expectations`;
  - в `prepare` задавать ожидания моков;
  - в `expectations` проверять результат и ошибку.
- Каждый тест-кейс запускать через `t.Run(tt.name, ...)`.
- Внутри каждого `t.Run` создавать отдельный `gomock.NewController(t)` и отдельные моки.
- После `prepare` вызывать тестируемый метод и передавать результат в `expectations`.
- Для генерации моков контрактов в `contract.go` использовать директиву:
  - `//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go`
- Эту директиву добавлять в любой `contract.go`, где объявлены интерфейсы зависимостей для unit-тестирования.
- Для repository/unit-тестов по умолчанию использовать black-box стиль с пакетом `${GOPACKAGE}_test`; white-box тесты допустимы только если без них невозможно проверить важное поведение через публичный API.
- Фреймворки: `stretchr/testify` для ассертов, `go.uber.org/mock/gomock` для моков.
- Smoke-тесты живут в `test/smoke/` под build-тегом `smoke` и ходят по HTTP в живой dev-стек; `go test ./...` не должен зависеть от сети.

---

# DevOps

Эта секция применяется к задачам, затрагивающим runtime, контейнеризацию, локальное окружение или деплой: `docker-compose.yaml`, `Makefile`, `Dockerfile`, `etc/caddy/`, `etc/ansible/`, `etc/grafana/`, `etc/prometheus/`, `.env`, `.env.prod`.

## Область ответственности

- Локальный запуск через Docker Compose и `make`-команды.
- Сборка контейнеров и production image.
- Reverse proxy и раздача статики через Caddy.
- Деплой и подготовка окружения через Ansible.
- Env-файлы и runtime-конфигурация.

## Источники истины

- Сначала ориентироваться на реальные файлы текущей ветки: `docker-compose.yaml`, `Makefile`, `Dockerfile`, `etc/caddy/`, `etc/ansible/`, `README.md`.
- Если меняются порты, service names, volume mounts, пути статики или env-переменные, синхронно обновлять все зависимые места: Compose, Caddy, Ansible (`docker-compose.yaml.j2` + `group_vars/all.yml`), Makefile и документацию.
- Не полагаться на устаревшие допущения о деплое: перед правками проверить, какой сценарий действительно поддерживается в текущей ветке.

## Локальная разработка

- Основной локальный сценарий — `make run`; очистка окружения — `make reset`.
- `.env` лежит в корне репозитория; Docker Compose использует его и для интерполяции (`${PG_USER}` и т.п.), и как `env_file` приложения.
- Dev-профиль в Compose использует `app-dev`, `caddy-dev`, `postgres`, `pgbouncer`, `prometheus`, `grafana`, `node-exporter`.
- `app-dev` работает через Docker target `dev` и hot reload (air).
- Вход в приложение локально — через Caddy на `http://localhost:8080`.
- Статические файлы в dev и prod раздаются из `public/`; маршрут `/s/*` обслуживает Caddy.

## Сборка и runtime

- `Dockerfile` содержит stage `dev`, `builder` и `prod`; production image должен оставаться минимальным и предсказуемым.
- Production binary собирается из `./cmd/service`; не менять эту точку входа без явной необходимости и связанных обновлений infra.
- Production Compose использует сервисы `app`, `caddy`, `postgres`, `pgbouncer`, `migration`, `prometheus`, `grafana`, `node-exporter`.
- Миграции применяются отдельным контейнером `migration`; изменения в schema flow должны оставаться совместимыми с этим сценарием.

## Деплой и инфраструктура

- В `etc/ansible/playbooks/deploy.yml` деплой собирает production image локально, загружает public assets и migrations, затем запускает `postgres`, `pgbouncer`, `migration`, `app`, observability-стек и пересоздаёт `caddy`.
- В `etc/ansible/playbooks/prepare.yml` подготавливается VPS под Docker Compose deploys; не ломать идемпотентность.
- Перед первым деплоем заполнить `etc/ansible/inventory.ini` (хост VPS), `etc/ansible/group_vars/all.yml` (`project_name`, `app_domain`, `app_dir`, `deploy_image`) и создать `.env.prod` в корне.
- Не переименовывать без необходимости `.env`, `.env.prod`, удалённые пути, имена сервисов и compose project name: на них завязаны playbooks и шаблоны.

## Безопасность и проверка

- Никогда не коммитить секреты из `.env`, `.env.prod` и производных артефактов.
- Предпочитать маленькие изменения с понятным rollback path.
- После devops-правок запускать минимально достаточную проверку: например, `docker compose config`, целевую `docker build`, `make run` или соответствующую проверку Ansible, если она доступна.

---

## Шаблон описания продукта (заменить в конкретном проекте)

Ниже пример секции, которую нужно обновить под конкретный pet-проект.

### Product context (example, replace me)

Project: **FocusFlow**

Коротко: сервис для планирования коротких фокус-сессий и трекинга прогресса.

Ключевые сценарии:
- Пользователь создает фокус-сессию с длительностью и тегом.
- Пользователь завершает сессию и получает запись в истории.
- Пользователь смотрит статистику за день/неделю.

Базовые доменные концепты:
- `Session` (entity)
- `Duration`, `Tag`, `Progress` (value objects)
- `SessionStatus` (value object/enum)

Внешние интеграции:
- `repository`: хранение сессий в Postgres
- `gateway`: интеграция с внешним таймером/нотификациями
- `presenter`: форматирование ответов для API/CLI

---

При старте нового проекта агент должен сначала переписать секцию **Product context** под текущий продукт, а затем предлагать структуру домена и usecase с учетом этих правил.
