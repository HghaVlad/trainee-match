import { http, HttpResponse } from 'msw'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { server } from '@/test/msw/server'
import { httpClient, SessionExpiredError } from './client'

beforeEach(() => {
  server.resetHandlers()
})

describe('single-flight 401 refresh', () => {
  it('calls /auth/refresh exactly once when 5 requests get 401', async () => {
    let refreshCallCount = 0
    let callCount = 0

    server.use(
      http.get('http://localhost:8080/api/test', () => {
        callCount++
        if (callCount <= 5) {
          return new HttpResponse(null, { status: 401 })
        }
        return HttpResponse.json({ ok: true })
      }),
      http.post('http://localhost:8080/auth/refresh', () => {
        refreshCallCount++
        return HttpResponse.json({ message: 'OK' })
      }),
    )

    const requests = Array.from({ length: 5 }, () =>
      httpClient.get('/api/test'),
    )

    const results = await Promise.all(requests)
    expect(refreshCallCount).toBe(1)
    results.forEach((r) => expect(r.data).toEqual({ ok: true }))
  })
})

describe('refresh fails → session:expired', () => {
  it('fires session:expired event and rejects with SessionExpiredError', async () => {
    server.use(
      http.get('http://localhost:8080/api/protected', () =>
        new HttpResponse(null, { status: 401 }),
      ),
      http.post('http://localhost:8080/auth/refresh', () =>
        new HttpResponse(null, { status: 401 }),
      ),
    )

    const eventFired = new Promise<void>((resolve) => {
      window.addEventListener('session:expired', () => resolve(), { once: true })
    })

    await expect(httpClient.get('/api/protected')).rejects.toBeInstanceOf(
      SessionExpiredError,
    )
    await eventFired
  })
})
