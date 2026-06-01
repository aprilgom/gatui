package widgets_test

import (
	"math"
	"testing"

	"gatui/buffer"
	"gatui/style"
	"gatui/widgets"
)

func TestSparkline_shouldNotPanicWhenMaxIsZero(t *testing.T) {
	buf := sparklineBuffer(6, 1)
	sparkline := widgets.NewSparkline().Data([]uint64{0, 0, 0})

	assertNotPanics(t, func() {
		sparkline.Render(buf.Area, buf)
	})

	assertLines(t, buf, []string{"   xxx"})
}

func TestSparkline_shouldNotPanicWhenMaxIsExplicitZero(t *testing.T) {
	buf := sparklineBuffer(6, 1)
	sparkline := widgets.NewSparkline().Data([]uint64{0, 1, 2}).Max(0)

	assertNotPanics(t, func() {
		sparkline.Render(buf.Area, buf)
	})

	assertLines(t, buf, []string{"   xxx"})
}

func TestSparkline_shouldDrawSingleRow(t *testing.T) {
	buf := sparklineBuffer(12, 1)
	sparkline := widgets.NewSparkline().Data([]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8})

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{" ▁▂▃▄▅▆▇█xxx"})
}

func TestSparkline_shouldDrawDoubleHeight(t *testing.T) {
	buf := sparklineBuffer(12, 2)
	sparkline := widgets.NewSparkline().Data([]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8})

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"     ▂▄▆█xxx",
		" ▂▄▆█████xxx",
	})
}

func TestSparkline_shouldHandleUint64MaxValue(t *testing.T) {
	buf := sparklineBuffer(1, 3)
	sparkline := widgets.NewSparkline().
		Data([]uint64{math.MaxUint64}).
		Max(math.MaxUint64)

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"█",
		"█",
		"█",
	})
}

func TestSparkline_shouldKeepIntegerPrecisionForLargeValues(t *testing.T) {
	buf := sparklineBuffer(1, 1)
	sparkline := widgets.NewSparkline().
		Data([]uint64{math.MaxUint64 - 1}).
		Max(math.MaxUint64)

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"▇"})
}

func TestSparkline_shouldRenderRightToLeft(t *testing.T) {
	buf := sparklineBuffer(12, 1)
	sparkline := widgets.NewSparkline().
		Data([]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}).
		Direction(widgets.RenderDirectionRightToLeft)

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"xxx█▇▆▅▄▃▂▁ "})
}

func TestSparkline_shouldRenderAbsentValueStyle(t *testing.T) {
	buf := sparklineBuffer(12, 1)
	sparkline := widgets.NewSparkline().
		AbsentValueStyle(style.NewStyle().Fg(style.Red)).
		AbsentValueSymbol("█").
		Bars([]widgets.SparklineBar{
			widgets.NewAbsentSparklineBar(),
			widgets.NewSparklineBar(1),
			widgets.NewSparklineBar(2),
			widgets.NewSparklineBar(3),
			widgets.NewSparklineBar(4),
			widgets.NewSparklineBar(5),
			widgets.NewSparklineBar(6),
			widgets.NewSparklineBar(7),
			widgets.NewSparklineBar(8),
		})

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"█▁▂▃▄▅▆▇█xxx"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red))
}

func TestSparkline_shouldRenderAbsentValueStyleDoubleHeight(t *testing.T) {
	buf := sparklineBuffer(12, 2)
	sparkline := widgets.NewSparkline().
		AbsentValueStyle(style.NewStyle().Fg(style.Red)).
		AbsentValueSymbol("█").
		Bars([]widgets.SparklineBar{
			widgets.NewAbsentSparklineBar(),
			widgets.NewSparklineBar(1),
			widgets.NewSparklineBar(2),
			widgets.NewSparklineBar(3),
			widgets.NewSparklineBar(4),
			widgets.NewSparklineBar(5),
			widgets.NewSparklineBar(6),
			widgets.NewSparklineBar(7),
			widgets.NewSparklineBar(8),
		})

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"█    ▂▄▆█xxx",
		"█▂▄▆█████xxx",
	})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red))
	assertCellStyle(t, buf, 0, 1, style.NewStyle().Fg(style.Red))
}

func TestSparkline_shouldRenderCustomAbsentValueSymbol(t *testing.T) {
	buf := sparklineBuffer(12, 1)
	sparkline := widgets.NewSparkline().
		AbsentValueSymbol("*").
		Bars([]widgets.SparklineBar{
			widgets.NewAbsentSparklineBar(),
			widgets.NewSparklineBar(1),
			widgets.NewSparklineBar(2),
			widgets.NewSparklineBar(3),
			widgets.NewSparklineBar(4),
			widgets.NewSparklineBar(5),
			widgets.NewSparklineBar(6),
			widgets.NewSparklineBar(7),
			widgets.NewSparklineBar(8),
		})

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"*▁▂▃▄▅▆▇█xxx"})
}

func TestSparkline_shouldRenderThreeLevelBarSet(t *testing.T) {
	buf := sparklineBuffer(12, 1)
	sparkline := widgets.NewSparkline().
		BarSet(widgets.ThreeLevelSparklineBarSet()).
		Data([]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8})

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"  ▄▄▄▄▄██xxx"})
}

func TestSparkline_shouldApplyBarStyles(t *testing.T) {
	buf := sparklineBuffer(3, 1)
	sparkline := widgets.NewSparkline().
		Style(style.NewStyle().Bg(style.Yellow)).
		Bars([]widgets.SparklineBar{
			widgets.NewSparklineBar(1).Style(style.NewStyle().Fg(style.Red)),
			widgets.NewSparklineBar(2).Style(style.NewStyle().Fg(style.Green)),
			widgets.NewSparklineBar(3).Style(style.NewStyle().Fg(style.Blue)),
		}).
		Max(3)

	sparkline.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 0, style.NewStyle().Bg(style.Yellow).Fg(style.Red))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Bg(style.Yellow).Fg(style.Green))
	assertCellStyle(t, buf, 2, 0, style.NewStyle().Bg(style.Yellow).Fg(style.Blue))
}

func sparklineBuffer(width, height int) *buffer.Buffer {
	lines := make([]string, height)
	for i := range lines {
		line := make([]rune, width)
		for j := range line {
			line[j] = 'x'
		}
		lines[i] = string(line)
	}
	return buffer.WithLines(lines)
}
