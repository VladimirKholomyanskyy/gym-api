package common

import "github.com/google/uuid"

func IsPageValid(page int32) bool {
	return page >= 1
}

func IsPageSizeValid(pageSize int32) bool {
	return pageSize >= 1 && pageSize <= 100
}

func IsUUIDValid(id string) bool {
	if _, err := uuid.Parse(id); err != nil {
		return false
	}
	return true
}
