package practice

import "testing"

func TestValidItems(t *testing.T) {
	generated := []generatedItem{
		{Type: "MULTIPLE_CHOICE", Prompt: "Chọn", PayloadJSON: `{"sentence":"She ___ a teacher.","options":[{"id":"a","text":"is"},{"id":"b","text":"are"}]}`, AnswerJSON: `{"correct":"a"}`},
		// Duplicate of the first.
		{Type: "MULTIPLE_CHOICE", Prompt: "Chọn", PayloadJSON: `{"sentence":"She ___ a teacher.","options":[{"id":"a","text":"is"},{"id":"b","text":"are"}]}`, AnswerJSON: `{"correct":"a"}`},
		// Correct id missing.
		{Type: "MULTIPLE_CHOICE", Prompt: "Chọn", PayloadJSON: `{"options":[{"id":"a","text":"is"},{"id":"b","text":"are"}]}`, AnswerJSON: `{"correct":"c"}`},
		{Type: "FILL_BLANK", Prompt: "Điền", PayloadJSON: `{"sentence":"They ___ (be) from Hanoi."}`, AnswerJSON: `{"accepted":["are"]}`},
		// Too long for A1.
		{Type: "FILL_BLANK", Prompt: "Điền", PayloadJSON: `{"sentence":"My brother and my sister ___ (be) at home with our parents every Sunday morning."}`, AnswerJSON: `{"accepted":["are"]}`},
		{Type: "REORDER", Prompt: "Sắp xếp", PayloadJSON: `{"tokens":["from","I","am","Vietnam","."]}`, AnswerJSON: `{"accepted":["I am from Vietnam."]}`},
		// Tokens do not match the answer.
		{Type: "REORDER", Prompt: "Sắp xếp", PayloadJSON: `{"tokens":["from","I","is","Vietnam","."]}`, AnswerJSON: `{"accepted":["I am from Vietnam."]}`},
		{Type: "ERROR_CORRECTION", Prompt: "Sửa lỗi", PayloadJSON: `{"tokens":["He","are","my","friend","."]}`, AnswerJSON: `{"wrongIndex":1,"accepted":["is"]}`},
		// Not valid JSON.
		{Type: "FILL_BLANK", Prompt: "Điền", PayloadJSON: `{"sentence":`, AnswerJSON: `{}`},
	}
	items := validItems(generated, "grammar:g1-to-be", 12)
	if len(items) != 4 {
		t.Fatalf("want 4 valid items, got %d: %+v", len(items), items)
	}
	for _, item := range items {
		if item.SkillTag != "grammar:g1-to-be" {
			t.Errorf("skill tag %q", item.SkillTag)
		}
	}
}
