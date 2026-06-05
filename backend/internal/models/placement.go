package models

import "time"

type PlacementTest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"size:255;not null" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	QuestionCount int       `gorm:"default:20" json:"questionCount"`
	IsActive      bool      `gorm:"default:true" json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	Questions []PlacementQuestion `gorm:"foreignKey:PlacementTestID" json:"questions,omitempty"`
}

type PlacementQuestion struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	PlacementTestID uint            `gorm:"index;not null" json:"placementTestId"`
	VocabularyID    *uint           `gorm:"index" json:"vocabularyId"`
	Question        string          `gorm:"type:text;not null" json:"question"`
	BandScore       *float64        `gorm:"type:decimal(3,1)" json:"bandScore"`
	Difficulty      DifficultyLevel `gorm:"type:varchar(30)" json:"difficulty"`
	OrderIndex      int             `gorm:"default:0" json:"orderIndex"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`

	PlacementTest PlacementTest     `gorm:"foreignKey:PlacementTestID" json:"placementTest,omitempty"`
	Vocabulary    *Vocabulary       `gorm:"foreignKey:VocabularyID" json:"vocabulary,omitempty"`
	Options       []PlacementOption `gorm:"foreignKey:PlacementQuestionID" json:"options,omitempty"`
}

type PlacementOption struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	PlacementQuestionID uint      `gorm:"index;not null" json:"placementQuestionId"`
	Label               string    `gorm:"size:5" json:"label"`
	Content             string    `gorm:"type:text;not null" json:"content"`
	IsCorrect           bool      `gorm:"default:false" json:"-"`
	OrderIndex          int       `gorm:"default:0" json:"orderIndex"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`

	PlacementQuestion PlacementQuestion `gorm:"foreignKey:PlacementQuestionID" json:"placementQuestion,omitempty"`
}

type PlacementAttempt struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          *uint      `gorm:"index" json:"userId"`
	PlacementTestID uint       `gorm:"index;not null" json:"placementTestId"`
	Score           int        `gorm:"default:0" json:"score"`
	TotalQuestions  int        `gorm:"default:0" json:"totalQuestions"`
	CorrectAnswers  int        `gorm:"default:0" json:"correctAnswers"`
	EstimatedBand   *float64   `gorm:"type:decimal(3,1)" json:"estimatedBand"`
	RecommendedBand *float64   `gorm:"type:decimal(3,1)" json:"recommendedBand"`
	StartedAt       time.Time  `gorm:"not null" json:"startedAt"`
	CompletedAt     *time.Time `json:"completedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`

	User          *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PlacementTest PlacementTest `gorm:"foreignKey:PlacementTestID" json:"placementTest,omitempty"`
}

type PlacementAttemptAnswer struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	PlacementAttemptID  uint      `gorm:"uniqueIndex:idx_placement_attempt_question;not null" json:"placementAttemptId"`
	PlacementQuestionID uint      `gorm:"uniqueIndex:idx_placement_attempt_question;not null" json:"placementQuestionId"`
	SelectedOptionID    *uint     `json:"selectedOptionId"`
	IsCorrect           bool      `gorm:"not null" json:"isCorrect"`
	AnsweredAt          time.Time `gorm:"not null" json:"answeredAt"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`

	PlacementAttempt  PlacementAttempt  `gorm:"foreignKey:PlacementAttemptID" json:"placementAttempt,omitempty"`
	PlacementQuestion PlacementQuestion `gorm:"foreignKey:PlacementQuestionID" json:"placementQuestion,omitempty"`
	SelectedOption    *PlacementOption  `gorm:"foreignKey:SelectedOptionID" json:"selectedOption,omitempty"`
}
