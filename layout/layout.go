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

func DefaultLayout() Layout {
	return NewLayout(Vertical)
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
