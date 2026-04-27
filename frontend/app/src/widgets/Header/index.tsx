import { Link } from 'react-router'
import { useSession } from '@/shared/session/useSession'
import { useSessionStore } from '@/shared/session/sessionStore'

export function Header() {
  const { isAuthed, role, user } = useSession()

  function handleLogout() {
    useSessionStore.getState().setAnon()
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
            <Link to="/me/profile">My Profile</Link>
            <Link to="/me/resumes">My Resumes</Link>
          </>
        )}
        {isAuthed && role === 'Company' && (
          <>
            <Link to="/company/me">Company</Link>
            <Link to="/company/vacancies">My Vacancies</Link>
            <Link to="/company/members">Members</Link>
          </>
        )}
      </nav>
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
