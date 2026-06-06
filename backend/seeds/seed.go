package seeds

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

const vocabularyPerLesson = 10

type topicSeed struct {
	Title       string
	Slug        string
	Description string
	Icon        string
	Emoji       string
	Color       string
	BandScore   float64
}

type lessonSeed struct {
	Title       string
	Slug        string
	Description string
}

type wordSeed struct {
	Word         string
	IPA          string
	PartOfSpeech string
	MeaningVI    string
	MeaningEN    string
	Difficulty   models.DifficultyLevel
	TargetBand   float64
	Synonyms     []string
	Collocations []string
}

type lessonRecord struct {
	Lesson    models.Lesson
	Topic     models.Topic
	BandLevel models.BandLevel
	WordCount int
}

func Run(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		course, err := seedCourse(tx)
		if err != nil {
			return err
		}

		bands, err := seedBandLevels(tx, course.ID)
		if err != nil {
			return err
		}

		topics, err := seedTopics(tx, bands)
		if err != nil {
			return err
		}

		lessons, err := seedLessons(tx, topics)
		if err != nil {
			return err
		}

		vocabularies, err := seedVocabularies(tx, lessons)
		if err != nil {
			return err
		}

		if err := seedQuiz(tx, lessons, vocabularies); err != nil {
			return err
		}

		if err := seedAchievements(tx); err != nil {
			return err
		}

		demoUser, err := seedDemoUser(tx)
		if err != nil {
			return err
		}

		if err := seedDemoProgress(tx, demoUser); err != nil {
			return err
		}

		return updateCourseTotals(tx, course.ID)
	})
}

func seedCourse(db *gorm.DB) (models.Course, error) {
	course := models.Course{
		Title:        "IELTS Vocabulary Roadmap",
		Description:  "A structured IELTS vocabulary path from foundation words to advanced academic language.",
		Slug:         "ielts-vocabulary-roadmap",
		IsPublished:  true,
		OrderIndex:   1,
		BandMin:      5.0,
		BandMax:      8.5,
		TotalWords:   0,
		TotalLessons: 0,
		TotalTopics:  0,
	}

	if err := db.Where("slug = ?", course.Slug).Assign(course).FirstOrCreate(&course).Error; err != nil {
		return models.Course{}, fmt.Errorf("seed course: %w", err)
	}

	return course, nil
}

func seedBandLevels(db *gorm.DB, courseID uint) (map[float64]models.BandLevel, error) {
	bandSeeds := []struct {
		BandScore   float64
		MinScore    float64
		MaxScore    float64
		Title       string
		Description string
		StatusLabel string
	}{
		{5.0, 5.0, 5.5, "Band 5.0 Foundation", "Build clear IELTS answers with essential topic vocabulary.", "Mostly completed"},
		{6.0, 6.0, 6.5, "Band 6.0 Developing", "Expand range and accuracy across common IELTS academic themes.", "In progress"},
		{7.0, 7.0, 7.5, "Band 7.0 Confident", "Use precise topic language for stronger writing and speaking responses.", "Next"},
		{8.0, 8.0, 8.5, "Band 8.0 Advanced", "Master nuanced academic vocabulary and high-scoring collocations.", "Locked"},
	}

	bands := make(map[float64]models.BandLevel, len(bandSeeds))
	for index, seed := range bandSeeds {
		minScore := seed.MinScore
		maxScore := seed.MaxScore
		statusLabel := seed.StatusLabel
		band := models.BandLevel{
			CourseID:    courseID,
			BandScore:   seed.BandScore,
			MinScore:    &minScore,
			MaxScore:    &maxScore,
			Title:       seed.Title,
			Description: seed.Description,
			StatusLabel: &statusLabel,
			OrderIndex:  index + 1,
		}

		if err := db.Where("course_id = ? AND band_score = ?", courseID, seed.BandScore).Assign(band).FirstOrCreate(&band).Error; err != nil {
			return nil, fmt.Errorf("seed band %.1f: %w", seed.BandScore, err)
		}
		bands[seed.BandScore] = band
	}

	return bands, nil
}

func seedTopics(db *gorm.DB, bands map[float64]models.BandLevel) (map[string]models.Topic, error) {
	topicSeeds := []topicSeed{
		{"Education", "education", "Vocabulary for schools, universities, learning methods, and academic success.", "graduation-cap", "ED", "indigo", 5.0},
		{"Health", "health", "Vocabulary for wellbeing, healthcare, lifestyle, and public health.", "heart-pulse", "HL", "rose", 5.0},
		{"Environment", "environment", "Language for climate, conservation, pollution, and sustainability.", "leaf", "EV", "emerald", 5.0},
		{"Technology", "technology", "Words for digital tools, innovation, automation, and online life.", "cpu", "TC", "blue", 6.0},
		{"Business", "business", "Terms for companies, markets, entrepreneurship, and workplace decisions.", "briefcase-business", "BU", "amber", 6.0},
		{"Society", "society", "Vocabulary for community life, social change, equality, and public services.", "users-round", "SO", "violet", 6.0},
		{"Media", "media", "Language for journalism, advertising, social media, and public opinion.", "newspaper", "ME", "cyan", 6.0},
		{"Science", "science", "Academic words for research, evidence, discovery, and scientific ethics.", "microscope", "SC", "sky", 7.0},
		{"Government", "government", "Vocabulary for policy, elections, public spending, and civic participation.", "landmark", "GV", "slate", 7.0},
		{"Global Issues", "global-issues", "Language for migration, poverty, conflict, aid, and international cooperation.", "globe-2", "GI", "orange", 7.0},
		{"Economics", "economics", "Terms for growth, inflation, employment, investment, and global markets.", "chart-no-axes-combined", "EC", "green", 7.0},
		{"Law", "law", "Vocabulary for legal systems, civil rights, crime, courts, and digital law.", "scale", "LW", "zinc", 8.0},
		{"Culture", "culture", "Language for heritage, traditions, arts, identity, and cultural exchange.", "palette", "CU", "pink", 8.0},
		{"Innovation", "innovation", "Words for research, startups, automation, design, and future technologies.", "lightbulb", "IN", "yellow", 8.0},
		{"Psychology", "psychology", "Terms for motivation, memory, behaviour, stress, and decision making.", "brain", "PS", "purple", 8.0},
	}

	orderByBand := map[float64]int{}
	topics := make(map[string]models.Topic, len(topicSeeds))
	for _, seed := range topicSeeds {
		band := bands[seed.BandScore]
		orderByBand[seed.BandScore]++
		topic := models.Topic{
			BandLevelID: band.ID,
			Title:       seed.Title,
			Slug:        seed.Slug,
			Description: seed.Description,
			Icon:        seed.Icon,
			Emoji:       seed.Emoji,
			Color:       seed.Color,
			OrderIndex:  orderByBand[seed.BandScore],
		}

		if err := db.Where("band_level_id = ? AND slug = ?", band.ID, seed.Slug).Assign(topic).FirstOrCreate(&topic).Error; err != nil {
			return nil, fmt.Errorf("seed topic %s: %w", seed.Slug, err)
		}
		topics[seed.Slug] = topic
	}

	return topics, nil
}

