import { useMemo } from 'react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
import {
  type ApplicationStatus,
  type HrSortQueryParameter,
} from '@/api/generated/application/schemas'
import { useGetCompaniesCompanyIdVacancies } from '@/api/generated/company/vacancy/vacancy'
import { APPLICATION_STATUS_LABEL } from '@/shared/ui/ApplicationStatusBadge'
import {
  ALL_STATUSES,
  DEFAULT_FILTERS,
  SORT_LABEL,
  type FiltersValue,
} from './filters'

interface Props {
  companyId: string
  value: FiltersValue
  onChange: (next: FiltersValue) => void
  hideVacancyFilter?: boolean
}

export function ApplicationsListFilters({
  companyId,
  value,
  onChange,
  hideVacancyFilter,
}: Props) {
  const vacanciesQ = useGetCompaniesCompanyIdVacancies(
    companyId,
    { limit: 100 },
    { query: { enabled: !hideVacancyFilter && !!companyId } },
  )

  const vacancyOptions = useMemo(
    () => vacanciesQ.data?.vacancies ?? [],
    [vacanciesQ.data],
  )

  function toggleStatus(s: ApplicationStatus, checked: boolean) {
    const next = checked
      ? Array.from(new Set([...value.statuses, s]))
      : value.statuses.filter((x) => x !== s)
    onChange({ ...value, statuses: next })
  }

  function reset() {
    onChange({ ...DEFAULT_FILTERS })
  }

  const statusButtonLabel =
    value.statuses.length === 0
      ? 'Все статусы'
      : `Статусы: ${value.statuses.length}`

  return (
    <div className="flex flex-wrap items-end gap-3">
      <div className="flex flex-col gap-1">
        <span className="text-xs text-muted-foreground">Статус</span>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="min-w-44 justify-between">
              {statusButtonLabel}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start">
            <DropdownMenuLabel>Фильтр по статусам</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {ALL_STATUSES.map((s) => (
              <DropdownMenuCheckboxItem
                key={s}
                checked={value.statuses.includes(s)}
                onCheckedChange={(c) => toggleStatus(s, Boolean(c))}
                onSelect={(e) => e.preventDefault()}
              >
                {APPLICATION_STATUS_LABEL[s]}
              </DropdownMenuCheckboxItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {!hideVacancyFilter && (
        <div className="flex flex-col gap-1">
          <span className="text-xs text-muted-foreground">Вакансия</span>
          <Select
            value={value.vacancyId ?? '__all__'}
            onValueChange={(v) =>
              onChange({
                ...value,
                vacancyId: v === '__all__' ? undefined : v,
              })
            }
          >
            <SelectTrigger className="w-64">
              <SelectValue placeholder="Все вакансии" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="__all__">Все вакансии</SelectItem>
              {vacancyOptions.map((v) => (
                <SelectItem key={v.id} value={v.id ?? ''}>
                  {v.title || v.id}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}

      <div className="flex flex-col gap-1">
        <span className="text-xs text-muted-foreground">Создано с</span>
        <Input
          type="date"
          value={value.createdFrom ?? ''}
          onChange={(e) =>
            onChange({
              ...value,
              createdFrom: e.target.value || undefined,
            })
          }
        />
      </div>
      <div className="flex flex-col gap-1">
        <span className="text-xs text-muted-foreground">Создано до</span>
        <Input
          type="date"
          value={value.createdTo ?? ''}
          onChange={(e) =>
            onChange({
              ...value,
              createdTo: e.target.value || undefined,
            })
          }
        />
      </div>

      <div className="flex flex-col gap-1">
        <span className="text-xs text-muted-foreground">Сортировка</span>
        <Select
          value={value.sort}
          onValueChange={(v) =>
            onChange({ ...value, sort: v as HrSortQueryParameter })
          }
        >
          <SelectTrigger className="w-56">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {(
              Object.keys(SORT_LABEL) as HrSortQueryParameter[]
            ).map((s) => (
              <SelectItem key={s} value={s}>
                {SORT_LABEL[s]}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <Button type="button" variant="ghost" onClick={reset}>
        Сбросить
      </Button>
    </div>
  )
}
