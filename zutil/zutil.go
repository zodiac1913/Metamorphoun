// Package used for common Dom utilities
package zutil

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func AsString(input interface{}) string {
	return fmt.Sprintf("%v", input)
}

func AsBool(input string) bool {
	switch input {
	case "true", "yes", "on", "1":
		return true
	case "false", "no", "off", "0":
		return false
	default:
		return false
	}
}

func AsInt(input string, ifNullValue ...int) int {
	nullVal := -1
	if len(ifNullValue) > 0 {
		nullVal = ifNullValue[0]
	}
	intValue, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Error converting string to int")
		return nullVal
	}
	return intValue
}

func AsInt16(input string, ifNullValue ...int16) int16 {
	var nullVal int16 = -1
	if len(ifNullValue) > 0 {
		int16val, erri1 := strconv.ParseInt(input, 10, 16)
		if erri1 != nil {
			fmt.Println("Error converting string to int:", erri1)
			return int16(-1)
		}
		nullVal = int16(int16val)
	}
	intValue, err := strconv.ParseInt(input, 10, 16)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return nullVal
	}
	// Convert int64 to int16 (check for potential overflow)
	if intValue < math.MinInt16 || intValue > math.MaxInt16 {
		fmt.Println("Value exceeds the range of int16")
		return nullVal
	}
	intValue16 := int16(intValue)
	return intValue16
}

func AsInt32(input string, ifNullValue ...int32) int32 {
	var nullVal int32 = -1
	if len(ifNullValue) > 0 {
		int16val, erri32 := strconv.ParseInt(input, 10, 16)
		if erri32 != nil {
			fmt.Println("Error converting string to int:", erri32)
			return int32(-1)
		}
		nullVal = int32(int16val)
	}
	intValue, err := strconv.ParseInt(input, 10, 16)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return nullVal
	}
	// Convert int64 to int16 (check for potential overflow)
	if intValue < math.MinInt32 || intValue > math.MaxInt32 {
		fmt.Println("Value exceeds the range of int16")
		return nullVal
	}
	intValue32 := int32(intValue)
	return intValue32
}

func AsInt64(input string, ifNullValue int64) int64 {
	if ifNullValue == 0 {
		ifNullValue = -1
	}
	intValue, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return ifNullValue
	}
	return intValue
}

// copyFile copies a file from src to dst. If dst does not exist, it will be created.
func CopyFile(src, dst string) error {
	// Open the source file
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	// Create the destination file
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	// Copy the contents from source to destination
	_, err = io.Copy(destination, source)
	return err
}

func IsInRange(value string, rangeOfStrings []string) bool {
	for _, str := range rangeOfStrings {
		if value == str {
			return true
		}
	}
	return false
}

func LoadImageFromURL(url string) (image.Image, error) {
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

// loadImg loads an image from the given file path (supports PNG and JPEG)
func LoadImg(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if strings.HasSuffix(strings.ToLower(path), ".png") {
		return png.Decode(f)
	}
	return jpeg.Decode(f)
}
