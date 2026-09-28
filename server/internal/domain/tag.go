package domain

import "time"

type Tag struct {
	ID                    string    `json:"id"`
	OwnerWorkspaceID      string    `json:"ownerWorkspaceId"`
	Name                  string    `json:"name"`
	AvailableWorkspaceIDs []string  `json:"availableWorkspaceIds"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	Version               int64     `json:"version"`
}

type TagWorkspaceAccess struct {
	TagID       string    `json:"tagId"`
	WorkspaceID string    `json:"workspaceId"`
	CreatedAt   time.Time `json:"createdAt"`
}
type NoteTag struct {
	NoteID    string    `json:"noteId"`
	TagID     string    `json:"tagId"`
	CreatedAt time.Time `json:"createdAt"`
}
type QuestionTag struct {
	QuestionID string    `json:"questionId"`
	TagID      string    `json:"tagId"`
	CreatedAt  time.Time `json:"createdAt"`
}
