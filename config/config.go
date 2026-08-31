package config

import (
	"Metamorphoun/enum"
	"Metamorphoun/zutil"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const AppVersion = "2026.08.29.1"
const PublishedOn = "2026-08-29.1"

//ugh

var GetFolderPath func(string) string

type PathLocType string

// Define the structure of your configuration ...
type Config struct {
	Version                         string  `json:"version"`
	Published                       string  `json:"published"`
	ServerAddress                   string  `json:"server_address"`
	ServerPort                      int     `json:"serverPort"`
	SourceCurrentBackgroundName     string  `json:"sourceCurrentBackgroundName"`
	SourceCurrentBackgroundFolder   string  `json:"sourceCurrentBackgroundFolder"`
	OriginalCurrentBackgroundName   string  `json:"originalCurrentBackgroundName"`
	OriginalCurrentBackgroundFolder string  `json:"originalCurrentBackgroundFolder"`
	CurrentBackgroundName           string  `json:"currentBackgroundName"`
	CurrentBackgroundFolder         string  `json:"currentBackgroundFolder"`
	BackgroundChangingBlock         bool    `json:"backgroundChangingBlock"`
	BackgroundChangeAttempt         int     `json:"backgroundChangeAttempt"`
	StartOnStartup                  bool    `json:"startOnStartup"`
	ChangeWallpaperOnStartup        bool    `json:"changeWallpaperOnStartup"`
	DifferentWallpaperPerScreen     bool    `json:"differentWallpaperPerScreen"`
	ChangeMinutes                   int32   `json:"changeMinutes"`
	Images                          []Image `json:"images"`
	WallpaperImageSizing            string  `json:"wallpaperImageSizing"`
	WallpaperFilterOriginal         bool    `json:"wallpaperFilterOriginal"`
	WallpaperFilterBlurSoft         bool    `json:"wallpaperFilterBlurSoft"`
	WallpaperFilterBlurHard         bool    `json:"wallpaperFilterBlurHard"`
	WallpaperFilterPixelate         bool    `json:"wallpaperFilterPixelate"`
	WallpaperFilterOilify           bool    `json:"wallpaperFilterOilify"`
	WallpaperFilterWavy             bool    `json:"wallpaperFilterWavy"`
	WallpaperFilterVortex           bool    `json:"wallpaperFilterVortex"`
	WallpaperFilterMosaic           bool    `json:"wallpaperFilterMosaic"`
	WallpaperFilterJigsawPuzzle     bool    `json:"wallpaperFilterJigsawPuzzle"`
	WallpaperFilterGraffiti         bool    `json:"wallpaperFilterGraffiti"`
	WallpaperFilterCartoon          bool    `json:"wallpaperFilterCartoon"`
	WallpaperFilterCyberpunk        bool    `json:"wallpaperFilterCyberpunk"`
	//WallpaperFilterSpiral           bool          `json:"wallpaperFilterSpiral"`
	WallpaperFilterMonochrome bool          `json:"wallpaperFilterMonochrome"`
	ShowTextOverlay           bool          `json:"showTextOverlay"`
	TextChangeMinutes         int           `json:"textChangeMinutes"`
	TextLibraries             []TextLibrary `json:"textLibraries"`
	TextFontFile              string        `json:"textFontFile"`
	//TextFontPath              string        `json:"textFontPath"`
	TextBoxLocation             string       `json:"textBoxLocation"`
	CurrentQuoteStatement       string       `json:"currentQuoteStatement"`
	CurrentQuoteAuthor          string       `json:"currentQuoteAuthor"`
	QuoteAppearanceRandom       bool         `json:"quoteAppearanceRandom"`
	QuoteFontRandom             bool         `json:"quoteFontRandom"`
	QuoteTextColor              string       `json:"quoteTextColor"`
	QuoteBackgroundColor        string       `json:"quoteBackgroundColor"`
	QuoteBackgroundOpacity      string       `json:"quoteBackgroundOpacity"`
	PicHistories                []PicHistory `json:"picHistories"`
	DarwinPerScreenPicHistories []PicHistory `json:"darwinPerScreenPicHistories"`
	PicUpdateCalled             bool         `json:"picUpdateCalled"`
	MBCMonth                    int          `json:"mbcMonth"`
	MBCMode                     bool         `json:"mbcMode"`
	MBCValue                    int          `json:"mbcValue"`
	QuoteFontSizeMin            float64      `json:"quoteFontSizeMin"`
	QuoteFontSizeMax            float64      `json:"quoteFontSizeMax"`
	// Add other configuration fields here
}
type Image struct {
	Use          bool   `json:"use"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	Location     string `json:"location"`
	Operation    string `json:"operation"`
	AllowDistort bool   `json:"allowDistort"`
	RequiresKey  bool   `json:"requiresKey"`
	APIKey       string `json:"apiKey,omitempty"`
	HasAPIKey    bool   `json:"hasApiKey,omitempty"`
	Inherent     bool   `json:"inherent"` // Indicates if the image is inherent to the system
}

type TextLibrary struct {
	Use      bool   `json:"use"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Location string `json:"location"`
	Citation string `json:"citation"`
	Creators string `json:"creators"`
	Info     string `json:"info"`
	Inherent bool   `json:"inherent"` // Indicates if the quote file is inherent to the system
}

type PicHistory struct {
	PicNum                int16              `json:"picNum"`
	OriginName            string             `json:"originName"`
	SaveName              string             `json:"saveName"`
	PerScreenPics         []PicHistory       `json:"perScreenPics,omitempty"`
	ImageItem             Image              `json:"imageItem"`
	Filter                string             `json:"filter"`
	FilterVortices        []PicHistoryVortex `json:"filterVortices"`
	FilterIntensity       float64            `json:"filterIntensity"`
	FilterX               float64            `json:"filterX"`
	FilterY               float64            `json:"filterY"`
	Sizing                string             `json:"sizing"`
	QuoteStatement        string             `json:"quoteStatement"`
	QuoteAuthor           string             `json:"quoteAuthor"`
	QuoteFont             string             `json:"quoteFont"`
	QuoteFontSize         float64            `json:"quoteFontSize"`
	QuoteTextColorR       uint8              `json:"quoteTextColorR"`
	QuoteTextColorG       uint8              `json:"quoteTextColorG"`
	QuoteTextColorB       uint8              `json:"quoteTextColorB"`
	QuoteBackgroundColorR uint8              `json:"quoteBackgroundColorR"`
	QuoteBackgroundColorG uint8              `json:"quoteBackgroundColorG"`
	QuoteBackgroundColorB uint8              `json:"quoteBackgroundColorB"`
	QuoteOpacity          uint64             `json:"quoteOpacity"`
	QuoteTextBoxWidth     float64            `json:"quoteTextBoxWidth"`
	QuoteTextBoxHeight    float64            `json:"quoteTextBoxHeight"`
	QuoteTextBoxX         float64            `json:"quoteTextBoxX"`
	QuoteTextBoxY         float64            `json:"quoteTextBoxY"`
}
type PicHistoryVortex struct {
	FilterIntensity float64 `json:"filterIntensity"`
	FilterQuadrant  string  `json:"filterQuadrant"`
	FilterX         float64 `json:"filterX"`
	FilterY         float64 `json:"filterY"`
}

var ConfigInstance *Config

// PerScreenSupported reports whether the current platform can display a
// different wallpaper on each screen. It defaults to true so Windows, macOS,
// and standard Linux desktops are unaffected. A platform-specific init() may
// set it to false (e.g. Omarchy/Hyprland, whose Quickshell background layer
// only accepts a single image). It is a runtime flag and is never persisted
// to the config file; the server injects it into the /configApi response so
// the web UI can hide the per-screen toggle where it does nothing.
var PerScreenSupported = true

var (
	loadedConfig *Config
	loadOnce     sync.Once
	loadError    error
	configMu     sync.RWMutex
)

func init() {
	// Load the configuration
}

// GetConfig returns the current Config instance
func GetConfig() *Config {
	configMu.Lock()
	defer configMu.Unlock()

	if loadedConfig != nil {
		loadedConfig.Version = AppVersion
		loadedConfig.Published = PublishedOn
		return loadedConfig
	}
	if ConfigInstance != nil {
		ConfigInstance.Version = AppVersion
		ConfigInstance.Published = PublishedOn
		loadedConfig = ConfigInstance
		return ConfigInstance
	}

	// Handle the case where loading failed and no config exists yet.
	fmt.Println("Warning: Config not loaded yet. Call LoadConfig first.")
	return &Config{} // Return a default empty config to avoid nil pointer
}

// OLD
//
//	func GetConfig() *Config {
//		cfg, err := LoadConfig()
//		if err != nil {
//			fmt.Println("Error loading config:", err)
//			// Handle the error (e.g., create a default config)
//			cfg = &Config{ServerAddress: "default_address"} // Set default values
//		}
//		ConfigInstance = cfg
//		return ConfigInstance
//	}
func GetConfigCopy() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return *ConfigInstance
}

