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
	"sync"
)

const AppVersion = "2025.6.25"
const PublishedOn = "2025-06-25"

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
	//WallpaperFilterSpiral           bool          `json:"wallpaperFilterSpiral"`
	WallpaperFilterMonochrome bool          `json:"wallpaperFilterMonochrome"`
	ShowTextOverlay           bool          `json:"showTextOverlay"`
	TextChangeMinutes         int           `json:"textChangeMinutes"`
	TextLibraries             []TextLibrary `json:"textLibraries"`
	TextFontFile              string        `json:"textFontFile"`
	//TextFontPath              string        `json:"textFontPath"`
	TextBoxLocation        string       `json:"textBoxLocation"`
	CurrentQuoteStatement  string       `json:"currentQuoteStatement"`
	CurrentQuoteAuthor     string       `json:"currentQuoteAuthor"`
	QuoteAppearanceRandom  bool         `json:"quoteAppearanceRandom"`
	QuoteFontRandom        bool         `json:"quoteFontRandom"`
	QuoteTextColor         string       `json:"quoteTextColor"`
	QuoteBackgroundColor   string       `json:"quoteBackgroundColor"`
	QuoteBackgroundOpacity string       `json:"quoteBackgroundOpacity"`
	PicHistories           []PicHistory `json:"picHistories"`
	PicUpdateCalled        bool         `json:"picUpdateCalled"`
	// Add other configuration fields here
}
type Image struct {
	Use       bool   `json:"use"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Location  string `json:"location"`
	Operation string `json:"operation"`
	Inherent  bool   `json:"inherent"` // Indicates if the image is inherent to the system
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

var (
	loadedConfig *Config
	loadOnce     sync.Once
	loadError    error
)

func init() {
	// Load the configuration
}

// GetConfig returns the current Config instance
func GetConfig() *Config {
	ConfigInstance.Version = AppVersion
	ConfigInstance.Published = PublishedOn
	if loadedConfig == nil {
		// Handle the case where loading failed, perhaps return a default or panic
		fmt.Println("Warning: Config not loaded yet. Call LoadConfig first.")
		return &Config{} // Return a default empty config to avoid nil pointer
	}
	return loadedConfig
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
	return *ConfigInstance
}

// SetConfig updates the Config instance and saves it to the file
func SetConfig(newConfig *Config) error {
	ConfigInstance = newConfig
	return SaveConfig(newConfig)
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
	//fmt.Println("UpdateConfigField:")
	//fmt.Println(propertyName)
	//fmt.Println(newValue)
	ConfigInstance = GetConfig()
	//fmt.Println("Config-BEFORE")
	//fmt.Println(ConfigInstance)
	switch propertyName {
	case "serverAddress":
		ConfigInstance.ServerAddress = newValue.(string)
	case "serverPort":
		toString := fmt.Sprintf("%v", newValue)
		ConfigInstance.ServerPort = zutil.AsInt(toString)
	case "startOnStartup":
		boolValue := zutil.AsBool(fmt.Sprintf("%v", newValue))
		ConfigInstance.StartOnStartup = boolValue
		fmt.Println("StartOnStartup-SET")
		if boolValue {
			err := AddToStartup()
			if err != nil {
				log.Println("Error adding to startup:", err)
			}
		} else {
			err := RemoveFromStartup()
			if err != nil {
				log.Println("Error adding to startup:", err)
			}
		}
	case "changeWallpaperOnStartup":
		boolValue := zutil.AsBool(fmt.Sprintf("%v", newValue))
		ConfigInstance.ChangeWallpaperOnStartup = boolValue
		fmt.Println("StartOnStartup-SET")
	case "changeMinutes":
		intValue := zutil.AsInt(fmt.Sprintf("%v", newValue))
		ConfigInstance.ChangeMinutes = int32(intValue)
	case "sourceCurrentBackgroundName":
		ConfigInstance.SourceCurrentBackgroundName = newValue.(string)
	case "sourceCurrentBackgroundFolder":
		ConfigInstance.SourceCurrentBackgroundFolder = newValue.(string)
	case "originalCurrentBackgroundName":
		ConfigInstance.OriginalCurrentBackgroundName = newValue.(string)
	case "originalCurrentBackgroundFolder":
		ConfigInstance.OriginalCurrentBackgroundFolder = newValue.(string)
	case "currentBackgroundName":
		ConfigInstance.CurrentBackgroundName = newValue.(string)
	case "currentBackgroundFolder":
		ConfigInstance.CurrentBackgroundFolder = newValue.(string)
	case "backgroundChangingBlock":
		ConfigInstance.BackgroundChangingBlock = newValue.(bool)
	case "currentQuoteStatement":
		ConfigInstance.CurrentQuoteStatement = newValue.(string)
	case "currentQuoteAuthor":
		ConfigInstance.CurrentQuoteAuthor = newValue.(string)
	case "showTextOverlay":
		boolValue := zutil.AsBool(fmt.Sprintf("%v", newValue))
		ConfigInstance.ShowTextOverlay = boolValue
	case "textChangeMinutes":
		intValue := zutil.AsInt(fmt.Sprintf("%v", newValue))
		ConfigInstance.TextChangeMinutes = int(intValue)
	// case "textFontPath":
	// 	ConfigInstance.TextFontPath = newValue.(string)
	case "textFontFile":
		ConfigInstance.TextFontFile = newValue.(string)
	case "textBoxLocation":
		ConfigInstance.TextBoxLocation = newValue.(string)
	case "quoteAppearanceRandom":
		boolValue := zutil.AsBool(fmt.Sprintf("%v", newValue))
		ConfigInstance.QuoteAppearanceRandom = boolValue
	case "quoteFontRandom":
		boolValue := zutil.AsBool(fmt.Sprintf("%v", newValue))
		ConfigInstance.QuoteFontRandom = boolValue
	case "quoteTextColor":
		ConfigInstance.QuoteTextColor = newValue.(string)
	case "quoteBackgroundColor":
		ConfigInstance.QuoteBackgroundColor = newValue.(string)
	case "quoteBackgroundOpacity":
		ConfigInstance.QuoteBackgroundOpacity = newValue.(string)
	case "wallpaperImageSizing":
		ConfigInstance.WallpaperImageSizing = newValue.(string)
	case "wallpaperFilterOriginal":
		ConfigInstance.WallpaperFilterOriginal = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterBlurSoft":
		ConfigInstance.WallpaperFilterBlurSoft = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterBlurHard":
		ConfigInstance.WallpaperFilterBlurHard = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterPixelate":
		ConfigInstance.WallpaperFilterPixelate = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterOilify":
		ConfigInstance.WallpaperFilterOilify = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterWavy":
		ConfigInstance.WallpaperFilterWavy = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterVortex":
		ConfigInstance.WallpaperFilterVortex = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterMosaic":
		ConfigInstance.WallpaperFilterMosaic = zutil.AsBool(fmt.Sprintf("%v", newValue))
	case "wallpaperFilterMonochrome":
		ConfigInstance.WallpaperFilterMonochrome = zutil.AsBool(fmt.Sprintf("%v", newValue))
	default:
		fmt.Printf("invalid field name: %s", propertyName)
		return fmt.Errorf("invalid field name: %s", propertyName)
	}
	//fmt.Println("Config-AFTER")
	//fmt.Println(ConfigInstance)
	return SaveConfig(ConfigInstance)
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
	ConfigInstance = GetConfig()
	var foundImage *Image
	for i, image := range ConfigInstance.Images {
		if image.Name == imageName {
			foundImage = &ConfigInstance.Images[i] // Use pointer assignment
			break                                  // Exit the loop after finding the image
		}
	}
	foundImage.Use = newValue
	return SaveConfig(ConfigInstance)
}
func AddImagesField(use bool, name string, title string,
	location string, operation string) error {
	ConfigInstance = GetConfig()
	ConfigInstance.Images = append(ConfigInstance.Images, Image{
		Use:       use,
		Name:      name,
		Title:     title,
		Location:  location,
		Operation: operation,
	})
	return SaveConfig(ConfigInstance)
}
func EditImagesField(use bool, name string, title string,
	location string, operation string) error {
	ConfigInstance = GetConfig()
	cfg := GetImageByName(name)
	if cfg.Inherent {
		fmt.Println("Cannot edit inherent image:", name)
		return fmt.Errorf("cannot edit inherent image: %s", name)
	} else {
		cfg.Use = use
		cfg.Title = title
		cfg.Location = location
		cfg.Operation = operation
		return SaveConfig(ConfigInstance)
	}
}

func UpdateQuotesField(quotesName string, newValue interface{}) error {
	ConfigInstance = GetConfig()
	var foundQuotes *TextLibrary
	for i, textLib := range ConfigInstance.TextLibraries {
		if textLib.Name == quotesName {
			foundQuotes = &ConfigInstance.TextLibraries[i] // Use pointer assignment
			break                                          // Exit the loop after finding the image
		}
	}
	foundQuotes.Use = zutil.AsBool(newValue.(string))
	return SaveConfig(ConfigInstance)

}

// AddPicHistory adds a new PicHistory to the stack, updates PicNum,
// and ensures the stack size does not exceed the limit.
func (cfg *Config) AddPicHistory(newPic PicHistory) error {
	ConfigInstance = GetConfig()
	// Prepend the new PicHistory to the stack
	ConfigInstance.PicHistories = append([]PicHistory{newPic}, ConfigInstance.PicHistories...)

	// Ensure the stack size does not exceed the limit (5 for now)
	if len(ConfigInstance.PicHistories) > 5 {
		ConfigInstance.PicHistories = ConfigInstance.PicHistories[:5]
	}

	// Update PicNum for all PicHistories in the stack
	for i := range ConfigInstance.PicHistories {
		ConfigInstance.PicHistories[i].PicNum = int16(i)
	}
	return SaveConfig(ConfigInstance)
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
		loadedConfig = &config
	})
	return loadedConfig, loadError
}

// SaveConfig writes the configuration to the JSON file
// SaveConfig would likely need to write back to the file if you make changes.
func SaveConfig(cfg *Config) error {
	configPath := GetFolderPath(enum.PathLoc.ConfigFile)
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func CreateConfig() error {
	wallpaperDir := GetFolderPath(enum.PathLoc.Pictures)
	wallpaperFavs := GetFolderPath(enum.PathLoc.Favorites) //filep@th.Join(usr.HomeDir, ".Metamorphoun", "Favorites")
	wallpaperFS := GetFolderPath(enum.PathLoc.Executable)  //filep@th.Join(exeDir, "static", "images")
	//staticWallpaperDir := shared.GetStaticImages()
	cfg := Config{
		Version:                         AppVersion,
		Published:                       PublishedOn,
		ServerAddress:                   "127.0.0.1",
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
		ChangeMinutes:                   15,
		Images: []Image{
			{
				Use:       false,
				Name:      "Favorites",
				Title:     "Favorites",
				Location:  wallpaperFavs,
				Operation: "Folder",
				Inherent:  true,
			},
			{
				Use:       true,
				Name:      "PDChristianArt",
				Title:     "Christian Images",
				Location:  wallpaperFS,
				Operation: "Folder",
				Inherent:  true,
			},
			{
				Use:       false,
				Name:      "Bing",
				Title:     "Bing Photo of the Day",
				Location:  "https://bing.gifposter.com",
				Operation: "Webpage",
				Inherent:  true,
			},
			{
				Use:       false,
				Name:      "Flickr",
				Title:     "DR Flickr Photos",
				Location:  "https://www.flickr.com/photos/202229109@N02",
				Operation: "WebPicPage",
				Inherent:  true,
			},
			{
				Use:       false,
				Name:      "NASA",
				Title:     "NASA's Astronomy Random Picture of the Day",
				Location:  "https://apod.nasa.gov/apod/random_apod.html",
				Operation: "Webpage",
				Inherent:  true,
			},
			{
				Use:       false,
				Name:      "UnSplash",
				Title:     "Photos from Unsplash.com",
				Location:  "https://unsplash.com",
				Operation: "WebPicPage",
				Inherent:  true,
			},
			{
				Use:       false,
				Name:      "PicSum",
				Title:     "Pictures from PicSum random photos API",
				Location:  "https://picsum.photos/1920/1080",
				Operation: "WebPicPage",
				Inherent:  true,
			},
			{
				Use:       true,
				Name:      "WallpapersLocal",
				Title:     "Wallpapers",
				Location:  wallpaperDir,
				Operation: "Folder",
				Inherent:  true,
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
				Location: "/quotes/macarthur.json",
				Citation: "https://www.goodreads.com/author/quotes/317613.Douglas_MacArthur",
				Creators: "GoodReads.com",
				Info:     "The right book in the right hands at the right time can change the world. Who We Are Goodreads is the world’s largest site for readers and book recommendations. Our mission is to help readers discover books they love and get more out of reading. Goodreads launched in January 2007.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "GeneralPattonQuotes",
				Title:    "General George S. Patton Quotes",
				Location: "/quotes/patton.json",
				Citation: "https://www.wearethemighty.com/lists/general-george-patton-quotes/",
				Creators: "We Are The Mighty",
				Info:     "We Are The Mighty is a veteran-led digital publisher and Emmy Award-winning media agency servicing brands with video production, marketing, advertising, and consulting services to engage with the military community. In addition to our digital publisher, we also run the Military Influencer Conference, the largest in-person event servicing our military community. WATM is owned by Recurrent Ventures and is a GSA-approved vendor.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "MarkTwainQuotes",
				Title:    "Quotes by Samuel Clemens (Mark Twain)",
				Location: "/quotes/markTwain.json",
				Citation: "https://parade.com/1216401/jessicasager/mark-twain-quotes/",
				Creators: "Parade",
				Info:     "The Parade brand has been delighting, enlightening and inspiring readers since it was founded in 1941. Through our access to A-list celebrities, top experts and today’s most intriguing and influential personalities, our team provides information, solutions, perspectives and advice on trending topics in entertainment, pop culture and lifestyle. We give you reasons to feel good about your life and the world around you through the stories we tell.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "WillRogers",
				Title:    "Will Rogers Quotes",
				Location: "/quotes/willRogers.json",
				Citation: "https://www.willrogers.com/quotes",
				Creators: "Will Rogers Memorial Museum",
				Info:     "The Will Rogers Memorial Museum is a 19,052-square-foot museum in Claremore, Oklahoma that memorializes entertainer Will Rogers. The museum houses artifacts, memorabilia, photographs, and manuscripts pertaining to Rogers' life, and documentaries, speeches, and movies starring Rogers are shown in a theater. The museum is one of five attractions operated by the Will Rogers Memorial Museums, Inc., a non-profit organization.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "DatabaseQuotes",
				Title:    "5000+ Famous Quotes",
				Location: "/quotes/JamesFTquotes.json",
				Citation: "https://github.com/JamesFT/Database-Quotes-JSON",
				Creators: "James F Thompson (JamesFT)",
				Info:     "#Database Quotes JSON ##JSON file with more than 5000+ famous quotes. Some example on how to work on this JSON quotes file",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "CelebrityQuotes",
				Title:    "Celebrity Quotes",
				Location: "/quotes/NasrulHazimQuotes.json",
				Citation: "https://gist.github.com/nasrulhazim/54b659e43b1035215cd0ba1d4577ee80",
				Creators: "Nasrul Hazim",
				Info:     "The Parade brand has been delighting, enlightening and inspiring readers since it was founded in 1941. Through our access to A-list celebrities, top experts and today’s most intriguing and influential personalities, our team provides information, solutions, perspectives and advice on trending topics in entertainment, pop culture and lifestyle. We give you reasons to feel good about your life and the world around you through the stories we tell.",
				Inherent: true,
			},
			{
				Use:      false,
				Name:     "CallOfDuty",
				Title:    "Quoted sayings in the Call of Duty series",
				Location: "/quotes/callOfDuty.json",
				Citation: "https://callofduty.fandom.com/wiki/Quoted_sayings_in_the_Call_of_Duty_series",
				Creators: "Fandom",
				Info:     "Our Mission -- We power fan experiences.  Our mission is to understand, inform, entertain, and celebrate fans by building the best entertainment and gaming communities, content, services, and experiences.",
				Inherent: true,
			},
		},
		PicHistories: []PicHistory{},
	}

	// Get the user's home directory
	//println(usr.Username)
	configPath := GetFolderPath(enum.PathLoc.ConfigFile)
	// Create the config directory if it doesn't exist
	err := os.MkdirAll(GetFolderPath(enum.PathLoc.Config), 0700) // Adjust permissions as needed
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(GetFolderPath(enum.PathLoc.Favorites), 0700) // Adjust permissions as needed
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal the config struct to JSON
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write the JSON data to the file
	//err = ioutil.WriteFile(configPath, data, 0600) // Adjust permissions as needed
	err = os.WriteFile(configPath, data, 0600) // Adjust permissions as needed
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func SetupSystemFolders() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("failed to get user home directory: %w", err)
	}
	metamorphounDirs := []string{"Favorites", "Logs"}
	for _, fldr := range metamorphounDirs {
		folderPath := filepath.Join(usr.HomeDir, ".Metamorphoun", fldr)

		_, err := os.Stat(folderPath)
		if os.IsNotExist(err) {
			fmt.Println("Folder does not exist.")
			err = os.MkdirAll(folderPath, 0755) // Adjust permissions as needed
			if err != nil {
				fmt.Printf("failed to create config directory: %w", err)
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
		fmt.Println("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(filepath.Join(wallpaperFavs, "Pictures", "WithOutQuotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)

		err = os.MkdirAll(filepath.Join(wallpaperFavs, "Quotes"), 0700) // Adjust permissions as needed
		if err != nil {
			fmt.Println("failed to create config directory: %w", err)
		}
	}
}
