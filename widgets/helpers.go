package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

func writeString(buf *buffer.Buffer, x, y int, value string, width int, cellStyle style.Style) {
	runes := []rune(value)
	for i := 0; i < width; i++ {
		symbol := " "
		if i < len(runes) {
			symbol = string(runes[i])
		}
		buf.SetCell(x+i, y, buffer.Cell{Symbol: symbol, Style: cellStyle})
	}
}

func writeStringWithin(buf *buffer.Buffer, x, y, right int, value string, cellStyle style.Style) int {
	for _, r := range value {
		if x >= right {
			return x
		}
		buf.SetCell(x, y, buffer.Cell{Symbol: string(r), Style: cellStyle})
		x++
	}
	return x
}

func lineWidth(line text.Line) int {
	return line.Width()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func cellsFromLine(line text.Line) []buffer.Cell {
	return cellsFromLineWithStyle(line, style.NewStyle())
}

func cellsFromLineWithStyle(line text.Line, textStyle style.Style) []buffer.Cell {
	var cells []buffer.Cell
	for _, grapheme := range line.StyledGraphemes(textStyle) {
		cells = append(cells, buffer.Cell{Symbol: grapheme.Symbol, Style: grapheme.Style})
	}
	return cells
}

func cellDisplayWidth(cell buffer.Cell) int {
	return cell.Width()
}

func cellsDisplayWidth(cells []buffer.Cell) int {
	width := 0
	for _, cell := range cells {
		width += cellDisplayWidth(cell)
	}
	return width
}

func alignedOffset(lineWidth, areaWidth int, alignment layout.Alignment) int {
	if lineWidth >= areaWidth {
		return 0
	}
	switch alignment {
	case layout.Center:
		return (areaWidth - lineWidth) / 2
	case layout.Right:
		return areaWidth - lineWidth
	default:
		return 0
	}
}
