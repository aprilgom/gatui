package terminal_test

import (
	"reflect"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/terminal"
	"gatui/text"
	"gatui/widgets"
)

type recordingBackend struct {
	size            layout.Size
	draws           [][]buffer.CellDiff
	flushCount      int
	clearCount      int
	hideCursorCount int
	showCursorCount int
	cursorPositions []layout.Position
}

func newRecordingBackend(width, height int) *recordingBackend {
	return &recordingBackend{size: layout.Size{Width: width, Height: height}}
}

func (b *recordingBackend) Size() (layout.Size, error) {
	return b.size, nil
}

func (b *recordingBackend) Draw(diffs []buffer.CellDiff) error {
	copied := make([]buffer.CellDiff, len(diffs))
	copy(copied, diffs)
	b.draws = append(b.draws, copied)
	return nil
}

func (b *recordingBackend) Flush() error {
	b.flushCount++
	return nil
}

func (b *recordingBackend) Clear() error {
	b.clearCount++
	return nil
}

func (b *recordingBackend) PollEvent() (terminal.Event, error) {
	return terminal.UnknownEvent{}, nil
}

func (b *recordingBackend) HideCursor() error {
	b.hideCursorCount++
	return nil
}

func (b *recordingBackend) ShowCursor() error {
	b.showCursorCount++
	return nil
}

func (b *recordingBackend) SetCursorPosition(pos layout.Position) error {
	b.cursorPositions = append(b.cursorPositions, pos)
	return nil
}

func TestTerminal_New_shouldCreateBuffersFromBackendSize(t *testing.T) {
	term, err := terminal.New(newRecordingBackend(5, 3))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	completed, err := term.Draw(func(frame *terminal.Frame) {
		if got, want := frame.Area(), layout.NewRect(0, 0, 5, 3); got != want {
			t.Fatalf("frame area = %#v, want %#v", got, want)
		}
		if got, want := frame.Buffer().Area, layout.NewRect(0, 0, 5, 3); got != want {
			t.Fatalf("buffer area = %#v, want %#v", got, want)
		}
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := completed.Area, layout.NewRect(0, 0, 5, 3); got != want {
		t.Fatalf("completed area = %#v, want %#v", got, want)
	}
	if got, want := completed.Count, 0; got != want {
		t.Fatalf("completed count = %d, want %d", got, want)
	}
}

func TestTerminal_Draw_shouldRenderWidgetAndFlushDiff(t *testing.T) {
	backend := newRecordingBackend(5, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	completed, err := term.Draw(func(frame *terminal.Frame) {
		frame.RenderWidget(widgets.NewParagraph(text.FromString("abc")), frame.Area())
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := completed.Buffer.Lines(), []string{"abc  "}; !reflect.DeepEqual(got, want) {
		t.Fatalf("completed lines = %#v, want %#v", got, want)
	}
	if got, want := len(backend.draws), 1; got != want {
		t.Fatalf("draw call count = %d, want %d", got, want)
	}
	if got, want := diffSymbols(backend.draws[0]), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("draw diff symbols = %#v, want %#v", got, want)
	}
	if got, want := backend.flushCount, 1; got != want {
		t.Fatalf("flush count = %d, want %d", got, want)
	}
}

func TestTerminal_Draw_shouldOnlySendChangedCellsOnSecondDraw(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	renderText(t, term, "abc")
	renderText(t, term, "axc")

	if got, want := backend.draws[1], []buffer.CellDiff{{X: 1, Y: 0, Cell: buffer.NewCell("x")}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("second draw diffs = %#v, want %#v", got, want)
	}
}

func TestTerminal_Draw_shouldSwapAndResetBuffers(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	renderText(t, term, "abc")
	renderText(t, term, "axc")
	renderText(t, term, "axc")

	if got, want := backend.draws[2], []buffer.CellDiff{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("third draw diffs = %#v, want %#v", got, want)
	}
}

func TestFrame_RenderWidget_shouldRenderIntoCurrentBuffer(t *testing.T) {
	term, err := terminal.New(newRecordingBackend(4, 1))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	completed, err := term.Draw(func(frame *terminal.Frame) {
		frame.RenderWidget(widgets.NewParagraph(text.FromString("go")), frame.Area())
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := completed.Buffer.Lines(), []string{"go  "}; !reflect.DeepEqual(got, want) {
		t.Fatalf("buffer lines = %#v, want %#v", got, want)
	}
}

func TestTerminal_Draw_shouldHideCursorWhenUnset(t *testing.T) {
	backend := newRecordingBackend(2, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	renderText(t, term, "a")

	if got, want := backend.hideCursorCount, 1; got != want {
		t.Fatalf("hide cursor count = %d, want %d", got, want)
	}
	if got, want := backend.showCursorCount, 0; got != want {
		t.Fatalf("show cursor count = %d, want %d", got, want)
	}
}

func TestTerminal_Draw_shouldShowAndPositionCursorWhenSet(t *testing.T) {
	backend := newRecordingBackend(2, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	_, err = term.Draw(func(frame *terminal.Frame) {
		frame.SetCursorPosition(layout.Position{X: 1, Y: 0})
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := backend.showCursorCount, 1; got != want {
		t.Fatalf("show cursor count = %d, want %d", got, want)
	}
	if got, want := backend.hideCursorCount, 0; got != want {
		t.Fatalf("hide cursor count = %d, want %d", got, want)
	}
	if got, want := backend.cursorPositions, []layout.Position{{X: 1, Y: 0}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cursor positions = %#v, want %#v", got, want)
	}
}

func TestTerminal_Resize_shouldResizeBothBuffers(t *testing.T) {
	backend := newRecordingBackend(5, 3)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	term.Resize(layout.NewRect(0, 0, 4, 2))
	completed, err := term.Draw(func(frame *terminal.Frame) {
		if got, want := frame.Area(), layout.NewRect(0, 0, 4, 2); got != want {
			t.Fatalf("frame area = %#v, want %#v", got, want)
		}
		frame.RenderWidget(widgets.NewParagraph(text.FromString("abcd\nxy")), frame.Area())
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := completed.Buffer.Area, layout.NewRect(0, 0, 4, 2); got != want {
		t.Fatalf("completed buffer area = %#v, want %#v", got, want)
	}
	if got, want := completed.Buffer.Lines(), []string{"abcd", "xy  "}; !reflect.DeepEqual(got, want) {
		t.Fatalf("completed lines = %#v, want %#v", got, want)
	}
}

func TestTerminal_Clear_shouldClearBackendAndForceFullRedraw(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	renderText(t, term, "abc")

	if err := term.Clear(); err != nil {
		t.Fatalf("Clear returned error: %v", err)
	}
	renderText(t, term, "abc")

	if got, want := backend.clearCount, 1; got != want {
		t.Fatalf("clear count = %d, want %d", got, want)
	}
	if got, want := diffSymbols(backend.draws[1]), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("second draw diff symbols = %#v, want %#v", got, want)
	}
}

func renderText(t *testing.T, term *terminal.Terminal, content string) {
	t.Helper()
	_, err := term.Draw(func(frame *terminal.Frame) {
		frame.RenderWidget(widgets.NewParagraph(text.FromString(content)), frame.Area())
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
}

func diffSymbols(diffs []buffer.CellDiff) []string {
	symbols := make([]string, 0, len(diffs))
	for _, diff := range diffs {
		symbols = append(symbols, diff.Cell.Symbol)
	}
	return symbols
}
