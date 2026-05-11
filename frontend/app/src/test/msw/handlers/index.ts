import { http, HttpResponse } from 'msw'

function uuid(index: number): string {
  const hex = index.toString(16).padStart(12, '0')
  return `00000000-0000-4000-8000-${hex}`
}

function daysAgo(days: number): string {
  const d = new Date()
  d.setDate(d.getDate() - days)
  return d.toISOString()
}

function monthsAgo(months: number): string {
  const d = new Date()
  d.setMonth(d.getMonth() - months)
  return d.toISOString()
}

let currentUser:
  | {
      id: string
      role: 'Candidate' | 'Company'
      username: string
      email: string
      first_name: string
      last_name: string
    }
  | null = null

function userFor(username: string): NonNullable<typeof currentUser> {
  const isCompany = username.toLowerCase().startsWith('company')
  return {
    id: isCompany ? uuid(2) : uuid(1),
    role: isCompany ? 'Company' : 'Candidate',
    username,
    email: `${username}@example.com`,
    first_name: 'Test',
    last_name: 'User',
  }
}

const COMPANY_ID = uuid(2)
const VACANCY_1_ID = uuid(10)
const VACANCY_2_ID = uuid(11)
const VACANCY_3_ID = uuid(12)

const mockApplications = [
  {
    id: uuid(100),
    status: 'submitted' as const,
    vacancyId: VACANCY_1_ID,
    vacancyTitle: 'Frontend Developer (React)',
    companyId: COMPANY_ID,
    companyName: 'TechCorp',
    createdAt: daysAgo(3),
    updatedAt: daysAgo(3),
  },
  {
    id: uuid(101),
    status: 'interview' as const,
    vacancyId: VACANCY_2_ID,
    vacancyTitle: 'Backend Go Developer',
    companyId: COMPANY_ID,
    companyName: 'TechCorp',
    createdAt: daysAgo(10),
    updatedAt: daysAgo(2),
  },
  {
    id: uuid(102),
    status: 'seen' as const,
    vacancyId: uuid(20),
    vacancyTitle: 'DevOps Engineer',
    companyId: uuid(3),
    companyName: 'CloudInc',
    createdAt: daysAgo(7),
    updatedAt: daysAgo(5),
  },
  {
    id: uuid(103),
    status: 'rejected' as const,
    vacancyId: uuid(21),
    vacancyTitle: 'Data Scientist',
    companyId: uuid(4),
    companyName: 'DataLab',
    createdAt: daysAgo(20),
    updatedAt: daysAgo(15),
  },
  {
    id: uuid(104),
    status: 'offer' as const,
    vacancyId: uuid(22),
    vacancyTitle: 'Junior UX Designer',
    companyId: uuid(5),
    companyName: 'DesignStudio',
    createdAt: daysAgo(30),
    updatedAt: daysAgo(1),
  },
  {
    id: uuid(105),
    status: 'withdrawn' as const,
    vacancyId: VACANCY_3_ID,
    vacancyTitle: 'QA Engineer',
    companyId: COMPANY_ID,
    companyName: 'TechCorp',
    createdAt: daysAgo(14),
    updatedAt: daysAgo(8),
  },
]

const mockCompanyVacancies = [
  {
    id: VACANCY_1_ID,
    title: 'Frontend Developer (React)',
    city: 'Moscow',
    createdAt: daysAgo(30),
    status: 'published' as const,
    workFormat: 'hybrid' as const,
    employmentType: 'fullTime' as const,
    salaryFrom: 150000,
    salaryTo: 250000,
    isPaid: true,
  },
  {
    id: VACANCY_2_ID,
    title: 'Backend Go Developer',
    city: 'Saint Petersburg',
    createdAt: daysAgo(14),
    status: 'published' as const,
    workFormat: 'remote' as const,
    employmentType: 'fullTime' as const,
    salaryFrom: 200000,
    salaryTo: 350000,
    isPaid: true,
  },
  {
    id: VACANCY_3_ID,
    title: 'QA Engineer',
    city: 'Moscow',
    createdAt: daysAgo(7),
    status: 'published' as const,
    workFormat: 'office' as const,
    employmentType: 'fullTime' as const,
    salaryFrom: 120000,
    salaryTo: 180000,
    isPaid: true,
  },
]

