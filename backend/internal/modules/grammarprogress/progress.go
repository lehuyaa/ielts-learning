// Package grammarprogress keeps each learner's mastery of grammar points up
// to date. It is shared by the lesson activity and grammar practice modules.
package grammarprogress

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

// Window is how many recent answers mastery is computed from.
const Window = 20

const tagPrefix = "grammar:"

// SlugFromSkillTag returns the grammar point slug of a "grammar:<slug>" tag.
func SlugFromSkillTag(skillTag string) (string, bool) {
	if !strings.HasPrefix(skillTag, tagPrefix) {
		return "", false
	}
	slug := strings.TrimPrefix(skillTag, tagPrefix)
	return slug, slug != ""
}

// Record adds one graded answer to the learner's mastery of the grammar point
// named by skillTag. Tags that are not grammar tags, or name an unknown
// grammar point, are ignored.
func Record(db *gorm.DB, userID uint, skillTag string, scoreRatio float64, isCorrect bool, at time.Time) error {
	slug, ok := SlugFromSkillTag(skillTag)
	if !ok {
		return nil
	}

	var point models.GrammarPoint
	err := db.Where("slug = ?", slug).First(&point).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("find grammar point: %w", err)
	}

	var progress models.UserGrammarProgress
	err = db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND grammar_point_id = ?", userID, point.ID).
		First(&progress).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("find grammar progress: %w", err)
	}
	progress.UserID = userID
	progress.GrammarPointID = point.ID

	recent := AddRecent(DecodeRecent(progress.RecentJSON), scoreRatio)
	encoded, _ := json.Marshal(recent)
	progress.RecentJSON = encoded
	progress.Mastery = Mastery(recent)
	if isCorrect {
		progress.CorrectCount++
	} else {
		progress.WrongCount++
	}
	progress.LastPracticed = &at

	if err := db.Save(&progress).Error; err != nil {
		return fmt.Errorf("save grammar progress: %w", err)
	}
	return nil
}

func DecodeRecent(raw []byte) []float64 {
	var recent []float64
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &recent)
	}
	return recent
}

// AddRecent appends a score and keeps only the last Window scores.
func AddRecent(recent []float64, scoreRatio float64) []float64 {
	if scoreRatio < 0 {
		scoreRatio = 0
	}
	if scoreRatio > 1 {
		scoreRatio = 1
	}
	recent = append(recent, scoreRatio)
	if len(recent) > Window {
		recent = recent[len(recent)-Window:]
	}
	return recent
}

// Mastery is the average of the recent scores as a percentage.
func Mastery(recent []float64) int {
	if len(recent) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range recent {
		sum += value
	}
	return int(math.Round(100 * sum / float64(len(recent))))
}
