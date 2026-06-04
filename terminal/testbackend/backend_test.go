package testbackend

import (
	"os"
	"os/exec"
	"slices"
	"strconv"
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/terminal"
)

func TestTestBackend_WindowSize_shouldReturnBufferSizeAndDefaultPixels(t *testing.T) {
	backend := New(80, 24)

	got, err := backend.WindowSize()
	if err != nil {
		t.Fatalf("WindowSize() error = %v", err)
	}

	want := terminal.WindowSize{
		ColumnsRows: layout.Size{Width: 80, Height: 24},
		Pixels:      layout.Size{Width: 640, Height: 480},
	}
	if got != want {
		t.Fatalf("WindowSize() = %+v, want %+v", got, want)
	}
}

func TestBackend_New_shouldCreateBlankBufferAtRequestedSize(t *testing.T) {
	backend := New(10, 2)

	gotSize, err := backend.Size()
	if err != nil {
		t.Fatalf("Size() error = %v", err)
	}
	if want := (layout.Size{Width: 10, Height: 2}); gotSize != want {
		t.Fatalf("Size() = %+v, want %+v", gotSize, want)
	}
	backend.AssertBufferLines(t, []string{
		"          ",
		"          ",
	})
	backend.AssertScrollbackEmpty(t)
}

func TestNoScrollBackend_WindowSize_shouldDelegateToWrappedBackend(t *testing.T) {
	backend := NewNoScroll(10, 5)

	got, err := backend.WindowSize()
	if err != nil {
		t.Fatalf("WindowSize() error = %v", err)
	}

	want := terminal.WindowSize{
		ColumnsRows: layout.Size{Width: 10, Height: 5},
		Pixels:      layout.Size{Width: 640, Height: 480},
	}
	if got != want {
		t.Fatalf("WindowSize() = %+v, want %+v", got, want)
	}
}

func TestBackend_SetSize_shouldResizeBackingBuffer(t *testing.T) {
	backend := New(10, 2)

	backend.SetSize(5, 5)

	gotSize, err := backend.Size()
	if err != nil {
		t.Fatalf("Size() error = %v", err)
	}
	if want := (layout.Size{Width: 5, Height: 5}); gotSize != want {
		t.Fatalf("Size() = %+v, want %+v", gotSize, want)
	}
	backend.AssertBufferLines(t, []string{
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
	})
}

func TestTestBackend_Buffer_shouldReturnCurrentBuffer(t *testing.T) {
	backend := WithLines([]string{
		"abc",
		"def",
	})

	got := backend.Buffer()

	want := buffer.WithLines([]string{
		"abc",
		"def",
	})
	if !buffersEqual(got, want) {
		t.Fatalf("Buffer() = %#v, want %#v", got, want)
	}
}

func TestBackend_Draw_shouldUpdateOnlyDiffCellsAndRecordCalls(t *testing.T) {
	backend := New(10, 2)
	cell := buffer.NewCell("a")

	if err := backend.Draw([]buffer.CellDiff{{X: 0, Y: 0, Cell: cell}}); err != nil {
		t.Fatalf("Draw() error = %v", err)
	}
	if err := backend.Draw([]buffer.CellDiff{{X: 0, Y: 1, Cell: cell}}); err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	backend.AssertBufferLines(t, []string{
		"a         ",
		"a         ",
	})
	gotDraws := backend.Draws()
	wantDraws := [][]buffer.CellDiff{
		{{X: 0, Y: 0, Cell: cell}},
		{{X: 0, Y: 1, Cell: cell}},
	}
	if !cellDiffBatchesEqual(gotDraws, wantDraws) {
		t.Fatalf("Draws() = %#v, want %#v", gotDraws, wantDraws)
	}
}

func cellDiffBatchesEqual(a, b [][]buffer.CellDiff) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !slices.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestTestBackend_Scrollback_shouldReturnScrollbackBuffer(t *testing.T) {
	backend := WithLines([]string{
		"aaaa",
		"bbbb",
		"cccc",
	})

	if err := backend.ScrollRegionUp(0, 3, 2); err != nil {
		t.Fatalf("ScrollRegionUp() error = %v", err)
	}

	got := backend.Scrollback()
	want := buffer.WithLines([]string{
		"aaaa",
		"bbbb",
	})
	if !buffersEqual(got, want) {
		t.Fatalf("Scrollback() = %#v, want %#v", got, want)
	}
}

func buffersEqual(a, b *buffer.Buffer) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Area == b.Area && slices.Equal(a.Cells, b.Cells)
}

