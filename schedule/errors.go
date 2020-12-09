package schedule

import "errors"

var (
	// ErrDocExistsAlready is thrown if document is already found
	ErrDocExistsAlready = errors.New("Document already exists")
	// ErrScheduleNotFound is thrown if schedule is not found
	ErrScheduleNotFound = errors.New("Schedule not found")
	// ErrAlreadyInGroup is returned if someone is already in a group
	ErrAlreadyInGroup = errors.New("User already in this group")
	// ErrGroupsNotFound is thrown if there are no groups found
	ErrGroupsNotFound = errors.New("There were no groups in this season")
)
