import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useQueryClient } from '@tanstack/react-query'
import { useToast } from '@/shared/hooks/use-toast'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/ui/form'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import {
  useGetCompaniesId,
  usePatchCompaniesId,
  useDeleteCompaniesId,
  getGetCompaniesIdQueryKey,
} from '@/api/generated/company/company/company'
import type { DtoCompanyResponse } from '@/api/generated/company/schemas'
import { AppError } from '@/shared/api/http/client'
import { useSession } from '@/shared/session/useSession'
import { refreshCompanies } from '@/shared/session/refreshCompanies'

const profileSchema = z.object({
  name: z.string().min(1, 'Введите название').max(200, 'Максимум 200 символов'),
  description: z
    .string()
    .max(5000, 'Максимум 5000 символов')
    .optional()
    .or(z.literal('')),
  website: z.string().url('Некорректный URL').optional().or(z.literal('')),
  logoKey: z
    .string()
    .max(2048, 'Максимум 2048 символов')
    .url('Должен быть URL')
    .optional()
    .or(z.literal('')),
})

type ProfileFormData = z.infer<typeof profileSchema>

interface Props {
  companyId: string
}

export function CompanyProfile({ companyId }: Props) {
  const { companies } = useSession()
  const isAdmin =
    companies.find((c) => c.id === companyId)?.role === 'admin'
  const query = useGetCompaniesId(companyId)
  const [editing, setEditing] = useState(false)

  if (query.isLoading) return <LoadingState />
  if (query.isError || !query.data) {
    return <ErrorState onRetry={() => query.refetch()} />
  }

  const company = query.data

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader className="flex-row items-start justify-between gap-4">
          <div>
            <CardTitle>{company.name ?? 'Компания'}</CardTitle>
            <CardDescription>
              {company.website ? (
                <a
                  href={company.website}
                  target="_blank"
                  rel="noreferrer noopener"
                  className="underline"
                >
                  {company.website}
                </a>
              ) : (
                'Сайт не указан'
              )}
            </CardDescription>
          </div>
          {isAdmin && !editing && (
            <Button onClick={() => setEditing(true)}>Редактировать</Button>
          )}
        </CardHeader>
        <CardContent>
          {editing ? (
            <CompanyProfileEditForm
              companyId={companyId}
              company={company}
              onCancel={() => setEditing(false)}
              onSaved={() => setEditing(false)}
            />
          ) : (
            <CompanyProfileView company={company} />
          )}
        </CardContent>
      </Card>

      {isAdmin && !editing && <CompanyDangerZone companyId={companyId} />}
    </div>
  )
}

function CompanyProfileView({ company }: { company: DtoCompanyResponse }) {
  return (
    <dl className="grid gap-3 sm:grid-cols-2">
      <Field label="Описание">
        {company.description ? (
          <p className="whitespace-pre-wrap">{company.description}</p>
        ) : (
          <span className="text-muted-foreground">—</span>
        )}
      </Field>
      <Field label="Логотип (URL)">
        {company.logoURL ? (
          <a
            href={company.logoURL}
            target="_blank"
            rel="noreferrer noopener"
            className="break-all underline"
          >
            {company.logoURL}
          </a>
        ) : (
          <span className="text-muted-foreground">—</span>
        )}
      </Field>
      <Field label="Открытых вакансий">
        {company.openVacanciesCount ?? 0}
      </Field>
      <Field label="Создана">
        {company.createdAt
          ? new Date(company.createdAt).toLocaleDateString()
          : '—'}
      </Field>
    </dl>
  )
}

function Field({
  label,
  children,
}: {
  label: string
  children: React.ReactNode
}) {
  return (
    <div>
      <dt className="text-sm font-medium text-muted-foreground">{label}</dt>
      <dd className="mt-1 text-sm">{children}</dd>
    </div>
  )
}

function CompanyProfileEditForm({
  companyId,
  company,
  onCancel,
  onSaved,
}: {
  companyId: string
  company: DtoCompanyResponse
  onCancel: () => void
  onSaved: () => void
}) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const patch = usePatchCompaniesId()
  const [serverError, setServerError] = useState<string | null>(null)

  const form = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: company.name ?? '',
      description: company.description ?? '',
      website: company.website ?? '',
      logoKey: company.logoURL ?? '',
    },
  })

  useEffect(() => {
    form.reset({
      name: company.name ?? '',
      description: company.description ?? '',
      website: company.website ?? '',
      logoKey: company.logoURL ?? '',
    })
  }, [company, form])

  async function onSubmit(values: ProfileFormData) {
    setServerError(null)
    try {
      await patch.mutateAsync({
        id: companyId,
        data: {
          name: values.name,
          description: values.description ? values.description : undefined,
          website: values.website ? values.website : undefined,
        },
      })
      await Promise.all([
        qc.invalidateQueries({ queryKey: getGetCompaniesIdQueryKey(companyId) }),
        refreshCompanies(),
      ])
      toast({ title: 'Профиль обновлён' })
      onSaved()
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось сохранить профиль'
          : 'Не удалось сохранить профиль'
      setServerError(msg)
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        noValidate
        className="space-y-4"
      >
        {serverError && (
          <p role="alert" className="text-sm text-destructive">
            {serverError}
          </p>
        )}
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Название</FormLabel>
              <FormControl>
                <Input maxLength={200} {...field} />
              </FormControl>
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
              <FormControl>
                <Input
                  type="url"
                  placeholder="https://example.com"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="logoKey"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Логотип (URL)</FormLabel>
              <FormControl>
                <Input
                  type="url"
                  placeholder="https://cdn.example.com/logo.png"
                  {...field}
                />
              </FormControl>
              <FormMessage />
              <p className="text-xs text-muted-foreground">
                Загрузка файла появится позже. Пока поле сохраняется только
                локально.
              </p>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Описание</FormLabel>
              <FormControl>
                <Textarea rows={5} maxLength={5000} {...field} />
              </FormControl>
              <FormMessage />
              <p className="text-xs text-muted-foreground">
                {(field.value ?? '').length} / 5000
              </p>
            </FormItem>
          )}
        />
        <div className="flex justify-end gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={patch.isPending}
          >
            Отмена
          </Button>
          <Button type="submit" disabled={patch.isPending}>
            {patch.isPending ? 'Сохранение…' : 'Сохранить'}
          </Button>
        </div>
      </form>
    </Form>
  )
}

function CompanyDangerZone({ companyId }: { companyId: string }) {
  const navigate = useNavigate()
  const { toast } = useToast()
  const del = useDeleteCompaniesId()
  const [open, setOpen] = useState(false)

  async function onConfirm() {
    try {
      await del.mutateAsync({ id: companyId })
      await refreshCompanies({ setActiveId: undefined })
      toast({ title: 'Компания удалена' })
      setOpen(false)
      navigate('/company')
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось удалить компанию'
          : 'Не удалось удалить компанию'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <Card className="border-destructive/50">
      <CardHeader>
        <CardTitle className="text-lg text-destructive">Опасная зона</CardTitle>
        <CardDescription>
          Удаление компании необратимо. Все вакансии и связи будут удалены.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button variant="destructive" onClick={() => setOpen(true)}>
          Удалить компанию
        </Button>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Удалить компанию?</DialogTitle>
              <DialogDescription>
                Это действие нельзя отменить.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
                disabled={del.isPending}
              >
                Отмена
              </Button>
              <Button
                type="button"
                variant="destructive"
                onClick={onConfirm}
                disabled={del.isPending}
              >
                {del.isPending ? 'Удаление…' : 'Удалить'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  )
}
