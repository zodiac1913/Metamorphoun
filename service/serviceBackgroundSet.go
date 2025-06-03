package service

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/morphLog"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/reujab/wallpaper"
)

func RecallBackground(caller string, pastImg int32) error {
	cfg := config.GetConfig()
	pic := cfg.PicHistories[pastImg]
	//if(!isFavoriteWithQuote) pic.Quote
	BackgroundSet(caller, pic)
	return nil
}

func BackgroundSet(caller string, currentPic config.PicHistory) error {
	println("BackgroundSet called from", caller)
	config.ConfigInstance.PicUpdateCalled = true
	var err error
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
		return BackgroundGenerate(caller, currentPic)
	}
	sourceExt := filepath.Ext(currentPic.OriginName)
	_ = sourceExt
	//Step 3: Stretch if set to fill the screen
	//To Stretch or not to Stretch that is the question
	sizingChoice := currentPic.Sizing
	img, currentPic = handleScaling(img, currentPic, sizingChoice, err)
	if img == nil {
		fmt.Println("Image is Empty 2")
	}
	//Step 4: Apply filters
	img, err = filterCurrentPic(currentPic, img)
	if img == nil {
		fmt.Println("Image is Empty 3")
	}
	if err != nil {
		fmt.Println("Image is Empty 1 wallpaper firing random")
		return BackgroundGenerate(caller, currentPic)
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
	//Step 5: Handle Quote
	if config.ConfigInstance.ShowTextOverlay && caller != "SystrayFavStoreNQ" {
		if specialCaseType != "WithQuotes" {
			currentPic, img, err = SetQuoteBlock(currentPic, img)
			if (err != nil) || img == nil {
				fmt.Println("Image is Empty 1 wallpaper firing random")
				return BackgroundGenerate(caller, currentPic)
			}
		}
	}
	//Step 6: Save the image
	removeAllPic0s()
	wallpaperMain := GetFolderPath(enum.PathLoc.Config)

	sourceExt = filepath.Ext(currentPic.OriginName)
	if sourceExt == "" {
		sourceExt = ".png"
	}
	if len(sourceExt) > 5 {
		sourceExt = UnUnsplash(currentPic.OriginName)
	}
	if caller == "SystrayFavStoreNQ" {
		now := time.Now()
		dt := now.Format("20060102_150405")
		fileName := filepath.Base(dt + sourceExt)
		fileLoc := filepath.Join(GetFolderPath(enum.PathLoc.FavWithoutQuote), fileName)
		saveImg(img, fileLoc)
		config.ConfigInstance.PicHistories[0] = currentPic
		return nil
	} else {
		currentPic.SaveName = filepath.Join(wallpaperMain, "pic0"+sourceExt)
		config.ConfigInstance.PicHistories[0] = currentPic
		fileLoc := currentPic.SaveName

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
		//Step 6: Save the image
		config.ConfigInstance.PicUpdateCalled = false

		return nil
	}
}

