package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/domain/repositories"
	"gorm.io/gorm"
)

type GORMCoachNoteRepository struct {
	db *gorm.DB
}

func NewGORMCoachNoteRepository(db *gorm.DB) repositories.CoachNoteRepository {
	return &GORMCoachNoteRepository{db: db}
}

func (r *GORMCoachNoteRepository) Create(ctx context.Context, note *models.CoachNote) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *GORMCoachNoteRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.CoachNote, int64, error) {
	var notes []models.CoachNote
	var total int64

	query := r.db.WithContext(ctx).Model(&models.CoachNote{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&notes).Error
	return notes, total, err
}
