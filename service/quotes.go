package service

import (
	"MorphPrototype/config"
	"MorphPrototype/morphLog"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type QService struct {
	interval time.Duration
	fn       func() error
}

type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

type Quote struct {
	Statement string `json:"statement"`
	Author    string `json:"author"`
	//Year      int    `json:"Year"`
}

// Start starts the service.
func (s *QService) Start() error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.fn(); err != nil {
			return fmt.Errorf("periodic run failed: %w", err)
		}
	}

	return nil
}

// NewService creates a new Service instance with an internafl function.
func StartChangeQuote(interval time.Duration) *QService {
	fmt.Println("Start Interval of", interval)
	return &QService{
		fn:       SetQuote,
		interval: interval,
	}
}

func SetQuote() error {
	// Load the background image
	config.GetConfig()

	cfg := config.GetConfig()
	onQLs := make([]config.TextLibrary, 0)
	for _, ql := range cfg.TextLibraries {
		if ql.Use {
			onQLs = append(onQLs, ql)
		}
	}
	if len(onQLs) < 1 {
		log.Println("Error: No Image choices selected. Select a image source")
		return nil
	}

	exePath, err := os.Executable()
	//quoteLibraries := strings.Split("biblekjv.json,JamesFTquotes.json,markTwain.json,NasrulHazimQuotes.json,patton.json,willRogers.json,callOfDuty.json", ",")

	randomIndex := rand.Intn(len(onQLs))
	qLibrary := onQLs[randomIndex]

	// Get the directory containing the executable
	exeDir := filepath.Dir(exePath)
	appFolder := filepath.Join(exeDir, "static")
	appFile := filepath.Join(appFolder, qLibrary.Location)

	// Read the config file
	quotesRaw, err := os.ReadFile(appFile)
	if err != nil {
		fmt.Println("failed to read config file: %w", err)
	}

	// Unmarshal the JSON data into a slice of Quotes
	var quotes []Quote
	err = json.Unmarshal(quotesRaw, &quotes)
	if err != nil {
		fmt.Println("failed to unmarshal config: %w", err)
	}

	fmt.Println("Quote List:", qLibrary.Name, "Quotes Count", err)
	// Get a random index within the range of quotes.
	if len(quotes) == 0 {
		fmt.Println("No quotes found.")
	}
	// Set random quote
	quote := quotes[rand.Intn(len(quotes))]
	config.UpdateConfigField("currentQuoteStatement", quote.Statement)
	config.UpdateConfigField("currentQuoteAuthor", quote.Author)
	fmt.Println("Quote:", quote.Statement)
	fmt.Println("Author:", quote.Author)

	lEntry := morphLog.LogItem{TimeStamp: time.Now().Format("20060102 15:04:05"),
		Message: "Selected Quote", Level: "INFO", Library: "quotes.go SetQuote()",
		Operation: "Setting Quote", Origin: qLibrary.Location, LocalFile: quote.Statement,
	}
	morphLog.UpdateLogs(lEntry)
	fmt.Println("new quote log entry:", lEntry)

	//service.AddQuote()
	//ChangeBackgroundRoutine()
	ChangeView()
	return nil
}
