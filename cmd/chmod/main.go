package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// chmod: changes file permissions, supports numeric and symbolic modes
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "chmod: usage: chmod MODE FILE...")
		os.Exit(1)
	}
	modeStr := os.Args[1]
	for _, fname := range os.Args[2:] {
		info, err := os.Stat(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "chmod: cannot stat:", fname, err)
			continue
		}
		origMode := info.Mode().Perm()
		var newMode os.FileMode
		if isNumericMode(modeStr) {
			m, err := strconv.ParseUint(modeStr, 8, 32)
			if err != nil {
				fmt.Fprintln(os.Stderr, "chmod: invalid mode:", modeStr)
				continue
			}
			newMode = os.FileMode(m)
		} else {
			newMode, err = applySymbolicMode(origMode, modeStr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "chmod:", err)
				continue
			}
		}
		if err := os.Chmod(fname, newMode); err != nil {
			fmt.Fprintln(os.Stderr, "chmod: cannot change permissions:", fname, err)
		}
	}
}

// isNumericMode returns true if the mode string is numeric
func isNumericMode(s string) bool {
	for _, c := range s {
		if c < '0' || c > '7' {
			return false
		}
	}
	return true
}

// applySymbolicMode applies a symbolic mode string to the original mode
func applySymbolicMode(orig os.FileMode, modeStr string) (os.FileMode, error) {
	mode := uint32(orig)
	for _, clause := range strings.Split(modeStr, ",") {
		clause = strings.TrimSpace(clause)
		if clause == "" {
			continue
		}
		who, op, perm, err := parseSymbolicClause(clause)
		if err != nil {
			return 0, err
		}
		for _, w := range who {
			var mask uint32
			if perm == 0 {
				continue
			}
			switch w {
			case 'u':
				mask = 0700
			case 'g':
				mask = 0070
			case 'o':
				mask = 0007
			case 'a':
				mask = 0777
			}
			shift := map[byte]uint{ 'u': 6, 'g': 3, 'o': 0, 'a': 0 }[w]
			p := (perm << shift) & mask
			switch op {
			case '+':
				mode |= p
			case '-':
				mode &^= p
			case '=':
				mode &^= mask
				mode |= p
			}
		}
	}
	return os.FileMode(mode), nil
}

// parseSymbolicClause parses a symbolic mode clause (e.g., u+x)
func parseSymbolicClause(clause string) ([]byte, byte, uint32, error) {
	who := []byte{}
	perm := uint32(0)
	op := byte(0)
	// Find op position
	opIdx := strings.IndexAny(clause, "+-=")
	if opIdx == -1 {
		return nil, 0, 0, fmt.Errorf("invalid symbolic mode: %s", clause)
	}
	whoStr := clause[:opIdx]
	if whoStr == "" {
		who = []byte{'a'}
	} else {
		who = []byte(whoStr)
	}
	op = clause[opIdx]
	permStr := clause[opIdx+1:]
	for _, c := range permStr {
		switch c {
		case 'r': perm |= 4
		case 'w': perm |= 2
		case 'x': perm |= 1
		}
	}
	return who, op, perm, nil
} 