package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/widgets"
)

func TestBarChart_shouldRenderBarsBelowMax(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 30, 10))
	barchart := widgets.NewBarChart().
		Block(widgets.BorderedBlock()).
		DataPairs([]widgets.BarData{
			{Label: "empty", Value: 0},
			{Label: "half", Value: 50},
			{Label: "almost", Value: 99},
			{Label: "full", Value: 100},
		}).
		Max(100).
		BarWidth(7).
		BarGap(0)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌────────────────────────────┐",
		"│              ▇▇▇▇▇▇▇███████│",
		"│              ██████████████│",
		"│              ██████████████│",
		"│       ▄▄▄▄▄▄▄██████████████│",
		"│       █████████████████████│",
		"│       █████████████████████│",
		"│       ██50█████99█████100██│",
		"│ empty  half  almost  full  │",
		"└────────────────────────────┘",
	})
}

func TestBarChart_shouldRenderGroupedBars(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 35, 11))
	barchart := widgets.NewBarChart().
		Block(widgets.BorderedBlock()).
		Data(
			widgets.NewBarGroup([]widgets.Bar{
				widgets.NewBar(10).
					Label("C1").
					Style(style.NewStyle().Fg(style.Red)).
					ValueStyle(style.NewStyle().Fg(style.Blue)),
				widgets.NewBar(20).
					Style(style.NewStyle().Fg(style.Green)).
					TextValue("20M"),
			}).Label("Mar"),
		).
		DataPairs([]widgets.BarData{{Label: "C1", Value: 50}, {Label: "C2", Value: 40}}).
		DataPairs([]widgets.BarData{{Label: "C1", Value: 60}, {Label: "C2", Value: 90}}).
		DataPairs([]widgets.BarData{{Label: "xx", Value: 10}, {Label: "xx", Value: 10}}).
		GroupGap(2).
		BarWidth(4).
		BarGap(1)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌─────────────────────────────────┐",
		"│                             ████│",
		"│                             ████│",
		"│                        ▅▅▅▅ ████│",
		"│            ▇▇▇▇        ████ ████│",
		"│            ████ ████   ████ ████│",
		"│     ▄▄▄▄   ████ ████   ████ ████│",
		"│▆10▆ 20M█   █50█ █40█   █60█ █90█│",
		"│ C1          C1   C2     C1   C2 │",
		"│Mar                              │",
		"└─────────────────────────────────┘",
	})
	for y := 1; y < 8; y++ {
		for x := 1; x < 5; x++ {
			if y != 7 || (x != 2 && x != 3) {
				assertCellStyle(t, buf, x, y, style.NewStyle().Fg(style.Red))
			}
			assertCellStyle(t, buf, x+5, y, style.NewStyle().Fg(style.Green))
		}
	}
	assertCellStyle(t, buf, 2, 7, style.NewStyle().Fg(style.Blue))
	assertCellStyle(t, buf, 3, 7, style.NewStyle().Fg(style.Blue))
}

func TestBarChart_shouldRenderEmptyChartWithoutPanic(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 3))

	widgets.NewBarChart().Block(widgets.BorderedBlock()).Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌──┐",
		"│  │",
		"└──┘",
	})
}

func TestBarChart_shouldRenderLabelsWithoutDivideByZero_whenMaxIsZero(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 4))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "zero", Value: 0}}).
		Max(0).
		BarWidth(4)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"        ",
		"        ",
		"        ",
		"zero    ",
	})
}

func TestBarChart_shouldLetPerBarStylesOverrideChartStyles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 4))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(10).
				Label("A").
				Style(style.NewStyle().Fg(style.Red)).
				ValueStyle(style.NewStyle().Fg(style.Blue)),
		})).
		Max(10).
		BarWidth(3).
		BarStyle(style.NewStyle().Fg(style.Green)).
		ValueStyle(style.NewStyle().Fg(style.Yellow))

	barchart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 1, style.NewStyle().Fg(style.Red))
	assertCellStyle(t, buf, 1, 2, style.NewStyle().Fg(style.Blue))
}
