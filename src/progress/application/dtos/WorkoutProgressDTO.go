package dtos

import "time"

type MarkWorkoutProgressRequest struct {
	UserID    uint      `json:"userId"`
	WorkoutID uint      `json:"workoutId" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Duration  int       `json:"duration"`
	Notes     string    `json:"notes"`
}

type WorkoutProgressResponse struct {
	ID        uint      `json:"id"`
	WorkoutID uint      `json:"workoutId"`
	Date      time.Time `json:"date"`
	Duration  int       `json:"duration"`
	Notes     string    `json:"notes"`
}

type GetWorkoutProgressResponse struct {
	Progress []WorkoutProgressResponse `json:"progress"`
	Total    int64                     `json:"total"`
}