// SetConfig updates the Config instance and saves it to the file
func SetConfig(newConfig *Config) error {
	configMu.Lock()
	defer configMu.Unlock()
	ConfigInstance = newConfig
	loadedConfig = newConfig
	return saveConfigUnlocked(newConfig)
}

func UpdateConfig(mutator func(cfg *Config) error) error {
	configMu.Lock()
	defer configMu.Unlock()

	if ConfigInstance == nil {
		if loadedConfig != nil {
			ConfigInstance = loadedConfig
		} else {
			return fmt.Errorf("config not loaded")
		}
	}

	if err := mutator(ConfigInstance); err != nil {
		return err
	}
	ConfigInstance.Version = AppVersion
	ConfigInstance.Published = PublishedOn
	loadedConfig = ConfigInstance
	return saveConfigUnlocked(ConfigInstance)
}

// create a function to load a config.ConfigInstance.Image by name
func GetImageByName(name string) *Image {
	for _, img := range ConfigInstance.Images {
		if img.Name == name {
			return &img
		}
	}
	return nil // Return nil if no image with the given name is found
}

func UpdateConfigField(propertyName string, newValue interface{}) error {
	return UpdateConfig(func(cfg *Config) error {
		ConfigInstance = cfg
		return SetConfigField(propertyName, newValue)
	})
}

func SetConfigField(fieldName string, value interface{}) error {
	v := reflect.ValueOf(ConfigInstance).Elem()
	f := CaseInsensitiveFieldByName(v, fieldName)
	if !f.IsValid() {
		return fmt.Errorf("no such field: %s", fieldName)
	}
	if !f.CanSet() {
		return fmt.Errorf("cannot set field: %s", fieldName)
	}

	val := reflect.ValueOf(value)

	// Handle string → numeric conversion
	if f.Kind() == reflect.Int32 && val.Kind() == reflect.String {
		parsed, err := strconv.Atoi(val.String())
		if err != nil {
			return fmt.Errorf("invalid int32 value for %s: %v", fieldName, err)
		}
		val = reflect.ValueOf(int32(parsed))
	}
	if f.Kind() == reflect.Float64 && val.Kind() == reflect.String {
		parsed, err := strconv.ParseFloat(val.String(), 64)
		if err != nil {
			return fmt.Errorf("invalid float64 value for %s: %v", fieldName, err)
		}
		val = reflect.ValueOf(parsed)
	}
	if f.Kind() == reflect.Int && val.Kind() == reflect.String {
		parsed, err := strconv.Atoi(val.String())
		if err != nil {
			return fmt.Errorf("invalid int value for %s: %v", fieldName, err)
		}
		val = reflect.ValueOf(parsed)
	}
	// Handle other type mismatches
	if val.Type() != f.Type() {
		if val.Type().ConvertibleTo(f.Type()) {
			val = val.Convert(f.Type())
		} else {
			return fmt.Errorf("cannot convert %s to %s", val.Type(), f.Type())
		}
	}

	f.Set(val)
	return nil
}

// func SetConfigField(fieldName string, value interface{}) error {

// 	v := reflect.ValueOf(ConfigInstance).Elem()
// 	f := CaseInsensitiveFieldByName(v, fieldName)
// 	if !f.IsValid() {
// 		return fmt.Errorf("no such field: %s", fieldName)
// 	}
// 	if !f.CanSet() {
// 		return fmt.Errorf("cannot set field: %s", fieldName)
// 	}
// 	val := reflect.ValueOf(value)
// 	// Convert value to the correct type if needed
// 	if val.Type() != f.Type() {
// 		val = val.Convert(f.Type())
// 	}
// 	f.Set(val)
// 	return nil
// }

