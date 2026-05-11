import { useState } from 'react'
import { z } from 'zod'
import { useNavigate, useSearchParams } from 'react-router'
import { FormWrapper } from '@/shared/ui/Form'
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/shared/ui/form'
import { Input } from '@/shared/ui/input'
import { Button } from '@/shared/ui/button'
import { usePostAuthLogin } from '@/api/generated/auth/auth/auth'
import { bootstrap } from '@/shared/session/bootstrap'
import { useSessionStore } from '@/shared/session/sessionStore'
import { AppError } from '@/shared/api/http/client'

const loginSchema = z.object({
  username: z.string().min(1, 'Введите имя пользователя'),
  password: z.string().min(1, 'Введите пароль'),
})

type LoginFormData = z.infer<typeof loginSchema>

export function LoginForm() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [error, setError] = useState<string | null>(null)
  const loginMutation = usePostAuthLogin()

  async function handleLogin(data: LoginFormData) {
    setError(null)
    try {
      await loginMutation.mutateAsync({ data })
      await bootstrap()
      const next = searchParams.get('next')
      if (next) {
        navigate(next)
        return
      }
      const role = useSessionStore.getState().user?.role
      navigate(role === 'Company' ? '/company' : '/me/profile')
    } catch (e) {
      if (e instanceof AppError && e.status === 401) {
        setError('Неверные данные')
      } else if (e instanceof AppError) {
        setError(e.message || 'Ошибка входа. Попробуйте позже.')
      } else {
        setError('Ошибка входа. Попробуйте позже.')
      }
    }
  }

  return (
    <FormWrapper
      schema={loginSchema}
      defaultValues={{ username: '', password: '' }}
      onSubmit={handleLogin}
    >
      {(form) => (
        <div className="space-y-4">
          {error && (
            <p role="alert" className="text-sm text-destructive">
              {error}
            </p>
          )}
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Имя пользователя</FormLabel>
                <FormControl>
                  <Input autoComplete="username" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Пароль</FormLabel>
                <FormControl>
                  <Input
                    type="password"
                    autoComplete="current-password"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button
            type="submit"
            className="w-full"
            disabled={form.formState.isSubmitting || loginMutation.isPending}
          >
            {form.formState.isSubmitting || loginMutation.isPending
              ? 'Вход...'
              : 'Войти'}
          </Button>
        </div>
      )}
    </FormWrapper>
  )
}