func seedLessons(db *gorm.DB, topics map[string]models.Topic) (map[string]models.Lesson, error) {
	lessons := make(map[string]models.Lesson)
	for _, topic := range sortedTopics(topics) {
		bandMin, bandMax := bandRangeForTopic(topic.Slug)
		timeLimit := 600
		for index, seed := range lessonSeedsForTopic(topic.Slug, topic.Title) {
			estimatedMinutes := 10 + index%4
			xpReward := 50 + index*5
			lesson := models.Lesson{
				TopicID:              topic.ID,
				Title:                seed.Title,
				Slug:                 seed.Slug,
				Description:          seed.Description,
				RequiredScore:        80,
				EstimatedMinutes:     estimatedMinutes,
				BandMin:              &bandMin,
				BandMax:              &bandMax,
				XPReward:             xpReward,
				QuizTimeLimitSeconds: &timeLimit,
				OrderIndex:           index + 1,
				IsPublished:          true,
			}

			if err := db.Where("topic_id = ? AND slug = ?", topic.ID, seed.Slug).Assign(lesson).FirstOrCreate(&lesson).Error; err != nil {
				return nil, fmt.Errorf("seed lesson %s: %w", seed.Slug, err)
			}
			lessons[seed.Slug] = lesson
		}
	}

	return lessons, nil
}

func seedVocabularies(db *gorm.DB, lessons map[string]models.Lesson) (map[string]models.Vocabulary, error) {
	var topics []models.Topic
	if err := db.Find(&topics).Error; err != nil {
		return nil, fmt.Errorf("load topics for vocabulary seed: %w", err)
	}

	topicByID := make(map[uint]models.Topic, len(topics))
	for _, topic := range topics {
		topicByID[topic.ID] = topic
	}

	vocabularies := make(map[string]models.Vocabulary)
	for _, lesson := range sortedLessons(lessons) {
		topic := topicByID[lesson.TopicID]
		words := wordBankForTopic(topic.Slug)
		for index := 0; index < vocabularyPerLesson; index++ {
			seed := words[(lesson.OrderIndex+index-1)%len(words)]
			targetBand := seed.TargetBand
			vocabulary := models.Vocabulary{
				Word:             seed.Word,
				Slug:             fmt.Sprintf("%s-%s", lesson.Slug, slugify(seed.Word)),
				IPA:              seed.IPA,
				PartOfSpeech:     seed.PartOfSpeech,
				MeaningVI:        seed.MeaningVI,
				MeaningEN:        seed.MeaningEN,
				ShortDefinition:  seed.MeaningEN,
				ExampleSentence:  fmt.Sprintf("%s This word is useful in the %s lesson.", exampleForWord(seed.Word, topic.Title), lesson.Title),
				ExampleMeaningVI: fmt.Sprintf("Người học có thể dùng từ này khi thảo luận về chủ đề %s.", topic.Title),
				ExampleSource:    "Demo seed data",
				SynonymsJSON:     mustJSON(seed.Synonyms),
				AntonymsJSON:     mustJSON([]string{}),
				CollocationsJSON: mustJSON(seed.Collocations),
				IELTSUsage:       fmt.Sprintf("Useful for IELTS Writing Task 2 and Speaking Part 3 responses about %s.", topic.Title),
				Difficulty:       seed.Difficulty,
				TargetBand:       &targetBand,
			}

			if err := db.Where("slug = ?", vocabulary.Slug).Assign(vocabulary).FirstOrCreate(&vocabulary).Error; err != nil {
				return nil, fmt.Errorf("seed vocabulary %s: %w", vocabulary.Slug, err)
			}

			link := models.LessonVocabulary{
				LessonID:     lesson.ID,
				VocabularyID: vocabulary.ID,
				OrderIndex:   index + 1,
				IsRequired:   true,
			}
			if err := db.Where("lesson_id = ? AND vocabulary_id = ?", lesson.ID, vocabulary.ID).Assign(link).FirstOrCreate(&link).Error; err != nil {
				return nil, fmt.Errorf("seed lesson vocabulary %s: %w", vocabulary.Slug, err)
			}

			vocabularies[vocabulary.Slug] = vocabulary
		}
	}

	return vocabularies, nil
}

func seedQuiz(db *gorm.DB, lessons map[string]models.Lesson, vocabularies map[string]models.Vocabulary) error {
	for lessonSlug, lesson := range lessons {
		words, err := findLessonVocabularies(db, lesson.ID)
		if err != nil {
			return fmt.Errorf("load quiz words for %s: %w", lessonSlug, err)
		}
		if len(words) < 4 {
			continue
		}

		for questionIndex := 0; questionIndex < 3; questionIndex++ {
			word := words[questionIndex]
			topicID := lesson.TopicID
			vocabularyID := word.ID
			timeLimit := 45
			question := models.QuizQuestion{
				LessonID:         lesson.ID,
				TopicID:          &topicID,
				VocabularyID:     &vocabularyID,
				Type:             models.QuizQuestionMeaningChoice,
				Question:         fmt.Sprintf("What is the best meaning of \"%s\"?", word.Word),
				Explanation:      fmt.Sprintf("\"%s\" means: %s", word.Word, word.MeaningEN),
				Points:           20,
				TimeLimitSeconds: &timeLimit,
				OrderIndex:       questionIndex + 1,
			}

			if err := db.Where("lesson_id = ? AND order_index = ?", lesson.ID, question.OrderIndex).Assign(question).FirstOrCreate(&question).Error; err != nil {
				return fmt.Errorf("seed quiz question for %s: %w", word.Word, err)
			}

			options := []models.QuizOption{
				{QuestionID: question.ID, Label: "A", Content: word.MeaningEN, IsCorrect: true, OrderIndex: 1},
				{QuestionID: question.ID, Label: "B", Content: words[(questionIndex+2)%len(words)].MeaningEN, IsCorrect: false, OrderIndex: 2},
				{QuestionID: question.ID, Label: "C", Content: words[(questionIndex+4)%len(words)].MeaningEN, IsCorrect: false, OrderIndex: 3},
				{QuestionID: question.ID, Label: "D", Content: words[(questionIndex+6)%len(words)].MeaningEN, IsCorrect: false, OrderIndex: 4},
			}
			for _, option := range options {
				if err := db.Where("question_id = ? AND order_index = ?", question.ID, option.OrderIndex).Assign(option).FirstOrCreate(&option).Error; err != nil {
					return fmt.Errorf("seed quiz option for %s: %w", word.Word, err)
				}
			}
		}
	}

	_ = vocabularies
	return nil
}

