package lessonactivity

import (
	"encoding/json"
	"hash/fnv"
	"math/rand"
	"sort"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/modules/mistakes"
)

// checkpointItemCount is content.itemCount of a CHECKPOINT activity: each
// attempt draws that many items from the activity's pool (Task 75). 0 means
// every exercise is used.
func checkpointItemCount(activity models.LessonActivity) int {
	if activity.Type != models.ActivityCheckpoint || len(activity.ContentJSON) == 0 {
		return 0
	}
	var content struct {
		ItemCount int `json:"itemCount"`
	}
	_ = json.Unmarshal(activity.ContentJSON, &content)
	return max(content.ItemCount, 0)
}

// drawAttempt is the attempt the next answer belongs to: the open one, or a
// new one.
func drawAttempt(progress models.UserActivityProgress, hasProgress bool) int {
	if attempt := currentAttempt(progress, hasProgress); attempt > 0 {
		return attempt
	}
	if hasProgress {
		return progress.Attempts + 1
	}
	return 1
}

// attemptExercises returns the exercises of an activity for one attempt:
// all of them, or a checkpoint's draw for that attempt.
func (s Service) attemptExercises(userID uint, activity models.LessonActivity, attempt int) ([]models.Exercise, error) {
	exercises, err := s.repository.FindExercises(activity.ID)
	if err != nil {
		return nil, err
	}
	count := checkpointItemCount(activity)
	if count == 0 || count >= len(exercises) {
		return exercises, nil
	}
	return drawExercises(exercises, count, drawSeed(userID, activity.ID, attempt)), nil
}

func drawSeed(userID, activityID uint, attempt int) int64 {
	h := fnv.New64a()
	_ = json.NewEncoder(h).Encode([]any{userID, activityID, attempt})
	return int64(h.Sum64() & (1<<62 - 1))
}

// drawExercises picks count items, the same for the same seed, spread over
// the skills of the pool (grammar, listening, sounds, vocabulary…) so a
// checkpoint tests the whole week.
func drawExercises(pool []models.Exercise, count int, seed int64) []models.Exercise {
	random := rand.New(rand.NewSource(seed))
	groups := map[string][]models.Exercise{}
	for _, exercise := range pool {
		kind := mistakes.SkillKind(exercise.SkillTag)
		groups[kind] = append(groups[kind], exercise)
	}
	kinds := make([]string, 0, len(groups))
	for kind := range groups {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	for _, kind := range kinds {
		items := groups[kind]
		random.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
	}

	picked := make([]models.Exercise, 0, count)
	for round := 0; len(picked) < count; round++ {
		added := false
		for _, kind := range kinds {
			if round < len(groups[kind]) && len(picked) < count {
				picked = append(picked, groups[kind][round])
				added = true
			}
		}
		if !added {
			break
		}
	}
	random.Shuffle(len(picked), func(i, j int) { picked[i], picked[j] = picked[j], picked[i] })
	for i := range picked {
		picked[i].OrderIndex = i + 1
	}
	return picked
}
