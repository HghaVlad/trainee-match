import { Outlet } from 'react-router'
import { Header } from '@/widgets/Header'
import { Toaster } from '@/shared/ui/toaster'

export function RootLayout() {
  return (
    <div>
      <Header />
      <main>
        <Outlet />
      </main>
      <Toaster />
    </div>
  )
}
