package vocabulary

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var (
	ErrVocabularyNotFound = errors.New("vocabulary not found")
	ErrUserNotFound       = errors.New("user not found")
)

type VocabularyTopicContext struct {
	VocabularyID uint
	TopicID      uint
	TopicTitle   string
	TopicSlug    string
	TopicIcon    string
	TopicEmoji   string
	TopicColor   string
	BandLevelID  uint
	BandScore    float64
	BandTitle    string
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindUser(userID uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func (r Repository) FindVocabularies(userID uint, query ListQuery) ([]models.Vocabulary, int64, error) {
	db := r.filteredVocabularyQuery(userID, query)

	var total int64
	if err := db.Distinct("vocabularies.id").Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count vocabularies: %w", err)
	}

	var vocabularies []models.Vocabulary
	err := db.
		Distinct("vocabularies.*").
		Order("vocabularies.word ASC, vocabularies.id ASC").
		Limit(query.Limit).
		Offset((query.Page - 1) * query.Limit).
		Find(&vocabularies).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find vocabularies: %w", err)
	}

	return vocabularies, total, nil
}

func (r Repository) FindVocabulary(vocabularyID uint) (models.Vocabulary, error) {
	var vocabulary models.Vocabulary
	err := r.db.First(&vocabulary, vocabularyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Vocabulary{}, ErrVocabularyNotFound
		}
		return models.Vocabulary{}, fmt.Errorf("find vocabulary: %w", err)
	}

	return vocabulary, nil
}

func (r Repository) FindVocabularyProgress(userID uint, vocabularyIDs []uint) ([]models.UserVocabularyProgress, error) {
	if len(vocabularyIDs) == 0 {
		return []models.UserVocabularyProgress{}, nil
	}

	var progress []models.UserVocabularyProgress
	err := r.db.
		Where("user_id = ? AND vocabulary_id IN ?", userID, vocabularyIDs).
		Find(&progress).Error
	if err != nil {
		return nil, fmt.Errorf("find vocabulary progress: %w", err)
	}

	return progress, nil
}

func (r Repository) FindTopicContext(vocabularyIDs []uint) (map[uint]VocabularyTopicContext, error) {
	contexts := make(map[uint]VocabularyTopicContext)
	if len(vocabularyIDs) == 0 {
		return contexts, nil
	}

	var rows []VocabularyTopicContext
	err := r.db.
		Table("lesson_vocabularies").
		Select(`
			lesson_vocabularies.vocabulary_id AS vocabulary_id,
			topics.id AS topic_id,
			topics.title AS topic_title,
			topics.slug AS topic_slug,
			topics.icon AS topic_icon,
			topics.emoji AS topic_emoji,
			topics.color AS topic_color,
			band_levels.id AS band_level_id,
			band_levels.band_score AS band_score,
			band_levels.title AS band_title
		`).
		Joins("JOIN lessons ON lessons.id = lesson_vocabularies.lesson_id").
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("lesson_vocabularies.vocabulary_id IN ?", vocabularyIDs).
		Where("lesson_vocabularies.is_required = ?", true).
		Order("lesson_vocabularies.vocabulary_id ASC, lessons.order_index ASC, lesson_vocabularies.order_index ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find vocabulary topic context: %w", err)
	}

	for _, row := range rows {
		if _, exists := contexts[row.VocabularyID]; exists {
			continue
		}
		contexts[row.VocabularyID] = row
	}

	return contexts, nil
}

func (r Repository) FindTopicsBySlugs(slugs []string) (map[string]models.Topic, error) {
	result := make(map[string]models.Topic)
	if len(slugs) == 0 {
		return result, nil
	}

	var topics []models.Topic
	if err := r.db.Where("slug IN ?", slugs).Find(&topics).Error; err != nil {
		return nil, fmt.Errorf("find topics by slug: %w", err)
	}

	for _, topic := range topics {
		result[topic.Slug] = topic
	}
	return result, nil
}

func (r Repository) FindLessonsByTopicIDsAndSlugs(topicIDs []uint, slugs []string) ([]models.Lesson, error) {
	if len(topicIDs) == 0 || len(slugs) == 0 {
		return nil, nil
	}

	var lessons []models.Lesson
	err := r.db.
		Where("topic_id IN ? AND slug IN ?", topicIDs, slugs).
		Find(&lessons).Error
	if err != nil {
		return nil, fmt.Errorf("find lessons by topic and slug: %w", err)
	}

	return lessons, nil
}

type ImportCommitStats struct {
	VocabulariesCreated int
	VocabulariesUpdated int
	LessonLinksCreated  int
	LessonLinksUpdated  int
}

// CommitImportRows persists every valid row in a single transaction, upserting
// the Vocabulary by slug and the LessonVocabulary link by (lessonId, vocabularyId).
// Invalid rows are skipped. Each committed row's VocabularyAction/LinkAction is
// set to CREATE or UPDATE so the caller can report a per-row outcome.
func (r Repository) CommitImportRows(rows []ImportRow) (ImportCommitStats, error) {
	var stats ImportCommitStats

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for i := range rows {
			row := &rows[i]
			if !row.Valid() {
				continue
			}

			vocabulary, vocabularyAction, err := upsertVocabulary(tx, *row)
			if err != nil {
				return fmt.Errorf("row %d: upsert vocabulary: %w", row.RowNumber, err)
			}
			row.VocabularyAction = vocabularyAction
			if vocabularyAction == "CREATE" {
				stats.VocabulariesCreated++
			} else {
				stats.VocabulariesUpdated++
			}

			linkAction, err := upsertLessonVocabulary(tx, row.LessonID, vocabulary.ID, row.OrderIndex, row.IsRequired)
			if err != nil {
				return fmt.Errorf("row %d: upsert lesson vocabulary: %w", row.RowNumber, err)
			}
			row.LinkAction = linkAction
			if linkAction == "CREATE" {
				stats.LessonLinksCreated++
			} else {
				stats.LessonLinksUpdated++
			}
		}

		return nil
	})
	if err != nil {
		return ImportCommitStats{}, err
	}

	return stats, nil
}

