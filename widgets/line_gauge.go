package widgets

import (
	"fmt"
	"math"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/textbuffer"
)

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
	if x >= right {
		return x
	}
	label := g.effectiveLineGaugeLabel()
	endX, _ := textbuffer.SetLine(buf, x, y, label, right-x)
	return endX
}

func (g LineGauge) effectiveLineGaugeLabel() text.Line {
	if g.label != nil {
		return *g.label
	}
	return text.LineFromString(fmt.Sprintf("%3.0f%%", g.ratio*100))
}
