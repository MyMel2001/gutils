package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TreeEntry struct {
	Mode string
	Hash []byte
	Path string
}

func writeTree(entries map[string]*IndexEntry) ([]byte, error) {
	var treeContent bytes.Buffer
	for path, entry := range entries {
		modeStr := strconv.FormatUint(uint64(entry.Mode), 8)
		treeContent.WriteString(fmt.Sprintf("%s %s\x00", modeStr, path))
		treeContent.Write(entry.Hash[:])
	}

	hasher := sha1.New()
	header := []byte(fmt.Sprintf("tree %d\x00", treeContent.Len()))
	hasher.Write(header)
	hasher.Write(treeContent.Bytes())
	hash := hasher.Sum(nil)

	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}
	hashStr := fmt.Sprintf("%x", hash)
	objectDir := filepath.Join(bruvPath, "objects", hashStr[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return nil, err
	}
	objectPath := filepath.Join(objectDir, hashStr[2:])

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(header)
	w.Write(treeContent.Bytes())
	w.Close()

	if err := os.WriteFile(objectPath, b.Bytes(), 0644); err != nil {
		return nil, err
	}

	return hash, nil
}

func cmdCommit(args []string) error {
	var message string
	if len(args) > 0 && (args[0] == "-m" || args[0] == "--message") {
		if len(args) < 2 {
			return fmt.Errorf("commit message not provided")
		}
		message = args[1]
	} else {
		return fmt.Errorf("commit message must be provided with -m or --message flag")
	}

	entries, err := readIndex()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return fmt.Errorf("nothing to commit, working tree clean")
	}

	treeHash, err := writeTree(entries)
	if err != nil {
		return err
	}

	parentHash, err := getParentCommitHash()
	if err != nil {
		return err
	}

	commitHash, err := writeCommit(treeHash, parentHash, message)
	if err != nil {
		return err
	}

	if err := updateCurrentBranch(commitHash); err != nil {
		return err
	}

	fmt.Printf("[%s] %s\n", fmt.Sprintf("%x", commitHash)[:7], strings.Split(message, "\n")[0])
	return nil
}

func getParentCommitHash() (string, error) {
	bruvPath, err := findBruvDir()
	if err != nil {
		return "", err
	}
	headPath := filepath.Join(bruvPath, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	refPath := strings.TrimSpace(strings.TrimPrefix(string(headContent), "ref: "))
	parentHashPath := filepath.Join(bruvPath, refPath)

	if _, err := os.Stat(parentHashPath); os.IsNotExist(err) {
		return "", nil // No parent commit yet
	}

	parentHash, err := os.ReadFile(parentHashPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(parentHash)), nil
}

func writeCommit(treeHash []byte, parentHash, message string) ([]byte, error) {
	var commitContent bytes.Buffer
	commitContent.WriteString(fmt.Sprintf("tree %x\n", treeHash))
	if parentHash != "" {
		commitContent.WriteString(fmt.Sprintf("parent %s\n", parentHash))
	}
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	author := fmt.Sprintf("%s <%s>", currentUser.Username, "user@example.com") // Placeholder email
	commitContent.WriteString(fmt.Sprintf("author %s %d +0000\n", author, time.Now().Unix()))
	commitContent.WriteString(fmt.Sprintf("committer %s %d +0000\n", author, time.Now().Unix()))
	commitContent.WriteString(fmt.Sprintf("\n%s\n", message))

	hasher := sha1.New()
	header := []byte(fmt.Sprintf("commit %d\x00", commitContent.Len()))
	hasher.Write(header)
	hasher.Write(commitContent.Bytes())
	hash := hasher.Sum(nil)

	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}
	hashStr := fmt.Sprintf("%x", hash)
	objectDir := filepath.Join(bruvPath, "objects", hashStr[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return nil, err
	}
	objectPath := filepath.Join(objectDir, hashStr[2:])

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(header)
	w.Write(commitContent.Bytes())
	w.Close()

	if err := os.WriteFile(objectPath, b.Bytes(), 0644); err != nil {
		return nil, err
	}

	return hash, nil
}

func updateCurrentBranch(commitHash []byte) error {
	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}
	headPath := filepath.Join(bruvPath, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return err
	}
	refPath := strings.TrimSpace(strings.TrimPrefix(string(headContent), "ref: "))
	branchPath := filepath.Join(bruvPath, refPath)

	return os.WriteFile(branchPath, []byte(fmt.Sprintf("%x\n", commitHash)), 0644)
} 