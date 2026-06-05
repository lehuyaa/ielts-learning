import { ArrowRight, CheckCircle2, Sparkles } from 'lucide-react'
import { Link } from 'react-router-dom'

import heroImage from '@/assets/hero.png'
import { heroBands } from '@/components/marketing/marketingContent'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

export function LandingHero() {
  return (
    <section className="grid items-center gap-12 py-16 lg:grid-cols-[1.05fr_0.95fr]">
      <div className="max-w-2xl">
        <Badge className="mb-6" variant="outline">
          <Sparkles className="size-3.5" aria-hidden="true" />
          IELTS vocabulary learning platform
        </Badge>

        <h1 className="text-5xl font-bold leading-tight tracking-normal text-slate-950 md:text-6xl">
          Build IELTS vocabulary with a clear band roadmap.
        </h1>

        <p className="mt-6 max-w-xl text-lg leading-8 text-slate-600">
          Follow structured lessons from your current level to your target band,
          review words with spaced repetition, and test yourself with
          IELTS-focused quizzes.
        </p>

        <div className="mt-8 flex flex-col gap-3 sm:flex-row">
          <Button asChild size="lg">
            <Link to="/register">
              Start Learning Free
              <ArrowRight aria-hidden="true" />
            </Link>
          </Button>
          <Button asChild size="lg" variant="outline">
            <Link to="/roadmap">View Roadmap</Link>
          </Button>
        </div>

        <div className="mt-12 flex flex-wrap items-center gap-3">
          {heroBands.map((band, index) => (
            <div className="flex items-center gap-3" key={band}>
              <span className="rounded-full border border-border bg-card px-4 py-2 text-sm font-semibold text-slate-700">
                Band {band}
              </span>
              {index < heroBands.length - 1 ? (
                <ArrowRight
                  className="size-4 text-muted-foreground"
                  aria-hidden="true"
                />
              ) : null}
            </div>
          ))}
        </div>
      </div>

      <div className="rounded-2xl border border-border bg-card p-6 shadow-xl shadow-slate-200/70 md:p-8">
        <div className="grid gap-8">
          <div className="grid place-items-center rounded-2xl bg-slate-950 p-8">
            <img
              alt="Layered roadmap illustration"
              className="h-auto w-full max-w-xs"
              src={heroImage}
            />
          </div>

          <div className="grid gap-4">
            {[
              'Choose your target band',
              'Learn topic-based vocabulary',
              'Review until words stick',
            ].map((item) => (
              <div
                className="flex items-center gap-3 rounded-xl border border-border bg-muted px-4 py-3"
                key={item}
              >
                <CheckCircle2 className="size-5 text-success" aria-hidden="true" />
                <span className="font-medium text-slate-700">{item}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
