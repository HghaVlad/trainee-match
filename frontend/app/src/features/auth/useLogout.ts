import { useNavigate } from 'react-router'
import { useSessionStore } from '@/shared/session/sessionStore'
import { usePostAuthLogout } from '@/api/generated/auth/auth/auth'
import { clearCompanyRoles } from '@/shared/session/types'

export function useLogout() {
  const navigate = useNavigate()
  const { user, setAnon } = useSessionStore((s) => ({
    user: s.user,
    setAnon: s.setAnon,
  }))
  const logoutMutation = usePostAuthLogout()

  return async function logout() {
    try {
      await logoutMutation.mutateAsync()
    } catch {
      // ignore network/server errors — always clear local session
    }
    // clear persisted roles on explicit logout (not on session expiry)
    if (user) clearCompanyRoles(user.id)
    setAnon()
    navigate('/')
  }
}
