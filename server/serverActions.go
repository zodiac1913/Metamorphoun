package server

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/service"
	"Metamorphoun/zutil"
	"fmt"
	"io"
	"net/http"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	//"net/http/pprof"
	"os"

	"github.com/tidwall/gjson"
	// Import the pprof package explicitly
	// Import the pprof package explicitly
	//_ "net/http/pprof"
)

func lastBackgroundApi(w http.ResponseWriter, r *http.Request) {
	configPath := config.GetFolderPath(enum.PathLoc.ConfigFile)
	// Read config file
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		http.Error(w, "Failed to read config file", http.StatusInternalServerError)
		return
	}
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	service.RecallBackground("RecallBackground", 1)
	// Write JSON data to response
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func nextBackgroundApi(w http.ResponseWriter, r *http.Request) {
	configPath := config.GetFolderPath(enum.PathLoc.ConfigFile)
	// Read config file
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		http.Error(w, "Failed to read config file", http.StatusInternalServerError)
		return
	}
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	config.ConfigInstance.BackgroundChangeAttempt++
	service.BackgroundGenerate("WebServerNext", config.PicHistory{})
	// Write JSON data to response
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func saveFavoriteApi(w http.ResponseWriter, r *http.Request) {
	//configPath := config.GetFolderPath(enum.PathLoc.ConfigFile)
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}

	// Read config file
	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		http.Error(w, "Failed to read config file", http.StatusInternalServerError)
		return
	}
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	//service.SaveFavorite("SaveFavorite")
	//{type:"Quote","save":"quote"}
	//Convert Json to struct
	fmt.Println("formApi-Received JSON:", string(jsonData))
	typ := gjson.GetBytes(jsonData, "type").String()
	save := gjson.GetBytes(jsonData, "save").String()
	now := time.Now()
	dt := now.Format("20060102_150405")
	wallpaperMain := GetFolderPath(enum.PathLoc.Config)
	service.MakeFavFolders()
	favPicFolderWithQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithQuotes")
	favPicFolderWithoutQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithOutQuotes")

	if typ == "BG" {
		if save == "quoteOnBG" {
			currImgWQ := config.ConfigInstance.PicHistories[0]
			if strings.HasPrefix(currImgWQ.OriginName, GetFolderPath(enum.PathLoc.Favorites)) {
				fmt.Println("This picture is already in your favorites, no need to save it again.")
			} else {
				fmt.Print("Current Image with Quote: ", currImgWQ.OriginName)
				wqExt := filepath.Ext(currImgWQ.OriginName)
				if len(wqExt) > 5 {
					wqExt = service.UnUnsplash(currImgWQ.OriginName)
				}
				if len(wqExt) < 1 {
					wqExt = ".png"
				}
				currentPicFile := filepath.Join(wallpaperMain, "pic0"+wqExt)
				picToSave := filepath.Join(favPicFolderWithQuote, dt+wqExt)
				zutil.CopyFile(currentPicFile, picToSave)
				OpenFolder("explorer", favPicFolderWithQuote)
			}
		} else if save == "noQuoteOnBG" {
			service.RecallBackground("SystrayFavStoreNQ", 0)
			time.Sleep(5 * time.Second)
			OpenFolder("explorer", favPicFolderWithoutQuote)
		}
	} else { //Quote Only
		if save == "quote" {
			cq := config.ConfigInstance.PicHistories[0]
			service.StoreQuote("{\"statement\": \"" + cq.QuoteStatement + "\", \"author\": \"" + cq.QuoteAuthor + "\"}")
		}
	}

	// Write JSON data to response
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
