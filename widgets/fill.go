package widgets

import (
	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
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

func (f Fill) Symbol(symbol string) Fill {
	if symbol == "" {
		symbol = " "
	}
	f.symbol = symbol
	return f
}

func (f Fill) Style(cellStyle style.Style) Fill {
	f.style = cellStyle
	return f
}

func (f Fill) Fg(color style.Color) Fill {
	f.style = f.style.Fg(color)
	return f
}

func (f Fill) Bg(color style.Color) Fill {
	f.style = f.style.Bg(color)
	return f
}

func (f Fill) Bold() Fill {
	f.style = f.style.AddModifier(style.ModifierBold)
	return f
}

func (f Fill) Dim() Fill {
	f.style = f.style.AddModifier(style.ModifierDim)
	return f
}

func (f Fill) Italic() Fill {
	f.style = f.style.AddModifier(style.ModifierItalic)
	return f
}

func (f Fill) Cyan() Fill {
	return f.Fg(style.Cyan)
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
