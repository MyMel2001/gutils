package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func cmdApprove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: bruv approve <merge-request-id>")
	}
	
	mergeRequestID := args[0]
	
	// Get repository info
	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}
	
	// Check if merge requests directory exists
	mergeRequestPath := filepath.Join(bruvPath, "merge-requests")
	if _, err := os.Stat(mergeRequestPath); os.IsNotExist(err) {
		return fmt.Errorf("no merge requests found")
	}
	
	// Find the merge request file
	requestFile := filepath.Join(mergeRequestPath, fmt.Sprintf("mr-%s.txt", mergeRequestID))
	if _, err := os.Stat(requestFile); os.IsNotExist(err) {
		return fmt.Errorf("merge request %s not found", mergeRequestID)
	}
	
	// Read the merge request
	requestContent, err := os.ReadFile(requestFile)
	if err != nil {
		return fmt.Errorf("could not read merge request: %w", err)
	}
	
	// Parse the merge request
	lines := strings.Split(string(requestContent), "\n")
	var sourceBranch, targetBranch, sourceHash, targetHash, status string
	
	for _, line := range lines {
		if strings.HasPrefix(line, "Source: ") {
			sourceBranch = strings.TrimPrefix(line, "Source: ")
		} else if strings.HasPrefix(line, "Target: ") {
			targetBranch = strings.TrimPrefix(line, "Target: ")
		} else if strings.HasPrefix(line, "SourceHash: ") {
			sourceHash = strings.TrimPrefix(line, "SourceHash: ")
		} else if strings.HasPrefix(line, "TargetHash: ") {
			targetHash = strings.TrimPrefix(line, "TargetHash: ")
		} else if strings.HasPrefix(line, "Status: ") {
			status = strings.TrimPrefix(line, "Status: ")
		}
	}
	
	if sourceBranch == "" || targetBranch == "" || sourceHash == "" || targetHash == "" {
		return fmt.Errorf("invalid merge request format")
	}
	
	if status != "pending" {
		return fmt.Errorf("merge request %s is already %s", mergeRequestID, status)
	}
	
	// Validate that branches still exist and haven't changed
	sourceRefPath := filepath.Join(bruvPath, "refs", "heads", sourceBranch)
	targetRefPath := filepath.Join(bruvPath, "refs", "heads", targetBranch)
	
	currentSourceHashBytes, err := os.ReadFile(sourceRefPath)
	if err != nil {
		return fmt.Errorf("source branch '%s' no longer exists", sourceBranch)
	}
	currentTargetHashBytes, err := os.ReadFile(targetRefPath)
	if err != nil {
		return fmt.Errorf("target branch '%s' no longer exists", targetBranch)
	}
	
	currentSourceHash := strings.TrimSpace(string(currentSourceHashBytes))
	currentTargetHash := strings.TrimSpace(string(currentTargetHashBytes))
	
	// Check if branches have diverged
	if currentSourceHash != sourceHash {
		return fmt.Errorf("source branch '%s' has new commits since merge request was created", sourceBranch)
	}
	
	if currentTargetHash != targetHash {
		return fmt.Errorf("target branch '%s' has new commits since merge request was created", targetBranch)
	}
	
	// Perform the actual merge by updating the target branch ref
	// This is a fast-forward merge in this simplified implementation
	if err := os.WriteFile(targetRefPath, []byte(sourceHash+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to update target branch: %w", err)
	}
	
	// Update the merge request status
	newContent := strings.Replace(string(requestContent), "Status: pending", "Status: approved", 1)
	newContent += fmt.Sprintf("ApprovedBy: owner\nApprovedAt: %d\n", time.Now().Unix())
	
	if err := os.WriteFile(requestFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("could not update merge request status: %w", err)
	}
	
	fmt.Printf("Merge request %s approved and completed successfully\n", mergeRequestID)
	fmt.Printf("Source: %s (%s) -> Target: %s (%s)\n", sourceBranch, sourceHash[:7], targetBranch, sourceHash[:7])
	fmt.Println("Fast-forward merge completed")
	
	return nil
}

// Helper function to read commit object (simplified version of existing functionality)
func readCommitObject(commitHash string) ([]byte, error) {
	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}
	
	objectPath := filepath.Join(bruvPath, "objects", commitHash[:2], commitHash[2:])
	f, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	r, err := zlib.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	
	decompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	
	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid object format")
	}
	
	return parts[1], nil
}