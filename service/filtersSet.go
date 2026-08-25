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

func DaliSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
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

// GraffitiItSet applies a graffiti/street-art effect using stored intensity from PicHistory
func GraffitiItSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	levels := int(currentPic.FilterIntensity)
	if levels < 2 {
		levels = 4
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	edgeThreshold := 80.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			edge := sobelEdge(img, x, y, width, height)

			if edge > edgeThreshold {
				newImg.Set(x, y, color.RGBA{20, 20, 20, 255})
			} else {
				r, g, b, a := img.At(x, y).RGBA()
				pr := posterize(uint8(r>>8), levels)
				pg := posterize(uint8(g>>8), levels)
				pb := posterize(uint8(b>>8), levels)

				maxC := math.Max(float64(pr), math.Max(float64(pg), float64(pb)))
				minC := math.Min(float64(pr), math.Min(float64(pg), float64(pb)))
				if maxC > 0 && maxC != minC {
					boost := 1.3
					mid := (maxC + minC) / 2
					pr = uint8(math.Min(255, mid+(float64(pr)-mid)*boost))
					pg = uint8(math.Min(255, mid+(float64(pg)-mid)*boost))
					pb = uint8(math.Min(255, mid+(float64(pb)-mid)*boost))
				}

				if rand.Intn(100) < 5 {
					noise := uint8(rand.Intn(30))
					pr = uint8(math.Min(255, float64(pr)+float64(noise)))
					pg = uint8(math.Min(255, float64(pg)+float64(noise)))
					pb = uint8(math.Min(255, float64(pb)+float64(noise)))
				}

				newImg.Set(x, y, color.RGBA{pr, pg, pb, uint8(a >> 8)})
			}
		}
	}
	return newImg, nil
}

// CyberpunkItSet recreates the offline focal neon treatment for a saved wallpaper.
func CyberpunkItSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	return cyberpunkFocus(img)
}

