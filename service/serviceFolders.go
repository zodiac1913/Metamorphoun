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
//TODO: Apparently I need to set the binary
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
func GetStaticImagesFolder(imgItem config.Image) (image.Image, string, error) {
	var allFilePaths []string
	normalizedPath := normalizePath(imgItem.Location)
	filePaths, err := getAllFilePaths(normalizedPath)
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", normalizedPath, err)
	} else {
		allFilePaths = append(allFilePaths, filePaths...)
	}

	//Filter to just files in the static/images folder
	var allFilePathsImages []string
	for _, filePath := range allFilePaths {
		if strings.Contains(filePath, "images") {
			allFilePathsImages = append(allFilePathsImages, filePath)
		}
	}

	for _, filePath := range allFilePathsImages {

		if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".png") {
			allFilePathsImages = append(allFilePathsImages, filePath)
		}
	}

	if len(allFilePathsImages) < 1 {
		fmt.Println("Error: No pictures found in folder", imgItem.Location, "for", imgItem.Operation)
		return nil, "", nil
	}

	fileRnd := rand.Intn(len(allFilePathsImages))
	pic := allFilePathsImages[fileRnd]

	img, err := loadImage(pic)
	if err != nil {
		fmt.Println("failed to fetch image from URL: %w", err)
		return nil, "", err
	}
	return img, pic, nil
}
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
