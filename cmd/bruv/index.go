package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Simplified from https://git-scm.com/docs/index-format
type IndexHeader struct {
	Signature [4]byte // "DIRC"
	Version   uint32
	Entries   uint32
}

type IndexEntry struct {
	CTimeSec  uint32
	CTimeNano uint32
	MTimeSec  uint32
	MTimeNano uint32
	Dev       uint32
	Ino       uint32
	Mode      uint32
	UID       uint32
	GID       uint32
	Size      uint32
	Hash      [20]byte
	Flags     uint16
	Path      string
}

func readIndex() (map[string]*IndexEntry, error) {
	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}
	indexPath := filepath.Join(bruvPath, "index")
	entries := make(map[string]*IndexEntry)

	f, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return nil, err
	}
	defer f.Close()

	var header IndexHeader
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		return nil, fmt.Errorf("reading index header: %w", err)
	}

	if string(header.Signature[:]) != "DIRC" {
		return nil, fmt.Errorf("not a bruv index file")
	}

	for i := 0; i < int(header.Entries); i++ {
		var entry IndexEntry
		var core struct {
			CTimeSec, CTimeNano, MTimeSec, MTimeNano, Dev, Ino, Mode, UID, GID, Size uint32
			Hash [20]byte
			Flags uint16
		}
		if err := binary.Read(f, binary.BigEndian, &core); err != nil {
			return nil, fmt.Errorf("reading index entry core: %w", err)
		}
		entry.CTimeSec, entry.CTimeNano, entry.MTimeSec, entry.MTimeNano = core.CTimeSec, core.CTimeNano, core.MTimeSec, core.MTimeNano
		entry.Dev, entry.Ino, entry.Mode, entry.UID, entry.GID, entry.Size = core.Dev, core.Ino, core.Mode, core.UID, core.GID, core.Size
		entry.Hash, entry.Flags = core.Hash, core.Flags

		pathLen := entry.Flags & 0xfff
		pathBytes := make([]byte, pathLen)
		if _, err := io.ReadFull(f, pathBytes); err != nil {
			return nil, err
		}
		entry.Path = string(pathBytes)
		entries[entry.Path] = &entry

		padding := (8 - (62+int(pathLen))%8) % 8
		if padding > 0 {
			if _, err := io.CopyN(io.Discard, f, int64(padding)); err != nil {
				return nil, err
			}
		}
	}
	return entries, nil
}

func writeIndex(entries map[string]*IndexEntry) error {
	bruvPath, err := findBruvDir()
	if err != nil {
		return err
	}
	indexPath := filepath.Join(bruvPath, "index")

	f, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer f.Close()

	header := IndexHeader{Version: 2, Entries: uint32(len(entries))}
	copy(header.Signature[:], "DIRC")
	if err := binary.Write(f, binary.BigEndian, &header); err != nil {
		return err
	}

	sortedPaths := make([]string, 0, len(entries))
	for path := range entries {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	for _, path := range sortedPaths {
		entry := entries[path]
		core := struct {
			CTimeSec, CTimeNano, MTimeSec, MTimeNano, Dev, Ino, Mode, UID, GID, Size uint32
			Hash [20]byte
			Flags uint16
		}{
			CTimeSec: entry.CTimeSec, CTimeNano: entry.CTimeNano, MTimeSec: entry.MTimeSec, MTimeNano: entry.MTimeNano,
			Dev: entry.Dev, Ino: entry.Ino, Mode: entry.Mode, UID: entry.UID, GID: entry.GID, Size: entry.Size,
			Hash: entry.Hash, Flags: entry.Flags,
		}
		if err := binary.Write(f, binary.BigEndian, &core); err != nil {
			return err
		}
		if _, err := f.WriteString(entry.Path); err != nil {
			return err
		}
		pathLen := len(entry.Path)
		padding := (8 - (62+pathLen)%8) % 8
		if _, err := f.Write(make([]byte, padding)); err != nil {
			return err
		}
	}
	return nil
}

func cmdAdd(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("nothing specified, nothing added")
	}

	entries, err := readIndex()
	if err != nil {
		return err
	}

	// Read .bruvignore configuration
	ignoreConfig, err := readBruvIgnore()
	if err != nil {
		return fmt.Errorf("reading .bruvignore: %w", err)
	}

	for _, path := range args {
		if strings.HasPrefix(path, ".bruv/") {
			continue
		}
		
		// Check if file should be ignored
		if ignoreConfig.isIgnored(path) {
			continue
		}
		
		err := addFileToIndex(entries, path)
		if err != nil {
			// Distinguish between file not existing and other errors
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "bruv: pathspec '%s' did not match any files\n", path)
				continue // continue to next argument
			}
			return fmt.Errorf("adding %s: %w", path, err)
		}
	}

	return writeIndex(entries)
}

func addFileToIndex(entries map[string]*IndexEntry, path string) error {
	lfsCfg, err := readLFSConfig()
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var objectContent []byte
	if lfsCfg.isLFS(path) {
		objectContent, err = writeLFSPointer(content)
		if err != nil {
			return err
		}
	} else {
		objectContent = content
	}

	hash, err := writeBlobObject(objectContent)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	var hashArr [20]byte
	copy(hashArr[:], hash)

	entries[path] = &IndexEntry{
		CTimeSec:  uint32(stat.ModTime().Unix()),
		MTimeSec:  uint32(stat.ModTime().Unix()),
		Mode:      uint32(stat.Mode()),
		Size:      uint32(stat.Size()),
		Hash:      hashArr,
		Flags:     uint16(len(path)),
		Path:      path,
	}

	return nil
}

func writeBlobObject(content []byte) ([]byte, error) {
	hasher := sha1.New()
	header := []byte(fmt.Sprintf("blob %d\x00", len(content)))
	hasher.Write(header)
	hasher.Write(content)
	hash := hasher.Sum(nil)

	bruvPath, err := findBruvDir()
	if err != nil {
		return nil, err
	}

	hashStr := fmt.Sprintf("%x", hash)
	objectDir := filepath.Join(bruvPath, "objects", hashStr[:2])
	objectPath := filepath.Join(objectDir, hashStr[2:])

	if _, err := os.Stat(objectPath); !os.IsNotExist(err) {
		return hash, nil
	}

	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return nil, err
	}

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(header)
	w.Write(content)
	w.Close()

	if err := os.WriteFile(objectPath, b.Bytes(), 0644); err != nil {
		return nil, err
	}

	return hash, nil
}

// Dummy timeStat for systems that don't support birth time
type timeStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (ts *timeStat) Name() string       { return ts.name }
func (ts *timeStat) Size() int64        { return ts.size }
func (ts *timeStat) Mode() os.FileMode  { return ts.mode }
func (ts *timeStat) ModTime() time.Time { return ts.modTime }
func (ts *timeStat) IsDir() bool        { return false }
func (ts *timeStat) Sys() interface{}   { return nil } 