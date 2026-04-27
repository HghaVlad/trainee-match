import { Link } from 'react-router'
import { useGetResume, usePostResume } from '@/api/generated/candidate/resume/resume'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { useState } from 'react'
import { useNavigate } from 'react-router'
import { AppError } from '@/shared/api/http/client'

const statusLabel: Record<number, string> = { 0: 'Черновик', 1: 'Опубликовано' }

export default function ResumesPage() {
  const { data, isLoading, error, refetch } = useGetResume()
  const create = usePostResume()
  const navigate = useNavigate()
  const [err, setErr] = useState<string | null>(null)

  async function onCreate() {
    setErr(null)
    try {
      const r = await create.mutateAsync({
        data: { name: 'Новое резюме', status: 0, data: {} },
      })
      if (r?.id) navigate(`/me/resumes/${r.id}`)
      else await refetch()
    } catch (e) {
      setErr(e instanceof AppError ? e.message : 'Не удалось создать резюме')
    }
  }

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />

  const items = data ?? []

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Мои резюме</h1>
        <Button onClick={onCreate} disabled={create.isPending}>
          {create.isPending ? 'Создание...' : 'Создать'}
        </Button>
      </div>
      {err && <p role="alert" className="text-sm text-destructive">{err}</p>}
      {items.length === 0 ? (
        <EmptyState title="Резюме пока нет" />
      ) : (
        <ul className="space-y-2">
          {items.map((r) => (
            <li key={r.id} className="rounded-lg border bg-card p-4">
              <Link
                to={`/me/resumes/${r.id ?? ''}`}
                className="text-lg font-medium text-primary underline"
              >
                {r.name ?? '—'}
              </Link>
              <p className="text-sm text-muted-foreground">
                {statusLabel[r.status ?? 0] ?? '—'}
              </p>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
