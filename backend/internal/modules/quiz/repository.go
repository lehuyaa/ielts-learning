package quiz

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
	achievementmodule "ielts-learning/backend/internal/modules/achievement"
	xpmodule "ielts-learning/backend/internal/modules/xp"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrLessonNotFound     = errors.New("lesson not found")
	ErrQuizNotFound       = errors.New("quiz not found")
	ErrInvalidQuizAnswers = errors.New("invalid quiz answers")
)

type Repository struct {
	db                 *gorm.DB
	xpService          xpmodule.Service
	achievementService achievementmodule.Service
}

func NewRepository(db *gorm.DB) Repository {
	xpRepository := xpmodule.NewRepository(db)
	xpService := xpmodule.NewService(xpRepository)
	return Repository{
		db:                 db,
		xpService:          xpService,
		achievementService: achievementmodule.NewService(achievementmodule.NewRepository(db), xpService),
	}
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

func (r Repository) FindLesson(lessonID uint) (models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.
		Preload("Topic.BandLevel").
		Where("is_published = ?", true).
		First(&lesson, lessonID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Lesson{}, ErrLessonNotFound
		}
		return models.Lesson{}, fmt.Errorf("find lesson: %w", err)
	}

	return lesson, nil
}

func (r Repository) FindTopicLessons(topicID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.
		Where("topic_id = ? AND is_published = ?", topicID, true).
		Order("order_index ASC, id ASC").
		Find(&lessons).Error
	if err != nil {
		return nil, fmt.Errorf("find topic lessons: %w", err)
	}

	return lessons, nil
}

func (r Repository) FindLessonProgress(userID uint, lessonIDs []uint) ([]models.UserLessonProgress, error) {
	if len(lessonIDs) == 0 {
		return []models.UserLessonProgress{}, nil
	}

	var progress []models.UserLessonProgress
	err := r.db.
		Where("user_id = ? AND lesson_id IN ?", userID, lessonIDs).
		Find(&progress).Error
	if err != nil {
		return nil, fmt.Errorf("find lesson progress: %w", err)
	}

	return progress, nil
}

func (r Repository) FindQuizQuestions(lessonID uint) ([]models.QuizQuestion, error) {
	var questions []models.QuizQuestion
	err := r.db.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC, id ASC")
		}).
		Where("lesson_id = ?", lessonID).
		Order("order_index ASC, id ASC").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("find quiz questions: %w", err)
	}
	if len(questions) == 0 {
		return nil, ErrQuizNotFound
	}

	return questions, nil
}

