import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { z } from 'zod'
import { useQueryClient } from '@tanstack/react-query'
import { useGetCompanies } from '@/api/generated/company/company/company'
import {
  useGetCompaniesCompanyIdVacancies,
  useGetCompaniesCompanyIdVacanciesVacancyId,
  usePostCompaniesCompanyIdVacancies,
  usePatchCompaniesCompanyIdVacanciesVacancyId,
  useDeleteCompaniesCompanyIdVacanciesVacancyId,
  usePostCompaniesCompanyIdVacanciesVacancyIdArchive,
  usePostCompaniesCompanyIdVacanciesVacancyIdPublish,
  getGetCompaniesCompanyIdVacanciesQueryKey,
} from '@/api/generated/company/vacancy/vacancy'
import { DtoVacancyCreateRequestWorkFormat } from '@/api/generated/company/schemas'
import { FormWrapper } from '@/shared/ui/Form'
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/shared/ui/form'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import { Button } from '@/shared/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { AppError } from '@/shared/api/http/client'

function useOwnedCompanyId() {
  const { data, isLoading, error } = useGetCompanies({ limit: 1 })
  return {
    companyId: data?.companies?.[0]?.id,
    isLoading,
    error,
  }
}

const vacancySchema = z.object({
  title: z.string().min(2).max(200),
  city: z.string().max(100).optional().or(z.literal('')),
  description: z.string().max(5000).optional().or(z.literal('')),
  workFormat: z.enum([
    DtoVacancyCreateRequestWorkFormat.onsite,
    DtoVacancyCreateRequestWorkFormat.remote,
    DtoVacancyCreateRequestWorkFormat.hybrid,
  ]),
  salaryFrom: z.string().optional().or(z.literal('')),
  salaryTo: z.string().optional().or(z.literal('')),
  isPaid: z.boolean().optional(),
})
type VacancyFormData = z.infer<typeof vacancySchema>

export function CompanyVacanciesPage() {
  const { companyId, isLoading, error } = useOwnedCompanyId()
  const [cursor, setCursor] = useState<string | undefined>(undefined)
  const list = useGetCompaniesCompanyIdVacancies(
    companyId ?? '',
    { limit: 20, cursor },
    { query: { enabled: Boolean(companyId) } },
  )

  if (isLoading) return <LoadingState />
  if (error || !companyId) return <ErrorState message="Сначала создайте профиль компании" />

  const items = list.data?.vacancies ?? []

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Мои вакансии</h1>
        <Link to="/company/vacancies/new">
          <Button>Создать</Button>
        </Link>
      </div>
      {list.isLoading && <LoadingState />}
      {list.error && <ErrorState onRetry={() => list.refetch()} />}
      {items.length === 0 && !list.isLoading && <EmptyState title="Вакансий пока нет" />}
      <ul className="space-y-2">
        {items.map((v) => (
          <li key={v.id} className="rounded-lg border bg-card p-4">
            <Link
              to={`/company/vacancies/${v.id ?? ''}`}
              className="text-lg font-medium text-primary underline"
            >
              {v.title ?? '—'}
            </Link>
            <p className="text-sm text-muted-foreground">
              {v.city ?? '—'} • {v.workFormat ?? ''}
            </p>
          </li>
        ))}
      </ul>
      {list.data?.nextCursor && (
        <Button variant="outline" onClick={() => setCursor(list.data.nextCursor)}>
          Загрузить ещё
        </Button>
      )}
    </div>
  )
}

