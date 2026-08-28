package handlers

import (
	"errors"
	"net/http"

	"noted.local/noted/internal/api/jsoncodec"
	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
	"noted.local/noted/internal/store/sqlite"
)

type deleteCommand struct {
	PreviewToken      string `json:"previewToken"`
	Version           int64  `json:"version"`
	QuestionDecisions []struct {
		QuestionID string `json:"questionId"`
		Action     string `json:"action"`
	} `json:"questionDecisions"`
	ChildDecisions []struct {
		ChildNoteID string `json:"childNoteId"`
		Action      string `json:"action"`
	} `json:"childDecisions"`
}

func (a *API) previewNoteDeletion(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Version int64 `json:"version"`
	}
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	if input.Version < 1 {
		writeValidation(w, r, errors.New("version is required"))
		return
	}
	preview, err := a.Deletion.Preview(r.Context(), r.PathValue("noteId"), input.Version)
	if err != nil {
		if errors.Is(err, sqlite.ErrConflict) {
			problem.Write(w, middleware.ID(r), problem.Conflict(problem.VersionConflict, "The note changed; review deletion again", nil))
		} else {
			writeStoreError(w, r, err)
		}
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, preview)
}
func (a *API) executeNoteDeletion(w http.ResponseWriter, r *http.Request) {
	var input deleteCommand
	if err := jsoncodec.Decode(r, &input, 1000000); err != nil {
		writeValidation(w, r, err)
		return
	}
	if input.PreviewToken == "" || input.Version < 1 {
		writeValidation(w, r, errors.New("previewToken and version are required"))
		return
	}
	questions := map[string]string{}
	for _, decision := range input.QuestionDecisions {
		questions[decision.QuestionID] = decision.Action
	}
	children := map[string]string{}
	for _, decision := range input.ChildDecisions {
		children[decision.ChildNoteID] = decision.Action
	}
	if err := a.Deletion.Execute(r.Context(), r.PathValue("noteId"), input.PreviewToken, input.Version, questions, children); err != nil {
		problem.Write(w, middleware.ID(r), problem.Conflict(problem.DeletionPreviewStale, "The deletion review is stale; preview it again", nil))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
