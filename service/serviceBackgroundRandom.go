package service

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/morphLog"
	"Metamorphoun/shared"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kbinani/screenshot"
	"github.com/reujab/wallpaper"
	"golang.org/x/image/draw"
)

var SetRandomQuote func(config.PicHistory, image.Image) (config.PicHistory, image.Image, error)
var SetPerScreenWallpapers func([]string) error

var ErrBackgroundSourceRetry = errors.New("background source requires reroll")

func BackgroundGenerate(caller string, currentPic config.PicHistory) error {
	println("BackgroundGenerate called from", caller)
	if config.ConfigInstance.BackgroundChangeAttempt > 3 {
		log.Println("Too many attempts in", caller)
		config.ConfigInstance.BackgroundChangeAttempt = 0
		return fmt.Errorf("Too many bad attempts")
	}
	config.ConfigInstance.PicUpdateCalled = true
	var img image.Image
	var err error
	picEmpty := false
	if currentPic.OriginName == "" {
		picEmpty, currentPic = clearPic(picEmpty, currentPic)
	}
	if picEmpty {
		var sourceExt string
		currentPic, img, sourceExt, err = generateRandomWallpaperAsset(currentPic)
		if err != nil {
			if errors.Is(err, ErrBackgroundSourceRetry) {
				fmt.Println("Background source could not provide an image, rerolling:", err)
				config.ConfigInstance.BackgroundChangeAttempt++
				return BackgroundGenerate(caller, config.PicHistory{})
			}
			fmt.Println("Error generating wallpaper:", err)
			config.ConfigInstance.BackgroundChangeAttempt++
			return BackgroundGenerate(caller, currentPic)
		}

		//Step 6: Save the image
		wallpaperMain := GetFolderPath(enum.PathLoc.Config)

		removeAllPic0s()
		if runtime.GOOS == "darwin" && config.ConfigInstance.DifferentWallpaperPerScreen && screenshot.NumActiveDisplays() > 1 {
			err = saveDarwinWallpapersForAllScreens(wallpaperMain, currentPic, img, sourceExt)
			if err != nil {
				fmt.Println("Failed to set individual macOS wallpapers:", err)
				config.ConfigInstance.BackgroundChangeAttempt++
				return BackgroundGenerate(caller, currentPic)
			}
			config.ConfigInstance.PicUpdateCalled = false
			config.ConfigInstance.BackgroundChangeAttempt = 0
			return nil
		}
		if runtime.GOOS == "linux" && config.ConfigInstance.DifferentWallpaperPerScreen && screenshot.NumActiveDisplays() > 1 && SetPerScreenWallpapers != nil {
			err = saveLinuxWallpapersForAllScreens(wallpaperMain, currentPic, img, sourceExt)
			if err != nil {
				fmt.Println("Failed to set individual Linux wallpapers:", err)
			} else {
				config.ConfigInstance.PicUpdateCalled = false
				config.ConfigInstance.BackgroundChangeAttempt = 0
				return nil
			}
		}
		if runtime.GOOS == "windows" && config.ConfigInstance.DifferentWallpaperPerScreen && screenshot.NumActiveDisplays() > 1 && SetPerScreenWallpapers != nil {
			err = saveWindowsWallpapersForAllScreens(wallpaperMain, currentPic, img, sourceExt)
			if err != nil {
				fmt.Println("Failed to set individual Windows wallpapers:", err)
			} else {
				config.ConfigInstance.PicUpdateCalled = false
				config.ConfigInstance.BackgroundChangeAttempt = 0
				return nil
			}
		}

		if runtime.GOOS == "darwin" {
			currentPic.SaveName = filepath.Join(wallpaperMain, "btrfly"+uuid.New().String()+sourceExt)
		} else {
			currentPic.SaveName = filepath.Join(wallpaperMain, "pic0"+sourceExt)
		}
		config.ConfigInstance.AddPicHistory(currentPic)

		fileLoc := ""
		if runtime.GOOS == "windows" {
			numDisplays := screenshot.NumActiveDisplays()
			for i := 0; i < numDisplays; i++ {
				currentPic.SaveName = filepath.Join(wallpaperMain, "pic"+fmt.Sprintf("%d", i)+sourceExt)
				fileLoc = currentPic.SaveName
				// Save the resulting image to the bufferPic path
				fmt.Println(currentPic.OriginName)
				if _, err := os.Stat(fileLoc); os.IsExist(err) {
					os.Remove(fileLoc)
				}
				if img == nil {
					fmt.Println("Image is Empty 6")
				}
				// if currentPic.ImageItem.Name == "PicSum" {
				// 	//Picsum images are not saved in the cache
				// 	saveImage(img, "imgPicSumCache.png")
				// }

				saveImg(img, fileLoc)

			}
		} else {
			fileLoc = currentPic.SaveName
			// Save the resulting image to the bufferPic path
			fmt.Println(currentPic.OriginName)
			if _, err := os.Stat(fileLoc); os.IsExist(err) {
				os.Remove(fileLoc)
			}
			if img == nil {
				fmt.Println("Image is Empty 7")
			}
			saveImg(img, fileLoc)

		}
		//_ = imgType

		// Set the wallpaper
		fmt.Println("Attempting to set wallpaper from path:", fileLoc)
		fmt.Println("Caller:", caller)
		BeepHighTwice()
		err = wallpaper.SetFromFile(fileLoc)
		if err != nil {
			fmt.Println("Failed to set wallpaper:", err)
		} else {
			fmt.Println("Wallpaper set successfully!")
		}

	}
	config.ConfigInstance.PicUpdateCalled = false
	//Step 6: Save the image
	config.ConfigInstance.BackgroundChangeAttempt = 0
	return nil
}

