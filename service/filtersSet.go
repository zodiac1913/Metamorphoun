package service

import (
	"Metamorphoun/config"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"path"

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

// MosaicFilter creates a mosaic effect on an image by breaking it into randomly sized tiles.
func MosaicSet(currentPic config.PicHistory, img image.Image, tileMinSizeRatio, tileMaxSizeRatio float64) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Calculate the minimum and maximum tile sizes based on the image dimensions
	tileMinSize := int(tileMinSizeRatio * float64(width))
	tileMaxSize := int(tileMaxSizeRatio * float64(width))

	// Make grout much more visible
	groutMinSize := 15 // Increased for visibility
	groutMaxSize := 30 // Increased for visibility

	// Use a very dark grout color for contrast
	groutColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// Create a new blank image with a dark background (grout)
	mosaic := imaging.New(width, height, groutColor)

	for x := 0; x < width; {
		tileWidth := rand.Intn(tileMaxSize-tileMinSize+1) + tileMinSize
		if x+tileWidth > width {
			tileWidth = width - x
		}

		groutSize := 0

		for y := 0; y < height; {
			tileHeight := rand.Intn(tileMaxSize-tileMinSize+1) + tileMinSize
			if y+tileHeight > height {
				tileHeight = height - y
			}

			subImgRect := image.Rect(x, y, x+tileWidth, y+tileHeight)
			tile := imaging.Crop(img, subImgRect)

			// Edge-only blur: create a mask for the edges and blend blurred tile only at the edges
			edgeBlurRadius := 8.0 // How much of the edge to blur
			blurredTile := imaging.Blur(tile, 4.0)
			mask := image.NewAlpha(tile.Bounds())
			for mx := 0; mx < tileWidth; mx++ {
				for my := 0; my < tileHeight; my++ {
					distToEdge := min(min(mx, tileWidth-1-mx), min(my, tileHeight-1-my))
					alpha := uint8(0)
					if float64(distToEdge) < edgeBlurRadius {
						alpha = uint8(255 * (1.0 - float64(distToEdge)/edgeBlurRadius))
					}
					mask.SetAlpha(mx, my, color.Alpha{A: alpha})
				}
			}
			// Manually blend blurred edges with original tile using the mask
			tileWithBlurredEdges := image.NewNRGBA(tile.Bounds())
			for mx := 0; mx < tileWidth; mx++ {
				for my := 0; my < tileHeight; my++ {
					ma := mask.AlphaAt(mx, my).A
					origColor := tile.At(mx, my).(color.NRGBA)
					blurColor := blurredTile.At(mx, my).(color.NRGBA)
					// Blend: ma=0 -> orig, ma=255 -> blur
					blend := func(a, b uint8, alpha uint8) uint8 {
						return uint8((int(a)*(255-int(alpha)) + int(b)*int(alpha)) / 255)
					}
					tileWithBlurredEdges.SetNRGBA(mx, my, color.NRGBA{
						R: blend(origColor.R, blurColor.R, ma),
						G: blend(origColor.G, blurColor.G, ma),
						B: blend(origColor.B, blurColor.B, ma),
						A: origColor.A, // preserve original alpha
					})
				}
			}

			groutSize = rand.Intn(groutMaxSize-groutMinSize+1) + groutMinSize

			// Paste the tile with blurred edges into the mosaic with grout
			mosaic = imaging.Paste(mosaic, tileWithBlurredEdges, image.Pt(x+groutSize, y+groutSize))

			y += tileHeight + groutSize
		}
		x += tileWidth + groutSize
	}
	configFldr := GetFolderPath("config")
	fileLoc := path.Join(configFldr, "mosaicEnd.jpg")
	saveImage(mosaic, fileLoc)
	return mosaic, nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
