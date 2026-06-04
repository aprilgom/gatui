package layout

import "fmt"

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) Position {
	return Position{X: maxInt(0, x), Y: maxInt(0, y)}
}

func (p Position) Offset(offset Offset) Position {
	return p.AddOffset(offset)
}

func (p Position) AddOffset(offset Offset) Position {
	return NewPosition(p.X+offset.X, p.Y+offset.Y)
}

func (p Position) SubOffset(offset Offset) Position {
	return NewPosition(p.X-offset.X, p.Y-offset.Y)
}

type Size struct {
	Width  int
	Height int
}

func NewSize(width, height int) Size {
	return Size{Width: maxInt(0, width), Height: maxInt(0, height)}
}

func (s Size) Area() int {
	return s.Width * s.Height
}

func (s Size) Tuple() (width, height int) {
	return s.Width, s.Height
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

func OffsetFromPosition(position Position) Offset {
	return NewOffset(position.X, position.Y)
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func (d Direction) Other() Direction {
	if d == Horizontal {
		return Vertical
	}
	return Horizontal
}

func (d Direction) String() string {
	switch d {
	case Horizontal:
		return "Horizontal"
	case Vertical:
		return "Vertical"
	default:
		return fmt.Sprintf("Direction(%d)", d)
	}
}

func ParseDirection(value string) (Direction, error) {
	switch value {
	case "Horizontal":
		return Horizontal, nil
	case "Vertical":
		return Vertical, nil
	default:
		return Horizontal, fmt.Errorf("invalid direction: %q", value)
	}
}

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
)

type Flex int

const (
	FlexLegacy Flex = iota
	FlexStart
	FlexEnd
	FlexCenter
	FlexSpaceBetween
	FlexSpaceAround
	FlexSpaceEvenly
)
