package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
	"noted.local/noted/internal/domain"
)

func parseQuestionQuery(r *http.Request) (domain.QuestionQuery, error) {
	values := r.URL.Query()
	workspace := values.Get("workspaceId")
	if workspace == "" {
		return domain.QuestionQuery{}, fmt.Errorf("workspaceId is required")
	}
	query := domain.QuestionQuery{WorkspaceID: workspace, Kind: domain.QuestionKind(values.Get("kind")), Sort: values.Get("sort"), Direction: values.Get("direction"), Cursor: values.Get("cursor")}
	if raw := values.Get("status"); raw != "" {
		for _, status := range strings.Split(raw, ",") {
			query.Statuses = append(query.Statuses, domain.QuestionStatus(status))
		}
	}
	for key, target := range map[string]**string{"topicId": &query.TopicID, "tagId": &query.TagID, "createdFrom": &query.CreatedFrom, "createdTo": &query.CreatedTo, "updatedFrom": &query.UpdatedFrom, "updatedTo": &query.UpdatedTo, "dueFrom": &query.DueFrom, "dueTo": &query.DueTo} {
		if value := values.Get(key); value != "" {
			*target = &value
		}
	}
	for key, target := range map[string]**bool{"hasAnswer": &query.HasAnswer, "hasLinkedNotes": &query.HasLinkedNotes} {
		if value := values.Get(key); value != "" {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return domain.QuestionQuery{}, fmt.Errorf("%s must be a boolean", key)
			}
			*target = &parsed
		}
	}
	if value := values.Get("pageSize"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return domain.QuestionQuery{}, fmt.Errorf("pageSize must be an integer")
		}
		query.PageSize = parsed
	}
	return domain.NormalizeQuestionQuery(query)
}

func (a *API) handleQuestionQuery(w http.ResponseWriter, r *http.Request) {
	query, err := parseQuestionQuery(r)
	if err != nil {
		problem.Write(w, middleware.ID(r), problem.Validation(err.Error()))
		return
	}
	page, err := a.Store.QueryQuestions(r.Context(), query)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsonEncode(w, http.StatusOK, page)
}
