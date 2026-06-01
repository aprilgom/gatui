package buffer_test

import (
	"reflect"
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
	if got, want := buf.Lines(), []string{"ab", "c "}; !reflect.DeepEqual(got, want) {
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

func TestBuffer_Lines_shouldSkipHiddenFlagEmojiCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	buf.SetSymbol(0, 0, "🇺🇸")
	buf.SetSymbol(1, 0, " ")
	buf.SetSymbol(2, 0, "a")

	if got, want := buf.Lines(), []string{"🇺🇸a"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
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
	if !reflect.DeepEqual(diff, want) {
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
	if !reflect.DeepEqual(diff, want) {
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
	if !reflect.DeepEqual(diff, want) {
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
	if !reflect.DeepEqual(diff, want) {
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
	if !reflect.DeepEqual(diff, want) {
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
	if !reflect.DeepEqual(diff, want) {
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
	if !reflect.DeepEqual(diff, want) {
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
