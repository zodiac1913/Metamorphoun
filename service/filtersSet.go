package service

import (
	"Metamorphoun/config"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

func BlurItSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	// Apply a Gaussian blur
	blurred := imaging.Blur(img, currentPic.FilterIntensity) // The second parameter is the sigma of the Gaussian kernel
	return blurred, nil
}

func PixelateItSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {

	pixelSize := int(currentPic.FilterIntensity)
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
	return newImg, nil
}

func OilifyItSet(currentPic config.PicHistory, img image.Image) (image.Image, error) { //img image.Image, radius int) image.Image {
	radius := int(currentPic.FilterIntensity)
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

func PicassoSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	screenInfo := GetScreenInfo()[0]
	screenWidth := screenInfo.Width
	screenHeight := screenInfo.Height
	intensity := float64(currentPic.FilterIntensity)

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
	return resizedImg, nil
}

func vortexEffectSet(currentPic config.PicHistory, img image.Image, quadrant string, level float64, centerX, centerY float64) (config.PicHistory, image.Image) {
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
	//picV := config.PicHistoryVortex{FilterIntensity: level, FilterX: centerX, FilterY: centerY}
	picV := config.PicHistoryVortex{FilterQuadrant: quadrant, FilterIntensity: level, FilterX: centerX, FilterY: centerY}
	currentPic.FilterVortices = append(currentPic.FilterVortices, []config.PicHistoryVortex{picV}...)
	//saveImage(newImg, "spiralEffectEnd.jpg")

	return currentPic, newImg
}

func applyVortexToQuadrantsSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	newImg := img
	fmt.Println("Applying vortex effect to quadrants:")
	for _, quadrant := range currentPic.FilterVortices {
		var centerX, centerY float64
		spiralLevel := quadrant.FilterIntensity
		fmt.Println("Spiral level:", spiralLevel)

		fmt.Println("Quadrant:", quadrant.FilterQuadrant)

		switch quadrant.FilterQuadrant {
		case "topLeft":
			centerX, centerY = quadrant.FilterX, quadrant.FilterY
		case "topRight":
			centerX, centerY = quadrant.FilterX, quadrant.FilterY
		case "bottomLeft":
			centerX, centerY = quadrant.FilterX, quadrant.FilterY
		case "bottomRight":
			centerX, centerY = quadrant.FilterX, quadrant.FilterY
		case "center":
			centerX, centerY = quadrant.FilterX, quadrant.FilterY
		}
		// Apply the spiral effect to the quadrant
		currentPic, newImg = vortexEffectSet(currentPic, newImg, quadrant.FilterQuadrant, spiralLevel, centerX, centerY)
	}
	//saveImage(newImg, "applySpiralToQuadrantsEnd.jpg")
	return newImg, nil
}

//MonochromeIt is constant so it goes to the random(not so random function)
// package filters

// import (
//     "image"
//     "image/color"
//     "math/rand"

//     "github.com/disintegration/imaging"
// )

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MosaicSet creates a mosaic effect on the input image.
func MosaicSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	rsRnd := float64((rand.Intn(50) + 50))
	reductionScale := float64(rsRnd / 100) //0.95
	maxJitter := (rand.Intn(2) + 1)
	numberOfTiles := (rand.Intn(50) + 35) * (rand.Intn(3) + 1)

	fmt.Println("Number of tiles:", numberOfTiles)

	bounds := img.Bounds()
	origWidth, origHeight := bounds.Dx(), bounds.Dy()

	//grout color components randomized
	groutColor := color.RGBA{R: uint8(rand.Intn(96)), G: uint8(rand.Intn(96)), B: uint8(rand.Intn(96)), A: 255}
	mosaic := image.NewNRGBA(bounds)
	draw.Draw(mosaic, mosaic.Bounds(), &image.Uniform{C: groutColor}, image.Point{}, draw.Src)

	scaledWidth := int(float64(origWidth) * reductionScale)
	scaledImg := imaging.Resize(img, scaledWidth, 0, imaging.Lanczos)
	scaledBounds := scaledImg.Bounds()
	scaledW, scaledH := scaledBounds.Dx(), scaledBounds.Dy()
	tileW := origWidth / numberOfTiles
	tileH := origHeight / numberOfTiles

	stretchRatioX := float64(origWidth) / float64(scaledW)
	stretchRatioY := float64(origHeight) / float64(scaledH)

	// tileMinSize := int(tileMinSizeRatio * float64(scaledW))
	// tileMaxSize := int(tileMaxSizeRatio * float64(scaledW))

	// if tileMinSize <= 1 {
	// 	tileMinSize = 2
	// }
	// if tileMinSize >= tileMaxSize {
	// 	tileMaxSize = tileMinSize + 20
	// }

	// Calculate number of tiles for even distribution.
	numTilesX := origWidth / tileW //(tileMinSize * 2) // Adjust multiplier to control density
	numTilesY := origHeight / tileH

	if numTilesX < 1 {
		numTilesX = 1
	}
	if numTilesY < 1 {
		numTilesY = 1
	}

	tileWidth := tileW  //origWidth / numTilesX
	tileHeight := tileH //origHeight / numTilesY

	for xTile := 0; xTile < numTilesX; xTile++ {
		x := int(float64(xTile) * float64(tileWidth))

		for yTile := 0; yTile < numTilesY; yTile++ {
			y := int(float64(yTile) * float64(tileHeight))

			scaledX := int(float64(x) / stretchRatioX)
			scaledY := int(float64(y) / stretchRatioY)

			tileWidthScaled := int(float64(tileWidth) / stretchRatioX)
			if scaledX+tileWidthScaled > scaledW {
				tileWidthScaled = scaledW - scaledX
			}

			tileHeightScaled := int(float64(tileHeight) / stretchRatioY)
			if scaledY+tileHeightScaled > scaledH {
				tileHeightScaled = scaledH - scaledY
			}

			cropRect := image.Rect(scaledX, scaledY, scaledX+tileWidthScaled, scaledY+tileHeightScaled)
			tile := imaging.Crop(scaledImg, cropRect)

			jitterX := rand.Intn(maxJitter*2+1) - maxJitter
			jitterY := rand.Intn(maxJitter*2+1) - maxJitter

			pastePoint := image.Pt(x+jitterX, y+jitterY)
			destRect := tile.Bounds().Add(pastePoint)

			draw.Draw(mosaic, destRect, tile, image.Point{}, draw.Over)
		}
	}
	//saveImage(mosaic, "mosaicEnd.jpg")
	return mosaic, nil
}
