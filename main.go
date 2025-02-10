package main

import (
	"Metamorphoun/config"
	"Metamorphoun/server"
	"Metamorphoun/service"
	"Metamorphoun/systemTray"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
)

var updateSignal chan struct{}

func main() {

	//service.ChangeView()
	//return
	//Load config file
	//Used to test fonts
	// usr, err := user.Current()
	// if err != nil {
	// 	fmt.Println("failed to get user home directory: %w", err)
	// }
	// outputPath := filepath.Join(usr.HomeDir, ".Metamorphoun", "fonts_sample.png")
	// errf := service.RenderFontsSample(outputPath)
	// if errf != nil {
	// 	fmt.Println("Error rendering fonts sample:", errf)
	// } else {
	// 	fmt.Println("Fonts sample image created successfully:", outputPath)
	// }
	configData, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		config.CreateConfig()
		configData, err = config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config complete failure:", err)
			panic("Bad")
		}
	}
	cfg := config.GetConfig()
	config.SetupSystemFolders()
	fmt.Println("Server Address:", cfg.ServerAddress)
	println("Server (in main)")
	println(configData.ServerAddress)
	println("port")
	println(configData.ServerPort)
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the update signal channel
	updateSignal = make(chan struct{})

	// Start the server in a separate goroutine
	go func() {
		if !server.Serve(configData.ServerAddress, configData.ServerPort) {
			println("Server failed to start")
		}
	}() // Start background changing
	println("start background change from main")

	if cfg.ChangeWallpaperOnStartup {
		//service.ChangeBackground()
		service.ChangeView()
	}

	go func() {
		timer := time.NewTicker(time.Duration(configData.ChangeMinutes) * time.Minute)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				service.ChangeView()
			case <-updateSignal:
				service.ChangeView()
				timer.Reset(time.Duration(configData.ChangeMinutes) * time.Minute)
			case <-ctx.Done():
				return
			}
		}
	}()

	//quote service
	//quotes.SetQuote()
	if cfg.ShowTextOverlay {
		go func() {
			serveQuotes := service.StartChangeQuote(time.Duration(configData.TextChangeMinutes) * time.Minute)
			println("quotes started 0 and ", configData.ChangeMinutes, " min timer")
			// Start the service
			if err := serveQuotes.Start(); err != nil {
				println(err)
			}

		}()

	}

	// System tray onExit function
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		// Signal all goroutines to stop
		cancel()
	}

	systray.Run(systemTray.MakeSystemTray, onExit)
	// Start the service
	// Prevent the main function from exiting
	<-ctx.Done()

}
func openFolder(title string, path string) error {
	var cmd *exec.Cmd
	cmd = exec.Command(title, path)
	return cmd.Start()
}
