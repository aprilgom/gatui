package widgets

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type Widget interface {
	Render(area layout.Rect, buf *buffer.Buffer)
}

type WidgetRef interface {
	RenderRef(area layout.Rect, buf *buffer.Buffer)
}

type Wrap struct {
	Trim bool
}

type Paragraph struct {
	text      text.Text
	wrap      *Wrap
	block     *Block
	alignment layout.Alignment
	scrollY   int
	scrollX   int
}

type Tabs struct {
	titles         []text.Line
	block          *Block
	style          style.Style
	highlightStyle style.Style
	selected       *int
	divider        string
	paddingLeft    string
	paddingRight   string
}

type Gauge struct {
	block      *Block
	ratio      float64
	label      *text.Span
	style      style.Style
	gaugeStyle style.Style
	useUnicode bool
}

type LineGauge struct {
	block          *Block
	ratio          float64
	label          *text.Line
	style          style.Style
	filledStyle    style.Style
	unfilledStyle  style.Style
	filledSymbol   string
	unfilledSymbol string
}

func NewGauge() Gauge {
	return Gauge{style: style.NewStyle(), gaugeStyle: style.NewStyle()}
}

func (g Gauge) Block(block Block) Gauge {
	g.block = &block
	return g
}

func (g Gauge) Percent(percent int) Gauge {
	if percent < 0 || percent > 100 {
		panic("gauge percent must be between 0 and 100")
	}
	g.ratio = float64(percent) / 100
	return g
}

func (g Gauge) Ratio(ratio float64) Gauge {
	if ratio < 0 || ratio > 1 {
		panic("gauge ratio must be between 0 and 1")
	}
	g.ratio = ratio
	return g
}

func (g Gauge) Label(label text.Span) Gauge {
	g.label = &label
	return g
}

func (g Gauge) LabelString(label string) Gauge {
	return g.Label(text.NewSpan(label))
}

func (g Gauge) Style(gaugeStyle style.Style) Gauge {
	g.style = gaugeStyle
	return g
}

func (g Gauge) GaugeStyle(gaugeStyle style.Style) Gauge {
	g.gaugeStyle = gaugeStyle
	return g
}

func (g Gauge) UseUnicode(useUnicode bool) Gauge {
	g.useUnicode = useUnicode
	return g
}

func (g Gauge) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, g.style)
	gaugeArea := area
	if g.block != nil {
		g.block.Render(area, buf)
		gaugeArea = g.block.Inner(area)
	}
	g.renderGauge(gaugeArea, buf)
}

func (g Gauge) renderGauge(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, g.gaugeStyle)

	label := g.effectiveLabel()
	labelRunes := []rune(label.Content)
	if len(labelRunes) > area.Width {
		labelRunes = labelRunes[:area.Width]
	}
	labelX := area.X + (area.Width-len(labelRunes))/2
	labelY := area.Y + area.Height/2

	filledWidth := float64(area.Width) * g.ratio
	end := area.X + int(math.Round(filledWidth))
	if g.useUnicode {
		end = area.X + int(math.Floor(filledWidth))
	}
	if end > area.X+area.Width {
		end = area.X + area.Width
	}

	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < end; x++ {
			if y == labelY && x >= labelX && x < labelX+len(labelRunes) {
				buf.SetCell(x, y, buffer.Cell{Symbol: " ", Style: g.swappedGaugeStyle()})
				continue
			}
			buf.SetCell(x, y, buffer.Cell{Symbol: "█", Style: g.gaugeStyle})
		}
		if g.useUnicode && g.ratio < 1 && end < area.X+area.Width {
			buf.SetCell(end, y, buffer.Cell{Symbol: unicodeBlock(filledWidth - math.Floor(filledWidth)), Style: g.gaugeStyle})
		}
	}

	for i, r := range labelRunes {
		x := labelX + i
		cellStyle := g.gaugeStyle.Patch(label.Style)
		if x < end {
			cellStyle = g.swappedGaugeStyle().Patch(label.Style)
		}
		buf.SetCell(x, labelY, buffer.Cell{Symbol: string(r), Style: cellStyle})
	}
}

