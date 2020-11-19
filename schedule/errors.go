package schedule

import "errors"

var (
	// ErrDocExistsAlready is thrown if document is already found
	ErrDocExistsAlready = errors.New("Document already exists")
	// ErrScheduleNotFound is thrown if schedule is not found
	ErrScheduleNotFound = errors.New("Schedule not found")
)
