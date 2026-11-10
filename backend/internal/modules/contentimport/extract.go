package contentimport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"ielts-learning/backend/internal/models"
	"ielts-learning/backend/internal/shared/ai"
)

// UploadedFile is one page image or PDF sent for extraction.
type UploadedFile struct {
	Name string
	Data []byte
}

// ExtractRequest describes where the extracted lesson goes.
type ExtractRequest struct {
	CourseSlug string
	CEFRLevel  string
	TopicSlug  string
	TopicTitle string
	Notes      string
	// Publish makes the activities visible to learners right away (default: DRAFT).
	Publish bool
	Files   []UploadedFile
}

// ExtractResult is a bundle draft (spec §8.2) to review in the import preview.
type ExtractResult struct {
	Bundle   json.RawMessage `json:"bundle" swaggertype:"object"`
	Summary  string          `json:"summary"`
	Warnings []string        `json:"warnings"`
}

const (
	maxExtractImages = 10
	maxExtractBytes  = 30 << 20
)

// What the vision model returns. Payload and answer are JSON strings because
// strict JSON schemas cannot describe free-form objects.
type extraction struct {
	LessonTitle   string `json:"lessonTitle"`
	GrammarPoints []struct {
		Slug      string `json:"slug"`
		Title     string `json:"title"`
		TitleVI   string `json:"titleVi"`
		SummaryVI string `json:"summaryVi"`
		Formula   string `json:"formula"`
		Examples  []struct {
			EN string `json:"en"`
			VI string `json:"vi"`
		} `json:"examples"`
		CommonMistakes []struct {
			Wrong  string `json:"wrong"`
			Right  string `json:"right"`
			NoteVI string `json:"noteVi"`
		} `json:"commonMistakes"`
	} `json:"grammarPoints"`
	Activities []struct {
		Type             string `json:"type"`
		Title            string `json:"title"`
		InstructionsVI   string `json:"instructionsVi"`
		GrammarPointSlug string `json:"grammarPointSlug"`
		ReadingText      string `json:"readingText"`
		Exercises        []struct {
			Type          string `json:"type"`
			Prompt        string `json:"prompt"`
			PayloadJSON   string `json:"payloadJson"`
			AnswerJSON    string `json:"answerJson"`
			ExplanationVI string `json:"explanationVi"`
			SkillTag      string `json:"skillTag"`
		} `json:"exercises"`
	} `json:"activities"`
}

var extractExerciseTypes = []string{
	"MULTIPLE_CHOICE", "TRUE_FALSE", "FILL_BLANK", "REORDER", "ERROR_CORRECTION", "TRANSFORM", "TRANSLATE_VI_EN", "IPA_TO_WORD",
}

func strictObject(properties map[string]any) map[string]any {
	required := make([]string, 0, len(properties))
	for name := range properties {
		required = append(required, name)
	}
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}
}

func arrayOf(item map[string]any) map[string]any {
	return map[string]any{"type": "array", "items": item}
}

var str = map[string]any{"type": "string"}

var extractSchema = ai.Schema{
	Name: "lesson_extraction",
	Schema: strictObject(map[string]any{
		"lessonTitle": str,
		"grammarPoints": arrayOf(strictObject(map[string]any{
			"slug": str, "title": str, "titleVi": str, "summaryVi": str, "formula": str,
			"examples":       arrayOf(strictObject(map[string]any{"en": str, "vi": str})),
			"commonMistakes": arrayOf(strictObject(map[string]any{"wrong": str, "right": str, "noteVi": str})),
		})),
		"activities": arrayOf(strictObject(map[string]any{
			"type":             map[string]any{"type": "string", "enum": []string{"GRAMMAR", "READING", "PRONUNCIATION"}},
			"title":            str,
			"instructionsVi":   str,
			"grammarPointSlug": str,
			"readingText":      str,
			"exercises": arrayOf(strictObject(map[string]any{
				"type":          map[string]any{"type": "string", "enum": extractExerciseTypes},
				"prompt":        str,
				"payloadJson":   str,
				"answerJson":    str,
				"explanationVi": str,
				"skillTag":      str,
			})),
		})),
	}),
}

