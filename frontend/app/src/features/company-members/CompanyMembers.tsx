import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useToast } from '@/shared/hooks/use-toast'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
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
import { CursorPagination } from '@/shared/ui/CursorPagination'
import { DataTable } from '@/shared/ui/DataTable'
import {
  useCompanyMembersQuery,
  getCompanyMembersQueryKey,
  type CompanyMember,
} from '@/shared/api/companies'
import {
  usePostCompaniesIdMembers,
  usePatchCompaniesIdMembersUserId,
  useDeleteCompaniesIdMembersUserId,
} from '@/api/generated/company/member/member'
import {
  DtoCompanyAddHrRequestRole,
  DtoCompanyUpdateMemberRequestRole,
  type DtoCompanyUpdateMemberRequestRole as MemberRoleType,
} from '@/api/generated/company/schemas'
import { AppError } from '@/shared/api/http/client'
import { useSession } from '@/shared/session/useSession'
import { refreshCompanies } from '@/shared/session/refreshCompanies'
import type { ColumnDef } from '@tanstack/react-table'

const PAGE_SIZE = 20

const ROLE_LABEL: Record<MemberRoleType, string> = {
  admin: 'Администратор',
  recruiter: 'Рекрутер',
}

interface Props {
  companyId: string
}

export function CompanyMembers({ companyId }: Props) {
  const { user, companies } = useSession()
  const isAdmin =
    companies.find((c) => c.id === companyId)?.role === 'admin'
  const currentUserId = user ? String(user.id) : undefined

  const [cursor, setCursor] = useState<string | undefined>(undefined)
  const params = { cursor, limit: PAGE_SIZE }
  const query = useCompanyMembersQuery(companyId, params)

  if (query.isLoading) return <LoadingState />
  if (query.isError || !query.data) {
    return <ErrorState onRetry={() => query.refetch()} />
  }

  const members = query.data.data

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader className="flex-row items-start justify-between gap-4">
          <div>
            <CardTitle>Команда</CardTitle>
            <CardDescription>
              Управление участниками компании.
            </CardDescription>
          </div>
          {isAdmin && <AddMemberButton companyId={companyId} />}
        </CardHeader>
        <CardContent>
          {members.length === 0 ? (
            <EmptyState
              title="Нет участников"
              description={
                isAdmin
                  ? 'Добавьте первого члена команды, чтобы начать совместную работу.'
                  : 'Список участников пуст.'
              }
            />
          ) : (
            <MembersTable
              companyId={companyId}
              members={members}
              isAdmin={isAdmin}
              currentUserId={currentUserId}
            />
          )}
          <div className="mt-4">
            <CursorPagination
              nextCursor={query.data.nextCursor}
              onNext={() =>
                setCursor(query.data?.nextCursor ?? undefined)
              }
              isLoading={query.isFetching}
            />
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

function MembersTable({
  companyId,
  members,
  isAdmin,
  currentUserId,
}: {
  companyId: string
  members: CompanyMember[]
  isAdmin: boolean
  currentUserId: string | undefined
}) {
  const columns: ColumnDef<CompanyMember>[] = [
    {
      header: 'Пользователь',
      accessorKey: 'username',
      cell: ({ row }) => row.original.username || row.original.userId,
    },
    {
      header: 'Роль',
      accessorKey: 'role',
      cell: ({ row }) =>
        isAdmin ? (
          <RoleSelectCell companyId={companyId} member={row.original} />
        ) : (
          ROLE_LABEL[row.original.role]
        ),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) => {
        if (!isAdmin) return null
        if (currentUserId && row.original.userId === currentUserId) return null
        return (
          <RemoveMemberButton companyId={companyId} member={row.original} />
        )
      },
    },
  ]

  return <DataTable columns={columns} data={members} />
}

