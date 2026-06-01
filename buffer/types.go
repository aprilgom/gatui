package buffer

import (
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

func (b *Buffer) SetCell(x, y int, cell Cell) {
	if b == nil || x < b.Area.X || y < b.Area.Y || x >= b.Area.X+b.Area.Width || y >= b.Area.Y+b.Area.Height {
		return
	}
	index := (y-b.Area.Y)*b.Area.Width + (x - b.Area.X)
	b.Cells[index] = cell
}
