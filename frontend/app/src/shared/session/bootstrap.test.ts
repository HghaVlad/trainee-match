import { describe, it, expect, beforeEach } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '@/test/msw/server'
import { bootstrap } from './bootstrap'
import { useSessionStore } from './sessionStore'

describe('bootstrap', () => {
  beforeEach(() => {
    useSessionStore.setState({ status: 'unknown', user: undefined })
    server.resetHandlers()
  })

  it('sets authed when /candidate/me returns 200', async () => {
    server.use(
      http.get('*/candidate/me', () =>
        HttpResponse.json({
          id: 1,
          username: 'testuser',
          email: 'test@test.com',
        }),
      ),
    )

    await bootstrap()

    const state = useSessionStore.getState()
    expect(state.status).toBe('authed')
    expect(state.user?.role).toBe('Candidate')
  })

  it('sets anon when /candidate/me returns 401', async () => {
    server.use(
      http.get('*/candidate/me', () => HttpResponse.json({}, { status: 401 })),
      http.post('*/auth/refresh', () => HttpResponse.json({}, { status: 401 })),
    )

    await bootstrap()

    const state = useSessionStore.getState()
    expect(state.status).toBe('anon')
  })
})
