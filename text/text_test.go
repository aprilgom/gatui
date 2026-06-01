package text_test

import (
	"testing"

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

func TestText_Width_shouldReuseLineWidth(t *testing.T) {
	got := text.NewText(
		text.NewLine(text.NewSpan("a"), text.NewSpan("コンピ")),
		text.LineFromString("short"),
	)

	if got.Width() != 7 {
		t.Fatalf("Width() = %d, want 7", got.Width())
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
