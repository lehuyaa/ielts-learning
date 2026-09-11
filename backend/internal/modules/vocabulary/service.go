package vocabulary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"ielts-learning/backend/internal/models"
)

var ErrImportFileInvalid = errors.New("invalid import file")

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

// Import parses an Excel file of vocabulary rows, validates each row
// (structurally, then against topic/lesson slugs and in-file duplicates),
// and, unless dryRun is set, commits every valid row in one transaction:
// upserting the Vocabulary by slug and linking it to the given lesson.
func (s Service) Import(reader io.Reader, dryRun bool) (ImportResultResponse, error) {
	rows, err := ParseImportExcel(reader)
	if err != nil {
		return ImportResultResponse{}, fmt.Errorf("%w: %s", ErrImportFileInvalid, err.Error())
	}

	topicSlugs := uniqueNonEmptyStrings(rowTopicSlugs(rows))
	topicsBySlug, err := s.repository.FindTopicsBySlugs(topicSlugs)
	if err != nil {
		return ImportResultResponse{}, err
	}

	topicIDs := make([]uint, 0, len(topicsBySlug))
	for _, topic := range topicsBySlug {
		topicIDs = append(topicIDs, topic.ID)
	}

	lessonSlugs := uniqueNonEmptyStrings(rowLessonSlugs(rows))
	lessons, err := s.repository.FindLessonsByTopicIDsAndSlugs(topicIDs, lessonSlugs)
	if err != nil {
		return ImportResultResponse{}, err
	}

	lessonsByTopicAndSlug := make(map[string]models.Lesson, len(lessons))
	for _, lesson := range lessons {
		lessonsByTopicAndSlug[lessonLookupKey(lesson.TopicID, lesson.Slug)] = lesson
	}

	seenRowByKey := make(map[string]int, len(rows))
	for i := range rows {
		row := &rows[i]
		if !row.Valid() {
			continue
		}

		topic, topicFound := topicsBySlug[row.TopicSlug]
		if !topicFound {
			row.addError(fmt.Sprintf("topic with slug %q was not found", row.TopicSlug))
			continue
		}
		row.TopicID = topic.ID

		lesson, lessonFound := lessonsByTopicAndSlug[lessonLookupKey(topic.ID, row.LessonSlug)]
		if !lessonFound {
			row.addError(fmt.Sprintf("lesson with slug %q was not found under topic %q", row.LessonSlug, row.TopicSlug))
			continue
		}
		row.LessonID = lesson.ID

		dedupeKey := lessonLookupKey(lesson.ID, row.Slug)
		if firstRowNumber, exists := seenRowByKey[dedupeKey]; exists {
			row.addError(fmt.Sprintf("word %q for this lesson is duplicated with row %d", row.Word, firstRowNumber))
			continue
		}
		seenRowByKey[dedupeKey] = row.RowNumber
	}

	result := ImportResultResponse{DryRun: dryRun}

	if !dryRun {
		stats, err := s.repository.CommitImportRows(rows)
		if err != nil {
			return ImportResultResponse{}, fmt.Errorf("commit import: %w", err)
		}
		result.VocabulariesCreated = stats.VocabulariesCreated
		result.VocabulariesUpdated = stats.VocabulariesUpdated
		result.LessonLinksCreated = stats.LessonLinksCreated
		result.LessonLinksUpdated = stats.LessonLinksUpdated
	}

	result.Summary = buildImportSummary(rows)
	result.Rows = importRowResults(rows)

	return result, nil
}

func rowTopicSlugs(rows []ImportRow) []string {
	values := make([]string, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.TopicSlug)
	}
	return values
}

func rowLessonSlugs(rows []ImportRow) []string {
	values := make([]string, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.LessonSlug)
	}
	return values
}

func uniqueNonEmptyStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func lessonLookupKey(topicOrLessonID uint, slug string) string {
	return fmt.Sprintf("%d|%s", topicOrLessonID, slug)
}

func buildImportSummary(rows []ImportRow) ImportSummary {
	summary := ImportSummary{TotalRows: len(rows)}
	for _, row := range rows {
		if row.Valid() {
			summary.ValidRows++
		} else {
			summary.InvalidRows++
		}
	}
	return summary
}

func importRowResults(rows []ImportRow) []ImportRowResult {
	results := make([]ImportRowResult, 0, len(rows))
	for _, row := range rows {
		results = append(results, ImportRowResult{
			Row:              row.RowNumber,
			Word:             row.Word,
			Slug:             row.Slug,
			TopicSlug:        row.TopicSlug,
			LessonSlug:       row.LessonSlug,
			Valid:            row.Valid(),
			Errors:           row.Errors,
			VocabularyAction: row.VocabularyAction,
			LinkAction:       row.LinkAction,
		})
	}
	return results
}

