package domain

import "testing"

func TestQuestionValidationDefaults(t *testing.T) {
	q := Question{QuestionText: "  What?  "}
	if err := q.Validate(); err != nil {
		t.Fatal(err)
	}
	if q.QuestionText != "What?" || q.Status != StatusUnanswered || q.Priority != PriorityNone {
		t.Fatalf("defaults not applied: %+v", q)
	}
	for _, mode := range []DisplayMode{DisplayExpanded, DisplayCollapsed, DisplayLink} {
		if err := ValidateDisplayMode(mode); err != nil {
			t.Fatal(err)
		}
	}
	if err := ValidateDisplayMode("bad"); err == nil {
		t.Fatal("bad display mode accepted")
	}
}

func TestQuestionKindDefaultsAndValidation(t *testing.T) {
	q := Question{QuestionText: "Why?"}
	if err := q.Validate(); err != nil {
		t.Fatal(err)
	}
	if q.Kind != KindQuestion {
		t.Fatalf("default kind=%q", q.Kind)
	}
	q.Kind = KindAnnotation
	if err := q.Validate(); err != nil {
		t.Fatal(err)
	}
	q.Kind = "other"
	if err := q.Validate(); err == nil {
		t.Fatal("invalid kind accepted")
	}
}
