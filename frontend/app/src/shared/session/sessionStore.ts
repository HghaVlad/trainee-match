import { create } from 'zustand'

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
  setAuthed: (user: SessionUser) => void
  setAnon: () => void
  setUnknown: () => void
}

export const useSessionStore = create<SessionState>()((set) => ({
  status: 'unknown',
  user: undefined,
  setAuthed: (user) => set({ status: 'authed', user }),
  setAnon: () => set({ status: 'anon', user: undefined }),
  setUnknown: () => set({ status: 'unknown', user: undefined }),
}))
