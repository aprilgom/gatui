package tcell

import (
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/terminal"

	tcelllib "github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type spyScreen struct {
	initCount       int
	finiCount       int
	clearCount      int
	showCount       int
	hideCursorCount int
	width           int
	height          int
	events          chan tcelllib.Event
	showCursorCalls []layout.Position
	contentCalls    []contentCall
}

type contentCall struct {
	x         int
	y         int
	primary   rune
	combining []rune
	style     tcelllib.Style
}

func newSpyScreen(width, height int) *spyScreen {
	return &spyScreen{width: width, height: height, events: make(chan tcelllib.Event, 1)}
}

func (s *spyScreen) Init() error {
	s.initCount++
	return nil
}

func (s *spyScreen) Fini() {
	s.finiCount++
	close(s.events)
}

func (s *spyScreen) Clear() {
	s.clearCount++
}

func (s *spyScreen) Show() {
	s.showCount++
}

func (s *spyScreen) HideCursor() {
	s.hideCursorCount++
}

func (s *spyScreen) ShowCursor(x, y int) {
	s.showCursorCalls = append(s.showCursorCalls, layout.Position{X: x, Y: y})
}

func (s *spyScreen) SetContent(x int, y int, primary rune, combining []rune, cellStyle tcelllib.Style) {
	copiedCombining := append([]rune(nil), combining...)
	s.contentCalls = append(s.contentCalls, contentCall{
		x:         x,
		y:         y,
		primary:   primary,
		combining: copiedCombining,
		style:     cellStyle,
	})
}

func (s *spyScreen) Fill(rune, tcelllib.Style) {}

func (s *spyScreen) Put(x int, y int, str string, style tcelllib.Style) (string, int) {
	return "", 0
}

func (s *spyScreen) PutStr(x int, y int, str string) {}

func (s *spyScreen) PutStrStyled(x int, y int, str string, style tcelllib.Style) {}

func (s *spyScreen) Get(x, y int) (string, tcelllib.Style, int) {
	return "", tcelllib.StyleDefault, 0
}

func (s *spyScreen) SetStyle(style tcelllib.Style) {}

func (s *spyScreen) SetCursorStyle(tcelllib.CursorStyle, ...color.Color) {}

func (s *spyScreen) Size() (int, int) {
	return s.width, s.height
}

func (s *spyScreen) EventQ() chan tcelllib.Event {
	return s.events
}

func (s *spyScreen) EnableMouse(...tcelllib.MouseFlags) {}

func (s *spyScreen) DisableMouse() {}

func (s *spyScreen) EnablePaste() {}

func (s *spyScreen) DisablePaste() {}

func (s *spyScreen) EnableFocus() {}

func (s *spyScreen) DisableFocus() {}

func (s *spyScreen) Colors() int {
	return 256
}

func (s *spyScreen) Sync() {}

func (s *spyScreen) CharacterSet() string {
	return "UTF-8"
}

func (s *spyScreen) RegisterRuneFallback(r rune, subst string) {}

func (s *spyScreen) UnregisterRuneFallback(r rune) {}

func (s *spyScreen) Resize(int, int, int, int) {}

func (s *spyScreen) Suspend() error {
	return nil
}

func (s *spyScreen) Resume() error {
	return nil
}

func (s *spyScreen) Beep() error {
	return nil
}

func (s *spyScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *spyScreen) LockRegion(x, y, width, height int, lock bool) {}

func (s *spyScreen) Tty() (tcelllib.Tty, bool) {
	return nil, false
}

func (s *spyScreen) SetTitle(string) {}

func (s *spyScreen) SetClipboard([]byte) {}

func (s *spyScreen) GetClipboard() {}

func (s *spyScreen) HasClipboard() bool {
	return false
}

func (s *spyScreen) ShowNotification(title string, body string) {}

func (s *spyScreen) KeyboardProtocol() tcelllib.KeyProtocol {
	return tcelllib.LegacyKeyboard
}

