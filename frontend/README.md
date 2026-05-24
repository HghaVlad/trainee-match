# trainee-match — Frontend

Веб-клиент для платформы подбора стажёров и компаний. React 19 + Vite + TypeScript с автогенерацией API-клиента из Swagger-спецификаций.

## Технологический стек

| Категория | Инструменты |
|-----------|-------------|
| Сборка | Vite 8, TypeScript 5 |
| UI | React 19, Tailwind CSS, shadcn/ui (Radix primitives) |
| Маршрутизация | React Router v7 |
| Состояние | TanStack Query 5 (серверное) + Zustand (сессия/UI) |
| Формы | React Hook Form + Zod |
| API | Axios + Orval (React Query hooks) |
| Тестирование | Vitest + Testing Library + MSW; Playwright (e2e) |
| Линтинг | ESLint (boundaries) + Prettier |

---

## 1. Запуск

### Локальная разработка (dev-сервер)

```bash
# Установка зависимостей
pnpm install

# Создание .env файла
cp app/.env.example app/.env
# Отредактируйте app/.env при необходимости

# Генерация API-клиента из Swagger
pnpm codegen

# Запуск dev-сервера
pnpm dev
# Откроется http://localhost:5173
```

**Переменные окружения** (определяются в `app/.env`, префикс `VITE_` для доступа в браузере):

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `VITE_USE_MSW` | `true` — включить MSW (мок API в браузере) | `false` |
| `VITE_BACKEND_URL` | Базовый URL бэкенда | `https://api.traineematch.space` |
| `VITE_AUTH_URL` | URL сервиса авторизации | `http://localhost:8000` |
| `VITE_CANDIDATE_URL` | URL сервиса кандидатов | `http://localhost:8081` |
| `VITE_COMPANY_URL` | URL сервиса компаний | `http://localhost:8088` |
| `VITE_APPLICATION_URL` | URL сервиса откликов | `http://localhost:8086` |

### Docker Compose

```bash
# Сборка и запуск
docker compose up --build

# С портом по умолчанию (8080)
docker compose up -d
```

Настройка через переменные окружения (`.env` в корне `frontend/`):

| Переменная | Описание |
|------------|----------|
| `VITE_API_URL` | Базовый URL API |
| `VITE_USE_MSW` | `true` — включить MSW |
| `SYNC_SWAGGER` | `true` — синхронизировать спецификации с бэкендом перед сборкой |
| `FRONTEND_PORT` | Внешний порт (по умолчанию 8080) |

Dockerfile собирает frontend в многоступенчатой схеме:
1. `deps` — установка pnpm и зависимостей
2. `codegen` — синхронизация Swagger + генерация API-клиента
3. `builder` — production-сборка Vite
4. `production` — nginx с отдачей статики

### Swagger-спецификации (`/swagger`)

Директория `swagger/` содержит локальные копии спецификаций всех микросервисов:

| Файл | Сервис |
|------|--------|
| `swagger-auth.yaml` | Авторизация |
| `swagger-candidate.yaml` | Кандидаты |
| `swagger-company.yaml` | Компании |
| `openapi-application.yaml` | Отклики (OpenAPI 3.0) |

### Скрипты (`/scripts`)

| Файл | Назначение |
|------|------------|
| `codegen.ts` | Конвертирует Swagger 2.0 → OpenAPI 3 → Orval (React Query hooks) |
| `sync-swagger.ts` | Синхронизирует спецификации с запущенными сервисами бэкенда |

```bash
# Синхронизация Swagger с бэкенд-сервисами (сервисы должны быть запущены)
pnpm sync:swagger

# Генерация / обновление API-клиента
pnpm codegen

# CI-проверка: ошибка если сгенерированные файлы отличаются от коммита
pnpm codegen:check
```

---

## 2. Архитектура Feature-Sliced Design (FSD)

Строго восходящие зависимости. Нижние слои не могут импортировать из верхних. Правило контролируется ESLint-плагином `boundaries`.

