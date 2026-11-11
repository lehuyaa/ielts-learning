import { createBrowserRouter } from 'react-router-dom'

import { AuthLayout } from '@/components/layout/AuthLayout'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import { MarketingLayout } from '@/components/layout/MarketingLayout'
import { AdminRoute } from '@/features/auth/AdminRoute'
import { ProtectedRoute } from '@/features/auth/ProtectedRoute'
import { AdminContentImportPage } from '@/pages/AdminContentImportPage'
import { AdminVocabularyImportPage } from '@/pages/AdminVocabularyImportPage'
import { AIConversationCreatePage } from '@/pages/AIConversationCreatePage'
import { AIConversationPage } from '@/pages/AIConversationPage'
import { AIConversationSessionPage } from '@/pages/AIConversationSessionPage'
import { SpeakingSessionPage } from '@/pages/SpeakingSessionPage'
import { MistakesPage } from '@/pages/MistakesPage'
import { MistakeReviewPage } from '@/pages/MistakeReviewPage'
import { PracticePage } from '@/pages/PracticePage'
import { DashboardPage } from '@/pages/DashboardPage'
import { FlashcardPage } from '@/pages/FlashcardPage'
import { LandingPage } from '@/pages/LandingPage'
import { LessonDetailPage } from '@/pages/LessonDetailPage'
import { LoginPage } from '@/pages/LoginPage'
import { ProfilePage } from '@/pages/ProfilePage'
import { QuizPage } from '@/pages/QuizPage'
import { RegisterPage } from '@/pages/RegisterPage'
import { GrammarDetailPage } from '@/pages/GrammarDetailPage'
import { GrammarPage } from '@/pages/GrammarPage'
import { PronunciationDetailPage } from '@/pages/PronunciationDetailPage'
import { AdminContentExtractPage } from '@/pages/AdminContentExtractPage'
import { AdminPassageImportPage } from '@/pages/AdminPassageImportPage'
import { PronunciationPage } from '@/pages/PronunciationPage'
import { LessonPlayerPage } from '@/pages/LessonPlayerPage'
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
            element: <QuizPage />,
          },
          {
            path: '/lessons/:lessonId/learn',
            element: <LessonPlayerPage />,
          },
          {
            path: '/grammar',
            element: <GrammarPage />,
          },
          {
            path: '/grammar/:slug',
            element: <GrammarDetailPage />,
          },
          {
            path: '/pronunciation',
            element: <PronunciationPage />,
          },
          {
            path: '/pronunciation/:symbol',
            element: <PronunciationDetailPage />,
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
          {
            path: '/ai-conversation',
            element: <AIConversationPage />,
          },
          {
            path: '/ai-conversation/create',
            element: <AIConversationCreatePage />,
          },
          {
            path: '/mistakes',
            element: <MistakesPage />,
          },
          {
            path: '/mistakes/review',
            element: <MistakeReviewPage />,
          },
          {
            path: '/practice/:setId',
            element: <PracticePage />,
          },
          {
            path: '/ai-conversation/sessions/:sessionId',
            element: <SpeakingSessionPage />,
          },
          {
            path: '/ai-conversation/:scenarioSlug',
            element: <AIConversationSessionPage />,
          },
          {
            element: <AdminRoute />,
            children: [
              {
                path: '/admin/vocabularies/import',
                element: <AdminVocabularyImportPage />,
              },
              {
                path: '/admin/content/import',
                element: <AdminContentImportPage />,
              },
              {
                path: '/admin/content/extract',
                element: <AdminContentExtractPage />,
              },
              {
                path: '/admin/passages/import',
                element: <AdminPassageImportPage />,
              },
            ],
          },
        ],
      },
    ],
  },
])
