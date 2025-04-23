package service

import (
	"Metamorphoun/config"
	"Metamorphoun/morphLog"
	"fmt"
	"image"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/reujab/wallpaper"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procBeep = kernel32.NewProc("Beep")
)

func Beep(frequency, duration int) {
	procBeep.Call(uintptr(frequency), uintptr(duration))
}

func UpdateQuote(caller string) error {
	println("UpdateQuote called from", caller)
	trackImage := false
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}
	currentPic := config.ConfigInstance.PicHistories[0]
	var err error
	usr, err := user.Current()
	wallpaperMain := filepath.Join(usr.HomeDir, ".Metamorphoun")
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
	//Step 2: get image from source (web/local)
	img, err = backgroundSetSource(currentPic)
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
			currentPic, img, err = setRandomQuote(currentPic, img)
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
	//wallpaperFavs := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites")

	sourceExt = filepath.Ext(currentPic.OriginName)
	if sourceExt == "" {
		sourceExt = ".png"
	}
	currentPic.SaveName = filepath.Join(wallpaperMain, "pic0"+sourceExt)
	config.ConfigInstance.PicHistories[0] = currentPic
	// fileLocBase := strings.Split(filepath.Base(currentPic.SaveName), ".")[0]
	// fileLocDir := filepath.Dir(currentPic.SaveName)
	// println(fileLocBase)
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
	//_ = imgType

	// Set the wallpaper
	fmt.Println("Attempting to set wallpaper from path:", fileLoc)
	//fmt.Println("Caller:", caller)
	if config.ConfigInstance.PicUpdateCalled {
		return nil
	}
	err = wallpaper.SetFromFile(fileLoc)
	if err != nil {
		fmt.Println("Failed to set wallpaper:", err)
	} else {
		fmt.Println("Wallpaper set successfully!")
	}
	BeepLowShort()
	return nil
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

// func clearPic(picEmpty bool, currentPic config.PicHistory, imageItem config.Image) (bool, config.PicHistory) {
// 	picEmpty = true
// 	currentPic = config.PicHistory{}
// 	currentPic.PicNum = 0
// 	currentPic.OriginName = ""
// 	currentPic.SaveName = ""
// 	currentPic.ImageItem = config.Image{}
// 	currentPic.ImageItem.Use = imageItem.Use
// 	currentPic.ImageItem.Name = imageItem.Name
// 	currentPic.ImageItem.Title = imageItem.Title
// 	currentPic.ImageItem.Location = imageItem.Location
// 	currentPic.ImageItem.Operation = imageItem.Operation
// 	currentPic.ImageItem.Inherent = imageItem.Inherent
// 	currentPic.Filter = ""
// 	currentPic.Sizing = ""
// 	currentPic.QuoteStatement = ""
// 	currentPic.QuoteAuthor = ""
// 	currentPic.QuoteFont = ""
// 	currentPic.QuoteFontSize = 0
// 	currentPic.QuoteTextColorR = 0
// 	currentPic.QuoteTextColorG = 0
// 	currentPic.QuoteTextColorB = 0
// 	currentPic.QuoteBackgroundColorR = 0
// 	currentPic.QuoteBackgroundColorG = 0
// 	currentPic.QuoteBackgroundColorB = 0
// 	currentPic.QuoteOpacity = 0
// 	currentPic.QuoteTextBoxWidth = 0
// 	currentPic.QuoteTextBoxHeight = 0
// 	currentPic.QuoteTextBoxX = 0
// 	currentPic.QuoteTextBoxY = 0
// 	return picEmpty, currentPic
// }

// func backgroundGenImageItem(currentPic config.PicHistory) (config.PicHistory, error) {
// 	cfg := config.GetConfig()
// 	onImages, shouldReturn, err := getConfigImages(cfg)
// 	if shouldReturn {
// 		return currentPic, err
// 	}
// 	randomIndex := rand.Intn(len(onImages))
// 	imgItem := onImages[randomIndex]
// 	currentPic.ImageItem = imgItem
// 	return currentPic, nil
// }