func (r Repository) SubmitQuiz(user models.User, lesson models.Lesson, questions []models.QuizQuestion, now time.Time, graded GradedQuiz) (SubmitResponse, error) {
	var response SubmitResponse

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var progress models.UserLessonProgress
		progressExists := true
		wasCompleted := false

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND lesson_id = ?", user.ID, lesson.ID).
			First(&progress).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("find quiz lesson progress: %w", err)
			}
			progressExists = false
			progress = models.UserLessonProgress{
				UserID:    user.ID,
				LessonID:  lesson.ID,
				Status:    models.LessonStatusInProgress,
				StartedAt: &now,
			}
		} else {
			wasCompleted = progress.Status == models.LessonStatusCompleted
		}

		session := models.QuizSession{
			UserID:               user.ID,
			LessonID:             lesson.ID,
			Status:               models.QuizSessionCompleted,
			CurrentQuestionIndex: len(questions),
			TotalQuestions:       len(questions),
			Points:               graded.EarnedPoints,
			CorrectAnswers:       graded.CorrectCount,
			WrongAnswers:         len(questions) - graded.CorrectCount,
			StartedAt:            now,
			FinishedAt:           &now,
			CreatedAt:            now,
			UpdatedAt:            now,
		}
		if err := tx.Create(&session).Error; err != nil {
			return fmt.Errorf("create quiz session: %w", err)
		}

		for _, result := range graded.Results {
			answer := models.QuizSessionAnswer{
				QuizSessionID:    session.ID,
				QuestionID:       result.QuestionID,
				SelectedOptionID: result.SelectedOptionID,
				CorrectOptionID:  &result.CorrectOptionID,
				IsCorrect:        result.IsCorrect,
				PointsAwarded:    result.PointsAwarded,
				AnsweredAt:       now,
				CreatedAt:        now,
				UpdatedAt:        now,
			}
			if err := tx.Create(&answer).Error; err != nil {
				return fmt.Errorf("create quiz answer: %w", err)
			}
		}

		newlyCompleted := graded.Passed && !wasCompleted
		xpAwarded := 0
		if newlyCompleted {
			xpAwarded = lesson.XPReward
		}

		attempt := models.UserQuizAttempt{
			UserID:         user.ID,
			LessonID:       lesson.ID,
			QuizSessionID:  &session.ID,
			Score:          graded.Score,
			Points:         graded.EarnedPoints,
			TotalQuestions: len(questions),
			CorrectAnswers: graded.CorrectCount,
			XPEarned:       xpAwarded,
			Passed:         graded.Passed,
			StartedAt:      &now,
			FinishedAt:     &now,
			CreatedAt:      now,
		}
		if err := tx.Create(&attempt).Error; err != nil {
			return fmt.Errorf("create quiz attempt: %w", err)
		}

		score := graded.Score
		progress.Score = &score
		if progress.BestScore == nil || score > *progress.BestScore {
			progress.BestScore = &score
		}
		if progress.StartedAt == nil {
			progress.StartedAt = &now
		}
		progress.LastStudiedAt = &now
		if progress.TotalWords == 0 {
			progress.TotalWords = countLessonVocabulary(tx, lesson.ID)
		}

		if graded.Passed {
			progress.Status = models.LessonStatusCompleted
			if progress.CompletedAt == nil {
				progress.CompletedAt = &now
			}
			if xpAwarded > progress.BestXP {
				progress.BestXP = xpAwarded
			}
		} else if progress.Status != models.LessonStatusCompleted {
			progress.Status = models.LessonStatusInProgress
		}

		if progressExists {
			if err := tx.Save(&progress).Error; err != nil {
				return fmt.Errorf("update quiz lesson progress: %w", err)
			}
		} else if err := tx.Create(&progress).Error; err != nil {
			return fmt.Errorf("create quiz lesson progress: %w", err)
		}

		if newlyCompleted {
			if err := unlockNextLesson(tx, user.ID, lesson, now); err != nil {
				return err
			}
		}

		if xpAwarded > 0 {
			sourceID := lesson.ID
			if _, err := r.xpService.AwardXP(tx, xpmodule.AwardInput{
				UserID:           user.ID,
				SourceType:       xpmodule.EventLessonCompleted,
				SourceID:         &sourceID,
				XP:               xpAwarded,
				Description:      fmt.Sprintf("Passed quiz for %s", lesson.Title),
				AwardedAt:        now,
				PreventDuplicate: true,
				TouchLastActive:  true,
			}); err != nil {
				return fmt.Errorf("award quiz completion xp: %w", err)
			}
		}

		if err := upsertDailyActivity(tx, user.ID, activityDate(now, user.Timezone), newlyCompleted, xpAwarded, now); err != nil {
			return err
		}

		if _, err := r.achievementService.CheckAndUnlockAchievementsTx(tx, user.ID); err != nil {
			return fmt.Errorf("check achievements after quiz submit: %w", err)
		}

		response = SubmitResponse{
			AttemptID:      attempt.ID,
			LessonID:       lesson.ID,
			Score:          graded.Score,
			RequiredScore:  lesson.RequiredScore,
			Passed:         graded.Passed,
			CorrectCount:   graded.CorrectCount,
			TotalQuestions: len(questions),
			EarnedXP:       xpAwarded,
			CompletedAt:    progress.CompletedAt,
			Results:        graded.Results,
		}

		return nil
	})
	if err != nil {
		return SubmitResponse{}, err
	}

	return response, nil
}

func countLessonVocabulary(tx *gorm.DB, lessonID uint) int {
	var count int64
	if err := tx.Model(&models.LessonVocabulary{}).
		Where("lesson_id = ? AND is_required = ?", lessonID, true).
		Count(&count).Error; err != nil {
		return 0
	}
	return int(count)
}

func unlockNextLesson(tx *gorm.DB, userID uint, lesson models.Lesson, now time.Time) error {
	var nextLesson models.Lesson
	err := tx.
		Where("topic_id = ? AND is_published = ? AND (order_index > ? OR (order_index = ? AND id > ?))", lesson.TopicID, true, lesson.OrderIndex, lesson.OrderIndex, lesson.ID).
		Order("order_index ASC, id ASC").
		First(&nextLesson).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("find next lesson: %w", err)
	}

	var progress models.UserLessonProgress
	err = tx.
		Where("user_id = ? AND lesson_id = ?", userID, nextLesson.ID).
		First(&progress).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("find next lesson progress: %w", err)
	}

	progress = models.UserLessonProgress{
		UserID:        userID,
		LessonID:      nextLesson.ID,
		Status:        models.LessonStatusUnlocked,
		LastStudiedAt: &now,
	}
	if err := tx.Create(&progress).Error; err != nil {
		return fmt.Errorf("unlock next lesson: %w", err)
	}

	return nil
}

func upsertDailyActivity(tx *gorm.DB, userID uint, date time.Time, lessonNewlyCompleted bool, xpAwarded int, now time.Time) error {
	lessonsDone := 0
	if lessonNewlyCompleted {
		lessonsDone = 1
	}

	activity := models.DailyActivity{
		UserID:       userID,
		Date:         date,
		QuizzesTaken: 1,
		LessonsDone:  lessonsDone,
		XPEarned:     xpAwarded,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "date"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"quizzes_taken": gorm.Expr("quizzes_taken + ?", 1),
			"lessons_done":  gorm.Expr("lessons_done + ?", lessonsDone),
			"xp_earned":     gorm.Expr("xp_earned + ?", xpAwarded),
			"updated_at":    now,
		}),
	}).Create(&activity).Error
	if err != nil {
		return fmt.Errorf("upsert quiz daily activity: %w", err)
	}

	return nil
}

func activityDate(now time.Time, timezone string) time.Time {
	location := time.UTC
	if timezone != "" {
		if loadedLocation, err := time.LoadLocation(timezone); err == nil {
			location = loadedLocation
		}
	}

	local := now.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
}
