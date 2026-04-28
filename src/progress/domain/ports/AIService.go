package ports

import "context"

type AIService interface {
	AdaptTraining(ctx context.Context, userID uint) error
	AdaptNutrition(ctx context.Context, userID uint) error
}
