package tcell

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/terminal"

	tcelllib "github.com/gdamore/tcell/v2"
)

type spyScreen struct {
	tcelllib.SimulationScreen

	initCount       int
	finiCount       int
	clearCount      int
	showCount       int
	hideCursorCount int
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
	screen := tcelllib.NewSimulationScreen("UTF-8")
	screen.SetSize(width, height)
	return &spyScreen{SimulationScreen: screen}
}

func (s *spyScreen) Init() error {
	s.initCount++
	return s.SimulationScreen.Init()
}

func (s *spyScreen) Fini() {
	s.finiCount++
	s.SimulationScreen.Fini()
}

func (s *spyScreen) Clear() {
	s.clearCount++
	s.SimulationScreen.Clear()
}

func (s *spyScreen) Show() {
	s.showCount++
	s.SimulationScreen.Show()
}

func (s *spyScreen) HideCursor() {
	s.hideCursorCount++
	s.SimulationScreen.HideCursor()
}

func (s *spyScreen) ShowCursor(x, y int) {
	s.showCursorCalls = append(s.showCursorCalls, layout.Position{X: x, Y: y})
	s.SimulationScreen.ShowCursor(x, y)
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
	s.SimulationScreen.SetContent(x, y, primary, combining, cellStyle)
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
	fg, bg, attrs := gotStyle.Decompose()
	if fg != tcelllib.ColorRed {
		t.Fatalf("foreground = %v, want %v", fg, tcelllib.ColorRed)
	}
	if bg != tcelllib.ColorBlue {
		t.Fatalf("background = %v, want %v", bg, tcelllib.ColorBlue)
	}
	wantAttrs := tcelllib.AttrBold | tcelllib.AttrItalic | tcelllib.AttrUnderline | tcelllib.AttrReverse
	if attrs&wantAttrs != wantAttrs {
		t.Fatalf("attrs = %v, want to include %v", attrs, wantAttrs)
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

func TestBackend_PollEvent_shouldConvertRuneKey(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	screen.InjectKey(tcelllib.KeyRune, 'a', tcelllib.ModNone)
	event, err := backend.PollEvent()
	if err != nil {
		t.Fatalf("PollEvent() error = %v", err)
	}

	key, ok := event.(terminal.KeyEvent)
	if !ok {
		t.Fatalf("PollEvent() = %T, want terminal.KeyEvent", event)
	}
	if key.Type() != terminal.EventKey {
		t.Fatalf("Type() = %v, want %v", key.Type(), terminal.EventKey)
	}
	if key.Code != terminal.KeyRune || key.Rune != 'a' {
		t.Fatalf("key = %+v, want Code KeyRune and Rune 'a'", key)
	}
}

func TestBackend_PollEvent_shouldConvertSpecialKeys(t *testing.T) {
	cases := []struct {
		name string
		key  tcelllib.Key
		want terminal.KeyCode
	}{
		{"Enter", tcelllib.KeyEnter, terminal.KeyEnter},
		{"Esc", tcelllib.KeyEsc, terminal.KeyEsc},
		{"Backspace", tcelllib.KeyBackspace, terminal.KeyBackspace},
		{"Backspace2", tcelllib.KeyBackspace2, terminal.KeyBackspace},
		{"Tab", tcelllib.KeyTab, terminal.KeyTab},
		{"Up", tcelllib.KeyUp, terminal.KeyUp},
		{"Down", tcelllib.KeyDown, terminal.KeyDown},
		{"Left", tcelllib.KeyLeft, terminal.KeyLeft},
		{"Right", tcelllib.KeyRight, terminal.KeyRight},
		{"Home", tcelllib.KeyHome, terminal.KeyHome},
		{"End", tcelllib.KeyEnd, terminal.KeyEnd},
		{"PgUp", tcelllib.KeyPgUp, terminal.KeyPgUp},
		{"PgDown", tcelllib.KeyPgDn, terminal.KeyPgDown},
		{"Delete", tcelllib.KeyDelete, terminal.KeyDelete},
		{"Insert", tcelllib.KeyInsert, terminal.KeyInsert},
		{"F1", tcelllib.KeyF1, terminal.KeyF1},
		{"F12", tcelllib.KeyF12, terminal.KeyF12},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			screen := newSpyScreen(10, 5)
			backend, err := NewWithScreen(screen)
			if err != nil {
				t.Fatalf("NewWithScreen() error = %v", err)
			}
			defer backend.Close()

			screen.InjectKey(tc.key, 0, tcelllib.ModNone)
			event, err := backend.PollEvent()
			if err != nil {
				t.Fatalf("PollEvent() error = %v", err)
			}

			key, ok := event.(terminal.KeyEvent)
			if !ok {
				t.Fatalf("PollEvent() = %T, want terminal.KeyEvent", event)
			}
			if key.Code != tc.want {
				t.Fatalf("Code = %v, want %v", key.Code, tc.want)
			}
		})
	}
}

