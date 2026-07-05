package notification

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var ErrNotificationNotFound = errors.New("notification not found")

type CreateInput struct {
	UserID   uint
	Type     models.NotificationType
	Title    string
	Message  string
	Metadata map[string]string
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateNotification(userID uint, notificationType models.NotificationType, title string, message string, metadata map[string]string) error {
	return s.repository.db.Transaction(func(tx *gorm.DB) error {
		return s.CreateNotificationTx(tx, CreateInput{
			UserID:   userID,
			Type:     notificationType,
			Title:    title,
			Message:  message,
			Metadata: metadata,
		})
	})
}

func (s Service) CreateNotificationTx(tx *gorm.DB, input CreateInput) error {
	notification := models.Notification{
		UserID:    input.UserID,
		Type:      input.Type,
		Title:     strings.TrimSpace(input.Title),
		Body:      strings.TrimSpace(input.Message),
		ActionURL: metadataActionURL(input.Metadata),
		CreatedAt: nowUTC(),
	}

	return s.repository.Create(tx, notification)
}

func (s Service) ListNotifications(userID uint, page int, limit int) (ListNotificationsResponse, error) {
	page = normalizePage(page)
	limit = normalizeLimit(limit)

	items, total, err := s.repository.FindByUser(userID, page, limit)
	if err != nil {
		return ListNotificationsResponse{}, err
	}

	unreadCount, err := s.repository.CountUnread(userID)
	if err != nil {
		return ListNotificationsResponse{}, err
	}

	return ListNotificationsResponse{
		Items:       toNotificationItems(items),
		Page:        page,
		Limit:       limit,
		Total:       int(total),
		UnreadCount: int(unreadCount),
	}, nil
}

func (s Service) MarkAsRead(userID uint, notificationID uint) (MarkAsReadResponse, error) {
	notification, err := s.repository.FindByIDForUser(userID, notificationID)
	if err != nil {
		return MarkAsReadResponse{}, err
	}

	if notification.ReadAt == nil {
		readAt := nowUTC()
		notification.ReadAt = &readAt
		if err := s.repository.Save(&notification); err != nil {
			return MarkAsReadResponse{}, err
		}
	}

	return MarkAsReadResponse{
		Notification: toNotificationItem(notification),
	}, nil
}

func metadataActionURL(metadata map[string]string) string {
	if metadata == nil {
		return ""
	}

	return strings.TrimSpace(metadata["actionUrl"])
}

func normalizePage(page int) int {
	if page < 1 {
		return 1
	}

	return page
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}

	return limit
}

func toNotificationItems(notifications []models.Notification) []NotificationItem {
	items := make([]NotificationItem, 0, len(notifications))
	for _, notification := range notifications {
		items = append(items, toNotificationItem(notification))
	}

	return items
}

func toNotificationItem(notification models.Notification) NotificationItem {
	return NotificationItem{
		ID:        notification.ID,
		Type:      notification.Type,
		Title:     notification.Title,
		Body:      notification.Body,
		ActionURL: notification.ActionURL,
		ReadAt:    notification.ReadAt,
		CreatedAt: notification.CreatedAt,
	}
}
