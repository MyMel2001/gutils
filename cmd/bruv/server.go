package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
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
		if len(parts) > 1 && parts[0] == "git-upload-pack" {
			// In a real implementation, we'd respect the repo path from the client
			// For now, we assume the server is run from the parent of the repo dir
			mainRefPath := filepath.Join("test-repo", ".bruv", "refs", "heads", "main")
			mainRef, err := os.ReadFile(mainRefPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not read main ref: %v\n", err)
				return
			}
			commitHash := strings.TrimSpace(string(mainRef))

			// Send the ref first, so the client knows what it's getting
			fmt.Fprintf(conn, "%s refs/heads/main\n", commitHash)
			// A real protocol would have a flush packet (0000) here
			
			// Now create and send the packfile
			packBuffer, err := createPackfile(commitHash)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not create packfile: %v\n", err)
				return
			}
			
			if _, err := conn.Write(packBuffer.Bytes()); err != nil {
				fmt.Fprintf(os.Stderr, "could not write packfile: %v\n", err)
			}
		}
	}
} 