package importexport

import (
	"bytes"
	"testing"
)

func TestArchivePathAndLimits(t *testing.T) {
	for _, path := range []string{"../manifest.json", "/absolute", "data/../notes.json"} {
		if err := ValidateMemberPath(path); err == nil {
			t.Errorf("unsafe path accepted: %s", path)
		}
	}
	if err := ValidateMemberPath("data/notes.json"); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadArchive(bytes.NewReader([]byte("not zip")), 1024); err == nil {
		t.Fatal("malformed archive accepted")
	}
}
