package named_errors

import (
	"testing"
)

func TestNamedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrNotFound", ErrNotFound, "entity not found"},
		{"ErrAccessDenied", ErrAccessDenied, "access denied"},
		{"ErrConflict", ErrConflict, "resource conflict or duplicate"},
		{"ErrInvalidInput", ErrInvalidInput, "invalid input data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message %q, got %q", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestNamedErrors_Uniqueness(t *testing.T) {
	errorList := []error{ErrNotFound, ErrAccessDenied, ErrConflict, ErrInvalidInput}

	for i, err1 := range errorList {
		for j, err2 := range errorList {
			if i != j && err1.Error() == err2.Error() {
				t.Errorf("Error %d and %d should not be equal", i, j)
			}
		}
	}
}

func TestNamedErrors_ErrorMessages(t *testing.T) {
	if ErrNotFound.Error() == "" {
		t.Error("ErrNotFound should have an error message")
	}

	if ErrAccessDenied.Error() == "" {
		t.Error("ErrAccessDenied should have an error message")
	}

	if ErrConflict.Error() == "" {
		t.Error("ErrConflict should have an error message")
	}

	if ErrInvalidInput.Error() == "" {
		t.Error("ErrInvalidInput should have an error message")
	}
}

func TestNamedErrors_ErrorTypes(t *testing.T) {
	if ErrNotFound == nil {
		t.Error("ErrNotFound should not be nil")
	}

	if ErrAccessDenied == nil {
		t.Error("ErrAccessDenied should not be nil")
	}

	if ErrConflict == nil {
		t.Error("ErrConflict should not be nil")
	}

	if ErrInvalidInput == nil {
		t.Error("ErrInvalidInput should not be nil")
	}
}

func TestNamedErrors_ErrorConsistency(t *testing.T) {
	notFound1 := ErrNotFound
	notFound2 := ErrNotFound

	if notFound1.Error() != notFound2.Error() {
		t.Error("Same error should have consistent messages")
	}

	accessDenied1 := ErrAccessDenied
	accessDenied2 := ErrAccessDenied

	if accessDenied1.Error() != accessDenied2.Error() {
		t.Error("Same error should have consistent messages")
	}
}

func TestNamedErrors_ErrorLength(t *testing.T) {
	if len(ErrNotFound.Error()) < 5 {
		t.Error("ErrNotFound message should be at least 5 characters")
	}

	if len(ErrAccessDenied.Error()) < 5 {
		t.Error("ErrAccessDenied message should be at least 5 characters")
	}

	if len(ErrConflict.Error()) < 5 {
		t.Error("ErrConflict message should be at least 5 characters")
	}

	if len(ErrInvalidInput.Error()) < 5 {
		t.Error("ErrInvalidInput message should be at least 5 characters")
	}
}