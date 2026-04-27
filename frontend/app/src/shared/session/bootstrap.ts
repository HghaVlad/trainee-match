import { httpClient } from '@/shared/api/http/client'
import { eventBus } from '@/shared/api/http/eventBus'
import { useSessionStore } from './sessionStore'
import type { SessionUser } from './sessionStore'
import { env } from '@/shared/config/env'
import { AppError } from '@/shared/api/http/errors'

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

export async function bootstrap(): Promise<void> {
  const { setAuthed, setAnon } = useSessionStore.getState()

  let user: SessionUser | null = null

  if (env.VITE_AUTH_ME_AVAILABLE) {
    user = await bootstrapViaAuthMe()
  } else {
    user = await bootstrapViaProbe()
  }

  if (user) {
    setAuthed(user)
  } else {
    setAnon()
  }
}

eventBus.on('session:expired', () => {
  useSessionStore.getState().setAnon()
})
