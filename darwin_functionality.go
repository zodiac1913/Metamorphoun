//go:build darwin
// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc -framework Foundation -framework ScreenCaptureKit
#cgo LDFLAGS: -framework Foundation -framework ScreenCaptureKit
void DummyCaptureInit();
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"Metamorphoun/config"
	"Metamorphoun/service"
	"Metamorphoun/shared"

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

func initCapture() {
	C.DummyCaptureInit()
}

func SetRandomQuote(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
	fmt.Println("running setRandomQuote")
	var err error

	//if(len(service.GetScreenInfo())==0) return currentPic, img, nil
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
		currentPic, err = service.GetQuote(currentPic)
		if err != nil {
			fmt.Println("Error getting quote:", err)
			return currentPic, img, err
		}
	}

	fmt.Println("Quote:", currentPic.QuoteStatement)
	fmt.Println("Author:", currentPic.QuoteAuthor)

	dc := gg.NewContextForImage(img)

	initialFontSize, fontPath, shouldReturn, currentPic, err := service.GetFontInfo(currentPic)
	if shouldReturn || err != nil {
		return currentPic, img, err
	}

	currentPic.QuoteFont = fontPath
	currentPic.QuoteFontSize = initialFontSize

	if err := dc.LoadFontFace(fontPath, initialFontSize); err != nil {
		fmt.Println("Error loading font:", err)
		return currentPic, img, err
	}

	authorText, wrappedQuoteText, _, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic := service.CalculateBoxInfo(screenWidth, screenHeight, currentPic, dc)

	textBlockX, textBlockY = service.LocateBox(textBlockX, screenWidth, textBlockY, screenHeight, textBoxWidth, textBoxHeight)

	red, green, blue, shouldReturn, currentPic, err := service.GetBackgroundColor(currentPic)
	if shouldReturn || err != nil {
		return currentPic, img, err
	}

	shouldReturn, currPic, err := service.GetOpacityAndSetBoxBackground(currentPic, dc, red, green, blue, textBlockX, textBlockY, textBoxWidth, textBoxHeight)
	if shouldReturn || err != nil {
		return currentPic, img, err
	}
	currentPic = currPic

	shouldReturn, currPic2, err := service.GetTextColor(red, green, blue, currentPic, dc)
	if shouldReturn || err != nil {
		return currentPic, img, err
	}
	currentPic = currPic2

	service.DrawQuoteText(dc, wrappedQuoteText, authorText, textBlockX, textBlockY, textBoxWidth)

	imgWithQuote := dc.Image()
	return currentPic, imgWithQuote, err
}

// ⚙️ macOS lock screen updater
func ChangeLockScreen(pic config.PicHistory) error {
	lockScreenPath := pic.SaveName

	// macOS doesn't have a direct public API for changing the lock screen image.
	// However, you can change the desktop background for all spaces with AppleScript.
	cmd := exec.Command("osascript", "-e",
		fmt.Sprintf(`tell application "System Events"
            set picture of every desktop to "%s"
        end tell`, lockScreenPath))

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to change background: %v", err)
	}

	log.Println("Desktop background updated successfully.")
	return nil
}
func PrintPlatformMessage() {
	fmt.Println("Running Mac OS-specific code")
}
func GetFolderPath(pathNeeded string) string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}

	// Common base paths
	metamorphounBase := filepath.Join(usr.HomeDir, "Library", "Application Support", "Metamorphoun")
	favPicFolderWithQuote := filepath.Join(metamorphounBase, "Favorites", "Pictures", "WithQuotes")
	favPicFolderWithoutQuote := filepath.Join(metamorphounBase, "Favorites", "Pictures", "WithOutQuotes")

	switch pathNeeded {
	case "fonts":
		return filepath.Join("/System", "Library", "Fonts")
	case "config":
		return metamorphounBase
	case "favorites":
		return filepath.Join(metamorphounBase, "Favorites")
	case "favwithquote":
		return favPicFolderWithQuote
	case "favwithoutquote":
		return favPicFolderWithoutQuote
	case "quotes":
		return filepath.Join(metamorphounBase, "Favorites", "Quotes")
	case "configfile":
		return filepath.Join(metamorphounBase, "config.json")
	case "pictures":
		return filepath.Join(usr.HomeDir, "Pictures")
	case "logs":
		return filepath.Join(metamorphounBase, "Logs")
	case "executable":
		exePath, errEP := os.Executable()
		if errEP != nil {
			fmt.Println("Error:", errEP)
		}
		return filepath.Dir(exePath)
	default:
		return filepath.Join("/usr", "local", "bin", "ZodiSoft", "Metamorphoun") // Local custom install path
	}
}
