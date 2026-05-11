import { describe, it, expect, beforeEach } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '@/test/msw/server'
import { bootstrap } from './bootstrap'
import { useSessionStore } from './sessionStore'

describe('bootstrap', () => {
  beforeEach(() => {
    useSessionStore.setState({
      status: 'unknown',
      user: undefined,
      companies: [],
      activeCompanyId: undefined,
    })
    server.resetHandlers()
  })

  it('sets authed when /auth/me returns 200', async () => {
    server.use(
      http.post('*/auth/me', () =>
        HttpResponse.json({
          id: '1',
          username: 'testuser',
          email: 'test@test.com',
          first_name: 'Test',
          last_name: 'User',
          role: 'Candidate',
        }),
      ),
    )

    await bootstrap()

    const state = useSessionStore.getState()
    expect(state.status).toBe('authed')
    expect(state.user?.role).toBe('Candidate')
    expect(state.user?.firstName).toBe('Test')
  })

  it('sets anon when /auth/me returns 401', async () => {
    server.use(
      http.post('*/auth/me', () => HttpResponse.json({}, { status: 401 })),
      http.post('*/auth/refresh', () => HttpResponse.json({}, { status: 401 })),
    )

    await bootstrap()

    const state = useSessionStore.getState()
    expect(state.status).toBe('anon')
  })
})
