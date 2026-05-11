import { Link } from 'react-router'
import { RegisterForm } from '@/features/auth'

export default function RegisterPage() {
  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 py-8">
      <div className="w-full max-w-md space-y-6 rounded-lg border bg-card p-8 shadow-sm">
        <div className="space-y-1 text-center">
          <h1 className="text-2xl font-bold">Регистрация</h1>
          <p className="text-sm text-muted-foreground">Создайте новый аккаунт</p>
        </div>
        <RegisterForm />
        <p className="text-center text-sm text-muted-foreground">
          Уже есть аккаунт?{' '}
          <Link to="/login" className="font-medium text-primary underline">
            Войти
          </Link>
        </p>
      </div>
    </div>
  )
}
