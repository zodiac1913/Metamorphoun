package service

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func RenderFontsSample(outputPath string) error {
	fontDir := "C:\\Windows\\Fonts"
	sampleText := "The quick brown fox jumps over the lazy dog"
	fontSize := 14.0
	lineHeight := 50
	margin := 10

	// Get all font files in the specified directory
	files, err := ioutil.ReadDir(fontDir)
	if err != nil {
		return fmt.Errorf("failed to read font directory: %v", err)
	}

	// Filter font files
	var fontFiles []string
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".ttf" || filepath.Ext(file.Name()) == ".otf") {
			fontFiles = append(fontFiles, filepath.Join(fontDir, file.Name()))
		}
	}

	// Create a new image context with enough height to fit all fonts
	imgHeight := (lineHeight + margin) * len(fontFiles)
	dc := gg.NewContext(2000, imgHeight)
	dc.SetColor(color.White)
	dc.Clear()
	dc.SetColor(color.Black)

	// Load the default font for printing file names
	defaultFontPath := "C:\\Windows\\Fonts\\DejaVuSans-Bold.ttf"
	defaultFontBytes, err := os.ReadFile(defaultFontPath)
	if err != nil {
		return fmt.Errorf("failed to read default font file: %v", err)
	}
	defaultFont, err := opentype.Parse(defaultFontBytes)
	if err != nil {
		return fmt.Errorf("failed to parse default font: %v", err)
	}
	defaultFace, err := opentype.NewFace(defaultFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("failed to create default font face: %v", err)
	}

	y := margin
	for _, fontPath := range fontFiles {
		// Load the font
		fontBytes, err := os.ReadFile(fontPath)
		if err != nil {
			fmt.Println("Error reading font file:", err)
			continue
		}
		fnt, err := opentype.Parse(fontBytes)
		if err != nil {
			fmt.Println("Error parsing font:", err)
			continue
		}
		face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
			Size:    fontSize,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			fmt.Println("Error creating font face:", err)
			continue
		}

		// Print the font file name using the default font
		dc.SetFontFace(defaultFace)
		dc.DrawStringAnchored(filepath.Base(fontPath), 100, float64(y)+fontSize/2, 0, 0.5)

		// Print the sample text using the current font
		dc.SetFontFace(face)
		dc.DrawStringAnchored(sampleText, 500, float64(y)+fontSize/2, 0, 0.5)

		y += lineHeight + margin
	}

	// Save the resulting image
	err = dc.SavePNG(outputPath)
	if err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}
