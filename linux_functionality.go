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
	"strings"
	"time"

	"github.com/fogleman/gg"
)

var mbcQuotes []byte

type linuxWallpaperModule struct {
	name     string
	supports func(linuxWallpaperContext) bool
	apply    func(linuxWallpaperContext, []string) error
}

type linuxWallpaperContext struct {
	distroID    string
	idLike      []string
	sessionType string
	desktops    []string
	outputs     []string
}

var linuxWallpaperModules = []linuxWallpaperModule{
	{
		name:     "gnome-x11-xwallpaper",
		supports: supportsGnomeX11,
		apply:    applyXwallpaperModule,
	},
	{
		name:     "kde-x11-xwallpaper",
		supports: supportsKDEX11,
		apply:    applyXwallpaperModule,
	},
	{
		name:     "cinnamon-x11-xwallpaper",
		supports: supportsCinnamonX11,
		apply:    applyXwallpaperModule,
	},
	{
		name:     "ubuntu-family-x11-xwallpaper",
		supports: supportsUbuntuFamilyX11,
		apply:    applyXwallpaperModule,
	},
}

func init() {
	service.SetPerScreenWallpapers = setLinuxPerScreenWallpapersImpl
	loadMBCQuotes()
}

func setLinuxPerScreenWallpapersImpl(wallpaperPaths []string) error {
	ctx, err := detectLinuxWallpaperContext()
	if err != nil {
		return err
	}
	for _, module := range linuxWallpaperModules {
		if module.supports(ctx) {
			return module.apply(ctx, wallpaperPaths)
		}
	}
	return fmt.Errorf("no Linux wallpaper module for distro %s desktop %v on %s", ctx.distroID, ctx.desktops, ctx.sessionType)
}

func detectLinuxWallpaperContext() (linuxWallpaperContext, error) {
	ctx := linuxWallpaperContext{}
	releaseData, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ctx, fmt.Errorf("failed to read /etc/os-release: %w", err)
	}

	for _, line := range strings.Split(string(releaseData), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], "\"")
		switch key {
		case "ID":
			ctx.distroID = strings.ToLower(value)
		case "ID_LIKE":
			ctx.idLike = strings.Fields(strings.ToLower(value))
		}
	}

	ctx.sessionType = detectLinuxSessionType()
	if ctx.sessionType != "x11" {
		return ctx, fmt.Errorf("Linux multi-screen wallpapers currently require X11; detected session type %q", ctx.sessionType)
	}
	ctx.desktops = detectLinuxDesktops()

	ctx.outputs, err = getXRandROutputs()
	if err != nil {
		return ctx, err
	}
	if len(ctx.outputs) == 0 {
		return ctx, fmt.Errorf("xrandr returned no connected outputs")
	}

	return ctx, nil
}

func detectLinuxSessionType() string {
	sessionType := strings.ToLower(os.Getenv("XDG_SESSION_TYPE"))
	if sessionType == "" && os.Getenv("DISPLAY") != "" {
		return "x11"
	}
	if sessionType == "" && os.Getenv("WAYLAND_DISPLAY") != "" {
		return "wayland"
	}
	return sessionType
}

func detectLinuxDesktops() []string {
	seen := make(map[string]struct{})
	desktops := make([]string, 0)
	for _, source := range []string{os.Getenv("XDG_CURRENT_DESKTOP"), os.Getenv("DESKTOP_SESSION")} {
		for _, desktop := range splitDesktopNames(source) {
			if _, exists := seen[desktop]; exists {
				continue
			}
			seen[desktop] = struct{}{}
			desktops = append(desktops, desktop)
		}
	}
	return desktops
}

func splitDesktopNames(source string) []string {
	parts := strings.FieldsFunc(strings.ToLower(source), func(r rune) bool {
		return r == ':' || r == ';' || r == ','
	})
	results := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.TrimPrefix(part, "x-")
		if part == "" {
			continue
		}
		results = append(results, part)
	}
	return results
}

func supportsGnomeX11(ctx linuxWallpaperContext) bool {
	return supportsDesktopFamilyX11(ctx, "gnome", "ubuntu")
}

func supportsKDEX11(ctx linuxWallpaperContext) bool {
	return supportsDesktopFamilyX11(ctx, "kde", "plasma")
}

func supportsCinnamonX11(ctx linuxWallpaperContext) bool {
	return supportsDesktopFamilyX11(ctx, "cinnamon")
}

func supportsUbuntuFamilyX11(ctx linuxWallpaperContext) bool {
	if ctx.sessionType != "x11" {
		return false
	}
	if ctx.distroID == "ubuntu" || ctx.distroID == "linuxmint" || ctx.distroID == "mint" {
		return true
	}
	for _, like := range ctx.idLike {
		if like == "ubuntu" || like == "debian" {
			return true
		}
	}
	return false
}

func supportsDesktopFamilyX11(ctx linuxWallpaperContext, families ...string) bool {
	if ctx.sessionType != "x11" {
		return false
	}
	for _, desktop := range ctx.desktops {
		for _, family := range families {
			if desktop == family || strings.HasPrefix(desktop, family+"-") {
				return true
			}
		}
	}
	return false
}

func applyXwallpaperModule(ctx linuxWallpaperContext, wallpaperPaths []string) error {
	if _, err := exec.LookPath("xwallpaper"); err != nil {
		return fmt.Errorf("xwallpaper is required for Linux multi-screen wallpapers: %w", err)
	}
	if len(wallpaperPaths) == 0 {
		return fmt.Errorf("no wallpaper paths provided")
	}

	args := make([]string, 0, len(ctx.outputs)*4)
	for index, output := range ctx.outputs {
		pathIndex := index
		if pathIndex >= len(wallpaperPaths) {
			pathIndex = len(wallpaperPaths) - 1
		}
		args = append(args, "--output", output, "--zoom", wallpaperPaths[pathIndex])
	}

	cmd := exec.Command("xwallpaper", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xwallpaper failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func getXRandROutputs() ([]string, error) {
	if _, err := exec.LookPath("xrandr"); err != nil {
		return nil, fmt.Errorf("xrandr is required to enumerate Linux displays: %w", err)
	}

	cmd := exec.Command("xrandr", "--query")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("xrandr failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	outputs := make([]string, 0)
	for _, line := range strings.Split(string(output), "\n") {
		if !strings.Contains(line, " connected") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		outputs = append(outputs, fields[0])
	}

	return outputs, nil
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

func hasExistingMetamorphounTab(string) bool {
	return false
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
