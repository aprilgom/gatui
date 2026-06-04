package text_test

import (
	"slices"
	"strings"
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
		{name: "halfwidth voiced mark", span: text.NewSpan("ﾞ"), want: buffer.CellWidth("ﾞ")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.span.Width(); got != tt.want {
				t.Fatalf("Width() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSpan_add(t *testing.T) {
	baseStyle := style.NewStyle().Fg(style.Green).Bg(style.Yellow)
	otherStyle := style.NewStyle().Fg(style.Red)
	original := text.StyledSpan("hello", baseStyle)

	got := original.Add(text.StyledSpan(" world", otherStyle))

	if got.Content != "hello world" {
		t.Fatalf("content = %q, want hello world", got.Content)
	}
	if got.Style != baseStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, baseStyle)
	}
	if original.Content != "hello" {
		t.Fatalf("Add mutated original content to %q", original.Content)
	}
}

func TestSpan_fromRefStrBorrowedCow(t *testing.T) {
	source := "borrowed"

	got := text.NewSpan(source)
	source = "changed"

	if source != "changed" {
		t.Fatalf("source = %q, want changed", source)
	}
	if got.Content != "borrowed" {
		t.Fatalf("content = %q, want borrowed", got.Content)
	}
	if got.Style != style.NewStyle() {
		t.Fatalf("style = %#v, want default", got.Style)
	}
}

func TestSpan_fromRefStringBorrowedCow(t *testing.T) {
	source := strings.Builder{}
	source.WriteString("borrowed string")
	content := source.String()

	got := text.NewSpan(content)
	source.Reset()
	source.WriteString("changed")

	if got.Content != "borrowed string" {
		t.Fatalf("content = %q, want borrowed string", got.Content)
	}
	if got.Style != style.NewStyle() {
		t.Fatalf("style = %#v, want default", got.Style)
	}
}

func TestSpan_rawString(t *testing.T) {
	got := text.NewSpan(strings.Join([]string{"raw", "string"}, " "))

	if got.Content != "raw string" {
		t.Fatalf("content = %q, want raw string", got.Content)
	}
	if got.Style != style.NewStyle() {
		t.Fatalf("style = %#v, want default", got.Style)
	}
}

func TestSpan_setContent(t *testing.T) {
	originalStyle := style.NewStyle().Fg(style.Blue)
	original := text.StyledSpan("old", originalStyle)

	got := original.SetContent("new")

	if got.Content != "new" {
		t.Fatalf("content = %q, want new", got.Content)
	}
	if got.Style != originalStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, originalStyle)
	}
	if original.Content != "old" {
		t.Fatalf("SetContent mutated original content to %q", original.Content)
	}
}

func TestSpan_setStyle(t *testing.T) {
	original := text.StyledSpan("content", style.NewStyle().Fg(style.Blue))
	replacement := style.NewStyle().Fg(style.Red).Bg(style.Yellow).AddModifier(style.ModifierBold)

	got := original.SetStyle(replacement)

	if got.Content != "content" {
		t.Fatalf("content = %q, want content", got.Content)
	}
	if got.Style != replacement {
		t.Fatalf("style = %#v, want %#v", got.Style, replacement)
	}
	if original.Style == replacement {
		t.Fatalf("SetStyle mutated original style")
	}
}

func TestSpan_styledString(t *testing.T) {
	spanStyle := style.NewStyle().Fg(style.Magenta).AddModifier(style.ModifierItalic)
	got := text.StyledSpan(strings.Join([]string{"styled", "string"}, " "), spanStyle)

	if got.Content != "styled string" {
		t.Fatalf("content = %q, want styled string", got.Content)
	}
	if got.Style != spanStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, spanStyle)
	}
}

func TestSpan_toSpan(t *testing.T) {
	spanStyle := style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierUnderlined)
	original := text.StyledSpan("identity", spanStyle)

	got := original.ToSpan()

	if got != original {
		t.Fatalf("ToSpan() = %#v, want %#v", got, original)
	}
}

func TestStyledGrapheme_New_shouldStoreSymbolAndStyle(t *testing.T) {
	graphemeStyle := style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierItalic)

	got := text.NewStyledGrapheme("a", graphemeStyle)

	if got.Symbol != "a" {
		t.Fatalf("symbol = %q, want a", got.Symbol)
	}
	if got.Style != graphemeStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, graphemeStyle)
	}
}

