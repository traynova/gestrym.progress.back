package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/domain/repositories"
	"gorm.io/gorm"
)

type GORMBodyMetricsRepository struct {
	db *gorm.DB
}

func NewGORMBodyMetricsRepository(db *gorm.DB) repositories.BodyMetricsRepository {
	return &GORMBodyMetricsRepository{db: db}
}

func (r *GORMBodyMetricsRepository) Create(ctx context.Context, metrics *models.BodyMetrics) error {
	return r.db.WithContext(ctx).Create(metrics).Error
}

func (r *GORMBodyMetricsRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.BodyMetrics, int64, error) {
	var metrics []models.BodyMetrics
	var total int64

	query := r.db.WithContext(ctx).Model(&models.BodyMetrics{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Order("date desc").Limit(limit).Offset(offset).Find(&metrics).Error
	return metrics, total, err
}

func (r *GORMBodyMetricsRepository) FindLatestByUserID(ctx context.Context, userID uint) (*models.BodyMetrics, error) {
	var metrics models.BodyMetrics
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date desc").First(&metrics).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &metrics, nil
}

func (r *GORMBodyMetricsRepository) FindEarliestByUserID(ctx context.Context, userID uint) (*models.BodyMetrics, error) {
	var metrics models.BodyMetrics
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date asc").First(&metrics).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &metrics, nil
}
