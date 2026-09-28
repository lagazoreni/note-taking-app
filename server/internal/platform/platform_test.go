package platform

import (
	"regexp"
	"testing"
	"time"
)

func TestNewIDIsLowercaseUUID(t *testing.T) {
	id, err := NewID()
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`).MatchString(id) {
		t.Fatalf("not canonical uuid: %q", id)
	}
}

func TestFixedClockAndTimestamp(t *testing.T) {
	want := time.Date(2026, time.August, 9, 11, 22, 33, 123456000, time.FixedZone("PDT", -7*60*60))
	clock := NewFixedClock(want)
	got := clock.Now()
	if !got.Equal(want.UTC()) {
		t.Fatalf("clock = %v, want %v", got, want.UTC())
	}
	if FormatTimestamp(got) != "2026-08-09T18:22:33.123456Z" {
		t.Fatalf("unexpected timestamp: %s", FormatTimestamp(got))
	}
}
