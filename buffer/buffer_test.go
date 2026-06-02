package buffer_test

import (
	"slices"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

func TestWithLines_shouldCreateBlankPaddedBuffer(t *testing.T) {
	buf := buffer.WithLines([]string{"ab", "c"})

	if got, want := buf.Area, layout.NewRect(0, 0, 2, 2); got != want {
		t.Fatalf("area = %#v, want %#v", got, want)
	}
	if got, want := buf.Lines(), []string{"ab", "c "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestNewCell_shouldCreateCellWithSymbolAndDefaultStyle(t *testing.T) {
	cell := buffer.NewCell("x")

	if got, want := cell.Symbol, "x"; got != want {
		t.Fatalf("symbol = %q, want %q", got, want)
	}
	if got, want := cell.Style, style.NewStyle(); got != want {
		t.Fatalf("style = %#v, want %#v", got, want)
	}
	if got, want := cell.DiffOption, buffer.CellDiffNone; got != want {
		t.Fatalf("diff option = %#v, want %#v", got, want)
	}
	if got, want := cell.ForcedWidth, 0; got != want {
		t.Fatalf("forced width = %d, want %d", got, want)
	}
}

func TestCell_DisplaySymbol_shouldTreatEmptyAsSpace(t *testing.T) {
	cell := buffer.Cell{}

	if got, want := cell.DisplaySymbol(), " "; got != want {
		t.Fatalf("display symbol = %q, want %q", got, want)
	}
}

func TestCell_SetSymbolAndSetChar_shouldUpdateSymbol(t *testing.T) {
	cell := buffer.NewCell("a")

	cell.SetSymbol("bc")
	if got, want := cell.Symbol, "bc"; got != want {
		t.Fatalf("symbol after SetSymbol = %q, want %q", got, want)
	}

	cell.SetChar('コ')
	if got, want := cell.Symbol, "コ"; got != want {
		t.Fatalf("symbol after SetChar = %q, want %q", got, want)
	}
}

func TestCell_AppendSymbol_shouldAppendZeroWidthMarks(t *testing.T) {
	cell := buffer.NewCell("a")

	cell.AppendSymbol("\u200B")

	if got, want := cell.Symbol, "a\u200B"; got != want {
		t.Fatalf("symbol = %q, want %q", got, want)
	}
}

func TestCell_SetStyle_shouldPatchExistingStyle(t *testing.T) {
	cell := buffer.NewCell("x")
	cell.SetStyle(style.NewStyle().Fg(style.Red))
	cell.SetStyle(style.NewStyle().Bg(style.Blue).AddModifier(style.ModifierBold))

	want := style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold)
	if cell.Style != want {
		t.Fatalf("style = %#v, want %#v", cell.Style, want)
	}
}

func TestCell_Reset_shouldReturnEmptyDefaultCell(t *testing.T) {
	cell := buffer.NewCell("x")
	cell.SetStyle(style.NewStyle().Fg(style.Red))
	cell.SetDiffOption(buffer.CellDiffAlwaysUpdate)
	cell.SetForcedWidth(3)

	cell.Reset()

	want := buffer.NewCell(" ")
	if cell != want {
		t.Fatalf("cell = %#v, want %#v", cell, want)
	}
}

func TestCellWidth_shouldMatchRatatuiDisplayWidth(t *testing.T) {
	tests := []struct {
		value string
		want  int
	}{
		{value: "あ", want: 2},
		{value: "", want: 0},
		{value: "ﾞ", want: 1},
		{value: "ﾟ", want: 1},
		{value: "ｶﾞ", want: 2},
		{value: "ﾊﾟ", want: 2},
		{value: "aﾞ", want: 2},
		{value: "あﾞ", want: 3},
		{value: "ｶ゙", want: 1},
		{value: "ガ", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := buffer.CellWidth(tt.value); got != tt.want {
				t.Fatalf("CellWidth(%q) = %d, want %d", tt.value, got, tt.want)
			}
		})
	}
}

func TestCell_Width_shouldUseDisplayWidth(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
		want   int
	}{
		{name: "ascii", symbol: "a", want: 1},
		{name: "cjk", symbol: "コ", want: 2},
		{name: "flag", symbol: "🇺🇸", want: 2},
		{name: "empty", symbol: "", want: 1},
		{name: "halfwidth voiced mark", symbol: "ﾞ", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := buffer.NewCell(tt.symbol)

			if got := cell.Width(); got != tt.want {
				t.Fatalf("width = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCell_Width_shouldMatchCellWidthOfDisplaySymbol(t *testing.T) {
	cell := buffer.NewCell("あﾞ")

	if got, want := cell.Width(), buffer.CellWidth(cell.DisplaySymbol()); got != want {
		t.Fatalf("width = %d, want CellWidth(DisplaySymbol()) = %d", got, want)
	}
}

func TestCell_Width_shouldPreferForcedWidth(t *testing.T) {
	cell := buffer.NewCell("a")
	cell.SetForcedWidth(4)

	if got, want := cell.Width(), 4; got != want {
		t.Fatalf("width = %d, want %d", got, want)
	}

	cell.SetForcedWidth(0)
	if got, want := cell.Width(), 1; got != want {
		t.Fatalf("width after clearing forced width = %d, want %d", got, want)
	}
}

func TestCell_SetDiffOption_shouldStoreOption(t *testing.T) {
	cell := buffer.NewCell("x")

	cell.SetDiffOption(buffer.CellDiffSkip)
	if got, want := cell.DiffOption, buffer.CellDiffSkip; got != want {
		t.Fatalf("diff option = %#v, want %#v", got, want)
	}
}

func TestFilled_shouldFillAreaWithCell(t *testing.T) {
	cell := buffer.NewCell("x")
	cell.SetStyle(style.NewStyle().Fg(style.Red))

	buf := buffer.Filled(layout.NewRect(0, 0, 2, 2), cell)

	if got, want := buf.Area, layout.NewRect(0, 0, 2, 2); got != want {
		t.Fatalf("area = %#v, want %#v", got, want)
	}
	for i, got := range buf.Cells {
		if got != cell {
			t.Fatalf("cell %d = %#v, want %#v", i, got, cell)
		}
	}
}

func TestBuffer_IndexOfAndPosOf_translateCoordinates(t *testing.T) {
	area := layout.NewRect(200, 100, 50, 80)
	buf := buffer.Empty(area)

	firstIndex, ok := buf.IndexOf(200, 100)
	if !ok {
		t.Fatal("IndexOf first cell ok = false, want true")
	}
	if got, want := firstIndex, 0; got != want {
		t.Fatalf("IndexOf first cell = %d, want %d", got, want)
	}
	firstPos, ok := buf.PosOf(0)
	if !ok {
		t.Fatal("PosOf first cell ok = false, want true")
	}
	if got, want := firstPos, layout.NewPosition(200, 100); got != want {
		t.Fatalf("PosOf first cell = %#v, want %#v", got, want)
	}

	lastIndex, ok := buf.IndexOf(249, 179)
	if !ok {
		t.Fatal("IndexOf last cell ok = false, want true")
	}
	if got, want := lastIndex, len(buf.Cells)-1; got != want {
		t.Fatalf("IndexOf last cell = %d, want %d", got, want)
	}
	lastPos, ok := buf.PosOf(len(buf.Cells) - 1)
	if !ok {
		t.Fatal("PosOf last cell ok = false, want true")
	}
	if got, want := lastPos, layout.NewPosition(249, 179); got != want {
		t.Fatalf("PosOf last cell = %#v, want %#v", got, want)
	}
}

func TestBuffer_IndexOf_returnsFalseOutOfBounds(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(10, 10, 10, 10))
	tests := []struct {
		name string
		x    int
		y    int
	}{
		{name: "left", x: 9, y: 10},
		{name: "top", x: 10, y: 9},
		{name: "right", x: 20, y: 10},
		{name: "bottom", x: 10, y: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := buf.IndexOf(tt.x, tt.y); ok {
				t.Fatalf("IndexOf(%d,%d) = (%d,true), want ok=false", tt.x, tt.y, got)
			}
		})
	}

	var nilBuf *buffer.Buffer
	if got, ok := nilBuf.IndexOf(10, 10); ok {
		t.Fatalf("nil IndexOf = (%d,true), want ok=false", got)
	}
}

func TestBuffer_PosOf_returnsFalseOutOfBounds(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 10))

	for _, index := range []int{-1, 100} {
		if got, ok := buf.PosOf(index); ok {
			t.Fatalf("PosOf(%d) = (%#v,true), want ok=false", index, got)
		}
	}

	var nilBuf *buffer.Buffer
	if got, ok := nilBuf.PosOf(0); ok {
		t.Fatalf("nil PosOf = (%#v,true), want ok=false", got)
	}
}

