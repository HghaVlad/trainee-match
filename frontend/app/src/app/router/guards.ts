import { redirect, type LoaderFunctionArgs } from 'react-router'
import { useSessionStore } from '@/shared/session/sessionStore'

export function requireAuth({ request }: LoaderFunctionArgs) {
  const { status } = useSessionStore.getState()
  if (status !== 'authed') {
    const url = new URL(request.url)
    return redirect(`/login?next=${encodeURIComponent(url.pathname)}`)
  }
  return null
}

export function requireRole(role: 'Candidate' | 'Company') {
  return function roleLoader({ request }: LoaderFunctionArgs) {
    const { status, user } = useSessionStore.getState()
    if (status !== 'authed') {
      const url = new URL(request.url)
      return redirect(`/login?next=${encodeURIComponent(url.pathname)}`)
    }
    if (user?.role !== role) {
      return redirect('/403')
    }
    return null
  }
}

export function redirectIfAuth() {
  const { status, user } = useSessionStore.getState()
  if (status !== 'authed' || !user) return null
  return redirect(user.role === 'Candidate' ? '/me/profile' : '/company/me')
}
