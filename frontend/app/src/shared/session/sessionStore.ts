import { create } from 'zustand'
import type { CompanyMembership } from './types'

export interface SessionUser {
  id: number | string
  role: 'Candidate' | 'Company'
  username: string
  email?: string
  firstName?: string
  lastName?: string
}

export type SessionStatus = 'unknown' | 'anon' | 'authed'

interface SessionState {
  status: SessionStatus
  user: SessionUser | undefined
  companies: CompanyMembership[]
  activeCompanyId: string | undefined
  setAuthed: (user: SessionUser) => void
  setAnon: () => void
  setUnknown: () => void
  setCompanies: (list: CompanyMembership[]) => void
  setActiveCompany: (id: string | undefined) => void
}

export const useSessionStore = create<SessionState>()((set) => ({
  status: 'unknown',
  user: undefined,
  companies: [],
  activeCompanyId: undefined,
  setAuthed: (user) => set({ status: 'authed', user }),
  setAnon: () =>
    set({
      status: 'anon',
      user: undefined,
      companies: [],
      activeCompanyId: undefined,
    }),
  setUnknown: () => set({ status: 'unknown', user: undefined }),
  setCompanies: (list) => set({ companies: list }),
  setActiveCompany: (id) => set({ activeCompanyId: id }),
}))
