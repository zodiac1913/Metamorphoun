package service

import (
	"Metamorphoun/config"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBackgroundUnSplashUsesOfficialWallpaperAPI(t *testing.T) {
	const apiKey = "test-access-key"
	downloadTracked := false

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/photos/random":
			assertUnsplashRequest(t, request, apiKey)
			if request.URL.Query().Get("topics") != "wallpapers" {
				t.Errorf("unexpected topic: %q", request.URL.Query().Get("topics"))
			}
			if request.URL.Query().Get("orientation") != "landscape" {
				t.Errorf("unexpected orientation: %q", request.URL.Query().Get("orientation"))
			}
			if request.URL.Query().Get("content_filter") != "high" {
				t.Errorf("unexpected content filter: %q", request.URL.Query().Get("content_filter"))
			}
			json.NewEncoder(response).Encode(map[string]any{
				"urls":  map[string]string{"full": server.URL + "/image.jpg"},
				"links": map[string]string{"download_location": server.URL + "/download"},
			})
		case "/download":
			assertUnsplashRequest(t, request, apiKey)
			downloadTracked = true
			response.WriteHeader(http.StatusOK)
		case "/image.jpg":
			response.Header().Set("Content-Type", "image/jpeg")
			img := image.NewRGBA(image.Rect(0, 0, 2, 2))
			img.Set(0, 0, color.White)
			if err := jpeg.Encode(response, img, nil); err != nil {
				t.Errorf("failed to encode test image: %v", err)
			}
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	previousURL := unsplashRandomPhotoURL
	unsplashRandomPhotoURL = server.URL + "/photos/random"
	defer func() { unsplashRandomPhotoURL = previousURL }()

	img, imageURL, err := GetBackgroundUnSplash(config.Image{APIKey: apiKey})
	if err != nil {
		t.Fatal(err)
	}
	if imageURL != server.URL+"/image.jpg" {
		t.Fatalf("unexpected image URL: %s", imageURL)
	}
	if img.Bounds() != image.Rect(0, 0, 2, 2) {
		t.Fatalf("unexpected image bounds: %v", img.Bounds())
	}
	if !downloadTracked {
		t.Fatal("expected Unsplash download tracking endpoint to be called")
	}
}

func assertUnsplashRequest(t *testing.T, request *http.Request, apiKey string) {
	t.Helper()
	if request.Header.Get("Authorization") != "Client-ID "+apiKey {
		t.Errorf("unexpected authorization header: %q", request.Header.Get("Authorization"))
	}
	if request.Header.Get("Accept-Version") != "v1" {
		t.Errorf("unexpected API version: %q", request.Header.Get("Accept-Version"))
	}
}
