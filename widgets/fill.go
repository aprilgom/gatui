package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

type Fill struct {
	symbol string
	style  style.Style
}

func NewFill(symbol string) Fill {
	if symbol == "" {
		symbol = " "
	}
	return Fill{symbol: symbol, style: style.NewStyle()}
}

func (f Fill) Style(cellStyle style.Style) Fill {
	f.style = cellStyle
	return f
}

func (f Fill) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.IsEmpty() {
		return
	}
	area = area.Intersection(buf.Area)
	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < area.X+area.Width; {
			endX, _ := buf.SetStringN(x, y, f.symbol, area.X+area.Width-x, f.style)
			if endX <= x {
				x++
				continue
			}
			x = endX
		}
	}
}
