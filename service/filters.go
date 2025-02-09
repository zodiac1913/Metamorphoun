package service

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func BlurIt(img image.Image, blurSigma float64) (image.Image, error) {
	// Apply a Gaussian blur
	blurred := imaging.Blur(img, blurSigma) // The second parameter is the sigma of the Gaussian kernel
	return blurred, nil
}

func PixelateIt(img image.Image, pixelSize int) (image.Image, error) {
	if pixelSize == 0 {
		pixelSize = rand.Intn(16) + 4
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y += pixelSize {
		for x := 0; x < width; x += pixelSize {
			var r, g, b, a uint32
			var count int

			for dy := 0; dy < pixelSize; dy++ {
				for dx := 0; dx < pixelSize; dx++ {
					if y+dy < height && x+dx < width {
						r1, g1, b1, a1 := img.At(x+dx, y+dy).RGBA()
						r += r1
						g += g1
						b += b1
						a += a1
						count++
					}
				}
			}

			if count > 0 {
				r /= uint32(count)
				g /= uint32(count)
				b /= uint32(count)
				a /= uint32(count)
			}

			for dy := 0; dy < pixelSize; dy++ {
				for dx := 0; dx < pixelSize; dx++ {
					if y+dy < height && x+dx < width {
						newImg.Set(x+dx, y+dy, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
					}
				}
			}
		}
	}

	// usr, err := user.Current()
	// if err != nil {
	// 	fmt.Println("failed to get user home directory:", err)
	// }

	// currentPicsFolder := filepath.Join(usr.HomeDir, ".Metamorphoun")

	// fileStep2 := filepath.Join(currentPicsFolder, "file6APixelateFiltered.png")
	// saveImg(newImg, fileStep2)

	return newImg, nil
}

func OilifyIt(img image.Image, radius int) (image.Image, error) { //img image.Image, radius int) image.Image {
	if radius == 0 {
		radius = rand.Intn(6) + 3
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			hist := make(map[color.Color]int)
			var mostCommonColor color.Color
			var maxCount int

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					px := x + dx
					py := y + dy

					if px >= 0 && px < width && py >= 0 && py < height {
						c := img.At(px, py)
						hist[c]++
						if hist[c] > maxCount {
							maxCount = hist[c]
							mostCommonColor = c
						}
					}
				}
			}

			newImg.Set(x, y, mostCommonColor)
		}
	}
	return newImg, nil
}

func WavyMeltIt(img image.Image, intensity float64) (image.Image, error) {
	// Load the original image
	if intensity == 0 {
		intensity = float64(rand.Intn(15) + 15)
		fmt.Println("-------------------- melt intensity value=", intensity, "]]]]]]]]]]]")
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offsetY := int(float64(y) + intensity*math.Sin(float64(x)/50))
			if offsetY >= 0 && offsetY < height {
				newImg.Set(x, y, img.At(x, offsetY))
			} else {
				newImg.Set(x, y, color.Transparent)
			}
		}
	}
	return newImg, nil
}

//Spiral Start

// subtleSpiralEffect applies a spiral distortion to an image with control over pull distance, max angle, and max distance.
func spiralEffect(img image.Image, pullDistance, maxAngle, maxDistance float64, centerX, centerY float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	// Convert maxAngle from degrees to radians
	maxAngleRad := maxAngle * math.Pi / 180.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= maxDistance {
				angle := math.Atan2(dy, dx) + (pullDistance/maxDistance)*distance*maxAngleRad

				sx := int(centerX + distance*math.Cos(angle))
				sy := int(centerY + distance*math.Sin(angle))

				// Ensure sx and sy are within bounds
				if sx >= 0 && sx < width && sy >= 0 && sy < height {
					newImg.Set(x, y, img.At(sx, sy))
				} else {
					newImg.Set(x, y, img.At(x, y)) // Keep the original pixel if out of bounds
				}
			} else {
				//newImg.Set(x, y, img.At(x, y)) // Keep the original pixel if outside maxDistance
				angle := math.Atan2(dy, dx) + (pullDistance/maxDistance)*distance*maxAngleRad

				sx := int(centerX + distance*math.Cos(angle))
				sy := int(centerY + distance*math.Sin(angle))

				// Ensure sx and sy are within bounds
				if sx >= 0 && sx < width && sy >= 0 && sy < height {
					newImg.Set(x, y, img.At(sx, sy))
				} else {
					newImg.Set(x, y, img.At(x, y)) // Keep the original pixel if out of bounds
				}

			}
		}
	}
	saveImage(newImg, "spiralEffectEnd.jpg")
	return newImg
}