func (s Service) List(userID uint, query ListQuery) (ListResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return ListResponse{}, err
	}

	vocabularies, total, err := s.repository.FindVocabularies(userID, query)
	if err != nil {
		return ListResponse{}, err
	}

	vocabularyIDs := vocabularyIDs(vocabularies)
	progressRows, err := s.repository.FindVocabularyProgress(userID, vocabularyIDs)
	if err != nil {
		return ListResponse{}, err
	}

	contexts, err := s.repository.FindTopicContext(vocabularyIDs)
	if err != nil {
		return ListResponse{}, err
	}

	items := listItemResponses(
		vocabularies,
		mapProgressByVocabulary(progressRows),
		contexts,
	)

	return ListResponse{
		Items: items,
		Pagination: PaginationResponse{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages(total, query.Limit),
		},
	}, nil
}

func (s Service) Get(userID uint, vocabularyID uint) (DetailResponse, error) {
	if _, err := s.repository.FindUser(userID); err != nil {
		return DetailResponse{}, err
	}

	vocabulary, err := s.repository.FindVocabulary(vocabularyID)
	if err != nil {
		return DetailResponse{}, err
	}

	progressRows, err := s.repository.FindVocabularyProgress(userID, []uint{vocabulary.ID})
	if err != nil {
		return DetailResponse{}, err
	}

	contexts, err := s.repository.FindTopicContext([]uint{vocabulary.ID})
	if err != nil {
		return DetailResponse{}, err
	}

	progressByVocabulary := mapProgressByVocabulary(progressRows)
	progress, hasProgress := progressByVocabulary[vocabulary.ID]
	context, hasContext := contexts[vocabulary.ID]

	return detailResponse(vocabulary, progress, hasProgress, context, hasContext), nil
}

func listItemResponses(
	vocabularies []models.Vocabulary,
	progressByVocabulary map[uint]models.UserVocabularyProgress,
	contexts map[uint]VocabularyTopicContext,
) []ListItemResponse {
	responses := make([]ListItemResponse, 0, len(vocabularies))
	for _, vocabulary := range vocabularies {
		progress, hasProgress := progressByVocabulary[vocabulary.ID]
		context, hasContext := contexts[vocabulary.ID]
		summary := progressSummary(progress, hasProgress)

		responses = append(responses, ListItemResponse{
			ID:              vocabulary.ID,
			Word:            vocabulary.Word,
			Slug:            vocabulary.Slug,
			IPA:             vocabulary.IPA,
			PartOfSpeech:    vocabulary.PartOfSpeech,
			MeaningVI:       vocabulary.MeaningVI,
			MeaningEN:       vocabulary.MeaningEN,
			ShortDefinition: vocabulary.ShortDefinition,
			Difficulty:      vocabulary.Difficulty,
			TargetBand:      vocabulary.TargetBand,
			Topic:           topicResponse(context, hasContext),
			Progress:        summary,
			Status:          summary.Status,
			MasteryScore:    summary.MasteryScore,
		})
	}

	return responses
}

func detailResponse(
	vocabulary models.Vocabulary,
	progress models.UserVocabularyProgress,
	hasProgress bool,
	context VocabularyTopicContext,
	hasContext bool,
) DetailResponse {
	detailProgress := progressDetail(progress, hasProgress)
	primaryMeaning := firstNonEmpty(vocabulary.MeaningEN, vocabulary.ShortDefinition, vocabulary.MeaningVI)
	secondaryMeaning := ""
	if primaryMeaning != vocabulary.MeaningVI {
		secondaryMeaning = vocabulary.MeaningVI
	}

	return DetailResponse{
		ID:               vocabulary.ID,
		Word:             vocabulary.Word,
		Slug:             vocabulary.Slug,
		IPA:              vocabulary.IPA,
		AudioURL:         vocabulary.AudioURL,
		PartOfSpeech:     vocabulary.PartOfSpeech,
		MeaningVI:        vocabulary.MeaningVI,
		MeaningEN:        vocabulary.MeaningEN,
		PrimaryMeaning:   primaryMeaning,
		SecondaryMeaning: secondaryMeaning,
		ExampleSentences: compactStrings([]string{vocabulary.ExampleSentence}),
		Synonyms:         stringSlice(vocabulary.SynonymsJSON),
		Antonyms:         stringSlice(vocabulary.AntonymsJSON),
		Collocations:     stringSlice(vocabulary.CollocationsJSON),
		IELTSUsage:       vocabulary.IELTSUsage,
		RelatedForms:     relatedForms(vocabulary),
		Difficulty:       vocabulary.Difficulty,
		TargetBand:       vocabulary.TargetBand,
		Frequency:        frequencyLabel(vocabulary.TargetBand),
		Rating:           ratingFromMastery(detailProgress.MasteryScore),
		Topic:            topicResponse(context, hasContext),
		BandLevel:        bandLevelResponse(context, hasContext),
		Progress:         detailProgress,
		UserProgress:     detailProgress,
		MasteryScore:     detailProgress.MasteryScore,
	}
}

