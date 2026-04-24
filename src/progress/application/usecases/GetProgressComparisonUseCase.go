package usecases

import (
	"context"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type GetProgressComparisonUseCase struct {
	metricsRepo repositories.BodyMetricsRepository
	photosRepo  repositories.ProgressPhotoRepository
}

func NewGetProgressComparisonUseCase(metricsRepo repositories.BodyMetricsRepository, photosRepo repositories.ProgressPhotoRepository) *GetProgressComparisonUseCase {
	return &GetProgressComparisonUseCase{
		metricsRepo: metricsRepo,
		photosRepo:  photosRepo,
	}
}

func (uc *GetProgressComparisonUseCase) Execute(ctx context.Context, userID uint) (*dtos.ProgressComparisonResponse, error) {
	firstMetric, _ := uc.metricsRepo.FindEarliestByUserID(ctx, userID)
	latestMetric, _ := uc.metricsRepo.FindLatestByUserID(ctx, userID)
	firstPhoto, _ := uc.photosRepo.FindEarliestByUserID(ctx, userID)
	latestPhoto, _ := uc.photosRepo.FindLatestByUserID(ctx, userID)

	response := &dtos.ProgressComparisonResponse{}

	if firstMetric != nil {
		response.FirstMetrics = &dtos.MetricResponse{
			ID:         firstMetric.ID,
			Date:       firstMetric.Date,
			Weight:     firstMetric.Weight,
			BodyFat:    firstMetric.BodyFat,
			MuscleMass: firstMetric.MuscleMass,
		}
	}

	if latestMetric != nil {
		response.LatestMetrics = &dtos.MetricResponse{
			ID:         latestMetric.ID,
			Date:       latestMetric.Date,
			Weight:     latestMetric.Weight,
			BodyFat:    latestMetric.BodyFat,
			MuscleMass: latestMetric.MuscleMass,
		}
	}

	if firstPhoto != nil {
		response.FirstPhoto = &dtos.PhotoResponse{
			ID:       firstPhoto.ID,
			Type:     firstPhoto.Type,
			ImageURL: firstPhoto.ImageURL,
			Date:     firstPhoto.Date,
		}
	}

	if latestPhoto != nil {
		response.LatestPhoto = &dtos.PhotoResponse{
			ID:       latestPhoto.ID,
			Type:     latestPhoto.Type,
			ImageURL: latestPhoto.ImageURL,
			Date:     latestPhoto.Date,
		}
	}

	return response, nil
}
