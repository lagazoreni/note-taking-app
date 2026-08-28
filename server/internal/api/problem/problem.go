package problem

import (
	"encoding/json"
	"errors"
	"net/http"

	"noted.local/noted/internal/platform"
)

type Code string

const (
	ValidationFailed       Code = "VALIDATION_FAILED"
	NotFound               Code = "NOT_FOUND"
	VersionConflict        Code = "VERSION_CONFLICT"
	NameConflict           Code = "NAME_CONFLICT"
	InvalidStateTransition Code = "INVALID_STATE_TRANSITION"
	DeletionPreviewStale   Code = "DELETION_PREVIEW_STALE"
	ImportInvalid          Code = "IMPORT_INVALID"
	ImportExpired          Code = "IMPORT_EXPIRED"
	PayloadTooLarge        Code = "PAYLOAD_TOO_LARGE"
	ServiceUnavailable     Code = "SERVICE_UNAVAILABLE"
	InternalError          Code = "INTERNAL_ERROR"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	Code        Code         `json:"code"`
	Message     string       `json:"message"`
	RequestID   string       `json:"requestId"`
	FieldErrors []FieldError `json:"fieldErrors,omitempty"`
	Details     any          `json:"details,omitempty"`
}

type Envelope struct {
	Error Error `json:"error"`
}

type APIError struct {
	Status      int
	Code        Code
	Message     string
	FieldErrors []FieldError
	Details     any
	Cause       error
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}
func (e *APIError) Unwrap() error { return e.Cause }
func New(status int, code Code, message string) *APIError {
	return &APIError{Status: status, Code: code, Message: message}
}
func Validation(message string, fields ...FieldError) *APIError {
	return &APIError{Status: http.StatusUnprocessableEntity, Code: ValidationFailed, Message: message, FieldErrors: fields}
}
func NotFoundError(message string) *APIError { return New(http.StatusNotFound, NotFound, message) }
func Conflict(code Code, message string, details any) *APIError {
	return &APIError{Status: http.StatusConflict, Code: code, Message: message, Details: details}
}

func StatusFor(code Code) int {
	switch code {
	case ValidationFailed, InvalidStateTransition:
		return http.StatusUnprocessableEntity
	case NotFound:
		return http.StatusNotFound
	case VersionConflict, NameConflict, DeletionPreviewStale, ImportExpired:
		return http.StatusConflict
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func FromError(err error) *APIError {
	if err == nil {
		return nil
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.Status == 0 {
			apiErr.Status = StatusFor(apiErr.Code)
		}
		return apiErr
	}
	return &APIError{Status: http.StatusInternalServerError, Code: InternalError, Message: "The server could not complete the request", Cause: err}
}

func Write(w http.ResponseWriter, requestID string, err error) {
	if requestID == "" {
		requestID, _ = platform.NewID()
	}
	apiErr := FromError(err)
	if apiErr == nil {
		return
	}
	status := apiErr.Status
	if status == 0 {
		status = StatusFor(apiErr.Code)
	}
	if apiErr.Code == "" {
		apiErr.Code = InternalError
	}
	if apiErr.Message == "" {
		apiErr.Message = "The server could not complete the request"
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Error: Error{
		Code: apiErr.Code, Message: apiErr.Message, RequestID: requestID,
		FieldErrors: apiErr.FieldErrors, Details: apiErr.Details,
	}})
}