// CaseInsensitiveFieldByName returns the struct field with the given name (case-insensitive).
func CaseInsensitiveFieldByName(v reflect.Value, name string) reflect.Value {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	lower := strings.ToLower(name)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if strings.ToLower(t.Field(i).Name) == lower {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

func AddToStartup() error {
	err := AddToStartup()
	if err != nil {
		log.Println("Error adding to startup:", err)
		return err
	}
	return nil
}

func RemoveFromStartup() error {
	err := RemoveFromStartup()
	if err != nil {
		log.Println("Error removing from startup:", err)
		return err
	}
	return nil
}

func UpdateImagesField(imageName string, newValue bool) error {
	return UpdateConfig(func(cfg *Config) error {
		for i, image := range cfg.Images {
			if image.Name == imageName {
				cfg.Images[i].Use = newValue
				return nil
			}
		}
		return fmt.Errorf("image source not found: %s", imageName)
	})
}

func UpdateImageAPIKey(imageName string, apiKey string) error {
	return UpdateConfig(func(cfg *Config) error {
		for i, image := range cfg.Images {
			if image.Name == imageName {
				if !image.RequiresKey {
					return fmt.Errorf("image source does not require an API key: %s", imageName)
				}
				cfg.Images[i].APIKey = strings.TrimSpace(apiKey)
				return nil
			}
		}
		return fmt.Errorf("image source not found: %s", imageName)
	})
}
func AddImagesField(use bool, name string, title string,
	location string, operation string) error {
	return UpdateConfig(func(cfg *Config) error {
		cfg.Images = append(cfg.Images, Image{
			Use:       use,
			Name:      name,
			Title:     title,
			Location:  location,
			Operation: operation,
		})
		return nil
	})
}
func EditImagesField(use bool, name string, title string,
	location string, operation string) error {
	return UpdateConfig(func(cfg *Config) error {
		for i := range cfg.Images {
			if cfg.Images[i].Name != name {
				continue
			}
			if cfg.Images[i].Inherent {
				fmt.Println("Cannot edit inherent image:", name)
				return fmt.Errorf("cannot edit inherent image: %s", name)
			}
			cfg.Images[i].Use = use
			cfg.Images[i].Title = title
			cfg.Images[i].Location = location
			cfg.Images[i].Operation = operation
			return nil
		}
		return fmt.Errorf("image source not found: %s", name)
	})
}

func UpdateQuotesField(quotesName string, newValue interface{}) error {
	return UpdateConfig(func(cfg *Config) error {
		for i, textLib := range cfg.TextLibraries {
			if textLib.Name == quotesName {
				cfg.TextLibraries[i].Use = zutil.AsBool(newValue.(string))
				return nil
			}
		}
		return fmt.Errorf("text library not found: %s", quotesName)
	})

}

// AddPicHistory adds a new PicHistory to the stack, updates PicNum,
// and ensures the stack size does not exceed the limit.
const picHistoryLimit = 10

func (cfg *Config) AddPicHistoryInPlace(newPic PicHistory) error {
	// Prepend the new PicHistory to the stack
	cfg.PicHistories = append([]PicHistory{newPic}, cfg.PicHistories...)

	// Ensure the stack size does not exceed the limit (10 for the History tab).
	// Entries that fall off the end have their locally-cached image files
	// (unique per-pic originals and saved wallpapers) deleted so they don't
	// accumulate on disk.
	if len(cfg.PicHistories) > picHistoryLimit {
		dropped := cfg.PicHistories[picHistoryLimit:]
		cfg.PicHistories = cfg.PicHistories[:picHistoryLimit]
		cleanupDroppedPicFiles(dropped, cfg.PicHistories)
	}

	// Update PicNum for all PicHistories in the stack
	for i := range cfg.PicHistories {
		cfg.PicHistories[i].PicNum = int16(i)
	}
	return nil
}

func (cfg *Config) AddPicHistory(newPic PicHistory) error {
	return UpdateConfig(func(currentCfg *Config) error {
		return currentCfg.AddPicHistoryInPlace(newPic)
	})
}

func collectReferencedPicPaths(paths map[string]bool, pic PicHistory) {
	if pic.OriginName != "" {
		paths[pic.OriginName] = true
	}
	if pic.SaveName != "" {
		paths[pic.SaveName] = true
	}
	for _, perScreenPic := range pic.PerScreenPics {
		collectReferencedPicPaths(paths, perScreenPic)
	}
}

// cleanupDroppedPicFiles deletes the locally-cached image files belonging to
// pic-history entries that have aged off the end of the stack. It only removes
// files inside the config directory (so user folders like Favorites/Pictures
// are never touched) and never removes a file still referenced by a retained
// entry.
func cleanupDroppedPicFiles(dropped []PicHistory, retained []PicHistory) {
	configDir, err := filepath.Abs(GetFolderPath(enum.PathLoc.Config))
	if err != nil || configDir == "" {
		return
	}

	// Collect paths still in use so we never delete a shared/duplicate file.
	stillUsed := make(map[string]bool)
	for _, pic := range retained {
		collectReferencedPicPaths(stillUsed, pic)
	}
	for _, pic := range ConfigInstance.DarwinPerScreenPicHistories {
		collectReferencedPicPaths(stillUsed, pic)
	}

	inConfigDir := func(path string) bool {
		absPath, absErr := filepath.Abs(path)
		if absErr != nil {
			return false
		}
		rel, relErr := filepath.Rel(configDir, absPath)
		return relErr == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
	}

	for _, pic := range dropped {
		picPaths := []string{pic.OriginName, pic.SaveName}
		for _, perScreenPic := range pic.PerScreenPics {
			picPaths = append(picPaths, perScreenPic.OriginName, perScreenPic.SaveName)
		}
		for _, path := range picPaths {
			if path == "" || stillUsed[path] {
				continue
			}
			if strings.HasPrefix(strings.ToLower(path), "http") {
				continue // remote URL, nothing on disk
			}
			if !inConfigDir(path) {
				continue // outside the config sandbox (e.g. Favorites, user Pictures)
			}
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				fmt.Println("AddPicHistory: failed to remove aged pic file:", path, err)
			}
		}
	}
}

// LoadConfig reads the configuration from the JSON file
func LoadConfig() (*Config, error) {
	loadOnce.Do(func() {
		// Get the user's home directory
		// usr, err := user.Current()
		// if err != nil {
		// 	loadError = fmt.Errorf("failed to get user home directory: %w", err)
		// 	return
		// }
		//pathLoc :=
		configPath := GetFolderPath(enum.PathLoc.ConfigFile)

		// Read the config file
		data, err := os.ReadFile(configPath)
		if err != nil {
			loadError = fmt.Errorf("failed to read config file: %w", err)
			return
		}

		// Unmarshal the JSON data into the Config struct
		var config Config
		err = json.Unmarshal(data, &config)
		if err != nil {
			loadError = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}
		configMu.Lock()
		loadedConfig = &config
		ConfigInstance = &config
		configMu.Unlock()
	})
	return loadedConfig, loadError
}

// MigrateConfig ensures the loaded config has all inherent items that a fresh
// install would have (new filters, text libraries, image sources, etc.).
// Call this once after LoadConfig so existing users pick up new additions.
func MigrateConfig(cfg *Config) bool {
	changed := false

	// --- Inherent TextLibraries -------------------------------------------
	canonical := canonicalTextLibraries()
	existingLibs := make(map[string]bool)
	for _, tl := range cfg.TextLibraries {
		existingLibs[tl.Name] = true
	}
	for _, tl := range canonical {
		if !existingLibs[tl.Name] {
			fmt.Println("MigrateConfig: adding missing text library:", tl.Name)
			cfg.TextLibraries = append(cfg.TextLibraries, tl)
			changed = true
		}
	}

	// --- Inherent Images --------------------------------------------------
	canonImgs := canonicalImages()
	existingImgs := make(map[string]bool)
	for _, img := range cfg.Images {
		existingImgs[img.Name] = true
	}
	for _, img := range canonImgs {
		if !existingImgs[img.Name] {
			fmt.Println("MigrateConfig: adding missing image source:", img.Name)
			cfg.Images = append(cfg.Images, img)
			changed = true
		}
	}
	changed = syncCanonicalImageMetadata(cfg.Images, canonImgs) || changed

	// --- Retired image sources --------------------------------------------
	// Remove sources that have been discontinued so existing users stop
	// seeing them. Unsplash was dropped because it competes with its own
	// wallpaper product.
	if removed := removeRetiredImageSources(cfg); removed {
		changed = true
	}

	// Repair stale embedded image locations from old build-cache binaries.
	for idx, img := range cfg.Images {
		if img.Inherent && (img.Name == "Christian" || img.Name == "Judaism") {
			expectedDir := filepath.Join(GetFolderPath(enum.PathLoc.Executable), "shared", "static", "images")
			expectedFolder := filepath.Join(expectedDir, img.Name+"PD")
			if _, err := os.Stat(expectedFolder); err == nil {
				if img.Location != expectedFolder {
					fmt.Println("MigrateConfig: repairing image location for", img.Name)
					cfg.Images[idx].Location = expectedFolder
					changed = true
				}
			}
		}
	}

	// --- Inherent filter keys ---------------------------------------------
	// Boolean filter flags added in newer versions are absent from configs
	// written by older builds. An absent JSON key unmarshals to false, which
	// is the correct default, but the key never gets written back unless we
	// force a save. Without the key persisted, the settings UI cannot show or
	// toggle the new filter. Detect any missing filter keys in the raw config
	// file and flag a re-save so SaveConfig serializes the full struct.
	if filterKeysMissingFromDisk() {
		fmt.Println("MigrateConfig: persisting newly added wallpaper filter keys")
		changed = true
	}

	// --- Version stamp ----------------------------------------------------
	if cfg.Version != AppVersion {
		cfg.Version = AppVersion
		cfg.Published = PublishedOn
		changed = true
	}

	return changed
}

// knownFilterKeys lists every wallpaper filter JSON key the current schema
// serializes. Keep this in sync with the WallpaperFilter* fields on Config.
var knownFilterKeys = []string{
	"wallpaperFilterOriginal",
	"wallpaperFilterBlurSoft",
	"wallpaperFilterBlurHard",
	"wallpaperFilterPixelate",
	"wallpaperFilterOilify",
	"wallpaperFilterWavy",
	"wallpaperFilterVortex",
	"wallpaperFilterMosaic",
	"wallpaperFilterJigsawPuzzle",
	"wallpaperFilterGraffiti",
	"wallpaperFilterCartoon",
	"wallpaperFilterCyberpunk",
	"wallpaperFilterMonochrome",
}

// filterKeysMissingFromDisk reports whether the on-disk config file is missing
// any known filter key. Booleans cannot distinguish "absent" from "false"
// after unmarshalling into the struct, so this inspects the raw JSON instead.
func filterKeysMissingFromDisk() bool {
	configPath := GetFolderPath(enum.PathLoc.ConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		// If the file is unreadable, let other migration steps / save handle it.
		return false
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return false
	}
	for _, key := range knownFilterKeys {
		if _, ok := raw[key]; !ok {
			return true
		}
	}
	return false
}

// retiredImageSourceNames lists image source names that have been removed from
// the product. MigrateConfig strips these from existing users' configs.
var retiredImageSourceNames = []string{"UnSplash"}

// removeRetiredImageSources deletes any retired image sources from the config's
// Images slice. Returns true if anything was removed.
func removeRetiredImageSources(cfg *Config) bool {
	retired := make(map[string]bool, len(retiredImageSourceNames))
	for _, name := range retiredImageSourceNames {
		retired[name] = true
	}
	filtered := cfg.Images[:0]
	removed := false
	for _, img := range cfg.Images {
		if retired[img.Name] {
			fmt.Println("MigrateConfig: removing retired image source:", img.Name)
			removed = true
			continue
		}
		filtered = append(filtered, img)
	}
	cfg.Images = filtered
	return removed
}

// canonicalTextLibraries returns the full set of inherent text libraries
// that a fresh install would have. Keep this in sync with CreateConfig.
func canonicalTextLibraries() []TextLibrary {
	return []TextLibrary{
		{Use: true, Name: "BibleVerses", Title: "King James Bible Verses", Location: "quotes/biblekjv.json", Citation: "https://aruljohn.com/Bible/", Creators: "Arul John", Info: "The King James Bible", Inherent: true},
		{Use: true, Name: "MBC Values", Title: "Manchester Baptist Church Core Values", Location: "quotes/mbc.json", Citation: "https://www.manchesterbaptist.org/", Creators: "MBC", Info: "Manchester Baptist", Inherent: true},
		{Use: true, Name: "AugustineQuotes", Title: "Augustine Quotes", Location: "quotes/augustine.json", Citation: "https://gracequotes.org/author-quote/augustine/", Creators: "Grace Quotes", Info: "Grace Quotes", Inherent: true},
		{Use: true, Name: "CharlesSpurgeonQuotes", Title: "Charles Spurgeon Quotes", Location: "quotes/charlesSpurgeon.json", Citation: "https://www.goodreads.com/author/quotes/2876959.Charles_Haddon_Spurgeon", Creators: "GoodReads Quotes", Info: "Goodreads", Inherent: true},
		{Use: true, Name: "RichardBaxterQuotes", Title: "Richard Baxter Quotes", Location: "quotes/richardBaxter.json", Citation: "https://gracequotes.org/author-quote/richard-baxter/", Creators: "Grace Quotes", Info: "Grace Quotes", Inherent: true},
		{Use: true, Name: "JohnCalvinQuotes", Title: "John Calvin Quotes", Location: "quotes/johnCalvin.json", Citation: "https://gracequotes.org/author-quote/john-calvin/", Creators: "Grace Quotes", Info: "Grace Quotes", Inherent: true},
		{Use: true, Name: "CSLewisQuotes", Title: "C.S. Lewis Quotes", Location: "quotes/csLewis.json", Citation: "https://gracequotes.org/author-quote/c-s-lewis/", Creators: "Grace Quotes", Info: "Grace Quotes", Inherent: true},
		{Use: true, Name: "MartinLutherQuotes", Title: "Martin Luther Quotes", Location: "quotes/martinLuther.json", Citation: "https://gracequotes.org/author-quote/martin-luther/", Creators: "Grace Quotes", Info: "Grace Quotes", Inherent: true},
		{Use: true, Name: "ChristianInspirations", Title: "Christian Inspirations", Location: "quotes/inspirations.json", Citation: "????", Creators: "Multiple", Info: "Multiple Sources", Inherent: true},
		{Use: true, Name: "TalmudQuotes", Title: "Talmud Quotes", Location: "quotes/21TalmudQuotes.json", Citation: "https://www.chabad.org", Creators: "Multiple", Info: "Multiple Sources", Inherent: true},
		{Use: false, Name: "GeneralMacArthurQuotes", Title: "General Douglas MacArthur Quotes", Location: "quotes/macarthur.json", Citation: "https://www.goodreads.com/author/quotes/317613.Douglas_MacArthur", Creators: "GoodReads.com", Info: "GoodReads", Inherent: true},
		{Use: false, Name: "GeneralPattonQuotes", Title: "General George S. Patton Quotes", Location: "quotes/patton.json", Citation: "https://www.wearethemighty.com/lists/general-george-patton-quotes/", Creators: "We Are The Mighty", Info: "We Are The Mighty", Inherent: true},
		{Use: false, Name: "MarkTwainQuotes", Title: "Quotes by Samuel Clemens (Mark Twain)", Location: "quotes/markTwain.json", Citation: "https://parade.com/1216401/jessicasager/mark-twain-quotes/", Creators: "Parade", Info: "Parade", Inherent: true},
		{Use: false, Name: "WillRogers", Title: "Will Rogers Quotes", Location: "quotes/willRogers.json", Citation: "https://www.willrogers.com/quotes", Creators: "Will Rogers Memorial Museum", Info: "Will Rogers Memorial Museum", Inherent: true},
		{Use: false, Name: "SenatorKennedy", Title: "Senator John Kennedy Quotes", Location: "quotes/senatorKennedy.json", Citation: "https://burningforsuccess.com/senator-john-kennedy-funny-quotes/", Creators: "burningforsuccess.com", Info: "Multiple Sources", Inherent: true},
		{Use: false, Name: "DatabaseQuotes", Title: "5000+ Famous Quotes", Location: "quotes/JamesFTquotes.json", Citation: "https://github.com/JamesFT/Database-Quotes-JSON", Creators: "James F Thompson (JamesFT)", Info: "Database Quotes JSON", Inherent: true},
		{Use: false, Name: "CelebrityQuotes", Title: "Celebrity Quotes", Location: "quotes/NasrulHazimQuotes.json", Citation: "https://gist.github.com/nasrulhazim/54b659e43b1035215cd0ba1d4577ee80", Creators: "Nasrul Hazim", Info: "Celebrity Quotes", Inherent: true},
		{Use: false, Name: "CallOfDuty", Title: "Quoted sayings in the Call of Duty series", Location: "quotes/callOfDuty.json", Citation: "https://callofduty.fandom.com/wiki/Quoted_sayings_in_the_Call_of_Duty_series", Creators: "Fandom", Info: "Fandom", Inherent: true},
	}
}

// canonicalImages returns the full set of inherent image sources.
// Keep this in sync with CreateConfig.
func canonicalImages() []Image {
	wallpaperDir := GetFolderPath(enum.PathLoc.Pictures)
	wallpaperFavs := GetFolderPath(enum.PathLoc.Favorites)
	wallpaperFS := GetFolderPath(enum.PathLoc.Executable)
	wallpaperChristian := filepath.Join(wallpaperFS, "shared", "static", "images", "ChristianPD")
	wallpaperJudaism := filepath.Join(wallpaperFS, "shared", "static", "images", "JudaismPD")
	return []Image{
		{Use: false, Name: "Favorites", Title: "Favorites", Location: wallpaperFavs, Operation: "Folder", AllowDistort: true, Inherent: true},
		{Use: true, Name: "Christian", Title: "Public Domain Christian Images", Location: wallpaperChristian, Operation: "Folder", AllowDistort: false, Inherent: true},
		{Use: true, Name: "Judaism", Title: "Public Domain Judaism Images", Location: wallpaperJudaism, Operation: "Folder", AllowDistort: false, Inherent: true},
		{Use: false, Name: "Bing", Title: "Bing Photo of the Day", Location: "https://bing.gifposter.com", Operation: "Webpage", AllowDistort: true, Inherent: true},
		{Use: false, Name: "Flickr", Title: "DR Flickr Photos", Location: "https://www.flickr.com/photos/202229109@N02", Operation: "WebPicPage", AllowDistort: true, Inherent: true},
		{Use: false, Name: "NASA", Title: "NASA's Astronomy Random Picture of the Day", Location: "https://apod.nasa.gov/apod/random_apod.html", Operation: "Webpage", AllowDistort: true, Inherent: true},
		{Use: false, Name: "PicSum", Title: "Pictures from PicSum random photos API", Location: "https://picsum.photos/1920/1080", Operation: "WebPicPage", AllowDistort: true, Inherent: true},
		{Use: false, Name: "Pexels", Title: "Photos from Pexels", Location: "https://www.pexels.com", Operation: "API", AllowDistort: true, RequiresKey: true, Inherent: true},
		{Use: true, Name: "WallpapersLocal", Title: "Wallpapers", Location: wallpaperDir, Operation: "Folder", AllowDistort: true, Inherent: true},
	}
}

func findCanonicalImage(images []Image, name string) (Image, bool) {
	for _, image := range images {
		if image.Name == name {
			return image, true
		}
	}
	return Image{}, false
}

func syncCanonicalImageMetadata(images []Image, canonicalImages []Image) bool {
	changed := false
	for idx, image := range images {
		canonical, ok := findCanonicalImage(canonicalImages, image.Name)
		if !ok {
			continue
		}
		if image.RequiresKey != canonical.RequiresKey {
			images[idx].RequiresKey = canonical.RequiresKey
			changed = true
		}
		if image.Inherent && (image.Title != canonical.Title || image.Location != canonical.Location || image.Operation != canonical.Operation) {
			images[idx].Title = canonical.Title
			images[idx].Location = canonical.Location
			images[idx].Operation = canonical.Operation
			changed = true
		}
	}
	return changed
}

// SaveConfig writes the configuration to the JSON file
// SaveConfig would likely need to write back to the file if you make changes.
func SaveConfig(cfg *Config) error {
	configMu.Lock()
	defer configMu.Unlock()
	if cfg == nil {
		cfg = ConfigInstance
	}
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}
	cfg.Version = AppVersion
	cfg.Published = PublishedOn
	loadedConfig = cfg
	ConfigInstance = cfg
	return saveConfigUnlocked(cfg)
}

