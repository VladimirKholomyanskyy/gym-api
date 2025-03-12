package account

import (
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"gorm.io/gorm"
)

type Profile struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ExternalID string `gorm:"uniqueIndex;not null"`
	Sex        *openapi.Sex
	Birthday   *time.Time
	Weight     *float64
	Height     *float64
	AvatarURL  *string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt
}

type Setting struct {
	ID                   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProfileID            string
	Language             string
	MeasurementUnits     openapi.MeasurementUnits
	Timezone             string
	NotificationsEnabled bool
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt
}
