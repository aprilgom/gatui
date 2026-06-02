package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/widgets"
)

func TestFill_shouldFillAreaWithSymbol(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 3))

	widgets.NewFill(".").Render(layout.NewRect(1, 1, 3, 1), buf)

	assertLines(t, buf, []string{
		"     ",
		" ... ",
		"     ",
	})
}

func TestFill_shouldPatchStyleOverExistingCells(t *testing.T) {
	buf := buffer.WithLines([]string{"abc", "def"})
	buf.SetBg(buf.Area, style.Blue)

	widgets.NewFill("x").
		Style(style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold)).
		Render(layout.NewRect(1, 0, 2, 2), buf)

	assertLines(t, buf, []string{"axx", "dxx"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Bg(style.Blue))
	wantStyle := style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold)
	assertCellStyle(t, buf, 1, 0, wantStyle)
	assertCellStyle(t, buf, 2, 1, wantStyle)
}

func TestFill_shouldRenderEmptySymbolAsSpace(t *testing.T) {
	buf := buffer.WithLines([]string{"abc"})

	widgets.NewFill("").Render(layout.NewRect(1, 0, 1, 1), buf)

	assertLines(t, buf, []string{"a c"})
}

func TestFill_shouldIgnoreEmptyArea(t *testing.T) {
	buf := buffer.WithLines([]string{"abc"})

	widgets.NewFill("x").Render(layout.NewRect(0, 0, 0, 1), buf)
	widgets.NewFill("x").Render(layout.NewRect(0, 0, 1, 0), buf)

	assertLines(t, buf, []string{"abc"})
}

func TestFill_shouldUseBufferWidthSemanticsForWideGraphemes(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	widgets.NewFill("コ").Render(buf.Area, buf)

	assertLines(t, buf, []string{"ココ "})
}