func cyberpunkFocus(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("cannot apply cyberpunk filter to an empty image")
	}

	// Downscaling makes the saliency calculation fast while favoring broad, coherent detail.
	analysis := imaging.Resize(img, 320, 0, imaging.Lanczos)
	analysisBounds := analysis.Bounds()
	analysisWidth, analysisHeight := analysisBounds.Dx(), analysisBounds.Dy()
	softened := imaging.Blur(analysis, 14)
	saliency := image.NewGray(analysisBounds)

	luminanceAt := func(source image.Image, x, y int) float64 {
		if x < analysisBounds.Min.X {
			x = analysisBounds.Min.X
		}
		if x >= analysisBounds.Max.X {
			x = analysisBounds.Max.X - 1
		}
		if y < analysisBounds.Min.Y {
			y = analysisBounds.Min.Y
		}
		if y >= analysisBounds.Max.Y {
			y = analysisBounds.Max.Y - 1
		}
		r, g, b, _ := source.At(x, y).RGBA()
		return 0.2126*float64(r>>8) + 0.7152*float64(g>>8) + 0.0722*float64(b>>8)
	}

	for y := analysisBounds.Min.Y; y < analysisBounds.Max.Y; y++ {
		for x := analysisBounds.Min.X; x < analysisBounds.Max.X; x++ {
			lum := luminanceAt(analysis, x, y)
			localContrast := math.Abs(lum-luminanceAt(softened, x, y)) / 128
			gx := luminanceAt(analysis, x+1, y) - luminanceAt(analysis, x-1, y)
			gy := luminanceAt(analysis, x, y+1) - luminanceAt(analysis, x, y-1)
			edge := math.Min(1, math.Sqrt(gx*gx+gy*gy)/180)
			r, g, b, _ := analysis.At(x, y).RGBA()
			maxChannel := math.Max(float64(r>>8), math.Max(float64(g>>8), float64(b>>8)))
			minChannel := math.Min(float64(r>>8), math.Min(float64(g>>8), float64(b>>8)))
			saturation := (maxChannel - minChannel) / 255
			centerX := (float64(x-analysisBounds.Min.X) / float64(analysisWidth)) - 0.5
			centerY := (float64(y-analysisBounds.Min.Y) / float64(analysisHeight)) - 0.5
			centerBias := math.Exp(-5 * (centerX*centerX + centerY*centerY))
			score := math.Min(1, 0.48*edge+0.24*localContrast+0.16*saturation+0.12*centerBias)
			saliency.SetGray(x, y, color.Gray{Y: uint8(score * 255)})
		}
	}

	// The blurred score chooses an area with sustained visual interest, not a single sharp edge.
	blurredSaliency := imaging.Blur(saliency, 14)
	focusX, focusY := analysisBounds.Min.X+analysisWidth/2, analysisBounds.Min.Y+analysisHeight/2
	var bestScore uint8
	for y := analysisBounds.Min.Y; y < analysisBounds.Max.Y; y++ {
		for x := analysisBounds.Min.X; x < analysisBounds.Max.X; x++ {
			scoreChannel, _, _, _ := blurredSaliency.At(x, y).RGBA()
			score := uint8(scoreChannel >> 8)
			if score > bestScore {
				bestScore = score
				focusX, focusY = x, y
			}
		}
	}

	focusCenterX := float64(bounds.Min.X) + ((float64(focusX-analysisBounds.Min.X) + 0.5) / float64(analysisWidth) * float64(width))
	focusCenterY := float64(bounds.Min.Y) + ((float64(focusY-analysisBounds.Min.Y) + 0.5) / float64(analysisHeight) * float64(height))
	focusWidth := float64(width) * (0.25 + rand.Float64()*0.50)
	focusHeight := float64(height) * (0.25 + rand.Float64()*0.50)
	if rand.Intn(100) < 30 {
		size := math.Min(focusWidth, focusHeight) * (0.85 + rand.Float64()*0.30)
		focusWidth = math.Min(float64(width)*0.75, size)
		focusHeight = math.Min(float64(height)*0.75, size)
	}
	halfFocusWidth := math.Max(1, focusWidth/2)
	halfFocusHeight := math.Max(1, focusHeight/2)
	focusCenterX = clampFloat(focusCenterX+(rand.Float64()-0.5)*halfFocusWidth*0.45, float64(bounds.Min.X)+halfFocusWidth*0.70, float64(bounds.Max.X)-halfFocusWidth*0.70)
	focusCenterY = clampFloat(focusCenterY+(rand.Float64()-0.5)*halfFocusHeight*0.45, float64(bounds.Min.Y)+halfFocusHeight*0.70, float64(bounds.Max.Y)-halfFocusHeight*0.70)
	focus := cyberpunkFocusSpec{
		centerX:    focusCenterX,
		centerY:    focusCenterY,
		halfWidth:  halfFocusWidth,
		halfHeight: halfFocusHeight,
		curvePower: 1.45 + rand.Float64()*3.75,
		softness:   0.10 + rand.Float64()*0.18,
		squared:    rand.Intn(100) < 38,
	}
	smoothed := imaging.Blur(img, 1.1)
	result := image.NewRGBA(bounds)
	edges := image.NewRGBA(bounds)
	contours := image.NewRGBA(bounds)
	figures := image.NewRGBA(bounds)
	drawCyberpunkContours(contours, bounds, focus)
	drawCyberpunkBackgroundFigures(figures, bounds, focus)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			red, green, blue := float64(r>>8), float64(g>>8), float64(b>>8)
			luminance := 0.2126*red + 0.7152*green + 0.0722*blue
			focusAmount := cyberpunkFocusAmount(float64(x), float64(y), focus)

			background := color.RGBA{R: uint8(luminance * 0.05), G: uint8(luminance * 0.11), B: uint8(14 + luminance*0.20), A: uint8(a >> 8)}
			neon := color.RGBA{
				R: uint8(math.Min(255, luminance*0.34+red*0.24+18)),
				G: uint8(math.Min(255, luminance*0.68+green*0.30+12)),
				B: uint8(math.Min(255, luminance*0.92+blue*0.34+36)),
				A: uint8(a >> 8),
			}
			result.SetRGBA(x, y, mixRGBA(background, neon, focusAmount))

			edge := sobelEdge(smoothed, x-bounds.Min.X, y-bounds.Min.Y, width, height)
			if edge > 85 && focusAmount > 0.08 {
				line := color.RGBA{R: 10, G: 235, B: 255, A: uint8(255 * focusAmount)}
				if red > blue {
					line = color.RGBA{R: 255, G: 24, B: 198, A: uint8(255 * focusAmount)}
				}
				edges.SetRGBA(x, y, line)
			}
		}
	}

	glow := imaging.Blur(edges, 6)
	contourGlow := imaging.Blur(contours, math.Max(4, math.Min(float64(width), float64(height))*0.006))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			focusAmount := cyberpunkFocusAmount(float64(x), float64(y), focus)
			darkness := 1 - focusAmount
			base := result.RGBAAt(x, y)
			fr, fg, fb, fa := figures.At(x, y).RGBA()
			if fa > 0 && darkness > 0.12 {
				amount := math.Min(0.72, (float64(fa>>8)/255)*darkness)
				base = mixRGBA(base, color.RGBA{R: uint8(fr >> 8), G: uint8(fg >> 8), B: uint8(fb >> 8), A: 255}, amount)
				result.SetRGBA(x, y, base)
			}
			cgr, cgg, cgb, cga := contourGlow.At(x, y).RGBA()
			if cga > 0 && darkness > 0.03 {
				amount := math.Min(0.58, (float64(cga>>8)/255)*darkness*0.65)
				base = mixRGBA(result.RGBAAt(x, y), color.RGBA{R: uint8(cgr >> 8), G: uint8(cgg >> 8), B: uint8(cgb >> 8), A: 255}, amount)
				result.SetRGBA(x, y, base)
			}
			glowColor := glow.At(x, y)
			gr, gg, gb, ga := glowColor.RGBA()
			if ga > 0 {
				result.SetRGBA(x, y, mixRGBA(base, color.RGBA{R: uint8(gr >> 8), G: uint8(gg >> 8), B: uint8(gb >> 8), A: 255}, 0.38))
			}
			cr, cg, cb, ca := contours.At(x, y).RGBA()
			if ca > 0 && darkness > 0.02 {
				amount := math.Min(0.92, (float64(ca>>8)/255)*math.Max(0.35, darkness))
				result.SetRGBA(x, y, mixRGBA(result.RGBAAt(x, y), color.RGBA{R: uint8(cr >> 8), G: uint8(cg >> 8), B: uint8(cb >> 8), A: 255}, amount))
			}
			_, _, _, edgeAlpha := edges.At(x, y).RGBA()
			if edgeAlpha > 0 {
				er, eg, eb, _ := edges.At(x, y).RGBA()
				result.SetRGBA(x, y, color.RGBA{R: uint8(er >> 8), G: uint8(eg >> 8), B: uint8(eb >> 8), A: 255})
			}
		}
	}

	return result, nil
}