const mockCompanyMembers = [
  { userId: uuid(30), username: 'ivan.petrov', email: 'ivan@techcorp.com', role: 'admin' as const, companyId: COMPANY_ID },
  { userId: uuid(31), username: 'anna.smirnova', email: 'anna@techcorp.com', role: 'admin' as const, companyId: COMPANY_ID },
  { userId: uuid(32), username: 'petr.sidorov', email: 'petr@techcorp.com', role: 'recruiter' as const, companyId: COMPANY_ID },
  { userId: uuid(33), username: 'elena.kozlov', email: 'elena@techcorp.com', role: 'recruiter' as const, companyId: COMPANY_ID },
  { userId: uuid(34), username: 'dmitry.volkov', email: 'dmitry@techcorp.com', role: 'recruiter' as const, companyId: COMPANY_ID },
]

const candidateNames = ['Алиса Иванова', 'Боб Смирнов', 'Чарли Козлов', 'Диана Петрова', 'Евгений Сидоров']
const candidateEmails = ['alisa@example.com', 'bob@example.com', 'charlie@example.com', 'diana@example.com', 'evgeny@example.com']
const candidateTelegrams = ['@alisa_i', '@bob_s', '@charlie_k', '@diana_p', '@evgeny_s']

function makeHrApplication(candidateIndex: number, appStatus: string, vacId: string, vacTitle: string) {
  const id = uuid(200 + candidateIndex)
  const name = candidateNames[candidateIndex % 5]
  const email = candidateEmails[candidateIndex % 5]
  const telegram = candidateTelegrams[candidateIndex % 5]
  return {
    id,
    status: appStatus,
    vacancyId: vacId,
    vacancyTitle: vacTitle,
    coverLetter: candidateIndex % 2 === 0 ? 'Я очень заинтересован в этой позиции. Имею соответствующий опыт и готов развиваться.' : null,
    snapshot: {
      fullName: name,
      email,
      telegram,
      resumeData: {
        lastName: name.split(' ')[1] || 'Фамилия',
        firstName: name.split(' ')[0] || 'Имя',
        middleName: '',
        dateOfBirth: '2000-01-15',
        email,
        phone: '+7-999-123-45-6' + (candidateIndex % 10),
        city: ['Москва', 'Санкт-Петербург', 'Казань', 'Новосибирск', 'Екатеринбург'][candidateIndex % 5],
        citizenship: 'РФ',
        education: [
          {
            level: 'bachelor',
            university: 'МГУ им. Ломоносова',
            faculty: 'ВМК',
            specialization: 'Прикладная математика и информатика',
            startYear: 2018,
            endYear: 2022,
            format: 'full-time',
          },
        ],
        workExperiences: [
          {
            position: ['Junior Developer', 'Стажёр', 'Frontend Developer', 'Аналитик', 'Тестировщик'][candidateIndex % 5],
            company: ['TechCorp', 'StartupX', 'BigCo', 'Digital Agency', 'Bank'][candidateIndex % 5],
            period: `${monthsAgo(18).slice(0, 7)} — настоящее время`,
            responsibilities: 'Разработка и поддержка веб-приложений, code review, участие в планировании.',
          },
        ],
        skillsList: ['JavaScript', 'TypeScript', 'React', 'Git', 'Docker'],
        desiredFormat: 'remote',
        englishLevel: 'B1',
      },
      createdAt: daysAgo(5 + candidateIndex),
    },
    statusHistory: [
      { status: 'submitted', changedByRole: 'candidate' as const, createdAt: daysAgo(5 + candidateIndex) },
    ],
    allowedActions: [],
    createdAt: daysAgo(5 + candidateIndex),
    updatedAt: daysAgo(1 + candidateIndex),
  }
}

const hrApplications = mockCompanyVacancies.flatMap((vac, vi) =>
  Array.from({ length: 2 + vi }, (_, ci) => makeHrApplication(vi * 3 + ci, ['submitted', 'interview', 'seen', 'rejected', 'offer'][(vi + ci) % 5], vac.id!, vac.title!)),
)

