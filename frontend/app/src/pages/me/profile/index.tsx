import { useState } from 'react'
import { useGetCandidateMe } from '@/api/generated/candidate/candidate/candidate'
import { CandidateProfileForm } from '@/features/candidate-profile'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { AppError } from '@/shared/api/http/client'

export default function CandidateProfilePage() {
  const [editing, setEditing] = useState(false)
  const { data, isLoading, error, refetch } = useGetCandidateMe({
    query: { retry: false },
  })

  if (isLoading) return <LoadingState />

  const notFound =
    error instanceof AppError ? error.status === 404 : false
  if (error && !notFound) {
    return <ErrorState message="Не удалось загрузить профиль" onRetry={() => refetch()} />
  }

  if (notFound || !data) {
    return (
      <div className="mx-auto max-w-xl p-6 space-y-4">
        <EmptyState
          title="Профиль ещё не создан"
          description="Заполните данные кандидата, чтобы откликаться на вакансии."
        />
        <CandidateProfileForm mode="create" onSuccess={() => setEditing(false)} />
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Мой профиль</h1>
      {editing ? (
        <CandidateProfileForm
          mode="edit"
          initial={data}
          onSuccess={() => setEditing(false)}
        />
      ) : (
        <div className="space-y-2 rounded-lg border bg-card p-4">
          <Field label="Телефон" value={data.phone} />
          <Field label="Telegram" value={data.telegram} />
          <Field label="Город" value={data.city} />
          <Field label="Дата рождения" value={data.birthday} />
          <Button onClick={() => setEditing(true)}>Редактировать</Button>
        </div>
      )}
    </div>
  )
}

function Field({ label, value }: { label: string; value?: string }) {
  return (
    <div className="flex justify-between">
      <span className="text-muted-foreground">{label}:</span>
      <span>{value ?? '—'}</span>
    </div>
  )
}
