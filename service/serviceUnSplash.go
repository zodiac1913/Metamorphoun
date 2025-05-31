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
