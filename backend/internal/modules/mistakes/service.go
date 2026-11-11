package mistakes

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ielts-learning/backend/internal/models"
)

// FixedStreak is how many correct reviews in a row mark a mistake FIXED.
const FixedStreak = 3

var (
	ErrNotFound      = errors.New("mistake not found")
	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidRating = errors.New("invalid rating")
	ErrNotReviewing  = errors.New("mistake is not in the review queue")
)

type Service struct {
	db  *gorm.DB
	now func() time.Time
}

func NewService(db *gorm.DB) Service {
	return Service{db: db, now: time.Now}
}

// List returns the mistakes of a status filter (open = OPEN + REVIEWING,
// fixed, all), grouped by skill tag, most frequent first.
func (s Service) List(userID uint, filter string) (ListResponse, error) {
	var rows []models.LearnerMistake
	query := s.db.Where("user_id = ?", userID)
	switch strings.ToLower(filter) {
	case "", "open":
		query = query.Where("status IN ?", []models.MistakeStatus{models.MistakeOpen, models.MistakeReviewing})
	case "fixed":
		query = query.Where("status = ?", models.MistakeFixed)
	case "all":
	default:
		return ListResponse{}, ErrInvalidStatus
	}
	if err := query.Order("occurrences DESC, last_seen_at DESC").Find(&rows).Error; err != nil {
		return ListResponse{}, fmt.Errorf("find mistakes: %w", err)
	}
	summary, err := s.summary(userID)
	if err != nil {
		return ListResponse{}, err
	}
	titles, err := s.grammarTitles()
	if err != nil {
		return ListResponse{}, err
	}

	byTag := map[string]*Group{}
	var order []string
	for _, row := range rows {
		group, ok := byTag[row.SkillTag]
		if !ok {
			group = &Group{SkillTag: row.SkillTag, Label: SkillLabel(row.SkillTag, titles), Kind: SkillKind(row.SkillTag), Practiceable: Practiceable(row.SkillTag), Items: []Item{}}
			byTag[row.SkillTag] = group
			order = append(order, row.SkillTag)
		}
		group.Count++
		group.Occurrences += row.Occurrences
		if row.LastSeenAt.After(group.LastSeenAt) {
			group.LastSeenAt = row.LastSeenAt
		}
		group.Items = append(group.Items, toItem(row, group.Label))
	}
	groups := make([]Group, 0, len(order))
	for _, tag := range order {
		groups = append(groups, *byTag[tag])
	}
	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].Occurrences != groups[j].Occurrences {
			return groups[i].Occurrences > groups[j].Occurrences
		}
		return groups[i].LastSeenAt.After(groups[j].LastSeenAt)
	})
	return ListResponse{Summary: summary, Groups: groups}, nil
}

func (s Service) summary(userID uint) (Summary, error) {
	var counts []struct {
		Status models.MistakeStatus
		Total  int
	}
	if err := s.db.Model(&models.LearnerMistake{}).Select("status, COUNT(*) AS total").
		Where("user_id = ?", userID).Group("status").Scan(&counts).Error; err != nil {
		return Summary{}, fmt.Errorf("count mistakes: %w", err)
	}
	var summary Summary
	for _, count := range counts {
		switch count.Status {
		case models.MistakeOpen:
			summary.Open = count.Total
		case models.MistakeReviewing:
			summary.Reviewing = count.Total
		case models.MistakeFixed:
			summary.Fixed = count.Total
		}
	}
	var due int64
	if err := s.dueQuery(userID).Model(&models.LearnerMistake{}).Count(&due).Error; err != nil {
		return Summary{}, fmt.Errorf("count due mistakes: %w", err)
	}
	summary.Due = int(due)
	return summary, nil
}

func (s Service) dueQuery(userID uint) *gorm.DB {
	return s.db.Where("user_id = ? AND status = ? AND (next_review_at IS NULL OR next_review_at <= ?)",
		userID, models.MistakeReviewing, s.now())
}

func (s Service) grammarTitles() (map[string]string, error) {
	var points []models.GrammarPoint
	if err := s.db.Select("slug", "title", "title_vi").Find(&points).Error; err != nil {
		return nil, fmt.Errorf("find grammar points: %w", err)
	}
	titles := make(map[string]string, len(points))
	for _, point := range points {
		titles[point.Slug] = firstNonEmpty(point.TitleVI, point.Title)
	}
	return titles, nil
}

