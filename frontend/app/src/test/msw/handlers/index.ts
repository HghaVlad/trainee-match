import { http, HttpResponse } from 'msw'

let currentUser:
  | { id: number; role: 'Candidate' | 'Company'; username: string; email: string; firstName: string; lastName: string }
  | null = null

function userFor(username: string): NonNullable<typeof currentUser> {
  const isCompany = username.toLowerCase().startsWith('company')
  return {
    id: isCompany ? 2 : 1,
    role: isCompany ? 'Company' : 'Candidate',
    username,
    email: `${username}@example.com`,
    firstName: 'Test',
    lastName: 'User',
  }
}

export const handlers = [
  http.post('/auth/login', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as { username?: string }
    currentUser = userFor(body.username ?? 'candidate')
    return HttpResponse.json({ message: 'OK' })
  }),
  http.post('/auth/logout', () => {
    currentUser = null
    return HttpResponse.json({ message: 'OK' })
  }),
  http.post('/auth/refresh', () => HttpResponse.json({ message: 'OK' })),

  http.get('/auth/me', () => {
    if (!currentUser) return new HttpResponse(null, { status: 401 })
    return HttpResponse.json(currentUser)
  }),

  http.get('/api/v1/candidate/me', () => {
    if (!currentUser || currentUser.role !== 'Candidate') {
      return new HttpResponse(null, { status: 401 })
    }
    return HttpResponse.json(currentUser)
  }),

  http.get('/api/v1/companies', () =>
    HttpResponse.json({ content: [], companies: [], nextCursor: null }),
  ),

  http.get('/api/v1/vacancies', () =>
    HttpResponse.json({ content: [], vacancies: [], nextCursor: null }),
  ),
]
