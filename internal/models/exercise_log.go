package models

import "time"

type ExerciseLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID  uint      `gorm:"not null" json:"session_id"`
	ExerciseID uint      `gorm:"not null" json:"exercise_id"`
	SetNumber  int       `gorm:"not null;check:set_number > 0" json:"set_number"`
	Reps       int       `gorm:"not null;check:reps >= 0" json:"reps"`
	Weight     float64   `gorm:"type:decimal(5,2);not null;check:weight >= 0" json:"weight"`
	LoggedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"logged_at"`
}

type LogExerciseRequest struct {
	ExerciseID    uint    `json:"exercise_id" binding:"required"`
	SetNumber     int     `json:"set_number" binding:"required"`
	RepsCompleted int     `json:"reps_completed" binding:"required"`
	WeightUsed    float64 `json:"weight_used" binding:"required"`
}

type LogExerciseResponse struct {
	LogID         uint    `json:"log_id"`
	ExerciseID    uint    `json:"exercise_id"`
	SetNumber     int     `json:"set_number"`
	RepsCompleted int     `json:"reps_completed"`
	WeightUsed    float64 `json:"weight_used"`
	LoggedAt      string  `json:"logged_at"`
}
