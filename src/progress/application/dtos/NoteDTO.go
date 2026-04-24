package dtos

import "time"

type CreateNoteRequest struct {
	UserID    uint   `json:"userId" binding:"required"`
	TrainerID uint   `json:"trainerId"`
	Message   string `json:"message" binding:"required"`
}

type NoteResponse struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	TrainerID uint      `json:"trainerId"`
	Date      time.Time `json:"date"`
}

type GetNotesResponse struct {
	Notes []NoteResponse `json:"notes"`
	Total int64          `json:"total"`
}