func seedAchievements(db *gorm.DB) error {
	achievements := []models.Achievement{
		{Code: "FIRST_LESSON", Title: "First Lesson", Description: "Complete your first vocabulary lesson.", Icon: "book-open-check", Category: "lesson", RequirementType: "completed_lessons", RequirementValue: 1, XPReward: 50, SortOrder: 1, IsActive: true},
		{Code: "WORD_COLLECTOR_100", Title: "100 Words Learned", Description: "Learn 100 IELTS vocabulary items.", Icon: "layers", Category: "vocabulary", RequirementType: "learned_words", RequirementValue: 100, XPReward: 150, SortOrder: 2, IsActive: true},
		{Code: "QUIZ_STARTER", Title: "Quiz Starter", Description: "Pass your first lesson quiz.", Icon: "circle-check", Category: "quiz", RequirementType: "passed_quizzes", RequirementValue: 1, XPReward: 75, SortOrder: 3, IsActive: true},
		{Code: "STREAK_7", Title: "Seven Day Streak", Description: "Study for seven days in a row.", Icon: "flame", Category: "streak", RequirementType: "current_streak", RequirementValue: 7, XPReward: 120, SortOrder: 4, IsActive: true},
		{Code: "BAND_7_READY", Title: "Band 7 Ready", Description: "Master your first Band 7 vocabulary lesson.", Icon: "trophy", Category: "roadmap", RequirementType: "band_7_lessons", RequirementValue: 1, XPReward: 150, SortOrder: 5, IsActive: true},
	}

	for _, achievement := range achievements {
		if err := db.Where("code = ?", achievement.Code).Assign(achievement).FirstOrCreate(&achievement).Error; err != nil {
			return fmt.Errorf("seed achievement %s: %w", achievement.Code, err)
		}
	}

	return nil
}

func seedDemoUser(db *gorm.DB) (models.User, error) {
	username := "demo"
	currentBand := 6.0
	startingBand := 5.0
	recommendedBand := 6.0
	now := demoNow()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("hash demo password: %w", err)
	}

	user := models.User{
		Email:           "demo@example.com",
		Name:            "Demo Learner",
		Username:        &username,
		PasswordHash:    string(passwordHash),
		Role:            models.UserRoleUser,
		TargetBand:      7.0,
		CurrentBand:     &currentBand,
		StartingBand:    &startingBand,
		RecommendedBand: &recommendedBand,
		TotalXP:         3650,
		Level:           8,
		CurrentStreak:   14,
		LongestStreak:   21,
		LastActiveAt:    &now,
		Timezone:        "Asia/Ho_Chi_Minh",
		Locale:          "en",
	}

	if err := db.Where("email = ?", user.Email).Assign(user).FirstOrCreate(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("seed demo user: %w", err)
	}

	return user, nil
}

func seedDemoProgress(db *gorm.DB, user models.User) error {
	records, err := loadLessonRecords(db)
	if err != nil {
		return err
	}

	now := demoNow()
	for _, record := range records {
		status, hasProgress := demoLessonStatus(record)
		if !hasProgress {
			continue
		}

		if err := seedUserLessonProgress(db, user.ID, record, status, now); err != nil {
			return err
		}
		if err := seedVocabularyProgress(db, user.ID, record, status, now); err != nil {
			return err
		}
		if err := seedQuizAttempts(db, user.ID, record, status, now); err != nil {
			return err
		}
		if err := seedXPEvents(db, user.ID, record, status, now); err != nil {
			return err
		}
	}

	if err := seedDailyActivities(db, user.ID, now); err != nil {
		return err
	}
	if err := seedUserAchievements(db, user.ID, now); err != nil {
		return err
	}

	return nil
}

func demoLessonStatus(record lessonRecord) (models.LessonStatus, bool) {
	switch record.BandLevel.BandScore {
	case 5.0:
		if record.Topic.Slug == "environment" && record.Lesson.OrderIndex == 7 {
			return models.LessonStatusInProgress, true
		}
		if record.Topic.Slug == "environment" && record.Lesson.OrderIndex == 8 {
			return models.LessonStatusUnlocked, true
		}
		return models.LessonStatusCompleted, true
	case 6.0:
		if record.Topic.Slug == "technology" {
			switch record.Lesson.OrderIndex {
			case 1, 2:
				return models.LessonStatusCompleted, true
			case 3:
				return models.LessonStatusInProgress, true
			case 4:
				return models.LessonStatusUnlocked, true
			default:
				return models.LessonStatusLocked, false
			}
		}
		if record.Topic.Slug == "business" {
			switch record.Lesson.OrderIndex {
			case 1, 2:
				return models.LessonStatusCompleted, true
			case 3:
				return models.LessonStatusInProgress, true
			case 4:
				return models.LessonStatusUnlocked, true
			default:
				return models.LessonStatusLocked, false
			}
		}
		if record.Topic.Slug == "society" {
			if record.Lesson.OrderIndex == 1 {
				return models.LessonStatusCompleted, true
			}
			if record.Lesson.OrderIndex == 2 {
				return models.LessonStatusUnlocked, true
			}
			return models.LessonStatusLocked, false
		}
		if record.Topic.Slug == "media" && record.Lesson.OrderIndex == 1 {
			return models.LessonStatusUnlocked, true
		}
	case 7.0:
		if record.Topic.Slug == "science" && record.Lesson.OrderIndex == 1 {
			return models.LessonStatusUnlocked, true
		}
	case 8.0:
		if record.Lesson.OrderIndex == 1 {
			return models.LessonStatusLocked, true
		}
	}

	return models.LessonStatusLocked, false
}

func seedUserLessonProgress(db *gorm.DB, userID uint, record lessonRecord, status models.LessonStatus, now time.Time) error {
	score, bestScore := scoreForStatus(status)
	startedAt := now.AddDate(0, 0, -record.Lesson.OrderIndex-2)
	lastStudiedAt := now.AddDate(0, 0, -record.Lesson.OrderIndex)
	completedAt := (*time.Time)(nil)
	wordsLearned := 0
	bestXP := 0

	switch status {
	case models.LessonStatusCompleted:
		done := lastStudiedAt
		completedAt = &done
		wordsLearned = record.WordCount
		bestXP = record.Lesson.XPReward
	case models.LessonStatusInProgress:
		wordsLearned = progressWords(record.WordCount, 45)
		bestXP = record.Lesson.XPReward / 2
	}

	progress := models.UserLessonProgress{
		UserID:        userID,
		LessonID:      record.Lesson.ID,
		Status:        status,
		Score:         score,
		BestScore:     bestScore,
		BestXP:        bestXP,
		WordsLearned:  wordsLearned,
		TotalWords:    record.WordCount,
		CompletedAt:   completedAt,
		StartedAt:     &startedAt,
		LastStudiedAt: &lastStudiedAt,
	}

	if status == models.LessonStatusLocked || status == models.LessonStatusUnlocked {
		progress.Score = nil
		progress.BestScore = nil
		progress.CompletedAt = nil
		progress.LastStudiedAt = nil
		if status == models.LessonStatusLocked {
			progress.StartedAt = nil
		}
	}

	if err := db.Where("user_id = ? AND lesson_id = ?", userID, record.Lesson.ID).Assign(progress).FirstOrCreate(&progress).Error; err != nil {
		return fmt.Errorf("seed lesson progress %s: %w", record.Lesson.Slug, err)
	}

	return nil
}

