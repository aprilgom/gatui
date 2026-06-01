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

func (r Rect) Area() int {
	return r.Width * r.Height
}

func (r Rect) Rows() []Rect {
	rows := make([]Rect, r.Height)
	for i := range rows {
		rows[i] = NewRect(r.X, r.Y+i, r.Width, 1)
	}
	return rows
}

func (r Rect) Columns() []Rect {
	columns := make([]Rect, r.Width)
	for i := range columns {
		columns[i] = NewRect(r.X+i, r.Y, 1, r.Height)
	}
	return columns
}

func (r Rect) Positions() []Position {
	positions := make([]Position, 0, r.Area())
	for y := r.Y; y < r.Bottom(); y++ {
		for x := r.X; x < r.Right(); x++ {
			positions = append(positions, NewPosition(x, y))
		}
	}
	return positions
}

func (r Rect) IsEmpty() bool {
	return r.Width == 0 || r.Height == 0
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
	return NewRect(maxInt(0, r.X+offset.X), maxInt(0, r.Y+offset.Y), r.Width, r.Height)
}

func (r Rect) Resize(size Size) Rect {
	return NewRect(r.X, r.Y, size.Width, size.Height)
}

func (r Rect) Union(other Rect) Rect {
	x1 := minInt(r.X, other.X)
	y1 := minInt(r.Y, other.Y)
	x2 := maxInt(r.Right(), other.Right())
	y2 := maxInt(r.Bottom(), other.Bottom())
	return NewRect(x1, y1, x2-x1, y2-y1)
}

func (r Rect) Intersection(other Rect) Rect {
	x1 := maxInt(r.X, other.X)
	y1 := maxInt(r.Y, other.Y)
	x2 := minInt(r.Right(), other.Right())
	y2 := minInt(r.Bottom(), other.Bottom())
	return NewRect(x1, y1, maxInt(0, x2-x1), maxInt(0, y2-y1))
}

func (r Rect) Intersects(other Rect) bool {
	return r.X < other.Right() &&
		r.Right() > other.X &&
		r.Y < other.Bottom() &&
		r.Bottom() > other.Y
}

func (r Rect) Contains(position Position) bool {
	return position.X >= r.X &&
		position.X < r.Right() &&
		position.Y >= r.Y &&
		position.Y < r.Bottom()
}

func (r Rect) Clamp(other Rect) Rect {
	width := minInt(r.Width, other.Width)
	height := minInt(r.Height, other.Height)
	x := clampInt(r.X, other.X, other.Right()-width)
	y := clampInt(r.Y, other.Y, other.Bottom()-height)
	return NewRect(x, y, width, height)
}

func (r Rect) AsPosition() Position {
	return NewPosition(r.X, r.Y)
}

func (r Rect) AsSize() Size {
	return NewSize(r.Width, r.Height)
}

func (r Rect) CenteredHorizontally(constraint Constraint) Rect {
	width := centeredLength(r.Width, constraint)
	return NewRect(r.X+(r.Width-width)/2, r.Y, width, r.Height)
}

func (r Rect) CenteredVertically(constraint Constraint) Rect {
	height := centeredLength(r.Height, constraint)
	return NewRect(r.X, r.Y+(r.Height-height)/2, r.Width, height)
}

func (r Rect) Centered(horizontal, vertical Constraint) Rect {
	return r.CenteredHorizontally(horizontal).CenteredVertically(vertical)
}
