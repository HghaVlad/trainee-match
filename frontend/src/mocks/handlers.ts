import { http } from 'msw';

interface Vacancy {
  id: number;
  title: string;
  companyName: string;
  status?: 'SENT' | 'REJECTED';
  city: string;
  format: 'remote' | 'office';
}

let vacancies: Vacancy[] = Array.from({ length: 57 }).map((_, i) => ({
  id: i + 1,
  title: i % 2 === 0 ? 'Frontend Developer' : 'Backend Developer',
  companyName: ['Google', 'Meta', 'Amazon'][i % 3],
  status: i % 4 === 0 ? 'REJECTED' : undefined,
  city: ['moscow', 'spb', 'kazan'][i % 3],
  format: i % 2 == 0 ? 'remote' : 'office',
}));

export const handlers = [
  // GET список вакансий
  http.get('*/vacancies', ({ request }) => {
    const url = new URL(request.url);

    const search = url.searchParams.get('search')?.toLowerCase() || '';
    const page = Number(url.searchParams.get('page') ?? 1);
    const size = Number(url.searchParams.get('size') ?? 5);

    // любые фильтры
    const city = url.searchParams.get('city');
    const format = url.searchParams.get('format');

    let filtered = [...vacancies];

    if (search) {
      filtered = filtered.filter(v =>
        v.title.toLowerCase().includes(search) ||
        v.companyName.toLowerCase().includes(search)
      );
    }

    if (city) {
      filtered = filtered.filter(v => v.city === city);
    }

    if (format) {
      filtered = filtered.filter(v => v.format === format);
    }

    const total = filtered.length;
    const start = (page - 1) * size;
    const end = start + size;
    const content = filtered.slice(start, end);

    // --- ответ ---
    return new Response(
      JSON.stringify({
        content,
        page,
        size,
        total,
        totalPages: Math.ceil(total / size),
      }),
      {
        headers: { 'Content-Type': 'application/json' },
      }
    );
  }),

  // Подать заявку
  http.post('*/applications/:id/apply', ({ params }) => {
    const id = Number(params.id);

    console.log("clicked");
    vacancies = vacancies.map(v =>
      v.id === id ? { ...v, status: 'SENT' } : v
    );

    return new Response(null, { status: 200 });
  }),

  http.delete('*/applications/:id/withdraw', ({ params }) => {
    const id = Number(params.id);

    vacancies = vacancies.map(v =>
      v.id === id ? { ...v, status: undefined } : v
    );

    return new Response(null, { status: 200 });
  })
];
