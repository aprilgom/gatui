package text_test

import (
	"reflect"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

func TestFromString_shouldSplitLinesAndPreserveContent(t *testing.T) {
	got := text.FromString("first\nsecond")

	if len(got.Lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(got.Lines))
	}
	if got.Lines[0].Spans[0].Content != "first" || got.Lines[1].Spans[0].Content != "second" {
		t.Fatalf("unexpected text: %#v", got)
	}
}

func TestText_StyledText_shouldSetTextStyle(t *testing.T) {
	textStyle := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierItalic)

	got := text.StyledText("a\nb", textStyle)

	if got.Style != textStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, textStyle)
	}
	if len(got.Lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(got.Lines))
	}
	if got.Lines[0].Spans[0].Content != "a" || got.Lines[1].Spans[0].Content != "b" {
		t.Fatalf("unexpected text: %#v", got)
	}
}

func TestText_WidthAndHeight_shouldUseDisplayWidth(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "ascii",
			content:    "The first line\nThe second line",
			wantWidth:  15,
			wantHeight: 2,
		},
		{
			name:       "unicode",
			content:    "コンピ\nabc",
			wantWidth:  6,
			wantHeight: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := text.FromString(tt.content)

			if got.Width() != tt.wantWidth {
				t.Fatalf("Width() = %d, want %d", got.Width(), tt.wantWidth)
			}
			if got.Height() != tt.wantHeight {
				t.Fatalf("Height() = %d, want %d", got.Height(), tt.wantHeight)
			}
		})
	}
}

