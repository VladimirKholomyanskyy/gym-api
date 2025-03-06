package repository

import "errors"

var (
	// ErrInvalidPagination is returned when pagination parameters are invalid
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)

// validatePagination checks if pagination parameters are valid
func validatePagination(page, pageSize int) error {
	if page < 1 || pageSize < 1 {
		return ErrInvalidPagination
	}
	return nil
}
