package layout

import "fmt"

type Layout struct {
	direction   Direction
	constraints []Constraint
	flex        Flex
	spacing     int
	margin      Margin
}

func NewLayout(direction Direction) Layout {
	return Layout{direction: direction}
}

func NewVerticalLayout(constraints ...Constraint) Layout {
	return NewLayout(Vertical).Constraints(constraints...)
}

func NewHorizontalLayout(constraints ...Constraint) Layout {
	return NewLayout(Horizontal).Constraints(constraints...)
}

func (l Layout) Constraints(constraints ...Constraint) Layout {
	l.constraints = append([]Constraint(nil), constraints...)
	return l
}

func (l Layout) Direction(direction Direction) Layout {
	l.direction = direction
	return l
}

func (l Layout) Flex(flex Flex) Layout {
	l.flex = flex
	return l
}

func (l Layout) Spacing(spacing int) Layout {
	l.spacing = spacing
	return l
}

func (l Layout) Margin(horizontal, vertical int) Layout {
	l.margin = NewMargin(horizontal, vertical)
	return l
}

func (l Layout) UniformMargin(margin int) Layout {
	return l.Margin(margin, margin)
}

func (l Layout) HorizontalMargin(horizontal int) Layout {
	l.margin.Horizontal = horizontal
	return l
}

func (l Layout) VerticalMargin(vertical int) Layout {
	l.margin.Vertical = vertical
	return l
}

func (l Layout) Split(area Rect) []Rect {
	rects, _ := l.SplitWithSpacers(area)
	return rects
}

func (l Layout) SplitN(area Rect, count int) []Rect {
	rects, err := l.TrySplitN(area, count)
	if err != nil {
		panic(err)
	}
	return rects
}

func (l Layout) TrySplitN(area Rect, count int) ([]Rect, error) {
	rects := l.Split(area)
	if len(rects) != count {
		return nil, fmt.Errorf("invalid number of rects: expected %d, found %d", count, len(rects))
	}
	return rects, nil
}

func (l Layout) SplitWithSpacers(area Rect) ([]Rect, []Rect) {
	area = area.Inner(l.margin)
	if len(l.constraints) == 0 {
		return []Rect{area}, []Rect{emptySpacer(area, l.direction, 0), emptySpacer(area, l.direction, axisLength(area, l.direction))}
	}

	solved := l.solveLayout(area)
	return solved.segments, solved.spacers
}

func (l Layout) splitSegments(area Rect) ([]Rect, []int, []int) {
	axisLength := area.Width
	if l.direction == Vertical {
		axisLength = area.Height
	}
	lengths := calculateLengths(maxInt(0, axisLength-l.spacingAllowance()), l.constraints, false)
	if l.flex == FlexSpaceBetween && len(lengths) == 1 {
		lengths[0] = axisLength
	}
	if l.flex == FlexLegacy && len(lengths) > 0 {
		occupied := spacedLength(lengths, l.spacing)
		if occupied < axisLength {
			lengths[len(lengths)-1] += axisLength - occupied
		}
	}
	offsets := flexOffsets(axisLength, lengths, l.flex, l.spacing)

	rects := make([]Rect, 0, len(l.constraints))

	for i, length := range lengths {
		width := area.Width
		height := area.Height
		x := area.X
		y := area.Y
		if l.direction == Horizontal {
			width = length
			x += offsets[i]
		} else {
			height = length
			y += offsets[i]
		}

		rect := Rect{X: x, Y: y, Width: width, Height: height}
		rects = append(rects, rect)
	}

	return rects, offsets, lengths
}

func (l Layout) spacingAllowance() int {
	if l.spacing == 0 || len(l.constraints) <= 1 {
		return 0
	}
	if l.spacing < 0 && (l.flex == FlexSpaceAround || l.flex == FlexSpaceEvenly) {
		return 0
	}
	return l.spacing * (len(l.constraints) - 1)
}

func (l Layout) spacerRects(area Rect, offsets []int, lengths []int) []Rect {
	spacers := make([]Rect, 0, len(lengths)+1)
	areaLength := axisLength(area, l.direction)
	previousEnd := 0

	for i, offset := range offsets {
		start := clampInt(offset, 0, areaLength)
		spacers = append(spacers, spacerRect(area, l.direction, previousEnd, maxInt(0, start-previousEnd)))

		end := clampInt(offset+lengths[i], 0, areaLength)
		if end > previousEnd {
			previousEnd = end
		}
	}

	spacers = append(spacers, spacerRect(area, l.direction, previousEnd, maxInt(0, areaLength-previousEnd)))
	return spacers
}
