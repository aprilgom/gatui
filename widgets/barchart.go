package widgets

import (
	"strconv"
	"unicode/utf8"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/symbols"
)

type BarData struct {
	Label string
	Value uint64
}

type Bar struct {
	label      string
	value      uint64
	textValue  string
	style      style.Style
	valueStyle style.Style
	labelStyle style.Style
}

func NewBar(value uint64) Bar {
	return Bar{value: value, textValue: uintToString(value), style: style.NewStyle(), valueStyle: style.NewStyle(), labelStyle: style.NewStyle()}
}

func BarWithLabel(label string, value uint64) Bar {
	return NewBar(value).Label(label)
}

func (b Bar) Label(label string) Bar {
	b.label = label
	return b
}

func (b Bar) TextValue(value string) Bar {
	b.textValue = value
	return b
}

func (b Bar) Style(barStyle style.Style) Bar {
	b.style = barStyle
	return b
}

func (b Bar) ValueStyle(valueStyle style.Style) Bar {
	b.valueStyle = valueStyle
	return b
}

func (b Bar) LabelStyle(labelStyle style.Style) Bar {
	b.labelStyle = labelStyle
	return b
}

type BarGroup struct {
	label          string
	labelStyle     style.Style
	labelAlignment layout.Alignment
	bars           []Bar
}

func NewBarGroup(bars []Bar) BarGroup {
	return BarGroup{
		labelStyle:     style.NewStyle(),
		labelAlignment: layout.Left,
		bars:           append([]Bar(nil), bars...),
	}
}

func (g BarGroup) Label(label string) BarGroup {
	g.label = label
	return g
}

func (g BarGroup) Bars(bars []Bar) BarGroup {
	g.bars = append([]Bar(nil), bars...)
	return g
}

func (g BarGroup) LabelStyle(labelStyle style.Style) BarGroup {
	g.labelStyle = labelStyle
	return g
}

func (g BarGroup) LabelAlignment(alignment layout.Alignment) BarGroup {
	g.labelAlignment = alignment
	return g
}

type BarSet = symbols.BarSet

var NineLevelBarSet = symbols.NineLevelBarSet

var ThreeLevelBarSet = symbols.ThreeLevelBarSet

type BarChart struct {
	groups     []BarGroup
	block      *Block
	direction  layout.Direction
	max        uint64
	barWidth   int
	barGap     int
	groupGap   int
	style      style.Style
	barStyle   style.Style
	valueStyle style.Style
	labelStyle style.Style
	barSet     BarSet
}

func NewBarChart() BarChart {
	return BarChart{
		direction:  layout.Vertical,
		barWidth:   1,
		barGap:     1,
		groupGap:   1,
		style:      style.NewStyle(),
		barStyle:   style.NewStyle(),
		valueStyle: style.NewStyle(),
		labelStyle: style.NewStyle(),
		barSet:     NineLevelBarSet,
	}
}

func NewBarChartWithBars(bars []Bar) BarChart {
	return NewBarChart().Data(NewBarGroup(bars))
}

func (c BarChart) Data(group BarGroup) BarChart {
	if len(group.bars) == 0 {
		return c
	}
	c.groups = append(c.groups, group)
	return c
}

func (c BarChart) DataPairs(data []BarData) BarChart {
	bars := make([]Bar, 0, len(data))
	for _, item := range data {
		bars = append(bars, BarWithLabel(item.Label, item.Value))
	}
	return c.Data(NewBarGroup(bars))
}

func (c BarChart) Block(block Block) BarChart {
	c.block = &block
	return c
}

func (c BarChart) Direction(direction layout.Direction) BarChart {
	c.direction = direction
	return c
}

func (c BarChart) Max(max uint64) BarChart {
	c.max = max
	return c
}

func (c BarChart) BarWidth(width int) BarChart {
	c.barWidth = width
	return c
}

func (c BarChart) BarGap(gap int) BarChart {
	c.barGap = gap
	return c
}

func (c BarChart) GroupGap(gap int) BarChart {
	c.groupGap = gap
	return c
}

func (c BarChart) Style(chartStyle style.Style) BarChart {
	c.style = chartStyle
	return c
}

func (c BarChart) BarStyle(barStyle style.Style) BarChart {
	c.barStyle = barStyle
	return c
}

func (c BarChart) ValueStyle(valueStyle style.Style) BarChart {
	c.valueStyle = valueStyle
	return c
}

func (c BarChart) LabelStyle(labelStyle style.Style) BarChart {
	c.labelStyle = labelStyle
	return c
}

func (c BarChart) BarSet(barSet BarSet) BarChart {
	c.barSet = barSet
	return c
}

func (c BarChart) Fg(color style.Color) BarChart {
	c.style = c.style.Fg(color)
	return c
}

func (c BarChart) Bg(color style.Color) BarChart {
	c.style = c.style.Bg(color)
	return c
}

func (c BarChart) Bold() BarChart {
	c.style = c.style.AddModifier(style.ModifierBold)
	return c
}

func (c BarChart) Italic() BarChart {
	c.style = c.style.AddModifier(style.ModifierItalic)
	return c
}