func generateRandomWallpaperAsset(currentPic config.PicHistory) (config.PicHistory, image.Image, string, error) {
	currentPic.PicNum = 0

	var err error
	var img image.Image

	currentPic, err = backgroundGenImageItem(currentPic)
	if err != nil {
		return currentPic, nil, "", err
	}

	currentPic, img, err = backgroundGenRandomSource(currentPic)
	if err != nil {
		return currentPic, img, "", err
	}
	if img == nil {
		return currentPic, nil, "", fmt.Errorf("image is empty after selecting a source")
	}

	img, currentPic = handleScaling(img, currentPic, config.ConfigInstance.WallpaperImageSizing, err)
	if img == nil {
		return currentPic, nil, "", fmt.Errorf("image is empty after scaling")
	}

	specialCaseType := "General"
	if currentPic.ImageItem.Name == "Favorites" && config.ConfigInstance.ShowTextOverlay {
		if strings.Contains(currentPic.OriginName, "WithQuotes") {
			specialCaseType = "WithQuotes"
		} else {
			specialCaseType = "WithoutQuotes"
		}
	}

	if currentPic.ImageItem.Name == "PicSum" {
		currentPicsFolder := GetFolderPath(enum.PathLoc.Config)
		picSumCach := filepath.Join(currentPicsFolder, "imgPicSumCache.png")
		if removeErr := os.Remove(picSumCach); removeErr != nil {
			fmt.Println("Error deleting pic0 file:", removeErr)
		}
		saveImage(img, "imgPicSumCache.png")
	}

	if specialCaseType != "WithQuotes" {
		currentPic, img, err = picTypeAndFilter(currentPic, img, "")
		if err != nil {
			return currentPic, img, "", err
		}
		if img == nil {
			return currentPic, nil, "", fmt.Errorf("image is empty after filtering")
		}
	}

	if config.ConfigInstance.ShowTextOverlay && specialCaseType != "WithQuotes" {
		currentPic, img, err = SetRandomQuote(currentPic, img)
		if err != nil {
			return currentPic, img, "", err
		}
		if img == nil {
			return currentPic, nil, "", fmt.Errorf("image is empty after applying quote")
		}
	}

	sourceExt := filepath.Ext(currentPic.OriginName)
	if sourceExt == "" {
		sourceExt = ".png"
	}
	if len(sourceExt) > 5 {
		sourceExt = UnUnsplash(currentPic.OriginName)
	}

	return currentPic, img, sourceExt, nil
}

