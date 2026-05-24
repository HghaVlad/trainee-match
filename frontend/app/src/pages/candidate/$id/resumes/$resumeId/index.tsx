import { useMemo, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import { useGetAdminResumesId, usePostAdminResumesIdArchive } from '@/api/generated/candidate/admin/admin'
import { useGetSkillList } from '@/api/generated/candidate/skill/skill'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Button } from '@/shared/ui/button'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'

function formatDate(iso?: string): string {
  if (!iso) return '—'
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString('ru-RU')
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

export default function CandidateResumeDetailPage() {
  const { id = '', resumeId = '' } = useParams<{ id: string; resumeId: string }>()
  const navigate = useNavigate()
  const { toast } = useToast()
  const qc = useQueryClient()
  const archive = usePostAdminResumesIdArchive()
  const { data: skillCatalog } = useGetSkillList()

  const { data, isLoading, error, refetch } = useGetAdminResumesId(resumeId, {
    query: { enabled: Boolean(resumeId) },
  })

  const skillNameMap = useMemo(() => {
    const map = new Map<string, string>()
    if (skillCatalog) {
      for (const s of skillCatalog) {
        if (s.id && s.name) map.set(s.id, s.name)
      }
    }
    return map
  }, [skillCatalog])

  async function onArchive() {
    try {
      await archive.mutateAsync({ id: resumeId })
      await qc.invalidateQueries({ queryKey: ['useGetAdminResumesId'] })
      toast({ title: 'Резюме архивировано' })
      navigate(`/candidates/${id}/resumes`)
    } catch (err) {
      const msg = err instanceof AppError ? err.message : 'Не удалось архивировать резюме'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  if (isLoading) return <LoadingState />
  if (error || !data) return <ErrorState onRetry={() => refetch()} />

  const rd = data.data

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link
          to={`/candidates/${id}/resumes`}
          className="text-sm text-muted-foreground underline"
        >
          ← Все резюме
        </Link>
      </div>

      <div className="flex flex-wrap items-start justify-between gap-3">
        <h1 className="text-2xl font-bold">{data.name ?? `Резюме #${resumeId}`}</h1>
        <Button variant="destructive" onClick={onArchive} disabled={archive.isPending}>
          {archive.isPending ? 'Архивация...' : 'Архивировать'}
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Резюме</CardTitle>
          <CardDescription>Полная информация о резюме</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 text-sm">
          {/* Personal info */}
          <div className="grid grid-cols-2 gap-x-6 gap-y-1">
            <InfoRow label="Имя" value={`${rd?.last_name ?? ''} ${rd?.first_name ?? ''} ${rd?.middle_name ?? ''}`.trim() || null} />
            <InfoRow label="Email" value={rd?.email} />
            <InfoRow label="Телефон" value={rd?.phone} />
            <InfoRow label="Город" value={rd?.city} />
            <InfoRow label="Дата рождения" value={formatDate(rd?.date_of_birth)} />
            <InfoRow label="Гражданство" value={rd?.citizenship} />
            <InfoRow label="Формат работы" value={rd?.desired_format ? (WORK_FORMAT_LABEL[rd.desired_format] ?? rd.desired_format) : null} />
            <InfoRow label="Уровень английского" value={rd?.english_level} />
          </div>

          {/* Education */}
          {rd?.education && rd.education.length > 0 && (
            <div>
              <h4 className="mb-1.5 font-medium">Образование</h4>
              <div className="space-y-2">
                {rd.education.map((edu, i) => (
                  <div key={i} className="rounded-md border p-2.5">
                    <p className="font-medium">{edu.university}</p>
                    <p className="text-muted-foreground">{edu.faculty}, {edu.specialization}</p>
                    <p className="text-muted-foreground">
                      {edu.start_year}–{edu.end_year} &middot; {edu.level} &middot; {edu.format}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Work experience */}
          {rd?.work_experiences && rd.work_experiences.length > 0 && (
            <div>
              <h4 className="mb-1.5 font-medium">Опыт работы</h4>
              <div className="space-y-2">
                {rd.work_experiences.map((we, i) => (
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
          {rd?.skills_list && rd.skills_list.length > 0 && (
            <div>
              <h4 className="mb-1.5 font-medium">Навыки</h4>
              <div className="flex flex-wrap gap-2">
                {rd.skills_list.map((s) => (
                  <span key={s} className="rounded-full border px-2 py-0.5 text-xs">
                    {skillNameMap.get(s) ?? s}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Additional info */}
          {rd?.additional_info && (
            <div>
              <h4 className="mb-1 font-medium">Дополнительная информация</h4>
              <p className="whitespace-pre-wrap text-muted-foreground">{rd.additional_info}</p>
            </div>
          )}

          {/* Portfolio link */}
          {rd?.portfolio_link && (
            <p>
              <span className="text-muted-foreground">Портфолио: </span>
              <a
                href={rd.portfolio_link}
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary underline"
              >
                {rd.portfolio_link}
              </a>
            </p>
          )}

          {/* Empty state */}
          {!rd?.last_name && (!rd?.education || rd.education.length === 0) && (!rd?.work_experiences || rd.work_experiences.length === 0) && (
            <p className="text-muted-foreground">Нет данных резюме.</p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}