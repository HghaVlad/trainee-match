import { useState } from 'react'
import { Link, Navigate, useParams } from 'react-router'
import { useGetCompaniesCompanyIdVacanciesVacancyId } from '@/api/generated/company/vacancy/vacancy'
import {
  useGetVacancyAnalyticsSummary,
  useGetVacancyStatusFunnel,
  useGetVacancyDynamics,
} from '@/api/generated/application/application-analytics/application-analytics'
import { IntervalQueryParameter } from '@/api/generated/application/schemas'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { useToast } from '@/shared/hooks/use-toast'
import {
  AnalyticsDateRange,
  type AnalyticsRangeValue,
  SummaryCards,
  StatusFunnelChart,
  DynamicsChart,
} from '@/features/analytics'
import { AppError } from '@/shared/api/http/client'

export default function VacancyAnalyticsPage() {
  const { companyId, vacancyId } = useParams<{
    companyId: string
    vacancyId: string
  }>()
  if (!companyId) return <Navigate to="/company" replace />
  if (!vacancyId)
    return <Navigate to={`/company/${companyId}/vacancies`} replace />
  return <View companyId={companyId} vacancyId={vacancyId} />
}

function describeError(e: unknown): string {
  if (e instanceof AppError) return e.message || 'Не удалось загрузить данные'
  return 'Не удалось загрузить данные'
}

function View({
  companyId,
  vacancyId,
}: {
  companyId: string
  vacancyId: string
}) {
  const { toast } = useToast()
  const [range, setRange] = useState<AnalyticsRangeValue>({
    interval: IntervalQueryParameter.day,
  })

  const detailQ = useGetCompaniesCompanyIdVacanciesVacancyId(
    companyId,
    vacancyId,
  )

  const summaryParams = {
    createdFrom: range.createdFrom,
    createdTo: range.createdTo,
  }
  const dynamicsParams = {
    createdFrom: range.createdFrom,
    createdTo: range.createdTo,
    interval: range.interval,
  }

  const summaryQ = useGetVacancyAnalyticsSummary(vacancyId, summaryParams)
  const funnelQ = useGetVacancyStatusFunnel(vacancyId, summaryParams)
  const dynamicsQ = useGetVacancyDynamics(vacancyId, dynamicsParams)

  function notify(e: unknown) {
    toast({
      title: 'Ошибка',
      description: describeError(e),
      variant: 'destructive',
    })
  }

  const title = detailQ.data?.title || 'Вакансия'

  return (
    <div className="mx-auto max-w-6xl space-y-6 p-6">
      <div className="space-y-1">
        <Link
          to={`/company/${companyId}/vacancies/${vacancyId}`}
          className="text-sm text-muted-foreground underline"
        >
          ← К вакансии
        </Link>
        <h1 className="text-2xl font-bold">Аналитика: {title}</h1>
        <p className="text-sm text-muted-foreground">
          Сводка, воронка и динамика откликов по вакансии.
        </p>
      </div>

      <Card>
        <CardHeader className="p-4 pb-2">
          <CardTitle className="text-base">Период</CardTitle>
        </CardHeader>
        <CardContent className="p-4 pt-0">
          <AnalyticsDateRange value={range} onChange={setRange} showInterval />
        </CardContent>
      </Card>

      <section className="space-y-2">
        <h2 className="text-lg font-semibold">Сводка</h2>
        {summaryQ.isLoading ? (
          <LoadingState />
        ) : summaryQ.isError || !summaryQ.data ? (
          <ErrorState
            onRetry={() => {
              summaryQ.refetch().catch(notify)
            }}
          />
        ) : (
          <SummaryCards summary={summaryQ.data.data} />
        )}
      </section>

      <section className="space-y-2">
        <h2 className="text-lg font-semibold">Воронка статусов</h2>
        <Card>
          <CardContent className="p-4">
            {funnelQ.isLoading ? (
              <LoadingState />
            ) : funnelQ.isError || !funnelQ.data ? (
              <ErrorState
                onRetry={() => {
                  funnelQ.refetch().catch(notify)
                }}
              />
            ) : (
              <StatusFunnelChart funnel={funnelQ.data.data} />
            )}
          </CardContent>
        </Card>
      </section>

      <section className="space-y-2">
        <h2 className="text-lg font-semibold">Динамика</h2>
        <Card>
          <CardContent className="p-4">
            {dynamicsQ.isLoading ? (
              <LoadingState />
            ) : dynamicsQ.isError || !dynamicsQ.data ? (
              <ErrorState
                onRetry={() => {
                  dynamicsQ.refetch().catch(notify)
                }}
              />
            ) : (
              <DynamicsChart points={dynamicsQ.data.data} />
            )}
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
