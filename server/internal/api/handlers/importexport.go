package handlers

import (
	"errors"
	"io"
	"net/http"
	"time"

	"noted.local/noted/internal/api/jsoncodec"
	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
	"noted.local/noted/internal/importexport"
	"noted.local/noted/internal/platform"
)

type importSession struct {
	Archive importexport.Archive
	Preview importexport.Preview
}

func (a *API) exportAllData(w http.ResponseWriter, r *http.Request) {
	data, err := importexport.Export(r.Context(), a.DB, "0.1.0", platform.SystemClock{})
	if err != nil {
		problem.Write(w, middleware.ID(r), err)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=noted-export.zip")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
func (a *API) validateImport(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(256 << 20); err != nil {
		problem.Write(w, middleware.ID(r), problem.New(http.StatusRequestEntityTooLarge, problem.PayloadTooLarge, "The import archive is too large"))
		return
	}
	file, _, err := r.FormFile("archive")
	if err != nil {
		writeValidation(w, r, errors.New("archive is required"))
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 256<<20+1))
	if err != nil {
		problem.Write(w, middleware.ID(r), problem.New(http.StatusUnprocessableEntity, problem.ImportInvalid, "The import archive could not be read"))
		return
	}
	if len(data) > 256<<20 {
		problem.Write(w, middleware.ID(r), problem.New(http.StatusRequestEntityTooLarge, problem.PayloadTooLarge, "The import archive is too large"))
		return
	}
	archive, err := importexport.ValidateArchive(data, 256<<20)
	if err != nil {
		problem.Write(w, middleware.ID(r), problem.New(http.StatusUnprocessableEntity, problem.ImportInvalid, "The import archive is invalid"))
		return
	}
	conflicts, err := importexport.FindConflicts(r.Context(), a.DB, archive)
	if err != nil {
		problem.Write(w, middleware.ID(r), problem.New(http.StatusUnprocessableEntity, problem.ImportInvalid, "The import could not be inspected"))
		return
	}
	id, _ := platform.NewID()
	preview := importexport.Preview{ImportID: id, ExpiresAt: time.Now().UTC().Add(24 * time.Hour), FormatVersion: archive.Manifest.FormatVersion, ArchiveSHA256: archive.SHA256, Counts: archive.Manifest.Counts, Conflicts: conflicts}
	a.importsMu.Lock()
	a.imports[id] = importSession{Archive: archive, Preview: preview}
	a.importsMu.Unlock()
	_ = jsoncodec.Encode(w, http.StatusOK, preview)
}
func (a *API) applyImport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("importId")
	a.importsMu.Lock()
	session, ok := a.imports[id]
	a.importsMu.Unlock()
	if !ok {
		problem.Write(w, middleware.ID(r), problem.NotFoundError("Import session was not found"))
		return
	}
	if time.Now().After(session.Preview.ExpiresAt) {
		a.importsMu.Lock()
		delete(a.imports, id)
		a.importsMu.Unlock()
		problem.Write(w, middleware.ID(r), problem.Conflict(problem.ImportExpired, "Import preview has expired", nil))
		return
	}
	var input struct {
		ConflictResolutions []importexport.Resolution `json:"conflictResolutions"`
	}
	if err := jsoncodec.Decode(r, &input, 1000000); err != nil {
		writeValidation(w, r, err)
		return
	}
	resolutions := map[string]string{}
	for _, conflict := range session.Preview.Conflicts {
		found := ""
		for _, resolution := range input.ConflictResolutions {
			if resolution.ConflictID == conflict.ConflictID {
				found = resolution.Resolution
				break
			}
		}
		if found != "keep_existing" && found != "use_imported" {
			problem.Write(w, middleware.ID(r), problem.Validation("Every import conflict requires a resolution"))
			return
		}
		resolutions[conflict.EntityID] = found
	}
	counts, err := importexport.Apply(r.Context(), a.DB, session.Archive, resolutions)
	if err != nil {
		problem.Write(w, middleware.ID(r), problem.New(http.StatusUnprocessableEntity, problem.ImportInvalid, "The import could not be applied; existing data was left unchanged"))
		return
	}
	a.importsMu.Lock()
	delete(a.imports, id)
	a.importsMu.Unlock()
	_ = jsoncodec.Encode(w, http.StatusOK, map[string]any{"importId": id, "state": "applied", "appliedCounts": counts})
}
func (a *API) cancelImport(w http.ResponseWriter, r *http.Request) {
	a.importsMu.Lock()
	_, ok := a.imports[r.PathValue("importId")]
	delete(a.imports, r.PathValue("importId"))
	a.importsMu.Unlock()
	if !ok {
		problem.Write(w, middleware.ID(r), problem.NotFoundError("Import session was not found"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
