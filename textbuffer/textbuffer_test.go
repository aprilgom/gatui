package textbuffer_test

import (
	"slices"
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/textbuffer"
)

func TestSetSpan_shouldWriteRawContent(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	endX, endY := textbuffer.SetSpan(buf, 1, 0, text.NewSpan("abc"), 4)

	if endX != 4 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (4,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{" abc "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestSetSpan_shouldWriteContentWithSpanStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	red := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold)

	textbuffer.SetSpan(buf, 0, 0, text.StyledSpan("ab", red), 4)

	for x := range 2 {
		cell, ok := buf.CellAt(x, 0)
		if !ok {
			t.Fatalf("missing cell at %d,0", x)
		}
		if cell.Style != red {
			t.Fatalf("cell %d style = %#v, want %#v", x, cell.Style, red)
		}
	}
}

func TestSetSpan_shouldClipToMaxWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	endX, endY := textbuffer.SetSpan(buf, 0, 0, text.NewSpan("abcdef"), 3)

	if endX != 3 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (3,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"abc  "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestSetLine_shouldHandleEmptyLine(t *testing.T) {
	buf := buffer.WithLines([]string{"abc"})

	endX, endY := textbuffer.SetLine(buf, 1, 0, text.NewLine(), 2)

	if endX != 1 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (1,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"abc"}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestSetLine_rawRatatuiCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantX   int
	}{
		{name: "empty", content: "", want: "     ", wantX: 0},
		{name: "one", content: "1", want: "1    ", wantX: 1},
		{name: "full", content: "12345", want: "12345", wantX: 5},
		{name: "overflow", content: "123456", want: "12345", wantX: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

			endX, endY := textbuffer.SetLine(buf, 0, 0, text.LineFromString(tt.content), 5)

			if endX != tt.wantX || endY != 0 {
				t.Fatalf("end = (%d,%d), want (%d,0)", endX, endY, tt.wantX)
			}
			if got, want := buf.Lines(), []string{tt.want}; !slices.Equal(got, want) {
				t.Fatalf("lines = %#v, want %#v", got, want)
			}
		})
	}
}

func TestSetLine_styledRatatuiCases(t *testing.T) {
	blue := style.NewStyle().Fg(style.Blue)
	tests := []struct {
		name         string
		content      string
		want         string
		styledCells  int
		defaultCells int
	}{
		{name: "empty", content: "", want: "     ", styledCells: 0, defaultCells: 5},
		{name: "one", content: "1", want: "1    ", styledCells: 1, defaultCells: 4},
		{name: "full", content: "12345", want: "12345", styledCells: 5, defaultCells: 0},
		{name: "overflow", content: "123456", want: "12345", styledCells: 5, defaultCells: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
			line := text.StyledLine(tt.content, blue)

			textbuffer.SetLine(buf, 0, 0, line, 5)

			if got, want := buf.Lines(), []string{tt.want}; !slices.Equal(got, want) {
				t.Fatalf("lines = %#v, want %#v", got, want)
			}
			for x := 0; x < tt.styledCells; x++ {
				cell, ok := buf.CellAt(x, 0)
				if !ok {
					t.Fatalf("missing cell at %d,0", x)
				}
				if cell.Style != blue {
					t.Fatalf("cell %d style = %#v, want %#v", x, cell.Style, blue)
				}
			}
			for x := tt.styledCells; x < tt.styledCells+tt.defaultCells; x++ {
				cell, ok := buf.CellAt(x, 0)
				if !ok {
					t.Fatalf("missing cell at %d,0", x)
				}
				if want := style.NewStyle(); cell.Style != want {
					t.Fatalf("cell %d style = %#v, want %#v", x, cell.Style, want)
				}
			}
		})
	}
}

