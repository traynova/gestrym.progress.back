package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gestrym-progress/src/progress/domain/ports"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/spf13/viper"
)

type StorageServiceAdapter struct {
	baseURL string
}

func NewStorageServiceAdapter() ports.StorageService {
	url := viper.GetString("STORAGE_SERVICE_URL")
	if url == "" {
		// En un entorno real, esto debería estar en las variables de entorno
		url = "https://gestrym-storage-back.onrender.com/gestrym-storage"
	}
	return &StorageServiceAdapter{baseURL: url}
}

func (a *StorageServiceAdapter) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, src)
	if err != nil {
		return "", err
	}
	writer.Close()

	uploadURL := fmt.Sprintf("%s/public/upload", a.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("storage service returned status: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
		URL string `json:"url"` // Fallback if data wrapper is not used
	}
	
	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	finalURL := result.URL
	if finalURL == "" {
		finalURL = result.Data.URL
	}

	if finalURL == "" {
		return "", fmt.Errorf("could not find URL in storage service response")
	}

	return finalURL, nil
}
