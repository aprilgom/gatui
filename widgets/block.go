package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type Block struct {
	title       text.Line
	borders     Borders
	padding     Padding
	style       style.Style
	borderStyle style.Style
	titleStyle  style.Style
}

type Padding struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

func NewPadding(left, right, top, bottom int) Padding {
	return Padding{Left: left, Right: right, Top: top, Bottom: bottom}
}

func PaddingZero() Padding {
	return Padding{}
}

func PaddingHorizontal(value int) Padding {
	return NewPadding(value, value, 0, 0)
}

func PaddingVertical(value int) Padding {
	return NewPadding(0, 0, value, value)
}

func PaddingUniform(value int) Padding {
	return NewPadding(value, value, value, value)
}

func PaddingProportional(value int) Padding {
	return NewPadding(value*2, value*2, value, value)
}

func PaddingSymmetric(horizontal, vertical int) Padding {
	return NewPadding(horizontal, horizontal, vertical, vertical)
}

func PaddingLeft(value int) Padding {
	return NewPadding(value, 0, 0, 0)
}

func PaddingRight(value int) Padding {
	return NewPadding(0, value, 0, 0)
}

func PaddingTop(value int) Padding {
	return NewPadding(0, 0, value, 0)
}

func PaddingBottom(value int) Padding {
	return NewPadding(0, 0, 0, value)
}

func NewBlock() Block {
	return Block{
		style:       style.NewStyle(),
		borderStyle: style.NewStyle(),
		titleStyle:  style.NewStyle(),
	}
}

func BorderedBlock() Block {
	return NewBlock().Borders(AllBorders)
}

func (b Block) Title(title text.Line) Block {
	b.title = title
	return b
}

func (b Block) Borders(borders Borders) Block {
	b.borders = borders
	return b
}

func (b Block) Padding(padding Padding) Block {
	b.padding = padding
	return b
}

func (b Block) Style(blockStyle style.Style) Block {
	b.style = blockStyle
	return b
}

func (b Block) BorderStyle(borderStyle style.Style) Block {
	b.borderStyle = borderStyle
	return b
}

func (b Block) TitleStyle(titleStyle style.Style) Block {
	b.titleStyle = titleStyle
	return b
}

func (b Block) Inner(area layout.Rect) layout.Rect {
	left := 0
	right := 0
	top := 0
	bottom := 0
	if b.borders.Has(LeftBorder) {
		left = 1
	}
	if b.borders.Has(RightBorder) {
		right = 1
	}
	if b.borders.Has(TopBorder) {
		top = 1
	}
	if b.borders.Has(BottomBorder) {
		bottom = 1
	}
	inner := area
	inner.X += left
	inner.Y += top
	inner.Width = saturatingSub(inner.Width, saturatingAdd(left, right))
	inner.Height = saturatingSub(inner.Height, saturatingAdd(top, bottom))
	inner.X += b.padding.Left
	inner.Y += b.padding.Top
	inner.Width = saturatingSub(inner.Width, saturatingAdd(b.padding.Left, b.padding.Right))
	inner.Height = saturatingSub(inner.Height, saturatingAdd(b.padding.Top, b.padding.Bottom))
	return inner
}

func (b Block) horizontalSpace() int {
	space := saturatingAdd(b.padding.Left, b.padding.Right)
	if b.borders.Has(LeftBorder) {
		space = saturatingAdd(space, 1)
	}
	if b.borders.Has(RightBorder) {
		space = saturatingAdd(space, 1)
	}
	return space
}

func (b Block) verticalSpace() int {
	space := saturatingAdd(b.padding.Top, b.padding.Bottom)
	if b.borders.Has(TopBorder) {
		space = saturatingAdd(space, 1)
	}
	if b.borders.Has(BottomBorder) {
		space = saturatingAdd(space, 1)
	}
	return space
}

func (b Block) Fg(color style.Color) Block {
	b.style = b.style.Fg(color)
	return b
}

func (b Block) Bg(color style.Color) Block {
	b.style = b.style.Bg(color)
	return b
}

