package models

import (
	"github.com/lib/pq"
)

type Exercise struct {
	ID              uint `gorm:"primaryKey"`
	Name            string
	PrimaryMuscle   string
	SecondaryMuscle pq.StringArray `gorm:"type:text[]"` // PostgreSQL array type
	Equipment       string
	Description     string
}
