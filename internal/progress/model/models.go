package model

import (
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"gorm.io/datatypes"
)

type WorkoutSession struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProfileID   string
	WorkoutID   string
	Snapshot    datatypes.JSON `gorm:"type:jsonb;not null"` // JSONB for workout snapshot
	StartedAt   time.Time
	CompletedAt *time.Time
	Logs        []ExerciseLog `gorm:"foreignKey:SessionID" json:"logs"` // Association
}

type ExerciseLog struct {
	common.Base
	ProfileID  string
	SessionID  string
	ExerciseID string
	SetNumber  int
	Reps       int
	Weight     float64
}

type WeightPerDay struct {
	Date        time.Time `json:"date"`
	TotalWeight float64   `json:"total_weight"`
}
