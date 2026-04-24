package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
)

type BodyMetricsRepository interface {
	Create(ctx context.Context, metrics *models.BodyMetrics) error
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.BodyMetrics, int64, error)
	FindLatestByUserID(ctx context.Context, userID uint) (*models.BodyMetrics, error)
	FindEarliestByUserID(ctx context.Context, userID uint) (*models.BodyMetrics, error)
}
