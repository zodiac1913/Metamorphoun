package service

import (
	"Metamorphoun/config"
	"Metamorphoun/morphLog"
	"Metamorphoun/zutil"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"net/http"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/kbinani/screenshot"
)

type screenInfo struct {
	Number int16 `json:"number"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

func placeQuote(img image.Image, currentPic config.PicHistory) (image.Image, config.PicHistory, error) {

	// Get screen size w32 only BOOO!!
	// screenWidth := w32.GetSystemMetrics(w32.SM_CXSCREEN)
	// screenHeight := w32.GetSystemMetrics(w32.SM_CYSCREEN)

	// Get the number of displays
	screenInfo := getScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height
	//Make Sure a Quote is loaded
	if config.ConfigInstance.CurrentQuoteStatement == "" && config.ConfigInstance.CurrentQuoteAuthor == "" {
		SetQuote("backgroundChange")
	}
	fmt.Println("Quote:", config.ConfigInstance.CurrentQuoteStatement)
	fmt.Println("Author:", config.ConfigInstance.CurrentQuoteAuthor)

	// Create a new context with the image dimensions
	dc := gg.NewContextForImage(img)

	// Set initial font size
	initialFontSize, fontPath, shouldReturn, currentPic, err := getFontInfo(currentPic)
	if shouldReturn {
		return img, currentPic, err
	}
	currentPic.QuoteFont = fontPath
	if err := dc.LoadFontFace(fontPath, initialFontSize); err != nil {
		fmt.Println("Error loading font:", err)
		return img, currentPic, err
	}

	// Set maximum dimensions for the text box (60% of the quadrant)
	authorText, wrappedQuoteText, quoteHeight, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic := calculateBoxInfo(screenWidth, screenHeight, currentPic, dc)

	textBlockX, textBlockY = locateBox(textBlockX, screenWidth, textBlockY, screenHeight, textBoxWidth, textBoxHeight)

	// Set transparent background for text block
	//Make Background color
	redColorBackground, greenColorBackground, blueColorBackground, shouldReturn, currentPic, err := getBackgroundColor(currentPic)
	if shouldReturn {
		return img, currentPic, err
	}

	shouldReturn, currPic, err := getOpacityAndSetBoxBackground(currentPic, dc, redColorBackground, greenColorBackground, blueColorBackground, textBlockX, textBlockY, textBoxWidth, textBoxHeight)
	if shouldReturn {
		return img, currentPic, err
	}
	currentPic = currPic
	// Set text color and draw text
	//Make Text color
	shouldReturn, currPic2, err := getTextColor(redColorBackground, greenColorBackground, blueColorBackground, currentPic, dc)
	if shouldReturn {
		return img, currPic2, err
	}
	currentPic = currPic2
	//dc.SetColor(color.White)

	dc.DrawStringWrapped(wrappedQuoteText, textBlockX+10, textBlockY+30, 0, 0, textBoxWidth-20, 1.5, gg.AlignLeft)

	// Calculate a line height buffer between the quote and the author
	lineHeight := 48.0                                    // Replace with the actual height of a line of text
	authorY := textBlockY + 30 + quoteHeight + lineHeight // Add a buffer between quote and author
	dc.DrawString(authorText, textBlockX+10, authorY+30)
	// Get the resulting image (THIS IS THE MAGIC OF THE NEW PIC CONTEXT.  Started with dc := gg.NewContextForImage(img) )
	imgWithQuote := dc.Image()

	// fileStep2 := filepath.Join(currentPicsFolder, "file5AQuoted.png")
	// saveImg(imgWithQuote, fileStep2)

	// Save the resulting image
	//img = dc.Image() // SavePNG(outputPath)

	// fileStep3 := filepath.Join(currentPicsFolder, "file5DcImage.png")
	// saveImg(imgWithQuote, fileStep3)
	return imgWithQuote, currentPic, nil
}

func placeQuoteExact(img image.Image, currentPic config.PicHistory) (image.Image, config.PicHistory, error) {

	//OK SHit. I need all the created values of the text box not just the
	// colors...size and position are needed

	// Get the number of displays
	//screenInfo := getScreenInfo()[0]
	//screenWidth := screenInfo.Width
	//screenHeight := screenInfo.Height

	// Create a new context with the image dimensions
	dc := gg.NewContextForImage(img)

	// Set font size

	initialFontSize := currentPic.QuoteFontSize
	fontPath := currentPic.QuoteFont
	if err := dc.LoadFontFace(fontPath, initialFontSize); err != nil {
		fmt.Println("Error loading font:", err)
		return img, currentPic, err
	}

	// "QuoteStatement": "If we really think that home is elsewhere and that this life is a “wandering to find home,” why should we not look forward to the arrival?",
	// "QuoteAuthor": "C.S. Lewis",
	// "QuoteFont": "C:\\Windows\\Fonts\\impact.ttf",
	// "QuoteTextColorR": 255,
	// "QuoteTextColorG": 247,
	// "QuoteTextColorB": 255,
	// "QuoteBackgroundColorR": 28,
	// "QuoteBackgroundColorG": 39,
	// "QuoteBackgroundColorB": 14,
	// "QuoteOpacity": 140

	// Set maximum dimensions for the text box (60% of the quadrant)
	authorText := currentPic.QuoteAuthor
	wrappedQuoteText := currentPic.QuoteStatement
	//quoteHeight:=0
	textBoxWidth := currentPic.QuoteTextBoxWidth
	//textBoxHeight:= currentPic.QuoteTextBoxHeight
	textBlockX := currentPic.QuoteTextBoxX
	textBlockY := currentPic.QuoteTextBoxY

	// Set transparent background for text block
	//redColorBackground:= currentPic.QuoteBackgroundColorR
	//greenColorBackground:= currentPic.QuoteBackgroundColorG
	//blueColorBackground:= currentPic.QuoteBackgroundColorB

	//Set text color
	//redColorText:= currentPic.QuoteTextColorR
	//greenColorText:= currentPic.QuoteTextColorG
	//blueColorText:= currentPic.QuoteTextColorB

	//Set Opacity
	//opacity:= currentPic.QuoteOpacity

	dc.DrawStringWrapped(wrappedQuoteText, textBlockX+10, textBlockY+30, 0, 0, textBoxWidth-20, 1.5, gg.AlignLeft)

	// Calculate a line height buffer between the quote and the author
	lineHeight := 48.0                                                               // Replace with the actual height of a line of text
	authorY := currentPic.QuoteTextBoxY + 30 + currentPic.QuoteFontSize + lineHeight // Add a buffer between quote and author
	dc.DrawString(authorText, textBlockX+10, authorY+30)
	// Get the resulting image (THIS IS THE MAGIC OF THE NEW PIC CONTEXT.  Started with dc := gg.NewContextForImage(img) )
	imgWithQuote := dc.Image()

	// fileStep2 := filepath.Join(currentPicsFolder, "file5AQuoted.png")
	// saveImg(imgWithQuote, fileStep2)

	// Save the resulting image
	//img = dc.Image() // SavePNG(outputPath)

	// fileStep3 := filepath.Join(currentPicsFolder, "file5DcImage.png")
	// saveImg(imgWithQuote, fileStep3)
	return imgWithQuote, currentPic, nil
}

func getScreenInfo() []screenInfo {
	var screenInfoRange []screenInfo
	displayCount := screenshot.NumActiveDisplays()
	fmt.Printf("Number of displays: %d\n", displayCount)
	for i := 0; i < displayCount; i++ {
		// Get the bounds of the display
		bounds := screenshot.GetDisplayBounds(i)
		width := bounds.Dx()  // Width of the display
		height := bounds.Dy() // Height of the display
		var screen screenInfo
		screen.Number = int16(i)
		screen.Width = width
		screen.Height = height
		screenInfoRange = append([]screenInfo{screen}, screenInfoRange...)
		fmt.Printf("Display %d: width = %d, height = %d\n", i, width, height)
	}
	return screenInfoRange
}
func getFontInfo(currentPic config.PicHistory) (float64, string, bool, config.PicHistory, error) {
	initialFontSize := 22.0
	fontPath := filepath.Join(config.ConfigInstance.TextFontPath, config.ConfigInstance.TextFontFile)
	// List of substrings to exclude
	excludedSubstrings := []string{
		"AmiriQuran.ttf", "EmojiOneColor-SVGinOT.ttf", "KacstBook.ttf", "KacstOffice.ttf",
		"MiriamCLM-Bold.ttf", "MiriamCLM-Book.ttf", "NotoKufi", "NotoNaskh", "NotoSans", "NotoSansArabic",
		"NotoSansArmenian", "NotoSansGeorgian", "NotoSansHebrew", "NotoSansLao", "NotoSerif",
		"SegoeIcons", "Marlett.ttf", "opens__", "segmdl2", "symbol.ttf", "webdings", "wingding",
	}

	// Get all font files in the specified path
	fontFiles, err := getFontFiles(config.ConfigInstance.TextFontPath)
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

		// Select a random valid font
		fileRnd := rand.Intn(len(validFontFiles))
		fontPath = validFontFiles[fileRnd]
		lEntry := morphLog.LogItem{
			TimeStamp: time.Now().Format("20060102 15:04:05"),
			Message:   "Random Font Picked:" + fontPath,
			Level:     "INFO",
			Library:   "AddQuote:Random Font",
			Operation: "Picked random font",
			Origin:    config.ConfigInstance.TextFontPath,
			LocalFile: fontPath,
		}
		morphLog.UpdateLogs(lEntry)
		fmt.Println("new log entry:", lEntry)
	} else {
		if zutil.IsInRange(fontPath, validFontFiles) {
			fontPath = filepath.Join(config.ConfigInstance.TextFontPath, config.ConfigInstance.TextFontFile)
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

func calculateBoxInfo(screenWidth int, screenHeight int, currentPic config.PicHistory, dc *gg.Context) (string, string, float64, float64, float64, float64, float64, config.PicHistory) {
	maxTextBoxWidth := float64(screenWidth) * 0.4   // 60% of half the screen width
	maxTextBoxHeight := float64(screenHeight) * 0.9 // 60% of half the screen height

	// Split the quote text into lines based on the estimated number of characters per line
	quoteText := `"` + currentPic.QuoteStatement + `"`
	authorText := currentPic.QuoteAuthor

	wrappedQuoteText := wordWrap(quoteText, maxTextBoxWidth, dc)

	// Measure the dimensions of the wrapped quote and author text
	quoteWidth, quoteHeight := dc.MeasureMultilineString(wrappedQuoteText, 2)
	authorWidth, authorHeight := dc.MeasureString(authorText)

	// Calculate the required width and height for the text box
	textBoxWidth := math.Min(math.Max(quoteWidth, authorWidth)+40, maxTextBoxWidth) // Add some padding
	textBoxHeight := math.Min(quoteHeight+authorHeight+60, maxTextBoxHeight)        // Add padding
	currentPic.QuoteTextBoxWidth = textBoxWidth
	currentPic.QuoteTextBoxHeight = textBoxHeight
	currentPic.QuoteTextBoxX = textBoxWidth
	currentPic.QuoteTextBoxY = textBoxHeight
	// Define the position for the text block based on the selected quadrant
	var textBlockX, textBlockY float64
	return authorText, wrappedQuoteText, quoteHeight, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic
}

func locateBox(textBlockX float64, screenWidth int, textBlockY float64, screenHeight int, textBoxWidth float64, textBoxHeight float64) (float64, float64) {
	textBoxLoc := config.ConfigInstance.TextBoxLocation
	validLocs := []string{"topLeft", "topRight", "bottomLeft", "bottomRight", "center"}
	if textBoxLoc == "random" {
		locRnd := rand.Intn(5)
		textBoxLoc = validLocs[locRnd]
	}

	switch textBoxLoc {
	case "topLeft":
		textBlockX = float64(screenWidth) * 0.05
		textBlockY = float64(screenHeight) * 0.1
	case "topRight":
		textBlockX = float64(screenWidth)*0.9 - textBoxWidth
		textBlockY = float64(screenHeight) * 0.1
	case "bottomLeft":
		textBlockX = float64(screenWidth) * 0.1
		textBlockY = float64(screenHeight)*0.8 - textBoxHeight
	case "bottomRight":
		textBlockX = float64(screenWidth)*0.9 - textBoxWidth
		textBlockY = float64(screenHeight)*0.8 - textBoxHeight
	case "center":
		textBlockX = (float64(screenWidth) - textBoxWidth) / 2
		textBlockY = (float64(screenHeight) - textBoxHeight) / 2
	}

	// Debug prints to verify the calculated positions
	fmt.Printf("Text block position: X=%.2f, Y=%.2f\n", textBlockX, textBlockY)
	fmt.Printf("Text box dimensions: Width=%.2f, Height=%.2f\n", textBoxWidth, textBoxHeight)
	return textBlockX, textBlockY
}

func getBackgroundColor(currentPic config.PicHistory) (uint8, uint8, uint8, bool, config.PicHistory, error) {
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
func getOpacityAndSetBoxBackground(currentPic config.PicHistory, dc *gg.Context, redColorBackground uint8, greenColorBackground uint8, blueColorBackground uint8, textBlockX float64, textBlockY float64, textBoxWidth float64, textBoxHeight float64) (bool, config.PicHistory, error) {
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
	currentPic.QuoteOpacity = opacity

	//fmt.Println("opacity", opacity)
	dc.SetColor(color.RGBA{redColorBackground, greenColorBackground, blueColorBackground, uint8(opacity)})
	dc.DrawRoundedRectangle(textBlockX, textBlockY, textBoxWidth+20, textBoxHeight+50, 10)
	dc.Fill()
	return false, currentPic, nil
}
func getTextColor(redColorBackground uint8, greenColorBackground uint8, blueColorBackground uint8, currentPic config.PicHistory, dc *gg.Context) (bool, config.PicHistory, error) {
	redColorText, greenColorText, blueColorText := uint8(0), uint8(0), uint8(0)
	if config.ConfigInstance.QuoteAppearanceRandom {
		prominentBGColor := "red"
		if redColorBackground >= greenColorBackground && redColorBackground >= blueColorBackground {
			prominentBGColor = "red"
		}
		if greenColorBackground >= redColorBackground && greenColorBackground >= blueColorBackground {
			prominentBGColor = "green"
		}
		if blueColorBackground >= redColorBackground && blueColorBackground >= greenColorBackground {
			prominentBGColor = "blue"
		}
		otherColorsModifier := uint8(0)
		if prominentBGColor == "red" {
			otherColorsModifier = (redColorBackground - greenColorBackground) + (redColorBackground - blueColorBackground)
		} else {
			if prominentBGColor == "green" {
				otherColorsModifier = (greenColorBackground - redColorBackground) + (greenColorBackground - blueColorBackground)
			} else {
				otherColorsModifier = (blueColorBackground - greenColorBackground) + (blueColorBackground - redColorBackground)
			}
		}
		redColorText = uint8(224 + rand.Intn(32))
		if prominentBGColor != "red" {
			if uint32(redColorText)+uint32(otherColorsModifier) > 255 {
				redColorText = uint8(255)
			} else {
				redColorText += otherColorsModifier
			}
		}
		greenColorText = uint8(224 + rand.Intn(32))
		if prominentBGColor != "green" {
			if uint32(greenColorText)+uint32(otherColorsModifier) > 255 {
				greenColorText = uint8(255)
			} else {
				greenColorText += otherColorsModifier
			}
		}
		blueColorText = uint8(224 + rand.Intn(32))
		if prominentBGColor != "blue" {
			if uint32(blueColorText)+uint32(otherColorsModifier) > 255 {
				blueColorText = uint8(255)
			} else {
				blueColorText += otherColorsModifier
			}
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

	fmt.Println("RGB for text: R-", redColorText, ",G-", greenColorText, ",B-", blueColorText, "")

	dc.SetColor(color.RGBA{redColorText, greenColorText, blueColorText, 255})
	return false, currentPic, nil
}

func saveImage(img image.Image, fileName string) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	currentPicsFolder := filepath.Join(usr.HomeDir, ".Metamorphoun")

	fileIn := filepath.Join(currentPicsFolder, fileName)
	saveImg(img, fileIn)
}
