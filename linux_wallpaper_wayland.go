//go:build linux
// +build linux

// linux_wallpaper_wayland.go
//
// Single-screen wallpaper setting for Wayland compositors.
//
// Probe order (first tool found and working wins):
//  1. omarchy-theme-bg-set  — Omarchy 4.x (Arch + Hyprland + Quickshell)
//  2. hyprctl hyprpaper     — Hyprland with hyprpaper daemon
//  3. swww                  — animated wallpaper daemon (wlroots compositors)
//  4. swaybg                — generic Wayland background tool (Sway, Wayfire, river…)
//  5. gsettings             — GNOME / Cinnamon on Wayland
//
// To add a new Wayland compositor:
//   - Add a probe block at the appropriate priority position below.
//   - If the compositor needs its own display-info query, add a
//     getDisplaysFrom<Compositor> function in linux_distro_<name>.go.
package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// setWaylandWallpaper iterates through known Wayland wallpaper tools and uses
// the first one that succeeds. All probes print a [wallpaper] log line so the
// terminal output makes it easy to see which tool was tried and why it failed.
func setWaylandWallpaper(filePath string) error {
	// -------------------------------------------------------------------------
	// 1. Omarchy / Quickshell (Omarchy 4.x — Arch Linux + Hyprland)
	//    omarchy-theme-bg-set updates ~/.local/state/omarchy/current/background
	//    and notifies the running Quickshell instance via IPC in one step.
	//    See: linux_distro_omarchy.go
	// -------------------------------------------------------------------------
	if _, err := exec.LookPath("omarchy-theme-bg-set"); err == nil {
		fmt.Println("[wallpaper] trying omarchy-theme-bg-set")
		out, err := exec.Command("omarchy-theme-bg-set", filePath).CombinedOutput()
		if err == nil {
			fmt.Println("[wallpaper] omarchy-theme-bg-set: OK")
			return nil
		}
		fmt.Printf("[wallpaper] omarchy-theme-bg-set failed: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Println("[wallpaper] omarchy-theme-bg-set: not found")
	}

	// -------------------------------------------------------------------------
	// 2. hyprpaper via hyprctl (Hyprland with hyprpaper daemon running)
	//    "reload" combines preload + set in one IPC call.
	//    Empty monitor field means "apply to all monitors".
	// -------------------------------------------------------------------------
	if path, err := exec.LookPath("hyprctl"); err == nil {
		fmt.Printf("[wallpaper] trying hyprctl (%s) hyprpaper reload\n", path)
		arg := "," + filePath
		out, err := exec.Command("hyprctl", "hyprpaper", "reload", arg).CombinedOutput()
		if err == nil {
			fmt.Println("[wallpaper] hyprctl hyprpaper reload: OK")
			return nil
		}
		fmt.Printf("[wallpaper] hyprctl hyprpaper reload failed: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Println("[wallpaper] hyprctl: not found")
	}

	// -------------------------------------------------------------------------
	// 3. swww (wlroots-based compositors — Wayfire, river, labwc…)
	//    swww-daemon must already be running (typically started at login).
	// -------------------------------------------------------------------------
	if path, err := exec.LookPath("swww"); err == nil {
		fmt.Printf("[wallpaper] trying swww (%s) img\n", path)
		out, err := exec.Command("swww", "img", filePath).CombinedOutput()
		if err == nil {
			fmt.Println("[wallpaper] swww img: OK")
			return nil
		}
		fmt.Printf("[wallpaper] swww img failed: %s\n", strings.TrimSpace(string(out)))
	} else {
		fmt.Println("[wallpaper] swww: not found")
	}

	// -------------------------------------------------------------------------
	// 4. swaybg (Sway, Wayfire, and any wlroots compositor without a daemon)
	//    swaybg must stay running to hold the background; we kill any prior
	//    instance first so only one copy owns the layer surface.
	// -------------------------------------------------------------------------
	if path, err := exec.LookPath("swaybg"); err == nil {
		fmt.Printf("[wallpaper] trying swaybg (%s)\n", path)
		_ = exec.Command("pkill", "-x", "swaybg").Run()
		cmd := exec.Command("swaybg", "-m", "fill", "-i", filePath)
		if err := cmd.Start(); err == nil {
			fmt.Println("[wallpaper] swaybg: OK")
			return nil
		} else {
			fmt.Printf("[wallpaper] swaybg start failed: %v\n", err)
		}
	} else {
		fmt.Println("[wallpaper] swaybg: not found")
	}

	// -------------------------------------------------------------------------
	// 5. gsettings (GNOME / Cinnamon on Wayland — GNOME Shell, Mutter)
	//    See: linux_distro_gnome.go
	// -------------------------------------------------------------------------
	fmt.Println("[wallpaper] trying gsettings (GNOME/Cinnamon on Wayland)")
	if err := tryGsettingsWallpaper(filePath); err == nil {
		fmt.Println("[wallpaper] gsettings: OK")
		return nil
	} else {
		fmt.Printf("[wallpaper] gsettings failed: %v\n", err)
	}

	return fmt.Errorf(
		"no Wayland wallpaper tool found " +
			"(tried: omarchy-theme-bg-set, hyprctl/hyprpaper, swww, swaybg, gsettings); " +
			"install the appropriate tool for your compositor")
}
