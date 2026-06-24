package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const metamorphounDirName = ".Metamorphoun"
const quoteFavoritesFileName = "quoteFavorites.json"

var defaultFavoriteQuote = Quote{
	Statement: "I don't like traffic cameras. In fact, I hate them. But that doesn't mean I can break the speed limit and run red lights to get to a New Orleans Saints game",
	Author:    "Senator John Kennedy",
}

func ensureFavoriteQuotesFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	favQuoteFolder := filepath.Join(usr.HomeDir, metamorphounDirName, "Favorites", "Quotes")
	if err := os.MkdirAll(favQuoteFolder, 0700); err != nil {
		return "", fmt.Errorf("failed to create Favorites Quotes directory: %w", err)
	}

	filePath := filepath.Join(favQuoteFolder, quoteFavoritesFileName)
	quotes, err := loadFavoriteQuotes(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		quotes = []Quote{defaultFavoriteQuote}
	}

	if len(quotes) == 0 {
		quotes = []Quote{defaultFavoriteQuote}
	}

	if err := saveFavoriteQuotes(filePath, quotes); err != nil {
		return "", err
	}

	return filePath, nil
}

func loadFavoriteQuotes(filePath string) ([]Quote, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return []Quote{}, nil
	}

	var quotes []Quote
	if err := json.Unmarshal([]byte(trimmed), &quotes); err == nil {
		return quotes, nil
	}

	legacyQuotes := make([]Quote, 0)
	for _, line := range strings.Split(trimmed, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var quote Quote
		if err := json.Unmarshal([]byte(line), &quote); err != nil {
			return nil, fmt.Errorf("failed to parse favorite quotes file: %w", err)
		}
		legacyQuotes = append(legacyQuotes, quote)
	}

	return legacyQuotes, nil
}

func saveFavoriteQuotes(filePath string, quotes []Quote) error {
	data, err := json.Marshal(quotes)
	if err != nil {
		return fmt.Errorf("failed to marshal favorite quotes: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write favorite quotes file: %w", err)
	}

	return nil
}

func StoreQuote(quoteRecord string) {
	filePath, err := ensureFavoriteQuotesFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	var newQuote Quote
	if err := json.Unmarshal([]byte(quoteRecord), &newQuote); err != nil {
		fmt.Println("failed to parse quote record:", err)
		return
	}

	quotes, err := loadFavoriteQuotes(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, quote := range quotes {
		if quote.Statement == newQuote.Statement && quote.Author == newQuote.Author {
			fmt.Println("Quote already exists in favorites, not saving.")
			return
		}
	}

	quotes = append(quotes, newQuote)
	if err := saveFavoriteQuotes(filePath, quotes); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Quote stored successfully in %s\n", filePath)
}

func MakeFavFolders() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
		return
	}

	err = os.MkdirAll(filepath.Join(usr.HomeDir, metamorphounDirName, "Favorites", "Pictures", "WithQuotes"), 0700)
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(filepath.Join(usr.HomeDir, metamorphounDirName, "Favorites", "Pictures", "WithOutQuotes"), 0700)
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
	if _, err := ensureFavoriteQuotesFile(); err != nil {
		fmt.Println(err)
	}
	err = os.MkdirAll(filepath.Join(usr.HomeDir, metamorphounDirName, "Favorites", "Quotes"), 0700)
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
}
