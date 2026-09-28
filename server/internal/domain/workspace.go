package domain

import (
	"fmt"
	"strings"
	"time"
)

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Version   int64     `json:"version"`
}

func ValidateName(value string, max int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("name is required")
	}
	if len([]rune(value)) > max {
		return "", fmt.Errorf("name must be at most %d characters", max)
	}
	return value, nil
}

func (w Workspace) Validate() error {
	name, err := ValidateName(w.Name, 100)
	if err != nil {
		return err
	}
	_ = name
	if w.Version < 1 {
		return fmt.Errorf("version must be positive")
	}
	return nil
}
