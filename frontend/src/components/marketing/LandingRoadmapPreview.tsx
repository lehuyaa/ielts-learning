import { ArrowDown, ArrowRight } from 'lucide-react'
import { Link } from 'react-router-dom'

import { roadmapPreview } from '@/components/marketing/marketingContent'
import { SectionHeader } from '@/components/marketing/SectionHeader'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

export function LandingRoadmapPreview() {
  return (
    <section className="py-16">
      <div className="space-y-12 rounded-2xl border border-border bg-card p-6 shadow-sm md:p-8">
        <SectionHeader
          eyebrow="Roadmap Preview"
          title="A visible path from Band 5.0 to Band 8.0."
          description="Each band contains focused lessons, topic vocabulary, quizzes, and review checkpoints."
        />

        <div className="grid gap-6 lg:grid-cols-4">
          {roadmapPreview.map((step, index) => (
            <div className="relative" key={step.band}>
              <article className="h-full rounded-2xl border border-border bg-background p-6">
                <div className="flex items-center justify-between gap-4">
                  <Badge>{step.band}</Badge>
                  <div className="grid size-10 place-items-center rounded-xl bg-accent text-primary">
                    <step.icon className="size-5" aria-hidden="true" />
                  </div>
                </div>
                <h3 className="mt-6 text-xl font-semibold tracking-normal">
                  {step.title}
                </h3>
                <p className="mt-3 text-sm font-medium text-muted-foreground">
                  {step.lessons}
                </p>
              </article>

              {index < roadmapPreview.length - 1 ? (
                <>
                  <ArrowRight
                    className="absolute -right-5 top-1/2 hidden size-4 text-muted-foreground lg:block"
                    aria-hidden="true"
                  />
                  <ArrowDown
                    className="mx-auto mt-3 size-4 text-muted-foreground lg:hidden"
                    aria-hidden="true"
                  />
                </>
              ) : null}
            </div>
          ))}
        </div>

        <div className="flex justify-center">
          <Button asChild variant="outline">
            <Link to="/roadmap">View Roadmap</Link>
          </Button>
        </div>
      </div>
    </section>
  )
}
