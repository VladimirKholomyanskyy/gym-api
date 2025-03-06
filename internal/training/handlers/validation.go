package handlers

import "github.com/google/uuid"

func isPageValid(page int32) bool {
	return page >= 1
}

func isPageSizeValid(pageSize int32) bool {
	return pageSize >= 1 && pageSize <= 100
}

func isUUIDValid(id string) bool {
	if _, err := uuid.Parse(id); err != nil {
		return false
	}
	return true
}
