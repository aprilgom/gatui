package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestLineGauge_shouldRenderProgressLines(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 6))

	widgets.NewLineGauge().
		FilledStyle(style.NewStyle().Fg(style.Green)).
		UnfilledStyle(style.NewStyle().Fg(style.White)).
		Ratio(0.43).
		Render(layout.NewRect(0, 0, 20, 1), buf)
	widgets.NewLineGauge().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Gauge 2"))).
		FilledStyle(style.NewStyle().Fg(style.Green)).
		FilledSymbol("━").
		UnfilledSymbol("━").
		Ratio(0.211_313_934_313_1).
		Render(layout.NewRect(0, 1, 20, 3), buf)
	widgets.NewLineGauge().
		UnfilledSymbol(" ").
		Ratio(0.50).
		Render(layout.NewRect(0, 4, 20, 1), buf)
	widgets.NewLineGauge().
		FilledSymbol("█").
		UnfilledSymbol("░").
		Ratio(0.80).
		Render(layout.NewRect(0, 5, 20, 1), buf)

	assertLines(t, buf, []string{
		" 43% ───────────────",
		"┌Gauge 2───────────┐",
		"│ 21% ━━━━━━━━━━━━━│",
		"└──────────────────┘",
		" 50% ───────        ",
		" 80% ████████████░░░",
	})
	for x := 5; x < 11; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.Green))
	}
	for x := 11; x < 20; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.White))
	}
	for x := 6; x < 8; x++ {
		assertCellStyle(t, buf, x, 2, style.NewStyle().Fg(style.Green))
	}
}

func TestLineGauge_shouldPanicOnInvalidRatio(t *testing.T) {
	assertPanics(t, func() { widgets.NewLineGauge().Ratio(-0.1) })
	assertPanics(t, func() { widgets.NewLineGauge().Ratio(1.1) })
}

func TestLineGauge_shouldIgnoreEmptyArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	widgets.NewLineGauge().Ratio(0.50).Render(layout.NewRect(0, 0, 0, 1), buf)
	widgets.NewLineGauge().Ratio(0.50).Render(layout.NewRect(0, 0, 1, 0), buf)

	assertLines(t, buf, []string{" "})
}

func TestLineGauge_shouldTruncateCustomLabelToAvailableWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))

	widgets.NewLineGauge().
		LabelString("abcdef").
		Ratio(0.50).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"abcd"})
}
