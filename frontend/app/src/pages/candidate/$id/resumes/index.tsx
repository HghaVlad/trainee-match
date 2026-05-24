import { useState } from 'react'
import { useParams, Link } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import {
  useGetAdminCandidatesIdResumes,
  usePostAdminResumesIdArchive,
} from '@/api/generated/candidate/admin/admin'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'

const PAGE_SIZE = 20

export default function CandidateResumesPage() {
  const { id = '' } = useParams<{ id: string }>()
  const qc = useQueryClient()
  const { toast } = useToast()
  const archive = usePostAdminResumesIdArchive()
  const [page, setPage] = useState(1)

  const { data, isLoading, error, refetch, isFetching } = useGetAdminCandidatesIdResumes(
    id,
    { page, size: PAGE_SIZE },
    { query: { enabled: Boolean(id) } },
  )

  async function onArchive(resumeId: string) {
    try {
      await archive.mutateAsync({ id: resumeId })
      await qc.invalidateQueries({ queryKey: ['useGetAdminCandidatesIdResumes'] })
      toast({ title: 'Резюме архивировано' })
    } catch (err) {
      const msg = err instanceof AppError ? err.message : 'Не удалось архивировать резюме'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />

  const hasMore = data != null && data.length >= PAGE_SIZE

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link
          to={`/candidates/${id}`}
          className="text-sm text-muted-foreground underline"
        >
          ← Профиль кандидата
        </Link>
      </div>

      <h1 className="text-2xl font-bold">Резюме кандидата</h1>

      {!data || data.length === 0 ? (
        <EmptyState title="Резюме не найдены" />
      ) : (
        <ul className="space-y-2">
          {data.map((resume) => (
            <li key={resume.id} className="rounded-lg border bg-card p-4">
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1">
                  <Link
                    to={`/candidates/${id}/resumes/${resume.id}`}
                    className="text-lg font-medium text-primary underline"
                  >
                    {resume.name ?? '—'}
                  </Link>
                  <p className="text-sm text-muted-foreground">
                    Статус: {resume.status ?? '—'} • Модерация:{' '}
                    {resume.moderation_status ?? '—'}
                  </p>
                </div>
                {resume.id && (
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => onArchive(resume.id as string)}
                    disabled={archive.isPending}
                  >
                    {archive.isPending ? '...' : 'Архивировать'}
                  </Button>
                )}
              </div>
            </li>
          ))}
        </ul>
      )}

      {data && data.length > 0 && (
        <div className="flex items-center justify-center gap-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1 || isFetching}
          >
            Назад
          </Button>
          <span className="text-sm text-muted-foreground">
            Страница {page}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => p + 1)}
            disabled={!hasMore || isFetching}
          >
            {isFetching ? 'Загрузка…' : 'Вперёд'}
          </Button>
        </div>
      )}
    </div>
  )
}
