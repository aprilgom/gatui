package widgets

import (
	"math"
	"strings"

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

type CanvasMarker string

const (
	CanvasMarkerDot       CanvasMarker = "dot"
	CanvasMarkerBlock     CanvasMarker = "block"
	CanvasMarkerBar       CanvasMarker = "bar"
	CanvasMarkerBraille   CanvasMarker = "braille"
	CanvasMarkerHalfBlock CanvasMarker = "half_block"
	CanvasMarkerQuadrant  CanvasMarker = "quadrant"
	CanvasMarkerSextant   CanvasMarker = "sextant"
	CanvasMarkerOctant    CanvasMarker = "octant"
)

func CanvasMarkerCustom(symbol string) CanvasMarker {
	runes := []rune(symbol)
	if len(runes) == 0 {
		return CanvasMarker("custom: ")
	}
	return CanvasMarker("custom:" + string(runes[0]))
}

func (m CanvasMarker) kind() CanvasMarker {
	if strings.HasPrefix(string(m), "custom:") {
		return "custom"
	}
	return m
}

func (m CanvasMarker) customSymbol() string {
	symbol, ok := strings.CutPrefix(string(m), "custom:")
	if !ok || symbol == "" {
		return " "
	}
	return symbol
}

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
	marker CanvasMarker
	grid   []canvasPixel
}