func seedVocabularyProgress(db *gorm.DB, userID uint, record lessonRecord, status models.LessonStatus, now time.Time) error {
	words, err := findLessonVocabularies(db, record.Lesson.ID)
	if err != nil {
		return fmt.Errorf("load vocab progress words for %s: %w", record.Lesson.Slug, err)
	}

	limit := len(words)
	if status == models.LessonStatusInProgress {
		limit = progressWords(len(words), 45)
	}
	if status == models.LessonStatusUnlocked || status == models.LessonStatusLocked {
		limit = 0
	}

	for index, word := range words[:limit] {
		vocabStatus := models.VocabularyStatusMastered
		reviewCount := 5 + index%4
		correctCount := reviewCount
		wrongCount := index % 2
		intervalDays := 7 + index
		rating := models.FlashcardRatingEasy
		firstLearnedAt := now.AddDate(0, 0, -14-index)
		learnedAt := now.AddDate(0, 0, -10-index)
		masteredAt := now.AddDate(0, 0, -5-index)
		lastReviewedAt := now.AddDate(0, 0, -index%4)
		nextReviewAt := now.AddDate(0, 0, 3+index%5)

		if status == models.LessonStatusInProgress {
			vocabStatus = models.VocabularyStatusLearning
			reviewCount = 1 + index%2
			correctCount = index % 2
			wrongCount = 1
			intervalDays = 1
			rating = models.FlashcardRatingHard
			masteredAt = time.Time{}
			nextReviewAt = now.AddDate(0, 0, 1)
		} else if index%5 == 0 {
			vocabStatus = models.VocabularyStatusReview
			rating = models.FlashcardRatingGood
			nextReviewAt = now
		}

		progress := models.UserVocabularyProgress{
			UserID:         userID,
			VocabularyID:   word.ID,
			Status:         vocabStatus,
			EaseFactor:     2.3 + float64(index%3)/10,
			IntervalDays:   intervalDays,
			ReviewCount:    reviewCount,
			CorrectCount:   correctCount,
			WrongCount:     wrongCount,
			LastRating:     &rating,
			FirstLearnedAt: &firstLearnedAt,
			LearnedAt:      &learnedAt,
			LastReviewedAt: &lastReviewedAt,
			NextReviewAt:   &nextReviewAt,
		}
		if !masteredAt.IsZero() {
			progress.MasteredAt = &masteredAt
		}

		if err := db.Where("user_id = ? AND vocabulary_id = ?", userID, word.ID).Assign(progress).FirstOrCreate(&progress).Error; err != nil {
			return fmt.Errorf("seed vocabulary progress %s: %w", word.Slug, err)
		}
	}

	if status == models.LessonStatusUnlocked {
		for _, word := range words[:minInt(3, len(words))] {
			nextReviewAt := now.AddDate(0, 0, 2)
			progress := models.UserVocabularyProgress{
				UserID:       userID,
				VocabularyID: word.ID,
				Status:       models.VocabularyStatusNew,
				EaseFactor:   2.5,
				NextReviewAt: &nextReviewAt,
			}
			if err := db.Where("user_id = ? AND vocabulary_id = ?", userID, word.ID).Assign(progress).FirstOrCreate(&progress).Error; err != nil {
				return fmt.Errorf("seed new vocabulary progress %s: %w", word.Slug, err)
			}
		}
	}

	return nil
}

func seedQuizAttempts(db *gorm.DB, userID uint, record lessonRecord, status models.LessonStatus, now time.Time) error {
	if status != models.LessonStatusCompleted && status != models.LessonStatusInProgress {
		return nil
	}

	startedAt := now.AddDate(0, 0, -record.Lesson.OrderIndex).Add(-20 * time.Minute)
	finishedAt := startedAt.Add(8 * time.Minute)
	duration := 480
	score := 92
	correctAnswers := 9
	xpEarned := record.Lesson.XPReward
	passed := true
	if status == models.LessonStatusInProgress {
		score = 55
		correctAnswers = 5
		xpEarned = 10
		passed = false
	}

	attempt := models.UserQuizAttempt{
		UserID:          userID,
		LessonID:        record.Lesson.ID,
		Score:           score,
		Points:          correctAnswers * 20,
		TotalQuestions:  10,
		CorrectAnswers:  correctAnswers,
		DurationSeconds: &duration,
		XPEarned:        xpEarned,
		Passed:          passed,
		StartedAt:       &startedAt,
		FinishedAt:      &finishedAt,
		CreatedAt:       finishedAt,
	}

	if err := db.Where("user_id = ? AND lesson_id = ? AND score = ? AND created_at = ?", userID, record.Lesson.ID, score, finishedAt).Assign(attempt).FirstOrCreate(&attempt).Error; err != nil {
		return fmt.Errorf("seed quiz attempt %s: %w", record.Lesson.Slug, err)
	}

	if status == models.LessonStatusCompleted && record.Lesson.OrderIndex%3 == 0 {
		failedStart := startedAt.AddDate(0, 0, -1)
		failedFinish := failedStart.Add(9 * time.Minute)
		failedDuration := 540
		failed := models.UserQuizAttempt{
			UserID:          userID,
			LessonID:        record.Lesson.ID,
			Score:           65,
			Points:          120,
			TotalQuestions:  10,
			CorrectAnswers:  6,
			DurationSeconds: &failedDuration,
			XPEarned:        5,
			Passed:          false,
			StartedAt:       &failedStart,
			FinishedAt:      &failedFinish,
			CreatedAt:       failedFinish,
		}
		if err := db.Where("user_id = ? AND lesson_id = ? AND score = ? AND created_at = ?", userID, record.Lesson.ID, failed.Score, failedFinish).Assign(failed).FirstOrCreate(&failed).Error; err != nil {
			return fmt.Errorf("seed failed quiz attempt %s: %w", record.Lesson.Slug, err)
		}
	}

	return nil
}

func seedXPEvents(db *gorm.DB, userID uint, record lessonRecord, status models.LessonStatus, now time.Time) error {
	if status != models.LessonStatusCompleted && status != models.LessonStatusInProgress {
		return nil
	}

	events := []models.UserXPEvent{
		{UserID: userID, SourceType: "FLASHCARD_REVIEW", SourceID: &record.Lesson.ID, XP: 15, Description: fmt.Sprintf("Reviewed flashcards in %s", record.Lesson.Title), CreatedAt: now.AddDate(0, 0, -record.Lesson.OrderIndex)},
	}
	if status == models.LessonStatusCompleted {
		events = append(events,
			models.UserXPEvent{UserID: userID, SourceType: "LESSON_COMPLETION", SourceID: &record.Lesson.ID, XP: record.Lesson.XPReward, Description: fmt.Sprintf("Completed %s", record.Lesson.Title), CreatedAt: now.AddDate(0, 0, -record.Lesson.OrderIndex)},
			models.UserXPEvent{UserID: userID, SourceType: "QUIZ_SUCCESS", SourceID: &record.Lesson.ID, XP: 25, Description: fmt.Sprintf("Passed quiz for %s", record.Lesson.Title), CreatedAt: now.AddDate(0, 0, -record.Lesson.OrderIndex)},
		)
	}

	for _, event := range events {
		if err := db.Where("user_id = ? AND source_type = ? AND source_id = ? AND description = ?", userID, event.SourceType, record.Lesson.ID, event.Description).Assign(event).FirstOrCreate(&event).Error; err != nil {
			return fmt.Errorf("seed xp event %s: %w", record.Lesson.Slug, err)
		}
	}

	return nil
}

