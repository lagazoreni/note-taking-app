package importexport

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
)

const (
	MaxCompressedBytes int64 = 268435456
	MaxExpandedBytes   int64 = 1073741824
	MaxMembers               = 1000
)

func ValidateMemberPath(name string) error {
	if name == "" || strings.HasPrefix(name, "/") || strings.Contains(name, "\\") || path.IsAbs(name) || path.Clean(name) != name || strings.HasPrefix(name, "../") || name == ".." {
		return fmt.Errorf("unsafe archive member path")
	}
	return nil
}
func ReadArchive(source any, maxBytes int64) (Archive, error) {
	var data []byte
	switch value := source.(type) {
	case []byte:
		data = value
	case io.Reader:
		var err error
		data, err = io.ReadAll(io.LimitReader(value, MaxCompressedBytes+1))
		if err != nil {
			return Archive{}, err
		}
	default:
		return Archive{}, fmt.Errorf("archive source must be bytes or a reader")
	}
	if maxBytes <= 0 || maxBytes > MaxCompressedBytes {
		maxBytes = MaxCompressedBytes
	}
	if int64(len(data)) > maxBytes {
		return Archive{}, fmt.Errorf("archive exceeds limit")
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return Archive{}, fmt.Errorf("invalid zip archive: %w", err)
	}
	if len(reader.File) > MaxMembers {
		return Archive{}, fmt.Errorf("archive has too many members")
	}
	files := map[string][]byte{}
	var expanded int64
	for _, file := range reader.File {
		if err := ValidateMemberPath(file.Name); err != nil {
			return Archive{}, err
		}
		if _, exists := files[file.Name]; exists {
			return Archive{}, fmt.Errorf("duplicate archive member")
		}
		if file.UncompressedSize64 > uint64(MaxExpandedBytes) || expanded+int64(file.UncompressedSize64) > MaxExpandedBytes {
			return Archive{}, fmt.Errorf("archive expands beyond limit")
		}
		handle, err := file.Open()
		if err != nil {
			return Archive{}, err
		}
		value, err := io.ReadAll(io.LimitReader(handle, MaxExpandedBytes-expanded+1))
		handle.Close()
		if err != nil {
			return Archive{}, err
		}
		expanded += int64(len(value))
		files[file.Name] = value
	}
	archive := Archive{Files: files, SHA256: sha(data)}
	if err := validateManifest(archive); err != nil {
		return Archive{}, err
	}
	if err := json.Unmarshal(files["manifest.json"], &archive.Manifest); err != nil {
		return Archive{}, err
	}
	return archive, nil
}
func ValidateArchive(data []byte, maxBytes int64) (Archive, error) {
	return ReadArchive(data, maxBytes)
}
func validateManifest(archive Archive) error {
	manifestData, ok := archive.Files["manifest.json"]
	if !ok {
		return fmt.Errorf("manifest.json is required")
	}
	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	allowedMembers := map[string]bool{"manifest.json": true, "checksums.txt": true}
	for _, file := range manifest.Files {
		allowedMembers[file.Path] = true
	}
	for member := range archive.Files {
		if !allowedMembers[member] {
			return fmt.Errorf("unknown archive member %s", member)
		}
	}
	if manifest.Format != "interactive-note-questions-export" || manifest.FormatVersion != 1 {
		return fmt.Errorf("unsupported export format")
	}
	if len(manifest.Files) != 10 {
		return fmt.Errorf("manifest must declare ten data files")
	}
	for _, entry := range manifest.Files {
		data, ok := archive.Files[entry.Path]
		if !ok {
			return fmt.Errorf("missing member %s", entry.Path)
		}
		if sum(data) != entry.SHA256 {
			return fmt.Errorf("checksum mismatch for %s", entry.Path)
		}
		if err := validateJSONArray(data); err != nil {
			return fmt.Errorf("invalid records in %s: %w", entry.Path, err)
		}
	}
	checksums, ok := archive.Files["checksums.txt"]
	if !ok {
		return fmt.Errorf("checksums.txt is required")
	}
	for _, line := range strings.Split(strings.TrimSpace(string(checksums)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 || len(parts[0]) != 64 {
			return fmt.Errorf("invalid checksums.txt")
		}
		if data, exists := archive.Files[parts[1]]; !exists || sum(data) != parts[0] {
			return fmt.Errorf("checksum entry mismatch")
		}
	}
	return nil
}
func FindConflicts(ctx context.Context, db *sql.DB, archive Archive) ([]Conflict, error) {
	conflicts := []Conflict{}
	for _, entry := range archive.Manifest.Files {
		table := map[string]string{"workspace": "workspaces", "topic": "topics", "note": "notes", "question": "questions", "tag": "tags", "reminder": "reminders"}[entry.Entity]
		if table == "" {
			continue
		}
		var records []map[string]any
		if err := json.Unmarshal(archive.Files[entry.Path], &records); err != nil {
			return nil, err
		}
		for _, record := range records {
			id := fmt.Sprint(record["id"])
			if id == "<nil>" || id == "" {
				continue
			}
			var count int
			if err := db.QueryRowContext(ctx, "SELECT count(*) FROM "+table+" WHERE id=?", id).Scan(&count); err != nil {
				return nil, err
			}
			if count == 0 {
				continue
			}
			var existing sql.NullInt64
			_ = db.QueryRowContext(ctx, "SELECT version FROM "+table+" WHERE id=?", id).Scan(&existing)
			var imported *int64
			if value, ok := record["version"].(float64); ok {
				parsed := int64(value)
				imported = &parsed
			}
			reason := "different_content"
			if existing.Valid && imported != nil && existing.Int64 > *imported {
				reason = "newer_existing"
			}
			conflicts = append(conflicts, Conflict{ConflictID: conflictID(entry.Entity, id), EntityType: entry.Entity, EntityID: id, Reason: reason, ExistingVersion: nullableInt64(existing), ImportedVersion: imported, AllowedResolutions: []string{"keep_existing", "use_imported"}})
		}
	}
	return conflicts, nil
}
func conflictID(entity, id string) string {
	value := sha256.Sum256([]byte(entity + ":" + id))
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", value[0:4], value[4:6], value[6:8], value[8:10], value[10:16])
}
func nullableInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	return &value.Int64
}

func validateJSONArray(data []byte) error {
	var value []map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return nil
}
func sha(data []byte) string { value := sha256.Sum256(data); return hex.EncodeToString(value[:]) }
