package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/term"
)

const tabWidth = 4

// Minimal vim-like editor: normal/insert mode, h/j/k/l, i/a/x, :w, :q, status bar
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "bse: usage: bse FILE")
		os.Exit(1)
	}
	filename := os.Args[1]
	lines := []string{""}
	if data, err := ioutil.ReadFile(filename); err == nil {
		lines = strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	}
	row, col := 0, 0
	topLine := 0
	leftCol := 0
	mode := "NORMAL"
	status := ""
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() { <-ch; term.Restore(int(os.Stdin.Fd()), oldState); os.Exit(0) }()
	stdin := bufio.NewReader(os.Stdin)
	cmd := ""
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	windowHeight := height - 1
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
				fmt.Print("\r\n")
			}
		}
		fmt.Printf("\x1b[7m--%s-- %s\x1b[0m\r\n", mode, status)
		cursorCol := visualCol(lines[row], col) - leftCol
		fmt.Printf("\x1b[%d;%dH", row-topLine+1, cursorCol+1)
		b, _ := stdin.ReadByte()
		if b == 0x1b {
			seq, _ := stdin.Peek(2)
			if len(seq) == 2 && seq[0] == '[' {
				stdin.Discard(2)
				switch seq[1] {
				case 'A':
					if row > 0 {
						row--
						col = min(col, len(lines[row]))
					}
				case 'B':
					if row < len(lines)-1 {
						row++
						col = min(col, len(lines[row]))
					}
				case 'C':
					if col < len(lines[row]) {
						col++
					}
				case 'D':
					if col > 0 {
						col--
					}
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
			case 'a':
				mode = "INSERT"
				if col < len(lines[row]) {
					col++
				}
			case 'x':
				if len(lines[row]) > 0 && col < len(lines[row]) {
					lines[row] = lines[row][:col] + lines[row][col+1:]
				}
			case ':':
				mode = "CMD"
				cmd = ":"
				status = cmd
			}
		} else if mode == "INSERT" {
			if b == 27 {
				mode = "NORMAL"
				continue
			}
			if b == 127 || b == 8 {
				if col > 0 {
					lines[row] = lines[row][:col-1] + lines[row][col:]
					col--
				}
				continue
			}
			if b == '\r' || b == '\n' {
				rest := lines[row][col:]
				lines[row] = lines[row][:col]
				lines = append(lines[:row+1], append([]string{""}, lines[row+1:]...)...)
				lines[row+1] = rest + lines[row+1]
				row++
				col = 0
				continue
			}
			lines[row] = lines[row][:col] + string(b) + lines[row][col:]
			col++
		} else if mode == "CMD" {
			if b == '\r' || b == '\n' {
				if cmd == ":q" {
					clearScreen()
					return
				}
				if cmd == ":w" {
					err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
					if err != nil {
						status = "write error"
					} else {
						status = "[wrote]"
					}
				}
				mode = "NORMAL"
				cmd = ""
				continue
			}
			if b == 127 || b == 8 {
				if len(cmd) > 1 {
					cmd = cmd[:len(cmd)-1]
				}
			}
			if b != 127 && b != 8 {
				cmd += string(b)
			}
			status = cmd
		}
		col = min(col, len(lines[row]))
	}
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
	fmt.Printf("%-*s\r\n", width, display)
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
	fmt.Print("\x1b[7m")
	fmt.Printf("%-*s", width, display)
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
