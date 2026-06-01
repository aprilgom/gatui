package terminal

import (
	"errors"

	"gatui/buffer"
	"gatui/layout"
	"gatui/widgets"
)

type Backend interface {
	Size() (layout.Size, error)
	Draw([]buffer.CellDiff) error
	Flush() error
	Clear() error
	HideCursor() error
	ShowCursor() error
	SetCursorPosition(layout.Position) error
}

type Terminal struct {
	backend  Backend
	previous *buffer.Buffer
	current  *buffer.Buffer
	area     layout.Rect
	count    int
}

type Frame struct {
	area           layout.Rect
	buffer         *buffer.Buffer
	cursorPosition *layout.Position
}

type CompletedFrame struct {
	Area   layout.Rect
	Buffer *buffer.Buffer
	Count  int
}

func New(backend Backend) (*Terminal, error) {
	if backend == nil {
		return nil, errors.New("terminal backend is nil")
	}
	size, err := backend.Size()
	if err != nil {
		return nil, err
	}
	area := layout.NewRect(0, 0, size.Width, size.Height)
	return &Terminal{
		backend:  backend,
		previous: buffer.Empty(area),
		current:  buffer.Empty(area),
		area:     area,
	}, nil
}

func (t *Terminal) Draw(render func(*Frame)) (*CompletedFrame, error) {
	t.current.Reset()
	frame := &Frame{area: t.area, buffer: t.current}
	if render != nil {
		render(frame)
	}

	diffs := t.previous.Diff(t.current)
	if err := t.backend.Draw(diffs); err != nil {
		return nil, err
	}
	if err := t.updateCursor(frame.cursorPosition); err != nil {
		return nil, err
	}
	if err := t.backend.Flush(); err != nil {
		return nil, err
	}

	completed := &CompletedFrame{Area: t.area, Buffer: t.current, Count: t.count}
	t.previous, t.current = t.current, t.previous
	t.count++
	return completed, nil
}

func (t *Terminal) Resize(area layout.Rect) {
	t.area = area
	t.previous.Resize(area)
	t.previous.Reset()
	t.current.Resize(area)
	t.current.Reset()
}

func (t *Terminal) Clear() error {
	if err := t.backend.Clear(); err != nil {
		return err
	}
	t.previous.Reset()
	return nil
}

func (t *Terminal) Backend() Backend {
	return t.backend
}

func (t *Terminal) updateCursor(pos *layout.Position) error {
	if pos == nil {
		return t.backend.HideCursor()
	}
	if err := t.backend.ShowCursor(); err != nil {
		return err
	}
	return t.backend.SetCursorPosition(*pos)
}

func (f *Frame) Area() layout.Rect {
	return f.area
}

func (f *Frame) Buffer() *buffer.Buffer {
	return f.buffer
}

func (f *Frame) RenderWidget(widget widgets.Widget, area layout.Rect) {
	if widget == nil {
		return
	}
	widget.Render(area, f.buffer)
}

func (f *Frame) SetCursorPosition(pos layout.Position) {
	f.cursorPosition = &pos
}
