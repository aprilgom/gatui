package layout_test

import (
	"strings"
	"testing"

	"gatui/buffer"
	"gatui/layout"
)

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
