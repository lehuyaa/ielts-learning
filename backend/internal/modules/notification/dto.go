package notification

import (
	"time"

	"ielts-learning/backend/internal/models"
)

type ListNotificationsResponse struct {
	Items       []NotificationItem `json:"items"`
	Page        int                `json:"page"`
	Limit       int                `json:"limit"`
	Total       int                `json:"total"`
	UnreadCount int                `json:"unreadCount"`
}

type NotificationItem struct {
	ID        uint                    `json:"id"`
	Type      models.NotificationType `json:"type"`
	Title     string                  `json:"title"`
	Body      string                  `json:"body"`
	ActionURL string                  `json:"actionUrl"`
	ReadAt    *time.Time              `json:"readAt"`
	CreatedAt time.Time               `json:"createdAt"`
}

type MarkAsReadResponse struct {
	Notification NotificationItem `json:"notification"`
}
