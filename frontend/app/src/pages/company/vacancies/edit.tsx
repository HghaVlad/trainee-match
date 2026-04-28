import { useState } from 'react'
import { Link, Navigate, useParams } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { useToast } from '@/shared/hooks/use-toast'
import {
  useGetCompaniesCompanyIdVacanciesVacancyId,
  usePatchCompaniesCompanyIdVacanciesVacancyId,
  getGetCompaniesCompanyIdVacanciesQueryKey,
  getGetCompaniesCompanyIdVacanciesVacancyIdQueryKey,
} from '@/api/generated/company/vacancy/vacancy'
import type {
  DtoVacancyCreateRequestWorkFormat,
  DtoVacancyUpdateRequestEmploymentType,
} from '@/api/generated/company/schemas'
import { AppError } from '@/shared/api/http/client'
import { useSession } from '@/shared/session/useSession'
import {
  VacancyActions,
  VacancyForm,
  VacancyStatusBadge,
  type VacancyFormPayload,
} from '@/features/company-vacancies'

export default function CompanyVacancyEditPage() {
  const { companyId, vacancyId } = useParams<{
    companyId: string
    vacancyId: string
  }>()
  if (!companyId) return <Navigate to="/company" replace />
  if (!vacancyId)
    return <Navigate to={`/company/${companyId}/vacancies`} replace />
  return <EditView companyId={companyId} vacancyId={vacancyId} />
}

function EditView({
  companyId,
  vacancyId,
}: {
  companyId: string
  vacancyId: string
}) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const { companies } = useSession()
  const isAdmin =
    companies.find((c) => c.id === companyId)?.role === 'admin'

  const detail = useGetCompaniesCompanyIdVacanciesVacancyId(
    companyId,
    vacancyId,
  )
  const updateMut = usePatchCompaniesCompanyIdVacanciesVacancyId()
  const [serverError, setServerError] = useState<string | null>(null)

  if (detail.isLoading) return <LoadingState />
  if (detail.isError || !detail.data) {
    return <ErrorState onRetry={() => detail.refetch()} />
  }

  const vacancy = detail.data

  async function onSubmit(payload: VacancyFormPayload) {
    setServerError(null)
    try {
      await updateMut.mutateAsync({
        companyId,
        vacancyId,
        data: {
          title: payload.title,
          description: payload.description,
          city: payload.city,
          workFormat: payload.workFormat,
          employmentType: payload.employmentType,
          salaryFrom: payload.salaryFrom,
          salaryTo: payload.salaryTo,
          isPaid: payload.isPaid,
        },
      })
      await Promise.all([
        qc.invalidateQueries({
          queryKey: getGetCompaniesCompanyIdVacanciesVacancyIdQueryKey(
            companyId,
            vacancyId,
          ),
        }),
        qc.invalidateQueries({
          queryKey: getGetCompaniesCompanyIdVacanciesQueryKey(companyId),
        }),
      ])
      toast({ title: 'Вакансия сохранена' })
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось сохранить вакансию'
          : 'Не удалось сохранить вакансию'
      setServerError(msg)
      toast({
        title: 'Ошибка',
        description: msg,
        variant: 'destructive',
      })
    }
  }

  return (
    <div className="mx-auto max-w-3xl space-y-4 p-6">
      <Link
        to={`/company/${companyId}/vacancies`}
        className="text-sm text-muted-foreground underline"
      >
        ← Все вакансии
      </Link>
      <Card>
        <CardHeader className="flex-row items-start justify-between gap-4">
          <div className="space-y-1">
            <CardTitle>{vacancy.title || 'Без названия'}</CardTitle>
            <CardDescription className="flex items-center gap-2">
              <VacancyStatusBadge status={vacancy.status} />
              {vacancy.publishedAt && (
                <span>
                  Опубликована{' '}
                  {new Date(vacancy.publishedAt).toLocaleDateString()}
                </span>
              )}
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          <VacancyForm
            mode="edit"
            isSubmitting={updateMut.isPending}
            serverError={serverError}
            onSubmit={onSubmit}
            defaultValues={{
              title: vacancy.title,
              description: vacancy.description,
              city: vacancy.city,
              workFormat: vacancy.workFormat as
                | DtoVacancyCreateRequestWorkFormat
                | undefined,
              employmentType: vacancy.employmentType as
                | DtoVacancyUpdateRequestEmploymentType
                | undefined,
              salaryFrom: vacancy.salaryFrom,
              salaryTo: vacancy.salaryTo,
              isPaid: vacancy.isPaid,
            }}
          />
          <VacancyActions
            companyId={companyId}
            vacancyId={vacancyId}
            status={vacancy.status}
            isAdmin={isAdmin}
            variant="detail"
          />
        </CardContent>
      </Card>
    </div>
  )
}
