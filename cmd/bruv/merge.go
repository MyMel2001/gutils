package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func cmdMerge(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: bruv merge <source-branch> <target-branch> [remote] [--select <path>...]")
	}
	
	// Parse arguments for selective merging
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
	
	if len(regularArgs) < 2 {
		return fmt.Errorf("usage: bruv merge <source-branch> <target-branch> [remote] [--select <path>...]")
	}
	
	sourceBranch := args[0]
	targetBranch := args[1]
	
	// Default remote if not specified
	remote := "origin"
	if len(args) > 2 {
		remote = args[2]
	}
	
	// Validate branch names
	if sourceBranch == targetBranch {
		return fmt.Errorf("source and target branches cannot be the same")
	}
	
	if targetBranch != "main" && targetBranch != "master" {
		return fmt.Errorf("merge requests are only allowed to main/master branch for security")
	}
	
	// Get current repository info
	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}
	
	// Check if source branch exists
	sourceRefPath := filepath.Join(bruvPath, "refs", "heads", sourceBranch)
	sourceHashBytes, err := os.ReadFile(sourceRefPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source branch '%s' does not exist", sourceBranch)
		}
		return fmt.Errorf("could not read source branch: %w", err)
	}
	
	// Check if target branch exists (typically main/master)
	targetRefPath := filepath.Join(bruvPath, "refs", "heads", targetBranch)
	targetHashBytes, err := os.ReadFile(targetRefPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("target branch '%s' does not exist", targetBranch)
		}
		return fmt.Errorf("could not read target branch: %w", err)
	}
	
	// Get the remote URL from config or use default
	remoteURL := getRemoteURL(remote)
	if remoteURL == "" {
		return fmt.Errorf("remote '%s' not found", remote)
	}
	
	// Parse remote URL
	parts := strings.Split(remoteURL, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid remote format '%s'. Use: host:port/repo", remoteURL)
	}
	host := parts[0]
	repoPath := parts[1]
	
	// Connect to remote server
	conn, err := net.Dial("tcp", host+":9418")
	if err != nil {
		return fmt.Errorf("could not connect to remote '%s': %w", remote, err)
	}
	defer conn.Close()
	
	// Send merge request
	fmt.Fprintf(conn, "git-merge-request %s\n", repoPath)
	
	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}
	
	if strings.HasPrefix(response, "error:") {
		return fmt.Errorf("remote rejected merge request: %s", strings.TrimSpace(response))
	}
	
	if !strings.HasPrefix(response, "ok") {
		return fmt.Errorf("remote rejected merge request: %s", strings.TrimSpace(response))
	}
	
	// Send merge details
	if len(selectPaths) > 0 {
		fmt.Fprintf(conn, "%s %s --select %s\n", sourceBranch, targetBranch, strings.Join(selectPaths, " "))
	} else {
		fmt.Fprintf(conn, "%s %s\n", sourceBranch, targetBranch)
	}
	
	// Read final response
	finalResponse, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("could not read final response: %w", err)
	}
	
	// Display merge request information
	fmt.Printf("Merge request submitted successfully!\n")
	fmt.Printf("Source branch: %s (%s)\n", sourceBranch, strings.TrimSpace(string(sourceHashBytes))[:7])
	fmt.Printf("Target branch: %s (%s)\n", targetBranch, strings.TrimSpace(string(targetHashBytes))[:7])
	fmt.Printf("Response: %s", finalResponse)
	fmt.Println("\nWaiting for repository owner approval...")
	fmt.Println("Use 'bruv approve <merge-request-id>' to approve this merge request (owner only)")
	
	return nil
}

// getRemoteURL returns the URL for a given remote name
func getRemoteURL(remoteName string) string {
	// In a real implementation, this would read from .bruv/config
	// For now, we'll use hardcoded mappings
	switch remoteName {
	case "origin":
		return "localhost:9418/test-repo"
	case "upstream":
		return "localhost:9418/upstream-repo"
	default:
		return ""
	}
}