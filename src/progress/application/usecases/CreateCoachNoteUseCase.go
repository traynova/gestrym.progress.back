package usecases

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type CreateCoachNoteUseCase struct {
	repo repositories.CoachNoteRepository
}

func NewCreateCoachNoteUseCase(repo repositories.CoachNoteRepository) *CreateCoachNoteUseCase {
	return &CreateCoachNoteUseCase{repo: repo}
}

func (uc *CreateCoachNoteUseCase) Execute(ctx context.Context, req dtos.CreateNoteRequest) error {
	note := &models.CoachNote{
		UserID:    req.UserID,
		TrainerID: req.TrainerID,
		Message:   req.Message,
	}
	return uc.repo.Create(ctx, note)
}
