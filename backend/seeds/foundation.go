package seeds

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

// foundationA1JSON describes the A1 Foundation band level: its topics and
// lessons in study order. It is generated together with the A1 vocabulary
// import file, so topic/lesson slugs always match the Excel rows.
//
//go:embed data/foundation_a1.json
var foundationA1JSON []byte

const (
	foundationQuizQuestionsPerLesson = 5
	foundationLessonXPReward         = 30
)

type foundationFile struct {
	BandLevel struct {
		BandScore float64 `json:"bandScore"`
		Title     string  `json:"title"`
	} `json:"bandLevel"`
	Topics []foundationTopic `json:"topics"`
}

type foundationTopic struct {
	Slug        string             `json:"slug"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Icon        string             `json:"icon"`
	Emoji       string             `json:"emoji"`
	Color       string             `json:"color"`
	OrderIndex  int                `json:"orderIndex"`
	Lessons     []foundationLesson `json:"lessons"`
}

type foundationLesson struct {
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	OrderIndex int    `json:"orderIndex"`
}

// RunFoundation seeds the A1 Foundation band level (topics + lessons) into the
// published course, without touching demo users or progress. It is idempotent:
// run it once before importing the A1 vocabulary Excel file, and again after the
// import to (re)build each lesson's quiz from its imported words.
func RunFoundation(db *gorm.DB) error {
	var file foundationFile
	if err := json.Unmarshal(foundationA1JSON, &file); err != nil {
		return fmt.Errorf("parse foundation data: %w", err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		course, err := seedCourse(tx)
		if err != nil {
			return err
		}

		band, err := seedFoundationBand(tx, course.ID, file)
		if err != nil {
			return err
		}

		lessons, err := seedFoundationTopicsAndLessons(tx, band, file.Topics)
		if err != nil {
			return err
		}

		if err := seedFoundationQuiz(tx, lessons); err != nil {
			return err
		}

		return updateCourseTotals(tx, course.ID)
	})
}

func seedFoundationBand(db *gorm.DB, courseID uint, file foundationFile) (models.BandLevel, error) {
	minScore := 1.0
	maxScore := 3.5
	statusLabel := "Start here"
	band := models.BandLevel{
		CourseID:    courseID,
		BandScore:   file.BandLevel.BandScore,
		MinScore:    &minScore,
		MaxScore:    &maxScore,
		Title:       file.BandLevel.Title,
		Description: "Nền tảng tiếng Anh cho người mất gốc: từ vựng A1 thông dụng nhất, có phiên âm IPA, nghĩa tiếng Việt và câu ví dụ.",
		StatusLabel: &statusLabel,
		// IELTS band levels use 1..4, so 0 puts A1 first on the roadmap.
		OrderIndex: 0,
	}

	if err := db.Where("course_id = ? AND band_score = ?", courseID, band.BandScore).Assign(band).FirstOrCreate(&band).Error; err != nil {
		return models.BandLevel{}, fmt.Errorf("seed foundation band: %w", err)
	}

	return band, nil
}

func seedFoundationTopicsAndLessons(db *gorm.DB, band models.BandLevel, topics []foundationTopic) ([]models.Lesson, error) {
	bandMin := 1.0
	bandMax := 3.5
	timeLimit := 600

	lessons := make([]models.Lesson, 0)
	for _, seed := range topics {
		topic := models.Topic{
			BandLevelID: band.ID,
			Title:       seed.Title,
			Slug:        seed.Slug,
			Description: seed.Description,
			Icon:        seed.Icon,
			Emoji:       seed.Emoji,
			Color:       seed.Color,
			OrderIndex:  seed.OrderIndex,
		}
		if err := db.Where("band_level_id = ? AND slug = ?", band.ID, seed.Slug).Assign(topic).FirstOrCreate(&topic).Error; err != nil {
			return nil, fmt.Errorf("seed foundation topic %s: %w", seed.Slug, err)
		}

		for _, lessonData := range seed.Lessons {
			lesson := models.Lesson{
				TopicID:              topic.ID,
				Title:                lessonData.Title,
				Slug:                 lessonData.Slug,
				Description:          fmt.Sprintf("%s: %s.", seed.Title, seed.Description),
				RequiredScore:        80,
				EstimatedMinutes:     10,
				BandMin:              &bandMin,
				BandMax:              &bandMax,
				XPReward:             foundationLessonXPReward,
				QuizTimeLimitSeconds: &timeLimit,
				OrderIndex:           lessonData.OrderIndex,
				IsPublished:          true,
			}
			if err := db.Where("topic_id = ? AND slug = ?", topic.ID, lessonData.Slug).Assign(lesson).FirstOrCreate(&lesson).Error; err != nil {
				return nil, fmt.Errorf("seed foundation lesson %s: %w", lessonData.Slug, err)
			}
			lessons = append(lessons, lesson)
		}
	}

	return lessons, nil
}

// seedFoundationQuiz builds Vietnamese-meaning questions ("X nghĩa là gì?")
// from the words already imported into each lesson. Lessons without at least
// four words are skipped, so running this before the import is harmless.
func seedFoundationQuiz(db *gorm.DB, lessons []models.Lesson) error {
	for _, lesson := range lessons {
		words, err := findLessonVocabularies(db, lesson.ID)
		if err != nil {
			return fmt.Errorf("load quiz words for %s: %w", lesson.Slug, err)
		}
		if len(words) < 4 {
			continue
		}

		questionCount := minInt(foundationQuizQuestionsPerLesson, len(words))
		for questionIndex := 0; questionIndex < questionCount; questionIndex++ {
			word := words[questionIndex]
			topicID := lesson.TopicID
			vocabularyID := word.ID
			timeLimit := 45
			question := models.QuizQuestion{
				LessonID:         lesson.ID,
				TopicID:          &topicID,
				VocabularyID:     &vocabularyID,
				Type:             models.QuizQuestionMeaningChoice,
				Question:         fmt.Sprintf("\"%s\" nghĩa là gì?", word.Word),
				Explanation:      fmt.Sprintf("\"%s\" %s nghĩa là: %s. Ví dụ: %s", word.Word, word.IPA, word.MeaningVI, word.ExampleSentence),
				Points:           20,
				TimeLimitSeconds: &timeLimit,
				OrderIndex:       questionIndex + 1,
			}
			if err := db.Where("lesson_id = ? AND order_index = ?", lesson.ID, question.OrderIndex).Assign(question).FirstOrCreate(&question).Error; err != nil {
				return fmt.Errorf("seed foundation quiz question for %s: %w", word.Word, err)
			}

			// The correct answer moves between A-D so it is not always first.
			correctPosition := questionIndex % 4
			distractorOffsets := []int{1, 2, 3}
			labels := []string{"A", "B", "C", "D"}
			distractor := 0
			for position := 0; position < 4; position++ {
				content := word.MeaningVI
				isCorrect := position == correctPosition
				if !isCorrect {
					content = words[(questionIndex+distractorOffsets[distractor])%len(words)].MeaningVI
					distractor++
				}
				option := models.QuizOption{
					QuestionID: question.ID,
					Label:      labels[position],
					Content:    content,
					IsCorrect:  isCorrect,
					OrderIndex: position + 1,
				}
				if err := db.Where("question_id = ? AND order_index = ?", question.ID, option.OrderIndex).Assign(option).FirstOrCreate(&option).Error; err != nil {
					return fmt.Errorf("seed foundation quiz option for %s: %w", word.Word, err)
				}
			}
		}
	}

	return nil
}
