package models

import "gorm.io/gorm"

type Workout struct {
	gorm.Model
	Name              string
	TrainingProgramID uint
	Exercises         []WorkoutExercise
}

type CreateWorkoutRequest struct {
	Name string `json:"name"`
}

type WorkoutResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