type canvasPixel struct {
	painted bool
	color   style.Color
	pattern uint8
	upper   *style.Color
	lower   *style.Color
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

type Circle struct {
	X      float64
	Y      float64
	Radius float64
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

func NewCircle(x, y, radius float64, color style.Color) Circle {
	return Circle{X: x, Y: y, Radius: radius, Color: color}
}

func (c Circle) Draw(painter *CanvasPainter) {
	for angle := 0; angle < 360; angle++ {
		radians := float64(angle) * math.Pi / 180
		x := c.X + c.Radius*math.Cos(radians)
		y := c.Y + c.Radius*math.Sin(radians)
		if gridX, gridY, ok := painter.GetPoint(x, y); ok {
			painter.Paint(gridX, gridY, c.Color)
		}
	}
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
	painter := newCanvasPainter(area.Width, area.Height, c.xMin, c.xMax, c.yMin, c.yMax, c.marker)
	for _, shape := range ctx.shapes {
		shape.Draw(painter)
	}
	c.renderShapes(area, buf, painter)
	for _, label := range ctx.labels {
		c.renderLabel(area, buf, label)
	}
}

func newCanvasPainter(width, height int, xMin, xMax, yMin, yMax float64, marker CanvasMarker) *CanvasPainter {
	return &CanvasPainter{
		width:  width,
		height: height,
		xMin:   xMin,
		xMax:   xMax,
		yMin:   yMin,
		yMax:   yMax,
		marker: marker,
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

	resX, resY := p.resolution()
	gridX := int(math.Round((x - p.xMin) * float64(resX-1) / (p.xMax - p.xMin)))
	gridY := int(math.Round((p.yMax - y) * float64(resY-1) / (p.yMax - p.yMin)))
	if gridX < 0 || gridY < 0 || gridX >= resX || gridY >= resY {
		return 0, 0, false
	}
	return gridX, gridY, true
}

func (p *CanvasPainter) Paint(x, y int, color style.Color) {
	if p == nil || x < 0 || y < 0 {
		return
	}
	switch p.marker.kind() {
	case CanvasMarkerBraille, CanvasMarkerQuadrant, CanvasMarkerSextant, CanvasMarkerOctant:
		cellW, cellH := p.marker.cellResolution()
		cellX, cellY := x/cellW, y/cellH
		if cellX >= p.width || cellY >= p.height {
			return
		}
		index := cellY*p.width + cellX
		pixel := p.grid[index]
		pixel.painted = true
		pixel.color = color
		pixel.pattern |= 1 << uint((x%cellW)+cellW*(y%cellH))
		p.grid[index] = pixel
	case CanvasMarkerHalfBlock:
		if x >= p.width || y >= p.height*2 {
			return
		}
		index := (y/2)*p.width + x
		pixel := p.grid[index]
		pixel.painted = true
		c := color
		if y%2 == 0 {
			pixel.upper = &c
		} else {
			pixel.lower = &c
		}
		p.grid[index] = pixel
	default:
		if x >= p.width || y >= p.height {
			return
		}
		p.grid[y*p.width+x] = canvasPixel{painted: true, color: color}
	}
}

func (p *CanvasPainter) resolution() (int, int) {
	cellW, cellH := p.marker.cellResolution()
	return p.width * cellW, p.height * cellH
}

func (m CanvasMarker) cellResolution() (int, int) {
	switch m.kind() {
	case CanvasMarkerBraille, CanvasMarkerOctant:
		return 2, 4
	case CanvasMarkerHalfBlock:
		return 1, 2
	case CanvasMarkerQuadrant:
		return 2, 2
	case CanvasMarkerSextant:
		return 2, 3
	default:
		return 1, 1
	}
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
			symbol, patch, ok := c.marker.renderPixel(pixel)
			if !ok {
				continue
			}
			cell.Symbol = symbol
			cell.Style = cell.Style.Patch(patch)
			buf.SetCell(cellX, cellY, cell)
		}
	}
}

func (m CanvasMarker) renderPixel(pixel canvasPixel) (string, style.Style, bool) {
	if !pixel.painted {
		return "", style.NewStyle(), false
	}
	switch m.kind() {
	case CanvasMarkerBlock:
		return "█", style.NewStyle().Fg(pixel.color).Bg(pixel.color), true
	case CanvasMarkerBar:
		return "▄", style.NewStyle().Fg(pixel.color), true
	case "custom":
		return m.customSymbol(), style.NewStyle().Fg(pixel.color), true
	case CanvasMarkerBraille:
		return string(brailleSymbol(pixel.pattern)), style.NewStyle().Fg(pixel.color), pixel.pattern != 0
	case CanvasMarkerHalfBlock:
		return halfBlockSymbol(pixel)
	case CanvasMarkerQuadrant:
		return string(quadrantSymbols[pixel.pattern]), style.NewStyle().Fg(pixel.color), pixel.pattern != 0
	case CanvasMarkerSextant:
		return sextantSymbol(pixel.pattern), style.NewStyle().Fg(pixel.color), pixel.pattern != 0
	case CanvasMarkerOctant:
		return octantSymbol(pixel.pattern), style.NewStyle().Fg(pixel.color), pixel.pattern != 0
	default:
		return "•", style.NewStyle().Fg(pixel.color), true
	}
}

func brailleSymbol(pattern uint8) rune {
	masks := [8]rune{0x01, 0x08, 0x02, 0x10, 0x04, 0x20, 0x40, 0x80}
	code := rune(0x2800)
	for i, mask := range masks {
		if pattern&(1<<uint(i)) != 0 {
			code += mask
		}
	}
	return code
}

func halfBlockSymbol(pixel canvasPixel) (string, style.Style, bool) {
	switch {
	case pixel.upper == nil && pixel.lower == nil:
		return "", style.NewStyle(), false
	case pixel.upper != nil && pixel.lower == nil:
		return "▀", style.NewStyle().Fg(*pixel.upper), true
	case pixel.upper == nil && pixel.lower != nil:
		return "▄", style.NewStyle().Fg(*pixel.lower), true
	case *pixel.upper == *pixel.lower:
		return "█", style.NewStyle().Fg(*pixel.upper).Bg(*pixel.lower), true
	default:
		return "▀", style.NewStyle().Fg(*pixel.upper).Bg(*pixel.lower), true
	}
}

var quadrantSymbols = [16]rune{' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛', '▗', '▚', '▐', '▜', '▄', '▙', '▟', '█'}

func sextantSymbol(pattern uint8) string {
	symbols := [64]string{
		" ", "🬀", "🬁", "🬂", "🬃", "🬄", "🬅", "🬆",
		"🬇", "🬈", "🬉", "🬊", "🬋", "🬌", "🬍", "🬎",
		"🬏", "🬐", "🬑", "🬒", "🬓", "▌", "🬔", "🬕",
		"🬖", "🬗", "🬘", "🬙", "🬚", "🬛", "🬜", "🬝",
		"🬞", "🬟", "🬠", "🬡", "🬢", "🬣", "🬤", "🬥",
		"🬦", "🬧", "▐", "🬨", "🬩", "🬪", "🬫", "🬬",
		"🬭", "🬮", "🬯", "🬰", "🬱", "🬲", "🬳", "🬴",
		"🬵", "🬶", "🬷", "🬸", "🬹", "🬺", "🬻", "█",
	}
	return symbols[pattern]
}

func octantSymbol(pattern uint8) string {
	symbols := [256]string{
		" ", "𜺨", "𜺫", "🮂", "𜴀", "▘", "𜴁", "𜴂",
		"𜴃", "𜴄", "▝", "𜴅", "𜴆", "𜴇", "𜴈", "▀",
		"𜴉", "𜴊", "𜴋", "𜴌", "🯦", "𜴍", "𜴎", "𜴏",
		"𜴐", "𜴑", "𜴒", "𜴓", "𜴔", "𜴕", "𜴖", "𜴗",
		"𜴘", "𜴙", "𜴚", "𜴛", "𜴜", "𜴝", "𜴞", "𜴟",
		"🯧", "𜴠", "𜴡", "𜴢", "𜴣", "𜴤", "𜴥", "𜴦",
		"𜴧", "𜴨", "𜴩", "𜴪", "𜴫", "𜴬", "𜴭", "𜴮",
		"𜴯", "𜴰", "𜴱", "𜴲", "𜴳", "𜴴", "𜴵", "🮅",
		"𜺣", "𜴶", "𜴷", "𜴸", "𜴹", "𜴺", "𜴻", "𜴼",
		"𜴽", "𜴾", "𜴿", "𜵀", "𜵁", "𜵂", "𜵃", "𜵄",
		"▖", "𜵅", "𜵆", "𜵇", "𜵈", "▌", "𜵉", "𜵊",
		"𜵋", "𜵌", "▞", "𜵍", "𜵎", "𜵏", "𜵐", "▛",
		"𜵑", "𜵒", "𜵓", "𜵔", "𜵕", "𜵖", "𜵗", "𜵘",
		"𜵙", "𜵚", "𜵛", "𜵜", "𜵝", "𜵞", "𜵟", "𜵠",
		"𜵡", "𜵢", "𜵣", "𜵤", "𜵥", "𜵦", "𜵧", "𜵨",
		"𜵩", "𜵪", "𜵫", "𜵬", "𜵭", "𜵮", "𜵯", "𜵰",
		"𜺠", "𜵱", "𜵲", "𜵳", "𜵴", "𜵵", "𜵶", "𜵷",
		"𜵸", "𜵹", "𜵺", "𜵻", "𜵼", "𜵽", "𜵾", "𜵿",
		"𜶀", "𜶁", "𜶂", "𜶃", "𜶄", "𜶅", "𜶆", "𜶇",
		"𜶈", "𜶉", "𜶊", "𜶋", "𜶌", "𜶍", "𜶎", "𜶏",
		"▗", "𜶐", "𜶑", "𜶒", "𜶓", "▚", "𜶔", "𜶕",
		"𜶖", "𜶗", "▐", "𜶘", "𜶙", "𜶚", "𜶛", "▜",
		"𜶜", "𜶝", "𜶞", "𜶟", "𜶠", "𜶡", "𜶢", "𜶣",
		"𜶤", "𜶥", "𜶦", "𜶧", "𜶨", "𜶩", "𜶪", "𜶫",
		"▂", "𜶬", "𜶭", "𜶮", "𜶯", "𜶰", "𜶱", "𜶲",
		"𜶳", "𜶴", "𜶵", "𜶶", "𜶷", "𜶸", "𜶹", "𜶺",
		"𜶻", "𜶼", "𜶽", "𜶾", "𜶿", "𜷀", "𜷁", "𜷂",
		"𜷃", "𜷄", "𜷅", "𜷆", "𜷇", "𜷈", "𜷉", "𜷊",
		"𜷋", "𜷌", "𜷍", "𜷎", "𜷏", "𜷐", "𜷑", "𜷒",
		"𜷓", "𜷔", "𜷕", "𜷖", "𜷗", "𜷘", "𜷙", "𜷚",
		"▄", "𜷛", "𜷜", "𜷝", "𜷞", "▙", "𜷟", "𜷠",
		"𜷡", "𜷢", "▟", "𜷣", "▆", "𜷤", "𜷥", "█",
	}
	return symbols[pattern]
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
