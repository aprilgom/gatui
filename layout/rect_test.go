package layout_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"gatui/buffer"
	"gatui/layout"
)

func TestPosition_NewAndOffset_shouldMatchRatatui(t *testing.T) {
	tests := []struct {
		name     string
		position layout.Position
		offset   layout.Offset
		want     layout.Position
	}{
		{
			name:     "positive",
			position: layout.NewPosition(2, 3),
			offset:   layout.NewOffset(5, 7),
			want:     layout.NewPosition(7, 10),
		},
		{
			name:     "negative",
			position: layout.NewPosition(10, 10),
			offset:   layout.NewOffset(-3, -4),
			want:     layout.NewPosition(7, 6),
		},
		{
			name:     "negative clamps to origin",
			position: layout.NewPosition(1, 1),
			offset:   layout.NewOffset(-5, -6),
			want:     layout.NewPosition(0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.position.Offset(tt.offset)
			if got != tt.want {
				t.Fatalf("position mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestOffset_FromPosition_convertsCoordinates(t *testing.T) {
	got := layout.OffsetFromPosition(layout.NewPosition(4, 9))
	want := layout.NewOffset(4, 9)

	if got != want {
		t.Fatalf("OffsetFromPosition() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestPosition_AddAndSubtractOffset(t *testing.T) {
	got := layout.NewPosition(10, 10).
		AddOffset(layout.NewOffset(-3, 4)).
		SubOffset(layout.NewOffset(5, 20))
	want := layout.NewPosition(2, 0)

	if got != want {
		t.Fatalf("position mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestPosition_AddAssignAndSubAssignOffset_goPattern(t *testing.T) {
	position := layout.NewPosition(10, 10)

	position = position.AddOffset(layout.NewOffset(-3, 4))
	position = position.SubOffset(layout.NewOffset(5, 20))

	if want := layout.NewPosition(2, 0); position != want {
		t.Fatalf("position mismatch\nwant: %#v\n got: %#v", want, position)
	}
}

func TestRect_FromPositionAndSize(t *testing.T) {
	got := layout.RectFromPositionAndSize(layout.NewPosition(1, 2), layout.NewSize(3, 4))
	want := layout.NewRect(1, 2, 3, 4)

	if got != want {
		t.Fatalf("RectFromPositionAndSize() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_FromSize(t *testing.T) {
	got := layout.RectFromSize(layout.NewSize(3, 4))
	want := layout.NewRect(0, 0, 3, 4)

	if got != want {
		t.Fatalf("RectFromSize() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_New_sizeTruncation(t *testing.T) {
	got := layout.NewRect(layout.MaxCoordinate-100, layout.MaxCoordinate-1000, 200, 2000)
	want := layout.Rect{X: layout.MaxCoordinate - 100, Y: layout.MaxCoordinate - 1000, Width: 100, Height: 1000}

	if got != want {
		t.Fatalf("NewRect() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_New_sizePreservation(t *testing.T) {
	got := layout.NewRect(10, 20, 200, 300)
	want := layout.Rect{X: 10, Y: 20, Width: 200, Height: 300}

	if got != want {
		t.Fatalf("NewRect() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_Offset_saturatesAtMaxCoordinate(t *testing.T) {
	tests := []struct {
		name   string
		rect   layout.Rect
		offset layout.Offset
		want   layout.Rect
	}{
		{
			name:   "offset",
			rect:   layout.NewRect(layout.MaxCoordinate-10, layout.MaxCoordinate-20, 10, 20),
			offset: layout.NewOffset(100, 100),
			want:   layout.Rect{X: layout.MaxCoordinate - 10, Y: layout.MaxCoordinate - 20, Width: 10, Height: 20},
		},
		{
			name:   "sub offset with negative values",
			rect:   layout.NewRect(layout.MaxCoordinate-10, layout.MaxCoordinate-20, 10, 20),
			offset: layout.NewOffset(-100, -100),
			want:   layout.Rect{X: layout.MaxCoordinate - 10, Y: layout.MaxCoordinate - 20, Width: 10, Height: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got layout.Rect
			if tt.name == "sub offset with negative values" {
				got = tt.rect.SubOffset(tt.offset)
			} else {
				got = tt.rect.Offset(tt.offset)
			}
			if got != tt.want {
				t.Fatalf("rect mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestRect_Resize_clampsAtBounds(t *testing.T) {
	got := layout.NewRect(layout.MaxCoordinate-100, layout.MaxCoordinate-1000, 50, 60).
		Resize(layout.NewSize(200, 2000))
	want := layout.Rect{X: layout.MaxCoordinate - 100, Y: layout.MaxCoordinate - 1000, Width: 100, Height: 1000}

	if got != want {
		t.Fatalf("Resize() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_Layout_returnsSplitRects(t *testing.T) {
	rect := layout.NewRect(0, 0, 10, 5)
	split := layout.NewHorizontalLayout(layout.Length(3), layout.Fill(1))
	want := []layout.Rect{
		layout.NewRect(0, 0, 3, 5),
		layout.NewRect(3, 0, 7, 5),
	}

	got := rect.Layout(split)

	if !slices.Equal(got, want) {
		t.Fatalf("Rect.Layout(%#v) mismatch\nwant: %#v\n got: %#v", split, want, got)
	}
}

func TestRect_SplitN_invalidNumberOfRectsPanics(t *testing.T) {
	rect := layout.NewRect(0, 0, 10, 5)
	split := layout.NewHorizontalLayout(layout.Length(3), layout.Fill(1))
	wantPanic := "invalid number of rects: expected 3, found 2"

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("Rect.SplitN(%#v, 3) did not panic", split)
		}
		if gotString := fmt.Sprint(got); gotString != wantPanic {
			t.Fatalf("Rect.SplitN(%#v, 3) panic = %q, want %q", split, gotString, wantPanic)
		}
	}()

	_ = rect.SplitN(split, 3)
}

func TestRect_TrySplitN_invalidNumberOfRectsReturnsError(t *testing.T) {
	rect := layout.NewRect(0, 0, 10, 5)
	split := layout.NewHorizontalLayout(layout.Length(3), layout.Fill(1))
	wantErr := "invalid number of rects: expected 3, found 2"

	got, err := rect.TrySplitN(split, 3)

	if err == nil {
		t.Fatalf("Rect.TrySplitN(%#v, 3) error = nil, want %q", split, wantErr)
	}
	if err.Error() != wantErr {
		t.Fatalf("Rect.TrySplitN(%#v, 3) error = %q, want %q", split, err, wantErr)
	}
	if got != nil {
		t.Fatalf("Rect.TrySplitN(%#v, 3) rects = %#v, want nil", split, got)
	}
}

func TestSize_Tuple(t *testing.T) {
	width, height := layout.NewSize(10, 20).Tuple()

	if width != 10 || height != 20 {
		t.Fatalf("Tuple() = (%d, %d), want (10, 20)", width, height)
	}
}

func TestRect_Rows_shouldIterateTopToBottom(t *testing.T) {
	tests := []struct {
		name string
		rect layout.Rect
		want []layout.Rect
	}{
		{
			name: "normal",
			rect: layout.NewRect(0, 0, 2, 3),
			want: []layout.Rect{
				layout.NewRect(0, 0, 2, 1),
				layout.NewRect(0, 1, 2, 1),
				layout.NewRect(0, 2, 2, 1),
			},
		},
		{
			name: "zero height",
			rect: layout.NewRect(0, 0, 2, 0),
			want: []layout.Rect{},
		},
		{
			name: "zero width",
			rect: layout.NewRect(0, 0, 0, 3),
			want: []layout.Rect{
				layout.NewRect(0, 0, 0, 1),
				layout.NewRect(0, 1, 0, 1),
				layout.NewRect(0, 2, 0, 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rect.Rows()
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Rows() mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestRect_RowsBack(t *testing.T) {
	got := layout.NewRect(0, 0, 2, 3).RowsReversed()
	want := []layout.Rect{
		layout.NewRect(0, 2, 2, 1),
		layout.NewRect(0, 1, 2, 1),
		layout.NewRect(0, 0, 2, 1),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("RowsReversed() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_RowsMeetInTheMiddle(t *testing.T) {
	rows := layout.NewRect(0, 0, 2, 4).Rows()

	got := consumeSliceEnds(rows, true, false, true, false)
	want := []layout.Rect{
		layout.NewRect(0, 0, 2, 1),
		layout.NewRect(0, 3, 2, 1),
		layout.NewRect(0, 1, 2, 1),
		layout.NewRect(0, 2, 2, 1),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("mixed row traversal mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_Columns_shouldIterateLeftToRight(t *testing.T) {
	tests := []struct {
		name string
		rect layout.Rect
		want []layout.Rect
	}{
		{
			name: "normal",
			rect: layout.NewRect(0, 0, 3, 2),
			want: []layout.Rect{
				layout.NewRect(0, 0, 1, 2),
				layout.NewRect(1, 0, 1, 2),
				layout.NewRect(2, 0, 1, 2),
			},
		},
		{
			name: "zero width",
			rect: layout.NewRect(0, 0, 0, 2),
			want: []layout.Rect{},
		},
		{
			name: "zero height",
			rect: layout.NewRect(0, 0, 3, 0),
			want: []layout.Rect{
				layout.NewRect(0, 0, 1, 0),
				layout.NewRect(1, 0, 1, 0),
				layout.NewRect(2, 0, 1, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rect.Columns()
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Columns() mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestRect_ColumnsBack(t *testing.T) {
	got := layout.NewRect(0, 0, 3, 2).ColumnsReversed()
	want := []layout.Rect{
		layout.NewRect(2, 0, 1, 2),
		layout.NewRect(1, 0, 1, 2),
		layout.NewRect(0, 0, 1, 2),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("ColumnsReversed() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_ColumnsMeetInTheMiddle(t *testing.T) {
	columns := layout.NewRect(0, 0, 4, 2).Columns()

	got := consumeSliceEnds(columns, true, false, true, false)
	want := []layout.Rect{
		layout.NewRect(0, 0, 1, 2),
		layout.NewRect(3, 0, 1, 2),
		layout.NewRect(1, 0, 1, 2),
		layout.NewRect(2, 0, 1, 2),
	}

	if !slices.Equal(got, want) {
		t.Fatalf("mixed column traversal mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_ColumnsMax(t *testing.T) {
	columns := layout.NewRect(0, 0, layout.MaxCoordinate, 1).Columns()

	got := columns[len(columns)-1]
	want := layout.NewRect(layout.MaxCoordinate-1, 0, 1, 1)

	if got != want {
		t.Fatalf("last column mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_ColumnsMin(t *testing.T) {
	columns := layout.NewRect(0, 0, layout.MaxCoordinate, 1).ColumnsReversed()

	got := columns[len(columns)-1]
	want := layout.NewRect(0, 0, 1, 1)

	if got != want {
		t.Fatalf("last reversed column mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_Positions_shouldIterateRowMajor(t *testing.T) {
	tests := []struct {
		name string
		rect layout.Rect
		want []layout.Position
	}{
		{
			name: "origin",
			rect: layout.NewRect(0, 0, 2, 2),
			want: []layout.Position{
				layout.NewPosition(0, 0),
				layout.NewPosition(1, 0),
				layout.NewPosition(0, 1),
				layout.NewPosition(1, 1),
			},
		},
		{
			name: "offset",
			rect: layout.NewRect(2, 3, 2, 2),
			want: []layout.Position{
				layout.NewPosition(2, 3),
				layout.NewPosition(3, 3),
				layout.NewPosition(2, 4),
				layout.NewPosition(3, 4),
			},
		},
		{
			name: "zero width",
			rect: layout.NewRect(0, 0, 0, 2),
			want: []layout.Position{},
		},
		{
			name: "zero height",
			rect: layout.NewRect(0, 0, 2, 0),
			want: []layout.Position{},
		},
		{
			name: "zero by zero",
			rect: layout.NewRect(0, 0, 0, 0),
			want: []layout.Position{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rect.Positions()
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Positions() mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func consumeSliceEnds[T any](items []T, fromFront ...bool) []T {
	front := 0
	back := len(items) - 1
	consumed := make([]T, 0, len(fromFront))

	for _, takeFront := range fromFront {
		if front > back {
			break
		}
		if takeFront {
			consumed = append(consumed, items[front])
			front++
			continue
		}
		consumed = append(consumed, items[back])
		back--
	}

	return consumed
}

func TestSize_NewAndArea_shouldMatchRatatui(t *testing.T) {
	tests := []struct {
		size layout.Size
		want int
	}{
		{layout.NewSize(10, 20), 200},
		{layout.NewSize(0, 20), 0},
		{layout.NewSize(10, 0), 0},
	}

	for _, tt := range tests {
		if got := tt.size.Area(); got != tt.want {
			t.Fatalf("Area() = %d, want %d for %#v", got, tt.want, tt.size)
		}
	}
}

func TestRect_GeometryHelpers_shouldMatchRatatui(t *testing.T) {
	rect := layout.NewRect(1, 2, 3, 4)

	if got, want := rect.Area(), 12; got != want {
		t.Fatalf("Area() = %d, want %d", got, want)
	}
	if rect.IsEmpty() {
		t.Fatalf("IsEmpty() = true, want false")
	}
	for _, empty := range []layout.Rect{
		layout.NewRect(1, 2, 0, 4),
		layout.NewRect(1, 2, 3, 0),
	} {
		if !empty.IsEmpty() {
			t.Fatalf("IsEmpty() = false, want true for %#v", empty)
		}
	}
	if got, want := rect.Resize(layout.NewSize(30, 40)), layout.NewRect(1, 2, 30, 40); got != want {
		t.Fatalf("Resize() mismatch\nwant: %#v\n got: %#v", want, got)
	}
	if got, want := rect.AsPosition(), layout.NewPosition(1, 2); got != want {
		t.Fatalf("AsPosition() mismatch\nwant: %#v\n got: %#v", want, got)
	}
	if got, want := rect.AsSize(), layout.NewSize(3, 4); got != want {
		t.Fatalf("AsSize() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_UnionIntersectsContains_shouldUseExclusiveRightAndBottomEdges(t *testing.T) {
	if got, want := layout.NewRect(1, 2, 3, 4).Union(layout.NewRect(2, 3, 4, 5)), layout.NewRect(1, 2, 5, 6); got != want {
		t.Fatalf("Union() mismatch\nwant: %#v\n got: %#v", want, got)
	}

	intersections := []struct {
		name string
		a    layout.Rect
		b    layout.Rect
		want bool
	}{
		{"overlap", layout.NewRect(1, 2, 3, 4), layout.NewRect(2, 3, 4, 5), true},
		{"corner touch", layout.NewRect(0, 0, 10, 10), layout.NewRect(10, 10, 20, 20), false},
		{"edge touch", layout.NewRect(0, 0, 10, 10), layout.NewRect(10, 0, 20, 10), false},
		{"contained", layout.NewRect(0, 0, 20, 20), layout.NewRect(5, 5, 10, 10), true},
	}
	for _, tt := range intersections {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Intersects(tt.b); got != tt.want {
				t.Fatalf("Intersects() = %v, want %v", got, tt.want)
			}
			if got := tt.b.Intersects(tt.a); got != tt.want {
				t.Fatalf("Intersects() reverse = %v, want %v", got, tt.want)
			}
		})
	}

	contains := []struct {
		name     string
		position layout.Position
		want     bool
	}{
		{"top left", layout.NewPosition(1, 2), true},
		{"bottom right inside", layout.NewPosition(3, 5), true},
		{"right edge", layout.NewPosition(4, 2), false},
		{"bottom edge", layout.NewPosition(1, 6), false},
	}
	rect := layout.NewRect(1, 2, 3, 4)
	for _, tt := range contains {
		t.Run(tt.name, func(t *testing.T) {
			if got := rect.Contains(tt.position); got != tt.want {
				t.Fatalf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_CenteredHelpers_shouldMatchRatatuiLayoutIntent(t *testing.T) {
	tests := []struct {
		name string
		got  layout.Rect
		want layout.Rect
	}{
		{
			name: "centered horizontally",
			got:  layout.NewRect(10, 20, 100, 50).CenteredHorizontally(layout.Length(20)),
			want: layout.NewRect(50, 20, 20, 50),
		},
		{
			name: "centered vertically",
			got:  layout.NewRect(10, 20, 100, 50).CenteredVertically(layout.Percentage(50)),
			want: layout.NewRect(10, 32, 100, 25),
		},
		{
			name: "centered both",
			got:  layout.NewRect(10, 20, 100, 50).Centered(layout.Ratio(1, 2), layout.Ratio(1, 5)),
			want: layout.NewRect(35, 40, 50, 10),
		},
		{
			name: "oversized constraint clamps to area",
			got:  layout.NewRect(10, 20, 100, 50).Centered(layout.Length(200), layout.Length(80)),
			want: layout.NewRect(10, 20, 100, 50),
		},
		{
			name: "min constraint fills single centered area like layout",
			got:  layout.NewRect(10, 20, 100, 50).CenteredHorizontally(layout.Min(20)),
			want: layout.NewRect(10, 20, 100, 50),
		},
		{
			name: "fill constraint fills single centered area like layout",
			got:  layout.NewRect(10, 20, 100, 50).CenteredVertically(layout.Fill(1)),
			want: layout.NewRect(10, 20, 100, 50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("rect mismatch\nwant: %#v\n got: %#v", tt.want, tt.got)
			}
		})
	}
}

func TestRect_Inner_shouldShrinkByMargin(t *testing.T) {
	base := layout.NewRect(2, 2, 10, 6)
	inner := base.Inner(layout.NewMargin(2, 1))

	buf := buffer.Empty(layout.NewRect(0, 0, 15, 10))
	fill(buf, base, "█")
	fill(buf, inner, "░")

	assertBufferLines(t, buf, []string{
		"               ",
		"               ",
		"  ██████████   ",
		"  ██░░░░░░██   ",
		"  ██░░░░░░██   ",
		"  ██░░░░░░██   ",
		"  ██░░░░░░██   ",
		"  ██████████   ",
		"               ",
		"               ",
	})
}

func TestRect_Outer_shouldExpandByMargin(t *testing.T) {
	base := layout.NewRect(4, 3, 6, 4)
	outer := base.Outer(layout.NewMargin(2, 1))

	buf := buffer.Empty(layout.NewRect(0, 0, 15, 10))
	fill(buf, outer, "░")
	fill(buf, base, "█")

	assertBufferLines(t, buf, []string{
		"               ",
		"               ",
		"  ░░░░░░░░░░   ",
		"  ░░██████░░   ",
		"  ░░██████░░   ",
		"  ░░██████░░   ",
		"  ░░██████░░   ",
		"  ░░░░░░░░░░   ",
		"               ",
		"               ",
	})
}

func TestRect_Offset_shouldMoveWithoutResizing(t *testing.T) {
	base := layout.NewRect(2, 2, 5, 3)
	moved := base.Offset(layout.NewOffset(4, 2))

	buf := buffer.Empty(layout.NewRect(0, 0, 15, 10))
	fill(buf, base, "░")
	fill(buf, moved, "█")

	assertBufferLines(t, buf, []string{
		"               ",
		"               ",
		"  ░░░░░        ",
		"  ░░░░░        ",
		"  ░░░░█████    ",
		"      █████    ",
		"      █████    ",
		"               ",
		"               ",
		"               ",
	})
}

func TestRect_Offset_negative(t *testing.T) {
	got := layout.NewRect(10, 10, 3, 4).Offset(layout.NewOffset(-4, -7))
	want := layout.NewRect(6, 3, 3, 4)

	if got != want {
		t.Fatalf("Offset() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_Offset_negativeSaturatesAtOrigin(t *testing.T) {
	got := layout.NewRect(2, 3, 4, 5).Offset(layout.NewOffset(-10, -20))
	want := layout.NewRect(0, 0, 4, 5)

	if got != want {
		t.Fatalf("Offset() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_SubOffset(t *testing.T) {
	tests := []struct {
		name   string
		rect   layout.Rect
		offset layout.Offset
		want   layout.Rect
	}{
		{
			name:   "positive offset moves left and up",
			rect:   layout.NewRect(10, 10, 3, 4),
			offset: layout.NewOffset(4, 7),
			want:   layout.NewRect(6, 3, 3, 4),
		},
		{
			name:   "negative offset moves right and down",
			rect:   layout.NewRect(10, 10, 3, 4),
			offset: layout.NewOffset(-4, -7),
			want:   layout.NewRect(14, 17, 3, 4),
		},
		{
			name:   "saturates at origin",
			rect:   layout.NewRect(2, 3, 4, 5),
			offset: layout.NewOffset(10, 20),
			want:   layout.NewRect(0, 0, 4, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rect.SubOffset(tt.offset)
			if got != tt.want {
				t.Fatalf("SubOffset() mismatch\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestRect_AddAssignOffset_goPattern(t *testing.T) {
	rect := layout.NewRect(10, 10, 3, 4)

	rect = rect.Offset(layout.NewOffset(-4, 7))

	if want := layout.NewRect(6, 17, 3, 4); rect != want {
		t.Fatalf("rect mismatch\nwant: %#v\n got: %#v", want, rect)
	}
}

func TestRect_SubAssignOffset_goPattern(t *testing.T) {
	rect := layout.NewRect(10, 10, 3, 4)

	rect = rect.SubOffset(layout.NewOffset(4, -7))

	if want := layout.NewRect(6, 17, 3, 4); rect != want {
		t.Fatalf("rect mismatch\nwant: %#v\n got: %#v", want, rect)
	}
}

func TestRect_Intersection_shouldReturnOverlappingArea(t *testing.T) {
	a := layout.NewRect(2, 2, 6, 4)
	b := layout.NewRect(5, 3, 6, 4)
	intersection := a.Intersection(b)

	buf := buffer.Empty(layout.NewRect(0, 0, 15, 10))
	fill(buf, a, "░")
	fill(buf, b, "▒")
	fill(buf, intersection, "█")

	assertBufferLines(t, buf, []string{
		"               ",
		"               ",
		"  ░░░░░░       ",
		"  ░░░███▒▒▒    ",
		"  ░░░███▒▒▒    ",
		"  ░░░███▒▒▒    ",
		"     ▒▒▒▒▒▒    ",
		"               ",
		"               ",
		"               ",
	})
}

func TestRect_Intersection_underflow(t *testing.T) {
	got := layout.NewRect(1, 1, 2, 2).Intersection(layout.NewRect(4, 4, 2, 2))
	want := layout.NewRect(4, 4, 0, 0)

	if got != want {
		t.Fatalf("Intersection() mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRect_Clamp_shouldMoveRectInsideArea(t *testing.T) {
	area := layout.NewRect(2, 2, 10, 6)
	rect := layout.NewRect(8, 5, 8, 4)
	clamped := rect.Clamp(area)

	buf := buffer.Empty(layout.NewRect(0, 0, 20, 12))
	fill(buf, area, "█")
	fill(buf, rect, "▒")
	fill(buf, clamped, "░")

	assertBufferLines(t, buf, []string{
		"                    ",
		"                    ",
		"  ██████████        ",
		"  ██████████        ",
		"  ██░░░░░░░░        ",
		"  ██░░░░░░░░▒▒▒▒    ",
		"  ██░░░░░░░░▒▒▒▒    ",
		"  ██░░░░░░░░▒▒▒▒    ",
		"        ▒▒▒▒▒▒▒▒    ",
		"                    ",
		"                    ",
		"                    ",
	})
}

func fill(buf *buffer.Buffer, area layout.Rect, symbol string) {
	for y := area.Top(); y < area.Bottom(); y++ {
		for x := area.Left(); x < area.Right(); x++ {
			buf.SetSymbol(x, y, symbol)
		}
	}
}

func assertBufferLines(t *testing.T, buf *buffer.Buffer, expected []string) {
	t.Helper()

	actual := strings.Join(buf.Lines(), "\n")
	want := strings.Join(expected, "\n")
	if actual != want {
		t.Fatalf("buffer mismatch\nwant:\n%s\n\ngot:\n%s", want, actual)
	}
}
