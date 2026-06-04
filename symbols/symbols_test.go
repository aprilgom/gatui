package symbols_test

import (
	"testing"

	"gatui/symbols"
)

func TestBorderSet_defaults(t *testing.T) {
	if symbols.PlainBorderSet.TopLeft != "┌" || symbols.RoundedBorderSet.TopLeft != "╭" || symbols.DoubleBorderSet.HorizontalTop != "═" {
		t.Fatalf("unexpected border presets: plain=%+v rounded=%+v double=%+v", symbols.PlainBorderSet, symbols.RoundedBorderSet, symbols.DoubleBorderSet)
	}
}

func TestBorderSet_extendedSetsShouldMatchRatatui(t *testing.T) {
	tests := []struct {
		name string
		set  symbols.BorderSet
		want string
	}{
		{
			name: "one eighth wide",
			set:  symbols.OneEighthWideBorderSet,
			want: "░░░░░░\n░▁▁▁▁░\n░▏░░▕░\n░▏░░▕░\n░▔▔▔▔░\n░░░░░░",
		},
		{
			name: "one eighth tall",
			set:  symbols.OneEighthTallBorderSet,
			want: "░░░░░░\n░▕▔▔▏░\n░▕░░▏░\n░▕░░▏░\n░▕▁▁▏░\n░░░░░░",
		},
		{
			name: "proportional wide",
			set:  symbols.ProportionalWideBorderSet,
			want: "░░░░░░\n░▄▄▄▄░\n░█░░█░\n░█░░█░\n░▀▀▀▀░\n░░░░░░",
		},
		{
			name: "proportional tall",
			set:  symbols.ProportionalTallBorderSet,
			want: "░░░░░░\n░█▀▀█░\n░█░░█░\n░█░░█░\n░█▄▄█░\n░░░░░░",
		},
		{
			name: "full",
			set:  symbols.FullBorderSet,
			want: "░░░░░░\n░████░\n░█░░█░\n░█░░█░\n░████░\n░░░░░░",
		},
		{
			name: "empty",
			set:  symbols.EmptyBorderSet,
			want: "░░░░░░\n░    ░\n░ ░░ ░\n░ ░░ ░\n░    ░\n░░░░░░",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderBorderSet(tt.set); got != tt.want {
				t.Fatalf("renderBorderSet(%s) =\n%s\nwant\n%s", tt.name, got, tt.want)
			}
		})
	}
}

