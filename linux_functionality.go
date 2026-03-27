//go:build linux
// +build linux

// linux_functionality.go
package main

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/morphLog"
	"Metamorphoun/service"
	"Metamorphoun/shared"
	"Metamorphoun/zutil"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
)

var mbcQuotes []byte

func init() {
	loadMBCQuotes()
}

func loadMBCQuotes() {
	mbcData, err := shared.GetStaticFSQuotes("quotes/mbc.json")
	if err != nil {
		fmt.Println("Error loading MBC quotes:", err)
		return
	}
	mbcQuotes = mbcData
}

func PrintPlatformMessage() {
	fmt.Println("Running Linux-specific code")
}

func AddToStartup() error {
	cronJob := "@reboot /path/to/your/application\n"
	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | crontab -u youruser -", cronJob))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}
	log.Println("Application added to Linux startup via cron.")
	return nil
}

func RemoveFromStartup() error {
	cmd := exec.Command("bash", "-c", "crontab -u youruser -l | grep '/path/to/your/application' && crontab -u youruser -e")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to remove cron job: %w", err)
	}
	log.Println("Application removed from Linux startup via cron.")
	return nil
}

func GetFolderPath(pathNeeded string) string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	favPicFolderWithQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithQuotes")
	favPicFolderWithoutQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithOutQuotes")

	if pathNeeded == "fonts" {
		return filepath.Join("/usr", "share", "fonts")
	} else if pathNeeded == "config" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun")
	} else if pathNeeded == "favorites" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites")
	} else if pathNeeded == "favwithquote" {
		return favPicFolderWithQuote
	} else if pathNeeded == "favwithoutquote" {
		return favPicFolderWithoutQuote
	} else if pathNeeded == "quotes" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Quotes")
	} else if pathNeeded == "configfile" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "config.json")
	} else if pathNeeded == "pictures" {
		return filepath.Join(usr.HomeDir, "Pictures")
	} else if pathNeeded == "logs" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Logs")
	} else if pathNeeded == "executable" {
		exePath, errEP := os.Executable()
		if errEP != nil {
			fmt.Println("Error:", errEP)
		}
		exeDir := filepath.Dir(exePath)
		return exeDir
	} else {
		return filepath.Join("usr", "bin", "ZodiSoft", "Metamorphoun")
	}
}

// Common font directories
var fontDirs = []string{
	"/usr/share/fonts",
	"/usr/local/share/fonts",
	"~/.local/share/fonts",
	"~/.fonts",
	"C:\\Windows\\Fonts",
}

func findFonts(currentPic config.PicHistory) (float64, string, bool, config.PicHistory, error) {
	var foundFonts []string
	minSize := config.ConfigInstance.QuoteFontSizeMin
	maxSize := config.ConfigInstance.QuoteFontSizeMax
	if minSize < 8 {
		minSize = 16
	}
	if maxSize < minSize {
		maxSize = minSize
	}
	initialFontSize := minSize
	if maxSize > minSize {
		initialFontSize = minSize + float64(rand.Intn(int(maxSize-minSize+1)))
	}
	fontPath := filepath.Join(GetFolderPath(enum.PathLoc.Fonts), config.ConfigInstance.TextFontFile)
	for _, dir := range fontDirs {
		expandedDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}

		// Walk through directory recursively
		filepath.Walk(expandedDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && (filepath.Ext(path) == ".ttf" || filepath.Ext(path) == ".otf") {
				foundFonts = append(foundFonts, path)
			}
			return nil
		})
	}

	if config.ConfigInstance.QuoteFontRandom {

		// Select a random valid font
		fileRnd := rand.Intn(len(foundFonts))
		fontPath := foundFonts[fileRnd]
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
		if zutil.IsInRange(fontPath, foundFonts) {
			fontPath = filepath.Join(GetFolderPath(enum.PathLoc.Fonts), config.ConfigInstance.TextFontFile)
		} else {
			fontPath = foundFonts[0]
		}

	}
	fmt.Println("Selected font:", fontPath)
	currentPic.QuoteFont = fontPath
	currentPic.QuoteStatement = config.ConfigInstance.CurrentQuoteStatement
	currentPic.QuoteAuthor = config.ConfigInstance.CurrentQuoteAuthor
	return initialFontSize, fontPath, false, currentPic, nil
}