func TestTestBackend_String_shouldRenderQuotedBufferRows(t *testing.T) {
	backend := WithLines([]string{
		"aaaa",
		"aaaa",
	})

	got := backend.String()
	want := "\"aaaa\"\n\"aaaa\"\n"
	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestTestBackend_String_shouldIncludeTrailingNewlinePerRow(t *testing.T) {
	backend := WithLines([]string{
		"aa",
		"bb",
	})

	got := backend.String()
	want := "\"aa\"\n\"bb\"\n"
	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestTestBackend_String_shouldShowWideCellOverwrites(t *testing.T) {
	backend := New(3, 1)
	backend.Buffer().SetSymbol(0, 0, "界")
	backend.Buffer().SetSymbol(1, 0, "x")

	got := backend.String()
	want := "\"界x \"\n"
	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestNoScrollBackend_String_shouldDelegateToWrappedBackend(t *testing.T) {
	backend := WithLinesNoScroll([]string{
		"aa",
		"bb",
	})

	got := backend.String()
	want := "\"aa\"\n\"bb\"\n"
	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestTestBackend_AssertBufferLines_shouldPassForMatchingLines(t *testing.T) {
	backend := WithLines([]string{
		"abc",
		"def",
	})

	backend.AssertBufferLines(t, []string{
		"abc",
		"def",
	})
}

func TestTestBackend_AssertBufferLines_shouldFailForMismatchedLines(t *testing.T) {
	if os.Getenv("GATUI_ASSERT_BUFFER_LINES_MISMATCH") == "1" {
		backend := WithLines([]string{
			"abc",
			"def",
		})
		backend.AssertBufferLines(t, []string{
			"abc",
			"xyz",
		})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestTestBackend_AssertBufferLines_shouldFailForMismatchedLines$")
	cmd.Env = append(os.Environ(), "GATUI_ASSERT_BUFFER_LINES_MISMATCH=1")
	if err := cmd.Run(); err == nil {
		t.Fatal("AssertBufferLines() unexpectedly passed")
	}
}

func TestTestBackend_AssertScrollbackLines_shouldPassForMatchingLines(t *testing.T) {
	backend := WithLines([]string{
		"aaaa",
		"bbbb",
		"cccc",
	})

	if err := backend.ScrollRegionUp(0, 3, 1); err != nil {
		t.Fatalf("ScrollRegionUp() error = %v", err)
	}

	backend.AssertScrollbackLines(t, []string{
		"aaaa",
	})
}

func TestTestBackend_AssertScrollbackLines_shouldFailForMismatchedLines(t *testing.T) {
	if os.Getenv("GATUI_ASSERT_SCROLLBACK_LINES_MISMATCH") == "1" {
		backend := New(10, 2)
		backend.AssertScrollbackLines(t, []string{
			"aaaaaaaaaa",
			"aaaaaaaaaa",
		})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestTestBackend_AssertScrollbackLines_shouldFailForMismatchedLines$")
	cmd.Env = append(os.Environ(), "GATUI_ASSERT_SCROLLBACK_LINES_MISMATCH=1")
	if err := cmd.Run(); err == nil {
		t.Fatal("AssertScrollbackLines() unexpectedly passed")
	}
}

func TestTestBackend_AssertScrollbackEmpty_shouldPassWhenEmpty(t *testing.T) {
	backend := New(4, 2)

	backend.AssertScrollbackEmpty(t)
}

func TestTestBackend_AssertCursorPosition_shouldPassForMatchingPosition(t *testing.T) {
	backend := New(4, 2)
	pos := layout.Position{X: 2, Y: 1}
	if err := backend.SetCursorPosition(pos); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	backend.AssertCursorPosition(t, pos)
}

func TestBackend_ClearRegion_beforeCursor(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 5, Y: 3}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.ClearRegion(terminal.ClearBeforeCursor); err != nil {
		t.Fatalf("ClearRegion(ClearBeforeCursor) error = %v", err)
	}

	want := []string{
		"          ",
		"          ",
		"          ",
		"      aaaa",
		"aaaaaaaaaa",
	}
	if got := backend.Lines(); !slices.Equal(got, want) {
		t.Fatalf("Lines() = %#v, want %#v", got, want)
	}
}

func TestBackend_Clear_shouldRecordClearAllAndBlankBuffer(t *testing.T) {
	backend := New(4, 2)
	cell := buffer.NewCell("a")
	if err := backend.Draw([]buffer.CellDiff{
		{X: 0, Y: 0, Cell: cell},
		{X: 0, Y: 1, Cell: cell},
	}); err != nil {
		t.Fatalf("Draw() error = %v", err)
	}

	if err := backend.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	if got, want := backend.ClearCount(), 1; got != want {
		t.Fatalf("ClearCount() = %d, want %d", got, want)
	}
	if got, want := backend.ClearRegions(), []terminal.ClearType{terminal.ClearAll}; !slices.Equal(got, want) {
		t.Fatalf("ClearRegions() = %#v, want %#v", got, want)
	}
	backend.AssertBufferLines(t, []string{
		"    ",
		"    ",
	})
}

func TestBackend_ClearRegion_all(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	})

	if err := backend.ClearRegion(terminal.ClearAll); err != nil {
		t.Fatalf("ClearRegion(ClearAll) error = %v", err)
	}

	if got, want := backend.ClearRegions(), []terminal.ClearType{terminal.ClearAll}; !slices.Equal(got, want) {
		t.Fatalf("ClearRegions() = %#v, want %#v", got, want)
	}
	backend.AssertBufferLines(t, []string{
		"          ",
		"          ",
		"          ",
		"          ",
		"          ",
	})
}

func TestBackend_ClearRegion_untilNewLine(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 3, Y: 0}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.ClearRegion(terminal.ClearUntilNewLine); err != nil {
		t.Fatalf("ClearRegion(ClearUntilNewLine) error = %v", err)
	}

	want := []string{
		"aaa       ",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
		"aaaaaaaaaa",
	}
	if got := backend.Lines(); !slices.Equal(got, want) {
		t.Fatalf("Lines() = %#v, want %#v", got, want)
	}
}

func TestTestBackend_AppendLines_notAtLastLineMovesCursorDownWithoutScrollback(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 0}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	for _, want := range []layout.Position{
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 3},
		{X: 4, Y: 4},
	} {
		if err := backend.AppendLines(1); err != nil {
			t.Fatalf("AppendLines(1) error = %v", err)
		}
		backend.AssertCursorPosition(t, want)
	}

	backend.AssertBufferLines(t, []string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	backend.AssertScrollbackEmpty(t)
}

func TestTestBackend_AppendLines_atLastLineScrollsOneLineToScrollback(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 4}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.AppendLines(1); err != nil {
		t.Fatalf("AppendLines(1) error = %v", err)
	}

	backend.AssertCursorPosition(t, layout.Position{X: 1, Y: 4})
	backend.AssertBufferLines(t, []string{
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
		"          ",
	})
	backend.AssertScrollbackLines(t, []string{"aaaaaaaaaa"})
}

