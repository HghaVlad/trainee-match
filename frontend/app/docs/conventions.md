# Code conventions

## Forbidden patterns (Must NOT)

These are enforced by review and, where possible, by ESLint. If you must violate one, add a `// reason: <why>` comment on the offending line and discuss in PR.

- ❌ **No manual DTO / type duplication**. Use what `src/api/generated/` exports.
- ❌ **No tokens in JS**. No `localStorage`, `sessionStorage`, `document.cookie`, or Zustand for `access`/`refresh`. Cookies only — see `docs/auth.md`.
- ❌ **No `any`, `as unknown as`, `@ts-ignore`** without `// reason: ...`. Use proper generics or `unknown` + type guards.
- ❌ **No duplicated server state** in Zustand. If it's in React Query cache, keep it there.
- ❌ **No inline `fetch` / `axios`** in components. Use generated hooks or wrappers from `shared/api/`.
- ❌ **No deep barrel re-exports** (`export *` chains > 1 level) — bundler perf and tree-shaking.
- ❌ **No custom router / forms / i18n frameworks**. Stack is fixed: React Router, RHF + Zod, no i18n yet.
- ❌ **No business logic in `src/shared/`**. Shared is mechanism, not policy.
- ❌ **No `console.log` in commits** (lint rule: `no-console` with `warn` allowed for `error`/`warn`).
- ❌ **No applications/interviews/offers code beyond `EmptyState` stubs** until backend endpoints exist (feature-flagged).

## Required patterns (Must)

### Forms

```tsx
// src/features/.../SomeForm.tsx
import { z } from 'zod'
import { FormWrapper } from '@/shared/ui/Form'

const schema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
})
type Data = z.infer<typeof schema>

export function SomeForm() {
  const m = usePostThing()
  return (
    <FormWrapper schema={schema} defaultValues={{ name: '', email: '' }}
                 onSubmit={(d) => m.mutateAsync({ data: d })}>
      {/* FormField + FormControl + FormMessage from shadcn */}
    </FormWrapper>
  )
}
```

`FormWrapper` accepts `schema: ZodType` and wires `zodResolver` + RHF context.

### Errors

```ts
import { AppError } from '@/shared/api/http/client'

try {
  await mutation.mutateAsync({ data })
} catch (e) {
  if (e instanceof AppError && e.status === 404) {
    // not found case
  } else if (e instanceof AppError) {
    setError(e.message)
  } else {
    setError('Что-то пошло не так')
  }
}
```

### Cache invalidation

After a mutation, invalidate using the generated query-key helper:

```ts
import { getGetCandidateMeQueryKey } from '@/api/generated/candidate/candidate/candidate'

queryClient.invalidateQueries({ queryKey: getGetCandidateMeQueryKey() })
```

### Cursor pagination

```ts
const [cursor, setCursor] = useState<string | undefined>(undefined)
const { data } = useGetVacancies({ cursor, limit: 20 })
// "Load more" button: onClick={() => setCursor(data?.nextCursor ?? undefined)}
```

### Lazy routes

Add new pages via `lazy()` + `Suspense`. See `src/app/router/router.tsx`.

### Imports

- Absolute via `@/...` alias for anything outside the current feature folder.
- Relative for files inside the same feature/page slice.
- Group order: external → `@/app` → `@/widgets` → `@/features` → `@/entities` → `@/shared` → `@/api/generated` → relative.

## Tooling

- `pnpm` only. `yarn` and `npm` lockfiles must not appear.
- Commits should be atomic and pass `pnpm typecheck && pnpm lint && pnpm test --run && pnpm build` locally before push.
- Generated code (`src/api/generated/`) is excluded from lint and typecheck and is checked in.
