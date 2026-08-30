//go:build linux
// +build linux

// linux_wallpaper_multiscreen.go
//
// Multi-screen wallpaper support for Linux X11 sessions.
//
// Two backends are registered via linuxWallpaperModules in linux_functionality.go:
//
//	gsettings-composite  — GNOME / Cinnamon: builds one spanned PNG covering
//	                       the entire virtual desktop and sets it via gsettings.
//	                       Required because these desktops own the root window
//	                       and overwrite anything set by xwallpaper.
//
//	xwallpaper           — Bare X11 WMs (i3, bspwm, openbox…): sets each
//	                       output independently via --output flags.
//
// To add a new multi-screen backend:
//   - Implement a supports(linuxWallpaperContext) bool function.
//   - Implement an apply(linuxWallpaperContext, []string) error function.
//   - Append a linuxWallpaperModule entry to linuxWallpaperModules in
//     linux_functionality.go.
package main

import (
	"Metamorphoun/enum"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

// ---------------------------------------------------------------------------
// gsettings composite backend (GNOME / Cinnamon)
// ---------------------------------------------------------------------------

// supportsCinnamonGsettings reports whether the Cinnamon desktop is running.
// Cinnamon owns the root window through its settings daemon, so the
// composite-via-gsettings backend is required; xwallpaper does not stick.
func supportsCinnamonGsettings(ctx linuxWallpaperContext) bool {
	return hasGsettings() && desktopMatches(ctx, "cinnamon")
}

// supportsGnomeGsettings reports whether a GNOME-based desktop is running.
func supportsGnomeGsettings(ctx linuxWallpaperContext) bool {
	return hasGsettings() && desktopMatches(ctx, "gnome", "ubuntu", "unity")
}

func applyCinnamonGsettingsComposite(ctx linuxWallpaperContext, wallpaperPaths []string) error {
	return applyGsettingsComposite(ctx, wallpaperPaths, "org.cinnamon.desktop.background")
}

func applyGnomeGsettingsComposite(ctx linuxWallpaperContext, wallpaperPaths []string) error {
	return applyGsettingsComposite(ctx, wallpaperPaths, "org.gnome.desktop.background")
}

// applyGsettingsComposite composites each display's wallpaper onto a single
// canvas sized to the full virtual desktop, writes it to disk, then points
// the desktop background schema at it using the "spanned" layout.
func applyGsettingsComposite(ctx linuxWallpaperContext, wallpaperPaths []string, schema string) error {
	if len(wallpaperPaths) == 0 {
		return fmt.Errorf("no wallpaper paths provided")
	}
	if len(ctx.outputs) == 0 {
		return fmt.Errorf("no connected outputs detected")
	}

	compositePath, err := buildCompositeWallpaper(ctx.outputs, wallpaperPaths)
	if err != nil {
		return err
	}

	uri := "file://" + compositePath
	if out, err := exec.Command("gsettings", "set", schema, "picture-options", "spanned").CombinedOutput(); err != nil {
		return fmt.Errorf("gsettings set picture-options failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if out, err := exec.Command("gsettings", "set", schema, "picture-uri", uri).CombinedOutput(); err != nil {
		return fmt.Errorf("gsettings set picture-uri failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	// GNOME 42+ dark-theme key; non-fatal on Cinnamon (key absent).
	_ = exec.Command("gsettings", "set", schema, "picture-uri-dark", uri).Run()
	return nil
}

// buildCompositeWallpaper assembles a single image spanning the entire virtual
// desktop. Each output's assigned wallpaper is scaled to fill that monitor's
// rectangle and drawn at the monitor's pixel offset.
// Returns the path to the written composite PNG.
func buildCompositeWallpaper(outputs []linuxOutput, wallpaperPaths []string) (string, error) {
	maxX, maxY := 0, 0
	for _, out := range outputs {
		if out.width <= 0 || out.height <= 0 {
			return "", fmt.Errorf("output %s has unknown geometry; cannot composite", out.name)
		}
		if right := out.x + out.width; right > maxX {
			maxX = right
		}
		if bottom := out.y + out.height; bottom > maxY {
			maxY = bottom
		}
	}
	if maxX <= 0 || maxY <= 0 {
		return "", fmt.Errorf("invalid virtual desktop size %dx%d", maxX, maxY)
	}

	canvas := image.NewRGBA(image.Rect(0, 0, maxX, maxY))

	for index, out := range outputs {
		pathIndex := index
		if pathIndex >= len(wallpaperPaths) {
			pathIndex = len(wallpaperPaths) - 1
		}
		srcFile, err := os.Open(wallpaperPaths[pathIndex])
		if err != nil {
			return "", fmt.Errorf("failed to open wallpaper %s: %w", wallpaperPaths[pathIndex], err)
		}
		srcImg, _, decodeErr := image.Decode(srcFile)
		srcFile.Close()
		if decodeErr != nil {
			return "", fmt.Errorf("failed to decode wallpaper %s: %w", wallpaperPaths[pathIndex], decodeErr)
		}

		filled := imaging.Fill(srcImg, out.width, out.height, imaging.Center, imaging.Lanczos)
		destRect := image.Rect(out.x, out.y, out.x+out.width, out.y+out.height)
		draw.Draw(canvas, destRect, filled, image.Point{}, draw.Src)
	}

	compositePath := filepath.Join(GetFolderPath(enum.PathLoc.Config), "linux-composite-wallpaper.png")
	outFile, err := os.Create(compositePath)
	if err != nil {
		return "", fmt.Errorf("failed to create composite wallpaper file: %w", err)
	}
	defer outFile.Close()
	if err := png.Encode(outFile, canvas); err != nil {
		return "", fmt.Errorf("failed to encode composite wallpaper: %w", err)
	}
	return compositePath, nil
}

// ---------------------------------------------------------------------------
// xwallpaper backend (bare X11 WMs)
// ---------------------------------------------------------------------------

func supportsGnomeX11(ctx linuxWallpaperContext) bool {
	return supportsDesktopFamilyX11(ctx, "gnome", "ubuntu")
}

func supportsKDEX11(ctx linuxWallpaperContext) bool {
	return supportsDesktopFamilyX11(ctx, "kde", "plasma")
}

func supportsCinnamonX11(ctx linuxWallpaperContext) bool {
	return supportsDesktopFamilyX11(ctx, "cinnamon")
}

func supportsUbuntuFamilyX11(ctx linuxWallpaperContext) bool {
	if ctx.sessionType != "x11" {
		return false
	}
	if ctx.distroID == "ubuntu" || ctx.distroID == "linuxmint" || ctx.distroID == "mint" {
		return true
	}
	for _, like := range ctx.idLike {
		if like == "ubuntu" || like == "debian" {
			return true
		}
	}
	return false
}

func supportsDesktopFamilyX11(ctx linuxWallpaperContext, families ...string) bool {
	if ctx.sessionType != "x11" {
		return false
	}
	for _, desktop := range ctx.desktops {
		for _, family := range families {
			if desktop == family || strings.HasPrefix(desktop, family+"-") {
				return true
			}
		}
	}
	return false
}

func applyXwallpaperModule(ctx linuxWallpaperContext, wallpaperPaths []string) error {
	if _, err := exec.LookPath("xwallpaper"); err != nil {
		return fmt.Errorf("xwallpaper is required for Linux multi-screen wallpapers: %w", err)
	}
	if len(wallpaperPaths) == 0 {
		return fmt.Errorf("no wallpaper paths provided")
	}

	args := make([]string, 0, len(ctx.outputs)*4)
	for index, output := range ctx.outputs {
		pathIndex := index
		if pathIndex >= len(wallpaperPaths) {
			pathIndex = len(wallpaperPaths) - 1
		}
		args = append(args, "--output", output.name, "--zoom", wallpaperPaths[pathIndex])
	}

	out, err := exec.Command("xwallpaper", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("xwallpaper failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// ---------------------------------------------------------------------------
// xrandr output enumeration (shared by multi-screen and display-info paths)
// ---------------------------------------------------------------------------

// getXRandROutputs enumerates connected X11 displays with their pixel geometry.
func getXRandROutputs() ([]linuxOutput, error) {
	if _, err := exec.LookPath("xrandr"); err != nil {
		return nil, fmt.Errorf("xrandr is required to enumerate Linux displays: %w", err)
	}

	out, err := exec.Command("xrandr", "--query").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("xrandr failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	outputs := make([]linuxOutput, 0)
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, " connected") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		o := linuxOutput{name: fields[0]}
		for _, field := range fields[1:] {
			if w, h, x, y, ok := parseXRandRGeometry(field); ok {
				o.width, o.height, o.x, o.y = w, h, x, y
				break
			}
		}
		outputs = append(outputs, o)
	}
	return outputs, nil
}

// parseXRandRGeometry parses an xrandr geometry token "1920x1080+1920+0".
func parseXRandRGeometry(token string) (width, height, x, y int, ok bool) {
	xIdx := strings.IndexByte(token, 'x')
	plusIdx := strings.IndexByte(token, '+')
	if xIdx <= 0 || plusIdx <= xIdx {
		return 0, 0, 0, 0, false
	}
	rest := token[plusIdx:]
	offsets := strings.Split(strings.TrimPrefix(rest, "+"), "+")
	if len(offsets) != 2 {
		return 0, 0, 0, 0, false
	}
	w, errW := strconv.Atoi(token[:xIdx])
	h, errH := strconv.Atoi(token[xIdx+1 : plusIdx])
	px, errX := strconv.Atoi(offsets[0])
	py, errY := strconv.Atoi(offsets[1])
	if errW != nil || errH != nil || errX != nil || errY != nil {
		return 0, 0, 0, 0, false
	}
	if w <= 0 || h <= 0 {
		return 0, 0, 0, 0, false
	}
	return w, h, px, py, true
}