func upsertVocabulary(tx *gorm.DB, row ImportRow) (models.Vocabulary, string, error) {
	synonymsJSON, _ := json.Marshal(row.Synonyms)
	antonymsJSON, _ := json.Marshal(row.Antonyms)
	collocationsJSON, _ := json.Marshal(row.Collocations)

	var vocabulary models.Vocabulary
	err := tx.Where("slug = ?", row.Slug).First(&vocabulary).Error
	switch {
	case err == nil:
		vocabulary.Word = row.Word
		vocabulary.IPA = row.IPA
		vocabulary.PartOfSpeech = row.PartOfSpeech
		vocabulary.MeaningVI = row.MeaningVI
		vocabulary.MeaningEN = row.MeaningEN
		vocabulary.ShortDefinition = row.ShortDefinition
		vocabulary.ExampleSentence = row.ExampleSentence
		vocabulary.ExampleMeaningVI = row.ExampleMeaningVI
		vocabulary.ExampleSource = row.ExampleSource
		vocabulary.SynonymsJSON = synonymsJSON
		vocabulary.AntonymsJSON = antonymsJSON
		vocabulary.CollocationsJSON = collocationsJSON
		vocabulary.Difficulty = row.Difficulty
		vocabulary.TargetBand = row.TargetBand

		if err := tx.Save(&vocabulary).Error; err != nil {
			return models.Vocabulary{}, "", err
		}
		return vocabulary, "UPDATE", nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		vocabulary = models.Vocabulary{
			Word:             row.Word,
			Slug:             row.Slug,
			IPA:              row.IPA,
			PartOfSpeech:     row.PartOfSpeech,
			MeaningVI:        row.MeaningVI,
			MeaningEN:        row.MeaningEN,
			ShortDefinition:  row.ShortDefinition,
			ExampleSentence:  row.ExampleSentence,
			ExampleMeaningVI: row.ExampleMeaningVI,
			ExampleSource:    row.ExampleSource,
			SynonymsJSON:     synonymsJSON,
			AntonymsJSON:     antonymsJSON,
			CollocationsJSON: collocationsJSON,
			Difficulty:       row.Difficulty,
			TargetBand:       row.TargetBand,
		}
		if err := tx.Create(&vocabulary).Error; err != nil {
			return models.Vocabulary{}, "", err
		}
		return vocabulary, "CREATE", nil
	default:
		return models.Vocabulary{}, "", err
	}
}

func upsertLessonVocabulary(tx *gorm.DB, lessonID uint, vocabularyID uint, orderIndex int, isRequired bool) (string, error) {
	var link models.LessonVocabulary
	err := tx.Where("lesson_id = ? AND vocabulary_id = ?", lessonID, vocabularyID).First(&link).Error
	switch {
	case err == nil:
		link.OrderIndex = orderIndex
		link.IsRequired = isRequired
		if err := tx.Save(&link).Error; err != nil {
			return "", err
		}
		return "UPDATE", nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		link = models.LessonVocabulary{
			LessonID:     lessonID,
			VocabularyID: vocabularyID,
			OrderIndex:   orderIndex,
			IsRequired:   isRequired,
		}
		if err := tx.Create(&link).Error; err != nil {
			return "", err
		}
		return "CREATE", nil
	default:
		return "", err
	}
}

func (r Repository) filteredVocabularyQuery(userID uint, query ListQuery) *gorm.DB {
	db := r.db.
		Model(&models.Vocabulary{}).
		Joins("LEFT JOIN user_vocabulary_progresses ON user_vocabulary_progresses.vocabulary_id = vocabularies.id AND user_vocabulary_progresses.user_id = ?", userID)

	needsTopicJoin := query.Q != ""
	if needsTopicJoin {
		db = db.
			Joins("LEFT JOIN lesson_vocabularies ON lesson_vocabularies.vocabulary_id = vocabularies.id").
			Joins("LEFT JOIN lessons ON lessons.id = lesson_vocabularies.lesson_id").
			Joins("LEFT JOIN topics ON topics.id = lessons.topic_id")
	}

	if query.Q != "" {
		like := "%" + query.Q + "%"
		db = db.Where(
			"vocabularies.word LIKE ? OR vocabularies.meaning_vi LIKE ? OR vocabularies.meaning_en LIKE ? OR vocabularies.short_definition LIKE ? OR topics.title LIKE ?",
			like,
			like,
			like,
			like,
			like,
		)
	}

	if query.Difficulty != nil {
		db = db.Where("vocabularies.difficulty = ?", *query.Difficulty)
	}

	if query.TargetBand != nil {
		db = db.Where("vocabularies.target_band = ?", *query.TargetBand)
	}

	if query.Status != nil {
		if *query.Status == models.VocabularyStatusNew {
			db = db.Where("user_vocabulary_progresses.id IS NULL OR user_vocabulary_progresses.status = ?", models.VocabularyStatusNew)
		} else {
			db = db.Where("user_vocabulary_progresses.status = ?", *query.Status)
		}
	}

	return db
}
