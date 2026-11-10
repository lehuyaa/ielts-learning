package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// GrammarPoint is a theory card (docs/foundation-spec.md §3.3). Exercises
// practise it through the skill tag "grammar:<slug>".
type GrammarPoint struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Slug               string         `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	CEFRLevel          string         `gorm:"column:cefr_level;size:5;index" json:"cefrLevel"`
	Title              string         `gorm:"size:255;not null" json:"title"`
	TitleVI            string         `gorm:"column:title_vi;size:255" json:"titleVi"`
	SummaryVI          string         `gorm:"column:summary_vi;type:text" json:"summaryVi"`
	Formula            string         `gorm:"type:text" json:"formula"`
	ExamplesJSON       datatypes.JSON `gorm:"type:json" json:"examples"`
	CommonMistakesJSON datatypes.JSON `gorm:"type:json" json:"commonMistakes"`
	OrderIndex         int            `gorm:"default:0" json:"orderIndex"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// SkillTag is the exercise tag that practises this grammar point.
func (g GrammarPoint) SkillTag() string {
	return "grammar:" + g.Slug
}

// UserGrammarProgress is a learner's mastery of one grammar point: the
// average score of their last answers tagged with it.
type UserGrammarProgress struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex:idx_user_grammar;not null" json:"userId"`
	GrammarPointID uint           `gorm:"uniqueIndex:idx_user_grammar;index;not null" json:"grammarPointId"`
	Mastery        int            `gorm:"default:0;not null" json:"mastery"`
	CorrectCount   int            `gorm:"default:0;not null" json:"correctCount"`
	WrongCount     int            `gorm:"default:0;not null" json:"wrongCount"`
	RecentJSON     datatypes.JSON `gorm:"type:json" json:"-"`
	LastPracticed  *time.Time     `gorm:"column:last_practiced_at" json:"lastPracticedAt"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}