```
app  →  pages  →  widgets  →  features  →  entities  →  shared
```

| Слой | Ответственность | Примеры |
|------|-----------------|---------|
| `app/` | Точка сборки: провайдеры, роутер, глобальные стили | `QueryProvider`, `ErrorBoundary`, `AppRouter` |
| `pages/` | Компоненты маршрутов — тонкие обёртки | `pages/login`, `pages/me/profile`, `pages/company/me` |
| `widgets/` | Составные блоки, переиспользуемые на страницах | `RootLayout`, `Header` |
| `features/` | Пользовательские сценарии с формами и мутациями | `auth`, `candidate-profile`, `company-vacancies` |
| `entities/` | Чистые доменные модели (редко; основной код в `api/generated`) | — |
| `shared/` | Кросс-слойные утилиты: ui-kit, http-клиент, сессия, хуки | `shared/ui`, `shared/api/http`, `shared/session` |

**Важные правила:**
- `shared/` НЕ содержит бизнес-логику. Только механизмы, не политики.
- `api/generated/` — автогенерируемый код; исключён из линтинга и проверки типов.
- `boundaries`-плагин запрещает циклические и нисходящие импорты.

### Поток данных

```
Component
  └─> Generated React Query hook  (src/api/generated/<service>/<tag>/<tag>.ts)
        └─> Orval mutator         (src/shared/api/http/client.ts)
              └─> Axios instance с interceptors (auth refresh, error mapping)
                    └─> Backend API (cookies отправляются автоматически)
```

### Стратегия состояния

Два хранилища. Без дублирования.

| Хранилище | Управляет | Примеры |
|-----------|-----------|---------|
| TanStack Query cache | Всё серверное состояние | профиль кандидата, список вакансий, участники компании |
| Zustand `sessionStore` | Сессия/UI-состояние, не живущее в Query | `user`, `isHydrated`, transient UI flags |

**Правила:**
- Если данные приходят с эндпоинта → только Query. Не зеркалировать в Zustand.
- Инвалидация через `queryClient.invalidateQueries({ queryKey: getXxxQueryKey() })` после мутаций.
- `sessionStore` гидрируется через `bootstrap()` в `src/shared/session/bootstrap.ts`.

### Аутентификация

Cookie-based сессия. Токены НИКОГДА не хранятся клиентом — только в `httpOnly` cookies бэкенда.

```
1. POST /auth/login { username, password }
     <- Set-Cookie: access=...; HttpOnly; Secure; SameSite=Lax
     <- Set-Cookie: refresh=...; HttpOnly; Secure; SameSite=Strict; Path=/auth/refresh

2. Frontend вызывает bootstrap() → GET /auth/me
     <- { id, username, role, ... }
   Результат сохраняется в Zustand sessionStore.

3. Все API-вызовы автоматически включают cookies (axios `withCredentials: true`).

4. При 401, axios interceptor:
     a. POST /auth/refresh
     b. Если 200: повтор оригинального запроса
     c. Если не 200: очистка sessionStore, редирект на /login?next=<current path>
```

Подробнее: [`app/docs/auth.md`](app/docs/auth.md)

---

## 3. Структура проекта