func TestBackend_PollEvent_shouldConvertKeyModifiers(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	screen.InjectKey(tcelllib.KeyRune, 'x', tcelllib.ModCtrl|tcelllib.ModAlt|tcelllib.ModShift)
	event, err := backend.PollEvent()
	if err != nil {
		t.Fatalf("PollEvent() error = %v", err)
	}

	key, ok := event.(terminal.KeyEvent)
	if !ok {
		t.Fatalf("PollEvent() = %T, want terminal.KeyEvent", event)
	}
	want := terminal.ModifierCtrl | terminal.ModifierAlt | terminal.ModifierShift
	if key.Modifiers != want {
		t.Fatalf("Modifiers = %v, want %v", key.Modifiers, want)
	}
}

func TestBackend_PollEvent_shouldConvertResize(t *testing.T) {
	screen := newSpyScreen(10, 5)
	backend, err := NewWithScreen(screen)
	if err != nil {
		t.Fatalf("NewWithScreen() error = %v", err)
	}
	defer backend.Close()

	if err := screen.PostEvent(tcelllib.NewEventResize(80, 24)); err != nil {
		t.Fatalf("PostEvent() error = %v", err)
	}
	event, err := backend.PollEvent()
	if err != nil {
		t.Fatalf("PollEvent() error = %v", err)
	}

	resize, ok := event.(terminal.ResizeEvent)
	if !ok {
		t.Fatalf("PollEvent() = %T, want terminal.ResizeEvent", event)
	}
	want := layout.Size{Width: 80, Height: 24}
	if resize.Size != want {
		t.Fatalf("Size = %+v, want %+v", resize.Size, want)
	}
}

func TestBackend_PollEvent_shouldConvertMouse(t *testing.T) {
	cases := []struct {
		name    string
		buttons tcelllib.ButtonMask
		want    terminal.MouseButton
	}{
		{"Left", tcelllib.Button1, terminal.MouseButtonLeft},
		{"Right", tcelllib.Button2, terminal.MouseButtonRight},
		{"Middle", tcelllib.Button3, terminal.MouseButtonMiddle},
		{"WheelUp", tcelllib.WheelUp, terminal.MouseWheelUp},
		{"WheelDown", tcelllib.WheelDown, terminal.MouseWheelDown},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			screen := newSpyScreen(10, 5)
			backend, err := NewWithScreen(screen)
			if err != nil {
				t.Fatalf("NewWithScreen() error = %v", err)
			}
			defer backend.Close()

			screen.InjectMouse(3, 4, tc.buttons, tcelllib.ModCtrl|tcelllib.ModShift)
			event, err := backend.PollEvent()
			if err != nil {
				t.Fatalf("PollEvent() error = %v", err)
			}

			mouse, ok := event.(terminal.MouseEvent)
			if !ok {
				t.Fatalf("PollEvent() = %T, want terminal.MouseEvent", event)
			}
			if mouse.Position != (layout.Position{X: 3, Y: 4}) {
				t.Fatalf("Position = %+v, want (3, 4)", mouse.Position)
			}
			if mouse.Button != tc.want {
				t.Fatalf("Button = %v, want %v", mouse.Button, tc.want)
			}
			wantModifiers := terminal.ModifierCtrl | terminal.ModifierShift
			if mouse.Modifiers != wantModifiers {
				t.Fatalf("Modifiers = %v, want %v", mouse.Modifiers, wantModifiers)
			}
		})
	}
}

func TestTerminalBackendInterface_shouldAcceptTcellBackend(t *testing.T) {
	var _ terminal.Backend = (*Backend)(nil)
}