func (g Gauge) effectiveLabel() text.Span {
	if g.label != nil {
		return *g.label
	}
	return text.NewSpan(strconv.Itoa(int(math.Round(g.ratio*100))) + "%")
}

func (g Gauge) swappedGaugeStyle() style.Style {
	swapped := g.gaugeStyle
	swapped.Foreground = g.gaugeStyle.Background
	swapped.Background = g.gaugeStyle.Foreground
	return swapped
}

func unicodeBlock(frac float64) string {
	switch int(math.Round(frac * 8)) {
	case 1:
		return "▏"
	case 2:
		return "▎"
	case 3:
		return "▍"
	case 4:
		return "▌"
	case 5:
		return "▋"
	case 6:
		return "▊"
	case 7:
		return "▉"
	case 8:
		return "█"
	default:
		return " "
	}
}

func NewLineGauge() LineGauge {
	return LineGauge{
		style:          style.NewStyle(),
		filledStyle:    style.NewStyle(),
		unfilledStyle:  style.NewStyle(),
		filledSymbol:   "─",
		unfilledSymbol: "─",
	}
}

func (g LineGauge) Block(block Block) LineGauge {
	g.block = &block
	return g
}

func (g LineGauge) Ratio(ratio float64) LineGauge {
	if ratio < 0 || ratio > 1 {
		panic("line gauge ratio must be between 0 and 1")
	}
	g.ratio = ratio
	return g
}

func (g LineGauge) Label(label text.Line) LineGauge {
	g.label = &label
	return g
}

func (g LineGauge) LabelString(label string) LineGauge {
	return g.Label(text.LineFromString(label))
}

func (g LineGauge) Style(lineGaugeStyle style.Style) LineGauge {
	g.style = lineGaugeStyle
	return g
}

func (g LineGauge) FilledStyle(filledStyle style.Style) LineGauge {
	g.filledStyle = filledStyle
	return g
}

func (g LineGauge) UnfilledStyle(unfilledStyle style.Style) LineGauge {
	g.unfilledStyle = unfilledStyle
	return g
}

func (g LineGauge) FilledSymbol(symbol string) LineGauge {
	g.filledSymbol = symbol
	return g
}

func (g LineGauge) UnfilledSymbol(symbol string) LineGauge {
	g.unfilledSymbol = symbol
	return g
}

func (g LineGauge) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, g.style)
	gaugeArea := area
	if g.block != nil {
		g.block.Render(area, buf)
		gaugeArea = g.block.Inner(area)
	}
	g.renderLineGauge(gaugeArea, buf)
}

func (g LineGauge) renderLineGauge(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	right := area.X + area.Width
	x := g.writeLabel(area.X, area.Y, right, buf)
	start := x + 1
	if start >= right {
		return
	}
	remainingWidth := right - start
	filledWidth := int(math.Floor(float64(remainingWidth) * g.ratio))
	filledEnd := start + filledWidth
	for col := start; col < filledEnd; col++ {
		buf.SetCell(col, area.Y, buffer.Cell{Symbol: g.filledSymbol, Style: g.filledStyle})
	}
	for col := filledEnd; col < right; col++ {
		buf.SetCell(col, area.Y, buffer.Cell{Symbol: g.unfilledSymbol, Style: g.unfilledStyle})
	}
}

func (g LineGauge) writeLabel(x, y, right int, buf *buffer.Buffer) int {
	for _, cell := range cellsFromLine(g.effectiveLineGaugeLabel()) {
		if x >= right {
			return x
		}
		buf.SetCell(x, y, cell)
		x++
	}
	return x
}

func (g LineGauge) effectiveLineGaugeLabel() text.Line {
	if g.label != nil {
		return *g.label
	}
	return text.LineFromString(fmt.Sprintf("%3.0f%%", g.ratio*100))
}

