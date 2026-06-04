package layout

import "fmt"

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) Position {
	return Position{X: maxInt(0, x), Y: maxInt(0, y)}
}

func PositionOrigin() Position {
	return NewPosition(0, 0)
}

func PositionMin() Position {
	return PositionOrigin()
}

func PositionMax() Position {
	return NewPosition(MaxCoordinate, MaxCoordinate)
}

func PositionFromTuple(x, y int) Position {
	return NewPosition(x, y)
}

func PositionFromRect(rect Rect) Position {
	return rect.AsPosition()
}

func (p Position) Tuple() (x, y int) {
	return p.X, p.Y
}

func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
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

func SizeZero() Size {
	return NewSize(0, 0)
}

func SizeMin() Size {
	return SizeZero()
}

func SizeMax() Size {
	return NewSize(MaxCoordinate, MaxCoordinate)
}

func SizeFromTuple(width, height int) Size {
	return NewSize(width, height)
}

func SizeFromRect(rect Rect) Size {
	return rect.AsSize()
}

func (s Size) Area() int {
	return s.Width * s.Height
}

func (s Size) Tuple() (width, height int) {
	return s.Width, s.Height
}

func (s Size) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

type Margin struct {
	Horizontal int
	Vertical   int
}

func NewMargin(horizontal, vertical int) Margin {
	return Margin{Horizontal: horizontal, Vertical: vertical}
}

func MarginFromInt(value int) Margin {
	return NewMargin(value, value)
}

func (m Margin) String() string {
	return fmt.Sprintf("%dx%d", m.Horizontal, m.Vertical)
}

type Offset struct {
	X int
	Y int
}

func NewOffset(x, y int) Offset {
	return Offset{X: x, Y: y}
}

func OffsetZero() Offset {
	return NewOffset(0, 0)
}

func OffsetMin() Offset {
	minimum := -maxIntValue() - 1
	return NewOffset(minimum, minimum)
}

func OffsetMax() Offset {
	maximum := maxIntValue()
	return NewOffset(maximum, maximum)
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

type HorizontalAlignment int

const (
	Left HorizontalAlignment = iota
	Center
	Right
)

type Alignment = HorizontalAlignment

func (a HorizontalAlignment) String() string {
	switch a {
	case Left:
		return "Left"
	case Center:
		return "Center"
	case Right:
		return "Right"
	default:
		return fmt.Sprintf("HorizontalAlignment(%d)", a)
	}
}

func ParseHorizontalAlignment(value string) (HorizontalAlignment, error) {
	switch value {
	case "Left":
		return Left, nil
	case "Center":
		return Center, nil
	case "Right":
		return Right, nil
	default:
		return Left, fmt.Errorf("invalid horizontal alignment: %q", value)
	}
}

func ParseAlignment(value string) (Alignment, error) {
	return ParseHorizontalAlignment(value)
}

type VerticalAlignment int

const (
	Top VerticalAlignment = iota
	VerticalCenter
	Bottom
)

func (a VerticalAlignment) String() string {
	switch a {
	case Top:
		return "Top"
	case VerticalCenter:
		return "Center"
	case Bottom:
		return "Bottom"
	default:
		return fmt.Sprintf("VerticalAlignment(%d)", a)
	}
}

func ParseVerticalAlignment(value string) (VerticalAlignment, error) {
	switch value {
	case "Top":
		return Top, nil
	case "Center":
		return VerticalCenter, nil
	case "Bottom":
		return Bottom, nil
	default:
		return Top, fmt.Errorf("invalid vertical alignment: %q", value)
	}
}

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
