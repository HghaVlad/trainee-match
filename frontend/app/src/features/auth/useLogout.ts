import { useNavigate } from 'react-router'
import { useSessionStore } from '@/shared/session/sessionStore'
import { usePostAuthLogout } from '@/api/generated/auth/auth/auth'

export function useLogout() {
  const navigate = useNavigate()
  const setAnon = useSessionStore((s) => s.setAnon)
  const logoutMutation = usePostAuthLogout()

  return async function logout() {
    try {
      await logoutMutation.mutateAsync()
    } catch {
      // ignore network/server errors — always clear local session
    }
    setAnon()
    navigate('/')
  }
}