func NewTabs(titles []text.Line) Tabs {
	selected := 0
	return Tabs{
		titles:         append([]text.Line(nil), titles...),
		style:          style.NewStyle(),
		highlightStyle: style.NewStyle().AddModifier(style.ModifierReversed),
		selected:       &selected,
		divider:        "│",
		paddingLeft:    " ",
		paddingRight:   " ",
	}
}

func TabsFromStrings(titles []string) Tabs {
	lines := make([]text.Line, 0, len(titles))
	for _, title := range titles {
		lines = append(lines, text.LineFromString(title))
	}
	return NewTabs(lines)
}

func (t Tabs) Block(block Block) Tabs {
	t.block = &block
	return t
}

func (t Tabs) Select(index int) Tabs {
	t.selected = &index
	return t
}

func (t Tabs) ClearSelection() Tabs {
	t.selected = nil
	return t
}

func (t Tabs) Style(tabStyle style.Style) Tabs {
	t.style = tabStyle
	return t
}

func (t Tabs) HighlightStyle(highlightStyle style.Style) Tabs {
	t.highlightStyle = highlightStyle
	return t
}

func (t Tabs) Divider(divider string) Tabs {
	t.divider = divider
	return t
}

func (t Tabs) Padding(left, right string) Tabs {
	t.paddingLeft = left
	t.paddingRight = right
	return t
}

func (t Tabs) PaddingLeft(left string) Tabs {
	t.paddingLeft = left
	return t
}

func (t Tabs) PaddingRight(right string) Tabs {
	t.paddingRight = right
	return t
}

func (t Tabs) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, t.style)
	tabsArea := area
	if t.block != nil {
		t.block.Render(area, buf)
		tabsArea = t.block.Inner(area)
	}
	if tabsArea.Width == 0 || tabsArea.Height == 0 {
		return
	}
	x := tabsArea.X
	right := tabsArea.X + tabsArea.Width
	for index, title := range t.titles {
		if x >= right {
			return
		}
		x = writeStringWithin(buf, x, tabsArea.Y, right, t.paddingLeft, t.style)
		selected := t.selected != nil && *t.selected == index
		x = t.renderTitle(buf, title, x, tabsArea.Y, right, selected)
		x = writeStringWithin(buf, x, tabsArea.Y, right, t.paddingRight, t.style)
		if index < len(t.titles)-1 {
			x = writeStringWithin(buf, x, tabsArea.Y, right, t.divider, t.style)
		}
	}
}

func (t Tabs) renderTitle(buf *buffer.Buffer, title text.Line, x, y, right int, selected bool) int {
	for _, span := range title.Spans {
		cellStyle := t.style.Patch(span.Style)
		if selected {
			cellStyle = cellStyle.Patch(t.highlightStyle)
		}
		for _, r := range span.Content {
			if x >= right {
				return x
			}
			buf.SetCell(x, y, buffer.Cell{Symbol: string(r), Style: cellStyle})
			x++
		}
	}
	return x
}

func NewParagraph(content text.Text) Paragraph {
	return Paragraph{text: content, alignment: layout.Left}
}

func (p Paragraph) Wrap(wrap Wrap) Paragraph {
	p.wrap = &wrap
	return p
}

func (p Paragraph) Block(block Block) Paragraph {
	p.block = &block
	return p
}

func (p Paragraph) Alignment(alignment layout.Alignment) Paragraph {
	p.alignment = alignment
	return p
}

func (p Paragraph) Scroll(y, x int) Paragraph {
	p.scrollY = y
	p.scrollX = x
	return p
}

func (p Paragraph) Fg(color style.Color) Paragraph {
	p.text = p.text.Fg(color)
	return p
}

func (p Paragraph) Bg(color style.Color) Paragraph {
	p.text = p.text.Bg(color)
	return p
}

func (p Paragraph) Bold() Paragraph {
	p.text = p.text.Bold()
	return p
}

func (p Paragraph) Italic() Paragraph {
	p.text = p.text.Italic()
	return p
}

