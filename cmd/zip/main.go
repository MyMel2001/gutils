package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func addFilesToZip(zipWriter *zip.Writer, basePath, baseInZip string) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		inZipPath := filepath.Join(baseInZip, relPath)
		if info.IsDir() {
			if inZipPath != "" {
				_, err := zipWriter.Create(inZipPath + "/")
				return err
			}
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		w, err := zipWriter.Create(inZipPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, file)
		return err
	})
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("zip: usage: zip <folder> <zipfile>")
		os.Exit(1)
	}
	folder := os.Args[1]
	zipPath := os.Args[2]
	zipFile, err := os.Create(zipPath)
	if err != nil {
		fmt.Printf("zip: failed to create zip: %v\n", err)
		os.Exit(1)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	err = addFilesToZip(zipWriter, folder, "")
	if err != nil {
		fmt.Printf("zip: failed to add files: %v\n", err)
		os.Exit(1)
	}
	err = zipWriter.Close()
	if err != nil {
		fmt.Printf("zip: failed to close zip: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("zip: archive created successfully")
}