func vocabularyIDs(vocabularies []models.Vocabulary) []uint {
	ids := make([]uint, 0, len(vocabularies))
	for _, vocabulary := range vocabularies {
		ids = append(ids, vocabulary.ID)
	}
	return ids
}

func mapProgressByVocabulary(rows []models.UserVocabularyProgress) map[uint]models.UserVocabularyProgress {
	result := make(map[uint]models.UserVocabularyProgress, len(rows))
	for _, row := range rows {
		result[row.VocabularyID] = row
	}
	return result
}

func progressSummary(progress models.UserVocabularyProgress, hasProgress bool) ProgressSummaryResponse {
	status := models.VocabularyStatusNew
	if hasProgress {
		status = progress.Status
	}
	masteryScore := masteryScore(progress, hasProgress)

	return ProgressSummaryResponse{
		Status:       status,
		ReviewCount:  progress.ReviewCount,
		NextReviewAt: progress.NextReviewAt,
		MasteryScore: masteryScore,
	}
}

func progressDetail(progress models.UserVocabularyProgress, hasProgress bool) ProgressDetailResponse {
	status := models.VocabularyStatusNew
	if hasProgress {
		status = progress.Status
	}

	return ProgressDetailResponse{
		Status:         status,
		ReviewCount:    progress.ReviewCount,
		CorrectCount:   progress.CorrectCount,
		WrongCount:     progress.WrongCount,
		NextReviewAt:   progress.NextReviewAt,
		LastReviewedAt: progress.LastReviewedAt,
		LearnedAt:      progress.LearnedAt,
		MasteryScore:   masteryScore(progress, hasProgress),
	}
}

func masteryScore(progress models.UserVocabularyProgress, hasProgress bool) int {
	if !hasProgress {
		return 0
	}
	if progress.Status == models.VocabularyStatusMastered {
		return 100
	}

	totalAnswers := progress.CorrectCount + progress.WrongCount
	if totalAnswers > 0 {
		return clampPercentage((progress.CorrectCount * 100) / totalAnswers)
	}

	switch progress.Status {
	case models.VocabularyStatusReview:
		return 70
	case models.VocabularyStatusLearning:
		return 40
	default:
		return 0
	}
}

func topicResponse(context VocabularyTopicContext, hasContext bool) *TopicResponse {
	if !hasContext {
		return nil
	}

	return &TopicResponse{
		ID:    context.TopicID,
		Title: context.TopicTitle,
		Slug:  context.TopicSlug,
		Icon:  firstNonEmpty(context.TopicIcon, context.TopicEmoji),
		Emoji: context.TopicEmoji,
		Color: context.TopicColor,
	}
}

func bandLevelResponse(context VocabularyTopicContext, hasContext bool) *BandLevelResponse {
	if !hasContext {
		return nil
	}

	return &BandLevelResponse{
		ID:        context.BandLevelID,
		BandScore: context.BandScore,
		Title:     context.BandTitle,
	}
}

func stringSlice(data []byte) []string {
	if len(data) == 0 {
		return []string{}
	}

	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return []string{}
	}

	return values
}

func relatedForms(vocabulary models.Vocabulary) []RelatedFormResponse {
	// The current schema does not store related word forms yet.
	return []RelatedFormResponse{}
}

func compactStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func frequencyLabel(targetBand *float64) string {
	if targetBand == nil {
		return "Medium"
	}
	switch {
	case *targetBand >= 7.0:
		return "Very High"
	case *targetBand >= 6.0:
		return "High"
	default:
		return "Medium"
	}
}

func ratingFromMastery(masteryScore int) int {
	switch {
	case masteryScore >= 90:
		return 5
	case masteryScore >= 70:
		return 4
	case masteryScore >= 45:
		return 3
	case masteryScore > 0:
		return 2
	default:
		return 1
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

func totalPages(total int64, limit int) int {
	if total == 0 {
		return 0
	}

	return int((total + int64(limit) - 1) / int64(limit))
}

func clampPercentage(value int) int {
	switch {
	case value < 0:
		return 0
	case value > 100:
		return 100
	default:
		return value
	}
}
