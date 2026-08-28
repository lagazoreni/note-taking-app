package domain

import "testing"

func TestQuestionLifecycleTransitions(t *testing.T) {
	cases := []struct {
		from, to QuestionStatus
		answer   *string
		ok       bool
	}{
		{StatusUnanswered, StatusAnswered, nil, false},
		{StatusUnanswered, StatusAnswered, stringPtr("found"), true},
		{StatusAnswered, StatusInProgress, nil, true},
		{StatusAnswered, StatusUnanswered, nil, true},
		{StatusAnswered, StatusAnswered, stringPtr(" "), false},
		{StatusInProgress, StatusDeferred, nil, true},
	}
	for _, tc := range cases {
		err := ValidateTransition(tc.from, tc.to, tc.answer)
		if (err == nil) != tc.ok {
			t.Errorf("%s -> %s answer=%v: err=%v", tc.from, tc.to, tc.answer, err)
		}
	}
}
func stringPtr(value string) *string { return &value }
