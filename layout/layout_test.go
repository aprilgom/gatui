package layout_test

import (
	"fmt"
	"slices"
	"testing"

	"gatui/layout"
)

func TestLayout_Split_shouldMatchRatatuiBasicConstraintPriority(t *testing.T) {
	tests := []struct {
		name        string
		constraints []layout.Constraint
		want        []layout.Rect
	}{
		{
			name:        "length min",
			constraints: []layout.Constraint{layout.Length(20), layout.Min(10)},
			want: []layout.Rect{
				layout.NewRect(10, 5, 20, 3),
				layout.NewRect(30, 5, 80, 3),
			},
		},
		{
			name:        "min length",
			constraints: []layout.Constraint{layout.Min(10), layout.Length(20)},
			want: []layout.Rect{
				layout.NewRect(10, 5, 80, 3),
				layout.NewRect(90, 5, 20, 3),
			},
		},
		{
			name:        "length percentage",
			constraints: []layout.Constraint{layout.Length(20), layout.Percentage(30)},
			want: []layout.Rect{
				layout.NewRect(10, 5, 20, 3),
				layout.NewRect(30, 5, 80, 3),
			},
		},
		{
			name:        "percentage length",
			constraints: []layout.Constraint{layout.Percentage(30), layout.Length(20)},
			want: []layout.Rect{
				layout.NewRect(10, 5, 30, 3),
				layout.NewRect(40, 5, 70, 3),
			},
		},
		{
			name:        "length ratio",
			constraints: []layout.Constraint{layout.Length(20), layout.Ratio(1, 4)},
			want: []layout.Rect{
				layout.NewRect(10, 5, 20, 3),
				layout.NewRect(30, 5, 80, 3),
			},
		},
		{
			name:        "ratio length",
			constraints: []layout.Constraint{layout.Ratio(1, 4), layout.Length(20)},
			want: []layout.Rect{
				layout.NewRect(10, 5, 25, 3),
				layout.NewRect(35, 5, 75, 3),
			},
		},
	}

	area := layout.NewRect(10, 5, 100, 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := layout.NewLayout(layout.Horizontal).
				Constraints(tt.constraints...).
				Split(area)

			if !slices.Equal(got, tt.want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestConstraint_shouldExposeMaxAndFillKinds(t *testing.T) {
	maxConstraint := layout.Max(10)
	if !maxConstraint.IsMax() {
		t.Fatalf("Max(10).IsMax() = false, want true")
	}
	if maxConstraint.Value() != 10 {
		t.Fatalf("Max(10).Value() = %d, want 10", maxConstraint.Value())
	}

	fillConstraint := layout.Fill(2)
	if !fillConstraint.IsFill() {
		t.Fatalf("Fill(2).IsFill() = false, want true")
	}
	if fillConstraint.Value() != 2 {
		t.Fatalf("Fill(2).Value() = %d, want 2", fillConstraint.Value())
	}
}

func TestConstraint_FromFills(t *testing.T) {
	got := layout.FromFills(1, 2, 3)
	want := []layout.Constraint{layout.Fill(1), layout.Fill(2), layout.Fill(3)}

	if !slices.Equal(got, want) {
		t.Fatalf("constraints mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestConstraint_FromLengths(t *testing.T) {
	got := layout.FromLengths(1, 2, 3)
	want := []layout.Constraint{layout.Length(1), layout.Length(2), layout.Length(3)}

	if !slices.Equal(got, want) {
		t.Fatalf("constraints mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestConstraint_FromMaxes(t *testing.T) {
	got := layout.FromMaxes(1, 2, 3)
	want := []layout.Constraint{layout.Max(1), layout.Max(2), layout.Max(3)}

	if !slices.Equal(got, want) {
		t.Fatalf("constraints mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestConstraint_FromMins(t *testing.T) {
	got := layout.FromMins(1, 2, 3)
	want := []layout.Constraint{layout.Min(1), layout.Min(2), layout.Min(3)}

	if !slices.Equal(got, want) {
		t.Fatalf("constraints mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestConstraint_FromPercentages(t *testing.T) {
	got := layout.FromPercentages(25, 50, 25)
	want := []layout.Constraint{layout.Percentage(25), layout.Percentage(50), layout.Percentage(25)}

	if !slices.Equal(got, want) {
		t.Fatalf("constraints mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestConstraint_FromRatios(t *testing.T) {
	got := layout.FromRatios([2]int{1, 4}, [2]int{1, 2}, [2]int{1, 4})
	want := []layout.Constraint{layout.Ratio(1, 4), layout.Ratio(1, 2), layout.Ratio(1, 4)}

	if !slices.Equal(got, want) {
		t.Fatalf("constraints mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestDirection_Other(t *testing.T) {
	if got := layout.Horizontal.Other(); got != layout.Vertical {
		t.Fatalf("Horizontal.Other() = %v, want %v", got, layout.Vertical)
	}
	if got := layout.Vertical.Other(); got != layout.Horizontal {
		t.Fatalf("Vertical.Other() = %v, want %v", got, layout.Horizontal)
	}
}

func TestNewVerticalLayout_shouldMatchRatatuiVerticalConstructor(t *testing.T) {
	got := layout.NewVerticalLayout(layout.Min(0)).
		Split(layout.NewRect(0, 0, 5, 10))
	want := layout.NewLayout(layout.Vertical).
		Constraints(layout.Min(0)).
		Split(layout.NewRect(0, 0, 5, 10))

	if !slices.Equal(got, want) {
		t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestNewHorizontalLayout_shouldMatchRatatuiHorizontalConstructor(t *testing.T) {
	got := layout.NewHorizontalLayout(layout.Min(0)).
		Split(layout.NewRect(0, 0, 10, 5))
	want := layout.NewLayout(layout.Horizontal).
		Constraints(layout.Min(0)).
		Split(layout.NewRect(0, 0, 10, 5))

	if !slices.Equal(got, want) {
		t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestLayout_SplitN_returnsExpectedRects(t *testing.T) {
	area := layout.NewRect(0, 0, 10, 5)
	split := layout.NewHorizontalLayout(layout.Length(3), layout.Fill(1))
	want := []layout.Rect{
		layout.NewRect(0, 0, 3, 5),
		layout.NewRect(3, 0, 7, 5),
	}

	got := split.SplitN(area, 2)

	if !slices.Equal(got, want) {
		t.Fatalf("SplitN(%#v, 2) mismatch\nwant: %#v\n got: %#v", area, want, got)
	}
}

func TestLayout_SplitN_invalidNumberOfRectsPanics(t *testing.T) {
	area := layout.NewRect(0, 0, 10, 5)
	split := layout.NewHorizontalLayout(layout.Length(3), layout.Fill(1))
	wantPanic := "invalid number of rects: expected 3, found 2"

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("SplitN(%#v, 3) did not panic", area)
		}
		if fmt.Sprint(got) != wantPanic {
			t.Fatalf("SplitN(%#v, 3) panic = %q, want %q", area, got, wantPanic)
		}
	}()

	_ = split.SplitN(area, 3)
}

func TestLayout_TrySplitN_invalidNumberOfRectsReturnsError(t *testing.T) {
	area := layout.NewRect(0, 0, 10, 5)
	split := layout.NewHorizontalLayout(layout.Length(3), layout.Fill(1))
	wantErr := "invalid number of rects: expected 3, found 2"

	got, err := split.TrySplitN(area, 3)

	if err == nil {
		t.Fatalf("TrySplitN(%#v, 3) error = nil, want %q", area, wantErr)
	}
	if err.Error() != wantErr {
		t.Fatalf("TrySplitN(%#v, 3) error = %q, want %q", area, err, wantErr)
	}
	if got != nil {
		t.Fatalf("TrySplitN(%#v, 3) rects = %#v, want nil", area, got)
	}
}

func TestLayoutConstructors_shouldCopyConstraintsFromCallerSlice(t *testing.T) {
	constraints := []layout.Constraint{layout.Length(2), layout.Fill(1)}
	constructed := layout.NewHorizontalLayout(constraints...)

	constraints[0] = layout.Length(8)

	got := constructed.Split(layout.NewRect(0, 0, 10, 1))
	want := []layout.Rect{
		layout.NewRect(0, 0, 2, 1),
		layout.NewRect(2, 0, 8, 1),
	}
	if !slices.Equal(got, want) {
		t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestLayout_Direction_shouldPreserveConstraintsAndChangeSplitAxis(t *testing.T) {
	horizontal := layout.NewVerticalLayout(layout.Length(5), layout.Fill(1)).
		Direction(layout.Horizontal).
		Split(layout.NewRect(0, 0, 10, 10))
	wantHorizontal := []layout.Rect{
		layout.NewRect(0, 0, 5, 10),
		layout.NewRect(5, 0, 5, 10),
	}
	if !slices.Equal(horizontal, wantHorizontal) {
		t.Fatalf("horizontal rects mismatch\nwant: %#v\n got: %#v", wantHorizontal, horizontal)
	}

	vertical := layout.NewHorizontalLayout(layout.Length(5), layout.Fill(1)).
		Direction(layout.Vertical).
		Split(layout.NewRect(0, 0, 10, 10))
	wantVertical := []layout.Rect{
		layout.NewRect(0, 0, 10, 5),
		layout.NewRect(0, 5, 10, 5),
	}
	if !slices.Equal(vertical, wantVertical) {
		t.Fatalf("vertical rects mismatch\nwant: %#v\n got: %#v", wantVertical, vertical)
	}
}

func TestLayout_Split_shouldApplyMarginsBeforeSplitting(t *testing.T) {
	tests := []struct {
		name   string
		layout layout.Layout
		want   []layout.Rect
	}{
		{
			name:   "uniform",
			layout: layout.NewVerticalLayout(layout.Fill(1)).UniformMargin(2),
			want:   []layout.Rect{layout.NewRect(2, 2, 6, 6)},
		},
		{
			name:   "horizontal",
			layout: layout.NewVerticalLayout(layout.Fill(1)).HorizontalMargin(2),
			want:   []layout.Rect{layout.NewRect(2, 0, 6, 10)},
		},
		{
			name:   "vertical",
			layout: layout.NewVerticalLayout(layout.Fill(1)).VerticalMargin(2),
			want:   []layout.Rect{layout.NewRect(0, 2, 10, 6)},
		},
		{
			name:   "axes",
			layout: layout.NewVerticalLayout(layout.Fill(1)).Margin(1, 2),
			want:   []layout.Rect{layout.NewRect(1, 2, 8, 6)},
		},
		{
			name:   "oversized",
			layout: layout.NewVerticalLayout(layout.Fill(1)).UniformMargin(6),
			want:   []layout.Rect{layout.NewRect(0, 0, 0, 0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.layout.Split(layout.NewRect(0, 0, 10, 10))
			if !slices.Equal(got, tt.want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestLayout_Split_shouldMatchRatatuiMaxAndFillConstraints(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		constraints []layout.Constraint
		want        []layout.Rect
	}{
		{
			name:        "max min",
			width:       100,
			constraints: []layout.Constraint{layout.Max(100), layout.Min(0)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 100, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name:        "min max",
			width:       100,
			constraints: []layout.Constraint{layout.Min(0), layout.Max(100)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(0, 0, 100, 1),
			},
		},
		{
			name:        "length max",
			width:       100,
			constraints: []layout.Constraint{layout.Length(25), layout.Max(100)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 25, 1),
				layout.NewRect(25, 0, 75, 1),
			},
		},
		{
			name:        "max length min",
			width:       100,
			constraints: []layout.Constraint{layout.Max(25), layout.Length(25), layout.Min(25)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 25, 1),
				layout.NewRect(25, 0, 25, 1),
				layout.NewRect(50, 0, 50, 1),
			},
		},
		{
			name:        "fill length fill equal",
			width:       100,
			constraints: []layout.Constraint{layout.Fill(1), layout.Length(10), layout.Fill(1)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 45, 1),
				layout.NewRect(45, 0, 10, 1),
				layout.NewRect(55, 0, 45, 1),
			},
		},
		{
			name:        "fill length fill weighted",
			width:       100,
			constraints: []layout.Constraint{layout.Fill(1), layout.Length(10), layout.Fill(2)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 30, 1),
				layout.NewRect(30, 0, 10, 1),
				layout.NewRect(40, 0, 60, 1),
			},
		},
		{
			name:        "zero fill around positive fill",
			width:       100,
			constraints: []layout.Constraint{layout.Fill(0), layout.Fill(1), layout.Fill(0)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(0, 0, 100, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name:        "length fill fill constrained width",
			width:       50,
			constraints: []layout.Constraint{layout.Length(10), layout.Fill(2), layout.Fill(1)},
			want: []layout.Rect{
				layout.NewRect(0, 0, 10, 1),
				layout.NewRect(10, 0, 27, 1),
				layout.NewRect(37, 0, 13, 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := layout.NewLayout(layout.Horizontal).
				Constraints(tt.constraints...).
				Split(layout.NewRect(0, 0, tt.width, 1))

			if !slices.Equal(got, tt.want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestLayout_Split_shouldAlignSingleFixedSegmentByFlex(t *testing.T) {
	tests := []struct {
		name string
		flex layout.Flex
		want layout.Rect
	}{
		{name: "legacy stretches segment", flex: layout.FlexLegacy, want: layout.NewRect(0, 0, 100, 1)},
		{name: "start", flex: layout.FlexStart, want: layout.NewRect(0, 0, 50, 1)},
		{name: "end", flex: layout.FlexEnd, want: layout.NewRect(50, 0, 50, 1)},
		{name: "center", flex: layout.FlexCenter, want: layout.NewRect(25, 0, 50, 1)},
		{name: "space between stretches segment", flex: layout.FlexSpaceBetween, want: layout.NewRect(0, 0, 100, 1)},
		{name: "space around centers segment", flex: layout.FlexSpaceAround, want: layout.NewRect(25, 0, 50, 1)},
		{name: "space evenly centers segment", flex: layout.FlexSpaceEvenly, want: layout.NewRect(25, 0, 50, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := layout.NewLayout(layout.Horizontal).
				Flex(tt.flex).
				Constraints(layout.Length(50)).
				Split(layout.NewRect(0, 0, 100, 1))

			want := []layout.Rect{tt.want}
			if !slices.Equal(got, want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", want, got)
			}
		})
	}
}

func TestLayout_Split_shouldAlignTwoFixedSegmentsByFlex(t *testing.T) {
	tests := []struct {
		name string
		flex layout.Flex
		want []layout.Rect
	}{
		{
			name: "legacy stretches last segment",
			flex: layout.FlexLegacy,
			want: []layout.Rect{
				layout.NewRect(0, 0, 25, 1),
				layout.NewRect(25, 0, 75, 1),
			},
		},
		{
			name: "start",
			flex: layout.FlexStart,
			want: []layout.Rect{
				layout.NewRect(0, 0, 25, 1),
				layout.NewRect(25, 0, 25, 1),
			},
		},
		{
			name: "center",
			flex: layout.FlexCenter,
			want: []layout.Rect{
				layout.NewRect(25, 0, 25, 1),
				layout.NewRect(50, 0, 25, 1),
			},
		},
		{
			name: "end",
			flex: layout.FlexEnd,
			want: []layout.Rect{
				layout.NewRect(50, 0, 25, 1),
				layout.NewRect(75, 0, 25, 1),
			},
		},
		{
			name: "space between",
			flex: layout.FlexSpaceBetween,
			want: []layout.Rect{
				layout.NewRect(0, 0, 25, 1),
				layout.NewRect(75, 0, 25, 1),
			},
		},
		{
			name: "space around",
			flex: layout.FlexSpaceAround,
			want: []layout.Rect{
				layout.NewRect(13, 0, 25, 1),
				layout.NewRect(63, 0, 25, 1),
			},
		},
		{
			name: "space evenly",
			flex: layout.FlexSpaceEvenly,
			want: []layout.Rect{
				layout.NewRect(17, 0, 25, 1),
				layout.NewRect(58, 0, 25, 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := layout.NewLayout(layout.Horizontal).
				Flex(tt.flex).
				Constraints(layout.Length(25), layout.Length(25)).
				Split(layout.NewRect(0, 0, 100, 1))

			if !slices.Equal(got, tt.want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestLayout_Split_shouldLetMinAndFillAbsorbSurplusBeforeFlexAlignment(t *testing.T) {
	tests := []struct {
		name        string
		constraints []layout.Constraint
	}{
		{name: "min", constraints: []layout.Constraint{layout.Min(25), layout.Min(25)}},
		{name: "fill", constraints: []layout.Constraint{layout.Fill(1), layout.Fill(1)}},
	}
	flexes := []layout.Flex{
		layout.FlexLegacy,
		layout.FlexStart,
		layout.FlexEnd,
		layout.FlexCenter,
		layout.FlexSpaceBetween,
		layout.FlexSpaceAround,
		layout.FlexSpaceEvenly,
	}
	defaultWant := []layout.Rect{
		layout.NewRect(0, 0, 50, 1),
		layout.NewRect(50, 0, 50, 1),
	}

	for _, tt := range tests {
		for _, flex := range flexes {
			t.Run(tt.name, func(t *testing.T) {
				got := layout.NewLayout(layout.Horizontal).
					Flex(flex).
					Constraints(tt.constraints...).
					Split(layout.NewRect(0, 0, 100, 1))

				want := defaultWant
				if tt.name == "min" && flex == layout.FlexLegacy {
					want = []layout.Rect{
						layout.NewRect(0, 0, 25, 1),
						layout.NewRect(25, 0, 75, 1),
					}
				}
				if !slices.Equal(got, want) {
					t.Fatalf("rects mismatch for flex %d\nwant: %#v\n got: %#v", flex, want, got)
				}
			})
		}
	}
}

func TestLayout_Split_shouldApplyFlexAlignmentToMaxConstraints(t *testing.T) {
	got := layout.NewLayout(layout.Horizontal).
		Flex(layout.FlexCenter).
		Constraints(layout.Max(25), layout.Max(25)).
		Split(layout.NewRect(0, 0, 100, 1))

	want := []layout.Rect{
		layout.NewRect(25, 0, 25, 1),
		layout.NewRect(50, 0, 25, 1),
	}
	if !slices.Equal(got, want) {
		t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestLayout_Split_shouldMatchRatatuiFlexSpacing(t *testing.T) {
	tests := []struct {
		name    string
		flex    layout.Flex
		spacing int
		want    []layout.Rect
	}{
		{
			name:    "start positive spacing",
			flex:    layout.FlexStart,
			spacing: 2,
			want: []layout.Rect{
				layout.NewRect(0, 0, 20, 1),
				layout.NewRect(22, 0, 20, 1),
				layout.NewRect(44, 0, 20, 1),
			},
		},
		{
			name:    "center positive spacing",
			flex:    layout.FlexCenter,
			spacing: 2,
			want: []layout.Rect{
				layout.NewRect(18, 0, 20, 1),
				layout.NewRect(40, 0, 20, 1),
				layout.NewRect(62, 0, 20, 1),
			},
		},
		{
			name:    "end positive spacing",
			flex:    layout.FlexEnd,
			spacing: 2,
			want: []layout.Rect{
				layout.NewRect(36, 0, 20, 1),
				layout.NewRect(58, 0, 20, 1),
				layout.NewRect(80, 0, 20, 1),
			},
		},
		{
			name:    "legacy positive spacing",
			flex:    layout.FlexLegacy,
			spacing: 2,
			want: []layout.Rect{
				layout.NewRect(0, 0, 20, 1),
				layout.NewRect(22, 0, 20, 1),
				layout.NewRect(44, 0, 56, 1),
			},
		},
		{
			name:    "space between positive spacing",
			flex:    layout.FlexSpaceBetween,
			spacing: 2,
			want: []layout.Rect{
				layout.NewRect(0, 0, 20, 1),
				layout.NewRect(40, 0, 20, 1),
				layout.NewRect(80, 0, 20, 1),
			},
		},
		{
			name:    "space evenly positive spacing",
			flex:    layout.FlexSpaceEvenly,
			spacing: 2,
			want: []layout.Rect{
				layout.NewRect(10, 0, 20, 1),
				layout.NewRect(40, 0, 20, 1),
				layout.NewRect(70, 0, 20, 1),
			},
		},
		{
			name:    "space around positive spacing",
			flex:    layout.FlexSpaceAround,
			spacing: 2,
			want: []layout.Rect{
				layout.NewRect(7, 0, 20, 1),
				layout.NewRect(40, 0, 20, 1),
				layout.NewRect(73, 0, 20, 1),
			},
		},
		{
			name:    "start negative overlap",
			flex:    layout.FlexStart,
			spacing: -1,
			want: []layout.Rect{
				layout.NewRect(0, 0, 20, 1),
				layout.NewRect(19, 0, 20, 1),
				layout.NewRect(38, 0, 20, 1),
			},
		},
		{
			name:    "center negative overlap",
			flex:    layout.FlexCenter,
			spacing: -1,
			want: []layout.Rect{
				layout.NewRect(21, 0, 20, 1),
				layout.NewRect(40, 0, 20, 1),
				layout.NewRect(59, 0, 20, 1),
			},
		},
		{
			name:    "end negative overlap",
			flex:    layout.FlexEnd,
			spacing: -1,
			want: []layout.Rect{
				layout.NewRect(42, 0, 20, 1),
				layout.NewRect(61, 0, 20, 1),
				layout.NewRect(80, 0, 20, 1),
			},
		},
		{
			name:    "legacy negative overlap",
			flex:    layout.FlexLegacy,
			spacing: -1,
			want: []layout.Rect{
				layout.NewRect(0, 0, 20, 1),
				layout.NewRect(19, 0, 20, 1),
				layout.NewRect(38, 0, 62, 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := layout.NewLayout(layout.Horizontal).
				Flex(tt.flex).
				Spacing(tt.spacing).
				Constraints(layout.Length(20), layout.Length(20), layout.Length(20)).
				Split(layout.NewRect(0, 0, 100, 1))

			if !slices.Equal(got, tt.want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestLayout_SplitWithSpacers_shouldMatchRatatuiSpacerRects(t *testing.T) {
	tests := []struct {
		name    string
		flex    layout.Flex
		spacing int
		want    []layout.Rect
	}{
		{
			name: "legacy no spacing",
			flex: layout.FlexLegacy,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(10, 0, 0, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name: "space between no spacing",
			flex: layout.FlexSpaceBetween,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(10, 0, 80, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name: "space evenly no spacing",
			flex: layout.FlexSpaceEvenly,
			want: []layout.Rect{
				layout.NewRect(0, 0, 27, 1),
				layout.NewRect(37, 0, 26, 1),
				layout.NewRect(73, 0, 27, 1),
			},
		},
		{
			name: "space around no spacing",
			flex: layout.FlexSpaceAround,
			want: []layout.Rect{
				layout.NewRect(0, 0, 20, 1),
				layout.NewRect(30, 0, 40, 1),
				layout.NewRect(80, 0, 20, 1),
			},
		},
		{
			name: "start no spacing",
			flex: layout.FlexStart,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(10, 0, 0, 1),
				layout.NewRect(20, 0, 80, 1),
			},
		},
		{
			name: "center no spacing",
			flex: layout.FlexCenter,
			want: []layout.Rect{
				layout.NewRect(0, 0, 40, 1),
				layout.NewRect(50, 0, 0, 1),
				layout.NewRect(60, 0, 40, 1),
			},
		},
		{
			name: "end no spacing",
			flex: layout.FlexEnd,
			want: []layout.Rect{
				layout.NewRect(0, 0, 80, 1),
				layout.NewRect(90, 0, 0, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name:    "legacy positive spacing",
			flex:    layout.FlexLegacy,
			spacing: 5,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(10, 0, 5, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name:    "start positive spacing",
			flex:    layout.FlexStart,
			spacing: 5,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(10, 0, 5, 1),
				layout.NewRect(25, 0, 75, 1),
			},
		},
		{
			name:    "center positive spacing",
			flex:    layout.FlexCenter,
			spacing: 5,
			want: []layout.Rect{
				layout.NewRect(0, 0, 38, 1),
				layout.NewRect(48, 0, 5, 1),
				layout.NewRect(63, 0, 37, 1),
			},
		},
		{
			name:    "end positive spacing",
			flex:    layout.FlexEnd,
			spacing: 5,
			want: []layout.Rect{
				layout.NewRect(0, 0, 75, 1),
				layout.NewRect(85, 0, 5, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name:    "start negative spacing clamps overlap",
			flex:    layout.FlexStart,
			spacing: -1,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(10, 0, 0, 1),
				layout.NewRect(19, 0, 81, 1),
			},
		},
		{
			name:    "center negative spacing clamps overlap",
			flex:    layout.FlexCenter,
			spacing: -1,
			want: []layout.Rect{
				layout.NewRect(0, 0, 41, 1),
				layout.NewRect(51, 0, 0, 1),
				layout.NewRect(60, 0, 40, 1),
			},
		},
		{
			name:    "end negative spacing clamps overlap",
			flex:    layout.FlexEnd,
			spacing: -1,
			want: []layout.Rect{
				layout.NewRect(0, 0, 81, 1),
				layout.NewRect(91, 0, 0, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name:    "legacy too much spacing",
			flex:    layout.FlexLegacy,
			spacing: 200,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(0, 0, 100, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		},
		{
			name:    "space evenly too much spacing",
			flex:    layout.FlexSpaceEvenly,
			spacing: 200,
			want: []layout.Rect{
				layout.NewRect(0, 0, 33, 1),
				layout.NewRect(33, 0, 34, 1),
				layout.NewRect(67, 0, 33, 1),
			},
		},
		{
			name:    "space around too much spacing",
			flex:    layout.FlexSpaceAround,
			spacing: 200,
			want: []layout.Rect{
				layout.NewRect(0, 0, 25, 1),
				layout.NewRect(25, 0, 50, 1),
				layout.NewRect(75, 0, 25, 1),
			},
		},
	}

	largeSpacingFlexes := []layout.Flex{
		layout.FlexStart,
		layout.FlexCenter,
		layout.FlexEnd,
		layout.FlexSpaceBetween,
	}
	for _, flex := range largeSpacingFlexes {
		tests = append(tests, struct {
			name    string
			flex    layout.Flex
			spacing int
			want    []layout.Rect
		}{
			name:    "packed too much spacing",
			flex:    flex,
			spacing: 200,
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(0, 0, 100, 1),
				layout.NewRect(100, 0, 0, 1),
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, spacers := layout.NewLayout(layout.Horizontal).
				Flex(tt.flex).
				Spacing(tt.spacing).
				Constraints(layout.Length(10), layout.Length(10)).
				SplitWithSpacers(layout.NewRect(0, 0, 100, 1))

			wantSegments := layout.NewLayout(layout.Horizontal).
				Flex(tt.flex).
				Spacing(tt.spacing).
				Constraints(layout.Length(10), layout.Length(10)).
				Split(layout.NewRect(0, 0, 100, 1))
			if !slices.Equal(segments, wantSegments) {
				t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", wantSegments, segments)
			}
			if len(spacers) != 3 {
				t.Fatalf("spacer count mismatch: want 3 got %d", len(spacers))
			}
			if !slices.Equal(spacers, tt.want) {
				t.Fatalf("spacers mismatch\nwant: %#v\n got: %#v", tt.want, spacers)
			}
		})
	}
}

func TestLayout_Split_fillOverlap(t *testing.T) {
	tests := []struct {
		name        string
		constraints []layout.Constraint
		flex        layout.Flex
		spacing     int
		want        []layout.Rect
	}{
		{
			name:        "fill space between overlap 10",
			constraints: []layout.Constraint{layout.Fill(1), layout.Fill(1)},
			flex:        layout.FlexSpaceBetween,
			spacing:     -10,
			want: []layout.Rect{
				layout.NewRect(0, 0, 55, 1),
				layout.NewRect(45, 0, 55, 1),
			},
		},
		{
			name:        "fill space evenly ignores overlap",
			constraints: []layout.Constraint{layout.Fill(1), layout.Fill(1)},
			flex:        layout.FlexSpaceEvenly,
			spacing:     -10,
			want: []layout.Rect{
				layout.NewRect(0, 0, 50, 1),
				layout.NewRect(50, 0, 50, 1),
			},
		},
		{
			name:        "fill length fill overlap 10",
			constraints: []layout.Constraint{layout.Fill(1), layout.Length(10), layout.Fill(1)},
			flex:        layout.FlexStart,
			spacing:     -10,
			want: []layout.Rect{
				layout.NewRect(0, 0, 55, 1),
				layout.NewRect(45, 0, 10, 1),
				layout.NewRect(45, 0, 55, 1),
			},
		},
		{
			name:        "fill length fill space around ignores overlap",
			constraints: []layout.Constraint{layout.Fill(1), layout.Length(10), layout.Fill(1)},
			flex:        layout.FlexSpaceAround,
			spacing:     -10,
			want: []layout.Rect{
				layout.NewRect(0, 0, 45, 1),
				layout.NewRect(45, 0, 10, 1),
				layout.NewRect(55, 0, 45, 1),
			},
		},
		{
			name:        "fill length fill overlap 1",
			constraints: []layout.Constraint{layout.Fill(1), layout.Length(10), layout.Fill(1)},
			flex:        layout.FlexSpaceBetween,
			spacing:     -1,
			want: []layout.Rect{
				layout.NewRect(0, 0, 46, 1),
				layout.NewRect(45, 0, 10, 1),
				layout.NewRect(54, 0, 46, 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := layout.NewLayout(layout.Horizontal).
				Flex(tt.flex).
				Spacing(tt.spacing).
				Constraints(tt.constraints...).
				Split(layout.NewRect(0, 0, 100, 1))

			if !slices.Equal(got, tt.want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestLayout_Split_flexSpacingLowerPriorityThanUserSpacing(t *testing.T) {
	got := layout.NewLayout(layout.Horizontal).
		Flex(layout.FlexCenter).
		Spacing(80).
		Constraints(layout.Length(10), layout.Length(10)).
		Split(layout.NewRect(0, 0, 100, 1))

	want := []layout.Rect{
		layout.NewRect(0, 0, 10, 1),
		layout.NewRect(90, 0, 10, 1),
	}
	if !slices.Equal(got, want) {
		t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestLayout_Split_percentageSpaceBetween(t *testing.T) {
	tests := []struct {
		name        string
		constraints []layout.Constraint
		want        string
	}{
		{name: "zero half", constraints: []layout.Constraint{layout.Percentage(0), layout.Percentage(50)}, want: "     bbbbb"},
		{name: "ten half", constraints: []layout.Constraint{layout.Percentage(10), layout.Percentage(50)}, want: "a    bbbbb"},
		{name: "quarter quarter", constraints: []layout.Constraint{layout.Percentage(25), layout.Percentage(25)}, want: "aaa     bb"},
		{name: "third full", constraints: []layout.Constraint{layout.Percentage(33), layout.Percentage(100)}, want: "aaabbbbbbb"},
		{name: "half half", constraints: []layout.Constraint{layout.Percentage(50), layout.Percentage(50)}, want: "aaaaabbbbb"},
		{name: "full half", constraints: []layout.Constraint{layout.Percentage(100), layout.Percentage(50)}, want: "aaaaabbbbb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderLayoutLetters(layout.NewLayout(layout.Horizontal).
				Flex(layout.FlexSpaceBetween).
				Constraints(tt.constraints...), 10)

			if got != tt.want {
				t.Fatalf("letters mismatch\nwant: %q\n got: %q", tt.want, got)
			}
		})
	}
}

func TestLayout_Split_ratioSpaceBetween(t *testing.T) {
	tests := []struct {
		name        string
		constraints []layout.Constraint
		want        string
	}{
		{name: "zero half", constraints: []layout.Constraint{layout.Ratio(0, 1), layout.Ratio(1, 2)}, want: "     bbbbb"},
		{name: "tenth half", constraints: []layout.Constraint{layout.Ratio(1, 10), layout.Ratio(1, 2)}, want: "a    bbbbb"},
		{name: "quarter quarter", constraints: []layout.Constraint{layout.Ratio(1, 4), layout.Ratio(1, 4)}, want: "aaa     bb"},
		{name: "third full", constraints: []layout.Constraint{layout.Ratio(1, 3), layout.Ratio(1, 1)}, want: "aaabbbbbbb"},
		{name: "half half", constraints: []layout.Constraint{layout.Ratio(1, 2), layout.Ratio(1, 2)}, want: "aaaaabbbbb"},
		{name: "full half", constraints: []layout.Constraint{layout.Ratio(1, 1), layout.Ratio(1, 2)}, want: "aaaaabbbbb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderLayoutLetters(layout.NewLayout(layout.Horizontal).
				Flex(layout.FlexSpaceBetween).
				Constraints(tt.constraints...), 10)

			if got != tt.want {
				t.Fatalf("letters mismatch\nwant: %q\n got: %q", tt.want, got)
			}
		})
	}
}

func renderLayoutLetters(l layout.Layout, width int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	cells := make([]rune, width)
	for i := range cells {
		cells[i] = ' '
	}

	for i, rect := range l.Split(layout.NewRect(0, 0, width, 1)) {
		letter := letters[i%len(letters)]
		start := max(0, rect.X)
		end := min(width, rect.X+rect.Width)
		for x := start; x < end; x++ {
			cells[x] = letter
		}
	}
	return string(cells)
}

func TestLayout_SplitWithSpacers_shouldApplyMarginsBeforeSplitting(t *testing.T) {
	segments, spacers := layout.NewHorizontalLayout(layout.Length(2), layout.Length(2)).
		Flex(layout.FlexSpaceBetween).
		Margin(1, 2).
		SplitWithSpacers(layout.NewRect(0, 0, 10, 6))

	wantSegments := []layout.Rect{
		layout.NewRect(1, 2, 2, 2),
		layout.NewRect(7, 2, 2, 2),
	}
	wantSpacers := []layout.Rect{
		layout.NewRect(1, 2, 0, 2),
		layout.NewRect(3, 2, 4, 2),
		layout.NewRect(9, 2, 0, 2),
	}

	if !slices.Equal(segments, wantSegments) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", wantSegments, segments)
	}
	if !slices.Equal(spacers, wantSpacers) {
		t.Fatalf("spacers mismatch\nwant: %#v\n got: %#v", wantSpacers, spacers)
	}
}

func TestLayout_Split_shouldHandleVerticalPercentageMin(t *testing.T) {
	area := layout.NewRect(2, 4, 8, 10)

	got := layout.NewLayout(layout.Vertical).
		Constraints(layout.Percentage(30), layout.Min(2), layout.Length(2)).
		Split(area)

	want := []layout.Rect{
		layout.NewRect(2, 4, 8, 3),
		layout.NewRect(2, 7, 8, 5),
		layout.NewRect(2, 12, 8, 2),
	}
	if !slices.Equal(got, want) {
		t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", want, got)
	}

	totalHeight := 0
	for i, rect := range got {
		totalHeight += rect.Height
		if rect.X != area.X || rect.Width != area.Width {
			t.Fatalf("rect %d cross-axis mismatch: %#v", i, rect)
		}
		if i > 0 && got[i-1].Bottom() != rect.Y {
			t.Fatalf("rect %d is not stacked after previous rect: prev=%#v current=%#v", i, got[i-1], rect)
		}
	}
	if totalHeight > area.Height {
		t.Fatalf("total height %d exceeds area height %d", totalHeight, area.Height)
	}
}

func TestLayout_Split_shouldHandleConstrainedWidthEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		constraints []layout.Constraint
		wantWidths  []int
	}{
		{
			name:        "multiple lengths leave remainder in last segment",
			width:       10,
			constraints: []layout.Constraint{layout.Length(2), layout.Length(3), layout.Length(1)},
			wantWidths:  []int{2, 3, 5},
		},
		{
			name:        "shortage preserves min before flexible constraints",
			width:       10,
			constraints: []layout.Constraint{layout.Length(8), layout.Min(4), layout.Percentage(50), layout.Ratio(1, 2)},
			wantWidths:  []int{6, 4, 0, 0},
		},
		{
			name:        "zero width returns zero width segments",
			width:       0,
			constraints: []layout.Constraint{layout.Length(2), layout.Min(1), layout.Percentage(50)},
			wantWidths:  []int{0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area := layout.NewRect(4, 6, tt.width, 2)
			got := layout.NewLayout(layout.Horizontal).
				Constraints(tt.constraints...).
				Split(area)

			if len(got) != len(tt.wantWidths) {
				t.Fatalf("rect count mismatch: want %d got %d", len(tt.wantWidths), len(got))
			}

			cursor := area.X
			for i, rect := range got {
				if rect.X != cursor || rect.Y != area.Y || rect.Width != tt.wantWidths[i] || rect.Height != area.Height {
					t.Fatalf("rect %d mismatch: want x=%d width=%d in area %#v, got %#v", i, cursor, tt.wantWidths[i], area, rect)
				}
				cursor += rect.Width
			}
			if cursor-area.X > area.Width {
				t.Fatalf("total width %d exceeds area width %d", cursor-area.X, area.Width)
			}
		})
	}
}
