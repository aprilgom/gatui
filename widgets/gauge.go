package widgets

import (
	"math"
	"strconv"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type Gauge struct {
	block      *Block
	ratio      float64
	label      *text.Span
	style      style.Style
	gaugeStyle style.Style
	useUnicode bool
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

func (g Gauge) Fg(color style.Color) Gauge {
	g.style = g.style.Fg(color)
	return g
}

func (g Gauge) Bg(color style.Color) Gauge {
	g.style = g.style.Bg(color)
	return g
}

func (g Gauge) Bold() Gauge {
	g.style = g.style.AddModifier(style.ModifierBold)
	return g
}

func (g Gauge) Dim() Gauge {
	g.style = g.style.AddModifier(style.ModifierDim)
	return g
}

func (g Gauge) Italic() Gauge {
	g.style = g.style.AddModifier(style.ModifierItalic)
	return g
}

func (g Gauge) Cyan() Gauge {
	return g.Fg(style.Cyan)
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
		cellStyle := g.style.Patch(label.Style)
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
