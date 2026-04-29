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
    <header
      style={{
        padding: '0.75rem 1rem',
        borderBottom: '1px solid #eee',
        display: 'flex',
        gap: '1rem',
        alignItems: 'center',
      }}
    >
      <Link to="/" style={{ fontWeight: 'bold' }}>
        trainee-match
      </Link>
      <nav style={{ display: 'flex', gap: '0.75rem', flex: 1 }}>
        <Link to="/vacancies">Vacancies</Link>
        <Link to="/companies">Companies</Link>
        {isAuthed && role === 'Candidate' && (
          <>
            <Link to="/me/profile">Profile</Link>
            <Link to="/me/resumes">Resumes</Link>
            <Link to="/me/applications">Applications</Link>
          </>
        )}
        {isAuthed && role === 'Company' && activeCompanyId && (
          <>
            <Link to={`${companyBase}/dashboard`}>Dashboard</Link>
            <Link to={`${companyBase}/vacancies`}>Vacancies</Link>
            <Link to={`${companyBase}/applications`}>Applications</Link>
            <Link to={`${companyBase}/members`}>Members</Link>
            <Link to={`${companyBase}/profile`}>Profile</Link>
          </>
        )}
        {isAuthed && role === 'Company' && !activeCompanyId && (
          <Link to="/company/new">Создать компанию</Link>
        )}
      </nav>
      {isAuthed && role === 'Company' && <CompanySwitcher />}
      {isAuthed ? (
        <>
          <span>{user?.username}</span>
          <button type="button" onClick={handleLogout}>
            Logout
          </button>
        </>
      ) : (
        <>
          <Link to="/login">Login</Link>
          <Link to="/register">Register</Link>
        </>
      )}
    </header>
  )
}
