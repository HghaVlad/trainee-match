import { useState } from 'react'
import { Link } from 'react-router'
import { useGetVacancies } from '@/api/generated/company/vacancy/vacancy'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'

export default function VacanciesPage() {
  const [cursor, setCursor] = useState<string | undefined>(undefined)
  const [salaryMin, setSalaryMin] = useState<string>('')
  const { data, isLoading, error, refetch } = useGetVacancies({
    limit: 20,
    cursor,
    ...(salaryMin ? { salary_min: Number(salaryMin) } : {}),
  })

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Вакансии</h1>
      <div className="flex gap-2">
        <Input
          type="number"
          placeholder="Зарплата от"
          value={salaryMin}
          onChange={(e) => {
            setCursor(undefined)
            setSalaryMin(e.target.value)
          }}
          className="max-w-xs"
        />
      </div>

      {isLoading && <LoadingState />}
      {error && <ErrorState onRetry={() => refetch()} />}
      {!isLoading && !error && (data?.vacancies?.length ?? 0) === 0 && (
        <EmptyState title="Вакансии не найдены" />
      )}

      <ul className="space-y-2">
        {(data?.vacancies ?? []).map((v) => (
          <li key={v.id} className="rounded-lg border bg-card p-4">
            <Link
              to={`/vacancies/${v.id ?? ''}`}
              className="text-lg font-medium text-primary underline"
            >
              {v.title ?? '—'}
            </Link>
            <p className="text-sm text-muted-foreground">
              {v.companyName ?? '—'} • {v.city ?? '—'}
            </p>
            {(v.salaryFrom ?? v.salaryTo) && (
              <p className="text-sm">
                {v.salaryFrom ?? ''}
                {v.salaryFrom && v.salaryTo ? '–' : ''}
                {v.salaryTo ?? ''} ₽
              </p>
            )}
          </li>
        ))}
      </ul>

      {data?.nextCursor && (
        <Button variant="outline" onClick={() => setCursor(data.nextCursor)}>
          Загрузить ещё
        </Button>
      )}
    </div>
  )
}
