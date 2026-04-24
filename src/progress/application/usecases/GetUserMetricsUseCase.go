package usecases

import (
	"context"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type GetUserMetricsUseCase struct {
	repo repositories.BodyMetricsRepository
}

func NewGetUserMetricsUseCase(repo repositories.BodyMetricsRepository) *GetUserMetricsUseCase {
	return &GetUserMetricsUseCase{repo: repo}
}

func (uc *GetUserMetricsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) (*dtos.GetMetricsResponse, error) {
	metrics, total, err := uc.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := &dtos.GetMetricsResponse{
		Metrics: make([]dtos.MetricResponse, len(metrics)),
		Total:   total,
	}

	for i, m := range metrics {
		response.Metrics[i] = dtos.MetricResponse{
			ID:         m.ID,
			Date:       m.Date,
			Weight:     m.Weight,
			BodyFat:    m.BodyFat,
			MuscleMass: m.MuscleMass,
		}
	}

	return response, nil
}
