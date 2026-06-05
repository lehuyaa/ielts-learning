package models

import "time"

type QuizQuestion struct {
	ID               uint             `gorm:"primaryKey" json:"id"`
	LessonID         uint             `gorm:"index;not null" json:"lessonId"`
	TopicID          *uint            `gorm:"index" json:"topicId"`
	VocabularyID     *uint            `gorm:"index" json:"vocabularyId"`
	Type             QuizQuestionType `gorm:"type:varchar(50);not null" json:"type"`
	Question         string           `gorm:"type:text;not null" json:"question"`
	Explanation      string           `gorm:"type:text" json:"explanation"`
	Points           int              `gorm:"default:20" json:"points"`
	TimeLimitSeconds *int             `json:"timeLimitSeconds"`
	OrderIndex       int              `gorm:"default:0" json:"orderIndex"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`

	Lesson     Lesson       `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Topic      *Topic       `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	Vocabulary *Vocabulary  `gorm:"foreignKey:VocabularyID" json:"vocabulary,omitempty"`
	Options    []QuizOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
}

type QuizOption struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionID uint      `gorm:"index;not null" json:"questionId"`
	Label      string    `gorm:"size:5" json:"label"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsCorrect  bool      `gorm:"default:false" json:"-"`
	OrderIndex int       `gorm:"default:0" json:"orderIndex"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	Question QuizQuestion `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

type QuizSession struct {
	ID                   uint              `gorm:"primaryKey" json:"id"`
	UserID               uint              `gorm:"index:idx_user_quiz_session_status;not null" json:"userId"`
	LessonID             uint              `gorm:"index;not null" json:"lessonId"`
	Status               QuizSessionStatus `gorm:"type:varchar(30);default:'IN_PROGRESS';index:idx_user_quiz_session_status" json:"status"`
	CurrentQuestionIndex int               `gorm:"default:0" json:"currentQuestionIndex"`
	TotalQuestions       int               `gorm:"default:0" json:"totalQuestions"`
	Points               int               `gorm:"default:0" json:"points"`
	CorrectAnswers       int               `gorm:"default:0" json:"correctAnswers"`
	WrongAnswers         int               `gorm:"default:0" json:"wrongAnswers"`
	StartedAt            time.Time         `gorm:"not null" json:"startedAt"`
	FinishedAt           *time.Time        `json:"finishedAt"`
	ExpiresAt            *time.Time        `json:"expiresAt"`
	CreatedAt            time.Time         `json:"createdAt"`
	UpdatedAt            time.Time         `json:"updatedAt"`

	User    User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Lesson  Lesson              `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Answers []QuizSessionAnswer `gorm:"foreignKey:QuizSessionID" json:"answers,omitempty"`
}

type QuizSessionAnswer struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	QuizSessionID    uint      `gorm:"uniqueIndex:idx_session_question;not null" json:"quizSessionId"`
	QuestionID       uint      `gorm:"uniqueIndex:idx_session_question;not null" json:"questionId"`
	SelectedOptionID *uint     `json:"selectedOptionId"`
	CorrectOptionID  *uint     `json:"correctOptionId"`
	IsCorrect        bool      `gorm:"not null" json:"isCorrect"`
	PointsAwarded    int       `gorm:"default:0" json:"pointsAwarded"`
	AnsweredAt       time.Time `gorm:"not null" json:"answeredAt"`
	TimeSpentSeconds *int      `json:"timeSpentSeconds"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	QuizSession    QuizSession  `gorm:"foreignKey:QuizSessionID" json:"quizSession,omitempty"`
	Question       QuizQuestion `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	SelectedOption *QuizOption  `gorm:"foreignKey:SelectedOptionID" json:"selectedOption,omitempty"`
	CorrectOption  *QuizOption  `gorm:"foreignKey:CorrectOptionID" json:"correctOption,omitempty"`
}

type UserQuizAttempt struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"index;not null" json:"userId"`
	LessonID        uint       `gorm:"index;not null" json:"lessonId"`
	QuizSessionID   *uint      `gorm:"index" json:"quizSessionId"`
	Score           int        `gorm:"not null" json:"score"`
	Points          int        `gorm:"default:0" json:"points"`
	TotalQuestions  int        `gorm:"not null" json:"totalQuestions"`
	CorrectAnswers  int        `gorm:"not null" json:"correctAnswers"`
	DurationSeconds *int       `json:"durationSeconds"`
	XPEarned        int        `gorm:"default:0" json:"xpEarned"`
	Passed          bool       `gorm:"default:false" json:"passed"`
	StartedAt       *time.Time `json:"startedAt"`
	FinishedAt      *time.Time `json:"finishedAt"`
	CreatedAt       time.Time  `json:"createdAt"`

	User        User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Lesson      Lesson       `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	QuizSession *QuizSession `gorm:"foreignKey:QuizSessionID" json:"quizSession,omitempty"`
}
