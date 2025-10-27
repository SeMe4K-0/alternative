package named_errors

import "errors"

var (
	ErrNotFound = errors.New("entity not found")

	ErrAccessDenied = errors.New("access denied")

	ErrConflict = errors.New("resource conflict or duplicate")

	ErrInvalidInput = errors.New("invalid input data")
)