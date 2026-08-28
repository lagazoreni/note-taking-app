package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestOrganizationUniquenessAndTagSharing(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), nil)
	one, _ := store.CreateWorkspace(context.Background(), "Work")
	two, _ := store.CreateWorkspace(context.Background(), "Studies")
	if _, err := store.CreateTopic(context.Background(), one.ID, "Research"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateTopic(context.Background(), one.ID, "research"); err == nil {
		t.Fatal("case-insensitive topic duplicate accepted")
	}
	tag, err := store.CreateTag(context.Background(), one.ID, "shared")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetTagAccess(context.Background(), tag.ID, []string{two.ID}, tag.Version, map[string]bool{}); err != nil {
		t.Fatal(err)
	}
	available, err := store.ListTags(context.Background(), two.ID, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(available) != 1 {
		t.Fatalf("available tags=%d", len(available))
	}
	_ = domain.PriorityNone
}
