package systemTray

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/icon"
	"Metamorphoun/server"
	"Metamorphoun/service"
	"time"

	"Metamorphoun/zutil"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/getlantern/systray"
	//"github.com/getlantern/systray/example/icon"
)

var GetFolderPath func(string) string

type PathLocType string

func MakeSystemTray() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Background Fun")
	systray.SetTooltip("Metamorphoun")
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	now := time.Now()
	dt := now.Format("20060102_150405")
	wallpaperMain := GetFolderPath(enum.PathLoc.Config)
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
		mFavStoreWQ := favPicsMenu.AddSubMenuItem("Store With Quote", "Store this pic with the quote that is on it")
		mFavStoreNQ := favPicsMenu.AddSubMenuItem("Store Without Quote", "Store this pic without the quote that is on it")

		//subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
		//subMenuBottom := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
		//subMenuBottom2 := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")

		mUrl := systray.AddMenuItem("Settings", "Configure your Metamorphoun")
		mNextBG := systray.AddMenuItem("Next Background", "Change to next background image")
		mLastBG := systray.AddMenuItem("Last Background", "Change to the last background image")
		//mShowCurrentPicture := systray.AddMenuItem("Current Info", "Show current picture information")
		systray.AddSeparator()
		mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

		// Sets the icon of a menu item. Only available on Mac.
		//mQuit.SetIcon(icon.Data)

		//systray.AddSeparator()
		//mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
		go func() {
			<-mQuitOrig.ClickedCh
			fmt.Println("Requesting quit")
			systray.Quit()
			fmt.Println("Finished quitting")
		}()

		//shown := true
		// toggle := func() {
		// 	if shown {
		// 		//subMenuBottom.Check()
		// 		//subMenuBottom2.Hide()
		// 		mQuitOrig.Hide()
		// 		//mEnabled.Hide()
		// 		shown = false
		// 	} else {
		// 		//subMenuBottom.Uncheck()
		// 		//subMenuBottom2.Show()
		// 		mQuitOrig.Show()
		// 		//mEnabled.Show()
		// 		shown = true
		// 	}
		// }

		for {
			select {
			case <-mFavStoreWQ.ClickedCh:
				currImgWQ := config.ConfigInstance.PicHistories[0]
				fmt.Print("Current Image with Quote: ", currImgWQ.OriginName)
				wqExt := filepath.Ext(currImgWQ.OriginName)
				if len(wqExt) > 5 {
					wqExt = service.UnUnsplash(currImgWQ.OriginName)
				}
				currentPicFile := filepath.Join(wallpaperMain, "pic0"+wqExt)
				picToSave := filepath.Join(favPicFolderWithQuote, dt+wqExt)
				zutil.CopyFile(currentPicFile, picToSave)
				server.OpenFolder("explorer", favPicFolderWithQuote)
			case <-mFavStoreNQ.ClickedCh:
				service.RecallBackground("SystrayFavStoreNQ", 0)
				time.Sleep(15 * time.Second)
				server.OpenFolder("explorer", favPicFolderWithoutQuote)
			case <-mNextBG.ClickedCh:
				service.BackgroundGenerate("SystrayNextBackground", config.PicHistory{})
			case <-mLastBG.ClickedCh:
				service.RecallBackground("RecallBackground", 1)
			// case <-mShowCurrentPicture.ClickedCh:
			// 	currPicInfo := "http://" + config.ConfigInstance.ServerAddress + ":" + zutil.AsString(config.ConfigInstance.ServerPort) + "/picInfo.html"
			// 	server.OpenFolder("explorer", currPicInfo)

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
				//case <-mToggle.ClickedCh:
				//	toggle()
				//case <-mQuit.ClickedCh:
				//	systray.Quit()
				//	fmt.Println("Quit2 now...")
				//	return
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
