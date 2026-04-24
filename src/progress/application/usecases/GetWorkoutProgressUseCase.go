package usecases

import (
	"context"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type GetWorkoutProgressUseCase struct {
	repo repositories.WorkoutProgressRepository
}

func NewGetWorkoutProgressUseCase(repo repositories.WorkoutProgressRepository) *GetWorkoutProgressUseCase {
	return &GetWorkoutProgressUseCase{repo: repo}
}

func (uc *GetWorkoutProgressUseCase) Execute(ctx context.Context, userID uint, limit, offset int) (*dtos.GetWorkoutProgressResponse, error) {
	progress, total, err := uc.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := &dtos.GetWorkoutProgressResponse{
		Progress: make([]dtos.WorkoutProgressResponse, len(progress)),
		Total:    total,
	}

	for i, p := range progress {
		response.Progress[i] = dtos.WorkoutProgressResponse{
			ID:        p.ID,
			WorkoutID: p.WorkoutID,
			Date:      p.Date,
			Duration:  p.Duration,
			Notes:     p.Notes,
		}
	}

	return response, nil
}
