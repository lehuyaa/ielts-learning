import { Outlet } from 'react-router-dom'

import { BrandMark } from '@/components/layout/BrandMark'

export function AuthLayout() {
  return (
    <main className="grid min-h-screen bg-background px-4 py-8 text-foreground">
      <div className="mx-auto grid w-full max-w-xl content-center gap-8">
        <BrandMark className="mx-auto" />
        <Outlet />
      </div>
    </main>
  )
}
