import { http, HttpResponse } from 'msw'

export const handlers = [
  http.post('/auth/login', () => HttpResponse.json({ message: 'OK' })),
  http.post('/auth/logout', () => HttpResponse.json({ message: 'OK' })),
  http.post('/auth/refresh', () => HttpResponse.json({ message: 'OK' })),

  http.get('/api/v1/candidate/me', () =>
    HttpResponse.json({
      id: 1,
      username: 'testuser',
      email: 'test@example.com',
      firstName: 'Test',
      lastName: 'User',
    }),
  ),

  http.get('/api/v1/companies', () =>
    HttpResponse.json({ content: [], nextCursor: null }),
  ),

  http.get('/api/v1/vacancies', () =>
    HttpResponse.json({ content: [], nextCursor: null }),
  ),
]