func TestTestBackend_AppendLines_multipleLinesNotAtLastLine(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 0}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.AppendLines(4); err != nil {
		t.Fatalf("AppendLines(4) error = %v", err)
	}

	backend.AssertCursorPosition(t, layout.Position{X: 1, Y: 4})
	backend.AssertBufferLines(t, []string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	backend.AssertScrollbackEmpty(t)
}

func TestTestBackend_AppendLines_multipleLinesPastLastLine(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 3}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.AppendLines(3); err != nil {
		t.Fatalf("AppendLines(3) error = %v", err)
	}

	backend.AssertCursorPosition(t, layout.Position{X: 1, Y: 4})
	backend.AssertBufferLines(t, []string{
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
		"          ",
		"          ",
	})
	backend.AssertScrollbackLines(t, []string{"aaaaaaaaaa", "bbbbbbbbbb"})
}

func TestTestBackend_AppendLines_cursorAtEndAppendsHeightLines(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 4}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.AppendLines(5); err != nil {
		t.Fatalf("AppendLines(5) error = %v", err)
	}

	backend.AssertCursorPosition(t, layout.Position{X: 1, Y: 4})
	backend.AssertBufferLines(t, []string{
		"          ",
		"          ",
		"          ",
		"          ",
		"          ",
	})
	backend.AssertScrollbackLines(t, []string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
}

func TestTestBackend_AppendLines_moreThanHeightLinesKeepsOnlyVisibleTail(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 4}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.AppendLines(8); err != nil {
		t.Fatalf("AppendLines(8) error = %v", err)
	}

	backend.AssertCursorPosition(t, layout.Position{X: 1, Y: 4})
	backend.AssertBufferLines(t, []string{
		"          ",
		"          ",
		"          ",
		"          ",
		"          ",
	})
	backend.AssertScrollbackLines(t, []string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
		"cccccccccc",
		"dddddddddd",
		"eeeeeeeeee",
		"          ",
		"          ",
		"          ",
	})
}

