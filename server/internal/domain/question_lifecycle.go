package domain

import (
	"fmt"
	"strings"
)

func ValidateTransition(from, to QuestionStatus, answer *string) error {
	if !ValidStatus(to) {
		return fmt.Errorf("invalid question status")
	}
	var trimmed string
	if answer != nil {
		trimmed = strings.TrimSpace(*answer)
	}
	if to == StatusAnswered && trimmed == "" {
		return fmt.Errorf("an answer is required before a question can be answered")
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
