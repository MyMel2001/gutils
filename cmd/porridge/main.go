package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	repoDir              = "porridge_repo"
	instDir              = "porridge_installed"
	sourcesFile          = "porridge_sources.txt"
	installedSourcesFile = "porridge_installed_sources.txt"
	downloadedGoDir      = "downloaded_go"
	cacheDir             = "porridge_cache"
)

// downloadFile downloads a file from a URL and returns its contents or an error
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// getInstalledSource returns the source URL for a package, if any
func getInstalledSource(pkg string) string {
	data, err := os.ReadFile(installedSourcesFile)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 && parts[0] == pkg {
			return parts[1]
		}
	}
	return ""
}

// setInstalledSource records the source URL for a package
func setInstalledSource(pkg, url string) {
	data, _ := os.ReadFile(installedSourcesFile)
	lines := strings.Split(string(data), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, pkg+" ") {
			lines[i] = pkg + " " + url
			found = true
		}
	}
	if !found {
		lines = append(lines, pkg+" "+url)
	}
	out := strings.Join(lines, "\n")
	_ = os.WriteFile(installedSourcesFile, []byte(out), 0644)
}

// unzipFolderFromZip extracts a folder from a zip archive in memory to a destination directory
func unzipFolderFromZip(zipData []byte, folder, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}
	for _, f := range r.File {
		if !strings.HasPrefix(f.Name, folder+"/") {
			continue
		}
		relPath := strings.TrimPrefix(f.Name, folder+"/")
		if relPath == "" {
			continue
		}
		outPath := filepath.Join(dest, relPath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0755)
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		outFile, err := os.Create(outPath)
		if err != nil {
			inFile.Close()
			return err
		}
		io.Copy(outFile, inFile)
		inFile.Close()
		outFile.Close()
	}
	return nil
}

// getZipPackages returns a map of package name to zip URL for all sources, and caches zips
func getZipPackages() map[string]string {
	pkgs := make(map[string]string)
	os.MkdirAll(cacheDir, 0755)
	data, err := ioutil.ReadFile(sourcesFile)
	if err != nil {
		return pkgs
	}
	sources := strings.Split(string(data), "\n")
	for _, src := range sources {
		src = strings.TrimSpace(src)
		if src == "" {
			continue
		}
		cacheZip := filepath.Join(cacheDir, filepath.Base(src))
		var zipData []byte
		if _, err := os.Stat(cacheZip); err == nil {
			zipData, _ = ioutil.ReadFile(cacheZip)
		} else {
			zipData, err = downloadFile(src)
			if err != nil {
				continue
			}
			ioutil.WriteFile(cacheZip, zipData, 0644)
		}
		r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
		if err != nil {
			continue
		}
		for _, f := range r.File {
			if f.FileInfo().IsDir() && !strings.Contains(f.Name, "/") {
				pkgName := strings.TrimSuffix(f.Name, "/")
				pkgs[pkgName] = src
			}
		}
	}
	return pkgs
}

