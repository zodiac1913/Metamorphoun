package service

import (
	"Metamorphoun/config"
	"fmt"
	"image"
	"math/rand"
	"net/http"

	"golang.org/x/net/html"
)

//New Way

func GetBackgroundUnSplash(imgItem config.Image) (image.Image, string, error) {
	//=====================================================Get Random Image from page
	wppArray, wppErr := extractUnSplashImageURLs(imgItem)
	if wppErr != nil {
		fmt.Println("Error: getting web urls on ", imgItem.Location, " for ", imgItem.Operation, " on ", wppErr.Error())
	}
	//choose image to use
	if len(wppArray) < 1 {
		fmt.Println("Error: No img links found on page ", imgItem.Location, " for ", imgItem.Operation, " on ", wppErr.Error())
		return nil, "", nil
	}
	wppRnd := rand.Intn(len(wppArray))
	pic := wppArray[wppRnd]

	println("Unsplash Pic:" + pic)
	img, err := loadFlickrImageFromURL(pic)
	if err != nil {
		fmt.Println("failed to fetch image from URL: %w", err)
		return nil, "", err
	}
	return img, pic, nil
}

func extractUnSplashImageURLs(imgItem config.Image) ([]string, error) {
	url := imgItem.Location
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var imageURLs []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			var src, itemprop string
			for _, a := range n.Attr {
				if a.Key == "src" {
					src = a.Val
				}
				if a.Key == "itemprop" {
					itemprop = a.Val
				}
			}
			if src != "" && itemprop != "" {
				imageURLs = append(imageURLs, src)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return imageURLs, nil
}

//-------------------------------------------- OLD WAY

// func ChangeBackgroundUnSplash(imgItem config.Image) error {
// 	background := ""
// 	wppArray, wppErr := extractUnSplashImageURLs(imgItem)
// 	if wppErr != nil {
// 		fmt.Println("Error: getting web urls on ", imgItem.Location, " for ", imgItem.Operation, " on ", wppErr.Error())
// 	}
// 	//choose image to use
// 	if len(wppArray) < 1 {
// 		fmt.Println("Error: No img links found on page ", imgItem.Location, " for ", imgItem.Operation, " on ", wppErr.Error())
// 		return nil
// 	}
// 	wppRnd := rand.Intn(len(wppArray))
// 	pic := wppArray[wppRnd]
// 	usr, err := user.Current()
// 	if err != nil {
// 		fmt.Println("failed to get user home directory:", err)
// 	}
// 	wppFolder := filepath.Join(usr.HomeDir, ".Metamorphoun", imgItem.Operation)
// 	background = downloadUnSplashPic(pic, imgItem.Operation+"Choice", wppFolder, imgItem)
// 	if background != "" {
// 		fmt.Println("Background downloaded")
// 	} else {
// 		fmt.Println("failed to download pic. Unknown Reason", err)
// 	}
// 	backgroundFile := filepath.Base(background)
// 	ffFolder := filepath.Dir(background)
// 	ext := filepath.Ext(filepath.Base(background))
// 	bufferFile := strings.Replace(backgroundFile, ext, "Buffer"+ext, -1)

// 	//filepath.Join(usr.HomeDir, ".Metamorphoun", imgItem.Operation)
// 	fmt.Println("Setting as Background")

// 	config.UpdateConfigField("sourceCurrentBackgroundName", bufferFile)
// 	config.UpdateConfigField("sourceCurrentBackgroundFolder", ffFolder)
// 	config.UpdateConfigField("originalCurrentBackgroundName", bufferFile)
// 	config.UpdateConfigField("originalCurrentBackgroundFolder", ffFolder)
// 	config.UpdateConfigField("currentBackgroundName", backgroundFile)
// 	config.UpdateConfigField("currentBackgroundFolder", ffFolder)
// 	lEntry := morphLog.LogItem{TimeStamp: time.Now().Format("20060102 15:04:05"),
// 		Message: "Changed Background", Level: "INFO", Library: imgItem.Location,
// 		Operation: imgItem.Operation, Origin: pic, LocalFile: filepath.Join(ffFolder, bufferFile),
// 	}
// 	morphLog.UpdateLogs(lEntry)
// 	fmt.Println("new log entry:", lEntry)
// 	//fmt.Println("Todays Logs:", logs)
// 	ChangeBackgroundRoutine()

// 	return nil
// }

// func downloadUnSplashPic(url string, fileName string, folder string, img config.Image) string {
// 	//flickr issue
// 	urlFiltered := ""
// 	ext := ""
// 	urlFiltered = url
// 	ext = "jpg"

// 	// Send a GET request to the URL
// 	resp, err := http.Get(urlFiltered)
// 	if err != nil {
// 		log.Printf("Error downloading image: %v", err)
// 		return ""
// 	}

// 	// Save the image to disk
// 	parts := strings.Split(urlFiltered, ".")
// 	if ext == "" {
// 		ext = strings.Replace(strings.Split(parts[len(parts)-1], "?")[0], "_sm", "", 0)
// 		ext = strings.Replace(ext, "_mb", "", 0)
// 	}
// 	//fname := folder + "\\" + fileName + "." + ext
// 	usr, err := user.Current()
// 	if err != nil {
// 		fmt.Println("failed to get user home directory:", err)
// 	}
// 	bufferFile := filepath.Join(usr.HomeDir, ".Metamorphoun", img.Operation, fileName+"Buffer."+ext)
// 	fname := filepath.Join(usr.HomeDir, ".Metamorphoun", img.Operation, fileName+"."+ext)
// 	//get rid of the old
// 	DeleteFile(bufferFile)
// 	DeleteFile(fname)

// 	//MAKE NEW FILE
// 	f, err := os.Create(bufferFile)
// 	if err != nil {
// 		log.Println("Error could not create file:", bufferFile, " Error:", err)
// 		return ""
// 	}
// 	defer f.Close()

// 	_, err = io.Copy(f, resp.Body)
// 	if err != nil {
// 		log.Printf("Error copying data: %v", err)
// 		return ""
// 	}
// 	copyFile(bufferFile, fname)
// 	return fname
// }
