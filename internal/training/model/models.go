package model

import (
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TrainingProgram struct {
	common.Base
	Name        string
	ProfileID   string
	Description string
}

type Workout struct {
	common.Base
	Name              string
	TrainingProgramID string
	Exercises         []WorkoutExercise `gorm:"constraint:OnDelete:CASCADE"`
	Position          int
}

type WorkoutExercise struct {
	common.Base
	WorkoutID  string
	ExerciseID string
	Sets       int
	Reps       int
	Position   int
	Exercise   Exercise `gorm:"constraint:OnDelete:CASCADE"`
}

type Exercise struct {
	ID              string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name            string
	PrimaryMuscle   string
	SecondaryMuscle pq.StringArray `gorm:"type:text[]"` // PostgreSQL array type
	Equipment       string
	Description     string
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *Exercise) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type ScheduledWorkout struct {
	common.Base
	ProfileID string
	WorkoutID string
	Workout   Workout `gorm:"constraint:OnDelete:CASCADE;"`
	Date      time.Time
	Notes     string
}

type CreateScheduledWorkoutInput struct {
	ProfileID string
	WorkoutID string
	Date      time.Time
	Notes     string
}

type UpdateScheduledWorkoutInput struct {
	ScheduledWorkoutID string
	ProfileID          string
	Date               *time.Time
	Notes              *string
}
