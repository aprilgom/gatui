package layout

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) Position {
	return Position{X: maxInt(0, x), Y: maxInt(0, y)}
}

func (p Position) Offset(offset Offset) Position {
	return NewPosition(p.X+offset.X, p.Y+offset.Y)
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
