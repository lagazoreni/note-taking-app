package handlers

import (
	"net/http"
	"time"

	"noted.local/noted/internal/api/jsoncodec"
	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
)

func (a *API) evaluateReminders(w http.ResponseWriter, r *http.Request) {
	items, err := a.Store.EvaluateReminders(r.Context(), time.Now().UTC())
	if err != nil {
		problem.Write(w, middleware.ID(r), err)
		return
	}
	_ = jsoncodec.Encode(w, http.StatusOK, items)
}