func (s *spyScreen) Terminal() (string, string) {
	return "", ""
}

func TestBackend_NewWithScreen_shouldInitInjectedScreen(t *testing.T) {
	screen := newSpyScreen(80, 24)

	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	if screen.initCount != 1 {
		t.Fatalf("Init count = %d, want 1", screen.initCount)
	}
}

func TestBackend_Size_shouldReturnScreenSize(t *testing.T) {
	screen := newSpyScreen(80, 24)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()
	screen.SetSize(80, 24)

	got, err := backend.Size()
	if err != nil {
		t.Fatalf("Size() error = %v", err)
	}

	want := layout.Size{Width: 80, Height: 24}
	if got != want {
		t.Fatalf("Size() = %+v, want %+v", got, want)
	}
}

func TestTcellBackend_WindowSize_shouldReturnScreenSizeAndZeroPixels(t *testing.T) {
	screen := newSpyScreen(80, 24)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()
	screen.SetSize(80, 24)

	got, err := backend.WindowSize()
	if err != nil {
		t.Fatalf("WindowSize() error = %v", err)
	}

	want := terminal.WindowSize{
		ColumnsRows: layout.Size{Width: 80, Height: 24},
		Pixels:      layout.Size{Width: 0, Height: 0},
	}
	if got != want {
		t.Fatalf("WindowSize() = %+v, want %+v", got, want)
	}
}

func TestBackend_shouldSatisfyTerminalInterface(t *testing.T) {
	var _ terminal.Backend = (*Backend)(nil)
}

func TestBackend_AppendLines_shouldMoveCursorToLastRowAndBlankScreen(t *testing.T) {
	screen := newSpyScreen(4, 3)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()
	screen.SetSize(4, 3)

	if err := backend.AppendLines(2); err != nil {
		t.Fatalf("AppendLines() error = %v", err)
	}

	if got, want := backend.cursorPosition, (layout.Position{X: 0, Y: 2}); got != want {
		t.Fatalf("cursor position = %#v, want %#v", got, want)
	}
	if got, want := screen.showCursorCalls[len(screen.showCursorCalls)-1], (layout.Position{X: 0, Y: 2}); got != want {
		t.Fatalf("show cursor call = %#v, want %#v", got, want)
	}
	if got, want := len(screen.contentCalls), 12; got != want {
		t.Fatalf("blank content calls = %d, want %d", got, want)
	}
}

func TestBackend_AppendLines_zeroCountNoop(t *testing.T) {
	screen := newSpyScreen(4, 3)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	if err := backend.AppendLines(0); err != nil {
		t.Fatalf("AppendLines() error = %v", err)
	}

	if got := len(screen.contentCalls); got != 0 {
		t.Fatalf("blank content calls = %d, want 0", got)
	}
}

func TestBackend_Draw_shouldSetContentForAsciiCell(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	err = backend.Draw([]buffer.CellDiff{{
		X:    1,
		Y:    2,
		Cell: buffer.NewCell("A"),
	}})
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	if len(screen.contentCalls) != 1 {
		t.Fatalf("SetContent calls = %d, want 1", len(screen.contentCalls))
	}
	got := screen.contentCalls[0]
	if got.x != 1 || got.y != 2 {
		t.Fatalf("SetContent position = (%d, %d), want (1, 2)", got.x, got.y)
	}
	if got.primary != 'A' {
		t.Fatalf("primary rune = %q, want %q", got.primary, 'A')
	}
	if len(got.combining) != 0 {
		t.Fatalf("combining runes = %q, want none", string(got.combining))
	}
}

func TestBackend_Draw_shouldSetContentForGraphemeCluster(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	err = backend.Draw([]buffer.CellDiff{{
		X:    1,
		Y:    2,
		Cell: buffer.NewCell("o\u0301"),
	}})
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	if len(screen.contentCalls) != 1 {
		t.Fatalf("SetContent calls = %d, want 1", len(screen.contentCalls))
	}
	got := screen.contentCalls[0]
	if got.primary != 'o' {
		t.Fatalf("primary rune = %q, want %q", got.primary, 'o')
	}
	if got, want := string(got.combining), "\u0301"; got != want {
		t.Fatalf("combining runes = %q, want %q", got, want)
	}
}

