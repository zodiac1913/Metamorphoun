package shared

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strings"
)

// This accesses the files stored inside this app (static files)
//
//go:embed static/*
var StaticFiles embed.FS

func ListStaticFiles() {
	files, err := fs.ReadDir(StaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		log.Println("Static file:", f.Name())
	}
}

// A function to retrieve files from staticFS
func GetStaticFSQuotes(filename string) ([]byte, error) {
	fmt.Println("Embedded files:", StaticFiles)
	filenameHasPrefixSlash := filename[0] == '/'
	filenamePrefix := "static/"
	if filenameHasPrefixSlash {
		filenamePrefix = "static"
	}
	if strings.Contains(filenamePrefix, "quotes") {
		filenamePrefix = "static/"
	}
	if strings.Contains(filenamePrefix, "/quotes") {
		filenamePrefix = "static"
	}
	fmt.Println("filenamePrefix:", filenamePrefix)
	fmt.Println("filename:", filename)
	data, err := StaticFiles.ReadFile(filenamePrefix + filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}
