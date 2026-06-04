package buffer

import (
	"gatui/style"
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

func NewCell(symbol string) Cell {
	return Cell{Symbol: symbol, Style: style.NewStyle()}
}

func (c Cell) DisplaySymbol() string {
	if c.Symbol == "" {
		return " "
	}
	return c.Symbol
}

func (c Cell) Equal(other Cell) bool {
	return c.DisplaySymbol() == other.DisplaySymbol() &&
		c.Style == other.Style &&
		c.DiffOption == other.DiffOption &&
		c.ForcedWidth == other.ForcedWidth
}

func (c Cell) Width() int {
	if c.ForcedWidth > 0 {
		return c.ForcedWidth
	}
	return CellWidth(c.DisplaySymbol())
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

func (c *Cell) SetFg(color style.Color) {
	c.SetStyle(style.NewStyle().Fg(color))
}

func (c *Cell) SetBg(color style.Color) {
	c.SetStyle(style.NewStyle().Bg(color))
}

func (c Cell) StyleValue() style.Style {
	return c.Style
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
