package service

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/morphLog"
	"Metamorphoun/zutil"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/kbinani/screenshot"
)

func UpdateQuote(caller string) error {
	println("UpdateQuote called from", caller)
	_, favoriteQuotesErr := ensureFavoriteQuotesFile()
	if favoriteQuotesErr != nil {
		return favoriteQuotesErr
	}
	trackImage := false
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}
	if runtime.GOOS == "darwin" && config.ConfigInstance.DifferentWallpaperPerScreen && screenshot.NumActiveDisplays() > 1 {
		return updateDarwinPerScreenQuotes()
	}
	currentPic := config.ConfigInstance.PicHistories[0]
	var err error
	wallpaperMain := GetFolderPath(enum.PathLoc.Config)
	var img image.Image
	if currentPic.OriginName == "" {
		morphLog.UpdateLogs(morphLog.LogItem{
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
			Message:   "Pic History is empty",
			Level:     "Error",
			Library:   "Service",
			Operation: "BackgroundSet",
			Origin:    "No OriginName",
			LocalFile: "serviceBackgroundSet.go",
		})
		log.Println("Pic History is empty")
		return nil
	}
	//Step 2: get image from source (web/local/PicSum)
	if currentPic.ImageItem.Name == "PicSum" {
		picSumCachePath := filepath.Join(GetFolderPath(enum.PathLoc.Config), "picsumPureCache.png")
		img, err = zutil.LoadImg(picSumCachePath)
		if err != nil {
			fmt.Println("Error loading PicSum cache image:", err)
			return err
		}
	} else {
		img, err = backgroundSetSource(currentPic)
	}
	if img == nil {
		fmt.Println("Image is Empty 1 wallpaper firing random")
		println(err)
		return UpdateQuote("UpdateQuote")
	}

	sourceExt := filepath.Ext(currentPic.OriginName)
	_ = sourceExt
	if trackImage {
		pureImage := filepath.Join(wallpaperMain, "qTrackstep2"+sourceExt)
		saveImg(img, pureImage)
	}
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}

	//Step 3: Stretch if set to fill the screen
	//To Stretch or not to Stretch that is the question
	sizingChoice := currentPic.Sizing
	img, currentPic = handleScaling(img, currentPic, sizingChoice, err)
	if img == nil {
		fmt.Println("Image is Empty 2")
	}

	if trackImage {
		stretchImage := filepath.Join(wallpaperMain, "qTrackstep3"+sourceExt)
		saveImg(img, stretchImage)
	}
	//Handle Favorite with quote
	specialCaseType := "General"
	if currentPic.ImageItem.Name == "Favorites" && config.ConfigInstance.ShowTextOverlay {
		if strings.Contains(currentPic.OriginName, "WithQuotes") {
			specialCaseType = "WithQuotes"
		} else {
			specialCaseType = "WithoutQuotes"
		}
	}
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}

	//Step 4: Apply filters

	if specialCaseType != "WithQuotes" {
		img, err = filterCurrentPic(currentPic, img)
		if img == nil {
			fmt.Println("Image is Empty 3")
		}
		if err != nil {
			fmt.Println("Image is Empty 1 wallpaper firing random")
			return UpdateQuote("UpdateQuote")
		}
		if trackImage {
			filteredImage := filepath.Join(wallpaperMain, "qTrackstep4"+sourceExt)
			saveImg(img, filteredImage)
		}
		if err != nil {
			fmt.Println("Error playing beep sound:", err)
		}
	}
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}

	//Step 5: Handle Quote
	if config.ConfigInstance.ShowTextOverlay {
		if specialCaseType != "WithQuotes" {
			currentPic, img, err = SetRandomQuote(currentPic, img)
			if (err != nil) || img == nil {
				_ = err
				fmt.Println("Image is Empty 1 wallpaper firing random")
				return UpdateQuote("UpdateQuote")
			}
			if trackImage {
				quoteImage := filepath.Join(wallpaperMain, "qTrackstep5"+sourceExt)
				saveImg(img, quoteImage)
			}
		}
	}
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}
	//Step 6: Save the image
	removeAllPic0s()
	_ = err

	sourceExt = filepath.Ext(currentPic.OriginName)
	if sourceExt == "" {
		sourceExt = ".png"
	}

	// Per-screen: each screen gets its own random pic + random quote (fully independent)
	if runtime.GOOS == "linux" && config.ConfigInstance.DifferentWallpaperPerScreen && len(GetScreenInfo()) > 1 && SetPerScreenWallpapers != nil {
		numDisplays := len(GetScreenInfo())
		wallpaperPaths := make([]string, 0, numDisplays)
		seenFingerprints := make(map[string]struct{}, numDisplays)

		// Screen 0 uses the image we already built above (random pic + quote already applied)
		firstPath := filepath.Join(wallpaperMain, fmt.Sprintf("linux-screen-0-q%s", sourceExt))
		saveImageForDisplay(img, firstPath, 0)
		wallpaperPaths = append(wallpaperPaths, firstPath)
		perScreenPics := []config.PicHistory{}
		currentPic.SaveName = firstPath
		perScreenPics = append(perScreenPics, currentPic)
		if fp, fpErr := wallpaperAssetFingerprint(currentPic, img); fpErr == nil {
			seenFingerprints[fp] = struct{}{}
		}

		// Each additional screen gets a completely independent random pic + quote
		for displayIndex := 1; displayIndex < numDisplays; displayIndex++ {
			nextPic, nextImg, nextExt, genErr := generateDistinctWallpaperAsset(seenFingerprints, 12)
			if genErr != nil || nextImg == nil {
				return fmt.Errorf("display %d: independent generation failed: %w", displayIndex, genErr)
			}
			nextPath := filepath.Join(wallpaperMain, fmt.Sprintf("linux-screen-%d-q%s", displayIndex, nextExt))
			saveImageForDisplay(nextImg, nextPath, displayIndex)
			nextPic.SaveName = nextPath
			wallpaperPaths = append(wallpaperPaths, nextPath)
			perScreenPics = append(perScreenPics, nextPic)
		}

		currentPic = attachPerScreenPics(currentPic, perScreenPics)
		if config.ConfigInstance.PicUpdateCalled {
			return nil
		}
		perScreenErr := SetPerScreenWallpapers(wallpaperPaths)
		if perScreenErr != nil {
			fmt.Println("UpdateQuote: per-screen wallpaper failed, falling back to single:", perScreenErr)
		} else {
			if err := config.UpdateConfig(func(cfg *config.Config) error {
				if len(cfg.PicHistories) == 0 {
					cfg.PicHistories = append(cfg.PicHistories, currentPic)
				} else {
					cfg.PicHistories[0] = currentPic
				}
				return nil
			}); err != nil {
				return err
			}
			fmt.Println("UpdateQuote: per-screen wallpapers set successfully!")
			BeepLowShort()
			return nil
		}
	}

	if runtime.GOOS == "windows" && config.ConfigInstance.DifferentWallpaperPerScreen {
		if SetPerScreenWallpapers == nil {
			return fmt.Errorf("UpdateQuote: Windows per-screen mode enabled but backend is not registered")
		}
		numDisplays := screenshot.NumActiveDisplays()
		if numDisplays < 2 {
			numDisplays = 8
		}
		wallpaperPaths := make([]string, 0, numDisplays)
		seenFingerprints := make(map[string]struct{}, numDisplays)

		firstPath := windowsWallpaperPath(wallpaperMain, 0, sourceExt)
		saveImageForDisplay(img, firstPath, 0)
		wallpaperPaths = append(wallpaperPaths, firstPath)
		perScreenPics := []config.PicHistory{}
		currentPic.SaveName = firstPath
		perScreenPics = append(perScreenPics, currentPic)
		if fp, fpErr := wallpaperAssetFingerprint(currentPic, img); fpErr == nil {
			seenFingerprints[fp] = struct{}{}
		}

		for displayIndex := 1; displayIndex < numDisplays; displayIndex++ {
			nextPic, nextImg, nextExt, genErr := generateDistinctWallpaperAsset(seenFingerprints, 12)
			if genErr != nil || nextImg == nil {
				return fmt.Errorf("display %d: independent generation failed: %w", displayIndex, genErr)
			}
			nextPath := windowsWallpaperPath(wallpaperMain, displayIndex, nextExt)
			saveImageForDisplay(nextImg, nextPath, displayIndex)
			nextPic.SaveName = nextPath
			wallpaperPaths = append(wallpaperPaths, nextPath)
			perScreenPics = append(perScreenPics, nextPic)
		}

		currentPic = attachPerScreenPics(currentPic, perScreenPics)
		if config.ConfigInstance.PicUpdateCalled {
			return nil
		}
		perScreenErr := SetPerScreenWallpapers(wallpaperPaths)
		if perScreenErr != nil {
			return fmt.Errorf("UpdateQuote: per-screen wallpaper failed: %w", perScreenErr)
		} else {
			if err := config.UpdateConfig(func(cfg *config.Config) error {
				if len(cfg.PicHistories) == 0 {
					cfg.PicHistories = append(cfg.PicHistories, currentPic)
				} else {
					cfg.PicHistories[0] = currentPic
				}
				return nil
			}); err != nil {
				return err
			}
			fmt.Println("UpdateQuote: per-screen wallpapers set successfully!")
			BeepLowShort()
			return nil
		}
	}

	// Fallback: single wallpaper for all screens
	currentPic.SaveName = retainedWallpaperPath(sourceExt)
	currentPic.PerScreenPics = nil
	fileLoc := currentPic.SaveName

	// Save the resulting image to the bufferPic path
	fmt.Println(currentPic.OriginName)
	if _, err := os.Stat(fileLoc); os.IsExist(err) {
		os.Remove(fileLoc)
	}
	if img == nil {
		fmt.Println("Image is Empty 6")
	}
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}
	saveImg(img, fileLoc)

	// Set the wallpaper
	fmt.Println("Attempting to set wallpaper from path:", fileLoc)
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}
	err = setWallpaper(fileLoc)
	if err != nil {
		fmt.Println("Failed to set wallpaper:", err)
	} else {
		fmt.Println("Wallpaper set successfully!")
	}
	if err := config.UpdateConfig(func(cfg *config.Config) error {
		if len(cfg.PicHistories) == 0 {
			cfg.PicHistories = append(cfg.PicHistories, currentPic)
		} else {
			cfg.PicHistories[0] = currentPic
		}
		return nil
	}); err != nil {
		return err
	}
	BeepLowShort()
	return nil
}

