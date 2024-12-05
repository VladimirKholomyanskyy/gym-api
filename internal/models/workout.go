package models

import "gorm.io/gorm"

type Workout struct {
	gorm.Model
	Name              string
	TrainingProgramID uint
	Exercises         []WorkoutExercise `gorm:"many2many:workout_exercises"` // Association with WorkoutExercise
}

type CreateWorkoutRequest struct {
	Name string `json:"name"`
}