func (b Block) Bold() Block {
	b.style = b.style.AddModifier(style.ModifierBold)
	return b
}

func (b Block) Italic() Block {
	b.style = b.style.AddModifier(style.ModifierItalic)
	return b
}

func (b Block) Cyan() Block {
	return b.Fg(style.Cyan)
}

func (b Block) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	if b.borders != NoBorders {
		b.renderBorders(area, buf)
	}
	titleX := area.X
	if b.borders.Has(LeftBorder) {
		titleX++
	}
	x := titleX
	for _, grapheme := range b.title.StyledGraphemes(b.style.Patch(b.titleStyle)) {
		width := buffer.CellWidth(grapheme.Symbol)
		if width == 0 {
			continue
		}
		if x+width > area.X+area.Width {
			return
		}
		b.setCell(buf, x, area.Y, grapheme.Symbol, grapheme.Style)
		for trailing := 1; trailing < width; trailing++ {
			b.setCell(buf, x+trailing, area.Y, " ", grapheme.Style)
		}
		x += width
	}
}

func (b Block) setCell(buf *buffer.Buffer, x, y int, symbol string, cellStyle style.Style) {
	cell, ok := buf.CellAt(x, y)
	if !ok {
		return
	}
	cell.Symbol = symbol
	cell.Style = cell.Style.Patch(cellStyle)
	buf.SetCell(x, y, cell)
}

func (b Block) renderBorders(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	borderStyle := b.style.Patch(b.borderStyle)
	right := area.X + area.Width - 1
	bottom := area.Y + area.Height - 1
	if b.borders.Has(TopBorder) {
		for x := area.X; x <= right; x++ {
			b.setCell(buf, x, area.Y, "─", borderStyle)
		}
	}
	if b.borders.Has(BottomBorder) && bottom != area.Y {
		for x := area.X; x <= right; x++ {
			b.setCell(buf, x, bottom, "─", borderStyle)
		}
	}
	if b.borders.Has(LeftBorder) {
		for y := area.Y; y <= bottom; y++ {
			b.setCell(buf, area.X, y, "│", borderStyle)
		}
	}
	if b.borders.Has(RightBorder) && right != area.X {
		for y := area.Y; y <= bottom; y++ {
			b.setCell(buf, right, y, "│", borderStyle)
		}
	}
	if b.borders.Has(TopBorder) && b.borders.Has(LeftBorder) {
		b.setCell(buf, area.X, area.Y, "┌", borderStyle)
	}
	if b.borders.Has(TopBorder) && b.borders.Has(RightBorder) && right != area.X {
		b.setCell(buf, right, area.Y, "┐", borderStyle)
	}
	if b.borders.Has(BottomBorder) && b.borders.Has(LeftBorder) && bottom != area.Y {
		b.setCell(buf, area.X, bottom, "└", borderStyle)
	}
	if b.borders.Has(BottomBorder) && b.borders.Has(RightBorder) && right != area.X && bottom != area.Y {
		b.setCell(buf, right, bottom, "┘", borderStyle)
	}
}

func saturatingAdd(a, b int) int {
	if b > 0 && a > maxIntValue-b {
		return maxIntValue
	}
	if b < 0 && a < minIntValue-b {
		return minIntValue
	}
	return a + b
}

func saturatingSub(a, b int) int {
	if b <= 0 {
		return saturatingAdd(a, -b)
	}
	if a <= b {
		return 0
	}
	return a - b
}

const (
	maxIntValue = int(^uint(0) >> 1)
	minIntValue = -maxIntValue - 1
)

type Borders uint8

const (
	NoBorders    Borders = 0
	TopBorder    Borders = 1 << 0
	RightBorder  Borders = 1 << 1
	BottomBorder Borders = 1 << 2
	LeftBorder   Borders = 1 << 3
	AllBorders   Borders = TopBorder | RightBorder | BottomBorder | LeftBorder
)

func (b Borders) Has(border Borders) bool {
	return b&border != 0
}