func SetRandomQuote(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
	var err error
	fmt.Println("running setRandomQuote")
	// Get the number of displays
	screenInfo := service.GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height

	if config.ConfigInstance.MBCMode {
		fmt.Println("mbc mode active, using MBC quotes")
		if len(mbcQuotes) == 0 {
			currentPic.QuoteStatement = "MBC Quotes not loaded"
			currentPic.QuoteAuthor = ""
		} else {
			var quotes []struct {
				Statement string `json:"statement"`
				Author    string `json:"author"`
			}
			err = json.Unmarshal(mbcQuotes, &quotes)
			if err != nil {
				fmt.Printf("JSON unmarshal failed: %v\n", err)
				currentPic.QuoteStatement = "MBC Quotes unmarshal failed"
				currentPic.QuoteAuthor = ""
			} else if len(quotes) > 0 {
				currentMonth := int(time.Now().Month())
				if config.ConfigInstance.MBCMonth != currentMonth {
					config.ConfigInstance.MBCMonth = currentMonth
					config.ConfigInstance.MBCValue++
					if config.ConfigInstance.MBCValue >= len(quotes) {
						config.ConfigInstance.MBCValue = 0
					}
					fmt.Println("Month changed — MBCValue now:", config.ConfigInstance.MBCValue)
				}
				idx := config.ConfigInstance.MBCValue % len(quotes)
				currentPic.QuoteStatement = quotes[idx].Statement
				currentPic.QuoteAuthor = quotes[idx].Author
				fmt.Println("Quote set to:", currentPic.QuoteStatement, "by", currentPic.QuoteAuthor)
			} else {
				currentPic.QuoteStatement = "MBC Quotes empty"
				currentPic.QuoteAuthor = ""
			}
		}
		config.UpdateConfigField("currentQuoteStatement", currentPic.QuoteStatement)
		config.UpdateConfigField("currentQuoteAuthor", currentPic.QuoteAuthor)
		if err := config.SaveConfig(config.ConfigInstance); err != nil {
			fmt.Println("Failed to save MBC config:", err)
		}
	} else {
		//Make Sure a Quote is loaded
		currentPic, err = service.GetQuote(currentPic)
		if err != nil {
			fmt.Println("Error getting quote:", err)
			return currentPic, img, err
		}
	}
	fmt.Println("Quote:", currentPic.QuoteStatement)
	fmt.Println("Author:", currentPic.QuoteAuthor)

	// Create a new context with the image dimensions
	dc := gg.NewContextForImage(img)

	// Set initial font size
	initialFontSize, fontPath, shouldReturn, currentPic, err := service.GetFontInfo(currentPic)
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
	authorText, wrappedQuoteText, _, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic := service.CalculateBoxInfo(screenWidth, screenHeight, currentPic, dc)

	textBlockX, textBlockY = service.LocateBox(textBlockX, screenWidth, textBlockY, screenHeight, textBoxWidth, textBoxHeight)

	// Set transparent background for text block
	//Make Background color
	redColorBackground, greenColorBackground, blueColorBackground, shouldReturn, currentPic, err := service.GetBackgroundColor(currentPic)
	if shouldReturn {
		return currentPic, img, err
	}

	shouldReturn, currPic, err := service.GetOpacityAndSetBoxBackground(currentPic, dc, redColorBackground, greenColorBackground, blueColorBackground, textBlockX, textBlockY, textBoxWidth, textBoxHeight)
	if shouldReturn {
		return currentPic, img, err
	}
	currentPic = currPic
	// Set text color and draw text
	//Make Text color
	shouldReturn, currPic2, err := service.GetTextColor(redColorBackground, greenColorBackground, blueColorBackground, currentPic, dc)
	if shouldReturn {
		return currentPic, img, err
	}
	currentPic = currPic2
	//dc.SetColor(color.White)

	service.DrawQuoteText(dc, wrappedQuoteText, authorText, textBlockX, textBlockY, textBoxWidth)

	imgWithQuote := dc.Image()
	return currentPic, imgWithQuote, err

}

func ChangeLockScreen(pic config.PicHistory) error {
	lockScreenPath := pic.SaveName
	// On Linux, changing the lock screen varies by desktop environment.
	// This uses gsettings for GNOME-based desktops.
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.screensaver", "picture-uri", "file://"+lockScreenPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to change lock screen image: %v", err)
		return err
	}
	log.Println("Lock screen image changed successfully.")
	return nil
}
