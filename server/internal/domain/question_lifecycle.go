package domain

import (
	"fmt"
	"strings"
)

func ValidateTransition(from, to QuestionStatus, answer, dueDate *string) error {
	if !ValidStatus(to) {
		return fmt.Errorf("invalid question status")
	}
	var trimmedAnswer string
	if answer != nil {
		trimmedAnswer = strings.TrimSpace(*answer)
	}
	if to == StatusAnswered && trimmedAnswer == "" {
		return fmt.Errorf("an answer is required before a question can be answered")
	}
	if to == StatusDeferred && (dueDate == nil || strings.TrimSpace(*dueDate) == "") {
		return fmt.Errorf("a resume date is required to defer a question")
	}
	return nil
}

func NormalizeAnswer(answer *string) *string {
	if answer == nil {
		return nil
	}
	value := strings.TrimSpace(*answer)
	if value == "" {
		return nil
	}
	return &value
}

func ReopenStatus(status QuestionStatus) QuestionStatus {
	if status == StatusAnswered {
		return StatusInProgress
	}
	return status
}
