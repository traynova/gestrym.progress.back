package dtos

import "time"

type WeightChartPoint struct {
	Date   time.Time `json:"date"`
	Weight float64   `json:"weight"`
}

type WeightChartResponse struct {
	Points []WeightChartPoint `json:"points"`
}
