import { useState } from 'react'
import { Link, Navigate, useNavigate, useParams } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import { useToast } from '@/shared/hooks/use-toast'
import {
  usePostCompaniesCompanyIdVacancies,
  getGetCompaniesCompanyIdVacanciesQueryKey,
} from '@/api/generated/company/vacancy/vacancy'
import { AppError } from '@/shared/api/http/client'
import { VacancyForm, type VacancyFormPayload } from '@/features/company-vacancies'

export default function CompanyVacancyNewPage() {
  const { companyId } = useParams<{ companyId: string }>()
  if (!companyId) return <Navigate to="/company" replace />
  return <CreateView companyId={companyId} />
}

function CreateView({ companyId }: { companyId: string }) {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { toast } = useToast()
  const createMut = usePostCompaniesCompanyIdVacancies()
  const [serverError, setServerError] = useState<string | null>(null)

  async function onSubmit(payload: VacancyFormPayload) {
    setServerError(null)
    try {
      const r = await createMut.mutateAsync({
        companyId,
        data: {
          title: payload.title,
          description: payload.description,
          city: payload.city,
          workFormat: payload.workFormat,
          salaryFrom: payload.salaryFrom,
          salaryTo: payload.salaryTo,
          isPaid: payload.isPaid,
        },
      })
      await qc.invalidateQueries({
        queryKey: getGetCompaniesCompanyIdVacanciesQueryKey(companyId),
      })
      toast({ title: 'Вакансия создана' })
      if (r?.id) {
        navigate(`/company/${companyId}/vacancies/${r.id}`)
      } else {
        navigate(`/company/${companyId}/vacancies`)
      }
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось создать вакансию'
          : 'Не удалось создать вакансию'
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
        <CardHeader>
          <CardTitle>Новая вакансия</CardTitle>
          <CardDescription>
            Заполните основные данные. После создания вы сможете её опубликовать.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <VacancyForm
            mode="create"
            onSubmit={onSubmit}
            isSubmitting={createMut.isPending}
            submitLabel="Создать"
            serverError={serverError}
            onCancel={() => navigate(`/company/${companyId}/vacancies`)}
          />
        </CardContent>
      </Card>
    </div>
  )
}
