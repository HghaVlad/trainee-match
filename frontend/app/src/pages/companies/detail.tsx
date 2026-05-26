import { useEffect, useMemo, useRef, useState } from 'react'
import { useParams, Link } from 'react-router'
import { keepPreviousData } from '@tanstack/react-query'
import { useGetCompaniesId } from '@/api/generated/company/company/company'
import { useGetVacanciesSearch } from '@/api/generated/company/vacancy/vacancy'
import { usePatchAdminCompaniesIdModeration } from '@/api/generated/company/admin-company/admin-company'
import { usePatchAdminVacanciesIdModeration } from '@/api/generated/company/admin-vacancy/admin-vacancy'
import type { DtoVacancyListItemResponse } from '@/api/generated/company/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Input } from '@/shared/ui/input'
import { Button } from '@/shared/ui/button'
import { Label } from '@/shared/ui/label'
import { Badge } from '@/shared/ui/badge'
import { SearchIcon, XIcon } from 'lucide-react'
import { useDebouncedValue } from '@/shared/hooks/useDebouncedValue'
import { useSession } from '@/shared/session/useSession'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'

function dash(value: unknown): string {
  if (value === undefined || value === null || value === '') return '—'
  return String(value)
}

function formatDate(iso?: string): string {
  if (!iso) return '—'
  try {
    return new Date(iso).toLocaleDateString('ru-RU')
  } catch {
    return iso
  }
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
  employmentType: string[]
  workFormat: string[]
  cities: string[]
}

