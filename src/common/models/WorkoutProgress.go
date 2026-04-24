package models

import (
	"time"
)

type WorkoutProgress struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	WorkoutID uint      `gorm:"not null;index" json:"workoutId"`
	Date      time.Time `gorm:"not null" json:"date"`
	Duration  int       `json:"duration"` // in minutes
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}