func TestSpan_Width_shouldUseDisplayWidth(t *testing.T) {
	tests := []struct {
		name string
		span text.Span
		want int
	}{
		{name: "ascii", span: text.NewSpan("My text"), want: 7},
		{name: "unicode", span: text.NewSpan("コンピ"), want: 6},
		{name: "empty", span: text.NewSpan(""), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.span.Width(); got != tt.want {
				t.Fatalf("Width() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSpan_PatchStyle_shouldPatchExistingStyle(t *testing.T) {
	base := text.StyledSpan("hi", style.NewStyle().
		Fg(style.Yellow).
		AddModifier(style.ModifierItalic))
	patch := style.NewStyle().
		Fg(style.Red).
		AddModifier(style.ModifierUnderlined)

	got := base.PatchStyle(patch)
	want := style.NewStyle().
		Fg(style.Red).
		AddModifier(style.ModifierItalic | style.ModifierUnderlined)

	if got.Style != want {
		t.Fatalf("style = %#v, want %#v", got.Style, want)
	}
}

func TestSpan_LeftLine_shouldPreserveSpanAndSetAlignment(t *testing.T) {
	greenItalic := style.NewStyle().
		Fg(style.Green).
		AddModifier(style.ModifierItalic)

	got := text.StyledSpan("Test Content", greenItalic).LeftLine()

	if len(got.Spans) != 1 {
		t.Fatalf("span count = %d, want 1", len(got.Spans))
	}
	if got.Spans[0].Content != "Test Content" {
		t.Fatalf("content = %q, want Test Content", got.Spans[0].Content)
	}
	if got.Spans[0].Style != greenItalic {
		t.Fatalf("span style = %#v, want %#v", got.Spans[0].Style, greenItalic)
	}
	if got.LineStyle != style.NewStyle() {
		t.Fatalf("line style = %#v, want default", got.LineStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Left {
		t.Fatalf("alignment = %#v, want Left", got.Alignment)
	}
}

func TestSpan_CenterLine_shouldPreserveSpanAndSetAlignment(t *testing.T) {
	greenItalic := style.NewStyle().
		Fg(style.Green).
		AddModifier(style.ModifierItalic)

	got := text.StyledSpan("Test Content", greenItalic).CenterLine()

	if len(got.Spans) != 1 {
		t.Fatalf("span count = %d, want 1", len(got.Spans))
	}
	if got.Spans[0].Content != "Test Content" {
		t.Fatalf("content = %q, want Test Content", got.Spans[0].Content)
	}
	if got.Spans[0].Style != greenItalic {
		t.Fatalf("span style = %#v, want %#v", got.Spans[0].Style, greenItalic)
	}
	if got.LineStyle != style.NewStyle() {
		t.Fatalf("line style = %#v, want default", got.LineStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Center {
		t.Fatalf("alignment = %#v, want Center", got.Alignment)
	}
}

func TestSpan_RightLine_shouldPreserveSpanAndSetAlignment(t *testing.T) {
	greenItalic := style.NewStyle().
		Fg(style.Green).
		AddModifier(style.ModifierItalic)

	got := text.StyledSpan("Test Content", greenItalic).RightLine()

	if len(got.Spans) != 1 {
		t.Fatalf("span count = %d, want 1", len(got.Spans))
	}
	if got.Spans[0].Content != "Test Content" {
		t.Fatalf("content = %q, want Test Content", got.Spans[0].Content)
	}
	if got.Spans[0].Style != greenItalic {
		t.Fatalf("span style = %#v, want %#v", got.Spans[0].Style, greenItalic)
	}
	if got.LineStyle != style.NewStyle() {
		t.Fatalf("line style = %#v, want default", got.LineStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Right {
		t.Fatalf("alignment = %#v, want Right", got.Alignment)
	}
}

func TestLine_Width_shouldSumSpanDisplayWidths(t *testing.T) {
	tests := []struct {
		name string
		line text.Line
		want int
	}{
		{
			name: "ascii spans",
			line: text.NewLine(
				text.StyledSpan("My", style.NewStyle().Fg(style.Yellow)),
				text.NewSpan(" text"),
			),
			want: 7,
		},
		{
			name: "mixed unicode and ascii spans",
			line: text.NewLine(
				text.NewSpan("コンピ"),
				text.NewSpan(" abc"),
			),
			want: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.Width(); got != tt.want {
				t.Fatalf("Width() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestLine_PatchStyle_shouldPatchExistingLineStyle(t *testing.T) {
	base := text.LineFromString("hi").Style(style.NewStyle().Fg(style.Yellow))
	patch := style.NewStyle().AddModifier(style.ModifierItalic)

	got := base.PatchStyle(patch)
	want := style.NewStyle().
		Fg(style.Yellow).
		AddModifier(style.ModifierItalic)

	if got.LineStyle != want {
		t.Fatalf("style = %#v, want %#v", got.LineStyle, want)
	}
}

func TestLine_PushSpan_shouldAppendSpanAndPreserveLineMetadata(t *testing.T) {
	got := text.LineFromString("Hello, ").Cyan().Center().PushSpan(text.NewSpan("world!"))
	wantStyle := style.NewStyle().Fg(style.Cyan)

	if len(got.Spans) != 2 {
		t.Fatalf("span count = %d, want 2", len(got.Spans))
	}
	if got.Spans[0].Content != "Hello, " || got.Spans[1].Content != "world!" {
		t.Fatalf("spans = %#v, want Hello, /world!", got.Spans)
	}
	if got.LineStyle != wantStyle {
		t.Fatalf("style = %#v, want %#v", got.LineStyle, wantStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Center {
		t.Fatalf("alignment = %#v, want Center", got.Alignment)
	}
}

func TestLine_AppendSpans_shouldAppendMultipleSpans(t *testing.T) {
	got := text.LineFromString("Hello, ").AppendSpans(
		text.NewSpan("world! "),
		text.NewSpan("How are you?"),
	)

	if len(got.Spans) != 3 {
		t.Fatalf("span count = %d, want 3", len(got.Spans))
	}
	wantContent := []string{"Hello, ", "world! ", "How are you?"}
	for i, want := range wantContent {
		if got.Spans[i].Content != want {
			t.Fatalf("span[%d] = %q, want %q", i, got.Spans[i].Content, want)
		}
	}
	if got.Width() != 26 {
		t.Fatalf("Width() = %d, want 26", got.Width())
	}
}

func TestText_Width_shouldReuseLineWidth(t *testing.T) {
	got := text.NewText(
		text.NewLine(text.NewSpan("a"), text.NewSpan("コンピ")),
		text.LineFromString("short"),
	)

	if got.Width() != 7 {
		t.Fatalf("Width() = %d, want 7", got.Width())
	}
}

func TestText_PushLine_shouldAppendLineAndPreserveTextStyle(t *testing.T) {
	got := text.FromString("A").Cyan().PushLine(text.LineFromString("B"))
	wantStyle := style.NewStyle().Fg(style.Cyan)

	if len(got.Lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(got.Lines))
	}
	if got.Lines[0].Spans[0].Content != "A" || got.Lines[1].Spans[0].Content != "B" {
		t.Fatalf("lines = %#v, want A/B", got.Lines)
	}
	if got.Style != wantStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, wantStyle)
	}
}

func TestText_shouldSupportAlignmentHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  text.Text
		want layout.Alignment
	}{
		{name: "left", got: text.FromString("hello").Left(), want: layout.Left},
		{name: "center", got: text.FromString("hello").Center(), want: layout.Center},
		{name: "right", got: text.FromString("hello").Right(), want: layout.Right},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.Alignment == nil || *tt.got.Alignment != tt.want {
				t.Fatalf("alignment = %#v, want %v", tt.got.Alignment, tt.want)
			}
		})
	}
}

func TestTextAlignment_shouldPreserveStyleAndMutationHelpers(t *testing.T) {
	got := text.FromString("A").Cyan().Center().PushLine(text.LineFromString("B"))
	wantStyle := style.NewStyle().Fg(style.Cyan)

	if got.Style != wantStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, wantStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Center {
		t.Fatalf("alignment = %#v, want Center", got.Alignment)
	}
	if len(got.Lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(got.Lines))
	}
}

func TestText_PushSpan_shouldAppendToLastLine(t *testing.T) {
	got := text.FromString("A").
		PushSpan(text.NewSpan("B")).
		PushSpan(text.NewSpan("C"))

	if len(got.Lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(got.Lines))
	}
	wantContent := []string{"A", "B", "C"}
	for i, want := range wantContent {
		if got.Lines[0].Spans[i].Content != want {
			t.Fatalf("span[%d] = %q, want %q", i, got.Lines[0].Spans[i].Content, want)
		}
	}
	if got.Width() != 3 || got.Height() != 1 {
		t.Fatalf("dimensions = %dx%d, want 3x1", got.Width(), got.Height())
	}
}

func TestText_PushSpan_shouldCreateLineWhenTextIsEmpty(t *testing.T) {
	got := text.NewText().PushSpan(text.NewSpan("Hello"))

	if len(got.Lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(got.Lines))
	}
	if len(got.Lines[0].Spans) != 1 {
		t.Fatalf("span count = %d, want 1", len(got.Lines[0].Spans))
	}
	if got.Lines[0].Spans[0].Content != "Hello" {
		t.Fatalf("content = %q, want Hello", got.Lines[0].Spans[0].Content)
	}
	if got.Width() != 5 || got.Height() != 1 {
		t.Fatalf("dimensions = %dx%d, want 5x1", got.Width(), got.Height())
	}
}

func TestText_PatchStyle_shouldPatchExistingTextStyle(t *testing.T) {
	base := text.StyledText("hi", style.NewStyle().
		Fg(style.Yellow).
		AddModifier(style.ModifierItalic))
	patch := style.NewStyle().
		Fg(style.Red).
		AddModifier(style.ModifierUnderlined)

	got := base.PatchStyle(patch)
	want := style.NewStyle().
		Fg(style.Red).
		AddModifier(style.ModifierItalic | style.ModifierUnderlined)

	if got.Style != want {
		t.Fatalf("style = %#v, want %#v", got.Style, want)
	}
}

func TestLine_shouldSupportStylizeAndAlignmentHelpers(t *testing.T) {
	got := text.LineFromString("hello").Cyan().Bold().Right()
	wantStyle := style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierBold)

	if got.LineStyle != wantStyle {
		t.Fatalf("style = %#v, want %#v", got.LineStyle, wantStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Right {
		t.Fatalf("alignment = %#v, want Right", got.Alignment)
	}
}

func TestLineStylize_shouldUpdateLineStyle(t *testing.T) {
	got := text.LineFromString("hi").Cyan().Bold()
	wantStyle := style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierBold)

	if got.LineStyle != wantStyle {
		t.Fatalf("style = %#v, want %#v", got.LineStyle, wantStyle)
	}
	if got.Spans[0].Style != style.NewStyle() {
		t.Fatalf("span style = %#v, want default", got.Spans[0].Style)
	}
}

func TestTextStylize_shouldUpdateTextStyle(t *testing.T) {
	got := text.FromString("hi").Cyan().Bold()
	wantStyle := style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierBold)

	if got.Style != wantStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, wantStyle)
	}
	if got.Lines[0].LineStyle != style.NewStyle() {
		t.Fatalf("line style = %#v, want default", got.Lines[0].LineStyle)
	}
	if got.Lines[0].Spans[0].Style != style.NewStyle() {
		t.Fatalf("span style = %#v, want default", got.Lines[0].Spans[0].Style)
	}
}

func TestSpan_Render_shouldDrawStyledContent(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 1))
	spanStyle := style.NewStyle().Fg(style.Green).Bg(style.Yellow)

	text.StyledSpan("test content", spanStyle).Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"test content   "})
	for x := 0; x < len("test content"); x++ {
		assertTextCellStyle(t, buf, x, 0, spanStyle)
	}
}

