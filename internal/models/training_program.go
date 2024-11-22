package models

import (
	"time"
)

type TrainingProgram struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	UserID      uint      `gorm:"not null" json:"user_id"` // Foreign key to the users table
	Description string    `json:"description,omitempty"`   // Optional field
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Association
	User User `gorm:"foreignKey:UserID" json:"-"`
}

type TrainingProgramCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TrainingProgramInput struct {
	Name        string
	Description string
	UserID      uint
}