```
frontend/
├── app/                        # Рабочая директория (pnpm workspace)
│   ├── src/
│   │   ├── app/               # app/ — точка сборки
│   │   │   ├── router/        # Роутер + guards
│   │   │   └── providers/     # Context providers
│   │   ├── pages/             # Страницы (thin shells)
│   │   │   ├── login/
│   │   │   ├── me/
│   │   │   ├── vacancies/
│   │   │   ├── companies/
│   │   │   └── company/
│   │   ├── widgets/           # Композитные блоки
│   │   │   ├── RootLayout/
│   │   │   ├── Header/
│   │   │   └── ...
│   │   ├── features/          # Фичи (бизнес-логика + UI)
│   │   │   ├── auth/
│   │   │   ├── candidate-profile/
│   │   │   ├── company-vacancies/
│   │   │   ├── company-members/
│   │   │   ├── company-profile/
│   │   │   ├── company-create/
│   │   │   ├── resume-default/
│   │   │   ├── skill-catalog/
│   │   │   ├── applications/
│   │   │   ├── hr-applications/
│   │   │   └── analytics/
│   │   ├── shared/            # Кросс-слойные утилиты
│   │   │   ├── ui/            # UI-кит (shadcn)
│   │   │   ├── api/http/      # Axios + interceptors
│   │   │   └── session/       # Zustand session store
│   │   ├── api/
│   │   │   └── generated/     # Orval output — НЕ РЕДАКТИРОВАТЬ
│   │   │       ├── auth/
│   │   │       ├── candidate/
│   │   │       └── company/
│   │   ├── test/
│   │   │   └── msw/           # MSW handlers для unit/integration
│   │   ├── assets/
│   │   ├── main.tsx
│   │   └── docs/              # Внутренняя документация
│   ├── e2e/                   # Playwright e2e-тесты
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   └── docs/                  # README, архитектура, конвенции
│
├── swagger/                   # Swagger 2.0 / OpenAPI 3 спецификации
│   ├── swagger-auth.yaml
│   ├── swagger-candidate.yaml
│   ├── swagger-company.yaml
│   └── openapi-application.yaml
│
├── scripts/                   # Скрипты сборки и codegen
│   ├── codegen.ts
│   └── sync-swagger.ts
│
├── Dockerfile
└── docker-compose.yml
```

---

## 4. Тестирование

### MSW (unit + integration)

```bash
pnpm test          # Запуск всех тестов
pnpm test:watch    # watch-режим
pnpm test:ui       # Vitest UI
pnpm coverage       # Покрытие кода
```

**Структура:**
- Handlers: `src/test/msw/handlers/`
- Тесты живут рядом с компонентами: `*.test.ts(x)`
- MSW перехватывает все HTTP-запросы. Тесты НЕ обращаются к реальному бэкенду.

**Пример теста:**

```tsx
// src/features/auth/LoginForm.test.tsx
import { server } from '@/test/msw/server'
import { http, HttpResponse } from 'msw'

it('shows error on failed login', async () => {
  server.use(
    http.post('/api/auth/login', () => HttpResponse.json({ message: 'Invalid credentials' }, { status: 401 }))
  )
  // ...
})
```

### E2E (Playwright)

```bash
# Запуск против dev-сервера с MSW
pnpm dev        # в одном терминале
pnpm test:e2e   # в другом

# UI-режим для отладки
pnpm test:e2e:ui
```

Тесты находятся в `e2e/` и запускаются против `http://localhost:5173` с `VITE_USE_MSW=true`.

---

## 5. Документация (`/frontend/app/docs`)

| Файл | Содержание |
|------|------------|
| `architecture.md` | FSD-слои, поток данных, стратегия состояния, роутинг |
| `codegen.md` | Pipeline Orval, правила генерации API-клиента |
| `auth.md` | Cookie-аутентификация, refresh interceptor, role guards |
| `conventions.md` | Code conventions, запрещённые паттерны, примеры |
| `screens-map.md` | Карта экранов, маршруты, модалки, статусы |

---

## Команды для контроля качества

```bash
pnpm typecheck    # tsc --noEmit
pnpm lint         # ESLint
pnpm test --run   # Vitest (unit + integration)
pnpm test:e2e     # Playwright (smoke)
pnpm build        # Production bundle
```

Перед коммитом рекомендуется: `pnpm typecheck && pnpm lint && pnpm test --run`

---

## Добавление новых фич

1. Добавить Swagger-спецификацию в `frontend/swagger/`
2. Добавить запись в массив `specs` в `scripts/codegen.ts`
3. Добавить target в `frontend/app/orval.config.ts`
4. Запустить `pnpm codegen`
5. Создать slice в `src/features/` по FSD
6. Написать MSW-handlers в `src/test/msw/handlers/`
7. Добавить маршрут в `src/app/router/router.tsx`
