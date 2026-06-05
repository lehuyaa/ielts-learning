import { CheckCircle2 } from 'lucide-react'
import { Link } from 'react-router-dom'

import { pricingPlans } from '@/components/marketing/marketingContent'
import { SectionHeader } from '@/components/marketing/SectionHeader'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

export function LandingPricing() {
  return (
    <section className="py-16">
      <div className="space-y-12">
        <SectionHeader
          eyebrow="Pricing"
          title="Start free while the full learning platform grows."
          description="Pricing is intentionally simple for the MVP. Premium features are marked as placeholders until future tasks define them."
        />

        <div className="mx-auto grid max-w-4xl gap-6 md:grid-cols-2">
          {pricingPlans.map((plan) => (
            <article
              className={cn(
                'rounded-2xl border bg-card p-6 shadow-sm md:p-8',
                plan.highlighted ? 'border-primary' : 'border-border',
              )}
              key={plan.name}
            >
              <div className="flex items-start justify-between gap-4">
                <div>
                  <h3 className="text-2xl font-bold tracking-normal">
                    {plan.name}
                  </h3>
                  <p className="mt-3 text-muted-foreground">{plan.description}</p>
                </div>
                {!plan.highlighted ? (
                  <Badge variant="warning">Placeholder</Badge>
                ) : null}
              </div>

              <p className="mt-8 text-4xl font-bold tracking-normal text-slate-950">
                {plan.price}
              </p>

              <ul className="mt-8 grid gap-3">
                {plan.features.map((feature) => (
                  <li className="flex items-center gap-3" key={feature}>
                    <CheckCircle2
                      className="size-5 text-success"
                      aria-hidden="true"
                    />
                    <span className="text-sm font-medium text-slate-700">
                      {feature}
                    </span>
                  </li>
                ))}
              </ul>

              <Button
                asChild={plan.highlighted}
                className="mt-8 w-full"
                disabled={!plan.highlighted}
                variant={plan.highlighted ? 'default' : 'outline'}
              >
                {plan.highlighted ? <Link to={plan.to}>{plan.cta}</Link> : plan.cta}
              </Button>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}
