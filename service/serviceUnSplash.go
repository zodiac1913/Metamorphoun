package service

import (
	"Metamorphoun/config"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"net/url"
)

var unsplashRandomPhotoURL = "https://api.unsplash.com/photos/random"

type unsplashPhoto struct {
	URLs struct {
		Full string `json:"full"`
		Raw  string `json:"raw"`
	} `json:"urls"`
	Links struct {
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
}

func GetBackgroundUnSplash(imgItem config.Image) (image.Image, string, error) {
	if imgItem.APIKey == "" {
		return nil, "", fmt.Errorf("Unsplash API key is not configured")
	}

	photo, err := getRandomUnsplashWallpaper(imgItem.APIKey)
	if err != nil {
		return nil, "", err
	}
	if err := trackUnsplashDownload(photo.Links.DownloadLocation, imgItem.APIKey); err != nil {
		return nil, "", err
	}

	imageURL := buildUnsplashImageURL(photo.URLs.Raw)
	if imageURL == "" {
		imageURL = photo.URLs.Full
	}
	if imageURL == "" {
		return nil, "", fmt.Errorf("Unsplash photo did not include an image URL")
	}

	img, err := loadUnsplashImageFromURL(imageURL)
	if err != nil {
		return nil, "", err
	}
	return img, imageURL, nil
}

func getRandomUnsplashWallpaper(apiKey string) (unsplashPhoto, error) {
	endpoint, err := url.Parse(unsplashRandomPhotoURL)
	if err != nil {
		return unsplashPhoto{}, err
	}
	query := endpoint.Query()
	query.Set("topics", "wallpapers")
	query.Set("orientation", "landscape")
	query.Set("content_filter", "high")
	endpoint.RawQuery = query.Encode()

	request, err := newUnsplashRequest(endpoint.String(), apiKey)
	if err != nil {
		return unsplashPhoto{}, fmt.Errorf("failed to create Unsplash request: %w", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return unsplashPhoto{}, fmt.Errorf("failed to fetch Unsplash photo: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return unsplashPhoto{}, fmt.Errorf("Unsplash API returned %s", response.Status)
	}

	var photo unsplashPhoto
	if err := json.NewDecoder(response.Body).Decode(&photo); err != nil {
		return unsplashPhoto{}, fmt.Errorf("failed to decode Unsplash response: %w", err)
	}
	return photo, nil
}

func trackUnsplashDownload(downloadLocation string, apiKey string) error {
	if downloadLocation == "" {
		return fmt.Errorf("Unsplash photo did not include a download tracking URL")
	}
	request, err := newUnsplashRequest(downloadLocation, apiKey)
	if err != nil {
		return fmt.Errorf("failed to create Unsplash download request: %w", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to track Unsplash download: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Unsplash download endpoint returned %s", response.Status)
	}
	return nil
}

func newUnsplashRequest(requestURL string, apiKey string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Client-ID "+apiKey)
	request.Header.Set("Accept-Version", "v1")
	return request, nil
}

func buildUnsplashImageURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	imageURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	query := imageURL.Query()
	query.Set("fm", "jpg")
	query.Set("fit", "max")
	query.Set("q", "85")
	query.Set("w", "3840")
	imageURL.RawQuery = query.Encode()
	return imageURL.String()
}

func loadUnsplashImageFromURL(imageURL string) (image.Image, error) {
	response, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Unsplash image: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unsplash image returned %s", response.Status)
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Unsplash image: %w", err)
	}
	return img, nil
}