type cyberpunkFocusSpec struct {
	centerX    float64
	centerY    float64
	halfWidth  float64
	halfHeight float64
	curvePower float64
	softness   float64
	squared    bool
}

type cyberpunkContourSpec struct {
	centerX    float64
	centerY    float64
	radiusX    float64
	radiusY    float64
	rotation   float64
	curvePower float64
	squared    bool
}

func cyberpunkFocusAmount(x, y float64, focus cyberpunkFocusSpec) float64 {
	dx := math.Abs(x-focus.centerX) / math.Max(1, focus.halfWidth)
	dy := math.Abs(y-focus.centerY) / math.Max(1, focus.halfHeight)
	distance := math.Max(dx, dy)
	if !focus.squared {
		power := math.Max(1, focus.curvePower)
		distance = math.Pow(math.Pow(dx, power)+math.Pow(dy, power), 1/power)
	}
	amount := 1 - smoothstep(1-focus.softness, 1+focus.softness, distance)
	return math.Max(0, math.Min(1, amount))
}

func smoothstep(edge0, edge1, value float64) float64 {
	if edge0 == edge1 {
		if value < edge0 {
			return 0
		}
		return 1
	}
	t := clampFloat((value-edge0)/(edge1-edge0), 0, 1)
	return t * t * (3 - 2*t)
}