func (p Paragraph) Cyan() Paragraph {
	return p.Fg(style.Cyan)
}

func (p Paragraph) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	textArea := area
	if p.block != nil {
		p.block.Render(area, buf)
		textArea = p.block.Inner(area)
	}
	lines := p.renderLines(textArea.Width)
	if p.scrollY < len(lines) {
		lines = lines[p.scrollY:]
	} else {
		lines = nil
	}
	for y := 0; y < textArea.Height && y < len(lines); y++ {
		line := lines[y]
		if p.scrollX > 0 && line.alignment == layout.Left {
			line = line.skip(p.scrollX)
		}
		offset := paragraphLineOffset(line.width(), textArea.Width, line.alignment)
		x := textArea.X + offset
		for _, cell := range line.cells {
			if x >= textArea.X+textArea.Width {
				break
			}
			buf.SetCell(x, textArea.Y+y, cell)
			x++
		}
	}
}

func paragraphLineOffset(lineWidth, areaWidth int, alignment layout.Alignment) int {
	if lineWidth >= areaWidth {
		return 0
	}
	switch alignment {
	case layout.Center:
		return areaWidth/2 - lineWidth/2
	case layout.Right:
		return areaWidth - lineWidth
	default:
		return 0
	}
}

func (p Paragraph) renderLines(width int) []renderLine {
	if width <= 0 {
		return nil
	}
	var lines []renderLine
	for _, line := range p.text.Lines {
		alignment := p.alignment
		if line.Alignment != nil {
			alignment = *line.Alignment
		}
		cells := cellsFromLine(line)
		if p.wrap == nil {
			lines = append(lines, renderLine{cells: append([]buffer.Cell(nil), cells...), alignment: alignment})
			continue
		}
		for _, wrapped := range wrapCells(cells, width, p.wrap.Trim) {
			lines = append(lines, renderLine{cells: wrapped, alignment: alignment})
		}
	}
	return lines
}

type Clear struct{}

func (Clear) Render(area layout.Rect, buf *buffer.Buffer) {
	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < area.X+area.Width; x++ {
			buf.SetCell(x, y, buffer.Cell{Symbol: " ", Style: style.NewStyle()})
		}
	}
}

type renderLine struct {
	cells     []buffer.Cell
	alignment layout.Alignment
}

func (l renderLine) width() int {
	return len(l.cells)
}

func (l renderLine) skip(count int) renderLine {
	if count >= len(l.cells) {
		return renderLine{alignment: l.alignment}
	}
	l.cells = l.cells[count:]
	return l
}

func wrapCells(cells []buffer.Cell, width int, trim bool) [][]buffer.Cell {
	var lines [][]buffer.Cell
	for len(cells) > 0 {
		if trim {
			cells = trimLeftCells(cells)
		}
		if len(cells) <= width {
			lines = append(lines, trimRightCells(append([]buffer.Cell(nil), cells...), trim))
			break
		}
		breakAt := width
		for i := width; i >= 0; i-- {
			if i < len(cells) && isSpaceCell(cells[i]) {
				breakAt = i
				break
			}
		}
		if breakAt == 0 {
			breakAt = width
		}
		line := append([]buffer.Cell(nil), cells[:breakAt]...)
		lines = append(lines, trimRightCells(line, trim))
		cells = cells[breakAt:]
	}
	if len(lines) == 0 {
		lines = append(lines, nil)
	}
	return lines
}

func trimLeftCells(cells []buffer.Cell) []buffer.Cell {
	for len(cells) > 0 && isSpaceCell(cells[0]) {
		cells = cells[1:]
	}
	return cells
}

func trimRightCells(cells []buffer.Cell, trim bool) []buffer.Cell {
	if !trim {
		return cells
	}
	for len(cells) > 0 && isSpaceCell(cells[len(cells)-1]) {
		cells = cells[:len(cells)-1]
	}
	return cells
}

func isSpaceCell(cell buffer.Cell) bool {
	return strings.TrimFunc(cell.Symbol, unicode.IsSpace) == ""
}
