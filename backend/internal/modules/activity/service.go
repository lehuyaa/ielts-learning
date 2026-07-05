package activity

import (
	"time"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

type Event struct {
	WordsLearnedDelta  int
	WordsReviewedDelta int
	QuizzesTakenDelta  int
	LessonsDoneDelta   int
	ActiveMinutesDelta int
	XPEarnedDelta      int
	OccurredAt         time.Time
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) RecordActivity(userID uint, event Event) error {
	return s.repository.db.Transaction(func(tx *gorm.DB) error {
		return s.RecordActivityTx(tx, userID, event)
	})
}

func (s Service) RecordActivityTx(tx *gorm.DB, userID uint, event Event) error {
	user, err := s.repository.FindUserForUpdate(tx, userID)
	if err != nil {
		return err
	}

	occurredAt := normalizeOccurredAt(event.OccurredAt)
	activityDate, location := activityDateForTimezone(occurredAt, user.Timezone)
	yesterdayDate := activityDate.AddDate(0, 0, -1)

	_, hadActivityToday, err := s.repository.FindDailyActivityByDate(tx, userID, activityDate)
	if err != nil {
		return err
	}

	hadYesterdayActivity := false
	if !hadActivityToday {
		_, hadYesterdayActivity, err = s.repository.FindDailyActivityByDate(tx, userID, yesterdayDate)
		if err != nil {
			return err
		}
	}

	activity := models.DailyActivity{
		UserID:        userID,
		Date:          activityDate,
		WordsLearned:  event.WordsLearnedDelta,
		WordsReviewed: event.WordsReviewedDelta,
		QuizzesTaken:  event.QuizzesTakenDelta,
		LessonsDone:   event.LessonsDoneDelta,
		ActiveMinutes: event.ActiveMinutesDelta,
		XPEarned:      event.XPEarnedDelta,
		CreatedAt:     occurredAt,
		UpdatedAt:     occurredAt,
	}
	if err := s.repository.UpsertDailyActivity(tx, activity, event); err != nil {
		return err
	}

	if !hadActivityToday {
		if hadYesterdayActivity {
			if user.CurrentStreak > 0 {
				user.CurrentStreak++
			} else {
				user.CurrentStreak = 1
			}
		} else {
			user.CurrentStreak = 1
		}
		if user.CurrentStreak > user.LongestStreak {
			user.LongestStreak = user.CurrentStreak
		}
	}

	lastActiveAt := occurredAt.In(location).UTC()
	user.LastActiveAt = &lastActiveAt

	return s.repository.SaveUser(tx, &user)
}

func normalizeOccurredAt(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}

	return value.UTC()
}

func activityDateForTimezone(now time.Time, timezone string) (time.Time, *time.Location) {
	location := time.UTC
	if timezone != "" {
		if loadedLocation, err := time.LoadLocation(timezone); err == nil {
			location = loadedLocation
		}
	}

	local := now.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location), location
}
