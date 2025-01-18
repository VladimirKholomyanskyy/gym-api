package training

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TrainingProgram struct {
	gorm.Model
	Name        string
	UserID      uint
	Description string
}

type Workout struct {
	gorm.Model
	Name              string
	TrainingProgramID uint
	Exercises         []WorkoutExercise
}

type WorkoutExercise struct {
	gorm.Model
	WorkoutID  uint
	ExerciseID uint
	Sets       int
	Reps       int
	Weight     float64
	Exercise   Exercise `gorm:"foreignKey:ID;references:ExerciseID"`
}

type Exercise struct {
	ID              uint `gorm:"primaryKey"`
	Name            string
	PrimaryMuscle   string
	SecondaryMuscle pq.StringArray `gorm:"type:text[]"` // PostgreSQL array type
	Equipment       string
	Description     string
}
