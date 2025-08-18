package service

import (
	"Metamorphoun/config"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

//New Way

func GetBackgroundBing(imgItem config.Image) (image.Image, string, error) {
	//=====================================================Get Random Image from page
	wppArray, wppErr := extractBingImageURLs(imgItem)
	if wppErr != nil {
		fmt.Println("Error: getting web urls on ", imgItem.Location, " for ", imgItem.Operation, " on ", wppErr.Error())
	}
	//choose image to use
	if len(wppArray) < 1 {
		//fmt.Println("Error: No img links found on page ", imgItem.Location, " for ", imgItem.Operation, " on ", wppErr.Error())
		return nil, "", nil
	}
	wppRnd := rand.Intn(len(wppArray))
	pic := wppArray[wppRnd]

	urlFiltered := ""
	//ext := ""
	urlFiltered1 := strings.Replace(pic, "jpg_sm", "jpg", -1)
	urlFiltered = strings.Replace(urlFiltered1, "jpg_mb", "jpg", -1)

	println("Bing Pic:" + urlFiltered)
	img, err := loadBingImageFromURL(urlFiltered)
	if err != nil {
		fmt.Println("failed to fetch image from URL: %w", err)
		return nil, "", err
	}
	return img, urlFiltered, nil
}

func loadBingImageFromURL(url string) (image.Image, error) {
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
func extractBingImageURLs(imgItem config.Image) ([]string, error) {
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
