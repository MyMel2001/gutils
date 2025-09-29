package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func cmdServe() error {
	ln, err := net.Listen("tcp", ":9418") // Git daemon port
	if err != nil {
		return err
	}
	defer ln.Close()
	fmt.Println("Bruv server listening on :9418")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %s\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			// Parse command and arguments
			command := parts[0]
			var args []string
			if len(parts) > 1 {
				args = parts[1:]
			}
			
			switch command {
			case "git-upload-pack":
				if len(args) > 0 {
					handleUploadPack(conn, args[0], args[1:]...)
				} else {
					handleUploadPack(conn, "", []string{}...)
				}
			case "git-receive-pack":
				if len(args) > 0 {
					handleReceivePack(conn, args[0], args[1:]...)
				} else {
					handleReceivePack(conn, "", []string{}...)
				}
			case "git-merge-request":
				if len(args) > 0 {
					handleMergeRequest(conn, args[0], args[1:]...)
				} else {
					handleMergeRequest(conn, "", []string{}...)
				}
			default:
				fmt.Fprintf(os.Stderr, "unknown git command: %s\n", command)
			}
		}
	}
}

func handleUploadPack(conn net.Conn, repoPath string, args ...string) {
	// In a real implementation, we'd respect the repo path from the client
	// For now, we assume the server is run from the parent of the repo dir
	mainRefPath := filepath.Join("test-repo", ".bruv", "refs", "heads", "main")
	mainRef, err := os.ReadFile(mainRefPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read main ref: %v\n", err)
		return
	}
	commitHash := strings.TrimSpace(string(mainRef))

	// Parse selective arguments
	selectPaths := []string{}
	for i, arg := range args {
		if arg == "--select" && i+1 < len(args) {
			// Collect all paths after --select
			for j := i + 1; j < len(args); j++ {
				if !strings.HasPrefix(args[j], "-") {
					selectPaths = append(selectPaths, args[j])
				} else {
					break
				}
			}
			break
		}
	}

	// Send the ref first, so the client knows what it's getting
	fmt.Fprintf(conn, "%s refs/heads/main\n", commitHash)
	// A real protocol would have a flush packet (0000) here
	
	// Now create and send the packfile
	var packBuffer *bytes.Buffer
	if len(selectPaths) > 0 {
		// Create selective packfile with only specified paths
		fmt.Printf("Selective pull request for paths: %v\n", selectPaths)
		packBuffer, err = createSelectivePackfile(commitHash, selectPaths)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create selective packfile: %v\n", err)
			return
		}
	} else {
		// Create full packfile
		packBuffer, err = createPackfile(commitHash)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create packfile: %v\n", err)
			return
		}
	}
	
	if _, err := conn.Write(packBuffer.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "could not write packfile: %v\n", err)
	}
}

func handleReceivePack(conn net.Conn, repoPath string, args ...string) {
	// Block direct pushes to master/main branch
	// In a real implementation, we'd parse the packfile to determine the target branch
	// For now, we'll reject all receive-pack operations to protect master/main
	fmt.Fprintf(conn, "error: Direct pushes to master/main branch are not allowed. Use 'bruv merge' command instead.\n")
	return
	
	// The original implementation below is disabled to prevent direct pushes
	
	/*
	// Send acknowledgment
	fmt.Fprintf(conn, "ok\n")
	
	// Read packfile data and save it
	bruvPath := filepath.Join("test-repo", ".bruv")
	packDir := filepath.Join(bruvPath, "objects", "pack")
	if err := os.MkdirAll(packDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "could not create pack directory: %v\n", err)
		return
	}
	
	packfilePath := filepath.Join(packDir, "temp.receive.pack")
	packfile, err := os.Create(packfilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create packfile: %v\n", err)
		return
	}
	defer packfile.Close()
	
	// Copy packfile data from connection
	_, err = io.Copy(packfile, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not receive packfile: %v\n", err)
		return
	}
	
	// Unpack the received packfile
	if err := unpackPackfile(packfilePath, bruvPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unpack packfile: %v\n", err)
		return
	}
	
	// Update refs (simplified - just update main branch)
	// In a real implementation, we'd parse the packfile to get the new commit hash
	// For now, we'll just clean up the temp file
	os.Remove(packfilePath)
	
	fmt.Println("Received and unpacked packfile successfully")
	*/
}