func saveDarwinWallpapersForAllScreens(wallpaperMain string, currentPic config.PicHistory, img image.Image, sourceExt string) error {
	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays < 2 {
		return fmt.Errorf("individual wallpapers requested without multiple displays")
	}

	seenFingerprints := make(map[string]struct{}, numDisplays)
	firstFingerprint, err := wallpaperAssetFingerprint(currentPic, img)
	if err != nil {
		return err
	}
	seenFingerprints[firstFingerprint] = struct{}{}

	wallpaperPaths := make([]string, 0, numDisplays)
	firstPath := darwinWallpaperPath(wallpaperMain, 0, sourceExt)
	if err := saveImageForDisplay(img, firstPath, 0); err != nil {
		return err
	}
	fmt.Printf("Display %d assigned wallpaper: %s\n", 0, firstPath)
	currentPic.SaveName = firstPath
	if err := config.ConfigInstance.AddPicHistory(currentPic); err != nil {
		return err
	}
	wallpaperPaths = append(wallpaperPaths, firstPath)

	for displayIndex := 1; displayIndex < numDisplays; displayIndex++ {
		_, nextImg, nextExt, err := generateDistinctWallpaperAsset(seenFingerprints, 8)
		if err != nil {
			return fmt.Errorf("display %d wallpaper generation failed: %w", displayIndex, err)
		}
		nextPath := darwinWallpaperPath(wallpaperMain, displayIndex, nextExt)
		if err := saveImageForDisplay(nextImg, nextPath, displayIndex); err != nil {
			return err
		}
		fmt.Printf("Display %d assigned wallpaper: %s\n", displayIndex, nextPath)
		wallpaperPaths = append(wallpaperPaths, nextPath)
	}

	return setDarwinWallpapers(wallpaperPaths)
}

func saveLinuxWallpapersForAllScreens(wallpaperMain string, currentPic config.PicHistory, img image.Image, sourceExt string) error {
	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays < 2 {
		return fmt.Errorf("individual wallpapers requested without multiple displays")
	}
	if SetPerScreenWallpapers == nil {
		return fmt.Errorf("no Linux per-screen wallpaper backend registered")
	}

	seenFingerprints := make(map[string]struct{}, numDisplays)
	firstFingerprint, err := wallpaperAssetFingerprint(currentPic, img)
	if err != nil {
		return err
	}
	seenFingerprints[firstFingerprint] = struct{}{}

	wallpaperPaths := make([]string, 0, numDisplays)
	firstPath := linuxWallpaperPath(wallpaperMain, 0, sourceExt)
	if err := saveImageForDisplay(img, firstPath, 0); err != nil {
		return err
	}
	fmt.Printf("Display %d assigned wallpaper: %s\n", 0, firstPath)
	currentPic.SaveName = firstPath
	if err := config.ConfigInstance.AddPicHistory(currentPic); err != nil {
		return err
	}
	wallpaperPaths = append(wallpaperPaths, firstPath)

	for displayIndex := 1; displayIndex < numDisplays; displayIndex++ {
		_, nextImg, nextExt, err := generateDistinctWallpaperAsset(seenFingerprints, 8)
		if err != nil {
			return fmt.Errorf("display %d wallpaper generation failed: %w", displayIndex, err)
		}
		nextPath := linuxWallpaperPath(wallpaperMain, displayIndex, nextExt)
		if err := saveImageForDisplay(nextImg, nextPath, displayIndex); err != nil {
			return err
		}
		fmt.Printf("Display %d assigned wallpaper: %s\n", displayIndex, nextPath)
		wallpaperPaths = append(wallpaperPaths, nextPath)
	}

	return SetPerScreenWallpapers(wallpaperPaths)
}

