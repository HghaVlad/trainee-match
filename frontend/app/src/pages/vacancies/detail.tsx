import { useParams, Link } from 'react-router'
import { useGetVacanciesVacancyId } from '@/api/generated/company/vacancy/vacancy'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { ApplyVacancyButton } from '@/features/applications'

const WORK_FORMAT_LABEL: Record<string, string> = {
  remote: 'Удалёнка',
  office: 'Офис',
  hybrid: 'Гибрид',
}

const EMPLOYMENT_LABEL: Record<string, string> = {
  full_time: 'Полная занятость',
  part_time: 'Частичная занятость',
  internship: 'Стажировка',
  contract: 'Контракт',
}

function formatSalary(from?: number, to?: number): string {
  if (from && to) return `${from.toLocaleString('ru-RU')}–${to.toLocaleString('ru-RU')} ₽`
  if (from) return `от ${from.toLocaleString('ru-RU')} ₽`
  if (to) return `до ${to.toLocaleString('ru-RU')} ₽`
  return 'не указана'
}

function formatDate(iso?: string): string {
  if (!iso) return '—'
  try {
    return new Date(iso).toLocaleDateString('ru-RU')
  } catch {
    return iso
  }
}

function Field({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="grid grid-cols-[200px_1fr] gap-2 border-b py-2 last:border-b-0">
      <dt className="text-sm text-muted-foreground">{label}</dt>
      <dd className="text-sm">{value}</dd>
    </div>
  )
}

export default function VacancyDetailPage() {
  const { vacancyId = '' } = useParams<{ vacancyId: string }>()
  const { data, isLoading, error, refetch } = useGetVacanciesVacancyId(vacancyId, {
    query: { enabled: Boolean(vacancyId) },
  })

  if (isLoading) return <LoadingState />
  if (error || !data) return <ErrorState onRetry={() => refetch()} />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <nav aria-label="breadcrumb" className="text-sm text-muted-foreground">
        <ol className="flex flex-wrap items-center gap-1">
          <li>
            <Link to="/companies" className="underline">Компании</Link>
          </li>
          <li aria-hidden>/</li>
          <li>
            <Link to="/vacancies" className="underline">Вакансии</Link>
          </li>
          <li aria-hidden>/</li>
          <li aria-current="page" className="text-foreground">{data.title ?? '—'}</li>
        </ol>
      </nav>

      <div className="sticky top-0 -mx-2 flex flex-wrap items-start justify-between gap-2 border-b bg-background/95 px-2 py-3 backdrop-blur">
        <div className="min-w-0 flex-1">
          <h1 className="text-2xl font-bold">{data.title ?? '—'}</h1>
          {data.companyId ? (
            <Link
              to={`/companies/${data.companyId}`}
              className="text-sm text-primary underline"
            >
              {data.companyName ?? 'Компания'}
            </Link>
          ) : (
            <p className="text-sm text-muted-foreground">{data.companyName ?? '—'}</p>
          )}
        </div>
        {vacancyId && <ApplyVacancyButton vacancyId={vacancyId} />}
      </div>

      <dl className="rounded-lg border bg-card p-4">
        <Field label="Город" value={data.city ?? '—'} />
        <Field
          label="Формат работы"
          value={data.workFormat ? WORK_FORMAT_LABEL[data.workFormat] ?? data.workFormat : '—'}
        />
        <Field
          label="Тип занятости"
          value={
            data.employmentType
              ? EMPLOYMENT_LABEL[data.employmentType] ?? data.employmentType
              : '—'
          }
        />
        <Field label="Зарплата" value={formatSalary(data.salaryFrom, data.salaryTo)} />
        <Field
          label="Часов в неделю"
          value={
            data.hoursPerWeekFrom || data.hoursPerWeekTo
              ? `${data.hoursPerWeekFrom ?? '—'} – ${data.hoursPerWeekTo ?? '—'}`
              : '—'
          }
        />
        <Field
          label="Длительность (дней)"
          value={
            data.durationFromDays || data.durationToDays
              ? `${data.durationFromDays ?? '—'} – ${data.durationToDays ?? '—'}`
              : '—'
          }
        />
        <Field label="Оплачиваемая" value={data.isPaid ? 'Да' : 'Нет'} />
        <Field label="Гибкий график" value={data.flexibleSchedule ? 'Да' : 'Нет'} />
        <Field
          label="С возможностью трудоустройства"
          value={data.internshipToOffer ? 'Да' : 'Нет'}
        />
        <Field label="Опубликовано" value={formatDate(data.publishedAt)} />
      </dl>

      {data.description && (
        <section className="rounded-lg border bg-card p-4">
          <h2 className="text-lg font-semibold">Описание</h2>
          <div className="prose mt-2 max-w-none whitespace-pre-wrap text-sm">
            {data.description}
          </div>
        </section>
      )}
    </div>
  )
}
