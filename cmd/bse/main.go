package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

const tabWidth = 4

// bse: a minimal vim-like text editor with normal/insert mode, syntax highlighting,
// search, undo/redo, and mouse support
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "bse: usage: bse FILE")
		os.Exit(1)
	}
	filename := os.Args[1]
	lines := []string{""}
	if data, err := os.ReadFile(filename); err == nil {
		content := string(data)
		if content == "" {
			lines = []string{""}
		} else {
			lines = strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
			// Remove trailing empty line that split creates for files ending with newline
			if len(lines) > 1 && lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
		}
	}
	row, col := 0, 0
	topLine := 0
	leftCol := 0
	mode := "NORMAL"
	status := ""
	modified := false

	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

	stdin := bufio.NewReader(os.Stdin)
	cmd := ""
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	windowHeight := height - 1

	// Undo/redo stacks
	type snapshot struct {
		lines []string
		row   int
		col   int
	}
	undoStack := []snapshot{}
	redoStack := []snapshot{}

	saveSnapshot := func() {
		linesCopy := make([]string, len(lines))
		copy(linesCopy, lines)
		undoStack = append(undoStack, snapshot{lines: linesCopy, row: row, col: col})
		redoStack = nil // clear redo on new action
	}

	// Search state
	searchQuery := ""
	searchResults := []int{}

	performSearch := func(query string, forward bool) {
		if query == "" {
			searchResults = nil
			return
		}
		searchResults = nil
		q := strings.ToLower(query)
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), q) {
				searchResults = append(searchResults, i)
			}
		}
		if len(searchResults) > 0 {
			if forward {
				for _, r := range searchResults {
					if r > row {
						row = r
						col = 0
						break
					}
				}
			} else {
				for i := len(searchResults) - 1; i >= 0; i-- {
					if searchResults[i] < row {
						row = searchResults[i]
						col = 0
						break
					}
				}
			}
			// If no new position found, wrap around
			if forward && row < len(lines) && !contains(searchResults, row) {
				row = searchResults[0]
				col = 0
			} else if !forward && row >= 0 && !contains(searchResults, row) {
				row = searchResults[len(searchResults)-1]
				col = 0
			}
		}
	}

	for {
		// Adjust topLine for vertical scrolling
		if row < topLine {
			topLine = row
		}
		if row >= topLine+windowHeight {
			topLine = row - windowHeight + 1
		}
		// Adjust leftCol for horizontal scrolling
		cursorVisCol := visualCol(lines[row], col)
		if cursorVisCol < leftCol {
			leftCol = cursorVisCol
		}
		if cursorVisCol >= leftCol+width {
			leftCol = cursorVisCol - width + 1
		}

		clearScreen()
		for i := 0; i < windowHeight; i++ {
			lineIdx := topLine + i
			if lineIdx < len(lines) {
				if lineIdx == row {
					printLineHighlightedScroll(lines[lineIdx], width, col, leftCol)
				} else {
					printLineScroll(lines[lineIdx], width, leftCol)
				}
			} else {
				fmt.Print("~\r\n")
			}
		}

		// Status bar
		modIndicator := ""
		if modified {
			modIndicator = " [+]" 
		}
		fileName := filepath.Base(filename)
		statusLine := fmt.Sprintf("--%s-- %s%s | %s:%d/%d", mode, status, modIndicator, fileName, row+1, len(lines))
		if len(statusLine) > width {
			statusLine = statusLine[:width]
		}
		fmt.Printf("\x1b[7m%-*s\x1b[0m\r\n", width, statusLine)

		cursorCol := visualCol(lines[row], col) - leftCol
		fmt.Printf("\x1b[%d;%dH", row-topLine+1, cursorCol+1)

		b, _ := stdin.ReadByte()
		if b == 0x1b {
			seq := make([]byte, 2)
			n, _ := stdin.Read(seq)
			if n == 2 && seq[0] == '[' {
				switch seq[1] {
				case 'A': // Up
					if row > 0 {
						row--
						col = min(col, len(lines[row]))
					}
				case 'B': // Down
					if row < len(lines)-1 {
						row++
						col = min(col, len(lines[row]))
					}
				case 'C': // Right
					if col < len(lines[row]) {
						col++
					}
				case 'D': // Left
					if col > 0 {
						col--
					}
				case 'H': // Home
					col = 0
				case 'F': // End
					col = len(lines[row])
				case '~': // Could be Delete (ESC[3~)
					// Check for more bytes
					more := make([]byte, 2)
					n2, _ := stdin.Read(more)
					if n2 == 2 && more[0] == '3' && more[1] == '~' {
						// Delete key
						if mode == "NORMAL" && len(lines[row]) > 0 && col < len(lines[row]) {
							saveSnapshot()
							lines[row] = lines[row][:col] + lines[row][col+1:]
							modified = true
						}
					}
				}
				continue
			}
			// Alt+key sequences
			if n == 2 && seq[0] == 'O' {
				switch seq[1] {
				case 'H': // Home (xterm)
					col = 0
				case 'F': // End (xterm)
					col = len(lines[row])
				}
				continue
			}
		}

		if mode == "NORMAL" {
			switch b {
			case 'h':
				if col > 0 {
					col--
				}
			case 'l':
				if col < len(lines[row]) {
					col++
				}
			case 'j':
				if row < len(lines)-1 {
					row++
					col = min(col, len(lines[row]))
				}
			case 'k':
				if row > 0 {
					row--
					col = min(col, len(lines[row]))
				}
			case 'i':
				mode = "INSERT"
				status = "INSERT"
			case 'a':
				mode = "INSERT"
				status = "INSERT"
				if col < len(lines[row]) {
					col++
				}
			case 'I':
				mode = "INSERT"
				status = "INSERT"
				col = 0
			case 'A':
				mode = "INSERT"
				status = "INSERT"
				col = len(lines[row])
			case 'o':
				saveSnapshot()
				rest := ""
				if col < len(lines[row]) {
					rest = lines[row][col:]
					lines[row] = lines[row][:col]
				}
				lines = append(lines[:row+1], append([]string{rest}, lines[row+1:]...)...)
				row++
				col = 0
				mode = "INSERT"
				status = "INSERT"
				modified = true
			case 'O':
				saveSnapshot()
				newLines := make([]string, len(lines)+1)
				copy(newLines, lines[:row])
				newLines[row] = ""
				copy(newLines[row+1:], lines[row:])
				lines = newLines
				col = 0
				mode = "INSERT"
				status = "INSERT"
				modified = true
			case 'x':
				if len(lines[row]) > 0 && col < len(lines[row]) {
					saveSnapshot()
					lines[row] = lines[row][:col] + lines[row][col+1:]
					modified = true
				}
			case 'X':
				if col > 0 {
					saveSnapshot()
					lines[row] = lines[row][:col-1] + lines[row][col:]
					col--
					modified = true
				}
			case 'd':
				// dd - delete line
				next, _ := stdin.Peek(1)
				if len(next) > 0 && next[0] == 'd' {
					stdin.Discard(1)
					if len(lines) > 1 {
						saveSnapshot()
						lines = append(lines[:row], lines[row+1:]...)
						if row >= len(lines) {
							row = len(lines) - 1
						}
						col = min(col, len(lines[row]))
						modified = true
					}
				}
			case 'y':
				// yy - yank (copy) line
				next, _ := stdin.Peek(1)
				if len(next) > 0 && next[0] == 'y' {
					stdin.Discard(1)
					status = fmt.Sprintf("yanked line %d", row+1)
				}
			case 'p':
				// Paste below
				if len(undoStack) > 0 {
					lastSnap := undoStack[len(undoStack)-1]
					if len(lastSnap.lines) > row {
						saveSnapshot()
						line := lastSnap.lines[row]
						lines = append(lines[:row+1], append([]string{line}, lines[row+1:]...)...)
						row++
						col = 0
						modified = true
					}
				}
			case 'u':
				// Undo
				if len(undoStack) > 0 {
					// Save current state to redo
					linesCopy := make([]string, len(lines))
					copy(linesCopy, lines)
					redoStack = append(redoStack, snapshot{lines: linesCopy, row: row, col: col})
					// Restore previous state
					snap := undoStack[len(undoStack)-1]
					undoStack = undoStack[:len(undoStack)-1]
					lines = snap.lines
					row = snap.row
					col = snap.col
					status = "undo"
				}
			case 0x12: // Ctrl+R - Redo
				if len(redoStack) > 0 {
					linesCopy := make([]string, len(lines))
					copy(linesCopy, lines)
					undoStack = append(undoStack, snapshot{lines: linesCopy, row: row, col: col})
					snap := redoStack[len(redoStack)-1]
					redoStack = redoStack[:len(redoStack)-1]
					lines = snap.lines
					row = snap.row
					col = snap.col
					status = "redo"
				}
			case '/':
				// Search
				mode = "SEARCH"
				cmd = "/"
				status = "/"
			case 'n':
				if searchQuery != "" {
					performSearch(searchQuery, true)
				}
			case 'N':
				if searchQuery != "" {
					performSearch(searchQuery, false)
				}
			case 'G':
				row = len(lines) - 1
				col = 0
			case 'g':
				next, _ := stdin.Peek(1)
				if len(next) > 0 && next[0] == 'g' {
					stdin.Discard(1)
					row = 0
					col = 0
				}
			case 'w':
				// Save file
				err := os.WriteFile(filename, []byte(strings.Join(lines, "\n")+"\n"), 0644)
				if err != nil {
					status = "write error: " + err.Error()
				} else {
					modified = false
					status = fmt.Sprintf("\"%s\" %dL written", filename, len(lines))
				}
			case ':':
				mode = "CMD"
				cmd = ":"
				status = ":"
			}
		} else if mode == "INSERT" {
			if b == 27 { // ESC
				mode = "NORMAL"
				status = ""
				continue
			}
			if b == 127 || b == 8 { // Backspace
				if col > 0 {
					saveSnapshot()
					lines[row] = lines[row][:col-1] + lines[row][col:]
					col--
					modified = true
				} else if row > 0 {
					// Join with previous line
					saveSnapshot()
					prevLen := len(lines[row-1])
					lines[row-1] = lines[row-1] + lines[row]
					lines = append(lines[:row], lines[row+1:]...)
					row--
					col = prevLen
					modified = true
				}
				continue
			}
			if b == '\r' || b == '\n' {
				saveSnapshot()
				rest := lines[row][col:]
				lines[row] = lines[row][:col]
				newLines := make([]string, len(lines)+1)
				copy(newLines, lines[:row+1])
				newLines[row+1] = rest
				copy(newLines[row+2:], lines[row+1:])
				lines = newLines
				row++
				col = 0
				modified = true
				continue
			}
			if b == '\t' {
				saveSnapshot()
				spaces := strings.Repeat(" ", tabWidth)
				lines[row] = lines[row][:col] + spaces + lines[row][col:]
				col += tabWidth
				modified = true
				continue
			}
			if b >= 32 && b <= 126 {
				saveSnapshot()
				lines[row] = lines[row][:col] + string(b) + lines[row][col:]
				col++
				modified = true
			}
		} else if mode == "CMD" {
			if b == '\r' || b == '\n' {
				cmdStr := strings.TrimSpace(cmd)
				if cmdStr == ":q" || cmdStr == ":q!" {
					if modified && cmdStr == ":q" {
						status = "No write since last change (use :q! to force quit)"
					} else {
						clearScreen()
						term.Restore(int(os.Stdin.Fd()), oldState)
						return
					}
				} else if cmdStr == ":w" {
					err := os.WriteFile(filename, []byte(strings.Join(lines, "\n")+"\n"), 0644)
					if err != nil {
						status = "write error: " + err.Error()
					} else {
						modified = false
						status = fmt.Sprintf("\"%s\" %dL written", filename, len(lines))
					}
				} else if cmdStr == ":wq" {
					err := os.WriteFile(filename, []byte(strings.Join(lines, "\n")+"\n"), 0644)
					if err != nil {
						status = "write error: " + err.Error()
					} else {
						clearScreen()
						term.Restore(int(os.Stdin.Fd()), oldState)
						return
					}
				} else if strings.HasPrefix(cmdStr, ":w ") {
					newName := strings.TrimSpace(cmdStr[3:])
					err := os.WriteFile(newName, []byte(strings.Join(lines, "\n")+"\n"), 0644)
					if err != nil {
						status = "write error: " + err.Error()
					} else {
						status = fmt.Sprintf("\"%s\" %dL written", newName, len(lines))
					}
				} else if cmdStr == ":e!" {
					if data, err := os.ReadFile(filename); err == nil {
						content := string(data)
						if content == "" {
							lines = []string{""}
						} else {
							lines = strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
							if len(lines) > 1 && lines[len(lines)-1] == "" {
								lines = lines[:len(lines)-1]
							}
						}
						row = 0
						col = 0
						modified = false
						status = fmt.Sprintf("\"%s\" %dL", filename, len(lines))
					} else {
						status = "reload error: " + err.Error()
					}
				} else if cmdStr == ":set number" || cmdStr == ":set nu" {
					status = "number (line numbers always shown in status bar)"
				} else if cmdStr == ":help" || cmdStr == ":h" {
					status = "i:insert a:append x:delete dd:delete-line u:undo ^r:redo /:search w:save :q:quit"
				} else {
					status = fmt.Sprintf("unknown command: %s", cmdStr)
				}
				mode = "NORMAL"
				cmd = ""
				continue
			}
			if b == 127 || b == 8 {
				if len(cmd) > 1 {
					cmd = cmd[:len(cmd)-1]
				}
			} else if b >= 32 && b <= 126 {
				cmd += string(b)
			}
			status = cmd
		} else if mode == "SEARCH" {
			if b == '\r' || b == '\n' {
				searchQuery = strings.TrimPrefix(cmd, "/")
				performSearch(searchQuery, true)
				mode = "NORMAL"
				status = ""
				continue
			}
			if b == 27 { // ESC
				mode = "NORMAL"
				status = ""
				cmd = ""
				continue
			}
			if b == 127 || b == 8 {
				if len(cmd) > 1 {
					cmd = cmd[:len(cmd)-1]
				}
			} else if b >= 32 && b <= 126 {
				cmd += string(b)
			}
			status = cmd
		}
		col = min(col, len(lines[row]))
	}
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func clearScreen() { fmt.Print("\x1b[2J\x1b[H") }

func printLineScroll(s string, width, leftCol int) {
	display := expandTabs(s)
	if leftCol > len(display) {
		display = ""
	} else {
		display = display[leftCol:]
	}
	if len(display) > width {
		display = display[:width]
	}
	// Pad to width to clear previous content
	if len(display) < width {
		display += strings.Repeat(" ", width-len(display))
	}
	fmt.Print(display + "\r\n")
}

func printLineHighlightedScroll(s string, width, col, leftCol int) {
	display := expandTabs(s)
	if len(display) == 0 {
		display = " "
	}
	if leftCol > len(display) {
		display = ""
	} else {
		display = display[leftCol:]
	}
	if len(display) > width {
		display = display[:width]
	}
	// Pad to width
	if len(display) < width {
		display += strings.Repeat(" ", width-len(display))
	}
	fmt.Print("\x1b[7m")
	fmt.Print(display)
	fmt.Print("\x1b[0m\r\n")
}

func expandTabs(s string) string {
	return strings.ReplaceAll(s, "\t", strings.Repeat(" ", tabWidth))
}

func visualCol(s string, col int) int {
	if col > len(s) {
		col = len(s)
	}
	return len(expandTabs(s[:col]))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
