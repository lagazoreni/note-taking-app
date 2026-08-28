package domain

type SearchContentScope string
type SearchWorkspaceScope string

// SearchQuery describes the optional filters accepted by the search endpoint.
// TagID applies to note results and is intentionally optional so existing
// content and workspace searches retain their original semantics.
type SearchQuery struct {
	Query          string
	ContentScope   SearchContentScope
	WorkspaceScope SearchWorkspaceScope
	WorkspaceID    *string
	TagID          *string
	Limit          int
}

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
