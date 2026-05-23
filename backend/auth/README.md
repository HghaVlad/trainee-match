# Auth service — сервис аутентификации

Сервис отвечает за регистрацию и аутентификацию пользователей, выдачу пары JWT-токенов в httpOnly cookies (
`access_token`, `refresh_token`) и получение данных пользователя по access-токену. В качестве identity provider
используется Keycloak (OpenID Connect), а события о создании пользователя публикуются через outbox в Kafka.
<br>Для всех usecases написаны unit-тесты. Моки генерируются через mockery (конфиг
в [mockey.yml](.mockery.yml))

## Быстрый старт

1) Создайте файл `.env` в `backend/auth` на основе [`.env.example`](.env.example).
2) Запустите зависимости и сервисы из папки `backend`:

```bash
cd "E:\Code works\GoLang\trainee-match\backend"
docker-compose up -d
```

Если `auth` не поднялся в compose, запустите вручную:

```bash
cd "E:\Code works\GoLang\trainee-match\backend\auth"
go run .\cmd\main.go
```

## Зависимости сервиса

- Keycloak — аутентификация/авторизация через OIDC, выдача токенов.
- Postgres — хранение outbox и технических данных.
- Kafka — доставка событий о создании пользователя.
- Schema Registry — хранение Avro-схем для событий.

## Инфраструктура

- Postgres для outbox и миграций.
- Kafka + Schema Registry для событий `user-created`.
- Keycloak как identity provider.
- Swagger спецификация в `docs/` (yaml/json).

## Структура проекта

- `cmd/` — точка входа сервиса.
- `internal/app/` — сборка приложения и wiring зависимостей.
- `internal/delivery/http/` — HTTP-роутинг, handlers, DTO, helpers.
- `internal/domain/` — доменные модели и события.
- `internal/usecase/` — сценарии регистрации/логина/рефреша/логаута/получения профиля.
    - `register/`, `login/`, `refresh_token/`, `logout/`, `get_user_me/` — отдельные use case.
    - `common/outbox/` — запись и релей событий из outbox.
- `internal/infra/` — инфраструктурные адаптеры.
    - `db/postgres/` — подключение, миграции, outbox-репозиторий.
    - `keycloack/` — клиент для Keycloak.
    - `message_broker/kafka/` — Kafka клиент/producer.
    - `message_broker/schemaregistry/` — Schema Registry, Avro-схемы.
- `docs/` — Swagger документация.

## Зависимости (библиотеки)

- `github.com/Nerzal/gocloak/v13` — клиент Keycloak.
- `github.com/go-chi/chi` — HTTP-роутер.
- `github.com/go-playground/validator/v10` — валидация входных DTO.
- `github.com/spf13/viper` — конфигурация.
- `github.com/jackc/pgx/v5` — драйвер Postgres.
- `github.com/golang-migrate/migrate/v4` — миграции.
- `github.com/twmb/franz-go` — Kafka клиент.
- `github.com/hamba/avro/v2` — Avro сериализация.
- `github.com/swaggo/swag` — генерация Swagger.
- `github.com/avito-tech/go-transaction-manager` — транзакции.
- `github.com/google/uuid` — идентификаторы.

## Основные возможности

- Регистрация пользователя в Keycloak.
- Логин и выдача пары токенов в httpOnly cookies.
- Обновление access-токена по refresh-токену.
- Логаут с отзывом refresh-токена и очисткой cookies.
- Получение данных пользователя по access-токену (`/me`).
- Публикация события `user-created` через outbox в Kafka.

## API и примеры использования

Полное описание методов и моделей доступно в Swagger-спецификации: [docs/swagger.yaml](docs/swagger.yaml).

Группы методов:

- Auth: регистрация, логин, рефреш, логаут, получение профиля.

Полный список ручек:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/me`

Аутентификация cookie-first: `access_token` и `refresh_token` устанавливаются в httpOnly cookies. `refresh` и `logout`
используют refresh-токен из cookies, `me` — access-токен из cookies.

## Установка и настройка

### Требования

- Go 1.25.x
- Postgres
- Kafka + Schema Registry
- Keycloak
- Docker (если поднимаете зависимости через `docker-compose`)

## Стиль и подход

- Тонкий handler → use case, без бизнес-логики в HTTP-слое.
- Cookie-first аутентификация (`access_token`, `refresh_token`).
- Интеграция с Keycloak через выделенный клиент в `internal/infra/keycloack`.
- События о создании пользователей публикуются через outbox.
