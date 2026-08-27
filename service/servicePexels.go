package service

import (
	"Metamorphoun/config"
	"encoding/json"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"net/url"
)

type pexelsResponse struct {
	Photos []struct {
		Source struct {
			Landscape string `json:"landscape"`
			Original  string `json:"original"`
		} `json:"src"`
	} `json:"photos"`
}

func GetBackgroundPexels(imgItem config.Image) (image.Image, string, error) {
	if imgItem.APIKey == "" {
		return nil, "", fmt.Errorf("Pexels API key is not configured")
	}

	endpoint, err := url.Parse("https://api.pexels.com/v1/curated")
	if err != nil {
		return nil, "", err
	}
	query := endpoint.Query()
	query.Set("page", fmt.Sprintf("%d", rand.Intn(100)+1))
	query.Set("per_page", "80")
	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Pexels request: %w", err)
	}
	request.Header.Set("Authorization", imgItem.APIKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch Pexels photos: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("Pexels API returned %s", response.Status)
	}

	var result pexelsResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, "", fmt.Errorf("failed to decode Pexels response: %w", err)
	}
	if len(result.Photos) == 0 {
		return nil, "", fmt.Errorf("Pexels API returned no photos")
	}

	photo := result.Photos[rand.Intn(len(result.Photos))]
	imageURL := photo.Source.Landscape
	if imageURL == "" {
		imageURL = photo.Source.Original
	}
	if imageURL == "" {
		return nil, "", fmt.Errorf("Pexels photo did not include an image URL")
	}

	img, err := loadPexelsImageFromURL(imageURL)
	if err != nil {
		return nil, "", err
	}
	return img, imageURL, nil
}

func loadPexelsImageFromURL(imageURL string) (image.Image, error) {
	response, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pexels image: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Pexels image returned %s", response.Status)
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Pexels image: %w", err)
	}
	return img, nil
}