func seedDailyActivities(db *gorm.DB, userID uint, now time.Time) error {
	for daysAgo := 0; daysAgo < 14; daysAgo++ {
		activityTime := now.AddDate(0, 0, -daysAgo)
		date := time.Date(activityTime.Year(), activityTime.Month(), activityTime.Day(), 0, 0, 0, 0, time.UTC)
		accuracy := 82.0 + float64(daysAgo%6)
		activity := models.DailyActivity{
			UserID:             userID,
			Date:               date,
			WordsLearned:       4 + daysAgo%7,
			WordsReviewed:      12 + daysAgo%9,
			QuizzesTaken:       1 + daysAgo%2,
			LessonsDone:        daysAgo % 3,
			XPEarned:           90 + daysAgo*8,
			AccuracyPercent:    &accuracy,
			ActiveMinutes:      18 + daysAgo%10,
			ChallengeProgress:  minInt(10, 4+daysAgo%8),
			ChallengeCompleted: daysAgo%3 != 0,
		}

		if err := db.Where("user_id = ? AND date = ?", userID, date).Assign(activity).FirstOrCreate(&activity).Error; err != nil {
			return fmt.Errorf("seed daily activity %s: %w", date.Format("2006-01-02"), err)
		}
	}

	return nil
}

func seedUserAchievements(db *gorm.DB, userID uint, now time.Time) error {
	unlocked := map[string]int{
		"FIRST_LESSON":       1,
		"WORD_COLLECTOR_100": 100,
		"QUIZ_STARTER":       1,
		"STREAK_7":           14,
	}

	for code, progressValue := range unlocked {
		var achievement models.Achievement
		if err := db.Where("code = ?", code).First(&achievement).Error; err != nil {
			return fmt.Errorf("load achievement %s: %w", code, err)
		}
		userAchievement := models.UserAchievement{
			UserID:        userID,
			AchievementID: achievement.ID,
			ProgressValue: progressValue,
			IsSeen:        code != "STREAK_7",
			UnlockedAt:    now.AddDate(0, 0, -progressValue%10),
		}

		if err := db.Where("user_id = ? AND achievement_id = ?", userID, achievement.ID).Assign(userAchievement).FirstOrCreate(&userAchievement).Error; err != nil {
			return fmt.Errorf("seed user achievement %s: %w", code, err)
		}
	}

	return nil
}

func updateCourseTotals(db *gorm.DB, courseID uint) error {
	var course models.Course
	if err := db.First(&course, courseID).Error; err != nil {
		return fmt.Errorf("load course totals: %w", err)
	}

	var totalTopics int64
	if err := db.Model(&models.Topic{}).
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).
		Count(&totalTopics).Error; err != nil {
		return fmt.Errorf("count topics: %w", err)
	}

	var totalLessons int64
	if err := db.Model(&models.Lesson{}).
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).
		Count(&totalLessons).Error; err != nil {
		return fmt.Errorf("count lessons: %w", err)
	}

	var totalWords int64
	if err := db.Model(&models.LessonVocabulary{}).
		Joins("JOIN lessons ON lessons.id = lesson_vocabularies.lesson_id").
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).
		Count(&totalWords).Error; err != nil {
		return fmt.Errorf("count words: %w", err)
	}

	course.TotalTopics = int(totalTopics)
	course.TotalLessons = int(totalLessons)
	course.TotalWords = int(totalWords)

	if err := db.Save(&course).Error; err != nil {
		return fmt.Errorf("update course totals: %w", err)
	}

	return nil
}

func loadLessonRecords(db *gorm.DB) ([]lessonRecord, error) {
	var lessons []models.Lesson
	if err := db.Preload("Topic.BandLevel").Find(&lessons).Error; err != nil {
		return nil, fmt.Errorf("load lessons for demo progress: %w", err)
	}

	records := make([]lessonRecord, 0, len(lessons))
	for _, lesson := range lessons {
		var wordCount int64
		if err := db.Model(&models.LessonVocabulary{}).Where("lesson_id = ?", lesson.ID).Count(&wordCount).Error; err != nil {
			return nil, fmt.Errorf("count lesson words %s: %w", lesson.Slug, err)
		}
		records = append(records, lessonRecord{
			Lesson:    lesson,
			Topic:     lesson.Topic,
			BandLevel: lesson.Topic.BandLevel,
			WordCount: int(wordCount),
		})
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].BandLevel.OrderIndex != records[j].BandLevel.OrderIndex {
			return records[i].BandLevel.OrderIndex < records[j].BandLevel.OrderIndex
		}
		if records[i].Topic.OrderIndex != records[j].Topic.OrderIndex {
			return records[i].Topic.OrderIndex < records[j].Topic.OrderIndex
		}
		return records[i].Lesson.OrderIndex < records[j].Lesson.OrderIndex
	})

	return records, nil
}

func findLessonVocabularies(db *gorm.DB, lessonID uint) ([]models.Vocabulary, error) {
	var words []models.Vocabulary
	err := db.Model(&models.Vocabulary{}).
		Joins("JOIN lesson_vocabularies ON lesson_vocabularies.vocabulary_id = vocabularies.id").
		Where("lesson_vocabularies.lesson_id = ?", lessonID).
		Order("lesson_vocabularies.order_index ASC").
		Find(&words).Error
	return words, err
}

