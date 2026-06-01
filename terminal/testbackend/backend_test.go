package testbackend

import (
	"os"
	"os/exec"
	"reflect"
	"testing"

	"gatui/buffer"
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

func TestTestBackend_Buffer_shouldReturnCurrentBuffer(t *testing.T) {
	backend := WithLines([]string{
		"abc",
		"def",
	})

	got := backend.Buffer()

	want := buffer.WithLines([]string{
		"abc",
		"def",
	})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Buffer() = %#v, want %#v", got, want)
	}
}

func TestTestBackend_Scrollback_shouldReturnScrollbackBuffer(t *testing.T) {
	backend := WithLines([]string{
		"aaaa",
		"bbbb",
		"cccc",
	})

	if err := backend.ScrollRegionUp(0, 3, 2); err != nil {
		t.Fatalf("ScrollRegionUp() error = %v", err)
	}

	got := backend.Scrollback()
	want := buffer.WithLines([]string{
		"aaaa",
		"bbbb",
	})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Scrollback() = %#v, want %#v", got, want)
	}
}

func TestTestBackend_AssertBufferLines_shouldPassForMatchingLines(t *testing.T) {
	backend := WithLines([]string{
		"abc",
		"def",
	})

	backend.AssertBufferLines(t, []string{
		"abc",
		"def",
	})
}

func TestTestBackend_AssertBufferLines_shouldFailForMismatchedLines(t *testing.T) {
	if os.Getenv("GATUI_ASSERT_BUFFER_LINES_MISMATCH") == "1" {
		backend := WithLines([]string{
			"abc",
			"def",
		})
		backend.AssertBufferLines(t, []string{
			"abc",
			"xyz",
		})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestTestBackend_AssertBufferLines_shouldFailForMismatchedLines$")
	cmd.Env = append(os.Environ(), "GATUI_ASSERT_BUFFER_LINES_MISMATCH=1")
	if err := cmd.Run(); err == nil {
		t.Fatal("AssertBufferLines() unexpectedly passed")
	}
}

func TestTestBackend_AssertScrollbackLines_shouldPassForMatchingLines(t *testing.T) {
	backend := WithLines([]string{
		"aaaa",
		"bbbb",
		"cccc",
	})

	if err := backend.ScrollRegionUp(0, 3, 1); err != nil {
		t.Fatalf("ScrollRegionUp() error = %v", err)
	}

	backend.AssertScrollbackLines(t, []string{
		"aaaa",
	})
}

func TestTestBackend_AssertScrollbackEmpty_shouldPassWhenEmpty(t *testing.T) {
	backend := New(4, 2)

	backend.AssertScrollbackEmpty(t)
}

func TestTestBackend_AssertCursorPosition_shouldPassForMatchingPosition(t *testing.T) {
	backend := New(4, 2)
	pos := layout.Position{X: 2, Y: 1}
	if err := backend.SetCursorPosition(pos); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	backend.AssertCursorPosition(t, pos)
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
