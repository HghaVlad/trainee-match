import { useMemo, useState } from 'react'
import { Link, useParams } from 'react-router'
import {
  useGetMyApplication,
  useGetMyApplicationHistory,
} from '@/api/generated/application/candidate-applications/candidate-applications'
import { CandidateAllowedAction } from '@/api/generated/application/schemas'
import type { ResumeData } from '@/api/generated/application/schemas/resumeData'
import { useGetSkillList } from '@/api/generated/candidate/skill/skill'
import { ApplicationStatusBadge } from '@/shared/ui/ApplicationStatusBadge'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Button } from '@/shared/ui/button'
import { STATUS_LABEL, WithdrawApplicationDialog } from '@/features/applications'

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString('ru-RU')
}

function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString()
}

const ROLE_LABEL: Record<string, string> = {
  candidate: 'Вы',
  hr: 'HR',
  system: 'Система',
}

const WORK_FORMAT_LABEL: Record<string, string> = {
  remote: 'Удалёнка',
  onsite: 'Офис',
  hybrid: 'Гибрид',
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

export default function MyApplicationDetailPage() {
  const { applicationId = '' } = useParams<{ applicationId: string }>()
  const detailQ = useGetMyApplication(applicationId)
  const historyQ = useGetMyApplicationHistory(applicationId)
  const { data: skillCatalog } = useGetSkillList()
  const [withdrawOpen, setWithdrawOpen] = useState(false)

  const skillNameMap = useMemo(() => {
    const map = new Map<string, string>()
    if (skillCatalog) {
      for (const s of skillCatalog) {
        if (s.id && s.name) map.set(s.id, s.name)
      }
    }
    return map
  }, [skillCatalog])

  if (detailQ.isLoading) return <LoadingState />
  if (detailQ.error || !detailQ.data)
    return <ErrorState onRetry={() => detailQ.refetch()} />

  const app = detailQ.data.data
  const history = historyQ.data?.data ?? app.statusHistory ?? []
  const resume: ResumeData = app.snapshot.resumeData
  const canWithdraw = (app.allowedActions ?? []).includes(
    CandidateAllowedAction.withdraw,
  )

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <Link
        to="/me/applications"
        className="text-sm text-muted-foreground underline"
      >
        ← Все отклики
      </Link>
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">{app.vacancyTitle}</h1>
          <p className="text-sm text-muted-foreground">{app.companyName}</p>
          <p className="text-xs text-muted-foreground">
            Создан: {formatDateTime(app.createdAt)}
          </p>
        </div>
        <ApplicationStatusBadge status={app.status} />
      </div>

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
          <CardTitle>Резюме</CardTitle>
          <CardDescription>Снимок на момент отклика</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 text-sm">
          {/* Personal info */}
          <div className="grid grid-cols-2 gap-x-6 gap-y-1">
            <InfoRow label="Имя" value={`${resume.lastName} ${resume.firstName} ${resume.middleName}`} />
            <InfoRow label="Email" value={resume.email} />
            <InfoRow label="Телефон" value={resume.phone} />
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

          {/* Skills */}
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

      <Card>
        <CardHeader>
          <CardTitle>История статусов</CardTitle>
        </CardHeader>
        <CardContent>
          {historyQ.isLoading ? (
            <p className="text-sm text-muted-foreground">Загрузка…</p>
          ) : history.length === 0 ? (
            <p className="text-sm text-muted-foreground">История пуста.</p>
          ) : (
            <ul className="space-y-2">
              {history.map((h, i) => (
                <li
                  key={`${h.createdAt}-${i}`}
                  className="flex items-center justify-between gap-3 text-sm"
                >
                  <span>
                    {STATUS_LABEL[h.status] ?? h.status}{' '}
                    <span className="text-muted-foreground">
                      · {ROLE_LABEL[h.changedByRole] ?? h.changedByRole}
                    </span>
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {formatDateTime(h.createdAt)}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>

      {canWithdraw && (
        <div>
          <Button
            type="button"
            variant="destructive"
            onClick={() => setWithdrawOpen(true)}
          >
            Отозвать отклик
          </Button>
          <WithdrawApplicationDialog
            applicationId={app.id}
            open={withdrawOpen}
            onOpenChange={setWithdrawOpen}
          />
        </div>
      )}
    </div>
  )
}
