package common

import "errors"

var (
	ErrAccessForbidden = errors.New("access to resource forbidden for user")
	ErrProgramNotFound = errors.New("training program not found")
	ErrUnauthorized    = errors.New("invalid or missing profile ID")
	ErrInvalidUUID     = errors.New("invalid UUID format")
)

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e ForbiddenError) Error() string {
	return e.Message
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{Message: msg}
}

func NewForbiddenError(msg string) error {
	return &ForbiddenError{Message: msg}
}

func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}
