// Package foundation summarises a learner's progress through a foundation
// level (Task 76): per-skill progress, lessons, streak and the checklist of
// docs/foundation-spec.md §1.
package foundation

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

// CourseSlug is the foundation course.
const CourseSlug = "foundation-english"

var ErrLevelNotFound = errors.New("level not found")

// Targets of the A1 checklist (spec §1).
const (
	situationTarget  = 36
	soundTarget      = 36 // of the 40 sounds, at accuracy >= soundGoodScore
	soundGoodScore   = 70
	soundMinAttempts = 3
	paragraphTarget  = 3
	scoreGood        = 70
)

// introScenarios are the role-plays of goal 1 (self, work, family, routine).
var introScenarios = []string{"a1-w1-introducing-yourself", "a1-w1-describing-your-job", "a1-w1-family", "a1-w2-daily-schedule"}

// tenseTags are the grammar points whose mistakes count as tense mistakes.
var tenseTags = map[string]bool{
	"grammar:g1-to-be": true, "grammar:g3-present-simple": true, "grammar:g4-present-continuous": true,
	"grammar:g8-past-simple": true, "grammar:g9-future": true,
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return Service{db: db}
}

type lessonRow struct {
	ID         uint
	Title      string
	TopicTitle string
	Status     *string
	Checkpoint bool
}

type activityRow struct {
	Type      string
	BestScore *int
	Passed    bool
}

func (s Service) Summary(userID uint, level string) (SummaryResponse, error) {
	level = strings.ToUpper(strings.TrimSpace(level))
	if level == "" {
		level = "A1"
	}
	var band models.BandLevel
	err := s.db.Joins("JOIN courses c ON c.id = band_levels.course_id").
		Where("c.slug = ? AND band_levels.cefr_level = ?", CourseSlug, level).Take(&band).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return SummaryResponse{}, ErrLevelNotFound
	}
	if err != nil {
		return SummaryResponse{}, fmt.Errorf("find level: %w", err)
	}

	var user models.User
	if err := s.db.Select("id", "current_streak", "longest_streak").Take(&user, userID).Error; err != nil {
		return SummaryResponse{}, fmt.Errorf("find user: %w", err)
	}

	lessons, err := s.coreLessons(userID, band.ID)
	if err != nil {
		return SummaryResponse{}, err
	}
	response := SummaryResponse{Level: level, Title: band.Title, Streak: StreakSummary{Current: user.CurrentStreak, Longest: user.LongestStreak}}
	for _, lesson := range lessons {
		completed := lesson.Status != nil && *lesson.Status == string(models.LessonStatusCompleted)
		if lesson.Checkpoint {
			response.Lessons.CheckpointsTotal++
			if completed {
				response.Lessons.CheckpointsPassed++
			}
		} else {
			response.Lessons.Total++
			if completed {
				response.Lessons.Completed++
			}
		}
		if !completed && response.NextLesson == nil {
			response.NextLesson = &LessonLink{ID: lesson.ID, Title: lesson.Title, TopicTitle: lesson.TopicTitle}
		}
	}
	response.Lessons.Target = max(response.Lessons.Total, situationTarget)

	skills, sounds, err := s.skills(userID, band.ID, level)
	if err != nil {
		return SummaryResponse{}, err
	}
	response.Skills = skills

	if level == "A1" {
		goals, err := s.checklist(userID, response.Lessons, sounds, skills)
		if err != nil {
			return SummaryResponse{}, err
		}
		response.Checklist = goals
		for _, goal := range goals {
			if goal.Done {
				response.GoalsDone++
			}
		}
	} else {
		response.Checklist = []Goal{}
	}
	return response, nil
}

func (s Service) coreLessons(userID, bandID uint) ([]lessonRow, error) {
	var rows []lessonRow
	err := s.db.Raw(`SELECT l.id, l.title, t.title AS topic_title, p.status,
			EXISTS (SELECT 1 FROM lesson_activities a WHERE a.lesson_id = l.id AND a.type = ?) AS checkpoint
		FROM lessons l
		JOIN topics t ON t.id = l.topic_id AND t.deleted_at IS NULL
		LEFT JOIN user_lesson_progresses p ON p.lesson_id = l.id AND p.user_id = ?
		WHERE t.band_level_id = ? AND t.track = ? AND l.is_published = ? AND l.deleted_at IS NULL
		ORDER BY t.order_index, t.id, l.order_index, l.id`,
		models.ActivityCheckpoint, userID, bandID, models.TopicTrackCore, true).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("find lessons: %w", err)
	}
	return rows, nil
}