// applySubtleSpiralToQuadrants applies the subtle spiral effect to selected quadrants of an image.
func applySpiralToQuadrants(img image.Image, quadrants []string, pullDistance, maxAngle, maxDistance float64) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := img
	randomSize := maxDistance
	randomAngle := maxAngle
	randomPull := pullDistance
	saveImage(img, "ToBeSpiraled.jpg")
	screenHeight15 := (float64(height) * 0.15)
	//screenHeight25 := (float64(height) * 0.25)
	screenHeight15Int := int32(screenHeight15)
	//screenHeight25Int := int32(screenHeight25)
	for _, quadrant := range quadrants {
		var centerX, centerY float64
		if randomPull == 0 {
			pd := rand.Int31n(screenHeight15Int) + 50
			pullDistance = float64(pd)
		}
		if randomAngle == 0 {
			ma := rand.Int31n(75) + 25
			maxAngle = float64(ma)
		}
		switch quadrant {
		case "topLeft":
			centerX, centerY = float64(width)*0.25, float64(height)*0.25
		case "topRight":
			centerX, centerY = float64(width)*0.75, float64(height)*0.25
		case "bottomLeft":
			centerX, centerY = float64(width)*0.25, float64(height)*0.75
		case "bottomRight":
			centerX, centerY = float64(width)*0.75, float64(height)*0.75
		case "center":
			centerX, centerY = float64(width)*0.5, float64(height)*0.5
		}
		randFloat := rand.Int31n(5)
		if randomSize == 0 {
			maxDistance = float64(randFloat + 2)
		}

		// Apply the subtle spiral effect to the quadrant
		newImg = spiralEffect(newImg, pullDistance, maxAngle, maxDistance, centerX, centerY)
	}
	saveImage(newImg, "applySpiralToQuadrantsEnd.jpg")
	return newImg, nil
}

// func spiralEffect(img image.Image, level float64, centerX, centerY float64) image.Image {
// 	bounds := img.Bounds()
// 	width, height := bounds.Dx(), bounds.Dy()
// 	newImg := image.NewRGBA(bounds)
// 	dc := gg.NewContextForRGBA(newImg)
// 	_ = dc
// 	for y := 0; y < height; y++ {
// 		for x := 0; x < width; x++ {
// 			dx := float64(x) - centerX
// 			dy := float64(y) - centerY
// 			distance := math.Sqrt(dx*dx + dy*dy)
// 			angle := math.Atan2(dy, dx) + level*distance

// 			sx := int(centerX + distance*math.Cos(angle))
// 			sy := int(centerY + distance*math.Sin(angle))

// 			// Ensure sx and sy are within bounds
// 			if sx >= 0 && sx < width && sy >= 0 && sy < height {
// 				newImg.Set(x, y, img.At(sx, sy))
// 			} else {
// 				newImg.Set(x, y, img.At(x, y)) // Keep the original pixel if out of bounds
// 			}
// 		}
// 	}
// 	saveImage(newImg, "spiralEffectEnd.jpg")

// 	return newImg
// }

// func applySpiralToQuadrants(img image.Image, quadrants []string) (image.Image, error) {
// 	saveImage(img, "ToBeSpiraled.jpg")
// 	bounds := img.Bounds()
// 	width, height := bounds.Dx(), bounds.Dy()
// 	newImg := img

// 	// Randomize the spiral effect level
// 	spiralLevel := rand.Float64() * 0.05 // Adjust the range as needed

// 	for _, quadrant := range quadrants {
// 		var centerX, centerY float64

// 		switch quadrant {
// 		case "topLeft":
// 			centerX, centerY = float64(width)*0.25, float64(height)*0.25
// 		case "topRight":
// 			centerX, centerY = float64(width)*0.75, float64(height)*0.25
// 		case "bottomLeft":
// 			centerX, centerY = float64(width)*0.25, float64(height)*0.75
// 		case "bottomRight":
// 			centerX, centerY = float64(width)*0.75, float64(height)*0.75
// 		case "center":
// 			centerX, centerY = float64(width)*0.5, float64(height)*0.5
// 		}

// 		// Apply the spiral effect to the quadrant
// 		newImg = spiralEffect(newImg, spiralLevel, centerX, centerY)
// 	}
// 	saveImage(newImg, "applySpiralToQuadrantsEnd.jpg")
// 	return newImg, nil
// }

