package widgets

import (
	"gatui/buffer"
	"gatui/layout"
)

type renderLine struct {
	cells     []buffer.Cell
	alignment layout.Alignment
}

func (l renderLine) width() int {
	return cellsDisplayWidth(l.cells)
}

func (l renderLine) skip(count int) renderLine {
	for len(l.cells) > 0 && count > 0 {
		width := cellDisplayWidth(l.cells[0])
		if width > count {
			break
		}
		count -= width
		l.cells = l.cells[1:]
	}
	return l
}
