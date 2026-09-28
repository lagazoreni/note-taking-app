package domain

import "fmt"

func ValidateHierarchyMove(noteID, parentID string, parents map[string]string) error {
	if noteID == parentID {
		return fmt.Errorf("a note cannot parent itself")
	}
	seen := map[string]bool{noteID: true}
	current := parentID
	for current != "" && current != "root" {
		if seen[current] {
			return fmt.Errorf("note hierarchy cycle")
		}
		seen[current] = true
		next, ok := parents[current]
		if !ok {
			break
		}
		current = next
	}
	return nil
}
func ValidateTagAccess(ownerWorkspaceID, targetWorkspaceID string, available []string) error {
	if ownerWorkspaceID == targetWorkspaceID {
		return nil
	}
	for _, id := range available {
		if id == targetWorkspaceID {
			return nil
		}
	}
	return fmt.Errorf("tag is not shared with this workspace")
}
