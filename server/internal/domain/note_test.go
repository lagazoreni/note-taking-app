package domain

import "testing"

func TestValidateNoteAndDirectiveIntegrity(t *testing.T) {
	if err := ValidateNoteText(" ", "body"); err == nil {
		t.Fatal("blank title should fail")
	}
	ids, err := ParseQuestionDirectives("One {{question:550e8400-e29b-41d4-a716-446655440000}}\n\n- item")
	if err != nil || len(ids) != 1 {
		t.Fatalf("parse directive: %v %#v", err, ids)
	}
	if ids[0] != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("unexpected id %s", ids[0])
	}
	if _, err := ParseQuestionDirectives("{{question:bad}}"); err == nil {
		t.Fatal("malformed directive should fail")
	}
	if _, err := ParseQuestionDirectives("{{question:550e8400-e29b-41d4-a716-446655440000}} {{question:550e8400-e29b-41d4-a716-446655440000}}"); err == nil {
		t.Fatal("duplicate directive should fail")
	}
}
