package domain

import (
	"fmt"
	"strings"
	"time"
)

type QuestionStatus string
type Priority string
type DisplayMode string
type QuestionKind string

const (
	StatusUnanswered QuestionStatus = "unanswered"
	StatusInProgress QuestionStatus = "in_progress"
	StatusDeferred   QuestionStatus = "deferred"
	StatusAnswered   QuestionStatus = "answered"
	PriorityNone     Priority       = "none"
	PriorityLow      Priority       = "low"
	PriorityMedium   Priority       = "medium"
	PriorityHigh     Priority       = "high"
	PriorityUrgent   Priority       = "urgent"
	DisplayExpanded  DisplayMode    = "expanded"
	DisplayCollapsed DisplayMode    = "collapsed"
	DisplayLink      DisplayMode    = "link"
	KindQuestion     QuestionKind   = "question"
	KindAnnotation   QuestionKind   = "annotation"
)

type Reminder struct {
	ID              string     `json:"id"`
	ScheduledAt     time.Time  `json:"scheduledAt"`
	State           string     `json:"state"`
	LastEvaluatedAt *time.Time `json:"lastEvaluatedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	Version         int64      `json:"version"`
}

type LinkedNoteSummary struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	DisplayMode DisplayMode `json:"displayMode"`
}

type QuestionSummary struct {
	ID           string         `json:"id"`
	WorkspaceID  string         `json:"workspaceId"`
	QuestionText string         `json:"questionText"`
	Kind         QuestionKind   `json:"kind"`
	Status       QuestionStatus `json:"status"`
	Priority     Priority       `json:"priority"`
	DueDate      *string        `json:"dueDate"`
	Version      int64          `json:"version"`
}

type Question struct {
	ID             string              `json:"id"`
	WorkspaceID    string              `json:"workspaceId"`
	QuestionText   string              `json:"questionText"`
	Kind           QuestionKind        `json:"kind"`
	AnswerMarkdown *string             `json:"answerMarkdown"`
	Status         QuestionStatus      `json:"status"`
	Priority       Priority            `json:"priority"`
	DueDate        *string             `json:"dueDate"`
	Reminder       *Reminder           `json:"reminder"`
	TagIDs         []string            `json:"tagIds"`
	LinkedNotes    []LinkedNoteSummary `json:"linkedNotes"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
	Version        int64               `json:"version"`
}

type QuestionWrite struct {
	WorkspaceID    string         `json:"workspaceId"`
	QuestionText   string         `json:"questionText"`
	Kind           QuestionKind   `json:"kind"`
	AnswerMarkdown *string        `json:"answerMarkdown"`
	Status         QuestionStatus `json:"status"`
	Priority       Priority       `json:"priority"`
	DueDate        *string        `json:"dueDate"`
	Reminder       *ReminderWrite `json:"reminder"`
	TagIDs         []string       `json:"tagIds"`
}

type ReminderWrite struct {
	ScheduledAt time.Time `json:"scheduledAt"`
}

func (q *Question) Validate() error {
	q.QuestionText = strings.TrimSpace(q.QuestionText)
	if q.QuestionText == "" {
		return fmt.Errorf("question text is required")
	}
	if len([]rune(q.QuestionText)) > 10000 {
		return fmt.Errorf("question text must be at most 10000 characters")
	}
	if q.Kind == "" {
		q.Kind = KindQuestion
	}
	if !ValidKind(q.Kind) {
		return fmt.Errorf("invalid question kind")
	}
	if q.Status == "" {
		q.Status = StatusUnanswered
	}
	if !ValidStatus(q.Status) {
		return fmt.Errorf("invalid question status")
	}
	if q.Priority == "" {
		q.Priority = PriorityNone
	}
	if !ValidPriority(q.Priority) {
		return fmt.Errorf("invalid question priority")
	}
	if q.AnswerMarkdown != nil {
		value := strings.TrimSpace(*q.AnswerMarkdown)
		if len([]rune(value)) > 5000000 {
			return fmt.Errorf("answer must be at most 5000000 characters")
		}
		if value == "" {
			q.AnswerMarkdown = nil
		} else {
			q.AnswerMarkdown = &value
		}
	}
	if q.Status == StatusAnswered && q.AnswerMarkdown == nil {
		return fmt.Errorf("an answer is required before a question can be answered")
	}
	if q.DueDate != nil && !validDate(*q.DueDate) {
		return fmt.Errorf("due date must use YYYY-MM-DD")
	}
	return nil
}

func ValidStatus(status QuestionStatus) bool {
	return status == StatusUnanswered || status == StatusInProgress || status == StatusDeferred || status == StatusAnswered
}
func ValidKind(kind QuestionKind) bool {
	return kind == KindQuestion || kind == KindAnnotation
}
func ValidPriority(priority Priority) bool {
	return priority == PriorityNone || priority == PriorityLow || priority == PriorityMedium || priority == PriorityHigh || priority == PriorityUrgent
}
func ValidateDisplayMode(mode DisplayMode) error {
	if mode != DisplayExpanded && mode != DisplayCollapsed && mode != DisplayLink {
		return fmt.Errorf("invalid display mode")
	}
	return nil
}
func validDate(value string) bool {
	if len(value) != 10 {
		return false
	}
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}