func TestStyledGrapheme_IsWhitespace_shouldMatchRatatuiRules(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
		want   bool
	}{
		{name: "space", symbol: " ", want: true},
		{name: "tab", symbol: "\t", want: true},
		{name: "zero width space", symbol: "\u200B", want: true},
		{name: "non breaking space", symbol: "\u00A0", want: false},
		{name: "letter", symbol: "a", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := text.NewStyledGrapheme(tt.symbol, style.NewStyle()).IsWhitespace()
			if got != tt.want {
				t.Fatalf("IsWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStyledGrapheme_Stylize_shouldPatchStyle(t *testing.T) {
	grapheme := text.NewStyledGrapheme("a", style.NewStyle().Fg(style.Yellow).Bg(style.Red))

	got := grapheme.Cyan()
	want := style.NewStyle().Fg(style.Cyan).Bg(style.Red)

	if got.Style != want {
		t.Fatalf("style = %#v, want %#v", got.Style, want)
	}
}

func TestStyledGrapheme_setStyle(t *testing.T) {
	original := text.NewStyledGrapheme("a", style.NewStyle().Fg(style.Yellow))
	replacement := style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold)

	got := original.SetStyle(replacement)

	if got.Symbol != "a" {
		t.Fatalf("symbol = %q, want a", got.Symbol)
	}
	if got.Style != replacement {
		t.Fatalf("style = %#v, want %#v", got.Style, replacement)
	}
	if original.Style == replacement {
		t.Fatalf("SetStyle mutated original grapheme")
	}
}

func TestSpan_StyledGraphemes_shouldPatchBaseStyleWithSpanStyle(t *testing.T) {
	span := text.StyledSpan("Test", style.NewStyle().Fg(style.Green).AddModifier(style.ModifierItalic))
	baseStyle := style.NewStyle().Fg(style.Red).Bg(style.Yellow)

	got := span.StyledGraphemes(baseStyle)
	wantStyle := style.NewStyle().Fg(style.Green).Bg(style.Yellow).AddModifier(style.ModifierItalic)
	want := []text.StyledGrapheme{
		text.NewStyledGrapheme("T", wantStyle),
		text.NewStyledGrapheme("e", wantStyle),
		text.NewStyledGrapheme("s", wantStyle),
		text.NewStyledGrapheme("t", wantStyle),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("StyledGraphemes() = %#v, want %#v", got, want)
	}
}

func TestSpan_StyledGraphemes_shouldUseGraphemeClustersAndFilterControl(t *testing.T) {
	span := text.NewSpan("🇺🇸a\nb")

	got := span.StyledGraphemes(style.NewStyle())
	want := []text.StyledGrapheme{
		text.NewStyledGrapheme("🇺🇸", style.NewStyle()),
		text.NewStyledGrapheme("a", style.NewStyle()),
		text.NewStyledGrapheme("b", style.NewStyle()),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("StyledGraphemes() = %#v, want %#v", got, want)
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

func TestSpan_ResetStyle_shouldResetStyleAndPreserveContent(t *testing.T) {
	got := text.NewSpan("test content").
		Fg(style.Green).
		Bg(style.Yellow).
		Italic().
		ResetStyle()
	wantStyle := style.ResetStyle().AddModifier(style.ModifierItalic)

	if got.Content != "test content" {
		t.Fatalf("content = %q, want test content", got.Content)
	}
	if got.Style != wantStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, wantStyle)
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

func TestLine_StyledGraphemes_shouldPatchBaseLineAndSpanStyles(t *testing.T) {
	baseStyle := style.NewStyle().Bg(style.White)
	line := text.LineFromSpans(
		text.StyledSpan("He", style.NewStyle().Fg(style.Red)),
		text.StyledSpan("ll", style.NewStyle().Fg(style.Green)),
		text.StyledSpan("o!", style.NewStyle().Fg(style.Blue)),
	).Italic()

	got := line.StyledGraphemes(baseStyle)
	redStyle := style.NewStyle().Fg(style.Red).Bg(style.White).AddModifier(style.ModifierItalic)
	greenStyle := style.NewStyle().Fg(style.Green).Bg(style.White).AddModifier(style.ModifierItalic)
	blueStyle := style.NewStyle().Fg(style.Blue).Bg(style.White).AddModifier(style.ModifierItalic)
	want := []text.StyledGrapheme{
		text.NewStyledGrapheme("H", redStyle),
		text.NewStyledGrapheme("e", redStyle),
		text.NewStyledGrapheme("l", greenStyle),
		text.NewStyledGrapheme("l", greenStyle),
		text.NewStyledGrapheme("o", blueStyle),
		text.NewStyledGrapheme("!", blueStyle),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("StyledGraphemes() = %#v, want %#v", got, want)
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
		{
			name: "crab emoji",
			line: text.LineFromString("🦀"),
			want: 2,
		},
		{
			name: "flag emoji",
			line: text.LineFromString("🇺🇸"),
			want: 2,
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

func TestLine_crabEmojiWidth(t *testing.T) {
	got := text.LineFromString("🦀")

	if got.Width() != 2 {
		t.Fatalf("Width() = %d, want 2", got.Width())
	}
}

func TestLine_flagEmojiWidth(t *testing.T) {
	got := text.LineFromString("🇺🇸")

	if got.Width() != 2 {
		t.Fatalf("Width() = %d, want 2", got.Width())
	}
}

func TestLineFromSpans_shouldCreateLineFromStyledSpans(t *testing.T) {
	red := style.NewStyle().Fg(style.Red)
	green := style.NewStyle().Fg(style.Green)

	got := text.LineFromSpans(
		text.StyledSpan("Hello,", red),
		text.StyledSpan(" world!", green),
	)

	wantSpans := []text.Span{
		text.StyledSpan("Hello,", red),
		text.StyledSpan(" world!", green),
	}
	if !slices.Equal(got.Spans, wantSpans) {
		t.Fatalf("spans = %#v, want %#v", got.Spans, wantSpans)
	}
	if got.LineStyle != style.NewStyle() {
		t.Fatalf("line style = %#v, want default", got.LineStyle)
	}
	if got.Alignment != nil {
		t.Fatalf("alignment = %#v, want nil", got.Alignment)
	}
}

func TestLine_fromSpan(t *testing.T) {
	span := text.StyledSpan("hello", style.NewStyle().Fg(style.Red))

	got := text.LineFromSpan(span)

	if len(got.Spans) != 1 || got.Spans[0] != span {
		t.Fatalf("spans = %#v, want %#v", got.Spans, []text.Span{span})
	}
	if got.LineStyle != style.NewStyle() {
		t.Fatalf("line style = %#v, want default", got.LineStyle)
	}
}

func TestLine_toLine(t *testing.T) {
	span := text.StyledSpan("hello", style.NewStyle().Fg(style.Green))

	got := span.ToLine()

	if len(got.Spans) != 1 || got.Spans[0] != span {
		t.Fatalf("spans = %#v, want %#v", got.Spans, []text.Span{span})
	}
}

func TestLine_String_shouldConcatenateSpanContent(t *testing.T) {
	got := text.LineFromSpans(
		text.StyledSpan("Hello,", style.NewStyle().Fg(style.Red)),
		text.StyledSpan(" world!", style.NewStyle().Fg(style.Green)),
	).String()

	if got != "Hello, world!" {
		t.Fatalf("String() = %q, want %q", got, "Hello, world!")
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

func TestLine_ResetStyle_shouldResetLineStyleAndPreserveSpansAndAlignment(t *testing.T) {
	span := text.NewSpan("test content").Fg(style.Cyan)
	got := text.NewLine(span).
		Fg(style.Green).
		Bg(style.Yellow).
		Italic().
		Center().
		ResetStyle()
	wantStyle := style.ResetStyle().AddModifier(style.ModifierItalic)

	if len(got.Spans) != 1 || got.Spans[0] != span {
		t.Fatalf("spans = %#v, want %#v", got.Spans, []text.Span{span})
	}
	if got.Alignment == nil || *got.Alignment != layout.Center {
		t.Fatalf("alignment = %#v, want Center", got.Alignment)
	}
	if got.LineStyle != wantStyle {
		t.Fatalf("line style = %#v, want %#v", got.LineStyle, wantStyle)
	}
}

func TestLine_addLine(t *testing.T) {
	baseStyle := style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierBold)
	left := text.LineFromSpans(text.NewSpan("A")).
		Style(baseStyle).
		Center()
	right := text.LineFromSpans(
		text.StyledSpan("B", style.NewStyle().Fg(style.Red)),
		text.StyledSpan("C", style.NewStyle().Fg(style.Green)),
	).Right()

	got := left.AddLine(right)

	if got.String() != "ABC" {
		t.Fatalf("String() = %q, want ABC", got.String())
	}
	if got.LineStyle != baseStyle {
		t.Fatalf("line style = %#v, want %#v", got.LineStyle, baseStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Center {
		t.Fatalf("alignment = %#v, want Center", got.Alignment)
	}
	if len(got.Spans) != 3 || got.Spans[1] != right.Spans[0] || got.Spans[2] != right.Spans[1] {
		t.Fatalf("spans = %#v, want appended right spans", got.Spans)
	}
}

func TestLine_Extend_shouldAppendSpansAndPreserveMetadata(t *testing.T) {
	baseStyle := style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierBold)
	got := text.LineFromString("A").
		Style(baseStyle).
		Right().
		Extend(
			text.StyledSpan("B", style.NewStyle().Fg(style.Red)),
			text.StyledSpan("C", style.NewStyle().Fg(style.Green)),
		)

	if got.String() != "ABC" {
		t.Fatalf("String() = %q, want ABC", got.String())
	}
	if got.LineStyle != baseStyle {
		t.Fatalf("line style = %#v, want %#v", got.LineStyle, baseStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Right {
		t.Fatalf("alignment = %#v, want Right", got.Alignment)
	}
	if got.Spans[1] != text.StyledSpan("B", style.NewStyle().Fg(style.Red)) ||
		got.Spans[2] != text.StyledSpan("C", style.NewStyle().Fg(style.Green)) {
		t.Fatalf("appended span styles = %#v, want red/green", got.Spans)
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

func TestLine_forLoopInto(t *testing.T) {
	line := text.LineFromSpans(text.NewSpan("a"), text.NewSpan("b"))
	var got strings.Builder

	for _, span := range line.Spans {
		got.WriteString(span.Content)
	}

	if got.String() != "ab" {
		t.Fatalf("range content = %q, want ab", got.String())
	}
}

func TestLine_forLoopMutRef(t *testing.T) {
	line := text.LineFromSpans(text.NewSpan("a"), text.NewSpan("b"))

	for i := range line.Spans {
		line.Spans[i] = line.Spans[i].Fg(style.Cyan)
	}

	want := style.NewStyle().Fg(style.Cyan)
	for i, span := range line.Spans {
		if span.Style != want {
			t.Fatalf("span[%d].Style = %#v, want %#v", i, span.Style, want)
		}
	}
}

func TestLine_forLoopRef(t *testing.T) {
	line := text.LineFromSpans(text.NewSpan("a"), text.NewSpan("b"))
	got := make([]string, 0, len(line.Spans))

	for i := range line.Spans {
		span := &line.Spans[i]
		got = append(got, span.Content)
	}

	if !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("range refs = %#v, want [a b]", got)
	}
}

func TestLine_styledCow(t *testing.T) {
	lineStyle := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierItalic)

	got := text.StyledLine("hello", lineStyle)

	if got.String() != "hello" {
		t.Fatalf("String() = %q, want hello", got.String())
	}
	if got.LineStyle != lineStyle {
		t.Fatalf("line style = %#v, want %#v", got.LineStyle, lineStyle)
	}
}

func TestLine_styledString(t *testing.T) {
	lineStyle := style.NewStyle().Bg(style.Blue).AddModifier(style.ModifierBold)

	got := text.StyledLine("hello", lineStyle)

	if len(got.Spans) != 1 || got.Spans[0].Content != "hello" {
		t.Fatalf("spans = %#v, want one hello span", got.Spans)
	}
	if got.LineStyle != lineStyle {
		t.Fatalf("line style = %#v, want %#v", got.LineStyle, lineStyle)
	}
}

func TestTextFromSpan_shouldCreateSingleLineText(t *testing.T) {
	span := text.StyledSpan("hello", style.NewStyle().Fg(style.Red))

	got := text.TextFromSpan(span)

	if len(got.Lines) != 1 || len(got.Lines[0].Spans) != 1 {
		t.Fatalf("text shape = %#v, want one line with one span", got)
	}
	if got.Lines[0].Spans[0] != span {
		t.Fatalf("span = %#v, want %#v", got.Lines[0].Spans[0], span)
	}
	if got.String() != "hello" {
		t.Fatalf("String() = %q, want hello", got.String())
	}
}

func TestTextFromLine_shouldCreateSingleLineText(t *testing.T) {
	line := text.LineFromString("hello").
		Style(style.NewStyle().Fg(style.Yellow)).
		Center()

	got := text.TextFromLine(line)

	if len(got.Lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(got.Lines))
	}
	if !lineEqual(got.Lines[0], line) {
		t.Fatalf("line = %#v, want %#v", got.Lines[0], line)
	}
}

func TestText_fromVecLine(t *testing.T) {
	lines := []text.Line{
		text.LineFromString("The first line"),
		text.LineFromString("The second line").Right(),
	}

	got := text.TextFromLines(lines)

	if !linesEqual(got.Lines, lines) {
		t.Fatalf("lines = %#v, want %#v", got.Lines, lines)
	}
	lines[0] = text.LineFromString("changed")
	if got.Lines[0].String() != "The first line" {
		t.Fatalf("TextFromLines retained source slice alias; first line = %q", got.Lines[0].String())
	}
}

func TestText_fromCow(t *testing.T) {
	got := text.FromString("The first line\nThe second line")

	want := []text.Line{
		text.LineFromString("The first line"),
		text.LineFromString("The second line"),
	}
	if !linesEqual(got.Lines, want) {
		t.Fatalf("lines = %#v, want %#v", got.Lines, want)
	}
}

func TestText_toText(t *testing.T) {
	original := text.FromString("identity\ntext").
		Cyan().
		Center()

	got := original.ToText()

	if got.Style != original.Style {
		t.Fatalf("style = %#v, want %#v", got.Style, original.Style)
	}
	if got.Alignment == nil || original.Alignment == nil || *got.Alignment != *original.Alignment {
		t.Fatalf("alignment = %#v, want %#v", got.Alignment, original.Alignment)
	}
	if !linesEqual(got.Lines, original.Lines) {
		t.Fatalf("lines = %#v, want %#v", got.Lines, original.Lines)
	}
}

func TestText_String_shouldJoinLinesWithNewline(t *testing.T) {
	got := text.NewText(
		text.LineFromSpans(
			text.StyledSpan("Hello,", style.NewStyle().Fg(style.Red)),
			text.StyledSpan(" world!", style.NewStyle().Fg(style.Green)),
		),
		text.LineFromString("second").Right(),
		text.NewLine(),
	).String()

	if got != "Hello, world!\nsecond\n" {
		t.Fatalf("String() = %q, want %q", got, "Hello, world!\nsecond\n")
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

func TestText_Extend_shouldAppendLinesAndPreserveTextMetadata(t *testing.T) {
	textStyle := style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierItalic)
	got := text.TextFromLine(text.LineFromString("A")).
		PatchStyle(textStyle).
		Center().
		Extend(
			text.LineFromString("B").Right(),
			text.LineFromString("C"),
		)

	if got.String() != "A\nB\nC" {
		t.Fatalf("String() = %q, want A\\nB\\nC", got.String())
	}
	if got.Style != textStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, textStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Center {
		t.Fatalf("alignment = %#v, want Center", got.Alignment)
	}
	if got.Lines[1].Alignment == nil || *got.Lines[1].Alignment != layout.Right {
		t.Fatalf("appended line alignment = %#v, want Right", got.Lines[1].Alignment)
	}
}

func TestText_AppendText_shouldAppendLinesAndIgnoreRightMetadata(t *testing.T) {
	leftStyle := style.NewStyle().Fg(style.Red)
	rightStyle := style.NewStyle().Fg(style.Green)
	left := text.FromString("left").PatchStyle(leftStyle).Left()
	right := text.FromString("right").PatchStyle(rightStyle).Right()

	got := left.AppendText(right)

	if got.String() != "left\nright" {
		t.Fatalf("String() = %q, want left\\nright", got.String())
	}
	if got.Style != leftStyle {
		t.Fatalf("style = %#v, want %#v", got.Style, leftStyle)
	}
	if got.Alignment == nil || *got.Alignment != layout.Left {
		t.Fatalf("alignment = %#v, want Left", got.Alignment)
	}
	if len(got.Lines) != 2 || got.Lines[1].String() != "right" {
		t.Fatalf("lines = %#v, want appended right line", got.Lines)
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

func TestText_pushLineEmpty(t *testing.T) {
	got := text.NewText().PushLine(text.LineFromString(""))

	if len(got.Lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(got.Lines))
	}
	if got.Lines[0].String() != "" {
		t.Fatalf("line string = %q, want empty", got.Lines[0].String())
	}
	if got.Height() != 1 || got.Width() != 0 {
		t.Fatalf("dimensions = %dx%d, want 0x1", got.Width(), got.Height())
	}
}

func TestText_forLoopInto(t *testing.T) {
	helloWorld := text.TextFromLines([]text.Line{
		text.StyledLine("Hello ", style.NewStyle().Fg(style.Blue)),
		text.StyledLine("world!", style.NewStyle().Fg(style.Green)),
	})
	var got strings.Builder

	for _, line := range helloWorld.Lines {
		got.WriteString(line.String())
	}

	if got.String() != "Hello world!" {
		t.Fatalf("range content = %q, want Hello world!", got.String())
	}
}

func TestText_forLoopMutRef(t *testing.T) {
	helloWorld := text.TextFromLines([]text.Line{
		text.LineFromString("Hello "),
		text.LineFromString("world!"),
	})

	for i := range helloWorld.Lines {
		helloWorld.Lines[i] = helloWorld.Lines[i].Cyan()
	}

	want := style.NewStyle().Fg(style.Cyan)
	for i, line := range helloWorld.Lines {
		if line.LineStyle != want {
			t.Fatalf("line[%d].LineStyle = %#v, want %#v", i, line.LineStyle, want)
		}
	}
}

func TestText_forLoopRef(t *testing.T) {
	helloWorld := text.TextFromLines([]text.Line{
		text.LineFromString("Hello "),
		text.LineFromString("world!"),
	})
	got := make([]string, 0, len(helloWorld.Lines))

	for i := range helloWorld.Lines {
		line := &helloWorld.Lines[i]
		got = append(got, line.String())
	}

	if !slices.Equal(got, []string{"Hello ", "world!"}) {
		t.Fatalf("range refs = %#v, want [Hello  world!]", got)
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

func TestText_ResetStyle_shouldResetTextStyleAndPreserveLinesAndAlignment(t *testing.T) {
	lines := []text.Line{
		text.LineFromString("first"),
		text.LineFromString("second").Right(),
	}
	got := text.NewText(lines...).
		Fg(style.Green).
		Bg(style.Yellow).
		Italic().
		Center().
		ResetStyle()
	wantStyle := style.ResetStyle().AddModifier(style.ModifierItalic)

	if !linesEqual(got.Lines, lines) {
		t.Fatalf("lines = %#v, want %#v", got.Lines, lines)
	}
	if got.Alignment == nil || *got.Alignment != layout.Center {
		t.Fatalf("alignment = %#v, want Center", got.Alignment)
	}
	if got.Style != wantStyle {
		t.Fatalf("text style = %#v, want %#v", got.Style, wantStyle)
	}
}

func TestText_renderCenteredEven(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))

	text.FromString("foo").Center().Render(buf.Area, buf)

	assertTextLines(t, buf, []string{" foo  "})
}

func TestText_renderCenteredEvenWithTruncation(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))

	text.FromString("123456789").Center().Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"234567"})
}

func TestText_renderCenteredOddWithTruncation(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	text.FromString("123456789").Center().Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"34567"})
}

func TestText_renderOnlyStylesLineArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	textStyle := style.NewStyle().Bg(style.Blue)

	text.StyledText("foo", textStyle).Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"foo  "})
	for x := range 5 {
		assertTextCellStyle(t, buf, x, 0, textStyle)
	}
}

func TestText_renderOutOfBounds(t *testing.T) {
	tests := []struct {
		name string
		area layout.Rect
	}{
		{name: "fully outside", area: layout.NewRect(20, 20, 10, 1)},
		{name: "zero width", area: layout.NewRect(0, 0, 0, 1)},
		{name: "zero height", area: layout.NewRect(0, 0, 5, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

			text.FromString("Hello, world!").Render(tt.area, buf)

			assertTextLines(t, buf, []string{"     "})
		})
	}
}

func TestText_renderRightAlignedWithTruncation(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	text.FromString("123456789").Right().Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"56789"})
}

func TestText_renderTruncates(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))
	textStyle := style.NewStyle().Bg(style.Blue)

	text.StyledText("foobar", textStyle).Render(layout.NewRect(0, 0, 3, 1), buf)

	assertTextLines(t, buf, []string{"foo   "})
	for x := range 6 {
		want := style.NewStyle()
		if x < 3 {
			want = textStyle
		}
		assertTextCellStyle(t, buf, x, 0, want)
	}
}

