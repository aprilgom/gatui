package layout

func axisLength(area Rect, direction Direction) int {
	if direction == Vertical {
		return area.Height
	}
	return area.Width
}

func spacerRect(area Rect, direction Direction, start int, length int) Rect {
	if direction == Vertical {
		return Rect{X: area.X, Y: area.Y + start, Width: area.Width, Height: length}
	}
	return Rect{X: area.X + start, Y: area.Y, Width: length, Height: area.Height}
}

func emptySpacer(area Rect, direction Direction, start int) Rect {
	return spacerRect(area, direction, start, 0)
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

func clampInt(value, low, high int) int {
	if high < low {
		return low
	}
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