func TestBuffer_IndexPosOf_handlesIndexesBeyondU16Max(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 256, 257))
	tests := []struct {
		x     int
		y     int
		index int
	}{
		{x: 255, y: 255, index: 65535},
		{x: 0, y: 256, index: 65536},
		{x: 1, y: 256, index: 65537},
		{x: 255, y: 256, index: 65791},
	}

	for _, tt := range tests {
		gotIndex, ok := buf.IndexOf(tt.x, tt.y)
		if !ok {
			t.Fatalf("IndexOf(%d,%d) ok = false, want true", tt.x, tt.y)
		}
		if gotIndex != tt.index {
			t.Fatalf("IndexOf(%d,%d) = %d, want %d", tt.x, tt.y, gotIndex, tt.index)
		}

		gotPos, ok := buf.PosOf(tt.index)
		if !ok {
			t.Fatalf("PosOf(%d) ok = false, want true", tt.index)
		}
		if want := layout.NewPosition(tt.x, tt.y); gotPos != want {
			t.Fatalf("PosOf(%d) = %#v, want %#v", tt.index, gotPos, want)
		}
	}
}

func TestBuffer_CellAt_matchesRatatuiCell(t *testing.T) {
	buf := buffer.WithLines([]string{"Hello", "World"})

	got, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatal("CellAt(0,0) ok = false, want true")
	}
	if want := buffer.NewCell("H"); got != want {
		t.Fatalf("CellAt(0,0) = %#v, want %#v", got, want)
	}
	if got, ok := buf.CellAt(10, 10); ok {
		t.Fatalf("CellAt(10,10) = (%#v,true), want ok=false", got)
	}
}