// Update moves a mistake into or out of the review queue, or marks it fixed.
func (s Service) Update(userID, id uint, req UpdateRequest) (Item, error) {
	now := s.now()
	var row models.LearnerMistake
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		row, err = findForUpdate(tx, userID, id)
		if err != nil {
			return err
		}
		updates := map[string]any{"status": req.Status, "updated_at": now}
		switch req.Status {
		case models.MistakeReviewing:
			if row.Status != models.MistakeReviewing {
				updates["review_streak"] = 0
				updates["interval_days"] = 0
				updates["next_review_at"] = now
				updates["fixed_at"] = nil
			}
		case models.MistakeFixed:
			updates["fixed_at"] = now
			updates["next_review_at"] = nil
		case models.MistakeOpen:
			updates["next_review_at"] = nil
			updates["fixed_at"] = nil
			updates["review_streak"] = 0
		default:
			return ErrInvalidStatus
		}
		if err := tx.Model(&row).Updates(updates).Error; err != nil {
			return fmt.Errorf("update mistake: %w", err)
		}
		return tx.Take(&row, row.ID).Error
	})
	if err != nil {
		return Item{}, err
	}
	return s.item(row)
}

// ReviewQueue returns the review cards due now (oldest first).
func (s Service) ReviewQueue(userID uint, limit int) (ReviewQueueResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	var rows []models.LearnerMistake
	if err := s.dueQuery(userID).Order("next_review_at IS NULL DESC, next_review_at ASC, id ASC").Limit(limit).Find(&rows).Error; err != nil {
		return ReviewQueueResponse{}, fmt.Errorf("find due mistakes: %w", err)
	}
	var total int64
	if err := s.dueQuery(userID).Model(&models.LearnerMistake{}).Count(&total).Error; err != nil {
		return ReviewQueueResponse{}, err
	}
	titles, err := s.grammarTitles()
	if err != nil {
		return ReviewQueueResponse{}, err
	}
	cards := make([]Item, 0, len(rows))
	for _, row := range rows {
		cards = append(cards, toItem(row, SkillLabel(row.SkillTag, titles)))
	}
	return ReviewQueueResponse{Cards: cards, Count: int(total)}, nil
}

// Review rates one card: GOOD/EASY extend the interval and the streak,
// HARD keeps the streak, AGAIN starts over. FixedStreak correct reviews in a
// row mark the mistake FIXED.
func (s Service) Review(userID, id uint, rating models.FlashcardRating) (ReviewResponse, error) {
	now := s.now()
	var row models.LearnerMistake
	fixed := false
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		row, err = findForUpdate(tx, userID, id)
		if err != nil {
			return err
		}
		if row.Status != models.MistakeReviewing {
			return ErrNotReviewing
		}
		streak, interval, next, ok := schedule(row.ReviewStreak, row.IntervalDays, rating, now)
		if !ok {
			return ErrInvalidRating
		}
		updates := map[string]any{"review_streak": streak, "interval_days": interval, "next_review_at": next, "updated_at": now}
		if streak >= FixedStreak {
			fixed = true
			updates["status"] = models.MistakeFixed
			updates["fixed_at"] = now
			updates["next_review_at"] = nil
		}
		if err := tx.Model(&row).Updates(updates).Error; err != nil {
			return fmt.Errorf("review mistake: %w", err)
		}
		return tx.Take(&row, row.ID).Error
	})
	if err != nil {
		return ReviewResponse{}, err
	}
	item, err := s.item(row)
	if err != nil {
		return ReviewResponse{}, err
	}
	return ReviewResponse{Item: item, Fixed: fixed}, nil
}

// schedule is the small spaced-repetition step for mistake cards.
func schedule(streak, interval int, rating models.FlashcardRating, now time.Time) (int, int, time.Time, bool) {
	switch rating {
	case models.FlashcardRatingAgain:
		return 0, 0, now.Add(10 * time.Minute), true
	case models.FlashcardRatingHard:
		interval = max(1, interval)
		return streak, interval, now.AddDate(0, 0, interval), true
	case models.FlashcardRatingGood:
		interval = max(1, interval*2+1)
		return streak + 1, interval, now.AddDate(0, 0, interval), true
	case models.FlashcardRatingEasy:
		interval = max(3, interval*3+1)
		return streak + 1, interval, now.AddDate(0, 0, interval), true
	default:
		return streak, interval, now, false
	}
}

func findForUpdate(tx *gorm.DB, userID, id uint) (models.LearnerMistake, error) {
	var row models.LearnerMistake
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", id, userID).Take(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, ErrNotFound
	}
	if err != nil {
		return row, fmt.Errorf("find mistake: %w", err)
	}
	return row, nil
}

func (s Service) item(row models.LearnerMistake) (Item, error) {
	titles, err := s.grammarTitles()
	if err != nil {
		return Item{}, err
	}
	return toItem(row, SkillLabel(row.SkillTag, titles)), nil
}

func toItem(row models.LearnerMistake, label string) Item {
	return Item{
		ID:            row.ID,
		Source:        row.Source,
		SkillTag:      row.SkillTag,
		Label:         label,
		WrongText:     row.WrongText,
		CorrectText:   row.CorrectText,
		ExplanationVI: row.ExplanationVI,
		Occurrences:   row.Occurrences,
		Status:        row.Status,
		ReviewStreak:  row.ReviewStreak,
		NextReviewAt:  row.NextReviewAt,
		FixedAt:       row.FixedAt,
		LastSeenAt:    row.LastSeenAt,
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