export const handlers = [
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
  http.post('/auth/me', () => {
    if (!currentUser) return new HttpResponse(null, { status: 401 })
    return HttpResponse.json(currentUser)
  }),

  http.get('/api/v1/candidate/me', () => {
    if (!currentUser || currentUser.role !== 'Candidate') {
      return new HttpResponse(null, { status: 401 })
    }
    return HttpResponse.json({
      id: currentUser.id,
      firstName: currentUser.first_name,
      lastName: currentUser.last_name,
      email: currentUser.email,
      username: currentUser.username,
      phone: '+7-999-123-45-67',
      city: 'Москва',
    })
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

  http.get('/api/v1/resume', () =>
    HttpResponse.json({
      content: [
        {
          id: uuid(50),
          title: 'My Resume',
          candidateName: 'Test User',
          createdAt: daysAgo(30),
          updatedAt: daysAgo(1),
        },
      ],
      resumes: [],
      nextCursor: null,
    }),
  ),
  http.post('/api/v1/resume', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: uuid(51), ...body }, { status: 201 })
  }),
  http.get('/api/v1/resume/:id', () =>
    HttpResponse.json({
      id: uuid(50),
      title: 'My Resume',
      candidateName: 'Test User',
      skills: [
        { id: 1, name: 'JavaScript' },
        { id: 2, name: 'TypeScript' },
        { id: 3, name: 'React' },
      ],
      workExperience: [
        {
          company: 'TechCorp',
          position: 'Frontend Developer',
          startDate: monthsAgo(18),
          current: true,
          description: 'Building UI components with React and TypeScript',
        },
      ],
      education: [
        {
          institution: 'Bauman Moscow State Technical University',
          degree: "Bachelor's",
          field: 'Computer Science',
          startDate: monthsAgo(48),
          endDate: monthsAgo(24),
        },
      ],
      createdAt: daysAgo(30),
      updatedAt: daysAgo(1),
    }),
  ),
  http.patch('/api/v1/resume/:id', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: uuid(50), ...body })
  }),

  http.get('/api/v1/skills', () =>
    HttpResponse.json({
      content: [
        { id: 1, name: 'JavaScript' },
        { id: 2, name: 'TypeScript' },
        { id: 3, name: 'React' },
        { id: 4, name: 'Node.js' },
        { id: 5, name: 'Python' },
        { id: 6, name: 'Go' },
        { id: 7, name: 'PostgreSQL' },
        { id: 8, name: 'Docker' },
        { id: 9, name: 'Kubernetes' },
        { id: 10, name: 'Figma' },
      ],
    }),
  ),

  http.get('/api/v1/applications', () =>
    HttpResponse.json({
      data: mockApplications,
      nextCursor: null,
      hasNext: false,
    }),
  ),

  http.get('/api/v1/applications/:id/history', ({ params }) => {
    const app = mockApplications.find((a) => a.id === params.id)
    if (!app) return new HttpResponse(null, { status: 404 })
    return HttpResponse.json({
      data: [
        { status: app.status, changedByRole: 'candidate', createdAt: app.updatedAt },
      ],
    })
  }),

  http.get('/api/v1/applications/:id', ({ params }) => {
    const app = mockApplications.find((a) => a.id === params.id)
    if (!app) return new HttpResponse(null, { status: 404 })

    return HttpResponse.json({
      data: {
        id: app.id,
        status: app.status,
        vacancyId: app.vacancyId,
        vacancyTitle: app.vacancyTitle,
        companyId: app.companyId,
        companyName: app.companyName,
        coverLetter: 'Я работаю с React более 3 лет и хотел бы присоединиться к вашей команде.',
        snapshot: {
          fullName: 'Test User',
          email: 'test.user@example.com',
          telegram: '@testuser',
          resumeData: {
            lastName: 'User',
            firstName: 'Test',
            middleName: '',
            dateOfBirth: '2000-05-20',
            email: 'test.user@example.com',
            phone: '+7-999-111-22-33',
            city: 'Москва',
            citizenship: 'РФ',
            education: [
              {
                level: 'bachelor',
                university: 'МГТУ им. Баумана',
                faculty: 'ИУ',
                specialization: 'Программная инженерия',
                startYear: 2018,
                endYear: 2022,
                format: 'full-time',
              },
            ],
            workExperiences: [
              {
                position: 'Frontend Developer',
                company: 'Previous Co',
                period: '2022 — настоящее время',
                responsibilities: 'Разработка React-приложений, оптимизация производительности.',
              },
            ],
            skillsList: ['JavaScript', 'TypeScript', 'React', 'Redux', 'Jest'],
            desiredFormat: 'hybrid',
            englishLevel: 'B2',
          },
          createdAt: app.createdAt,
        },
        statusHistory: [
          { status: app.status, changedByRole: 'candidate', createdAt: app.updatedAt },
        ],
        allowedActions:
          app.status === 'submitted' || app.status === 'seen' || app.status === 'interview'
            ? ['withdraw']
            : [],
        createdAt: app.createdAt,
        updatedAt: app.updatedAt,
      },
    })
  }),

  http.post('/api/v1/applications', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json(
      {
        data: {
          id: uuid(200),
          status: 'submitted',
          vacancyId: body.vacancyId ?? uuid(10),
          vacancyTitle: 'New Vacancy',
          companyId: COMPANY_ID,
          companyName: 'TechCorp',
          coverLetter: body.coverLetter ?? null,
          snapshot: {
            fullName: 'Test User',
            email: 'test.user@example.com',
            resumeData: {
              lastName: 'User',
              firstName: 'Test',
              middleName: '',
              dateOfBirth: '2000-05-20',
              email: 'test.user@example.com',
              phone: '+7-999-111-22-33',
              city: 'Москва',
              citizenship: 'РФ',
              education: [],
              workExperiences: [],
              skillsList: [],
              desiredFormat: 'remote',
              englishLevel: 'B1',
            },
            createdAt: new Date().toISOString(),
          },
          statusHistory: [],
          allowedActions: ['withdraw'],
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        },
      },
      { status: 201 },
    )
  }),

  http.get('/api/v1/companies', () =>
    HttpResponse.json({
      companies: [
        {
          id: COMPANY_ID,
          logoUrl: 'https://via.placeholder.com/64',
          name: 'TechCorp',
          openVacanciesCount: 3,
        },
        {
          id: uuid(3),
          logoUrl: 'https://via.placeholder.com/64',
          name: 'CloudInc',
          openVacanciesCount: 5,
        },
        {
          id: uuid(4),
          logoUrl: 'https://via.placeholder.com/64',
          name: 'DataLab',
          openVacanciesCount: 2,
        },
      ],
      nextCursor: null,
    }),
  ),

  http.get('/api/v1/companies/:id', ({ params }) => {
    const isOwn = params.id === COMPANY_ID || params.id === currentUser?.id
    return HttpResponse.json({
      id: params.id,
      name: isOwn ? 'TechCorp' : 'External Company',
      description: isOwn
        ? 'TechCorp — ведущая технологическая компания, специализирующаяся на разработке ПО, IT-консалтинге и цифровой трансформации.'
        : 'An external company description.',
      logoURL: 'https://via.placeholder.com/128',
      website: 'https://techcorp.example.com',
      openVacanciesCount: isOwn ? 3 : 1,
      createdAt: monthsAgo(24),
      updatedAt: daysAgo(7),
    })
  }),

  http.get('/api/v1/companies/me', () => {
    if (!currentUser || currentUser.role !== 'Company') {
      return new HttpResponse(null, { status: 401 })
    }
    return HttpResponse.json({
      id: currentUser.id,
      name: 'TechCorp',
      description: 'TechCorp — ведущая технологическая компания, специализирующаяся на разработке ПО.',
      logoURL: 'https://via.placeholder.com/128',
      website: 'https://techcorp.example.com',
      openVacanciesCount: 3,
      createdAt: monthsAgo(24),
      updatedAt: daysAgo(7),
    })
  }),

  http.post('/api/v1/companies', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: uuid(99), ...body }, { status: 201 })
  }),

  http.patch('/api/v1/companies/:id', async ({ request, params }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: params.id, ...body })
  }),

  http.get('/api/v1/companies/:id/members', () =>
    HttpResponse.json({
      members: mockCompanyMembers,
      hasMore: false,
    }),
  ),

  http.post('/api/v1/companies/:id/members', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: uuid(40), ...body }, { status: 201 })
  }),

  http.patch('/api/v1/companies/:id/members/:userId', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ ...body })
  }),

  http.delete('/api/v1/companies/:id/members/:userId', () =>
    HttpResponse.json({ message: 'OK' }),
  ),

  http.get('/api/v1/vacancies', () =>
    HttpResponse.json({
      content: mockCompanyVacancies.map((v) => ({
        id: v.id,
        title: v.title,
        city: v.city,
        salaryFrom: v.salaryFrom,
        salaryTo: v.salaryTo,
        workFormat: v.workFormat,
        employmentType: v.employmentType,
        companyId: COMPANY_ID,
        companyName: 'TechCorp',
        createdAt: v.createdAt,
      })),
      vacancies: [],
      nextCursor: null,
    }),
  ),

  http.get('/api/v1/vacancies/:id', ({ params }) => {
    const vac = mockCompanyVacancies.find((v) => v.id === params.id)
    if (!vac) return new HttpResponse(null, { status: 404 })
    return HttpResponse.json({
      id: vac.id,
      title: vac.title,
      description: `We are looking for a talented ${vac.title} to join our team. You will work on cutting-edge projects with modern technologies.`,
      city: vac.city,
      salaryFrom: vac.salaryFrom,
      salaryTo: vac.salaryTo,
      workFormat: vac.workFormat,
      employmentType: vac.employmentType,
      isPaid: vac.isPaid,
      companyId: COMPANY_ID,
      companyName: 'TechCorp',
      createdAt: vac.createdAt,
      status: vac.status,
    })
  }),

  http.get('/api/v1/companies/:id/vacancies', ({ params }) => {
    if (params.id === COMPANY_ID || params.id === currentUser?.id) {
      return HttpResponse.json({
        vacancies: mockCompanyVacancies,
        nextCursor: null,
      })
    }
    return HttpResponse.json({
      vacancies: [
        {
          id: uuid(50),
          title: 'Software Engineer',
          city: 'Moscow',
          createdAt: daysAgo(10),
          status: 'published',
          workFormat: 'remote',
          employmentType: 'fullTime',
          salaryFrom: 180000,
          salaryTo: 300000,
          isPaid: true,
        },
      ],
      nextCursor: null,
    })
  }),

  http.post('/api/v1/companies/:id/vacancies', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: uuid(60), ...body }, { status: 201 })
  }),

  http.patch('/api/v1/companies/:id/vacancies/:vacancyId', async ({ request }) => {
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>
    return HttpResponse.json({ id: uuid(10), ...body })
  }),

  http.delete('/api/v1/companies/:id/vacancies/:vacancyId', () =>
    HttpResponse.json({ message: 'OK' }),
  ),

  http.post('/api/v1/companies/:id/vacancies/:vacancyId/publish', () =>
    HttpResponse.json({ id: uuid(10), status: 'published' }),
  ),

  http.post('/api/v1/companies/:id/vacancies/:vacancyId/archive', () =>
    HttpResponse.json({ id: uuid(10), status: 'archived' }),
  ),

  http.get('/api/v1/hr/companies/:companyId/analytics/summary', () =>
    HttpResponse.json({
      data: {
        totalApplications: 47,
        activeApplications: 12,
        submittedCount: 8,
        seenCount: 3,
        interviewCount: 1,
        rejectedCount: 25,
        offerCount: 2,
        withdrawnCount: 8,
        conversionToInterview: 0.255,
        conversionToOffer: 0.043,
      },
    }),
  ),

  http.get('/api/v1/hr/companies/:companyId/analytics/dynamics', () =>
    HttpResponse.json({
      data: Array.from({ length: 12 }, (_, i) => ({
        bucketStart: monthsAgo(11 - i),
        createdCount: Math.floor(Math.random() * 8) + 1,
        withdrawnCount: Math.floor(Math.random() * 3),
        seenCount: Math.floor(Math.random() * 5),
        interviewCount: Math.floor(Math.random() * 3),
        rejectedCount: Math.floor(Math.random() * 4),
        offerCount: Math.floor(Math.random() * 2),
      })),
    }),
  ),

  http.get('/api/v1/hr/vacancies/:vacancyId/analytics/summary', () =>
    HttpResponse.json({
      data: {
        totalApplications: 15,
        activeApplications: 4,
        submittedCount: 3,
        seenCount: 1,
        interviewCount: 0,
        rejectedCount: 8,
        offerCount: 1,
        withdrawnCount: 2,
        conversionToInterview: 0.067,
        conversionToOffer: 0.067,
      },
    }),
  ),

  http.get('/api/v1/hr/vacancies/:vacancyId/analytics/dynamics', () =>
    HttpResponse.json({
      data: Array.from({ length: 6 }, (_, i) => ({
        bucketStart: daysAgo(30 - i * 5),
        createdCount: Math.floor(Math.random() * 4) + 1,
        withdrawnCount: Math.floor(Math.random() * 2),
        seenCount: Math.floor(Math.random() * 3),
        interviewCount: Math.floor(Math.random() * 2),
        rejectedCount: Math.floor(Math.random() * 3),
        offerCount: Math.floor(Math.random() * 1),
      })),
    }),
  ),

  http.get('/api/v1/hr/companies/:companyId/applications', () =>
    HttpResponse.json({
      data: hrApplications.map((app) => ({
        id: app.id,
        status: app.status,
        vacancyId: app.vacancyId,
        vacancyTitle: app.vacancyTitle,
        snapshot: {
          fullName: app.snapshot.fullName,
          email: app.snapshot.email,
          telegram: app.snapshot.telegram,
          createdAt: app.snapshot.createdAt,
        },
        createdAt: app.createdAt,
        updatedAt: app.updatedAt,
      })),
      nextCursor: null,
      hasNext: false,
    }),
  ),

  http.get('/api/v1/hr/applications/:applicationId/history', ({ params }) => {
    const app = hrApplications.find((a) => a.id === params.applicationId)
    if (!app) return new HttpResponse(null, { status: 404 })
    return HttpResponse.json({
      data: [
        { status: app.status, changedByRole: 'candidate', changedByUserId: null, comment: null, createdAt: app.updatedAt },
      ],
    })
  }),

  http.get('/api/v1/hr/applications/:applicationId', ({ params }) => {
    const app = hrApplications.find((a) => a.id === params.applicationId)
    if (!app) return new HttpResponse(null, { status: 404 })

    return HttpResponse.json({
      data: {
        id: app.id,
        status: app.status,
        vacancyId: app.vacancyId,
        vacancyTitle: app.vacancyTitle,
        coverLetter: app.coverLetter,
        snapshot: app.snapshot,
        statusHistory: app.statusHistory,
        allowedActions: app.allowedActions,
        createdAt: app.createdAt,
        updatedAt: app.updatedAt,
      },
    })
  }),

  http.get('/api/v1/hr/vacancies/:vacancyId/applications', ({ params }) => {
    const apps = hrApplications.filter((a) => a.vacancyId === params.vacancyId)
    return HttpResponse.json({
      data: apps.map((app) => ({
        id: app.id,
        status: app.status,
        vacancyId: app.vacancyId,
        vacancyTitle: app.vacancyTitle,
        snapshot: {
          fullName: app.snapshot.fullName,
          email: app.snapshot.email,
          telegram: app.snapshot.telegram,
          createdAt: app.snapshot.createdAt,
        },
        createdAt: app.createdAt,
        updatedAt: app.updatedAt,
      })),
      nextCursor: null,
      hasNext: false,
    })
  }),

  http.post('/api/v1/hr/applications/:applicationId/status', async ({ request, params }) => {
    const body = (await request.json().catch(() => ({}))) as { status?: string }
    const app = hrApplications.find((a) => a.id === params.applicationId)
    if (!app) return new HttpResponse(null, { status: 404 })
    return HttpResponse.json({
      data: {
        ...app,
        status: body.status ?? app.status,
        updatedAt: new Date().toISOString(),
      },
    })
  }),
]
