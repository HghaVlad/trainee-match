import { useState } from 'react'
import { useNavigate } from 'react-router'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
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
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/ui/form'
import { usePostCompanies } from '@/api/generated/company/company/company'
import { AppError } from '@/shared/api/http/client'
import { addLocalCompany, refreshCompanies } from '@/shared/session/refreshCompanies'

const createSchema = z.object({
  name: z.string().min(1, 'Введите название').max(200, 'Максимум 200 символов'),
  description: z.string().max(5000, 'Максимум 5000 символов').optional().or(z.literal('')),
  website: z.string().url('Некорректный URL').optional().or(z.literal('')),
})

type CreateCompanyFormData = z.infer<typeof createSchema>

export function CompanyCreateForm() {
  const navigate = useNavigate()
  const { toast } = useToast()
  const [serverError, setServerError] = useState<string | null>(null)
  const create = usePostCompanies()

  const form = useForm<CreateCompanyFormData>({
    resolver: zodResolver(createSchema),
    defaultValues: { name: '', description: '', website: '' },
  })

  async function onSubmit(values: CreateCompanyFormData) {
    setServerError(null)
    try {
      const result = await create.mutateAsync({
        data: {
          name: values.name,
          description: values.description ? values.description : undefined,
          website: values.website ? values.website : undefined,
        },
      })
      const newId = result?.id
      if (!newId) {
        throw new AppError('NO_ID', 'Сервер не вернул id компании', 500)
      }
      addLocalCompany(
        {
          id: newId,
          name: values.name,
          openVacanciesCount: 0,
          createdAt: new Date().toISOString(),
          role: 'admin',
        },
        true,
      )
      void refreshCompanies({ setActiveId: newId }).catch(() => undefined)
      toast({ title: 'Компания создана' })
      navigate(`/company/${newId}/dashboard`)
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось создать компанию'
          : 'Не удалось создать компанию'
      setServerError(msg)
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Создание компании</CardTitle>
        <CardDescription>
          Заполните основные данные. Логотип и расширенный профиль можно будет
          добавить позже.
        </CardDescription>
      </CardHeader>
      <CardContent>
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
                    <Input maxLength={200} autoComplete="organization" {...field} />
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
                  <FormLabel>Сайт (необязательно)</FormLabel>
                  <FormControl>
                    <Input
                      type="url"
                      placeholder="https://example.com"
                      autoComplete="url"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Описание (необязательно)</FormLabel>
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
                onClick={() => navigate('/')}
                disabled={create.isPending}
              >
                Отмена
              </Button>
              <Button type="submit" disabled={create.isPending}>
                {create.isPending ? 'Создание…' : 'Создать'}
              </Button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  )
}
