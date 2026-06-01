package terminal_test

import (
	"errors"
	"reflect"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/terminal"
	"gatui/terminal/testbackend"
	"gatui/text"
	"gatui/widgets"
)

type recordingBackend struct {
	size            layout.Size
	draws           [][]buffer.CellDiff
	flushCount      int
	clearCount      int
	clearRegions    []terminal.ClearType
	hideCursorCount int
	showCursorCount int
	cursorPositions []layout.Position
	cursorPosition  layout.Position
	appendLines     []int
	operations      []string
}

func newRecordingBackend(width, height int) *recordingBackend {
	return &recordingBackend{size: layout.Size{Width: width, Height: height}}
}

func (b *recordingBackend) Size() (layout.Size, error) {
	return b.size, nil
}

func (b *recordingBackend) SetSize(width, height int) {
	b.size = layout.Size{Width: width, Height: height}
}

func (b *recordingBackend) Draw(diffs []buffer.CellDiff) error {
	copied := make([]buffer.CellDiff, len(diffs))
	copy(copied, diffs)
	b.draws = append(b.draws, copied)
	b.operations = append(b.operations, "draw")
	return nil
}

func (b *recordingBackend) Flush() error {
	b.flushCount++
	b.operations = append(b.operations, "backend-flush")
	return nil
}

func (b *recordingBackend) Clear() error {
	b.clearCount++
	return b.ClearRegion(terminal.ClearAll)
}

func (b *recordingBackend) ClearRegion(clearType terminal.ClearType) error {
	b.clearRegions = append(b.clearRegions, clearType)
	b.operations = append(b.operations, "clear-region")
	return nil
}

func (b *recordingBackend) HideCursor() error {
	b.hideCursorCount++
	b.operations = append(b.operations, "hide-cursor")
	return nil
}

func (b *recordingBackend) ShowCursor() error {
	b.showCursorCount++
	b.operations = append(b.operations, "show-cursor")
	return nil
}

func (b *recordingBackend) SetCursorPosition(pos layout.Position) error {
	b.cursorPositions = append(b.cursorPositions, pos)
	b.cursorPosition = pos
	b.operations = append(b.operations, "set-cursor")
	return nil
}

func (b *recordingBackend) GetCursorPosition() (layout.Position, error) {
	return b.cursorPosition, nil
}

func (b *recordingBackend) AppendLines(count int) error {
	b.appendLines = append(b.appendLines, count)
	b.operations = append(b.operations, "append-lines")
	return nil
}

