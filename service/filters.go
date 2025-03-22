package service

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
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

// -------------------------FINAL PICASSO---------------------------------------
func Picasso(img image.Image, intensity float64) (image.Image, error) {
	screenInfo := getScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height
	if intensity == 0 {
		intensity = float64(rand.Intn(15) + 15)
		fmt.Printf("Melt intensity value = %.2f\n", intensity)
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	fmt.Printf("Original image dimensions: width=%d, height=%d\n", width, height)

	// Step 1: Calculate maximum distortion size
	maxDistortion := int(intensity) // Maximum vertical distortion in pixels
	fmt.Printf("Calculated max distortion size: %d pixels\n", maxDistortion)

	// Step 2: Distort the image
	distortedImg := image.NewNRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			distortedY := y + int(intensity*math.Sin(float64(x)/50))
			if distortedY >= 0 && distortedY < height {
				distortedImg.Set(x, distortedY, img.At(x, y))
			} else {
				distortedImg.Set(x, y, color.Transparent)
			}
		}
	}

	// Step 3: Clip the interior to remove distorted edges
	// New dimensions after clipping
	clippedWidth := width - 2*maxDistortion
	clippedHeight := height - 2*maxDistortion
	fmt.Printf("Clipped image dimensions: width=%d, height=%d\n", clippedWidth, clippedHeight)

	// Define clipping rectangle (centered)
	clipRect := image.Rect(maxDistortion, maxDistortion, width-maxDistortion, height-maxDistortion)
	clippedImg := distortedImg.SubImage(clipRect).(*image.NRGBA)

	// Step 4: Resize the image back to screen size
	resizedImg := image.NewNRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	draw.BiLinear.Scale(resizedImg, resizedImg.Rect, clippedImg, clippedImg.Bounds(), draw.Over, nil)
	// resizedImg := image.NewNRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	// for y := 0; y < screenHeight; y++ {
	// 	for x := 0; x < screenWidth; x++ {
	// 		// Map screen coordinates to clipped image coordinates
	// 		srcX := x * clippedWidth / screenWidth
	// 		srcY := y * clippedHeight / screenHeight
	// 		resizedImg.Set(x, y, clippedImg.At(srcX, srcY))
	// 	}
	// }

	return resizedImg, nil
}

func vortexEffect(img image.Image, level float64, centerX, centerY float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)
	dc := gg.NewContextForRGBA(newImg)
	_ = dc
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distance := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx) + level*distance

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
	saveImage(newImg, "spiralEffectEnd.jpg")

	return newImg
}

func applyVortexToQuadrants(img image.Image, quadrants []string) (image.Image, error) {
	//saveImage(img, "ToBeSpiraled.jpg")
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := img

	// Randomize the spiral effect level
	spiralLevel := rand.Float64() * 0.2 // Adjust the range as needed

	for _, quadrant := range quadrants {
		var centerX, centerY float64

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

		// Apply the spiral effect to the quadrant
		newImg = vortexEffect(newImg, spiralLevel, centerX, centerY)
	}
	//saveImage(newImg, "applySpiralToQuadrantsEnd.jpg")
	return newImg, nil
}

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
