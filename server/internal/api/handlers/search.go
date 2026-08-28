package handlers

import (
	"net/http"

	"noted.local/noted/internal/api/jsoncodec"
	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
	"noted.local/noted/internal/domain"
)

func (a *API) search(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	query := values.Get("q")
	content := domain.SearchContentScope(values.Get("contentScope"))
	scope := domain.SearchWorkspaceScope(values.Get("workspaceScope"))
	var workspace *string
	if value := values.Get("workspaceId"); value != "" {
		workspace = &value
	}
	if query == "" || !validContentScope(content) || !validWorkspaceScope(scope) || (scope == domain.SearchCurrent && workspace == nil) {
		problem.Write(w, middleware.ID(r), problem.Validation("q, contentScope, workspaceScope, and current workspace are required"))
		return
	}
	items, err := a.Store.Search(r.Context(), query, content, scope, workspace, 50)
	if err != nil {
		problem.Write(w, middleware.ID(r), problem.Validation(err.Error()))
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, map[string]any{"items": items, "nextCursor": nil})
}
func validContentScope(value domain.SearchContentScope) bool {
	return value == domain.SearchNotes || value == domain.SearchQuestions || value == domain.SearchAnswers || value == domain.SearchEverything
}
func validWorkspaceScope(value domain.SearchWorkspaceScope) bool {
	return value == domain.SearchCurrent || value == domain.SearchAll
}
func jsonEncode(w http.ResponseWriter, status int, value any) error {
	return jsoncodec.Encode(w, status, value)
}
