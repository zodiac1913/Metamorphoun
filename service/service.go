package service

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

const (
	SPI_SETDESKWALLPAPER = 20
	SPIF_UPDATEINIFILE   = 0x01
)

var GetFolderPath func(string) string

type PathLocType string

var ChangeLockScreen func(pic config.PicHistory) error

// Service represents a service that runs an internal function periodically.
type Service struct {
	interval time.Duration
	fn       func() error
}

// NewService creates a new Service instance with an internafl function.
// func StartChangeBackground(interval time.Duration) *Service {
// 	fmt.Println("Start Interval of", interval)
// 	return &Service{
// 		fn:       ChangeView,
// 		interval: interval,

// 	}
// }

// Start starts the service.
func (s *Service) Start() error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run the function immediately
	if err := s.fn(); err != nil {
		return fmt.Errorf("initial run failed: %w", err)
	}

	for range ticker.C {
		if err := s.fn(); err != nil {
			return fmt.Errorf("periodic run failed: %w", err)
		}
	}

	return nil
}

func removeAllPic0s() error {
	//Delete all files in the picture folder with pic0*.*
	wallpaperMain := GetFolderPath(enum.PathLoc.Config)
	//wallpaperFavs := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites")
	pic0Files, err := filepath.Glob(filepath.Join(wallpaperMain, "pic0*.*"))
	if err != nil {
		fmt.Println("Error finding pic0 files:", err)
	}
	for _, file := range pic0Files {
		err = os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting pic0 file:", err)
		}
	}
	pic1Files, err := filepath.Glob(filepath.Join(wallpaperMain, "btrfly*.*"))
	if err != nil {
		fmt.Println("Error finding btrfly files:", err)
	}
	for _, file := range pic1Files {
		err = os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting btrfly file:", err)
		}
	}
	return nil
}

// Choose the scaling choice and scale image
func handleScaling(img image.Image, currentPic config.PicHistory, choice string, err error) (image.Image, config.PicHistory) {

	if len(GetScreenInfo()) < 1 {
		// config.ConfigInstance.BackgroundChangeAttempt++
		// return BackgroundGenerate("handleScaling", currentPic)
		_ = BackgroundGenerate("handleScaling", currentPic)
		return nil, currentPic
	}
	screenInfo := GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height

	// Create a base image with the screen size
	dc := gg.NewContext(screenWidth, screenHeight)
	imWidth := img.Bounds().Dx()
	imHeight := img.Bounds().Dy()

	// Scale the image if necessary
	if screenHeight > imWidth || screenWidth > imHeight || screenHeight < imWidth || screenWidth < imHeight {
		if choice == "backdrop" {
			currentPic.Sizing = "backdrop"
			img, err = centerOnSmokeyBackdrop(img, *dc)
		} else {
			currentPic.Sizing = "stretch"
			img, err = scaleToScreen(img, *dc)
		}
	}
	// Draw the scaled image onto the context
	dc.DrawImage(img, 0, 0)
	return img, currentPic
}

// get the imageItems for random selection
func getConfigImages(cfg *config.Config) ([]config.Image, bool, error) {
	onImages := make([]config.Image, 0)
	for _, obj := range cfg.Images {
		if obj.Use {
			onImages = append(onImages, obj)
		}
	}
	if len(onImages) < 1 {
		log.Println("Error: No Image choices selected. Select a image source")
		return nil, true, nil
	}
	return onImages, false, nil
}

func centerOnSmokeyBackdrop(img image.Image, dc gg.Context) (image.Image, error) {
	// Get screen size
	screenInfo := GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height

	// Load the smokey background image
	backgroundPath := "static/pics/smokey.jpg" // Path to your smokey background image
	bgImage, err := loadImage(backgroundPath)
	if err != nil {
		fmt.Println("Error determining image type:", err)
		return nil, err
	}

	// Draw the background image stretched to fit the screen size
	bgWidth := bgImage.Bounds().Dx()
	bgHeight := bgImage.Bounds().Dy()
	scaleX := float64(screenWidth) / float64(bgWidth)
	scaleY := float64(screenHeight) / float64(bgHeight)
	dc.Scale(scaleX, scaleY)
	dc.DrawImage(bgImage, 0, 0)
	dc.Scale(1/scaleX, 1/scaleY) // Reset scaling

	// Apply a semi-transparent black overlay to create a smokey effect
	dc.SetColor(color.RGBA{0, 0, 0, 128})
	dc.DrawRectangle(0, 0, float64(screenWidth), float64(screenHeight))
	dc.Fill()

	// Calculate the position to center the smaller image
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	offsetX := (screenWidth - imgWidth) / 2
	offsetY := (screenHeight - imgHeight) / 2

	// Draw the smaller image on the base image
	dc.DrawImage(img, offsetX, offsetY)
	return img, nil
}

func scaleToScreen(img image.Image, dc gg.Context) (image.Image, error) {
	// Get screen size
	screenInfo := GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height

	// Get the dimensions of the image
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	fmt.Println("Image Width:", imgWidth)
	fmt.Println("Image Height:", imgHeight)

	// Determine the scaling factor based on the wallpaperImageSizing configuration
	scaleX := 1.0
	scaleY := 1.0
	if config.ConfigInstance.WallpaperImageSizing == "stretch" {
		scaleX = float64(screenWidth) / float64(imgWidth)
		scaleY = float64(screenHeight) / float64(imgHeight)
	}
	fmt.Println("ScaleX: {}", scaleX)
	fmt.Println("ScaleY: {}", scaleY)

	// Create a new context with the screen size
	//dc := gg.NewContext(screenWidth, screenHeight)

	// Resize the image using the calculated scaling factors
	resizedImg := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	// Draw the resized image onto the context
	dc.DrawImage(resizedImg, 0, 0)

	currentPicsFolder := GetFolderPath(enum.PathLoc.Config)
	fmt.Println(currentPicsFolder)
	return resizedImg, nil
}

