package handlers

import (
	"net/http"

	"noted.local/noted/internal/api/jsoncodec"
	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
)

// shared_question.go keeps canonical-link transport concerns separate from the
// basic capture handlers. The note update endpoint remains the atomic public
// contract; this small search endpoint powers the existing-question picker.
func (a *API) searchQuestions(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.URL.Query().Get("workspaceId")
	query := r.URL.Query().Get("q")
	if workspaceID == "" || query == "" {
		problem.Write(w, middleware.ID(r), problem.Validation("workspaceId and q are required"))
		return
	}
	items, err := a.Store.FindQuestions(r.Context(), workspaceID, query, 50)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, map[string]any{"items": items, "nextCursor": nil})
}
