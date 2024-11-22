package models

import "time"

type Workout struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	Name              string            `gorm:"type:varchar(255);not null" json:"name"`
	TrainingProgramID uint              `gorm:"not null;index"` // Foreign key to training_programs
	CreatedAt         time.Time         `gorm:"autoCreateTime"`
	Exercises         []WorkoutExercise `gorm:"foreignKey:WorkoutID"` // Association with WorkoutExercise
}

type CreateWorkoutRequest struct {
	Name string `json:"id"`
}

type WorkoutInput struct {
	Name              string
	TrainingProgramID uint
	UserID            uint
}
