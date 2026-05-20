import { useState } from 'react'
import { useNavigate } from 'react-router'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { useToast } from '@/shared/hooks/use-toast'
import { usePostAdminNew } from '@/api/generated/auth/admin/admin'
import { AppError } from '@/shared/api/http/client'

export default function AdminUsersPage() {
  const [userId, setUserId] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()
  const { toast } = useToast()
  const promoteUser = usePostAdminNew()

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!userId.trim()) return

    setIsLoading(true)

    try {
      await promoteUser.mutateAsync({ data: { user_id: userId } })
      toast({ title: 'Пользователь повышен до администратора' })
      setUserId('')
    } catch (err) {
      const msg = err instanceof AppError ? err.message : 'Не удалось повысить пользователя'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-md p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Button variant="ghost" size="sm" onClick={() => navigate('/admin/skills')}>
          ← Назад
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Назначить администратором</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-4">
            <div>
              <Input
                placeholder="User ID"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                disabled={isLoading}
              />
            </div>

            <Button type="submit" disabled={isLoading || !userId.trim()}>
              {isLoading ? 'Назначение...' : 'Назначить администратором'}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Информация</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Введите ID пользователя, чтобы назначить его администратором платформы.
            Пользователь получит доступ к админ-панели.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}