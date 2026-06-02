package widgets

import (
	"math"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

func TestBlock_new(t *testing.T) {
	block := NewBlock()

	if block.borders != NoBorders {
		t.Fatalf("NewBlock().borders = %v, want %v", block.borders, NoBorders)
	}
	if block.padding != PaddingZero() {
		t.Fatalf("NewBlock().padding = %#v, want %#v", block.padding, PaddingZero())
	}
	if block.style != style.NewStyle() {
		t.Fatalf("NewBlock().style = %#v, want %#v", block.style, style.NewStyle())
	}

	bordered := BorderedBlock()
	if bordered.borders != AllBorders {
		t.Fatalf("BorderedBlock().borders = %v, want %v", bordered.borders, AllBorders)
	}
}

func TestBlock_titleStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 3))
	block := BorderedBlock().
		Title(text.LineFromString("Title")).
		TitleStyle(style.NewStyle().Fg(style.Red))

	block.Render(buf.Area, buf)

	for x := 1; x <= 5; x++ {
		cell, ok := buf.CellAt(x, 0)
		if !ok {
			t.Fatalf("missing cell at (%d,0)", x)
		}
		if cell.Style != style.NewStyle().Fg(style.Red) {
			t.Fatalf("cell(%d,0).Style = %#v, want red foreground", x, cell.Style)
		}
	}
}

func TestBlock_styleIntoWorksFromUserView(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 3))
	block := BorderedBlock().
		Style(style.NewStyle().Bg(style.Green)).
		BorderStyle(style.NewStyle().Fg(style.Cyan)).
		Title(text.LineFromString("T")).
		TitleStyle(style.NewStyle().Fg(style.Red))

	block.Render(buf.Area, buf)

	assertBlockCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Cyan).Bg(style.Green))
	assertBlockCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).Bg(style.Green))
	assertBlockCellStyle(t, buf, 1, 1, style.NewStyle())
}

func TestBlock_leftTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))

	NewBlock().Title(text.LineFromString("L12")).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L12       "})
}

func TestBlock_leftTitleTruncated(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))

	NewBlock().Title(text.LineFromString("L1234567890")).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L123456789"})
}

func TestBlock_centerTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Center).
		Title(text.LineFromString("C12"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"   C12    "})
}

func TestBlock_centerTitleTruncated(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Center).
		Title(text.LineFromString("C1234567890"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"C123456789"})
}

func TestBlock_centerTitleTruncatesLeftTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		Title(text.LineFromString("L1234")).
		Title(text.LineFromString("C5678").Center())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L1C5678   "})
}

func TestBlock_rightTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Right).
		Title(text.LineFromString("R12"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"       R12"})
}

func TestBlock_rightTitleTruncated(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Right).
		Title(text.LineFromString("R1234567890"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"R123456789"})
}

func TestBlock_rightTitleTruncatesLeftTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		Title(text.LineFromString("L12345")).
		Title(text.LineFromString("R67890").Right())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L123R67890"})
}

func TestBlock_rightTitleTruncatesCenterTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		Title(text.LineFromString("C12345").Center()).
		Title(text.LineFromString("R67890").Right())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"  C1R67890"})
}

func TestBlock_titleAlignmentOverridesBlockTitleAlignment(t *testing.T) {
	tests := []struct {
		blockAlignment layout.Alignment
		lineAlignment  func(text.Line) text.Line
		want           string
	}{
		{blockAlignment: layout.Right, lineAlignment: func(line text.Line) text.Line { return line.Left() }, want: "test    "},
		{blockAlignment: layout.Left, lineAlignment: func(line text.Line) text.Line { return line.Center() }, want: "  test  "},
		{blockAlignment: layout.Center, lineAlignment: func(line text.Line) text.Line { return line.Right() }, want: "    test"},
	}
	for _, tt := range tests {
		buf := buffer.Empty(layout.NewRect(0, 0, 8, 1))
		block := NewBlock().
			TitleAlignment(tt.blockAlignment).
			Title(tt.lineAlignment(text.LineFromString("test")))

		block.Render(buf.Area, buf)

		assertBlockLines(t, buf, []string{tt.want})
	}
}

func TestBlock_titleStyleOverridesBlockTitleStyle(t *testing.T) {
	for _, alignment := range []layout.Alignment{layout.Left, layout.Center, layout.Right} {
		buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
		block := NewBlock().
			TitleAlignment(alignment).
			TitleStyle(style.NewStyle().Fg(style.Green).Bg(style.Red)).
			Title(text.LineFromString("test").Fg(style.Yellow))

		block.Render(buf.Area, buf)

		for x := range 4 {
			assertBlockCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.Yellow).Bg(style.Red))
		}
	}
}

func TestBlock_titlePosition(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 2))
	block := NewBlock().
		TitlePosition(TitlePositionBottom).
		Title(text.LineFromString("test"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"    ",
		"test",
	})
}