func updateDarwinPerScreenQuotes() error {
	numDisplays := screenshot.NumActiveDisplays()
	perScreenPics := config.ConfigInstance.DarwinPerScreenPicHistories
	if len(perScreenPics) < numDisplays {
		return BackgroundGenerate("UpdateQuotePerScreenDarwin", config.PicHistory{})
	}

	wallpaperMain := GetFolderPath(enum.PathLoc.Config)
	wallpaperPaths := make([]string, 0, numDisplays)
	updatedPics := make([]config.PicHistory, 0, numDisplays)
	usedQuotes := make(map[string]struct{}, numDisplays)

	for displayIndex := 0; displayIndex < numDisplays; displayIndex++ {
		currentPic := perScreenPics[displayIndex]
		currentPic, img, sourceExt, err := buildDistinctDarwinQuoteImage(currentPic, wallpaperMain, usedQuotes)
		if err != nil {
			return err
		}

		fileLoc := currentPic.SaveName
		if fileLoc == "" {
			fileLoc = darwinWallpaperPath(wallpaperMain, displayIndex, sourceExt)
		}
		if _, err := os.Stat(fileLoc); err == nil {
			_ = os.Remove(fileLoc)
		}
		if err := saveImageForDisplay(img, fileLoc, displayIndex); err != nil {
			return err
		}

		currentPic.SaveName = fileLoc
		wallpaperPaths = append(wallpaperPaths, fileLoc)
		updatedPics = append(updatedPics, currentPic)
		if !config.ConfigInstance.MBCMode && config.ConfigInstance.ShowTextOverlay && currentPic.QuoteStatement != "" {
			usedQuotes[currentPic.QuoteStatement] = struct{}{}
		}
	}

	if err := setDarwinWallpapers(wallpaperPaths); err != nil {
		return err
	}

	updatedRoot := attachPerScreenPics(updatedPics[0], updatedPics)
	if err := config.UpdateConfig(func(cfg *config.Config) error {
		cfg.DarwinPerScreenPicHistories = clonePerScreenPics(updatedPics)
		if len(cfg.PicHistories) == 0 {
			cfg.PicHistories = append(cfg.PicHistories, updatedRoot)
		} else {
			cfg.PicHistories[0] = updatedRoot
		}
		return nil
	}); err != nil {
		return err
	}
	BeepLowShort()
	return nil
}

