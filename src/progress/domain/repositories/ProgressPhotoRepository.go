package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
)

type ProgressPhotoRepository interface {
	Create(ctx context.Context, photo *models.ProgressPhoto) error
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.ProgressPhoto, int64, error)
	FindLatestByUserID(ctx context.Context, userID uint) (*models.ProgressPhoto, error)
	FindEarliestByUserID(ctx context.Context, userID uint) (*models.ProgressPhoto, error)
}
