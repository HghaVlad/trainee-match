import * as React from 'react'
import { Link } from 'react-router'
import {
  useGetHrApplication,
  useGetHrApplicationHistory,
} from '@/api/generated/application/hr-applications/hr-applications'
import type { HrAllowedAction } from '@/api/generated/application/schemas'
import type { ResumeData } from '@/api/generated/application/schemas/resumeData'
import { useGetSkillList } from '@/api/generated/candidate/skill/skill'
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

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString('ru-RU')
}

function InfoRow({ label, value }: { label: string; value: string | null | undefined }) {
  if (!value) return null
  return (
    <p>
      <span className="text-muted-foreground">{label}: </span>
      {value}
    </p>
  )
}

const WORK_FORMAT_LABEL: Record<string, string> = {
  remote: 'Удалёнка',
  onsite: 'Офис',
  hybrid: 'Гибрид',
}

function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString()
}

export function ApplicationDetail({ companyId, applicationId }: Props) {
  const detailQ = useGetHrApplication(applicationId)
  const historyQ = useGetHrApplicationHistory(applicationId)
  const [pendingAction, setPendingAction] = React.useState<HrAllowedAction | null>(null)
  const { data: skillCatalog } = useGetSkillList()
  const skillNameMap = React.useMemo(() => {
    const map = new Map<string, string>()
    if (skillCatalog) {
      for (const s of skillCatalog) {
        if (s.id && s.name) map.set(s.id, s.name)
      }
    }
    return map
  }, [skillCatalog])

  if (detailQ.isLoading) return <LoadingState />
  if (detailQ.error || !detailQ.data) {
    return <ErrorState onRetry={() => detailQ.refetch()} />
  }

  const app = detailQ.data.data
  const history = historyQ.data?.data ?? app.statusHistory ?? []
  const resume: ResumeData = app.snapshot.resumeData
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
        <CardContent className="space-y-4 text-sm">
          {/* Personal info */}
          <div className="grid grid-cols-2 gap-x-6 gap-y-1">
            <InfoRow label="Имя" value={`${resume.lastName} ${resume.firstName} ${resume.middleName}`} />
            <InfoRow label="Email" value={resume.email} />
            <InfoRow label="Телефон" value={resume.phone} />
            <InfoRow label="Telegram" value={app.snapshot.telegram} />
            <InfoRow label="Город" value={resume.city} />
            <InfoRow label="Дата рождения" value={formatDate(resume.dateOfBirth)} />
            <InfoRow label="Гражданство" value={resume.citizenship} />
            <InfoRow label="Формат работы" value={WORK_FORMAT_LABEL[resume.desiredFormat] ?? resume.desiredFormat} />
            <InfoRow label="Уровень английского" value={resume.englishLevel} />
          </div>

          {/* Education */}
          {resume.education.length > 0 && (
            <div>
              <h4 className="mb-1.5 font-medium">Образование</h4>
              <div className="space-y-2">
                {resume.education.map((edu, i) => (
                  <div key={i} className="rounded-md border p-2.5">
                    <p className="font-medium">{edu.university}</p>
                    <p className="text-muted-foreground">{edu.faculty}, {edu.specialization}</p>
                    <p className="text-muted-foreground">
                      {edu.startYear}–{edu.endYear} &middot; {edu.level} &middot; {edu.format}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Work experience */}
          {resume.workExperiences.length > 0 && (
            <div>
              <h4 className="mb-1.5 font-medium">Опыт работы</h4>
              <div className="space-y-2">
                {resume.workExperiences.map((we, i) => (
                  <div key={i} className="rounded-md border p-2.5">
                    <p className="font-medium">{we.position}</p>
                    <p className="text-muted-foreground">{we.company} &middot; {we.period}</p>
                    {we.responsibilities && (
                      <p className="mt-1 whitespace-pre-wrap text-muted-foreground">{we.responsibilities}</p>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Skills (resolved from skill catalog) */}
          {resume.skillsList.length > 0 && (
            <div>
              <h4 className="mb-1.5 font-medium">Навыки</h4>
              <div className="flex flex-wrap gap-2">
                {resume.skillsList.map((s) => (
                  <span
                    key={s}
                    className="rounded-full border px-2 py-0.5 text-xs"
                  >
                    {skillNameMap.get(s) ?? s}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Additional info */}
          {resume.additionalInfo && (
            <div>
              <h4 className="mb-1 font-medium">Дополнительная информация</h4>
              <p className="whitespace-pre-wrap text-muted-foreground">{resume.additionalInfo}</p>
            </div>
          )}

          {/* Portfolio link */}
          {resume.portfolioLink && (
            <p>
              <span className="text-muted-foreground">Портфолио: </span>
              <a
                href={resume.portfolioLink}
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary underline"
              >
                {resume.portfolioLink}
              </a>
            </p>
          )}

          {/* Empty state */}
          {!resume.lastName && resume.education.length === 0 && resume.workExperiences.length === 0 && (
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
