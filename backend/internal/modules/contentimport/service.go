package contentimport

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

var ErrInvalidBundle = errors.New("invalid content bundle")

const (
	defaultLessonXP       = 50
	defaultLessonMinutes  = 20
	defaultRequiredScore  = 70
	defaultPassScore      = 70
	defaultActivityXP     = 10
	defaultExercisePoints = 10
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

// ParseBundle decodes a bundle file, rejecting unknown fields so typos in
// field names are caught instead of silently ignored.
func ParseBundle(data []byte) (Bundle, error) {
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var bundle Bundle
	if err := decoder.Decode(&bundle); err != nil {
		return Bundle{}, fmt.Errorf("%w: %v", ErrInvalidBundle, err)
	}
	return bundle, nil
}

// Import validates a bundle and, when commit is true, writes every valid item
// in one transaction. Invalid items (and their children) are skipped.
func (s Service) Import(data []byte, commit bool) (ImportResult, error) {
	bundle, err := ParseBundle(data)
	if err != nil {
		return ImportResult{}, err
	}

	result := ImportResult{BundleErrors: []string{}, Warnings: []string{}, Items: []ItemResult{}}
	for name, raw := range map[string]json.RawMessage{
		"grammarPoints": bundle.GrammarPoints, "phonemes": bundle.Phonemes,
		"passages": bundle.Passages, "scenarios": bundle.Scenarios,
	} {
		if len(raw) > 0 && string(raw) != "[]" && string(raw) != "null" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s are not imported yet and were ignored", name))
		}
	}

	if bundle.CourseSlug == "" || bundle.CEFRLevel == "" {
		result.BundleErrors = append(result.BundleErrors, "courseSlug and cefrLevel are required")
		return result, nil
	}
	band, err := s.repository.FindBandLevel(bundle.CourseSlug, bundle.CEFRLevel)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result.BundleErrors = append(result.BundleErrors, fmt.Sprintf(
			"level %s of course %q not found — run `go run ./cmd/seed -only-foundation` first", bundle.CEFRLevel, bundle.CourseSlug))
		return result, nil
	}
	if err != nil {
		return ImportResult{}, err
	}

	knownVocabulary, err := s.repository.ExistingVocabularySlugs(collectVocabularySlugs(bundle))
	if err != nil {
		return ImportResult{}, err
	}

	run := func(repository Repository) error {
		return s.apply(repository, band, bundle, knownVocabulary, commit, &result)
	}
	if commit {
		err = s.repository.Transaction(func(tx Repository) error {
			if err := run(tx); err != nil {
				return err
			}
			return tx.UpdateCourseTotals(band.CourseID)
		})
	} else {
		err = run(s.repository)
	}
	if err != nil {
		return ImportResult{}, err
	}

	result.Committed = commit
	return result, nil
}

