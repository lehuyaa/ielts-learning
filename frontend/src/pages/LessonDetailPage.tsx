import { useParams } from 'react-router-dom'

import { LessonHeader } from '@/features/lesson/LessonHeader'
import { LessonSidebar } from '@/features/lesson/LessonSidebar'
import { VocabularyPreviewList } from '@/features/lesson/VocabularyPreviewList'
import { mockLesson } from '@/features/lesson/mockLesson'

export function LessonDetailPage() {
  const { lessonId = mockLesson.id } = useParams()
  const learnedCount = mockLesson.vocabulary.filter((item) => item.learned).length
  const progressPercentage =
    mockLesson.vocabulary.length > 0
      ? (learnedCount / mockLesson.vocabulary.length) * 100
      : 0

  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <LessonHeader
        learnedCount={learnedCount}
        lesson={mockLesson}
        lessonId={lessonId}
        progressPercentage={progressPercentage}
      />

      <main className="mx-auto grid max-w-[1500px] gap-8 px-4 py-8 lg:grid-cols-[1fr_360px] lg:px-8">
        <section className="min-w-0">
          <div className="mb-6 rounded-2xl border border-[#e6e6f3] bg-white p-6 shadow-sm">
            <p className="text-sm font-bold uppercase tracking-normal text-[#676982]">
              {mockLesson.topic} - Band {mockLesson.bandRange}
            </p>
            <h2 className="mt-3 text-3xl font-bold tracking-normal text-[#10111f]">
              {mockLesson.title}
            </h2>
            <p className="mt-3 max-w-3xl text-lg font-medium text-[#676982]">
              {mockLesson.description}
            </p>
          </div>
          <VocabularyPreviewList vocabulary={mockLesson.vocabulary} />
        </section>

        <LessonSidebar
          learnedCount={learnedCount}
          lesson={mockLesson}
          lessonId={lessonId}
          progressPercentage={progressPercentage}
        />
      </main>
    </div>
  )
}