func TestBuffer_CellRef_allowsMutation(t *testing.T) {
	buf := buffer.WithLines([]string{"Hello", "World"})

	cell, ok := buf.CellRef(0, 0)
	if !ok {
		t.Fatal("CellRef(0,0) ok = false, want true")
	}
	cell.SetSymbol("Y")
	cell.SetStyle(style.NewStyle().Fg(style.Red))

	got, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatal("CellAt(0,0) after mutation ok = false, want true")
	}
	want := buffer.NewCell("Y")
	want.SetStyle(style.NewStyle().Fg(style.Red))
	if got != want {
		t.Fatalf("mutated cell = %#v, want %#v", got, want)
	}
}

func TestBuffer_CellRef_returnsFalseOutOfBounds(t *testing.T) {
	buf := buffer.WithLines([]string{"Hello", "World"})

	if got, ok := buf.CellRef(10, 10); ok {
		t.Fatalf("CellRef(10,10) = (%#v,true), want ok=false", got)
	}

	var nilBuf *buffer.Buffer
	if got, ok := nilBuf.CellRef(0, 0); ok {
		t.Fatalf("nil CellRef = (%#v,true), want ok=false", got)
	}
}

func TestBuffer_Reset_shouldResetAllCellsButKeepArea(t *testing.T) {
	area := layout.NewRect(1, 2, 2, 1)
	buf := buffer.Empty(area)
	cell := buffer.NewCell("x")
	cell.SetStyle(style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold))
	cell.SetDiffOption(buffer.CellDiffAlwaysUpdate)
	cell.SetForcedWidth(3)
	buf.SetCell(1, 2, cell)
	buf.SetCell(2, 2, cell)

	buf.Reset()

	if got := buf.Area; got != area {
		t.Fatalf("area = %#v, want %#v", got, area)
	}
	if got, want := len(buf.Cells), area.Width*area.Height; got != want {
		t.Fatalf("cell count = %d, want %d", got, want)
	}
	for i, got := range buf.Cells {
		if want := buffer.NewCell(" "); got != want {
			t.Fatalf("cell %d = %#v, want %#v", i, got, want)
		}
	}
}

