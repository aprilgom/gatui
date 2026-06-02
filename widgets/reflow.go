package widgets

import (
	"strings"
	"unicode"

	"gatui/buffer"
	"gatui/layout"
)

type wrappedLine struct {
	cells     []buffer.Cell
	width     int
	alignment layout.Alignment
}

func (l wrappedLine) skip(count int) wrappedLine {
	for len(l.cells) > 0 && count > 0 {
		width := cellDisplayWidth(l.cells[0])
		if width > count {
			break
		}
		count -= width
		l.width -= width
		l.cells = l.cells[1:]
	}
	if l.width < 0 {
		l.width = 0
	}
	return l
}

type wordWrapper struct {
	width int
	trim  bool
}

func newWordWrapper(width int, trim bool) wordWrapper {
	return wordWrapper{width: width, trim: trim}
}

func (w wordWrapper) wrap(cells []buffer.Cell, alignment layout.Alignment) []wrappedLine {
	if w.width <= 0 {
		return nil
	}

	cells = append([]buffer.Cell(nil), cells...)
	if w.trim && len(trimLeftCells(cells)) == 0 {
		return []wrappedLine{{alignment: alignment}}
	}
	var lines []wrappedLine
	for len(cells) > 0 {
		if w.trim {
			cells = trimLeftCells(cells)
			if len(cells) == 0 {
				break
			}
		}
		if cellsDisplayWidth(cells) <= w.width {
			lineCells := trimRightCells(append([]buffer.Cell(nil), cells...), w.trim)
			lines = append(lines, newWrappedLine(lineCells, alignment))
			break
		}

		fitEnd := cellsThatFit(cells, w.width)
		if fitEnd == 0 {
			fitEnd = 1
		}

		breakAt := fitEnd
		if fitEnd < len(cells) && isSpaceCell(cells[fitEnd]) {
			breakAt = fitEnd
		} else if spaceAt, ok := lastWhitespaceBreak(cells[:fitEnd], w.trim); ok {
			breakAt = spaceAt
		}
		if breakAt == 0 {
			breakAt = fitEnd
		}

		lineCells := append([]buffer.Cell(nil), cells[:breakAt]...)
		lineCells = trimRightCells(lineCells, w.trim)
		if len(lineCells) > 0 || !w.trim {
			lines = append(lines, newWrappedLine(lineCells, alignment))
		}
		cells = cells[breakAt:]
	}
	return lines
}

type lineTruncator struct {
	width int
}

func newLineTruncator(width int) lineTruncator {
	return lineTruncator{width: width}
}

func (t lineTruncator) truncate(cells []buffer.Cell, alignment layout.Alignment) wrappedLine {
	if t.width <= 0 {
		return wrappedLine{alignment: alignment}
	}
	end := cellsThatFitPositiveWidth(cells, t.width)
	return newWrappedLine(append([]buffer.Cell(nil), cells[:end]...), alignment)
}

func newWrappedLine(cells []buffer.Cell, alignment layout.Alignment) wrappedLine {
	return wrappedLine{cells: cells, width: cellsDisplayWidth(cells), alignment: alignment}
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

func cellsThatFitPositiveWidth(cells []buffer.Cell, width int) int {
	used := 0
	end := 0
	for i, cell := range cells {
		cellWidth := cellDisplayWidth(cell)
		if cellWidth == 0 {
			continue
		}
		if used+cellWidth > width {
			break
		}
		used += cellWidth
		end = i + 1
	}
	return end
}

func lastWhitespaceBreak(cells []buffer.Cell, trim bool) (int, bool) {
	for i := len(cells) - 1; i >= 0; i-- {
		if !isSpaceCell(cells[i]) {
			continue
		}
		if !trim && !hasNonSpaceBefore(cells, i) {
			continue
		}
		return i, true
	}
	return 0, false
}

func hasNonSpaceBefore(cells []buffer.Cell, index int) bool {
	for i := range index {
		if !isSpaceCell(cells[i]) {
			return true
		}
	}
	return false
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
