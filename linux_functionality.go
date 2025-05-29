//go:build linux
// +build linux

// linux_functionality.go
package main

import "fmt"

func PrintPlatformMessage() {
	fmt.Println("Running Linux-specific code")
}

func AddToStartup() error {
	return nil
}
func RemoveFromStartup() error {
	return nil
}
