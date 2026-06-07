import { AlertCircle } from 'lucide-react'
import { useNavigate, useParams } from 'react-router-dom'

import { APIError } from '@/api/api'
import { LessonHeader } from '@/features/lesson/LessonHeader'
import { LessonSidebar } from '@/features/lesson/LessonSidebar'
import { VocabularyPreviewList } from '@/features/lesson/VocabularyPreviewList'
import { useLessonDetail } from '@/features/lesson/hooks/useLessonDetail'
import { useStartLesson } from '@/features/lesson/hooks/useStartLesson'
import type {
  LessonDetailResponse,
  LessonDetailViewModel,
  LessonDifficulty,
} from '@/types/lesson'

const loadingLesson: LessonDetailViewModel = {
  id: 'loading',
  title: 'Lesson',
  topic: 'Topic',
  topicId: 0,
  bandRange: 'Loading',
  description: 'Loading lesson details and vocabulary.',
  estimatedMinutes: 0,
  requiredScore: 0,
  xpReward: 0,
  score: null,
  bestScore: null,
  status: 'UNLOCKED',
  lockedReason: null,
  progressPercentage: 0,
  vocabulary: [],
}

export function LessonDetailPage() {
  const { lessonId } = useParams()
  const navigate = useNavigate()
  const lessonQuery = useLessonDetail(lessonId)
  const lesson = lessonQuery.data
    ? mapLessonDetailToViewModel(lessonQuery.data)
    : loadingLesson
  const startLessonMutation = useStartLesson({
    lessonId,
    topicId: lesson.topicId || undefined,
  })

  const learnedCount = lesson.vocabulary.filter((item) => item.learned).length
  const progressPercentage = lesson.progressPercentage
  const errorMessage = getLessonErrorMessage(lessonQuery.error)
  const startErrorMessage = getLessonErrorMessage(startLessonMutation.error)
  const isEmpty =
    !lessonQuery.isLoading &&
    !errorMessage &&
    Boolean(lessonQuery.data) &&
    lesson.vocabulary.length === 0
  const isLocked = lesson.status === 'LOCKED'
  const isStartDisabled =
    !lessonId || isLocked || startLessonMutation.isPending || lessonQuery.isLoading

  async function handleStart(target: 'flashcards' | 'quiz') {
    if (!lessonId || isLocked) {
      return
    }

    try {
      if (lesson.status !== 'COMPLETED') {
        await startLessonMutation.mutateAsync()
      }

      void navigate(`/lessons/${lessonId}/${target}`)
    } catch {
      // React Query exposes the mutation error for the page-level message.
    }
  }

  return (
    <div className="min-h-screen bg-[#f8f8ff] text-[#10111f]">
      <LessonHeader
        learnedCount={learnedCount}
        isStartDisabled={isStartDisabled}
        isStarting={startLessonMutation.isPending}
        lesson={lesson}
        onStartFlashcards={() => {
          void handleStart('flashcards')
        }}
        progressPercentage={progressPercentage}
      />

      <main className="mx-auto grid max-w-[1320px] gap-6 px-4 py-6 lg:grid-cols-[1fr_320px] lg:px-6">
        {lessonQuery.isLoading ? (
          <LessonStateMessage
            description="Preparing vocabulary, progress, and lesson metadata."
            title="Loading lesson"
          />
        ) : null}

        {errorMessage ? (
          <LessonStateMessage
            description={errorMessage}
            title="Lesson unavailable"
            tone="error"
          />
        ) : null}

        {isEmpty ? (
          <>
            <LessonMainContent lesson={lesson} />
            <LessonStateMessage
              description="This lesson exists, but no vocabulary has been added yet."
              title="No vocabulary yet"
            />
          </>
        ) : null}

        {!lessonQuery.isLoading && !errorMessage && !isEmpty ? (
          <>
            <LessonMainContent lesson={lesson} />

            <LessonSidebar
              learnedCount={learnedCount}
              isStartDisabled={isStartDisabled}
              isStarting={startLessonMutation.isPending}
              lesson={lesson}
              onStartFlashcards={() => {
                void handleStart('flashcards')
              }}
              onStartQuiz={() => {
                void handleStart('quiz')
              }}
              progressPercentage={progressPercentage}
            />
          </>
        ) : null}

        {startErrorMessage ? (
          <LessonStateMessage
            description={startErrorMessage}
            title="Could not start lesson"
            tone="error"
          />
        ) : null}
      </main>
    </div>
  )
}

