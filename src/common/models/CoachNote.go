package models

import (
	"time"
)

type CoachNote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	TrainerID uint      `gorm:"not null;index" json:"trainerId"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
