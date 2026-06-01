package widgets_test

import (
	"reflect"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestParagraph_shouldPreserveStylizedSpanStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	paragraph := widgets.NewParagraph(text.NewText(text.NewLine(text.NewSpan("Text").Cyan())))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Text      "})
	for x := 0; x < 4; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.Cyan))
	}
}

func TestBlock_shouldRenderBorderTitleAndTitleStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 10))
	block := widgets.BorderedBlock().
		Title(text.NewLine(text.StyledSpan("Title", style.NewStyle().Fg(style.LightBlue))))

	block.Render(layout.NewRect(0, 0, 8, 8), buf)

	assertLines(t, buf, []string{
		"┌Title─┐  ",
		"│      │  ",
		"│      │  ",
		"│      │  ",
		"│      │  ",
		"│      │  ",
		"│      │  ",
		"└──────┘  ",
		"          ",
		"          ",
	})
	for x := 1; x <= 5; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.LightBlue))
	}
}

func TestBlock_shouldBeStylizable(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 3))

	widgets.BorderedBlock().Cyan().Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌──┐",
		"│  │",
		"└──┘",
	})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Cyan))
	assertCellStyle(t, buf, 3, 2, style.NewStyle().Fg(style.Cyan))
}

func TestParagraph_shouldRenderInsideBlockWithWrapAlignmentAndScroll(t *testing.T) {
	content := text.FromString("The library is based on immediate rendering.\nLittle line")
	paragraph := widgets.NewParagraph(content).
		Block(widgets.BorderedBlock()).
		Wrap(widgets.Wrap{Trim: true}).
		Alignment(layout.Center)
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 6))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌──────────────────┐",
		"│  The library is  │",
		"│based on immediate│",
		"│    rendering.    │",
		"│   Little line    │",
		"└──────────────────┘",
	})

	scrolled := widgets.NewParagraph(text.FromString("Paragraph can scroll horizontally!")).
		Block(widgets.BorderedBlock()).
		Scroll(0, 10)
	scrollBuf := buffer.Empty(layout.NewRect(0, 0, 20, 3))

	scrolled.Render(scrollBuf.Area, scrollBuf)

	assertLines(t, scrollBuf, []string{
		"┌──────────────────┐",
		"│can scroll horizon│",
		"└──────────────────┘",
	})
}

func TestClear_shouldResetCellsToBlankDefaultStyle(t *testing.T) {
	buf := buffer.WithLines([]string{"abcd", "efgh"})
	buf.SetFg(layout.NewRect(1, 0, 2, 2), style.Red)

	widgets.Clear{}.Render(layout.NewRect(1, 0, 2, 2), buf)

	assertLines(t, buf, []string{"a  d", "e  h"})
	assertCellStyle(t, buf, 1, 0, style.NewStyle())
	assertCellStyle(t, buf, 2, 1, style.NewStyle())
}

func assertLines(t *testing.T, buf *buffer.Buffer, expected []string) {
	t.Helper()
	if actual := buf.Lines(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("lines mismatch\nactual:   %#v\nexpected: %#v", actual, expected)
	}
}

func assertCellStyle(t *testing.T, buf *buffer.Buffer, x, y int, expected style.Style) {
	t.Helper()
	cell, ok := buf.CellAt(x, y)
	if !ok {
		t.Fatalf("expected cell at (%d,%d)", x, y)
	}
	if cell.Style != expected {
		t.Fatalf("style at (%d,%d) = %#v, want %#v", x, y, cell.Style, expected)
	}
}
