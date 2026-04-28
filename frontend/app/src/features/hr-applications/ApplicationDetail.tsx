import { useState } from 'react'
import { Link } from 'react-router'
import {
  useGetHrApplication,
  useGetHrApplicationHistory,
} from '@/api/generated/application/hr-applications/hr-applications'
import { type HrAllowedAction } from '@/api/generated/application/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Button } from '@/shared/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import { ApplicationStatusBadge } from '@/shared/ui/ApplicationStatusBadge'
import { HistoryTimeline } from './HistoryTimeline'
import { ChangeStatusDialog } from './ChangeStatusDialog'
import { ACTION_LABEL, isDestructiveAction } from './actionMap'

interface Props {
  companyId: string
  applicationId: string
}

interface ResumeView {
  title?: string
  content?: string
  skills: string[]
}

function readResume(data: { [k: string]: unknown }): ResumeView {
  const title = typeof data.title === 'string' ? data.title : undefined
  const content =
    typeof data.content === 'string'
      ? data.content
      : typeof data.description === 'string'
        ? data.description
        : undefined
  const rawSkills = data.skills
  const skills: string[] = Array.isArray(rawSkills)
    ? rawSkills.filter((s): s is string => typeof s === 'string')
    : []
  return { title, content, skills }
}

function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString()
}

export function ApplicationDetail({ companyId, applicationId }: Props) {
  const detailQ = useGetHrApplication(applicationId)
  const historyQ = useGetHrApplicationHistory(applicationId)
  const [pendingAction, setPendingAction] = useState<HrAllowedAction | null>(
    null,
  )

  if (detailQ.isLoading) return <LoadingState />
  if (detailQ.error || !detailQ.data) {
    return <ErrorState onRetry={() => detailQ.refetch()} />
  }

  const app = detailQ.data.data
  const history = historyQ.data?.data ?? app.statusHistory ?? []
  const resume = readResume(app.snapshot.resumeData)
  const allowedActions = app.allowedActions ?? []
  const hasActions = allowedActions.length > 0

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <Link
            to={`/company/${companyId}/applications`}
            className="text-sm text-muted-foreground underline"
          >
            ← Все отклики
          </Link>
          <h1 className="mt-1 text-2xl font-bold">{app.snapshot.fullName}</h1>
          <p className="text-sm text-muted-foreground">
            Вакансия:{' '}
            <Link
              to={`/company/${companyId}/vacancies/${app.vacancyId}`}
              className="text-primary underline"
            >
              {app.vacancyTitle}
            </Link>
          </p>
          <p className="text-xs text-muted-foreground">
            Создан: {formatDateTime(app.createdAt)}
          </p>
        </div>
        <ApplicationStatusBadge status={app.status} />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Кандидат</CardTitle>
        </CardHeader>
        <CardContent className="space-y-1 text-sm">
          <p>
            <span className="text-muted-foreground">Имя: </span>
            {app.snapshot.fullName}
          </p>
          {app.snapshot.email && (
            <p>
              <span className="text-muted-foreground">Email: </span>
              {app.snapshot.email}
            </p>
          )}
          {app.snapshot.telegram && (
            <p>
              <span className="text-muted-foreground">Telegram: </span>
              {app.snapshot.telegram}
            </p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Резюме</CardTitle>
          <CardDescription>Снимок на момент отклика</CardDescription>
        </CardHeader>
        <CardContent className="space-y-3 text-sm">
          {resume.title && <p className="text-base font-medium">{resume.title}</p>}
          {resume.content && (
            <p className="whitespace-pre-wrap">{resume.content}</p>
          )}
          {resume.skills.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {resume.skills.map((s) => (
                <span
                  key={s}
                  className="rounded-full border px-2 py-0.5 text-xs"
                >
                  {s}
                </span>
              ))}
            </div>
          )}
          {!resume.title && !resume.content && resume.skills.length === 0 && (
            <p className="text-muted-foreground">Нет данных резюме.</p>
          )}
        </CardContent>
      </Card>

      {app.coverLetter && (
        <Card>
          <CardHeader>
            <CardTitle>Сопроводительное письмо</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="whitespace-pre-wrap text-sm">{app.coverLetter}</p>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Действия</CardTitle>
          <CardDescription>
            {hasActions
              ? 'Доступные переходы статуса.'
              : 'Финальный статус: дальнейшие переходы недоступны.'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {hasActions ? (
            <div className="flex flex-wrap gap-2">
              {allowedActions.map((a) => (
                <Button
                  key={a}
                  variant={isDestructiveAction(a) ? 'destructive' : 'default'}
                  onClick={() => setPendingAction(a)}
                >
                  {ACTION_LABEL[a]}
                </Button>
              ))}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">Нет доступных действий.</p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>История</CardTitle>
        </CardHeader>
        <CardContent>
          <HistoryTimeline items={history} isLoading={historyQ.isLoading} />
        </CardContent>
      </Card>

      <ChangeStatusDialog
        open={pendingAction !== null}
        onOpenChange={(next) => {
          if (!next) setPendingAction(null)
        }}
        applicationId={app.id}
        companyId={companyId}
        vacancyId={app.vacancyId}
        action={pendingAction}
      />
    </div>
  )
}
