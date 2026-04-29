//go:build windows
// +build windows

package main

import (
	"Metamorphoun/config"
	"Metamorphoun/service"
	"Metamorphoun/shared"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/sys/windows/registry"
)

var mbcQuotes []byte

func init() {
	loadMBCQuotes()
}

func loadMBCQuotes() {
	fmt.Println("Starting to load MBC quotes...")

	// Try embedded files first
	fmt.Println("Trying embedded files...")
	mbcData, err := shared.GetStaticFSQuotes("quotes/mbc.json")
	if err != nil {
		fmt.Printf("Embedded loading failed: %v\n", err)
		// Fallback to file system for development
		fmt.Println("Trying file system fallback...")
		mbcFilePath := filepath.Join("shared", "static", "quotes", "mbc.json")
		fmt.Printf("File path: %s\n", mbcFilePath)
		mbcData, err = os.ReadFile(mbcFilePath)
		if err != nil {
			fmt.Printf("File system loading also failed: %v\n", err)
			return
		}
		mbcQuotes = mbcData
		fmt.Printf("Successfully loaded from file system: %d bytes\n", len(mbcData))
	} else {
		fmt.Printf("Successfully loaded from embedded: %d bytes\n", len(mbcData))
	}
	mbcQuotes = mbcData
}

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
		return filepath.Join("C:\\", "Windows", "Fonts")
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
		staticImagesPath := filepath.Join(exeDir, "shared", "static", "images")
		if _, err := os.Stat(staticImagesPath); os.IsNotExist(err) {
			if cwd, err := os.Getwd(); err == nil {
				cwdStatic := filepath.Join(cwd, "shared", "static", "images")
				if _, err := os.Stat(cwdStatic); err == nil {
					return cwd
				}
			}
		}
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
				// Check if the month changed
				currentMonth := int(time.Now().Month())
				fmt.Println("Current month:", currentMonth, "MBCMonth:", config.ConfigInstance.MBCMonth)
				if config.ConfigInstance.MBCMonth != currentMonth {
					// Advance MBCMonth (wrap 13 -> 1)
					config.ConfigInstance.MBCMonth = currentMonth
					// Advance MBCValue by 1, wrap around if past end
					config.ConfigInstance.MBCValue++
					if config.ConfigInstance.MBCValue >= len(quotes) {
						config.ConfigInstance.MBCValue = 0
					}
					fmt.Println("Month changed — MBCValue now:", config.ConfigInstance.MBCValue)
				}
				// Use the current MBCValue (safe mod in case config was hand-edited)
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
		fmt.Println("MBCValue:", config.ConfigInstance.MBCValue)
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
	// Get the path to the lock screen image
	lockScreenPath := pic.SaveName

	// Use the Set-ItemProperty cmdlet in PowerShell to change the lock screen background
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Set-ItemProperty -Path "HKCU:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\PersonalizationCSP" -Name "LockScreenImageFilename" -Value "%s";
                        Set-ItemProperty -Path "HKCU:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\PersonalizationCSP" -Name "LockScreenImageClipartEnabled" -Value 0`, lockScreenPath))

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to change lock screen image: %v", err)
	}

	log.Println("Lock screen image changed successfully.")
	return nil
}
