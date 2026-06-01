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
	marker          CanvasMarker
	xMin            float64
	xMax            float64
	yMin            float64
	yMax            float64
	paint           func(*CanvasContext)
}

type CanvasShape interface {
	Draw(*CanvasPainter)
}

type CanvasMarker int

const (
	CanvasMarkerDot CanvasMarker = iota
	CanvasMarkerBlock
)

type CanvasContext struct {
	labels []CanvasLabel
	shapes []CanvasShape
}

type CanvasLabel struct {
	X    float64
	Y    float64
	Span text.Span
}

type CanvasPoint struct {
	X float64
	Y float64
}

type CanvasPainter struct {
	width  int
	height int
	xMin   float64
	xMax   float64
	yMin   float64
	yMax   float64
	grid   []canvasPixel
}

type canvasPixel struct {
	painted bool
	color   style.Color
}

type Points struct {
	Coords []CanvasPoint
	Color  style.Color
}

type CanvasLine struct {
	X1    float64
	Y1    float64
	X2    float64
	Y2    float64
	Color style.Color
}

type Rectangle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Color  style.Color
}

func NewCanvas() Canvas {
	return Canvas{
		backgroundColor: style.Default,
		marker:          CanvasMarkerDot,
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

func (c Canvas) Marker(marker CanvasMarker) Canvas {
	c.marker = marker
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

func (ctx *CanvasContext) Draw(shape CanvasShape) {
	if ctx == nil || shape == nil {
		return
	}
	ctx.shapes = append(ctx.shapes, shape)
}

func NewPoints(coords []CanvasPoint, color style.Color) Points {
	return Points{Coords: coords, Color: color}
}

func (p Points) Draw(painter *CanvasPainter) {
	for _, coord := range p.Coords {
		if x, y, ok := painter.GetPoint(coord.X, coord.Y); ok {
			painter.Paint(x, y, p.Color)
		}
	}
}

func NewCanvasLine(x1, y1, x2, y2 float64, color style.Color) CanvasLine {
	return CanvasLine{X1: x1, Y1: y1, X2: x2, Y2: y2, Color: color}
}

func (l CanvasLine) Draw(painter *CanvasPainter) {
	x1, y1, x2, y2, ok := painter.clipLine(l.X1, l.Y1, l.X2, l.Y2)
	if !ok {
		return
	}
	startX, startY, ok := painter.GetPoint(x1, y1)
	if !ok {
		return
	}
	endX, endY, ok := painter.GetPoint(x2, y2)
	if !ok {
		return
	}
	painter.drawLine(startX, startY, endX, endY, l.Color)
}

func NewRectangle(x, y, width, height float64, color style.Color) Rectangle {
	return Rectangle{X: x, Y: y, Width: width, Height: height, Color: color}
}

func (r Rectangle) Draw(painter *CanvasPainter) {
	x2 := r.X + r.Width
	y2 := r.Y + r.Height
	NewCanvasLine(r.X, r.Y, x2, r.Y, r.Color).Draw(painter)
	NewCanvasLine(x2, r.Y, x2, y2, r.Color).Draw(painter)
	NewCanvasLine(x2, y2, r.X, y2, r.Color).Draw(painter)
	NewCanvasLine(r.X, y2, r.X, r.Y, r.Color).Draw(painter)
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
	painter := newCanvasPainter(area.Width, area.Height, c.xMin, c.xMax, c.yMin, c.yMax)
	for _, shape := range ctx.shapes {
		shape.Draw(painter)
	}
	c.renderShapes(area, buf, painter)
	for _, label := range ctx.labels {
		c.renderLabel(area, buf, label)
	}
}

func newCanvasPainter(width, height int, xMin, xMax, yMin, yMax float64) *CanvasPainter {
	return &CanvasPainter{
		width:  width,
		height: height,
		xMin:   xMin,
		xMax:   xMax,
		yMin:   yMin,
		yMax:   yMax,
		grid:   make([]canvasPixel, width*height),
	}
}

func (p *CanvasPainter) GetPoint(x, y float64) (int, int, bool) {
	if p == nil || p.width == 0 || p.height == 0 || p.xMin >= p.xMax || p.yMin >= p.yMax {
		return 0, 0, false
	}
	if x < p.xMin || x > p.xMax || y < p.yMin || y > p.yMax {
		return 0, 0, false
	}

	gridX := int(math.Floor((x - p.xMin) / (p.xMax - p.xMin) * float64(p.width)))
	gridY := p.height - 1 - int(math.Floor((y-p.yMin)/(p.yMax-p.yMin)*float64(p.height)))
	if gridX == p.width {
		gridX = p.width - 1
	}
	if gridY < 0 {
		gridY = 0
	}
	if gridX < 0 || gridY < 0 || gridX >= p.width || gridY >= p.height {
		return 0, 0, false
	}
	return gridX, gridY, true
}

func (p *CanvasPainter) Paint(x, y int, color style.Color) {
	if p == nil || x < 0 || y < 0 || x >= p.width || y >= p.height {
		return
	}
	p.grid[y*p.width+x] = canvasPixel{painted: true, color: color}
}

func (p *CanvasPainter) drawLine(x1, y1, x2, y2 int, color style.Color) {
	dx := absInt(x2 - x1)
	dy := -absInt(y2 - y1)
	stepX := -1
	if x1 < x2 {
		stepX = 1
	}
	stepY := -1
	if y1 < y2 {
		stepY = 1
	}
	err := dx + dy

	for {
		p.Paint(x1, y1, color)
		if x1 == x2 && y1 == y2 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += stepX
		}
		if e2 <= dx {
			err += dx
			y1 += stepY
		}
	}
}

func (p *CanvasPainter) clipLine(x1, y1, x2, y2 float64) (float64, float64, float64, float64, bool) {
	const (
		inside = 0
		left   = 1
		right  = 2
		bottom = 4
		top    = 8
	)

	outCode := func(x, y float64) int {
		code := inside
		if x < p.xMin {
			code |= left
		} else if x > p.xMax {
			code |= right
		}
		if y < p.yMin {
			code |= bottom
		} else if y > p.yMax {
			code |= top
		}
		return code
	}

	code1 := outCode(x1, y1)
	code2 := outCode(x2, y2)
	for {
		if code1|code2 == 0 {
			return x1, y1, x2, y2, true
		}
		if code1&code2 != 0 {
			return 0, 0, 0, 0, false
		}

		codeOut := code1
		if codeOut == 0 {
			codeOut = code2
		}

		x, y := 0.0, 0.0
		switch {
		case codeOut&top != 0:
			if y2 == y1 {
				return 0, 0, 0, 0, false
			}
			x = x1 + (x2-x1)*(p.yMax-y1)/(y2-y1)
			y = p.yMax
		case codeOut&bottom != 0:
			if y2 == y1 {
				return 0, 0, 0, 0, false
			}
			x = x1 + (x2-x1)*(p.yMin-y1)/(y2-y1)
			y = p.yMin
		case codeOut&right != 0:
			if x2 == x1 {
				return 0, 0, 0, 0, false
			}
			y = y1 + (y2-y1)*(p.xMax-x1)/(x2-x1)
			x = p.xMax
		case codeOut&left != 0:
			if x2 == x1 {
				return 0, 0, 0, 0, false
			}
			y = y1 + (y2-y1)*(p.xMin-x1)/(x2-x1)
			x = p.xMin
		}

		if codeOut == code1 {
			x1 = x
			y1 = y
			code1 = outCode(x1, y1)
		} else {
			x2 = x
			y2 = y
			code2 = outCode(x2, y2)
		}
	}
}

func (c Canvas) renderShapes(area layout.Rect, buf *buffer.Buffer, painter *CanvasPainter) {
	if painter == nil {
		return
	}
	for y := 0; y < painter.height; y++ {
		for x := 0; x < painter.width; x++ {
			pixel := painter.grid[y*painter.width+x]
			if !pixel.painted {
				continue
			}
			cellX := area.X + x
			cellY := area.Y + y
			cell, ok := buf.CellAt(cellX, cellY)
			if !ok {
				continue
			}
			switch c.marker {
			case CanvasMarkerBlock:
				cell.Symbol = "█"
				cell.Style = cell.Style.Patch(style.NewStyle().Fg(pixel.color).Bg(pixel.color))
			default:
				cell.Symbol = "•"
				cell.Style = cell.Style.Patch(style.NewStyle().Fg(pixel.color))
			}
			buf.SetCell(cellX, cellY, cell)
		}
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

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
