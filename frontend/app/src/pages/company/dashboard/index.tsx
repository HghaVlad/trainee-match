import { useState } from 'react'
import { Navigate, useParams } from 'react-router'
import {
  useGetCompanyAnalyticsSummary,
  useGetCompanyStatusFunnel,
  useGetCompanyDynamics,
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

export default function CompanyDashboardPage() {
  const { companyId } = useParams<{ companyId: string }>()
  if (!companyId) return <Navigate to="/company" replace />
  return <Dashboard companyId={companyId} />
}

function describeError(e: unknown): string {
  if (e instanceof AppError) return e.message || 'Не удалось загрузить данные'
  return 'Не удалось загрузить данные'
}

function Dashboard({ companyId }: { companyId: string }) {
  const { toast } = useToast()
  const [range, setRange] = useState<AnalyticsRangeValue>({
    interval: IntervalQueryParameter.day,
  })

  const summaryParams = {
    createdFrom: range.createdFrom,
    createdTo: range.createdTo,
  }
  const dynamicsParams = {
    createdFrom: range.createdFrom,
    createdTo: range.createdTo,
    interval: range.interval,
  }

  const summaryQ = useGetCompanyAnalyticsSummary(companyId, summaryParams)
  const funnelQ = useGetCompanyStatusFunnel(companyId, summaryParams)
  const dynamicsQ = useGetCompanyDynamics(companyId, dynamicsParams)

  function notify(e: unknown) {
    toast({
      title: 'Ошибка',
      description: describeError(e),
      variant: 'destructive',
    })
  }

  return (
    <div className="mx-auto max-w-6xl space-y-6 p-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-bold">Аналитика компании</h1>
        <p className="text-sm text-muted-foreground">
          Сводка, воронка и динамика откликов.
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
