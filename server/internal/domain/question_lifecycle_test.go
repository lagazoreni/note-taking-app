package domain

import "testing"

func TestQuestionLifecycleTransitions(t *testing.T) {
	cases := []struct {
		name            string
		from, to        QuestionStatus
		answer, dueDate *string
		ok              bool
	}{
		{
			name:    "answered requires an answer",
			from:    StatusUnanswered,
			to:      StatusAnswered,
			answer:  nil,
			dueDate: nil,
			ok:      false,
		},
		{
			name:    "answered does not require a due date",
			from:    StatusUnanswered,
			to:      StatusAnswered,
			answer:  stringPtr("found"),
			dueDate: nil,
			ok:      true,
		},
		{
			name:    "answered rejects a blank answer",
			from:    StatusAnswered,
			to:      StatusAnswered,
			answer:  stringPtr(" "),
			dueDate: nil,
			ok:      false,
		},
		{
			name:    "in progress allows a null due date",
			from:    StatusAnswered,
			to:      StatusInProgress,
			answer:  nil,
			dueDate: nil,
			ok:      true,
		},
		{
			name:    "unanswered allows a null due date",
			from:    StatusInProgress,
			to:      StatusUnanswered,
			answer:  nil,
			dueDate: nil,
			ok:      true,
		},
		{
			name:    "deferred requires a due date",
			from:    StatusInProgress,
			to:      StatusDeferred,
			answer:  nil,
			dueDate: nil,
			ok:      false,
		},
		{
			name:    "deferred rejects an empty due date",
			from:    StatusInProgress,
			to:      StatusDeferred,
			answer:  nil,
			dueDate: stringPtr(""),
			ok:      false,
		},
		{
			name:    "deferred accepts a date",
			from:    StatusInProgress,
			to:      StatusDeferred,
			answer:  nil,
			dueDate: stringPtr("2026-04-15"),
			ok:      true,
		},
	}
	for _, tc := range cases {
		err := ValidateTransition(tc.from, tc.to, tc.answer, tc.dueDate)
		if (err == nil) != tc.ok {
			t.Errorf("%s (%s -> %s answer=%v dueDate=%v): err=%v", tc.name, tc.from, tc.to, tc.answer, tc.dueDate, err)
		}
	}
}
func stringPtr(value string) *string { return &value }
