import { useSessionStore, type SessionStatus, type SessionUser } from './sessionStore'

export function useSession() {
  return useSessionStore((s) => ({
    status: s.status,
    user: s.user,
    isAuthed: s.status === 'authed',
    isAnon: s.status === 'anon',
    isLoading: s.status === 'unknown',
    role: s.user?.role,
  }))
}

export function useSessionStatus(): SessionStatus {
  return useSessionStore((s) => s.status)
}

export function useSessionUser(): SessionUser | undefined {
  return useSessionStore((s) => s.user)
}
