package importexport

import "time"

type FileEntry struct {
	Path    string `json:"path"`
	Entity  string `json:"entity"`
	SHA256  string `json:"sha256"`
	Records int    `json:"records"`
}
type Counts struct {
	Workspaces         int `json:"workspaces"`
	Topics             int `json:"topics"`
	Notes              int `json:"notes"`
	Questions          int `json:"questions"`
	NoteQuestions      int `json:"noteQuestions"`
	Tags               int `json:"tags"`
	TagWorkspaceAccess int `json:"tagWorkspaceAccess"`
	NoteTags           int `json:"noteTags"`
	QuestionTags       int `json:"questionTags"`
	Reminders          int `json:"reminders"`
}
type Manifest struct {
	Format             string      `json:"format"`
	FormatVersion      int         `json:"formatVersion"`
	ExportedAt         time.Time   `json:"exportedAt"`
	ApplicationVersion string      `json:"applicationVersion"`
	Files              []FileEntry `json:"files"`
	Counts             Counts      `json:"counts"`
}
type Preview struct {
	ImportID      string     `json:"importId"`
	ExpiresAt     time.Time  `json:"expiresAt"`
	FormatVersion int        `json:"formatVersion"`
	ArchiveSHA256 string     `json:"archiveSha256"`
	Counts        Counts     `json:"counts"`
	Conflicts     []Conflict `json:"conflicts"`
}
type Conflict struct {
	ConflictID         string   `json:"conflictId"`
	EntityType         string   `json:"entityType"`
	EntityID           string   `json:"entityId"`
	Reason             string   `json:"reason"`
	ExistingVersion    *int64   `json:"existingVersion"`
	ImportedVersion    *int64   `json:"importedVersion"`
	AllowedResolutions []string `json:"allowedResolutions"`
}
type Resolution struct {
	ConflictID string `json:"conflictId"`
	Resolution string `json:"resolution"`
}
type Archive struct {
	Manifest Manifest
	Files    map[string][]byte
	SHA256   string
}
