package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// chmod: changes file permissions, supports numeric and symbolic modes (including u+s, g+s, o+t, and u=g)
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
		origMode := info.Mode()
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

// applySymbolicMode applies a symbolic mode string to the original mode (supports u+s, g+s, o+t, and u=g)
func applySymbolicMode(orig os.FileMode, modeStr string) (os.FileMode, error) {
	mode := uint32(orig.Perm())
	// Handle special bits
	setuid := orig&os.ModeSetuid != 0
	setgid := orig&os.ModeSetgid != 0
	sticky := orig&os.ModeSticky != 0

	for _, clause := range strings.Split(modeStr, ",") {
		clause = strings.TrimSpace(clause)
		if clause == "" {
			continue
		}
		who, op, perm, special, copyFrom, err := parseSymbolicClause(clause)
		if err != nil {
			return 0, err
		}
		// Handle permission copying (e.g., u=g)
		if copyFrom != 0 {
			copyPerm := func(w byte) uint32 {
				switch w {
				case 'u':
					return (mode & 0700) >> 6
				case 'g':
					return (mode & 0070) >> 3
				case 'o':
					return (mode & 0007)
				}
				return 0
			}
			for _, w := range who {
				shift := map[byte]uint{ 'u': 6, 'g': 3, 'o': 0 }[w]
				mode &^= 7 << shift
				mode |= (copyPerm(copyFrom) & 7) << shift
			}
			continue
		}
		// Handle special bits (setuid, setgid, sticky)
		for _, w := range who {
			if special != 0 {
				switch special {
				case 's':
					if w == 'u' {
						if op == '+' || op == '=' { setuid = true } else if op == '-' { setuid = false }
					}
					if w == 'g' {
						if op == '+' || op == '=' { setgid = true } else if op == '-' { setgid = false }
					}
				case 't':
					if w == 'o' || w == 'a' {
						if op == '+' || op == '=' { sticky = true } else if op == '-' { sticky = false }
					}
				}
			}
			if perm == 0 {
				continue
			}
			var mask uint32
			switch w {
			case 'u': mask = 0700
			case 'g': mask = 0070
			case 'o': mask = 0007
			case 'a': mask = 0777
			}
			shift := map[byte]uint{ 'u': 6, 'g': 3, 'o': 0, 'a': 0 }[w]
			p := (perm << shift) & mask
			switch op {
			case '+': mode |= p
			case '-': mode &^= p
			case '=': mode &^= mask; mode |= p
			}
		}
	}
	finalMode := os.FileMode(mode)
	if setuid { finalMode |= os.ModeSetuid }
	if setgid { finalMode |= os.ModeSetgid }
	if sticky { finalMode |= os.ModeSticky }
	return finalMode, nil
}

// parseSymbolicClause parses a symbolic mode clause (e.g., u+x, g-s, o=t, u=g)
func parseSymbolicClause(clause string) ([]byte, byte, uint32, byte, byte, error) {
	who := []byte{}
	perm := uint32(0)
	special := byte(0)
	copyFrom := byte(0)
	op := byte(0)
	// Find op position
	opIdx := strings.IndexAny(clause, "+-=")
	if opIdx == -1 {
		return nil, 0, 0, 0, 0, fmt.Errorf("invalid symbolic mode: %s", clause)
	}
	whoStr := clause[:opIdx]
	if whoStr == "" {
		who = []byte{'a'}
	} else {
		who = []byte(whoStr)
	}
	op = clause[opIdx]
	permStr := clause[opIdx+1:]
	// Permission copying (e.g., u=g)
	if len(permStr) == 1 && (permStr[0] == 'u' || permStr[0] == 'g' || permStr[0] == 'o') {
		copyFrom = permStr[0]
		return who, op, 0, 0, copyFrom, nil
	}
	for _, c := range permStr {
		switch c {
		case 'r': perm |= 4
		case 'w': perm |= 2
		case 'x': perm |= 1
		case 's': special = 's'
		case 't': special = 't'
		}
	}
	return who, op, perm, special, copyFrom, nil
} 