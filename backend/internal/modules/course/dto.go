package course

import "ielts-learning/backend/internal/models"

type ListResponse struct {
	ActiveCourseID uint             `json:"activeCourseId"`
	Items          []CourseResponse `json:"items"`
}

type CourseResponse struct {
	ID               uint              `json:"id"`
	Title            string            `json:"title"`
	Slug             string            `json:"slug"`
	Description      string            `json:"description"`
	Kind             models.CourseKind `json:"kind"`
	BandMin          float64           `json:"bandMin"`
	BandMax          float64           `json:"bandMax"`
	TotalWords       int               `json:"totalWords"`
	TotalLessons     int               `json:"totalLessons"`
	TotalTopics      int               `json:"totalTopics"`
	CompletedLessons int               `json:"completedLessons"`
	IsActive         bool              `json:"isActive"`
}

type SetActiveCourseRequest struct {
	CourseID uint `json:"courseId"`
}

type SetActiveCourseResponse struct {
	ActiveCourseID uint `json:"activeCourseId"`
}
