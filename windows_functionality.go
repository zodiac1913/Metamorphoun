//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// Add to startup registry
const (
	runKeyCurrentUser = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName           = "Metamorphoun"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procBeep = kernel32.NewProc("Beep")
)

func PrintPlatformMessage() {
	fmt.Println("Running Windows-specific code")
}

func Beep(frequency, duration int) {
	procBeep.Call(uintptr(frequency), uintptr(duration))
}

func AddToStartup() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyCurrentUser, registry.WRITE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Ensure the path doesn't have any surrounding quotes that could cause issues
	exePath = strings.Trim(exePath, "\"")

	err = key.SetStringValue(appName, fmt.Sprintf("\"%s\"", exePath))
	if err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	log.Printf("%s added to Windows startup for the current user.", appName)
	return nil
}
func RemoveFromStartup() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyCurrentUser, registry.WRITE)
	if err != nil {
		// If the key doesn't exist, it's already removed or never added.
		if err == registry.ErrNotExist {
			log.Printf("%s startup entry not found for the current user.", appName)
			return nil
		}
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue(appName)
	if err != nil {
		// If the value doesn't exist, it's already removed.
		if err == registry.ErrNotExist {
			log.Printf("%s startup entry not found for the current user.", appName)
			return nil
		}
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	log.Printf("%s removed from Windows startup for the current user.", appName)
	return nil
}
func GetFolderPath(pathNeeded string) string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory:", err)
	}
	favPicFolderWithQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithQuotes")
	favPicFolderWithoutQuote := filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Pictures", "WithOutQuotes")

	if pathNeeded == "fonts" {
		return filepath.Join("C:", "Windows", "Fonts")
	} else if pathNeeded == "config" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun")
	} else if pathNeeded == "favorites" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites")
	} else if pathNeeded == "favwithquote" {
		return favPicFolderWithQuote
	} else if pathNeeded == "favwithoutquote" {
		return favPicFolderWithoutQuote
	} else if pathNeeded == "quotes" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Favorites", "Quotes")
	} else if pathNeeded == "configfile" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "config.json")
	} else if pathNeeded == "pictures" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Pictures")
	} else if pathNeeded == "logs" {
		return filepath.Join(usr.HomeDir, ".Metamorphoun", "Logs")
	} else if pathNeeded == "executable" {
		exePath, errEP := os.Executable()
		if errEP != nil {
			fmt.Println("Error:", errEP)
		}
		exeDir := filepath.Dir(exePath)
		return exeDir
	} else {
		return filepath.Join("C:", "Programs", "ZodiSoft", "Metamorphoun")
	}
}
