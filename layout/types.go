package layout

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func NewRect(x, y, width, height int) Rect {
	return Rect{X: x, Y: y, Width: width, Height: height}
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
	kind  constraintKind
	value int
}

type constraintKind int

const (
	constraintLength constraintKind = iota
	constraintMin
)

func Length(value int) Constraint {
	return Constraint{kind: constraintLength, value: value}
}

func Min(value int) Constraint {
	return Constraint{kind: constraintMin, value: value}
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

	rects := make([]Rect, 0, len(l.constraints))
	cursorX := area.X
	cursorY := area.Y
	remainingWidth := area.Width
	remainingHeight := area.Height

	for _, constraint := range l.constraints {
		width := remainingWidth
		height := remainingHeight
		if constraint.kind == constraintLength {
			if l.direction == Horizontal {
				width = minInt(constraint.value, remainingWidth)
			} else {
				height = minInt(constraint.value, remainingHeight)
			}
		}

		rect := Rect{X: cursorX, Y: cursorY, Width: width, Height: height}
		rects = append(rects, rect)

		if l.direction == Horizontal {
			cursorX += width
			remainingWidth -= width
		} else {
			cursorY += height
			remainingHeight -= height
		}
	}

	return rects
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