func clampFloat(value, minValue, maxValue float64) float64 {
	if minValue > maxValue {
		return (minValue + maxValue) / 2
	}
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func randomCyberpunkColor(alpha uint8) color.RGBA {
	palette := []color.RGBA{
		{R: 0, G: 244, B: 255, A: alpha},
		{R: 255, G: 28, B: 205, A: alpha},
		{R: 255, G: 238, B: 38, A: alpha},
		{R: 0, G: 255, B: 137, A: alpha},
		{R: 255, G: 83, B: 35, A: alpha},
		{R: 126, G: 58, B: 255, A: alpha},
	}
	return palette[rand.Intn(len(palette))]
}

func drawCyberpunkContours(layer *image.RGBA, bounds image.Rectangle, focus cyberpunkFocusSpec) {
	dc := gg.NewContextForRGBA(layer)
	minDimension := math.Max(1, math.Min(float64(bounds.Dx()), float64(bounds.Dy())))
	count := 5 + rand.Intn(7)
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		push := 0.85 + rand.Float64()*0.75
		contour := cyberpunkContourSpec{
			centerX:    focus.centerX + math.Cos(angle)*focus.halfWidth*push,
			centerY:    focus.centerY + math.Sin(angle)*focus.halfHeight*push,
			radiusX:    focus.halfWidth * (0.40 + rand.Float64()*1.05),
			radiusY:    focus.halfHeight * (0.26 + rand.Float64()*0.92),
			rotation:   angle + (rand.Float64()-0.5)*0.9,
			curvePower: 1.4 + rand.Float64()*4.6,
			squared:    rand.Intn(100) < 42,
		}
		bands := 1 + rand.Intn(3)
		for band := 0; band < bands; band++ {
			strokeWidth := minDimension*(0.0025+rand.Float64()*0.0065) + float64(band)*minDimension*0.002
			contourColor := randomCyberpunkColor(uint8(145 + rand.Intn(95)))
			dc.SetColor(contourColor)
			dc.SetLineWidth(strokeWidth)
			bandContour := contour
			bandContour.radiusX += float64(band) * strokeWidth * 2.7
			bandContour.radiusY += float64(band) * strokeWidth * 2.7
			traceCyberpunkContour(dc, bandContour)
			dc.Stroke()
		}
	}
}

func traceCyberpunkContour(dc *gg.Context, contour cyberpunkContourSpec) {
	if contour.squared {
		corner := math.Min(contour.radiusX, contour.radiusY) * (0.04 + rand.Float64()*0.16)
		points := [][2]float64{
			{-contour.radiusX + corner, -contour.radiusY}, {contour.radiusX - corner, -contour.radiusY}, {contour.radiusX, -contour.radiusY + corner}, {contour.radiusX, contour.radiusY - corner},
			{contour.radiusX - corner, contour.radiusY}, {-contour.radiusX + corner, contour.radiusY}, {-contour.radiusX, contour.radiusY - corner}, {-contour.radiusX, -contour.radiusY + corner},
		}
		for index, point := range points {
			x, y := rotatePoint(point[0], point[1], contour.rotation)
			if index == 0 {
				dc.MoveTo(contour.centerX+x, contour.centerY+y)
			} else {
				dc.LineTo(contour.centerX+x, contour.centerY+y)
			}
		}
		dc.ClosePath()
		return
	}

	steps := 96
	power := math.Max(1, contour.curvePower)
	for step := 0; step <= steps; step++ {
		angle := (float64(step) / float64(steps)) * math.Pi * 2
		cosine := math.Cos(angle)
		sine := math.Sin(angle)
		px := math.Copysign(math.Pow(math.Abs(cosine), 2/power)*contour.radiusX, cosine)
		py := math.Copysign(math.Pow(math.Abs(sine), 2/power)*contour.radiusY, sine)
		rx, ry := rotatePoint(px, py, contour.rotation)
		if step == 0 {
			dc.MoveTo(contour.centerX+rx, contour.centerY+ry)
		} else {
			dc.LineTo(contour.centerX+rx, contour.centerY+ry)
		}
	}
	dc.ClosePath()
}

func rotatePoint(x, y, rotation float64) (float64, float64) {
	cosine := math.Cos(rotation)
	sine := math.Sin(rotation)
	return x*cosine - y*sine, x*sine + y*cosine
}

func drawCyberpunkBackgroundFigures(layer *image.RGBA, bounds image.Rectangle, focus cyberpunkFocusSpec) {
	dc := gg.NewContextForRGBA(layer)
	width, height := float64(bounds.Dx()), float64(bounds.Dy())
	minDimension := math.Max(1, math.Min(width, height))
	drawCyberpunkCelestialBodies(dc, bounds, focus, width, height, minDimension)
	drawCyberpunkBuildings(dc, bounds, focus, width, height)
}

