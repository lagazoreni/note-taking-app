package domain

import "testing"

func TestOrganizationRules(t *testing.T) {
	if err := ValidateHierarchyMove("parent", "child", map[string]string{"child": "parent", "parent": "root"}); err == nil {
		t.Fatal("cycle was accepted")
	}
	if err := ValidateHierarchyMove("child", "root", map[string]string{"child": "parent"}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTagAccess("owner", "owner", nil); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTagAccess("owner", "other", nil); err == nil {
		t.Fatal("unshared tag accepted")
	}
}
