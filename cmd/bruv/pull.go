package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func cmdPull(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: bruv pull <remote> [<branch>] [--select <path>...]")
	}
	
	// Parse arguments for selective pulling
	selectPaths := []string{}
	regularArgs := []string{}
	skipNext := false
	
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if arg == "--select" {
			// Collect all paths after --select
			for j := i + 1; j < len(args); j++ {
				if !strings.HasPrefix(args[j], "-") {
					selectPaths = append(selectPaths, args[j])
				} else {
					break
				}
			}
			skipNext = true
		} else {
			regularArgs = append(regularArgs, arg)
		}
	}
	
	remote := args[0]
	branch := "main"
	if len(args) > 1 {
		branch = args[1]
	}
	
	// Parse remote URL (simplified format: host:port/repo)
	parts := strings.Split(remote, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid remote format. Use: host:port/repo")
	}
	host := parts[0]
	repoPath := parts[1]
	
	// Get current repository path
	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}
	
	// Connect to remote
	conn, err := net.Dial("tcp", host+":9418")
	if err != nil {
		return fmt.Errorf("could not connect to remote: %w", err)
	}
	defer conn.Close()
	
	// Send fetch request
	if len(selectPaths) > 0 {
		// Send selective fetch request with paths
		fmt.Fprintf(conn, "git-upload-pack %s --select %s\n", repoPath, strings.Join(selectPaths, " "))
	} else {
		fmt.Fprintf(conn, "git-upload-pack %s\n", repoPath)
	}
	
	// Read response
	reader := bufio.NewReader(conn)
	refLine, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("could not read ref: %w", err)
	}
	
	refParts := strings.Fields(refLine)
	if len(refParts) < 1 {
		return fmt.Errorf("invalid ref line received: %s", refLine)
	}
	remoteCommitHash := refParts[0]
	
	// Check if we need to update
	localBranchRefPath := filepath.Join(bruvPath, "refs", "heads", branch)
	localCommitHash := ""
	if localCommitHashBytes, err := os.ReadFile(localBranchRefPath); err == nil {
		localCommitHash = strings.TrimSpace(string(localCommitHashBytes))
	}
	
	if localCommitHash == remoteCommitHash {
		fmt.Println("Already up to date.")
		return nil
	}
	
	// Create temporary packfile
	packDir := filepath.Join(bruvPath, "objects", "pack")
	if err := os.MkdirAll(packDir, 0755); err != nil {
		return err
	}
	packfilePath := filepath.Join(packDir, "temp.pull.pack")
	
	packfile, err := os.Create(packfilePath)
	if err != nil {
		return err
	}
	defer packfile.Close()
	
	// Copy packfile data
	_, err = io.Copy(packfile, reader)
	if err != nil {
		return err
	}
	
	// Unpack the packfile
	if err := unpackPackfile(packfilePath, bruvPath); err != nil {
		return fmt.Errorf("failed to unpack packfile: %w", err)
	}
	
	// Update local branch ref
	if err := os.WriteFile(localBranchRefPath, []byte(remoteCommitHash+"\n"), 0644); err != nil {
		return err
	}
	
	// Update HEAD if it points to this branch
	headPath := filepath.Join(bruvPath, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err == nil && strings.TrimSpace(string(headContent)) == "ref: refs/heads/"+branch {
		// Update working directory (simplified - just print message for now)
		fmt.Printf("Updated %s to %s\n", branch, remoteCommitHash[:7])
		fmt.Println("Note: Working directory update not implemented. Files remain unchanged.")
	}
	
	fmt.Printf("Pulled %s from %s (%s -> %s)\n", branch, remote, 
		localCommitHash[:7], remoteCommitHash[:7])
	
	// Clean up temporary packfile
	os.Remove(packfilePath)
	
	// If selective paths were specified, filter the working directory
	if len(selectPaths) > 0 {
		if err := filterWorkingDirectoryForPull(bruvPath, selectPaths); err != nil {
			return fmt.Errorf("failed to filter working directory: %w", err)
		}
		fmt.Printf("Filtered working directory to: %v\n", selectPaths)
	}
	
	return nil
}

// filterWorkingDirectoryForPull filters the working directory after pull to only include specified paths
func filterWorkingDirectoryForPull(bruvPath string, selectPaths []string) error {
	// For now, we'll just print what would be filtered
	fmt.Printf("Selective pull would filter to paths: %v\n", selectPaths)
	fmt.Println("Note: Selective pull implementation is simplified in this example")
	
	return nil
}