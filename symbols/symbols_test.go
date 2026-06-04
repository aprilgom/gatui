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

func TestScrollbarSymbols_defaults(t *testing.T) {
	if symbols.HorizontalScrollbarSet.Track != "═" || symbols.HorizontalScrollbarSet.Begin != "◄" || symbols.HorizontalScrollbarSet.End != "►" || symbols.HorizontalScrollbarSet.Thumb != "█" {
		t.Fatalf("unexpected horizontal scrollbar set: %+v", symbols.HorizontalScrollbarSet)
	}
	if symbols.VerticalScrollbarSet.Track != "║" || symbols.VerticalScrollbarSet.Begin != "▲" || symbols.VerticalScrollbarSet.End != "▼" || symbols.VerticalScrollbarSet.Thumb != "█" {
		t.Fatalf("unexpected vertical scrollbar set: %+v", symbols.VerticalScrollbarSet)
	}
}
