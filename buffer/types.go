package buffer

import (
	"strings"

	"gatui/layout"
	"gatui/style"

	"github.com/rivo/uniseg"
)

type Cell struct {
	Symbol      string
	Style       style.Style
	DiffOption  CellDiffOption
	ForcedWidth int
}

type CellDiffOption int

const (
	CellDiffNone CellDiffOption = iota
	CellDiffSkip
	CellDiffAlwaysUpdate
	CellDiffForcedWidth
)

type Buffer struct {
	Area  layout.Rect
	Cells []Cell
}

type CellDiff struct {
	X    int
	Y    int
	Cell Cell
}

func NewCell(symbol string) Cell {
	return Cell{Symbol: symbol, Style: style.NewStyle()}
}

func (c Cell) DisplaySymbol() string {
	if c.Symbol == "" {
		return " "
	}
	return c.Symbol
}

func (c Cell) Width() int {
	if c.ForcedWidth > 0 {
		return c.ForcedWidth
	}
	return uniseg.StringWidth(c.DisplaySymbol())
}

func (c *Cell) SetSymbol(symbol string) {
	c.Symbol = symbol
}

func (c *Cell) SetChar(char rune) {
	c.SetSymbol(string(char))
}

func (c *Cell) AppendSymbol(symbol string) {
	c.Symbol += symbol
}

func (c *Cell) SetStyle(cellStyle style.Style) {
	c.Style = c.Style.Patch(cellStyle)
}

func (c *Cell) Reset() {
	*c = NewCell(" ")
}

func (c *Cell) SetDiffOption(option CellDiffOption) {
	c.DiffOption = option
}

func (c *Cell) SetForcedWidth(width int) {
	c.ForcedWidth = width
}

func Empty(area layout.Rect) *Buffer {
	cells := make([]Cell, area.Width*area.Height)
	for i := range cells {
		cells[i] = NewCell(" ")
	}
	return &Buffer{Area: area, Cells: cells}
}

func WithLines(lines []string) *Buffer {
	width := 0
	for _, line := range lines {
		if lineWidth := uniseg.StringWidth(line); lineWidth > width {
			width = lineWidth
		}
	}
	buf := Empty(layout.NewRect(0, 0, width, len(lines)))
	for y, line := range lines {
		x := 0
		graphemes := uniseg.NewGraphemes(line)
		for graphemes.Next() {
			symbol := graphemes.Str()
			buf.SetSymbol(x, y, symbol)
			x += uniseg.StringWidth(symbol)
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

func (b *Buffer) Diff(next *Buffer) []CellDiff {
	if b.Area.X != next.Area.X || b.Area.Y != next.Area.Y || b.Area.Width != next.Area.Width {
		panic("buffer areas must have the same x, y, and width")
	}

	height := b.Area.Height
	if next.Area.Height < height {
		height = next.Area.Height
	}
	diffs := make([]CellDiff, 0)
	for y := 0; y < height; y++ {
		for x := 0; x < b.Area.Width; x++ {
			index := y*b.Area.Width + x
			previous := b.Cells[index]
			current := next.Cells[index]
			width := current.Width()

			if current.DiffOption == CellDiffSkip {
				continue
			}
			if current.DiffOption == CellDiffAlwaysUpdate || current != previous {
				diffs = append(diffs, CellDiff{
					X:    next.Area.X + x,
					Y:    next.Area.Y + y,
					Cell: current,
				})
			}
			if current.DiffOption == CellDiffForcedWidth || current.ForcedWidth > 0 || width > 1 {
				x += width - 1
			}
		}
	}
	return diffs
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
