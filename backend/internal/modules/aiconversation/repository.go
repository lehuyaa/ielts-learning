package aiconversation

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Create(scenario *models.AIConversationScenario) error {
	if err := r.db.Create(scenario).Error; err != nil {
		return fmt.Errorf("create ai conversation scenario: %w", err)
	}

	return nil
}

func (r Repository) FindByUser(userID uint) ([]models.AIConversationScenario, error) {
	var scenarios []models.AIConversationScenario
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at ASC, id ASC").
		Find(&scenarios).Error; err != nil {
		return nil, fmt.Errorf("find ai conversation scenarios by user: %w", err)
	}

	return scenarios, nil
}

func (r Repository) ExistsBySlug(userID uint, slug string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.AIConversationScenario{}).
		Where("user_id = ? AND slug = ?", userID, slug).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check ai conversation scenario slug: %w", err)
	}

	return count > 0, nil
}

func (r Repository) FindByUserAndSlug(userID uint, slug string) (models.AIConversationScenario, error) {
	var scenario models.AIConversationScenario
	if err := r.db.
		Where("user_id = ? AND slug = ?", userID, slug).
		First(&scenario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AIConversationScenario{}, ErrScenarioNotFound
		}
		return models.AIConversationScenario{}, fmt.Errorf("find ai conversation scenario by user and slug: %w", err)
	}

	return scenario, nil
}

func (r Repository) CreateMessage(message *models.AIConversationMessage) error {
	if err := r.db.Create(message).Error; err != nil {
		return fmt.Errorf("create ai conversation message: %w", err)
	}

	return nil
}

func (r Repository) FindMessagesByScenario(scenarioID uint) ([]models.AIConversationMessage, error) {
	var messages []models.AIConversationMessage
	if err := r.db.
		Where("scenario_id = ?", scenarioID).
		Order("created_at ASC, id ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("find ai conversation messages by scenario: %w", err)
	}

	return messages, nil
}
