package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"Metamorphoun/config"
	"Metamorphoun/service"
)

func TestLiveDesktopWallpaperAssignment(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(home, "Library", "Application Support", "Metamorphoun", "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	paths := make([]string, len(cfg.DarwinPerScreenPicHistories))
	for i, pic := range cfg.DarwinPerScreenPicHistories {
		paths[i] = pic.SaveName
		if _, err := os.Stat(paths[i]); err != nil {
			t.Fatal(err)
		}
	}
	if len(paths) == 0 {
		t.Fatal("no current per-screen wallpaper assets")
	}
	if err := service.SetPerScreenWallpapers(paths); err != nil {
		t.Fatal(err)
	}
	output, err := exec.Command("swift", "-e", `import AppKit; import Foundation; let paths = NSScreen.screens.map { NSWorkspace.shared.desktopImageURL(for: $0)?.path ?? "" }; print(String(data: try JSONSerialization.data(withJSONObject: paths), encoding: .utf8)!)`).CombinedOutput()
	if err != nil {
		t.Fatalf("read live desktop URLs: %v: %s", err, output)
	}
	var actual []string
	if err := json.Unmarshal(output, &actual); err != nil {
		t.Fatalf("decode live desktop URLs: %v: %s", err, output)
	}
	if !reflect.DeepEqual(actual, paths) {
		t.Fatalf("desktop URLs = %v; want %v", actual, paths)
	}
	t.Logf("verified all %d desktop image URLs match their assigned assets", len(paths))
}
