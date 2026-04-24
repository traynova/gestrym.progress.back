package models

import (
	"time"
)

type ProgressPhoto struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	Type      string    `gorm:"size:50;not null" json:"type"` // front, back, side
	ImageURL  string    `gorm:"size:500;not null" json:"imageUrl"`
	Date      time.Time `gorm:"not null" json:"date"`
	CreatedAt time.Time `json:"createdAt"`
}
