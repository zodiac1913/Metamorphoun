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
