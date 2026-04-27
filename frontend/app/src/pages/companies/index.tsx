import { useState } from 'react'
import { Link } from 'react-router'
import { useGetCompanies } from '@/api/generated/company/company/company'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'

export default function CompaniesPage() {
  const [cursor, setCursor] = useState<string | undefined>(undefined)
  const { data, isLoading, error, refetch } = useGetCompanies({
    limit: 20,
    cursor,
  })

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />

  const items = data?.companies ?? []
  if (items.length === 0) return <EmptyState title="Компании не найдены" />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Компании</h1>
      <ul className="space-y-2">
        {items.map((c) => (
          <li key={c.id} className="rounded-lg border bg-card p-4">
            <Link
              to={`/companies/${c.id ?? ''}`}
              className="text-lg font-medium text-primary underline"
            >
              {c.name ?? '—'}
            </Link>
            <p className="text-sm text-muted-foreground">
              Открытых вакансий: {c.openVacanciesCount ?? 0}
            </p>
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
