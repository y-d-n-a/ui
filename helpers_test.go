
package ui

import "testing"

func TestIsEnter(t *testing.T) {
	if !isEnter([]byte{'\r'}) {
		t.Error("expected \\r to be enter")
	}
	if !isEnter([]byte{'\n'}) {
		t.Error("expected \\n to be enter")
	}
	if isEnter([]byte{'a'}) {
		t.Error("expected 'a' not to be enter")
	}
}

func TestIsBackspace(t *testing.T) {
	if !isBackspace([]byte{127}) {
		t.Error("expected 127 to be backspace")
	}
	if !isBackspace([]byte{8}) {
		t.Error("expected 8 to be backspace")
	}
	if isBackspace([]byte{'a'}) {
		t.Error("expected 'a' not to be backspace")
	}
}

func TestIsCtrlC(t *testing.T) {
	if !isCtrlC([]byte{3}) {
		t.Error("expected byte 3 to be ctrl+c")
	}
	if isCtrlC([]byte{'c'}) {
		t.Error("expected 'c' not to be ctrl+c")
	}
}

func TestIsEsc(t *testing.T) {
	if !isEsc([]byte{27}) {
		t.Error("expected byte 27 to be esc")
	}
	if isEsc([]byte{'e'}) {
		t.Error("expected 'e' not to be esc")
	}
}

func TestIsPrintable(t *testing.T) {
	for _, b := range []byte("abcABC123 !@#") {
		if !isPrintable([]byte{b}) {
			t.Errorf("expected byte %d (%q) to be printable", b, b)
		}
	}
	if isPrintable([]byte{0}) {
		t.Error("expected null byte not to be printable")
	}
	if isPrintable([]byte{127}) {
		t.Error("expected 127 not to be printable")
	}
}

func TestIsArrowUp(t *testing.T) {
	if !isArrowUp([]byte{27, '[', 'A'}) {
		t.Error("expected escape sequence to be arrow up")
	}
	if isArrowUp([]byte{27, '[', 'B'}) {
		t.Error("expected down sequence not to be arrow up")
	}
}

func TestIsArrowDown(t *testing.T) {
	if !isArrowDown([]byte{27, '[', 'B'}) {
		t.Error("expected escape sequence to be arrow down")
	}
	if isArrowDown([]byte{27, '[', 'A'}) {
		t.Error("expected up sequence not to be arrow down")
	}
}

func TestIsSpace(t *testing.T) {
	if !isSpace([]byte{' '}) {
		t.Error("expected space byte to be space")
	}
	if isSpace([]byte{'a'}) {
		t.Error("expected 'a' not to be space")
	}
}