func (c BarChart) Cyan() BarChart {
	return c.Fg(style.Cyan)
}

func (c BarChart) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	chartArea := area
	if c.block != nil {
		c.block.Render(area, buf)
		chartArea = c.block.Inner(area)
	}
	if chartArea.Width == 0 || chartArea.Height == 0 || len(c.groups) == 0 || c.barWidth <= 0 {
		return
	}
	buf.SetStyle(chartArea, c.style)
	if c.direction == layout.Horizontal {
		c.renderHorizontal(chartArea, buf)
		return
	}
	c.renderVertical(chartArea, buf)
}

func (c BarChart) renderVertical(chartArea layout.Rect, buf *buffer.Buffer) {
	hasBarLabels := c.hasBarLabels()
	hasGroupLabels := c.hasGroupLabels()
	labelRows := 0
	if chartArea.Height >= 2 && hasBarLabels {
		labelRows = 1
	}
	if hasGroupLabels && (chartArea.Height >= 3 || !hasBarLabels && chartArea.Height >= 2) {
		labelRows++
	}
	max := c.effectiveMax()
	barHeight := chartArea.Height - labelRows
	barLabelY := chartArea.Y + barHeight
	groupLabelY := barLabelY
	if hasBarLabels && hasGroupLabels {
		groupLabelY = barLabelY + 1
	}
	x := chartArea.X
	right := chartArea.X + chartArea.Width
	totalBars := c.totalBars()
	for groupIndex, group := range c.groups {
		if groupIndex > 0 {
			x += nonNegative(c.groupGap) + nonNegative(c.barGap)
		}
		groupStart := x
		for barIndex, bar := range group.bars {
			if barIndex > 0 {
				x += nonNegative(c.barGap)
			}
			if x >= right {
				break
			}
			width := c.barWidth
			if x+width > right {
				width = right - x
			}
			c.renderBar(buf, x, chartArea.Y, width, barHeight, max, bar, totalBars, hasGroupLabels)
			if hasBarLabels {
				c.renderCentered(buf, x, barLabelY, width, bar.label, c.style.Patch(c.labelStyle).Patch(bar.labelStyle))
			}
			x += c.barWidth
		}
		if group.label != "" && groupStart < right {
			groupEnd := minInt(x, right)
			if groupEnd > groupStart {
				c.renderAligned(buf, groupStart, groupLabelY, groupEnd-groupStart, group.label, group.labelAlignment, c.style.Patch(c.labelStyle).Patch(group.labelStyle))
			}
		}
	}
}

func (c BarChart) renderHorizontal(chartArea layout.Rect, buf *buffer.Buffer) {
	max := c.effectiveMax()
	labelWidth := c.horizontalLabelWidth(chartArea)
	barAreaWidth := chartArea.Width - labelWidth
	if barAreaWidth <= 0 {
		return
	}
	y := chartArea.Y
	bottom := chartArea.Y + chartArea.Height
	for groupIndex, group := range c.groups {
		if groupIndex > 0 {
			y += nonNegative(c.groupGap)
		}
		for barIndex, bar := range group.bars {
			if barIndex > 0 {
				y += nonNegative(c.barGap)
			}
			if y >= bottom {
				return
			}
			height := c.barWidth
			if y+height > bottom {
				height = bottom - y
			}
			if labelWidth > 0 && bar.label != "" {
				writeStringWithin(buf, chartArea.X, y+height/2, chartArea.X+labelWidth-1, bar.label, c.style.Patch(c.labelStyle).Patch(bar.labelStyle))
			}
			c.renderHorizontalBar(buf, chartArea.X+labelWidth, y, barAreaWidth, height, max, bar)
			y += c.barWidth
		}
		if c.groupGap > 0 && group.label != "" && y < bottom {
			c.renderAligned(buf, chartArea.X, y, chartArea.Width, group.label, group.labelAlignment, c.style.Patch(c.labelStyle).Patch(group.labelStyle))
		}
	}
}

func (c BarChart) horizontalLabelWidth(area layout.Rect) int {
	maxLabelWidth := 0
	for _, group := range c.groups {
		for _, bar := range group.bars {
			if width := buffer.CellWidth(bar.label); width > maxLabelWidth {
				maxLabelWidth = width
			}
		}
	}
	if maxLabelWidth == 0 {
		return 0
	}
	width := maxLabelWidth + 1
	if width > area.Width {
		return area.Width
	}
	return width
}

func (c BarChart) renderHorizontalBar(buf *buffer.Buffer, x, y, width, height int, max uint64, bar Bar) {
	if width <= 0 || height <= 0 {
		return
	}
	barLength := scaledTicks(bar.value, max, width)
	barStyle := c.style.Patch(c.barStyle).Patch(bar.style)
	if barLength > 0 {
		buf.SetStyle(layout.NewRect(x, y, barLength, height), barStyle)
		for dy := range height {
			for dx := range barLength {
				buf.SetCell(x+dx, y+dy, buffer.Cell{Symbol: c.barSet.Full, Style: barStyle})
			}
		}
	}
	value := bar.textValue
	if value == "" {
		value = uintToString(bar.value)
	}
	valueY := y + height/2
	c.renderHorizontalValue(buf, x, valueY, width, barLength, value, barStyle.Patch(c.valueStyle).Patch(bar.valueStyle), barStyle)
}

