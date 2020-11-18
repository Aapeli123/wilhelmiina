package schedule

import "errors"

var (
	// ErrDocExistsAlready is thrown if document is already found
	ErrDocExistsAlready = errors.New("Document already exists")
)
