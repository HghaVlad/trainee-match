import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import type { AnalyticsSummary } from '@/api/generated/application/schemas'

interface Props {
  summary: AnalyticsSummary
}

function formatPercent(ratio: number): string {
  if (!Number.isFinite(ratio)) return '—'
  return `${(ratio * 100).toFixed(1)}%`
}

interface KpiItem {
  label: string
  value: string
  hint?: string
}

export function SummaryCards({ summary }: Props) {
  const items: KpiItem[] = [
    { label: 'Всего откликов', value: String(summary.totalApplications) },
    {
      label: 'Активные',
      value: String(summary.activeApplications),
      hint: 'submitted + seen + interview',
    },
    { label: 'Отправлено', value: String(summary.submittedCount) },
    { label: 'Просмотрено', value: String(summary.seenCount) },
    { label: 'Интервью', value: String(summary.interviewCount) },
    { label: 'Офферы', value: String(summary.offerCount) },
    { label: 'Отказы', value: String(summary.rejectedCount) },
    { label: 'Отозваны', value: String(summary.withdrawnCount) },
    {
      label: 'Конверсия в интервью',
      value: formatPercent(summary.conversionToInterview),
    },
    {
      label: 'Конверсия в оффер',
      value: formatPercent(summary.conversionToOffer),
    },
  ]

  return (
    <div className="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-5">
      {items.map((it) => (
        <Card key={it.label}>
          <CardHeader className="p-4 pb-1">
            <CardTitle className="text-xs font-medium text-muted-foreground">
              {it.label}
            </CardTitle>
          </CardHeader>
          <CardContent className="p-4 pt-0">
            <div className="text-2xl font-semibold tabular-nums">
              {it.value}
            </div>
            {it.hint && (
              <div className="mt-1 text-xs text-muted-foreground">{it.hint}</div>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
