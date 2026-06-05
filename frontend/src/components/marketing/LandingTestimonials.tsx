import { Quote } from 'lucide-react'

import { testimonials } from '@/components/marketing/marketingContent'
import { SectionHeader } from '@/components/marketing/SectionHeader'

export function LandingTestimonials() {
  return (
    <section className="py-16">
      <div className="space-y-12">
        <SectionHeader
          eyebrow="Testimonials"
          title="Built for learners who want structure before test day."
          description="Placeholder stories show the kind of progress LexPath is designed to support."
        />

        <div className="grid gap-6 lg:grid-cols-3">
          {testimonials.map((testimonial) => (
            <article
              className="rounded-2xl border border-border bg-card p-6 shadow-sm md:p-8"
              key={testimonial.name}
            >
              <Quote className="size-6 text-primary" aria-hidden="true" />
              <p className="mt-6 leading-7 text-slate-700">
                “{testimonial.feedback}”
              </p>
              <div className="mt-8">
                <p className="font-semibold text-slate-950">{testimonial.name}</p>
                <p className="mt-1 text-sm text-muted-foreground">
                  {testimonial.targetBand}
                </p>
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}
