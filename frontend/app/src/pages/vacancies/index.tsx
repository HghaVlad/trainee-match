import { useEffect, useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router'
import { useQueryClient, keepPreviousData } from '@tanstack/react-query'
import { useGetVacanciesSearch } from '@/api/generated/company/vacancy/vacancy'
import { usePatchAdminVacanciesIdModeration } from '@/api/generated/company/admin-vacancy/admin-vacancy'
import type { GetVacanciesSearchParams } from '@/api/generated/company/schemas'
import type { DtoVacancyListItemResponse } from '@/api/generated/company/schemas'
import { SearchIcon, XIcon } from 'lucide-react'
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
import { useDebouncedValue } from '@/shared/hooks/useDebouncedValue'
import { useSession } from '@/shared/session/useSession'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'

const UUID_RE =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

function isValidUuid(v: string): boolean {
  return UUID_RE.test(v.trim())
}

interface FilterState {
  searchQuery: string
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
  employmentType: string[]
  workFormat: string[]
  cities: string[]
  companyIds: string[]
}

const EMPTY: FilterState = {
  searchQuery: '',
  order: 'relevance',
  salaryMin: '',
  salaryMax: '',
  hoursMin: '',
  hoursMax: '',
  durationMin: '',
  durationMax: '',
  isPaid: false,
  internshipToOffer: false,
  flexibleSchedule: false,
  employmentType: [],
  workFormat: [],
  cities: [],
  companyIds: [],
}

const SALARY_STEP = 1000
const HOURS_MIN = 1
const HOURS_MAX = 80
const DURATION_MIN = 1
const DURATION_MAX = 1800

const WORK_FORMAT_LABEL: Record<string, string> = {
  remote: 'Удалёнка',
  office: 'Офис',
  onsite: 'Офис',
  hybrid: 'Гибрид',
}

const EMPLOYMENT_TYPE_LABEL: Record<string, string> = {
  internship: 'Стажировка',
  full_time: 'Полная занятость',
  part_time: 'Частичная занятость',
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

function toDigits(s: string, max?: number): string {
  const d = s.replace(/\D/g, '').replace(/^0+/, '')
  if (d === '') return ''
  const n = Number(d)
  if (max !== undefined && n > max) return String(max)
  return d
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

function toParams(f: FilterState, cursor: string | undefined): GetVacanciesSearchParams {
  const p: GetVacanciesSearchParams = { limit: 20 }
  if (cursor) p.cursor = cursor
  if (f.searchQuery) p.query = f.searchQuery
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
  if (f.cities.length > 0) p.city = f.cities
  const validIds = f.companyIds.filter(isValidUuid)
  if (validIds.length > 0) p.company_id = validIds
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
  const invalidIds = f.companyIds.filter((id) => !isValidUuid(id))
  if (invalidIds.length > 0) {
    errors.push(`ID компании должен быть в формате UUID (${invalidIds.join(', ')})`)
  }
  return errors
}

export default function VacanciesPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [filters, setFilters] = useState<FilterState>(() => ({
    ...EMPTY,
    companyIds: searchParams.getAll('company_id'),
    cities: searchParams.getAll('city'),
  }))
  const [cursor, setCursor] = useState<string | undefined>(undefined)
  const { user } = useSession()
  const { toast } = useToast()
  const qc = useQueryClient()
  const archive = usePatchAdminVacanciesIdModeration()
  const isPlatformAdmin = user?.role === 'admin'

  async function onArchive(vacancyId: string, vacancyTitle?: string) {
    const removed = allVacancies.find((v) => v.id === vacancyId)
    setAllVacancies((prev) => prev.filter((v) => v.id !== vacancyId))
    try {
      await archive.mutateAsync({ id: vacancyId, data: { status: 'hidden' } })
      toast({
        title: `Скрыто: ${vacancyTitle ?? vacancyId}`,
        description: `ID: ${vacancyId}`,
      })
    } catch (e) {
      if (removed) setAllVacancies((prev) => [removed, ...prev])
      const msg = e instanceof AppError ? e.message : 'Не удалось скрыть вакансию'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  function update<K extends keyof FilterState>(key: K, value: FilterState[K]) {
    setCursor(undefined)
    setFilters((prev) => ({ ...prev, [key]: value }))
  }

  function addToList(key: 'cities' | 'companyIds', value: string) {
    const trimmed = value.trim()
    if (!trimmed) return
    setCursor(undefined)
    setFilters((prev) =>
      prev[key].includes(trimmed)
        ? prev
        : { ...prev, [key]: [...prev[key], trimmed] },
    )
  }

  function removeFromList(key: 'cities' | 'companyIds', value: string) {
    setCursor(undefined)
    setFilters((prev) => ({
      ...prev,
      [key]: prev[key].filter((v) => v !== value),
    }))
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

  function toggleEmploymentType(value: string) {
    setCursor(undefined)
    setFilters((prev) => ({
      ...prev,
      employmentType: prev.employmentType.includes(value)
        ? prev.employmentType.filter((v) => v !== value)
        : [...prev.employmentType, value],
    }))
  }

  const errors = filterErrors(filters)
  const debouncedFilters = useDebouncedValue(filters, 400)

  const { data, isLoading, error, refetch, isFetching } = useGetVacanciesSearch(
    toParams(debouncedFilters, cursor),
    {
      query: {
        enabled: filterErrors(debouncedFilters).length === 0,
        placeholderData: keepPreviousData,
      },
    },
  )

  const [allVacancies, setAllVacancies] = useState<DtoVacancyListItemResponse[]>([])

  useEffect(() => {
    const items = data?.vacancies
    if (!items) return

    if (cursor) {
      setAllVacancies((prev) => {
        const existingIds = new Set(prev.map((v) => v.id))
        const newItems = items.filter((v) => !existingIds.has(v.id))
        return newItems.length > 0 ? [...prev, ...newItems] : prev
      })
    } else {
      setAllVacancies(items)
    }
  }, [data, cursor])

  const filteredVacancies = useMemo(() => {
    if (filters.employmentType.length === 0) return allVacancies
    return allVacancies.filter((v) =>
      v.employmentType ? filters.employmentType.includes(v.employmentType) : false,
    )
  }, [allVacancies, filters.employmentType])

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
          <div className="sm:col-span-3">
            <Label htmlFor="f-search">Поиск по ключевым словам</Label>
            <div className="relative">
              <SearchIcon className="pointer-events-none absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                id="f-search"
                type="search"
                placeholder="Название вакансии, компания…"
                value={filters.searchQuery}
                onChange={(e) => update('searchQuery', e.target.value)}
                className="pl-8 pr-8"
              />
              {filters.searchQuery && (
                <button
                  type="button"
                  onClick={() => update('searchQuery', '')}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  aria-label="Очистить поиск"
                >
                  <XIcon className="h-4 w-4" />
                </button>
              )}
            </div>
          </div>
          <div>
            <Label htmlFor="f-order">Сортировка</Label>
            <Select value={filters.order} onValueChange={(v) => update('order', v)}>
              <SelectTrigger id="f-order">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="relevance">По релевантности</SelectItem>
                <SelectItem value="published_at_desc">Сначала новые</SelectItem>
                <SelectItem value="salary_desc">Зарплата по убыванию</SelectItem>
                <SelectItem value="salary_asc">Зарплата по возрастанию</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label htmlFor="f-city">Город</Label>
            <ChipsInput
              inputId="f-city"
              ariaLabel="Город"
              values={filters.cities}
              onAdd={(v) => addToList('cities', v)}
              onRemove={(v) => removeFromList('cities', v)}
              placeholder="Москва, Enter"
            />
          </div>
          <div>
            <Label htmlFor="f-company">ID компании</Label>
            <ChipsInput
              inputId="f-company"
              ariaLabel="ID компании"
              values={filters.companyIds}
              onAdd={(v) => addToList('companyIds', v)}
              onRemove={(v) => removeFromList('companyIds', v)}
              placeholder="UUID, Enter"
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
              onChange={(e) => update('salaryMin', toDigits(e.target.value))}
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
              onChange={(e) => update('salaryMax', toDigits(e.target.value))}
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
              onChange={(e) => update('hoursMin', toDigits(e.target.value, HOURS_MAX))}
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
              onChange={(e) => update('hoursMax', toDigits(e.target.value, HOURS_MAX))}
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
onChange={(e) => update('durationMin', toDigits(e.target.value, DURATION_MAX))}
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
onChange={(e) => update('durationMax', toDigits(e.target.value, DURATION_MAX))}
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
              { v: 'onsite', l: 'Офис' },
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

        <div>
          <Label className="block mb-1">Тип занятости</Label>
          <div className="flex flex-wrap gap-3">
            {Object.entries(EMPLOYMENT_TYPE_LABEL).map(([v, l]) => (
              <label key={v} className="flex items-center gap-1 text-sm">
                <input
                  type="checkbox"
                  checked={filters.employmentType.includes(v)}
                  onChange={() => toggleEmploymentType(v)}
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
      {!isLoading && !error && filteredVacancies.length === 0 && (
        <EmptyState title="Вакансии не найдены" />
      )}
      {isFetching && !isLoading && (
        <p className="text-xs text-muted-foreground">Обновление…</p>
      )}

      <ul className="space-y-2">
        {filteredVacancies.map((v) => {
          const target = v.id ? `/vacancies/${v.id}` : undefined
          return (
            <li
              key={v.id ?? `${v.title ?? 'noid'}-${v.companyName ?? ''}`}
              className="rounded-lg border bg-card p-4 space-y-2"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1">
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
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">
                    {salaryRange(v.salaryFrom, v.salaryTo)}
                  </span>
                  {isPlatformAdmin && v.id && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => onArchive(v.id as string, v.title)}
                      disabled={archive.isPending}
                    >
                      {archive.isPending ? '...' : 'Скрыть'}
                    </Button>
                  )}
                </div>
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
                  <Badge variant="outline">{EMPLOYMENT_TYPE_LABEL[v.employmentType] ?? v.employmentType}</Badge>
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

interface ChipsInputProps {
  inputId: string
  ariaLabel: string
  values: string[]
  onAdd: (value: string) => void
  onRemove: (value: string) => void
  placeholder?: string
}

function ChipsInput({
  inputId,
  ariaLabel,
  values,
  onAdd,
  onRemove,
  placeholder,
}: ChipsInputProps) {
  const [draft, setDraft] = useState('')

  function commit() {
    if (draft.trim()) {
      onAdd(draft)
      setDraft('')
    }
  }

  function onKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault()
      commit()
    } else if (e.key === 'Backspace' && draft === '' && values.length > 0) {
      e.preventDefault()
      onRemove(values[values.length - 1] as string)
    }
  }

  return (
    <div className="flex flex-wrap items-center gap-1 rounded-md border border-input bg-background px-2 py-1 focus-within:ring-2 focus-within:ring-ring">
      {values.map((v) => (
        <span
          key={v}
          className="inline-flex items-center gap-1 rounded bg-secondary px-2 py-0.5 text-xs text-secondary-foreground"
        >
          {v}
          <button
            type="button"
            aria-label={`Удалить ${v}`}
            className="text-muted-foreground hover:text-foreground"
            onClick={() => onRemove(v)}
          >
            ×
          </button>
        </span>
      ))}
      <input
        id={inputId}
        aria-label={ariaLabel}
        type="text"
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onKeyDown={onKeyDown}
        onBlur={commit}
        placeholder={values.length === 0 ? placeholder : ''}
        className="flex-1 min-w-[8rem] bg-transparent py-1 text-sm outline-none"
      />
    </div>
  )
}
