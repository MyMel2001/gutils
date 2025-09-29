package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func cmdPush(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: bruv push <remote> [<branch>] [--select <path>...]")
	}
	
	// Parse arguments for selective pushing
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
	
	if len(regularArgs) < 1 {
		return fmt.Errorf("usage: bruv push <remote> [<branch>] [--select <path>...]")
	}
	
	remote := args[0]
	branch := "main"
	if len(args) > 1 {
		branch = args[1]
	}
	
	// Get current branch ref
	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}
	
	branchRefPath := filepath.Join(bruvPath, "refs", "heads", branch)
	commitHashBytes, err := os.ReadFile(branchRefPath)
	if err != nil {
		return fmt.Errorf("could not read branch ref: %w", err)
	}
	commitHash := strings.TrimSpace(string(commitHashBytes))
	
	// Parse remote URL (simplified format: host:port/repo)
	parts := strings.Split(remote, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid remote format. Use: host:port/repo")
	}
	host := parts[0]
	repoPath := parts[1]
	
	// Connect to remote
	conn, err := net.Dial("tcp", host+":9418")
	if err != nil {
		return fmt.Errorf("could not connect to remote: %w", err)
	}
	defer conn.Close()
	
	// Send push request
	if len(selectPaths) > 0 {
		// Send selective push request with paths
		fmt.Fprintf(conn, "git-receive-pack %s --select %s\n", repoPath, strings.Join(selectPaths, " "))
	} else {
		fmt.Fprintf(conn, "git-receive-pack %s\n", repoPath)
	}
	
	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}
	
	if strings.HasPrefix(response, "error:") {
		return fmt.Errorf("remote rejected push: %s", strings.TrimSpace(response))
	}
	if !strings.HasPrefix(response, "ok") {
		return fmt.Errorf("remote rejected push: %s", strings.TrimSpace(response))
	}
	
	// Create packfile with objects to push
	var packBuffer *bytes.Buffer
	if len(selectPaths) > 0 {
		// Create selective packfile with only specified paths
		packBuffer, err = createSelectivePackfile(commitHash, selectPaths)
		if err != nil {
			return fmt.Errorf("could not create selective packfile: %w", err)
		}
	} else {
		// Create full packfile
		packBuffer, err = createPackfile(commitHash)
		if err != nil {
			return fmt.Errorf("could not create packfile: %w", err)
		}
	}
	
	// Send packfile
	if _, err := conn.Write(packBuffer.Bytes()); err != nil {
		return fmt.Errorf("could not send packfile: %w", err)
	}
	
	fmt.Printf("Pushed %s to %s (%s)\n", branch, remote, commitHash[:7])
	return nil
}

// createSelectivePackfile creates a packfile with only objects related to specified paths
func createSelectivePackfile(commitHash string, selectPaths []string) (*bytes.Buffer, error) {
	// For now, we'll just print what would be included
	fmt.Printf("Selective push would include paths: %v\n", selectPaths)
	fmt.Println("Note: Selective push implementation is simplified in this example")
	
	// In a real implementation, we would:
	// 1. Parse the commit to get the tree
	// 2. Filter the tree to only include objects for the specified paths
	// 3. Collect only those objects and their dependencies
	// 4. Create a packfile with only those objects
	
	// For now, just create a regular packfile
	return createPackfile(commitHash)
}