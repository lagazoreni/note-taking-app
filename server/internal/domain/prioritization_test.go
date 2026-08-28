package domain

import "testing"

func TestPrioritizationValidationAndOverdue(t *testing.T) {
	for _, priority := range []Priority{PriorityNone, PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent} {
		if !ValidPriority(priority) {
			t.Errorf("priority %s rejected", priority)
		}
	}
	if ValidatePriority("critical") == nil {
		t.Fatal("invalid priority accepted")
	}
	if err := ValidateQuestionSort("updatedAt", "desc"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateQuestionSort("random", "desc"); err == nil {
		t.Fatal("invalid sort accepted")
	}
	if !IsOverdue("2026-08-08", StatusUnanswered, "2026-08-09") {
		t.Fatal("past active due date not overdue")
	}
	if IsOverdue("2026-08-08", StatusAnswered, "2026-08-09") {
		t.Fatal("answered question marked overdue")
	}
}
