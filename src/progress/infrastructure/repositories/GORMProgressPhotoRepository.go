package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/domain/repositories"
	"gorm.io/gorm"
)

type GORMProgressPhotoRepository struct {
	db *gorm.DB
}

func NewGORMProgressPhotoRepository(db *gorm.DB) repositories.ProgressPhotoRepository {
	return &GORMProgressPhotoRepository{db: db}
}

func (r *GORMProgressPhotoRepository) Create(ctx context.Context, photo *models.ProgressPhoto) error {
	return r.db.WithContext(ctx).Create(photo).Error
}

func (r *GORMProgressPhotoRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.ProgressPhoto, int64, error) {
	var photos []models.ProgressPhoto
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ProgressPhoto{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Order("date desc").Limit(limit).Offset(offset).Find(&photos).Error
	return photos, total, err
}

func (r *GORMProgressPhotoRepository) FindLatestByUserID(ctx context.Context, userID uint) (*models.ProgressPhoto, error) {
	var photo models.ProgressPhoto
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date desc").First(&photo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &photo, nil
}

func (r *GORMProgressPhotoRepository) FindEarliestByUserID(ctx context.Context, userID uint) (*models.ProgressPhoto, error) {
	var photo models.ProgressPhoto
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date asc").First(&photo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &photo, nil
}
