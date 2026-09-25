package usecases

import (
	"context"
	"fmt"
	"gestrym-progress/src/common/models"
	"gestrym-progress/src/common/utils"
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
	logger         utils.ILogger
}

func NewUploadProgressPhotoUseCase(repo repositories.ProgressPhotoRepository, storageService ports.StorageService, aiService ports.AIService) *UploadProgressPhotoUseCase {
	return &UploadProgressPhotoUseCase{
		repo:           repo,
		storageService: storageService,
		aiService:      aiService,
		logger:         utils.NewLogger(),
	}
}

func (uc *UploadProgressPhotoUseCase) Execute(ctx context.Context, req dtos.UploadPhotoRequest, file *multipart.FileHeader) error {
	uc.logger.Info("[PHOTO_UPLOAD] iniciando upload | userID:%d type:%s filename:%s size:%d",
		req.UserID, req.Type, file.Filename, file.Size)

	imageURL, err := uc.storageService.UploadFile(ctx, file)
	if err != nil {
		uc.logger.Error("[PHOTO_UPLOAD] error en storage.UploadFile | userID:%d filename:%s | error:%v",
			req.UserID, file.Filename, err)
		return fmt.Errorf("error al subir archivo al storage: %w", err)
	}
	uc.logger.Info("[PHOTO_UPLOAD] storage respondió OK | imageURL/collectionID:%q", imageURL)

	photo := &models.ProgressPhoto{
		UserID:   req.UserID,
		Type:     req.Type,
		ImageURL: imageURL,
		Date:     req.Date,
	}

	uc.logger.Info("[PHOTO_UPLOAD] guardando en DB | userID:%d imageURL:%q", req.UserID, imageURL)
	err = uc.repo.Create(ctx, photo)
	if err != nil {
		uc.logger.Error("[PHOTO_UPLOAD] error al guardar en DB | userID:%d | error:%v", req.UserID, err)
		return fmt.Errorf("error al guardar foto en base de datos: %w", err)
	}
	uc.logger.Info("[PHOTO_UPLOAD] foto guardada exitosamente | photoID:%d userID:%d", photo.ID, req.UserID)

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
