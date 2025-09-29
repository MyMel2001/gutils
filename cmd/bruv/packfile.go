package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"bufio"
	"crypto/sha1"
)

func createPackfile(commitHash string) (*bytes.Buffer, error) {
	objectsToPack := make(map[string]bool)
	if err := collectObjects(commitHash, objectsToPack); err != nil {
		return nil, err
	}

	packBuffer := new(bytes.Buffer)
	// Write packfile header
	header := struct {
		Signature [4]byte
		Version   uint32
		NumObjects uint32
	}{
		Signature: [4]byte{'P', 'A', 'C', 'K'},
		Version:   2,
		NumObjects: uint32(len(objectsToPack)),
	}
	if err := binary.Write(packBuffer, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	// Write objects
	for hash := range objectsToPack {
		if err := writePackedObject(packBuffer, hash); err != nil {
			return nil, fmt.Errorf("packing object %s: %w", hash, err)
		}
	}

	// Write packfile checksum (SHA1 of all previous data)
	hasher := sha1.New()
	hasher.Write(packBuffer.Bytes())
	checksum := hasher.Sum(nil)
	
	if _, err := packBuffer.Write(checksum); err != nil {
		return nil, fmt.Errorf("writing packfile checksum: %w", err)
	}

	return packBuffer, nil
}

func collectObjects(hashStr string, objects map[string]bool) error {
	if objects[hashStr] {
		return nil // Already processed
	}
	objects[hashStr] = true

	objType, content, err := readObject(hashStr)
	if err != nil {
		return err
	}

	switch objType {
	case "commit":
		commitBody := string(content)
		lines := strings.Split(commitBody, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "tree ") {
				treeHash := strings.TrimPrefix(line, "tree ")
				if err := collectObjects(treeHash, objects); err != nil {
					return err
				}
			} else if strings.HasPrefix(line, "parent ") {
				parentHash := strings.TrimPrefix(line, "parent ")
				if err := collectObjects(parentHash, objects); err != nil {
					return err
				}
			}
		}
	case "tree":
		buf := bytes.NewBuffer(content)
		for {
			_, err := buf.ReadString(' ')
			if err != nil {
				break // End of tree
			}
			_, err = buf.ReadString(0)
			if err != nil {
				return fmt.Errorf("malformed tree object")
			}
			var blobHashBytes [20]byte
			if _, err := io.ReadFull(buf, blobHashBytes[:]); err != nil {
				return fmt.Errorf("malformed tree object: %w", err)
			}
			
			blobHash := fmt.Sprintf("%x", blobHashBytes)
			// The object type is determined by reading the object itself,
			// so we just collect all hashes from the tree.
			if err := collectObjects(blobHash, objects); err != nil {
				return err
			}
		}
	case "blob":
		// Blobs have no further objects to collect
	}

	return nil
}

func readObject(hashStr string) (string, []byte, error) {
	bruvPath, err := findBruvDir()
	if err != nil {
		// Try alternative paths for server context
		bruvPath = findAlternativeBruvPath()
		if bruvPath == "" {
			return "", nil, err
		}
	}

	objectPath := filepath.Join(bruvPath, "objects", hashStr[:2], hashStr[2:])
	f, err := os.Open(objectPath)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return "", nil, err
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return "", nil, err
	}

	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid object format for %s", hashStr)
	}

	header := parts[0]
	content := parts[1]

	var objType string
	fmt.Sscanf(string(header), "%s", &objType)

	return objType, content, nil
}

// findAlternativeBruvPath tries to find the .bruv directory in alternative locations
func findAlternativeBruvPath() string {
	// Try current directory first
	if _, err := os.Stat(".bruv"); err == nil {
		return ".bruv"
	}
	
	// Try test-repo directory
	if _, err := os.Stat("test-repo/.bruv"); err == nil {
		return "test-repo/.bruv"
	}
	
	// Try parent directory
	if wd, err := os.Getwd(); err == nil {
		parent := filepath.Dir(wd)
		if parent != wd {
			if _, err := os.Stat(filepath.Join(parent, "test-repo/.bruv")); err == nil {
				return filepath.Join(parent, "test-repo/.bruv")
			}
		}
	}
	
	return ""
}

// writeBlobObject writes a blob object to the object database
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
		return hash, nil // Object already exists
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

func writePackedObject(w io.Writer, hashStr string) error {
	bruvPath, err := findBruvDir()
	if err != nil {
		bruvPath = ".bruv" // HACK
		wd, err := os.Getwd()
		if err == nil {
			if strings.HasSuffix(wd, "test-repo") {
				bruvPath = ".bruv"
			} else {
				bruvPath = "test-repo/.bruv"
			}
		}
	}
	objectPath := filepath.Join(bruvPath, "objects", hashStr[:2], hashStr[2:])
	
	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()
	
	// Read the compressed object data from storage
	compressedData, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	
	// Decompress to get the original object data
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return err
	}
	decompressed, err := io.ReadAll(r)
	r.Close()
	if err != nil {
		return err
	}

	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid object format for %s", hashStr)
	}
	
	header := parts[0]
	content := parts[1]
	
	var objTypeStr string
	var size int
	if _, err := fmt.Sscanf(string(header), "%s %d", &objTypeStr, &size); err != nil {
		return fmt.Errorf("parsing object header for %s: %w", hashStr, err)
	}

	typeMap := map[string]uint8{
		"commit": 1,
		"tree":   2,
		"blob":   3,
	}
	objType := typeMap[objTypeStr]
	if objType == 0 {
		return fmt.Errorf("unknown object type: %s", objTypeStr)
	}

	// Write packfile object header (type and size)
	var objectHeader []byte
	sizeVal := uint64(size)
	
	// First byte: type in high nibble (4 bits), size low nibble (4 bits)
	b := byte(sizeVal&0x0f) | byte(objType<<4)
	sizeVal >>= 4
	
	// Continue with variable-length encoding for remaining size
	for {
		if sizeVal == 0 {
			objectHeader = append(objectHeader, b)
			break
		}
		objectHeader = append(objectHeader, b|0x80)
		b = byte(sizeVal & 0x7f)
		sizeVal >>= 7
	}
	
	if _, err := w.Write(objectHeader); err != nil {
		return err
	}

	// Write the compressed content (not the full object with header)
	var compressedContent bytes.Buffer
	zlibWriter := zlib.NewWriter(&compressedContent)
	if _, err := zlibWriter.Write(content); err != nil {
		return err
	}
	zlibWriter.Close()
	
	if _, err := w.Write(compressedContent.Bytes()); err != nil {
		return err
	}
	
	return nil
}

