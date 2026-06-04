package widgets_test

import (
	"math"
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
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

func TestBarChart_default(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))

	widgets.NewBarChart().Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"          ",
		"          ",
		"          ",
	})
}

func TestBarChart_constructorsIgnoreEmptyGroups(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 3))
	barchart := widgets.NewBarChartWithBars(nil).
		Data(widgets.NewBarGroup(nil).Label("empty")).
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(1).Label("A"),
			widgets.NewBar(2).Label("B"),
		}))

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"  █     ",
		"1 2     ",
		"A B     ",
	})
}

func TestBarChart_block(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 5))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "foo", Value: 1}, {Label: "bar", Value: 2}}).
		Block(widgets.BorderedBlock().Title(text.LineFromString("Block")))

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌Block───┐",
		"│  █     │",
		"│1 2     │",
		"│f b     │",
		"└────────┘",
	})
}

func TestBarChart_canBeStylized(t *testing.T) {
	got := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}}).
		Fg(style.Black).
		Bg(style.White).
		Bold()
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))
	got.Render(buf.Area, buf)

	if gotStyle := buf.Cells[0].Style; gotStyle != style.NewStyle().Fg(style.Black).Bg(style.White).AddModifier(style.ModifierBold) {
		t.Fatalf("style = %#v", gotStyle)
	}
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

func TestBarChart_shouldRenderDataWithDefaultNineLevelBarSet(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 5))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{
			{Label: "A", Value: 1},
			{Label: "B", Value: 2},
			{Label: "C", Value: 3},
		}).
		Max(4).
		BarGap(1)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"        ",
		"    █   ",
		"  █ █   ",
		"1 2 3   ",
		"A B C   ",
	})
}

func TestBarChart_shouldUseConfiguredBarWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 5))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}, {Label: "B", Value: 2}}).
		Max(2).
		BarWidth(2).
		BarGap(1)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"   ██   ",
		"   ██   ",
		"██ ██   ",
		"1█ 2█   ",
		"A  B    ",
	})
}

func TestBarChart_shouldUseConfiguredBarGap(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 4))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}, {Label: "B", Value: 2}}).
		Max(2).
		BarGap(2)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"   █    ",
		"▄  █    ",
		"1  2    ",
		"A  B    ",
	})
}

func TestBarChart_shouldUseConfiguredBarSet(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 5))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{
			{Label: "A", Value: 1},
			{Label: "B", Value: 2},
			{Label: "C", Value: 3},
		}).
		Max(4).
		BarGap(1).
		BarSet(widgets.ThreeLevelBarSet)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"        ",
		"    █   ",
		"  █ █   ",
		"1 2 3   ",
		"A B C   ",
	})
}

func TestBarChart_shouldUseNineLevelBarSetPreset(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 2))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{
			{Label: "0", Value: 0},
			{Label: "1", Value: 1},
			{Label: "2", Value: 2},
			{Label: "3", Value: 3},
			{Label: "4", Value: 4},
		}).
		Max(8).
		BarGap(1).
		BarSet(widgets.NineLevelBarSet)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"  ▁ ▂ ▃ ▄ ",
		"0 1 2 3 4 ",
	})
}

func TestBarChart_shouldApplyValueStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}}).
		ValueStyle(style.NewStyle().Fg(style.Red))

	barchart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 1, style.NewStyle().Fg(style.Red))
}

func TestBarChart_shouldApplyLabelStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}}).
		LabelStyle(style.NewStyle().Fg(style.Blue))

	barchart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 2, style.NewStyle().Fg(style.Blue))
}

func TestBarChart_shouldApplyChartAreaStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 3))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}}).
		Style(style.NewStyle().Bg(style.Yellow))

	barchart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 3, 0, style.NewStyle().Bg(style.Yellow))
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Bg(style.Yellow))
	assertCellStyle(t, buf, 0, 1, style.NewStyle().Bg(style.Yellow))
	assertCellStyle(t, buf, 0, 2, style.NewStyle().Bg(style.Yellow))
}

func TestBarChart_shouldIgnoreEmptyGroups(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 3))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup(nil).Label("empty")).
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}})

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"█       ",
		"1       ",
		"A       ",
	})
}

func TestBarChart_shouldRenderSingleLine(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}})

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{"█   "})
}

func TestBarChart_shouldRenderTwoLines(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 2))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}})

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"█   ",
		"A   ",
	})
}

func TestBarChart_shouldRenderThreeLines(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 3))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}})

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"█   ",
		"1   ",
		"A   ",
	})
}

func TestBarChart_shouldRenderFourLines(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 4))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: 1}})

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"█   ",
		"█   ",
		"1   ",
		"A   ",
	})
}

func TestBarChart_threeLinesDoubleWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 26, 3))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.BarWithLabel("a", 0),
			widgets.BarWithLabel("b", 1),
			widgets.BarWithLabel("c", 2),
			widgets.BarWithLabel("d", 3),
			widgets.BarWithLabel("e", 4),
			widgets.BarWithLabel("f", 5),
			widgets.BarWithLabel("g", 6),
			widgets.BarWithLabel("h", 7),
			widgets.BarWithLabel("i", 8),
		}).Label("Group").LabelAlignment(layout.Center)).
		BarWidth(2).
		BarSet(widgets.NineLevelBarSet)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"   1▁ 2▂ 3▃ 4▄ 5▅ 6▆ 7▇ 8█",
		"a  b  c  d  e  f  g  h  i ",
		"          Group           ",
	})
}

