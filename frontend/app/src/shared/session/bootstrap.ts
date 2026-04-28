import { httpClient } from '@/shared/api/http/client'
import { useSessionStore } from './sessionStore'
import type { SessionUser } from './sessionStore'
import { env } from '@/shared/config/env'
import { AppError } from '@/shared/api/http/client'
import { fetchCompaniesMe } from '@/shared/api/companies/companiesMe'
import { readActiveCompanyId, writeActiveCompanyId } from './types'

interface AuthMeResponse {
  id: number
  role: 'Candidate' | 'Company'
  username: string
  email: string
  firstName: string
  lastName: string
}

interface CandidateMeResponse {
  id: number
  username: string
  email?: string
  firstName?: string
  lastName?: string
}

async function bootstrapViaAuthMe(): Promise<SessionUser | null> {
  try {
    const data = await httpClient
      .get<AuthMeResponse>('/auth/me')
      .then((r) => r.data)
    return {
      id: data.id,
      role: data.role,
      username: data.username,
      email: data.email,
      firstName: data.firstName,
      lastName: data.lastName,
    }
  } catch {
    return null
  }
}

async function bootstrapViaProbe(): Promise<SessionUser | null> {
  try {
    const data = await httpClient
      .get<CandidateMeResponse>('/api/v1/candidate/me')
      .then((r) => r.data)
    return {
      id: data.id,
      role: 'Candidate',
      username: data.username,
      email: data.email,
      firstName: data.firstName,
      lastName: data.lastName,
    }
  } catch (e) {
    if (e instanceof AppError && e.status === 401) {
      return null
    }
    return null
  }
}

async function loadCompaniesForUser(user: SessionUser): Promise<void> {
  const { setCompanies, setActiveCompany } = useSessionStore.getState()
  try {
    const { data } = await fetchCompaniesMe({ limit: 100 })
    setCompanies(data)
    if (data.length === 0) {
      setActiveCompany(undefined)
      writeActiveCompanyId(user.id, undefined)
      return
    }
    const stored = readActiveCompanyId(user.id)
    const restored = stored && data.some((c) => c.id === stored) ? stored : data[0]!.id
    setActiveCompany(restored)
    writeActiveCompanyId(user.id, restored)
  } catch {
    setCompanies([])
    setActiveCompany(undefined)
  }
}

export async function bootstrap(): Promise<void> {
  const { setAuthed, setAnon } = useSessionStore.getState()

  const user = env.VITE_AUTH_ME_AVAILABLE
    ? await bootstrapViaAuthMe()
    : await bootstrapViaProbe()

  if (!user) {
    setAnon()
    return
  }

  setAuthed(user)
  if (user.role === 'Company') {
    await loadCompaniesForUser(user)
  }
}

if (typeof window !== 'undefined') {
  window.addEventListener('session:expired', () => {
    useSessionStore.getState().setAnon()
  })
}
