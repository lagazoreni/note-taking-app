package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"noted.local/noted/internal/api/jsoncodec"
	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/service"
	"noted.local/noted/internal/store/sqlite"
)

type API struct {
	Store     *sqlite.CaptureStore
	DB        *sql.DB
	Deletion  *service.NoteDeletionService
	importsMu sync.Mutex
	imports   map[string]importSession
}

func NewAPI(db *sql.DB) http.Handler {
	store := sqlite.NewCaptureStore(db, platform.SystemClock{}, nil)
	api := &API{DB: db, Store: store, Deletion: service.NewNoteDeletionService(store, platform.SystemClock{}), imports: make(map[string]importSession)}
	return api.Routes()
}

func NewAPIWithStore(store *sqlite.CaptureStore) http.Handler {
	return (&API{Store: store, Deletion: service.NewNoteDeletionService(store, platform.SystemClock{}), imports: make(map[string]importSession)}).Routes()
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/workspaces", a.listWorkspaces)
	mux.HandleFunc("POST /api/v1/workspaces", a.createWorkspace)
	mux.HandleFunc("GET /api/v1/topics", a.listTopics)
	mux.HandleFunc("POST /api/v1/topics", a.createTopic)
	mux.HandleFunc("PUT /api/v1/topics/{topicId}", a.updateTopic)
	mux.HandleFunc("DELETE /api/v1/topics/{topicId}", a.deleteTopic)
	mux.HandleFunc("GET /api/v1/tags", a.listTags)
	mux.HandleFunc("POST /api/v1/tags", a.createTag)
	mux.HandleFunc("PUT /api/v1/tags/{tagId}", a.updateTag)
	mux.HandleFunc("PUT /api/v1/tags/{tagId}/workspace-access", a.setTagAccess)
	mux.HandleFunc("GET /api/v1/workspaces/{workspaceId}", a.getWorkspace)
	mux.HandleFunc("PUT /api/v1/workspaces/{workspaceId}", a.updateWorkspace)
	mux.HandleFunc("GET /api/v1/notes", a.listNotes)
	mux.HandleFunc("POST /api/v1/notes", a.createNote)
	mux.HandleFunc("GET /api/v1/notes/{noteId}", a.getNote)
	mux.HandleFunc("PUT /api/v1/notes/{noteId}", a.updateNote)
	mux.HandleFunc("POST /api/v1/notes/{noteId}/deletion-preview", a.previewNoteDeletion)
	mux.HandleFunc("POST /api/v1/notes/{noteId}/delete", a.executeNoteDeletion)
	mux.HandleFunc("GET /api/v1/questions", a.handleQuestionQuery)
	mux.HandleFunc("GET /api/v1/questions/search", a.searchQuestions)
	mux.HandleFunc("GET /api/v1/search", a.search)
	mux.HandleFunc("POST /api/v1/reminders/evaluate", a.evaluateReminders)
	mux.HandleFunc("POST /api/v1/export", a.exportAllData)
	mux.HandleFunc("POST /api/v1/imports/validate", a.validateImport)
	mux.HandleFunc("POST /api/v1/imports/{importId}/apply", a.applyImport)
	mux.HandleFunc("DELETE /api/v1/imports/{importId}", a.cancelImport)
	mux.HandleFunc("POST /api/v1/questions", a.createQuestion)
	mux.HandleFunc("GET /api/v1/questions/{questionId}", a.getQuestion)
	mux.HandleFunc("PUT /api/v1/questions/{questionId}", a.updateQuestion)
	return mux
}

func (a *API) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	items, err := a.Store.ListWorkspaces(r.Context())
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, map[string]any{"items": items})
}

type workspaceWrite struct {
	Name    string `json:"name"`
	Version int64  `json:"version"`
}

func (a *API) createWorkspace(w http.ResponseWriter, r *http.Request) {
	var input workspaceWrite
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	value, err := a.Store.CreateWorkspace(r.Context(), input.Name)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusCreated, value)
}
func (a *API) getWorkspace(w http.ResponseWriter, r *http.Request) {
	value, err := a.Store.GetWorkspace(r.Context(), r.PathValue("workspaceId"))
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}
func (a *API) updateWorkspace(w http.ResponseWriter, r *http.Request) {
	var input workspaceWrite
	if err := jsoncodec.Decode(r, &input, 100000); err != nil {
		writeValidation(w, r, err)
		return
	}
	if input.Version < 1 {
		writeValidation(w, r, errors.New("version must be positive"))
		return
	}
	value, err := a.Store.UpdateWorkspace(r.Context(), r.PathValue("workspaceId"), input.Name, input.Version)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}

