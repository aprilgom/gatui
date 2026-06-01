package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

type RenderDirection int

const (
	RenderDirectionLeftToRight RenderDirection = iota
	RenderDirectionRightToLeft
)

type SparklineBar struct {
	value uint64
	style style.Style
}

func NewSparklineBar(value uint64) SparklineBar {
	return SparklineBar{value: value, style: style.NewStyle()}
}

func (b SparklineBar) Style(barStyle style.Style) SparklineBar {
	b.style = barStyle
	return b
}

type Sparkline struct {
	data      []SparklineBar
	max       *uint64
	block     *Block
	style     style.Style
	direction RenderDirection
}

func NewSparkline() Sparkline {
	return Sparkline{
		style:     style.NewStyle(),
		direction: RenderDirectionLeftToRight,
	}
}

func (s Sparkline) Data(values []uint64) Sparkline {
	bars := make([]SparklineBar, 0, len(values))
	for _, value := range values {
		bars = append(bars, NewSparklineBar(value))
	}
	s.data = bars
	return s
}

func (s Sparkline) Bars(values []SparklineBar) Sparkline {
	s.data = append([]SparklineBar(nil), values...)
	return s
}

func (s Sparkline) Max(max uint64) Sparkline {
	s.max = &max
	return s
}

func (s Sparkline) Block(block Block) Sparkline {
	s.block = &block
	return s
}

func (s Sparkline) Style(sparklineStyle style.Style) Sparkline {
	s.style = sparklineStyle
	return s
}

func (s Sparkline) Direction(direction RenderDirection) Sparkline {
	s.direction = direction
	return s
}

func (s Sparkline) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	sparkArea := area
	if s.block != nil {
		s.block.Render(area, buf)
		sparkArea = s.block.Inner(area)
	}
	s.renderSparkline(sparkArea, buf)
}

func (s Sparkline) renderSparkline(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 || len(s.data) == 0 {
		return
	}
	maxValue := s.effectiveMax()
	maxIndex := minInt(area.Width, len(s.data))
	for i := 0; i < maxIndex; i++ {
		bar := s.data[i]
		x := area.X + i
		if s.direction == RenderDirectionRightToLeft {
			x = area.X + area.Width - i - 1
		}
		height := scaleSparklineHeight(bar.value, maxValue, area.Height)
		cellStyle := s.style.Patch(bar.style)
		for rowFromBottom := 0; rowFromBottom < area.Height; rowFromBottom++ {
			symbol := sparklineSymbolForHeight(height)
			if height > 8 {
				height -= 8
			} else {
				height = 0
			}
			y := area.Y + area.Height - 1 - rowFromBottom
			buf.SetCell(x, y, buffer.Cell{Symbol: symbol, Style: cellStyle})
		}
	}
}

func (s Sparkline) effectiveMax() uint64 {
	if s.max != nil {
		return *s.max
	}
	var maxValue uint64
	for _, bar := range s.data {
		if bar.value > maxValue {
			maxValue = bar.value
		}
	}
	return maxValue
}

func scaleSparklineHeight(value, maxValue uint64, maxHeight int) uint64 {
	if maxValue == 0 || maxHeight <= 0 {
		return 0
	}
	maxTicks := uint64(maxHeight * 8)
	ticks := value * maxTicks / maxValue
	if ticks > maxTicks {
		return maxTicks
	}
	return ticks
}

func sparklineSymbolForHeight(height uint64) string {
	switch height {
	case 0:
		return " "
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
		return "█"
	}
}
