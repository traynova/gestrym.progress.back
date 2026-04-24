package usecases

import (
	"context"
	"gestrym-progress/src/progress/application/dtos"
	"gestrym-progress/src/progress/domain/repositories"
)

type GetWeightChartUseCase struct {
	repo repositories.BodyMetricsRepository
}

func NewGetWeightChartUseCase(repo repositories.BodyMetricsRepository) *GetWeightChartUseCase {
	return &GetWeightChartUseCase{repo: repo}
}

func (uc *GetWeightChartUseCase) Execute(ctx context.Context, userID uint) (*dtos.WeightChartResponse, error) {
	metrics, _, err := uc.repo.FindByUserID(ctx, userID, 100, 0) // Obtener las últimas 100 entradas
	if err != nil {
		return nil, err
	}

	points := make([]dtos.WeightChartPoint, len(metrics))
	// Invertimos el orden para que sea cronológico ascendente para la gráfica
	for i := range metrics {
		m := metrics[len(metrics)-1-i]
		points[i] = dtos.WeightChartPoint{
			Date:   m.Date,
			Weight: m.Weight,
		}
	}

	return &dtos.WeightChartResponse{Points: points}, nil
}