func TestBackend_Draw_shouldUseDisplaySymbolForEmptyCell(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	err = backend.Draw([]buffer.CellDiff{{
		X:    1,
		Y:    2,
		Cell: buffer.NewCell(""),
	}})
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	if len(screen.contentCalls) != 1 {
		t.Fatalf("SetContent calls = %d, want 1", len(screen.contentCalls))
	}
	got := screen.contentCalls[0]
	if got.primary != ' ' {
		t.Fatalf("primary rune = %q, want space", got.primary)
	}
}

func TestBackend_Draw_shouldMapStyle(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	cellStyle := style.NewStyle().
		Fg(style.Red).
		Bg(style.Blue).
		AddModifier(style.ModifierBold | style.ModifierItalic | style.ModifierUnderlined | style.ModifierReversed)
	err = backend.Draw([]buffer.CellDiff{{
		X:    1,
		Y:    2,
		Cell: buffer.Cell{Symbol: "A", Style: cellStyle},
	}})
	if err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	if len(screen.contentCalls) != 1 {
		t.Fatalf("SetContent calls = %d, want 1", len(screen.contentCalls))
	}
	gotStyle := screen.contentCalls[0].style
	if fg := gotStyle.GetForeground(); fg != color.Red {
		t.Fatalf("foreground = %v, want %v", fg, color.Red)
	}
	if bg := gotStyle.GetBackground(); bg != color.Blue {
		t.Fatalf("background = %v, want %v", bg, color.Blue)
	}
	if !gotStyle.HasBold() {
		t.Fatal("style should include bold")
	}
	if !gotStyle.HasItalic() {
		t.Fatal("style should include italic")
	}
	if !gotStyle.HasUnderline() {
		t.Fatal("style should include underline")
	}
	if !gotStyle.HasReverse() {
		t.Fatal("style should include reverse")
	}
}

func TestBackend_ClearAndFlush_shouldCallScreenClearAndShow(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	if err := backend.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}
	if err := backend.Flush(); err != nil {
		t.Fatalf("Flush() error = %v", err)
	}

	if screen.clearCount != 1 {
		t.Fatalf("Clear count = %d, want 1", screen.clearCount)
	}
	if screen.showCount != 1 {
		t.Fatalf("Show count = %d, want 1", screen.showCount)
	}
}

func TestBackend_ClearRegion_afterCursorBlanksFromCursorToScreenEnd(t *testing.T) {
	screen := newSpyScreen(4, 3)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()
	screen.SetSize(4, 3)
	if err := backend.SetCursorPosition(layout.Position{X: 2, Y: 1}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}
	screen.contentCalls = nil

	if err := backend.ClearRegion(terminal.ClearAfterCursor); err != nil {
		t.Fatalf("ClearRegion(ClearAfterCursor) error = %v", err)
	}

	wantPositions := []layout.Position{{X: 2, Y: 1}, {X: 3, Y: 1}, {X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}}
	if got := contentCallPositions(screen.contentCalls); !positionsEqual(got, wantPositions) {
		t.Fatalf("SetContent positions = %+v, want %+v", got, wantPositions)
	}
	for _, call := range screen.contentCalls {
		if call.primary != ' ' {
			t.Fatalf("SetContent primary = %q, want space", call.primary)
		}
	}
}

