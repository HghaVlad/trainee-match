# Application Service — сервис заявок

Сервис заявок — компонент системы, который управляет откликами кандидатов на вакансии, процессом рассмотрения HR и
аналитикой по заявкам. Данные хранятся в PostgreSQL, а проекции синхронизируются через Kafka события от смежных
сервисов. HTTP API описан в [`api/contract.yaml`](./api/contract.yaml) и используется для генерации обработчиков.

## Быстрый старт

1) Создайте файл `.env` в `backend/application` пример из [example env](.env.example)

2) Запустите зависимости и сервисы из папки `backend`:

```bash
cd "trainee-match\backend"
docker-compose up -d
```

## Зависимости сервиса

- `auth service` — выдача JWT для аутентификации.
- Keycloak (JWK endpoint) — проверка подписи токенов.
- Postgres — хранение заявок и проекций.
- Kafka + Schema Registry — обработка событий и DLQ.

## Инфраструктура

- Postgres для данных заявок и проекций.
- Kafka для межсервисных событий и DLQ.
- Schema Registry для Avro-схем.
- Keycloak как JWK провайдер.

## Структура проекта

- `cmd/` — точка входа сервиса.
- `internal/app/` — сборка приложения и wiring зависимостей.
- `internal/config/` — загрузка конфигурации (Viper + env).
- `internal/transport/http/` — HTTP-роутинг, middleware, handlers, OAPI.
- `internal/domain/` — доменные модели и валидация.
- `internal/usecase/` — сценарии заявок, HR, аналитики и проекций.
  - `application/` — основные сценарии: создание, просмотр, списки, изменение статусов.
  - `analytics/` — summary и dynamics отчеты.
  - `projection/` — обработка событий смежных сервисов и обновление проекций.
  - `common/` — общие утилиты и обработчики событий.
- `internal/infrastructure/` — PostgreSQL, Kafka, Schema Registry, утилиты.
- `api/contract.yaml` — OpenAPI спецификация.

## Зависимости (библиотеки)

- `github.com/go-chi/chi/v5` — HTTP-роутер.
- `github.com/spf13/viper` — конфигурация.
- `github.com/jackc/pgx/v5` — драйвер Postgres.
- `github.com/twmb/franz-go` — Kafka клиент.
- `github.com/hamba/avro/v2` — Avro сериализация.
- `github.com/avito-tech/go-transaction-manager` — транзакции.
- `github.com/google/uuid` — UUID.

## Основные возможности

- Кандидат: создание заявки, список, детали, отзыв заявки.
- HR: список заявок, детали, смена статуса.
- Аналитика: summary и dynamics по компании/вакансии.
- Проекции: обработка событий из смежных сервисов.

## API и примеры использования

Полное описание всех методов и моделей доступно в [`api/contract.yaml`](api/contract.yaml).

Кандидатские заявки:

- `POST /api/v1/applications`
- `GET /api/v1/applications`
- `GET /api/v1/applications/{applicationId}`
- `GET /api/v1/applications/{applicationId}/history`
- `POST /api/v1/applications/{applicationId}/withdraw`

HR заявки:

- `GET /api/v1/hr/companies/{companyId}/applications`
- `GET /api/v1/hr/vacancies/{vacancyId}/applications`
- `GET /api/v1/hr/applications/{applicationId}`
- `GET /api/v1/hr/applications/{applicationId}/history`
- `POST /api/v1/hr/applications/{applicationId}/status`

Аналитика:

- `GET /api/v1/hr/companies/{companyId}/analytics/summary`
- `GET /api/v1/hr/vacancies/{vacancyId}/analytics/summary`
- `GET /api/v1/hr/companies/{companyId}/analytics/dynamics`
- `GET /api/v1/hr/vacancies/{vacancyId}/analytics/dynamics`

## Аутентификация

Все защищенные методы требуют аутентификации через JWT (проверка по JWK). Токены хранится в куках.

## События и проекции

Сервис потребляет события для построения проекций кандидатов, резюме, вакансий и участников компаний. Примеры событий:

- candidate upserted
- resume upserted / deleted / archived
- vacancy published / updated / archived / moderation updated
- company updated / deleted / moderation updated
- company member added / removed

## Запуск вручную

```bash
cd "trainee-match\backend\application"
go run .\cmd\main.go
```

## Примечания

- Снимок данных резюме и кандидата фиксируется при создании заявки и не меняется.
