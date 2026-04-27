import { createBrowserRouter, RouterProvider } from 'react-router'
import { lazy, Suspense } from 'react'
import { requireAuth, requireRole, redirectIfAuth } from './guards'
import { RootLayout } from '@/widgets/RootLayout'

function Placeholder({ name }: { name: string }) {
  return (
    <div style={{ padding: '2rem' }}>
      <h1>{name}</h1>
    </div>
  )
}

const LoginPage = lazy(() => import('@/pages/login'))
const RegisterPage = lazy(() => import('@/pages/register'))
const NotFoundPage = lazy(() => import('@/pages/NotFound'))
const ForbiddenPage = lazy(() => import('@/pages/Forbidden'))
const CandidateProfilePage = lazy(() => import('@/pages/me/profile'))
const ResumesPage = lazy(() => import('@/pages/me/resumes'))
const ResumeEditPage = lazy(() => import('@/pages/me/resumes/edit'))
const CompaniesPage = lazy(() => import('@/pages/companies'))
const CompanyDetailPage = lazy(() => import('@/pages/companies/detail'))
const VacanciesPage = lazy(() => import('@/pages/vacancies'))
const VacancyDetailPage = lazy(() => import('@/pages/vacancies/detail'))
const CompanyMePage = lazy(() => import('@/pages/company/me'))
const CompanyVacanciesPage = lazy(() => import('@/pages/company/vacancies'))
const CompanyVacancyEditPage = lazy(() => import('@/pages/company/vacancies/edit'))
const ApplicationsStub = lazy(() => import('@/pages/stubs/Applications'))
const InterviewsStub = lazy(() => import('@/pages/stubs/Interviews'))
const OffersStub = lazy(() => import('@/pages/stubs/Offers'))

function lazyEl(El: React.LazyExoticComponent<() => React.JSX.Element>) {
  return (
    <Suspense fallback={null}>
      <El />
    </Suspense>
  )
}

const router = createBrowserRouter([
  {
    element: <RootLayout />,
    children: [
      {
        path: '/login',
        loader: redirectIfAuth,
        element: (
          <Suspense fallback={null}>
            <LoginPage />
          </Suspense>
        ),
      },
      {
        path: '/register',
        loader: redirectIfAuth,
        element: (
          <Suspense fallback={null}>
            <RegisterPage />
          </Suspense>
        ),
      },
      { path: '/', element: <Placeholder name="Home / Vacancies" /> },
      { path: '/vacancies', element: lazyEl(VacanciesPage) },
      { path: '/vacancies/:id', element: lazyEl(VacancyDetailPage) },
      { path: '/companies', element: lazyEl(CompaniesPage) },
      { path: '/companies/:id', element: lazyEl(CompanyDetailPage) },
      {
        path: '/me',
        loader: requireAuth,
        children: [
          { path: 'profile', element: lazyEl(CandidateProfilePage) },
          { path: 'resumes', element: lazyEl(ResumesPage) },
          { path: 'resumes/:id', element: lazyEl(ResumeEditPage) },
        ],
      },
      {
        path: '/company',
        loader: requireRole('Company'),
        children: [
          { path: 'me', element: lazyEl(CompanyMePage) },
          { path: 'members', element: lazyEl(CompanyMePage) },
          { path: 'vacancies', element: lazyEl(CompanyVacanciesPage) },
          { path: 'vacancies/:id', element: lazyEl(CompanyVacancyEditPage) },
        ],
      },
      { path: '/applications', element: lazyEl(ApplicationsStub) },
      { path: '/interviews', element: lazyEl(InterviewsStub) },
      { path: '/offers', element: lazyEl(OffersStub) },
      {
        path: '/403',
        element: (
          <Suspense fallback={null}>
            <ForbiddenPage />
          </Suspense>
        ),
      },
      {
        path: '*',
        element: (
          <Suspense fallback={null}>
            <NotFoundPage />
          </Suspense>
        ),
      },
    ],
  },
])

export function AppRouter() {
  return <RouterProvider router={router} />
}
