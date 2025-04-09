package service

import (
	"Metamorphoun/config"
	"Metamorphoun/morphLog"
	"Metamorphoun/zutil"
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
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/reujab/wallpaper"
	"golang.org/x/image/draw"
)

const (
	SPI_SETDESKWALLPAPER = 20
	SPIF_UPDATEINIFILE   = 0x01
)

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

// The new way of doing this
func ChangeView(caller string) error {
	fmt.Println(caller)
	//Make the pic
	var img image.Image
	var url string
	var err error
	var currentPicsFolder string
	isFavWithQuote := false
	cfg := config.GetConfig()
	currentPic := config.PicHistory{}
	var filteredImg image.Image
	filterChoice := ""
	sourceExt := ""
	sizingChoice := ""
	if len(config.ConfigInstance.PicHistories) < 1 {
		config.ConfigInstance.PicHistories = append([]config.PicHistory{currentPic}, config.ConfigInstance.PicHistories...)
	}
	currentPicInPlace := config.ConfigInstance.PicHistories[0]
	sourceExt = filepath.Ext(currentPicInPlace.SaveName)
	//from here start saving data in
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	currentPicsFolder = filepath.Join(usr.HomeDir, ".Metamorphoun")
	currentPic.SaveName = filepath.Join(currentPicsFolder, "pic0"+sourceExt)

	if caller == "quoteUpdate" {
		currentPic = currentPicInPlace
		if strings.Contains(currentPic.OriginName, "Favorites\\Pictures\\WithQuotes") {
			isFavWithQuote = true
		}
		if strings.HasPrefix(currentPic.OriginName, "http") {
			img, err = zutil.LoadImageFromURL(currentPic.OriginName)
		} else {
			img, err = loadImage(currentPic.OriginName)
		}
		if img == nil {
			fmt.Println("Image is Empty 1 quote")
			lEntry := morphLog.LogItem{TimeStamp: time.Now().Format("20060102 15:04:05"),
				Message: "Selected Quote", Level: "ERROR", Library: "service.go ChangeView",
				Operation: "Setting Quote", Origin: "Pic from " + currentPic.ImageItem.Title + " selected image origin: " +
					currentPic.OriginName + " rendered nil.", LocalFile: ""}
			morphLog.UpdateLogs(lEntry)
			return nil
		}

		if err != nil {
			fmt.Println("failed to fetch image from URL: %w", err)
		}
		saveImg(img, filepath.Join(currentPicsFolder, "pic0"+sourceExt))
		sizingChoice = currentPic.Sizing
		filterChoice = currentPic.Filter
	} else {
		var shouldReturn bool
		onImages, shouldReturn, err := getConfigImages(cfg)
		if shouldReturn {
			return err
		}
		randomIndex := rand.Intn(len(onImages))
		imgItem := onImages[randomIndex]
		//Start Configure Image History
		currentPic.PicNum = 0
		currentPic.ImageItem = imgItem

		img, url, shouldReturn, err = getPicFromRandomSource(imgItem, img, url, err)
		if img == nil {
			fmt.Println("Image is Empty 1 wallpaper")
			return nil
		}
		if shouldReturn {
			return err
		}
		if img == nil {
			//Try next time
			fmt.Println("[ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR]")
			fmt.Println(imgItem.Name + " has NO files! Turn it off or add files")
			fmt.Println("[ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR][ERROR]")
			return nil
		}
		sourceExt = filepath.Ext(url)
		if imgItem.Name == "UnSplash" {
			sourceExt = ".jpg"
		}
		fmt.Println(sourceExt)
		currentPic.OriginName = url // Get screen size
		if strings.Contains(currentPic.OriginName, "Favorites\\Pictures\\WithQuotes") {
			isFavWithQuote = true
		}
		sizingChoice = config.ConfigInstance.WallpaperImageSizing
		filterChoice = ""
	}
	img, currentPic = handleScaling(img, currentPic, sizingChoice, err)
	if img == nil {
		fmt.Println("Image is Empty 2")
	}

	if isFavWithQuote {
		filteredImg = img
	} else {
		filteredImg, filterChoice, err = applyFilter(img, filterChoice)
	}
	if filteredImg == nil {
		fmt.Println("Image is Empty 3")
	}

	// fileStep4 := filepath.Join(currentPicsFolder, "file4BFiltered.png")
	// saveImg(filteredImg, fileStep4)
	// favPicFolderWithQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithQuotes")
	// favPicFolderWithoutQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithOutQuotes")
	// favWithQuotes := strings.Contains(currentPic.OriginName,"WithQuotes")
	//favWithoutQuotes := strings.Contains(currentPic.OriginName, "WithOutQuotes")
	if config.ConfigInstance.ShowTextOverlay && !isFavWithQuote {

		if !isFavWithQuote || currentPic.ImageItem.Name != "Favorites" {
			filteredImg, currentPic, err = placeQuote(filteredImg, currentPic)
			if err != nil {
				fmt.Println("Error determining adding font:", err)
				return err
			}
			if filteredImg == nil {
				fmt.Println("Image is Empty 4")
			}
			// fileStep5 := filepath.Join(currentPicsFolder, "file5BQuoted.png")
			// saveImg(filteredImg, fileStep5)

		}
	}
	img = filteredImg
	if img == nil {
		fmt.Println("Image is Empty 5")
	}
	currentPic.Filter = filterChoice
	currentPic.Sizing = config.ConfigInstance.WallpaperImageSizing
	config.ConfigInstance.AddPicHistory(currentPic)
	if sourceExt == "" {
		sourceExt = ".png"
	}
	fileLocBase := strings.Split(filepath.Base(currentPic.SaveName), ".")[0]
	fileLocDir := filepath.Dir(currentPic.SaveName)
	println(fileLocBase)
	fileLoc := filepath.Join(fileLocDir, fileLocBase+sourceExt)

	// Save the resulting image to the bufferPic path
	fmt.Println(currentPic.OriginName)
	if _, err := os.Stat(fileLoc); os.IsExist(err) {
		os.Remove(fileLoc)
	}
	if img == nil {
		fmt.Println("Image is Empty 6")
	}
	saveImg(img, fileLoc)
	//_ = imgType

	// Set the wallpaper
	fmt.Println("Attempting to set wallpaper from path:", fileLoc)
	fmt.Println("Caller:", caller)
	err = wallpaper.SetFromFile(fileLoc)
	if err != nil {
		fmt.Println("Failed to set wallpaper:", err)
	} else {
		fmt.Println("Wallpaper set successfully!")
	}

	return nil
}

