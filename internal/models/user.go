package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Age      int
	Weight   float64
	Height   float64
}

type CreateUserRequest struct {
	Username string  `json:"name"`
	Email    string  `json:"email"`
	Age      int     `json:"age"`
	Weight   float64 `json:"weight"`
	Height   float64 `json:"height"`
}
