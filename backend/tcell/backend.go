package tcell

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/terminal"

	tcelllib "github.com/gdamore/tcell/v2"
)

var _ terminal.Backend = (*Backend)(nil)

type Backend struct {
	screen tcelllib.Screen
}

func New() (*Backend, error) {
	screen, err := tcelllib.NewScreen()
	if err != nil {
		return nil, err
	}
	return NewWithScreen(screen)
}

func NewWithScreen(screen tcelllib.Screen) (*Backend, error) {
	if err := screen.Init(); err != nil {
		return nil, err
	}
	return &Backend{screen: screen}, nil
}

func (b *Backend) Close() {
	b.screen.Fini()
}

func (b *Backend) Size() (layout.Size, error) {
	width, height := b.screen.Size()
	return layout.Size{Width: width, Height: height}, nil
}

func (b *Backend) Draw(diffs []buffer.CellDiff) error {
	for _, diff := range diffs {
		symbol := diff.Cell.DisplaySymbol()
		runes := []rune(symbol)
		primary := ' '
		var combining []rune
		if len(runes) > 0 {
			primary = runes[0]
			if len(runes) > 1 {
				combining = runes[1:]
			}
		}
		b.screen.SetContent(diff.X, diff.Y, primary, combining, convertStyle(diff.Cell.Style))
	}
	return nil
}

func (b *Backend) Flush() error {
	b.screen.Show()
	return nil
}

func (b *Backend) Clear() error {
	b.screen.Clear()
	return nil
}

func (b *Backend) PollEvent() (terminal.Event, error) {
	return convertEvent(b.screen.PollEvent()), nil
}

func (b *Backend) HideCursor() error {
	b.screen.HideCursor()
	return nil
}

func (b *Backend) ShowCursor() error {
	b.screen.ShowCursor(0, 0)
	return nil
}

func (b *Backend) SetCursorPosition(pos layout.Position) error {
	b.screen.ShowCursor(pos.X, pos.Y)
	return nil
}

func convertEvent(event tcelllib.Event) terminal.Event {
	switch event := event.(type) {
	case *tcelllib.EventKey:
		return terminal.KeyEvent{
			Code:      convertKeyCode(event.Key()),
			Rune:      event.Rune(),
			Modifiers: convertModifiers(event.Modifiers()),
		}
	case *tcelllib.EventMouse:
		x, y := event.Position()
		return terminal.MouseEvent{
			Position:  layout.Position{X: x, Y: y},
			Button:    convertMouseButton(event.Buttons()),
			Modifiers: convertModifiers(event.Modifiers()),
		}
	case *tcelllib.EventResize:
		width, height := event.Size()
		return terminal.ResizeEvent{Size: layout.Size{Width: width, Height: height}}
	default:
		return terminal.UnknownEvent{}
	}
}

func convertKeyCode(key tcelllib.Key) terminal.KeyCode {
	switch key {
	case tcelllib.KeyRune:
		return terminal.KeyRune
	case tcelllib.KeyEnter:
		return terminal.KeyEnter
	case tcelllib.KeyEsc:
		return terminal.KeyEsc
	case tcelllib.KeyBackspace, tcelllib.KeyBackspace2:
		return terminal.KeyBackspace
	case tcelllib.KeyTab:
		return terminal.KeyTab
	case tcelllib.KeyUp:
		return terminal.KeyUp
	case tcelllib.KeyDown:
		return terminal.KeyDown
	case tcelllib.KeyLeft:
		return terminal.KeyLeft
	case tcelllib.KeyRight:
		return terminal.KeyRight
	case tcelllib.KeyHome:
		return terminal.KeyHome
	case tcelllib.KeyEnd:
		return terminal.KeyEnd
	case tcelllib.KeyPgUp:
		return terminal.KeyPgUp
	case tcelllib.KeyPgDn:
		return terminal.KeyPgDown
	case tcelllib.KeyDelete:
		return terminal.KeyDelete
	case tcelllib.KeyInsert:
		return terminal.KeyInsert
	case tcelllib.KeyF1:
		return terminal.KeyF1
	case tcelllib.KeyF2:
		return terminal.KeyF2
	case tcelllib.KeyF3:
		return terminal.KeyF3
	case tcelllib.KeyF4:
		return terminal.KeyF4
	case tcelllib.KeyF5:
		return terminal.KeyF5
	case tcelllib.KeyF6:
		return terminal.KeyF6
	case tcelllib.KeyF7:
		return terminal.KeyF7
	case tcelllib.KeyF8:
		return terminal.KeyF8
	case tcelllib.KeyF9:
		return terminal.KeyF9
	case tcelllib.KeyF10:
		return terminal.KeyF10
	case tcelllib.KeyF11:
		return terminal.KeyF11
	case tcelllib.KeyF12:
		return terminal.KeyF12
	default:
		return terminal.KeyUnknown
	}
}

func convertModifiers(modifiers tcelllib.ModMask) terminal.KeyModifier {
	var converted terminal.KeyModifier
	if modifiers&tcelllib.ModCtrl != 0 {
		converted |= terminal.ModifierCtrl
	}
	if modifiers&tcelllib.ModAlt != 0 {
		converted |= terminal.ModifierAlt
	}
	if modifiers&tcelllib.ModShift != 0 {
		converted |= terminal.ModifierShift
	}
	return converted
}

func convertMouseButton(buttons tcelllib.ButtonMask) terminal.MouseButton {
	switch {
	case buttons&tcelllib.Button1 != 0:
		return terminal.MouseButtonLeft
	case buttons&tcelllib.Button2 != 0:
		return terminal.MouseButtonRight
	case buttons&tcelllib.Button3 != 0:
		return terminal.MouseButtonMiddle
	case buttons&tcelllib.WheelUp != 0:
		return terminal.MouseWheelUp
	case buttons&tcelllib.WheelDown != 0:
		return terminal.MouseWheelDown
	default:
		return terminal.MouseButtonNone
	}
}

func convertStyle(cellStyle style.Style) tcelllib.Style {
	tcellStyle := tcelllib.StyleDefault.
		Foreground(convertColor(cellStyle.Foreground)).
		Background(convertColor(cellStyle.Background))

	if cellStyle.Modifiers&style.ModifierBold != 0 {
		tcellStyle = tcellStyle.Bold(true)
	}
	if cellStyle.Modifiers&style.ModifierDim != 0 {
		tcellStyle = tcellStyle.Dim(true)
	}
	if cellStyle.Modifiers&style.ModifierItalic != 0 {
		tcellStyle = tcellStyle.Italic(true)
	}
	if cellStyle.Modifiers&style.ModifierUnderlined != 0 {
		tcellStyle = tcellStyle.Underline(true)
	}
	if cellStyle.Modifiers&style.ModifierReversed != 0 {
		tcellStyle = tcellStyle.Reverse(true)
	}

	return tcellStyle
}

func convertColor(color style.Color) tcelllib.Color {
	switch color {
	case style.Black:
		return tcelllib.ColorBlack
	case style.Red:
		return tcelllib.ColorRed
	case style.Green:
		return tcelllib.ColorGreen
	case style.Yellow:
		return tcelllib.ColorYellow
	case style.Blue:
		return tcelllib.ColorBlue
	case style.Magenta:
		return tcelllib.ColorFuchsia
	case style.Cyan:
		return tcelllib.ColorAqua
	case style.White:
		return tcelllib.ColorWhite
	case style.LightBlue:
		return tcelllib.ColorLightBlue
	case style.LightGreen:
		return tcelllib.ColorLightGreen
	default:
		return tcelllib.ColorDefault
	}
}
