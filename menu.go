
package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)

const (
	reset      = "\033[0m"
	bold       = "\033[1m"
	dimWhite   = "\033[38;2;180;180;180m"
	fullWhite  = "\033[38;2;255;255;255m"
	arrowColor = "\033[38;2;255;180;0m"
	treeGreen  = "\033[38;2;80;220;120m"
	clearLine  = "\033[2K\r"
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
)

// Option is a selectable menu item.
type Option struct {
	Label string
	Value string
}

// Run displays an interactive fuzzy menu and returns the selected Option.
func Run(prompt string, options []Option) (Option, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return Option{}, fmt.Errorf("raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	query := ""
	cursor := 0

	ranked := func() []Option {
		if query == "" {
			out := make([]Option, len(options))
			copy(out, options)
			return out
		}
		q := strings.ToLower(query)
		type scored struct {
			opt   Option
			score int
			idx   int
		}
		var items []scored
		for i, o := range options {
			lower := strings.ToLower(o.Label)
			var score int
			if strings.HasPrefix(lower, q) {
				score = 0
			} else if strings.Contains(lower, q) {
				score = 1
			} else {
				score = 2
			}
			items = append(items, scored{o, score, i})
		}
		sort.SliceStable(items, func(a, b int) bool {
			if items[a].score != items[b].score {
				return items[a].score < items[b].score
			}
			return items[a].idx < items[b].idx
		})
		out := make([]Option, len(items))
		for i, it := range items {
			out[i] = it.opt
		}
		return out
	}

	termHeight := func() int {
		_, h, err := term.GetSize(fd)
		if err != nil || h <= 0 {
			return 24
		}
		return h
	}

	lineCount := 0

	erase := func() {
		for i := 0; i < lineCount; i++ {
			fmt.Printf("\033[A%s", clearLine)
		}
	}

	render := func(sorted []Option) {
		lines := 0

		fmt.Printf("%s%s%s%s\r\n", bold, fullWhite, prompt, reset)
		lines++

		fmt.Printf("%s  > %s%s%s\r\n", dimWhite, fullWhite, query, reset)
		lines++

		maxItems := termHeight() - 3
		if maxItems < 1 {
			maxItems = 1
		}

		visible := sorted
		truncated := 0
		if len(sorted) > maxItems {
			truncated = len(sorted) - maxItems
			visible = sorted[:maxItems]
		}

		for i, opt := range visible {
			fmt.Print(clearLine)
			if i == cursor {
				rendered := greenTreeChars(opt.Label, query)
				fmt.Printf(
					"%s%s ❯ %s%s%s\r\n",
					arrowColor, bold,
					reset,
					rendered,
					reset,
				)
			} else {
				fmt.Printf("     %s%s%s\r\n", dimWhite, opt.Label, reset)
			}
			lines++
		}

		if truncated > 0 {
			fmt.Printf("%s  ... %d more (type to filter)%s\r\n", dimWhite, truncated, reset)
			lines++
		}

		lineCount = lines
	}

	sorted := ranked()
	render(sorted)

	buf := make([]byte, 8)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return Option{}, err
		}
		b := buf[:n]

		erase()

		switch {
		case isCtrlC(b) || isEsc(b):
			fmt.Print("\r\n")
			return Option{}, fmt.Errorf("cancelled")

		case isArrowUp(b):
			if cursor > 0 {
				cursor--
			}

		case isArrowDown(b):
			maxVisible := termHeight() - 3
			if maxVisible < 1 {
				maxVisible = 1
			}
			limit := len(sorted)
			if limit > maxVisible {
				limit = maxVisible
			}
			if cursor < limit-1 {
				cursor++
			}

		case isEnter(b):
			sorted = ranked()
			if len(sorted) == 0 {
				sorted = options
			}
			if cursor >= len(sorted) {
				cursor = len(sorted) - 1
			}
			if cursor < 0 {
				cursor = 0
			}
			selected := sorted[cursor]
			fmt.Printf("%s%s%s\r\n", fullWhite, selected.Label, reset)
			return selected, nil

		case isBackspace(b):
			if len(query) > 0 {
				query = query[:len(query)-1]
				cursor = 0
			}

		default:
			if isPrintable(b) {
				query += string(b[:n])
				cursor = 0
			}
		}

		sorted = ranked()
		render(sorted)
	}
}

