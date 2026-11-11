package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ActivityType is one skill step inside a CORE lesson (docs/foundation-spec.md §4).
type ActivityType string

const (
	ActivityVocabulary    ActivityType = "VOCABULARY"
	ActivityPronunciation ActivityType = "PRONUNCIATION"
	ActivityGrammar       ActivityType = "GRAMMAR"
	ActivityListening     ActivityType = "LISTENING"
	ActivityReading       ActivityType = "READING"
	ActivityWriting       ActivityType = "WRITING"
	ActivitySpeaking      ActivityType = "SPEAKING"
	ActivityCheckpoint    ActivityType = "CHECKPOINT"
)

type ContentStatus string

const (
	ContentDraft     ContentStatus = "DRAFT"
	ContentPublished ContentStatus = "PUBLISHED"
)

type ActivityProgressStatus string

const (
	ActivityInProgress ActivityProgressStatus = "IN_PROGRESS"
	ActivityPassed     ActivityProgressStatus = "PASSED"
	ActivityFailed     ActivityProgressStatus = "FAILED"
)

// ExerciseType is a scored item format (docs/foundation-spec.md §5).
type ExerciseType string

const (
	ExerciseMultipleChoice  ExerciseType = "MULTIPLE_CHOICE"
	ExerciseTrueFalse       ExerciseType = "TRUE_FALSE"
	ExerciseFillBlank       ExerciseType = "FILL_BLANK"
	ExerciseReorder         ExerciseType = "REORDER"
	ExerciseErrorCorrection ExerciseType = "ERROR_CORRECTION"
	ExerciseTransform       ExerciseType = "TRANSFORM"
	ExerciseTranslateViEn   ExerciseType = "TRANSLATE_VI_EN"
	ExerciseDictation       ExerciseType = "DICTATION"
	ExerciseListenGapFill   ExerciseType = "LISTEN_GAP_FILL"
	ExerciseListenPickImage ExerciseType = "LISTEN_PICK_IMAGE"
	ExerciseListenNumbers   ExerciseType = "LISTEN_NUMBERS"
	ExerciseMinimalPair     ExerciseType = "MINIMAL_PAIR"
	ExerciseIPAToWord       ExerciseType = "IPA_TO_WORD"
	ExerciseWordToIPA       ExerciseType = "WORD_TO_IPA"
	ExerciseEndingSound     ExerciseType = "ENDING_SOUND"
	ExerciseWordStress      ExerciseType = "WORD_STRESS"
	ExerciseSpeak           ExerciseType = "SPEAK"
)

// LessonActivity is an ordered skill step of a CORE lesson.
type LessonActivity struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	LessonID       uint           `gorm:"uniqueIndex:idx_lesson_activity_slug;index:idx_lesson_activity_order;not null" json:"lessonId"`
	Slug           string         `gorm:"size:150;uniqueIndex:idx_lesson_activity_slug;not null" json:"slug"`
	Type           ActivityType   `gorm:"type:varchar(30);not null" json:"type"`
	Title          string         `gorm:"size:255;not null" json:"title"`
	InstructionsVI string         `gorm:"column:instructions_vi;type:text" json:"instructionsVi"`
	GrammarPointID *uint          `gorm:"index" json:"grammarPointId"`
	PassageID      *uint          `gorm:"index" json:"passageId"`
	ScenarioID     *uint          `gorm:"index" json:"scenarioId"`
	ContentJSON    datatypes.JSON `gorm:"type:json" json:"content,omitempty"`
	PassScore      int            `gorm:"default:70;not null" json:"passScore"`
	XPReward       int            `gorm:"column:xp_reward;default:10;not null" json:"xpReward"`
	IsRequired     bool           `gorm:"default:true;not null" json:"isRequired"`
	OrderIndex     int            `gorm:"default:0;index:idx_lesson_activity_order" json:"orderIndex"`
	Status         ContentStatus  `gorm:"type:varchar(20);default:'PUBLISHED';not null" json:"status"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Lesson    Lesson     `gorm:"foreignKey:LessonID" json:"-"`
	Exercises []Exercise `gorm:"foreignKey:ActivityID" json:"-"`
}

// UserActivityProgress tracks one learner's attempts at an activity.
type UserActivityProgress struct {
	ID            uint                   `gorm:"primaryKey" json:"id"`
	UserID        uint                   `gorm:"uniqueIndex:idx_user_activity;not null" json:"userId"`
	ActivityID    uint                   `gorm:"uniqueIndex:idx_user_activity;index;not null" json:"activityId"`
	Status        ActivityProgressStatus `gorm:"type:varchar(20);default:'IN_PROGRESS';not null" json:"status"`
	Score         *int                   `json:"score"`
	BestScore     *int                   `json:"bestScore"`
	Attempts      int                    `gorm:"default:0;not null" json:"attempts"`
	PassedAt      *time.Time             `json:"passedAt"`
	LastAttemptAt *time.Time             `json:"lastAttemptAt"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

// Exercise is a scored item inside an activity. PayloadJSON is public;
// AnswerJSON holds the private answer key and must never be serialised
// before the learner answers the item.
type Exercise struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ActivityID    uint           `gorm:"index:idx_exercise_activity_order;not null" json:"activityId"`
	Type          ExerciseType   `gorm:"type:varchar(40);not null" json:"type"`
	Prompt        string         `gorm:"type:text;not null" json:"prompt"`
	PayloadJSON   datatypes.JSON `gorm:"type:json;not null" json:"payload"`
	AnswerJSON    datatypes.JSON `gorm:"type:json;not null" json:"-"`
	ExplanationVI string         `gorm:"column:explanation_vi;type:text" json:"-"`
	SkillTag      string         `gorm:"size:80;index" json:"skillTag"`
	Points        int            `gorm:"default:10;not null" json:"points"`
	OrderIndex    int            `gorm:"default:0;index:idx_exercise_activity_order" json:"orderIndex"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

// UserExerciseAnswer stores the learner's latest answer to an exercise
// within one activity attempt.
type UserExerciseAnswer struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"uniqueIndex:idx_user_exercise_attempt;index;not null" json:"userId"`
	ExerciseID      uint           `gorm:"uniqueIndex:idx_user_exercise_attempt;index;not null" json:"exerciseId"`
	ActivityAttempt int            `gorm:"uniqueIndex:idx_user_exercise_attempt;not null" json:"activityAttempt"`
	AnswerJSON      datatypes.JSON `gorm:"type:json" json:"answer"`
	IsCorrect       bool           `gorm:"not null" json:"isCorrect"`
	ScoreRatio      float64        `gorm:"type:decimal(4,3);not null" json:"scoreRatio"`
	AnsweredAt      time.Time      `gorm:"not null" json:"answeredAt"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}
