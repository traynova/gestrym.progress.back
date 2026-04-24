package ports

import (
	"context"
	"mime/multipart"
)

type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}
