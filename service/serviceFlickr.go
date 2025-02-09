package service

import (
	"MorphPrototype/config"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

//New Way

func GetBackgroundFlickr(imgItem config.Image) (image.Image, string, error) {
	//=====================================================Get Random Image from page
	wppArray, wppErr := extractFlickrImageURLs(imgItem)
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

	urlFiltered := ""
	//ext := ""
	urlFiltered1 := strings.Replace(pic, "jpg_sm", "jpg", -1)
	urlFiltered = strings.Replace(urlFiltered1, "jpg_mb", "jpg", -1)
	urlFiltered = strings.Replace(urlFiltered, "//live.", "https://live.", -1)

	println("Flickr Pic:" + urlFiltered)
	img, err := loadFlickrImageFromURL(urlFiltered)
	if err != nil {
		fmt.Println("failed to fetch image from URL: %w", err)
		return nil, "", err
	}
	return img, urlFiltered, nil
}

func loadFlickrImageFromURL(url string) (image.Image, error) {
	// Fetch the image from the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()
	// Decode the image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

// Get urls from the page for pics
func extractFlickrImageURLs(imgItem config.Image) ([]string, error) {
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
			pattern := "(?i)logo"
			regex := regexp.MustCompile(pattern) //Get rid of natgeos goofy logos
			for _, a := range n.Attr {
				if a.Key == "src" {
					if !strings.HasPrefix(a.Val, "data:") &&
						!strings.HasPrefix(a.Val, "__PUBLIC__") &&
						!regex.MatchString(a.Val) {
						imageURLs = append(imageURLs, a.Val)
						break
					}
				}
			}
		}
		// Check if the child node exists before iterating
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)

	return imageURLs, nil
}

// -------------------------------------------Old Way

// func ChangeBackgroundFlickr(imgItem config.Image) error {
// 	background := ""
// 	wppArray, wppErr := extractFlickrImageURLs(imgItem)
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
// 	println("Bing Pic:" + pic)
// 	usr, err := user.Current()
// 	if err != nil {
// 		fmt.Println("failed to get user home directory:", err)
// 	}
// 	wppFolder := filepath.Join(usr.HomeDir, ".Metamorphoun", imgItem.Operation)
// 	background = downloadFlickrPic(pic, imgItem.Operation+"Choice", wppFolder, imgItem)
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

// func downloadFlickrPic(url string, fileName string, folder string, img config.Image) string {
// 	//flickr issue
// 	urlFiltered := ""
// 	ext := ""
// 	urlFiltered = strings.Replace(url, "//", "http://", 1)
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
