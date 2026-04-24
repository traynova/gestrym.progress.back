package usecases

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type CreateBodyMetricsUseCase struct {
	repo repositories.BodyMetricsRepository
}

func NewCreateBodyMetricsUseCase(repo repositories.BodyMetricsRepository) *CreateBodyMetricsUseCase {
	return &CreateBodyMetricsUseCase{repo: repo}
}

func (uc *CreateBodyMetricsUseCase) Execute(ctx context.Context, req dtos.CreateMetricsRequest) error {
	metrics := &models.BodyMetrics{
		UserID:     req.UserID,
		Date:       req.Date,
		Weight:     req.Weight,
		Height:     req.Height,
		BodyFat:    req.BodyFat,
		MuscleMass: req.MuscleMass,
	}
	return uc.repo.Create(ctx, metrics)
}