func lessonSeedsForTopic(topicSlug string, topicTitle string) []lessonSeed {
	custom := map[string][]string{
		"education":     {"School Systems", "Study Habits", "Higher Education", "Exams and Assessment", "Online Learning", "Educational Inequality", "Lifelong Learning", "Student Wellbeing"},
		"health":        {"Healthy Lifestyle", "Public Health", "Mental Wellbeing", "Medical Services", "Nutrition", "Disease Prevention", "Healthcare Technology", "Work-Life Balance"},
		"environment":   {"Climate Change", "Pollution Control", "Conservation", "Renewable Energy", "Urban Sustainability", "Waste Management", "Water Scarcity", "Biodiversity"},
		"technology":    {"Digital Communication", "Internet & Connectivity", "Artificial Intelligence", "Data Privacy & Security", "E-commerce", "Mobile Technology", "Software & Applications", "Cloud Computing", "Innovation & Research", "Future Technologies"},
		"business":      {"Entrepreneurship", "Marketing", "Consumer Behaviour", "Workplace Strategy", "Global Trade", "Finance Basics", "Corporate Responsibility", "Leadership"},
		"society":       {"Community Life", "Social Change", "Equality", "Urbanisation", "Family Structures", "Public Services", "Crime Prevention", "Cultural Diversity"},
		"media":         {"News & Journalism", "Advertising", "Social Media", "Digital Entertainment", "Media Literacy", "Censorship", "Public Opinion", "Online Influence"},
		"science":       {"Scientific Research", "Space Exploration", "Genetics", "Medical Science", "Energy Systems", "Environmental Science", "Data & Evidence", "Ethical Research"},
		"government":    {"Public Policy", "Elections", "Local Services", "Taxation", "Public Spending", "Regulation", "Civic Participation", "International Relations"},
		"global-issues": {"Migration", "Poverty", "Human Rights", "Food Security", "Conflict Resolution", "Global Health", "Climate Cooperation", "International Aid"},
		"economics":     {"Economic Growth", "Employment", "Inflation", "Inequality", "Consumer Spending", "Investment", "Productivity", "Global Markets"},
		"law":           {"Legal Systems", "Crime & Punishment", "Civil Rights", "Contracts", "Courts", "Policing", "International Law", "Digital Law"},
		"culture":       {"Traditions", "Arts", "Heritage", "Language Identity", "Festivals", "Cultural Exchange", "Popular Culture", "Museums"},
		"innovation":    {"Creative Thinking", "Product Design", "Startups", "Research Funding", "Future Transport", "Green Innovation", "Automation", "Disruptive Change"},
		"psychology":    {"Motivation", "Memory", "Behaviour", "Stress", "Personality", "Decision Making", "Social Influence", "Learning Mindset"},
	}

	titles := custom[topicSlug]
	lessons := make([]lessonSeed, 0, len(titles))
	for _, title := range titles {
		lessons = append(lessons, lessonSeed{
			Title:       title,
			Slug:        slugify(title),
			Description: fmt.Sprintf("Learn IELTS vocabulary for %s, including definitions, examples, collocations, and quiz practice.", title),
		})
	}

	_ = topicTitle
	return lessons
}

func wordBankForTopic(topicSlug string) []wordSeed {
	banks := map[string][]wordSeed{
		"education": {
			w("curriculum", "/kəˈrɪkjələm/", "noun", "chương trình học", "The subjects and content taught in a course.", models.DifficultyIntermediate, 5.5, "syllabus", "school curriculum"),
			w("literacy", "/ˈlɪtərəsi/", "noun", "khả năng đọc viết", "The ability to read and write effectively.", models.DifficultyBeginner, 5.0, "reading ability", "digital literacy"),
			w("assessment", "/əˈsesmənt/", "noun", "sự đánh giá", "A method of judging progress or ability.", models.DifficultyIntermediate, 6.0, "evaluation", "assessment criteria"),
			w("scholarship", "/ˈskɒlərʃɪp/", "noun", "học bổng", "Financial support awarded for study.", models.DifficultyBeginner, 5.0, "grant", "full scholarship"),
			w("competence", "/ˈkɒmpɪtəns/", "noun", "năng lực", "The ability to do something well.", models.DifficultyIntermediate, 6.0, "proficiency", "language competence"),
			w("pedagogy", "/ˈpedəɡɒdʒi/", "noun", "phương pháp sư phạm", "The method and practice of teaching.", models.DifficultyAdvanced, 7.0, "teaching method", "modern pedagogy"),
			w("attendance", "/əˈtendəns/", "noun", "sự tham dự", "Being present at a class or event.", models.DifficultyBeginner, 5.0, "presence", "attendance record"),
			w("tuition", "/tjuˈɪʃən/", "noun", "học phí", "Money paid for instruction.", models.DifficultyBeginner, 5.0, "fees", "tuition fees"),
			w("mentor", "/ˈmentɔːr/", "noun", "người cố vấn", "An experienced person who guides someone.", models.DifficultyBeginner, 5.5, "advisor", "academic mentor"),
			w("discipline", "/ˈdɪsəplɪn/", "noun", "kỷ luật", "Controlled behaviour or a field of study.", models.DifficultyIntermediate, 5.5, "order", "classroom discipline"),
		},
		"health": {
			w("prevention", "/prɪˈvenʃən/", "noun", "sự phòng ngừa", "Action taken to stop something happening.", models.DifficultyIntermediate, 6.0, "avoidance", "disease prevention"),
			w("diagnosis", "/ˌdaɪəɡˈnəʊsɪs/", "noun", "sự chẩn đoán", "Identification of an illness or problem.", models.DifficultyIntermediate, 6.5, "identification", "accurate diagnosis"),
			w("sedentary", "/ˈsedntri/", "adjective", "ít vận động", "Involving much sitting and little exercise.", models.DifficultyAdvanced, 7.0, "inactive", "sedentary lifestyle"),
			w("nutrition", "/njuˈtrɪʃən/", "noun", "dinh dưỡng", "Food needed for health and growth.", models.DifficultyIntermediate, 6.0, "diet", "balanced nutrition"),
			w("immunity", "/ɪˈmjuːnəti/", "noun", "khả năng miễn dịch", "Protection against disease.", models.DifficultyIntermediate, 6.0, "resistance", "natural immunity"),
			w("therapy", "/ˈθerəpi/", "noun", "liệu pháp", "Treatment for physical or mental illness.", models.DifficultyBeginner, 5.5, "treatment", "therapy session"),
			w("chronic", "/ˈkrɒnɪk/", "adjective", "mãn tính", "Continuing for a long time.", models.DifficultyAdvanced, 7.0, "long-term", "chronic illness"),
			w("hygiene", "/ˈhaɪdʒiːn/", "noun", "vệ sinh", "Practices that keep people healthy and clean.", models.DifficultyBeginner, 5.0, "cleanliness", "personal hygiene"),
			w("wellbeing", "/ˌwelˈbiːɪŋ/", "noun", "sức khỏe tinh thần và thể chất", "General health and happiness.", models.DifficultyIntermediate, 6.0, "welfare", "mental wellbeing"),
			w("vaccination", "/ˌvæksɪˈneɪʃən/", "noun", "tiêm chủng", "Giving a vaccine to protect against disease.", models.DifficultyIntermediate, 6.5, "immunisation", "vaccination programme"),
		},
		"environment": {
			w("sustainable", "/səˈsteɪnəbl/", "adjective", "bền vững", "Able to continue without damaging resources.", models.DifficultyIntermediate, 6.5, "eco-friendly", "sustainable development"),
			w("emission", "/ɪˈmɪʃən/", "noun", "khí thải", "Gas or substance released into the air.", models.DifficultyIntermediate, 6.0, "release", "carbon emissions"),
			w("conservation", "/ˌkɒnsəˈveɪʃən/", "noun", "sự bảo tồn", "Protection of nature and resources.", models.DifficultyIntermediate, 6.5, "preservation", "wildlife conservation"),
			w("biodiversity", "/ˌbaɪəʊdaɪˈvɜːsəti/", "noun", "đa dạng sinh học", "The variety of living things in an area.", models.DifficultyAdvanced, 7.0, "biological variety", "protect biodiversity"),
			w("renewable", "/rɪˈnjuːəbl/", "adjective", "có thể tái tạo", "Naturally replaced and not exhausted.", models.DifficultyIntermediate, 6.0, "reusable", "renewable energy"),
			w("pollution", "/pəˈluːʃən/", "noun", "ô nhiễm", "Damage caused by harmful substances.", models.DifficultyBeginner, 5.0, "contamination", "air pollution"),
			w("habitat", "/ˈhæbɪtæt/", "noun", "môi trường sống", "The natural home of an animal or plant.", models.DifficultyIntermediate, 6.0, "environment", "habitat loss"),
			w("scarce", "/skeəs/", "adjective", "khan hiếm", "Not enough for demand.", models.DifficultyIntermediate, 6.0, "limited", "scarce resources"),
			w("mitigate", "/ˈmɪtɪɡeɪt/", "verb", "giảm nhẹ", "To make something less harmful.", models.DifficultyAdvanced, 7.0, "reduce", "mitigate impact"),
			w("ecosystem", "/ˈiːkəʊsɪstəm/", "noun", "hệ sinh thái", "A community of organisms and their environment.", models.DifficultyIntermediate, 6.0, "ecological system", "marine ecosystem"),
		},
		"technology": {
			w("innovation", "/ˌɪnəˈveɪʃən/", "noun", "sự đổi mới", "A new idea, method, or product.", models.DifficultyIntermediate, 6.0, "invention", "technological innovation"),
			w("automation", "/ˌɔːtəˈmeɪʃən/", "noun", "tự động hóa", "The use of machines or software to do work.", models.DifficultyIntermediate, 6.5, "mechanisation", "industrial automation"),
			w("algorithm", "/ˈælɡərɪðəm/", "noun", "thuật toán", "A set of rules used by a computer.", models.DifficultyAdvanced, 7.0, "procedure", "recommendation algorithm"),
			w("privacy", "/ˈprɪvəsi/", "noun", "quyền riêng tư", "The right to keep personal information secret.", models.DifficultyIntermediate, 6.0, "confidentiality", "protect privacy"),
			w("device", "/dɪˈvaɪs/", "noun", "thiết bị", "A machine or tool for a specific purpose.", models.DifficultyBeginner, 5.0, "gadget", "mobile device"),
			w("platform", "/ˈplætfɔːrm/", "noun", "nền tảng", "A digital service or system.", models.DifficultyIntermediate, 6.0, "system", "learning platform"),
			w("cybersecurity", "/ˌsaɪbəsɪˈkjʊərəti/", "noun", "an ninh mạng", "Protection of computer systems and data.", models.DifficultyAdvanced, 7.0, "digital security", "cybersecurity threat"),
			w("accessible", "/əkˈsesəbl/", "adjective", "dễ tiếp cận", "Easy to reach, use, or understand.", models.DifficultyIntermediate, 6.0, "available", "accessible information"),
			w("disruptive", "/dɪsˈrʌptɪv/", "adjective", "mang tính đột phá", "Causing major change to an industry or system.", models.DifficultyAdvanced, 7.0, "transformative", "disruptive innovation"),
			w("connectivity", "/ˌkɒnekˈtɪvəti/", "noun", "khả năng kết nối", "The ability of systems to connect and communicate.", models.DifficultyIntermediate, 6.5, "connection", "internet connectivity"),
		},
	}

	if words, ok := banks[topicSlug]; ok {
		return words
	}

	return advancedWordBank(topicSlug)
}