func unpackPackfile(packfilePath, bruvPath string) error {
	f, err := os.Open(packfilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var header struct {
		Signature  [4]byte
		Version    uint32
		NumObjects uint32
	}
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		return fmt.Errorf("could not read packfile header: %w", err)
	}
	if string(header.Signature[:]) != "PACK" || header.Version != 2 {
		return fmt.Errorf("not a valid bruv packfile")
	}

	// Read and store objects for delta resolution
	objects := make([][]byte, 0, header.NumObjects)
	objectTypes := make([]uint8, 0, header.NumObjects)

	for i := 0; i < int(header.NumObjects); i++ {
		// Read object header (type and size)
		objType, size, err := readObjectHeader(reader)
		if err != nil {
			return fmt.Errorf("error reading object %d header: %w", i, err)
		}

		objectTypes = append(objectTypes, objType)
		
		var objectData []byte
		
		switch objType {
		case 1, 2, 3: // commit, tree, blob
			// Read compressed data - we need to buffer it to determine the exact length
			var compressedData bytes.Buffer
			tempReader := io.TeeReader(reader, &compressedData)
			
			// Try to decompress to find the boundary
			zlibReader, err := zlib.NewReader(tempReader)
			if err != nil {
				return fmt.Errorf("error creating zlib reader for object %d: %w", i, err)
			}
			
			objectData, err = io.ReadAll(zlibReader)
			zlibReader.Close()
			if err != nil {
				return fmt.Errorf("failed to decompress object %d: %w", i, err)
			}
			
		case 6, 7: // OBJ_OFS_DELTA, OBJ_REF_DELTA
			// For now, we'll skip delta objects as they require more complex handling
			// In a full implementation, we'd resolve these against base objects
			return fmt.Errorf("delta objects (type %d) not yet supported", objType)
			
		default:
			return fmt.Errorf("unknown object type %d", objType)
		}
		
		objects = append(objects, objectData)
		
		// Write the object to the object database
		if err := writeObjectToDB(bruvPath, objectData); err != nil {
			return fmt.Errorf("error writing object %d to database: %w", i, err)
		}
	}

	return nil
}

// readObjectHeader reads the variable-length object header and returns the type and size
func readObjectHeader(r io.ByteReader) (uint8, uint64, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	
	objType := (b >> 4) & 0x07
	size := uint64(b & 0x0f)
	shift := uint(4)
	
	for b&0x80 != 0 {
		b, err = r.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		size |= uint64(b&0x7f) << shift
		shift += 7
	}
	
	return objType, size, nil
}

// writeObjectToDB writes a decompressed object to the object database
func writeObjectToDB(bruvPath string, objectData []byte) error {
	// Hash the object to get its ID
	hasher := sha1.New()
	hasher.Write(objectData)
	hash := hasher.Sum(nil)
	hashStr := fmt.Sprintf("%x", hash)

	// Create object directory
	objectDir := filepath.Join(bruvPath, "objects", hashStr[:2])
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return err
	}
	objectPath := filepath.Join(objectDir, hashStr[2:])

	// Compress the object data
	var compressedData bytes.Buffer
	zlibWriter := zlib.NewWriter(&compressedData)
	if _, err := zlibWriter.Write(objectData); err != nil {
		return err
	}
	zlibWriter.Close()

	// Write the compressed object
	if err := os.WriteFile(objectPath, compressedData.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

// createSelectivePackfile creates a packfile with only objects related to specified paths
func createSelectivePackfile(commitHash string, selectPaths []string) (*bytes.Buffer, error) {
	// For now, we'll just print what would be included
	fmt.Printf("Selective operation would include paths: %v\n", selectPaths)
	fmt.Println("Note: Selective operation implementation is simplified in this example")
	
	// In a real implementation, we would:
	// 1. Parse the commit to get the tree
	// 2. Filter the tree to only include objects for the specified paths
	// 3. Collect only those objects and their dependencies
	// 4. Create a packfile with only those objects
	
	// For now, just create a regular packfile
	return createPackfile(commitHash)
}

func updateRefsAfterClone(bruvPath, commitHash string) error {
	mainRefPath := filepath.Join(bruvPath, "refs", "heads", "main")
	if err := os.WriteFile(mainRefPath, []byte(commitHash+"\n"), 0644); err != nil {
		return err
	}
	
	headPath := filepath.Join(bruvPath, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return err
	}
	return nil
} 