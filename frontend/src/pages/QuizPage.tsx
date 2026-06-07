import { useParams } from 'react-router-dom'

import { mockQuiz } from '@/features/quiz/mockQuiz'
import { QuizSession } from '@/features/quiz/QuizSession'

export function QuizPage() {
  const { lessonId = String(mockQuiz.lessonId) } = useParams()

  return <QuizSession lessonId={lessonId} quiz={mockQuiz} />
}