func lineEqual(a, b text.Line) bool {
	if a.LineStyle != b.LineStyle {
		return false
	}
	if (a.Alignment == nil) != (b.Alignment == nil) {
		return false
	}
	if a.Alignment != nil && *a.Alignment != *b.Alignment {
		return false
	}
	return slices.Equal(a.Spans, b.Spans)
}

func linesEqual(a, b []text.Line) bool {
	return slices.EqualFunc(a, b, lineEqual)
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
	for x := range len("test content") {
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

func TestResetStyle_Render_shouldApplyResetColors(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))
	buf.SetStyle(buf.Area, style.NewStyle().Fg(style.Green).Bg(style.Yellow))

	text.NewSpan("hi").ResetStyle().Render(buf.Area, buf)

	assertTextCellStyle(t, buf, 0, 0, style.ResetStyle())
	assertTextCellStyle(t, buf, 1, 0, style.ResetStyle())
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

func TestSpan_renderOutOfBounds(t *testing.T) {
	tests := []struct {
		name string
		area layout.Rect
		want []string
	}{
		{
			name: "outside right",
			area: layout.NewRect(5, 0, 2, 1),
			want: []string{"    "},
		},
		{
			name: "outside below",
			area: layout.NewRect(0, 2, 4, 1),
			want: []string{"    "},
		},
		{
			name: "partially overlaps right edge",
			area: layout.NewRect(2, 0, 4, 1),
			want: []string{"  he"},
		},
		{
			name: "zero width",
			area: layout.NewRect(0, 0, 0, 1),
			want: []string{"    "},
		},
		{
			name: "zero height",
			area: layout.NewRect(0, 0, 4, 0),
			want: []string{"    "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))

			text.NewSpan("hello").Render(tt.area, buf)

			assertTextLines(t, buf, tt.want)
		})
	}
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

