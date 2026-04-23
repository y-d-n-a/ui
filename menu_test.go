
package ui

import (
	"testing"
)

func TestStripANSI(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"\033[1mhello\033[0m", "hello"},
		{"\033[38;2;255;255;255mworld\033[0m", "world"},
		{"plain", "plain"},
		{"", ""},
		{"\033[1m\033[0m", ""},
	}
	for _, c := range cases {
		got := stripANSI(c.in)
		if got != c.want {
			t.Errorf("stripANSI(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGreenTreeChars_NoQuery(t *testing.T) {
	result := greenTreeChars("hello", "")
	if result == "" {
		t.Error("expected non-empty result")
	}
	plain := stripANSI(result)
	if plain != "hello" {
		t.Errorf("expected plain text 'hello', got %q", plain)
	}
}

func TestGreenTreeChars_WithQuery(t *testing.T) {
	result := greenTreeChars("foobar", "foo")
	plain := stripANSI(result)
	if plain != "foobar" {
		t.Errorf("expected plain text 'foobar', got %q", plain)
	}
}

func TestGreenTreeChars_TreeRunes(t *testing.T) {
	result := greenTreeChars("│ item", "")
	plain := stripANSI(result)
	if plain != "│ item" {
		t.Errorf("expected plain text '│ item', got %q", plain)
	}
}