func (a *API) listNotes(w http.ResponseWriter, r *http.Request) {
	workspace := r.URL.Query().Get("workspaceId")
	if workspace == "" {
		writeValidation(w, r, errors.New("workspaceId is required"))
		return
	}
	items, err := a.listNotesBasic(r.Context(), workspace)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, map[string]any{"items": items, "nextCursor": nil})
}
func (a *API) listNotesBasic(ctx context.Context, workspace string) ([]domain.Note, error) {
	rows, err := a.DB.QueryContext(ctx, "SELECT id FROM notes WHERE workspace_id=? ORDER BY updated_at DESC,id DESC LIMIT 201", workspace)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Note, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		note, err := a.Store.GetNote(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, note)
	}
	return result, rows.Err()
}
func (a *API) createNote(w http.ResponseWriter, r *http.Request) {
	var input domain.NoteWrite
	if err := jsoncodec.Decode(r, &input, 5*1024*1024); err != nil {
		writeValidation(w, r, err)
		return
	}
	value, err := a.Store.CreateNote(r.Context(), input)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusCreated, value)
}
func (a *API) getNote(w http.ResponseWriter, r *http.Request) {
	value, err := a.Store.GetNote(r.Context(), r.PathValue("noteId"))
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}
func (a *API) updateNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		domain.NoteWrite
		Version int64 `json:"version"`
	}
	if err := jsoncodec.Decode(r, &input, 5*1024*1024); err != nil {
		writeValidation(w, r, err)
		return
	}
	if input.Version < 1 {
		writeValidation(w, r, errors.New("version must be positive"))
		return
	}
	value, err := a.Store.UpdateNote(r.Context(), r.PathValue("noteId"), input.NoteWrite, input.Version)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}

func (a *API) listQuestions(w http.ResponseWriter, r *http.Request) {
	query, err := parseQuestionQuery(r)
	if err != nil {
		writeValidation(w, r, err)
		return
	}
	page, err := a.Store.QueryQuestions(r.Context(), query)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, page)
}
func (a *API) listQuestionsBasic(ctx context.Context, workspace, status string) ([]domain.Question, error) {
	query := "SELECT id FROM questions WHERE workspace_id=?"
	args := []any{workspace}
	if status != "" {
		query += " AND status=?"
		args = append(args, status)
	}
	query += " ORDER BY updated_at DESC,id DESC LIMIT 201"
	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Question, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		value, err := a.Store.GetQuestion(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func (a *API) createQuestion(w http.ResponseWriter, r *http.Request) {
	var input domain.QuestionWrite
	if err := jsoncodec.Decode(r, &input, 5*1024*1024); err != nil {
		writeValidation(w, r, err)
		return
	}
	value, err := a.Store.CreateQuestion(r.Context(), input)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	if input.Reminder != nil {
		if _, err := a.Store.SetReminder(r.Context(), value.ID, input.Reminder.ScheduledAt); err != nil {
			writeStoreError(w, r, err)
			return
		}
		value, err = a.Store.GetQuestion(r.Context(), value.ID)
		if err != nil {
			writeStoreError(w, r, err)
			return
		}
	}
	_ = jsoncodec.Encode(w, http.StatusCreated, value)
}
func (a *API) getQuestion(w http.ResponseWriter, r *http.Request) {
	value, err := a.Store.GetQuestion(r.Context(), r.PathValue("questionId"))
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}
func (a *API) updateQuestion(w http.ResponseWriter, r *http.Request) {
	var input struct {
		domain.QuestionWrite
		Version int64 `json:"version"`
	}
	if err := jsoncodec.Decode(r, &input, 5*1024*1024); err != nil {
		writeValidation(w, r, err)
		return
	}
	if input.Version < 1 {
		writeValidation(w, r, errors.New("version must be positive"))
		return
	}
	value, err := a.Store.UpdateQuestion(r.Context(), r.PathValue("questionId"), input.QuestionWrite, input.Version)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	if input.Reminder != nil {
		if _, err := a.Store.SetReminder(r.Context(), value.ID, input.Reminder.ScheduledAt); err != nil {
			writeStoreError(w, r, err)
			return
		}
	}
	value, err = a.Store.GetQuestion(r.Context(), value.ID)
	if err != nil {
		writeStoreError(w, r, err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, value)
}

func writeValidation(w http.ResponseWriter, r *http.Request, err error) {
	message := strings.TrimSpace(err.Error())
	if strings.Contains(strings.ToLower(message), "request body too large") || strings.Contains(strings.ToLower(message), "body too large") {
		problem.Write(w, middleware.ID(r), problem.New(http.StatusRequestEntityTooLarge, problem.PayloadTooLarge, "The request body is too large"))
		return
	}
	problem.Write(w, middleware.ID(r), problem.Validation(message))
}
func writeStoreError(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr error
	switch {
	case errors.Is(err, sqlite.ErrNotFound):
		apiErr = problem.NotFoundError("The requested record was not found")
	case errors.Is(err, sqlite.ErrNameConflict):
		apiErr = problem.Conflict(problem.NameConflict, "That name is already in use", nil)
	case errors.Is(err, sqlite.ErrConflict):
		apiErr = problem.Conflict(problem.VersionConflict, "The record changed; reload before saving", nil)
	case errors.Is(err, sqlite.ErrWorkspaceBoundary), errors.Is(err, sqlite.ErrInvalidLink):
		apiErr = problem.Validation("The records must belong to the same workspace")
	default:
		message := strings.ToLower(err.Error())
		if strings.Contains(message, "resume date") {
			apiErr = problem.Validation(err.Error(), problem.FieldError{Field: "dueDate", Message: err.Error()})
		} else if strings.Contains(message, "required") || strings.Contains(message, "invalid") || strings.Contains(message, "directive") || strings.Contains(message, "must be") || strings.Contains(message, "confirmation") || strings.Contains(message, "already linked") {
			apiErr = problem.Validation(err.Error())
		} else {
			apiErr = err
		}
	}
	problem.Write(w, middleware.ID(r), apiErr)
}
func parseVersion(value string) int64 { v, _ := strconv.ParseInt(value, 10, 64); return v }