// handleMergeRequest handles merge request submissions from clients
func handleMergeRequest(conn net.Conn, repoPath string, args ...string) {
	// Send acknowledgment
	fmt.Fprintf(conn, "ok\n")
	
	// Read merge request data
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		mergeInfo := scanner.Text()
		parts := strings.Fields(mergeInfo)
		if len(parts) >= 2 {
			// Parse source and target branches
			sourceBranch := parts[0]
			targetBranch := parts[1]
			
			// Parse selective arguments from the merge request
			selectPaths := []string{}
			for i := 2; i < len(parts); i++ {
				if parts[i] == "--select" && i+1 < len(parts) {
					// Collect all paths after --select
					for j := i + 1; j < len(parts); j++ {
						if !strings.HasPrefix(parts[j], "-") {
							selectPaths = append(selectPaths, parts[j])
						} else {
							break
						}
					}
					break
				}
			}
			
			if len(selectPaths) > 0 {
				fmt.Printf("Selective merge request for paths: %v\n", selectPaths)
			}
			
			// Validate branches
			if targetBranch != "main" && targetBranch != "master" {
				fmt.Fprintf(conn, "error: merge requests only allowed to main/master branch\n")
				return
			}
			
			// Check if branches exist
			bruvPath := filepath.Join("test-repo", ".bruv")
			sourceRefPath := filepath.Join(bruvPath, "refs", "heads", sourceBranch)
			targetRefPath := filepath.Join(bruvPath, "refs", "heads", targetBranch)
			
			if _, err := os.Stat(sourceRefPath); os.IsNotExist(err) {
				fmt.Fprintf(conn, "error: source branch '%s' does not exist\n", sourceBranch)
				return
			}
			
			if _, err := os.Stat(targetRefPath); os.IsNotExist(err) {
				fmt.Fprintf(conn, "error: target branch '%s' does not exist\n", targetBranch)
				return
			}
			
			// Store merge request
			mergeRequestPath := filepath.Join(bruvPath, "merge-requests")
			if err := os.MkdirAll(mergeRequestPath, 0755); err != nil {
				fmt.Fprintf(conn, "error: could not create merge request directory\n")
				return
			}
			
			// Generate unique merge request ID
			timestamp := fmt.Sprintf("%d", time.Now().Unix())
			requestFile := filepath.Join(mergeRequestPath, fmt.Sprintf("mr-%s.txt", timestamp))
			
			// Get commit hashes for both branches
			sourceHash, _ := os.ReadFile(sourceRefPath)
			targetHash, _ := os.ReadFile(targetRefPath)
			
			requestContent := fmt.Sprintf(`Source: %s
Target: %s
SourceHash: %s
TargetHash: %s
Status: pending
Timestamp: %s
Requester: user
`, strings.TrimSpace(sourceBranch), strings.TrimSpace(targetBranch),
   strings.TrimSpace(string(sourceHash)), strings.TrimSpace(string(targetHash)), timestamp)
			
			if err := os.WriteFile(requestFile, []byte(requestContent), 0644); err != nil {
				fmt.Fprintf(conn, "error: could not save merge request\n")
				return
			}
			
			fmt.Fprintf(conn, "Merge request created: MR-%s\n", timestamp)
			fmt.Printf("Received merge request MR-%s: %s -> %s\n", timestamp, sourceBranch, targetBranch)
		} else {
			fmt.Fprintf(conn, "error: invalid merge request format\n")
		}
	}
}