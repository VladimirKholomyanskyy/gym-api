package models

import (
	"github.com/lib/pq"
)

type Exercise struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"unique;not null" json:"name"`
	PrimaryMuscle   string         `gorm:"not null" json:"primary_muscle"`
	SecondaryMuscle pq.StringArray `gorm:"type:text[]" json:"secondary_muscle"` // PostgreSQL array type
	Equipment       string         `gorm:"not null" json:"equipment"`
	Description     string         `json:"description"`
}
