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
	"net/textproto"
	"strings"

	"github.com/spf13/viper"
)

type StorageServiceAdapter struct {
	baseURL string
}

func NewStorageServiceAdapter() ports.StorageService {
	url := viper.GetString("STORAGE_SERVICE_URL")
	if url == "" {
		url = "http://gestrym-storage:8080"
	}
	// Normalizar: quitar trailing slash
	url = strings.TrimRight(url, "/")
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

	// Detectar el Content-Type real del archivo
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		// Leer los primeros 512 bytes para detección automática
		buf := make([]byte, 512)
		n, _ := src.Read(buf)
		mimeType = http.DetectContentType(buf[:n])
		// Reiniciar el reader (volver al inicio)
		if seeker, ok := src.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		} else {
			// Si no es seekable, recrear desde el FileHeader
			src.Close()
			src, err = file.Open()
			if err != nil {
				return "", fmt.Errorf("error reabriendo archivo para detección MIME: %w", err)
			}
		}
	}

	// Crear la parte multipart con el Content-Type correcto
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="files"; filename="%s"`, file.Filename))
	h.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, src)
	if err != nil {
		return "", err
	}
	writer.Close()

	uploadURL := fmt.Sprintf("%s/gestrym-storage/internal/files/upload", a.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// API Key requerida por el grupo /internal
	if apiKey := viper.GetString("STORAGE_SERVICE_API_KEY"); apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("storage service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		CollectionID string `json:"collection_id"`
		Data         struct {
			URL string `json:"url"`
		} `json:"data"`
		URL string `json:"url"` // Fallback
	}

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error parsing storage response: %v | body: %s", err, string(respBody))
	}

	// El storage devuelve collection_id; intentamos URL como fallback
	finalURL := result.URL
	if finalURL == "" {
		finalURL = result.Data.URL
	}
	if finalURL == "" && result.CollectionID != "" {
		// Devolvemos el collection_id para que el caller pueda obtener las URLs
		finalURL = result.CollectionID
	}

	if finalURL == "" {
		return "", fmt.Errorf("could not find URL or collection_id in storage response: %s", string(respBody))
	}

	return finalURL, nil
}
