import { Link } from 'react-router'
import { LoginForm } from '@/features/auth'

export default function LoginPage() {
  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4">
      <div className="w-full max-w-md space-y-6 rounded-lg border bg-card p-8 shadow-sm">
        <div className="space-y-1 text-center">
          <h1 className="text-2xl font-bold">Вход</h1>
          <p className="text-sm text-muted-foreground">
            Войдите в свой аккаунт
          </p>
        </div>
        <LoginForm />
        <p className="text-center text-sm text-muted-foreground">
          Нет аккаунта?{' '}
          <Link to="/register" className="font-medium text-primary underline">
            Зарегистрироваться
          </Link>
        </p>
      </div>
    </div>
  )
}
