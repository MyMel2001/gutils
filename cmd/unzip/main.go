package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("unzip: usage: unzip <zipfile> <destination>")
		os.Exit(1)
	}
	zipPath := os.Args[1]
	dest := os.Args[2]
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Printf("unzip: failed to open zip: %v\n", err)
		os.Exit(1)
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		outPath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0755)
		inFile, err := f.Open()
		if err != nil {
			fmt.Printf("unzip: failed to open file in zip: %v\n", err)
			os.Exit(1)
		}
		outFile, err := os.Create(outPath)
		if err != nil {
			inFile.Close()
			fmt.Printf("unzip: failed to create file: %v\n", err)
			os.Exit(1)
		}
		_, err = io.Copy(outFile, inFile)
		inFile.Close()
		outFile.Close()
		if err != nil {
			fmt.Printf("unzip: failed to write file: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("unzip: extraction complete")
}
