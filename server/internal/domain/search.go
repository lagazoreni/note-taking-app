package domain

type SearchContentScope string
type SearchWorkspaceScope string

const (
	SearchNotes      SearchContentScope   = "notes"
	SearchQuestions  SearchContentScope   = "questions"
	SearchAnswers    SearchContentScope   = "answers"
	SearchEverything SearchContentScope   = "everything"
	SearchCurrent    SearchWorkspaceScope = "current"
	SearchAll        SearchWorkspaceScope = "all"
)

type SearchResult struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	WorkspaceID   string `json:"workspaceId"`
	WorkspaceName string `json:"workspaceName"`
	Title         string `json:"title,omitempty"`
	Snippet       string `json:"snippet"`
	Destination   string `json:"destination"`
}