// Prompt asks the user to type a value with a label, returns the entered string.
func Prompt(label string) (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	fmt.Print(showCursor)
	defer fmt.Print(hideCursor)

	fmt.Printf("%s%s%s ", bold+fullWhite, label, reset)

	value := ""
	buf := make([]byte, 8)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "", err
		}
		b := buf[:n]

		switch {
		case isCtrlC(b) || isEsc(b):
			fmt.Print("\r\n")
			return "", fmt.Errorf("cancelled")
		case isEnter(b):
			fmt.Print(hideCursor)
			fmt.Print("\r\n")
			return value, nil
		case isBackspace(b):
			if len(value) > 0 {
				value = value[:len(value)-1]
				fmt.Print("\b \b")
			}
		default:
			if isPrintable(b) {
				ch := string(b[:n])
				value += ch
				fmt.Printf("%s%s%s", fullWhite, ch, reset)
			}
		}
	}
}

// PromptSecret asks the user to type a value without echoing input to the terminal.
func PromptSecret(label string) (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	fmt.Print(showCursor)
	defer fmt.Print(hideCursor)

	fmt.Printf("%s%s%s ", bold+fullWhite, label, reset)

	value := ""
	buf := make([]byte, 8)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "", err
		}
		b := buf[:n]

		switch {
		case isCtrlC(b) || isEsc(b):
			fmt.Print("\r\n")
			return "", fmt.Errorf("cancelled")
		case isEnter(b):
			fmt.Print(hideCursor)
			fmt.Print("\r\n")
			return value, nil
		case isBackspace(b):
			if len(value) > 0 {
				value = value[:len(value)-1]
				fmt.Print("\b \b")
			}
		default:
			if isPrintable(b) {
				value += string(b[:n])
				fmt.Print("*")
			}
		}
	}
}

func greenTreeChars(label, query string) string {
	plain := stripANSI(label)

	matchStart, matchEnd := -1, -1
	if query != "" {
		lp := strings.ToLower(plain)
		lq := strings.ToLower(query)
		idx := strings.Index(lp, lq)
		if idx >= 0 {
			matchStart = idx
			matchEnd = idx + len(query)
		}
	}

	treeRunes := map[rune]bool{
		'│': true, '├': true, '└': true, '─': true,
		'▸': true, '⬡': true, '❯': true,
	}

	var sb strings.Builder
	for i, r := range plain {
		inMatch := matchStart >= 0 && i >= matchStart && i < matchEnd
		isTree := treeRunes[r]

		switch {
		case inMatch:
			sb.WriteString(treeGreen)
			sb.WriteRune(r)
			sb.WriteString(fullWhite)
		case isTree:
			sb.WriteString(treeGreen)
			sb.WriteRune(r)
			sb.WriteString(fullWhite)
		default:
			sb.WriteRune(r)
		}
	}

	return fullWhite + sb.String()
}

func stripANSI(s string) string {
	var sb strings.Builder
	inEsc := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if s[i] == 'm' {
				inEsc = false
			}
			continue
		}
		sb.WriteByte(s[i])
	}
	return sb.String()
}

func isEnter(b []byte) bool     { return len(b) == 1 && (b[0] == '\r' || b[0] == '\n') }
func isBackspace(b []byte) bool { return len(b) == 1 && (b[0] == 127 || b[0] == 8) }
func isCtrlC(b []byte) bool     { return len(b) == 1 && b[0] == 3 }
func isEsc(b []byte) bool       { return len(b) == 1 && b[0] == 27 }
func isPrintable(b []byte) bool { return len(b) == 1 && b[0] >= 32 && b[0] < 127 }
func isArrowUp(b []byte) bool   { return len(b) == 3 && b[0] == 27 && b[1] == '[' && b[2] == 'A' }
func isArrowDown(b []byte) bool { return len(b) == 3 && b[0] == 27 && b[1] == '[' && b[2] == 'B' }
