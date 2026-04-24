package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
)

type WorkoutProgressRepository interface {
	Create(ctx context.Context, progress *models.WorkoutProgress) error
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.WorkoutProgress, int64, error)
}
