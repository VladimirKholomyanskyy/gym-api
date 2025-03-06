package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TrainingProgram struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	ProfileID   string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type Workout struct {
	ID                string `gorm:"primaryKey"`
	Name              string
	TrainingProgramID string
	Exercises         []WorkoutExercise `gorm:"constraint:OnDelete:CASCADE"`
	Position          int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

type WorkoutExercise struct {
	ID         string `gorm:"primaryKey"`
	WorkoutID  string
	ExerciseID string
	Sets       int
	Reps       int
	Position   int
	Exercise   Exercise `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type Exercise struct {
	ID              string `gorm:"primaryKey"`
	Name            string
	PrimaryMuscle   string
	SecondaryMuscle pq.StringArray `gorm:"type:text[]"` // PostgreSQL array type
	Equipment       string
	Description     string
}

type ScheduledWorkout struct {
	ID        string `gorm:"primaryKey"`
	ProfileID string
	WorkoutID string
	Workout   Workout `gorm:"constraint:OnDelete:CASCADE;"`
	Date      time.Time
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
