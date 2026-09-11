import { landingFeatures } from '@/components/marketing/marketingContent'
import { SectionHeader } from '@/components/marketing/SectionHeader'
import { CardDescription, CardTitle } from '@/components/ui/card'

export function LandingFeatures() {
  return (
    <section className="py-16">
      <div className="space-y-12">
        <SectionHeader
          eyebrow="Features"
          title="Everything learners need to stop memorizing random word lists."
          description="LexPath connects roadmap structure, practice, reviews, and progress into one focused IELTS vocabulary workflow."
        />

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {landingFeatures.map((feature) => (
            <article
              className="rounded-2xl border border-border bg-card p-6 shadow-sm md:p-8"
              key={feature.title}
            >
              <div className="grid size-12 place-items-center rounded-xl bg-accent text-primary">
                <feature.icon className="size-5" aria-hidden="true" />
              </div>
              <CardTitle
                as="h3"
                className="mt-6 text-xl leading-normal text-slate-950"
              >
                {feature.title}
              </CardTitle>
              <CardDescription className="mt-3 leading-6">
                {feature.description}
              </CardDescription>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}