func TestSetLine_shouldWritePartialExactAndOverflowWidths(t *testing.T) {
	tests := []struct {
		name     string
		maxWidth int
		wantEndX int
		wantLine string
	}{
		{name: "partial", maxWidth: 2, wantEndX: 2, wantLine: "ab   "},
		{name: "exact", maxWidth: 5, wantEndX: 5, wantLine: "abcde"},
		{name: "overflow", maxWidth: 8, wantEndX: 5, wantLine: "abcde"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
			line := text.NewLine(text.NewSpan("ab"), text.NewSpan("cde"), text.NewSpan("fg"))

			endX, endY := textbuffer.SetLine(buf, 0, 0, line, tt.maxWidth)

			if endX != tt.wantEndX || endY != 0 {
				t.Fatalf("end = (%d,%d), want (%d,0)", endX, endY, tt.wantEndX)
			}
			if got, want := buf.Lines(), []string{tt.wantLine}; !slices.Equal(got, want) {
				t.Fatalf("lines = %#v, want %#v", got, want)
			}
		})
	}
}

func TestSetLine_shouldPatchLineStyleIntoSpanStyles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	lineStyle := style.NewStyle().Fg(style.Red).Bg(style.Blue)
	spanStyle := style.NewStyle().Fg(style.Green).AddModifier(style.ModifierItalic)
	line := text.NewLine(text.NewSpan("a"), text.StyledSpan("b", spanStyle)).Style(lineStyle)

	textbuffer.SetLine(buf, 0, 0, line, 4)

	first, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatal("missing cell at 0,0")
	}
	if first.Style != lineStyle {
		t.Fatalf("first style = %#v, want %#v", first.Style, lineStyle)
	}
	second, ok := buf.CellAt(1, 0)
	if !ok {
		t.Fatal("missing cell at 1,0")
	}
	wantSecond := lineStyle.Patch(spanStyle)
	if second.Style != wantSecond {
		t.Fatalf("second style = %#v, want %#v", second.Style, wantSecond)
	}
}

func TestSetLine_shouldNotWriteWideGraphemeIntoOneRemainingCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	line := text.NewLine(text.NewSpan("ab"), text.NewSpan("コ"))

	endX, endY := textbuffer.SetLine(buf, 0, 0, line, 3)

	if endX != 2 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (2,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"ab "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestSetLine_shouldFollowSetStringNZeroWidthMarkSemantics(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	line := text.NewLine(text.NewSpan("a\u200B"), text.NewSpan("b"))

	endX, endY := textbuffer.SetLine(buf, 0, 0, line, 3)

	if endX != 2 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (2,0)", endX, endY)
	}
	cell, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatal("missing cell at 0,0")
	}
	if got, want := cell.Symbol, "a\u200B"; got != want {
		t.Fatalf("symbol = %q, want %q", got, want)
	}
}

func TestSetSpanAndLine_shouldNoOpForNilBufferOutsideAreaAndInvalidWidth(t *testing.T) {
	if endX, endY := textbuffer.SetSpan(nil, 4, 5, text.NewSpan("abc"), 3); endX != 4 || endY != 5 {
		t.Fatalf("nil SetSpan end = (%d,%d), want (4,5)", endX, endY)
	}
	if endX, endY := textbuffer.SetLine(nil, 4, 5, text.LineFromString("abc"), 3); endX != 4 || endY != 5 {
		t.Fatalf("nil SetLine end = (%d,%d), want (4,5)", endX, endY)
	}

	buf := buffer.WithLines([]string{"abc"})
	if endX, endY := textbuffer.SetSpan(buf, 3, 0, text.NewSpan("x"), 1); endX != 3 || endY != 0 {
		t.Fatalf("outside SetSpan end = (%d,%d), want (3,0)", endX, endY)
	}
	if endX, endY := textbuffer.SetLine(buf, 0, 0, text.LineFromString("x"), 0); endX != 0 || endY != 0 {
		t.Fatalf("zero-width SetLine end = (%d,%d), want (0,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"abc"}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}
