package vocabulary

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"ielts-learning/backend/internal/models"
)

const maxImportRows = 2000

// ImportRow is one parsed and progressively-validated line from an uploaded
// vocabulary import file. Errors accumulate as parsing and, later, database
// lookups run; a row with no errors is safe to commit.
type ImportRow struct {
	RowNumber        int
	Word             string
	Slug             string
	IPA              string
	PartOfSpeech     string
	MeaningVI        string
	MeaningEN        string
	ShortDefinition  string
	ExampleSentence  string
	ExampleMeaningVI string
	ExampleSource    string
	Synonyms         []string
	Antonyms         []string
	Collocations     []string
	Difficulty       models.DifficultyLevel
	TargetBand       *float64
	TopicSlug        string
	LessonSlug       string
	OrderIndex       int
	IsRequired       bool
	TopicID          uint
	LessonID         uint
	Errors           []string
	VocabularyAction string
	LinkAction       string
}

func (row ImportRow) Valid() bool {
	return len(row.Errors) == 0
}

func (row *ImportRow) addError(message string) {
	row.Errors = append(row.Errors, message)
}

func ParseImportExcel(reader io.Reader) ([]ImportRow, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("open excel file: %w", err)
	}
	defer file.Close()

	sheetName := file.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("excel file has no sheets")
	}

	allRows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("read excel rows: %w", err)
	}
	if len(allRows) == 0 {
		return nil, fmt.Errorf("excel file is empty")
	}

	columnIndex := indexHeader(allRows[0])
	if _, ok := columnIndex["word"]; !ok {
		return nil, fmt.Errorf("missing required column: word")
	}

	dataRows := allRows[1:]
	if len(dataRows) > maxImportRows {
		return nil, fmt.Errorf("file has %d data rows, which exceeds the limit of %d", len(dataRows), maxImportRows)
	}

	rows := make([]ImportRow, 0, len(dataRows))
	for i, raw := range dataRows {
		rowNumber := i + 2 // +1 for 1-based indexing, +1 for the header row
		if isBlankRow(raw) {
			continue
		}
		rows = append(rows, parseImportRow(rowNumber, raw, columnIndex))
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("excel file has no data rows")
	}

	return rows, nil
}

func indexHeader(header []string) map[string]int {
	index := make(map[string]int, len(header))
	for i, cell := range header {
		key := normalizeHeaderKey(cell)
		if key == "" {
			continue
		}
		index[key] = i
	}
	return index
}

func normalizeHeaderKey(value string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), " ", ""))
}

func isBlankRow(raw []string) bool {
	for _, value := range raw {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func importCell(raw []string, columnIndex map[string]int, key string) string {
	idx, ok := columnIndex[key]
	if !ok || idx >= len(raw) {
		return ""
	}
	return strings.TrimSpace(raw[idx])
}

func parseImportRow(rowNumber int, raw []string, columnIndex map[string]int) ImportRow {
	row := ImportRow{
		RowNumber:        rowNumber,
		Word:             importCell(raw, columnIndex, "word"),
		IPA:              importCell(raw, columnIndex, "ipa"),
		PartOfSpeech:     importCell(raw, columnIndex, "partofspeech"),
		MeaningVI:        importCell(raw, columnIndex, "meaningvi"),
		MeaningEN:        importCell(raw, columnIndex, "meaningen"),
		ShortDefinition:  importCell(raw, columnIndex, "shortdefinition"),
		ExampleSentence:  importCell(raw, columnIndex, "examplesentence"),
		ExampleMeaningVI: importCell(raw, columnIndex, "examplemeaningvi"),
		ExampleSource:    importCell(raw, columnIndex, "examplesource"),
		TopicSlug:        slugify(importCell(raw, columnIndex, "topicslug")),
		LessonSlug:       slugify(importCell(raw, columnIndex, "lessonslug")),
		IsRequired:       true,
	}

	row.Synonyms = splitImportList(importCell(raw, columnIndex, "synonyms"))
	row.Antonyms = splitImportList(importCell(raw, columnIndex, "antonyms"))
	row.Collocations = splitImportList(importCell(raw, columnIndex, "collocations"))

	if row.Word == "" {
		row.addError("word is required")
	}
	if row.MeaningVI == "" {
		row.addError("meaningVi is required")
	}
	if row.TopicSlug == "" {
		row.addError("topicSlug is required")
	}
	if row.LessonSlug == "" {
		row.addError("lessonSlug is required")
	}

	if slugValue := importCell(raw, columnIndex, "slug"); slugValue != "" {
		row.Slug = slugify(slugValue)
	} else if row.Word != "" {
		row.Slug = slugify(row.Word)
	}
	if row.Slug == "" && row.Word != "" {
		row.addError("word could not be converted to a valid slug")
	}

	if difficultyValue := importCell(raw, columnIndex, "difficulty"); difficultyValue != "" {
		parsed, ok := parseDifficultyValue(difficultyValue)
		if !ok {
			row.addError("difficulty must be BEGINNER, INTERMEDIATE, or ADVANCED")
		} else {
			row.Difficulty = parsed
		}
	} else {
		row.Difficulty = models.DifficultyIntermediate
	}

	if targetBandValue := importCell(raw, columnIndex, "targetband"); targetBandValue != "" {
		parsed, err := strconv.ParseFloat(targetBandValue, 64)
		if err != nil {
			row.addError("targetBand must be a number")
		} else {
			row.TargetBand = &parsed
		}
	}

	if orderIndexValue := importCell(raw, columnIndex, "orderindex"); orderIndexValue != "" {
		parsed, err := strconv.Atoi(orderIndexValue)
		if err != nil {
			row.addError("orderIndex must be an integer")
		} else {
			row.OrderIndex = parsed
		}
	}

	if isRequiredValue := importCell(raw, columnIndex, "isrequired"); isRequiredValue != "" {
		parsed, err := strconv.ParseBool(strings.ToLower(isRequiredValue))
		if err != nil {
			row.addError("isRequired must be true or false")
		} else {
			row.IsRequired = parsed
		}
	}

	return row
}

func parseDifficultyValue(value string) (models.DifficultyLevel, bool) {
	switch models.DifficultyLevel(strings.ToUpper(strings.TrimSpace(value))) {
	case models.DifficultyBeginner:
		return models.DifficultyBeginner, true
	case models.DifficultyIntermediate:
		return models.DifficultyIntermediate, true
	case models.DifficultyAdvanced:
		return models.DifficultyAdvanced, true
	default:
		return "", false
	}
}

func splitImportList(value string) []string {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

var slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	lowered := strings.ToLower(strings.TrimSpace(value))
	slug := slugInvalidChars.ReplaceAllString(lowered, "-")
	return strings.Trim(slug, "-")
}