const extractSystemPrompt = `You turn pages of an English textbook for Vietnamese adult learners into lesson data for a learning app.
Extract the grammar explanations and the exercises you can see. Do not invent exercises that are not on the pages, except that you may split one book exercise into several items (one item per question).
Write instructions, explanations and grammar summaries in Vietnamese; keep all English content exactly as printed.

grammarPoints: one per grammar topic explained on the pages. slug like "g3-present-simple" (lowercase, hyphens). summaryVi: when and how to use it, 2-4 sentences. formula: e.g. "S + V(s/es) + O". 3-5 examples, up to 5 common mistakes Vietnamese learners make.
activities: GRAMMAR for grammar exercises (grammarPointSlug = the slug it practises, or ""), READING when the pages contain a text with questions (put the full text in readingText), PRONUNCIATION for IPA exercises. Use "" for unused string fields.

Each exercise: type, prompt (short Vietnamese instruction), payloadJson and answerJson as JSON strings, explanationVi (why the answer is right), skillTag ("grammar:<slug>", "reading" or "sound:<ipa>").
Formats (payloadJson / answerJson):
- MULTIPLE_CHOICE: {"options":[{"id":"a","text":"..."},...]} / {"correct":"a"}
- TRUE_FALSE: {"statement":"..."} / {"value":true}
- FILL_BLANK: {"sentence":"She ___ (work) in a bank."} / {"accepted":["works"]}  (exactly one ___)
- REORDER: {"tokens":["name","is","My","Huy"]} (shuffled words) / {"accepted":["My name is Huy."]}
- ERROR_CORRECTION: {"tokens":["He","work","here","."]} / {"wrongIndex":1,"accepted":["works"]}
- TRANSFORM: {"sentence":"I like coffee.","instruction":"NEGATIVE"} (NEGATIVE|QUESTION|PAST|THIRD_PERSON) / {"accepted":["I don't like coffee.","I do not like coffee."],"aiFallback":true}
- TRANSLATE_VI_EN: {"sourceVi":"..."} / {"reference":["..."],"aiGraded":true}
- IPA_TO_WORD: {"ipa":"/wɜːrk/","options":[{"id":"a","word":"work"},...]} / {"correct":"a"}
Use the answer key on the pages when there is one; otherwise give the correct answer yourself and list all acceptable variants.`

// Extract sends page images or a PDF to the vision model and turns the reply
// into a bundle draft. Every activity is DRAFT, so nothing reaches learners
// before it is reviewed and published.
func (s Service) Extract(ctx context.Context, request ExtractRequest) (ExtractResult, error) {
	if s.ai == nil || !s.ai.Configured() {
		return ExtractResult{}, ErrAIUnavailable
	}
	parts, err := extractParts(request)
	if err != nil {
		return ExtractResult{}, err
	}

	var reply extraction
	if err := s.ai.CompleteJSONWithParts(ctx, extractSystemPrompt, parts, extractSchema, &reply); err != nil {
		return ExtractResult{}, fmt.Errorf("%w: %v", ErrAIUnavailable, err)
	}
	result, err := buildExtractBundle(request, reply)
	if err != nil {
		return ExtractResult{}, err
	}

	// Importing would overwrite theory cards that already exist.
	var slugs []string
	for _, point := range reply.GrammarPoints {
		slugs = append(slugs, slugify(point.Slug))
	}
	existing, err := s.repository.GrammarPointIDs(slugs)
	if err != nil {
		return ExtractResult{}, err
	}
	for _, slug := range slugs {
		if _, ok := existing[slug]; ok {
			result.Warnings = append(result.Warnings, fmt.Sprintf(
				"grammar point %q already exists — importing this draft replaces its theory card; delete it from \"grammarPoints\" to keep the current one", slug))
		}
	}
	return result, nil
}

func extractParts(request ExtractRequest) ([]ai.Part, error) {
	if len(request.Files) == 0 {
		return nil, fmt.Errorf("%w: upload page images or one PDF", ErrExtractInput)
	}
	total := 0
	pdfs := 0
	parts := []ai.Part{ai.TextPart(fmt.Sprintf(
		"CEFR level: %s. Lesson topic: %s. %s", defaultString(request.CEFRLevel, "A1"), request.TopicTitle, request.Notes))}
	for _, file := range request.Files {
		total += len(file.Data)
		mimeType := http.DetectContentType(file.Data)
		switch {
		case mimeType == "application/pdf" || strings.EqualFold(filepath.Ext(file.Name), ".pdf"):
			pdfs++
			parts = append(parts, ai.PDFPart(filepath.Base(file.Name), file.Data))
		case strings.HasPrefix(mimeType, "image/"):
			parts = append(parts, ai.ImagePart(mimeType, file.Data))
		default:
			return nil, fmt.Errorf("%w: %s is not an image or PDF", ErrExtractInput, file.Name)
		}
	}
	if pdfs > 1 || (pdfs == 1 && len(request.Files) > 1) {
		return nil, fmt.Errorf("%w: upload one PDF, or up to %d images", ErrExtractInput, maxExtractImages)
	}
	if len(request.Files) > maxExtractImages {
		return nil, fmt.Errorf("%w: at most %d images", ErrExtractInput, maxExtractImages)
	}
	if total > maxExtractBytes {
		return nil, fmt.Errorf("%w: files are too large (max 30MB in total)", ErrExtractInput)
	}
	return parts, nil
}