// func spiralEffect(img image.Image, level float64, centerX, centerY float64) image.Image {
// 	saveImage(img, "ToBeSpiraled.jpg")
// 	bounds := img.Bounds()
// 	width, height := bounds.Dx(), bounds.Dy()
// 	newImg := image.NewRGBA(bounds)
// 	dc := gg.NewContextForRGBA(newImg)
// 	_ = dc
// 	for y := 0; y < height; y++ {
// 		for x := 0; x < width; x++ {
// 			dx := float64(x) - centerX
// 			dy := float64(y) - centerY
// 			distance := math.Sqrt(dx*dx + dy*dy)
// 			angle := math.Atan2(dy, dx) + level*distance
// 			sx := int(centerX + distance*math.Cos(angle))
// 			sy := int(centerY + distance*math.Sin(angle))

// 			if sx >= 0 && sx < width && sy >= 0 && sy < height {
// 				newImg.Set(x, y, img.At(sx, sy))
// 			} else {
// 				newImg.Set(x, y, color.Transparent)
// 			}
// 		}
// 	}
// 	saveImage(newImg, "spiralEffectEnd.jpg")
// 	return newImg
// }

// func applySpiralToQuadrants(img image.Image, quadrants []string) (image.Image, error) {
// 	saveImage(img, "ToBeSpiraled.jpg")
// 	bounds := img.Bounds()
// 	width, height := bounds.Dx(), bounds.Dy()
// 	newImg := img

// 	// Randomize the spiral effect level
// 	spiralLevel := rand.Float64() * 0.05 // Adjust the range as needed

// 	for _, quadrant := range quadrants {
// 		var centerX, centerY float64

// 		switch quadrant {
// 		case "topLeft":
// 			centerX, centerY = float64(width)*0.25, float64(height)*0.25
// 		case "topRight":
// 			centerX, centerY = float64(width)*0.75, float64(height)*0.25
// 		case "bottomLeft":
// 			centerX, centerY = float64(width)*0.25, float64(height)*0.75
// 		case "bottomRight":
// 			centerX, centerY = float64(width)*0.75, float64(height)*0.75
// 		case "center":
// 			centerX, centerY = float64(width)*0.5, float64(height)*0.5
// 		}

// 		// Apply the spiral effect to the quadrant
// 		newImg = spiralEffect(newImg, spiralLevel, centerX, centerY)
// 	}
// 	saveImage(newImg, "applySpiralToQuadrantsEnd.jpg")
// 	return newImg, nil
// }

//Spiral End

