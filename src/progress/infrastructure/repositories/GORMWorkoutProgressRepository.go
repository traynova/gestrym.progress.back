package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/domain/repositories"
	"gorm.io/gorm"
)

type GORMWorkoutProgressRepository struct {
	db *gorm.DB
}

func NewGORMWorkoutProgressRepository(db *gorm.DB) repositories.WorkoutProgressRepository {
	return &GORMWorkoutProgressRepository{db: db}
}

func (r *GORMWorkoutProgressRepository) Create(ctx context.Context, progress *models.WorkoutProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}

func (r *GORMWorkoutProgressRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.WorkoutProgress, int64, error) {
	var progress []models.WorkoutProgress
	var total int64

	query := r.db.WithContext(ctx).Model(&models.WorkoutProgress{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Order("date desc").Limit(limit).Offset(offset).Find(&progress).Error
	return progress, total, err
}
