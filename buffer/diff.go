package buffer

const variationSelector16 = '\ufe0f'

type CellDiff struct {
	X    int
	Y    int
	Cell Cell
}

func (b *Buffer) Diff(next *Buffer) []CellDiff {
	if b.Area.X != next.Area.X || b.Area.Y != next.Area.Y || b.Area.Width != next.Area.Width {
		panic("buffer areas must have the same x, y, and width")
	}

	height := b.Area.Height
	if next.Area.Height < height {
		height = next.Area.Height
	}
	diffs := make([]CellDiff, 0)
	for y := 0; y < height; y++ {
		for x := 0; x < b.Area.Width; x++ {
			index := y*b.Area.Width + x
			previous := b.Cells[index]
			current := next.Cells[index]
			width := current.Width()

			if current.DiffOption == CellDiffSkip {
				continue
			}
			if current.DiffOption == CellDiffAlwaysUpdate || current != previous {
				diffs = append(diffs, CellDiff{
					X:    next.Area.X + x,
					Y:    next.Area.Y + y,
					Cell: current,
				})
				if width > 1 && containsRune(current.Symbol, variationSelector16) {
					trailingEnd := minInt(x+width, b.Area.Width)
					for trailingX := x + 1; trailingX < trailingEnd; trailingX++ {
						trailingIndex := y*b.Area.Width + trailingX
						trailingPrevious := b.Cells[trailingIndex]
						trailingCurrent := next.Cells[trailingIndex]
						if trailingCurrent.DiffOption != CellDiffSkip && trailingPrevious.Symbol != trailingCurrent.Symbol {
							diffs = append(diffs, CellDiff{
								X:    next.Area.X + trailingX,
								Y:    next.Area.Y + y,
								Cell: trailingCurrent,
							})
						}
					}
				}
			}
			if current.DiffOption == CellDiffForcedWidth || current.ForcedWidth > 0 || width > 1 {
				x += width - 1
			}
		}
	}
	return diffs
}

func containsRune(value string, target rune) bool {
	for _, r := range value {
		if r == target {
			return true
		}
	}
	return false
}