func MonochromeIt(img image.Image) (image.Image, error) {
	// Create a new grayscale image
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Convert each pixel to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grayValue := uint8((r*299 + g*587 + b*114) / 1000 >> 8) // Standard formula for luminance
			grayImg.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}
	return grayImg, nil
}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//										OLD AND IN THE WAY
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Blur(currentPicturePath string, bufferPic string, blurSigma float64) (string, error) {
	src, err := imaging.Open(currentPicturePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Apply a Gaussian blur
	blurred := imaging.Blur(src, blurSigma) // The second parameter is the sigma of the Gaussian kernel

	// Save the resulting image
	err = imaging.Save(blurred, currentPicturePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	bufferPicExt := filepath.Ext(bufferPic)
	preBufferPath := strings.Replace(bufferPic, "Buffer"+bufferPicExt, "PreBuffer"+bufferPicExt, 1)
	DeleteFile(preBufferPath)
	copyFile(bufferPic, preBufferPath)
	DeleteFile(bufferPic)
	copyFile(currentPicturePath, bufferPic)
	return currentPicturePath, nil
}

func Pixelate(currentPicturePath string, bufferPic string, pixelSize int) (string, error) {
	// screenWidth := w32.GetSystemMetrics(w32.SM_CXSCREEN)
	// screenHeight := w32.GetSystemMetrics(w32.SM_CYSCREEN)
	if pixelSize == 0 {
		pixelSize = rand.Intn(16) + 4
	}
	// Load the original image
	file, err := os.Open(currentPicturePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return currentPicturePath, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return currentPicturePath, err
	}

	// Create a base image with the screen size
	//dc := gg.NewContext(screenWidth, screenHeight)

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y += pixelSize {
		for x := 0; x < width; x += pixelSize {
			var r, g, b, a uint32
			var count int

			for dy := 0; dy < pixelSize; dy++ {
				for dx := 0; dx < pixelSize; dx++ {
					if y+dy < height && x+dx < width {
						r1, g1, b1, a1 := img.At(x+dx, y+dy).RGBA()
						r += r1
						g += g1
						b += b1
						a += a1
						count++
					}
				}
			}

			if count > 0 {
				r /= uint32(count)
				g /= uint32(count)
				b /= uint32(count)
				a /= uint32(count)
			}

			for dy := 0; dy < pixelSize; dy++ {
				for dx := 0; dx < pixelSize; dx++ {
					if y+dy < height && x+dx < width {
						newImg.Set(x+dx, y+dy, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
					}
				}
			}
		}
	}
	bufferPicExt := filepath.Ext(bufferPic)
	preBufferPath := strings.Replace(bufferPic, "Buffer"+bufferPicExt, "PreBuffer"+bufferPicExt, 1)
	DeleteFile(preBufferPath)
	copyFile(bufferPic, preBufferPath)
	DeleteFile(bufferPic)
	// Save the pixelated image to the bufferPic path
	outFile, err := os.Create(bufferPic)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return currentPicturePath, err
	}
	defer outFile.Close()

	err = png.Encode(outFile, newImg)
	if err != nil {
		fmt.Println("Error saving pixelated image:", err)
		return currentPicturePath, err
	}
	DeleteFile(currentPicturePath)
	copyFile(bufferPic, currentPicturePath)
	return currentPicturePath, nil
}

func Oilify(currentPicturePath string, bufferPic string, radius int) (string, error) { //img image.Image, radius int) image.Image {

	if radius == 0 {
		radius = rand.Intn(6) + 3
	}
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	// fmt.Println("     ||||||||||||||||||||||    radius:", radius, "     ||||||||||||||||||||")
	// fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	// Load the original image
	file, err := os.Open(currentPicturePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return currentPicturePath, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return currentPicturePath, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			hist := make(map[color.Color]int)
			var mostCommonColor color.Color
			var maxCount int

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					px := x + dx
					py := y + dy

					if px >= 0 && px < width && py >= 0 && py < height {
						c := img.At(px, py)
						hist[c]++
						if hist[c] > maxCount {
							maxCount = hist[c]
							mostCommonColor = c
						}
					}
				}
			}

			newImg.Set(x, y, mostCommonColor)
		}
	}

	bufferPicExt := filepath.Ext(bufferPic)
	preBufferPath := strings.Replace(bufferPic, "Buffer"+bufferPicExt, "PreBuffer"+bufferPicExt, 1)
	DeleteFile(preBufferPath)
	copyFile(bufferPic, preBufferPath)
	DeleteFile(bufferPic)
	// Save the pixelated image to the bufferPic path
	outFile, err := os.Create(bufferPic)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return currentPicturePath, err
	}
	defer outFile.Close()

	err = png.Encode(outFile, newImg)
	if err != nil {
		fmt.Println("Error saving pixelated image:", err)
		return currentPicturePath, err
	}
	DeleteFile(currentPicturePath)
	copyFile(bufferPic, currentPicturePath)
	return currentPicturePath, nil

}

func MeltEffect(currentPicturePath string, bufferPic string, intensity float64) (string, error) {
	// Load the original image

	if intensity == 0 {
		intensity = float64(rand.Intn(15) + 15)
		fmt.Println("[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]")
		fmt.Println("-------------------- melt intensity value=", intensity, "]]]]]]]]]]]")
		fmt.Println("[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]")
	}

	file, err := os.Open(currentPicturePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return currentPicturePath, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return currentPicturePath, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offsetY := int(float64(y) + intensity*math.Sin(float64(x)/50))
			if offsetY >= 0 && offsetY < height {
				newImg.Set(x, y, img.At(x, offsetY))
			} else {
				newImg.Set(x, y, color.Transparent)
			}
		}
	}

	// Save the resulting image to the bufferPic path
	outFile, err := os.Create(bufferPic)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return currentPicturePath, err
	}
	defer outFile.Close()

	err = png.Encode(outFile, newImg)
	if err != nil {
		fmt.Println("Error saving melted image:", err)
		return currentPicturePath, err
	}

	return bufferPic, nil
}

func Monochrome(currentPicturePath string, bufferPic string) (string, error) {
	// Load the original image
	file, err := os.Open(currentPicturePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return currentPicturePath, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return currentPicturePath, err
	}

	// Create a new grayscale image
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Convert each pixel to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grayValue := uint8((r*299 + g*587 + b*114) / 1000 >> 8) // Standard formula for luminance
			grayImg.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	// Save the monochrome image to the bufferPic path
	outFile, err := os.Create(bufferPic)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return currentPicturePath, err
	}
	defer outFile.Close()

	err = png.Encode(outFile, grayImg)
	if err != nil {
		fmt.Println("Error saving monochrome image:", err)
		return currentPicturePath, err
	}

	return bufferPic, nil
}
