import { Link } from 'react-router-dom'

import { BrandMark } from '@/components/layout/BrandMark'

const footerLinks = [
  {
    label: 'Home',
    to: '/',
  },
  {
    label: 'Roadmap',
    to: '/roadmap',
  },
  {
    label: 'Register',
    to: '/register',
  },
  {
    label: 'Login',
    to: '/login',
  },
]

const socialLinks = ['LinkedIn', 'YouTube', 'X']

export function LandingFooter() {
  return (
    <footer className="border-t border-border py-16">
      <div className="grid gap-8 md:grid-cols-[1fr_auto] md:items-start">
        <div className="space-y-4">
          <BrandMark />
          <p className="max-w-md text-sm leading-6 text-muted-foreground">
            IELTS vocabulary learning with structured lessons, spaced
            repetition, quizzes, and visible progress toward your target band.
          </p>
          <p className="text-sm text-muted-foreground">
            Copyright 2026 LexPath. All rights reserved.
          </p>
        </div>

        <div className="grid gap-8 sm:grid-cols-2">
          <nav className="grid gap-3">
            <p className="text-sm font-semibold text-slate-950">Navigation</p>
            {footerLinks.map((link) => (
              <Link
                className="cursor-pointer text-sm text-muted-foreground transition-colors hover:text-primary"
                key={link.to}
                to={link.to}
              >
                {link.label}
              </Link>
            ))}
          </nav>

          <div className="grid gap-3">
            <p className="text-sm font-semibold text-slate-950">Social</p>
            {socialLinks.map((link) => (
              <span className="text-sm text-muted-foreground" key={link}>
                {link}
              </span>
            ))}
          </div>
        </div>
      </div>
    </footer>
  )
}
