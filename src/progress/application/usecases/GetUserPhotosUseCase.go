package usecases

import (
	"context"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type GetUserPhotosUseCase struct {
	repo repositories.ProgressPhotoRepository
}

func NewGetUserPhotosUseCase(repo repositories.ProgressPhotoRepository) *GetUserPhotosUseCase {
	return &GetUserPhotosUseCase{repo: repo}
}

func (uc *GetUserPhotosUseCase) Execute(ctx context.Context, userID uint, limit, offset int) (*dtos.GetPhotosResponse, error) {
	photos, total, err := uc.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := &dtos.GetPhotosResponse{
		Photos: make([]dtos.PhotoResponse, len(photos)),
		Total:   total,
	}

	for i, p := range photos {
		response.Photos[i] = dtos.PhotoResponse{
			ID:       p.ID,
			Type:     p.Type,
			ImageURL: p.ImageURL,
			Date:     p.Date,
		}
	}

	return response, nil
}
