package dtos

import "time"

type CreateMetricsRequest struct {
	UserID     uint      `json:"userId"`
	Date       time.Time `json:"date" binding:"required"`
	Weight     float64   `json:"weight" binding:"required"`
	Height     float64   `json:"height" binding:"required"`
	BodyFat    float64   `json:"bodyFat"`
	MuscleMass float64   `json:"muscleMass"`
}

type MetricResponse struct {
	ID         uint      `json:"id"`
	Date       time.Time `json:"date"`
	Weight     float64   `json:"weight"`
	BodyFat    float64   `json:"bodyFat"`
	MuscleMass float64   `json:"muscleMass,omitempty"`
}

type GetMetricsResponse struct {
	Metrics []MetricResponse `json:"metrics"`
	Total   int64            `json:"total"`
}
