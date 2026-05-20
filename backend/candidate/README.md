# Candidate Service — сервис профилей и резюме кандидатов

Сервис кандидатов — самостоятельный компонент системы, предоставляющий REST API для CRUD-операций над профилями
кандидатов, резюме и навыками.
Данные хранятся в PostgreSQL. При любом изменении сущностей сервис отправляет события в Kafka, обеспечивая
консистентность в смежных сервисах.
<br> Авторизация через JWT, для всех usecases написаны unit-тесты. Моки генерируются через mockery (конфиг
в [mockey.yml](.mockery.yml)) События в Kafka
через [avro схемы](internal/infrastructure/messagebroker/schemaregistry/avroschemas)

## Быстрый старт

1) Создайте файл `.env` в `backend/candidate` (пример из [example env](.env.example))

2) Запустите зависимости и сервисы из папки `backend`:

```bash
cd "E:\Code works\GoLang\trainee-match\backend"
docker-compose up -d
```

## Зависимости сервиса

- `auth service` — выдает `access_token` в куках для аутентификации.
- Keycloak (JWK endpoint) — проверка подписи токенов.
- Postgres — хранение профилей и резюме.
- Kafka + Schema Registry — публикация событий через outbox.

## Инфраструктура

- Postgres для данных кандидатов и резюме.
- Kafka для событий (outbox → продюсер).
- Schema Registry для схем Avro.
- Keycloak как провайдер JWK.

## Структура проекта

- `cmd/` — точка входа сервиса.
- `internal/app/` — сборка приложения и wiring зависимостей.
- `internal/delivery/http/` — HTTP-роутинг, middleware, handlers, DTO.
- `internal/domain/` — доменные модели и валидация.
- `internal/usecase/` — бизнес-логика и сценарии. Здесь лежат отдельные use case (create/update/get), оркестрация работы
  с репозиториями, outbox и валидация на уровне сценария.
    - `create_candidate/`, `update_candidate/`, `get_candidate_by_user_id/` — операции с профилем кандидата.
    - `create_resume/`, `update_resume/`, `get_resume/`, `remove_resume/` — операции с резюме.
    - `get_skill/` — чтение справочника навыков.
    - `common/` — общий код use case (outbox, утилиты).
- `internal/infrastructure/` — инфраструктурные детали: PostgreSQL (репозитории и миграции), Kafka/Schema Registry,
  outbox, клиентские обертки и интеграционные адаптеры.
    - `db/postgres/` — подключение, миграции, репозитории.
    - `messagebroker/kafka/` — Kafka продюсер/клиент.
    - `messagebroker/schemaregistry/` — работа со Schema Registry.
- `docs/` — Swagger/OpenAPI артефакты.

## Зависимости (библиотеки)

- `github.com/go-chi/chi/v5` — HTTP-роутер.
- `github.com/spf13/viper`, `github.com/mitchellh/mapstructure` — конфигурация.
- `github.com/jackc/pgx/v5` — драйвер Postgres.
- `github.com/golang-migrate/migrate/v4` — миграции.
- `github.com/lestrrat-go/jwx/v3` — проверка JWT/JWK.
- `github.com/twmb/franz-go` — Kafka клиент.
- `github.com/hamba/avro/v2` — Avro сериализация.
- `github.com/swaggo/swag`, `github.com/swaggo/http-swagger` — автогенерация и UI Swagger.
- `google.golang.org/grpc` — gRPC.
- `github.com/stretchr/testify` — unit-тесты.
- `github.com/avito-tech/go-transaction-manager` — транзакции.
- `github.com/google/uuid` — генерация UUID.

## Основные возможности

- Профиль кандидата: создание, получение "себя", обновление данных профиля.
- Резюме: создание, список резюме кандидата, получение по ID, обновление, удаление.
- Навыки получение навыка по ID и список доступных навыков.
- События в Kafka — изменения кандидатов и резюме публикуются через outbox (для других сервисов).

## API и примеры использования

Полное описание всех методов, моделей и возможность "потрогать" запросы доступны в [Swagger UI](docs/swagger.yaml) .

Группы методов:

- Candidate: управление профилем кандидата (создать, получить "себя", обновить).
- Resume: управление резюме кандидата (создать, список, получить по ID, обновить, удалить).
- Skill: справочник навыков (получить по ID, список).

Полный список ручек:

- `GET /api/v1/candidate/me`
- `POST /api/v1/candidate/`
- `PATCH /api/v1/candidate/`
- `POST /api/v1/resume/`
- `GET /api/v1/resume/`
- `GET /api/v1/resume/{id}`
- `PATCH /api/v1/resume/{id}`
- `GET /api/v1/skill/{id}`
- `GET /api/v1/skill/list`
- `GET /swagger/*`

Важно: все защищенные методы требуют аутентификации через cookie `access_token` (его выдает `auth` сервис).

## Установка и настройка

### Требования

- Go 1.25.x
- Postgres
- Kafka + Schema Registry
- Keycloak (или совместимый JWK endpoint)
- Docker (если поднимаете зависимости через `docker-compose`)

### Конфигурация

Сервис читает переменные окружения или файл `.env` в корне `backend/candidate`.

Запуск вручную :

```bash
cd "E:\Code works\GoLang\trainee-match\backend\candidate"
go run .\cmd\main.go
```

## Стиль и подход

- Handlers декодируют и валидируют DTO, затем вызывают use case.
- Use case лежат в `internal/usecase/*` и не знают про транспорт.
- Аутентификация cookie-first: `access_token` читается из cookies и валидируется по JWK.
- Изменения публикуются через outbox, а не напрямую из хендлеров.
