package account

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ExternalID string
	Age        int
	Weight     float64
	Height     float64
}

type CreateUserRequest struct {
	ExternalID string  `json:"external_id"`
	Email      string  `json:"email"`
	Age        int     `json:"age"`
	Weight     float64 `json:"weight"`
	Height     float64 `json:"height"`
}
