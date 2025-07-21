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
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"Metamorphoun/config"
	"Metamorphoun/service"

	"github.com/fogleman/gg"
)

func initCapture() {
	C.DummyCaptureInit()
}

func SetRandomQuote(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
	fmt.Println("running setRandomQuote")

	//if(len(service.GetScreenInfo())==0) return currentPic, img, nil
	screenInfo := service.GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height

	currentPic, err := service.GetQuote(currentPic)
	if err != nil {
		fmt.Println("Error getting quote:", err)
		return currentPic, img, err
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

	authorText, wrappedQuoteText, quoteHeight, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic := service.CalculateBoxInfo(screenWidth, screenHeight, currentPic, dc)

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

	dc.DrawStringWrapped(wrappedQuoteText, textBlockX+10, textBlockY+30, 0, 0, textBoxWidth-20, 1.5, gg.AlignLeft)

	lineHeight := 48.0
	authorY := textBlockY + 30 + quoteHeight + lineHeight
	dc.DrawString(authorText, textBlockX+10, authorY+30)

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
		return filepath.Join(metamorphounBase, "Pictures")
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