func CallMakeView(pastImg int32, isFavorite bool, isFavoriteWithQuote bool) error {
	cfg := config.GetConfig()
	pic := cfg.PicHistories[pastImg]
	picType := ""
	if isFavorite {
		quoteType := "WithOutQuotes"
		if isFavorite {
			if isFavoriteWithQuote {
				quoteType = "WithQuotes"
			}
			usr, err := user.Current()
			if err != nil {
				fmt.Println("failed to get user home directory: %w", err)
			}
			favPicFolder := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", quoteType)
			now := time.Now()
			dt := now.Format("20060102_150405")
			lastDotIndex := strings.LastIndex(pic.SaveName, ".")
			ext := "png"
			if lastDotIndex != -1 {
				ext := pic.SaveName[lastDotIndex+1:]
				fmt.Println(ext)
			}
			fileName := dt + "." + ext
			pic.SaveName = filepath.Join(favPicFolder, fileName)

		}
		if quoteType == "WithOutQuotes" {
			pic.QuoteStatement = ""
			pic.QuoteAuthor = ""
			pic.QuoteFont = ""
			pic.QuoteTextColorR = 0
			pic.QuoteTextColorG = 0
			pic.QuoteTextColorB = 0
			pic.QuoteBackgroundColorR = 0
			pic.QuoteBackgroundColorG = 0
			pic.QuoteBackgroundColorB = 0
			pic.QuoteOpacity = 0
			picType = "WithOutQuotes"
		} else {
			picType = "WithQuotes"
		}
	}
	//if(!isFavoriteWithQuote) pic.Quote
	MakeView(pic, picType)
	return nil
}

// This is not good.  If its with quote it should just copy and rename the file to the WithQuotes dir.
// if it is without quotes it does need to recreate the image without a quote.  that will be more difficult
// but must be done correctly.  Finally this need a separate track I think for its third function
// of rerendering a past image
func MakeView(pic config.PicHistory, picType string) error {
	//Make the pic
	var img image.Image
	var err error
	currentPic := pic
	var filteredImg image.Image
	filterChoice := ""
	sizingChoice := ""
	if strings.HasPrefix(currentPic.OriginName, "http") {
		resp, err := http.Get(currentPic.OriginName)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		// Decode the image
		img, _, err = image.Decode(resp.Body)
		if err != nil {
			return err
		}
	} else { //local
		img, err = loadImage(currentPic.OriginName)
		if err != nil {
			fmt.Println("failed to fetch image from URL: %w", err)
			return err
		}
	}

	//currentPicInPlace := config.ConfigInstance.PicHistories[0]
	// currentPic.OriginName = pic.OriginName
	// currentPic.SaveName = pic.SaveName
	//var shouldReturn bool
	//Start Configure Image History
	currentPic.PicNum = 0
	//currentPic.OriginName = url // Get screen size
	sizingChoice = pic.Sizing
	filterChoice = pic.Filter
	img, currentPic = handleScaling(img, currentPic, sizingChoice, err)
	filteredImg, filterChoice, err = applyFilter(img, filterChoice)

	if picType == "WithQuotes" {
		filteredImg = img
		filteredImg, currentPic, err = placeQuoteExact(filteredImg, currentPic)
		if err != nil {
			fmt.Println("Error determining adding font:", err)
			return err
		}

	} else {
		//No need to place quote
		filteredImg = img
	}

	img = filteredImg
	config.ConfigInstance.AddPicHistory(currentPic)
	fileLoc := currentPic.SaveName

	// Save the resulting image to the bufferPic path
	fmt.Println(currentPic.OriginName)
	saveImg(img, fileLoc)
	//_ = imgType

	// Set the wallpaper
	fmt.Println("Attempting to set wallpaper from path:", fileLoc)
	if _, err := os.Stat(fileLoc); os.IsNotExist(err) {
		fmt.Println("Error: Wallpaper file does not exist at path:", fileLoc)
		return nil
	}

	err = wallpaper.SetFromFile(fileLoc)
	if err != nil {
		fmt.Println("Failed to set wallpaper:", err)
	} else {
		fmt.Println("Wallpaper set successfully!")
	}

	return nil

}

