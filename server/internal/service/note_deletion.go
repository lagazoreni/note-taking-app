package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
)

type DeletionPreview struct {
	Token                   string                   `json:"previewToken"`
	ExpiresAt               time.Time                `json:"expiresAt"`
	NoteID                  string                   `json:"noteId"`
	NoteVersion             int64                    `json:"noteVersion"`
	SinglyLinkedQuestions   []domain.QuestionSummary `json:"singlyLinkedQuestions"`
	MultiplyLinkedQuestions []domain.QuestionSummary `json:"multiplyLinkedQuestions"`
	ChildNotes              []sqlite.ChildNote       `json:"childNotes"`
}
type deletionState struct{ Preview DeletionPreview }
type NoteDeletionService struct {
	Store    *sqlite.CaptureStore
	Clock    platform.Clock
	TTL      time.Duration
	mu       sync.Mutex
	previews map[string]deletionState
}

func NewNoteDeletionService(store *sqlite.CaptureStore, clock platform.Clock) *NoteDeletionService {
	if clock == nil {
		clock = platform.SystemClock{}
	}
	return &NoteDeletionService{Store: store, Clock: clock, TTL: 15 * time.Minute, previews: map[string]deletionState{}}
}
func (s *NoteDeletionService) Preview(ctx context.Context, noteID string, version int64) (DeletionPreview, error) {
	impact, err := s.Store.NoteDeletionImpact(ctx, noteID)
	if err != nil {
		return DeletionPreview{}, err
	}
	if impact.Note.Version != version {
		return DeletionPreview{}, sqlite.ErrConflict
	}
	token, err := platform.NewID()
	if err != nil {
		return DeletionPreview{}, err
	}
	preview := DeletionPreview{Token: token, ExpiresAt: s.Clock.Now().Add(s.TTL), NoteID: noteID, NoteVersion: version, SinglyLinkedQuestions: impact.Singly, MultiplyLinkedQuestions: impact.Multiple, ChildNotes: impact.Children}
	s.mu.Lock()
	s.previews[token] = deletionState{Preview: preview}
	s.mu.Unlock()
	return preview, nil
}
func (s *NoteDeletionService) Execute(ctx context.Context, noteID, token string, version int64, questionDecisions map[string]string, childDecisions map[string]string) error {
	s.mu.Lock()
	state, ok := s.previews[token]
	if ok {
		delete(s.previews, token)
	}
	s.mu.Unlock()
	if !ok || state.Preview.NoteID != noteID || state.Preview.NoteVersion != version || s.Clock.Now().After(state.Preview.ExpiresAt) {
		return fmt.Errorf("deletion preview is stale")
	}
	for _, q := range state.Preview.SinglyLinkedQuestions {
		if questionDecisions[q.ID] != "keep_unlinked" && questionDecisions[q.ID] != "delete" {
			return fmt.Errorf("an explicit decision is required for question %s", q.ID)
		}
	}
	for _, child := range state.Preview.ChildNotes {
		if childDecisions[child.ID] == "" {
			return fmt.Errorf("an explicit child decision is required for note %s", child.ID)
		}
	}
	return s.Store.DeleteNote(ctx, noteID, version, questionDecisions, childDecisions)
}
