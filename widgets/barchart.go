package widgets

import (
	"strconv"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
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
	label string
	bars  []Bar
}

func NewBarGroup(bars []Bar) BarGroup {
	return BarGroup{bars: append([]Bar(nil), bars...)}
}

func (g BarGroup) Label(label string) BarGroup {
	g.label = label
	return g
}

func (g BarGroup) Bars(bars []Bar) BarGroup {
	g.bars = append([]Bar(nil), bars...)
	return g
}

type BarChart struct {
	groups     []BarGroup
	block      *Block
	max        uint64
	barWidth   int
	barGap     int
	groupGap   int
	barStyle   style.Style
	valueStyle style.Style
	labelStyle style.Style
}

func NewBarChart() BarChart {
	return BarChart{
		barWidth:   1,
		barGap:     1,
		groupGap:   1,
		barStyle:   style.NewStyle(),
		valueStyle: style.NewStyle(),
		labelStyle: style.NewStyle(),
	}
}

func (c BarChart) Data(group BarGroup) BarChart {
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
	labelRows := 1
	if c.hasGroupLabels() {
		labelRows = 2
	}
	if chartArea.Height <= labelRows {
		return
	}
	max := c.effectiveMax()
	barHeight := chartArea.Height - labelRows
	barLabelY := chartArea.Y + barHeight
	groupLabelY := barLabelY
	if labelRows == 2 {
		groupLabelY = barLabelY + 1
	}
	x := chartArea.X
	right := chartArea.X + chartArea.Width
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
			c.renderBar(buf, x, chartArea.Y, width, barHeight, max, bar)
			c.renderCentered(buf, x, barLabelY, width, bar.label, c.labelStyle.Patch(bar.labelStyle))
			x += c.barWidth
		}
		if group.label != "" && groupStart < right {
			writeStringWithin(buf, groupStart, groupLabelY, right, group.label, c.labelStyle)
		}
	}
}

func (c BarChart) renderBar(buf *buffer.Buffer, x, y, width, height int, max uint64, bar Bar) {
	if width <= 0 || height <= 0 || max == 0 || bar.value == 0 {
		return
	}
	eighths := int((bar.value * uint64(height) * 8) / max)
	if eighths > height*8 {
		eighths = height * 8
	}
	barStyle := c.barStyle.Patch(bar.style)
	buf.SetStyle(layout.NewRect(x, y, width, height), barStyle)
	for rowFromBottom := 0; rowFromBottom < height; rowFromBottom++ {
		rowEighths := eighths - rowFromBottom*8
		if rowEighths <= 0 {
			continue
		}
		symbol := "█"
		if rowEighths < 8 {
			symbol = partialBarSymbol(rowEighths)
		}
		for dx := 0; dx < width; dx++ {
			buf.SetCell(x+dx, y+height-1-rowFromBottom, buffer.Cell{Symbol: symbol, Style: barStyle})
		}
	}
	value := bar.textValue
	if value == "" {
		value = uintToString(bar.value)
	}
	valueStyle := barStyle.Patch(c.valueStyle).Patch(bar.valueStyle)
	c.renderCentered(buf, x, y+height-1, width, value, valueStyle)
}

func (c BarChart) renderCentered(buf *buffer.Buffer, x, y, width int, value string, cellStyle style.Style) {
	runes := []rune(value)
	if len(runes) > width {
		runes = runes[:width]
	}
	offset := (width - len(runes)) / 2
	for i, r := range runes {
		buf.SetCell(x+offset+i, y, buffer.Cell{Symbol: string(r), Style: cellStyle})
	}
}

func (c BarChart) hasGroupLabels() bool {
	for _, group := range c.groups {
		if group.label != "" {
			return true
		}
	}
	return false
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

func partialBarSymbol(eighths int) string {
	switch eighths {
	case 1:
		return "▁"
	case 2:
		return "▂"
	case 3:
		return "▃"
	case 4:
		return "▄"
	case 5:
		return "▅"
	case 6:
		return "▆"
	case 7:
		return "▇"
	default:
		return " "
	}
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