func SetQuoteBlock(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
	config.GetConfig()
	screenInfo := GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height
	_ = screenWidth
	_ = screenHeight

	config.UpdateConfigField("Quote Statement", currentPic.QuoteStatement)
	config.UpdateConfigField("Quote Author", currentPic.QuoteAuthor)

	lEntry := morphLog.LogItem{TimeStamp: time.Now().Format("20060102 15:04:05"),
		Message: "Selected Quote", Level: "INFO", Library: "quotes.go SetQuote()",
		Operation: "Setting Quote", Origin: "Pic Histories", LocalFile: currentPic.QuoteStatement,
	}
	morphLog.UpdateLogs(lEntry)
	fmt.Println("new quote log entry:", lEntry)

	//----------------------
	// Create a new context with the image dimensions
	dc := gg.NewContextForImage(img)

	// Set initial font size
	initialFontSize, fontPath, shouldReturn, currentPic, err := GetFontInfo(currentPic)
	if shouldReturn {
		return currentPic, img, err
	}
	currentPic.QuoteFont = fontPath
	currentPic.QuoteFontSize = initialFontSize
	if err := dc.LoadFontFace(fontPath, initialFontSize); err != nil {
		fmt.Println("Error loading font:", err)
		return currentPic, img, err
	}

	// Set maximum dimensions for the text box (60% of the quadrant)

	authorText := currentPic.QuoteAuthor
	//wrappedQuoteText?
	//quoteHeight:=
	textBoxWidth := currentPic.QuoteTextBoxWidth
	//textBoxHeight := currentPic.QuoteTextBoxHeight
	textBlockX := currentPic.QuoteTextBoxX
	textBlockY := currentPic.QuoteTextBoxY

	// Set transparent background for text block
	//Make Background color
	redColorBackground := currentPic.QuoteBackgroundColorR
	greenColorBackground := currentPic.QuoteBackgroundColorG
	blueColorBackground := currentPic.QuoteBackgroundColorB

	opacity := currentPic.QuoteOpacity
	_ = opacity
	// Set text color and draw text
	//Make Text color
	shouldReturn, currPic2, err := GetTextColor(redColorBackground, greenColorBackground, blueColorBackground, currentPic, dc)
	if shouldReturn {
		return currentPic, img, err
	}
	currentPic = currPic2
	//dc.SetColor(color.White)

	dc.DrawStringWrapped(currentPic.QuoteStatement, textBlockX, textBlockY, 0, 0, textBoxWidth, 1.5, gg.AlignLeft)

	// Calculate a line height buffer between the quote and the author
	lineHeight := 48.0                                                 // Replace with the actual height of a line of text
	authorY := textBlockY + currentPic.QuoteTextBoxHeight + lineHeight // Add a buffer between quote and author
	dc.DrawString(authorText, textBlockX+10, authorY+30)
	// Get the resulting image (THIS IS THE MAGIC OF THE NEW PIC CONTEXT.  Started with dc := gg.NewContextForImage(img) )
	imgWithQuote := dc.Image()
	return currentPic, imgWithQuote, err
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

func backgroundSetSource(currentPic config.PicHistory) (image.Image, error) {
	var img image.Image
	var err error
	url := currentPic.OriginName

	if currentPic.ImageItem.Name == "Bing" {
		img, err = loadBingImageFromURL(url)
	} else if currentPic.ImageItem.Name == "Flickr" {
		img, err = loadFlickrImageFromURL(url)
	} else if currentPic.ImageItem.Name == "NASA" {
		img, err = loadNASAImageFromURL(url)
	} else if currentPic.ImageItem.Name == "UnSplash" {
		img, err = loadNASAImageFromURL(url)
	} else {
		//WallpapersLocal && Favorites
		img, err := loadImage(url)
		if err != nil {
			fmt.Println("failed to fetch image from URL: %w", err)
			return nil, err
		}
		return img, nil

	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return img, nil
}

func filterCurrentPic(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	var err error
	switch currentPic.Filter {
	case "blurSoft":
		img, err = BlurItSet(currentPic, img)
	case "blurHard":
		img, err = BlurItSet(currentPic, img)
	case "pixelate":
		img, err = PixelateItSet(currentPic, img)
	case "oilify":
		img, err = OilifyItSet(currentPic, img)
	case "wavy":
		img, err = PicassoSet(currentPic, img)
	case "vortex":
		img, err = applyVortexToQuadrantsSet(currentPic, img) //, pullDistance, maxAngle, maxDistance
	case "monochrome":
		currentPic, img, err = MonochromeItNfo(currentPic, img)
	default: //Original
		err = nil
		//Do Nothing
	}
	if err != nil {
		fmt.Println("Error saving image:", err)
		return img, err
	}

	currentPicsFolder := GetFolderPath(enum.PathLoc.Config)
	fmt.Println(currentPicsFolder)
	return img, nil

}
