package course

import "ielts-learning/backend/internal/models"

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) List(userID uint) (ListResponse, error) {
	user, err := s.repository.FindUser(userID)
	if err != nil {
		return ListResponse{}, err
	}

	courses, err := s.repository.FindPublishedCourses()
	if err != nil {
		return ListResponse{}, err
	}

	completed, err := s.repository.CountCompletedLessonsByCourse(userID)
	if err != nil {
		return ListResponse{}, err
	}

	activeID := ActiveCourseID(user, courses)
	items := make([]CourseResponse, 0, len(courses))
	for _, course := range courses {
		items = append(items, toCourseResponse(course, completed[course.ID], course.ID == activeID))
	}

	return ListResponse{ActiveCourseID: activeID, Items: items}, nil
}

func (s Service) SetActiveCourse(userID uint, req SetActiveCourseRequest) (SetActiveCourseResponse, error) {
	course, err := s.repository.FindPublishedCourse(req.CourseID)
	if err != nil {
		return SetActiveCourseResponse{}, err
	}

	if err := s.repository.UpdateActiveCourse(userID, course.ID); err != nil {
		return SetActiveCourseResponse{}, err
	}

	return SetActiveCourseResponse{ActiveCourseID: course.ID}, nil
}

// ActiveCourseID returns the user's chosen course when it is still published,
// otherwise the first published course (courses must be in display order).
func ActiveCourseID(user models.User, courses []models.Course) uint {
	if user.ActiveCourseID != nil {
		for _, course := range courses {
			if course.ID == *user.ActiveCourseID {
				return course.ID
			}
		}
	}
	if len(courses) > 0 {
		return courses[0].ID
	}

	return 0
}

func toCourseResponse(course models.Course, completedLessons int, isActive bool) CourseResponse {
	return CourseResponse{
		ID:               course.ID,
		Title:            course.Title,
		Slug:             course.Slug,
		Description:      course.Description,
		Kind:             course.Kind,
		BandMin:          course.BandMin,
		BandMax:          course.BandMax,
		TotalWords:       course.TotalWords,
		TotalLessons:     course.TotalLessons,
		TotalTopics:      course.TotalTopics,
		CompletedLessons: completedLessons,
		IsActive:         isActive,
	}
}
