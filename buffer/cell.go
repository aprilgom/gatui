package buffer

import (
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