func saveImg(img image.Image, fileName string) {
	// Save the resulting image to the bufferPic path
	//Deal with unsplash bs
	ext := filepath.Ext(fileName)
	if len(ext) > 5 {
		ext = UnUnsplash(ext)
	}
	fmt.Println("Image type:", ext)
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		fmt.Println("Error saving melted image:", err)
	}

}

func UnUnsplash(url string) string {
	// Split the string at "?" to separate the base URL and query parameters
	parts := strings.SplitN(url, "?", 2)

	if len(parts) < 2 {
		fmt.Println("Query parameters not found")
		return ".png"
	}

	queryString := parts[1]
	params := strings.Split(queryString, "&")

	for _, param := range params {
		if strings.HasPrefix(param, "fm=") {
			ext := strings.TrimPrefix(param, "fm=")
			if ext != "" {
				return "." + ext
			}
		}
	}

	fmt.Println("Parameter 'fm' not found")
	return ".png"
}

func DeleteFilesInFolder(folderPath string) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the folder itself
		if path == folderPath {
			return nil
		}
		// If the item is a file, delete it
		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				fmt.Println("Error deleting file:", err)
			} else {
				fmt.Println("Deleted file:", path)
			}
		}
		return nil
	})
	return err
}
func loadImage(picPath string) (image.Image, error) {
	im, err := gg.LoadImage(picPath)
	if err != nil {
		return nil, fmt.Errorf("bg-text overlay: Error loading background image: %w", err)
	}
	return im, nil
}

// ConvertHexToRGB converts a hex color code to RGB values
func ConvertHexToRGB(hex string) (uint8, uint8, uint8, error) {
	// Remove the leading '#' if present
	hex = strings.TrimPrefix(hex, "#")

	// Parse the hex string
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color code")
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}

	return uint8(r), uint8(g), uint8(b), nil
}

// getFontFiles returns a slice of all font file paths in the given directory
func getFontFiles(dir string) ([]string, error) {
	var fontFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".ttf") || strings.HasSuffix(file.Name(), ".otf")) {
			fontFiles = append(fontFiles, filepath.Join(dir, file.Name()))
		} else if file.IsDir() {
			subdir := filepath.Join(dir, file.Name())
			subfiles, err := getFontFiles(subdir)
			if err != nil {
				return nil, err
			}
			fontFiles = append(fontFiles, subfiles...)
		}
	}
	return fontFiles, nil
}

// func getFontFiles(dir string) ([]string, error) {
// 	var fontFiles []string
// 	files, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, file := range files {
// 		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".tts") || strings.HasSuffix(file.Name(), ".ttf") || strings.HasSuffix(file.Name(), ".otf")) {
// 			if file.Name() != "random" {
// 				fontFiles = append(fontFiles, filepath.Join(dir, file.Name()))
// 			}
// 		}
// 	}

// 	return fontFiles, nil
// }

func wordWrap(text string, maxWidth float64, dc *gg.Context) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var wrappedText string
	var line string

	for _, word := range words {
		testLine := line + word + " "
		testWidth, _ := dc.MeasureString(testLine)
		if testWidth > maxWidth {
			wrappedText += line + "\n"
			line = word + " "
		} else {
			line = testLine
		}
	}
	wrappedText += line
	return wrappedText
}

// ----------------------Utilities

func DeleteFile(fname string) string {
	dErr := os.Remove(fname)
	if dErr != nil {
		log.Println("Error: could not delete file:", fname, "Error:", dErr)
		return ""
	}
	return fname
}
func getFilePaths(directory string) ([]string, error) {
	filePaths := make([]string, 0)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking directory: %v", err)
			return err
		}
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return filePaths, nil
}

// Use platform-specific commands:   .

// Windows:

// Use exec.Command("explorer", folderPath) to open the folder in Windows Explorer.
// macOS:

// Use exec.Command("open", folderPath) to open the folder in Finder.
// Linux:

// Use exec.Command("xdg-open", folderPath) (or other file managers like nautilus, dolphin, etc.) to open the folder in the default file manager.

func OpenFolder(title string, path string) error {
	var cmd *exec.Cmd
	exec.Command("explorer", path)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error opening folder:", err)
	}
	return nil
}

func PicEncode(w io.Writer, m image.Image) error {
	err := png.Encode(w, m)
	if err != nil {
		fmt.Println("PngEncode Error")
		fmt.Println(err)
		err = gif.Encode(w, m, nil)
		if err != nil {
			fmt.Println("GifEncode Error")
			fmt.Println(err)
		}
	}
	return nil
}
func saveImage(img image.Image, fileName string) {
	currentPicsFolder := GetFolderPath(enum.PathLoc.Config)

	fileIn := filepath.Join(currentPicsFolder, fileName)
	saveImg(img, fileIn)
}

type screenInfo struct {
	Number int16 `json:"number"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}
