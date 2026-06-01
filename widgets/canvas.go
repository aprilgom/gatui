package widgets

import (
	"math"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type Canvas struct {
	backgroundColor style.Color
	xMin            float64
	xMax            float64
	yMin            float64
	yMax            float64
	paint           func(*CanvasContext)
}

type CanvasContext struct {
	labels []CanvasLabel
}

type CanvasLabel struct {
	X    float64
	Y    float64
	Span text.Span
}

func NewCanvas() Canvas {
	return Canvas{
		backgroundColor: style.Default,
		xMin:            0,
		xMax:            1,
		yMin:            0,
		yMax:            1,
	}
}

func (c Canvas) BackgroundColor(color style.Color) Canvas {
	c.backgroundColor = color
	return c
}

func (c Canvas) XBounds(min, max float64) Canvas {
	c.xMin = min
	c.xMax = max
	return c
}

func (c Canvas) YBounds(min, max float64) Canvas {
	c.yMin = min
	c.yMax = max
	return c
}

func (c Canvas) Paint(paint func(*CanvasContext)) Canvas {
	c.paint = paint
	return c
}

func (ctx *CanvasContext) Print(x, y float64, span text.Span) {
	if ctx == nil {
		return
	}
	ctx.labels = append(ctx.labels, CanvasLabel{X: x, Y: y, Span: span})
}

func (c Canvas) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}

	buf.SetBg(area, c.backgroundColor)
	if c.xMin >= c.xMax || c.yMin >= c.yMax {
		return
	}

	ctx := &CanvasContext{}
	if c.paint != nil {
		c.paint(ctx)
	}
	for _, label := range ctx.labels {
		c.renderLabel(area, buf, label)
	}
}

func (c Canvas) renderLabel(area layout.Rect, buf *buffer.Buffer, label CanvasLabel) {
	if label.X < c.xMin || label.X > c.xMax || label.Y < c.yMin || label.Y > c.yMax {
		return
	}

	x := area.X + int(math.Floor((label.X-c.xMin)/(c.xMax-c.xMin)*float64(area.Width)))
	y := area.Y + area.Height - 1 - int(math.Floor((label.Y-c.yMin)/(c.yMax-c.yMin)*float64(area.Height)))
	if x < area.X || y < area.Y || x >= area.Right() || y >= area.Bottom() {
		return
	}

	for offset, r := range []rune(label.Span.Content) {
		cellX := x + offset
		if cellX >= area.Right() {
			return
		}
		cell, ok := buf.CellAt(cellX, y)
		if !ok {
			continue
		}
		cell.Symbol = string(r)
		cell.Style = cell.Style.Patch(label.Span.Style)
		buf.SetCell(cellX, y, cell)
	}
}