func (s Service) apply(repo Repository, band models.BandLevel, bundle Bundle, knownVocabulary map[string]bool, commit bool, result *ImportResult) error {
	topicSlugs := map[string]bool{}
	for topicIndex, topicInput := range bundle.Topics {
		item := ItemResult{Kind: KindTopic, Path: topicInput.Slug, Title: topicInput.Title, Errors: validateTopic(topicInput)}
		if topicSlugs[topicInput.Slug] {
			item.Errors = append(item.Errors, "duplicate topic slug in this file")
		}
		topicSlugs[topicInput.Slug] = true

		topic, exists, err := repo.FindTopic(band.ID, topicInput.Slug)
		if err != nil {
			return err
		}
		if !s.record(result, &item, exists) {
			s.skipChildren(result, topicInput.Slug, topicInput.Lessons)
			continue
		}
		result.Summary.Topics++

		topic.BandLevelID = band.ID
		topic.Slug = topicInput.Slug
		topic.Title = topicInput.Title
		topic.Description = topicInput.Description
		topic.Icon = topicInput.Icon
		topic.Emoji = topicInput.Emoji
		topic.Color = topicInput.Color
		topic.Track = models.TopicTrack(defaultString(topicInput.Track, string(models.TopicTrackCore)))
		topic.OrderIndex = defaultInt(topicInput.OrderIndex, topicIndex+1)
		if commit {
			if err := repo.Save(&topic); err != nil {
				return err
			}
		}

		if err := s.applyLessons(repo, topic, exists, topicInput, knownVocabulary, commit, result); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) applyLessons(repo Repository, topic models.Topic, topicExists bool, topicInput TopicInput, knownVocabulary map[string]bool, commit bool, result *ImportResult) error {
	lessonSlugs := map[string]bool{}
	for lessonIndex, lessonInput := range topicInput.Lessons {
		path := topicInput.Slug + "/" + lessonInput.Slug
		item := ItemResult{Kind: KindLesson, Path: path, Title: lessonInput.Title, Errors: validateLesson(lessonInput)}
		if lessonSlugs[lessonInput.Slug] {
			item.Errors = append(item.Errors, "duplicate lesson slug in this topic")
		}
		lessonSlugs[lessonInput.Slug] = true

		var lesson models.Lesson
		exists := false
		if topicExists || commit {
			var err error
			lesson, exists, err = repo.FindLesson(topic.ID, lessonInput.Slug)
			if err != nil {
				return err
			}
		}
		if !s.record(result, &item, exists) {
			s.skipActivities(result, path, lessonInput.Activities)
			continue
		}
		result.Summary.Lessons++

		bandMin, bandMax := 1.0, 1.5
		lesson.TopicID = topic.ID
		lesson.Slug = lessonInput.Slug
		lesson.Title = lessonInput.Title
		lesson.Description = lessonInput.Description
		lesson.RequiredScore = defaultInt(lessonInput.RequiredScore, defaultRequiredScore)
		lesson.EstimatedMinutes = defaultInt(lessonInput.EstimatedMinutes, defaultLessonMinutes)
		lesson.XPReward = defaultInt(lessonInput.XPReward, defaultLessonXP)
		lesson.OrderIndex = defaultInt(lessonInput.OrderIndex, lessonIndex+1)
		lesson.IsPublished = true
		if lesson.BandMin == nil {
			lesson.BandMin, lesson.BandMax = &bandMin, &bandMax
		}
		if commit {
			if err := repo.Save(&lesson); err != nil {
				return err
			}
		}

		if err := s.applyActivities(repo, lesson, exists, path, lessonInput.Activities, knownVocabulary, commit, result); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) applyActivities(repo Repository, lesson models.Lesson, lessonExists bool, lessonPath string, activities []ActivityInput, knownVocabulary map[string]bool, commit bool, result *ImportResult) error {
	activitySlugs := map[string]bool{}
	for activityIndex, input := range activities {
		path := lessonPath + "/" + input.Slug
		item := ItemResult{Kind: KindActivity, Path: path, Title: input.Title, Errors: validateActivity(input, knownVocabulary)}
		item.Status = defaultString(input.Status, string(models.ContentPublished))
		if activitySlugs[input.Slug] {
			item.Errors = append(item.Errors, "duplicate activity slug in this lesson")
		}
		activitySlugs[input.Slug] = true

		var activity models.LessonActivity
		exists := false
		if lessonExists || commit {
			var err error
			activity, exists, err = repo.FindActivity(lesson.ID, input.Slug)
			if err != nil {
				return err
			}
		}
		if !s.record(result, &item, exists) {
			s.skipExercises(result, path, input.Exercises)
			continue
		}
		result.Summary.Activities++

		activity.LessonID = lesson.ID
		activity.Slug = input.Slug
		activity.Type = models.ActivityType(input.Type)
		activity.Title = input.Title
		activity.InstructionsVI = input.InstructionsVI
		activity.ContentJSON = nil
		if len(input.Content) > 0 {
			activity.ContentJSON = datatypes.JSON(input.Content)
		}
		activity.PassScore = defaultPtr(input.PassScore, defaultPassScore)
		activity.XPReward = defaultPtr(input.XPReward, defaultActivityXP)
		activity.IsRequired = input.IsRequired == nil || *input.IsRequired
		activity.OrderIndex = activityIndex + 1
		activity.Status = models.ContentStatus(item.Status)
		if commit {
			if err := repo.SaveActivity(&activity); err != nil {
				return err
			}
		}

		if err := s.applyExercises(repo, activity, exists, path, input.Exercises, commit, result); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) applyExercises(repo Repository, activity models.LessonActivity, activityExists bool, activityPath string, exercises []ExerciseInput, commit bool, result *ImportResult) error {
	existing := map[int]uint{}
	if activityExists {
		var err error
		existing, err = repo.ExerciseOrderIndexes(activity.ID)
		if err != nil {
			return err
		}
	}

	for index, input := range exercises {
		orderIndex := index + 1
		item := ItemResult{
			Kind:   KindExercise,
			Path:   fmt.Sprintf("%s#%d", activityPath, orderIndex),
			Title:  fmt.Sprintf("%s · %s", input.Type, input.Prompt),
			Errors: validateExercise(input),
		}
		existingID, exists := existing[orderIndex]
		if !s.record(result, &item, exists) {
			continue
		}
		result.Summary.Exercises++
		if !commit {
			continue
		}

		exercise := models.Exercise{
			ID:            existingID,
			ActivityID:    activity.ID,
			Type:          models.ExerciseType(input.Type),
			Prompt:        input.Prompt,
			PayloadJSON:   datatypes.JSON(input.Payload),
			AnswerJSON:    datatypes.JSON(input.Answer),
			ExplanationVI: input.ExplanationVI,
			SkillTag:      input.SkillTag,
			Points:        defaultInt(input.Points, defaultExercisePoints),
			OrderIndex:    orderIndex,
		}
		if err := repo.SaveExercise(&exercise); err != nil {
			return err
		}
	}

	if commit && activityExists {
		return repo.DeleteExercisesAfter(activity.ID, len(exercises))
	}
	return nil
}

// record adds an item to the result and reports whether it is valid.
func (s Service) record(result *ImportResult, item *ItemResult, exists bool) bool {
	if item.Errors == nil {
		item.Errors = []string{}
	}
	switch {
	case len(item.Errors) > 0:
		item.Action = ActionSkip
		result.Summary.Invalid++
	case exists:
		item.Action = ActionUpdate
		result.Summary.Updated++
	default:
		item.Action = ActionCreate
		result.Summary.Created++
	}
	result.Items = append(result.Items, *item)
	return len(item.Errors) == 0
}

func (s Service) skipChildren(result *ImportResult, topicPath string, lessons []LessonInput) {
	for _, lesson := range lessons {
		path := topicPath + "/" + lesson.Slug
		s.addSkipped(result, KindLesson, path, lesson.Title, "parent topic is invalid")
		s.skipActivities(result, path, lesson.Activities)
	}
}

func (s Service) skipActivities(result *ImportResult, lessonPath string, activities []ActivityInput) {
	for _, activity := range activities {
		path := lessonPath + "/" + activity.Slug
		s.addSkipped(result, KindActivity, path, activity.Title, "parent lesson is invalid")
		s.skipExercises(result, path, activity.Exercises)
	}
}

func (s Service) skipExercises(result *ImportResult, activityPath string, exercises []ExerciseInput) {
	for index, exercise := range exercises {
		s.addSkipped(result, KindExercise, fmt.Sprintf("%s#%d", activityPath, index+1),
			fmt.Sprintf("%s · %s", exercise.Type, exercise.Prompt), "parent activity is invalid")
	}
}

func (s Service) addSkipped(result *ImportResult, kind string, path string, title string, reason string) {
	result.Summary.Invalid++
	result.Items = append(result.Items, ItemResult{Kind: kind, Path: path, Title: title, Action: ActionSkip, Errors: []string{reason}})
}

func collectVocabularySlugs(bundle Bundle) []string {
	var slugs []string
	for _, topic := range bundle.Topics {
		for _, lesson := range topic.Lessons {
			for _, activity := range lesson.Activities {
				slugs = append(slugs, vocabularySlugs(activity.Content)...)
			}
		}
	}
	return slugs
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func defaultInt(value int, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func defaultPtr(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}
