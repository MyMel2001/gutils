package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// tar: create (-c) or extract (-x) tar archives, with optional gzip compression
func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "tar: usage: tar -c|-x -f ARCHIVE[.tar[.gz|.tgz]] [FILES...]")
		os.Exit(1)
	}
	mode := os.Args[1]
	archive := ""
	files := []string{}
	if os.Args[2] == "-f" {
		archive = os.Args[3]
		files = os.Args[4:]
	} else {
		fmt.Fprintln(os.Stderr, "tar: usage: tar -c|-x -f ARCHIVE [FILES...]")
		os.Exit(1)
	}
	isGz := strings.HasSuffix(archive, ".gz") || strings.HasSuffix(archive, ".tgz")
	if mode == "-c" {
		createTar(archive, files, isGz)
	} else if mode == "-x" {
		extractTar(archive, isGz)
	} else {
		fmt.Fprintln(os.Stderr, "tar: unknown mode", mode)
		os.Exit(1)
	}
}

func createTar(archive string, files []string, gzipIt bool) {
	f, err := os.Create(archive)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tar: cannot create archive:", err)
		os.Exit(1)
	}
	defer f.Close()
	var w io.Writer = f
	var gw *gzip.Writer
	if gzipIt {
		gw = gzip.NewWriter(f)
		defer gw.Close()
		w = gw
	}
	tw := tar.NewWriter(w)
	defer tw.Close()
	for _, file := range files {
		filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintln(os.Stderr, "tar: error:", err)
				return nil
			}
			hdr, err := tar.FileInfoHeader(info, "")
			if err != nil {
				fmt.Fprintln(os.Stderr, "tar: error:", err)
				return nil
			}
			hdr.Name = path
			err = tw.WriteHeader(hdr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "tar: error:", err)
				return nil
			}
			if !info.IsDir() {
				f, err := os.Open(path)
				if err != nil {
					fmt.Fprintln(os.Stderr, "tar: error:", err)
					return nil
				}
				io.Copy(tw, f)
				f.Close()
			}
			return nil
		})
	}
}

func extractTar(archive string, gzipIt bool) {
	f, err := os.Open(archive)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tar: cannot open archive:", err)
		os.Exit(1)
	}
	defer f.Close()
	var r io.Reader = f
	var gr *gzip.Reader
	if gzipIt {
		gr, err = gzip.NewReader(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tar: error reading gzip:", err)
			os.Exit(1)
		}
		defer gr.Close()
		r = gr
	}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "tar: error:", err)
			break
		}
		name := hdr.Name
		if strings.HasSuffix(name, "/") || hdr.FileInfo().IsDir() {
			os.MkdirAll(name, hdr.FileInfo().Mode())
			continue
		}
		os.MkdirAll(filepath.Dir(name), 0755)
		out, err := os.Create(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tar: error:", err)
			continue
		}
		io.Copy(out, tr)
		out.Chmod(hdr.FileInfo().Mode())
		out.Close()
	}
}
