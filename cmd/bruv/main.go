package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"net/url"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "bruv: missing command\n")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "init":
		if err := cmdInit(args); err != nil {
			fmt.Fprintf(os.Stderr, "bruv: %s\n", err)
			os.Exit(1)
		}
	case "hash-object":
		if err := cmdHashObject(args); err != nil {
			fmt.Fprintf(os.Stderr, "bruv: %s\n", err)
			os.Exit(1)
		}
	case "cat-file":
		if err := cmdCatFile(args); err != nil {
			fmt.Fprintf(os.Stderr, "bruv: %s\n", err)
			os.Exit(1)
		}
	case "add":
		if err := cmdAdd(args); err != nil {
			fmt.Fprintf(os.Stderr, "bruv: %s\n", err)
			os.Exit(1)
		}
	case "commit":
		if err := cmdCommit(args); err != nil {
			fmt.Fprintf(os.Stderr, "bruv: %s\n", err)
			os.Exit(1)
		}
	case "serve":
		if err := cmdServe(); err != nil {
			fmt.Fprintf(os.Stderr, "bruv: %s\n", err)
			os.Exit(1)
		}
	case "clone":
		if err := cmdClone(args); err != nil {
			fmt.Fprintf(os.Stderr, "bruv: %s\n", err)
			os.Exit(1)
		}
	case "help":
		cmdHelp()
	default:
		fmt.Fprintf(os.Stderr, "bruv: unknown command: %s\n", command)
		cmdHelp()
		os.Exit(1)
	}
}

func cmdClone(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: bruv clone <url> <directory>")
	}
	url := args[0]
	dir := args[1]

	parsedURL, err := url.Parse(url)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	host := parsedURL.Host
	if host == "" {
		host = parsedURL.Path
	}
	repoPath := strings.TrimPrefix(parsedURL.Path, "/")

	conn, err := net.Dial("tcp", host+":9418")
	if err != nil {
		return fmt.Errorf("failed to connect to %s:9418: %w", host, err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "git-upload-pack %s\n", repoPath)

	// We're expecting a ref and then the packfile
	reader := bufio.NewReader(conn)
	refLine, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read ref from server: %w", err)
	}
	if refLine == "" {
		return fmt.Errorf("empty ref line received from server")
	}
	fmt.Print(refLine) // "d350900... refs/heads/main"

	// Initialize the new repository
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// We call cmdInit with the directory argument
	if err := cmdInit([]string{dir}); err != nil {
		return err
	}
	
	bruvPath := filepath.Join(dir, ".bruv")
	packDir := filepath.Join(bruvPath, "objects", "pack")
	if err := os.MkdirAll(packDir, 0755); err != nil {
		return err
	}
	packfilePath := filepath.Join(packDir, "temp.pack")

	packfile, err := os.Create(packfilePath)
	if err != nil {
		return err
	}
	defer packfile.Close()

	n, err := io.Copy(packfile, reader)
	if err != nil {
		return fmt.Errorf("failed to copy packfile data: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("no data received for packfile")
	}

	fmt.Printf("Clone complete. Packfile saved to %s\n", packfilePath)

	// Unpack the packfile
	if err := unpackPackfile(packfilePath, bruvPath); err != nil {
		return fmt.Errorf("failed to unpack packfile: %w", err)
	}
	fmt.Println("Packfile unpacked.")

	// Update HEAD and refs
	refParts := strings.Fields(refLine)
	if len(refParts) < 1 {
		return fmt.Errorf("invalid ref line received: %s", refLine)
	}
	commitHash := refParts[0]
	if err := updateRefsAfterClone(bruvPath, commitHash); err != nil {
		return fmt.Errorf("failed to update refs: %w", err)
	}

	fmt.Println("Repository successfully cloned.")

	return nil
}

func cmdCatFile(args []string) error {
	printType := false
	prettyPrint := false
	var hashStr string

	for _, arg := range args {
		switch arg {
		case "-t":
			printType = true
		case "-p":
			prettyPrint = true
		default:
			if !strings.HasPrefix(arg, "-") {
				hashStr = arg
			}
		}
	}

	if hashStr == "" {
		return fmt.Errorf("usage: bruv cat-file [-t | -p] <object>")
	}

	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}

	objectPath := filepath.Join(bruvPath, "objects", hashStr[:2], hashStr[2:])
	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return err
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid object format")
	}

	header := parts[0]
	content := parts[1]

	var objType string
	var size int
	fmt.Sscanf(string(header), "%s %d", &objType, &size)

	if printType {
		fmt.Println(objType)
		return nil
	}

	if prettyPrint {
		switch objType {
		case "blob":
			fmt.Print(string(content))
		case "tree":
			// This is a simplified tree listing
			buf := bytes.NewBuffer(content)
			for {
				mode, err := buf.ReadString(' ')
				if err != nil {
					break
				}
				path, err := buf.ReadString(0)
				if err != nil {
					break
				}
				var hash [20]byte
				if _, err := io.ReadFull(buf, hash[:]); err != nil {
					break
				}
				fmt.Printf("%s %s %x\t%s\n", mode[:len(mode)-1], "blob", hash, path[:len(path)-1])
			}
		case "commit":
			fmt.Print(string(content))
		default:
			return fmt.Errorf("unsupported object type for pretty-print: %s", objType)
		}
		return nil
	}

	// For raw content, just print it
	fmt.Print(string(decompressed))

	return nil
}

func cmdHashObject(args []string) error {
	write := false
	paths := []string{}
	for _, arg := range args {
		if arg == "-w" {
			write = true
		} else if !strings.HasPrefix(arg, "-") {
			paths = append(paths, arg)
		}
	}

	if len(paths) < 1 {
		return fmt.Errorf("usage: bruv hash-object [-w] <file>")
	}
	path := paths[0]

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	hasher := sha1.New()
	header := []byte(fmt.Sprintf("blob %d\x00", len(content)))
	hasher.Write(header)
	hasher.Write(content)
	hash := hasher.Sum(nil)

	fmt.Printf("%x\n", hash)

	if write {
		if _, err := writeBlobObject(content); err != nil {
			return err
		}
	}

	return nil
}

func findBruvDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		bruvPath := filepath.Join(dir, ".bruv")
		if _, err := os.Stat(bruvPath); err == nil {
			return bruvPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not a bruv repository (or any of the parent directories): .bruv")
		}
		dir = parent
	}
}

func cmdInit(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	bruvPath := filepath.Join(dir, ".bruv")
	if _, err := os.Stat(bruvPath); !os.IsNotExist(err) {
		return fmt.Errorf("bruv repository already exists in %s", bruvPath)
	}

	fmt.Printf("Initializing empty Bruv repository in %s\n", bruvPath)

	dirs := []string{
		"objects",
		"refs/heads",
		"refs/tags",
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(bruvPath, d), 0755); err != nil {
			return err
		}
	}

	headPath := filepath.Join(bruvPath, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return err
	}

	// Create a minimal config file
	configPath := filepath.Join(bruvPath, "config")
	configContent := []byte("[core]\n\trepositoryformatversion = 0\n\tfilemode = true\n\tbare = false\n")
	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		return err
	}

	return nil
} 