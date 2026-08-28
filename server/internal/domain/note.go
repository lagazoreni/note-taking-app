package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type NoteQuestionWrite struct {
	QuestionID  string      `json:"questionId"`
	DisplayMode DisplayMode `json:"displayMode"`
	Position    int         `json:"position"`
}

type NoteQuestion struct {
	NoteQuestionWrite
	Question QuestionSummary `json:"question"`
}

type Note struct {
	ID            string         `json:"id"`
	WorkspaceID   string         `json:"workspaceId"`
	TopicID       *string        `json:"topicId"`
	ParentNoteID  *string        `json:"parentNoteId"`
	Title         string         `json:"title"`
	BodyMarkdown  string         `json:"bodyMarkdown"`
	QuestionLinks []NoteQuestion `json:"questionLinks"`
	TagIDs        []string       `json:"tagIds"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	Version       int64          `json:"version"`
}

type NoteWrite struct {
	WorkspaceID   string              `json:"workspaceId"`
	TopicID       *string             `json:"topicId"`
	ParentNoteID  *string             `json:"parentNoteId"`
	Title         string              `json:"title"`
	BodyMarkdown  string              `json:"bodyMarkdown"`
	QuestionLinks []NoteQuestionWrite `json:"questionLinks"`
	TagIDs        []string            `json:"tagIds"`
}

var directivePattern = regexp.MustCompile(`\{\{question:([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})\}\}`)

func ValidateNoteText(title, body string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("note title is required")
	}
	if len([]rune(strings.TrimSpace(title))) > 300 {
		return fmt.Errorf("note title must be at most 300 characters")
	}
	if len([]rune(body)) > 5000000 {
		return fmt.Errorf("note body is too large")
	}
	return nil
}

func ParseQuestionDirectives(body string) ([]string, error) {
	if strings.Contains(body, "{{question:") {
		// A directive-shaped token must be complete; do not silently retain malformed IDs as prose.
		for offset := 0; offset < len(body); {
			index := strings.Index(body[offset:], "{{question:")
			if index < 0 {
				break
			}
			index += offset
			end := strings.Index(body[index:], "}}")
			if end < 0 {
				return nil, fmt.Errorf("malformed question directive")
			}
			end += index + 2
			token := body[index:end]
			if !directivePattern.MatchString(token) {
				return nil, fmt.Errorf("malformed question directive")
			}
			offset = end
		}
	}
	matches := directivePattern.FindAllStringSubmatch(body, -1)
	ids := make([]string, 0, len(matches))
	seen := map[string]bool{}
	for _, match := range matches {
		id := strings.ToLower(match[1])
		if seen[id] {
			return nil, fmt.Errorf("question %s is linked more than once", id)
		}
		seen[id] = true
		ids = append(ids, id)
	}
	return ids, nil
}

func (w NoteWrite) Validate() error {
	if err := ValidateNoteText(w.Title, w.BodyMarkdown); err != nil {
		return err
	}
	for index, link := range w.QuestionLinks {
		if link.Position != index {
			return fmt.Errorf("question link positions must start at zero and be contiguous")
		}
		if err := ValidateDisplayMode(link.DisplayMode); err != nil {
			return err
		}
	}
	return nil
}
