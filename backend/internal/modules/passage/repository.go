package passage

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

var (
	ErrPassageNotFound    = errors.New("passage not found")
	ErrVocabularyNotFound = errors.New("vocabulary not found")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindPassage(id uint) (models.Passage, error) {
	var passage models.Passage
	err := r.db.First(&passage, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Passage{}, ErrPassageNotFound
	}
	if err != nil {
		return models.Passage{}, fmt.Errorf("find passage: %w", err)
	}
	return passage, nil
}

// FindWord returns the best vocabulary row for the candidate forms: earlier
// candidates first, A1 Foundation words before other courses' words.
func (r Repository) FindWord(candidates []string, slug string) (models.Vocabulary, bool, error) {
	var rows []models.Vocabulary
	query := r.db.Model(&models.Vocabulary{})
	if slug != "" {
		query = query.Where("slug = ?", slug)
	} else {
		if len(candidates) == 0 {
			return models.Vocabulary{}, false, nil
		}
		query = query.Where("LOWER(word) IN ?", candidates)
	}
	if err := query.Find(&rows).Error; err != nil {
		return models.Vocabulary{}, false, fmt.Errorf("find word: %w", err)
	}
	if len(rows) == 0 {
		return models.Vocabulary{}, false, nil
	}

	rank := func(row models.Vocabulary) int {
		score := len(candidates) * 2
		for i, candidate := range candidates {
			if equalFold(row.Word, candidate) {
				score = i * 2
				break
			}
		}
		if len(row.Slug) < 3 || row.Slug[:3] != "a1-" {
			score++
		}
		return score
	}
	best := rows[0]
	for _, row := range rows[1:] {
		if rank(row) < rank(best) {
			best = row
		}
	}
	return best, true, nil
}

func (r Repository) HasProgress(userID uint, vocabularyID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserVocabularyProgress{}).
		Where("user_id = ? AND vocabulary_id = ?", userID, vocabularyID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("find vocabulary progress: %w", err)
	}
	return count > 0, nil
}

func (r Repository) VocabularyExists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Vocabulary{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("find vocabulary: %w", err)
	}
	return count > 0, nil
}

// AddToReviews creates a NEW progress row due now, unless one exists.
// It returns the row and whether it was created.
func (r Repository) AddToReviews(userID uint, vocabularyID uint, now time.Time) (models.UserVocabularyProgress, bool, error) {
	var existing models.UserVocabularyProgress
	err := r.db.Where("user_id = ? AND vocabulary_id = ?", userID, vocabularyID).Take(&existing).Error
	if err == nil {
		return existing, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.UserVocabularyProgress{}, false, fmt.Errorf("find vocabulary progress: %w", err)
	}
	row := models.UserVocabularyProgress{
		UserID:       userID,
		VocabularyID: vocabularyID,
		Status:       models.VocabularyStatusNew,
		NextReviewAt: &now,
	}
	// A concurrent add hits the unique index; that is fine.
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
		return models.UserVocabularyProgress{}, false, fmt.Errorf("add flashcard: %w", err)
	}
	return row, true, nil
}

func equalFold(a string, b string) bool {
	return len(a) == len(b) && (a == b || lower(a) == lower(b))
}

func lower(value string) string {
	out := []byte(value)
	for i, c := range out {
		if c >= 'A' && c <= 'Z' {
			out[i] = c + 32
		}
	}
	return string(out)
}
