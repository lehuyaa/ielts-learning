// Package aiusage counts AI calls per user per day (docs/foundation-spec.md §6.7).
package aiusage

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

// Take counts one call of kind for today and reports whether it is allowed
// (false once limit calls were made; a limit <= 0 allows nothing).
func Take(db *gorm.DB, userID uint, kind string, limit int, now time.Time) (bool, error) {
	if limit <= 0 {
		return false, nil
	}
	day := now.Format("2006-01-02")
	allowed := false
	err := db.Transaction(func(tx *gorm.DB) error {
		row := models.AIUsageDaily{UserID: userID, Day: day, Kind: kind}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND day = ? AND kind = ?", userID, day, kind).
			Take(&row).Error; err != nil {
			return err
		}
		if row.Count >= limit {
			return nil
		}
		allowed = true
		return tx.Model(&row).Update("count", gorm.Expr("count + 1")).Error
	})
	return allowed, err
}