// porridge: minimal functional Portage-like package manager
func main() {
	os.MkdirAll(repoDir, 0755)
	os.MkdirAll(instDir, 0755)
	os.MkdirAll(downloadedGoDir, 0755)

	if len(os.Args) < 2 {
		fmt.Println("porridge: a Portage-like package manager")
		fmt.Println("Usage: porridge <command> [args]")
		fmt.Println("Commands: install, remove, search, sync, update, upgrade, fetchgo")
		os.Exit(0)
	}
	cmd := os.Args[1]
	switch cmd {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("porridge: install <package|url>")
			return
		}
		pkg := os.Args[2]
		if strings.HasPrefix(pkg, "http://") || strings.HasPrefix(pkg, "https://") {
			// Install from URL
			parts := strings.Split(pkg, "/")
			name := parts[len(parts)-1]
			instPkg := filepath.Join(instDir, name)
			if _, err := os.Stat(instPkg); err == nil {
				fmt.Printf("porridge: package '%s' already installed\n", name)
				return
			}
			data, err := downloadFile(pkg)
			if err != nil {
				fmt.Printf("porridge: failed to download '%s': %v\n", pkg, err)
				return
			}
			os.WriteFile(instPkg, data, 0644)
			setInstalledSource(name, pkg)
			fmt.Printf("porridge: installed '%s' from url\n", name)
			return
		}
		repoPkg := filepath.Join(repoDir, pkg)
		if _, err := os.Stat(repoPkg); err != nil {
			// Try zip sources
			zipPkgs := getZipPackages()
			if url, ok := zipPkgs[pkg]; ok {
				cacheZip := filepath.Join(cacheDir, filepath.Base(url))
				zipData, err := ioutil.ReadFile(cacheZip)
				if err != nil {
					fmt.Printf("porridge: failed to read cached zip for '%s': %v\n", pkg, err)
					return
				}
				instPkgDir := filepath.Join(instDir, pkg)
				if _, err := os.Stat(instPkgDir); err == nil {
					fmt.Printf("porridge: package '%s' already installed\n", pkg)
					return
				}
				os.MkdirAll(instPkgDir, 0755)
				err = unzipFolderFromZip(zipData, pkg, instPkgDir)
				if err != nil {
					fmt.Printf("porridge: failed to extract '%s' from zip: %v\n", pkg, err)
					return
				}
				setInstalledSource(pkg, url)
				fmt.Printf("porridge: installed '%s' from zip source\n", pkg)
				return
			}
			fmt.Printf("porridge: package '%s' not found in repo or zip sources\n", pkg)
			return
		}
		instPkg := filepath.Join(instDir, pkg)
		if _, err := os.Stat(instPkg); err == nil {
			fmt.Printf("porridge: package '%s' already installed\n", pkg)
			return
		}
		data, _ := os.ReadFile(repoPkg)
		os.WriteFile(instPkg, data, 0644)
		setInstalledSource(pkg, "local")
		fmt.Printf("porridge: installed '%s'\n", pkg)
	case "upgrade":
		if len(os.Args) < 3 {
			fmt.Println("porridge: upgrade <package>")
			return
		}
		pkg := os.Args[2]
		instPkg := filepath.Join(instDir, pkg)
		if _, err := os.Stat(instPkg); err != nil {
			fmt.Printf("porridge: package '%s' not installed\n", pkg)
			return
		}
		src := getInstalledSource(pkg)
		if src == "" || src == "local" {
			fmt.Printf("porridge: no remote source for '%s'\n", pkg)
			return
		}
		data, err := downloadFile(src)
		if err != nil {
			fmt.Printf("porridge: failed to download '%s': %v\n", src, err)
			return
		}
		os.WriteFile(instPkg, data, 0644)
		fmt.Printf("porridge: upgraded '%s'\n", pkg)
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("porridge: remove <package>")
			return
		}
		pkg := os.Args[2]
		instPkg := filepath.Join(instDir, pkg)
		if _, err := os.Stat(instPkg); err != nil {
			fmt.Printf("porridge: package '%s' not installed\n", pkg)
			return
		}
		os.Remove(instPkg)
		fmt.Printf("porridge: removed '%s'\n", pkg)
		// Remove from installed sources
		data, _ := os.ReadFile(installedSourcesFile)
		lines := strings.Split(string(data), "\n")
		newLines := []string{}
		for _, line := range lines {
			if !strings.HasPrefix(line, pkg+" ") {
				newLines = append(newLines, line)
			}
		}
		_ = os.WriteFile(installedSourcesFile, []byte(strings.Join(newLines, "\n")), 0644)
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("porridge: search <query>")
			return
		}
		query := strings.ToLower(os.Args[2])
		files, _ := os.ReadDir(repoDir)
		found := false
		for _, f := range files {
			if strings.Contains(strings.ToLower(f.Name()), query) {
				fmt.Println(f.Name())
				found = true
			}
		}
		// Search zip sources
		zipPkgs := getZipPackages()
		for name := range zipPkgs {
			if strings.Contains(strings.ToLower(name), query) {
				fmt.Println(name)
				found = true
			}
		}
		if !found {
			fmt.Println("porridge: no packages found")
		}
	case "sync":
		pkgs := []string{""}
		for _, pkg := range pkgs {
			f := filepath.Join(repoDir, pkg)
			os.WriteFile(f, []byte("package: "+pkg), 0644)
		}
		fmt.Println("porridge: repo synced")
		// Optionally sync meta repos (no-op, as we fetch live)
	case "update":
		files, _ := os.ReadDir(instDir)
		for _, f := range files {
			os.Chtimes(filepath.Join(instDir, f.Name()), time.Now(), time.Now())
			fmt.Printf("porridge: updated '%s'\n", f.Name())
		}
	case "fetchgo":
		if len(os.Args) < 3 {
			fmt.Println("porridge: fetchgo [--force] <url>")
			return
		}
		force := false
		urlIdx := 2
		if os.Args[2] == "--force" {
			force = true
			if len(os.Args) < 4 {
				fmt.Println("porridge: fetchgo [--force] <url>")
				return
			}
			urlIdx = 3
		}
		url := os.Args[urlIdx]
		if !strings.HasSuffix(url, ".go") {
			fmt.Println("porridge: only .go files can be fetched with fetchgo")
			return
		}
		parts := strings.Split(url, "/")
		name := parts[len(parts)-1]
		outPath := filepath.Join(downloadedGoDir, name)
		if _, err := os.Stat(outPath); err == nil && !force {
			fmt.Printf("porridge: file '%s' already exists in downloaded_go (use --force to overwrite)\n", name)
			return
		}
		data, err := downloadFile(url)
		if err != nil {
			fmt.Printf("porridge: failed to download '%s': %v\n", url, err)
			return
		}
		os.WriteFile(outPath, data, 0644)
		fmt.Printf("porridge: downloaded '%s' to downloaded_go/\n", name)
	default:
		fmt.Printf("porridge: unknown command '%s'\n", cmd)
	}
}
