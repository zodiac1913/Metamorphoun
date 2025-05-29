//go:build windows
// +build windows

package main

import "fmt"

func PrintPlatformMessage() {
	fmt.Println("Running Windows-specific code")
}
