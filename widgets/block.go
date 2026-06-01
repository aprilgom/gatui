package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type Block struct {
	title   text.Line
	borders Borders
	style   style.Style
}

func NewBlock() Block {
	return Block{style: style.NewStyle()}
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

func (b Block) Inner(area layout.Rect) layout.Rect {
	if b.borders == NoBorders {
		return area
	}
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
	inner.Width = maxInt(0, inner.Width-left-right)
	inner.Height = maxInt(0, inner.Height-top-bottom)
	return inner
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
	for _, span := range b.title.Spans {
		for _, r := range span.Content {
			if x >= area.X+area.Width {
				return
			}
			buf.SetCell(x, area.Y, buffer.Cell{Symbol: string(r), Style: b.style.Patch(span.Style)})
			x++
		}
	}
}

func (b Block) renderBorders(area layout.Rect, buf *buffer.Buffer) {
	right := area.X + area.Width - 1
	bottom := area.Y + area.Height - 1
	if b.borders.Has(TopBorder) {
		for x := area.X; x <= right; x++ {
			buf.SetCell(x, area.Y, buffer.Cell{Symbol: "─", Style: b.style})
		}
	}
	if b.borders.Has(BottomBorder) && bottom != area.Y {
		for x := area.X; x <= right; x++ {
			buf.SetCell(x, bottom, buffer.Cell{Symbol: "─", Style: b.style})
		}
	}
	if b.borders.Has(LeftBorder) {
		for y := area.Y; y <= bottom; y++ {
			buf.SetCell(area.X, y, buffer.Cell{Symbol: "│", Style: b.style})
		}
	}
	if b.borders.Has(RightBorder) && right != area.X {
		for y := area.Y; y <= bottom; y++ {
			buf.SetCell(right, y, buffer.Cell{Symbol: "│", Style: b.style})
		}
	}
	if b.borders.Has(TopBorder) && b.borders.Has(LeftBorder) {
		buf.SetCell(area.X, area.Y, buffer.Cell{Symbol: "┌", Style: b.style})
	}
	if b.borders.Has(TopBorder) && b.borders.Has(RightBorder) && right != area.X {
		buf.SetCell(right, area.Y, buffer.Cell{Symbol: "┐", Style: b.style})
	}
	if b.borders.Has(BottomBorder) && b.borders.Has(LeftBorder) && bottom != area.Y {
		buf.SetCell(area.X, bottom, buffer.Cell{Symbol: "└", Style: b.style})
	}
	if b.borders.Has(BottomBorder) && b.borders.Has(RightBorder) && right != area.X && bottom != area.Y {
		buf.SetCell(right, bottom, buffer.Cell{Symbol: "┘", Style: b.style})
	}
}

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