func drawCyberpunkCelestialBodies(dc *gg.Context, bounds image.Rectangle, focus cyberpunkFocusSpec, width, height, minDimension float64) {
	for i := 0; i < 3+rand.Intn(5); i++ {
		x := float64(bounds.Min.X) + rand.Float64()*width
		y := float64(bounds.Min.Y) + rand.Float64()*height*0.70
		if cyberpunkFocusAmount(x, y, focus) > 0.25 {
			continue
		}
		radius := minDimension * (0.018 + rand.Float64()*0.055)
		if rand.Intn(100) < 42 {
			drawCyberpunkMoon(dc, x, y, radius*(1.15+rand.Float64()*0.65))
		} else {
			bodyColor := randomCyberpunkColor(82)
			dc.SetColor(bodyColor)
			dc.SetLineWidth(math.Max(1.5, radius*0.08))
			dc.DrawCircle(x, y, radius)
			dc.Stroke()
			dc.SetColor(randomCyberpunkColor(58))
			dc.SetLineWidth(math.Max(1, radius*0.04))
			dc.DrawEllipse(x, y, radius*(1.45+rand.Float64()*0.80), radius*(0.24+rand.Float64()*0.22))
			dc.Stroke()
		}
	}
}

func drawCyberpunkBuildings(dc *gg.Context, bounds image.Rectangle, focus cyberpunkFocusSpec, width, height float64) {
	buildingCount := 4 + rand.Intn(8)
	for i := 0; i < buildingCount; i++ {
		buildingWidth := width * (0.018 + rand.Float64()*0.045)
		buildingHeight := height * (0.10 + rand.Float64()*0.28)
		x := float64(bounds.Min.X) + rand.Float64()*(width-buildingWidth)
		y := float64(bounds.Max.Y) - buildingHeight
		if cyberpunkFocusAmount(x+buildingWidth/2, y+buildingHeight/2, focus) > 0.40 {
			continue
		}
		dc.SetColor(color.RGBA{R: 3, G: 8, B: 22, A: 170})
		dc.DrawRectangle(x, y, buildingWidth, buildingHeight)
		dc.Fill()
		windowColor := randomCyberpunkColor(135)
		dc.SetColor(windowColor)
		columns := 2 + rand.Intn(3)
		rows := 4 + rand.Intn(8)
		windowWidth := buildingWidth / float64(columns*3)
		windowHeight := buildingHeight / float64(rows*5)
		for row := 0; row < rows; row++ {
			for col := 0; col < columns; col++ {
				if rand.Intn(100) < 45 {
					continue
				}
				wx := x + buildingWidth*0.20 + float64(col)*buildingWidth/float64(columns)
				wy := y + buildingHeight*0.12 + float64(row)*buildingHeight/float64(rows)
				dc.DrawRectangle(wx, wy, windowWidth, windowHeight)
				dc.Fill()
			}
		}
	}
}

func drawCyberpunkMoon(dc *gg.Context, x, y, radius float64) {
	moonColor := color.RGBA{R: 24, G: 202, B: 255, A: 132}
	hotEdge := color.RGBA{R: 114, G: 247, B: 255, A: 210}
	dc.SetColor(moonColor)
	dc.DrawCircle(x, y, radius)
	dc.Fill()
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 16, A: 205})
	dc.DrawCircle(x+radius*(0.30+rand.Float64()*0.30), y-radius*(0.04+rand.Float64()*0.12), radius*(0.78+rand.Float64()*0.18))
	dc.Fill()
	dc.SetColor(hotEdge)
	dc.SetLineWidth(math.Max(1.5, radius*0.11))
	dc.DrawArc(x, y, radius, math.Pi*0.58, math.Pi*1.42)
	dc.Stroke()
	dc.SetColor(color.RGBA{R: 255, G: 36, B: 210, A: 92})
	dc.SetLineWidth(math.Max(1, radius*0.045))
	dc.DrawArc(x, y, radius*1.16, math.Pi*0.56, math.Pi*1.44)
	dc.Stroke()
}

func mixRGBA(left, right color.RGBA, amount float64) color.RGBA {
	amount = math.Max(0, math.Min(1, amount))
	return color.RGBA{
		R: uint8(float64(left.R)*(1-amount) + float64(right.R)*amount),
		G: uint8(float64(left.G)*(1-amount) + float64(right.G)*amount),
		B: uint8(float64(left.B)*(1-amount) + float64(right.B)*amount),
		A: uint8(float64(left.A)*(1-amount) + float64(right.A)*amount),
	}
}

