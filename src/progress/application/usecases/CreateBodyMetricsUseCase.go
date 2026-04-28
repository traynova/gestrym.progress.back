package usecases

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/ports"
	"gestrym-progress/src/progress/domain/repositories"
	"log"
	"math"
)

type CreateBodyMetricsUseCase struct {
	repo      repositories.BodyMetricsRepository
	aiService ports.AIService
}

func NewCreateBodyMetricsUseCase(repo repositories.BodyMetricsRepository, aiService ports.AIService) *CreateBodyMetricsUseCase {
	return &CreateBodyMetricsUseCase{
		repo:      repo,
		aiService: aiService,
	}
}

func (uc *CreateBodyMetricsUseCase) Execute(ctx context.Context, req dtos.CreateMetricsRequest) error {
	// Get latest metrics for comparison (Threshold Bonus)
	latest, _ := uc.repo.FindLatestByUserID(ctx, req.UserID)

	metrics := &models.BodyMetrics{
		UserID:     req.UserID,
		Date:       req.Date,
		Weight:     req.Weight,
		Height:     req.Height,
		BodyFat:    req.BodyFat,
		MuscleMass: req.MuscleMass,
	}

	err := uc.repo.Create(ctx, metrics)
	if err != nil {
		return err
	}

	// Trigger AI Adaptation (Async and ignore failure)
	go func() {
		// Threshold check: trigger only if weight change > 0.5kg (example threshold)
		shouldTrigger := true
		if latest != nil {
			weightDiff := math.Abs(metrics.Weight - latest.Weight)
			if weightDiff < 0.5 {
				shouldTrigger = false
			}
		}

		if shouldTrigger {
			bgCtx := context.Background()
			errTr := uc.aiService.AdaptTraining(bgCtx, metrics.UserID)
			if errTr != nil {
				log.Printf("Error triggering AI training adaptation: %v", errTr)
			}
			errNu := uc.aiService.AdaptNutrition(bgCtx, metrics.UserID)
			if errNu != nil {
				log.Printf("Error triggering AI nutrition adaptation: %v", errNu)
			}
		}
	}()

	return nil
}
