package aiconversation

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrActivityMismatch = errors.New("activity does not use this scenario")
)

func (r Repository) DB() *gorm.DB {
	return r.db
}

// FindScenarioForUser returns the user's own scenario with that slug, else
// the system scenario.
func (r Repository) FindScenarioForUser(userID uint, slug string) (models.AIConversationScenario, error) {
	var scenario models.AIConversationScenario
	err := r.db.Where("slug = ? AND (user_id = ? OR is_system = ?)", slug, userID, true).
		Order("is_system ASC").First(&scenario).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.AIConversationScenario{}, ErrScenarioNotFound
	}
	if err != nil {
		return models.AIConversationScenario{}, fmt.Errorf("find scenario: %w", err)
	}
	return scenario, nil
}

func (r Repository) FindScenarioByID(id uint) (models.AIConversationScenario, error) {
	var scenario models.AIConversationScenario
	if err := r.db.First(&scenario, id).Error; err != nil {
		return models.AIConversationScenario{}, fmt.Errorf("find scenario: %w", err)
	}
	return scenario, nil
}

// FindActivity returns a lesson activity (for sessions started from a lesson).
func (r Repository) FindActivity(id uint) (models.LessonActivity, error) {
	var activity models.LessonActivity
	err := r.db.First(&activity, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.LessonActivity{}, ErrActivityMismatch
	}
	if err != nil {
		return models.LessonActivity{}, fmt.Errorf("find activity: %w", err)
	}
	return activity, nil
}

func (r Repository) CreateSession(session *models.AIConversationSession) error {
	if err := r.db.Create(session).Error; err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r Repository) SaveSession(session *models.AIConversationSession) error {
	if err := r.db.Save(session).Error; err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

// FindOpenSession returns the user's unfinished session of a scenario (and activity).
func (r Repository) FindOpenSession(userID uint, scenarioID uint, activityID *uint) (models.AIConversationSession, bool, error) {
	query := r.db.Where("user_id = ? AND scenario_id = ? AND status = ?", userID, scenarioID, models.SessionInProgress)
	if activityID != nil {
		query = query.Where("activity_id = ?", *activityID)
	} else {
		query = query.Where("activity_id IS NULL")
	}
	var session models.AIConversationSession
	err := query.Order("id DESC").Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.AIConversationSession{}, false, nil
	}
	if err != nil {
		return models.AIConversationSession{}, false, fmt.Errorf("find open session: %w", err)
	}
	return session, true, nil
}

func (r Repository) FindSession(userID uint, id uint, forUpdate bool) (models.AIConversationSession, error) {
	query := r.db.Where("id = ? AND user_id = ?", id, userID)
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var session models.AIConversationSession
	err := query.Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.AIConversationSession{}, ErrSessionNotFound
	}
	if err != nil {
		return models.AIConversationSession{}, fmt.Errorf("find session: %w", err)
	}
	return session, nil
}

func (r Repository) FindSessionMessages(sessionID uint) ([]models.AIConversationMessage, error) {
	var messages []models.AIConversationMessage
	if err := r.db.Where("session_id = ?", sessionID).Order("id ASC").Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("find session messages: %w", err)
	}
	return messages, nil
}

// SessionStats summarises a user's finished sessions per scenario.
type SessionStats struct {
	ScenarioID uint `json:"scenarioId"`
	Sessions   int  `json:"sessions"`
	BestScore  *int `json:"bestScore"`
}

func (r Repository) FinishedStats(userID uint) (map[uint]SessionStats, error) {
	var rows []SessionStats
	err := r.db.Model(&models.AIConversationSession{}).
		Select("scenario_id, COUNT(*) AS sessions, MAX(score) AS best_score").
		Where("user_id = ? AND status = ?", userID, models.SessionFinished).
		Group("scenario_id").Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("session stats: %w", err)
	}
	stats := make(map[uint]SessionStats, len(rows))
	for _, row := range rows {
		stats[row.ScenarioID] = row
	}
	return stats, nil
}

// hintsAllowed: never in checkpoints; otherwise unless content.hintsAllowed is false.
func hintsAllowed(activity *models.LessonActivity) bool {
	if activity == nil {
		return true
	}
	if activity.Type == models.ActivityCheckpoint {
		return false
	}
	var content struct {
		HintsAllowed *bool `json:"hintsAllowed"`
	}
	if len(activity.ContentJSON) > 0 && json.Unmarshal(activity.ContentJSON, &content) == nil && content.HintsAllowed != nil {
		return *content.HintsAllowed
	}
	return true
}
