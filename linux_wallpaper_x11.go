//go:build linux
// +build linux

// linux_wallpaper_x11.go
//
// Single-screen wallpaper setting for X11 sessions.
//
// Probe order (first tool found and working wins):
//  1. gsettings   — GNOME / Cinnamon / Budgie (X11 session)
//  2. xwallpaper  — lightweight, supports modern X11 WMs (i3, bspwm, openbox…)
//  3. feh         — classic; many minimal setups rely on it
//  4. nitrogen    — common graphical X11 wallpaper manager
//
// To add a new X11 desktop/WM:
//   - Add a probe block at the appropriate priority position below.
//   - If it needs its own schema or tool, add a linux_distro_<name>.go file.
package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// setX11Wallpaper iterates through known X11 wallpaper tools and uses the
// first one that succeeds.
func setX11Wallpaper(filePath string) error {
	// -------------------------------------------------------------------------
	// 1. gsettings — GNOME / Cinnamon / Budgie on X11
	//    Schema detection is handled in linux_distro_gnome.go.
	// -------------------------------------------------------------------------
	fmt.Println("[wallpaper] trying gsettings (GNOME/Cinnamon/Budgie on X11)")
	if err := tryGsettingsWallpaper(filePath); err == nil {
		fmt.Println("[wallpaper] gsettings: OK")
		return nil
	} else {
		fmt.Printf("[wallpaper] gsettings: %v\n", err)
	}

	// -------------------------------------------------------------------------
	// 2. xwallpaper — bare X11 WMs (i3, bspwm, openbox, dwm…)
	// -------------------------------------------------------------------------
	if path, err := exec.LookPath("xwallpaper"); err == nil {
		fmt.Printf("[wallpaper] trying xwallpaper (%s)\n", path)
		out, err := exec.Command("xwallpaper", "--zoom", filePath).CombinedOutput()
		if err == nil {
			fmt.Println("[wallpaper] xwallpaper: OK")
			return nil
		}
		fmt.Printf("[wallpaper] xwallpaper failed: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Println("[wallpaper] xwallpaper: not found")
	}

	// -------------------------------------------------------------------------
	// 3. feh — classic wallpaper setter; present on many minimal distros
	// -------------------------------------------------------------------------
	if path, err := exec.LookPath("feh"); err == nil {
		fmt.Printf("[wallpaper] trying feh (%s)\n", path)
		out, err := exec.Command("feh", "--bg-scale", filePath).CombinedOutput()
		if err == nil {
			fmt.Println("[wallpaper] feh: OK")
			return nil
		}
		fmt.Printf("[wallpaper] feh failed: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Println("[wallpaper] feh: not found")
	}

	// -------------------------------------------------------------------------
	// 4. nitrogen — graphical wallpaper manager common on Xfce / LXDE setups
	// -------------------------------------------------------------------------
	if path, err := exec.LookPath("nitrogen"); err == nil {
		fmt.Printf("[wallpaper] trying nitrogen (%s)\n", path)
		out, err := exec.Command("nitrogen", "--set-scaled", filePath).CombinedOutput()
		if err == nil {
			fmt.Println("[wallpaper] nitrogen: OK")
			return nil
		}
		fmt.Printf("[wallpaper] nitrogen failed: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Println("[wallpaper] nitrogen: not found")
	}

	return fmt.Errorf(
		"no X11 wallpaper tool found " +
			"(tried: gsettings, xwallpaper, feh, nitrogen); " +
			"install the appropriate tool for your desktop environment")
}
