import { useState } from 'react'
import { z } from 'zod'
import { useQueryClient } from '@tanstack/react-query'
import {
  useGetCompanies,
  usePostCompanies,
  usePatchCompaniesId,
  getGetCompaniesQueryKey,
} from '@/api/generated/company/company/company'
import {
  usePostCompaniesIdMembers,
  useDeleteCompaniesIdMembersUserId,
  usePatchCompaniesIdMembersUserId,
} from '@/api/generated/company/member/member'
import {
  DtoCompanyAddHrRequestRole,
  DtoCompanyUpdateMemberRequestRole,
} from '@/api/generated/company/schemas'
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

const companySchema = z.object({
  name: z.string().min(2, 'Минимум 2 символа').max(100),
  description: z.string().max(2000).optional().or(z.literal('')),
  website: z.string().url('Некорректный URL').optional().or(z.literal('')),
})
type CompanyFormData = z.infer<typeof companySchema>

const memberSchema = z.object({
  userID: z.string().min(1, 'Введите userID'),
  role: z.enum([
    DtoCompanyAddHrRequestRole.recruiter,
    DtoCompanyAddHrRequestRole.admin,
  ]),
})
type MemberFormData = z.infer<typeof memberSchema>

export default function CompanyMePage() {
  const qc = useQueryClient()
  const { data, isLoading, error, refetch } = useGetCompanies({ limit: 1 })
  const [editing, setEditing] = useState(false)
  const [error2, setError] = useState<string | null>(null)

  const createMut = usePostCompanies()
  const updateMut = usePatchCompaniesId()

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />

  const company = data?.companies?.[0]

  async function onCreate(values: CompanyFormData) {
    setError(null)
    try {
      await createMut.mutateAsync({
        data: {
          name: values.name,
          description: values.description || undefined,
          website: values.website || undefined,
        },
      })
      await qc.invalidateQueries({ queryKey: getGetCompaniesQueryKey() })
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Ошибка создания')
    }
  }

  async function onUpdate(values: CompanyFormData) {
    if (!company?.id) return
    setError(null)
    try {
      await updateMut.mutateAsync({
        id: company.id,
        data: {
          name: values.name,
          description: values.description || undefined,
          website: values.website || undefined,
        },
      })
      await qc.invalidateQueries({ queryKey: getGetCompaniesQueryKey() })
      setEditing(false)
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Ошибка сохранения')
    }
  }

  if (!company) {
    return (
      <div className="mx-auto max-w-xl p-6 space-y-4">
        <EmptyState title="Создайте профиль компании" />
        {error2 && <p role="alert" className="text-sm text-destructive">{error2}</p>}
        <FormWrapper
          schema={companySchema}
          defaultValues={{ name: '', description: '', website: '' }}
          onSubmit={onCreate}
        >
          {(form) => <CompanyFields form={form} submitting={createMut.isPending} />}
        </FormWrapper>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-xl p-6 space-y-6">
      <div className="space-y-2 rounded-lg border bg-card p-4">
        {editing ? (
          <>
            {error2 && <p role="alert" className="text-sm text-destructive">{error2}</p>}
            <FormWrapper
              schema={companySchema}
              defaultValues={{
                name: company.name ?? '',
                description: '',
                website: '',
              }}
              onSubmit={onUpdate}
            >
              {(form) => <CompanyFields form={form} submitting={updateMut.isPending} />}
            </FormWrapper>
            <Button variant="ghost" onClick={() => setEditing(false)}>Отмена</Button>
          </>
        ) : (
          <>
            <h1 className="text-2xl font-bold">{company.name}</h1>
            <p className="text-sm text-muted-foreground">
              Открытых вакансий: {company.openVacanciesCount ?? 0}
            </p>
            <Button onClick={() => setEditing(true)}>Редактировать</Button>
          </>
        )}
      </div>
      {company.id && <MembersBlock companyId={company.id} />}
    </div>
  )
}

function CompanyFields({ form, submitting }: { form: ReturnType<typeof import('react-hook-form').useForm<CompanyFormData>>; submitting: boolean }) {
  return (
    <div className="space-y-4">
      <FormField
        control={form.control}
        name="name"
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
        name="website"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Сайт</FormLabel>
            <FormControl><Input type="url" placeholder="https://..." {...field} /></FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name="description"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Описание</FormLabel>
            <FormControl><Textarea rows={4} {...field} /></FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <Button type="submit" disabled={submitting}>
        {submitting ? 'Сохранение...' : 'Сохранить'}
      </Button>
    </div>
  )
}

function MembersBlock({ companyId }: { companyId: string }) {
  const qc = useQueryClient()
  const [error, setError] = useState<string | null>(null)
  const addMut = usePostCompaniesIdMembers()
  const delMut = useDeleteCompaniesIdMembersUserId()
  const patchMut = usePatchCompaniesIdMembersUserId()

  async function onAdd(values: MemberFormData) {
    setError(null)
    try {
      await addMut.mutateAsync({
        id: companyId,
        data: { userID: values.userID, role: values.role },
      })
      await qc.invalidateQueries({ queryKey: getGetCompaniesQueryKey() })
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Не удалось добавить')
    }
  }

  async function changeRole(userId: string, role: DtoCompanyUpdateMemberRequestRole) {
    try {
      await patchMut.mutateAsync({ id: companyId, userId, data: { role } })
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Не удалось обновить роль')
    }
  }

  async function remove(userId: string) {
    if (!window.confirm('Удалить участника?')) return
    try {
      await delMut.mutateAsync({ id: companyId, userId })
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Не удалось удалить')
    }
  }

  void changeRole
  void remove

  return (
    <div className="space-y-3 rounded-lg border bg-card p-4">
      <h2 className="text-xl font-semibold">Команда</h2>
      {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
      <FormWrapper
        schema={memberSchema}
        defaultValues={{ userID: '', role: DtoCompanyAddHrRequestRole.recruiter }}
        onSubmit={onAdd}
      >
        {(form) => (
          <div className="flex items-end gap-2">
            <FormField
              control={form.control}
              name="userID"
              render={({ field }) => (
                <FormItem className="flex-1">
                  <FormLabel>UserID</FormLabel>
                  <FormControl><Input {...field} /></FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="role"
              render={({ field }) => (
                <FormItem className="w-40">
                  <FormLabel>Роль</FormLabel>
                  <Select value={field.value} onValueChange={field.onChange}>
                    <FormControl>
                      <SelectTrigger><SelectValue /></SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value={DtoCompanyAddHrRequestRole.recruiter}>recruiter</SelectItem>
                      <SelectItem value={DtoCompanyAddHrRequestRole.admin}>admin</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" disabled={addMut.isPending}>Добавить</Button>
          </div>
        )}
      </FormWrapper>
      <p className="text-xs text-muted-foreground">
        Список участников возвращается отдельным эндпоинтом, который не описан в swagger; добавление/удаление работает через POST/DELETE.
      </p>
    </div>
  )
}
