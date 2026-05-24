import { Link } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import {
  getGetResumeIdQueryKey,
  getGetResumeQueryKey,
  useGetResume,
  usePatchResumeId,
  usePostResume,
} from '@/api/generated/candidate/resume/resume'
import { useGetCandidateMe } from '@/api/generated/candidate/candidate/candidate'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { Badge } from '@/shared/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import { AppError } from '@/shared/api/http/client'
import { DefaultResumeStar, useDefaultResumeId } from '@/features/resume-default'
import { useSession } from '@/shared/session/useSession'
import { useToast } from '@/shared/hooks/use-toast'
import type { DtoCandidateResponse } from '@/api/generated/candidate/schemas'

type ResumeStatusValue = 'draft' | 'published' | number | string | undefined

function isPublishedStatus(status: ResumeStatusValue): boolean {
  if (typeof status === 'number') return status === 1
  if (typeof status === 'string') return status === 'published'
  return false
}

export default function ResumesPage() {
  const { data, isLoading, error, refetch } = useGetResume()
  const { data: candidate, isLoading: candidateLoading } = useGetCandidateMe({ query: { retry: false } })
  const { user } = useSession()
  const create = usePostResume()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const [err, setErr] = useState<string | null>(null)
  const { defaultResumeId, setDefaultResumeId } = useDefaultResumeId()

  useEffect(() => {
    if (!defaultResumeId || !data) return
    const found = data.find((r) => r.id === defaultResumeId)
    if (!found || !isPublishedStatus(found.status)) {
      setDefaultResumeId(undefined)
    }
  }, [data, defaultResumeId, setDefaultResumeId])

  async function onCreate() {
    if (candidateLoading || !candidate) return
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
      await qc.invalidateQueries({ queryKey: getGetResumeQueryKey() })
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
            <ResumeRow
              key={r.id}
              id={r.id ?? ''}
              name={r.name ?? '—'}
              status={r.status as ResumeStatusValue}
            />
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

function ResumeStatusBadge({ published }: { published: boolean }) {
  if (published) return <Badge>Опубликовано</Badge>
  return <Badge variant="outline">Черновик</Badge>
}

function ResumeRow({
  id,
  name,
  status,
}: {
  id: string
  name: string
  status: ResumeStatusValue
}) {
  const { defaultResumeId, setDefaultResumeId } = useDefaultResumeId()
  const isPublished = isPublishedStatus(status)
  const isDefault = defaultResumeId === id
  const qc = useQueryClient()
  const { toast } = useToast()
  const patch = usePatchResumeId()
  const [confirmPublish, setConfirmPublish] = useState(false)
  const [confirmDraft, setConfirmDraft] = useState(false)

  useEffect(() => {
    if (isDefault && !isPublished) {
      setDefaultResumeId(undefined)
    }
  }, [isDefault, isPublished, setDefaultResumeId])

  async function onPublish() {
    try {
      await patch.mutateAsync({
        id,
        data: { status: 'published' },
      })
      await Promise.all([
        qc.invalidateQueries({ queryKey: getGetResumeQueryKey() }),
        qc.invalidateQueries({ queryKey: getGetResumeIdQueryKey(id) }),
      ])
      toast({ title: 'Резюме опубликовано' })
      setConfirmPublish(false)
    } catch (e) {
      toast({
        title: 'Ошибка',
        description: e instanceof AppError ? e.message : 'Не удалось опубликовать резюме',
        variant: 'destructive',
      })
    }
  }

  async function onDraft() {
    try {
      await patch.mutateAsync({
        id,
        data: { status: 'draft' },
      })
      if (isDefault) setDefaultResumeId(undefined)
      await Promise.all([
        qc.invalidateQueries({ queryKey: getGetResumeQueryKey() }),
        qc.invalidateQueries({ queryKey: getGetResumeIdQueryKey(id) }),
      ])
      toast({ title: 'Резюме переведено в черновик' })
      setConfirmDraft(false)
    } catch (e) {
      toast({
        title: 'Ошибка',
        description: e instanceof AppError ? e.message : 'Не удалось изменить статус резюме',
        variant: 'destructive',
      })
    }
  }

  return (
    <li
      data-testid="resume-row"
      data-resume-id={id}
      className="flex items-start justify-between gap-3 rounded-lg border bg-card p-4"
    >
      <div className="min-w-0 flex-1 space-y-2">
        <Link
          to={`/me/resumes/${id}`}
          className="text-lg font-medium text-primary underline"
        >
          {name}
        </Link>
        <div className="flex flex-wrap items-center gap-2">
          <ResumeStatusBadge published={isPublished} />
          {isDefault && (
            <Badge
              variant="secondary"
              className="border-yellow-300/50 bg-yellow-100 text-yellow-900 dark:bg-yellow-500/20 dark:text-yellow-200"
            >
              Основное
            </Badge>
          )}
        </div>
      </div>
      <div className="flex flex-col items-end gap-2">
        {id && <DefaultResumeStar resumeId={id} isPublished={isPublished} />}
        {isPublished ? (
          <Button
            type="button"
            size="sm"
            variant="outline"
            onClick={() => setConfirmDraft(true)}
            disabled={patch.isPending}
          >
            В черновик
          </Button>
        ) : (
          <Button
            type="button"
            size="sm"
            variant="outline"
            onClick={() => setConfirmPublish(true)}
            disabled={patch.isPending}
          >
            Опубликовать
          </Button>
        )}
      </div>

      <ConfirmDialog
        open={confirmPublish}
        onOpenChange={setConfirmPublish}
        title="Опубликовать резюме?"
        description="Опубликованное резюме можно сделать основным и подавать вакансии."
        confirmLabel={patch.isPending ? 'Публикация…' : 'Опубликовать'}
        confirmDisabled={patch.isPending}
        onConfirm={onPublish}
      />
      <ConfirmDialog
        open={confirmDraft}
        onOpenChange={setConfirmDraft}
        title="Перевести в черновик?"
        description={
          isDefault
            ? 'Резюме перестанет быть основным и его нельзя будет подать.'
            : 'Резюме нельзя будет подать в вакансию, пока оно в черновике.'
        }
        confirmLabel={patch.isPending ? 'Сохранение…' : 'В черновик'}
        confirmDisabled={patch.isPending}
        onConfirm={onDraft}
      />
    </li>
  )
}

function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel,
  confirmDisabled,
  onConfirm,
}: {
  open: boolean
  onOpenChange: (next: boolean) => void
  title: string
  description: string
  confirmLabel: string
  confirmDisabled: boolean
  onConfirm: () => void
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={confirmDisabled}
          >
            Отмена
          </Button>
          <Button type="button" onClick={onConfirm} disabled={confirmDisabled}>
            {confirmLabel}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
