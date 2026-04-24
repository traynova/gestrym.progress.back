package dtos

type ProgressComparisonResponse struct {
	FirstPhoto    *PhotoResponse  `json:"firstPhoto"`
	LatestPhoto   *PhotoResponse  `json:"latestPhoto"`
	FirstMetrics  *MetricResponse `json:"firstMetrics"`
	LatestMetrics *MetricResponse `json:"latestMetrics"`
}
