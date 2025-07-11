package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type lfsConfig struct {
	patterns map[string]bool
}

func readLFSConfig() (*lfsConfig, error) {
	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}
	attributesPath := filepath.Join(bruvPath, "..", ".bruvattributes")

	config := &lfsConfig{patterns: make(map[string]bool)}

	f, err := os.Open(attributesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // No .bruvattributes file, no LFS patterns
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "filter=lfs" {
			config.patterns[fields[0]] = true
		}
	}

	return config, scanner.Err()
}

func (c *lfsConfig) isLFS(path string) bool {
	for pattern := range c.patterns {
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func writeLFSPointer(content []byte) ([]byte, error) {
	lfsHash := sha256.Sum256(content)
	lfsHashStr := fmt.Sprintf("%x", lfsHash)

	pointerContent := fmt.Sprintf("version https://git-lfs.github.com/spec/v1\noid sha256:%s\nsize %d\n", lfsHashStr, len(content))

	// Store the actual file content in the LFS object store
	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}
	lfsObjectDir := filepath.Join(bruvPath, "lfs", "objects", lfsHashStr[:2], lfsHashStr[2:4])
	if err := os.MkdirAll(lfsObjectDir, 0755); err != nil {
		return nil, err
	}
	lfsObjectPath := filepath.Join(lfsObjectDir, lfsHashStr[4:])
	if err := os.WriteFile(lfsObjectPath, content, 0644); err != nil {
		return nil, err
	}

	return []byte(pointerContent), nil
} 