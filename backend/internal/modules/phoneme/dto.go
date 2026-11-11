package phoneme

import (
	"time"

	"ielts-learning/backend/internal/modules/lessonactivity"
)

type ListResponse struct {
	Items []PhonemeItem `json:"items"`
	// PracticeTopics are always-open topics with pronunciation practice (Sound Lab).
	PracticeTopics []PracticeTopic `json:"practiceTopics"`
}

type PracticeTopic struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	LessonCount int    `json:"lessonCount"`
}

// PhonemeItem is one cell of the IPA chart.
type PhonemeItem struct {
	ID          uint             `json:"id"`
	Symbol      string           `json:"symbol"`
	Kind        string           `json:"kind"`
	Voiced      *bool            `json:"voiced"`
	NameVI      string           `json:"nameVi"`
	AudioURL    string           `json:"audioUrl"`
	ExampleWord *ExampleWord     `json:"exampleWord"`
	OrderIndex  int              `json:"orderIndex"`
	Progress    ProgressResponse `json:"progress"`
}

type ExampleWord struct {
	Word     string `json:"word"`
	IPA      string `json:"ipa"`
	AudioURL string `json:"audioUrl,omitempty"`
}

// ProgressResponse is the learner's pronunciation accuracy for a sound
// (filled by pronunciation scoring, Task 71).
type ProgressResponse struct {
	Accuracy        int        `json:"accuracy"`
	Attempts        int        `json:"attempts"`
	Practised       bool       `json:"practised"`
	LastPracticedAt *time.Time `json:"lastPracticedAt"`
}

type DetailResponse struct {
	Phoneme  lessonactivity.PhonemeCard `json:"phoneme"`
	Progress ProgressResponse           `json:"progress"`
	// Lessons whose pronunciation activity teaches this sound.
	Lessons []LessonLink `json:"lessons"`
	// Previous and Next are neighbours on the chart, for paging.
	Previous *string `json:"previous"`
	Next     *string `json:"next"`
}

type LessonLink struct {
	LessonID      uint   `json:"lessonId"`
	LessonTitle   string `json:"lessonTitle"`
	TopicTitle    string `json:"topicTitle"`
	ActivityID    uint   `json:"activityId"`
	ActivityTitle string `json:"activityTitle"`
}
