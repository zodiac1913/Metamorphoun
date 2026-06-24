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
	"math"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/fogleman/gg"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var mbcQuotes []byte

func init() {
	service.SetPerScreenWallpapers = setWindowsPerScreenWallpapersImpl
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
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procBeep             = kernel32.NewProc("Beep")
	ole32                = windows.NewLazySystemDLL("ole32.dll")
	procCoInitializeEx   = ole32.NewProc("CoInitializeEx")
	procCoUninitialize   = ole32.NewProc("CoUninitialize")
	procCoCreateInstance = ole32.NewProc("CoCreateInstance")
)

const (
	coinitApartmentThreaded = 0x2
	clsctxInprocServer      = 0x1
	wallpaperPositionFill   = 4
	sFalse                  = 0x00000001
)

var (
	clsidDesktopWallpaper = windows.GUID{Data1: 0xC2CF3110, Data2: 0x460E, Data3: 0x4FC1, Data4: [8]byte{0xB9, 0xD0, 0x8A, 0x1C, 0x0C, 0xD8, 0x57, 0xD1}}
	iidDesktopWallpaper   = windows.GUID{Data1: 0xB92B56A9, Data2: 0x8B55, Data3: 0x4E14, Data4: [8]byte{0x9A, 0x89, 0x01, 0x99, 0xBB, 0xB6, 0xF9, 0x3B}}
)

type desktopWallpaper struct {
	lpVtbl *desktopWallpaperVtbl
}

type desktopWallpaperVtbl struct {
	queryInterface            uintptr
	addRef                    uintptr
	release                   uintptr
	setWallpaper              uintptr
	getWallpaper              uintptr
	getMonitorDevicePathAt    uintptr
	getMonitorDevicePathCount uintptr
	getMonitorRECT            uintptr
	setBackgroundColor        uintptr
	getBackgroundColor        uintptr
	setPosition               uintptr
	getPosition               uintptr
	setSlideshow              uintptr
	getSlideshow              uintptr
	setSlideshowOptions       uintptr
	getSlideshowOptions       uintptr
	advanceSlideshow          uintptr
	getStatus                 uintptr
	enable                    uintptr
}

func setWindowsPerScreenWallpapersImpl(wallpaperPaths []string) error {
	if len(wallpaperPaths) == 0 {
		return fmt.Errorf("no wallpaper paths provided")
	}

	uninit, err := coInitialize()
	if err != nil {
		return err
	}
	defer uninit()

	desktop, err := createDesktopWallpaper()
	if err != nil {
		return err
	}
	defer desktop.Release()

	if err := desktop.SetPosition(wallpaperPositionFill); err != nil {
		return err
	}

	count, err := desktop.GetMonitorDevicePathCount()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("IDesktopWallpaper reported zero monitors")
	}

	for monitorIndex := uint32(0); monitorIndex < count; monitorIndex++ {
		monitorID, err := desktop.GetMonitorDevicePathAt(monitorIndex)
		if err != nil {
			return err
		}
		pathIndex := int(math.Min(float64(monitorIndex), float64(len(wallpaperPaths)-1)))
		if err := desktop.SetWallpaper(monitorID, wallpaperPaths[pathIndex]); err != nil {
			return err
		}
	}

	return nil
}

func coInitialize() (func(), error) {
	hr, _, _ := procCoInitializeEx.Call(0, uintptr(coinitApartmentThreaded))
	switch uint32(hr) {
	case 0, sFalse:
		return func() {
			procCoUninitialize.Call()
		}, nil
	default:
		return nil, fmt.Errorf("CoInitializeEx failed: HRESULT 0x%08x", uint32(hr))
	}
}

func createDesktopWallpaper() (*desktopWallpaper, error) {
	var instance *desktopWallpaper
	hr, _, _ := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsidDesktopWallpaper)),
		0,
		uintptr(clsctxInprocServer),
		uintptr(unsafe.Pointer(&iidDesktopWallpaper)),
		uintptr(unsafe.Pointer(&instance)),
	)
	if uint32(hr) != 0 {
		return nil, fmt.Errorf("CoCreateInstance for IDesktopWallpaper failed: HRESULT 0x%08x", uint32(hr))
	}
	return instance, nil
}

func (d *desktopWallpaper) Release() {
	if d == nil || d.lpVtbl == nil {
		return
	}
	syscall.SyscallN(d.lpVtbl.release, uintptr(unsafe.Pointer(d)))
}

func (d *desktopWallpaper) SetPosition(position uint32) error {
	hr, _, _ := syscall.SyscallN(d.lpVtbl.setPosition, uintptr(unsafe.Pointer(d)), uintptr(position))
	if uint32(hr) != 0 {
		return fmt.Errorf("IDesktopWallpaper::SetPosition failed: HRESULT 0x%08x", uint32(hr))
	}
	return nil
}

func (d *desktopWallpaper) GetMonitorDevicePathCount() (uint32, error) {
	var count uint32
	hr, _, _ := syscall.SyscallN(d.lpVtbl.getMonitorDevicePathCount, uintptr(unsafe.Pointer(d)), uintptr(unsafe.Pointer(&count)))
	if uint32(hr) != 0 {
		return 0, fmt.Errorf("IDesktopWallpaper::GetMonitorDevicePathCount failed: HRESULT 0x%08x", uint32(hr))
	}
	return count, nil
}

func (d *desktopWallpaper) GetMonitorDevicePathAt(index uint32) (string, error) {
	var monitorID *uint16
	hr, _, _ := syscall.SyscallN(d.lpVtbl.getMonitorDevicePathAt, uintptr(unsafe.Pointer(d)), uintptr(index), uintptr(unsafe.Pointer(&monitorID)))
	if uint32(hr) != 0 {
		return "", fmt.Errorf("IDesktopWallpaper::GetMonitorDevicePathAt failed: HRESULT 0x%08x", uint32(hr))
	}
	if monitorID == nil {
		return "", fmt.Errorf("IDesktopWallpaper returned nil monitor ID for index %d", index)
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(monitorID))
	return windows.UTF16PtrToString(monitorID), nil
}

func (d *desktopWallpaper) SetWallpaper(monitorID string, wallpaperPath string) error {
	monitorIDPtr, err := windows.UTF16PtrFromString(monitorID)
	if err != nil {
		return fmt.Errorf("invalid monitor ID: %w", err)
	}
	wallpaperPtr, err := windows.UTF16PtrFromString(wallpaperPath)
	if err != nil {
		return fmt.Errorf("invalid wallpaper path: %w", err)
	}
	hr, _, _ := syscall.SyscallN(d.lpVtbl.setWallpaper, uintptr(unsafe.Pointer(d)), uintptr(unsafe.Pointer(monitorIDPtr)), uintptr(unsafe.Pointer(wallpaperPtr)))
	if uint32(hr) != 0 {
		return fmt.Errorf("IDesktopWallpaper::SetWallpaper failed for %s: HRESULT 0x%08x", wallpaperPath, uint32(hr))
	}
	return nil
}

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