// Choose the scaling choice and scale image
func handleScaling(img image.Image, currentPic config.PicHistory, choice string, err error) (image.Image, config.PicHistory) {

	screenInfo := getScreenInfo()[0]
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

// Select a randome pic from the reandom source
func getPicFromRandomSource(imgItem config.Image, img image.Image, url string, err error) (image.Image, string, bool, error) {
	if imgItem.Name == "Bing" {
		img, url, err = GetBackgroundBing(imgItem)
	} else if imgItem.Name == "Flickr" {
		img, url, err = GetBackgroundFlickr(imgItem)
	} else if imgItem.Name == "NASA" {
		img, url, err = GetBackgroundNASA(imgItem)
	} else if imgItem.Name == "UnSplash" {
		img, url, err = GetBackgroundUnSplash(imgItem)
	} else {
		//WallpapersLocal && Favorites
		img, url, err = GetBackgroundFolder(imgItem)
	}
	if err != nil {
		fmt.Println(err)
		return nil, "", true, err
	}
	return img, url, false, nil
}

// get the imageItems for randomw selection
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
	screenInfo := getScreenInfo()[0]
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
	screenInfo := getScreenInfo()[0]
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

	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}

	currentPicsFolder := filepath.Join(usr.HomeDir, ".Metamorphoun")
	fmt.Println(currentPicsFolder)
	// fileStep2 := filepath.Join(currentPicsFolder, "file2APostScaling.png")
	// saveImg(resizedImg, fileStep2)

	return resizedImg, nil
}

func applyFilter(img image.Image, filterChoice string) (image.Image, string, error) {
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
	if config.ConfigInstance.WallpaperFilterOriginal {
		filters = append(filters, "oilify")
	}
	if config.ConfigInstance.WallpaperFilterWavy {
		filters = append(filters, "wavy")
	}
	if config.ConfigInstance.WallpaperFilterMonochrome {
		filters = append(filters, "monochrome")
	}
	if config.ConfigInstance.WallpaperFilterVortex {
		filters = append(filters, "vortex")
	}
	//if Original is on than weight it more
	if config.ConfigInstance.WallpaperFilterOriginal {
		filters = append(filters, "original")
		filters = append(filters, "original")
	}

	filtersRndNum := rand.Intn(len(filters))
	imageFilter := filters[filtersRndNum]
	//-------------------------------------------TESTING!!! FORCE FILTER
	//imageFilter = "spiral"
	var err error
	if filterChoice != "" {
		imageFilter = filterChoice
	}
	switch imageFilter {
	case "blurSoft":
		img, err = BlurIt(img, 2.5)
	case "blurHard":
		img, err = BlurIt(img, 7.5)
	case "pixelate":
		img, err = PixelateIt(img, 0)
	case "oilify":
		img, err = OilifyIt(img, 0)
	case "wavy":
		img, err = Picasso(img, 0)
	case "vortex":
		quadrants := []string{"topLeft", "topRight", "bottomLeft", "bottomRight", "center"}
		img, err = applyVortexToQuadrants(img, quadrants) //, pullDistance, maxAngle, maxDistance
	case "monochrome":
		img, err = MonochromeIt(img)
	default: //Original
		err = nil
		//Do Nothing
	}
	if err != nil {
		fmt.Println("Error saving image:", err)
		return img, imageFilter, err
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}

	currentPicsFolder := filepath.Join(usr.HomeDir, ".Metamorphoun")
	fmt.Println(currentPicsFolder)
	// fileStep2 := filepath.Join(currentPicsFolder, "file4AFiltered"+imageFilter+".png")
	// saveImg(img, fileStep2)

	return img, imageFilter, nil

}

func saveImg(img image.Image, fileName string) {
	// Save the resulting image to the bufferPic path
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
func getImageType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	return format, nil
}

// getFontFiles returns a slice of all font file paths in the given directory
func getFontFiles(dir string) ([]string, error) {
	var fontFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".tts") || strings.HasSuffix(file.Name(), ".ttf") || strings.HasSuffix(file.Name(), ".otf")) {
			if file.Name() != "random" {
				fontFiles = append(fontFiles, filepath.Join(dir, file.Name()))
			}
		}
	}

	return fontFiles, nil
}
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
