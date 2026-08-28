package testsupport

import (
	"context"
	"fmt"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
)

func SeedPerformanceData(t *testing.T, ctx context.Context, store *sqlite.CaptureStore, workspaceID string, notes, questions int) {
	t.Helper()
	for i := 0; i < questions; i++ {
		if _, err := store.CreateQuestion(ctx, domain.QuestionWrite{WorkspaceID: workspaceID, QuestionText: fmt.Sprintf("Performance question %06d", i), Status: domain.StatusUnanswered, Priority: domain.PriorityNone}); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < notes; i++ {
		if _, err := store.CreateNote(ctx, domain.NoteWrite{WorkspaceID: workspaceID, Title: fmt.Sprintf("Performance note %06d", i), BodyMarkdown: "fixture"}); err != nil {
			t.Fatal(err)
		}
	}
}

var _ = platform.SystemClock{}