func (c BarChart) renderHorizontalValue(buf *buffer.Buffer, x, y, width, barLength int, value string, valueStyle, overflowStyle style.Style) {
	if width <= 0 || value == "" {
		return
	}
	if barLength > 0 {
		endX, _ := buf.SetStringN(x, y, value, minInt(width, barLength), valueStyle)
		remaining := width - (endX - x)
		if remaining > 0 {
			buf.SetStringN(endX, y, valueDisplaySuffix(value, endX-x), remaining, overflowStyle)
		}
		return
	}
	buf.SetStringN(x, y, value, width, overflowStyle)
}

func valueDisplaySuffix(value string, skippedWidth int) string {
	if skippedWidth <= 0 {
		return value
	}
	for index, r := range value {
		width := buffer.CellWidth(string(r))
		if width > skippedWidth {
			return value[index:]
		}
		skippedWidth -= width
		if skippedWidth == 0 {
			return value[index+utf8.RuneLen(r):]
		}
	}
	return ""
}

func (c BarChart) renderBar(buf *buffer.Buffer, x, y, width, height int, max uint64, bar Bar, totalBars int, hasGroupLabels bool) {
	if width <= 0 || height <= 0 || max == 0 || bar.value == 0 {
		return
	}
	totalEighths := height * 8
	eighths := scaledTicks(bar.value, max, totalEighths)
	barStyle := c.style.Patch(c.barStyle).Patch(bar.style)
	buf.SetStyle(layout.NewRect(x, y, width, height), barStyle)
	for rowFromBottom := range height {
		rowEighths := eighths - rowFromBottom*8
		if rowEighths <= 0 {
			continue
		}
		symbol := c.barSet.Full
		if rowEighths < 8 {
			symbol = c.partialBarSymbol(rowEighths)
		}
		for dx := range width {
			buf.SetCell(x+dx, y+height-1-rowFromBottom, buffer.Cell{Symbol: symbol, Style: barStyle})
		}
	}
	value := bar.textValue
	if value == "" {
		value = uintToString(bar.value)
	}
	if width == 1 && hasGroupLabels && eighths < 8 {
		return
	}
	if height == 1 && !c.shouldRenderValueInSingleLine(width, totalBars, max, bar.value, value) {
		return
	}
	valueStyle := barStyle.Patch(c.valueStyle).Patch(bar.valueStyle)
	c.renderCentered(buf, x, y+height-1, width, value, valueStyle)
}

func (c BarChart) renderCentered(buf *buffer.Buffer, x, y, width int, value string, cellStyle style.Style) {
	c.renderAligned(buf, x, y, width, value, layout.Center, cellStyle)
}

func (c BarChart) renderAligned(buf *buffer.Buffer, x, y, width int, value string, alignment layout.Alignment, cellStyle style.Style) {
	if width <= 0 || value == "" {
		return
	}
	valueWidth := buffer.CellWidth(value)
	offset := alignedOffset(valueWidth, width, alignment)
	buf.SetStringN(x+offset, y, value, width-offset, cellStyle)
}

func (c BarChart) hasGroupLabels() bool {
	for _, group := range c.groups {
		if group.label != "" {
			return true
		}
	}
	return false
}

func (c BarChart) hasBarLabels() bool {
	for _, group := range c.groups {
		for _, bar := range group.bars {
			if bar.label != "" {
				return true
			}
		}
	}
	return false
}

func (c BarChart) totalBars() int {
	total := 0
	for _, group := range c.groups {
		total += len(group.bars)
	}
	return total
}

func (c BarChart) shouldRenderValueInSingleLine(width, totalBars int, max, value uint64, textValue string) bool {
	return width > 1 || (totalBars > 1 && value >= max && buffer.CellWidth(textValue) <= width)
}

func (c BarChart) effectiveMax() uint64 {
	if c.max > 0 {
		return c.max
	}
	var max uint64
	for _, group := range c.groups {
		for _, bar := range group.bars {
			if bar.value > max {
				max = bar.value
			}
		}
	}
	return max
}

func (c BarChart) partialBarSymbol(eighths int) string {
	switch eighths {
	case 1:
		return c.barSet.OneEighth
	case 2:
		return c.barSet.OneQuarter
	case 3:
		return c.barSet.ThreeEighths
	case 4:
		return c.barSet.Half
	case 5:
		return c.barSet.FiveEighths
	case 6:
		return c.barSet.ThreeQuarters
	case 7:
		return c.barSet.SevenEighths
	default:
		return c.barSet.Empty
	}
}

func scaledTicks(value, max uint64, total int) int {
	if total <= 0 || max == 0 || value == 0 {
		return 0
	}
	if value >= max {
		return total
	}
	var ticks int
	var remainder uint64
	for range total {
		next := remainder + value
		if next < remainder || next >= max {
			ticks++
			next -= max
		}
		remainder = next
	}
	return ticks
}

func uintToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func nonNegative(value int) int {
	if value < 0 {
		return 0
	}
	return value
}