func TestBlock_titleTopBottom(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 3))
	block := BorderedBlock().
		TitleTop(text.LineFromString("A").Left()).
		TitleTop(text.LineFromString("B").Center()).
		TitleTop(text.LineFromString("C").Right()).
		TitleBottom(text.LineFromString("D").Left()).
		TitleBottom(text.LineFromString("E").Center()).
		TitleBottom(text.LineFromString("F").Right())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"┌A───B───C┐",
		"│         │",
		"└D───E───F┘",
	})
}

func TestBlock_hasTitleAtPosition(t *testing.T) {
	block := NewBlock()
	if block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("empty block has top title")
	}
	if block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("empty block has bottom title")
	}

	block = NewBlock().TitleTop(text.LineFromString("test"))
	if !block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("TitleTop block missing top title")
	}
	if block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("TitleTop block has bottom title")
	}

	block = NewBlock().TitleBottom(text.LineFromString("test"))
	if block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("TitleBottom block has top title")
	}
	if !block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("TitleBottom block missing bottom title")
	}

	block = NewBlock().
		TitleTop(text.LineFromString("test")).
		TitleBottom(text.LineFromString("test"))
	if !block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("mixed block missing top title")
	}
	if !block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("mixed block missing bottom title")
	}

	block = NewBlock().
		Title(text.LineFromString("top")).
		TitlePosition(TitlePositionBottom).
		Title(text.LineFromString("bottom"))
	if !block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("default-position block missing top title")
	}
	if !block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("default-position block missing bottom title")
	}
}

func TestBlock_titlesAreaHandlesEmptyAreaWithoutPanicking(t *testing.T) {
	block := NewBlock()

	got := block.titlesArea(layout.NewRect(0, 0, 0, 0), TitlePositionBottom)

	if got != layout.NewRect(0, 0, 0, 1) {
		t.Fatalf("titlesArea = %#v, want %#v", got, layout.NewRect(0, 0, 0, 1))
	}
}

func TestBlock_titlesAreaSaturatesWhenLeftBorderOffsetOverflows(t *testing.T) {
	block := NewBlock().Borders(LeftBorder)

	got := block.titlesArea(layout.NewRect(layout.MaxCoordinate, 0, 1, 1), TitlePositionTop)

	if got != layout.NewRect(layout.MaxCoordinate, 0, 1, 1) {
		t.Fatalf("titlesArea = %#v, want %#v", got, layout.NewRect(layout.MaxCoordinate, 0, 1, 1))
	}
}

func TestBlock_renderRightAlignedEmptyTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 3))
	block := NewBlock().
		TitleAlignment(layout.Right).
		Title(text.LineFromString(""))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"               ",
		"               ",
		"               ",
	})
}

func TestBlock_renderInMinimalBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	BorderedBlock().Title(text.LineFromString("Title")).Render(buf.Area, buf)

	cell, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatalf("missing cell at (0,0)")
	}
	if cell.Symbol != "┌" {
		t.Fatalf("cell(0,0).Symbol = %q, want %q", cell.Symbol, "┌")
	}
}

func TestBlock_renderInZeroSizeBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))

	BorderedBlock().Title(text.LineFromString("Title")).Render(buf.Area, buf)
}

func TestBlock_renderCornersHandlesEmptyAreaWithoutPanicking(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))

	BorderedBlock().renderBorders(buf.Area, buf)
}

func TestBlock_renderSidesHandlesEmptyAreaWithoutPanicking(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))

	NewBlock().Borders(LeftBorder|RightBorder).renderBorders(buf.Area, buf)
}

func TestBlock_innerSaturatesWhenPaddingSumOverflows(t *testing.T) {
	block := BorderedBlock().Padding(NewPadding(math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt))

	inner := block.Inner(layout.NewRect(1, 2, 3, 4))

	if inner.Width != 0 || inner.Height != 0 {
		t.Fatalf("inner size = %dx%d, want 0x0", inner.Width, inner.Height)
	}
}

func TestBlock_verticalSpaceSaturatesWhenSpaceOverflows(t *testing.T) {
	block := NewBlock().
		Borders(TopBorder | BottomBorder).
		Padding(NewPadding(0, 0, math.MaxInt, math.MaxInt))

	if got := block.verticalSpace(); got != math.MaxInt {
		t.Fatalf("verticalSpace = %d, want %d", got, math.MaxInt)
	}
}

func assertBlockCellStyle(t *testing.T, buf *buffer.Buffer, x, y int, want style.Style) {
	t.Helper()
	cell, ok := buf.CellAt(x, y)
	if !ok {
		t.Fatalf("missing cell at (%d,%d)", x, y)
	}
	if cell.Style != want {
		t.Fatalf("cell(%d,%d).Style = %#v, want %#v", x, y, cell.Style, want)
	}
}

func assertBlockLines(t *testing.T, buf *buffer.Buffer, want []string) {
	t.Helper()
	if got := buf.Lines(); !equalStringSlices(got, want) {
		t.Fatalf("buffer lines = %#v, want %#v", got, want)
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