func advancedWordBank(topicSlug string) []wordSeed {
	switch topicSlug {
	case "business":
		return []wordSeed{w("entrepreneur", "/ˌɒntrəprəˈnɜːr/", "noun", "doanh nhân", "A person who starts and runs a business.", models.DifficultyAdvanced, 7.0, "business owner", "successful entrepreneur"), w("revenue", "/ˈrevənjuː/", "noun", "doanh thu", "Income earned by a company.", models.DifficultyIntermediate, 6.0, "income", "annual revenue"), w("investment", "/ɪnˈvestmənt/", "noun", "sự đầu tư", "Money used to gain future profit.", models.DifficultyIntermediate, 6.0, "funding", "foreign investment"), w("strategy", "/ˈstrætədʒi/", "noun", "chiến lược", "A plan to achieve a goal.", models.DifficultyIntermediate, 6.0, "plan", "business strategy"), w("consumer", "/kənˈsjuːmə/", "noun", "người tiêu dùng", "A person who buys goods or services.", models.DifficultyBeginner, 5.5, "buyer", "consumer behaviour"), w("brand", "/brænd/", "noun", "thương hiệu", "A product identity recognized by customers.", models.DifficultyBeginner, 5.5, "label", "brand loyalty"), w("profitability", "/ˌprɒfɪtəˈbɪləti/", "noun", "khả năng sinh lời", "The ability to make profit.", models.DifficultyAdvanced, 7.0, "profit potential", "improve profitability"), w("negotiation", "/nɪˌɡəʊʃiˈeɪʃən/", "noun", "sự đàm phán", "Discussion to reach an agreement.", models.DifficultyIntermediate, 6.5, "discussion", "business negotiation"), w("stakeholder", "/ˈsteɪkhəʊldə/", "noun", "bên liên quan", "A person or group affected by decisions.", models.DifficultyAdvanced, 7.0, "interested party", "key stakeholder"), w("liability", "/ˌlaɪəˈbɪləti/", "noun", "trách nhiệm pháp lý", "Legal or financial responsibility.", models.DifficultyAdvanced, 7.0, "responsibility", "financial liability")}
	case "society":
		return []wordSeed{w("community", "/kəˈmjuːnəti/", "noun", "cộng đồng", "People living together or sharing interests.", models.DifficultyBeginner, 5.0, "society", "local community"), w("inequality", "/ˌɪnɪˈkwɒləti/", "noun", "bất bình đẳng", "An unfair difference between groups.", models.DifficultyIntermediate, 6.5, "disparity", "social inequality"), w("integration", "/ˌɪntɪˈɡreɪʃən/", "noun", "sự hòa nhập", "The process of becoming part of a group.", models.DifficultyAdvanced, 7.0, "inclusion", "social integration"), w("diversity", "/daɪˈvɜːsəti/", "noun", "sự đa dạng", "Variety among people or things.", models.DifficultyIntermediate, 6.0, "variety", "cultural diversity"), w("urbanisation", "/ˌɜːbənaɪˈzeɪʃən/", "noun", "đô thị hóa", "The growth of cities.", models.DifficultyAdvanced, 7.0, "city growth", "rapid urbanisation"), w("welfare", "/ˈwelfeə/", "noun", "phúc lợi", "Health, comfort, and support.", models.DifficultyIntermediate, 6.0, "wellbeing", "social welfare"), w("poverty", "/ˈpɒvəti/", "noun", "nghèo đói", "The state of lacking basic resources.", models.DifficultyBeginner, 5.5, "deprivation", "reduce poverty"), w("migration", "/maɪˈɡreɪʃən/", "noun", "di cư", "Movement from one place to another.", models.DifficultyIntermediate, 6.5, "relocation", "rural migration"), w("cohesion", "/kəʊˈhiːʒən/", "noun", "sự gắn kết", "Unity among members of a group.", models.DifficultyAdvanced, 7.0, "unity", "social cohesion"), w("participation", "/pɑːˌtɪsɪˈpeɪʃən/", "noun", "sự tham gia", "Taking part in an activity.", models.DifficultyIntermediate, 6.0, "involvement", "public participation")}
	case "media":
		return []wordSeed{w("journalism", "/ˈdʒɜːnəlɪzəm/", "noun", "nghề báo", "The work of collecting and reporting news.", models.DifficultyIntermediate, 6.5, "reporting", "investigative journalism"), w("broadcast", "/ˈbrɔːdkɑːst/", "verb", "phát sóng", "To send out programmes by radio, TV, or online.", models.DifficultyBeginner, 5.5, "transmit", "broadcast news"), w("advertising", "/ˈædvətaɪzɪŋ/", "noun", "quảng cáo", "Promoting products or services.", models.DifficultyBeginner, 5.5, "promotion", "online advertising"), w("censorship", "/ˈsensəʃɪp/", "noun", "sự kiểm duyệt", "Control of information by authority.", models.DifficultyAdvanced, 7.0, "restriction", "media censorship"), w("bias", "/ˈbaɪəs/", "noun", "thiên kiến", "An unfair preference or viewpoint.", models.DifficultyIntermediate, 6.5, "prejudice", "political bias"), w("audience", "/ˈɔːdiəns/", "noun", "khán giả", "People who watch, listen, or read.", models.DifficultyBeginner, 5.0, "viewers", "target audience"), w("coverage", "/ˈkʌvərɪdʒ/", "noun", "phạm vi đưa tin", "The reporting of an event or issue.", models.DifficultyIntermediate, 6.0, "reporting", "media coverage"), w("credibility", "/ˌkredəˈbɪləti/", "noun", "độ tin cậy", "The quality of being trusted.", models.DifficultyAdvanced, 7.0, "trustworthiness", "source credibility"), w("viral", "/ˈvaɪrəl/", "adjective", "lan truyền nhanh", "Spreading quickly online.", models.DifficultyIntermediate, 6.0, "popular", "viral content"), w("influence", "/ˈɪnfluəns/", "noun", "sự ảnh hưởng", "The power to affect opinions or behaviour.", models.DifficultyIntermediate, 6.0, "impact", "media influence")}
	default:
		return fallbackAdvancedWords(topicSlug)
	}
}

