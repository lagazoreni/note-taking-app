package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"noted.local/noted/internal/api/jsoncodec"
	"noted.local/noted/internal/domain"
)

func (a *API) listTopics(w http.ResponseWriter, r *http.Request) {
	workspace := r.URL.Query().Get("workspaceId")
	if workspace == "" {
		writeValidation(w, r, errors.New("workspaceId is required"))
		return
	}
	items, err := a.Store.ListTopics(r.Context(), workspace)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, map[string]any{"items": items})
}
func (a *API) createTopic(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkspaceID string `json:"workspaceId"`
		Name        string `json:"name"`
	}
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	value, err := a.Store.CreateTopic(r.Context(), input.WorkspaceID, input.Name)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusCreated, value)
}
func (a *API) updateTopic(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Version int64  `json:"version"`
	}
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	value, err := a.Store.UpdateTopic(r.Context(), r.PathValue("topicId"), input.Name, input.Version)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}
func (a *API) deleteTopic(w http.ResponseWriter, r *http.Request) {
	version, _ := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	if version < 1 {
		writeValidation(w, r, errors.New("version is required"))
		return
	}
	if err := a.Store.DeleteTopic(r.Context(), r.PathValue("topicId"), version); err != nil {
		writeStoreError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (a *API) listTags(w http.ResponseWriter, r *http.Request) {
	workspace := r.URL.Query().Get("workspaceId")
	if workspace == "" {
		writeValidation(w, r, errors.New("workspaceId is required"))
		return
	}
	include := r.URL.Query().Get("includeShared") != "false"
	items, err := a.Store.ListTags(r.Context(), workspace, include)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, map[string]any{"items": items})
}
func (a *API) createTag(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OwnerWorkspaceID string `json:"ownerWorkspaceId"`
		Name             string `json:"name"`
	}
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	value, err := a.Store.CreateTag(r.Context(), input.OwnerWorkspaceID, input.Name)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusCreated, value)
}
func (a *API) updateTag(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Version int64  `json:"version"`
	}
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	value, err := a.Store.UpdateTag(r.Context(), r.PathValue("tagId"), input.Name, input.Version)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}
func (a *API) setTagAccess(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkspaceIDs     []string `json:"workspaceIds"`
		Version          int64    `json:"version"`
		RemovalDecisions []struct {
			WorkspaceID       string `json:"workspaceId"`
			RemoveAssignments bool   `json:"removeAssignments"`
		} `json:"removalDecisions"`
	}
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	decisions := map[string]bool{}
	for _, decision := range input.RemovalDecisions {
		decisions[decision.WorkspaceID] = decision.RemoveAssignments
	}
	value, err := a.Store.SetTagAccess(r.Context(), r.PathValue("tagId"), input.WorkspaceIDs, input.Version, decisions)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}

var _ domain.Topic