type soundStats struct {
	Good  int
	Total int
}

func (s Service) skills(userID, bandID uint, level string) ([]Skill, soundStats, error) {
	var skills []Skill

	// Vocabulary: words of the level's lessons.
	var vocab struct{ Total, Learned, Mastered int }
	err := s.db.Raw(`SELECT COUNT(DISTINCT lv.vocabulary_id) AS total,
			COUNT(DISTINCT CASE WHEN p.status IS NOT NULL AND p.status <> ? THEN lv.vocabulary_id END) AS learned,
			COUNT(DISTINCT CASE WHEN p.status = ? THEN lv.vocabulary_id END) AS mastered
		FROM lesson_vocabularies lv
		JOIN lessons l ON l.id = lv.lesson_id AND l.deleted_at IS NULL
		JOIN topics t ON t.id = l.topic_id AND t.deleted_at IS NULL
		LEFT JOIN user_vocabulary_progresses p ON p.vocabulary_id = lv.vocabulary_id AND p.user_id = ?
		WHERE t.band_level_id = ?`, models.VocabularyStatusNew, models.VocabularyStatusMastered, userID, bandID).Scan(&vocab).Error
	if err != nil {
		return nil, soundStats{}, fmt.Errorf("count vocabulary: %w", err)
	}
	skills = append(skills, Skill{
		Key: "VOCABULARY", Label: "Từ vựng", Value: percent(vocab.Mastered, vocab.Total), Practised: vocab.Learned > 0,
		Done: vocab.Mastered, Total: vocab.Total, Link: "/roadmap",
		Detail: fmt.Sprintf("%d/%d từ đã thuộc · %d từ đang học", vocab.Mastered, vocab.Total, vocab.Learned-vocab.Mastered),
	})

	// Grammar: mastery of the level's grammar points.
	var grammar []struct {
		Mastery   *int
		Practised bool
	}
	err = s.db.Raw(`SELECT p.mastery, p.id IS NOT NULL AS practised FROM grammar_points g
		LEFT JOIN user_grammar_progresses p ON p.grammar_point_id = g.id AND p.user_id = ?
		WHERE g.cefr_level = ? AND g.deleted_at IS NULL`, userID, level).Scan(&grammar).Error
	if err != nil {
		return nil, soundStats{}, fmt.Errorf("grammar mastery: %w", err)
	}
	sum, practised, strong := 0, 0, 0
	for _, row := range grammar {
		if row.Practised && row.Mastery != nil {
			sum += *row.Mastery
			practised++
			if *row.Mastery >= scoreGood {
				strong++
			}
		}
	}
	skills = append(skills, Skill{
		Key: "GRAMMAR", Label: "Ngữ pháp", Value: average(sum, practised), Practised: practised > 0,
		Done: strong, Total: len(grammar), Link: "/grammar",
		Detail: fmt.Sprintf("Đã luyện %d/%d điểm ngữ pháp · %d điểm thành thạo ≥ %d%%", practised, len(grammar), strong, scoreGood),
	})

	// Sounds: pronunciation accuracy of the 40 sound cards.
	var phonemes []struct {
		Accuracy *int
		Attempts *int
	}
	if err := s.db.Raw(`SELECT p.accuracy, p.attempts FROM phonemes ph
		LEFT JOIN user_phoneme_progresses p ON p.phoneme_id = ph.id AND p.user_id = ?`, userID).Scan(&phonemes).Error; err != nil {
		return nil, soundStats{}, fmt.Errorf("sound accuracy: %w", err)
	}
	sum, practised = 0, 0
	sounds := soundStats{Total: len(phonemes)}
	for _, row := range phonemes {
		if row.Accuracy == nil || row.Attempts == nil || *row.Attempts == 0 {
			continue
		}
		sum += *row.Accuracy
		practised++
		if *row.Attempts >= soundMinAttempts && *row.Accuracy >= soundGoodScore {
			sounds.Good++
		}
	}
	skills = append(skills, Skill{
		Key: "SOUNDS", Label: "Phát âm", Value: average(sum, practised), Practised: practised > 0,
		Done: sounds.Good, Total: sounds.Total, Link: "/pronunciation",
		Detail: fmt.Sprintf("Đã chấm %d/%d âm · %d âm đạt ≥ %d", practised, sounds.Total, sounds.Good, soundGoodScore),
	})

	// Listening, reading, writing, speaking: scores of the level's activities.
	var activities []activityRow
	err = s.db.Raw(`SELECT a.type, p.best_score, p.passed_at IS NOT NULL AS passed
		FROM lesson_activities a
		JOIN lessons l ON l.id = a.lesson_id AND l.is_published = ? AND l.deleted_at IS NULL
		JOIN topics t ON t.id = l.topic_id AND t.deleted_at IS NULL
		LEFT JOIN user_activity_progresses p ON p.activity_id = a.id AND p.user_id = ?
		WHERE t.band_level_id = ? AND t.track = ? AND a.status = ? AND a.type IN ?`,
		true, userID, bandID, models.TopicTrackCore, models.ContentPublished,
		[]models.ActivityType{models.ActivityListening, models.ActivityReading, models.ActivityWriting, models.ActivitySpeaking}).Scan(&activities).Error
	if err != nil {
		return nil, soundStats{}, fmt.Errorf("activity scores: %w", err)
	}
	for _, kind := range []struct{ key, label, link string }{
		{"LISTENING", "Nghe", "/roadmap"}, {"READING", "Đọc", "/roadmap"},
		{"WRITING", "Viết", "/roadmap"}, {"SPEAKING", "Nói", "/ai-conversation"},
	} {
		total, passed, attempted, scoreSum := 0, 0, 0, 0
		for _, row := range activities {
			if row.Type != kind.key {
				continue
			}
			total++
			if row.Passed {
				passed++
			}
			if row.BestScore != nil {
				attempted++
				scoreSum += *row.BestScore
			}
		}
		skills = append(skills, Skill{
			Key: kind.key, Label: kind.label, Value: average(scoreSum, attempted), Practised: attempted > 0,
			Done: passed, Total: total, Link: kind.link,
			Detail: fmt.Sprintf("Đạt %d/%d hoạt động · điểm trung bình %d", passed, total, average(scoreSum, attempted)),
		})
	}
	return skills, sounds, nil
}

