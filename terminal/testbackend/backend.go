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
	hideCursorCount int
	showCursorCount int
	cursorPositions []layout.Position
	events          []terminal.Event
}

func New(width, height int) *Backend {
	return &Backend{size: layout.Size{Width: width, Height: height}}
}

func (b *Backend) Size() (layout.Size, error) {
	return b.size, nil
}

func (b *Backend) Draw(diffs []buffer.CellDiff) error {
	copied := make([]buffer.CellDiff, len(diffs))
	copy(copied, diffs)
	b.draws = append(b.draws, copied)
	return nil
}

func (b *Backend) Flush() error {
	b.flushCount++
	return nil
}

func (b *Backend) Clear() error {
	b.clearCount++
	return nil
}

func (b *Backend) PollEvent() (terminal.Event, error) {
	if len(b.events) == 0 {
		return terminal.UnknownEvent{}, nil
	}
	event := b.events[0]
	b.events = b.events[1:]
	return event, nil
}

func (b *Backend) HideCursor() error {
	b.hideCursorCount++
	return nil
}

func (b *Backend) ShowCursor() error {
	b.showCursorCount++
	return nil
}

func (b *Backend) SetCursorPosition(pos layout.Position) error {
	b.cursorPositions = append(b.cursorPositions, pos)
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

func (b *Backend) HideCursorCount() int {
	return b.hideCursorCount
}

func (b *Backend) ShowCursorCount() int {
	return b.showCursorCount
}

func (b *Backend) CursorPositions() []layout.Position {
	return append([]layout.Position(nil), b.cursorPositions...)
}

func (b *Backend) PushEvent(event terminal.Event) {
	b.events = append(b.events, event)
}
