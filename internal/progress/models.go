package progress

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

type ExerciseLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID  uint      `gorm:"not null" json:"session_id"`
	ExerciseID uint      `gorm:"not null" json:"exercise_id"`
	SetNumber  int       `gorm:"not null;check:set_number > 0" json:"set_number"`
	Reps       int       `gorm:"not null;check:reps >= 0" json:"reps"`
	Weight     float64   `gorm:"type:decimal(5,2);not null;check:weight >= 0" json:"weight"`
	LoggedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"logged_at"`
}