function RoleSelectCell({
  companyId,
  member,
}: {
  companyId: string
  member: CompanyMember
}) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const patch = usePatchCompaniesIdMembersUserId()

  async function onChange(next: string) {
    const role = next as MemberRoleType
    if (role === member.role) return
    try {
      await patch.mutateAsync({
        id: companyId,
        userId: member.userId,
        data: { role },
      })
      await Promise.all([
        qc.invalidateQueries({
          queryKey: ['company-members', companyId],
        }),
        refreshCompanies(),
      ])
      toast({ title: 'Роль обновлена' })
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось обновить роль'
          : 'Не удалось обновить роль'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <Select
      value={member.role}
      onValueChange={onChange}
      disabled={patch.isPending}
    >
      <SelectTrigger className="w-44">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value={DtoCompanyUpdateMemberRequestRole.admin}>
          {ROLE_LABEL.admin}
        </SelectItem>
        <SelectItem value={DtoCompanyUpdateMemberRequestRole.recruiter}>
          {ROLE_LABEL.recruiter}
        </SelectItem>
      </SelectContent>
    </Select>
  )
}

function RemoveMemberButton({
  companyId,
  member,
}: {
  companyId: string
  member: CompanyMember
}) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const del = useDeleteCompaniesIdMembersUserId()
  const [open, setOpen] = useState(false)

  async function onConfirm() {
    try {
      await del.mutateAsync({ id: companyId, userId: member.userId })
      await Promise.all([
        qc.invalidateQueries({
          queryKey: ['company-members', companyId],
        }),
        refreshCompanies(),
      ])
      toast({ title: 'Участник удалён' })
      setOpen(false)
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось удалить участника'
          : 'Не удалось удалить участника'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <>
      <Button
        size="sm"
        variant="ghost"
        className="text-destructive"
        onClick={() => setOpen(true)}
      >
        Удалить
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Удалить участника?</DialogTitle>
            <DialogDescription>
              {member.username || member.userId} больше не сможет
              управлять компанией.
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
    </>
  )
}

const addSchema = z.object({
  userID: z.string().min(1, 'Введите идентификатор пользователя').max(128),
  role: z.enum([
    DtoCompanyAddHrRequestRole.recruiter,
    DtoCompanyAddHrRequestRole.admin,
  ]),
})

type AddFormData = z.infer<typeof addSchema>

function AddMemberButton({ companyId }: { companyId: string }) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const add = usePostCompaniesIdMembers()
  const [open, setOpen] = useState(false)
  const [serverError, setServerError] = useState<string | null>(null)

  const form = useForm<AddFormData>({
    resolver: zodResolver(addSchema),
    defaultValues: {
      userID: '',
      role: DtoCompanyAddHrRequestRole.recruiter,
    },
  })

  function handleOpenChange(next: boolean) {
    if (next) {
      form.reset({
        userID: '',
        role: DtoCompanyAddHrRequestRole.recruiter,
      })
      setServerError(null)
    }
    setOpen(next)
  }

  async function onSubmit(values: AddFormData) {
    setServerError(null)
    try {
      await add.mutateAsync({
        id: companyId,
        data: { userID: values.userID, role: values.role },
      })
      await Promise.all([
        qc.invalidateQueries({
          queryKey: getCompanyMembersQueryKey(companyId),
        }),
        qc.invalidateQueries({
          queryKey: ['company-members', companyId],
        }),
        refreshCompanies(),
      ])
      toast({ title: 'Участник добавлен' })
      setOpen(false)
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось добавить участника'
          : 'Не удалось добавить участника'
      setServerError(msg)
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <>
      <Button onClick={() => handleOpenChange(true)}>Добавить</Button>
      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Добавить участника</DialogTitle>
            <DialogDescription>
              Укажите идентификатор пользователя и роль в компании.
            </DialogDescription>
          </DialogHeader>
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
                name="userID"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>UserID</FormLabel>
                    <FormControl>
                      <Input autoFocus {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="role"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Роль</FormLabel>
                    <FormControl>
                      <Select
                        value={field.value}
                        onValueChange={field.onChange}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem
                            value={DtoCompanyAddHrRequestRole.recruiter}
                          >
                            {ROLE_LABEL.recruiter}
                          </SelectItem>
                          <SelectItem
                            value={DtoCompanyAddHrRequestRole.admin}
                          >
                            {ROLE_LABEL.admin}
                          </SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => handleOpenChange(false)}
                  disabled={add.isPending}
                >
                  Отмена
                </Button>
                <Button type="submit" disabled={add.isPending}>
                  {add.isPending ? 'Добавление…' : 'Добавить'}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </>
  )
}
