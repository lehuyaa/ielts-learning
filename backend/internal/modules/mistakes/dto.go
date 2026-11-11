package mistakes

import (
	"time"

	"ielts-learning/backend/internal/models"
)

// ListResponse is the mistake log grouped by skill tag (Task 73).
type ListResponse struct {
	Summary Summary `json:"summary"`
	Groups  []Group `json:"groups"`
}

type Summary struct {
	Open      int `json:"open"`
	Reviewing int `json:"reviewing"`
	Fixed     int `json:"fixed"`
	// Due is how many review cards are due now.
	Due int `json:"due"`
}

type Group struct {
	SkillTag string `json:"skillTag"`
	Label    string `json:"label"`
	Kind     string `json:"kind"`
	// Practiceable: POST /practice/generate can build 10 exercises for it.
	Practiceable bool      `json:"practiceable"`
	Count        int       `json:"count"`
	Occurrences  int       `json:"occurrences"`
	LastSeenAt   time.Time `json:"lastSeenAt"`
	Items        []Item    `json:"items"`
}

type Item struct {
	ID            uint                 `json:"id"`
	Source        models.MistakeSource `json:"source"`
	SkillTag      string               `json:"skillTag"`
	Label         string               `json:"label"`
	WrongText     string               `json:"wrongText"`
	CorrectText   string               `json:"correctText"`
	ExplanationVI string               `json:"explanationVi"`
	Occurrences   int                  `json:"occurrences"`
	Status        models.MistakeStatus `json:"status"`
	ReviewStreak  int                  `json:"reviewStreak"`
	NextReviewAt  *time.Time           `json:"nextReviewAt"`
	FixedAt       *time.Time           `json:"fixedAt"`
	LastSeenAt    time.Time            `json:"lastSeenAt"`
}

// UpdateRequest changes a mistake's status: REVIEWING adds it to the review
// queue, FIXED marks it learnt, OPEN takes it out of the queue.
type UpdateRequest struct {
	Status models.MistakeStatus `json:"status"`
}

type ReviewQueueResponse struct {
	Cards []Item `json:"cards"`
	Count int    `json:"count"`
}

type ReviewRequest struct {
	Rating models.FlashcardRating `json:"rating"`
}

type ReviewResponse struct {
	Item Item `json:"item"`
	// Fixed is true when this review completed the streak.
	Fixed bool `json:"fixed"`
}
