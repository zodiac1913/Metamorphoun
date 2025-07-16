package service

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// fontDirs: array of font folder paths (cross-platform)
// outputPaths: array of output image paths (one per folder)
func RenderFontsSample(fontDirs []string) error {
	if len(fontDirs) == 0 {
		return fmt.Errorf("fontDirs must be non-empty")
	}

	defaultFontPath := getDefaultFontPath()
	defaultFontBytes, err := os.ReadFile(defaultFontPath)
	if err != nil {
		return fmt.Errorf("failed to read default font file: %v", err)
	}
	defaultFont, err := opentype.Parse(defaultFontBytes)
	if err != nil {
		return fmt.Errorf("failed to parse default font: %v", err)
	}

	fontSize := 14.0
	lineHeight := 50
	margin := 10
	var failedFonts []string

	for _, fontDir := range fontDirs {
		var fontFiles []string
		println("Processing font directory:", fontDir)

		// Recursively walk the directory to find all .ttf and .otf files
		err := filepath.Walk(fontDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && (filepath.Ext(path) == ".ttf" || filepath.Ext(path) == ".otf") {
				println("Found font file:", path)
				fontFiles = append(fontFiles, path)

			}
			return nil
		})
		if err != nil {
			fmt.Printf("failed to walk font directory %s: %v\n", fontDir, err)
			continue
		}

		imgHeight := (lineHeight + margin) * len(fontFiles)
		if imgHeight < 100 {
			imgHeight = 100
		}
		dc := gg.NewContext(2000, imgHeight)
		dc.SetColor(color.White)
		dc.Clear()
		dc.SetColor(color.Black)

		defaultFace, err := opentype.NewFace(defaultFont, &opentype.FaceOptions{
			Size:    fontSize,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			fmt.Printf("failed to create default font face: %v\n", err)
			continue
		}

		y := margin
		for _, fontPath := range fontFiles {
			fontBytes, err := os.ReadFile(fontPath)
			if err != nil {
				failedFonts = append(failedFonts, fontPath+" (read error)")
				continue
			}
			fnt, err := opentype.Parse(fontBytes)
			if err != nil {
				failedFonts = append(failedFonts, fontPath+" (parse error)")
				continue
			}
			face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
				Size:    fontSize,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				failedFonts = append(failedFonts, fontPath+" (face error)")
				continue
			}

			dc.SetFontFace(defaultFace)
			dc.SetColor(color.Black) // Ensure font name is black
			dc.DrawStringAnchored(filepath.Base(fontPath), 100, float64(y)+fontSize/2, 0, 0.5)

			dc.SetFontFace(face)
			dc.SetColor(color.Black) // Ensure sample text is black
			dc.DrawStringAnchored("The quick brown fox jumps over the lazy dog", 500, float64(y)+fontSize/2, 0, 0.5)
			y += lineHeight + margin
		}

		// Generate output filename from folder path
		outName := folderPathToFilename(fontDir)
		err = dc.SavePNG(outName)
		if err != nil {
			fmt.Printf("failed to save image %s: %v\n", outName, err)
		} else {
			fmt.Printf("Saved: %s\n", outName)
		}
	}

	if len(failedFonts) > 0 {
		fmt.Println("Some fonts could not be processed:")
		for _, f := range failedFonts {
			fmt.Println(" -", f)
		}
	}

	return nil
}

// Helper to convert a folder path to a safe filename
func folderPathToFilename(path string) string {
	name := filepath.Clean(path)
	name = filepath.ToSlash(name)
	name = name[1:] // remove leading slash if present
	//name = filepath.Base(name)
	// Replace slashes and backslashes with underscores
	for _, c := range []string{"/", "\\", ":", " "} {
		name = strings.Replace(name, c, "_", -1)
	}
	return name + ".png"
}

// Helper to get a default font path for each OS
func getDefaultFontPath() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\Windows\Fonts\arial.ttf`
	case "darwin":
		return `/System/Library/Fonts/NewYorkItalic.ttf`
	default: // linux
		return `/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf`
	}
}

// package service

// import (
// 	"Metamorphoun/enum"
// 	"fmt"
// 	"image/color"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"

// 	"github.com/fogleman/gg"
// 	"golang.org/x/image/font"
// 	"golang.org/x/image/font/opentype"
// )

// func RenderFontsSample(outputPaths []string) error {
// 	if len(outputPaths) == 0 {
// 		return fmt.Errorf("no output paths provided")
// 	}
// 	for _, outputPath := range outputPaths {
// 		fontDir := GetFolderPath(enum.PathLoc.Fonts) //"C:\\Windows\\Fonts"
// 		sampleText := "The quick brown fox jumps over the lazy dog"
// 		fontSize := 14.0
// 		lineHeight := 50
// 		margin := 10

// 		// Get all font files in the specified directory
// 		files, err := ioutil.ReadDir(fontDir)
// 		if err != nil {
// 			return fmt.Errorf("failed to read font directory: %v", err)
// 		}

// 		// Filter font files
// 		var fontFiles []string
// 		for _, file := range files {
// 			if !file.IsDir() && (filepath.Ext(file.Name()) == ".ttf" || filepath.Ext(file.Name()) == ".otf") {
// 				fontFiles = append(fontFiles, filepath.Join(fontDir, file.Name()))
// 			}
// 		}

// 		// Create a new image context with enough height to fit all fonts
// 		imgHeight := (lineHeight + margin) * len(fontFiles)
// 		dc := gg.NewContext(2000, imgHeight)
// 		dc.SetColor(color.White)
// 		dc.Clear()
// 		dc.SetColor(color.Black)

// 		// Load the default font for printing file names
// 		defaultFontPath := "C:\\Windows\\Fonts\\DejaVuSans-Bold.ttf"
// 		defaultFontBytes, err := os.ReadFile(defaultFontPath)
// 		if err != nil {
// 			return fmt.Errorf("failed to read default font file: %v", err)
// 		}
// 		defaultFont, err := opentype.Parse(defaultFontBytes)
// 		if err != nil {
// 			return fmt.Errorf("failed to parse default font: %v", err)
// 		}
// 		defaultFace, err := opentype.NewFace(defaultFont, &opentype.FaceOptions{
// 			Size:    fontSize,
// 			DPI:     72,
// 			Hinting: font.HintingFull,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create default font face: %v", err)
// 		}

// 		y := margin
// 		for _, fontPath := range fontFiles {
// 			// Load the font
// 			fontBytes, err := os.ReadFile(fontPath)
// 			if err != nil {
// 				fmt.Println("Error reading font file:", err)
// 				continue
// 			}
// 			fnt, err := opentype.Parse(fontBytes)
// 			if err != nil {
// 				fmt.Println("Error parsing font:", err)
// 				continue
// 			}
// 			face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
// 				Size:    fontSize,
// 				DPI:     72,
// 				Hinting: font.HintingFull,
// 			})
// 			if err != nil {
// 				fmt.Println("Error creating font face:", err)
// 				continue
// 			}

// 			// Print the font file name using the default font
// 			dc.SetFontFace(defaultFace)
// 			dc.DrawStringAnchored(filepath.Base(fontPath), 100, float64(y)+fontSize/2, 0, 0.5)

// 			// Print the sample text using the current font
// 			dc.SetFontFace(face)
// 			dc.DrawStringAnchored(sampleText, 500, float64(y)+fontSize/2, 0, 0.5)

// 			y += lineHeight + margin
// 		}

// 		// Save the resulting image
// 		err = dc.SavePNG(outputPath)
// 		if err != nil {
// 			return fmt.Errorf("failed to save image: %v", err)
// 		}
// 	}
// 	return nil
// }
