package widgets

import (
	"strings"
	"unicode"

	"gatui/buffer"
)

func wrapCells(cells []buffer.Cell, width int, trim bool) [][]buffer.Cell {
	var lines [][]buffer.Cell
	for len(cells) > 0 {
		if trim {
			cells = trimLeftCells(cells)
		}
		if cellsDisplayWidth(cells) <= width {
			lines = append(lines, trimRightCells(append([]buffer.Cell(nil), cells...), trim))
			break
		}
		breakAt := cellsThatFit(cells, width)
		for i := breakAt; i >= 0; i-- {
			if i < len(cells) && isSpaceCell(cells[i]) {
				breakAt = i
				break
			}
		}
		if breakAt == 0 {
			breakAt = cellsThatFit(cells, width)
		}
		line := append([]buffer.Cell(nil), cells[:breakAt]...)
		lines = append(lines, trimRightCells(line, trim))
		cells = cells[breakAt:]
	}
	if len(lines) == 0 {
		lines = append(lines, nil)
	}
	return lines
}

func cellsThatFit(cells []buffer.Cell, width int) int {
	used := 0
	for i, cell := range cells {
		cellWidth := cellDisplayWidth(cell)
		if used+cellWidth > width {
			return i
		}
		used += cellWidth
	}
	return len(cells)
}

func trimLeftCells(cells []buffer.Cell) []buffer.Cell {
	for len(cells) > 0 && isSpaceCell(cells[0]) {
		cells = cells[1:]
	}
	return cells
}

func trimRightCells(cells []buffer.Cell, trim bool) []buffer.Cell {
	if !trim {
		return cells
	}
	for len(cells) > 0 && isSpaceCell(cells[len(cells)-1]) {
		cells = cells[:len(cells)-1]
	}
	return cells
}

func isSpaceCell(cell buffer.Cell) bool {
	return strings.TrimFunc(cell.Symbol, unicode.IsSpace) == ""
}
