import { CheckCircle2, Lock } from 'lucide-react'

import { RoadmapTopicCard } from '@/features/roadmap/RoadmapTopicCard'
import type { RoadmapBand, RoadmapLessonStatus } from '@/features/roadmap/types'
import { cn } from '@/lib/utils'

type RoadmapBandSectionProps = {
  band: RoadmapBand
}

export function RoadmapBandSection({ band }: RoadmapBandSectionProps) {
  const isLocked = band.status === 'locked'
  const topics = getScreenshotTopicOrder(band)

  if (isLocked) {
    return (
      <section className="relative py-10">
        <div className="mx-auto max-w-sm rounded-2xl border border-[#e3e4f8] bg-[#f8f8ff] p-5">
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-4">
              <Lock className="size-6 text-[#6d7088]" aria-hidden="true" />
              <div>
                <h2 className="text-xl font-bold tracking-normal text-[#676982]">
                  Band {band.band}
                </h2>
                <p className="mt-1 text-base font-medium text-[#676982]">
                  Complete previous band
                </p>
              </div>
            </div>
            <span className="rounded-full bg-[#e9e8ff] px-3 py-1 text-base font-bold text-[#3f35e8]">
              Locked
            </span>
          </div>
        </div>
      </section>
    )
  }

  return (
    <section className="relative py-10">
      <BandNode
        band={band.band}
        status={bandStatusLabel(band.status)}
        statusVariant={band.status}
        topicCount={band.topicCount}
      />

      <div className="mx-auto mt-6 grid max-w-[720px] gap-5 lg:grid-cols-2">
        {topics.map((topic) => (
          <RoadmapTopicCard key={topic.id} topic={topic} />
        ))}
      </div>
    </section>
  )
}

function getScreenshotTopicOrder(band: RoadmapBand) {
  if (band.band !== '6.0') {
    return band.topics
  }

  const order = ['technology', 'environment', 'government', 'economy']

  return [...band.topics].sort(
    (first, second) =>
      topicOrderIndex(order, first.id) - topicOrderIndex(order, second.id),
  )
}

function topicOrderIndex(order: string[], topicId: string) {
  const index = order.indexOf(topicId)
  return index === -1 ? order.length : index
}

function bandStatusLabel(status: RoadmapLessonStatus) {
  switch (status) {
    case 'completed':
      return 'Complete'
    case 'in-progress':
      return 'In Progress'
    case 'unlocked':
      return 'Unlocked'
    case 'locked':
    default:
      return 'Locked'
  }
}

type BandNodeProps = {
  band: string
  status: string
  statusVariant: RoadmapLessonStatus
  topicCount: number
}

function BandNode({ band, status, statusVariant, topicCount }: BandNodeProps) {
  return (
    <div
      className={cn(
        'mx-auto max-w-sm rounded-2xl border p-5',
        statusVariant === 'completed' && 'border-[#dce2f3] bg-white',
        statusVariant !== 'completed' && 'border-[#add2ff] bg-[#eef7ff]',
      )}
    >
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          {statusVariant === 'completed' ? (
            <CheckCircle2 className="size-6 text-success" aria-hidden="true" />
          ) : (
            <span className="size-6 rounded-full bg-primary" aria-hidden="true" />
          )}
          <div>
            <h2 className="text-xl font-bold tracking-normal text-[#10111f]">
              Band {band}
            </h2>
            <p className="mt-1 text-base font-medium text-[#676982]">
              {topicCount} topics
            </p>
          </div>
        </div>
        <span
          className={cn(
            'rounded-full px-3 py-1 text-base font-bold',
            statusVariant === 'completed' && 'bg-[#eef2f7] text-[#40516b]',
            statusVariant !== 'completed' && 'bg-[#dbeafe] text-[#1d4ed8]',
          )}
        >
          {status}
        </span>
      </div>
    </div>
  )
}