func TestBorderSymbols_extendedConstantsShouldMatchRatatui(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"quadrant top left", symbols.BorderQuadrantTopLeft, "▘"},
		{"quadrant top right", symbols.BorderQuadrantTopRight, "▝"},
		{"quadrant bottom left", symbols.BorderQuadrantBottomLeft, "▖"},
		{"quadrant bottom right", symbols.BorderQuadrantBottomRight, "▗"},
		{"quadrant top half", symbols.BorderQuadrantTopHalf, "▀"},
		{"quadrant bottom half", symbols.BorderQuadrantBottomHalf, "▄"},
		{"quadrant left half", symbols.BorderQuadrantLeftHalf, "▌"},
		{"quadrant right half", symbols.BorderQuadrantRightHalf, "▐"},
		{"quadrant top left bottom left bottom right", symbols.BorderQuadrantTopLeftBottomLeftBottomRight, "▙"},
		{"quadrant top left top right bottom left", symbols.BorderQuadrantTopLeftTopRightBottomLeft, "▛"},
		{"quadrant top left top right bottom right", symbols.BorderQuadrantTopLeftTopRightBottomRight, "▜"},
		{"quadrant top right bottom left bottom right", symbols.BorderQuadrantTopRightBottomLeftBottomRight, "▟"},
		{"quadrant top left bottom right", symbols.BorderQuadrantTopLeftBottomRight, "▚"},
		{"quadrant top right bottom left", symbols.BorderQuadrantTopRightBottomLeft, "▞"},
		{"quadrant block", symbols.BorderQuadrantBlock, "█"},
		{"one eighth top", symbols.BorderOneEighthTop, "▔"},
		{"one eighth bottom", symbols.BorderOneEighthBottom, "▁"},
		{"one eighth left", symbols.BorderOneEighthLeft, "▏"},
		{"one eighth right", symbols.BorderOneEighthRight, "▕"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Fatalf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func renderBorderSet(set symbols.BorderSet) string {
	return "░░░░░░\n" +
		"░" + set.TopLeft + set.HorizontalTop + set.HorizontalTop + set.TopRight + "░\n" +
		"░" + set.VerticalLeft + "░░" + set.VerticalRight + "░\n" +
		"░" + set.VerticalLeft + "░░" + set.VerticalRight + "░\n" +
		"░" + set.BottomLeft + set.HorizontalBottom + set.HorizontalBottom + set.BottomRight + "░\n" +
		"░░░░░░"
}

func TestBorderSymbolMerge_matchesExistingMergedBorders(t *testing.T) {
	tests := []struct {
		name     string
		strategy symbols.MergeStrategy
		prev     string
		next     string
		want     string
	}{
		{name: "replace", strategy: symbols.MergeStrategyReplace, prev: "│", next: "━", want: "━"},
		{name: "exact plain cross", strategy: symbols.MergeStrategyExact, prev: "│", next: "─", want: "┼"},
		{name: "exact mixed plain thick cross", strategy: symbols.MergeStrategyExact, prev: "│", next: "━", want: "┿"},
		{name: "exact falls back to replace", strategy: symbols.MergeStrategyExact, prev: "┘", next: "╔", want: "╔"},
		{name: "fuzzy double thick", strategy: symbols.MergeStrategyFuzzy, prev: "┘", next: "╔", want: "╬"},
		{name: "fuzzy rounded plain", strategy: symbols.MergeStrategyFuzzy, prev: "┘", next: "╭", want: "┼"},
		{name: "fuzzy dashed thick", strategy: symbols.MergeStrategyFuzzy, prev: "╎", next: "╍", want: "┿"},
		{name: "non border previous wins when next is border", strategy: symbols.MergeStrategyExact, prev: "a", next: "╭", want: "a"},
		{name: "non border next wins", strategy: symbols.MergeStrategyExact, prev: "┌", next: "a", want: "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := symbols.MergeBorderSymbols(tt.strategy, tt.prev, tt.next); got != tt.want {
				t.Fatalf("MergeBorderSymbols(%s, %q, %q) = %q, want %q", tt.strategy, tt.prev, tt.next, got, tt.want)
			}
		})
	}
}

func TestBarSet_presets(t *testing.T) {
	if symbols.NineLevelBarSet.Empty != " " || symbols.NineLevelBarSet.OneEighth != "▁" || symbols.NineLevelBarSet.Full != "█" {
		t.Fatalf("unexpected nine-level bar set: %+v", symbols.NineLevelBarSet)
	}
	if symbols.ThreeLevelBarSet.OneEighth != "▄" || symbols.ThreeLevelBarSet.FiveEighths != "█" {
		t.Fatalf("unexpected three-level bar set: %+v", symbols.ThreeLevelBarSet)
	}
}