func TestTerminalBackendInterface_shouldNotRequireEventPolling(t *testing.T) {
	var _ terminal.Backend = (*recordingBackend)(nil)
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

func TestTerminal_NewWithOptions_fullscreenUsesBackendSize(t *testing.T) {
	term, err := terminal.NewWithOptions(newRecordingBackend(5, 3), terminal.DefaultTerminalOptions())
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
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
}

func TestTerminal_NewWithOptions_fixedUsesProvidedArea(t *testing.T) {
	area := layout.NewRect(2, 1, 4, 2)
	term, err := terminal.NewWithOptions(newRecordingBackend(10, 5), terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(area),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	completed, err := term.Draw(func(frame *terminal.Frame) {
		if got, want := frame.Area(), area; got != want {
			t.Fatalf("frame area = %#v, want %#v", got, want)
		}
		if got, want := frame.Buffer().Area, area; got != want {
			t.Fatalf("buffer area = %#v, want %#v", got, want)
		}
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := completed.Area, area; got != want {
		t.Fatalf("completed area = %#v, want %#v", got, want)
	}
	if got, want := completed.Buffer.Area, area; got != want {
		t.Fatalf("completed buffer area = %#v, want %#v", got, want)
	}
}

func TestTerminal_NewWithOptions_inlineAnchorsToCursorWhenSpaceAvailable(t *testing.T) {
	backend := newRecordingBackend(10, 10)
	backend.cursorPosition = layout.Position{X: 0, Y: 3}

	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 3, 10, 4); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	if got, want := backend.appendLines, []int{3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("append lines = %#v, want %#v", got, want)
	}
}

func TestTerminal_NewWithOptions_inlineShiftsUpWhenNearBottom(t *testing.T) {
	backend := newRecordingBackend(10, 10)
	backend.cursorPosition = layout.Position{X: 0, Y: 8}

	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 6, 10, 4); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	if got, want := backend.appendLines, []int{3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("append lines = %#v, want %#v", got, want)
	}
}

func TestTerminal_NewWithOptions_inlineClampsHeightToTerminal(t *testing.T) {
	backend := newRecordingBackend(10, 3)
	backend.cursorPosition = layout.Position{X: 0, Y: 0}

	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(10),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 0, 10, 3); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	if got, want := backend.appendLines, []int{9}; !reflect.DeepEqual(got, want) {
		t.Fatalf("append lines = %#v, want %#v", got, want)
	}
}

func TestTerminal_Draw_fixedViewportUsesFixedFrameArea(t *testing.T) {
	area := layout.NewRect(2, 1, 2, 1)
	backend := newRecordingBackend(5, 3)
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(area),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	completed, err := term.Draw(func(frame *terminal.Frame) {
		if got, want := frame.Area(), area; got != want {
			t.Fatalf("frame area = %#v, want %#v", got, want)
		}
		frame.Buffer().SetSymbol(area.X, area.Y, "z")
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := completed.Buffer.Lines(), []string{"z "}; !reflect.DeepEqual(got, want) {
		t.Fatalf("completed lines = %#v, want %#v", got, want)
	}
	if got, want := backend.draws[0], []buffer.CellDiff{{X: 2, Y: 1, Cell: buffer.NewCell("z")}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("draw diffs = %#v, want %#v", got, want)
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

func TestTerminal_TryDraw_shouldReturnCallbackErrorWithoutMutatingTerminal(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	renderText(t, term, "abc")
	backend.draws = nil
	backend.flushCount = 0
	backend.hideCursorCount = 0
	backend.showCursorCount = 0
	backend.cursorPositions = nil

	callbackErr := errors.New("render failed")
	completed, err := term.TryDraw(func(frame *terminal.Frame) error {
		frame.RenderWidget(widgets.NewParagraph(text.FromString("xyz")), frame.Area())
		frame.SetCursorPosition(layout.Position{X: 1, Y: 0})
		return callbackErr
	})

	if !errors.Is(err, callbackErr) {
		t.Fatalf("TryDraw error = %v, want %v", err, callbackErr)
	}
	if completed != nil {
		t.Fatalf("completed = %#v, want nil", completed)
	}
	if got := len(backend.draws); got != 0 {
		t.Fatalf("draw count = %d, want 0", got)
	}
	if got := backend.flushCount; got != 0 {
		t.Fatalf("flush count = %d, want 0", got)
	}
	if got := backend.hideCursorCount + backend.showCursorCount + len(backend.cursorPositions); got != 0 {
		t.Fatalf("cursor backend calls = %d, want 0", got)
	}

	renderText(t, term, "abc")
	if got, want := backend.draws[0], []buffer.CellDiff{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("next draw diffs = %#v, want %#v", got, want)
	}
}

func TestTerminal_Draw_shouldCallAutoresizeBeforeRendering(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	backend.SetSize(4, 2)

	completed, err := term.Draw(func(frame *terminal.Frame) {
		if got, want := frame.Area(), layout.NewRect(0, 0, 4, 2); got != want {
			t.Fatalf("frame area = %#v, want %#v", got, want)
		}
		frame.RenderWidget(widgets.NewParagraph(text.FromString("abcd\nxy")), frame.Area())
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	if got, want := completed.Area, layout.NewRect(0, 0, 4, 2); got != want {
		t.Fatalf("completed area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Autoresize_shouldResizeWhenBackendSizeChanges(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	backend.SetSize(5, 2)

	if err := term.Autoresize(); err != nil {
		t.Fatalf("Autoresize returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 0, 5, 2); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	completed, err := term.Draw(nil)
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
	if got, want := completed.Buffer.Area, layout.NewRect(0, 0, 5, 2); got != want {
		t.Fatalf("buffer area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Autoresize_shouldNoopWhenSizeUnchanged(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := term.Autoresize(); err != nil {
		t.Fatalf("Autoresize returned error: %v", err)
	}
	if got, want := term.Area(), layout.NewRect(0, 0, 3, 1); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Autoresize_fullscreenTracksBackendSize(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FullscreenViewport(),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	backend.SetSize(5, 2)

	if err := term.Autoresize(); err != nil {
		t.Fatalf("Autoresize returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 0, 5, 2); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	completed, err := term.Draw(nil)
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
	if got, want := completed.Buffer.Area, layout.NewRect(0, 0, 5, 2); got != want {
		t.Fatalf("buffer area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Autoresize_inlineTracksBackendSize(t *testing.T) {
	backend := newRecordingBackend(10, 10)
	backend.cursorPosition = layout.Position{X: 0, Y: 4}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	backend.SetSize(12, 8)
	backend.cursorPosition = layout.Position{X: 0, Y: 5}

	if err := term.Autoresize(); err != nil {
		t.Fatalf("Autoresize returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 4, 12, 4); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Autoresize_fixedViewportNoop(t *testing.T) {
	backend := newRecordingBackend(10, 5)
	area := layout.NewRect(1, 1, 3, 2)
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(area),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	backend.SetSize(20, 10)

	if err := term.Autoresize(); err != nil {
		t.Fatalf("Autoresize returned error: %v", err)
	}

	if got, want := term.Area(), area; got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	completed, err := term.Draw(nil)
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
	if got, want := completed.Buffer.Area, area; got != want {
		t.Fatalf("buffer area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Flush_shouldDrawCurrentDiffOnly(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	frame := term.Frame()
	frame.Buffer().SetSymbol(1, 0, "x")
	if err := term.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}

	if got, want := backend.draws, [][]buffer.CellDiff{{{X: 1, Y: 0, Cell: buffer.NewCell("x")}}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("draws = %#v, want %#v", got, want)
	}
	if got := backend.flushCount; got != 0 {
		t.Fatalf("backend flush count = %d, want 0", got)
	}
}

func TestTerminal_InsertBefore_nonInlineViewportNoop(t *testing.T) {
	backend := testbackend.New(3, 2)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	frame := term.Frame()
	frame.Buffer().SetSymbol(0, 0, "x")
	if err := term.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}

	area := term.Area()
	err = term.InsertBefore(1, func(buf *buffer.Buffer) {
		setLine(buf, 0, "zzz")
	})
	if err != nil {
		t.Fatalf("InsertBefore returned error: %v", err)
	}

	if got, want := term.Area(), area; got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	if got, want := backend.Lines(), []string{"x  ", "   "}; !reflect.DeepEqual(got, want) {
		t.Fatalf("backend lines = %#v, want %#v", got, want)
	}
}

func TestTerminal_InsertBefore_inlinePushesViewportDownWhenSpaceAvailable(t *testing.T) {
	backend := testbackend.WithLines([]string{
		"0000000000",
		"1111111111",
		"2222222222",
		"3333333333",
		"4444444444",
		"5555555555",
		"6666666666",
		"7777777777",
		"8888888888",
		"9999999999",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 3}); err != nil {
		t.Fatalf("SetCursorPosition returned error: %v", err)
	}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	err = term.InsertBefore(1, func(buf *buffer.Buffer) {
		setLine(buf, 0, "INSERTLINE")
	})
	if err != nil {
		t.Fatalf("InsertBefore returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 4, 10, 4); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	wantLines := []string{
		"0000000000",
		"1111111111",
		"2222222222",
		"INSERTLINE",
		"          ",
		"          ",
		"          ",
		"          ",
		"          ",
		"          ",
	}
	if got := backend.Lines(); !reflect.DeepEqual(got, wantLines) {
		t.Fatalf("backend lines = %#v, want %#v", got, wantLines)
	}
}

func TestTerminal_InsertBefore_inlineScrollsWhenViewportIsAtBottom(t *testing.T) {
	backend := testbackend.WithLines([]string{
		"0000000000",
		"1111111111",
		"2222222222",
		"3333333333",
		"4444444444",
		"5555555555",
		"6666666666",
		"7777777777",
		"8888888888",
		"9999999999",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 6}); err != nil {
		t.Fatalf("SetCursorPosition returned error: %v", err)
	}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	err = term.InsertBefore(2, func(buf *buffer.Buffer) {
		setLine(buf, 0, "INSERTED1")
		setLine(buf, 1, "INSERTED2")
	})
	if err != nil {
		t.Fatalf("InsertBefore returned error: %v", err)
	}

	if got, want := term.Area(), layout.NewRect(0, 6, 10, 4); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	wantLines := []string{
		"2222222222",
		"3333333333",
		"4444444444",
		"5555555555",
		"INSERTED1 ",
		"INSERTED2 ",
		"          ",
		"          ",
		"          ",
		"          ",
	}
	if got := backend.Lines(); !reflect.DeepEqual(got, wantLines) {
		t.Fatalf("backend lines = %#v, want %#v", got, wantLines)
	}
}

func TestTerminal_InsertBefore_thenDrawRepaintsClearedViewport(t *testing.T) {
	backend := testbackend.New(10, 10)
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 6}); err != nil {
		t.Fatalf("SetCursorPosition returned error: %v", err)
	}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	drawRows(t, term, "AAAAAAAAAA")

	err = term.InsertBefore(2, func(buf *buffer.Buffer) {
		setLine(buf, 0, "INSERTED00")
		setLine(buf, 1, "INSERTED01")
	})
	if err != nil {
		t.Fatalf("InsertBefore returned error: %v", err)
	}
	drawRows(t, term, "BBBBBBBBBB")

	wantLines := []string{
		"          ",
		"          ",
		"          ",
		"          ",
		"INSERTED00",
		"INSERTED01",
		"BBBBBBBBBB",
		"BBBBBBBBBB",
		"BBBBBBBBBB",
		"BBBBBBBBBB",
	}
	if got := backend.Lines(); !reflect.DeepEqual(got, wantLines) {
		t.Fatalf("backend lines = %#v, want %#v", got, wantLines)
	}
}

func TestTerminal_SwapBuffers_shouldPrepareNextFrame(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	frame := term.Frame()
	frame.Buffer().SetSymbol(0, 0, "a")
	if err := term.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}
	term.SwapBuffers()
	backend.draws = nil

	renderText(t, term, "a")

	if got, want := backend.draws[0], []buffer.CellDiff{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("next draw diffs = %#v, want %#v", got, want)
	}
}

func TestTerminal_DirectCursorMethods_shouldCallBackend(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := term.HideCursor(); err != nil {
		t.Fatalf("HideCursor returned error: %v", err)
	}
	if err := term.ShowCursor(); err != nil {
		t.Fatalf("ShowCursor returned error: %v", err)
	}
	if err := term.SetCursorPosition(layout.Position{X: 2, Y: 0}); err != nil {
		t.Fatalf("SetCursorPosition returned error: %v", err)
	}

	if got, want := backend.hideCursorCount, 1; got != want {
		t.Fatalf("hide cursor count = %d, want %d", got, want)
	}
	if got, want := backend.showCursorCount, 1; got != want {
		t.Fatalf("show cursor count = %d, want %d", got, want)
	}
	if got, want := backend.cursorPositions, []layout.Position{{X: 2, Y: 0}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cursor positions = %#v, want %#v", got, want)
	}
}

func TestTerminal_Backend_returnsSharedBackendReference(t *testing.T) {
	backend := testbackend.New(3, 2)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if got, want := term.Backend(), terminal.Backend(backend); got != want {
		t.Fatalf("backend reference = %#v, want %#v", got, want)
	}
}

func TestTerminal_Size_queriesUnderlyingBackendSize(t *testing.T) {
	backend := newRecordingBackend(3, 2)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	backend.SetSize(4, 3)

	got, err := term.Size()
	if err != nil {
		t.Fatalf("Size returned error: %v", err)
	}

	if want := (layout.Size{Width: 4, Height: 3}); got != want {
		t.Fatalf("size = %#v, want %#v", got, want)
	}
	if got, want := term.Area(), layout.NewRect(0, 0, 3, 2); got != want {
		t.Fatalf("terminal area = %#v, want unchanged %#v", got, want)
	}
}

func TestTerminal_GetCursorPosition_queriesBackend(t *testing.T) {
	backend := newRecordingBackend(10, 5)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	backend.cursorPosition = layout.Position{X: 7, Y: 2}

	got, err := term.GetCursorPosition()
	if err != nil {
		t.Fatalf("GetCursorPosition returned error: %v", err)
	}

	if want := (layout.Position{X: 7, Y: 2}); got != want {
		t.Fatalf("cursor position = %#v, want %#v", got, want)
	}
}

func TestTerminal_HideCursor_updatesTerminalState(t *testing.T) {
	backend := testbackend.New(10, 5)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := term.HideCursor(); err != nil {
		t.Fatalf("HideCursor returned error: %v", err)
	}

	if backend.CursorVisible() {
		t.Fatalf("cursor visible = true, want false")
	}
}

func TestTerminal_ShowCursor_updatesTerminalState(t *testing.T) {
	backend := testbackend.New(10, 5)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := term.HideCursor(); err != nil {
		t.Fatalf("HideCursor returned error: %v", err)
	}
	if err := term.ShowCursor(); err != nil {
		t.Fatalf("ShowCursor returned error: %v", err)
	}

	if !backend.CursorVisible() {
		t.Fatalf("cursor visible = false, want true")
	}
}

func TestTerminal_SetCursorPosition_updatesBackendAndTracking(t *testing.T) {
	backend := testbackend.New(10, 5)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := term.SetCursorPosition(layout.Position{X: 3, Y: 4}); err != nil {
		t.Fatalf("SetCursorPosition returned error: %v", err)
	}

	if got, want := backend.CursorPosition(), (layout.Position{X: 3, Y: 4}); got != want {
		t.Fatalf("cursor position = %#v, want %#v", got, want)
	}
}

func TestTerminal_Resize_inlineUsesDirectCursorTracking(t *testing.T) {
	backend := newRecordingBackend(10, 10)
	backend.cursorPosition = layout.Position{X: 0, Y: 4}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	if err := term.SetCursorPosition(layout.Position{X: 0, Y: 5}); err != nil {
		t.Fatalf("SetCursorPosition returned error: %v", err)
	}
	backend.cursorPosition = layout.Position{X: 0, Y: 6}

	term.Resize(layout.NewRect(0, 0, 10, 12))

	if got, want := term.Area(), layout.NewRect(0, 5, 10, 4); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Draw_shouldUseTryDrawSuccessOrder(t *testing.T) {
	backend := newRecordingBackend(3, 1)
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

	if got, want := backend.operations, []string{"draw", "show-cursor", "set-cursor", "backend-flush"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("operations = %#v, want %#v", got, want)
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

	backend.SetSize(4, 2)
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

func TestTerminal_Resize_fixedViewportChangesAreaAndBuffers(t *testing.T) {
	backend := newRecordingBackend(5, 3)
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(layout.NewRect(1, 1, 2, 1)),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	area := layout.NewRect(0, 0, 3, 2)

	term.Resize(area)

	if got, want := term.Area(), area; got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
	completed, err := term.Draw(func(frame *terminal.Frame) {
		if got, want := frame.Area(), area; got != want {
			t.Fatalf("frame area = %#v, want %#v", got, want)
		}
		if got, want := frame.Buffer().Area, area; got != want {
			t.Fatalf("frame buffer area = %#v, want %#v", got, want)
		}
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
	if got, want := completed.Buffer.Area, area; got != want {
		t.Fatalf("completed buffer area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Resize_inlineRecomputesOriginUsingPreviousCursorOffset(t *testing.T) {
	backend := newRecordingBackend(10, 10)
	backend.cursorPosition = layout.Position{X: 0, Y: 4}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	if got, want := term.Area(), layout.NewRect(0, 4, 10, 4); got != want {
		t.Fatalf("initial terminal area = %#v, want %#v", got, want)
	}

	if _, err := term.Draw(func(frame *terminal.Frame) {
		frame.SetCursorPosition(layout.Position{X: 0, Y: 5})
	}); err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
	backend.cursorPosition = layout.Position{X: 0, Y: 6}
	backend.SetSize(10, 12)

	term.Resize(layout.NewRect(0, 0, 10, 12))

	if got, want := term.Area(), layout.NewRect(0, 5, 10, 4); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Resize_inlineClampsHeightToTerminalHeight(t *testing.T) {
	backend := newRecordingBackend(10, 10)
	backend.cursorPosition = layout.Position{X: 0, Y: 0}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(10),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	backend.SetSize(10, 3)

	term.Resize(layout.NewRect(0, 0, 10, 3))

	if got, want := term.Area(), layout.NewRect(0, 0, 10, 3); got != want {
		t.Fatalf("terminal area = %#v, want %#v", got, want)
	}
}

func TestTerminal_Resize_inlinePreservesBackendCursorAcrossRepeatedResizes(t *testing.T) {
	backend := newRecordingBackend(10, 10)
	backend.cursorPosition = layout.Position{X: 0, Y: 4}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.InlineViewport(4),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}
	if _, err := term.Draw(func(frame *terminal.Frame) {
		frame.SetCursorPosition(layout.Position{X: 0, Y: 5})
	}); err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
	backend.cursorPosition = layout.Position{X: 0, Y: 6}

	term.Resize(layout.NewRect(0, 0, 10, 12))
	if got, want := term.Area(), layout.NewRect(0, 5, 10, 4); got != want {
		t.Fatalf("first resize area = %#v, want %#v", got, want)
	}
	if got, want := backend.cursorPosition, (layout.Position{X: 0, Y: 6}); got != want {
		t.Fatalf("cursor after first resize = %#v, want %#v", got, want)
	}

	term.Resize(layout.NewRect(0, 0, 10, 14))
	if got, want := term.Area(), layout.NewRect(0, 6, 10, 4); got != want {
		t.Fatalf("second resize area = %#v, want %#v", got, want)
	}
	if got, want := backend.cursorPosition, (layout.Position{X: 0, Y: 6}); got != want {
		t.Fatalf("cursor after second resize = %#v, want %#v", got, want)
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

	if got, want := backend.clearRegions, []terminal.ClearType{terminal.ClearAll}; !reflect.DeepEqual(got, want) {
		t.Fatalf("clear regions = %#v, want %#v", got, want)
	}
	if got, want := diffSymbols(backend.draws[1]), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("second draw diff symbols = %#v, want %#v", got, want)
	}
}

func TestTerminal_Clear_fullscreenClearsBackendAndResetsBackBuffer(t *testing.T) {
	backend := newRecordingBackend(3, 2)
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	renderText(t, term, "abc\ndef")
	backend.draws = nil

	if err := term.Clear(); err != nil {
		t.Fatalf("Clear returned error: %v", err)
	}
	renderText(t, term, "abc\ndef")

	if got, want := backend.clearRegions, []terminal.ClearType{terminal.ClearAll}; !reflect.DeepEqual(got, want) {
		t.Fatalf("clear regions = %#v, want %#v", got, want)
	}
	if got, want := len(backend.draws[0]), 6; got != want {
		t.Fatalf("redraw diff count after Clear = %d, want %d", got, want)
	}
}

func TestTerminal_Clear_fixedFullWidthAtBottomClearsAfterViewportOrigin(t *testing.T) {
	backend := newRecordingBackend(10, 3)
	backend.cursorPosition = layout.Position{X: 2, Y: 0}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(layout.NewRect(0, 1, 10, 2)),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	if err := term.Clear(); err != nil {
		t.Fatalf("Clear returned error: %v", err)
	}

	if got, want := backend.clearRegions, []terminal.ClearType{terminal.ClearAfterCursor}; !reflect.DeepEqual(got, want) {
		t.Fatalf("clear regions = %#v, want %#v", got, want)
	}
	if got, want := backend.cursorPositions, []layout.Position{{X: 0, Y: 1}, {X: 2, Y: 0}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cursor positions = %#v, want %#v", got, want)
	}
}

func TestTerminal_Clear_fixedFullWidthNotAtBottomClearsViewportRowsOnly(t *testing.T) {
	backend := newRecordingBackend(10, 4)
	backend.cursorPosition = layout.Position{X: 1, Y: 0}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(layout.NewRect(0, 1, 10, 2)),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	if err := term.Clear(); err != nil {
		t.Fatalf("Clear returned error: %v", err)
	}

	if got, want := backend.clearRegions, []terminal.ClearType{terminal.ClearCurrentLine, terminal.ClearCurrentLine}; !reflect.DeepEqual(got, want) {
		t.Fatalf("clear regions = %#v, want %#v", got, want)
	}
	if got, want := backend.cursorPositions, []layout.Position{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 1, Y: 0}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cursor positions = %#v, want %#v", got, want)
	}
}

func TestTerminal_Clear_fixedNonFullWidthClearsViewportCellsOnly(t *testing.T) {
	backend := newRecordingBackend(10, 4)
	backend.cursorPosition = layout.Position{X: 3, Y: 0}
	term, err := terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(layout.NewRect(1, 1, 3, 2)),
	})
	if err != nil {
		t.Fatalf("NewWithOptions returned error: %v", err)
	}

	if err := term.Clear(); err != nil {
		t.Fatalf("Clear returned error: %v", err)
	}

	wantDraw := []buffer.CellDiff{
		{X: 1, Y: 1, Cell: buffer.NewCell(" ")},
		{X: 2, Y: 1, Cell: buffer.NewCell(" ")},
		{X: 3, Y: 1, Cell: buffer.NewCell(" ")},
		{X: 1, Y: 2, Cell: buffer.NewCell(" ")},
		{X: 2, Y: 2, Cell: buffer.NewCell(" ")},
		{X: 3, Y: 2, Cell: buffer.NewCell(" ")},
	}
	if got := backend.draws; !reflect.DeepEqual(got, [][]buffer.CellDiff{wantDraw}) {
		t.Fatalf("draws = %#v, want %#v", got, [][]buffer.CellDiff{wantDraw})
	}
	if got, want := backend.cursorPositions, []layout.Position{{X: 3, Y: 0}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cursor positions = %#v, want %#v", got, want)
	}
}

func TestTerminal_Clear_preservesBackendCursorPosition(t *testing.T) {
	backend := newRecordingBackend(3, 1)
	backend.cursorPosition = layout.Position{X: 2, Y: 0}
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := term.Clear(); err != nil {
		t.Fatalf("Clear returned error: %v", err)
	}

	if got, want := backend.cursorPosition, (layout.Position{X: 2, Y: 0}); got != want {
		t.Fatalf("backend cursor position = %#v, want %#v", got, want)
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

func drawRows(t *testing.T, term *terminal.Terminal, content string) {
	t.Helper()
	_, err := term.Draw(func(frame *terminal.Frame) {
		area := frame.Area()
		for y := area.Y; y < area.Bottom(); y++ {
			setLine(frame.Buffer(), y, content)
		}
	})
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}
}

func setLine(buf *buffer.Buffer, y int, content string) {
	for x, r := range content {
		buf.SetSymbol(x, y, string(r))
	}
}

func diffSymbols(diffs []buffer.CellDiff) []string {
	symbols := make([]string, 0, len(diffs))
	for _, diff := range diffs {
		symbols = append(symbols, diff.Cell.Symbol)
	}
	return symbols
}
