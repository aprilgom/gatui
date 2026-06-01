package buffer

import (
	"unicode"

	"gatui/layout"
)

func unionRect(a, b layout.Rect) layout.Rect {
	if a.Width == 0 || a.Height == 0 {
		return b
	}
	if b.Width == 0 || b.Height == 0 {
		return a
	}
	x1 := minInt(a.X, b.X)
	y1 := minInt(a.Y, b.Y)
	x2 := maxInt(a.X+a.Width, b.X+b.Width)
	y2 := maxInt(a.Y+a.Height, b.Y+b.Height)
	return layout.NewRect(x1, y1, x2-x1, y2-y1)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func containsControl(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