func TestTestBackend_AppendLines_truncatesLargeScrollbackToTail(t *testing.T) {
	backend := New(10, 5)
	const rowCount = 65535 + 10

	for row := 0; row <= rowCount; row++ {
		if row > 4 {
			if err := backend.SetCursorPosition(layout.Position{X: 0, Y: 4}); err != nil {
				t.Fatalf("SetCursorPosition() error = %v", err)
			}
			if err := backend.AppendLines(1); err != nil {
				t.Fatalf("AppendLines(1) error = %v", err)
			}
		}
		writeRow(t, backend, 4, row)
	}

	backend.AssertBufferLines(t, []string{
		"     65541",
		"     65542",
		"     65543",
		"     65544",
		"     65545",
	})

	scrollback := backend.ScrollbackLines()
	if got, want := len(scrollback), 65535; got != want {
		t.Fatalf("len(ScrollbackLines()) = %d, want %d", got, want)
	}
	if got, want := scrollback[:5], []string{
		"         6",
		"         7",
		"         8",
		"         9",
		"        10",
	}; !slices.Equal(got, want) {
		t.Fatalf("first scrollback lines = %#v, want %#v", got, want)
	}
	if got, want := scrollback[len(scrollback)-5:], []string{
		"     65536",
		"     65537",
		"     65538",
		"     65539",
		"     65540",
	}; !slices.Equal(got, want) {
		t.Fatalf("last scrollback lines = %#v, want %#v", got, want)
	}
}

func writeRow(t testing.TB, backend *Backend, y, row int) {
	t.Helper()
	line := []rune(formatRow(row))
	diffs := make([]buffer.CellDiff, len(line))
	for x, r := range line {
		diffs[x] = buffer.CellDiff{X: x, Y: y, Cell: buffer.NewCell(string(r))}
	}
	if err := backend.Draw(diffs); err != nil {
		t.Fatalf("Draw() error = %v", err)
	}
}

func formatRow(row int) string {
	value := "          " + strconv.Itoa(row)
	return value[len(value)-10:]
}

func TestTestBackend_AppendLines_zeroNoop(t *testing.T) {
	backend := WithLines([]string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
	})
	if err := backend.SetCursorPosition(layout.Position{X: 1, Y: 1}); err != nil {
		t.Fatalf("SetCursorPosition() error = %v", err)
	}

	if err := backend.AppendLines(0); err != nil {
		t.Fatalf("AppendLines(0) error = %v", err)
	}

	backend.AssertCursorPosition(t, layout.Position{X: 1, Y: 1})
	backend.AssertBufferLines(t, []string{
		"aaaaaaaaaa",
		"bbbbbbbbbb",
	})
	backend.AssertScrollbackEmpty(t)
	if got, want := backend.AppendLinesCalls(), []int{0}; !slices.Equal(got, want) {
		t.Fatalf("AppendLinesCalls() = %#v, want %#v", got, want)
	}
}

func TestTestBackend_ScrollRegionUp_table(t *testing.T) {
	const (
		a = "aaaaaaaaaa"
		b = "bbbbbbbbbb"
		c = "cccccccccc"
		d = "dddddddddd"
		e = "eeeeeeeeee"
		s = "          "
	)
	tests := []struct {
		name                string
		startY, endY, count int
		wantScrollback      []string
		wantLines           []string
	}{
		{name: "full screen zero", startY: 0, endY: 5, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "full screen partial", startY: 0, endY: 5, count: 2, wantScrollback: []string{a, b}, wantLines: []string{c, d, e, s, s}},
		{name: "full screen height", startY: 0, endY: 5, count: 5, wantScrollback: []string{a, b, c, d, e}, wantLines: []string{s, s, s, s, s}},
		{name: "full screen past height", startY: 0, endY: 5, count: 7, wantScrollback: []string{a, b, c, d, e, s, s}, wantLines: []string{s, s, s, s, s}},
		{name: "top partial zero", startY: 0, endY: 3, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "top partial scroll", startY: 0, endY: 3, count: 2, wantScrollback: []string{a, b}, wantLines: []string{c, s, s, d, e}},
		{name: "top partial height", startY: 0, endY: 3, count: 3, wantScrollback: []string{a, b, c}, wantLines: []string{s, s, s, d, e}},
		{name: "top partial past height", startY: 0, endY: 3, count: 4, wantScrollback: []string{a, b, c, s}, wantLines: []string{s, s, s, d, e}},
		{name: "middle partial zero", startY: 1, endY: 4, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "middle partial scroll", startY: 1, endY: 4, count: 2, wantLines: []string{a, d, s, s, e}},
		{name: "middle partial height", startY: 1, endY: 4, count: 3, wantLines: []string{a, s, s, s, e}},
		{name: "middle partial past height", startY: 1, endY: 4, count: 4, wantLines: []string{a, s, s, s, e}},
		{name: "empty at top zero", startY: 0, endY: 0, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "empty at top scroll", startY: 0, endY: 0, count: 2, wantScrollback: []string{s, s}, wantLines: []string{a, b, c, d, e}},
		{name: "empty middle zero", startY: 2, endY: 2, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "empty middle scroll", startY: 2, endY: 2, count: 2, wantLines: []string{a, b, c, d, e}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := WithLines([]string{a, b, c, d, e})
			if err := backend.SetCursorPosition(layout.Position{X: 4, Y: 3}); err != nil {
				t.Fatalf("SetCursorPosition() error = %v", err)
			}

			if err := backend.ScrollRegionUp(tt.startY, tt.endY, tt.count); err != nil {
				t.Fatalf("ScrollRegionUp(%d, %d, %d) error = %v", tt.startY, tt.endY, tt.count, err)
			}

			backend.AssertScrollbackLines(t, tt.wantScrollback)
			backend.AssertBufferLines(t, tt.wantLines)
			backend.AssertCursorPosition(t, layout.Position{X: 4, Y: 3})
			if got, want := backend.ScrollRegionUpCalls(), [][3]int{{tt.startY, tt.endY, tt.count}}; !slices.Equal(got, want) {
				t.Fatalf("ScrollRegionUpCalls() = %#v, want %#v", got, want)
			}
		})
	}
}

