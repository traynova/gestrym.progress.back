package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gestrym-progress/src/progress/domain/ports"
	"net/http"

	"github.com/spf13/viper"
)

type AIServiceAdapter struct {
	baseURL string
}

func NewAIServiceAdapter() ports.AIService {
	url := viper.GetString("AI_SERVICE_URL")
	if url == "" {
		url = "http://gestrym-ai-service:8080/ai"
	}
	return &AIServiceAdapter{baseURL: url}
}

func (a *AIServiceAdapter) AdaptTraining(ctx context.Context, userID uint) error {
	url := fmt.Sprintf("%s/adapt-training-plan", a.baseURL)
	return a.callAI(ctx, url, userID)
}

func (a *AIServiceAdapter) AdaptNutrition(ctx context.Context, userID uint) error {
	url := fmt.Sprintf("%s/adapt-meal-plan", a.baseURL)
	return a.callAI(ctx, url, userID)
}

func (a *AIServiceAdapter) callAI(ctx context.Context, url string, userID uint) error {
	body, err := json.Marshal(map[string]interface{}{
		"user_id": userID,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("ai service returned status: %d", resp.StatusCode)
	}

	return nil
}
