import { http } from 'msw';

// Исходные вакансии
let vacancies = [
  { id: 1, title: 'Frontend Developer', companyName: 'Google', status: undefined },
  { id: 2, title: 'Backend Developer', companyName: 'Meta', status: 'SENT' },
  { id: 3, title: 'QA Engineer', companyName: 'Amazon', status: 'REJECTED' },
];

export const handlers = [
  // GET список вакансий
  http.get('*/vacancies', () => {
    return new Response(JSON.stringify(vacancies), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });
  }),

  // Подать заявку
  http.post('*/applications/:id/apply', async ({ params }) => {
    const id = Number(params.id);
    vacancies = vacancies.map(v =>
      v.id === id ? { ...v, status: 'SENT' } : v
    );

    return new Response(null, { status: 200 });
  }),

  // Отозвать заявку
  http.delete('*/applications/:id/withdraw', async ({ params }) => {
    const id = Number(params.id);
    vacancies = vacancies.map(v =>
      v.id === id ? { ...v, status: undefined } : v
    );

    return new Response(null, { status: 200 });
  }),
];
