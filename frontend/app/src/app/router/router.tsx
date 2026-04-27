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
      { path: '/vacancies', element: <Placeholder name="Vacancy List" /> },
      { path: '/vacancies/:id', element: <Placeholder name="Vacancy Detail" /> },
      { path: '/companies', element: <Placeholder name="Companies" /> },
      { path: '/companies/:id', element: <Placeholder name="Company Detail" /> },
      {
        path: '/me',
        loader: requireAuth,
        children: [
          { path: 'profile', element: (
            <Suspense fallback={null}>
              <CandidateProfilePage />
            </Suspense>
          ) },
          { path: 'resumes', element: <Placeholder name="Resumes" /> },
          { path: 'resumes/:id', element: <Placeholder name="Resume Edit" /> },
        ],
      },
      {
        path: '/company',
        loader: requireRole('Company'),
        children: [
          { path: 'me', element: <Placeholder name="Company Profile" /> },
          { path: 'members', element: <Placeholder name="Company Members" /> },
          { path: 'vacancies', element: <Placeholder name="Company Vacancies" /> },
          { path: 'vacancies/:id', element: <Placeholder name="Vacancy Edit" /> },
        ],
      },
      { path: '/applications', element: <Placeholder name="Applications (Coming Soon)" /> },
      { path: '/interviews', element: <Placeholder name="Interviews (Coming Soon)" /> },
      { path: '/offers', element: <Placeholder name="Offers (Coming Soon)" /> },
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