func buildDistinctDarwinQuoteImage(currentPic config.PicHistory, wallpaperMain string, usedQuotes map[string]struct{}) (config.PicHistory, image.Image, string, error) {
	const maxAttempts = 8

	if config.ConfigInstance.MBCMode || !config.ConfigInstance.ShowTextOverlay {
		return buildUpdatedQuoteImage(currentPic, wallpaperMain)
	}

	var fallbackPic config.PicHistory
	var fallbackImg image.Image
	var fallbackExt string
	var fallbackReady bool
	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		candidatePic, candidateImg, candidateExt, err := buildUpdatedQuoteImage(currentPic, wallpaperMain)
		if err != nil {
			lastErr = err
			continue
		}
		if candidatePic.QuoteStatement == "" {
			return candidatePic, candidateImg, candidateExt, nil
		}
		if _, exists := usedQuotes[candidatePic.QuoteStatement]; !exists {
			return candidatePic, candidateImg, candidateExt, nil
		}
		if !fallbackReady {
			fallbackPic = candidatePic
			fallbackImg = candidateImg
			fallbackExt = candidateExt
			fallbackReady = true
		}
	}

	if fallbackReady {
		return fallbackPic, fallbackImg, fallbackExt, nil
	}
	if lastErr != nil {
		return currentPic, nil, "", lastErr
	}
	return buildUpdatedQuoteImage(currentPic, wallpaperMain)
}

