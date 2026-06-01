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
