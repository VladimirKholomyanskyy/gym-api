package customerrors

import (
	"errors"
	"fmt"
)

var (
	ErrAccessForbidden = errors.New("access to resource forbidden for user")
	ErrUnauthorized    = errors.New("invalid or missing profile ID")
	ErrInvalidUUID     = errors.New("invalid UUID format")
	ErrEntityNotFound  = errors.New("requested entity not found")
)

type ErrInvalidPosition struct {
	position int
	total    int64
}

func (e ErrInvalidPosition) Error() string {
	return fmt.Sprintf("invalid position %d, must be between 1 and %d", e.position, e.total)
}

func NewErrInvalidPosition(position int, total int64) error {
	return &ErrInvalidPosition{position: position, total: total}
}
