package service

import (
	"Metamorphoun/config"
	"Metamorphoun/shared"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const metamorphounDirName = ".Metamorphoun"
const quoteFavoritesFileName = "quoteFavorites.json"

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
	if err := initializeFavoriteQuotes(filePath, config.GetConfig()); err != nil {
		return "", err
	}

	return filePath, nil
}

func initializeFavoriteQuotes(filePath string, cfg *config.Config) error {
	quotes, err := loadFavoriteQuotes(filePath)
	if err == nil && len(quotes) > 0 {
		return nil
	}

	quotes, err = buildInitialFavoriteQuotes(cfg)
	if err != nil {
		return err
	}
	return saveFavoriteQuotes(filePath, quotes)
}

func buildInitialFavoriteQuotes(cfg *config.Config) ([]Quote, error) {
	favorites := make([]Quote, 0, len(cfg.TextLibraries))
	for _, library := range cfg.TextLibraries {
		quotesRaw, err := readQuoteLibrary(library)
		if err != nil {
			fmt.Printf("Skipping quote source %s while creating favorites: %v\n", library.Name, err)
			continue
		}

		var quotes []Quote
		if err := json.Unmarshal(quotesRaw, &quotes); err != nil {
			fmt.Printf("Skipping quote source %s while creating favorites: %v\n", library.Name, err)
			continue
		}
		if len(quotes) > 0 {
			favorites = append(favorites, quotes[rand.Intn(len(quotes))])
		}
	}

	if len(favorites) == 0 {
		return nil, fmt.Errorf("cannot create favorite quotes: no configured quote sources contained valid quotes")
	}
	return favorites, nil
}

func readQuoteLibrary(library config.TextLibrary) ([]byte, error) {
	if library.Inherent {
		return shared.GetStaticFSQuotes(library.Location)
	}
	return os.ReadFile(library.Location)
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

	tempFile, err := os.CreateTemp(filepath.Dir(filePath), quoteFavoritesFileName+".*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary favorite quotes file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write favorite quotes file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close favorite quotes file: %w", err)
	}
	if err := os.Rename(tempPath, filePath); err != nil {
		if removeErr := os.Remove(filePath); removeErr != nil && !os.IsNotExist(removeErr) {
			return fmt.Errorf("failed to replace favorite quotes file: %w", err)
		}
		if retryErr := os.Rename(tempPath, filePath); retryErr != nil {
			return fmt.Errorf("failed to write favorite quotes file: %w", retryErr)
		}
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
