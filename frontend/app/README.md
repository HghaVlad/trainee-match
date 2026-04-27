# trainee-match — Frontend

Web client for trainee/company matching. React 19 + Vite + TypeScript with auto-generated API client from Swagger.

## Tech stack

- **Build**: Vite 7, TypeScript 5
- **UI**: React 19, Tailwind CSS, shadcn/ui (Radix primitives)
- **Routing**: React Router v7
- **Data**: TanStack Query 5 (server state) + Zustand (session/UI state only)
- **Forms**: React Hook Form + Zod
- **API**: Axios via Orval-generated React Query hooks
- **Tests**: Vitest + Testing Library + MSW; Playwright for e2e
- **Lint**: ESLint (boundaries) + Prettier

## Quick start

```bash
pnpm install
pnpm codegen     # generate API client from swagger specs
pnpm dev         # http://localhost:5173 (MSW enabled in dev when VITE_USE_MSW=true)
```

Quality gates:

```bash
pnpm typecheck   # tsc --noEmit
pnpm lint        # eslint
pnpm test        # vitest (unit + integration with MSW)
pnpm test:e2e    # playwright (smoke)
pnpm build       # production bundle
```

## Project structure (Feature-Sliced Design)

```
src/
  app/        composition root: providers, router, global styles
  pages/      route components — thin shells that compose widgets/features
  widgets/    composite blocks (RootLayout, Header, ...)
  features/   user-facing features (auth, candidate-profile, ...)
  entities/   domain models without a feature surface (rare)
  shared/     cross-cutting code: ui-kit, http client, session, hooks
  api/
    generated/  Orval output — DO NOT EDIT
  test/       MSW handlers + vitest setup
e2e/          Playwright specs
```

ESLint `boundaries` plugin enforces upward-only imports. `shared/` may not import from any other layer.

## Environment variables

Defined in `.env`, prefixed `VITE_` to be exposed to the browser:

| Var | Purpose |
|-----|---------|
| `VITE_API_URL` | Backend base URL. Empty = same-origin (proxy via Vite in dev). |
| `VITE_USE_MSW` | `true` enables MSW in browser; used for dev-without-backend and Playwright. |

## Backend integration

API client is generated from Swagger 2.0 specs. The pipeline converts Swagger → OpenAPI 3 → Orval (axios + react-query). See [`docs/codegen.md`](./docs/codegen.md).

Auth is cookie-based. Tokens are NEVER stored client-side. See [`docs/auth.md`](./docs/auth.md).

## Documentation

- [`docs/architecture.md`](./docs/architecture.md) — FSD layers, data flow, state strategy
- [`docs/codegen.md`](./docs/codegen.md) — Orval pipeline and rules
- [`docs/auth.md`](./docs/auth.md) — Cookie auth, refresh interceptor, role guards
- [`docs/conventions.md`](./docs/conventions.md) — Code conventions and forbidden patterns
