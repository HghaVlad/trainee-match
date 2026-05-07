import { useState } from 'react'
import { Link, useSearchParams } from 'react-router'
import { useGetVacancies } from '@/api/generated/company/vacancy/vacancy'
import type { GetVacanciesParams } from '@/api/generated/company/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Badge } from '@/shared/ui/badge'
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

const SALARY_STEP = 1000
const HOURS_MIN = 1
const HOURS_MAX = 168
const DURATION_MIN = 1
const DURATION_MAX = 730

const WORK_FORMAT_LABEL: Record<string, string> = {
  remote: 'Удалёнка',
  office: 'Офис',
  onsite: 'Офис',
  hybrid: 'Гибрид',
}

function dash(value: unknown, suffix?: string): string {
  if (value === undefined || value === null || value === '') return '—'
  return suffix ? `${value} ${suffix}` : String(value)
}

function salaryRange(from?: number, to?: number): string {
  if (typeof from !== 'number' && typeof to !== 'number') return '—'
  if (typeof from === 'number' && typeof to === 'number') {
    if (from === to) return `${from.toLocaleString('ru-RU')} ₽`
    return `${from.toLocaleString('ru-RU')} – ${to.toLocaleString('ru-RU')} ₽`
  }
  if (typeof from === 'number') return `от ${from.toLocaleString('ru-RU')} ₽`
  return `до ${(to as number).toLocaleString('ru-RU')} ₽`
}

