// Package mistakes records learner mistakes (docs/foundation-spec.md §3.4),
// counted when the same correction repeats.
package mistakes

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

// Upsert saves a mistake or, for a known (user, skill tag, correction),
// counts it again and reopens it.
func Upsert(db *gorm.DB, mistake models.LearnerMistake) error {
	mistake.SkillTag = truncate(mistake.SkillTag, 80)
	mistake.WrongText = truncate(mistake.WrongText, 500)
	mistake.CorrectText = truncate(mistake.CorrectText, 500)
	if mistake.Occurrences == 0 {
		mistake.Occurrences = 1
	}
	if mistake.Status == "" {
		mistake.Status = models.MistakeOpen
	}
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "skill_tag"}, {Name: "correct_text"}},
		DoUpdates: clause.Assignments(map[string]any{
			"occurrences":    gorm.Expr("occurrences + 1"),
			"wrong_text":     mistake.WrongText,
			"explanation_vi": mistake.ExplanationVI,
			"source":         mistake.Source,
			// A repeated mistake reopens; one being reviewed stays in the
			// queue but starts its streak again and is due now.
			"status":         gorm.Expr("CASE WHEN status = ? THEN status ELSE ? END", models.MistakeReviewing, models.MistakeOpen),
			"review_streak":  0,
			"fixed_at":       nil,
			"next_review_at": gorm.Expr("CASE WHEN status = ? THEN ? ELSE next_review_at END", models.MistakeReviewing, mistake.LastSeenAt),
			"last_seen_at":   mistake.LastSeenAt,
			"updated_at":     mistake.LastSeenAt,
		}),
	}).Create(&mistake).Error
	if err != nil {
		return fmt.Errorf("save learner mistake: %w", err)
	}
	return nil
}

func truncate(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}