// func backgroundSetSource(currentPic config.PicHistory) (image.Image, error) {
// 	var img image.Image
// 	var err error
// 	url := currentPic.OriginName

// 	if currentPic.ImageItem.Name == "Bing" {
// 		img, err = loadBingImageFromURL(url)
// 	} else if currentPic.ImageItem.Name == "Flickr" {
// 		img, err = loadFlickrImageFromURL(url)
// 	} else if currentPic.ImageItem.Name == "NASA" {
// 		img, err = loadNASAImageFromURL(url)
// 	} else if currentPic.ImageItem.Name == "UnSplash" {
// 		img, err = loadNASAImageFromURL(url)
// 	} else {
// 		//WallpapersLocal && Favorites
// 		img, err := loadImage(url)
// 		if err != nil {
// 			fmt.Println("failed to fetch image from URL: %w", err)
// 			return nil, err
// 		}
// 		return img, nil

// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	return img, nil
// }

// func filterCurrentPic(currentPic config.PicHistory, img image.Image) (image.Image, error) {
// 	var err error
// 	switch currentPic.Filter {
// 	case "blurSoft":
// 		img, err = BlurItSet(currentPic, img)
// 	case "blurHard":
// 		img, err = BlurItSet(currentPic, img)
// 	case "pixelate":
// 		img, err = PixelateItSet(currentPic, img)
// 	case "oilify":
// 		img, err = OilifyItSet(currentPic, img)
// 	case "wavy":
// 		img, err = PicassoSet(currentPic, img)
// 	case "vortex":
// 		img, err = applyVortexToQuadrantsSet(currentPic, img) //, pullDistance, maxAngle, maxDistance
// 	case "monochrome":
// 		currentPic, img, err = MonochromeItNfo(currentPic, img)
// 	default: //Original
// 		err = nil
// 		//Do Nothing
// 	}
// 	if err != nil {
// 		fmt.Println("Error saving image:", err)
// 		return img, err
// 	}

// 	usr, err := user.Current()
// 	if err != nil {
// 		fmt.Println("failed to get user home directory:", err)
// 	}

// 	currentPicsFolder := filepath.Join(usr.HomeDir, ".Metamorphoun")
// 	fmt.Println(currentPicsFolder)
// 	return img, nil

// }

// func setRandomQuote(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {

// 	// Get screen size w32 only BOOO!!
// 	// screenWidth := w32.GetSystemMetrics(w32.SM_CXSCREEN)
// 	// screenHeight := w32.GetSystemMetrics(w32.SM_CYSCREEN)

// 	// Get the number of displays
// 	screenInfo := getScreenInfo()[0]
// 	screenWidth := screenInfo.Width
// 	screenHeight := screenInfo.Height
// 	//Make Sure a Quote is loaded
// 	if config.ConfigInstance.CurrentQuoteStatement == "" && config.ConfigInstance.CurrentQuoteAuthor == "" {
// 		SetQuote("backgroundChange")
// 	}
// 	fmt.Println("Quote:", config.ConfigInstance.CurrentQuoteStatement)
// 	fmt.Println("Author:", config.ConfigInstance.CurrentQuoteAuthor)

// 	// Create a new context with the image dimensions
// 	dc := gg.NewContextForImage(img)

// 	// Set initial font size
// 	initialFontSize, fontPath, shouldReturn, currentPic, err := getFontInfo(currentPic)
// 	if shouldReturn {
// 		return currentPic, img, err
// 	}
// 	currentPic.QuoteFont = fontPath
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

// 	// fileStep2 := filepath.Join(currentPicsFolder, "file5AQuoted.png")
// 	// saveImg(imgWithQuote, fileStep2)

// 	// Save the resulting image
// 	//img = dc.Image() // SavePNG(outputPath)

// 	// fileStep3 := filepath.Join(currentPicsFolder, "file5DcImage.png")
// 	// saveImg(imgWithQuote, fileStep3)
// 	return currentPic, imgWithQuote, err

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