func buildUpdatedQuoteImage(currentPic config.PicHistory, wallpaperMain string) (config.PicHistory, image.Image, string, error) {
	var err error
	var img image.Image
	if currentPic.OriginName == "" {
		return currentPic, nil, "", fmt.Errorf("pic history is empty")
	}

	if currentPic.ImageItem.Name == "PicSum" {
		picSumCachePath := filepath.Join(wallpaperMain, "picsumPureCache.png")
		img, err = zutil.LoadImg(picSumCachePath)
		if err != nil {
			return currentPic, nil, "", fmt.Errorf("error loading PicSum cache image: %w", err)
		}
	} else {
		img, err = backgroundSetSource(currentPic)
		if err != nil {
			return currentPic, nil, "", err
		}
	}
	if img == nil {
		return currentPic, nil, "", fmt.Errorf("image is empty after loading source")
	}

	img, currentPic = handleScaling(img, currentPic, currentPic.Sizing, err)
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

	if specialCaseType != "WithQuotes" {
		img, err = filterCurrentPic(currentPic, img)
		if err != nil {
			return currentPic, nil, "", err
		}
		if img == nil {
			return currentPic, nil, "", fmt.Errorf("image is empty after filtering")
		}
	}

	if config.ConfigInstance.ShowTextOverlay && specialCaseType != "WithQuotes" {
		currentPic, img, err = SetRandomQuote(currentPic, img)
		if err != nil {
			return currentPic, nil, "", err
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

func BeepLowShort() {
	switch runtime.GOOS {
	case "windows":
		//frequency := 2000 // Frequency in Hertz
		//duration := 400   // Duration in milliseconds
		//Beep(frequency, duration)
	default:
		//time.Sleep(time.Millisecond * 100) // Small delay between beeps
	}
}
func BeepHighTwice() {
	switch runtime.GOOS {
	case "windows":
		//frequency := 8000 // Frequency in Hertz
		//duration := 800   // Duration in milliseconds
		//Beep(frequency, duration)
		//time.Sleep(time.Millisecond * 100) // Small delay between beeps
		//Beep(frequency, duration)
	default:
		//time.Sleep(time.Millisecond * 100) // Small delay between beeps
	}
}
func GetScreenInfo() []ScreenInfo {
	// Use the platform-provided implementation when registered (e.g. Linux
	// Wayland via hyprctl, or Linux X11 via xrandr). Falls back to the
	// kbinani/screenshot library for Windows, macOS, and any platform that
	// hasn't registered an override.
	if GetDisplayInfo != nil {
		return GetDisplayInfo()
	}
	var screenInfoRange []ScreenInfo
	displayCount := screenshot.NumActiveDisplays()
	for i := 0; i < displayCount; i++ {
		// Get the bounds of the display
		bounds := screenshot.GetDisplayBounds(i)
		width := bounds.Dx()  // Width of the display
		height := bounds.Dy() // Height of the display
		var screen ScreenInfo
		screen.Number = int16(i)
		screen.Width = width
		screen.Height = height
		screenInfoRange = append([]ScreenInfo{screen}, screenInfoRange...)
	}
	return screenInfoRange
}
func GetFontInfo(currentPic config.PicHistory) (float64, string, bool, config.PicHistory, error) {
	// Use configured font size range; fall back to sensible defaults
	minSize := config.ConfigInstance.QuoteFontSizeMin
	maxSize := config.ConfigInstance.QuoteFontSizeMax
	if minSize < 8 {
		minSize = 16
	}
	if maxSize < minSize {
		maxSize = minSize
	}
	// Pick a random size in the range
	initialFontSize := minSize
	if maxSize > minSize {
		initialFontSize = minSize + float64(rand.Intn(int(maxSize-minSize+1)))
	}
	fmt.Printf("Font size range: %.0f–%.0f, picked: %.0f\n", minSize, maxSize, initialFontSize)

	// If graffiti filter is active, use the bundled Permanent Marker font
	if currentPic.Filter == "graffiti" {
		exePath, err := os.Executable()
		if err == nil {
			graffitiFont := filepath.Join(filepath.Dir(exePath), "shared", "static", "fonts", "PermanentMarker-Regular.ttf")
			if _, statErr := os.Stat(graffitiFont); statErr == nil {
				fmt.Println("Graffiti filter active — using Permanent Marker font:", graffitiFont)
				currentPic.QuoteFont = graffitiFont
				currentPic.QuoteStatement = config.ConfigInstance.CurrentQuoteStatement
				currentPic.QuoteAuthor = config.ConfigInstance.CurrentQuoteAuthor
				return initialFontSize, graffitiFont, false, currentPic, nil
			}
			fmt.Println("Graffiti font not found at:", graffitiFont, "— falling back to normal font selection")
		}
	}

	fontPath := GetFolderPath(enum.PathLoc.Fonts) //filepath.Join(GetFolderPath(enum.PathLoc.Fonts), config.ConfigInstance.TextFontFile)
	// List of substrings to exclude
	excludedSubstrings := []string{
		"AmiriQuran.ttf", "EmojiOneColor-SVGinOT.ttf", "KacstBook.ttf", "KacstOffice.ttf", "constani.ttf",
		"MiriamCLM", "NotoKufi", "NotoNaskh", "NotoSans", "NotoSansArabic",
		"Noto", "SegoeIcons", "Marlett.ttf", "opens__", "segmdl2", "symbol.ttf", "webdings", "wingding",
		"Gubbi.ttf", "Navilu.ttf", "DroidSansFallbackFull.ttf", "Mukti.ttf", "Muktibold.ttf",
		"padmaa-Medium-0.5.ttf", "Saab.ttf", "Kalapi.ttf", "utkal.ttf", "Pothana2000.ttf",
		"vemana2000.ttf", "opens___.ttf", "constanb", "SamYak", "LakkiReddy", "Ponnala.ttf",
		"RaviPrakash.ttf", "Raghu", "Lohit", "holomdl2.ttf", "constanz.ttf", "FrankRuehlCLM-Medium.otf",
		"corbeli.ttf", "constan.ttf", "SansSerifCollection.ttf", "corbel.ttf",

		//Mac
		"Braille", "STIXNonUniBolIta.otf", "SF", "Bodoni Ornaments.ttf", "Comic Sans MS Bold.ttf", "Ornaments",
		"NotoSans", "NotoSerif", "STIX", "Webdings", "Symbol", "Dingbats", "Gurmukhi",
	}

	// Get all font files in the specified path

	fontFiles, err := getFontFiles(fontPath)
	if err != nil {
		fmt.Println("Error getting font files", http.StatusInternalServerError)
		return 0, "", true, currentPic, err
	}

	// Filter out fonts that contain any of the excluded substrings
	var validFontFiles []string
	for _, fontFile := range fontFiles {
		exclude := false
		for _, substr := range excludedSubstrings {
			if strings.Contains(strings.ToLower(fontFile), strings.ToLower(substr)) {
				exclude = true
				break
			}
		}
		if !exclude {
			validFontFiles = append(validFontFiles, fontFile)
		}
	}

	if len(validFontFiles) == 0 {
		return 0, "", true, currentPic, fmt.Errorf("no valid fonts found")
	}

	if config.ConfigInstance.QuoteFontRandom {

		// Select a random valid font, retrying if it can't be loaded
		maxAttempts := 10
		for attempt := 0; attempt < maxAttempts; attempt++ {
			fileRnd := rand.Intn(len(validFontFiles))
			fontPath = validFontFiles[fileRnd]

			// Validate the font actually loads before committing
			dc := gg.NewContext(1, 1)
			if err := dc.LoadFontFace(fontPath, 12); err != nil {
				fmt.Printf("Font failed to load (attempt %d): %s — %v\n", attempt+1, fontPath, err)
				// Remove this bad font so we don't pick it again
				validFontFiles = append(validFontFiles[:fileRnd], validFontFiles[fileRnd+1:]...)
				if len(validFontFiles) == 0 {
					return 0, "", true, currentPic, fmt.Errorf("no loadable fonts found")
				}
				continue
			}
			break
		}

		lEntry := morphLog.LogItem{
			TimeStamp: time.Now().Format("20060102 15:04:05"),
			Message:   "Random Font Picked:" + fontPath,
			Level:     "INFO",
			Library:   "AddQuote:Random Font",
			Operation: "Picked random font",
			Origin:    GetFolderPath(enum.PathLoc.Fonts),
			LocalFile: fontPath,
		}
		morphLog.UpdateLogs(lEntry)
		fmt.Println("new log entry:", lEntry)
	} else {
		if zutil.IsInRange(fontPath, validFontFiles) {
			fontPath = filepath.Join(GetFolderPath(enum.PathLoc.Fonts), config.ConfigInstance.TextFontFile)
		} else {
			fontPath = validFontFiles[0]
		}

	}
	fmt.Println("Selected font:", fontPath)
	currentPic.QuoteFont = fontPath
	currentPic.QuoteStatement = config.ConfigInstance.CurrentQuoteStatement
	currentPic.QuoteAuthor = config.ConfigInstance.CurrentQuoteAuthor
	return initialFontSize, fontPath, false, currentPic, nil
}

func CalculateBoxInfo(screenWidth int, screenHeight int, currentPic config.PicHistory, dc *gg.Context) (string, []string, float64, float64, float64, float64, float64, config.PicHistory) {
	// A quadrant is one quarter of the screen
	quadW := float64(screenWidth) / 2.0
	quadH := float64(screenHeight) / 2.0

	quoteText := `"` + currentPic.QuoteStatement + `"`
	authorText := currentPic.QuoteAuthor
	lineSpacing := 1.5
	padX := 20.0
	padY := 20.0

	// The maximum width text can occupy inside the box
	boxInterior := quadW - (padX * 2)

	// MeasureString can underreport by up to ~65% for certain fonts
	// (italic, monospace, decorative). We apply a safety multiplier so that
	// our wrap decisions are based on a pessimistic width estimate.
	const measureFudge = 1.1

	// safeWidth returns the fudged width of a string
	safeWidth := func(s string) float64 {
		w, _ := dc.MeasureString(s)
		return w * measureFudge
	}

	// wrapText splits text into lines that fit within boxInterior
	wrapText := func(text string) []string {
		words := strings.Fields(text)
		var lines []string
		current := ""
		for _, word := range words {
			candidate := current
			if candidate != "" {
				candidate += " "
			}
			candidate += word
			if safeWidth(candidate) > boxInterior && current != "" {
				lines = append(lines, current)
				current = word
			} else {
				current = candidate
			}
		}
		if current != "" {
			lines = append(lines, current)
		}
		return lines
	}

	// Wrap both quote and author
	quoteLines := wrapText(quoteText)
	authorLines := wrapText(authorText)

	// Build the full text block: quote + blank gap + author
	var allLines []string
	allLines = append(allLines, quoteLines...)
	allLines = append(allLines, "") // gap
	allLines = append(allLines, authorLines...)
	fullText := strings.Join(allLines, "\n")

	// Measure the full text block height
	_, measuredH := dc.MeasureMultilineString(fullText, lineSpacing)

	// Find the widest line (using fudged measurement) to size the box snugly
	widest := 0.0
	for _, line := range allLines {
		w := safeWidth(line)
		if w > widest {
			widest = w
		}
	}

	// Box width = widest fudged line + padding, capped to quadrant
	textBoxWidth := math.Min(widest+(padX*2), quadW)
	textBoxHeight := math.Min(measuredH+(padY*2), quadH)

	fmt.Printf("Quadrant: %.0fx%.0f | Interior: %.0f | Box: %.0fx%.0f | Lines: %d\n",
		quadW, quadH, boxInterior, textBoxWidth, textBoxHeight, len(allLines))

	currentPic.QuoteTextBoxWidth = textBoxWidth
	currentPic.QuoteTextBoxHeight = textBoxHeight
	currentPic.QuoteTextBoxX = textBoxWidth
	currentPic.QuoteTextBoxY = textBoxHeight

	var textBlockX, textBlockY float64
	// Return all wrapped lines (quote + gap + author) as the second return value.
	// authorText is kept for backward compat but DrawQuoteText will use allLines.
	return authorText, allLines, measuredH, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic
}

func LocateBox(textBlockX float64, screenWidth int, textBlockY float64, screenHeight int, textBoxWidth float64, textBoxHeight float64) (float64, float64) {
	textBoxLoc := config.ConfigInstance.TextBoxLocation
	validLocs := []string{"topLeft", "topRight", "bottomLeft", "bottomRight", "center"}
	if textBoxLoc == "random" {
		locRnd := rand.Intn(5)
		textBoxLoc = validLocs[locRnd]
	}

	sw := float64(screenWidth)
	sh := float64(screenHeight)
	halfW := sw / 2.0
	halfH := sh / 2.0
	// Center the box within its quadrant
	margin := 20.0

	switch textBoxLoc {
	case "topLeft":
		textBlockX = (halfW-textBoxWidth)/2.0 + margin
		textBlockY = (halfH-textBoxHeight)/2.0 + margin
	case "topRight":
		textBlockX = halfW + (halfW-textBoxWidth)/2.0 - margin
		textBlockY = (halfH-textBoxHeight)/2.0 + margin
	case "bottomLeft":
		textBlockX = (halfW-textBoxWidth)/2.0 + margin
		textBlockY = halfH + (halfH-textBoxHeight)/2.0 - margin
	case "bottomRight":
		textBlockX = halfW + (halfW-textBoxWidth)/2.0 - margin
		textBlockY = halfH + (halfH-textBoxHeight)/2.0 - margin
	case "center":
		textBlockX = (sw - textBoxWidth) / 2
		textBlockY = (sh - textBoxHeight) / 2
	}

	// Clamp to screen bounds
	if textBlockX < 0 {
		textBlockX = margin
	}
	if textBlockY < 0 {
		textBlockY = margin
	}
	if textBlockX+textBoxWidth > sw {
		textBlockX = sw - textBoxWidth - margin
	}
	if textBlockY+textBoxHeight > sh {
		textBlockY = sh - textBoxHeight - margin
	}

	fmt.Printf("Text block position: X=%.2f, Y=%.2f\n", textBlockX, textBlockY)
	fmt.Printf("Text box dimensions: Width=%.2f, Height=%.2f\n", textBoxWidth, textBoxHeight)
	return textBlockX, textBlockY
}

// DrawQuoteText draws all pre-wrapped lines (quote + gap + author) inside the box.
// Lines are drawn one by one so no re-wrapping can occur.
func DrawQuoteText(dc *gg.Context, allLines []string, authorText string, textBlockX, textBlockY float64, textBoxWidth float64) {
	padX := 20.0
	padY := 20.0
	lineSpacing := 1.5

	// Compute the exact line step
	_, h1 := dc.MeasureMultilineString("Mg", lineSpacing)
	_, h2 := dc.MeasureMultilineString("Mg\nMg", lineSpacing)
	lineStep := h2 - h1
	if lineStep < 1 {
		lineStep = h1
	}

	// First baseline
	_, ascent := dc.MeasureString("Mg")
	x := textBlockX + padX
	y := textBlockY + padY + ascent

	for _, line := range allLines {
		dc.DrawString(line, x, y)
		y += lineStep
	}
}

func GetBackgroundColor(currentPic config.PicHistory) (uint8, uint8, uint8, bool, config.PicHistory, error) {
	redColorBackground, greenColorBackground, blueColorBackground := uint8(0), uint8(0), uint8(0)
	if config.ConfigInstance.QuoteAppearanceRandom {
		redColorBackground = uint8(rand.Intn(72))
		greenColorBackground = uint8(rand.Intn(64))
		blueColorBackground = uint8(rand.Intn(64))
	} else {
		bgR, bgG, bgB, err := ConvertHexToRGB(config.ConfigInstance.QuoteBackgroundColor)
		if err != nil {
			fmt.Println("Error converting hex color to RGB:", err)
			return 0, 0, 0, true, currentPic, nil
		}
		redColorBackground = bgR
		greenColorBackground = bgG
		blueColorBackground = bgB
	}

	currentPic.QuoteBackgroundColorR = redColorBackground
	currentPic.QuoteBackgroundColorG = greenColorBackground
	currentPic.QuoteBackgroundColorB = blueColorBackground
	return redColorBackground, greenColorBackground, blueColorBackground, false, currentPic, nil
}
func GetOpacityAndSetBoxBackground(currentPic config.PicHistory, dc *gg.Context, redColorBackground uint8, greenColorBackground uint8, blueColorBackground uint8, textBlockX float64, textBlockY float64, textBoxWidth float64, textBoxHeight float64) (bool, config.PicHistory, error) {
	opacity, errO := strconv.ParseUint(config.ConfigInstance.QuoteBackgroundOpacity, 10, 8)
	if opacity < 110 {
		opacity = uint64(110)
	}
	if errO != nil {
		fmt.Println("Error parsing opacity:", errO)
		return true, currentPic, nil
	}
	//Where did this go
	if config.ConfigInstance.QuoteAppearanceRandom {
		opacity = 110 + uint64(rand.Intn(144))
	}
	config.ConfigInstance.QuoteBackgroundOpacity = zutil.AsString(opacity)

	//If Mosaic the background opacity needs to be higher to be visible
	if currentPic.Filter == "mosaic" {
		if opacity < 180 {
			opacity = 180 + uint64(rand.Intn(60))
		}
	}

	currentPic.QuoteOpacity = opacity

	//fmt.Println("opacity", opacity)
	dc.SetColor(color.RGBA{redColorBackground, greenColorBackground, blueColorBackground, uint8(opacity)})
	dc.DrawRoundedRectangle(textBlockX, textBlockY, textBoxWidth, textBoxHeight, 10)
	dc.Fill()
	DrawCyberpunkQuoteBorder(currentPic, dc, textBlockX, textBlockY, textBoxWidth, textBoxHeight)
	return false, currentPic, nil
}

func DrawCyberpunkQuoteBorder(currentPic config.PicHistory, dc *gg.Context, textBlockX, textBlockY, textBoxWidth, textBoxHeight float64) {
	if currentPic.Filter != "cyberpunk" {
		return
	}
	borderCount := 1 + rand.Intn(2)
	colors := []color.RGBA{
		{R: 0, G: 232, B: 255, A: 235},
		{R: 255, G: 35, B: 214, A: 225},
		{R: 255, G: 242, B: 47, A: 210},
		{R: 44, G: 255, B: 144, A: 210},
	}
	for i := 0; i < borderCount; i++ {
		offset := float64(i) * 7
		lineWidth := 2.5 + rand.Float64()*2.5
		dc.SetColor(colors[rand.Intn(len(colors))])
		dc.SetLineWidth(lineWidth)
		dc.DrawRoundedRectangle(textBlockX-offset, textBlockY-offset, textBoxWidth+offset*2, textBoxHeight+offset*2, 10+offset)
		dc.Stroke()
	}
}
func calculateLuminance(r, g, b uint8) float64 {
	return 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
}
func GetTextColor(redColorBackground uint8, greenColorBackground uint8, blueColorBackground uint8, currentPic config.PicHistory, dc *gg.Context) (bool, config.PicHistory, error) {
	var redColorText, greenColorText, blueColorText uint8

	if config.ConfigInstance.QuoteAppearanceRandom {
		luminance := calculateLuminance(redColorBackground, greenColorBackground, blueColorBackground)

		if luminance < 128 {
			// Background is dark, use white text
			redColorText, greenColorText, blueColorText = 255, 255, 255
		} else {
			// Background is light, use black text
			redColorText, greenColorText, blueColorText = 0, 0, 0
		}

		currentPic.QuoteTextColorR = redColorText
		currentPic.QuoteTextColorG = greenColorText
		currentPic.QuoteTextColorB = blueColorText
	} else {
		bgR, bgG, bgB, err := ConvertHexToRGB(config.ConfigInstance.QuoteTextColor)
		if err != nil {
			fmt.Println("Error converting hex color to RGB:", err)
			return true, currentPic, nil
		}
		redColorText = bgR
		greenColorText = bgG
		blueColorText = bgB

		currentPic.QuoteTextColorR = redColorText
		currentPic.QuoteTextColorG = greenColorText
		currentPic.QuoteTextColorB = blueColorText
	}

	fmt.Printf("RGB for text: R-%d, G-%d, B-%d\n", redColorText, greenColorText, blueColorText)

	dc.SetColor(color.RGBA{redColorText, greenColorText, blueColorText, 255})
	return false, currentPic, nil
}

// func GetTextColor(redColorBackground uint8, greenColorBackground uint8, blueColorBackground uint8, currentPic config.PicHistory, dc *gg.Context) (bool, config.PicHistory, error) {
// 	redColorText, greenColorText, blueColorText := uint8(0), uint8(0), uint8(0)
// 	if config.ConfigInstance.QuoteAppearanceRandom {
// 		prominentBGColor := "red"
// 		if redColorBackground >= greenColorBackground && redColorBackground >= blueColorBackground {
// 			prominentBGColor = "red"
// 		}
// 		if greenColorBackground >= redColorBackground && greenColorBackground >= blueColorBackground {
// 			prominentBGColor = "green"
// 		}
// 		if blueColorBackground >= redColorBackground && blueColorBackground >= greenColorBackground {
// 			prominentBGColor = "blue"
// 		}
// 		otherColorsModifier := uint8(0)
// 		if prominentBGColor == "red" {
// 			otherColorsModifier = (redColorBackground - greenColorBackground) + (redColorBackground - blueColorBackground)
// 		} else {
// 			if prominentBGColor == "green" {
// 				otherColorsModifier = (greenColorBackground - redColorBackground) + (greenColorBackground - blueColorBackground)
// 			} else {
// 				otherColorsModifier = (blueColorBackground - greenColorBackground) + (blueColorBackground - redColorBackground)
// 			}
// 		}
// 		redColorText = uint8(224 + rand.Intn(32))
// 		if prominentBGColor != "red" {
// 			if uint32(redColorText)+uint32(otherColorsModifier) > 255 {
// 				redColorText = uint8(255)
// 			} else {
// 				redColorText += otherColorsModifier
// 			}
// 		}
// 		greenColorText = uint8(224 + rand.Intn(32))
// 		if prominentBGColor != "green" {
// 			if uint32(greenColorText)+uint32(otherColorsModifier) > 255 {
// 				greenColorText = uint8(255)
// 			} else {
// 				greenColorText += otherColorsModifier
// 			}
// 		}
// 		blueColorText = uint8(224 + rand.Intn(32))
// 		if prominentBGColor != "blue" {
// 			if uint32(blueColorText)+uint32(otherColorsModifier) > 255 {
// 				blueColorText = uint8(255)
// 			} else {
// 				blueColorText += otherColorsModifier
// 			}
// 		}

// 		currentPic.QuoteTextColorR = redColorText
// 		currentPic.QuoteTextColorG = greenColorText
// 		currentPic.QuoteTextColorB = blueColorText

// 	} else {
// 		bgR, bgG, bgB, err := ConvertHexToRGB(config.ConfigInstance.QuoteTextColor)
// 		if err != nil {
// 			fmt.Println("Error converting hex color to RGB:", err)
// 			return true, currentPic, nil
// 		}
// 		redColorText = bgR
// 		greenColorText = bgG
// 		blueColorText = bgB

// 		currentPic.QuoteTextColorR = redColorText
// 		currentPic.QuoteTextColorG = greenColorText
// 		currentPic.QuoteTextColorB = blueColorText

// 	}

// 	fmt.Println("RGB for text: R-", redColorText, ",G-", greenColorText, ",B-", blueColorText, "")

// 	dc.SetColor(color.RGBA{redColorText, greenColorText, blueColorText, 255})
// 	return false, currentPic, nil
// }

type QService struct {
	interval time.Duration
	fn       func(string) error
	param    string
}

type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

type Quote struct {
	Statement string `json:"statement"`
	Author    string `json:"author"`
	//Year      int    `json:"Year"`
}

func (qs *QService) Start() error {
	ticker := time.NewTicker(qs.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := qs.fn(qs.param); err != nil {
				return err
			}
		}
	}
}

func StartChangeQuote(interval time.Duration) *QService {
	fmt.Println("Start Interval of", interval)
	return &QService{
		fn:       UpdateQuote,
		interval: interval,
		//param:    param,
	}
}
