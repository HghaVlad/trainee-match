import { useState } from 'react'
import { Link } from 'react-router'
import { useGetVacancies } from '@/api/generated/company/vacancy/vacancy'
import type { GetVacanciesParams } from '@/api/generated/company/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'

interface FilterState {
  order: string
  salaryMin: string
  salaryMax: string
  hoursMin: string
  hoursMax: string
  durationMin: string
  durationMax: string
  isPaid: boolean
  internshipToOffer: boolean
  flexibleSchedule: boolean
  workFormat: string[]
  city: string
  companyId: string
}

const EMPTY: FilterState = {
  order: 'published_at_desc',
  salaryMin: '',
  salaryMax: '',
  hoursMin: '',
  hoursMax: '',
  durationMin: '',
  durationMax: '',
  isPaid: false,
  internshipToOffer: false,
  flexibleSchedule: false,
  workFormat: [],
  city: '',
  companyId: '',
}

function toParams(f: FilterState, cursor: string | undefined): GetVacanciesParams {
  const p: GetVacanciesParams = { limit: 20 }
  if (cursor) p.cursor = cursor
  if (f.order) p.order = f.order
  if (f.salaryMin) p.salary_min = Number(f.salaryMin)
  if (f.salaryMax) p.salary_max = Number(f.salaryMax)
  if (f.hoursMin) p.hours_min = Number(f.hoursMin)
  if (f.hoursMax) p.hours_max = Number(f.hoursMax)
  if (f.durationMin) p.duration_min = Number(f.durationMin)
  if (f.durationMax) p.duration_max = Number(f.durationMax)
  if (f.isPaid) p.is_paid = true
  if (f.internshipToOffer) p.internship_to_offer = true
  if (f.flexibleSchedule) p.flexible_schedule = true
  if (f.workFormat.length > 0) p.work_format = f.workFormat
  if (f.city.trim()) p.city = [f.city.trim()]
  if (f.companyId.trim()) p.company_id = [f.companyId.trim()]
  return p
}

export default function VacanciesPage() {
  const [filters, setFilters] = useState<FilterState>(EMPTY)
  const [cursor, setCursor] = useState<string | undefined>(undefined)

  function update<K extends keyof FilterState>(key: K, value: FilterState[K]) {
    setCursor(undefined)
    setFilters((prev) => ({ ...prev, [key]: value }))
  }

  function toggleWorkFormat(value: string) {
    setCursor(undefined)
    setFilters((prev) => ({
      ...prev,
      workFormat: prev.workFormat.includes(value)
        ? prev.workFormat.filter((v) => v !== value)
        : [...prev.workFormat, value],
    }))
  }

  const { data, isLoading, error, refetch } = useGetVacancies(toParams(filters, cursor))

  return (
    <div className="mx-auto max-w-5xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Вакансии</h1>

      <section className="rounded-lg border bg-card p-4 space-y-4">
        <h2 className="text-sm font-semibold text-muted-foreground">Фильтры</h2>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <div>
            <Label htmlFor="f-order">Сортировка</Label>
            <Select value={filters.order} onValueChange={(v) => update('order', v)}>
              <SelectTrigger id="f-order">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="published_at_desc">Сначала новые</SelectItem>
                <SelectItem value="salary_desc">Зарплата по убыванию</SelectItem>
                <SelectItem value="salary_asc">Зарплата по возрастанию</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label htmlFor="f-city">Город</Label>
            <Input
              id="f-city"
              value={filters.city}
              onChange={(e) => update('city', e.target.value)}
              placeholder="Москва"
            />
          </div>
          <div>
            <Label htmlFor="f-company">ID компании</Label>
            <Input
              id="f-company"
              value={filters.companyId}
              onChange={(e) => update('companyId', e.target.value)}
              placeholder="UUID"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
          <div>
            <Label htmlFor="f-salary-min">Зарплата от, ₽</Label>
            <Input
              id="f-salary-min"
              type="number"
              value={filters.salaryMin}
              onChange={(e) => update('salaryMin', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="f-salary-max">Зарплата до, ₽</Label>
            <Input
              id="f-salary-max"
              type="number"
              value={filters.salaryMax}
              onChange={(e) => update('salaryMax', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="f-hours-min">Часов/неделю от</Label>
            <Input
              id="f-hours-min"
              type="number"
              value={filters.hoursMin}
              onChange={(e) => update('hoursMin', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="f-hours-max">Часов/неделю до</Label>
            <Input
              id="f-hours-max"
              type="number"
              value={filters.hoursMax}
              onChange={(e) => update('hoursMax', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="f-duration-min">Длительность от (дней)</Label>
            <Input
              id="f-duration-min"
              type="number"
              value={filters.durationMin}
              onChange={(e) => update('durationMin', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="f-duration-max">Длительность до (дней)</Label>
            <Input
              id="f-duration-max"
              type="number"
              value={filters.durationMax}
              onChange={(e) => update('durationMax', e.target.value)}
            />
          </div>
        </div>

        <div>
          <Label className="block mb-1">Формат работы</Label>
          <div className="flex flex-wrap gap-3">
            {[
              { v: 'remote', l: 'Удалёнка' },
              { v: 'office', l: 'Офис' },
              { v: 'hybrid', l: 'Гибрид' },
            ].map(({ v, l }) => (
              <label key={v} className="flex items-center gap-1 text-sm">
                <input
                  type="checkbox"
                  checked={filters.workFormat.includes(v)}
                  onChange={() => toggleWorkFormat(v)}
                />
                {l}
              </label>
            ))}
          </div>
        </div>

        <div className="flex flex-wrap gap-4">
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={filters.isPaid}
              onChange={(e) => update('isPaid', e.target.checked)}
            />
            Только оплачиваемые
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={filters.internshipToOffer}
              onChange={(e) => update('internshipToOffer', e.target.checked)}
            />
            С возможностью трудоустройства
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={filters.flexibleSchedule}
              onChange={(e) => update('flexibleSchedule', e.target.checked)}
            />
            Гибкий график
          </label>
        </div>

        <div className="flex justify-end">
          <Button
            type="button"
            variant="outline"
            onClick={() => {
              setFilters(EMPTY)
              setCursor(undefined)
            }}
          >
            Сбросить
          </Button>
        </div>
      </section>

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