func buildExtractBundle(request ExtractRequest, reply extraction) (ExtractResult, error) {
	var warnings []string
	topicSlug := slugify(request.TopicSlug)
	if topicSlug == "" {
		topicSlug = "extracted-" + slugify(defaultString(request.TopicTitle, reply.LessonTitle))
	}
	lessonTitle := defaultString(reply.LessonTitle, request.TopicTitle)
	lessonSlug := topicSlug + "-" + defaultString(slugify(lessonTitle), "lesson")

	bundle := map[string]any{
		"version":    1,
		"courseSlug": defaultString(request.CourseSlug, "foundation-english"),
		"cefrLevel":  defaultString(request.CEFRLevel, "A1"),
	}

	grammarPoints := []map[string]any{}
	for i, point := range reply.GrammarPoints {
		examples := []map[string]string{}
		for _, example := range point.Examples {
			examples = append(examples, map[string]string{"en": example.EN, "vi": example.VI})
		}
		mistakes := []map[string]string{}
		for _, mistake := range point.CommonMistakes {
			mistakes = append(mistakes, map[string]string{"wrong": mistake.Wrong, "right": mistake.Right, "noteVi": mistake.NoteVI})
		}
		grammarPoints = append(grammarPoints, map[string]any{
			"slug": slugify(point.Slug), "cefrLevel": bundle["cefrLevel"], "title": point.Title, "titleVi": point.TitleVI,
			"summaryVi": point.SummaryVI, "formula": point.Formula, "examples": examples, "commonMistakes": mistakes,
			"orderIndex": 100 + i,
		})
	}
	if len(grammarPoints) > 0 {
		bundle["grammarPoints"] = grammarPoints
	}

	status := string(models.ContentDraft)
	if request.Publish {
		status = string(models.ContentPublished)
	}
	passages := []map[string]any{}
	activities := []map[string]any{}
	exerciseCount := 0
	for i, activity := range reply.Activities {
		slug := fmt.Sprintf("%s-%d", strings.ToLower(activity.Type), i+1)
		item := map[string]any{
			"slug": slug, "type": activity.Type, "title": defaultString(activity.Title, activity.Type),
			"instructionsVi": activity.InstructionsVI, "status": status,
		}
		if activity.GrammarPointSlug != "" {
			item["grammarPointSlug"] = slugify(activity.GrammarPointSlug)
		}
		if strings.TrimSpace(activity.ReadingText) != "" {
			passageSlug := fmt.Sprintf("%s-reading-%d", lessonSlug, i+1)
			passages = append(passages, map[string]any{
				"slug": passageSlug, "genre": string(models.GenreArticle), "title": item["title"],
				"bodyText": strings.TrimSpace(activity.ReadingText), "source": "extracted from uploaded pages",
			})
			item["passageSlug"] = passageSlug
			item["content"] = map[string]any{"tapToTranslate": true}
		}

		exercises := []map[string]any{}
		for j, exercise := range activity.Exercises {
			payload, payloadErr := parseObject(exercise.PayloadJSON)
			answer, answerErr := parseObject(exercise.AnswerJSON)
			if payloadErr != nil || answerErr != nil {
				warnings = append(warnings, fmt.Sprintf("%s #%d: the AI returned invalid JSON; fix it before importing", slug, j+1))
			}
			exercises = append(exercises, map[string]any{
				"type": exercise.Type, "prompt": exercise.Prompt, "payload": payload, "answer": answer,
				"explanationVi": exercise.ExplanationVI, "skillTag": exercise.SkillTag,
			})
		}
		exerciseCount += len(exercises)
		item["exercises"] = exercises
		activities = append(activities, item)
	}
	if len(passages) > 0 {
		bundle["passages"] = passages
	}
	bundle["topics"] = []map[string]any{{
		"slug": topicSlug, "title": defaultString(request.TopicTitle, lessonTitle), "track": string(models.TopicTrackVocabulary),
		"lessons": []map[string]any{{"slug": lessonSlug, "title": lessonTitle, "activities": activities}},
	}}
	if len(activities) == 0 {
		warnings = append(warnings, "no exercises were found on these pages")
	}

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return ExtractResult{}, fmt.Errorf("encode bundle: %w", err)
	}
	return ExtractResult{
		Bundle:   data,
		Summary:  fmt.Sprintf("%d grammar points, %d activities, %d exercises (%s)", len(grammarPoints), len(activities), exerciseCount, status),
		Warnings: append([]string{}, warnings...),
	}, nil
}

func parseObject(raw string) (map[string]any, error) {
	object := map[string]any{}
	if strings.TrimSpace(raw) == "" {
		return object, fmt.Errorf("empty")
	}
	if err := json.Unmarshal([]byte(raw), &object); err != nil {
		return map[string]any{}, err
	}
	return object, nil
}
