package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestGauge_shouldRenderPercentAndRatioWithUnicode(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 40, 6))
	gaugeStyle := style.NewStyle().Bg(style.Blue).Fg(style.Red)

	widgets.NewGauge().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Percentage"))).
		GaugeStyle(gaugeStyle).
		UseUnicode(true).
		Percent(43).
		Render(layout.NewRect(2, 0, 36, 3), buf)
	widgets.NewGauge().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Ratio"))).
		GaugeStyle(gaugeStyle).
		UseUnicode(true).
		Ratio(0.511_313_934_313_1).
		Render(layout.NewRect(2, 3, 36, 3), buf)

	assertLines(t, buf, []string{
		"  ┌Percentage────────────────────────┐  ",
		"  │██████████████▋43%                │  ",
		"  └──────────────────────────────────┘  ",
		"  ┌Ratio─────────────────────────────┐  ",
		"  │███████████████51%                │  ",
		"  └──────────────────────────────────┘  ",
	})
	for x := 3; x < 37; x++ {
		assertCellStyle(t, buf, x, 1, gaugeStyle)
	}
	for x := 3; x < 18; x++ {
		assertCellStyle(t, buf, x, 4, gaugeStyle)
	}
	for x := 18; x < 20; x++ {
		assertCellStyle(t, buf, x, 4, style.NewStyle().Fg(style.Blue).Bg(style.Red))
	}
	for x := 20; x < 37; x++ {
		assertCellStyle(t, buf, x, 4, gaugeStyle)
	}
}

func TestGauge_shouldRenderWithoutUnicode(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 40, 6))

	widgets.NewGauge().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Percentage"))).
		Percent(43).
		UseUnicode(false).
		Render(layout.NewRect(2, 0, 36, 3), buf)
	widgets.NewGauge().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Ratio"))).
		Ratio(0.211_313_934_313_1).
		UseUnicode(false).
		Render(layout.NewRect(2, 3, 36, 3), buf)

	assertLines(t, buf, []string{
		"  ┌Percentage────────────────────────┐  ",
		"  │███████████████43%                │  ",
		"  └──────────────────────────────────┘  ",
		"  ┌Ratio─────────────────────────────┐  ",
		"  │███████        21%                │  ",
		"  └──────────────────────────────────┘  ",
	})
}

func TestGauge_shouldApplyStyles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 5))
	gaugeStyle := style.NewStyle().Fg(style.Blue).Bg(style.Red)

	widgets.NewGauge().
		Block(widgets.BorderedBlock().Title(text.NewLine(text.StyledSpan("Test", style.NewStyle().Fg(style.Red))))).
		GaugeStyle(gaugeStyle).
		Percent(43).
		Label(text.StyledSpan("43%", style.NewStyle().Fg(style.Green).AddModifier(style.ModifierBold))).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌Test──────┐",
		"│████      │",
		"│███43%    │",
		"│████      │",
		"└──────────┘",
	})
	for x := 1; x <= 4; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.Red))
	}
	for y := 1; y <= 3; y++ {
		for x := 1; x <= 10; x++ {
			if y == 2 && x >= 4 && x <= 6 {
				continue
			}
			assertCellStyle(t, buf, x, y, gaugeStyle)
		}
	}
	assertCellStyle(t, buf, 4, 2, style.NewStyle().Fg(style.Green).Bg(style.Blue).AddModifier(style.ModifierBold))
	for x := 5; x <= 6; x++ {
		assertCellStyle(t, buf, x, 2, style.NewStyle().Fg(style.Green).Bg(style.Red).AddModifier(style.ModifierBold))
	}
}

func TestGauge_shouldSupportLargeLabels(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))

	widgets.NewGauge().
		Percent(43).
		LabelString("43333333333333333333333333333%").
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"4333333333"})
}

func TestGauge_shouldPanicOnInvalidPercent(t *testing.T) {
	assertPanics(t, func() { widgets.NewGauge().Percent(-1) })
	assertPanics(t, func() { widgets.NewGauge().Percent(101) })
}

func TestGauge_shouldPanicOnInvalidRatio(t *testing.T) {
	assertPanics(t, func() { widgets.NewGauge().Ratio(-0.1) })
	assertPanics(t, func() { widgets.NewGauge().Ratio(1.1) })
}

func TestGauge_shouldIgnoreEmptyArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	widgets.NewGauge().Percent(50).Render(layout.NewRect(0, 0, 0, 1), buf)
	widgets.NewGauge().Percent(50).Render(layout.NewRect(0, 0, 1, 0), buf)

	assertLines(t, buf, []string{" "})
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}
