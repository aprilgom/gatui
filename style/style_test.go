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

func TestStyle_ResetStyle_shouldSetForegroundAndBackgroundReset(t *testing.T) {
	got := style.ResetStyle()
	want := style.Style{
		Foreground: style.Reset,
		Background: style.Reset,
		Modifiers:  0,
	}

	if got != want {
		t.Fatalf("ResetStyle() = %#v, want %#v", got, want)
	}
}

func TestStyle_addModifierCanBeStylized(t *testing.T) {
	got := style.NewStyle().AddModifier(style.ModifierBold | style.ModifierItalic)

	want := style.NewStyle().AddModifier(style.ModifierBold | style.ModifierItalic)
	if got != want {
		t.Fatalf("AddModifier() = %#v, want %#v", got, want)
	}
}

func TestStyle_removeModifierCanBeStylized(t *testing.T) {
	got := style.NewStyle().
		AddModifier(style.ModifierBold | style.ModifierItalic | style.ModifierUnderlined).
		RemoveModifier(style.ModifierItalic | style.ModifierUnderlined)

	want := style.NewStyle().AddModifier(style.ModifierBold)
	if got != want {
		t.Fatalf("RemoveModifier() = %#v, want %#v", got, want)
	}
}

func TestStyle_hasModifierChecks(t *testing.T) {
	got := style.NewStyle().AddModifier(style.ModifierBold | style.ModifierItalic)

	if !got.HasModifier(style.ModifierBold) {
		t.Fatal("HasModifier(Bold) = false, want true")
	}
	if !got.HasModifier(style.ModifierBold | style.ModifierItalic) {
		t.Fatal("HasModifier(Bold|Italic) = false, want true")
	}
	if got.HasModifier(style.ModifierUnderlined) {
		t.Fatal("HasModifier(Underlined) = true, want false")
	}
}

func TestStyle_combineIndividualModifiers(t *testing.T) {
	got := style.NewStyle().
		AddModifier(style.ModifierBold).
		AddModifier(style.ModifierItalic).
		AddModifier(style.ModifierUnderlined)

	want := style.NewStyle().AddModifier(style.ModifierBold | style.ModifierItalic | style.ModifierUnderlined)
	if got != want {
		t.Fatalf("combined modifiers = %#v, want %#v", got, want)
	}
}

func TestStyle_fromColor(t *testing.T) {
	got := style.StyleFromColor(style.Red)
	want := style.NewStyle().Fg(style.Red)

	if got != want {
		t.Fatalf("StyleFromColor() = %#v, want %#v", got, want)
	}
}

func TestStyle_fromColorColor(t *testing.T) {
	got := style.StyleFromColors(style.Red, style.Blue)
	want := style.NewStyle().Fg(style.Red).Bg(style.Blue)

	if got != want {
		t.Fatalf("StyleFromColors() = %#v, want %#v", got, want)
	}
}

func TestStyle_fromColorModifier(t *testing.T) {
	got := style.StyleFromColor(style.Red).AddModifier(style.ModifierBold)
	want := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold)

	if got != want {
		t.Fatalf("StyleFromColor().AddModifier() = %#v, want %#v", got, want)
	}
}

func TestStyle_fromModifier(t *testing.T) {
	got := style.StyleFromModifier(style.ModifierItalic)
	want := style.NewStyle().AddModifier(style.ModifierItalic)

	if got != want {
		t.Fatalf("StyleFromModifier() = %#v, want %#v", got, want)
	}
}

func TestStyle_fromModifierModifier(t *testing.T) {
	got := style.StyleFromModifier(style.ModifierBold).AddModifier(style.ModifierItalic)
	want := style.NewStyle().AddModifier(style.ModifierBold | style.ModifierItalic)

	if got != want {
		t.Fatalf("StyleFromModifier().AddModifier() = %#v, want %#v", got, want)
	}
}

func TestStyle_fromColorColorModifier(t *testing.T) {
	got := style.StyleFromColors(style.Red, style.Blue).AddModifier(style.ModifierBold)
	want := style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold)

	if got != want {
		t.Fatalf("StyleFromColors().AddModifier() = %#v, want %#v", got, want)
	}
}

func TestStyle_fromColorColorModifierModifier(t *testing.T) {
	got := style.StyleFromColors(style.Red, style.Blue).
		AddModifier(style.ModifierBold).
		AddModifier(style.ModifierItalic)
	want := style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold | style.ModifierItalic)

	if got != want {
		t.Fatalf("StyleFromColors().AddModifier().AddModifier() = %#v, want %#v", got, want)
	}
}

