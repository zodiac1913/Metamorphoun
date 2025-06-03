//go:build windows
// +build windows

package main

import (
	"Metamorphoun/config"
	"Metamorphoun/service"
	"fmt"
	"image"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fogleman/gg"
	"golang.org/x/sys/windows/registry"
)

// Add to startup registry
const (
	runKeyCurrentUser = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName           = "Metamorphoun"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procBeep = kernel32.NewProc("Beep")
)

func PrintPlatformMessage() {
	fmt.Println("Running Windows-specific code")
}

func Beep(frequency, duration int) {
	procBeep.Call(uintptr(frequency), uintptr(duration))
}

func AddToStartup() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyCurrentUser, registry.WRITE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Ensure the path doesn't have any surrounding quotes that could cause issues
	exePath = strings.Trim(exePath, "\"")

	err = key.SetStringValue(appName, fmt.Sprintf("\"%s\"", exePath))
	if err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	log.Printf("%s added to Windows startup for the current user.", appName)
	return nil
}
func RemoveFromStartup() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyCurrentUser, registry.WRITE)
	if err != nil {
		// If the key doesn't exist, it's already removed or never added.
		if err == registry.ErrNotExist {
			log.Printf("%s startup entry not found for the current user.", appName)
			return nil
		}
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue(appName)
	if err != nil {
		// If the value doesn't exist, it's already removed.
		if err == registry.ErrNotExist {
			log.Printf("%s startup entry not found for the current user.", appName)
			return nil
		}
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	log.Printf("%s removed from Windows startup for the current user.", appName)
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
		return filepath.Join("C:", "Windows", "Fonts")
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
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Pictures")
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
		return filepath.Join("C:", "Programs", "ZodiSoft", "Metamorphoun")
	}
}
func SetRandomQuote(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
	var err error
	fmt.Println("running setRandomQuote")
	// Get the number of displays
	screenInfo := service.GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height
	//Make Sure a Quote is loaded
	currentPic, err = service.GetQuote(currentPic)
	if err != nil {
		fmt.Println("Error getting quote:", err)
		return currentPic, img, err
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
	authorText, wrappedQuoteText, quoteHeight, textBoxWidth, textBoxHeight, textBlockX, textBlockY, currentPic := service.CalculateBoxInfo(screenWidth, screenHeight, currentPic, dc)

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

	dc.DrawStringWrapped(wrappedQuoteText, textBlockX+10, textBlockY+30, 0, 0, textBoxWidth-20, 1.5, gg.AlignLeft)

	// Calculate a line height buffer between the quote and the author
	lineHeight := 48.0                                    // Replace with the actual height of a line of text
	authorY := textBlockY + 30 + quoteHeight + lineHeight // Add a buffer between quote and author
	dc.DrawString(authorText, textBlockX+10, authorY+30)
	// Get the resulting image (THIS IS THE MAGIC OF THE NEW PIC CONTEXT.  Started with dc := gg.NewContextForImage(img) )
	imgWithQuote := dc.Image()
	return currentPic, imgWithQuote, err

}
