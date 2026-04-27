# Architecture

## Layers (Feature-Sliced Design)

Strictly upward dependencies. Lower layers cannot import from higher layers. Enforced via `eslint-plugin-boundaries`.

```
app  →  pages  →  widgets  →  features  →  entities  →  shared
```

| Layer | Responsibility | Examples |
|-------|---------------|----------|
| `app/` | Composition root: providers, router, global styles | `QueryProvider`, `ErrorBoundary`, `AppRouter` |
| `pages/` | Route components — thin shells | `pages/login`, `pages/me/profile`, `pages/company/me` |
| `widgets/` | Composite blocks reused across pages | `RootLayout`, `Header` |
| `features/` | User-facing features with their own forms / mutations | `auth`, `candidate-profile` |
| `entities/` | Pure domain models (rare; most domain stays in `api/generated`) | — |
| `shared/` | Cross-cutting utilities, ui-kit, http client, session, hooks | `shared/ui`, `shared/api/http`, `shared/session` |

`shared/` MUST NOT contain business logic. `api/generated/` is auto-generated and excluded from lint/typecheck owned by app.

## Data flow

```
Component
  └─> Generated React Query hook  (src/api/generated/<service>/<tag>/<tag>.ts)
        └─> Orval mutator         (src/shared/api/http/client.ts)
              └─> Axios instance with interceptors (auth refresh, error mapping)
                    └─> Backend API (cookies sent automatically)
```

- All requests go through ONE axios instance owned by the mutator.
- The mutator throws `AppError` on non-2xx; components check `e instanceof AppError && e.status === N`.
- 401 triggers a refresh attempt; on refresh failure the user is logged out and bounced to `/login`.

## State strategy

Two stores. Never duplicate.

| Store | Owns | Examples |
|-------|------|----------|
| TanStack Query cache | All server state | candidate profile, vacancies list, company members |
| Zustand `sessionStore` | Session/UI state that can't live in Query | `user`, `isHydrated`, transient UI flags |

Rules:
- If a piece of data is fetched from an endpoint → it lives in Query, full stop. Do not mirror it into Zustand.
- Invalidate via `queryClient.invalidateQueries({ queryKey: getXxxQueryKey() })` after mutations. Use the `getXxxQueryKey` helper exported by Orval, never hand-typed string keys.
- `sessionStore` is hydrated by `bootstrap()` in `src/shared/session/bootstrap.ts`, called once from `main.tsx` and after `/auth/login`.

## Routing

`src/app/router/router.tsx` declares all routes. Lazy-load every page via `React.lazy` + `Suspense`. Three guards live in `src/app/router/guards.ts`:

- `requireAuth` — bounces unauthenticated users to `/login?next=…`
- `requireRole(role)` — 403 if user role doesn't match
- `redirectIfAuth` — sends already-logged-in users away from `/login` and `/register`

Guards are React Router loaders, so they run before render and avoid flicker.

## Error / loading UX

- `<ErrorBoundary>` wraps `<AppRouter>` in `main.tsx` and renders a fallback `ErrorState`.
- Toasts via shadcn `<Toaster>` mounted once at the root (`src/main.tsx`).
- Loading: page-level `Suspense fallback={null}` plus per-component skeletons in `shared/ui/skeleton`.
- Empty states: `shared/ui/EmptyState` with consistent icon + title + CTA.

## Testing

- **Unit / integration**: Vitest + Testing Library + MSW. Handlers in `src/test/msw/handlers/`.
- **E2E**: Playwright in `e2e/`. Runs against `pnpm dev` with MSW enabled (`VITE_USE_MSW=true`).
- Tests must not hit a real backend.

## Build

`pnpm build` runs `tsc -b && vite build`. Generated client is excluded from app `tsconfig`. Output goes to `dist/`. Source maps disabled in production by default.