func TestBuffer_Resize_shouldGrowWithBlankCells(t *testing.T) {
	buf := buffer.WithLines([]string{"ab"})

	buf.Resize(layout.NewRect(0, 0, 3, 2))

	if got, want := buf.Area, layout.NewRect(0, 0, 3, 2); got != want {
		t.Fatalf("area = %#v, want %#v", got, want)
	}
	if got, want := len(buf.Cells), 6; got != want {
		t.Fatalf("cell count = %d, want %d", got, want)
	}
	for i := 2; i < len(buf.Cells); i++ {
		if got, want := buf.Cells[i], buffer.NewCell(" "); got != want {
			t.Fatalf("new cell %d = %#v, want %#v", i, got, want)
		}
	}
}

func TestBuffer_Resize_shouldShrinkCells(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 2))

	buf.Resize(layout.NewRect(0, 0, 2, 1))

	if got, want := buf.Area, layout.NewRect(0, 0, 2, 1); got != want {
		t.Fatalf("area = %#v, want %#v", got, want)
	}
	if got, want := len(buf.Cells), 2; got != want {
		t.Fatalf("cell count = %d, want %d", got, want)
	}
}

func TestBuffer_Merge_shouldUnionAreasAndOverlayOther(t *testing.T) {
	tests := []struct {
		name     string
		one      layout.Rect
		two      layout.Rect
		wantArea layout.Rect
		want     []string
	}{
		{
			name:     "stacked",
			one:      layout.NewRect(0, 0, 2, 2),
			two:      layout.NewRect(0, 2, 2, 2),
			wantArea: layout.NewRect(0, 0, 2, 4),
			want:     []string{"11", "11", "22", "22"},
		},
		{
			name:     "offset",
			one:      layout.NewRect(2, 2, 2, 2),
			two:      layout.NewRect(0, 0, 2, 2),
			wantArea: layout.NewRect(0, 0, 4, 4),
			want:     []string{"22  ", "22  ", "  11", "  11"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			one := buffer.Filled(tt.one, buffer.NewCell("1"))
			two := buffer.Filled(tt.two, buffer.NewCell("2"))

			one.Merge(two)

			if got := one.Area; got != tt.wantArea {
				t.Fatalf("area = %#v, want %#v", got, tt.wantArea)
			}
			if got := one.Lines(); !slices.Equal(got, tt.want) {
				t.Fatalf("lines = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestBuffer_Merge_shouldHandleOffset(t *testing.T) {
	one := buffer.Filled(layout.NewRect(3, 3, 2, 2), buffer.NewCell("1"))
	two := buffer.Filled(layout.NewRect(1, 1, 3, 4), buffer.NewCell("2"))

	one.Merge(two)

	if got, want := one.Area, layout.NewRect(1, 1, 4, 4); got != want {
		t.Fatalf("area = %#v, want %#v", got, want)
	}
	if got, want := one.Lines(), []string{"222 ", "222 ", "2221", "2221"}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_Merge_shouldPreserveDiffOptions(t *testing.T) {
	oneCell := buffer.NewCell("1")
	oneCell.SetDiffOption(buffer.CellDiffSkip)
	oneCell.SetForcedWidth(2)
	oneCell.SetStyle(style.NewStyle().Fg(style.Red))
	twoCell := buffer.NewCell("2")
	twoCell.SetDiffOption(buffer.CellDiffAlwaysUpdate)
	twoCell.SetForcedWidth(3)
	twoCell.SetStyle(style.NewStyle().Bg(style.Blue))
	one := buffer.Filled(layout.NewRect(0, 0, 2, 2), oneCell)
	two := buffer.Filled(layout.NewRect(0, 1, 2, 2), twoCell)

	one.Merge(two)

	cases := []struct {
		x, y int
		want buffer.Cell
	}{
		{x: 0, y: 0, want: oneCell},
		{x: 0, y: 1, want: twoCell},
		{x: 0, y: 2, want: twoCell},
	}
	for _, tt := range cases {
		got, ok := one.CellAt(tt.x, tt.y)
		if !ok {
			t.Fatalf("missing cell at %d,%d", tt.x, tt.y)
		}
		if got != tt.want {
			t.Fatalf("cell at %d,%d = %#v, want %#v", tt.x, tt.y, got, tt.want)
		}
	}
}

func TestBuffer_Merge_shouldMakeDiffIdempotent(t *testing.T) {
	prev := buffer.Filled(layout.NewRect(0, 0, 2, 2), buffer.NewCell("1"))
	next := buffer.Filled(layout.NewRect(0, 0, 2, 2), buffer.NewCell("2"))

	prev.Merge(next)

	if diff := prev.Diff(next); len(diff) != 0 {
		t.Fatalf("diff = %#v, want empty", diff)
	}
}

func TestBuffer_Lines_shouldSkipHiddenFlagEmojiCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	buf.SetSymbol(0, 0, "🇺🇸")
	buf.SetSymbol(1, 0, " ")
	buf.SetSymbol(2, 0, "a")

	if got, want := buf.Lines(), []string{"🇺🇸a"}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetString_shouldWriteAsciiWithStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	red := style.NewStyle().Fg(style.Red)

	endX, endY := buf.SetString(1, 0, "abc", red)

	if endX != 4 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (4,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{" abc "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
	for x := 1; x <= 3; x++ {
		cell, ok := buf.CellAt(x, 0)
		if !ok {
			t.Fatalf("missing cell at %d,0", x)
		}
		if cell.Style != red {
			t.Fatalf("cell %d style = %#v, want %#v", x, cell.Style, red)
		}
	}
}

func TestBuffer_SetString_shouldClipAtRightBoundary(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))

	endX, endY := buf.SetString(2, 0, "abcd", style.NewStyle())

	if endX != 4 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (4,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"  ab"}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetStringN_shouldClipToMaxWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	endX, endY := buf.SetStringN(0, 0, "abcdef", 3, style.NewStyle())

	if endX != 3 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (3,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"abc  "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetString_shouldHandleWideGraphemes(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	red := style.NewStyle().Fg(style.Red)

	endX, endY := buf.SetString(0, 0, "コ🙂a", red)

	if endX != 5 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (5,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"コ🙂a"}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
	for _, x := range []int{0, 2, 4} {
		cell, ok := buf.CellAt(x, 0)
		if !ok {
			t.Fatalf("missing cell at %d,0", x)
		}
		if cell.Style != red {
			t.Fatalf("cell %d style = %#v, want %#v", x, cell.Style, red)
		}
	}
	for _, x := range []int{1, 3} {
		cell, ok := buf.CellAt(x, 0)
		if !ok {
			t.Fatalf("missing cell at %d,0", x)
		}
		if want := buffer.NewCell(" "); cell != want {
			t.Fatalf("trailing cell %d = %#v, want %#v", x, cell, want)
		}
	}
}

func TestBuffer_SetString_shouldNotWriteWideGraphemeIntoOneRemainingCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))

	endX, endY := buf.SetString(2, 0, "コ", style.NewStyle().Fg(style.Red))

	if endX != 2 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (2,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"   "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetStringN_shouldAppendZeroWidthMarkToPreviousCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))

	endX, endY := buf.SetStringN(0, 0, "a\u200Bb", 3, style.NewStyle())

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
	if got, want := buf.Lines(), []string{"a\u200Bb "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetStringN_shouldTreatHalfwidthVoicedMarksAsWidthOne(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))

	endX, endY := buf.SetStringN(0, 0, "aﾞb", 3, style.NewStyle())

	if endX != 3 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (3,0)", endX, endY)
	}
	cell, ok := buf.CellAt(1, 0)
	if !ok {
		t.Fatal("missing cell at 1,0")
	}
	if got, want := cell.Symbol, "ﾞ"; got != want {
		t.Fatalf("symbol = %q, want %q", got, want)
	}
	if got, want := buf.Lines(), []string{"aﾞb"}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetStringN_shouldSkipControlSequences(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	endX, endY := buf.SetStringN(0, 0, "a\nb\tc", 5, style.NewStyle())

	if endX != 3 || endY != 0 {
		t.Fatalf("end = (%d,%d), want (3,0)", endX, endY)
	}
	if got, want := buf.Lines(), []string{"abc  "}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetStringN_controlSequenceRenderedFull(t *testing.T) {
	text := "I \x1b[0;36mwas\x1b[0m here!"
	buf := buffer.Filled(layout.NewRect(0, 0, 25, 3), buffer.NewCell("x"))

	buf.SetString(1, 1, text, style.NewStyle())

	if got, want := buf.Lines(), []string{
		"xxxxxxxxxxxxxxxxxxxxxxxxx",
		"xI [0;36mwas[0m here!xxxx",
		"xxxxxxxxxxxxxxxxxxxxxxxxx",
	}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetStringN_controlSequenceRenderedPartially(t *testing.T) {
	text := "I \x1b[0;36mwas\x1b[0m here!"
	buf := buffer.Filled(layout.NewRect(0, 0, 11, 3), buffer.NewCell("x"))

	buf.SetString(1, 1, text, style.NewStyle())

	if got, want := buf.Lines(), []string{
		"xxxxxxxxxxx",
		"xI [0;36mwa",
		"xxxxxxxxxxx",
	}; !slices.Equal(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestBuffer_SetString_shouldNoOpOutsideArea(t *testing.T) {
	buf := buffer.WithLines([]string{"abc"})

	tests := []struct {
		name string
		x    int
		y    int
	}{
		{name: "left", x: -1, y: 0},
		{name: "right", x: 3, y: 0},
		{name: "above", x: 0, y: -1},
		{name: "below", x: 0, y: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endX, endY := buf.SetString(tt.x, tt.y, "xyz", style.NewStyle().Fg(style.Red))

			if endX != tt.x || endY != tt.y {
				t.Fatalf("end = (%d,%d), want (%d,%d)", endX, endY, tt.x, tt.y)
			}
			if got, want := buf.Lines(), []string{"abc"}; !slices.Equal(got, want) {
				t.Fatalf("lines = %#v, want %#v", got, want)
			}
		})
	}
}

func TestSetStyleHelpers_shouldPatchExistingCells(t *testing.T) {
	buf := buffer.WithLines([]string{"ab"})
	buf.SetFg(layout.NewRect(0, 0, 1, 1), style.Red)
	buf.SetModifier(layout.NewRect(0, 0, 1, 1), style.ModifierBold)

	cell, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatal("expected cell")
	}
	want := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold)
	if cell.Style != want {
		t.Fatalf("style = %#v, want %#v", cell.Style, want)
	}
}

func TestBuffer_Diff_shouldReturnNoDiffForIdenticalBuffers(t *testing.T) {
	buf := buffer.WithLines([]string{"hello"})

	diff := buf.Diff(buf)

	if len(diff) != 0 {
		t.Fatalf("diff = %#v, want empty", diff)
	}
}

func TestBuffer_Diff_shouldReturnSingleCellChange(t *testing.T) {
	prev := buffer.WithLines([]string{"hello"})
	next := buffer.WithLines([]string{"hallo"})

	diff := prev.Diff(next)

	want := []buffer.CellDiff{{X: 1, Y: 0, Cell: buffer.NewCell("a")}}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_shouldReturnAllChangedCells(t *testing.T) {
	prev := buffer.WithLines([]string{"aaa"})
	next := buffer.WithLines([]string{"bbb"})

	diff := prev.Diff(next)

	want := []buffer.CellDiff{
		{X: 0, Y: 0, Cell: buffer.NewCell("b")},
		{X: 1, Y: 0, Cell: buffer.NewCell("b")},
		{X: 2, Y: 0, Cell: buffer.NewCell("b")},
	}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_shouldSkipCellsMarkedSkip(t *testing.T) {
	prev := buffer.WithLines([]string{"abc"})
	next := buffer.WithLines([]string{"xyz"})
	cell, ok := next.CellAt(1, 0)
	if !ok {
		t.Fatal("expected cell")
	}
	cell.SetDiffOption(buffer.CellDiffSkip)
	next.SetCell(1, 0, cell)

	diff := prev.Diff(next)

	want := []buffer.CellDiff{
		{X: 0, Y: 0, Cell: buffer.NewCell("x")},
		{X: 2, Y: 0, Cell: buffer.NewCell("z")},
	}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_shouldAlwaysUpdateMarkedCells(t *testing.T) {
	prev := buffer.WithLines([]string{"abc"})
	next := buffer.WithLines([]string{"abc"})
	cell, ok := next.CellAt(1, 0)
	if !ok {
		t.Fatal("expected cell")
	}
	cell.SetDiffOption(buffer.CellDiffAlwaysUpdate)
	next.SetCell(1, 0, cell)

	diff := prev.Diff(next)

	want := []buffer.CellDiff{{X: 1, Y: 0, Cell: cell}}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_shouldSkipTrailingForcedWidthCells(t *testing.T) {
	prev := buffer.WithLines([]string{"abcd"})
	next := buffer.WithLines([]string{"xbcd"})
	cell, ok := next.CellAt(0, 0)
	if !ok {
		t.Fatal("expected cell")
	}
	cell.SetDiffOption(buffer.CellDiffForcedWidth)
	cell.SetForcedWidth(2)
	next.SetCell(0, 0, cell)

	diff := prev.Diff(next)

	want := []buffer.CellDiff{{X: 0, Y: 0, Cell: cell}}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_shouldHandleMultiWidthCells(t *testing.T) {
	prev := buffer.WithLines([]string{
		"┌Title─┐  ",
		"└──────┘  ",
	})
	next := buffer.WithLines([]string{
		"┌称号──┐  ",
		"└──────┘  ",
	})

	diff := prev.Diff(next)

	want := []buffer.CellDiff{
		{X: 1, Y: 0, Cell: buffer.NewCell("称")},
		{X: 3, Y: 0, Cell: buffer.NewCell("号")},
		{X: 5, Y: 0, Cell: buffer.NewCell("─")},
	}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_shouldHandleMultiWidthOffset(t *testing.T) {
	prev := buffer.WithLines([]string{"┌称号──┐"})
	next := buffer.WithLines([]string{"┌─称号─┐"})

	diff := prev.Diff(next)

	want := []buffer.CellDiff{
		{X: 1, Y: 0, Cell: buffer.NewCell("─")},
		{X: 2, Y: 0, Cell: buffer.NewCell("称")},
		{X: 4, Y: 0, Cell: buffer.NewCell("号")},
	}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_clearsTrailingCellForWideGrapheme(t *testing.T) {
	prev := buffer.WithLines([]string{"ab"})
	next := buffer.WithLines([]string{"  "})
	next.SetString(0, 0, "⌨️", style.NewStyle())

	diff := prev.Diff(next)

	want := []buffer.CellDiff{
		{X: 0, Y: 0, Cell: buffer.NewCell("⌨️")},
		{X: 1, Y: 0, Cell: buffer.NewCell(" ")},
	}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_ignoresStyleOnlyChangesInTrailingCells(t *testing.T) {
	prev := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	prev.SetString(0, 0, "  ", style.NewStyle().Fg(style.LightBlue))
	prev.SetString(2, 0, "x", style.NewStyle())

	next := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	next.SetString(0, 0, "⚠️", style.NewStyle().Fg(style.Reset))
	next.SetString(2, 0, "x", style.NewStyle())

	diff := prev.Diff(next)

	want := []buffer.CellDiff{{X: 0, Y: 0, Cell: buffer.NewCell("⚠️")}}
	want[0].Cell.SetStyle(style.NewStyle().Fg(style.Reset))
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_vs16TrailingCellUnchanged(t *testing.T) {
	prev := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	prev.SetString(0, 0, "⌨️", style.NewStyle())
	prev.SetString(2, 0, "ab", style.NewStyle())

	next := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	next.SetString(0, 0, "⌨️", style.NewStyle().Fg(style.Red))
	next.SetString(2, 0, "ab", style.NewStyle())

	diff := prev.Diff(next)

	wantCell := buffer.NewCell("⌨️")
	wantCell.SetStyle(style.NewStyle().Fg(style.Red))
	want := []buffer.CellDiff{{X: 0, Y: 0, Cell: wantCell}}
	if !slices.Equal(diff, want) {
		t.Fatalf("diff = %#v, want %#v", diff, want)
	}
}

func TestBuffer_Diff_shouldPanicForIncompatibleAreas(t *testing.T) {
	tests := []struct {
		name string
		next layout.Rect
	}{
		{name: "x", next: layout.NewRect(1, 0, 3, 1)},
		{name: "y", next: layout.NewRect(0, 1, 3, 1)},
		{name: "width", next: layout.NewRect(0, 0, 4, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev := buffer.Empty(layout.NewRect(0, 0, 3, 1))
			next := buffer.Empty(tt.next)
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic")
				}
			}()

			_ = prev.Diff(next)
		})
	}
}
