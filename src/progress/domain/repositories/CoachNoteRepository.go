package repositories

import (
	"context"
	"gestrym-progress/src/common/models"
)

type CoachNoteRepository interface {
	Create(ctx context.Context, note *models.CoachNote) error
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]models.CoachNote, int64, error)
}
