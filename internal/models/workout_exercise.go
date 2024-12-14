package models

import "gorm.io/gorm"

type WorkoutExercise struct {
	gorm.Model
	WorkoutID  uint
	ExerciseID uint
	Sets       int
	Reps       int
	Weight     float64
}

type CreateWorkoutExerciseRequest struct {
	WorkoutID  uint    `json:"workout_id"`
	ExerciseID uint    `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
}

type WorkoutExerciseResponse struct {
	ID         uint `json:"id"`
	WorkoutID  uint `json:"workout_id"`
	ExerciseID uint `json:"exercise_id"`
	Sets       int  `json:"sets"`
	Reps       int  `json:"reps"`
}
