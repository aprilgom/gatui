package widgets_test

import (
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/symbols"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
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

func TestLineGauge_defaultSymbolsShouldUseLineHorizontal(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 1))

	widgets.NewLineGauge().LabelString("").Ratio(0.50).Render(buf.Area, buf)

	for x := 1; x < 8; x++ {
		assertCellSymbol(t, buf, x, 0, symbols.LineHorizontal)
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

func TestLineGauge_shouldPreserveStyledMultiSpanLabel(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))
	label := text.NewLine(
		text.StyledSpan("ab", style.NewStyle().Fg(style.Red)),
		text.StyledSpan("cd", style.NewStyle().Fg(style.Green)),
	).Style(style.NewStyle().Bg(style.Blue))

	widgets.NewLineGauge().
		Label(label).
		Ratio(0.50).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"abcd ─"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red).Bg(style.Blue))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).Bg(style.Blue))
	assertCellStyle(t, buf, 2, 0, style.NewStyle().Fg(style.Green).Bg(style.Blue))
	assertCellStyle(t, buf, 3, 0, style.NewStyle().Fg(style.Green).Bg(style.Blue))
}

func TestLineGauge_canBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 16, 1))
	baseStyle := style.NewStyle().
		Fg(style.Black).
		Bg(style.White).
		AddModifier(style.ModifierBold | style.ModifierDim | style.ModifierItalic)
	filledStyle := style.NewStyle().Fg(style.Green)
	unfilledStyle := style.NewStyle().Fg(style.Red)

	widgets.NewLineGauge().
		Fg(style.Black).
		Bg(style.White).
		Bold().
		Dim().
		Italic().
		FilledStyle(filledStyle).
		UnfilledStyle(unfilledStyle).
		Label(text.NewLine(
			text.StyledSpan("ab", style.NewStyle().Fg(style.Blue)),
			text.NewSpan("cd"),
		).Style(style.NewStyle().Bg(style.Cyan))).
		Ratio(0.50).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"abcd ───────────"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Blue).Bg(style.Cyan).AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Blue).Bg(style.Cyan).AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
	assertCellStyle(t, buf, 2, 0, style.NewStyle().Fg(style.Black).Bg(style.Cyan).AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
	assertCellStyle(t, buf, 3, 0, style.NewStyle().Fg(style.Black).Bg(style.Cyan).AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
	assertCellStyle(t, buf, 4, 0, baseStyle)
	for x := 5; x < 10; x++ {
		assertCellStyle(t, buf, x, 0, filledStyle)
	}
	for x := 10; x < 16; x++ {
		assertCellStyle(t, buf, x, 0, unfilledStyle)
	}
}
