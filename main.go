//main.go

package main

import (
	"Metamorphoun/config"
	"Metamorphoun/morphLog"
	"Metamorphoun/server"
	"Metamorphoun/service"
	"Metamorphoun/systemTray"
	"context"
	"fmt"
	"image"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

var updateSignal chan struct{}

func main() {
	config.GetFolderPath = getFolderPathImpl
	morphLog.GetFolderPath = getFolderPathImpl
	service.GetFolderPath = getFolderPathImpl
	server.GetFolderPath = getFolderPathImpl
	systemTray.GetFolderPath = getFolderPathImpl
	service.SetRandomQuote = setRandomQuoteImpl
	cfg := loadOrCreateConfig()
	normalizeConfigDefaults(cfg)
	PrintPlatformMessage()

	config.SetupSystemFolders()
	fmt.Println("Server Address:", cfg.ServerAddress)
	println("Server (in main)")
	println(cfg.ServerAddress)
	println("port")
	println(cfg.ServerPort)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the update signal channel
	updateSignal = make(chan struct{})

	startServer(cfg)

	settingsURL := "http://localhost:" + strconv.Itoa(config.ConfigInstance.ServerPort)
	if !waitForSettingsServer(settingsURL, 10*time.Second) {
		fmt.Println("Settings server was not ready; browser launch skipped")
	} else if !server.HasRecentHeartbeat(10 * time.Second) {
		fmt.Println("No browser heartbeat detected — opening browser tab")
		if err := server.OpenFolder("explorer", settingsURL); err != nil {
			fmt.Println("Failed to open settings page:", err)
		}
	} else {
		fmt.Println("Browser heartbeat detected — tab already open, skipping browser open")
	}
	startBackgroundServices(ctx, cfg)

	onExit := func() {
		now := time.Now()
		os.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		cancel()
	}
	systray.Run(systemTray.MakeSystemTray, onExit)
	<-ctx.Done()
}

func waitForSettingsServer(settingsURL string, timeout time.Duration) bool {
	client := http.Client{Timeout: time.Second}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		response, err := client.Get(settingsURL)
		if err == nil {
			response.Body.Close()
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func loadOrCreateConfig() *config.Config {
	configData, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config not found, creating new config:", err)
		configData, err = config.CreateConfig()
		if err != nil {
			fmt.Println("Error creating config:", err)
			panic("Failed to create config")
		}
	}

	config.ConfigInstance = configData
	if config.MigrateConfig(config.ConfigInstance) {
		config.SaveConfig(config.ConfigInstance)
		fmt.Println("Config migrated to", config.AppVersion)
	}

	return configData
}

func normalizeConfigDefaults(cfg *config.Config) {
	saveNeeded := false

	if cfg.MBCMonth < 1 || cfg.MBCMonth > 12 {
		cfg.MBCMonth = int(time.Now().Month())
		saveNeeded = true
	}
	if cfg.QuoteFontSizeMin < 8 {
		cfg.QuoteFontSizeMin = 16
		saveNeeded = true
	}
	if cfg.QuoteFontSizeMax < cfg.QuoteFontSizeMin {
		cfg.QuoteFontSizeMax = 28
		saveNeeded = true
	}

	if saveNeeded {
		config.SaveConfig(cfg)
	}
}

func startBackgroundServices(ctx context.Context, cfg *config.Config) {
	startWallpaperScheduler(ctx, cfg)
	startQuoteService(cfg)
}

func startServer(cfg *config.Config) {
	go func() {
		if !server.Serve(*cfg) {
			println("Server failed to start")
		}
	}()
	println("start background change from main")
	println("start background change from main")
}

func startWallpaperScheduler(ctx context.Context, cfg *config.Config) {
	if cfg.ChangeWallpaperOnStartup {
		pic := config.PicHistory{}
		config.ConfigInstance.BackgroundChangeAttempt = 0
		service.BackgroundGenerate("ChangeOnStartup", pic)
	}

	go func() {
		timer := time.NewTicker(time.Duration(cfg.ChangeMinutes) * time.Minute)
		pic := config.PicHistory{}
		config.ConfigInstance.BackgroundChangeAttempt = 0
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				service.BackgroundGenerate("ChangeOnStartup", pic)
			case <-updateSignal:
				service.BackgroundGenerate("ChangeOnStartup", pic)
				timer.Reset(time.Duration(cfg.ChangeMinutes) * time.Minute)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func startQuoteService(cfg *config.Config) {
	if cfg.ShowTextOverlay && !cfg.MBCMode {
		go func() {
			serveQuotes := service.StartChangeQuote(time.Duration(cfg.TextChangeMinutes) * time.Minute)
			println("quotes started 0 and ", cfg.TextChangeMinutes, " min timer")
			if err := serveQuotes.Start(); err != nil {
				println(err)
			}
		}()
	}
}

func openFolder(title string, path string) error {
	var cmd *exec.Cmd
	cmd = exec.Command(title, path)
	return cmd.Start()
}

func getFolderPathImpl(pathNeeded string) string {
	return GetFolderPath(pathNeeded)
}

func setRandomQuoteImpl(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
	return SetRandomQuote(currentPic, img)
}
