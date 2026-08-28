package domain

import (
	"fmt"
	"strings"
)

type QuestionQuery struct {
	WorkspaceID    string
	Kind           QuestionKind
	Statuses       []QuestionStatus
	TopicID        *string
	TagID          *string
	CreatedFrom    *string
	CreatedTo      *string
	UpdatedFrom    *string
	UpdatedTo      *string
	DueFrom        *string
	DueTo          *string
	HasAnswer      *bool
	HasLinkedNotes *bool
	Sort           string
	Direction      string
	Cursor         string
	PageSize       int
}

type QuestionPage struct {
	Items      []Question `json:"items"`
	NextCursor *string    `json:"nextCursor"`
}

func ValidatePriority(value Priority) error {
	if !ValidPriority(value) {
		return fmt.Errorf("invalid priority")
	}
	return nil
}
func ValidateQuestionSort(sort, direction string) error {
	switch sort {
	case "", "createdAt", "updatedAt", "dueDate", "priority", "questionText", "status":
	default:
		return fmt.Errorf("invalid sort")
	}
	if direction != "" && direction != "asc" && direction != "desc" {
		return fmt.Errorf("invalid sort direction")
	}
	return nil
}
func NormalizeQuestionQuery(query QuestionQuery) (QuestionQuery, error) {
	if query.Kind != "" && !ValidKind(query.Kind) {
		return query, fmt.Errorf("invalid question kind")
	}
	if query.Sort == "" {
		query.Sort = "updatedAt"
	}
	if query.Direction == "" {
		query.Direction = "desc"
	}
	if query.PageSize <= 0 {
		query.PageSize = 50
	}
	if query.PageSize > 200 {
		return query, fmt.Errorf("page size must be at most 200")
	}
	if err := ValidateQuestionSort(query.Sort, query.Direction); err != nil {
		return query, err
	}
	for i, status := range query.Statuses {
		query.Statuses[i] = QuestionStatus(strings.TrimSpace(string(status)))
		if !ValidStatus(query.Statuses[i]) {
			return query, fmt.Errorf("invalid question status")
		}
	}
	return query, nil
}
func IsOverdue(dueDate string, status QuestionStatus, today string) bool {
	return dueDate != "" && dueDate < today && status != StatusAnswered
}