func TestSpan_Render_shouldPatchExistingStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))
	buf.SetStyle(buf.Area, style.NewStyle().AddModifier(style.ModifierItalic))
	spanStyle := style.NewStyle().Fg(style.Green)

	text.StyledSpan("hi", spanStyle).Render(buf.Area, buf)

	want := style.NewStyle().Fg(style.Green).AddModifier(style.ModifierItalic)
	assertTextCellStyle(t, buf, 0, 0, want)
	assertTextCellStyle(t, buf, 1, 0, want)
}

func TestSpan_Render_shouldTruncateToAreaWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	text.NewSpan("test content").Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"test "})
}

func TestSpan_Render_shouldRenderWideSymbolAndClearHiddenCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 1))

	text.NewSpan("test 😃 content").Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"test 😃 content"})
	cell, ok := buf.CellAt(6, 0)
	if !ok {
		t.Fatal("cell at (6,0) missing")
	}
	if cell.Symbol != " " || cell.Style != style.NewStyle() {
		t.Fatalf("wide hidden cell = %#v, want blank default cell", cell)
	}
}

func TestSpan_Render_shouldAppendLeadingZeroWidthToFirstCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))

	text.NewSpan("\u200Babc").Render(buf.Area, buf)

	assertTextCellSymbols(t, buf, []string{"\u200Ba", "b", "c"})
}