func (s Service) checklist(userID uint, lessons LessonsSummary, sounds soundStats, skills []Skill) ([]Goal, error) {
	// 1. Introduce yourself, work, family, routine: the matching role-plays done.
	var introDone int64
	if err := s.db.Raw(`SELECT COUNT(DISTINCT sc.slug) FROM ai_conversation_sessions se
		JOIN ai_conversation_scenarios sc ON sc.id = se.scenario_id
		WHERE se.user_id = ? AND se.status = ? AND se.task_completed = ? AND sc.slug IN ?`,
		userID, models.SessionFinished, true, introScenarios).Scan(&introDone).Error; err != nil {
		return nil, fmt.Errorf("intro role-plays: %w", err)
	}

	// 4. Understand slow speech and short texts: listening + reading activities passed.
	var comprehension Skill
	for _, skill := range skills {
		if skill.Key == "LISTENING" || skill.Key == "READING" {
			comprehension.Done += skill.Done
			comprehension.Total += skill.Total
			comprehension.Value += skill.Value
		}
	}
	comprehension.Value /= 2

	// 5. A 5-7 sentence paragraph with no tense mistakes.
	paragraphs, err := s.cleanParagraphs(userID)
	if err != nil {
		return nil, err
	}

	comprehensionTarget := max(comprehension.Total, 1)
	goals := []Goal{
		{
			Key: "introduce", Title: "Giới thiệu bản thân, công việc, gia đình và thói quen hằng ngày bằng câu đơn giản",
			Done: int(introDone) >= len(introScenarios), Progress: percent(int(introDone), len(introScenarios)),
			Detail: fmt.Sprintf("Hoàn thành %d/%d hội thoại: giới thiệu bản thân, công việc, gia đình, lịch hằng ngày", introDone, len(introScenarios)),
			Hint:   "Làm hoạt động \"Nói chuyện với AI\" trong các bài này.", Link: "/ai-conversation",
		},
		{
			Key: "situations", Title: "Xử lý 36 tình huống hằng ngày (công sở, mua sắm, đi lại, sức khỏe, điện thoại)",
			Done: lessons.Completed >= situationTarget, Progress: percent(lessons.Completed, situationTarget),
			Detail: fmt.Sprintf("Đã xong %d/%d bài tình huống (hiện có %d bài)", lessons.Completed, situationTarget, lessons.Total),
			Hint:   "Học tiếp lộ trình, mỗi bài là một tình huống.", Link: "/roadmap",
		},
		{
			Key: "ipa", Title: "Đọc và phát âm bất kỳ từ nào từ IPA, kể cả âm cuối",
			Done: sounds.Good >= soundTarget, Progress: percent(sounds.Good, soundTarget),
			Detail: fmt.Sprintf("%d/%d âm đạt độ chính xác ≥ %d (cần %d)", sounds.Good, sounds.Total, soundGoodScore, soundTarget),
			Hint:   "Luyện nói trên thẻ âm và chế độ \"Nói theo\".", Link: "/pronunciation",
		},
		{
			Key: "comprehension", Title: "Hiểu lời nói chậm, rõ và văn bản ngắn (tin nhắn, thông báo, thực đơn, lịch)",
			Done:     comprehension.Total > 0 && comprehension.Done >= comprehension.Total && comprehension.Value >= scoreGood,
			Progress: percent(comprehension.Done, comprehensionTarget),
			Detail:   fmt.Sprintf("Đạt %d/%d hoạt động nghe và đọc · điểm trung bình %d", comprehension.Done, comprehension.Total, comprehension.Value),
			Hint:     "Làm hết phần Nghe và Đọc của mỗi bài.", Link: "/roadmap",
		},
		{
			Key: "writing", Title: "Viết đoạn 5-7 câu không sai thì",
			Done: paragraphs >= paragraphTarget, Progress: percent(paragraphs, paragraphTarget),
			Detail: fmt.Sprintf("%d/%d đoạn văn 5-7 câu không có lỗi thì", paragraphs, paragraphTarget),
			Hint:   "Làm hoạt động Viết trong bài học; lỗi thì là lỗi to be, hiện tại, quá khứ, tương lai.", Link: "/roadmap",
		},
	}
	return goals, nil
}