// CartoonSet applies a cartoon / cel-shading effect to the image.
// It posterizes the colours to a small number of levels for that flat,
// hand-painted look, then overlays dark edge lines detected via a Sobel filter.
func CartoonSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	// Slight blur first to reduce noise before edge detection
	smoothed := imaging.Blur(img, 2.0)

	bounds := smoothed.Bounds()
	W := bounds.Dx()
	H := bounds.Dy()

	// --- Posterize: reduce each channel to a handful of levels -------------
	levels := 6 + rand.Intn(4) // 6-9 colour levels per channel
	step := 256.0 / float64(levels)

	posterized := image.NewRGBA(bounds)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			r, g, b, a := smoothed.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			pr := uint8(math.Floor(float64(r>>8)/step) * step)
			pg := uint8(math.Floor(float64(g>>8)/step) * step)
			pb := uint8(math.Floor(float64(b>>8)/step) * step)
			posterized.SetRGBA(x+bounds.Min.X, y+bounds.Min.Y, color.RGBA{pr, pg, pb, uint8(a >> 8)})
		}
	}

	// --- Edge detection (Sobel) on the smoothed greyscale -----------------
	grey := imaging.Grayscale(smoothed)
	edges := image.NewRGBA(bounds)

	// Sobel kernels
	pixel := func(px, py int) float64 {
		if px < bounds.Min.X {
			px = bounds.Min.X
		}
		if py < bounds.Min.Y {
			py = bounds.Min.Y
		}
		if px >= bounds.Min.X+W {
			px = bounds.Min.X + W - 1
		}
		if py >= bounds.Min.Y+H {
			py = bounds.Min.Y + H - 1
		}
		r, _, _, _ := grey.At(px, py).RGBA()
		return float64(r >> 8)
	}

	edgeThreshold := 28.0 + rand.Float64()*20.0 // 28-48, varies per image
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			ox := x + bounds.Min.X
			oy := y + bounds.Min.Y
			// Gx
			gx := -pixel(ox-1, oy-1) - 2*pixel(ox-1, oy) - pixel(ox-1, oy+1) +
				pixel(ox+1, oy-1) + 2*pixel(ox+1, oy) + pixel(ox+1, oy+1)
			// Gy
			gy := -pixel(ox-1, oy-1) - 2*pixel(ox, oy-1) - pixel(ox+1, oy-1) +
				pixel(ox-1, oy+1) + 2*pixel(ox, oy+1) + pixel(ox+1, oy+1)

			mag := math.Sqrt(gx*gx + gy*gy)
			if mag > edgeThreshold {
				edges.SetRGBA(ox, oy, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	// --- Composite: posterized image with edge overlay --------------------
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, posterized, bounds.Min, draw.Src)
	// Draw edges on top (only black pixels)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			ox := x + bounds.Min.X
			oy := y + bounds.Min.Y
			er, _, _, ea := edges.At(ox, oy).RGBA()
			if ea > 0 && er == 0 {
				result.SetRGBA(ox, oy, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	return result, nil
}

// JigsawPuzzleSet applies a jigsaw puzzle effect with realistic interlocking pieces
// JigsawPuzzleSet overlays black interlocking jigsaw-puzzle lines on the image.
// Each interior edge gets a smooth tab (bump) that alternates direction so
// neighbouring pieces interlock.  The grid targets roughly 30–50 pieces per
// axis depending on FilterIntensity (low = fewer/bigger, high = more/smaller).
func JigsawPuzzleSet(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	W := float64(bounds.Dx())
	H := float64(bounds.Dy())

	// --- grid size ----------------------------------------------------------
	// Map FilterIntensity (typically 1-10) to columns in the 30-50 range.
	cols := int(30 + currentPic.FilterIntensity*2)
	if cols < 20 {
		cols = 20
	}
	if cols > 60 {
		cols = 60
	}
	pieceW := W / float64(cols)
	rows := int(math.Round(H / pieceW)) // keep pieces roughly square
	if rows < 10 {
		rows = 10
	}
	pieceH := H / float64(rows)

	// --- line thickness -----------------------------------------------------
	lineW := 3.0 + rand.Float64()*3.0 // 3-6 px, varies per image

	// --- deterministic tab directions per edge ------------------------------
	// hTabs[row][col] = true means the tab on the horizontal edge between
	// row-1 and row bumps downward; false = upward.  Similar for vTabs.
	hTabs := make([][]bool, rows+1)
	for r := range hTabs {
		hTabs[r] = make([]bool, cols)
		for c := range hTabs[r] {
			hTabs[r][c] = rand.Intn(2) == 0
		}
	}
	vTabs := make([][]bool, rows)
	for r := range vTabs {
		vTabs[r] = make([]bool, cols+1)
		for c := range vTabs[r] {
			vTabs[r][c] = rand.Intn(2) == 0
		}
	}

	// --- draw on a gg context -----------------------------------------------
	dc := gg.NewContextForImage(img)
	dc.SetColor(color.Black)
	dc.SetLineWidth(lineW)
	dc.SetLineCapButt()

	// drawHTab draws a classic jigsaw nub on a horizontal edge.
	// Each nub gets randomized proportions for a natural look.
	drawHEdge := func(col, row int, down bool) {
		x0 := float64(col) * pieceW
		y0 := float64(row) * pieceH
		seg := pieceW
		d := 1.0
		if !down {
			d = -1.0
		}

		// Randomize nub position along the edge (30-45% from left)
		neckStart := seg * (0.30 + rand.Float64()*0.15)
		nubWidth := seg * (0.22 + rand.Float64()*0.10) // nub spans 22-32% of edge
		neckEnd := neckStart + nubWidth

		// Randomize how far the nub sticks out (25-38% of piece height)
		nubH := pieceH * (0.25 + rand.Float64()*0.13)
		// Randomize neck pinch (3-8% of segment)
		neckW := seg * (0.03 + rand.Float64()*0.05)

		dc.MoveTo(x0, y0)
		dc.LineTo(x0+neckStart, y0)

		dc.CubicTo(
			x0+neckStart+neckW, y0,
			x0+neckStart-neckW, y0+d*nubH*0.4,
			x0+neckStart-neckW, y0+d*nubH*0.6,
		)
		dc.CubicTo(
			x0+neckStart-neckW, y0+d*nubH,
			x0+neckEnd+neckW, y0+d*nubH,
			x0+neckEnd+neckW, y0+d*nubH*0.6,
		)
		dc.CubicTo(
			x0+neckEnd+neckW, y0+d*nubH*0.4,
			x0+neckEnd-neckW, y0,
			x0+neckEnd, y0,
		)

		dc.LineTo(x0+seg, y0)
	}

	// drawVTab draws a classic jigsaw nub on a vertical edge.
	// Same mushroom shape rotated 90°, with per-tab randomization.
	drawVEdge := func(col, row int, right bool) {
		x0 := float64(col) * pieceW
		y0 := float64(row) * pieceH
		seg := pieceH
		d := 1.0
		if !right {
			d = -1.0
		}

		// Randomize nub position along the edge (30-45% from top)
		neckStart := seg * (0.30 + rand.Float64()*0.15)
		nubHeight := seg * (0.22 + rand.Float64()*0.10)
		neckEnd := neckStart + nubHeight

		// Randomize protrusion (25-38% of piece width)
		nubW := pieceW * (0.25 + rand.Float64()*0.13)
		// Randomize neck pinch
		neckH := seg * (0.03 + rand.Float64()*0.05)

		dc.MoveTo(x0, y0)
		dc.LineTo(x0, y0+neckStart)

		dc.CubicTo(
			x0, y0+neckStart+neckH,
			x0+d*nubW*0.4, y0+neckStart-neckH,
			x0+d*nubW*0.6, y0+neckStart-neckH,
		)
		dc.CubicTo(
			x0+d*nubW, y0+neckStart-neckH,
			x0+d*nubW, y0+neckEnd+neckH,
			x0+d*nubW*0.6, y0+neckEnd+neckH,
		)
		dc.CubicTo(
			x0+d*nubW*0.4, y0+neckEnd+neckH,
			x0, y0+neckEnd-neckH,
			x0, y0+neckEnd,
		)

		dc.LineTo(x0, y0+seg)
	}

	// Interior horizontal edges (skip top=0 and bottom=rows)
	for row := 1; row < rows; row++ {
		for col := 0; col < cols; col++ {
			drawHEdge(col, row, hTabs[row][col])
		}
	}
	dc.Stroke()

	// Interior vertical edges (skip left=0 and right=cols)
	for row := 0; row < rows; row++ {
		for col := 1; col < cols; col++ {
			drawVEdge(col, row, vTabs[row][col])
		}
	}
	dc.Stroke()

	// Outer border
	border := lineW / 2.0
	dc.DrawRectangle(border, border, W-lineW, H-lineW)
	dc.Stroke()

	return dc.Image(), nil
}
