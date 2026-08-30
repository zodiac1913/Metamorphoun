//go:build linux
// +build linux

// linux_distro_gnome.go
//
// GNOME-family desktop support (GNOME Shell, Cinnamon, Budgie, Unity, Ubuntu).
//
// These desktops manage the root window through their settings daemon and set
// the wallpaper via the gsettings background schema, on both X11 and Wayland.
// Single-screen setting is done here; multi-screen compositing lives in
// linux_wallpaper_multiscreen.go (which reuses the same schema names).
package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// tryGsettingsWallpaper sets the wallpaper via gsettings, covering GNOME,
// Cinnamon, Budgie, and Ubuntu on both X11 and Wayland. Returns nil on
// success; a non-nil error if gsettings is unavailable or no recognised
// background schema exists.
func tryGsettingsWallpaper(filePath string) error {
	if _, err := exec.LookPath("gsettings"); err != nil {
		return fmt.Errorf("gsettings not found")
	}

	schema := gnomeLikeSchema()
	if schema == "" {
		return fmt.Errorf("no recognised gsettings background schema")
	}

	uri := "file://" + filePath
	if out, err := exec.Command("gsettings", "set", schema, "picture-options", "zoom").CombinedOutput(); err != nil {
		return fmt.Errorf("gsettings set picture-options: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if out, err := exec.Command("gsettings", "set", schema, "picture-uri", uri).CombinedOutput(); err != nil {
		return fmt.Errorf("gsettings set picture-uri: %w: %s", err, strings.TrimSpace(string(out)))
	}
	// GNOME 42+ dark-theme variant; non-fatal if the key doesn't exist (Cinnamon).
	_ = exec.Command("gsettings", "set", schema, "picture-uri-dark", uri).Run()
	return nil
}

// gnomeLikeSchema returns the gsettings background schema name for the current
// desktop, or an empty string if none is recognised. It first checks the
// desktop-environment env vars, then falls back to probing the schemas that
// gsettings actually exposes (so it works even when XDG_CURRENT_DESKTOP is
// unset, e.g. when launched from a service).
func gnomeLikeSchema() string {
	desktops := detectLinuxDesktops()
	for _, de := range desktops {
		switch {
		case de == "cinnamon":
			return "org.cinnamon.desktop.background"
		case de == "gnome", de == "ubuntu", de == "unity", de == "budgie",
			strings.HasPrefix(de, "gnome-"), strings.HasPrefix(de, "ubuntu-"):
			return "org.gnome.desktop.background"
		}
	}
	out, err := exec.Command("gsettings", "list-schemas").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "org.gnome.desktop.background" {
			return "org.gnome.desktop.background"
		}
		if line == "org.cinnamon.desktop.background" {
			return "org.cinnamon.desktop.background"
		}
	}
	return ""
}

// hasGsettings reports whether the gsettings tool is available.
func hasGsettings() bool {
	_, err := exec.LookPath("gsettings")
	return err == nil
}

// desktopMatches reports whether any detected desktop belongs to one of the
// given families. Does not require X11, so it works for both X11 and Wayland
// GNOME-family sessions.
func desktopMatches(ctx linuxWallpaperContext, families ...string) bool {
	for _, desktop := range ctx.desktops {
		for _, family := range families {
			if desktop == family || strings.HasPrefix(desktop, family+"-") {
				return true
			}
		}
	}
	return false
}
