package models

import (
	"gorm.io/gorm"
)

type TrainingProgram struct {
	gorm.Model
	Name        string
	UserID      uint
	Description string
}

type CreateTrainingProgramRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateTrainingProgramResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
}