type LessonMainContentProps = {
  lesson: LessonDetailViewModel
}

function LessonMainContent({ lesson }: LessonMainContentProps) {
  return (
    <section className="min-w-0">
      <div className="mb-5 rounded-2xl border border-[#e6e6f3] bg-white p-5 shadow-sm">
        <p className="text-sm font-bold uppercase tracking-normal text-[#676982]">
          {lesson.topic} - Band {lesson.bandRange}
        </p>
        <h2 className="mt-3 text-2xl font-bold tracking-normal text-[#10111f]">
          {lesson.title}
        </h2>
        <p className="mt-3 max-w-3xl text-base font-medium text-[#676982]">
          {lesson.description}
        </p>
      </div>
      <VocabularyPreviewList vocabulary={lesson.vocabulary} />
    </section>
  )
}

function mapLessonDetailToViewModel(
  response: LessonDetailResponse,
): LessonDetailViewModel {
  return {
    id: String(response.lesson.id),
    title: response.lesson.title,
    topic: response.topic.title,
    topicId: response.topic.id,
    bandRange: response.bandLevel.bandScore.toFixed(1),
    description: response.lesson.description,
    estimatedMinutes: response.lesson.estimatedMinutes,
    requiredScore: response.lesson.requiredScore,
    xpReward: response.lesson.xpReward,
    score: response.progress.score,
    bestScore: response.progress.bestScore,
    status: response.lesson.status,
    lockedReason: response.lesson.lockedReason,
    progressPercentage: response.lesson.progressPercentage,
    vocabulary: response.vocabularies.map((vocabulary) => ({
      id: String(vocabulary.id),
      word: vocabulary.word,
      ipa: vocabulary.ipa,
      partOfSpeech: vocabulary.partOfSpeech,
      shortDefinition:
        vocabulary.shortDefinition ||
        vocabulary.meaningEn ||
        vocabulary.meaningVi,
      definition:
        vocabulary.meaningEn ||
        vocabulary.shortDefinition ||
        vocabulary.meaningVi,
      example: vocabulary.exampleSentence,
      difficulty: mapDifficulty(vocabulary.difficulty),
      band: vocabulary.targetBand
        ? `Band ${vocabulary.targetBand.toFixed(1)}`
        : response.bandLevel.title,
      learned: vocabulary.learned,
    })),
  }
}

function mapDifficulty(
  difficulty: LessonDifficulty,
): LessonDetailViewModel['vocabulary'][number]['difficulty'] {
  switch (difficulty) {
    case 'BEGINNER':
      return 'Beginner'
    case 'ADVANCED':
      return 'Advanced'
    case 'INTERMEDIATE':
    default:
      return 'Intermediate'
  }
}

function getLessonErrorMessage(error: Error | null) {
  if (!error) {
    return null
  }

  if (error instanceof APIError) {
    if (error.status === 404) {
      return 'This lesson could not be found.'
    }
    if (error.status === 403 || error.code === 'LESSON_LOCKED') {
      return 'Complete previous lessons to unlock this lesson.'
    }

    return error.message
  }

  return 'Unable to load this lesson right now.'
}

type LessonStateMessageProps = {
  title: string
  description: string
  tone?: 'default' | 'error'
}

function LessonStateMessage({
  title,
  description,
  tone = 'default',
}: LessonStateMessageProps) {
  return (
    <section className="rounded-2xl border border-[#e6e6f3] bg-white p-6 text-center shadow-sm lg:col-span-2">
      {tone === 'error' ? (
        <AlertCircle
          className="mx-auto size-8 text-destructive"
          aria-hidden="true"
        />
      ) : null}
      <h2 className="text-xl font-bold tracking-normal text-[#10111f]">
        {title}
      </h2>
      <p className="mt-3 text-base font-medium text-[#676982]">{description}</p>
    </section>
  )
}
