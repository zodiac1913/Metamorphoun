package service

import (
	"Metamorphoun/config"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type immichRandomSearchRequest struct {
	Size     int      `json:"size,omitempty"`
	Type     string   `json:"type,omitempty"`
	AlbumIDs []string `json:"albumIds,omitempty"`
}

type immichAssetResponse struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	OriginalMimeType string `json:"originalMimeType"`
	OriginalFileName string `json:"originalFileName"`
}

func GetBackgroundImmich(imgItem config.Image) (image.Image, string, error) {
	if strings.TrimSpace(imgItem.APIKey) == "" {
		return nil, "", fmt.Errorf("Immich API key is not configured")
	}

	baseURL, albumIDs, err := normalizeImmichAPIBase(imgItem.Location)
	if err != nil {
		return nil, "", err
	}

	requestBody := immichRandomSearchRequest{
		Size: 1,
		Type: "IMAGE",
	}
	if len(albumIDs) > 0 {
		requestBody.AlbumIDs = albumIDs
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode Immich search request: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+"/search/random", bytes.NewReader(payload))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Immich search request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", strings.TrimSpace(imgItem.APIKey))

	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query Immich random search: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return nil, "", fmt.Errorf("Immich random search returned %s: %s", response.Status, strings.TrimSpace(string(body)))
	}

	var assets []immichAssetResponse
	if err := json.NewDecoder(response.Body).Decode(&assets); err != nil {
		return nil, "", fmt.Errorf("failed to decode Immich search response: %w", err)
	}
	if len(assets) == 0 {
		return nil, "", fmt.Errorf("Immich random search returned no assets")
	}

	asset, err := selectImmichImageAsset(assets)
	if err != nil {
		return nil, "", err
	}

	imageURL := fmt.Sprintf("%s/assets/%s/original", baseURL, asset.ID)
	img, err := loadImmichImageFromURL(imageURL, strings.TrimSpace(imgItem.APIKey))
	if err != nil {
		return nil, "", err
	}

	return img, imageURL, nil
}

func normalizeImmichAPIBase(rawLocation string) (string, []string, error) {
	trimmedLocation := strings.TrimSpace(rawLocation)
	if trimmedLocation == "" {
		return "", nil, fmt.Errorf("Immich location is empty")
	}

	parsedURL, err := url.Parse(trimmedLocation)
	if err != nil {
		return "", nil, fmt.Errorf("invalid Immich server URL: %w", err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", nil, fmt.Errorf("Immich server URL must include http:// or https://")
	}

	albumIDs := splitImmichAlbumIDs(parsedURL.Query().Get("albumIds"))
	if len(albumIDs) == 0 {
		albumIDs = splitImmichAlbumIDs(parsedURL.Query().Get("albumId"))
	}
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""

	cleanPath := strings.TrimRight(parsedURL.Path, "/")
	switch {
	case cleanPath == "":
		cleanPath = "/api"
	case cleanPath == "/api":
		// already normalized
	case strings.HasSuffix(cleanPath, "/api"):
		// already normalized
	default:
		cleanPath += "/api"
	}
	parsedURL.Path = cleanPath

	return strings.TrimRight(parsedURL.String(), "/"), albumIDs, nil
}

func splitImmichAlbumIDs(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	albumIDs := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			albumIDs = append(albumIDs, trimmed)
		}
	}
	return albumIDs
}

func selectImmichImageAsset(assets []immichAssetResponse) (immichAssetResponse, error) {
	for _, asset := range assets {
		if strings.EqualFold(asset.Type, "IMAGE") || strings.HasPrefix(strings.ToLower(asset.OriginalMimeType), "image/") {
			return asset, nil
		}
	}
	return immichAssetResponse{}, fmt.Errorf("Immich random search returned no image assets")
}

func loadImmichImageFromURL(imageURL string, apiKey string) (image.Image, error) {
	request, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Immich asset request: %w", err)
	}
	request.Header.Set("x-api-key", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Immich image: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return nil, fmt.Errorf("Immich image download returned %s: %s", response.Status, strings.TrimSpace(string(body)))
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Immich image: %w", err)
	}
	return img, nil
}
