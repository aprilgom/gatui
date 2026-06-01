package testbackend

import (
	"reflect"
	"testing"

	"gatui/layout"
	"gatui/terminal"
)

func TestTestBackend_WindowSize_shouldReturnBufferSizeAndDefaultPixels(t *testing.T) {
	backend := New(80, 24)

	got, err := backend.WindowSize()
	if err != nil {
		t.Fatalf("WindowSize() error = %v", err)
	}

	want := terminal.WindowSize{
		ColumnsRows: layout.Size{Width: 80, Height: 24},
		Pixels:      layout.Size{Width: 640, Height: 480},
	}
	if got != want {
		t.Fatalf("WindowSize() = %+v, want %+v", got, want)
	}
}

func TestNoScrollBackend_WindowSize_shouldDelegateToWrappedBackend(t *testing.T) {
	backend := NewNoScroll(10, 5)

	got, err := backend.WindowSize()
	if err != nil {
		t.Fatalf("WindowSize() error = %v", err)
	}

	want := terminal.WindowSize{
		ColumnsRows: layout.Size{Width: 10, Height: 5},
		Pixels:      layout.Size{Width: 640, Height: 480},
	}
	if got != want {
		t.Fatalf("WindowSize() = %+v, want %+v", got, want)
	}
}

func TestBackend_ClearRegion_beforeCursor(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 5, Y: 3}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.ClearRegion(terminal.ClearBeforeCursor); err != nil {
		t.Fatalf("ClearRegion(ClearBeforeCursor) error = %v", err)
	}

	want := []string{
		"          ",
		"          ",
		"          ",
		"      aaaa",
		"aaaaaaaaaa",
	}
	if got := backend.Lines(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Lines() = %#v, want %#v", got, want)
	}
}

func TestBackend_ClearRegion_untilNewLine(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 3, Y: 0}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.ClearRegion(terminal.ClearUntilNewLine); err != nil {
		t.Fatalf("ClearRegion(ClearUntilNewLine) error = %v", err)
	}

	want := []string{
		"aaa       ",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	}
	if got := backend.Lines(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Lines() = %#v, want %#v", got, want)
	}
}
