package testbackend

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/terminal"
)

type Backend struct {
	size            layout.Size
	draws           [][]buffer.CellDiff
	flushCount      int
	clearCount      int
	clearRegions    []terminal.ClearType
	hideCursorCount int
	showCursorCount int
	cursorPositions []layout.Position
	cursorPosition  layout.Position
	cursorVisible   bool
	appendLines     []int
	cells           *buffer.Buffer
}

func New(width, height int) *Backend {
	area := layout.NewRect(0, 0, width, height)
	return &Backend{size: layout.Size{Width: width, Height: height}, cursorVisible: true, cells: buffer.Empty(area)}
}

func WithLines(lines []string) *Backend {
	cells := buffer.WithLines(lines)
	return &Backend{
		size:          layout.Size{Width: cells.Area.Width, Height: cells.Area.Height},
		cursorVisible: true,
		cells:         cells,
	}
}

func (b *Backend) Size() (layout.Size, error) {
	return b.size, nil
}

func (b *Backend) SetSize(width, height int) {
	b.size = layout.Size{Width: width, Height: height}
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, width, height))
		return
	}
	b.cells.Resize(layout.NewRect(0, 0, width, height))
}

func (b *Backend) Draw(diffs []buffer.CellDiff) error {
	copied := make([]buffer.CellDiff, len(diffs))
	copy(copied, diffs)
	b.draws = append(b.draws, copied)
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	for _, diff := range diffs {
		b.cells.SetCell(diff.X, diff.Y, diff.Cell)
	}
	return nil
}

func (b *Backend) Flush() error {
	b.flushCount++
	return nil
}

func (b *Backend) Clear() error {
	b.clearCount++
	return b.ClearRegion(terminal.ClearAll)
}

func (b *Backend) ClearRegion(clearType terminal.ClearType) error {
	b.clearRegions = append(b.clearRegions, clearType)
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	switch clearType {
	case terminal.ClearAll:
		b.cells.Reset()
	case terminal.ClearAfterCursor:
		for y := b.cursorPosition.Y; y < b.size.Height; y++ {
			startX := 0
			if y == b.cursorPosition.Y {
				startX = b.cursorPosition.X
			}
			for x := startX; x < b.size.Width; x++ {
				b.cells.SetCell(x, y, buffer.NewCell(" "))
			}
		}
	case terminal.ClearCurrentLine:
		for x := 0; x < b.size.Width; x++ {
			b.cells.SetCell(x, b.cursorPosition.Y, buffer.NewCell(" "))
		}
	}
	return nil
}

func (b *Backend) HideCursor() error {
	b.hideCursorCount++
	b.cursorVisible = false
	return nil
}

func (b *Backend) ShowCursor() error {
	b.showCursorCount++
	b.cursorVisible = true
	return nil
}

func (b *Backend) SetCursorPosition(pos layout.Position) error {
	b.cursorPositions = append(b.cursorPositions, pos)
	b.cursorPosition = pos
	return nil
}

func (b *Backend) GetCursorPosition() (layout.Position, error) {
	return b.cursorPosition, nil
}

func (b *Backend) AppendLines(count int) error {
	b.appendLines = append(b.appendLines, count)
	if count <= 0 {
		return nil
	}
	if b.cells == nil {
		b.cells = buffer.Empty(layout.NewRect(0, 0, b.size.Width, b.size.Height))
	}
	scroll := b.cursorPosition.Y + count - (b.size.Height - 1)
	if scroll > b.size.Height {
		scroll = b.size.Height
	}
	if scroll > 0 {
		for y := 0; y < b.size.Height-scroll; y++ {
			for x := 0; x < b.size.Width; x++ {
				cell, _ := b.cells.CellAt(x, y+scroll)
				b.cells.SetCell(x, y, cell)
			}
		}
		for y := b.size.Height - scroll; y < b.size.Height; y++ {
			for x := 0; x < b.size.Width; x++ {
				b.cells.SetCell(x, y, buffer.NewCell(" "))
			}
		}
	}
	b.cursorPosition = layout.Position{X: 0, Y: b.size.Height - 1}
	return nil
}

func (b *Backend) Draws() [][]buffer.CellDiff {
	draws := make([][]buffer.CellDiff, len(b.draws))
	for i := range b.draws {
		draws[i] = make([]buffer.CellDiff, len(b.draws[i]))
		copy(draws[i], b.draws[i])
	}
	return draws
}

func (b *Backend) FlushCount() int {
	return b.flushCount
}

func (b *Backend) ClearCount() int {
	return b.clearCount
}

func (b *Backend) ClearRegions() []terminal.ClearType {
	return append([]terminal.ClearType(nil), b.clearRegions...)
}

func (b *Backend) Lines() []string {
	if b.cells == nil {
		return nil
	}
	return b.cells.Lines()
}

func (b *Backend) HideCursorCount() int {
	return b.hideCursorCount
}

func (b *Backend) ShowCursorCount() int {
	return b.showCursorCount
}

func (b *Backend) CursorPositions() []layout.Position {
	return append([]layout.Position(nil), b.cursorPositions...)
}

func (b *Backend) CursorVisible() bool {
	return b.cursorVisible
}

func (b *Backend) CursorPosition() layout.Position {
	return b.cursorPosition
}

func (b *Backend) AppendLinesCalls() []int {
	return append([]int(nil), b.appendLines...)
}
