package terminal

import (
	"errors"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
)

type TerminalOptions struct {
	Viewport Viewport
}

type Terminal struct {
	backend        Backend
	previous       *buffer.Buffer
	current        *buffer.Buffer
	area           layout.Rect
	lastKnownArea  layout.Rect
	viewport       Viewport
	count          int
	cursorPosition *layout.Position
}

type CompletedFrame struct {
	Area   layout.Rect
	Buffer *buffer.Buffer
	Count  int
}

func DefaultTerminalOptions() TerminalOptions {
	return TerminalOptions{Viewport: FullscreenViewport()}
}

func New(backend Backend) (*Terminal, error) {
	return NewWithOptions(backend, DefaultTerminalOptions())
}

func NewWithOptions(backend Backend, options TerminalOptions) (*Terminal, error) {
	if backend == nil {
		return nil, errors.New("terminal backend is nil")
	}
	area := options.Viewport.area
	lastKnownArea := area
	cursorPosition := (*layout.Position)(nil)
	switch options.Viewport.kind {
	case viewportFullscreen:
		size, err := backend.Size()
		if err != nil {
			return nil, err
		}
		area = layout.NewRect(0, 0, size.Width, size.Height)
		lastKnownArea = area
	case viewportInline:
		size, err := backend.Size()
		if err != nil {
			return nil, err
		}
		lastKnownArea = layout.NewRect(0, 0, size.Width, size.Height)
		var cursor layout.Position
		area, cursor, err = computeInlineArea(backend, options.Viewport.height, size, 0)
		if err != nil {
			return nil, err
		}
		cursorPosition = &cursor
	}
	return &Terminal{
		backend:        backend,
		previous:       buffer.Empty(area),
		current:        buffer.Empty(area),
		area:           area,
		lastKnownArea:  lastKnownArea,
		viewport:       options.Viewport,
		cursorPosition: cursorPosition,
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
	frame := &Frame{area: t.area, buffer: t.current, count: t.count}
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

	completed := &CompletedFrame{Area: t.lastKnownArea, Buffer: t.previous, Count: t.count}
	t.count++
	return completed, nil
}

func (t *Terminal) Autoresize() error {
	if t.viewport.kind == viewportFixed {
		return nil
	}
	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	area := layout.NewRect(0, 0, size.Width, size.Height)
	if area == t.lastKnownArea {
		return nil
	}
	return t.Resize(area)
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

func (t *Terminal) CurrentBuffer() *buffer.Buffer {
	return t.current
}

func (t *Terminal) Frame() *Frame {
	return &Frame{area: t.area, buffer: t.current, count: t.count}
}

func (t *Terminal) Resize(area layout.Rect) error {
	nextArea := area
	var cursorToRestore *layout.Position
	if t.viewport.kind == viewportInline {
		inlineArea, cursor, err := t.resizeInlineArea(area)
		if err != nil {
			return err
		}
		nextArea = inlineArea
		cursorToRestore = &cursor
	}
	t.lastKnownArea = area
	if nextArea.Width < t.area.Width {
		nextArea.Y = 0
		if err := t.backend.ClearRegion(ClearAll); err != nil {
			return err
		}
	}
	t.setViewportArea(nextArea)
	if t.viewport.kind == viewportFixed {
		t.viewport.area = nextArea
	}
	if err := t.clearViewport(); err != nil {
		return err
	}
	if cursorToRestore != nil {
		if err := t.backend.SetCursorPosition(*cursorToRestore); err != nil {
			return err
		}
	}
	return nil
}

func (t *Terminal) setViewportArea(area layout.Rect) {
	t.area = area
	t.previous.Resize(area)
	t.current.Resize(area)
}

func (t *Terminal) Backend() Backend {
	return t.backend
}

func (t *Terminal) Size() (layout.Size, error) {
	return t.backend.Size()
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

func (t *Terminal) GetCursorPosition() (layout.Position, error) {
	return t.backend.GetCursorPosition()
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
