package course

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var (
	ErrCourseNotFound = errors.New("course not found")
	ErrUserNotFound   = errors.New("user not found")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindPublishedCourses() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.
		Where("is_published = ?", true).
		Order("order_index ASC, id ASC").
		Find(&courses).Error
	if err != nil {
		return nil, fmt.Errorf("find published courses: %w", err)
	}

	return courses, nil
}

func (r Repository) FindPublishedCourse(courseID uint) (models.Course, error) {
	var course models.Course
	err := r.db.Where("is_published = ?", true).First(&course, courseID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Course{}, ErrCourseNotFound
		}
		return models.Course{}, fmt.Errorf("find course: %w", err)
	}

	return course, nil
}

func (r Repository) FindUser(userID uint) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func (r Repository) UpdateActiveCourse(userID uint, courseID uint) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("active_course_id", courseID)
	if result.Error != nil {
		return fmt.Errorf("update active course: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

type completedLessonCount struct {
	CourseID uint
	Total    int
}

// CountCompletedLessonsByCourse returns completed lessons per course for a user.
func (r Repository) CountCompletedLessonsByCourse(userID uint) (map[uint]int, error) {
	var rows []completedLessonCount
	err := r.db.Model(&models.UserLessonProgress{}).
		Select("band_levels.course_id AS course_id, COUNT(*) AS total").
		Joins("JOIN lessons ON lessons.id = user_lesson_progresses.lesson_id AND lessons.deleted_at IS NULL").
		Joins("JOIN topics ON topics.id = lessons.topic_id AND topics.deleted_at IS NULL").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id AND band_levels.deleted_at IS NULL").
		Where("user_lesson_progresses.user_id = ? AND user_lesson_progresses.status = ?", userID, models.LessonStatusCompleted).
		Group("band_levels.course_id").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("count completed lessons by course: %w", err)
	}

	result := make(map[uint]int, len(rows))
	for _, row := range rows {
		result[row.CourseID] = row.Total
	}

	return result, nil
}
