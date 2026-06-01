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
	backend        Backend
	previous       *buffer.Buffer
	current        *buffer.Buffer
	area           layout.Rect
	count          int
	cursorPosition *layout.Position
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
	return t.TryDraw(func(frame *Frame) error {
		if render != nil {
			render(frame)
		}
		return nil
	})
}

func (t *Terminal) TryDraw(render func(*Frame) error) (*CompletedFrame, error) {
	if err := t.Autoresize(); err != nil {
		return nil, err
	}
	snapshot := append([]buffer.Cell(nil), t.current.Cells...)
	t.current.Reset()
	frame := &Frame{area: t.area, buffer: t.current}
	if render != nil {
		if err := render(frame); err != nil {
			copy(t.current.Cells, snapshot)
			return nil, err
		}
	}

	if err := t.Flush(); err != nil {
		return nil, err
	}
	if err := t.updateCursor(frame.cursorPosition); err != nil {
		return nil, err
	}
	t.SwapBuffers()
	if err := t.backend.Flush(); err != nil {
		return nil, err
	}

	completed := &CompletedFrame{Area: t.area, Buffer: t.previous, Count: t.count}
	t.count++
	return completed, nil
}

func (t *Terminal) Autoresize() error {
	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	area := layout.NewRect(0, 0, size.Width, size.Height)
	if area == t.area {
		return nil
	}
	t.Resize(area)
	return nil
}

func (t *Terminal) Flush() error {
	return t.backend.Draw(t.previous.Diff(t.current))
}

func (t *Terminal) SwapBuffers() {
	t.previous.Reset()
	t.previous, t.current = t.current, t.previous
}

func (t *Terminal) Area() layout.Rect {
	return t.area
}

func (t *Terminal) Frame() *Frame {
	return &Frame{area: t.area, buffer: t.current}
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

func (t *Terminal) HideCursor() error {
	if err := t.backend.HideCursor(); err != nil {
		return err
	}
	t.cursorPosition = nil
	return nil
}

func (t *Terminal) ShowCursor() error {
	return t.backend.ShowCursor()
}

func (t *Terminal) SetCursorPosition(pos layout.Position) error {
	if err := t.backend.SetCursorPosition(pos); err != nil {
		return err
	}
	t.cursorPosition = &pos
	return nil
}

func (t *Terminal) updateCursor(pos *layout.Position) error {
	if pos == nil {
		return t.HideCursor()
	}
	if err := t.ShowCursor(); err != nil {
		return err
	}
	return t.SetCursorPosition(*pos)
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
