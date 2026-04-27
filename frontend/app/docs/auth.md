# Authentication

Cookie-based session. The frontend NEVER touches access or refresh tokens directly — they live in `httpOnly` cookies set by the backend.

## Flow

```
1. POST /auth/login { username, password }
     <- Set-Cookie: access=...; HttpOnly; Secure; SameSite=Lax
     <- Set-Cookie: refresh=...; HttpOnly; Secure; SameSite=Strict; Path=/auth/refresh

2. Frontend calls bootstrap() → GET /auth/me
     <- { id, username, role, ... }
   Result is stored in Zustand sessionStore.

3. Subsequent API calls automatically include cookies (axios `withCredentials: true`).

4. On any 401, axios interceptor:
     a. POST /auth/refresh
     b. If 200: retry the original request once
     c. If non-200: clear sessionStore, redirect to /login?next=<current path>
```

## Files

| File | Role |
|------|------|
| `src/shared/api/http/client.ts` | Axios instance, `mutatorFn`, refresh interceptor, `AppError` |
| `src/shared/session/sessionStore.ts` | Zustand store: `{ user, isHydrated }` |
| `src/shared/session/bootstrap.ts` | Calls `/auth/me`, populates session, marks `isHydrated` |
| `src/app/router/guards.ts` | `requireAuth`, `requireRole`, `redirectIfAuth` |
| `src/features/auth/LoginForm.tsx` | Login form, calls `bootstrap()` after success, role-redirects |

## Why cookies, not JWT in localStorage

- `httpOnly` cookies are not reachable from JS → safe from XSS token theft.
- `SameSite` mitigates CSRF; backend additionally enforces a CSRF token on state-changing requests if needed.
- No token refresh logic in components — it's a pure interceptor concern.

## Forbidden patterns

- ❌ Reading `document.cookie` to extract tokens.
- ❌ Storing tokens in `localStorage` / `sessionStorage` / Zustand.
- ❌ Hand-writing `Authorization: Bearer ...` headers.
- ❌ Calling `/auth/refresh` from components — only the interceptor.

## Role-based routing

```ts
// src/app/router/router.tsx
{
  path: '/company',
  loader: requireRole('Company'),
  children: [...],
}
```

`requireRole` reads the current user from `sessionStore`. If `isHydrated` is false, it awaits `bootstrap()`. If role mismatch → redirects to `/forbidden`. If unauthenticated → `/login?next=...`.

After login, `LoginForm` reads `useSessionStore.getState().user?.role` and routes:

- `Company` → `/company/me`
- everything else → `/me/profile`

## Logout

`POST /auth/logout` (backend invalidates refresh token and clears cookies) → `useSessionStore.getState().reset()` → `navigate('/login')`. Implemented in `Header` user menu.

## Local development

When `VITE_USE_MSW=true`, MSW intercepts `/auth/login`, `/auth/me`, `/auth/refresh`, `/auth/logout` — no backend required.