func TestSpan_Render_shouldAppendMiddleZeroWidthToPreviousCell(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{name: "second", content: "a\u200Bbc", want: []string{"a\u200B", "b", "c"}},
		{name: "middle", content: "ab\u200Bc", want: []string{"a", "b\u200B", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))

			text.NewSpan(tt.content).Render(buf.Area, buf)

			assertTextCellSymbols(t, buf, tt.want)
		})
	}
}

func TestSpan_Render_shouldAppendTrailingZeroWidthToLastCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))

	text.NewSpan("abc\u200B").Render(buf.Area, buf)

	assertTextCellSymbols(t, buf, []string{"a", "b", "c\u200B"})
}

func TestSpan_Render_shouldHandleLeftToRightMarkAtBufferEnd(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	text.NewSpan("Hello\u200E").Render(buf.Area, buf)

	assertTextCellSymbols(t, buf, []string{"H", "e", "l", "l", "o\u200E"})
}

func TestSpan_Render_shouldIgnoreNewlineDuringRender(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))

	text.NewSpan("a\nb").Render(buf.Area, buf)

	assertTextCellSymbols(t, buf, []string{"a", "b"})
}

func TestSpan_Render_shouldTruncateWideSymbolAsWholeSymbol(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))

	text.NewSpan("test 😃 content").Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"test  "})
}

func TestSpan_Render_shouldTruncateOverflowingAreaToBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 1))

	text.NewSpan("test content").Render(layout.NewRect(10, 0, 20, 1), buf)

	assertTextLines(t, buf, []string{"          test "})
}

func TestLine_Render_shouldRespectAlignmentAndTruncation(t *testing.T) {
	tests := []struct {
		name string
		line text.Line
		want string
	}{
		{name: "center", line: text.LineFromString("foo").Center(), want: " foo "},
		{name: "right", line: text.LineFromString("foo").Right(), want: "  foo"},
		{name: "right truncation", line: text.LineFromString("123456789").Right(), want: "56789"},
		{name: "center truncation", line: text.LineFromString("123456789").Center(), want: "34567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

			tt.line.Render(buf.Area, buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_RenderWithAlignment_shouldUseFallbackWhenLineAlignmentAbsent(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	fallback := layout.Right

	text.LineFromString("foo").RenderWithAlignment(buf.Area, buf, &fallback)

	assertTextLines(t, buf, []string{"  foo"})
}

func TestText_Render_shouldRenderRowsWithTextStyleAndAlignment(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 2))
	textStyle := style.NewStyle().Bg(style.Green)

	text.StyledText("foo\nbar", textStyle).Center().Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"  foo  ", "  bar  "})
	for y := 0; y < buf.Area.Height; y++ {
		for x := 0; x < buf.Area.Width; x++ {
			assertTextCellStyle(t, buf, x, y, textStyle)
		}
	}
}

func TestText_Render_shouldPreferLineAlignmentOverTextAlignment(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))
	content := text.NewText(
		text.LineFromString("foo"),
		text.LineFromString("bar").Center(),
	).Right()

	content.Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"  foo", " bar "})
}

func assertTextLines(t *testing.T, buf *buffer.Buffer, expected []string) {
	t.Helper()
	if actual := buf.Lines(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("lines = %#v, want %#v", actual, expected)
	}
}

func assertTextCellStyle(t *testing.T, buf *buffer.Buffer, x, y int, expected style.Style) {
	t.Helper()
	cell, ok := buf.CellAt(x, y)
	if !ok {
		t.Fatalf("cell at (%d,%d) missing", x, y)
	}
	if cell.Style != expected {
		t.Fatalf("style at (%d,%d) = %#v, want %#v", x, y, cell.Style, expected)
	}
}

func assertTextCellSymbols(t *testing.T, buf *buffer.Buffer, expected []string) {
	t.Helper()
	if len(buf.Cells) != len(expected) {
		t.Fatalf("cell count = %d, want %d", len(buf.Cells), len(expected))
	}
	for i, want := range expected {
		if got := buf.Cells[i].Symbol; got != want {
			t.Fatalf("cell %d symbol = %#v, want %#v", i, got, want)
		}
	}
}
