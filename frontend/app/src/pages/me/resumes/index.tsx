import { Link } from 'react-router'
import { useGetResume, usePostResume } from '@/api/generated/candidate/resume/resume'
import { useGetCandidateMe } from '@/api/generated/candidate/candidate/candidate'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import { AppError } from '@/shared/api/http/client'
import { DefaultResumeStar, useDefaultResumeId } from '@/features/resume-default'
import { useSession } from '@/shared/session/useSession'
import type { DtoCandidateResponse } from '@/api/generated/candidate/schemas'

const statusLabel: Record<number, string> = { 0: 'Черновик', 1: 'Опубликовано' }

export default function ResumesPage() {
  const { data, isLoading, error, refetch } = useGetResume()
  const { data: candidate } = useGetCandidateMe({ query: { retry: false } })
  const { user } = useSession()
  const create = usePostResume()
  const navigate = useNavigate()
  const [err, setErr] = useState<string | null>(null)
  const { defaultResumeId, setDefaultResumeId } = useDefaultResumeId()

  useEffect(() => {
    if (!defaultResumeId || !data) return
    const found = data.find((r) => r.id === defaultResumeId)
    if (!found || (found.status ?? 0) !== 1) {
      setDefaultResumeId(undefined)
    }
  }, [data, defaultResumeId, setDefaultResumeId])

  async function onCreate() {
    setErr(null)
    const missing = collectMissing(user, candidate)
    if (missing) {
      setErr(missing)
      return
    }
    try {
      const r = await create.mutateAsync({
        data: {
          name: 'Новое резюме',
          status: 'draft',
          data: {
            first_name: user!.firstName ?? user!.username,
            last_name: user!.lastName ?? user!.username,
            email: user!.email!,
            phone: candidate!.phone!,
            city: candidate!.city!,
            citizenship: 'Россия',
            date_of_birth: candidate!.birthday,
            desired_format: 'remote',
            english_level: 'B1',
          },
        } as unknown as Parameters<typeof create.mutateAsync>[0]['data'],
      })
      if (r?.id) navigate(`/me/resumes/${r.id}`)
      else await refetch()
    } catch (e) {
      setErr(e instanceof AppError ? e.message : 'Не удалось создать резюме')
    }
  }

  if (isLoading) return <LoadingState />
  const notFound = error instanceof AppError && error.status === 404
  if (error && !notFound) return <ErrorState onRetry={() => refetch()} />

  const items = notFound ? [] : data ?? []

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
            <ResumeRow key={r.id} id={r.id ?? ''} name={r.name ?? '—'} status={r.status ?? 0} />
          ))}
        </ul>
      )}
    </div>
  )
}

function collectMissing(
  user: ReturnType<typeof useSession>['user'],
  candidate?: DtoCandidateResponse,
): string | null {
  if (!user?.email) {
    return 'В учётной записи отсутствует email — обратитесь к администратору.'
  }
  if (!candidate?.phone || !candidate.city) {
    return 'Сначала заполните профиль в /me/profile (нужны телефон и город).'
  }
  return null
}

function ResumeRow({ id, name, status }: { id: string; name: string; status: number }) {
  const { defaultResumeId, setDefaultResumeId } = useDefaultResumeId()
  const isPublished = status === 1

  useEffect(() => {
    if (defaultResumeId === id && !isPublished) {
      setDefaultResumeId(undefined)
    }
  }, [defaultResumeId, id, isPublished, setDefaultResumeId])

  return (
    <li className="flex items-start justify-between gap-3 rounded-lg border bg-card p-4">
      <div className="min-w-0 flex-1">
        <Link
          to={`/me/resumes/${id}`}
          className="text-lg font-medium text-primary underline"
        >
          {name}
        </Link>
        <p className="text-sm text-muted-foreground">
          {statusLabel[status] ?? '—'}
        </p>
      </div>
      {id && <DefaultResumeStar resumeId={id} isPublished={isPublished} />}
    </li>
  )
}
