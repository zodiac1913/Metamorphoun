package service

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

func RepairRetainedHistoryAssets() error {
	cfg := config.GetConfig()
	if cfg == nil || len(cfg.PicHistories) == 0 {
		return nil
	}

	limit := len(cfg.PicHistories)
	if limit > 10 {
		limit = 10
	}

	for historyIndex := 0; historyIndex < limit; historyIndex++ {
		if err := repairPicHistoryEntry(&cfg.PicHistories[historyIndex]); err != nil {
			return fmt.Errorf("repair history %d: %w", historyIndex, err)
		}
	}

	return nil
}

func repairPicHistoryEntry(pic *config.PicHistory) error {
	if pic == nil {
		return nil
	}
	if len(pic.PerScreenPics) > 0 {
		for index := range pic.PerScreenPics {
			if err := repairAlteredPicFile(&pic.PerScreenPics[index]); err != nil {
				return err
			}
		}
	}
	return repairAlteredPicFile(pic)
}

func repairAlteredPicFile(pic *config.PicHistory) error {
	if pic == nil || pic.SaveName == "" {
		return nil
	}
	if strings.HasPrefix(strings.ToLower(pic.SaveName), "http") {
		return nil
	}
	if _, err := os.Stat(pic.SaveName); err == nil {
		return nil
	}

	img, err := rebuildAlteredImage(*pic)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(pic.SaveName), 0700); err != nil {
		return err
	}
	saveImg(img, pic.SaveName)
	return nil
}

func rebuildAlteredImage(currentPic config.PicHistory) (image.Image, error) {
	img, err := backgroundSetSource(currentPic)
	if err != nil {
		return nil, err
	}
	if img == nil {
		return nil, fmt.Errorf("history source image is empty")
	}

	img, currentPic = handleScaling(img, currentPic, currentPic.Sizing, err)
	if img == nil {
		return nil, fmt.Errorf("history image is empty after scaling")
	}

	img, err = filterCurrentPic(currentPic, img)
	if err != nil {
		return nil, err
	}
	if img == nil {
		return nil, fmt.Errorf("history image is empty after filtering")
	}

	if shouldDrawStoredQuote(currentPic) {
		img, err = drawStoredQuote(currentPic, img)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

func shouldDrawStoredQuote(currentPic config.PicHistory) bool {
	if currentPic.QuoteStatement == "" {
		return false
	}
	if currentPic.ImageItem.Name == "Favorites" && strings.Contains(currentPic.OriginName, "WithQuotes") {
		return false
	}
	return true
}

func drawStoredQuote(currentPic config.PicHistory, img image.Image) (image.Image, error) {
	dc := gg.NewContextForImage(img)
	fontPath := currentPic.QuoteFont
	if fontPath == "" {
		fontPath = filepath.Join(GetFolderPath(enum.PathLoc.Fonts), config.ConfigInstance.TextFontFile)
	}
	fontSize := currentPic.QuoteFontSize
	if fontSize < 8 {
		fontSize = math.Max(config.ConfigInstance.QuoteFontSizeMin, 16)
	}
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return nil, err
	}

	textBoxWidth := currentPic.QuoteTextBoxWidth
	textBoxHeight := currentPic.QuoteTextBoxHeight
	textBlockX := currentPic.QuoteTextBoxX
	textBlockY := currentPic.QuoteTextBoxY
	if textBoxWidth < 40 {
		textBoxWidth = 420
	}
	if textBoxHeight < 40 {
		textBoxHeight = 140
	}

	dc.SetColor(color.RGBA{
		R: currentPic.QuoteBackgroundColorR,
		G: currentPic.QuoteBackgroundColorG,
		B: currentPic.QuoteBackgroundColorB,
		A: uint8(currentPic.QuoteOpacity),
	})
	dc.DrawRoundedRectangle(textBlockX, textBlockY, textBoxWidth, textBoxHeight, 10)
	dc.Fill()
	DrawCyberpunkQuoteBorder(currentPic, dc, textBlockX, textBlockY, textBoxWidth, textBoxHeight)

	dc.SetColor(color.RGBA{
		R: currentPic.QuoteTextColorR,
		G: currentPic.QuoteTextColorG,
		B: currentPic.QuoteTextColorB,
		A: 255,
	})
	wrappedLines := dc.WordWrap(`"`+currentPic.QuoteStatement+`"`, textBoxWidth-40)
	DrawQuoteText(dc, wrappedLines, currentPic.QuoteAuthor, textBlockX, textBlockY, textBoxWidth)
	return dc.Image(), nil
}
