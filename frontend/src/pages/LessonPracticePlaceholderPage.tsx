import { Link, useParams } from 'react-router-dom'

import { Button } from '@/components/ui/button'

type LessonPracticePlaceholderPageProps = {
  mode: 'flashcards' | 'quiz'
}

const modeTitle = {
  flashcards: 'Flashcards',
  quiz: 'Quiz',
}

export function LessonPracticePlaceholderPage({
  mode,
}: LessonPracticePlaceholderPageProps) {
  const { lessonId = '1' } = useParams()

  return (
    <div className="grid min-h-[60vh] place-items-center px-4">
      <section className="w-full max-w-md rounded-2xl border border-[#e6e6f3] bg-white p-8 text-center shadow-sm">
        <p className="text-sm font-bold uppercase tracking-normal text-[#676982]">
          Lesson {lessonId}
        </p>
        <h1 className="mt-3 text-3xl font-bold tracking-normal text-[#10111f]">
          {modeTitle[mode]} coming soon
        </h1>
        <p className="mt-3 text-base font-medium text-[#676982]">
          This route is ready for the future {modeTitle[mode].toLowerCase()}{' '}
          task.
        </p>
        <Button asChild className="mt-8 rounded-full">
          <Link to={`/lessons/${lessonId}`}>Back to lesson</Link>
        </Button>
      </section>
    </div>
  )
}
