package service

import (
	"Metamorphoun/config"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeFavoriteQuotesRebuildsInvalidFileFromEverySource(t *testing.T) {
	tempDir := t.TempDir()
	favoritesPath := filepath.Join(tempDir, quoteFavoritesFileName)
	if err := os.WriteFile(favoritesPath, []byte{0, 0, 0}, 0644); err != nil {
		t.Fatal(err)
	}

	firstSource := writeQuoteSource(t, tempDir, "first.json", []Quote{
		{Statement: "First source quote", Author: "First author"},
	})
	secondSource := writeQuoteSource(t, tempDir, "second.json", []Quote{
		{Statement: "Second source quote", Author: "Second author"},
	})
	cfg := &config.Config{TextLibraries: []config.TextLibrary{
		{Name: "First", Location: firstSource},
		{Name: "Second", Location: secondSource},
	}}

	if err := initializeFavoriteQuotes(favoritesPath, cfg); err != nil {
		t.Fatal(err)
	}

	favorites, err := loadFavoriteQuotes(favoritesPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(favorites) != 2 {
		t.Fatalf("expected one favorite from each source, got %d", len(favorites))
	}
	if favorites[0].Statement != "First source quote" || favorites[1].Statement != "Second source quote" {
		t.Fatalf("unexpected favorites: %#v", favorites)
	}
}

func writeQuoteSource(t *testing.T, directory string, name string, quotes []Quote) string {
	t.Helper()
	path := filepath.Join(directory, name)
	if err := saveFavoriteQuotes(path, quotes); err != nil {
		t.Fatal(err)
	}
	return path
}