function clampNumeric(
  raw: string,
  { min, max, step }: { min?: number; max?: number; step?: number },
): string {
  if (raw === '') return ''
  let n = Number(raw)
  if (Number.isNaN(n)) return ''
  if (typeof min === 'number' && n < min) n = min
  if (typeof max === 'number' && n > max) n = max
  if (step && step > 1) n = Math.round(n / step) * step
  return String(n)
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

function filterErrors(f: FilterState): string[] {
  const errors: string[] = []
  const sMin = Number(f.salaryMin)
  const sMax = Number(f.salaryMax)
  if (f.salaryMin && f.salaryMax && sMin > sMax) {
    errors.push('Зарплата «от» больше «до»')
  }
  const hMin = Number(f.hoursMin)
  const hMax = Number(f.hoursMax)
  if (f.hoursMin && f.hoursMax && hMin > hMax) {
    errors.push('Часов «от» больше «до»')
  }
  const dMin = Number(f.durationMin)
  const dMax = Number(f.durationMax)
  if (f.durationMin && f.durationMax && dMin > dMax) {
    errors.push('Длительность «от» больше «до»')
  }
  return errors
}

export default function VacanciesPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [filters, setFilters] = useState<FilterState>(() => ({
    ...EMPTY,
    companyId: searchParams.get('company_id') ?? '',
    city: searchParams.get('city') ?? '',
  }))
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

  const errors = filterErrors(filters)
  const { data, isLoading, error, refetch } = useGetVacancies(
    toParams(filters, cursor),
    { query: { enabled: errors.length === 0 } },
  )

  function reset() {
    setFilters(EMPTY)
    setCursor(undefined)
    setSearchParams({})
  }

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
              min={0}
              step={SALARY_STEP}
              inputMode="numeric"
              value={filters.salaryMin}
              onChange={(e) => update('salaryMin', e.target.value)}
              onBlur={(e) =>
                update(
                  'salaryMin',
                  clampNumeric(e.target.value, { min: 0, step: SALARY_STEP }),
                )
              }
            />
          </div>
          <div>
            <Label htmlFor="f-salary-max">Зарплата до, ₽</Label>
            <Input
              id="f-salary-max"
              type="number"
              min={0}
              step={SALARY_STEP}
              inputMode="numeric"
              value={filters.salaryMax}
              onChange={(e) => update('salaryMax', e.target.value)}
              onBlur={(e) =>
                update(
                  'salaryMax',
                  clampNumeric(e.target.value, { min: 0, step: SALARY_STEP }),
                )
              }
            />
          </div>
          <div>
            <Label htmlFor="f-hours-min">Часов/неделю от</Label>
            <Input
              id="f-hours-min"
              type="number"
              min={HOURS_MIN}
              max={HOURS_MAX}
              step={1}
              inputMode="numeric"
              value={filters.hoursMin}
              onChange={(e) => update('hoursMin', e.target.value)}
              onBlur={(e) =>
                update(
                  'hoursMin',
                  clampNumeric(e.target.value, {
                    min: HOURS_MIN,
                    max: HOURS_MAX,
                    step: 1,
                  }),
                )
              }
            />
          </div>
          <div>
            <Label htmlFor="f-hours-max">Часов/неделю до</Label>
            <Input
              id="f-hours-max"
              type="number"
              min={HOURS_MIN}
              max={HOURS_MAX}
              step={1}
              inputMode="numeric"
              value={filters.hoursMax}
              onChange={(e) => update('hoursMax', e.target.value)}
              onBlur={(e) =>
                update(
                  'hoursMax',
                  clampNumeric(e.target.value, {
                    min: HOURS_MIN,
                    max: HOURS_MAX,
                    step: 1,
                  }),
                )
              }
            />
          </div>
          <div>
            <Label htmlFor="f-duration-min">Длительность от (дней)</Label>
            <Input
              id="f-duration-min"
              type="number"
              min={DURATION_MIN}
              max={DURATION_MAX}
              step={1}
              inputMode="numeric"
              value={filters.durationMin}
              onChange={(e) => update('durationMin', e.target.value)}
              onBlur={(e) =>
                update(
                  'durationMin',
                  clampNumeric(e.target.value, {
                    min: DURATION_MIN,
                    max: DURATION_MAX,
                    step: 1,
                  }),
                )
              }
            />
          </div>
          <div>
            <Label htmlFor="f-duration-max">Длительность до (дней)</Label>
            <Input
              id="f-duration-max"
              type="number"
              min={DURATION_MIN}
              max={DURATION_MAX}
              step={1}
              inputMode="numeric"
              value={filters.durationMax}
              onChange={(e) => update('durationMax', e.target.value)}
              onBlur={(e) =>
                update(
                  'durationMax',
                  clampNumeric(e.target.value, {
                    min: DURATION_MIN,
                    max: DURATION_MAX,
                    step: 1,
                  }),
                )
              }
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

        {errors.length > 0 && (
          <ul className="text-sm text-destructive">
            {errors.map((e) => (
              <li key={e}>{e}</li>
            ))}
          </ul>
        )}

        <div className="flex justify-end">
          <Button type="button" variant="outline" onClick={reset}>
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
        {(data?.vacancies ?? []).map((v) => {
          const target = v.id ? `/vacancies/${v.id}` : undefined
          return (
            <li
              key={v.id ?? `${v.title ?? 'noid'}-${v.companyName ?? ''}`}
              className="rounded-lg border bg-card p-4 space-y-2"
            >
              <div className="flex items-start justify-between gap-4">
                {target ? (
                  <Link
                    to={target}
                    className="text-lg font-medium text-primary underline"
                  >
                    {v.title ?? '—'}
                  </Link>
                ) : (
                  <span className="text-lg font-medium">{v.title ?? '—'}</span>
                )}
                <span className="text-sm text-muted-foreground">
                  {salaryRange(v.salaryFrom, v.salaryTo)}
                </span>
              </div>
              <p className="text-sm text-muted-foreground">
                {dash(v.companyName)} • {dash(v.city)}
              </p>
              <div className="flex flex-wrap gap-2 text-xs">
                {v.workFormat && (
                  <Badge variant="secondary">
                    {WORK_FORMAT_LABEL[v.workFormat] ?? v.workFormat}
                  </Badge>
                )}
                {v.isPaid !== undefined && (
                  <Badge variant={v.isPaid ? 'default' : 'outline'}>
                    {v.isPaid ? 'Оплачиваемая' : 'Без оплаты'}
                  </Badge>
                )}
                {v.employmentType && (
                  <Badge variant="outline">{v.employmentType}</Badge>
                )}
                {v.publishedAt && (
                  <span className="text-xs text-muted-foreground">
                    Опубликована{' '}
                    {new Date(v.publishedAt).toLocaleDateString('ru-RU')}
                  </span>
                )}
              </div>
            </li>
          )
        })}
      </ul>

      {data?.nextCursor && (
        <Button variant="outline" onClick={() => setCursor(data.nextCursor)}>
          Загрузить ещё
        </Button>
      )}
    </div>
  )
}
