package service

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/morphLog"
	"Metamorphoun/shared"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reujab/wallpaper"
)

var SetRandomQuote func(config.PicHistory, image.Image) (config.PicHistory, image.Image, error)

func BackgroundGenerate(caller string, currentPic config.PicHistory) error {
	println("BackgroundGenerate called from", caller)
	if config.ConfigInstance.BackgroundChangeAttempt > 3 {
		log.Println("Too many attempts in", caller)
		return fmt.Errorf("Too many bad attempts")
	}
	config.ConfigInstance.PicUpdateCalled = true
	var img image.Image
	picEmpty := false
	if currentPic.OriginName == "" {
		picEmpty, currentPic = clearPic(picEmpty, currentPic)
	}
	if picEmpty {
		// create a new pic history object
		currentPic.PicNum = 0
		//Step 1: get Images Source
		currentPic, err := backgroundGenImageItem(currentPic)
		if err != nil {
			//Failure to get image item
			println("Failure to get image item..rerun")
			config.ConfigInstance.BackgroundChangeAttempt++
			BackgroundGenerate(caller, currentPic)
			return nil
		}
		//Step 2: get image from source (web/local)
		currentPic, img, err = backgroundGenRandomSource(currentPic)
		if img == nil {
			fmt.Println("Image is Empty 1 wallpaper")
			println(err)
			config.ConfigInstance.BackgroundChangeAttempt++
			return BackgroundGenerate(caller, currentPic)
		}
		sourceExt := filepath.Ext(currentPic.OriginName)
		_ = sourceExt
		//Step 3: Stretch if set to fill the screen
		//To Stretch or not to Stretch that is the question
		sizingChoice := config.ConfigInstance.WallpaperImageSizing
		img, currentPic = handleScaling(img, currentPic, sizingChoice, err)
		if img == nil {
			fmt.Println("Image is Empty 2")
		}
		//Step 4: Apply filters
		//Handle Favorite with quote
		specialCaseType := "General"
		if currentPic.ImageItem.Name == "Favorites" && config.ConfigInstance.ShowTextOverlay {
			if strings.Contains(currentPic.OriginName, "WithQuotes") {
				specialCaseType = "WithQuotes"
			} else {
				specialCaseType = "WithoutQuotes"
			}
		}
		//Step 4: Apply filters
		if specialCaseType != "WithQuotes" {
			currentPic, img, err = picTypeAndFilter(currentPic, img, "")
			if img == nil {
				fmt.Println("Image is Empty 3")
			}
			if err != nil {
				fmt.Println("Error applying filter:", err)
				config.ConfigInstance.BackgroundChangeAttempt++
				return BackgroundGenerate(caller, currentPic)
			}
		}
		//Step 5: Handle Quote
		if config.ConfigInstance.ShowTextOverlay {
			if specialCaseType != "WithQuotes" {
				currentPic, img, err = SetRandomQuote(currentPic, img)
				if (err != nil) || img == nil {
					fmt.Println("Error applying quote:")
					config.ConfigInstance.BackgroundChangeAttempt++
					return BackgroundGenerate(caller, currentPic)
				}
			}
		}
		//Step 6: Save the image
		wallpaperMain := GetFolderPath(enum.PathLoc.Config)

		sourceExt = filepath.Ext(currentPic.OriginName)
		if sourceExt == "" {
			sourceExt = ".png"
		}
		if len(sourceExt) > 5 {
			sourceExt = UnUnsplash(currentPic.OriginName)
		}
		if runtime.GOOS == "darwin" {
			oldFn := config.ConfigInstance.PicHistories[1].SaveName
			err = os.Remove(oldFn)
			if err != nil {
				fmt.Println("Error deleting pic0 file:", err)
			}

			fn := uuid.New()
			currentPic.SaveName = filepath.Join(wallpaperMain, "btrfly"+fn.String()+sourceExt)
		} else {
			currentPic.SaveName = filepath.Join(wallpaperMain, "pic0"+sourceExt)
		}
		config.ConfigInstance.AddPicHistory(currentPic)

		fileLoc := currentPic.SaveName
		removeAllPic0s()
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
	} else if currentPic.ImageItem.Name == "PDChristianArt" {
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
		filters = append(filters, "Dali")
	}
	if config.ConfigInstance.WallpaperFilterMosaic {
		filters = append(filters, "mosaic")
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
	case "monochrome":
		currentPic, img, err = MonochromeItNfo(currentPic, img)
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
		return currentPic, nil
	}

	randomIndex := rand.Intn(len(onQLs))
	qLibrary := onQLs[randomIndex]

	quotesRaw := []byte{}
	err := error(nil)
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
