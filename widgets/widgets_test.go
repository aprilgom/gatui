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

func TestParagraph_shouldApplyWidgetStyleToEntireArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))
	paragraph := widgets.NewParagraph(text.FromString("hi")).
		Style(style.NewStyle().Bg(style.Green))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"hi   ",
		"     ",
	})
	assertAllCellsStyle(t, buf, style.NewStyle().Bg(style.Green))
}

func TestParagraph_shouldPatchSpanStyleOverWidgetStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	line := text.NewLine(text.StyledSpan("hi", style.NewStyle().Fg(style.Red)))
	paragraph := widgets.NewParagraph(text.NewText(line)).
		Style(style.NewStyle().Fg(style.Yellow).Bg(style.Green))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{"hi   "})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red).Bg(style.Green))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).Bg(style.Green))
	assertCellStyle(t, buf, 2, 0, style.NewStyle().Fg(style.Yellow).Bg(style.Green))
}

func TestParagraph_shouldRenderStyledLines(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 4))
	paragraph := widgets.NewParagraph(text.NewText(
		text.LineFromString("unformatted"),
		text.StyledLine("bold text", style.NewStyle().AddModifier(style.ModifierBold)),
		text.StyledLine("cyan text", style.NewStyle().Fg(style.Cyan)),
		text.StyledLine("dim text", style.NewStyle().AddModifier(style.ModifierDim)),
	))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"unformatted",
		"bold text  ",
		"cyan text  ",
		"dim text   ",
	})
	for x := 0; x < len("bold text"); x++ {
		assertCellStyle(t, buf, x, 1, style.NewStyle().AddModifier(style.ModifierBold))
	}
	for x := 0; x < len("cyan text"); x++ {
		assertCellStyle(t, buf, x, 2, style.NewStyle().Fg(style.Cyan))
	}
	for x := 0; x < len("dim text"); x++ {
		assertCellStyle(t, buf, x, 3, style.NewStyle().AddModifier(style.ModifierDim))
	}
}

func TestParagraph_shouldPatchSpanStyleOverLineStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	line := text.NewLine(text.StyledSpan("hi", style.NewStyle().Fg(style.Red))).
		Style(style.NewStyle().Fg(style.Yellow).Bg(style.Green))
	paragraph := widgets.NewParagraph(text.NewText(line))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{"hi   "})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red).Bg(style.Green))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).Bg(style.Green))
}

func TestParagraph_shouldPatchTextLineAndSpanStylesInOrder(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	line := text.NewLine(text.StyledSpan("hi", style.NewStyle().Fg(style.Red))).
		Style(style.NewStyle().AddModifier(style.ModifierBold))
	content := text.NewText(line).PatchStyle(style.NewStyle().Fg(style.Yellow))
	paragraph := widgets.NewParagraph(content).
		Style(style.NewStyle().Bg(style.Blue))

	paragraph.Render(buf.Area, buf)

	wantStyle := style.NewStyle().
		Fg(style.Red).
		Bg(style.Blue).
		AddModifier(style.ModifierBold)
	assertLines(t, buf, []string{"hi   "})
	assertCellStyle(t, buf, 0, 0, wantStyle)
	assertCellStyle(t, buf, 1, 0, wantStyle)
}

func TestParagraph_shouldApplyWidgetStyleBehindBlock(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 3))
	paragraph := widgets.NewParagraph(text.FromString("")).
		Style(style.NewStyle().Bg(style.Green)).
		Block(widgets.BorderedBlock().Fg(style.Cyan))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌───┐",
		"│   │",
		"└───┘",
	})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Cyan).Bg(style.Green))
	assertCellStyle(t, buf, 2, 1, style.NewStyle().Bg(style.Green))
}

func TestParagraphStylize_shouldUpdateWidgetStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))
	paragraph := widgets.NewParagraph(text.FromString("hi")).Cyan().Bold()

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"hi   ",
		"     ",
	})
	assertAllCellsStyle(t, buf, style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierBold))
}

func TestParagraph_LineCount_shouldReturnTextHeightWithoutWrap(t *testing.T) {
	paragraph := widgets.NewParagraph(text.FromString("one\ntwo"))

	if got := paragraph.LineCount(20); got != 2 {
		t.Fatalf("LineCount(20) = %d, want 2", got)
	}
}

func TestParagraph_LineCount_shouldWrapContent(t *testing.T) {
	paragraph := widgets.NewParagraph(text.FromString("Hello World")).
		Wrap(widgets.Wrap{Trim: false})

	if got := paragraph.LineCount(20); got != 1 {
		t.Fatalf("LineCount(20) = %d, want 1", got)
	}
	if got := paragraph.LineCount(10); got != 2 {
		t.Fatalf("LineCount(10) = %d, want 2", got)
	}
}

func TestParagraph_LineCount_shouldAccountForBlockPaddingAndBorders(t *testing.T) {
	paragraph := widgets.NewParagraph(text.FromString("Hello World")).
		Block(widgets.BorderedBlock().Padding(widgets.PaddingVertical(1))).
		Wrap(widgets.Wrap{Trim: false})

	if got := paragraph.LineCount(20); got != 5 {
		t.Fatalf("LineCount(20) = %d, want 5", got)
	}
	if got := paragraph.LineCount(10); got != 6 {
		t.Fatalf("LineCount(10) = %d, want 6", got)
	}
}