func saveWindowsWallpapersForAllScreens(wallpaperMain string, currentPic config.PicHistory, img image.Image, sourceExt string) error {
	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays < 2 {
		return fmt.Errorf("individual wallpapers requested without multiple displays")
	}
	if SetPerScreenWallpapers == nil {
		return fmt.Errorf("no Windows per-screen wallpaper backend registered")
	}

	seenFingerprints := make(map[string]struct{}, numDisplays)
	firstFingerprint, err := wallpaperAssetFingerprint(currentPic, img)
	if err != nil {
		return err
	}
	seenFingerprints[firstFingerprint] = struct{}{}

	wallpaperPaths := make([]string, 0, numDisplays)
	firstPath := windowsWallpaperPath(wallpaperMain, 0, sourceExt)
	if err := saveImageForDisplay(img, firstPath, 0); err != nil {
		return err
	}
	fmt.Printf("Display %d assigned wallpaper: %s\n", 0, firstPath)
	currentPic.SaveName = firstPath
	if err := config.ConfigInstance.AddPicHistory(currentPic); err != nil {
		return err
	}
	wallpaperPaths = append(wallpaperPaths, firstPath)

	for displayIndex := 1; displayIndex < numDisplays; displayIndex++ {
		_, nextImg, nextExt, err := generateDistinctWallpaperAsset(seenFingerprints, 8)
		if err != nil {
			return fmt.Errorf("display %d wallpaper generation failed: %w", displayIndex, err)
		}
		nextPath := windowsWallpaperPath(wallpaperMain, displayIndex, nextExt)
		if err := saveImageForDisplay(nextImg, nextPath, displayIndex); err != nil {
			return err
		}
		fmt.Printf("Display %d assigned wallpaper: %s\n", displayIndex, nextPath)
		wallpaperPaths = append(wallpaperPaths, nextPath)
	}

	return SetPerScreenWallpapers(wallpaperPaths)
}

func darwinWallpaperPath(wallpaperMain string, displayIndex int, sourceExt string) string {
	if sourceExt == "" {
		sourceExt = ".png"
	}
	return filepath.Join(wallpaperMain, fmt.Sprintf("btrfly-screen-%d-%s%s", displayIndex, uuid.New().String(), sourceExt))
}

func linuxWallpaperPath(wallpaperMain string, displayIndex int, sourceExt string) string {
	if sourceExt == "" {
		sourceExt = ".png"
	}
	return filepath.Join(wallpaperMain, fmt.Sprintf("linux-screen-%d-%s%s", displayIndex, uuid.New().String(), sourceExt))
}

func windowsWallpaperPath(wallpaperMain string, displayIndex int, sourceExt string) string {
	if sourceExt == "" {
		sourceExt = ".png"
	}
	return filepath.Join(wallpaperMain, fmt.Sprintf("win-screen-%d-%s%s", displayIndex, uuid.New().String(), sourceExt))
}

func saveImageForDisplay(img image.Image, filePath string, displayIndex int) error {
	if img == nil {
		return fmt.Errorf("cannot save empty image for display %d", displayIndex)
	}

	bounds := screenshot.GetDisplayBounds(displayIndex)
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		saveImg(img, filePath)
		return nil
	}

	fittedImg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.CatmullRom.Scale(fittedImg, fittedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	saveImg(fittedImg, filePath)
	return nil
}

func generateDistinctWallpaperAsset(seenFingerprints map[string]struct{}, maxAttempts int) (config.PicHistory, image.Image, string, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		nextPic, nextImg, nextExt, err := generateRandomWallpaperAsset(config.PicHistory{})
		if err != nil {
			lastErr = err
			continue
		}

		fingerprint, err := wallpaperAssetFingerprint(nextPic, nextImg)
		if err != nil {
			lastErr = err
			continue
		}

		if _, exists := seenFingerprints[fingerprint]; exists {
			lastErr = fmt.Errorf("duplicate wallpaper asset detected")
			continue
		}

		seenFingerprints[fingerprint] = struct{}{}
		return nextPic, nextImg, nextExt, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("failed to generate a distinct wallpaper asset")
	}
	return config.PicHistory{}, nil, "", fmt.Errorf("failed to generate distinct wallpaper after %d attempts: %w", maxAttempts, lastErr)
}

