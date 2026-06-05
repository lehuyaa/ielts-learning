package models

type UserRole string

const (
	UserRoleUser  UserRole = "USER"
	UserRoleAdmin UserRole = "ADMIN"
)

type LessonStatus string

const (
	LessonStatusLocked     LessonStatus = "LOCKED"
	LessonStatusUnlocked   LessonStatus = "UNLOCKED"
	LessonStatusInProgress LessonStatus = "IN_PROGRESS"
	LessonStatusCompleted  LessonStatus = "COMPLETED"
)

type VocabularyStatus string

const (
	VocabularyStatusNew      VocabularyStatus = "NEW"
	VocabularyStatusLearning VocabularyStatus = "LEARNING"
	VocabularyStatusReview   VocabularyStatus = "REVIEW"
	VocabularyStatusMastered VocabularyStatus = "MASTERED"
)

type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "BEGINNER"
	DifficultyIntermediate DifficultyLevel = "INTERMEDIATE"
	DifficultyAdvanced     DifficultyLevel = "ADVANCED"
)

type QuizQuestionType string

const (
	QuizQuestionMeaningChoice  QuizQuestionType = "MEANING_CHOICE"
	QuizQuestionWordChoice     QuizQuestionType = "WORD_CHOICE"
	QuizQuestionUsageChoice    QuizQuestionType = "USAGE_CHOICE"
	QuizQuestionSentenceChoice QuizQuestionType = "SENTENCE_CHOICE"
)

type FlashcardRating string

const (
	FlashcardRatingAgain FlashcardRating = "AGAIN"
	FlashcardRatingHard  FlashcardRating = "HARD"
	FlashcardRatingGood  FlashcardRating = "GOOD"
	FlashcardRatingEasy  FlashcardRating = "EASY"
)

type QuizSessionStatus string

const (
	QuizSessionInProgress QuizSessionStatus = "IN_PROGRESS"
	QuizSessionCompleted  QuizSessionStatus = "COMPLETED"
	QuizSessionAbandoned  QuizSessionStatus = "ABANDONED"
)

type NotificationType string

const (
	NotificationAchievement NotificationType = "ACHIEVEMENT"
	NotificationReviewDue   NotificationType = "REVIEW_DUE"
	NotificationStreak      NotificationType = "STREAK"
	NotificationSystem      NotificationType = "SYSTEM"
)

type SubscriptionStatus string

const (
	SubscriptionActive   SubscriptionStatus = "ACTIVE"
	SubscriptionTrialing SubscriptionStatus = "TRIALING"
	SubscriptionPastDue  SubscriptionStatus = "PAST_DUE"
	SubscriptionCanceled SubscriptionStatus = "CANCELED"
)
