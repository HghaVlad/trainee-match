import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import { useGetAdminResumesId, usePostAdminResumesIdArchive } from '@/api/generated/candidate/admin/admin'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Button } from '@/shared/ui/button'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'

export default function CandidateResumeDetailPage() {
  const { id = '', resumeId = '' } = useParams<{ id: string; resumeId: string }>()
  const navigate = useNavigate()
  const { toast } = useToast()
  const qc = useQueryClient()
  const archive = usePostAdminResumesIdArchive()

  const { data, isLoading, error, refetch } = useGetAdminResumesId(resumeId, {
    query: { enabled: Boolean(resumeId) },
  })

  async function onArchive() {
    try {
      await archive.mutateAsync({ id: resumeId })
      await qc.invalidateQueries({ queryKey: ['useGetAdminResumesId'] })
      toast({ title: 'Резюме архивировано' })
      navigate(`/candidate/${id}/resumes`)
    } catch (err) {
      const msg = err instanceof AppError ? err.message : 'Не удалось архивировать резюме'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  if (isLoading) return <LoadingState />
  if (error || !data) return <ErrorState onRetry={() => refetch()} />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link
          to={`/candidate/${id}/resumes`}
          className="text-sm text-muted-foreground underline"
        >
          ← Все резюме
        </Link>
      </div>

      <h1 className="text-2xl font-bold">{data.name ?? `Резюме #${resumeId}`}</h1>

      <div className="rounded-lg border bg-card p-4 space-y-2">
        <p><strong>ID:</strong> {data.id}</p>
        <p><strong>Название:</strong> {data.name ?? '—'}</p>
        <p><strong>Статус:</strong> {data.status ?? '—'}</p>
        <p><strong>Модерация:</strong> {data.moderation_status ?? '—'}</p>
        <p><strong>Кандидат:</strong> {data.candidate_id ?? id}</p>
      </div>

      <div className="flex gap-2">
        <Button variant="destructive" onClick={onArchive} disabled={archive.isPending}>
          {archive.isPending ? 'Архивация...' : 'Архивировать'}
        </Button>
      </div>
    </div>
  )
}