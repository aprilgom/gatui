package tcell

import (
	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/terminal"

	tcelllib "github.com/gdamore/tcell/v3"
	tcellcolor "github.com/gdamore/tcell/v3/color"
)

var _ terminal.Backend = (*Backend)(nil)

type Backend struct {
	screen         tcelllib.Screen
	cursorPosition layout.Position
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

func (b *Backend) WindowSize() (terminal.WindowSize, error) {
	size, err := b.Size()
	if err != nil {
		return terminal.WindowSize{}, err
	}
	return terminal.WindowSize{
		ColumnsRows: size,
		Pixels:      layout.Size{Width: 0, Height: 0},
	}, nil
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
	return b.ClearRegion(terminal.ClearAll)
}

func (b *Backend) ClearRegion(clearType terminal.ClearType) error {
	switch clearType {
	case terminal.ClearAll:
		b.screen.Clear()
	case terminal.ClearAfterCursor:
		size, err := b.Size()
		if err != nil {
			return err
		}
		for y := b.cursorPosition.Y; y < size.Height; y++ {
			startX := 0
			if y == b.cursorPosition.Y {
				startX = b.cursorPosition.X
			}
			for x := startX; x < size.Width; x++ {
				b.screen.SetContent(x, y, ' ', nil, tcelllib.StyleDefault)
			}
		}
	case terminal.ClearBeforeCursor:
		size, err := b.Size()
		if err != nil {
			return err
		}
		for y := 0; y <= b.cursorPosition.Y && y < size.Height; y++ {
			endX := size.Width - 1
			if y == b.cursorPosition.Y {
				endX = b.cursorPosition.X
			}
			for x := 0; x <= endX && x < size.Width; x++ {
				b.screen.SetContent(x, y, ' ', nil, tcelllib.StyleDefault)
			}
		}
	case terminal.ClearCurrentLine:
		size, err := b.Size()
		if err != nil {
			return err
		}
		for x := 0; x < size.Width; x++ {
			b.screen.SetContent(x, b.cursorPosition.Y, ' ', nil, tcelllib.StyleDefault)
		}
	case terminal.ClearUntilNewLine:
		size, err := b.Size()
		if err != nil {
			return err
		}
		for x := b.cursorPosition.X; x < size.Width; x++ {
			b.screen.SetContent(x, b.cursorPosition.Y, ' ', nil, tcelllib.StyleDefault)
		}
	}
	return nil
}

func (b *Backend) GetCursorPosition() (layout.Position, error) {
	return b.cursorPosition, nil
}

func (b *Backend) AppendLines(count int) error {
	if count <= 0 {
		return nil
	}
	size, err := b.Size()
	if err != nil {
		return err
	}
	if size.Height <= 0 {
		return nil
	}
	b.cursorPosition = layout.Position{X: 0, Y: size.Height - 1}
	for y := 0; y < size.Height; y++ {
		for x := 0; x < size.Width; x++ {
			b.screen.SetContent(x, y, ' ', nil, tcelllib.StyleDefault)
		}
	}
	b.screen.ShowCursor(b.cursorPosition.X, b.cursorPosition.Y)
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
	b.cursorPosition = pos
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
		return tcellcolor.Black
	case style.Red:
		return tcellcolor.Red
	case style.Green:
		return tcellcolor.Green
	case style.Yellow:
		return tcellcolor.Yellow
	case style.Blue:
		return tcellcolor.Blue
	case style.Magenta:
		return tcellcolor.Fuchsia
	case style.Cyan:
		return tcellcolor.Aqua
	case style.White:
		return tcellcolor.White
	case style.LightBlue:
		return tcellcolor.LightBlue
	case style.LightGreen:
		return tcellcolor.LightGreen
	default:
		return tcellcolor.Default
	}
}
