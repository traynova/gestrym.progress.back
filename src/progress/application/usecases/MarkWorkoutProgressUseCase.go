package usecases

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type MarkWorkoutProgressUseCase struct {
	repo repositories.WorkoutProgressRepository
}

func NewMarkWorkoutProgressUseCase(repo repositories.WorkoutProgressRepository) *MarkWorkoutProgressUseCase {
	return &MarkWorkoutProgressUseCase{repo: repo}
}

func (uc *MarkWorkoutProgressUseCase) Execute(ctx context.Context, req dtos.MarkWorkoutProgressRequest) error {
	progress := &models.WorkoutProgress{
		UserID:    req.UserID,
		WorkoutID: req.WorkoutID,
		Date:      req.Date,
		Duration:  req.Duration,
		Notes:     req.Notes,
	}
	return uc.repo.Create(ctx, progress)
}