func TestBarChart_twoLinesWithoutBarLabels(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 17, 3))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(0),
			widgets.NewBar(1),
			widgets.NewBar(2),
			widgets.NewBar(3),
			widgets.NewBar(4),
			widgets.NewBar(5),
			widgets.NewBar(6),
			widgets.NewBar(7),
			widgets.NewBar(8),
		}).Label("Group").LabelAlignment(layout.Center))

	barchart.Render(layout.NewRect(0, 1, buf.Area.Width, 2), buf)

	assertLines(t, buf, []string{
		"                 ",
		"  ▁ ▂ ▃ ▄ ▅ ▆ ▇ 8",
		"      Group      ",
	})
}

func TestBarChart_oneLineWithMoreBars(t *testing.T) {
	bars := make([]widgets.Bar, 0, 30)
	for i := range uint64(30) {
		bars = append(bars, widgets.NewBar(i))
	}
	buf := buffer.Empty(layout.NewRect(0, 0, 59, 1))
	barchart := widgets.NewBarChart().Data(widgets.NewBarGroup(bars))

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{"        ▁ ▁ ▁ ▁ ▂ ▂ ▂ ▃ ▃ ▃ ▃ ▄ ▄ ▄ ▄ ▅ ▅ ▅ ▆ ▆ ▆ ▆ ▇ ▇ ▇ █"})
}

func TestBarChart_firstBarOfTheGroupIsHalfOutsideView(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 6))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "a", Value: 1}, {Label: "b", Value: 2}}).
		DataPairs([]widgets.BarData{{Label: "a", Value: 1}, {Label: "b", Value: 2}}).
		BarWidth(2)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"   ██  ",
		"   ██  ",
		"▄▄ ██  ",
		"██ ██  ",
		"1█ 2█  ",
		"a  b   ",
	})
}

func TestBarChart_shouldRenderUint64MaxValue(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{{Label: "A", Value: math.MaxUint64}})

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"█  ",
		"1  ",
		"A  ",
	})
}

func TestBarChart_shouldKeepIntegerPrecisionForLargeValues(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 10))
	barchart := widgets.NewBarChart().
		DataPairs([]widgets.BarData{
			{Label: "A", Value: 9007199254740992},
			{Label: "B", Value: 9007199254740993},
		}).
		Max(9007199254740993).
		BarGap(1)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"▇ █  ",
		"█ █  ",
		"█ █  ",
		"█ █  ",
		"█ █  ",
		"█ █  ",
		"█ █  ",
		"█ █  ",
		"9 9  ",
		"A B  ",
	})
}

func TestBarChart_shouldRenderHorizontalBars(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 8))
	barchart := buildHorizontalTestBarChart()

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"2█   ",
		"3██  ",
		"4███ ",
		"G1   ",
		"3██  ",
		"4███ ",
		"5████",
		"G2   ",
	})
}

func TestBarChart_shouldRenderHorizontalBarsWithoutGroupLabel_whenHeightIsShort(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 7))
	barchart := buildHorizontalTestBarChart()

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"2█   ",
		"3██  ",
		"4███ ",
		"G1   ",
		"3██  ",
		"4███ ",
		"5████",
	})
}

func TestBarChart_shouldRenderOnlyVisibleHorizontalBars_whenHeightIsVeryShort(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))
	barchart := buildHorizontalTestBarChart()

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"2█   ",
		"3██  ",
		"4███ ",
		"G1   ",
		"3██  ",
	})
}

func TestBarChart_shouldKeepHorizontalValueStyleWithinActualBar_whenValueTextExceedsBarWithoutBarStyle(t *testing.T) {
	assertHorizontalValueTextExceedsBar(t, style.Default, false)
}

func TestBarChart_shouldKeepHorizontalValueStyleWithinActualBar_whenValueTextExceedsBarWithBarStyle(t *testing.T) {
	assertHorizontalValueTextExceedsBar(t, style.White, true)
}

func TestBarChart_shouldRenderHorizontalBarLabels(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	barchart := widgets.NewBarChart().
		Direction(layout.Horizontal).
		BarGap(0).
		DataPairs([]widgets.BarData{
			{Label: "Jan", Value: 10},
			{Label: "Feb", Value: 20},
			{Label: "Mar", Value: 5},
		})

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"Jan 10█   ",
		"Feb 20████",
		"Mar 5     ",
	})
}

func TestBarChart_shouldRenderHorizontalMultibyteValueTextWithoutPanic(t *testing.T) {
	textValue := "\u202f"
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 5))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(0).TextValue(textValue),
			widgets.NewBar(1).TextValue(textValue),
			widgets.NewBar(2).TextValue(textValue),
			widgets.NewBar(3).TextValue(textValue),
			widgets.NewBar(4).TextValue(textValue),
		})).
		BarGap(0).
		Direction(layout.Horizontal)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"\u202f   ",
		"\u202f   ",
		"\u202f█  ",
		"\u202f██ ",
		"\u202f███",
	})
}

