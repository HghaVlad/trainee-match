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
  // Auth
  http.post('/auth/login', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as { username?: string }
    currentUser = userFor(body.username ?? 'candidate')
    return HttpResponse.json({ message: 'OK' })
  }),
  http.post('/auth/register', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as { username?: string; role?: string }
    currentUser = userFor(body.username ?? 'newuser')
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

  // Candidate
  http.get('/api/v1/candidate/me', () => {
    if (!currentUser || currentUser.role !== 'Candidate') {
      return new HttpResponse(null, { status: 401 })
    }
    return HttpResponse.json(currentUser)
  }),
  http.post('/api/v1/candidate', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    if (currentUser) currentUser = { ...currentUser, ...body }
    return HttpResponse.json(currentUser)
  }),
  http.patch('/api/v1/candidate', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    if (currentUser) currentUser = { ...currentUser, ...body }
    return HttpResponse.json(currentUser)
  }),

  // Resume
  http.get('/api/v1/resume', () => HttpResponse.json({ content: [], resumes: [], nextCursor: null })),
  http.post('/api/v1/resume', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 1, ...body }, { status: 201 })
  }),
  http.get('/api/v1/resume/:id', () => HttpResponse.json({ id: 1, title: 'Test Resume' })),
  http.patch('/api/v1/resume/:id', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 1, ...body })
  }),

  // Skills
  http.get('/api/v1/skills', () =>
    HttpResponse.json({
      content: [
        { id: 1, name: 'JavaScript' },
        { id: 2, name: 'TypeScript' },
        { id: 3, name: 'React' },
        { id: 4, name: 'Node.js' },
        { id: 5, name: 'Python' },
      ],
    }),
  ),

  // Companies (public)
  http.get('/api/v1/companies', () =>
    HttpResponse.json({ content: [], companies: [], nextCursor: null }),
  ),
  http.get('/api/v1/companies/:id', () =>
    HttpResponse.json({ id: 1, name: 'Test Company', description: 'Test' }),
  ),

  // Company (owned)
  http.get('/api/v1/companies/me', () => {
    if (!currentUser || currentUser.role !== 'Company') {
      return new HttpResponse(null, { status: 401 })
    }
    return HttpResponse.json({ id: currentUser.id, name: currentUser.username, description: '' })
  }),
  http.post('/api/v1/companies', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 3, ...body }, { status: 201 })
  }),
  http.patch('/api/v1/companies/:id', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 1, ...body })
  }),

  // Company members
  http.get('/api/v1/companies/:id/members', () =>
    HttpResponse.json({ content: [], members: [] }),
  ),
  http.post('/api/v1/companies/:id/members', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 1, ...body }, { status: 201 })
  }),
  http.patch('/api/v1/companies/:id/members/:userId', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 1, ...body })
  }),
  http.delete('/api/v1/companies/:id/members/:userId', () =>
    HttpResponse.json({ message: 'OK' }),
  ),

  // Vacancies (public)
  http.get('/api/v1/vacancies', () =>
    HttpResponse.json({ content: [], vacancies: [], nextCursor: null }),
  ),
  http.get('/api/v1/vacancies/:id', () =>
    HttpResponse.json({
      id: 1,
      title: 'Test Vacancy',
      description: 'Test description',
      salary_min: 50000,
      salary_max: 100000,
      workFormat: 'fullTime',
    }),
  ),

  // Company vacancies
  http.get('/api/v1/companies/:id/vacancies', () =>
    HttpResponse.json({ content: [], vacancies: [], nextCursor: null }),
  ),
  http.post('/api/v1/companies/:id/vacancies', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 1, ...body }, { status: 201 })
  }),
  http.patch('/api/v1/companies/:id/vacancies/:vacancyId', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: 1, ...body })
  }),
  http.delete('/api/v1/companies/:id/vacancies/:vacancyId', () =>
    HttpResponse.json({ message: 'OK' }),
  ),
  http.post('/api/v1/companies/:id/vacancies/:vacancyId/publish', () =>
    HttpResponse.json({ id: 1, status: 'published' }),
  ),
  http.post('/api/v1/companies/:id/vacancies/:vacancyId/archive', () =>
    HttpResponse.json({ id: 1, status: 'archived' }),
  ),
]