func TestLineSet_shouldMatchRatatuiNormalRoundedDoubleThickAndDashed(t *testing.T) {
	tests := []struct {
		name string
		got  symbols.LineSet
		want symbols.LineSet
	}{
		{name: "normal", got: symbols.NormalLineSet, want: symbols.LineSet{
			Vertical: "│", Horizontal: "─", TopRight: "┐", TopLeft: "┌", BottomRight: "┘", BottomLeft: "└",
			VerticalLeft: "┤", VerticalRight: "├", HorizontalDown: "┬", HorizontalUp: "┴", Cross: "┼",
		}},
		{name: "rounded", got: symbols.RoundedLineSet, want: symbols.LineSet{
			Vertical: "│", Horizontal: "─", TopRight: "╮", TopLeft: "╭", BottomRight: "╯", BottomLeft: "╰",
			VerticalLeft: "┤", VerticalRight: "├", HorizontalDown: "┬", HorizontalUp: "┴", Cross: "┼",
		}},
		{name: "double", got: symbols.DoubleLineSet, want: symbols.LineSet{
			Vertical: "║", Horizontal: "═", TopRight: "╗", TopLeft: "╔", BottomRight: "╝", BottomLeft: "╚",
			VerticalLeft: "╣", VerticalRight: "╠", HorizontalDown: "╦", HorizontalUp: "╩", Cross: "╬",
		}},
		{name: "thick", got: symbols.ThickLineSet, want: symbols.LineSet{
			Vertical: "┃", Horizontal: "━", TopRight: "┓", TopLeft: "┏", BottomRight: "┛", BottomLeft: "┗",
			VerticalLeft: "┫", VerticalRight: "┣", HorizontalDown: "┳", HorizontalUp: "┻", Cross: "╋",
		}},
		{name: "light double dashed", got: symbols.LightDoubleDashedLineSet, want: symbols.LineSet{
			Vertical: "╎", Horizontal: "╌", TopRight: "┐", TopLeft: "┌", BottomRight: "┘", BottomLeft: "└",
			VerticalLeft: "┤", VerticalRight: "├", HorizontalDown: "┬", HorizontalUp: "┴", Cross: "┼",
		}},
		{name: "heavy double dashed", got: symbols.HeavyDoubleDashedLineSet, want: symbols.LineSet{
			Vertical: "╏", Horizontal: "╍", TopRight: "┓", TopLeft: "┏", BottomRight: "┛", BottomLeft: "┗",
			VerticalLeft: "┫", VerticalRight: "┣", HorizontalDown: "┳", HorizontalUp: "┻", Cross: "╋",
		}},
		{name: "light triple dashed", got: symbols.LightTripleDashedLineSet, want: symbols.LineSet{
			Vertical: "┆", Horizontal: "┄", TopRight: "┐", TopLeft: "┌", BottomRight: "┘", BottomLeft: "└",
			VerticalLeft: "┤", VerticalRight: "├", HorizontalDown: "┬", HorizontalUp: "┴", Cross: "┼",
		}},
		{name: "heavy triple dashed", got: symbols.HeavyTripleDashedLineSet, want: symbols.LineSet{
			Vertical: "┇", Horizontal: "┅", TopRight: "┓", TopLeft: "┏", BottomRight: "┛", BottomLeft: "┗",
			VerticalLeft: "┫", VerticalRight: "┣", HorizontalDown: "┳", HorizontalUp: "┻", Cross: "╋",
		}},
		{name: "light quadruple dashed", got: symbols.LightQuadrupleDashedLineSet, want: symbols.LineSet{
			Vertical: "┊", Horizontal: "┈", TopRight: "┐", TopLeft: "┌", BottomRight: "┘", BottomLeft: "└",
			VerticalLeft: "┤", VerticalRight: "├", HorizontalDown: "┬", HorizontalUp: "┴", Cross: "┼",
		}},
		{name: "heavy quadruple dashed", got: symbols.HeavyQuadrupleDashedLineSet, want: symbols.LineSet{
			Vertical: "┋", Horizontal: "┉", TopRight: "┓", TopLeft: "┏", BottomRight: "┛", BottomLeft: "┗",
			VerticalLeft: "┫", VerticalRight: "┣", HorizontalDown: "┳", HorizontalUp: "┻", Cross: "╋",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("LineSet = %+v, want %+v", tt.got, tt.want)
			}
		})
	}
	if symbols.LineVertical != symbols.NormalLineSet.Vertical || symbols.LineHorizontal != symbols.NormalLineSet.Horizontal || symbols.LineBottomLeft != symbols.NormalLineSet.BottomLeft {
		t.Fatalf("normal line constants should match NormalLineSet: %+v", symbols.NormalLineSet)
	}
}

