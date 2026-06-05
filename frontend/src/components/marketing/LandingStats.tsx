import { landingStats } from '@/components/marketing/marketingContent'

export function LandingStats() {
  return (
    <section className="py-16">
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {landingStats.map((stat) => (
          <div
            className="rounded-2xl border border-border bg-card p-6 shadow-sm md:p-8"
            key={stat.label}
          >
            <p className="text-3xl font-bold tracking-normal text-primary">
              {stat.value}
            </p>
            <p className="mt-2 text-sm font-medium text-muted-foreground">
              {stat.label}
            </p>
          </div>
        ))}
      </div>
    </section>
  )
}