func TestBackend_ClearRegion_beforeCursorBlanksFromScreenStartThroughCursor(t *testing.T) {
	screen := newSpyScreen(4, 3)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()
	screen.SetSize(4, 3)
	if err := backend.SetCursorPosition(layout.Position{X: 2, Y: 1}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}
	screen.contentCalls = nil

	if err := backend.ClearRegion(terminal.ClearBeforeCursor); err != nil {
		t.Fatalf("ClearRegion(ClearBeforeCursor) error = %v", err)
	}

	wantPositions := []layout.Position{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}}
	if got := contentCallPositions(screen.contentCalls); !positionsEqual(got, wantPositions) {
		t.Fatalf("SetContent positions = %+v, want %+v", got, wantPositions)
	}
	for _, call := range screen.contentCalls {
		if call.primary != ' ' {
			t.Fatalf("SetContent primary = %q, want space", call.primary)
		}
	}
}

func TestBackend_ClearRegion_currentLineBlanksCursorRow(t *testing.T) {
	screen := newSpyScreen(4, 3)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()
	screen.SetSize(4, 3)
	if err := backend.SetCursorPosition(layout.Position{X: 2, Y: 1}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}
	screen.contentCalls = nil

	if err := backend.ClearRegion(terminal.ClearCurrentLine); err != nil {
		t.Fatalf("ClearRegion(ClearCurrentLine) error = %v", err)
	}

	wantPositions := []layout.Position{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}}
	if got := contentCallPositions(screen.contentCalls); !positionsEqual(got, wantPositions) {
		t.Fatalf("SetContent positions = %+v, want %+v", got, wantPositions)
	}
}

func TestBackend_ClearRegion_untilNewLineBlanksFromCursorToLineEnd(t *testing.T) {
	screen := newSpyScreen(4, 3)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()
	screen.SetSize(4, 3)
	if err := backend.SetCursorPosition(layout.Position{X: 2, Y: 1}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}
	screen.contentCalls = nil

	if err := backend.ClearRegion(terminal.ClearUntilNewLine); err != nil {
		t.Fatalf("ClearRegion(ClearUntilNewLine) error = %v", err)
	}

	wantPositions := []layout.Position{{X: 2, Y: 1}, {X: 3, Y: 1}}
	if got := contentCallPositions(screen.contentCalls); !positionsEqual(got, wantPositions) {
		t.Fatalf("SetContent positions = %+v, want %+v", got, wantPositions)
	}
	for _, call := range screen.contentCalls {
		if call.primary != ' ' {
			t.Fatalf("SetContent primary = %q, want space", call.primary)
		}
	}
}

func TestBackend_CursorMethods_shouldHideShowAndPositionCursor(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	if err := backend.HideCursor(); err != nil {
		t.Fatalf("HideCursor() error = %v", err)
	}
	if err := backend.ShowCursor(); err != nil {
		t.Fatalf("ShowCursor() error = %v", err)
	}
	if err := backend.SetCursorPosition(layout.Position{X: 3, Y: 4}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if screen.hideCursorCount != 1 {
		t.Fatalf("HideCursor count = %d, want 1", screen.hideCursorCount)
	}
	wantCalls := []layout.Position{{X: 0, Y: 0}, {X: 3, Y: 4}}
	if len(screen.showCursorCalls) != len(wantCalls) {
		t.Fatalf("ShowCursor calls = %+v, want %+v", screen.showCursorCalls, wantCalls)
	}
	for i := range wantCalls {
		if screen.showCursorCalls[i] != wantCalls[i] {
			t.Fatalf("ShowCursor calls = %+v, want %+v", screen.showCursorCalls, wantCalls)
		}
	}
}

func contentCallPositions(calls []contentCall) []layout.Position {
	positions := make([]layout.Position, 0, len(calls))
	for _, call := range calls {
		positions = append(positions, layout.Position{X: call.x, Y: call.y})
	}
	return positions
}

func positionsEqual(a, b []layout.Position) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBackend_Close_shouldFinalizeScreen(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}

	backend.Close()

	if screen.finiCount != 1 {
		t.Fatalf("Fini count = %d, want 1", screen.finiCount)
	}
}

func TestTerminalBackendInterface_shouldAcceptTcellBackend(t *testing.T) {
	var _ terminal.Backend = (*Backend)(nil)
}
