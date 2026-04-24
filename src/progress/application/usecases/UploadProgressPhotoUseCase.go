package usecases

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/ports"
	"gestrym-progress/src/progress/domain/repositories"
	"mime/multipart"
)

type UploadProgressPhotoUseCase struct {
	repo           repositories.ProgressPhotoRepository
	storageService ports.StorageService
}

func NewUploadProgressPhotoUseCase(repo repositories.ProgressPhotoRepository, storageService ports.StorageService) *UploadProgressPhotoUseCase {
	return &UploadProgressPhotoUseCase{
		repo:           repo,
		storageService: storageService,
	}
}

func (uc *UploadProgressPhotoUseCase) Execute(ctx context.Context, req dtos.UploadPhotoRequest, file *multipart.FileHeader) error {
	imageURL, err := uc.storageService.UploadFile(ctx, file)
	if err != nil {
		return err
	}

	photo := &models.ProgressPhoto{
		UserID:   req.UserID,
		Type:     req.Type,
		ImageURL: imageURL,
		Date:     req.Date,
	}

	return uc.repo.Create(ctx, photo)
}
