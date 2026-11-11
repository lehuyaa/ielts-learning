package foundation

// SummaryResponse is the learner's progress through one foundation level
// (Task 76): skills, lessons, streak and the "level reached" checklist.
type SummaryResponse struct {
	Level      string         `json:"level"`
	Title      string         `json:"title"`
	Lessons    LessonsSummary `json:"lessons"`
	Streak     StreakSummary  `json:"streak"`
	Skills     []Skill        `json:"skills"`
	Checklist  []Goal         `json:"checklist"`
	GoalsDone  int            `json:"goalsDone"`
	NextLesson *LessonLink    `json:"nextLesson"`
}

type LessonsSummary struct {
	// Situation lessons (CORE, without the weekly checkpoints).
	Completed int `json:"completed"`
	Total     int `json:"total"`
	// Target is the planned number of situation lessons for the level.
	Target            int `json:"target"`
	CheckpointsPassed int `json:"checkpointsPassed"`
	CheckpointsTotal  int `json:"checkpointsTotal"`
}

type StreakSummary struct {
	Current int `json:"current"`
	Longest int `json:"longest"`
}

// Skill is one skill bar. Value is the quality of what was practised
// (0-100); Done / Total is how much of the level was covered.
type Skill struct {
	Key       string `json:"key"`
	Label     string `json:"label"`
	Value     int    `json:"value"`
	Practised bool   `json:"practised"`
	Done      int    `json:"done"`
	Total     int    `json:"total"`
	Detail    string `json:"detail"`
	Link      string `json:"link"`
}

// Goal is one line of the "A1 reached" checklist (spec §1).
type Goal struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Done     bool   `json:"done"`
	Progress int    `json:"progress"`
	Detail   string `json:"detail"`
	Hint     string `json:"hint"`
	Link     string `json:"link"`
}

type LessonLink struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	TopicTitle string `json:"topicTitle"`
}
