package usecases

import (
	"context"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type GetUserNotesUseCase struct {
	repo repositories.CoachNoteRepository
}

func NewGetUserNotesUseCase(repo repositories.CoachNoteRepository) *GetUserNotesUseCase {
	return &GetUserNotesUseCase{repo: repo}
}

func (uc *GetUserNotesUseCase) Execute(ctx context.Context, userID uint, limit, offset int) (*dtos.GetNotesResponse, error) {
	notes, total, err := uc.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := &dtos.GetNotesResponse{
		Notes: make([]dtos.NoteResponse, len(notes)),
		Total: total,
	}

	for i, n := range notes {
		response.Notes[i] = dtos.NoteResponse{
			ID:        n.ID,
			Message:   n.Message,
			TrainerID: n.TrainerID,
			Date:      n.CreatedAt,
		}
	}

	return response, nil
}
