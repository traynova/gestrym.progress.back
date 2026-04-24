package models

import (
	"time"
)

type BodyMetrics struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"userId"`
	Date       time.Time `gorm:"not null" json:"date"`
	Weight     float64   `gorm:"not null" json:"weight"`
	Height     float64   `gorm:"not null" json:"height"`
	BodyFat    float64   `json:"bodyFat"`
	MuscleMass float64   `json:"muscleMass"`
	CreatedAt  time.Time `json:"createdAt"`
}
