import { redirect } from 'react-router'
import { useSessionStore } from './sessionStore'
import type { SessionUser } from './sessionStore'
import { AppError } from '@/shared/api/http/client'
import { readActiveCompanyId, writeActiveCompanyId } from './types'
import { refreshCompanies } from './refreshCompanies'
import { getAuthMe } from '@/api/generated/auth/auth/auth'
import type { DtoUserResponse } from '@/api/generated/auth/schemas'

function toSessionUser(data: DtoUserResponse): SessionUser | null {
  if (!data.id || !data.username || !data.role) return null
  if (data.role !== 'Candidate' && data.role !== 'Company' && data.role !== 'PlatformAdmin') return null
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
    const data = await getAuthMe()
    return toSessionUser(data)
  } catch (e) {
    if (e instanceof AppError && e.status === 401) return null
    return null
  }
}

async function loadCompaniesForUser(user: SessionUser): Promise<void> {
  const { setActiveCompany } = useSessionStore.getState()
  try {
    await refreshCompanies()
    const { companies } = useSessionStore.getState()
    const stored = readActiveCompanyId(user.id)
    if (companies.length === 0) {
      if (stored) {
        setActiveCompany(stored)
      } else {
        setActiveCompany(undefined)
      }
      return
    }
    const restored =
      stored && companies.some((c) => c.id === stored) ? stored : companies[0]!.id
    setActiveCompany(restored)
    writeActiveCompanyId(user.id, restored)
  } catch {
    useSessionStore.getState().setCompanies([])
    useSessionStore.getState().setActiveCompany(undefined)
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

export function sessionRedirect(request: Request): Response | null {
  const url = new URL(request.url)
  const next = encodeURIComponent(url.pathname + url.search)
  return redirect(`/login?next=${next}`)
}

if (typeof window !== 'undefined') {
  window.addEventListener('session:expired', () => {
    useSessionStore.getState().setAnon()
  })
}
