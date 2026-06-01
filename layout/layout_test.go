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
