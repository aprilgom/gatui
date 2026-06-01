package style_test

import (
	"testing"

	"gatui/style"
)

func TestStylePatch_shouldMergeOnlySpecifiedColorsAndModifiers(t *testing.T) {
	base := style.NewStyle().Fg(style.Red).Bg(style.Black).AddModifier(style.ModifierBold)
	other := style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierItalic)

	got := base.Patch(other)

	want := style.Style{
		Foreground: style.Cyan,
		Background: style.Black,
		Modifiers:  style.ModifierBold | style.ModifierItalic,
	}
	if got != want {
		t.Fatalf("Patch() = %#v, want %#v", got, want)
	}
}
