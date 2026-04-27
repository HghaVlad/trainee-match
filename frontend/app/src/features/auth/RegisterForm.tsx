import { useState } from 'react'
import { z } from 'zod'
import { useNavigate } from 'react-router'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { usePostAuthRegister } from '@/api/generated/auth/auth/auth'
import { DtoRegisterUserRequestRole } from '@/api/generated/auth/schemas'
import { AppError } from '@/shared/api/http/client'

const registerSchema = z.object({
  username: z.string().min(3, 'Минимум 3 символа').max(50),
  password: z.string().min(8, 'Минимум 8 символов'),
  email: z.string().email('Некорректный email').max(254),
  first_name: z.string().min(2, 'Минимум 2 символа').max(50),
  last_name: z.string().min(2, 'Минимум 2 символа').max(50),
  role: z.enum([
    DtoRegisterUserRequestRole.Candidate,
    DtoRegisterUserRequestRole.Company,
  ]),
})

type RegisterFormData = z.infer<typeof registerSchema>

export function RegisterForm() {
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)
  const registerMutation = usePostAuthRegister()

  async function handleRegister(data: RegisterFormData) {
    setError(null)
    try {
      await registerMutation.mutateAsync({ data })
      navigate('/login')
    } catch (e) {
      if (e instanceof AppError) {
        setError(e.message || 'Ошибка регистрации. Попробуйте позже.')
      } else {
        setError('Ошибка регистрации. Попробуйте позже.')
      }
    }
  }

  return (
    <FormWrapper
      schema={registerSchema}
      defaultValues={{
        username: '',
        password: '',
        email: '',
        first_name: '',
        last_name: '',
        role: DtoRegisterUserRequestRole.Candidate,
      }}
      onSubmit={handleRegister}
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
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input type="email" autoComplete="email" {...field} />
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
                    autoComplete="new-password"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="first_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Имя</FormLabel>
                <FormControl>
                  <Input autoComplete="given-name" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="last_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Фамилия</FormLabel>
                <FormControl>
                  <Input autoComplete="family-name" {...field} />
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
                <Select
                  value={field.value}
                  onValueChange={field.onChange}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Выберите роль" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value={DtoRegisterUserRequestRole.Candidate}>
                      Кандидат
                    </SelectItem>
                    <SelectItem value={DtoRegisterUserRequestRole.Company}>
                      Компания
                    </SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button
            type="submit"
            className="w-full"
            disabled={form.formState.isSubmitting || registerMutation.isPending}
          >
            {form.formState.isSubmitting || registerMutation.isPending
              ? 'Регистрация...'
              : 'Зарегистрироваться'}
          </Button>
        </div>
      )}
    </FormWrapper>
  )
}