func wallpaperAssetFingerprint(pic config.PicHistory, img image.Image) (string, error) {
	if img == nil {
		return "", fmt.Errorf("cannot fingerprint empty wallpaper image")
	}

	var imageBuffer bytes.Buffer
	if err := png.Encode(&imageBuffer, img); err != nil {
		return "", fmt.Errorf("failed to encode wallpaper image fingerprint: %w", err)
	}

	hasher := sha256.New()
	hasher.Write(imageBuffer.Bytes())
	hasher.Write([]byte("\n"))
	hasher.Write([]byte(pic.QuoteStatement))
	hasher.Write([]byte("\n"))
	hasher.Write([]byte(pic.QuoteAuthor))

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func setDarwinWallpapers(wallpaperPaths []string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("individual wallpaper assignment is only supported on macOS")
	}
	if len(wallpaperPaths) == 0 {
		return fmt.Errorf("no wallpaper paths provided")
	}

	escapedPaths := make([]string, 0, len(wallpaperPaths))
	for _, path := range wallpaperPaths {
		escapedPath := strings.ReplaceAll(path, "\\", "\\\\")
		escapedPath = strings.ReplaceAll(escapedPath, "\"", "\\\"")
		escapedPaths = append(escapedPaths, fmt.Sprintf("\"%s\"", escapedPath))
	}

	script := fmt.Sprintf(`set wallpaperPaths to {%s}
tell application "System Events"
	set desktopCount to count of desktops
	repeat with desktopIndex from 1 to desktopCount
		set pathIndex to ((desktopIndex - 1) mod (count of wallpaperPaths)) + 1
		set picture of desktop desktopIndex to item pathIndex of wallpaperPaths
	end repeat
end tell`, strings.Join(escapedPaths, ", "))

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set macOS wallpapers: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

func clearPic(picEmpty bool, currentPic config.PicHistory) (bool, config.PicHistory) {
	picEmpty = true
	currentPic = config.PicHistory{}
	currentPic.PicNum = 0
	currentPic.OriginName = ""
	currentPic.SaveName = ""
	currentPic.ImageItem = config.Image{}
	currentPic.ImageItem.Use = true
	currentPic.ImageItem.Name = ""
	currentPic.ImageItem.Title = ""
	currentPic.ImageItem.Location = ""
	currentPic.ImageItem.Operation = ""
	currentPic.ImageItem.Inherent = false
	currentPic.Filter = ""
	currentPic.Sizing = ""
	currentPic.QuoteStatement = ""
	currentPic.QuoteAuthor = ""
	currentPic.QuoteFont = ""
	currentPic.QuoteFontSize = 0
	currentPic.QuoteTextColorR = 0
	currentPic.QuoteTextColorG = 0
	currentPic.QuoteTextColorB = 0
	currentPic.QuoteBackgroundColorR = 0
	currentPic.QuoteBackgroundColorG = 0
	currentPic.QuoteBackgroundColorB = 0
	currentPic.QuoteOpacity = 0
	currentPic.QuoteTextBoxWidth = 0
	currentPic.QuoteTextBoxHeight = 0
	currentPic.QuoteTextBoxX = 0
	currentPic.QuoteTextBoxY = 0
	return picEmpty, currentPic
}

func backgroundGenImageItem(currentPic config.PicHistory) (config.PicHistory, error) {
	cfg := config.GetConfig()
	onImages, shouldReturn, err := getConfigImages(cfg)
	if shouldReturn {
		return currentPic, err
	}
	randomIndex := rand.Intn(len(onImages))
	imgItem := onImages[randomIndex]
	currentPic.ImageItem = imgItem
	return currentPic, nil
}

func backgroundGenRandomSource(currentPic config.PicHistory) (config.PicHistory, image.Image, error) {
	var img image.Image
	var err error
	var url string
	if currentPic.ImageItem.Name == "Bing" {
		img, url, err = GetBackgroundBing(currentPic.ImageItem)
	} else if currentPic.ImageItem.Name == "Flickr" {
		img, url, err = GetBackgroundFlickr(currentPic.ImageItem)
	} else if currentPic.ImageItem.Name == "NASA" {
		img, url, err = GetBackgroundNASA(currentPic.ImageItem)
	} else if currentPic.ImageItem.Name == "UnSplash" {
		img, url, err = GetBackgroundUnSplash(currentPic.ImageItem)
	} else if currentPic.ImageItem.Name == "PicSum" {
		img, url, err = GetBackgroundPicSum(currentPic.ImageItem)
	} else if currentPic.ImageItem.Name == "ChristianPD" {
		img, url, err = GetStaticImagesFolder(currentPic.ImageItem)
	} else if currentPic.ImageItem.Name == "JudaismPD" {
		img, url, err = GetStaticImagesFolder(currentPic.ImageItem)
	} else {
		//WallpapersLocal && Favorites
		img, url, err = GetBackgroundFolder(currentPic.ImageItem)
	}
	if err != nil {
		fmt.Println(err)
		return currentPic, nil, err
	}
	currentPic.OriginName = url
	return currentPic, img, nil
}

func picTypeAndFilter(currentPic config.PicHistory, img image.Image, filterChoice string) (config.PicHistory, image.Image, error) {
	filters := []string{}
	if config.ConfigInstance.WallpaperFilterOriginal {
		filters = append(filters, "original")
	}
	if config.ConfigInstance.WallpaperFilterBlurSoft {
		filters = append(filters, "blurSoft")
	}
	if config.ConfigInstance.WallpaperFilterBlurHard {
		filters = append(filters, "blurHard")
	}
	if config.ConfigInstance.WallpaperFilterPixelate {
		filters = append(filters, "pixelate")
	}
	if config.ConfigInstance.WallpaperFilterOilify {
		filters = append(filters, "oilify")
	}
	if config.ConfigInstance.WallpaperFilterWavy {
		if !currentPic.ImageItem.AllowDistort {
			filters = append(filters, "oilify")
		} else {
			filters = append(filters, "Dali")
		}
	}
	if config.ConfigInstance.WallpaperFilterMosaic {
		filters = append(filters, "mosaic")
	}
	if config.ConfigInstance.WallpaperFilterJigsawPuzzle {
		filters = append(filters, "jigsawpuzzle")
	}
	if config.ConfigInstance.WallpaperFilterCartoon {
		filters = append(filters, "cartoon")
	}
	if config.ConfigInstance.WallpaperFilterMonochrome {
		filters = append(filters, "monochrome")
	}
	if config.ConfigInstance.WallpaperFilterGraffiti {
		filters = append(filters, "graffiti")
	}
	if config.ConfigInstance.WallpaperFilterVortex {

		if !currentPic.ImageItem.AllowDistort {
			filters = append(filters, "mosaic")
		} else {
			filters = append(filters, "vortex")
		}
	}
	//if Original is on than weight it more
	if config.ConfigInstance.WallpaperFilterOriginal {
		filters = append(filters, "original")
		filters = append(filters, "original")
	}

	// Ensure filters list is not empty; default to "original"
	if len(filters) == 0 {
		filters = append(filters, "original")
	}

	filtersRndNum := rand.Intn(len(filters))
	imageFilter := filters[filtersRndNum]
	currentPic.Filter = imageFilter
	//-------------------------------------------TESTING!!! FORCE FILTER
	//imageFilter = "spiral"
	var err error
	if filterChoice != "" {
		imageFilter = filterChoice
	}
	switch imageFilter {
	case "blurSoft":
		currentPic, img, err = BlurItNfo(currentPic, img, 2.5)
	case "blurHard":
		currentPic, img, err = BlurItNfo(currentPic, img, 7.5)
	case "pixelate":
		currentPic, img, err = PixelateItNfo(currentPic, img, 0)
	case "oilify":
		currentPic, img, err = OilifyItNfo(currentPic, img, 0)
	case "Dali":
		currentPic, img, err = DaliNfo(currentPic, img, 0)
	case "vortex":
		quadrants := []string{"topLeft", "topRight", "bottomLeft", "bottomRight", "center"}
		currentPic, img, err = applyVortexToQuadrantsNfo(currentPic, img, quadrants) //, pullDistance, maxAngle, maxDistance
	case "mosaic":
		img, err = MosaicSet(currentPic, img) //(img image.Image, tileMinSize int, tileMaxSize int)
	case "jigsawpuzzle":
		img, err = JigsawPuzzleSet(currentPic, img)
	case "cartoon":
		img, err = CartoonSet(currentPic, img)
	case "monochrome":
		currentPic, img, err = MonochromeItNfo(currentPic, img)
	case "graffiti":
		currentPic, img, err = GraffitiItNfo(currentPic, img, 0)
	default: //Original
		err = nil
		//Do Nothing
	}
	if err != nil {
		fmt.Println("Error saving image:", err)
		return currentPic, img, err
	}
	currentPicsFolder := GetFolderPath(enum.PathLoc.Config)
	fmt.Println(currentPicsFolder)
	return currentPic, img, nil

}

// func setRandomQuote(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
// 	var err error
// 	fmt.Println("running setRandomQuote")
// 	// Get the number of displays
// 	screenInfo := getScreenInfo()[0]
// 	screenWidth := screenInfo.Width
// 	screenHeight := screenInfo.Height
// 	//Make Sure a Quote is loaded
// 	currentPic, err = GetQuote(currentPic)
// 	if err != nil {
// 		fmt.Println("Error getting quote:", err)
// 		return currentPic, img, err
// 	}
// 	fmt.Println("Quote:", currentPic.QuoteStatement)
// 	fmt.Println("Author:", currentPic.QuoteAuthor)

// 	// Create a new context with the image dimensions
// 	dc := gg.NewContextForImage(img)

// 	// Set initial font size
// 	initialFontSize, fontPath, shouldReturn, currentPic, err := getFontInfo(currentPic)
// 	if shouldReturn {
// 		return currentPic, img, err
// 	}
// 	currentPic.QuoteFont = fontPath
// 	currentPic.QuoteFontSize = initialFontSize
// 	if err := dc.LoadFontFace(fontPath, initialFontSize); err != nil {
// 		fmt.Println("Error loading font:", err)
// 		return currentPic, img, err
// 	}

// 	// Set maximum dimensions for the text box (60% of the quadrant)
// 	authorText, wrappedQuoteText, quoteHeight, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic := calculateBoxInfo(screenWidth, screenHeight, currentPic, dc)

// 	textBlockX, textBlockY = locateBox(textBlockX, screenWidth, textBlockY, screenHeight, textBoxWidth, textBoxHeight)

// 	// Set transparent background for text block
// 	//Make Background color
// 	redColorBackground, greenColorBackground, blueColorBackground, shouldReturn, currentPic, err := getBackgroundColor(currentPic)
// 	if shouldReturn {
// 		return currentPic, img, err
// 	}

// 	shouldReturn, currPic, err := getOpacityAndSetBoxBackground(currentPic, dc, redColorBackground, greenColorBackground, blueColorBackground, textBlockX, textBlockY, textBoxWidth, textBoxHeight)
// 	if shouldReturn {
// 		return currentPic, img, err
// 	}
// 	currentPic = currPic
// 	// Set text color and draw text
// 	//Make Text color
// 	shouldReturn, currPic2, err := getTextColor(redColorBackground, greenColorBackground, blueColorBackground, currentPic, dc)
// 	if shouldReturn {
// 		return currentPic, img, err
// 	}
// 	currentPic = currPic2
// 	//dc.SetColor(color.White)

// 	dc.DrawStringWrapped(wrappedQuoteText, textBlockX+10, textBlockY+30, 0, 0, textBoxWidth-20, 1.5, gg.AlignLeft)

// 	// Calculate a line height buffer between the quote and the author
// 	lineHeight := 48.0                                    // Replace with the actual height of a line of text
// 	authorY := textBlockY + 30 + quoteHeight + lineHeight // Add a buffer between quote and author
// 	dc.DrawString(authorText, textBlockX+10, authorY+30)
// 	// Get the resulting image (THIS IS THE MAGIC OF THE NEW PIC CONTEXT.  Started with dc := gg.NewContextForImage(img) )
// 	imgWithQuote := dc.Image()
// 	return currentPic, imgWithQuote, err

// }

func GetQuote(currentPic config.PicHistory) (config.PicHistory, error) {
	fmt.Println("GetQuote called")
	config.GetConfig()
	cfg := config.GetConfig()
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
		return currentPic, err
	}

	onQLs := make([]config.TextLibrary, 0)
	for _, ql := range cfg.TextLibraries {
		if ql.Use {
			onQLs = append(onQLs, ql)
		}
	}

	favQuoteFolder := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Quotes") //, "quoteFavorites.json"
	if _, err := os.Stat(favQuoteFolder); os.IsNotExist(err) {
		//Ignore
	} else {
		filePath, err := ensureFavoriteQuotesFile()
		if err != nil {
			fmt.Println(err)
			return currentPic, err
		}
		third := len(onQLs) / 3
		if third < 1 {
			third = 1
		}
		// Inject a new record with filePath every 'third' records
		for i := third - 1; i < len(onQLs); i += third + 1 {
			newRec := onQLs[i] // copy the current record
			newRec.Use = true
			newRec.Name = "Favorites"
			newRec.Location = filePath
			newRec.Citation = "Favs"
			newRec.Creators = "User"
			newRec.Info = "Generated On the Fly for User"
			newRec.Inherent = false
			onQLs = append(onQLs[:i+1], append([]config.TextLibrary{newRec}, onQLs[i+1:]...)...)
		}

	}

	if len(onQLs) < 1 {
		log.Println("Error: No Image choices selected. Select a image source")
		return currentPic, nil
	}

	randomIndex := rand.Intn(len(onQLs))
	qLibrary := onQLs[randomIndex]

	quotesRaw := []byte{}
	err = error(nil)
	if qLibrary.Inherent {
		quotesRaw, err = shared.GetStaticFSQuotes(qLibrary.Location)
		if err != nil {
			fmt.Println("failed to get static file:", err)
			return currentPic, err
		}
	} else {
		quotesRaw, err = os.ReadFile(qLibrary.Location)
		if err != nil {
			fmt.Println("failed to read file:", err)
			return currentPic, err
		}
	}

	// // Read the config file
	// quotesRaw, err := os.ReadFile(appFile)
	// if err != nil {
	// 	fmt.Println("failed to read config file: %w", err)
	// }

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
	fmt.Println("--------------------LOG---------------------")
	fmt.Println("Quote:", quotes)
	quote := quotes[rand.Intn(len(quotes))]
	config.UpdateConfigField("currentQuoteStatement", quote.Statement)
	config.UpdateConfigField("currentQuoteAuthor", quote.Author)
	currentPic.QuoteStatement = quote.Statement
	currentPic.QuoteAuthor = quote.Author

	fmt.Println("Quote:", quote.Statement)
	fmt.Println("Author:", quote.Author)

	lEntry := morphLog.LogItem{TimeStamp: time.Now().Format("20060102 15:04:05"),
		Message: "Selected Quote", Level: "INFO", Library: "quotes.go SetQuote()",
		Operation: "Setting Quote", Origin: qLibrary.Location, LocalFile: quote.Statement,
	}
	morphLog.UpdateLogs(lEntry)
	fmt.Println("new quote log entry:", lEntry)

	return currentPic, nil
}
