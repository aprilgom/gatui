package widgets

import (
	"slices"
	"strings"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/text"
)

func TestWordWrapper_shouldKeepShortLines(t *testing.T) {
	lines := newWordWrapper(20, true).wrap(cellsFromLine(text.LineFromString("short")), layout.Left)

	assertWrappedLines(t, lines, []string{"short"}, []int{5})
}

func TestLineTruncator_shouldKeepShortLines(t *testing.T) {
	line := newLineTruncator(20).truncate(cellsFromLine(text.LineFromString("short")), layout.Left)

	assertWrappedLines(t, []wrappedLine{line}, []string{"short"}, []int{5})
}

func TestWordWrapper_shouldSplitLongWordsByWidth(t *testing.T) {
	lines := newWordWrapper(3, true).wrap(cellsFromLine(text.LineFromString("abcdefg")), layout.Left)

	assertWrappedLines(t, lines, []string{"abc", "def", "g"}, []int{3, 3, 1})
}

func TestLineTruncator_shouldReturnOnlyFirstWidthOfLongLine(t *testing.T) {
	line := newLineTruncator(3).truncate(cellsFromLine(text.LineFromString("abcdefg")), layout.Left)

	assertWrappedLines(t, []wrappedLine{line}, []string{"abc"}, []int{3})
}

func TestWordWrapper_shouldWrapAtWordBoundaryAndTrimSpaces(t *testing.T) {
	lines := newWordWrapper(10, true).wrap(cellsFromLine(text.LineFromString("The   library wraps")), layout.Left)

	assertWrappedLines(t, lines, []string{"The", "library", "wraps"}, []int{3, 7, 5})
}

func TestWordWrapper_shouldPreserveIndentationWhenTrimIsFalse(t *testing.T) {
	lines := newWordWrapper(6, false).wrap(cellsFromLine(text.LineFromString("  indented text")), layout.Left)

	assertWrappedLines(t, lines, []string{"  inde", "nted", " text"}, []int{6, 4, 5})
}

func TestReflow_shouldReturnEmptyForZeroWidth(t *testing.T) {
	wrapped := newWordWrapper(0, true).wrap(cellsFromLine(text.LineFromString("abc")), layout.Left)
	truncated := newLineTruncator(0).truncate(cellsFromLine(text.LineFromString("abc")), layout.Left)

	if len(wrapped) != 0 {
		t.Fatalf("wrapped line count = %d, want 0", len(wrapped))
	}
	assertWrappedLines(t, []wrappedLine{truncated}, []string{""}, []int{0})
}

func TestWordWrapper_shouldDropWhitespaceOnlyFragmentsAtWidthOneWhenTrimmed(t *testing.T) {
	lines := newWordWrapper(1, true).wrap(cellsFromLine(text.LineFromString(" a b ")), layout.Left)

	assertWrappedLines(t, lines, []string{"a", "b"}, []int{1, 1})
}

func TestReflow_shouldUseDisplayWidthForDoubleWidthSymbols(t *testing.T) {
	cells := cellsFromLine(text.LineFromString("aコンb"))
	wrapped := newWordWrapper(3, true).wrap(cells, layout.Left)
	truncated := newLineTruncator(3).truncate(cells, layout.Left)

	assertWrappedLines(t, wrapped, []string{"aコ", "ンb"}, []int{3, 3})
	assertWrappedLines(t, []wrappedLine{truncated}, []string{"aコ"}, []int{3})
}

func TestReflow_shouldPreserveAlignment(t *testing.T) {
	lines := newWordWrapper(3, true).wrap(cellsFromLine(text.LineFromString("abcdef")), layout.Right)

	for _, line := range lines {
		if line.alignment != layout.Right {
			t.Fatalf("alignment = %v, want Right", line.alignment)
		}
	}
}

func TestReflow_shouldKeepTrailingZeroWidthOnWrappedShortLineAndDropItWhenTruncated(t *testing.T) {
	cells := cellsFromLine(text.LineFromString("abc\u200b"))
	wrapped := newWordWrapper(3, true).wrap(cells, layout.Left)
	truncated := newLineTruncator(3).truncate(cells, layout.Left)

	assertWrappedLines(t, wrapped, []string{"abc\u200b"}, []int{3})
	assertWrappedLines(t, []wrappedLine{truncated}, []string{"abc"}, []int{3})
}

func assertWrappedLines(t *testing.T, lines []wrappedLine, wantSymbols []string, wantWidths []int) {
	t.Helper()
	gotSymbols := make([]string, 0, len(lines))
	gotWidths := make([]int, 0, len(lines))
	for _, line := range lines {
		gotSymbols = append(gotSymbols, cellsString(line.cells))
		gotWidths = append(gotWidths, line.width)
	}
	if !slices.Equal(gotSymbols, wantSymbols) {
		t.Fatalf("symbols = %#v, want %#v", gotSymbols, wantSymbols)
	}
	if !slices.Equal(gotWidths, wantWidths) {
		t.Fatalf("widths = %#v, want %#v", gotWidths, wantWidths)
	}
}

func cellsString(cells []buffer.Cell) string {
	var value strings.Builder
	for _, cell := range cells {
		value.WriteString(cell.Symbol)
	}
	return value.String()
}
