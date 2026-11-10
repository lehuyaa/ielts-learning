package lessonactivity

import (
	"testing"
	"time"

	"ielts-learning/backend/internal/models"
)

func TestActivityLocks(t *testing.T) {
	passedAt := time.Now()
	activities := []models.LessonActivity{
		{ID: 1, IsRequired: true},
		{ID: 2, IsRequired: false},
		{ID: 3, IsRequired: true},
		{ID: 4, IsRequired: true},
	}

	tests := []struct {
		name     string
		progress map[uint]models.UserActivityProgress
		want     map[uint]bool
	}{
		{
			name:     "nothing done: only the first is open",
			progress: map[uint]models.UserActivityProgress{},
			want:     map[uint]bool{1: false, 2: true, 3: true, 4: true},
		},
		{
			name: "optional activities do not block",
			progress: map[uint]models.UserActivityProgress{
				1: {Status: models.ActivityPassed, PassedAt: &passedAt},
			},
			want: map[uint]bool{1: false, 2: false, 3: false, 4: true},
		},
		{
			name: "retrying a passed activity keeps the next ones open",
			progress: map[uint]models.UserActivityProgress{
				1: {Status: models.ActivityInProgress, PassedAt: &passedAt},
				3: {Status: models.ActivityFailed},
			},
			want: map[uint]bool{1: false, 2: false, 3: false, 4: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := activityLocks(activities, tt.progress)
			for id, want := range tt.want {
				if got[id] != want {
					t.Errorf("activity %d locked = %v, want %v", id, got[id], want)
				}
			}
		})
	}
}
