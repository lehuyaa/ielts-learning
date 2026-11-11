package lessonactivity

import (
	"errors"
	"testing"

	"ielts-learning/backend/internal/models"
)

func TestGrade(t *testing.T) {
	tests := []struct {
		name        string
		typ         models.ExerciseType
		payload     string
		key         string
		answer      string
		wantCorrect bool
		wantRatio   float64
		wantNeedsAI bool
		wantErr     error
	}{
		// choice-based types
		{name: "multiple choice correct", typ: models.ExerciseMultipleChoice, key: `{"correct":"b"}`, answer: `{"optionId":"b"}`, wantCorrect: true, wantRatio: 1},
		{name: "multiple choice wrong", typ: models.ExerciseMultipleChoice, key: `{"correct":"b"}`, answer: `{"optionId":"a"}`},
		{name: "minimal pair", typ: models.ExerciseMinimalPair, key: `{"correct":"a"}`, answer: `{"optionId":"a"}`, wantCorrect: true, wantRatio: 1},
		{name: "ending sound uses option text", typ: models.ExerciseEndingSound, key: `{"correct":"/t/"}`, answer: `{"optionId":"/t/"}`, wantCorrect: true, wantRatio: 1},
		{name: "choice with text answer is a shape error", typ: models.ExerciseIPAToWord, key: `{"correct":"a"}`, answer: `{"text":"a"}`, wantErr: ErrInvalidAnswerShape},
		{name: "true false", typ: models.ExerciseTrueFalse, key: `{"value":false}`, answer: `{"value":false}`, wantCorrect: true, wantRatio: 1},
		{name: "true false wrong", typ: models.ExerciseTrueFalse, key: `{"value":false}`, answer: `{"value":true}`},
		{name: "word stress", typ: models.ExerciseWordStress, key: `{"stressIndex":1}`, answer: `{"stressIndex":1}`, wantCorrect: true, wantRatio: 1},

		// typed text and normalisation
		{name: "fill blank ignores case and spaces", typ: models.ExerciseFillBlank, key: `{"accepted":["works"]}`, answer: `{"text":"  Works "}`, wantCorrect: true, wantRatio: 1},
		{name: "fill blank wrong", typ: models.ExerciseFillBlank, key: `{"accepted":["works"]}`, answer: `{"text":"work"}`},
		{name: "contraction equals full form", typ: models.ExerciseTransform, key: `{"accepted":["I do not like coffee."]}`, answer: `{"text":"I don't like coffee"}`, wantCorrect: true, wantRatio: 1},
		{name: "curly apostrophe", typ: models.ExerciseTransform, key: `{"accepted":["I'm not tired."]}`, answer: `{"text":"I’m not tired!"}`, wantCorrect: true, wantRatio: 1},
		{name: "transform falls back to AI", typ: models.ExerciseTransform, key: `{"accepted":["Do you like coffee?"],"aiFallback":true}`, answer: `{"text":"You like coffee?"}`, wantNeedsAI: true},
		{name: "translate accepts reference", typ: models.ExerciseTranslateViEn, key: `{"reference":["I work in Hanoi."],"aiGraded":true}`, answer: `{"text":"i work in hanoi"}`, wantCorrect: true, wantRatio: 1},
		{name: "translate empty answer never asks AI", typ: models.ExerciseTranslateViEn, key: `{"reference":["I work in Hanoi."],"aiGraded":true}`, answer: `{"text":"  "}`},
		{name: "reorder joins tokens and punctuation", typ: models.ExerciseReorder, key: `{"accepted":["My name is Huy."]}`, answer: `{"tokens":["My","name","is","Huy","."]}`, wantCorrect: true, wantRatio: 1},
		{name: "reorder wrong order", typ: models.ExerciseReorder, key: `{"accepted":["My name is Huy."]}`, answer: `{"tokens":["Huy","is","my","name"]}`},
		{name: "reorder without tokens", typ: models.ExerciseReorder, key: `{"accepted":["My name is Huy."]}`, answer: `{"tokens":[]}`, wantErr: ErrInvalidAnswerShape},

		// error correction partial credit
		{name: "error correction full", typ: models.ExerciseErrorCorrection, key: `{"wrongIndex":1,"accepted":["works"]}`, answer: `{"index":1,"text":"works"}`, wantCorrect: true, wantRatio: 1},
		{name: "error correction index only", typ: models.ExerciseErrorCorrection, key: `{"wrongIndex":1,"accepted":["works"]}`, answer: `{"index":1,"text":"working"}`, wantRatio: 0.5},
		{name: "error correction wrong index", typ: models.ExerciseErrorCorrection, key: `{"wrongIndex":1,"accepted":["works"]}`, answer: `{"index":0,"text":"works"}`},

		// listening
		{name: "dictation exact", typ: models.ExerciseDictation, key: `{"accepted":["I usually take the bus."]}`, answer: `{"text":"i usually take the bus"}`, wantCorrect: true, wantRatio: 1},
		{name: "dictation one word missing", typ: models.ExerciseDictation, key: `{"accepted":["I usually take the bus to work."]}`, answer: `{"text":"I usually take bus to work"}`, wantRatio: 0.857},
		{name: "dictation one wrong in ten counts as correct", typ: models.ExerciseDictation, key: `{"accepted":["one two three four five six seven eight nine ten"]}`, answer: `{"text":"one two three four five six seven eight nine tan"}`, wantCorrect: true, wantRatio: 0.9},
		{name: "dictation nonsense floors at zero", typ: models.ExerciseDictation, key: `{"accepted":["Hi"]}`, answer: `{"text":"a b c d"}`},
		{name: "gap fill all", typ: models.ExerciseListenGapFill, key: `{"blanks":{"1":["nine","9"],"2":["bus"]}}`, answer: `{"blanks":{"1":"9","2":"Bus"}}`, wantCorrect: true, wantRatio: 1},
		{name: "gap fill number word", typ: models.ExerciseListenGapFill, key: `{"blanks":{"1":["9"]}}`, answer: `{"blanks":{"1":"nine"}}`, wantCorrect: true, wantRatio: 1},
		{name: "gap fill half", typ: models.ExerciseListenGapFill, key: `{"blanks":{"1":["nine"],"2":["bus"]}}`, answer: `{"blanks":{"1":"nine","2":"car"}}`, wantRatio: 0.5},
		{name: "numbers time words", typ: models.ExerciseListenNumbers, payload: `{"format":"TIME"}`, key: `{"accepted":["7:30"]}`, answer: `{"text":"seven thirty"}`, wantCorrect: true, wantRatio: 1},
		{name: "numbers time dot", typ: models.ExerciseListenNumbers, payload: `{"format":"TIME"}`, key: `{"accepted":["7:30"]}`, answer: `{"text":"07.30"}`, wantCorrect: true, wantRatio: 1},
		{name: "numbers time oclock", typ: models.ExerciseListenNumbers, payload: `{"format":"TIME"}`, key: `{"accepted":["9:00"]}`, answer: `{"text":"nine o'clock"}`, wantCorrect: true, wantRatio: 1},
		{name: "numbers price", typ: models.ExerciseListenNumbers, payload: `{"format":"PRICE"}`, key: `{"accepted":["$15"]}`, answer: `{"text":"15.00"}`, wantCorrect: true, wantRatio: 1},
		{name: "numbers price words", typ: models.ExerciseListenNumbers, payload: `{"format":"PRICE"}`, key: `{"accepted":["$25"]}`, answer: `{"text":"twenty-five dollars"}`, wantCorrect: true, wantRatio: 1},
		{name: "numbers phone", typ: models.ExerciseListenNumbers, payload: `{"format":"PHONE"}`, key: `{"accepted":["0912 345 678"]}`, answer: `{"text":"0912-345-678"}`, wantCorrect: true, wantRatio: 1},
		{name: "numbers wrong", typ: models.ExerciseListenNumbers, payload: `{"format":"NUMBER"}`, key: `{"accepted":["15"]}`, answer: `{"text":"50"}`},

		// speaking self-check and errors
		{name: "speak self check", typ: models.ExerciseSpeak, key: `{"minAccuracy":70}`, answer: `{"selfCheck":true}`, wantCorrect: true, wantRatio: 1},
		{name: "bad answer json", typ: models.ExerciseFillBlank, key: `{"accepted":["a"]}`, answer: `not json`, wantErr: ErrInvalidAnswerShape},
		{name: "missing accepted answers", typ: models.ExerciseFillBlank, key: `{}`, answer: `{"text":"a"}`, wantErr: ErrInvalidAnswerKey},
		{name: "unknown type", typ: models.ExerciseType("NOPE"), key: `{}`, answer: `{}`, wantErr: ErrInvalidAnswerKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := tt.payload
			if payload == "" {
				payload = `{}`
			}
			got, err := Grade(tt.typ, []byte(payload), []byte(tt.key), []byte(tt.answer))
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.IsCorrect != tt.wantCorrect {
				t.Errorf("IsCorrect = %v, want %v", got.IsCorrect, tt.wantCorrect)
			}
			if got.ScoreRatio != tt.wantRatio {
				t.Errorf("ScoreRatio = %v, want %v", got.ScoreRatio, tt.wantRatio)
			}
			if got.NeedsAI != tt.wantNeedsAI {
				t.Errorf("NeedsAI = %v, want %v", got.NeedsAI, tt.wantNeedsAI)
			}
		})
	}
}

func TestNormalizeText(t *testing.T) {
	cases := map[string]string{
		"  Hello,   WORLD!  ":  "hello world",
		"I'm fine.":            "i am fine",
		"It's  a   cat .":      "it is a cat",
		"We can't go":          "we cannot go",
		"We can not go":        "we cannot go",
		"Huy's laptop":         "huy's laptop",
		"Where’s the station?": "where is the station",
	}
	for input, want := range cases {
		if got := normalizeText(input); got != want {
			t.Errorf("normalizeText(%q) = %q, want %q", input, got, want)
		}
	}
}
