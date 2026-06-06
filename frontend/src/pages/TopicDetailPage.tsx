import { useParams } from 'react-router-dom'

import { LessonList } from '@/features/topic/LessonList'
import { TopicHeader } from '@/features/topic/TopicHeader'
import { TopicHero } from '@/features/topic/TopicHero'
import { mockTopic } from '@/features/topic/mockTopic'

export function TopicDetailPage() {
  const { topicId = mockTopic.id } = useParams()
  const topic = { ...mockTopic, id: topicId }

  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <TopicHeader topic={topic} />

      <main className="mx-auto grid max-w-[1280px] gap-8 px-4 py-11 md:px-8">
        <TopicHero topic={topic} />
        <LessonList topic={topic} />
      </main>
    </div>
  )
}
