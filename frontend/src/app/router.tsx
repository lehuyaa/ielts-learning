import { createBrowserRouter } from 'react-router-dom'

import { AuthLayout } from '@/components/layout/AuthLayout'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import { MarketingLayout } from '@/components/layout/MarketingLayout'
import { ProtectedRoute } from '@/features/auth/ProtectedRoute'
import { DashboardPage } from '@/pages/DashboardPage'
import { FlashcardPage } from '@/pages/FlashcardPage'
import { LandingPage } from '@/pages/LandingPage'
import { LessonDetailPage } from '@/pages/LessonDetailPage'
import { LessonPracticePlaceholderPage } from '@/pages/LessonPracticePlaceholderPage'
import { LoginPage } from '@/pages/LoginPage'
import { ProfilePage } from '@/pages/ProfilePage'
import { RegisterPage } from '@/pages/RegisterPage'
import { ReviewsPage } from '@/pages/ReviewsPage'
import { RoadmapPage } from '@/pages/RoadmapPage'
import { TopicDetailPage } from '@/pages/TopicDetailPage'
import { VocabularyDetailPage } from '@/pages/VocabularyDetailPage'
import { VocabularyPage } from '@/pages/VocabularyPage'

export const router = createBrowserRouter([
  {
    element: <MarketingLayout />,
    children: [
      {
        path: '/',
        element: <LandingPage />,
      },
    ],
  },
  {
    element: <AuthLayout />,
    children: [
      {
        path: '/login',
        element: <LoginPage />,
      },
      {
        path: '/register',
        element: <RegisterPage />,
      },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <DashboardLayout />,
        children: [
          {
            path: '/dashboard',
            element: <DashboardPage />,
          },
          {
            path: '/roadmap',
            element: <RoadmapPage />,
          },
          {
            path: '/topics/:topicId',
            element: <TopicDetailPage />,
          },
          {
            path: '/lessons/:lessonId',
            element: <LessonDetailPage />,
          },
          {
            path: '/lessons/:lessonId/flashcards',
            element: <FlashcardPage />,
          },
          {
            path: '/lessons/:lessonId/quiz',
            element: <LessonPracticePlaceholderPage mode="quiz" />,
          },
          {
            path: '/reviews',
            element: <ReviewsPage />,
          },
          {
            path: '/vocabulary',
            element: <VocabularyPage />,
          },
          {
            path: '/vocabulary/:vocabularyId',
            element: <VocabularyDetailPage />,
          },
          {
            path: '/profile',
            element: <ProfilePage />,
          },
        ],
      },
    ],
  },
])
