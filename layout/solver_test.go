package layout

import (
	"slices"
	"testing"
)

func TestLayout_cacheSize(t *testing.T) {
	InitLayoutCache(1)
	t.Cleanup(func() { InitLayoutCache(0) })

	split := NewHorizontalLayout(Length(4), Fill(1)).
		Flex(FlexSpaceBetween).
		Spacing(2)
	area := NewRect(0, 0, 12, 3)
	wantSegments := []Rect{
		NewRect(0, 0, 4, 3),
		NewRect(6, 0, 6, 3),
	}
	wantSpacers := []Rect{
		NewRect(0, 0, 0, 3),
		NewRect(4, 0, 2, 3),
		NewRect(12, 0, 0, 3),
	}

	for i := range 3 {
		segments, spacers := split.SplitWithSpacers(area)
		if !slices.Equal(segments, wantSegments) {
			t.Fatalf("iteration %d segments mismatch\nwant: %#v\n got: %#v", i, wantSegments, segments)
		}
		if !slices.Equal(spacers, wantSpacers) {
			t.Fatalf("iteration %d spacers mismatch\nwant: %#v\n got: %#v", i, wantSpacers, spacers)
		}
	}

	InitLayoutCache(0)
	segments, spacers := split.SplitWithSpacers(area)
	if !slices.Equal(segments, wantSegments) {
		t.Fatalf("disabled cache segments mismatch\nwant: %#v\n got: %#v", wantSegments, segments)
	}
	if !slices.Equal(spacers, wantSpacers) {
		t.Fatalf("disabled cache spacers mismatch\nwant: %#v\n got: %#v", wantSpacers, spacers)
	}
}

func TestLayout_strengthIsValid(t *testing.T) {
	strengths := []struct {
		name     string
		strength solverStrength
	}{
		{name: "area bounds", strength: strengthAreaBounds},
		{name: "segment sizing", strength: strengthSegmentSizing},
		{name: "flex spacing", strength: strengthFlexSpacing},
		{name: "spacer growth", strength: strengthSpacerGrowth},
	}

	for _, tt := range strengths {
		if !tt.strength.isValid() {
			t.Fatalf("%s strength %v is invalid", tt.name, tt.strength)
		}
	}
}
