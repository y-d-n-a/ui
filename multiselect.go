
package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)

// RunMulti displays an interactive multi-select menu.
// The user toggles items with space and confirms with enter.
func RunMulti(prompt string, options []Option) ([]Option, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	query := ""
	cursor := 0

	selectedOrder := []string{}
	selectedSet := map[string]bool{}

	toggle := func(val string) {
		if selectedSet[val] {
			selectedSet[val] = false
			for i, v := range selectedOrder {
				if v == val {
					selectedOrder = append(selectedOrder[:i], selectedOrder[i+1:]...)
					break
				}
			}
		} else {
			selectedSet[val] = true
			selectedOrder = append(selectedOrder, val)
		}
	}

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

		hint := fmt.Sprintf("%s  [space] toggle  [enter] confirm  [%d selected]%s", dimWhite, len(selectedOrder), reset)
		fmt.Printf("%s\r\n", hint)
		lines++

		fmt.Printf("%s  > %s%s%s\r\n", dimWhite, fullWhite, query, reset)
		lines++

		maxItems := termHeight() - 4
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
			check := " "
			if selectedSet[opt.Value] {
				check = "✓"
			}
			if i == cursor {
				rendered := greenTreeChars(opt.Label, query)
				fmt.Printf(
					"%s%s ❯ %s[%s] %s%s\r\n",
					arrowColor, bold,
					reset,
					check,
					rendered,
					reset,
				)
			} else {
				fmt.Printf("     %s[%s] %s%s\r\n", dimWhite, check, opt.Label, reset)
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
			return nil, err
		}
		b := buf[:n]

		erase()

		switch {
		case isCtrlC(b) || isEsc(b):
			fmt.Print("\r\n")
			return nil, fmt.Errorf("cancelled")

		case isArrowUp(b):
			if cursor > 0 {
				cursor--
			}

		case isArrowDown(b):
			maxVisible := termHeight() - 4
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

		case isSpace(b):
			if cursor < len(sorted) {
				toggle(sorted[cursor].Value)
			}

		case isEnter(b):
			if len(selectedOrder) == 0 {
				sorted = ranked()
				if len(sorted) > 0 {
					if cursor >= len(sorted) {
						cursor = len(sorted) - 1
					}
					sel := sorted[cursor]
					fmt.Printf("%s%s%s\r\n", fullWhite, sel.Label, reset)
					return []Option{sel}, nil
				}
				return nil, fmt.Errorf("nothing selected")
			}
			result := make([]Option, 0, len(selectedOrder))
			valToOpt := make(map[string]Option, len(options))
			for _, o := range options {
				valToOpt[o.Value] = o
			}
			for _, v := range selectedOrder {
				result = append(result, valToOpt[v])
			}
			fmt.Print("\r\n")
			return result, nil

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

func isSpace(b []byte) bool { return len(b) == 1 && b[0] == ' ' }
