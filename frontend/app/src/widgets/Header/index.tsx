import { Link, useNavigate } from 'react-router'
import { useSession } from '@/shared/session/useSession'
import { useSessionStore } from '@/shared/session/sessionStore'
import { postAuthLogout } from '@/api/generated/auth/auth/auth'
import { CompanySwitcher } from './CompanySwitcher'

export function Header() {
  const { isAuthed, role, user, activeCompanyId } = useSession()
  const navigate = useNavigate()
  const companyBase = activeCompanyId ? `/company/${activeCompanyId}` : '/company'

  async function handleLogout() {
    try {
      await postAuthLogout()
    } catch {
      void 0
    }
    useSessionStore.getState().setAnon()
    navigate('/login', { replace: true })
  }

  return (
    <header className="flex flex-wrap items-center gap-4 border-b border-border px-4 py-3">
      <Link to="/" className="text-base font-bold">
        trainee-match
      </Link>
      <nav className="flex flex-wrap items-center gap-3 text-sm">
        <Link to="/vacancies">Вакансии</Link>
        <Link to="/companies">Компании</Link>
        {isAuthed && role === 'Candidate' && (
          <>
            <Link to="/me/profile">Профиль</Link>
            <Link to="/me/resumes">Резюме</Link>
            <Link to="/me/applications">Отклики</Link>
          </>
        )}
        {isAuthed && role === 'Company' && activeCompanyId && (
          <>
            <Link to={`${companyBase}/dashboard`}>Дашборд</Link>
            <Link to={`${companyBase}/vacancies`}>Мои вакансии</Link>
            <Link to={`${companyBase}/applications`}>Отклики</Link>
            <Link to={`${companyBase}/members`}>Команда</Link>
            <Link to={`${companyBase}/profile`}>Профиль компании</Link>
          </>
        )}
        {isAuthed && role === 'Company' && !activeCompanyId && (
          <Link to="/company/new">Создать компанию</Link>
        )}
        {isAuthed && role === 'admin' && (
          <>
            <Link to="/admin/skills">Навыки</Link>
            <Link to="/admin/users">Пользователи</Link>
            <Link to="/candidates">Кандидаты</Link>
          </>
        )}
      </nav>
      <div className="ml-auto flex flex-wrap items-center gap-3 text-sm">
        {isAuthed && role === 'Company' && <CompanySwitcher />}
        {isAuthed ? (
          <>
            <span className="text-muted-foreground">{user?.username}</span>
            <button
              type="button"
              onClick={handleLogout}
              className="underline"
            >
              Выйти
            </button>
          </>
        ) : (
          <>
            <Link to="/login">Войти</Link>
            <Link to="/register">Регистрация</Link>
          </>
        )}
      </div>
    </header>
  )
}
