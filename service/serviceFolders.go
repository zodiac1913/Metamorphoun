package service

import (
	"Metamorphoun/config"
	"fmt"
	"image"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

//              New Way

func GetBackgroundFolder(imgItem config.Image) (image.Image, string, error) {
	var allFilePaths []string
	normalizedPath := normalizePath(imgItem.Location)
	filePaths, err := getAllFilePaths(normalizedPath)
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", normalizedPath, err)
	} else {
		allFilePaths = append(allFilePaths, filePaths...)
	}

	if len(allFilePaths) < 1 {
		fmt.Println("Error: No pictures found in folder", imgItem.Location, "for", imgItem.Operation)
		return nil, "", nil
	}

	fileRnd := rand.Intn(len(allFilePaths))
	pic := allFilePaths[fileRnd]

	img, err := loadImage(pic)
	if err != nil {
		fmt.Println("failed to fetch image from URL: %w", err)
		return nil, "", err
	}
	return img, pic, nil
}

// ----------------------------------------- OLD WAY
// func ChangeBackgroundFolder(imgItem config.Image) error {
// 	var allFilePaths []string
// 	normalizedPath := normalizePath(imgItem.Location)
// 	filePaths, err := getAllFilePaths(normalizedPath)
// 	if err != nil {
// 		fmt.Printf("Error walking the path %v: %v\n", normalizedPath, err)
// 	} else {
// 		allFilePaths = append(allFilePaths, filePaths...)
// 	}

// 	if len(allFilePaths) < 1 {
// 		fmt.Println("Error: No pictures found in folder", imgItem.Location, "for", imgItem.Operation)
// 		return nil
// 	}

// 	fileRnd := rand.Intn(len(allFilePaths))
// 	pic := allFilePaths[fileRnd]

// 	usr, err := user.Current()
// 	if err != nil {
// 		fmt.Println("failed to get user home directory:", err)
// 	}
// 	ffFolder := filepath.Join(usr.HomeDir, ".Metamorphoun", "MyPics")
// 	ext := filepath.Ext(filepath.Base(pic))
// 	bufferFile := filepath.Join(usr.HomeDir, ".Metamorphoun", "MyPics", "MyPicBuffer"+ext)
// 	DeleteFile(bufferFile)
// 	copyFile(pic, bufferFile)

// 	fmt.Println("Setting as Background")

// 	config.UpdateConfigField("sourceCurrentBackgroundName", filepath.Base(bufferFile))
// 	config.UpdateConfigField("sourceCurrentBackgroundFolder", ffFolder)
// 	config.UpdateConfigField("originalCurrentBackgroundName", filepath.Base(bufferFile))
// 	config.UpdateConfigField("originalCurrentBackgroundFolder", ffFolder)
// 	config.UpdateConfigField("currentBackgroundName", "MyPic"+ext)
// 	config.UpdateConfigField("currentBackgroundFolder", ffFolder)

// 	fmt.Println("Copying file to", bufferFile)
// 	picFile := filepath.Join(config.ConfigInstance.CurrentBackgroundFolder,
// 		config.ConfigInstance.CurrentBackgroundName)
// 	DeleteFile(picFile)
// 	copyFile(bufferFile, picFile)

// 	lEntry := morphLog.LogItem{TimeStamp: time.Now().Format("20060102 15:04:05"),
// 		Message: "Changed Background", Level: "INFO", Library: imgItem.Location,
// 		Operation: imgItem.Operation, Origin: pic, LocalFile: filepath.Join(ffFolder, "MyPic"+ext),
// 	}
// 	morphLog.UpdateLogs(lEntry)
// 	fmt.Println("new log entry:", lEntry)
// 	//fmt.Println("Todays Logs:", logs)
// 	ChangeBackgroundRoutine()
// 	return nil
// }

func getAllFilePaths(root string) ([]string, error) {
	var filePaths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), "jpg") ||
				strings.HasSuffix(info.Name(), "png") ||
				strings.HasSuffix(info.Name(), "bmp") ||
				strings.HasSuffix(info.Name(), "gif") {
				filePaths = append(filePaths, path)
			}
		}
		return nil
	})
	return filePaths, err
}

func normalizePath(path string) string {
	convertedPath := strings.ReplaceAll(path, `\`, "/")
	return filepath.Clean(convertedPath)
}
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}
	err = destinationFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

// func normalizePath(path string) string {
// 	// Convert Windows path to Linux-style path by replacing backslashes with slashes
// 	convertedPath := strings.ReplaceAll(path, `\`, "/")
// 	convertedPath1 := strings.ReplaceAll(convertedPath, `\\`, "/")
// 	// Use filepath.Clean to clean up the path and make it platform-independent
// 	return filepath.Clean(convertedPath1)
// }
