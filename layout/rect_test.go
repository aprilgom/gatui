package layout_test

import (
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
