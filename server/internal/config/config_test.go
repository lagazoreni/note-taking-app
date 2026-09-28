package config

import (
	"strings"
	"testing"
	"time"
)

func TestFromEnvDefaultsAndOverrides(t *testing.T) {
	cfg, err := FromEnv(map[string]string{})
	if err != nil {
		t.Fatalf("defaults: %v", err)
	}
	if cfg.Addr != "127.0.0.1:8080" || cfg.DBFile != "notes.db" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	custom, err := FromEnv(map[string]string{
		"NOTED_ADDR": "0.0.0.0:9000", "NOTED_DATA_DIR": "/tmp/noted",
		"NOTED_DB_FILE": "custom.sqlite", "NOTED_LOG_LEVEL": "debug",
		"NOTED_IMPORT_MAX_BYTES": "1024", "NOTED_IMPORT_TTL": "2h",
	})
	if err != nil {
		t.Fatalf("custom: %v", err)
	}
	if custom.Addr != "0.0.0.0:9000" || custom.ImportMaxBytes != 1024 || custom.ImportTTL != 2*time.Hour {
		t.Fatalf("unexpected custom config: %+v", custom)
	}
}

func TestFromEnvRejectsUnsafeValues(t *testing.T) {
	cases := []map[string]string{
		{"NOTED_ADDR": ""},
		{"NOTED_DATA_DIR": ""},
		{"NOTED_IMPORT_MAX_BYTES": "0"},
		{"NOTED_IMPORT_TTL": "not-a-duration"},
		{"NOTED_LOG_LEVEL": "verbose"},
		{"NOTED_DB_FILE": "../escape.db"},
	}
	for _, env := range cases {
		if _, err := FromEnv(env); err == nil {
			t.Errorf("expected error for %v", env)
		} else if strings.TrimSpace(err.Error()) == "" {
			t.Errorf("empty error for %v", env)
		}
	}
}
