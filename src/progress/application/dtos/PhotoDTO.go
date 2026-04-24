package dtos

import "time"

type UploadPhotoRequest struct {
	UserID uint      `form:"userId"`
	Type   string    `form:"type" binding:"required,oneof=front back side"`
	Date   time.Time `form:"date" binding:"required" time_format:"2006-01-02"`
}

type PhotoResponse struct {
	ID       uint      `json:"id"`
	Type     string    `json:"type"`
	ImageURL string    `json:"imageUrl"`
	Date     time.Time `json:"date"`
}

type GetPhotosResponse struct {
	Photos []PhotoResponse `json:"photos"`
	Total  int64           `json:"total"`
}
