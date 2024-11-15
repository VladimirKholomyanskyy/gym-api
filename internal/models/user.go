package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Age       int       `json:"age" gorm:"not null;check:age >= 0"`
	Weight    float64   `json:"weight" gorm:"not null;check:weight > 0"`
	Height    float64   `json:"height" gorm:"not null;check:height > 0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