func TestBlockSymbols_shouldMatchRatatuiNineLevelBlocks(t *testing.T) {
	want := []string{" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
	got := []string{
		symbols.BlockEmpty,
		symbols.BlockOneEighth,
		symbols.BlockOneQuarter,
		symbols.BlockThreeEighths,
		symbols.BlockHalf,
		symbols.BlockFiveEighths,
		symbols.BlockThreeQuarters,
		symbols.BlockSevenEighths,
		symbols.BlockFull,
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("block symbol %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestShadeSymbols_shouldMatchRatatuiShadeConstants(t *testing.T) {
	got := []string{symbols.ShadeEmpty, symbols.ShadeLight, symbols.ShadeMedium, symbols.ShadeDark, symbols.ShadeFull}
	want := []string{" ", "░", "▒", "▓", "█"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("shade symbol %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSparklineBarSet_presets(t *testing.T) {
	if symbols.NineLevelSparklineBarSet().OneEighth != "▁" || symbols.NineLevelSparklineBarSet().Full != "█" {
		t.Fatalf("unexpected nine-level sparkline set: %+v", symbols.NineLevelSparklineBarSet())
	}
	if symbols.ThreeLevelSparklineBarSet().OneEighth != " " || symbols.ThreeLevelSparklineBarSet().ThreeQuarters != "▄" {
		t.Fatalf("unexpected three-level sparkline set: %+v", symbols.ThreeLevelSparklineBarSet())
	}
}

func TestCanvasMarker_renderSymbols(t *testing.T) {
	if symbols.CanvasMarkerCustom("xy") != symbols.CanvasMarker("custom:x") {
		t.Fatalf("custom marker should keep first rune")
	}
	if symbols.BrailleSymbol(0b00000001) != '⠁' {
		t.Fatalf("unexpected braille symbol")
	}
	if symbols.QuadrantSymbol(0b1111) != '█' {
		t.Fatalf("unexpected quadrant symbol")
	}
	if symbols.SextantSymbol(0b111111) != "█" || symbols.OctantSymbol(0xff) != "█" {
		t.Fatalf("unexpected block subdivision symbols")
	}
}

func TestPixelSymbols_shouldExposeRatatuiPseudoPixelLookups(t *testing.T) {
	if symbols.Quadrants[0] != ' ' || symbols.Quadrants[3] != '▀' || symbols.Quadrants[15] != '█' {
		t.Fatalf("unexpected quadrant lookup: %q %q %q", symbols.Quadrants[0], symbols.Quadrants[3], symbols.Quadrants[15])
	}
	if symbols.Sextants[0] != " " || symbols.Sextants[21] != "▌" || symbols.Sextants[42] != "▐" || symbols.Sextants[63] != "█" {
		t.Fatalf("unexpected sextant lookup")
	}
	if symbols.Octants[0] != " " || symbols.Octants[5] != "▘" || symbols.Octants[15] != "▀" || symbols.Octants[255] != "█" {
		t.Fatalf("unexpected octant lookup")
	}
}

func TestHalfBlockSymbols_shouldMatchRatatuiConstants(t *testing.T) {
	if symbols.HalfBlockUpper != "▀" || symbols.HalfBlockLower != "▄" || symbols.HalfBlockFull != "█" {
		t.Fatalf("unexpected half-block constants: %q %q %q", symbols.HalfBlockUpper, symbols.HalfBlockLower, symbols.HalfBlockFull)
	}
}

func TestIsCanvasDatasetSymbol_shouldRecognizeOnlyDatasetMarkers(t *testing.T) {
	tests := []struct {
		symbol string
		want   bool
	}{
		{symbol: symbols.CanvasDotSymbol, want: true},
		{symbol: symbols.CanvasBlockSymbol, want: true},
		{symbol: symbols.CanvasBarSymbol, want: true},
		{symbol: string(symbols.BrailleSymbol(0xff)), want: true},
		{symbol: "x", want: false},
		{symbol: "xy", want: false},
		{symbol: "", want: false},
	}
	for _, tt := range tests {
		if got := symbols.IsCanvasDatasetSymbol(tt.symbol); got != tt.want {
			t.Fatalf("IsCanvasDatasetSymbol(%q) = %v, want %v", tt.symbol, got, tt.want)
		}
	}
}

func TestScrollbarSymbols_defaults(t *testing.T) {
	if symbols.HorizontalScrollbarSet.Track != "═" || symbols.HorizontalScrollbarSet.Begin != "◄" || symbols.HorizontalScrollbarSet.End != "►" || symbols.HorizontalScrollbarSet.Thumb != "█" {
		t.Fatalf("unexpected horizontal scrollbar set: %+v", symbols.HorizontalScrollbarSet)
	}
	if symbols.VerticalScrollbarSet.Track != "║" || symbols.VerticalScrollbarSet.Begin != "▲" || symbols.VerticalScrollbarSet.End != "▼" || symbols.VerticalScrollbarSet.Thumb != "█" {
		t.Fatalf("unexpected vertical scrollbar set: %+v", symbols.VerticalScrollbarSet)
	}
}
