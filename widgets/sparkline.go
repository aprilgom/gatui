package widgets

import (
	"fmt"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/symbols"
	"math/bits"
)

type RenderDirection int

const (
	RenderDirectionLeftToRight RenderDirection = iota
	RenderDirectionRightToLeft
)

func (d RenderDirection) String() string {
	switch d {
	case RenderDirectionLeftToRight:
		return "LeftToRight"
	case RenderDirectionRightToLeft:
		return "RightToLeft"
	default:
		return fmt.Sprintf("RenderDirection(%d)", d)
	}
}

func ParseRenderDirection(value string) (RenderDirection, error) {
	switch value {
	case "LeftToRight":
		return RenderDirectionLeftToRight, nil
	case "RightToLeft":
		return RenderDirectionRightToLeft, nil
	default:
		return RenderDirectionLeftToRight, fmt.Errorf("unknown RenderDirection %q", value)
	}
}

type SparklineBar struct {
	value   uint64
	present bool
	style   style.Style
}

func NewSparklineBar(value uint64) SparklineBar {
	return SparklineBar{value: value, present: true, style: style.NewStyle()}
}

func NewAbsentSparklineBar() SparklineBar {
	return SparklineBar{style: style.NewStyle()}
}

func (b SparklineBar) Style(barStyle style.Style) SparklineBar {
	b.style = barStyle
	return b
}

type Sparkline struct {
	data              []SparklineBar
	max               *uint64
	block             *Block
	style             style.Style
	absentValueStyle  style.Style
	absentValueSymbol string
	barSet            SparklineBarSet
	direction         RenderDirection
}

type SparklineBarSet = symbols.SparklineBarSet

func NineLevelSparklineBarSet() SparklineBarSet {
	return symbols.NineLevelSparklineBarSet()
}

func ThreeLevelSparklineBarSet() SparklineBarSet {
	return symbols.ThreeLevelSparklineBarSet()
}

func NewSparkline() Sparkline {
	return Sparkline{
		style:             style.NewStyle(),
		absentValueStyle:  style.NewStyle(),
		absentValueSymbol: " ",
		barSet:            NineLevelSparklineBarSet(),
		direction:         RenderDirectionLeftToRight,
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

func (s Sparkline) Fg(color style.Color) Sparkline {
	s.style = s.style.Fg(color)
	return s
}

func (s Sparkline) Bg(color style.Color) Sparkline {
	s.style = s.style.Bg(color)
	return s
}

func (s Sparkline) Bold() Sparkline {
	s.style = s.style.AddModifier(style.ModifierBold)
	return s
}

func (s Sparkline) Dim() Sparkline {
	s.style = s.style.AddModifier(style.ModifierDim)
	return s
}

func (s Sparkline) Italic() Sparkline {
	s.style = s.style.AddModifier(style.ModifierItalic)
	return s
}

func (s Sparkline) Cyan() Sparkline {
	return s.Fg(style.Cyan)
}

func (s Sparkline) AbsentValueStyle(absentValueStyle style.Style) Sparkline {
	s.absentValueStyle = absentValueStyle
	return s
}

func (s Sparkline) AbsentValueSymbol(symbol string) Sparkline {
	s.absentValueSymbol = symbol
	return s
}

func (s Sparkline) BarSet(barSet SparklineBarSet) Sparkline {
	s.barSet = barSet
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
	for i := range maxIndex {
		bar := s.data[i]
		x := area.X + i
		if s.direction == RenderDirectionRightToLeft {
			x = area.X + area.Width - i - 1
		}
		if !bar.present {
			cellStyle := s.style.Patch(s.absentValueStyle)
			for rowFromBottom := 0; rowFromBottom < area.Height; rowFromBottom++ {
				y := area.Y + area.Height - 1 - rowFromBottom
				buf.SetCell(x, y, buffer.Cell{Symbol: s.absentValueSymbol, Style: cellStyle})
			}
			continue
		}
		height := scaleSparklineHeight(bar.value, maxValue, area.Height)
		cellStyle := s.style.Patch(bar.style)
		for rowFromBottom := 0; rowFromBottom < area.Height; rowFromBottom++ {
			symbol := s.barSet.SymbolForHeight(height)
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
		if bar.present && bar.value > maxValue {
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
	if value >= maxValue {
		return maxTicks
	}
	hi, lo := bits.Mul64(value, maxTicks)
	ticks, _ := bits.Div64(hi, lo, maxValue)
	if ticks > maxTicks {
		return maxTicks
	}
	return ticks
}