func fallbackAdvancedWords(topicSlug string) []wordSeed {
	wordsByTopic := map[string][]string{
		"science":       {"hypothesis", "experiment", "evidence", "analysis", "discovery", "innovation", "ethics", "genetics", "astronomy", "methodology"},
		"government":    {"policy", "legislation", "election", "regulation", "taxation", "governance", "bureaucracy", "referendum", "diplomacy", "accountability"},
		"global-issues": {"migration", "poverty", "conflict", "cooperation", "aid", "inequality", "security", "sanitation", "resilience", "humanitarian"},
		"economics":     {"inflation", "employment", "productivity", "investment", "recession", "subsidy", "consumption", "inequality", "market", "growth"},
		"law":           {"legislation", "justice", "prosecution", "verdict", "contract", "liability", "evidence", "jurisdiction", "compliance", "rights"},
		"culture":       {"heritage", "tradition", "identity", "ritual", "diversity", "custom", "expression", "preservation", "museum", "festival"},
		"innovation":    {"prototype", "automation", "breakthrough", "research", "startup", "patent", "efficiency", "invention", "scalability", "adaptation"},
		"psychology":    {"motivation", "memory", "behaviour", "stress", "personality", "perception", "resilience", "cognition", "empathy", "habit"},
	}

	baseWords := wordsByTopic[topicSlug]
	words := make([]wordSeed, 0, len(baseWords))
	for index, word := range baseWords {
		difficulty := models.DifficultyIntermediate
		targetBand := 6.5
		if index%3 == 0 {
			difficulty = models.DifficultyAdvanced
			targetBand = 7.5
		}
		words = append(words, w(word, fmt.Sprintf("/%s/", slugify(word)), "noun", fmt.Sprintf("thuật ngữ học thuật về %s", word), fmt.Sprintf("An academic term used to discuss %s in IELTS topics.", word), difficulty, targetBand, "concept", fmt.Sprintf("%s topic", word)))
	}

	return words
}

func w(word string, ipa string, partOfSpeech string, meaningVI string, meaningEN string, difficulty models.DifficultyLevel, targetBand float64, synonym string, collocation string) wordSeed {
	return wordSeed{
		Word:         word,
		IPA:          ipa,
		PartOfSpeech: partOfSpeech,
		MeaningVI:    meaningVI,
		MeaningEN:    meaningEN,
		Difficulty:   difficulty,
		TargetBand:   targetBand,
		Synonyms:     []string{synonym},
		Collocations: []string{collocation, fmt.Sprintf("%s vocabulary", word)},
	}
}

func exampleForWord(word string, topicTitle string) string {
	return fmt.Sprintf("Candidates often use \"%s\" when explaining %s issues in IELTS Writing Task 2.", word, topicTitle)
}

func sortedTopics(topics map[string]models.Topic) []models.Topic {
	values := make([]models.Topic, 0, len(topics))
	for _, topic := range topics {
		values = append(values, topic)
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].BandLevelID != values[j].BandLevelID {
			return values[i].BandLevelID < values[j].BandLevelID
		}
		return values[i].OrderIndex < values[j].OrderIndex
	})
	return values
}

func sortedLessons(lessons map[string]models.Lesson) []models.Lesson {
	values := make([]models.Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		values = append(values, lesson)
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].TopicID != values[j].TopicID {
			return values[i].TopicID < values[j].TopicID
		}
		return values[i].OrderIndex < values[j].OrderIndex
	})
	return values
}

func bandRangeForTopic(topicSlug string) (float64, float64) {
	switch topicSlug {
	case "education", "health", "environment":
		return 5.0, 5.5
	case "technology", "business", "society", "media":
		return 6.0, 6.5
	case "science", "government", "global-issues", "economics":
		return 7.0, 7.5
	default:
		return 8.0, 8.5
	}
}

func scoreForStatus(status models.LessonStatus) (*int, *int) {
	switch status {
	case models.LessonStatusCompleted:
		score := 90
		bestScore := 95
		return &score, &bestScore
	case models.LessonStatusInProgress:
		score := 45
		bestScore := 45
		return &score, &bestScore
	default:
		return nil, nil
	}
}

func progressWords(total int, percent int) int {
	if total == 0 {
		return 0
	}
	learned := total * percent / 100
	if learned == 0 {
		return 1
	}
	return learned
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func demoNow() time.Time {
	return time.Date(2026, 6, 6, 10, 0, 0, 0, time.UTC)
}

func mustJSON(value []string) datatypes.JSON {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return datatypes.JSON(data)
}

func slugify(value string) string {
	slug := ""
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z':
			slug += string(char)
		case char >= 'A' && char <= 'Z':
			slug += string(char + 32)
		case char >= '0' && char <= '9':
			slug += string(char)
		case char == ' ' || char == '-' || char == '_' || char == '&':
			if len(slug) > 0 && slug[len(slug)-1] != '-' {
				slug += "-"
			}
		}
	}

	if len(slug) > 0 && slug[len(slug)-1] == '-' {
		return slug[:len(slug)-1]
	}

	return slug
}