func TestParagraph_LineWidth_shouldReturnLongestDisplayWidth(t *testing.T) {
	paragraph := widgets.NewParagraph(text.FromString("Hello World\nhi\nHello World!!!"))

	if got := paragraph.LineWidth(); got != 14 {
		t.Fatalf("LineWidth() = %d, want 14", got)
	}
}

func TestParagraph_LineWidth_shouldAccountForUnicodeDisplayWidth(t *testing.T) {
	paragraph := widgets.NewParagraph(text.FromString("コンピ"))

	if got := paragraph.LineWidth(); got != 6 {
		t.Fatalf("LineWidth() = %d, want 6", got)
	}
}

func TestParagraph_LineWidth_shouldAccountForBlockPaddingAndBorders(t *testing.T) {
	paragraph := widgets.NewParagraph(text.FromString("abc")).
		Block(widgets.BorderedBlock().Padding(widgets.PaddingHorizontal(2)))

	if got := paragraph.LineWidth(); got != 9 {
		t.Fatalf("LineWidth() = %d, want 9", got)
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

func TestPaddingConstructors_shouldMatchRatatuiSemantics(t *testing.T) {
	assertPadding := func(name string, actual, expected widgets.Padding) {
		t.Helper()
		if actual != expected {
			t.Fatalf("%s = %#v, want %#v", name, actual, expected)
		}
	}

	assertPadding("PaddingZero", widgets.PaddingZero(), widgets.NewPadding(0, 0, 0, 0))
	assertPadding("PaddingHorizontal", widgets.PaddingHorizontal(1), widgets.NewPadding(1, 1, 0, 0))
	assertPadding("PaddingVertical", widgets.PaddingVertical(1), widgets.NewPadding(0, 0, 1, 1))
	assertPadding("PaddingUniform", widgets.PaddingUniform(1), widgets.NewPadding(1, 1, 1, 1))
	assertPadding("PaddingProportional", widgets.PaddingProportional(1), widgets.NewPadding(2, 2, 1, 1))
	assertPadding("PaddingSymmetric", widgets.PaddingSymmetric(1, 2), widgets.NewPadding(1, 1, 2, 2))
	assertPadding("PaddingLeft", widgets.PaddingLeft(1), widgets.NewPadding(1, 0, 0, 0))
	assertPadding("PaddingRight", widgets.PaddingRight(1), widgets.NewPadding(0, 1, 0, 0))
	assertPadding("PaddingTop", widgets.PaddingTop(1), widgets.NewPadding(0, 0, 1, 0))
	assertPadding("PaddingBottom", widgets.PaddingBottom(1), widgets.NewPadding(0, 0, 0, 1))
}

func TestBlockInner_shouldAccountForBordersAndPadding(t *testing.T) {
	block := widgets.BorderedBlock().Padding(widgets.NewPadding(2, 2, 1, 1))

	actual := block.Inner(layout.NewRect(0, 0, 22, 12))
	expected := layout.NewRect(3, 2, 16, 8)

	if actual != expected {
		t.Fatalf("inner = %#v, want %#v", actual, expected)
	}
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
		"│    Little line   │",
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

func TestParagraph_shouldPreserveTrailingNBSP(t *testing.T) {
	line := text.NewLine(text.NewSpan("NBSP"), text.NewSpan("\u00a0"))
	paragraph := widgets.NewParagraph(text.NewText(line)).
		Block(widgets.BorderedBlock())
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 3))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌──────────────────┐",
		"│NBSP\u00a0             │",
		"└──────────────────┘",
	})
	assertCellSymbol(t, buf, 5, 1, "\u00a0")
	if cell, _ := buf.CellAt(5, 1); cell.Symbol == " " {
		t.Fatalf("symbol at (5,1) was normalized to a regular space")
	}
}

func TestParagraph_shouldRenderDoubleWidthGraphemes(t *testing.T) {
	content := text.FromString("コンピュータ上で文字を扱う場合、典型的には文字による通信を行う場合にその両端点では、")
	paragraph := widgets.NewParagraph(content).
		Block(widgets.BorderedBlock()).
		Wrap(widgets.Wrap{Trim: true})
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 10))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌────────┐",
		"│コンピュ│",
		"│ータ上で│",
		"│文字を扱│",
		"│う場合、│",
		"│典型的に│",
		"│は文字に│",
		"│よる通信│",
		"│を行う場│",
		"└────────┘",
	})
	assertCellSymbol(t, buf, 1, 1, "コ")
	assertCellSymbol(t, buf, 2, 1, " ")
	assertCellSymbol(t, buf, 7, 1, "ュ")
	assertCellSymbol(t, buf, 8, 1, " ")
}

