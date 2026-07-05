package activity

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindUserForUpdate(tx *gorm.DB, userID uint) (models.User, error) {
	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
		return models.User{}, fmt.Errorf("find activity user: %w", err)
	}

	return user, nil
}

func (r Repository) FindDailyActivityByDate(tx *gorm.DB, userID uint, date time.Time) (models.DailyActivity, bool, error) {
	var activity models.DailyActivity
	err := tx.
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		First(&activity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DailyActivity{}, false, nil
		}
		return models.DailyActivity{}, false, fmt.Errorf("find daily activity by date: %w", err)
	}

	return activity, true, nil
}

func (r Repository) UpsertDailyActivity(tx *gorm.DB, activity models.DailyActivity, event Event) error {
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "date"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"words_learned":  gorm.Expr("words_learned + ?", event.WordsLearnedDelta),
			"words_reviewed": gorm.Expr("words_reviewed + ?", event.WordsReviewedDelta),
			"quizzes_taken":  gorm.Expr("quizzes_taken + ?", event.QuizzesTakenDelta),
			"lessons_done":   gorm.Expr("lessons_done + ?", event.LessonsDoneDelta),
			"active_minutes": gorm.Expr("active_minutes + ?", event.ActiveMinutesDelta),
			"xp_earned":      gorm.Expr("xp_earned + ?", event.XPEarnedDelta),
			"updated_at":     activity.UpdatedAt,
		}),
	}).Create(&activity).Error; err != nil {
		return fmt.Errorf("upsert daily activity: %w", err)
	}

	return nil
}

func (r Repository) SaveUser(tx *gorm.DB, user *models.User) error {
	if err := tx.Save(user).Error; err != nil {
		return fmt.Errorf("save activity user: %w", err)
	}

	return nil
}
