package service

import (
	"Metamorphoun/config"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

//New Way

func GetBackgroundNASA(imgItem config.Image) (image.Image, string, error) {
	//=====================================================Get Random Image from page
	background, year, month, day := calcDate()
	fmt.Println(background)
	// Generate final HTML string
	finalHTML := "https://apod.nasa.gov/apod/ap" + ("" + year) + month + day + ".html"
	//fmt.Sprintf("<iframe src=\"https://apod.nasa.gov/apod/ap%s%s%s.html\" id=\"apod\"></iframe>", year, month, day)

	fmt.Println(finalHTML)

	wppArray, wppErr := extractNASAImageURLs(finalHTML)
	if wppErr != nil {
		fmt.Println("Error: getting web urls on ", imgItem.Location, " for ", imgItem.Operation, " on ", wppErr.Error())
	}
	// //choose image to use
	if len(wppArray) < 1 {
		fmt.Println("Error: No img links found on page ")
		return nil, "", fmt.Errorf("no image links found on page")
	}
	pic := wppArray[0]
	if !strings.HasPrefix(pic, "http") {
		pic = "https://apod.nasa.gov/apod/" + pic
	}
	bkgd, err := loadNASAImageFromURL(pic)
	if err != nil {
		fmt.Println("failed to fetch image from URL: %w", err)
		return nil, "", err
	}

	return bkgd, pic, nil
}

func loadNASAImageFromURL(url string) (image.Image, error) {
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
func extractNASAImageURLs(url string) ([]string, error) {
	// Step 1: Extract iframe src from the main page
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
			for _, a := range n.Attr {
				if a.Key == "src" {
					if !strings.HasPrefix(a.Val, "data:") &&
						!strings.HasPrefix(a.Val, "__PUBLIC__") {
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

func calcDate() (string, string, string, string) {
	background := ""

	now := time.Now()

	// Minimum date: 1995-06-16 (adjusted for UTC)
	min := time.Date(1995, 6, 16, 0, 0, 0, 0, time.UTC)

	// Maximum date: Today (adjusted for UTC and 18:59:59.999)
	max := time.Date(now.Year(), now.Month(), now.Day(), 18, 59, 59, 999*int(time.Millisecond), time.UTC)

	// Adjust max for Eastern Time (subtract 5 hours)
	max = max.Add(-5 * time.Hour)

	// Generate random date within range
	randomDate := min.Add(time.Duration(rand.Int63n(int64(max.Sub(min)))))

	// Missing date range (1995-06-17 to 1995-06-19)
	missingMin := time.Date(1995, 6, 17, 0, 0, 0, 0, time.UTC)
	missingMax := time.Date(1995, 6, 19, 23, 59, 59, 999*int(time.Millisecond), time.UTC)

	// Regenerate random date if it falls in the missing range
	for randomDate.After(orEqual(missingMin, randomDate)) && randomDate.Before(orEqual(randomDate, missingMax)) {
		randomDate = min.Add(time.Duration(rand.Int63n(int64(max.Sub(min)))))
	}
	// Format date components with zero-padding
	year := strconv.Itoa(randomDate.Year() % 100)
	month := strconv.Itoa(int(randomDate.Month())) // Convert month to string
	day := strconv.Itoa(randomDate.Day())          // Convert day to string

	// Zero-pad month and day if necessary
	if len(month) == 1 {
		month = "0" + month
	}
	if len(day) == 1 {
		day = "0" + day
	}
	return background, year, month, day
}

// Helper function to handle time comparisons where either operand might be nil
func orEqual(t1, t2 time.Time) time.Time {
	if t1.IsZero() {
		return t2
	}
	return t1
}
