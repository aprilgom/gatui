package layout

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func NewRect(x, y, width, height int) Rect {
	width = maxInt(0, width)
	height = maxInt(0, height)
	return Rect{X: x, Y: y, Width: width, Height: height}
}

func (r Rect) Left() int {
	return r.X
}

func (r Rect) Right() int {
	return r.X + r.Width
}

func (r Rect) Top() int {
	return r.Y
}

func (r Rect) Bottom() int {
	return r.Y + r.Height
}

func (r Rect) Inner(margin Margin) Rect {
	width := r.Width - margin.Horizontal*2
	height := r.Height - margin.Vertical*2
	if width < 0 || height < 0 {
		return Rect{}
	}

	return NewRect(r.X+margin.Horizontal, r.Y+margin.Vertical, width, height)
}

func (r Rect) Outer(margin Margin) Rect {
	x := r.X - margin.Horizontal
	y := r.Y - margin.Vertical
	return NewRect(x, y, r.Right()+margin.Horizontal-x, r.Bottom()+margin.Vertical-y)
}

func (r Rect) Offset(offset Offset) Rect {
	return NewRect(r.X+offset.X, r.Y+offset.Y, r.Width, r.Height)
}

func (r Rect) Intersection(other Rect) Rect {
	x1 := maxInt(r.X, other.X)
	y1 := maxInt(r.Y, other.Y)
	x2 := minInt(r.Right(), other.Right())
	y2 := minInt(r.Bottom(), other.Bottom())
	return NewRect(x1, y1, maxInt(0, x2-x1), maxInt(0, y2-y1))
}

func (r Rect) Clamp(other Rect) Rect {
	width := minInt(r.Width, other.Width)
	height := minInt(r.Height, other.Height)
	x := clampInt(r.X, other.X, other.Right()-width)
	y := clampInt(r.Y, other.Y, other.Bottom()-height)
	return NewRect(x, y, width, height)
}

type Position struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}

type Margin struct {
	Horizontal int
	Vertical   int
}

func NewMargin(horizontal, vertical int) Margin {
	return Margin{Horizontal: horizontal, Vertical: vertical}
}

type Offset struct {
	X int
	Y int
}

func NewOffset(x, y int) Offset {
	return Offset{X: x, Y: y}
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
)

type Constraint struct {
	kind        constraintKind
	value       int
	denominator int
}

type constraintKind int

const (
	constraintLength constraintKind = iota
	constraintMin
	constraintPercentage
	constraintRatio
)

func Length(value int) Constraint {
	return Constraint{kind: constraintLength, value: value}
}

func Min(value int) Constraint {
	return Constraint{kind: constraintMin, value: value}
}

func Percentage(percent int) Constraint {
	return Constraint{kind: constraintPercentage, value: percent}
}

func Ratio(numerator, denominator int) Constraint {
	return Constraint{kind: constraintRatio, value: numerator, denominator: denominator}
}

func (c Constraint) IsLength() bool {
	return c.kind == constraintLength
}

func (c Constraint) IsPercentage() bool {
	return c.kind == constraintPercentage
}

func (c Constraint) IsRatio() bool {
	return c.kind == constraintRatio
}

func (c Constraint) Value() int {
	return c.value
}

func (c Constraint) Denominator() int {
	return c.denominator
}

type Layout struct {
	direction   Direction
	constraints []Constraint
}

func NewLayout(direction Direction) Layout {
	return Layout{direction: direction}
}

func (l Layout) Constraints(constraints ...Constraint) Layout {
	l.constraints = append([]Constraint(nil), constraints...)
	return l
}

func (l Layout) Split(area Rect) []Rect {
	if len(l.constraints) == 0 {
		return []Rect{area}
	}

	axisLength := area.Width
	if l.direction == Vertical {
		axisLength = area.Height
	}
	lengths := calculateLengths(axisLength, l.constraints)

	rects := make([]Rect, 0, len(l.constraints))
	cursorX := area.X
	cursorY := area.Y

	for _, length := range lengths {
		width := area.Width
		height := area.Height
		if l.direction == Horizontal {
			width = length
		} else {
			height = length
		}

		rect := Rect{X: cursorX, Y: cursorY, Width: width, Height: height}
		rects = append(rects, rect)

		if l.direction == Horizontal {
			cursorX += length
		} else {
			cursorY += length
		}
	}

	return rects
}

func calculateLengths(areaLength int, constraints []Constraint) []int {
	areaLength = maxInt(0, areaLength)
	lengths := make([]int, len(constraints))
	total := 0
	hasMin := false

	for i, constraint := range constraints {
		length := constraintLengthValue(areaLength, constraint)
		lengths[i] = length
		total += length
		if constraint.kind == constraintMin {
			hasMin = true
		}
	}

	switch {
	case total < areaLength:
		surplus := areaLength - total
		if hasMin {
			for i, constraint := range constraints {
				if constraint.kind == constraintMin {
					lengths[i] += surplus
					break
				}
			}
		} else if len(lengths) > 0 {
			lengths[len(lengths)-1] += surplus
		}
	case total > areaLength:
		shrinkLengths(lengths, constraints, total-areaLength, false)
		shrinkLengths(lengths, constraints, sumInts(lengths)-areaLength, true)
	}

	return lengths
}

func constraintLengthValue(areaLength int, constraint Constraint) int {
	switch constraint.kind {
	case constraintLength:
		return clampInt(constraint.value, 0, areaLength)
	case constraintMin:
		return clampInt(constraint.value, 0, areaLength)
	case constraintPercentage:
		percent := clampInt(constraint.value, 0, 100)
		return areaLength * percent / 100
	case constraintRatio:
		if constraint.denominator <= 0 {
			return areaLength
		}
		return clampInt(areaLength*constraint.value/constraint.denominator, 0, areaLength)
	default:
		return 0
	}
}

func shrinkLengths(lengths []int, constraints []Constraint, shortage int, includeMin bool) {
	for i := len(lengths) - 1; i >= 0 && shortage > 0; i-- {
		if constraints[i].kind == constraintMin && !includeMin {
			continue
		}
		reduction := minInt(lengths[i], shortage)
		lengths[i] -= reduction
		shortage -= reduction
	}
}

func sumInts(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
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
