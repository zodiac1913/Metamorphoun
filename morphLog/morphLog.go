package morphLog

import (
	"Metamorphoun/enum"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var GetFolderPath func(string) string

type PathLocType string

// LogItem struct definition
type LogItem struct {
	TimeStamp string `json:"timeStamp"`
	Message   string `json:"message"`
	Level     string `json:"level"`
	Library   string `json:"library"`
	Operation string `json:"operation"`
	Origin    string `json:"origin"`
	LocalFile string `json:"localFile"`
}

func GetLogs(logType string) ([]LogItem, error) {
	// Construct the path to the JSON file
	logFilePath := filepath.Join(GetFolderPath(enum.PathLoc.Logs), logType+".json")

	// Read and parse the log file
	logItems, err := getLog(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	return logItems, nil
}

func UpdateLogs(entry LogItem) []LogItem {
	// Construct the path to the JSON file
	fileName := fmt.Sprintf("log%s.json", time.Now().Format("20060102"))
	folderLoc := GetFolderPath(enum.PathLoc.Logs)
	logFilePath := filepath.Join(folderLoc, fileName)

	// Read and parse the log file
	logItems, err := getLog(logFilePath)
	if err != nil {
		fmt.Errorf("failed to read log file: %w", err)
	}

	// Insert the new log entry at the first position
	logItems = append([]LogItem{entry}, logItems...)

	// Marshal the updated log items to JSON
	data, err := json.Marshal(logItems)
	if err != nil {
		fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write the updated JSON data back to the file
	err = os.WriteFile(logFilePath, data, 0644)
	if err != nil {
		fmt.Errorf("failed to write file: %w", err)
	}

	return logItems
}

// getLog reads a JSON file and returns a slice of LogItem structs
func getLog(filename string) ([]LogItem, error) {
	// Read the JSON file
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		// Create a blank array of LogItems
		logItems := []LogItem{}
		// Marshal the blank array to JSON
		data, err := json.Marshal(logItems)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		// Create the file and write the JSON data
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		return logItems, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the JSON data into a slice of LogItem structs
	var logItems []LogItem
	err = json.Unmarshal(data, &logItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return logItems, nil
}
