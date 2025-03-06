package account

import (
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"gorm.io/gorm"
)

type Profile struct {
	ID         string `gorm:"primaryKey"`
	ExternalID string `gorm:"uniqueIndex"`
	Sex        *openapi.Sex
	Birthday   *time.Time
	Weight     *float64       // Stored in kg
	Height     *float64       // Stored in meters
	AvatarURL  *string        // Missing from original model
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Setting struct {
	ID                   string                   `gorm:"primaryKey"`
	ProfileID            string                   `gorm:"index"`
	Language             string                   `gorm:"default:'en'"`
	MeasurementUnits     openapi.MeasurementUnits `gorm:"default:'metric'"`
	Timezone             string                   `gorm:"default:'UTC'"`
	NotificationsEnabled bool                     `gorm:"default:true"`
	CreatedAt            time.Time                `gorm:"autoCreateTime"`
	UpdatedAt            time.Time                `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt           `gorm:"index"`
}
