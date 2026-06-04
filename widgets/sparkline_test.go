package widgets_test

import (
	"math"
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/widgets"
)

func TestRenderDirection_String_shouldMatchRatatui(t *testing.T) {
	tests := []struct {
		direction widgets.RenderDirection
		want      string
	}{
		{direction: widgets.RenderDirectionLeftToRight, want: "LeftToRight"},
		{direction: widgets.RenderDirectionRightToLeft, want: "RightToLeft"},
	}

	for _, tt := range tests {
		if got := tt.direction.String(); got != tt.want {
			t.Fatalf("%#v.String() = %q, want %q", tt.direction, got, tt.want)
		}
	}
}

func TestParseRenderDirection_shouldMatchRatatui(t *testing.T) {
	tests := []struct {
		value string
		want  widgets.RenderDirection
	}{
		{value: "LeftToRight", want: widgets.RenderDirectionLeftToRight},
		{value: "RightToLeft", want: widgets.RenderDirectionRightToLeft},
	}

	for _, tt := range tests {
		got, err := widgets.ParseRenderDirection(tt.value)
		if err != nil {
			t.Fatalf("ParseRenderDirection(%q) returned error: %v", tt.value, err)
		}
		if got != tt.want {
			t.Fatalf("ParseRenderDirection(%q) = %#v, want %#v", tt.value, got, tt.want)
		}
	}
}

func TestRenderDirection_String_unknownShouldBeStable(t *testing.T) {
	got := widgets.RenderDirection(99).String()
	if got != "RenderDirection(99)" {
		t.Fatalf("RenderDirection(99).String() = %q, want %q", got, "RenderDirection(99)")
	}
}

func TestParseRenderDirection_unknownShouldReturnError(t *testing.T) {
	for _, value := range []string{"", "lefttoright", "Forward"} {
		if got, err := widgets.ParseRenderDirection(value); err == nil {
			t.Fatalf("ParseRenderDirection(%q) = %#v, want error", value, got)
		}
	}
}

func TestSparkline_canBeStylized(t *testing.T) {
	buf := sparklineBuffer(1, 1)
	sparkline := widgets.NewSparkline().
		Data([]uint64{1}).
		Fg(style.Red).
		Bg(style.Blue).
		Bold().
		Dim().
		Italic().
		Cyan()

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"█"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().
		Fg(style.Cyan).
		Bg(style.Blue).
		AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
}

func TestSparkline_canBeCreatedFromSliceOfUint64(t *testing.T) {
	buf := sparklineBuffer(4, 1)
	data := []uint64{1, 2, 3}
	sparkline := widgets.NewSparkline().Data(data)
	data[1] = 99

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"▂▅█x"})
}

func TestSparkline_canBeCreatedFromArrayOfUint64(t *testing.T) {
	buf := sparklineBuffer(4, 1)
	data := []uint64{0, 1, 2}
	sparkline := widgets.NewSparkline().Data(data)
	data[1] = 9

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{" ▄█x"})
}

func TestSparkline_canBeCreatedFromVecOfUint64(t *testing.T) {
	buf := sparklineBuffer(4, 1)
	data := []uint64{0, 1, 2}
	sparkline := widgets.NewSparkline().Data(append([]uint64(nil), data...))
	data = append(data, 9)
	if got, want := len(data), 4; got != want {
		t.Fatalf("len(data) = %d, want %d", got, want)
	}

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{" ▄█x"})
}

func TestSparkline_canBeCreatedFromSliceOfOptionalUint64(t *testing.T) {
	buf := sparklineBuffer(4, 1)
	bars := []widgets.SparklineBar{
		widgets.NewSparklineBar(1),
		widgets.NewAbsentSparklineBar(),
		widgets.NewSparklineBar(3),
	}
	sparkline := widgets.NewSparkline().AbsentValueSymbol("*").Bars(bars)
	bars[1] = widgets.NewSparklineBar(2)

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{"▂*█x"})
}

func TestSparkline_canBeCreatedFromArrayOfOptionalUint64(t *testing.T) {
	buf := sparklineBuffer(4, 1)
	bars := []widgets.SparklineBar{
		widgets.NewSparklineBar(0),
		widgets.NewAbsentSparklineBar(),
		widgets.NewSparklineBar(2),
	}
	sparkline := widgets.NewSparkline().AbsentValueSymbol("*").Bars(bars)
	bars[1] = widgets.NewSparklineBar(1)

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{" *█x"})
}

func TestSparkline_canBeCreatedFromVecOfOptionalUint64(t *testing.T) {
	buf := sparklineBuffer(4, 1)
	bars := []widgets.SparklineBar{
		widgets.NewSparklineBar(0),
		widgets.NewAbsentSparklineBar(),
		widgets.NewSparklineBar(2),
	}
	sparkline := widgets.NewSparkline().AbsentValueSymbol("*").Bars(append([]widgets.SparklineBar(nil), bars...))
	bars = append(bars, widgets.NewSparklineBar(9))
	if got, want := len(bars), 4; got != want {
		t.Fatalf("len(bars) = %d, want %d", got, want)
	}

	sparkline.Render(buf.Area, buf)

	assertLines(t, buf, []string{" *█x"})
}

func TestSparkline_renderInMinimalBuffer(t *testing.T) {
	t.Run("one_by_one", func(t *testing.T) {
		buf := sparklineBuffer(1, 1)
		sparkline := widgets.NewSparkline().Data([]uint64{1, 2, 3})

		assertNotPanics(t, func() {
			sparkline.Render(buf.Area, buf)
		})

		assertLines(t, buf, []string{"▂"})
	})

	t.Run("block_with_empty_inner_area", func(t *testing.T) {
		buf := sparklineBuffer(1, 1)
		sparkline := widgets.NewSparkline().
			Block(widgets.NewBlock().Borders(widgets.AllBorders)).
			Data([]uint64{1})

		assertNotPanics(t, func() {
			sparkline.Render(buf.Area, buf)
		})

		assertLines(t, buf, []string{"┌"})
	})
}

func TestSparkline_renderInZeroSizeBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))
	sparkline := widgets.NewSparkline().Data([]uint64{1})

	assertNotPanics(t, func() {
		sparkline.Render(buf.Area, buf)
		sparkline.Render(layout.NewRect(0, 0, 1, 0), buf)
		sparkline.Render(layout.NewRect(0, 0, 0, 1), buf)
	})
}

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
