//go:build linux
// +build linux

// linux_functionality.go
package main

import "fmt"

func PrintPlatformMessage() {
	fmt.Println("Running Linux-specific code")
}
