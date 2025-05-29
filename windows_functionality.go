//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
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
