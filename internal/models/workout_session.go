package models

import (
	"time"

	"gorm.io/datatypes"
)

type WorkoutSession struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	WorkoutID   uint           `gorm:"not null" json:"workout_id"`
	Snapshot    datatypes.JSON `gorm:"type:jsonb;not null" json:"snapshot"` // JSONB for workout snapshot
	StartedAt   time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"started_at"`
	CompletedAt *time.Time     `gorm:"default:null" json:"completed_at"` // Nullable
	Logs        []ExerciseLog  `gorm:"foreignKey:SessionID" json:"logs"` // Association
}

type StartWorkoutSessionRequest struct {
	WorkoutID uint `json:"workout_id"`
}

type StartWorkoutSessionResponse struct {
	SessionID uint   `json:"session_id"`
	StartedAt string `json:"started_at"`
}

type GetWorkoutSessionResponse struct {
	SessionID       uint    `json:"session_id"`
	StartedAt       string  `json:"started_at"`
	CompletedAt     string  `json:"completed_at"`
	WorkoutSnapshot Workout `json:"workout_snapshot"`
}
type GetWorkoutSessionLogsResponse struct {
	SessionID       uint                   `json:"session_id"`
	WorkoutSnapshot map[string]interface{} `json:"workout_snapshot"`
	Logs            []WorkoutLogResponse   `json:"logs"`
	StartedAt       string                 `json:"started_at"`
	CompletedAt     *string                `json:"completed_at"`
}

type WorkoutLogResponse struct {
	ExerciseName  uint    `json:"exercise_name"`
	SetNumber     int     `json:"set_number"`
	RepsCompleted int     `json:"reps_completed"`
	WeightUsed    float64 `json:"weight_used"`
	LoggedAt      string  `json:"logged_at"`
}
