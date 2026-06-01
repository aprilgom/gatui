package buffer

import (
	"strings"

	"gatui/layout"
	"gatui/style"

	"github.com/rivo/uniseg"
)

type Buffer struct {
	Area  layout.Rect
	Cells []Cell
}

func Empty(area layout.Rect) *Buffer {
	return Filled(area, NewCell(" "))
}

func Filled(area layout.Rect, cell Cell) *Buffer {
	cells := make([]Cell, area.Width*area.Height)
	for i := range cells {
		cells[i] = cell
	}
	return &Buffer{Area: area, Cells: cells}
}

func WithLines(lines []string) *Buffer {
	width := 0
	for _, line := range lines {
		if lineWidth := CellWidth(line); lineWidth > width {
			width = lineWidth
		}
	}
	buf := Empty(layout.NewRect(0, 0, width, len(lines)))
	for y, line := range lines {
		x := 0
		graphemes := uniseg.NewGraphemes(line)
		for graphemes.Next() {
			for _, symbol := range cellWidthSymbols(graphemes.Str()) {
				buf.SetSymbol(x, y, symbol)
				x += CellWidth(symbol)
			}
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

func (b *Buffer) SetString(x, y int, value string, cellStyle style.Style) (endX, endY int) {
	if b == nil || x < b.Area.X || y < b.Area.Y || x >= b.Area.X+b.Area.Width || y >= b.Area.Y+b.Area.Height {
		return x, y
	}
	return b.SetStringN(x, y, value, b.Area.X+b.Area.Width-x, cellStyle)
}

func (b *Buffer) SetStringN(x, y int, value string, maxWidth int, cellStyle style.Style) (endX, endY int) {
	if b == nil || maxWidth <= 0 || x < b.Area.X || y < b.Area.Y || x >= b.Area.X+b.Area.Width || y >= b.Area.Y+b.Area.Height {
		return x, y
	}

	remainingWidth := minInt(b.Area.X+b.Area.Width-x, maxWidth)
	graphemes := uniseg.NewGraphemes(value)
writeLoop:
	for graphemes.Next() {
		for _, symbol := range cellWidthSymbols(graphemes.Str()) {
			if containsControl(symbol) {
				continue
			}

			width := CellWidth(symbol)
			if width == 0 {
				if x > b.Area.X {
					index := (y-b.Area.Y)*b.Area.Width + (x - 1 - b.Area.X)
					b.Cells[index].AppendSymbol(symbol)
				}
				continue
			}
			if width > remainingWidth {
				break writeLoop
			}

			index := (y-b.Area.Y)*b.Area.Width + (x - b.Area.X)
			b.Cells[index].SetSymbol(symbol)
			b.Cells[index].SetStyle(cellStyle)
			x++
			remainingWidth--

			for trailing := 1; trailing < width; trailing++ {
				index := (y-b.Area.Y)*b.Area.Width + (x - b.Area.X)
				b.Cells[index].Reset()
				x++
				remainingWidth--
			}
		}
	}

	return x, y
}

func cellWidthSymbols(symbol string) []string {
	for _, r := range symbol {
		if isHalfwidthVoicingMark(r) {
			return splitHalfwidthVoicingMarks(symbol)
		}
	}
	return []string{symbol}
}

func splitHalfwidthVoicingMarks(symbol string) []string {
	symbols := make([]string, 0, len(symbol))
	var b strings.Builder
	for _, r := range symbol {
		if isHalfwidthVoicingMark(r) {
			if b.Len() > 0 {
				symbols = append(symbols, b.String())
				b.Reset()
			}
			symbols = append(symbols, string(r))
			continue
		}
		b.WriteRune(r)
	}
	if b.Len() > 0 {
		symbols = append(symbols, b.String())
	}
	return symbols
}

func isHalfwidthVoicingMark(r rune) bool {
	return r == '\uff9e' || r == '\uff9f'
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

func (b *Buffer) Reset() {
	if b == nil {
		return
	}
	for i := range b.Cells {
		b.Cells[i].Reset()
	}
}

func (b *Buffer) Resize(area layout.Rect) {
	if b == nil {
		return
	}
	length := area.Width * area.Height
	if len(b.Cells) > length {
		b.Cells = b.Cells[:length]
	} else {
		for len(b.Cells) < length {
			b.Cells = append(b.Cells, NewCell(" "))
		}
	}
	b.Area = area
}

func (b *Buffer) Merge(other *Buffer) {
	if b == nil || other == nil {
		return
	}
	area := unionRect(b.Area, other.Area)
	cells := make([]Cell, area.Width*area.Height)
	for i := range cells {
		cells[i] = NewCell(" ")
	}

	for y := b.Area.Y; y < b.Area.Y+b.Area.Height; y++ {
		for x := b.Area.X; x < b.Area.X+b.Area.Width; x++ {
			cell, ok := b.CellAt(x, y)
			if !ok {
				continue
			}
			index := (y-area.Y)*area.Width + (x - area.X)
			cells[index] = cell
		}
	}
	for y := other.Area.Y; y < other.Area.Y+other.Area.Height; y++ {
		for x := other.Area.X; x < other.Area.X+other.Area.Width; x++ {
			cell, ok := other.CellAt(x, y)
			if !ok {
				continue
			}
			index := (y-area.Y)*area.Width + (x - area.X)
			cells[index] = cell
		}
	}

	b.Area = area
	b.Cells = cells
}

func (b *Buffer) Lines() []string {
	if b == nil || b.Area.Height == 0 {
		return nil
	}

	lines := make([]string, b.Area.Height)
	for y := 0; y < b.Area.Height; y++ {
		var builder strings.Builder
		for x := 0; x < b.Area.Width; x++ {
			symbol := b.Cells[y*b.Area.Width+x].DisplaySymbol()
			builder.WriteString(symbol)
			if width := b.Cells[y*b.Area.Width+x].Width(); width > 1 {
				for skipped := 0; skipped < width-1 && x+1 < b.Area.Width; skipped++ {
					next := b.Cells[y*b.Area.Width+x+1].Symbol
					if next != "" && next != " " {
						break
					}
					x++
				}
			}
		}
		lines[y] = builder.String()
	}
	return lines
}
