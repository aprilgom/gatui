package layout_test

import (
	"reflect"
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

			if !reflect.DeepEqual(got, tt.want) {
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

			if !reflect.DeepEqual(got, tt.want) {
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
			if !reflect.DeepEqual(got, want) {
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

			if !reflect.DeepEqual(got, tt.want) {
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
	want := []layout.Rect{
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

				if !reflect.DeepEqual(got, want) {
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
	if !reflect.DeepEqual(got, want) {
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("rects mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
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
	if !reflect.DeepEqual(got, want) {
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
