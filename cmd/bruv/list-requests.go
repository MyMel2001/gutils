package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cmdListRequests(args []string) error {
	// Get repository info
	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}
	
	// Check if merge requests directory exists
	mergeRequestPath := filepath.Join(bruvPath, "merge-requests")
	if _, err := os.Stat(mergeRequestPath); os.IsNotExist(err) {
		fmt.Println("No merge requests found")
		return nil
	}
	
	// Read all merge request files
	files, err := os.ReadDir(mergeRequestPath)
	if err != nil {
		return fmt.Errorf("could not read merge requests directory: %w", err)
	}
	
	if len(files) == 0 {
		fmt.Println("No merge requests found")
		return nil
	}
	
	fmt.Println("Merge Requests:")
	fmt.Println("===============")
	
	pendingCount := 0
	approvedCount := 0
	
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "mr-") || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		
		requestFile := filepath.Join(mergeRequestPath, file.Name())
		requestContent, err := os.ReadFile(requestFile)
		if err != nil {
			continue
		}
		
		// Parse the merge request
		lines := strings.Split(string(requestContent), "\n")
		var sourceBranch, targetBranch, status, timestamp, requester string
		
		for _, line := range lines {
			if strings.HasPrefix(line, "Source: ") {
				sourceBranch = strings.TrimPrefix(line, "Source: ")
			} else if strings.HasPrefix(line, "Target: ") {
				targetBranch = strings.TrimPrefix(line, "Target: ")
			} else if strings.HasPrefix(line, "Status: ") {
				status = strings.TrimPrefix(line, "Status: ")
			} else if strings.HasPrefix(line, "Timestamp: ") {
				timestamp = strings.TrimPrefix(line, "Timestamp: ")
			} else if strings.HasPrefix(line, "Requester: ") {
				requester = strings.TrimPrefix(line, "Requester: ")
			}
		}
		
		if sourceBranch == "" || targetBranch == "" || status == "" {
			continue
		}
		
		// Extract merge request ID from filename
		requestID := strings.TrimSuffix(strings.TrimPrefix(file.Name(), "mr-"), ".txt")
		
		fmt.Printf("MR-%s: %s -> %s [%s]\n", requestID, sourceBranch, targetBranch, status)
		
		if status == "pending" {
			pendingCount++
			fmt.Printf("  Use 'bruv approve %s' to approve this merge request\n", requestID)
		} else if status == "approved" {
			approvedCount++
		}
		
		if requester != "" {
			fmt.Printf("  Requested by: %s\n", requester)
		}
		
		fmt.Println()
	}
	
	fmt.Printf("Summary: %d pending, %d approved\n", pendingCount, approvedCount)
	
	return nil
}