func saveConfigUnlocked(cfg *Config) error {
	configPath := GetFolderPath(enum.PathLoc.ConfigFile)
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	configDir := filepath.Dir(configPath)
	tempFile, err := os.CreateTemp(configDir, "config-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp config file: %w", err)
	}
	tempPath := tempFile.Name()
	if _, err = tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to write temp config file: %w", err)
	}
	if err = tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to close temp config file: %w", err)
	}
	if err = os.Rename(tempPath, configPath); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to replace config file: %w", err)
	}
	if err = os.Chmod(configPath, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func CreateConfig() (*Config, error) {
	wallpaperDir := GetFolderPath(enum.PathLoc.Pictures)
	wallpaperFavs := GetFolderPath(enum.PathLoc.Favorites) //filep@th.Join(usr.HomeDir, ".Metamorphoun", "Favorites")
	wallpaperFS := GetFolderPath(enum.PathLoc.Executable)  //filep@th.Join(exeDir, "static", "images")
	_ = wallpaperFS
	wallpaperChristian := filepath.Join(wallpaperFS, "shared", "static", "images", "ChristianPD")
	wallpaperJudaism := filepath.Join(wallpaperFS, "shared", "static", "images", "JudaismPD")
	//staticWallpaperDir := shared.GetStaticImages()
	cfg := Config{
		Version:                         AppVersion,
		Published:                       PublishedOn,
		ServerAddress:                   "localhost",
		ServerPort:                      3000,
		SourceCurrentBackgroundName:     "",
		SourceCurrentBackgroundFolder:   "",
		OriginalCurrentBackgroundName:   "",
		OriginalCurrentBackgroundFolder: "",
		CurrentBackgroundName:           "",
		CurrentBackgroundFolder:         "",
		BackgroundChangingBlock:         false,
		StartOnStartup:                  true,
		ChangeWallpaperOnStartup:        true,
		DifferentWallpaperPerScreen:     false,
		ChangeMinutes:                   15,
		Images: []Image{
			{
				Use:          false,
				Name:         "Favorites",
				Title:        "Favorites",
				Location:     wallpaperFavs,
				Operation:    "Folder",
				AllowDistort: true,
				Inherent:     true,
			},
			{
				Use:          true,
				Name:         "Christian",
				Title:        "Public Domain Christian Images",
				Location:     wallpaperChristian,
				Operation:    "Folder",
				AllowDistort: false,
				Inherent:     true,
			},
			{
				Use:          true,
				Name:         "Judaism",
				Title:        "Public Domain Judaism Images",
				Location:     wallpaperJudaism,
				Operation:    "Folder",
				AllowDistort: false,
				Inherent:     true,
			},
			{
				Use:          false,
				Name:         "Bing",
				Title:        "Bing Photo of the Day",
				Location:     "https://bing.gifposter.com",
				Operation:    "Webpage",
				AllowDistort: true,
				Inherent:     true,
			},
			{
				Use:          false,
				Name:         "Flickr",
				Title:        "DR Flickr Photos",
				Location:     "https://www.flickr.com/photos/202229109@N02",
				Operation:    "WebPicPage",
				AllowDistort: true,
				Inherent:     true,
			},
			{
				Use:          false,
				Name:         "NASA",
				Title:        "NASA's Astronomy Random Picture of the Day",
				Location:     "https://apod.nasa.gov/apod/random_apod.html",
				Operation:    "Webpage",
				AllowDistort: true,
				Inherent:     true,
			},
			{
				Use:          false,
				Name:         "PicSum",
				Title:        "Pictures from PicSum random photos API",
				Location:     "https://picsum.photos/1920/1080",
				Operation:    "WebPicPage",
				AllowDistort: true,
				Inherent:     true,
			},
			{
				Use:          false,
				Name:         "Pexels",
				Title:        "Photos from Pexels",
				Location:     "https://www.pexels.com",
				Operation:    "API",
				AllowDistort: true,
				RequiresKey:  true,
				Inherent:     true,
			},
			{
				Use:          true,
				Name:         "WallpapersLocal",
				Title:        "Wallpapers",
				Location:     wallpaperDir,
				Operation:    "Folder",
				AllowDistort: true,
				Inherent:     true,
			},
		},
		ShowTextOverlay:   false,
		TextChangeMinutes: 5,
		// TextFontPath:              "C:\\Windows\\Fonts\\",
		TextFontFile:              "DejaVuSans-Bold.ttf",
		TextBoxLocation:           "TopRight",
		WallpaperImageSizing:      "",
		WallpaperFilterOriginal:   true,
		WallpaperFilterBlurSoft:   false,
		WallpaperFilterBlurHard:   false,
		WallpaperFilterPixelate:   false,
		WallpaperFilterOilify:     false,
		WallpaperFilterWavy:       false,
		WallpaperFilterVortex:     false,
		WallpaperFilterGraffiti:   false,
		WallpaperFilterCyberpunk:  false,
		WallpaperFilterMonochrome: false,
		QuoteAppearanceRandom:     false,
		QuoteFontRandom:           false,
		QuoteTextColor:            "#FFFFFF",
		QuoteBackgroundColor:      "#000000",
		QuoteBackgroundOpacity:    "110",
		TextLibraries: []TextLibrary{
			{
				Use:      true,
				Name:     "BibleVerses",
				Title:    "King James Bible Verses",
				Location: "quotes/biblekjv.json",
				Citation: "https://aruljohn.com/Bible/",
				Creators: "Arul John",
				Info:     "The King James Bible",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "MBC Values",
				Title:    "Manchester Baptist Church Core Values",
				Location: "quotes/mbc.json",
				Citation: "https://www.manchesterbaptist.org/",
				Creators: "MBC",
				Info:     "Manchester Baptist",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "AugustineQuotes",
				Title:    "Augustine Quotes",
				Location: "quotes/augustine.json",
				Citation: "https://gracequotes.org/author-quote/augustine/",
				Creators: "Grace Quotes",
				Info:     "‘Grace Quotes’ is a growing database containing over 10,000 great Christian quotes arranged over hundreds of topics. The material is from theologically sound, well-respected pastors, authors and Christian heroes from across the centuries.",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "CharlesSpurgeonQuotes",
				Title:    "Charles Spurgeon Quotes",
				Location: "quotes/charlesSpurgeon.json",
				Citation: "https://www.goodreads.com/author/quotes/2876959.Charles_Haddon_Spurgeon",
				Creators: "GoodReads Quotes",
				Info:     "Goodreads is the world’s largest site for readers and book recommendations. Our mission is to help readers discover books they love and get more out of reading. Goodreads launched in January 2007",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "RichardBaxterQuotes",
				Title:    "Richard Baxter Quotes",
				Location: "quotes/richardBaxter.json",
				Citation: "https://gracequotes.org/author-quote/richard-baxter/",
				Creators: "Grace Quotes",
				Info:     "‘Grace Quotes’ is a growing database containing over 10,000 great Christian quotes arranged over hundreds of topics. The material is from theologically sound, well-respected pastors, authors and Christian heroes from across the centuries.",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "JohnCalvinQuotes",
				Title:    "John Calvin Quotes",
				Location: "quotes/johnCalvin.json",
				Citation: "https://gracequotes.org/author-quote/john-calvin/",
				Creators: "Grace Quotes",
				Info:     "‘Grace Quotes’ is a growing database containing over 10,000 great Christian quotes arranged over hundreds of topics. The material is from theologically sound, well-respected pastors, authors and Christian heroes from across the centuries.",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "CSLewisQuotes",
				Title:    "C.S. Lewis Quotes",
				Location: "quotes/csLewis.json",
				Citation: "https://gracequotes.org/author-quote/c-s-lewis/",
				Creators: "Grace Quotes",
				Info:     "‘Grace Quotes’ is a growing database containing over 10,000 great Christian quotes arranged over hundreds of topics. The material is from theologically sound, well-respected pastors, authors and Christian heroes from across the centuries.",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "MartinLutherQuotes",
				Title:    "Martin Luther Quotes",
				Location: "quotes/martinLuther.json",
				Citation: "https://gracequotes.org/author-quote/martin-luther/",
				Creators: "Grace Quotes",
				Info:     "‘Grace Quotes’ is a growing database containing over 10,000 great Christian quotes arranged over hundreds of topics. The material is from theologically sound, well-respected pastors, authors and Christian heroes from across the centuries.",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "ChristianInspirations",
				Title:    "Christian Inspirations",
				Location: "quotes/inspirations.json",
				Citation: "????",
				Creators: "Multiple",
				Info:     "Multiple Sources",
				Inherent: true,
			},
			{
				Use:      true,
				Name:     "TalmudQuotes",
				Title:    "Talmud Quotes",
				Location: "quotes/21TalmudQuotes.json",
				Citation: "https://www.chabad.org",
				Creators: "Multiple",
				Info:     "Multiple Sources",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "GeneralMacArthurQuotes",
				Title:    "General Douglas MacArthur Quotes",
				Location: "quotes/macarthur.json",
				Citation: "https://www.goodreads.com/author/quotes/317613.Douglas_MacArthur",
				Creators: "GoodReads.com",
				Info:     "The right book in the right hands at the right time can change the world. Who We Are Goodreads is the world’s largest site for readers and book recommendations. Our mission is to help readers discover books they love and get more out of reading. Goodreads launched in January 2007.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "GeneralPattonQuotes",
				Title:    "General George S. Patton Quotes",
				Location: "quotes/patton.json",
				Citation: "https://www.wearethemighty.com/lists/general-george-patton-quotes/",
				Creators: "We Are The Mighty",
				Info:     "We Are The Mighty is a veteran-led digital publisher and Emmy Award-winning media agency servicing brands with video production, marketing, advertising, and consulting services to engage with the military community. In addition to our digital publisher, we also run the Military Influencer Conference, the largest in-person event servicing our military community. WATM is owned by Recurrent Ventures and is a GSA-approved vendor.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "MarkTwainQuotes",
				Title:    "Quotes by Samuel Clemens (Mark Twain)",
				Location: "quotes/markTwain.json",
				Citation: "https://parade.com/1216401/jessicasager/mark-twain-quotes/",
				Creators: "Parade",
				Info:     "The Parade brand has been delighting, enlightening and inspiring readers since it was founded in 1941. Through our access to A-list celebrities, top experts and today’s most intriguing and influential personalities, our team provides information, solutions, perspectives and advice on trending topics in entertainment, pop culture and lifestyle. We give you reasons to feel good about your life and the world around you through the stories we tell.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "WillRogers",
				Title:    "Will Rogers Quotes",
				Location: "quotes/willRogers.json",
				Citation: "https://www.willrogers.com/quotes",
				Creators: "Will Rogers Memorial Museum",
				Info:     "The Will Rogers Memorial Museum is a 19,052-square-foot museum in Claremore, Oklahoma that memorializes entertainer Will Rogers. The museum houses artifacts, memorabilia, photographs, and manuscripts pertaining to Rogers' life, and documentaries, speeches, and movies starring Rogers are shown in a theater. The museum is one of five attractions operated by the Will Rogers Memorial Museums, Inc., a non-profit organization.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "DatabaseQuotes",
				Title:    "5000+ Famous Quotes",
				Location: "quotes/JamesFTquotes.json",
				Citation: "https://github.com/JamesFT/Database-Quotes-JSON",
				Creators: "James F Thompson (JamesFT)",
				Info:     "#Database Quotes JSON ##JSON file with more than 5000+ famous quotes. Some example on how to work on this JSON quotes file",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "CelebrityQuotes",
				Title:    "Celebrity Quotes",
				Location: "quotes/NasrulHazimQuotes.json",
				Citation: "https://gist.github.com/nasrulhazim/54b659e43b1035215cd0ba1d4577ee80",
				Creators: "Nasrul Hazim",
				Info:     "The Parade brand has been delighting, enlightening and inspiring readers since it was founded in 1941. Through our access to A-list celebrities, top experts and today’s most intriguing and influential personalities, our team provides information, solutions, perspectives and advice on trending topics in entertainment, pop culture and lifestyle. We give you reasons to feel good about your life and the world around you through the stories we tell.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "CallOfDuty",
				Title:    "Quoted sayings in the Call of Duty series",
				Location: "quotes/callOfDuty.json",
				Citation: "https://callofduty.fandom.com/wiki/Quoted_sayings_in_the_Call_of_Duty_series",
				Creators: "Fandom",
				Info:     "Our Mission -- We power fan experiences.  Our mission is to understand, inform, entertain, and celebrate fans by building the best entertainment and gaming communities, content, services, and experiences.",
				Inherent: true,
			},
		},
		MBCMonth:         0,     //set to current when MBCMode is enabled
		MBCMode:          false, //When On the MBC traits will replace quotes
		MBCValue:         0,
		QuoteFontSizeMin: 16,
		QuoteFontSizeMax: 28,
		PicHistories:     []PicHistory{},
	}

	// Get the user's home directory
	//println(usr.Username)
	configPath := GetFolderPath(enum.PathLoc.ConfigFile)
	// Create the config directory if it doesn't exist
	err := os.MkdirAll(GetFolderPath(enum.PathLoc.Config), 0700) // Adjust permissions as needed
	if err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(GetFolderPath(enum.PathLoc.Favorites), 0700) // Adjust permissions as needed
	if err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal the config struct to JSON
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write the JSON data to the file
	//err = os.WriteFile(configPath, data, 0600) // Adjust permissions as needed
	err = os.WriteFile(configPath, data, 0600) // Adjust permissions as needed
	if err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	configMu.Lock()
	loadedConfig = &cfg
	ConfigInstance = &cfg
	configMu.Unlock()
	return &cfg, nil
}

func SetupSystemFolders() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	metamorphounDirs := []string{"Favorites", "Logs"}
	for _, fldr := range metamorphounDirs {
		folderPath := filepath.Join(usr.HomeDir, ".Metamorphoun", fldr)

		_, err := os.Stat(folderPath)
		if os.IsNotExist(err) {
			fmt.Println("Folder does not exist.")
			err = os.MkdirAll(folderPath, 0755) // Adjust permissions as needed
			if err != nil {
				fmt.Println("failed to create config directory:", err)
			}
			if fldr == "Quotes" {
				//copy in common quotes
				//simple
				exePath, err := os.Executable()
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				// Get the directory containing the executable
				exeDir := filepath.Dir(exePath)
				fmt.Println("Executable Path:", exePath)
				fmt.Println("Executable Directory:", exeDir)

				appFolder := filepath.Join(exeDir, "static", "quotes")
				appFile := filepath.Join(appFolder, "simple.json")
				userFolder := folderPath
				userFileMMDir := filepath.Join(userFolder, ".Metamorphoun", "Quotes", "simple.json")
				err1 := zutil.CopyFile(appFile, userFileMMDir)
				if err1 != nil {
					fmt.Println("Error copying file:", err1)
				} else {
					fmt.Println("File copied successfully!")
				}
			}
		} else if err != nil {
			fmt.Println("Error checking folder:", err)
		} else {
			fmt.Println("Folder exists.")
		}
	}
	//add favorites subfolders
	wallpaperFavs := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites")

	err = os.MkdirAll(filepath.Join(wallpaperFavs, "Pictures", "WithQuotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory:", err)
	}
	err = os.MkdirAll(filepath.Join(wallpaperFavs, "Pictures", "WithOutQuotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory:", err)

		err = os.MkdirAll(filepath.Join(wallpaperFavs, "Quotes"), 0700) // Adjust permissions as needed
		if err != nil {
			fmt.Println("failed to create config directory:", err)
		}
	}
}
