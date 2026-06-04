package style_test

import (
	"testing"

	"github.com/aprilgom/gatui/style"
)

func TestColor_String_shouldMatchRatatuiDisplay(t *testing.T) {
	tests := []struct {
		name  string
		color style.Color
		want  string
	}{
		{name: "default", color: style.Default, want: "Default"},
		{name: "reset", color: style.Reset, want: "Reset"},
		{name: "black", color: style.Black, want: "Black"},
		{name: "red", color: style.Red, want: "Red"},
		{name: "green", color: style.Green, want: "Green"},
		{name: "yellow", color: style.Yellow, want: "Yellow"},
		{name: "blue", color: style.Blue, want: "Blue"},
		{name: "magenta", color: style.Magenta, want: "Magenta"},
		{name: "cyan", color: style.Cyan, want: "Cyan"},
		{name: "white", color: style.White, want: "White"},
		{name: "light blue", color: style.LightBlue, want: "LightBlue"},
		{name: "light green", color: style.LightGreen, want: "LightGreen"},
		{name: "gray", color: style.Gray, want: "Gray"},
		{name: "dark gray", color: style.DarkGray, want: "DarkGray"},
		{name: "light red", color: style.LightRed, want: "LightRed"},
		{name: "light yellow", color: style.LightYellow, want: "LightYellow"},
		{name: "light magenta", color: style.LightMagenta, want: "LightMagenta"},
		{name: "light cyan", color: style.LightCyan, want: "LightCyan"},
		{name: "light white", color: style.LightWhite, want: "LightWhite"},
		{name: "indexed", color: style.IndexedColor(10), want: "10"},
		{name: "rgb", color: style.RGBColor(255, 0, 0), want: "#FF0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.String(); got != tt.want {
				t.Fatalf("%#v.String() = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}

func TestParseColor_shouldParseNamedIndexedAndRGBColors(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  style.Color
	}{
		{name: "default", value: "default", want: style.Default},
		{name: "reset", value: "reset", want: style.Reset},
		{name: "black", value: "black", want: style.Black},
		{name: "red", value: "red", want: style.Red},
		{name: "green", value: "green", want: style.Green},
		{name: "yellow", value: "yellow", want: style.Yellow},
		{name: "blue", value: "blue", want: style.Blue},
		{name: "magenta", value: "magenta", want: style.Magenta},
		{name: "cyan", value: "cyan", want: style.Cyan},
		{name: "gray", value: "gray", want: style.Gray},
		{name: "dark gray", value: "darkgray", want: style.DarkGray},
		{name: "light red", value: "lightred", want: style.LightRed},
		{name: "light green", value: "lightgreen", want: style.LightGreen},
		{name: "light yellow", value: "lightyellow", want: style.LightYellow},
		{name: "light blue", value: "lightblue", want: style.LightBlue},
		{name: "light magenta", value: "lightmagenta", want: style.LightMagenta},
		{name: "light cyan", value: "lightcyan", want: style.LightCyan},
		{name: "white", value: "white", want: style.White},
		{name: "light white alias", value: "lightwhite", want: style.White},
		{name: "light black alias", value: "lightblack", want: style.DarkGray},
		{name: "light gray alias", value: "lightgray", want: style.White},
		{name: "grey alias", value: "grey", want: style.Gray},
		{name: "silver alias", value: "silver", want: style.Gray},
		{name: "spaces ignored", value: "light black", want: style.DarkGray},
		{name: "dashes ignored", value: "light-white", want: style.White},
		{name: "underscores ignored", value: "light_gray", want: style.White},
		{name: "bright alias", value: "bright-black", want: style.DarkGray},
		{name: "camel case", value: "LightRed", want: style.LightRed},
		{name: "indexed", value: "10", want: style.IndexedColor(10)},
		{name: "rgb uppercase", value: "#FF0000", want: style.RGBColor(255, 0, 0)},
		{name: "rgb lowercase", value: "#00ff00", want: style.RGBColor(0, 255, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := style.ParseColor(tt.value)
			if err != nil {
				t.Fatalf("ParseColor(%q) returned error: %v", tt.value, err)
			}
			if got != tt.want {
				t.Fatalf("ParseColor(%q) = %#v, want %#v", tt.value, got, tt.want)
			}
		})
	}
}

func TestParseColor_unknownShouldReturnError(t *testing.T) {
	badColors := []string{
		"invalid_color",
		"abcdef0",
		" bcdefa",
		"#abcdef00",
		"#1🦀2",
		"resets",
		"lightblackk",
		"256",
	}

	for _, badColor := range badColors {
		if got, err := style.ParseColor(badColor); err == nil {
			t.Fatalf("ParseColor(%q) = %#v, want error", badColor, got)
		}
	}
}

func TestStylePatch_shouldTreatDefaultAsUnspecifiedAndResetAsExplicit(t *testing.T) {
	base := style.NewStyle().Fg(style.Red).Bg(style.Blue)

	defaultPatch := base.Patch(style.NewStyle())
	if defaultPatch != base {
		t.Fatalf("Patch(Default colors) = %#v, want %#v", defaultPatch, base)
	}

	resetPatch := base.Patch(style.ResetStyle())
	want := style.Style{Foreground: style.Reset, Background: style.Reset}
	if resetPatch != want {
		t.Fatalf("Patch(ResetStyle()) = %#v, want %#v", resetPatch, want)
	}
}