func TestLine_Render(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 1))
	lineStyle := style.NewStyle().AddModifier(style.ModifierItalic)
	blue := style.NewStyle().Fg(style.Blue)
	green := style.NewStyle().Fg(style.Green)
	line := text.LineFromSpans(
		text.StyledSpan("Hello ", blue),
		text.StyledSpan("world!", green),
	).Style(lineStyle)

	line.Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"Hello world!   "})
	for x := range 15 {
		want := lineStyle
		switch {
		case x < 6:
			want = lineStyle.Patch(blue)
		case x < 12:
			want = lineStyle.Patch(green)
		}
		assertTextCellStyle(t, buf, x, 0, want)
	}
}

func TestLine_RenderOnlyStylesFirstLine(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 2))
	lineStyle := style.NewStyle().AddModifier(style.ModifierItalic)
	blue := style.NewStyle().Fg(style.Blue)
	green := style.NewStyle().Fg(style.Green)
	line := text.LineFromSpans(
		text.StyledSpan("Hello ", blue),
		text.StyledSpan("world!", green),
	).Style(lineStyle)

	line.Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"Hello world!        ", "                    "})
	for x := range 20 {
		want := lineStyle
		switch {
		case x < 6:
			want = lineStyle.Patch(blue)
		case x < 12:
			want = lineStyle.Patch(green)
		}
		assertTextCellStyle(t, buf, x, 0, want)
		assertTextCellStyle(t, buf, x, 1, style.NewStyle())
	}
}

func TestLine_RenderOnlyStylesLineArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 1))
	lineStyle := style.NewStyle().AddModifier(style.ModifierItalic)
	blue := style.NewStyle().Fg(style.Blue)
	green := style.NewStyle().Fg(style.Green)
	line := text.LineFromSpans(
		text.StyledSpan("Hello ", blue),
		text.StyledSpan("world!", green),
	).Style(lineStyle)

	line.Render(layout.NewRect(0, 0, 15, 1), buf)

	assertTextLines(t, buf, []string{"Hello world!        "})
	for x := range 20 {
		want := style.NewStyle()
		if x < 15 {
			want = lineStyle
		}
		switch {
		case x < 6:
			want = lineStyle.Patch(blue)
		case x < 12:
			want = lineStyle.Patch(green)
		}
		assertTextCellStyle(t, buf, x, 0, want)
	}
}

func TestLine_RenderOutOfBounds(t *testing.T) {
	line := text.LineFromSpans(
		text.StyledSpan("Hello ", style.NewStyle().Fg(style.Blue)),
		text.StyledSpan("world!", style.NewStyle().Fg(style.Green)),
	).Style(style.NewStyle().AddModifier(style.ModifierItalic))

	tests := []struct {
		name string
		area layout.Rect
		want string
	}{
		{name: "fully outside", area: layout.NewRect(20, 20, 10, 1), want: "     "},
		{name: "zero width", area: layout.NewRect(0, 0, 0, 1), want: "     "},
		{name: "zero height", area: layout.NewRect(0, 0, 5, 0), want: "     "},
		{name: "partially overlapping", area: layout.NewRect(3, 0, 10, 1), want: "   He"},
		{name: "wide grapheme clipped", area: layout.NewRect(4, 0, 1, 1), want: "     "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
			renderLine := line
			if tt.name == "wide grapheme clipped" {
				renderLine = text.LineFromString("🦀")
			}

			renderLine.Render(tt.area, buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_Regression1032(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 83, 1))
	line := text.LineFromString(
		"🦀 RFC8628 OAuth 2.0 Device Authorization GrantでCLIからGithubのaccess tokenを取得する",
	)

	line.Render(buf.Area, buf)

	assertTextLines(t, buf, []string{
		"🦀 RFC8628 OAuth 2.0 Device Authorization GrantでCLIからGithubのaccess tokenを取得 ",
	})
}

func TestLine_Render_emptyLineAppliesStyleToFirstRowOnly(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 2))
	lineStyle := style.NewStyle().Fg(style.Red)

	text.NewLine().Style(lineStyle).Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"    ", "    "})
	for x := range 4 {
		assertTextCellStyle(t, buf, x, 0, lineStyle)
		assertTextCellStyle(t, buf, x, 1, style.NewStyle())
	}
}

