package models

import "time"

type UserVocabularyProgress struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	UserID         uint             `gorm:"uniqueIndex:idx_user_vocab;not null" json:"userId"`
	VocabularyID   uint             `gorm:"uniqueIndex:idx_user_vocab;not null" json:"vocabularyId"`
	Status         VocabularyStatus `gorm:"type:varchar(30);default:'NEW';not null" json:"status"`
	EaseFactor     float64          `gorm:"type:decimal(5,2);default:2.5" json:"easeFactor"`
	IntervalDays   int              `gorm:"default:0" json:"intervalDays"`
	ReviewCount    int              `gorm:"default:0" json:"reviewCount"`
	CorrectCount   int              `gorm:"default:0" json:"correctCount"`
	WrongCount     int              `gorm:"default:0" json:"wrongCount"`
	LastRating     *FlashcardRating `gorm:"type:varchar(20)" json:"lastRating"`
	FirstLearnedAt *time.Time       `json:"firstLearnedAt"`
	LearnedAt      *time.Time       `json:"learnedAt"`
	MasteredAt     *time.Time       `json:"masteredAt"`
	LastReviewedAt *time.Time       `json:"lastReviewedAt"`
	NextReviewAt   *time.Time       `gorm:"index" json:"nextReviewAt"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`

	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Vocabulary Vocabulary `gorm:"foreignKey:VocabularyID" json:"vocabulary,omitempty"`
}

type UserLessonProgress struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	UserID        uint         `gorm:"uniqueIndex:idx_user_lesson;not null" json:"userId"`
	LessonID      uint         `gorm:"uniqueIndex:idx_user_lesson;not null" json:"lessonId"`
	Status        LessonStatus `gorm:"type:varchar(30);default:'UNLOCKED';not null" json:"status"`
	Score         *int         `json:"score"`
	BestScore     *int         `json:"bestScore"`
	BestXP        int          `gorm:"default:0" json:"bestXp"`
	WordsLearned  int          `gorm:"default:0" json:"wordsLearned"`
	TotalWords    int          `gorm:"default:0" json:"totalWords"`
	CompletedAt   *time.Time   `json:"completedAt"`
	StartedAt     *time.Time   `json:"startedAt"`
	LastStudiedAt *time.Time   `json:"lastStudiedAt"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`

	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Lesson Lesson `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
}

type DailyActivity struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `gorm:"uniqueIndex:idx_user_daily_activity;not null" json:"userId"`
	Date               time.Time `gorm:"type:date;uniqueIndex:idx_user_daily_activity;not null" json:"date"`
	WordsLearned       int       `gorm:"default:0" json:"wordsLearned"`
	WordsReviewed      int       `gorm:"default:0" json:"wordsReviewed"`
	QuizzesTaken       int       `gorm:"default:0" json:"quizzesTaken"`
	LessonsDone        int       `gorm:"default:0" json:"lessonsDone"`
	XPEarned           int       `gorm:"default:0" json:"xpEarned"`
	AccuracyPercent    *float64  `gorm:"type:decimal(5,2)" json:"accuracyPercent"`
	ActiveMinutes      int       `gorm:"default:0" json:"activeMinutes"`
	ChallengeProgress  int       `gorm:"default:0" json:"challengeProgress"`
	ChallengeCompleted bool      `gorm:"default:false" json:"challengeCompleted"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
