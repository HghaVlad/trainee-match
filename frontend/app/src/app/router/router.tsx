import { createBrowserRouter, RouterProvider } from 'react-router'
import { lazy, Suspense } from 'react'
import {
  redirectIfAuth,
  requireCompanyAdmin,
  requireCompanyMember,
  requireRole,
  resolveActiveCompany,
} from './guards'
import { RootLayout } from '@/widgets/RootLayout'

function Placeholder({ name }: { name: string }) {
  return (
    <div style={{ padding: '2rem' }}>
      <h1>{name}</h1>
      <p style={{ color: '#666' }}>Screen pending implementation.</p>
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
const CompanyDashboardPage = lazy(() => import('@/pages/company/dashboard'))
const CompanyVacancyAnalyticsPage = lazy(
  () => import('@/pages/company/vacancies/analytics'),
)
const CompanyVacanciesPage = lazy(() => import('@/pages/company/vacancies'))
const CompanyVacancyNewPage = lazy(() => import('@/pages/company/vacancies/new'))
const CompanyVacancyEditPage = lazy(() => import('@/pages/company/vacancies/edit'))
const MyApplicationsPage = lazy(() => import('@/pages/me/applications'))
const MyApplicationDetailPage = lazy(() => import('@/pages/me/applications/detail'))
const CompanyNewPage = lazy(() => import('@/pages/company/new'))
const CompanyProfilePage = lazy(() => import('@/pages/company/profile'))
const CompanyMembersPage = lazy(() => import('@/pages/company/members'))
const CompanyApplicationsPage = lazy(
  () => import('@/pages/company/applications'),
)
const CompanyApplicationDetailPage = lazy(
  () => import('@/pages/company/applications/detail'),
)
const CompanyVacancyApplicationsPage = lazy(
  () => import('@/pages/company/vacancies/applications'),
)

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
        element: lazyEl(LoginPage),
      },
      {
        path: '/register',
        loader: redirectIfAuth,
        element: lazyEl(RegisterPage),
      },
      { path: '/', element: <Placeholder name="Home" /> },
      { path: '/vacancies', element: lazyEl(VacanciesPage) },
      { path: '/vacancies/:vacancyId', element: lazyEl(VacancyDetailPage) },
      { path: '/companies', element: lazyEl(CompaniesPage) },
      { path: '/companies/:companyId', element: lazyEl(CompanyDetailPage) },
      {
        path: '/me',
        loader: requireRole('Candidate'),
        children: [
          { path: 'profile', element: lazyEl(CandidateProfilePage) },
          { path: 'resumes', element: lazyEl(ResumesPage) },
          { path: 'resumes/new', element: lazyEl(ResumeEditPage) },
          { path: 'resumes/:resumeId', element: lazyEl(ResumeEditPage) },
          {
            path: 'applications',
            element: lazyEl(MyApplicationsPage),
          },
          {
            path: 'applications/:applicationId',
            element: lazyEl(MyApplicationDetailPage),
          },
        ],
      },
      {
        path: '/company',
        loader: requireRole('Company'),
        children: [
          { index: true, loader: resolveActiveCompany, element: null },
            { path: 'new', element: lazyEl(CompanyNewPage) },
          {
            path: ':companyId',
            loader: requireCompanyMember,
            children: [
              { index: true, element: <Placeholder name="Company" /> },
              { path: 'dashboard', element: lazyEl(CompanyDashboardPage) },
              { path: 'profile', element: lazyEl(CompanyProfilePage) },
              {
                path: 'members',
                element: lazyEl(CompanyMembersPage),
              },
              {
                path: 'vacancies',
                children: [
                  { index: true, element: lazyEl(CompanyVacanciesPage) },
                  { path: 'new', element: lazyEl(CompanyVacancyNewPage) },
                  {
                    path: ':vacancyId',
                    children: [
                      { index: true, element: lazyEl(CompanyVacancyEditPage) },
                      {
                        path: 'applications',
                        element: lazyEl(CompanyVacancyApplicationsPage),
                      },
                      {
                        path: 'analytics',
                        element: lazyEl(CompanyVacancyAnalyticsPage),
                      },
                    ],
                  },
                ],
              },
              {
                path: 'applications',
                children: [
                  {
                    index: true,
                    element: lazyEl(CompanyApplicationsPage),
                  },
                  {
                    path: ':applicationId',
                    element: lazyEl(CompanyApplicationDetailPage),
                  },
                ],
              },
            ],
          },
          {
            path: ':companyId/admin',
            loader: requireCompanyAdmin,
            children: [
              {
                path: 'danger',
                element: <Placeholder name="Company Danger Zone" />,
              },
            ],
          },
        ],
      },
      {
        path: '/403',
        element: lazyEl(ForbiddenPage),
      },
      {
        path: '*',
        element: lazyEl(NotFoundPage),
      },
    ],
  },
])

export function AppRouter() {
  return <RouterProvider router={router} />
}
