package systemTray

import (
	"MorphPrototype/config"
	"MorphPrototype/icon"
	"MorphPrototype/server"
	"MorphPrototype/service"
	"MorphPrototype/zutil"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/getlantern/systray"
	//"github.com/getlantern/systray/example/icon"
)

func MakeSystemTray() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Background Fun")
	systray.SetTooltip("Metamorphoun")
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	makeFavFolders()
	favPicFolderWithQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithQuotes")
	favPicFolderWithoutQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithOutQuotes")

	//favQuotesFolder := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Quotes")

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Background Fun")
		systray.SetTooltip("Metamorphoun")
		// mChange := systray.AddMenuItem("Change Me", "Change Me")
		// mChecked := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
		// mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		// Sets the icon of a menu item. Only available on Mac.
		//mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		//systray.AddMenuItem("Ignored", "Ignored")

		favPicsMenu := systray.AddMenuItem("Favorite (Pics)", "Store Favorite Pictures")
		mStoreWQ := favPicsMenu.AddSubMenuItem("Store With Quote", "Store this pic with the quote that is on it")
		mStoreNQ := favPicsMenu.AddSubMenuItem("Store Without Quote", "Store this pic without the quote that is on it")

		//subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
		//subMenuBottom := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
		//subMenuBottom2 := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")

		mUrl := systray.AddMenuItem("Settings", "Configure your Metamorphoun")
		mNextBG := systray.AddMenuItem("Next Background", "Change to next background image")
		mShowCurrentPicture := systray.AddMenuItem("Current Info", "Show current picture information")
		mQuit := systray.AddMenuItem("Quit", "Shutdown Metamorphoun")

		// Sets the icon of a menu item. Only available on Mac.
		//mQuit.SetIcon(icon.Data)

		systray.AddSeparator()
		mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
		mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
		go func() {
			<-mQuitOrig.ClickedCh
			fmt.Println("Requesting quit")
			systray.Quit()
			fmt.Println("Finished quitting")
		}()

		shown := true
		toggle := func() {
			if shown {
				//subMenuBottom.Check()
				//subMenuBottom2.Hide()
				mQuitOrig.Hide()
				//mEnabled.Hide()
				shown = false
			} else {
				//subMenuBottom.Uncheck()
				//subMenuBottom2.Show()
				mQuitOrig.Show()
				//mEnabled.Show()
				shown = true
			}
		}

		for {
			select {
			case <-mStoreWQ.ClickedCh:
				currentPicFile := filepath.Join(config.ConfigInstance.CurrentBackgroundFolder, config.ConfigInstance.CurrentBackgroundName)
				picToSave := filepath.Join(favPicFolderWithQuote, config.ConfigInstance.SourceCurrentBackgroundName)
				zutil.CopyFile(currentPicFile, picToSave)
				server.OpenFolder("explorer", favPicFolderWithQuote)
			case <-mStoreNQ.ClickedCh:
				currentPicFile := filepath.Join(config.ConfigInstance.OriginalCurrentBackgroundFolder, config.ConfigInstance.OriginalCurrentBackgroundName)
				picToSave := filepath.Join(favPicFolderWithoutQuote, config.ConfigInstance.SourceCurrentBackgroundName)
				zutil.CopyFile(currentPicFile, picToSave)
				server.OpenFolder("explorer", favPicFolderWithoutQuote)
			case <-mNextBG.ClickedCh:
				service.ChangeView()
			case <-mShowCurrentPicture.ClickedCh:
				currPicInfo := "http://" + config.ConfigInstance.ServerAddress + ":" + zutil.AsString(config.ConfigInstance.ServerPort) + "/picInfo.html"
				server.OpenFolder("explorer", currPicInfo)

			// case <-mChange.ClickedCh:
			// 	mChange.SetTitle("I've Changed")
			// case <-mChecked.ClickedCh:
			// 	if mChecked.Checked() {
			// 		mChecked.Uncheck()
			// 		mChecked.SetTitle("Unchecked")
			// 	} else {
			// 		mChecked.Check()
			// 		mChecked.SetTitle("Checked")
			// 	}
			// case <-mEnabled.ClickedCh:
			// 	mEnabled.SetTitle("Disabled")
			// 	mEnabled.Disable()
			case <-mUrl.ClickedCh:
				//open.Run("https://www.getlantern.org")
				urlSettings := "http://" + config.ConfigInstance.ServerAddress + ":" + zutil.AsString(config.ConfigInstance.ServerPort)
				server.OpenFolder("explorer", urlSettings)
			//case <-subMenuBottom2.ClickedCh:
			//	panic("panic button pressed")
			//case <-subMenuBottom.ClickedCh:
			//	toggle()
			case <-mToggle.ClickedCh:
				toggle()
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}

func makeFavFolders() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	//wallpaperFavs := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites")

	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithQuotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithOutQuotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
	err = os.MkdirAll(filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Quotes"), 0700) // Adjust permissions as needed
	if err != nil {
		fmt.Println("failed to create config directory: %w", err)
	}
}
