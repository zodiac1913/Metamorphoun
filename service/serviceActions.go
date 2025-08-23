package service

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func StoreQuote(quoteRecord string) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
		return
	}
	favQuoteFolder := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Quotes")
	if _, err := os.Stat(favQuoteFolder); os.IsNotExist(err) {
		fmt.Println("Favorites Quotes folder does not exist, creating it...")
		err = os.MkdirAll(favQuoteFolder, 0700) // Adjust permissions as needed
		if err != nil {
			fmt.Println("failed to create Favorites Quotes directory: %w", err)
			return
		}
	}

	// Create a new file with the current timestamp
	fileName := fmt.Sprintf("quoteFavorites.json")
	filePath := filepath.Join(favQuoteFolder, fileName)
	// Check if file exists and if quoteRecord is already present
	if data, err := os.ReadFile(filePath); err == nil {
		if strings.Contains(string(data), quoteRecord) {
			fmt.Println("Quote already exists in favorites, not saving.")
			return
		}
	}

	// Append the quoteRecord to the file (with a newline for separation)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open quote file:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(quoteRecord + "\n"); err != nil {
		fmt.Println("failed to write quote to file:", err)
		return
	}
	fmt.Printf("Quote stored successfully in %s\n", filePath)
}
func MakeFavFolders() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	//wallpaperFavs := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites")

	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithQuotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithOutQuotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Quotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
}