func TestParagraph_shouldRenderMixedWidthGraphemes(t *testing.T) {
	content := text.FromString("aコンピュータ上で文字を扱う場合、")
	paragraph := widgets.NewParagraph(content).
		Block(widgets.BorderedBlock()).
		Wrap(widgets.Wrap{Trim: true})
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 7))

	paragraph.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌────────┐",
		"│aコンピ │",
		"│ュータ上│",
		"│で文字を│",
		"│扱う場合│",
		"│、      │",
		"└────────┘",
	})
	assertCellSymbol(t, buf, 1, 1, "a")
	assertCellSymbol(t, buf, 2, 1, "コ")
	assertCellSymbol(t, buf, 3, 1, " ")
	assertCellSymbol(t, buf, 8, 1, " ")
}

func TestParagraph_shouldScrollHorizontallyByDisplayWidth(t *testing.T) {
	content := text.FromString("段落现在可以水平滚动了！\nParagraph can scroll horizontally!\nLittle line")
	paragraph := widgets.NewParagraph(content).Block(widgets.BorderedBlock())

	leftBuf := buffer.Empty(layout.NewRect(0, 0, 20, 10))
	paragraph.Alignment(layout.Left).Scroll(0, 7).Render(leftBuf.Area, leftBuf)
	assertLines(t, leftBuf, []string{
		"┌──────────────────┐",
		"│在可以水平滚动了！│",
		"│ph can scroll hori│",
		"│line              │",
		"│                  │",
		"│                  │",
		"│                  │",
		"│                  │",
		"│                  │",
		"└──────────────────┘",
	})

	rightBuf := buffer.Empty(layout.NewRect(0, 0, 20, 10))
	paragraph.Alignment(layout.Right).Scroll(0, 7).Render(rightBuf.Area, rightBuf)
	assertLines(t, rightBuf, []string{
		"┌──────────────────┐",
		"│段落现在可以水平滚│",
		"│Paragraph can scro│",
		"│       Little line│",
		"│                  │",
		"│                  │",
		"│                  │",
		"│                  │",
		"│                  │",
		"└──────────────────┘",
	})
}

func TestParagraph_shouldWorkWithBlockPadding(t *testing.T) {
	const sampleString = "The library is based on the principle of immediate rendering with intermediate buffers. This means that at each new frame you should build all widgets that are supposed to be part of the UI."
	block := widgets.BorderedBlock().Padding(widgets.NewPadding(2, 2, 1, 1))
	paragraph := widgets.NewParagraph(text.NewText(text.LineFromString(sampleString))).
		Block(block).
		Wrap(widgets.Wrap{Trim: true})

	buf := buffer.Empty(layout.NewRect(0, 0, 22, 12))
	paragraph.Alignment(layout.Left).Render(buf.Area, buf)
	assertLines(t, buf, []string{
		"┌────────────────────┐",
		"│                    │",
		"│  The library is    │",
		"│  based on the      │",
		"│  principle of      │",
		"│  immediate         │",
		"│  rendering with    │",
		"│  intermediate      │",
		"│  buffers. This     │",
		"│  means that at     │",
		"│                    │",
		"└────────────────────┘",
	})

	buf = buffer.Empty(layout.NewRect(0, 0, 22, 12))
	paragraph.Alignment(layout.Right).Render(buf.Area, buf)
	assertLines(t, buf, []string{
		"┌────────────────────┐",
		"│                    │",
		"│    The library is  │",
		"│      based on the  │",
		"│      principle of  │",
		"│         immediate  │",
		"│    rendering with  │",
		"│      intermediate  │",
		"│     buffers. This  │",
		"│     means that at  │",
		"│                    │",
		"└────────────────────┘",
	})

	paragraphWithLineAlignment := widgets.NewParagraph(text.NewText(
		text.LineFromString("This is always centered.").Center(),
		text.LineFromString(sampleString),
	)).
		Block(block).
		Wrap(widgets.Wrap{Trim: true}).
		Alignment(layout.Right)
	buf = buffer.Empty(layout.NewRect(0, 0, 22, 14))
	paragraphWithLineAlignment.Render(buf.Area, buf)
	assertLines(t, buf, []string{
		"┌────────────────────┐",
		"│                    │",
		"│   This is always   │",
		"│      centered.     │",
		"│    The library is  │",
		"│      based on the  │",
		"│      principle of  │",
		"│         immediate  │",
		"│    rendering with  │",
		"│      intermediate  │",
		"│     buffers. This  │",
		"│     means that at  │",
		"│                    │",
		"└────────────────────┘",
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

func assertAllCellsStyle(t *testing.T, buf *buffer.Buffer, expected style.Style) {
	t.Helper()
	for y := buf.Area.Y; y < buf.Area.Y+buf.Area.Height; y++ {
		for x := buf.Area.X; x < buf.Area.X+buf.Area.Width; x++ {
			assertCellStyle(t, buf, x, y, expected)
		}
	}
}

func assertCellSymbol(t *testing.T, buf *buffer.Buffer, x, y int, expected string) {
	t.Helper()
	cell, ok := buf.CellAt(x, y)
	if !ok {
		t.Fatalf("expected cell at (%d,%d)", x, y)
	}
	if cell.Symbol != expected {
		t.Fatalf("symbol at (%d,%d) = %q, want %q", x, y, cell.Symbol, expected)
	}
}

func assertNotPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()
	fn()
}
