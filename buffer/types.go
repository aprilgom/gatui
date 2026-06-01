package buffer

import (
	"strings"

	"gatui/layout"
	"gatui/style"
)

type Cell struct {
	Symbol string
	Style  style.Style
}

type Buffer struct {
	Area  layout.Rect
	Cells []Cell
}

func Empty(area layout.Rect) *Buffer {
	cells := make([]Cell, area.Width*area.Height)
	for i := range cells {
		cells[i] = Cell{Symbol: " ", Style: style.NewStyle()}
	}
	return &Buffer{Area: area, Cells: cells}
}

func WithLines(lines []string) *Buffer {
	width := 0
	for _, line := range lines {
		if len([]rune(line)) > width {
			width = len([]rune(line))
		}
	}
	buf := Empty(layout.NewRect(0, 0, width, len(lines)))
	for y, line := range lines {
		for x, r := range []rune(line) {
			buf.SetSymbol(x, y, string(r))
		}
	}
	return buf
}

func (b *Buffer) CellAt(x, y int) (Cell, bool) {
	if b == nil || x < b.Area.X || y < b.Area.Y || x >= b.Area.X+b.Area.Width || y >= b.Area.Y+b.Area.Height {
		return Cell{}, false
	}
	index := (y-b.Area.Y)*b.Area.Width + (x - b.Area.X)
	return b.Cells[index], true
}

func (b *Buffer) SetCell(x, y int, cell Cell) {
	if b == nil || x < b.Area.X || y < b.Area.Y || x >= b.Area.X+b.Area.Width || y >= b.Area.Y+b.Area.Height {
		return
	}
	index := (y-b.Area.Y)*b.Area.Width + (x - b.Area.X)
	b.Cells[index] = cell
}

func (b *Buffer) SetSymbol(x, y int, symbol string) {
	if b == nil || x < b.Area.X || y < b.Area.Y || x >= b.Area.X+b.Area.Width || y >= b.Area.Y+b.Area.Height {
		return
	}
	index := (y-b.Area.Y)*b.Area.Width + (x - b.Area.X)
	b.Cells[index].Symbol = symbol
}

func (b *Buffer) SetStyle(area layout.Rect, cellStyle style.Style) {
	if b == nil {
		return
	}
	area = area.Intersection(b.Area)
	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < area.X+area.Width; x++ {
			index := (y-b.Area.Y)*b.Area.Width + (x - b.Area.X)
			b.Cells[index].Style = b.Cells[index].Style.Patch(cellStyle)
		}
	}
}

func (b *Buffer) SetFg(area layout.Rect, color style.Color) {
	b.SetStyle(area, style.NewStyle().Fg(color))
}

func (b *Buffer) SetBg(area layout.Rect, color style.Color) {
	b.SetStyle(area, style.NewStyle().Bg(color))
}

func (b *Buffer) SetModifier(area layout.Rect, modifier style.Modifier) {
	b.SetStyle(area, style.NewStyle().AddModifier(modifier))
}

func (b *Buffer) Lines() []string {
	if b == nil || b.Area.Height == 0 {
		return nil
	}

	lines := make([]string, b.Area.Height)
	for y := 0; y < b.Area.Height; y++ {
		var builder strings.Builder
		for x := 0; x < b.Area.Width; x++ {
			symbol := b.Cells[y*b.Area.Width+x].Symbol
			if symbol == "" {
				symbol = " "
			}
			builder.WriteString(symbol)
		}
		lines[y] = builder.String()
	}
	return lines
}