const EMPTY_FILTERS: FilterState = {
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

function buildParams(
  f: FilterState,
  companyId: string,
  query: string | undefined,
  cursor: string | undefined,
) {
  const p: Record<string, unknown> = {
    company_id: [companyId],
    limit: 10,
  }
  if (query) p.query = query
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
  if (f.cities.length > 0) p.city = f.cities
  return p
}

const SALARY_STEP = 1000
const HOURS_MIN = 1
const HOURS_MAX = 80
const DURATION_MIN = 1
const DURATION_MAX = 1800

export default function CompanyDetailPage() {
  const { companyId = '' } = useParams<{ companyId: string }>()
  const id = companyId
  const { data: company, isLoading, error, refetch } = useGetCompaniesId(id, {
    query: { enabled: Boolean(id), retry: false },
  })
  const notFound = error instanceof AppError && error.status === 404
  const { user } = useSession()
  const { toast } = useToast()
  const archive = usePatchAdminCompaniesIdModeration()
  const archiveVacancy = usePatchAdminVacanciesIdModeration()
  const isPlatformAdmin = user?.role === 'admin'

  const [searchQuery, setSearchQuery] = useState('')
  const [cursor, setCursor] = useState<string | undefined>(undefined)
  const inputRef = useRef<HTMLInputElement>(null)
  const debouncedQuery = useDebouncedValue(searchQuery, 300)

  const [filters, setFilters] = useState<FilterState>(EMPTY_FILTERS)

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

  function toggleEmploymentType(value: string) {
    setCursor(undefined)
    setFilters((prev) => ({
      ...prev,
      employmentType: prev.employmentType.includes(value)
        ? prev.employmentType.filter((v) => v !== value)
        : [...prev.employmentType, value],
    }))
  }

  function addCity(value: string) {
    const trimmed = value.trim()
    if (!trimmed) return
    setCursor(undefined)
    setFilters((prev) =>
      prev.cities.includes(trimmed)
        ? prev
        : { ...prev, cities: [...prev.cities, trimmed] },
    )
  }

  function removeCity(value: string) {
    setCursor(undefined)
    setFilters((prev) => ({
      ...prev,
      cities: prev.cities.filter((v) => v !== value),
    }))
  }

  const params = buildParams(filters, id, debouncedQuery || undefined, cursor)

  const vacancyQ = useGetVacanciesSearch(
    params as Parameters<typeof useGetVacanciesSearch>[0],
    { query: { enabled: Boolean(id), placeholderData: keepPreviousData } },
  )

  const [allVacancies, setAllVacancies] = useState<DtoVacancyListItemResponse[]>([])

  useEffect(() => {
    const items = vacancyQ.data?.vacancies
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
  }, [vacancyQ.data, cursor])

  const filteredVacancies = useMemo(() => {
    if (filters.employmentType.length === 0) return allVacancies
    return allVacancies.filter((v) =>
      v.employmentType ? filters.employmentType.includes(v.employmentType) : false,
    )
  }, [allVacancies, filters.employmentType])

  async function onArchiveVacancy(vacancyId: string, vacancyTitle?: string) {
    const removed = allVacancies.find((v) => v.id === vacancyId)
    setAllVacancies((prev) => prev.filter((v) => v.id !== vacancyId))
    try {
      await archiveVacancy.mutateAsync({ id: vacancyId, data: { status: 'hidden' as const } })
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

  function resetFilters() {
    setFilters(EMPTY_FILTERS)
    setCursor(undefined)
    setSearchQuery('')
  }

  if (isLoading) return <LoadingState />
  if (notFound) return <ErrorState title="Компания не найдена" message="Возможно, она была скрыта или удалена." onRetry={() => refetch()} />
  if (error || !company) return <ErrorState onRetry={() => refetch()} />

  const openCount = company.openVacanciesCount ?? 0
  const vacanciesHref = company.id
    ? `/vacancies?company_id=${encodeURIComponent(company.id)}`
    : '/vacancies'

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-6">
      <Link to="/companies" className="text-sm text-muted-foreground underline">
        ← Все компании
      </Link>
      <div className="flex items-start justify-between gap-4">
        <h1 className="text-2xl font-bold">{company.name ?? '—'}</h1>
        {isPlatformAdmin && (
          <Button
            variant="outline"
            size="sm"
            onClick={async () => {
              try {
                await archive.mutateAsync({ id, data: { status: 'hidden' as const } })
                toast({
                  title: `Скрыто: ${company.name ?? id}`,
                  description: `ID: ${id}`,
                })
              } catch (e) {
                const msg = e instanceof AppError ? e.message : 'Не удалось скрыть компанию'
                toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
              }
            }}
            disabled={archive.isPending}
          >
            {archive.isPending ? '...' : 'Скрыть'}
          </Button>
        )}
      </div>

      <dl className="grid grid-cols-1 gap-3 rounded-lg border bg-card p-4 sm:grid-cols-[200px,1fr]">
        <dt className="text-sm font-medium text-muted-foreground">Название</dt>
        <dd>{dash(company.name)}</dd>

        <dt className="text-sm font-medium text-muted-foreground">Сайт</dt>
        <dd>
          {company.website ? (
            <a
              href={company.website}
              target="_blank"
              rel="noreferrer"
              className="text-primary underline"
            >
              {company.website}
            </a>
          ) : (
            '—'
          )}
        </dd>

        <dt className="text-sm font-medium text-muted-foreground">Описание</dt>
        <dd className="whitespace-pre-line">{dash(company.description)}</dd>

        <dt className="text-sm font-medium text-muted-foreground">
          Открытых вакансий
        </dt>
        <dd>
          {openCount > 0 ? (
            <Link to={vacanciesHref} className="text-primary underline">
              {openCount}
            </Link>
          ) : (
            '0'
          )}
        </dd>
      </dl>

      <section className="space-y-3">
        <h2 className="text-xl font-semibold">Вакансии компании</h2>

        <section className="rounded-lg border bg-card p-4 space-y-4">
          <h3 className="text-sm font-semibold text-muted-foreground">Фильтры</h3>

          <div className="relative">
            <SearchIcon className="pointer-events-none absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              ref={inputRef}
              type="search"
              placeholder="Поиск вакансий…"
              value={searchQuery}
              onChange={(e) => {
                setCursor(undefined)
                setAllVacancies([])
                setSearchQuery(e.target.value)
              }}
              className="w-full pl-8 pr-8"
            />
            {searchQuery && (
              <button
                type="button"
                onClick={() => {
                  setCursor(undefined)
                  setAllVacancies([])
                  setSearchQuery('')
                  inputRef.current?.focus()
                }}
                className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                aria-label="Очистить поиск"
              >
                <XIcon className="h-4 w-4" />
              </button>
            )}
          </div>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
            <div>
              <Label>Сортировка</Label>
              <select
                value={filters.order}
                onChange={(e) => update('order', e.target.value)}
                className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm"
              >
                <option value="relevance">По релевантности</option>
                <option value="published_at_desc">Сначала новые</option>
                <option value="salary_desc">Зарплата по убыванию</option>
                <option value="salary_asc">Зарплата по возрастанию</option>
              </select>
            </div>
            <div className="sm:col-span-2">
              <Label>Город</Label>
              <ChipsInput
                values={filters.cities}
                onAdd={addCity}
                onRemove={removeCity}
                placeholder="Москва, Enter"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
            <div>
              <Label>Зарплата от, ₽</Label>
              <Input
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
              <Label>Зарплата до, ₽</Label>
              <Input
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
              <Label>Часов/неделю от</Label>
              <Input
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
                    clampNumeric(e.target.value, { min: HOURS_MIN, max: HOURS_MAX, step: 1 }),
                  )
                }
              />
            </div>
            <div>
              <Label>Часов/неделю до</Label>
              <Input
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
                    clampNumeric(e.target.value, { min: HOURS_MIN, max: HOURS_MAX, step: 1 }),
                  )
                }
              />
            </div>
            <div>
              <Label>Длительность от (дней)</Label>
              <Input
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
                    clampNumeric(e.target.value, { min: DURATION_MIN, max: DURATION_MAX, step: 1 }),
                  )
                }
              />
            </div>
            <div>
              <Label>Длительность до (дней)</Label>
              <Input
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
                    clampNumeric(e.target.value, { min: DURATION_MIN, max: DURATION_MAX, step: 1 }),
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

          <div className="flex justify-end">
            <Button type="button" variant="outline" onClick={resetFilters}>
              Сбросить
            </Button>
          </div>
        </section>

        <div className="min-h-[250px]">
          {vacancyQ.isLoading && allVacancies.length === 0 ? (
            <LoadingState />
          ) : vacancyQ.isError ? (
            <ErrorState onRetry={() => vacancyQ.refetch()} />
          ) : filteredVacancies.length === 0 ? (
            <EmptyState
              title="Вакансий нет"
              description="По заданным параметрам ничего не найдено."
            />
          ) : (
            <ul className="space-y-2">
              {filteredVacancies.map((v) => (
                <li
                  key={v.id}
                  className="rounded-lg border bg-card p-4 space-y-2"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <Link
                        to={`/vacancies/${v.id ?? ''}`}
                        className="text-lg font-medium text-primary underline"
                      >
                        {v.title ?? '—'}
                      </Link>
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      <span className="text-sm text-muted-foreground">
                        {salaryRange(v.salaryFrom, v.salaryTo)}
                      </span>
                      {isPlatformAdmin && v.id && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => onArchiveVacancy(v.id as string, v.title)}
                          disabled={archiveVacancy.isPending}
                        >
                          {archiveVacancy.isPending ? '...' : 'Скрыть'}
                        </Button>
                      )}
                    </div>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {dash(v.city)}
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
                        Опубликована {formatDate(v.publishedAt)}
                      </span>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        {vacancyQ.data?.nextCursor && allVacancies.length > 0 && (
          <Button
            variant="outline"
            onClick={() => setCursor(vacancyQ.data!.nextCursor!)}
            disabled={vacancyQ.isFetching}
          >
            {vacancyQ.isFetching ? 'Загрузка…' : 'Загрузить ещё'}
          </Button>
        )}
      </section>
    </div>
  )
}

interface ChipsInputProps {
  values: string[]
  onAdd: (value: string) => void
  onRemove: (value: string) => void
  placeholder?: string
}

function ChipsInput({ values, onAdd, onRemove, placeholder }: ChipsInputProps) {
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
        aria-label={placeholder ?? 'Добавить'}
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
