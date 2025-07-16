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
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
)

var updateSignal chan struct{}

// var fontFldrs = []string{
// 	"/usr/share/fonts",
// 	"/usr/local/share/fonts",
// 	"~/.local/share/fonts",
// 	"~/.fonts",
// 	"C:\\Windows\\Fonts",
// }

func main() {
	//top!!!
	// var ff []string
	// ff = append(ff, `/System/Library/Fonts/`)
	// service.RenderFontsSample(ff)
	//service.ChangeLockScreen = changeLockScreenImpl
	config.GetFolderPath = getFolderPathImpl
	morphLog.GetFolderPath = getFolderPathImpl
	service.GetFolderPath = getFolderPathImpl
	server.GetFolderPath = getFolderPathImpl
	systemTray.GetFolderPath = getFolderPathImpl
	service.SetRandomQuote = setRandomQuoteImpl
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
	config.ConfigInstance = configData
	//top!!!
	// Common logic
	PrintPlatformMessage()

	//service.ChangeView()
	//return
	//Load config file
	//Used to test fonts

	cfg := configData // Now cfg points to the single loaded instance

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

	// Start the server in a separate goroutine
	go func() {
		if !server.Serve(*cfg) { // Pass the config pointer
			println("Server failed to start")
		}
	}()
	println("start background change from main")

	// // Start the server in a separate goroutine
	// go func() {
	// 	if !server.Serve(configData.ServerAddress, configData.ServerPort) {
	// 		println("Server failed to start")
	// 	}
	// }() // Start background changing
	// println("start background change from main")

	// if cfg.ChangeWallpaperOnStartup {
	// 	//service.ChangeBackground()
	// 	//service.ChangeView("backgroundChange")
	// 	pic := config.PicHistory{}
	// 	service.BackgroundGenerate("ChangeOnStartup", pic)
	// }
	println("start background change from main")

	if cfg.ChangeWallpaperOnStartup {
		pic := config.PicHistory{}
		config.ConfigInstance.BackgroundChangeAttempt = 0
		service.BackgroundGenerate("ChangeOnStartup", pic)
	}

	// go func() {
	// 	timer := time.NewTicker(time.Duration(configData.ChangeMinutes) * time.Minute)
	// 	pic := config.PicHistory{}
	// 	defer timer.Stop()
	// 	for {
	// 		select {
	// 		case <-timer.C:
	// 			//service.ChangeView("backgroundChange")
	// 			service.BackgroundGenerate("ChangeOnStartup", pic)
	// 		case <-updateSignal:
	// 			//service.ChangeView("backgroundChange")
	// 			service.BackgroundGenerate("ChangeOnStartup", pic)
	// 			timer.Reset(time.Duration(configData.ChangeMinutes) * time.Minute)
	// 		case <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

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

	//quote service
	//quotes.SetQuote()
	if cfg.ShowTextOverlay {
		go func() {
			serveQuotes := service.StartChangeQuote(time.Duration(cfg.TextChangeMinutes) * time.Minute)
			println("quotes started 0 and ", cfg.TextChangeMinutes, " min timer")
			if err := serveQuotes.Start(); err != nil {
				println(err)
			}
		}()
	}
	// if cfg.ShowTextOverlay {
	// 	go func() {
	// 		serveQuotes := service.StartChangeQuote(time.Duration(configData.TextChangeMinutes) * time.Minute)
	// 		println("quotes started 0 and ", configData.TextChangeMinutes, " min timer")
	// 		// Start the service
	// 		if err := serveQuotes.Start(); err != nil {
	// 			println(err)
	// 		}

	// 	}()

	// }

	// if runtime.GOOS == "windows" {

	// 	// System tray onExit function
	// 	onExit := func() {
	// 		now := time.Now()
	// 		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	// 		// Signal all goroutines to stop
	// 		cancel()
	// 	}

	// 	systray.Run(systemTray.MakeSystemTray, onExit)
	// 	// Start the service
	// 	// Prevent the main function from exiting
	// 	<-ctx.Done()
	// } else {
	// 	//Linux
	// 	go func() {
	// 		linuxGui.MakeGui()
	// 	}()
	// }
	//	service.ChangeLockScreen(config.ConfigInstance.PicHistories[0]) // Initialize the ChangeLockScreen function
	//if runtime.GOOS == "windows" {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		cancel()
	}
	systray.Run(systemTray.MakeSystemTray, onExit)
	<-ctx.Done()
	//Perhaps check if the systray fails via err and then run the gui
	//} else {
	//	go func() {
	//		linuxGui.MakeGui()
	//	}()
	//}
}
func openFolder(title string, path string) error {
	var cmd *exec.Cmd
	cmd = exec.Command(title, path)
	return cmd.Start()
}

//	func GetFolderPath(pathNeeded string) string {
//		return GetFolderPath(pathNeeded)
//	}
func getFolderPathImpl(pathNeeded string) string {
	return GetFolderPath(pathNeeded)
}

func setRandomQuoteImpl(currentPic config.PicHistory, img image.Image) (config.PicHistory, image.Image, error) {
	return SetRandomQuote(currentPic, img)
}

// func changeLockScreenImpl(pic config.PicHistory) error {
// 	return ChangeLockScreen(pic)
// }
