package contentimport

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildExtractBundleParses(t *testing.T) {
	var reply extraction
	err := json.Unmarshal([]byte(`{
	  "lessonTitle": "Present simple",
	  "grammarPoints": [{"slug":"G3 Present Simple","title":"Present simple","titleVi":"Hiện tại đơn","summaryVi":"...","formula":"S + V(s)",
	    "examples":[{"en":"He works.","vi":"Anh ấy làm việc."}],"commonMistakes":[{"wrong":"He work.","right":"He works.","noteVi":"thêm s"}]}],
	  "activities": [
	    {"type":"GRAMMAR","title":"Exercise 1","instructionsVi":"Điền","grammarPointSlug":"g3-present-simple","readingText":"",
	     "exercises":[{"type":"FILL_BLANK","prompt":"Điền","payloadJson":"{\"sentence\":\"She ___ (work).\"}","answerJson":"{\"accepted\":[\"works\"]}","explanationVi":"she → works","skillTag":"grammar:g3-present-simple"},
	                  {"type":"TRUE_FALSE","prompt":"Đúng/sai","payloadJson":"not json","answerJson":"{\"value\":true}","explanationVi":"","skillTag":"grammar:g3"}]},
	    {"type":"READING","title":"A day","instructionsVi":"Đọc","grammarPointSlug":"","readingText":"Tom gets up at 7.","exercises":[]}
	  ]}`), &reply)
	if err != nil {
		t.Fatal(err)
	}

	result, err := buildExtractBundle(ExtractRequest{TopicTitle: "Book unit 5", TopicSlug: "book-unit-5"}, reply)
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := ParseBundle(result.Bundle)
	if err != nil {
		t.Fatalf("bundle does not parse: %v\n%s", err, result.Bundle)
	}
	if bundle.GrammarPoints[0].Slug != "g3-present-simple" {
		t.Errorf("grammar slug = %q", bundle.GrammarPoints[0].Slug)
	}
	activities := bundle.Topics[0].Lessons[0].Activities
	if len(activities) != 2 || activities[0].Status != "DRAFT" || activities[1].PassageSlug == "" {
		t.Errorf("activities = %+v", activities)
	}
	if len(bundle.Passages) != 1 || bundle.Passages[0].BodyText != "Tom gets up at 7." {
		t.Errorf("passages = %+v", bundle.Passages)
	}
	if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "#2") {
		t.Errorf("warnings = %v", result.Warnings)
	}
}

func TestExtractPartsRejectsMixedFiles(t *testing.T) {
	pdf := UploadedFile{Name: "a.pdf", Data: []byte("%PDF-1.4 test")}
	png := UploadedFile{Name: "a.png", Data: []byte("\x89PNG\r\n\x1a\n0000")}
	if _, err := extractParts(ExtractRequest{Files: []UploadedFile{pdf, png}}); err == nil {
		t.Fatal("PDF + image accepted")
	}
	if _, err := extractParts(ExtractRequest{Files: []UploadedFile{{Name: "a.txt", Data: []byte("hello")}}}); err == nil {
		t.Fatal("text file accepted")
	}
	if parts, err := extractParts(ExtractRequest{Files: []UploadedFile{png, png}}); err != nil || len(parts) != 3 {
		t.Fatalf("two images: %v %d", err, len(parts))
	}
}
