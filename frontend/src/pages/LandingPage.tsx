import { LandingFeatures } from '@/components/marketing/LandingFeatures'
import { LandingFooter } from '@/components/marketing/LandingFooter'
import { LandingHero } from '@/components/marketing/LandingHero'
import { LandingPricing } from '@/components/marketing/LandingPricing'
import { LandingRoadmapPreview } from '@/components/marketing/LandingRoadmapPreview'
import { LandingStats } from '@/components/marketing/LandingStats'
import { LandingTestimonials } from '@/components/marketing/LandingTestimonials'

export function LandingPage() {
  return (
    <>
      <LandingHero />
      <LandingStats />
      <LandingFeatures />
      <LandingRoadmapPreview />
      <LandingTestimonials />
      <LandingPricing />
      <LandingFooter />
    </>
  )
}
