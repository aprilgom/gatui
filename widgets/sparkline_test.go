package widgets_test

import (
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

func TestSparkline_shouldRenderRightToLeft(t *testing.T) {
	buf := sparklineBuffer(12, 1)
	sparkline := widgets.NewSparkline().
		Data([]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}).
		Direction(widgets.RenderDirectionRightToLeft)

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"xxx█▇▆▅▄▃▂▁ "})
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