func TestBarChart_groupLabelStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(2),
		}).Label("G1").LabelStyle(style.NewStyle().Fg(style.Red))).
		GroupGap(1).
		Direction(layout.Horizontal).
		LabelStyle(style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierBold))

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{"2████", "G1   "})
	assertCellStyle(t, buf, 0, 1, style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold))
	assertCellStyle(t, buf, 1, 1, style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold))
}

func TestBarChart_groupLabelCenter(t *testing.T) {
	group := widgets.NewBarGroup([]widgets.Bar{
		widgets.BarWithLabel("a", 1),
		widgets.BarWithLabel("b", 2),
		widgets.BarWithLabel("c", 3),
		widgets.BarWithLabel("c", 4),
	})
	buf := buffer.Empty(layout.NewRect(0, 0, 13, 5))
	barchart := widgets.NewBarChart().
		Data(group.Label("G1").LabelAlignment(layout.Center)).
		Data(group.Label("G2").LabelAlignment(layout.Center)).
		GroupGap(0)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"    ▂ █     ▂",
		"  ▄ █ █   ▄ █",
		"▆ 2 3 4 ▆ 2 3",
		"a b c c a b c",
		"  G1     G2  ",
	})
}

func TestBarChart_groupLabelRight(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(2),
			widgets.NewBar(5),
		}).Label("G").LabelAlignment(layout.Right))

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"  █",
		"▆ 5",
		"  G",
	})
}

func TestBarChart_unicodeAsValue(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 5))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(123).Label("B1").TextValue("写"),
			widgets.NewBar(321).Label("B2").TextValue("写"),
			widgets.NewBar(333).Label("B2").TextValue("写"),
		})).
		BarWidth(3).
		BarGap(1)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"    ▆▆▆ ███",
		"    ███ ███",
		"▃▃▃ ███ ███",
		"写█ 写█ 写█",
		"B1  B2  B2 ",
	})
}

func TestBarChart_newWithBars(t *testing.T) {
	bars := []widgets.Bar{widgets.BarWithLabel("Red", 1), widgets.BarWithLabel("Green", 2)}
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 3))
	widgets.NewBarChartWithBars(bars).Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"  █  ",
		"1 2  ",
		"R G  ",
	})
}

func TestBar_new(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 2))
	widgets.NewBarChartWithBars([]widgets.Bar{widgets.NewBar(7)}).Render(buf.Area, buf)

	assertLines(t, buf, []string{"█", "7"})
}

func TestBar_stylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 2))
	widgets.NewBarChartWithBars([]widgets.Bar{
		widgets.NewBar(7).Style(style.NewStyle().Fg(style.Red)),
	}).Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red))
}

func TestBar_withLabel(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 2))
	widgets.NewBarChartWithBars([]widgets.Bar{widgets.BarWithLabel("A", 7)}).Render(buf.Area, buf)

	assertLines(t, buf, []string{"█", "A"})
}

func TestBarGroup_new(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))
	widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(1),
			widgets.NewBar(2),
		}).Label("G")).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"  █",
		"1 2",
		"G  ",
	})
}

func buildHorizontalTestBarChart() widgets.BarChart {
	return widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(2),
			widgets.NewBar(3),
			widgets.NewBar(4),
		}).Label("G1")).
		Data(widgets.NewBarGroup([]widgets.Bar{
			widgets.NewBar(3),
			widgets.NewBar(4),
			widgets.NewBar(5),
		}).Label("G2")).
		GroupGap(1).
		Direction(layout.Horizontal).
		BarGap(0)
}

func assertHorizontalValueTextExceedsBar(t *testing.T, barColor style.Color, hasBarColor bool) {
	t.Helper()
	bar := widgets.NewBar(2).
		TextValue("label").
		ValueStyle(style.NewStyle().Fg(style.Red))
	if hasBarColor {
		bar = bar.Style(style.NewStyle().Fg(barColor))
	}
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))
	barchart := widgets.NewBarChart().
		Data(widgets.NewBarGroup([]widgets.Bar{
			bar,
			widgets.NewBar(5),
		})).
		Direction(layout.Horizontal).
		BarStyle(style.NewStyle().Fg(style.Yellow)).
		ValueStyle(style.NewStyle().AddModifier(style.ModifierItalic)).
		BarGap(0)

	barchart.Render(buf.Area, buf)

	assertLines(t, buf, []string{"label", "5████"})
	expectedBarColor := style.Yellow
	if hasBarColor {
		expectedBarColor = barColor
	}
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red).AddModifier(style.ModifierItalic))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).AddModifier(style.ModifierItalic))
	for x := 2; x < 5; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().Fg(expectedBarColor))
	}
	assertCellStyle(t, buf, 0, 1, style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierItalic))
	for x := 1; x < 5; x++ {
		assertCellStyle(t, buf, x, 1, style.NewStyle().Fg(style.Yellow))
	}
}
