package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type ignoreConfig struct {
	patterns []string
}

func readBruvIgnore() (*ignoreConfig, error) {
	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}
	
	ignorePath := filepath.Join(bruvPath, "..", ".bruvignore")
	
	config := &ignoreConfig{patterns: []string{}}
	
	f, err := os.Open(ignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // No .bruvignore file, no ignore patterns
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
		config.patterns = append(config.patterns, line)
	}
	
	return config, scanner.Err()
}

func (c *ignoreConfig) isIgnored(path string) bool {
	for _, pattern := range c.patterns {
		// Check if the pattern matches the filename
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}
		
		// Check if the pattern matches the full path
		matched, err = filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}
		
		// Check if the path contains the pattern as a directory
		if strings.Contains(path, pattern+string(filepath.Separator)) {
			return true
		}
	}
	
	// Always ignore .bruv directory
	if strings.HasPrefix(path, ".bruv"+string(filepath.Separator)) {
		return true
	}
	
	return false
}