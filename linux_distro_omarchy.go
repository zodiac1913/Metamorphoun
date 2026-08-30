//go:build linux
// +build linux

// linux_distro_omarchy.go
//
// Omarchy-specific support (Arch Linux + Hyprland + Quickshell).
//
// Omarchy 4.x does not use hyprpaper, swww, swaybg, feh, or any traditional
// wallpaper daemon.  The Quickshell desktop shell owns the background layer
// directly via the wlr-layer-shell Wayland protocol.  Wallpapers are set by:
//
//  1. Updating the symlink ~/.local/state/omarchy/current/background
//  2. Notifying the live Quickshell process via its IPC socket
//
// omarchy-theme-bg-set (in /usr/share/omarchy/bin) does both steps atomically,
// so that is the tool called by linux_wallpaper_wayland.go.
//
// This file provides display-geometry queries for Omarchy via hyprctl, which
// avoids the XGB/X11 errors produced by kbinani/screenshot on pure Wayland.
package main

import (
	"Metamorphoun/service"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// isOmarchy reports whether the running system is Omarchy (Arch + Hyprland +
// Quickshell). Detection is intentionally lenient and checks several signals:
//
//   - the omarchy CLI is on PATH, or
//   - /usr/share/omarchy exists (the Omarchy install tree), or
//   - OMARCHY_PATH is set in the environment.
//
// Omarchy's Quickshell background layer accepts only a single image, so the
// "different picture per screen" feature cannot work there. Callers use this
// to disable that option.
func isOmarchy() bool {
	if _, err := exec.LookPath("omarchy-theme-bg-set"); err == nil {
		return true
	}
	if os.Getenv("OMARCHY_PATH") != "" {
		return true
	}
	if info, err := os.Stat("/usr/share/omarchy"); err == nil && info.IsDir() {
		return true
	}
	// As a final signal, check /etc/os-release for an omarchy marker.
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lower := strings.ToLower(string(data))
		if strings.Contains(lower, "omarchy") {
			return true
		}
	}
	return false
}

// getDisplaysFromHyprctl queries Hyprland's IPC for monitor geometry.
// Called by getLinuxDisplayInfo in linux_functionality.go when the session
// is detected as Wayland.
//
// Uses `hyprctl monitors -j` which returns a JSON array; we only need
// width and height per monitor.
func getDisplaysFromHyprctl() ([]service.ScreenInfo, error) {
	if _, err := exec.LookPath("hyprctl"); err != nil {
		return nil, fmt.Errorf("hyprctl not found")
	}

	out, err := exec.Command("hyprctl", "monitors", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("hyprctl monitors -j failed: %w", err)
	}

	type hyprMonitor struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	var monitors []hyprMonitor
	if err := json.Unmarshal(out, &monitors); err != nil {
		return nil, fmt.Errorf("hyprctl monitors JSON parse failed: %w", err)
	}

	screens := make([]service.ScreenInfo, 0, len(monitors))
	for i, m := range monitors {
		if m.Width > 0 && m.Height > 0 {
			screens = append(screens, service.ScreenInfo{
				Number: int16(i),
				Width:  m.Width,
				Height: m.Height,
			})
		}
	}
	fmt.Printf("[display] hyprctl: found %d monitor(s)\n", len(screens))
	return screens, nil
}

// getOutputsFromHyprctl queries Hyprland for full monitor geometry including
// pixel offsets (x, y) and scale. Unlike getDisplaysFromHyprctl this returns
// linuxOutput values so a spanned composite covering the whole virtual desktop
// can be built for the per-screen wallpaper mode.
func getOutputsFromHyprctl() ([]linuxOutput, error) {
	if _, err := exec.LookPath("hyprctl"); err != nil {
		return nil, fmt.Errorf("hyprctl not found")
	}

	out, err := exec.Command("hyprctl", "monitors", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("hyprctl monitors -j failed: %w", err)
	}

	// Hyprland's width/height are the raw mode resolution; x/y are the layout
	// offset; scale divides the mode into logical pixels. The background layer
	// is sized in logical pixels, so divide by scale to match what Quickshell
	// actually paints.
	type hyprMonitor struct {
		Name   string  `json:"name"`
		Width  int     `json:"width"`
		Height int     `json:"height"`
		X      int     `json:"x"`
		Y      int     `json:"y"`
		Scale  float64 `json:"scale"`
	}
	var monitors []hyprMonitor
	if err := json.Unmarshal(out, &monitors); err != nil {
		return nil, fmt.Errorf("hyprctl monitors JSON parse failed: %w", err)
	}

	outputs := make([]linuxOutput, 0, len(monitors))
	for _, m := range monitors {
		scale := m.Scale
		if scale <= 0 {
			scale = 1
		}
		w := int(float64(m.Width) / scale)
		h := int(float64(m.Height) / scale)
		if w <= 0 || h <= 0 {
			continue
		}
		outputs = append(outputs, linuxOutput{
			name:   m.Name,
			width:  w,
			height: h,
			x:      m.X,
			y:      m.Y,
		})
	}
	return outputs, nil
}