func TestColor_fromANSIColor(t *testing.T) {
	got, ok := style.ColorFromANSIColor(style.LightRed)
	if !ok {
		t.Fatal("ColorFromANSIColor(LightRed) ok = false, want true")
	}
	if got != style.LightRed {
		t.Fatalf("ColorFromANSIColor(LightRed) = %#v, want %#v", got, style.LightRed)
	}
}

func TestColor_fromIndexedColor(t *testing.T) {
	got := style.IndexedColor(42)
	want, ok := style.ColorFromIndexedColor(42)
	if !ok {
		t.Fatal("ColorFromIndexedColor(42) ok = false, want true")
	}
	if got != want {
		t.Fatalf("IndexedColor(42) = %#v, want %#v", got, want)
	}
}

func TestColor_fromRGBColor(t *testing.T) {
	got := style.RGBColor(1, 2, 3)
	want, ok := style.ColorFromRGBColor(1, 2, 3)
	if !ok {
		t.Fatal("ColorFromRGBColor(1,2,3) ok = false, want true")
	}
	if got != want {
		t.Fatalf("RGBColor(1,2,3) = %#v, want %#v", got, want)
	}
}

func TestColor_fromUint32(t *testing.T) {
	got, ok := style.ColorFromUint32(0x010203)
	if !ok {
		t.Fatal("ColorFromUint32(0x010203) ok = false, want true")
	}
	want := style.RGBColor(1, 2, 3)
	if got != want {
		t.Fatalf("ColorFromUint32(0x010203) = %#v, want %#v", got, want)
	}
}

func TestColor_fromInvalidColors(t *testing.T) {
	if got, ok := style.ColorFromANSIColor(style.IndexedColor(1)); ok {
		t.Fatalf("ColorFromANSIColor(IndexedColor(1)) = %#v, true; want false", got)
	}
	if got, ok := style.ColorFromUint32(0x01000000); ok {
		t.Fatalf("ColorFromUint32(0x01000000) = %#v, true; want false", got)
	}
}

func TestStylize_fg(t *testing.T) {
	got := style.NewStyle().Fg(style.Red)
	want := style.NewStyle().Fg(style.Red)

	if got != want {
		t.Fatalf("Fg() = %#v, want %#v", got, want)
	}
}

func TestStyle_fgCanBeStylized(t *testing.T) {
	got := style.NewStyle().Fg(style.Red)
	want := style.StyleFromColor(style.Red)

	if got != want {
		t.Fatalf("Fg() = %#v, want %#v", got, want)
	}
}

func TestStylize_bg(t *testing.T) {
	got := style.NewStyle().Bg(style.Blue)
	want := style.NewStyle().Bg(style.Blue)

	if got != want {
		t.Fatalf("Bg() = %#v, want %#v", got, want)
	}
}

func TestStyle_bgCanBeStylized(t *testing.T) {
	got := style.NewStyle().Bg(style.Blue)
	want := style.NewStyle().Bg(style.Blue)

	if got != want {
		t.Fatalf("Bg() = %#v, want %#v", got, want)
	}
}

func TestStylize_colorModifier(t *testing.T) {
	got := style.StyleFromColor(style.Red).AddModifier(style.ModifierBold)
	want := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold)

	if got != want {
		t.Fatalf("color plus modifier style = %#v, want %#v", got, want)
	}
}

func TestStylize_fgBg(t *testing.T) {
	got := style.NewStyle().Fg(style.Red).Bg(style.Blue)
	want := style.NewStyle().Fg(style.Red).Bg(style.Blue)

	if got != want {
		t.Fatalf("Fg().Bg() = %#v, want %#v", got, want)
	}
}

func TestStylize_allChained(t *testing.T) {
	got := style.NewStyle().
		Fg(style.Red).
		Bg(style.Blue).
		AddModifier(style.ModifierBold).
		AddModifier(style.ModifierItalic)
	want := style.NewStyle().
		Fg(style.Red).
		Bg(style.Blue).
		AddModifier(style.ModifierBold | style.ModifierItalic)

	if got != want {
		t.Fatalf("chained style = %#v, want %#v", got, want)
	}
}

func TestStylize_repeatedAttributes(t *testing.T) {
	got := style.NewStyle().
		Fg(style.Red).
		Fg(style.Green).
		Bg(style.Blue).
		Bg(style.Yellow).
		AddModifier(style.ModifierBold).
		AddModifier(style.ModifierBold)
	want := style.NewStyle().
		Fg(style.Green).
		Bg(style.Yellow).
		AddModifier(style.ModifierBold)

	if got != want {
		t.Fatalf("repeated attributes = %#v, want %#v", got, want)
	}
}
