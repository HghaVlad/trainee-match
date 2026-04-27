# Codegen Pipeline

API client code is auto-generated from Swagger 2.0 specs.

## Pipeline

```
frontend/swagger/*.yaml (Swagger 2.0)
  | swagger2openapi --patch
frontend/app/.codegen-cache/openapi/*.json (OpenAPI 3)
  | orval (axios + react-query + zod)
frontend/app/src/api/generated/{auth,candidate,company}/
```

## Commands

```bash
pnpm codegen        # Generate / update API client
pnpm codegen:check  # CI: fail if generated files drift
```

The driver script lives in `frontend/scripts/codegen.ts` and is run from `frontend/app/` via `tsx`.

## Targets

`orval.config.ts` defines three targets sharing one mutator:

- `auth`     -> `src/api/generated/auth`
- `candidate`-> `src/api/generated/candidate`
- `company`  -> `src/api/generated/company`

All targets use `mode: tags-split`, `client: react-query`, `httpClient: axios`, and the mutator
`src/shared/api/http/client.ts#mutatorFn`. The mutator is a stub here and will be fully wired
(401 refresh, error mapping) in T6.

## Adding a service

1. Drop a Swagger 2.0 spec into `frontend/swagger/`.
2. Add an entry to the `specs` array in `frontend/scripts/codegen.ts`.
3. Add a target in `frontend/app/orval.config.ts`.
4. Run `pnpm codegen`.

## Rules

- Do not edit files under `src/api/generated/`. They are excluded from `tsc` (see `tsconfig.app.json`)
  and from ESLint (see `eslint.config.js`).
- Do not write hand-rolled DTO types that duplicate generated schemas.
- Do not enable `mock: true` in orval; mocks are owned by MSW in tests.
- Use `pnpm` only.
