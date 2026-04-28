package usecases

import (
	"context"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/ports"
	"gestrym-progress/src/progress/domain/repositories"
	"log"
	"mime/multipart"
)

type UploadProgressPhotoUseCase struct {
	repo           repositories.ProgressPhotoRepository
	storageService ports.StorageService
	aiService      ports.AIService
}

func NewUploadProgressPhotoUseCase(repo repositories.ProgressPhotoRepository, storageService ports.StorageService, aiService ports.AIService) *UploadProgressPhotoUseCase {
	return &UploadProgressPhotoUseCase{
		repo:           repo,
		storageService: storageService,
		aiService:      aiService,
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

	err = uc.repo.Create(ctx, photo)
	if err != nil {
		return err
	}

	// Trigger AI Adaptation (Async and ignore failure)
	go func() {
		bgCtx := context.Background()
		errTr := uc.aiService.AdaptTraining(bgCtx, photo.UserID)
		if errTr != nil {
			log.Printf("Error triggering AI training adaptation from photo: %v", errTr)
		}
		errNu := uc.aiService.AdaptNutrition(bgCtx, photo.UserID)
		if errNu != nil {
			log.Printf("Error triggering AI nutrition adaptation from photo: %v", errNu)
		}
	}()

	return nil
}