// cleanParagraphs counts completed writing submissions of 5+ sentences whose
// first AI report found no tense mistake.
func (s Service) cleanParagraphs(userID uint) (int, error) {
	var rows []models.WritingSubmission
	if err := s.db.Select("id", "text", "report_json").
		Where("user_id = ? AND status = ?", userID, models.WritingCompleted).Find(&rows).Error; err != nil {
		return 0, fmt.Errorf("find writing: %w", err)
	}
	count := 0
	for _, row := range rows {
		if sentences(row.Text) < 5 {
			continue
		}
		var report struct {
			Mistakes []struct {
				SkillTag string `json:"skillTag"`
			} `json:"mistakes"`
		}
		_ = json.Unmarshal(row.ReportJSON, &report)
		clean := true
		for _, mistake := range report.Mistakes {
			if tenseTags[mistake.SkillTag] {
				clean = false
				break
			}
		}
		if clean {
			count++
		}
	}
	return count, nil
}

func sentences(text string) int {
	count := 0
	for _, part := range strings.FieldsFunc(text, func(r rune) bool { return r == '.' || r == '!' || r == '?' }) {
		if strings.IndexFunc(part, func(r rune) bool { return r >= 'A' && r <= 'z' }) >= 0 {
			count++
		}
	}
	return count
}

func percent(done, total int) int {
	if total <= 0 {
		return 0
	}
	return min(100, int(math.Round(float64(done)*100/float64(total))))
}

func average(sum, count int) int {
	if count == 0 {
		return 0
	}
	return int(math.Round(float64(sum) / float64(count)))
}
