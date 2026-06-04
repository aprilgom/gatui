package widgets_test

import (
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/widgets"
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

func TestFill_shouldFillAreaWithSymbolAndStyle(t *testing.T) {
	buf := buffer.WithLines([]string{"abc", "def"})
	fillStyle := style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold)

	widgets.NewFill("x").Style(fillStyle).Render(layout.NewRect(1, 0, 2, 2), buf)

	assertLines(t, buf, []string{"axx", "dxx"})
	assertCellStyle(t, buf, 1, 0, fillStyle)
	assertCellStyle(t, buf, 2, 1, fillStyle)
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

func TestFill_acceptsOwnedStringSymbol(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))
	symbol := string([]rune{'•'})

	widgets.NewFill(symbol).Render(buf.Area, buf)

	assertLines(t, buf, []string{"••"})
}

func TestFill_symbolSetterReplacesSymbol(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))

	widgets.NewFill("a").Symbol("b").Render(buf.Area, buf)

	assertLines(t, buf, []string{"bb"})
}

func TestFill_symbolSetterNormalizesEmptySymbol(t *testing.T) {
	buf := buffer.WithLines([]string{"abc"})

	widgets.NewFill("a").Symbol("").Render(layout.NewRect(1, 0, 1, 1), buf)

	assertLines(t, buf, []string{"a c"})
}

func TestFill_stylizeShorthandWorks(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))

	widgets.NewFill("*").Fg(style.Blue).Bold().Render(buf.Area, buf)

	assertLines(t, buf, []string{"**"})
	wantStyle := style.NewStyle().Fg(style.Blue).AddModifier(style.ModifierBold)
	assertCellStyle(t, buf, 0, 0, wantStyle)
	assertCellStyle(t, buf, 1, 0, wantStyle)
}

func TestFill_shouldIgnoreEmptyArea(t *testing.T) {
	buf := buffer.WithLines([]string{"abc"})

	widgets.NewFill("x").Render(layout.NewRect(0, 0, 0, 1), buf)
	widgets.NewFill("x").Render(layout.NewRect(0, 0, 1, 0), buf)

	assertLines(t, buf, []string{"abc"})
}

func TestFill_shouldClipAreaToBuffer(t *testing.T) {
	buf := buffer.WithLines([]string{"abcd", "efgh"})

	widgets.NewFill("x").Render(layout.NewRect(2, 1, 5, 3), buf)

	assertLines(t, buf, []string{"abcd", "efxx"})
}

func TestFill_shouldIgnoreFullyOutOfBoundsArea(t *testing.T) {
	buf := buffer.WithLines([]string{"abcd", "efgh"})

	widgets.NewFill("x").Render(layout.NewRect(5, 0, 2, 2), buf)
	widgets.NewFill("x").Render(layout.NewRect(0, 3, 2, 2), buf)

	assertLines(t, buf, []string{"abcd", "efgh"})
}

func TestFill_shouldRenderWithOffsetBufferArea(t *testing.T) {
	buf := buffer.Filled(layout.NewRect(3, 2, 4, 2), buffer.NewCell("."))

	widgets.NewFill("x").Render(layout.NewRect(1, 1, 4, 3), buf)

	assertLines(t, buf, []string{"xx..", "xx.."})
}

func TestFill_shouldUseBufferWidthSemanticsForWideGraphemes(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	widgets.NewFill("コ").Render(buf.Area, buf)

	assertLines(t, buf, []string{"ココ "})
}
