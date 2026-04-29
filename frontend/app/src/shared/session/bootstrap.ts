import { useSessionStore } from './sessionStore'
import type { SessionUser } from './sessionStore'
import { AppError } from '@/shared/api/http/client'
import { fetchCompaniesMe } from '@/shared/api/companies/companiesMe'
import { readActiveCompanyId, writeActiveCompanyId } from './types'
import { postAuthMe } from '@/api/generated/auth/auth/auth'
import type { DtoUserResponse } from '@/api/generated/auth/schemas'

function toSessionUser(data: DtoUserResponse): SessionUser | null {
  if (!data.id || !data.username || !data.role) return null
  if (data.role !== 'Candidate' && data.role !== 'Company') return null
  return {
    id: data.id,
    role: data.role,
    username: data.username,
    email: data.email,
    firstName: data.first_name,
    lastName: data.last_name,
  }
}

async function fetchCurrentUser(): Promise<SessionUser | null> {
  try {
    const data = await postAuthMe()
    return toSessionUser(data)
  } catch (e) {
    if (e instanceof AppError && e.status === 401) return null
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

  const user = await fetchCurrentUser()

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