export function CompanyVacancyEditPage() {
  const { id = '' } = useParams<{ id: string }>()
  const isNew = id === 'new'
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { companyId } = useOwnedCompanyId()
  const [error, setError] = useState<string | null>(null)

  const detail = useGetCompaniesCompanyIdVacanciesVacancyId(
    companyId ?? '',
    isNew ? '' : id,
    { query: { enabled: Boolean(companyId) && !isNew } },
  )
  const createMut = usePostCompaniesCompanyIdVacancies()
  const updateMut = usePatchCompaniesCompanyIdVacanciesVacancyId()
  const deleteMut = useDeleteCompaniesCompanyIdVacanciesVacancyId()
  const archiveMut = usePostCompaniesCompanyIdVacanciesVacancyIdArchive()
  const publishMut = usePostCompaniesCompanyIdVacanciesVacancyIdPublish()

  if (!companyId) return <ErrorState message="Сначала создайте профиль компании" />
  if (!isNew && detail.isLoading) return <LoadingState />
  if (!isNew && (detail.error || !detail.data)) return <ErrorState onRetry={() => detail.refetch()} />

  const initial = !isNew ? detail.data : undefined

  async function onSubmit(values: VacancyFormData) {
    setError(null)
    const payload = {
      title: values.title,
      city: values.city || undefined,
      description: values.description || undefined,
      workFormat: values.workFormat,
      salaryFrom: values.salaryFrom ? Number(values.salaryFrom) : undefined,
      salaryTo: values.salaryTo ? Number(values.salaryTo) : undefined,
      isPaid: values.isPaid,
    }
    try {
      if (isNew) {
        const r = await createMut.mutateAsync({ companyId: companyId!, data: payload })
        await qc.invalidateQueries({
          queryKey: getGetCompaniesCompanyIdVacanciesQueryKey(companyId!),
        })
        if (r?.id) navigate(`/company/vacancies/${r.id}`)
      } else {
        await updateMut.mutateAsync({
          companyId: companyId!,
          vacancyId: id,
          data: payload,
        })
        await detail.refetch()
      }
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Ошибка сохранения')
    }
  }

  async function onAction(kind: 'publish' | 'archive' | 'delete') {
    if (!companyId || isNew) return
    try {
      if (kind === 'publish') await publishMut.mutateAsync({ companyId, vacancyId: id })
      if (kind === 'archive') await archiveMut.mutateAsync({ companyId, vacancyId: id })
      if (kind === 'delete') {
        if (!window.confirm('Удалить вакансию?')) return
        await deleteMut.mutateAsync({ companyId, vacancyId: id })
        navigate('/company/vacancies')
        return
      }
      await detail.refetch()
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Ошибка действия')
    }
  }

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <Link to="/company/vacancies" className="text-sm text-muted-foreground underline">
        ← Все вакансии
      </Link>
      <h1 className="text-2xl font-bold">
        {isNew ? 'Новая вакансия' : initial?.title ?? '—'}
      </h1>
      {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
      <FormWrapper
        schema={vacancySchema}
        defaultValues={{
          title: initial?.title ?? '',
          city: initial?.city ?? '',
          description: initial?.description ?? '',
          workFormat:
            (initial?.workFormat as DtoVacancyCreateRequestWorkFormat | undefined) ??
            DtoVacancyCreateRequestWorkFormat.remote,
          salaryFrom: initial?.salaryFrom ? String(initial.salaryFrom) : '',
          salaryTo: initial?.salaryTo ? String(initial.salaryTo) : '',
          isPaid: initial?.isPaid ?? true,
        }}
        onSubmit={onSubmit}
      >
        {(form) => (
          <div className="space-y-4">
            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Название</FormLabel>
                  <FormControl><Input {...field} /></FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="city"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Город</FormLabel>
                  <FormControl><Input {...field} /></FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="workFormat"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Формат работы</FormLabel>
                  <Select value={field.value} onValueChange={field.onChange}>
                    <FormControl>
                      <SelectTrigger><SelectValue /></SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value={DtoVacancyCreateRequestWorkFormat.onsite}>Офис</SelectItem>
                      <SelectItem value={DtoVacancyCreateRequestWorkFormat.remote}>Удалёнка</SelectItem>
                      <SelectItem value={DtoVacancyCreateRequestWorkFormat.hybrid}>Гибрид</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex gap-2">
              <FormField
                control={form.control}
                name="salaryFrom"
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormLabel>Зарплата от</FormLabel>
                    <FormControl><Input type="number" {...field} /></FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="salaryTo"
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormLabel>До</FormLabel>
                    <FormControl><Input type="number" {...field} /></FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Описание</FormLabel>
                  <FormControl><Textarea rows={6} {...field} /></FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button
              type="submit"
              disabled={form.formState.isSubmitting || createMut.isPending || updateMut.isPending}
            >
              Сохранить
            </Button>
          </div>
        )}
      </FormWrapper>
      {!isNew && (
        <div className="flex gap-2 pt-4 border-t">
          <Button variant="outline" onClick={() => onAction('publish')}>Опубликовать</Button>
          <Button variant="outline" onClick={() => onAction('archive')}>Архивировать</Button>
          <Button variant="destructive" onClick={() => onAction('delete')}>Удалить</Button>
        </div>
      )}
    </div>
  )
}

export default CompanyVacanciesPage