func TestTestBackend_ScrollRegionDown_table(t *testing.T) {
	const (
		a = "aaaaaaaaaa"
		b = "bbbbbbbbbb"
		c = "cccccccccc"
		d = "dddddddddd"
		e = "eeeeeeeeee"
		s = "          "
	)
	tests := []struct {
		name                string
		startY, endY, count int
		wantLines           []string
	}{
		{name: "full screen zero", startY: 0, endY: 5, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "full screen partial", startY: 0, endY: 5, count: 2, wantLines: []string{s, s, a, b, c}},
		{name: "full screen height", startY: 0, endY: 5, count: 5, wantLines: []string{s, s, s, s, s}},
		{name: "full screen past height", startY: 0, endY: 5, count: 7, wantLines: []string{s, s, s, s, s}},
		{name: "top partial zero", startY: 0, endY: 3, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "top partial scroll", startY: 0, endY: 3, count: 2, wantLines: []string{s, s, a, d, e}},
		{name: "top partial height", startY: 0, endY: 3, count: 3, wantLines: []string{s, s, s, d, e}},
		{name: "top partial past height", startY: 0, endY: 3, count: 4, wantLines: []string{s, s, s, d, e}},
		{name: "middle partial zero", startY: 1, endY: 4, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "middle partial scroll", startY: 1, endY: 4, count: 2, wantLines: []string{a, s, s, b, e}},
		{name: "middle partial height", startY: 1, endY: 4, count: 3, wantLines: []string{a, s, s, s, e}},
		{name: "middle partial past height", startY: 1, endY: 4, count: 4, wantLines: []string{a, s, s, s, e}},
		{name: "empty at top zero", startY: 0, endY: 0, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "empty at top scroll", startY: 0, endY: 0, count: 2, wantLines: []string{a, b, c, d, e}},
		{name: "empty middle zero", startY: 2, endY: 2, count: 0, wantLines: []string{a, b, c, d, e}},
		{name: "empty middle scroll", startY: 2, endY: 2, count: 2, wantLines: []string{a, b, c, d, e}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := WithLines([]string{a, b, c, d, e})
			if err := backend.SetCursorPosition(layout.Position{X: 4, Y: 3}); err != nil {
				t.Fatalf("SetCursorPosition() error = %v", err)
			}

			if err := backend.ScrollRegionDown(tt.startY, tt.endY, tt.count); err != nil {
				t.Fatalf("ScrollRegionDown(%d, %d, %d) error = %v", tt.startY, tt.endY, tt.count, err)
			}

			backend.AssertScrollbackEmpty(t)
			backend.AssertBufferLines(t, tt.wantLines)
			backend.AssertCursorPosition(t, layout.Position{X: 4, Y: 3})
			if got, want := backend.ScrollRegionDownCalls(), [][3]int{{tt.startY, tt.endY, tt.count}}; !slices.Equal(got, want) {
				t.Fatalf("ScrollRegionDownCalls() = %#v, want %#v", got, want)
			}
		})
	}
}
