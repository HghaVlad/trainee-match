import { useShallow } from 'zustand/react/shallow'
import { useSessionStore, type SessionStatus, type SessionUser } from './sessionStore'
import type { CompanyMembership } from './types'

export function useSession() {
  return useSessionStore(
    useShallow((s) => ({
      status: s.status,
      user: s.user,
      isAuthed: s.status === 'authed',
      isAnon: s.status === 'anon',
      isLoading: s.status === 'unknown',
      role: s.user?.role,
      companies: s.companies,
      activeCompanyId: s.activeCompanyId,
      activeCompany: s.companies.find((c) => c.id === s.activeCompanyId),
    })),
  )
}

export function useSessionStatus(): SessionStatus {
  return useSessionStore((s) => s.status)
}

export function useSessionUser(): SessionUser | undefined {
  return useSessionStore((s) => s.user)
}

export function useCompanies(): CompanyMembership[] {
  return useSessionStore((s) => s.companies)
}

export function useActiveCompanyId(): string | undefined {
  return useSessionStore((s) => s.activeCompanyId)
}