func TestLine_Render_shouldTruncateEmojiLeftAlignment(t *testing.T) {
	tests := []struct {
		width int
		want  string
	}{
		{width: 4, want: "1234"},
		{width: 5, want: "1234 "},
		{width: 6, want: "1234🦀"},
		{width: 7, want: "1234🦀7"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, tt.width, 1))

			text.LineFromString("1234🦀7890").Left().Render(buf.Area, buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_Render_shouldTruncateEmojiRightAlignment(t *testing.T) {
	tests := []struct {
		width int
		want  string
	}{
		{width: 4, want: "7890"},
		{width: 5, want: " 7890"},
		{width: 6, want: "🦀7890"},
		{width: 7, want: "4🦀7890"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, tt.width, 1))

			text.LineFromString("1234🦀7890").Right().Render(buf.Area, buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_Render_shouldTruncateEmojiCenterAlignment(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
		want    string
	}{
		{name: "ab crab cd width 1", content: "ab🦀cd", width: 1, want: " "},
		{name: "ab crab cd width 2", content: "ab🦀cd", width: 2, want: "🦀"},
		{name: "ab crab cd width 3", content: "ab🦀cd", width: 3, want: "b🦀"},
		{name: "ab crab cd width 4", content: "ab🦀cd", width: 4, want: "b🦀c"},
		{name: "ab crab cdef width 2", content: "ab🦀cdef", width: 2, want: " c"},
		{name: "ab crab cdef width 3", content: "ab🦀cdef", width: 3, want: "🦀c"},
		{name: "ab crab cdef width 5", content: "ab🦀cdef", width: 5, want: "b🦀cd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, tt.width, 1))

			text.LineFromString(tt.content).Center().Render(buf.Area, buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_Render_shouldTruncateAwayFromOriginWithoutOverwritingOutsideArea(t *testing.T) {
	tests := []struct {
		name      string
		alignment layout.Alignment
		want      string
	}{
		{name: "left", alignment: layout.Left, want: "XXa🦀bcXXX"},
		{name: "center", alignment: layout.Center, want: "XX🦀bc🦀XX"},
		{name: "right", alignment: layout.Right, want: "XXXbc🦀dXX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.WithLines([]string{"XXXXXXXXXX"})
			line := text.NewLine(text.NewSpan("a🦀b"), text.NewSpan("c🦀d"))
			switch tt.alignment {
			case layout.Center:
				line = line.Center()
			case layout.Right:
				line = line.Right()
			default:
				line = line.Left()
			}

			line.Render(layout.NewRect(2, 0, 6, 1), buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_Render_shouldRightAlignMultiSpanWithWideRuneSkip(t *testing.T) {
	tests := []struct {
		width int
		want  string
	}{
		{width: 4, want: "c🦀d"},
		{width: 5, want: "bc🦀d"},
		{width: 6, want: " bc🦀d"},
		{width: 7, want: "🦀bc🦀d"},
		{width: 8, want: "a🦀bc🦀d"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, tt.width, 1))
			line := text.NewLine(text.NewSpan("a🦀b"), text.NewSpan("c🦀d")).Right()

			line.Render(buf.Area, buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_Render_shouldTruncateFlagEmoji(t *testing.T) {
	tests := []struct {
		width int
		want  string
	}{
		{width: 1, want: " "},
		{width: 2, want: "🇺🇸"},
		{width: 3, want: "🇺🇸1"},
		{width: 4, want: "🇺🇸12"},
		{width: 5, want: "🇺🇸123"},
		{width: 6, want: "🇺🇸1234"},
		{width: 7, want: "🇺🇸1234 "},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, tt.width, 1))

			text.LineFromString("🇺🇸1234").Left().Render(buf.Area, buf)

			assertTextLines(t, buf, []string{tt.want})
		})
	}
}

func TestLine_Render_shouldTruncateVeryLongLineOfManySpans(t *testing.T) {
	line := veryLongLineOfManySpans().Left()
	buf := buffer.Empty(layout.NewRect(0, 0, 32, 1))

	line.Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"This is some content with a some"})

	buf = buffer.Empty(layout.NewRect(0, 0, 32, 1))
	line.Right().Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"horribly long Line over u16::MAX"})
}

func TestLine_Render_shouldTruncateVeryLongSingleSpanLine(t *testing.T) {
	line := text.LineFromString(veryLongLineContent()).Left()
	buf := buffer.Empty(layout.NewRect(0, 0, 32, 1))

	line.Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"This is some content with a some"})

	buf = buffer.Empty(layout.NewRect(0, 0, 32, 1))
	line.Right().Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"horribly long Line over u16::MAX"})
}

func TestLine_Render_shouldIgnoreNewlines(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 1))

	text.LineFromString("Hello\nworld!").Render(buf.Area, buf)

	assertTextLines(t, buf, []string{"Helloworld!"})
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
	if actual := buf.Lines(); !slices.Equal(actual, expected) {
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

func veryLongLineOfManySpans() text.Line {
	part := veryLongLinePart()
	line := text.NewLine()
	for line.Width() < 65536 {
		line = line.PushSpan(text.NewSpan(part))
	}
	return line.PushSpan(text.NewSpan("horribly long Line over u16::MAX"))
}

func veryLongLineContent() string {
	part := veryLongLinePart()
	var builder strings.Builder
	for builder.Len() < 65536 {
		builder.WriteString(part)
	}
	builder.WriteString("horribly long Line over u16::MAX")
	return builder.String()
}

func veryLongLinePart() string {
	return "This is some content with a somewhat "
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
