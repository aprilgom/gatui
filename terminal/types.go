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
	ClearRegion(ClearType) error
	GetCursorPosition() (layout.Position, error)
	AppendLines(count int) error
	HideCursor() error
	ShowCursor() error
	SetCursorPosition(layout.Position) error
}

type ClearType int

const (
	ClearAll ClearType = iota
	ClearAfterCursor
	ClearCurrentLine
)

type viewportKind int

const (
	viewportFullscreen viewportKind = iota
	viewportFixed
	viewportInline
)

type TerminalOptions struct {
	Viewport Viewport
}

type Viewport struct {
	kind   viewportKind
	area   layout.Rect
	height int
}

type Terminal struct {
	backend        Backend
	previous       *buffer.Buffer
	current        *buffer.Buffer
	area           layout.Rect
	viewport       Viewport
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

func FullscreenViewport() Viewport {
	return Viewport{kind: viewportFullscreen}
}

func FixedViewport(area layout.Rect) Viewport {
	return Viewport{kind: viewportFixed, area: area}
}

func InlineViewport(height int) Viewport {
	return Viewport{kind: viewportInline, height: height}
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
	cursorPosition := (*layout.Position)(nil)
	switch options.Viewport.kind {
	case viewportFullscreen:
		size, err := backend.Size()
		if err != nil {
			return nil, err
		}
		area = layout.NewRect(0, 0, size.Width, size.Height)
	case viewportInline:
		size, err := backend.Size()
		if err != nil {
			return nil, err
		}
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
	if t.viewport.kind == viewportFixed {
		return nil
	}
	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	area := layout.NewRect(0, 0, size.Width, size.Height)
	if area == t.area {
		return nil
	}
	return t.Resize(area)
}

func (t *Terminal) Flush() error {
	return t.backend.Draw(t.previous.Diff(t.current))
}

func (t *Terminal) InsertBefore(height int, render func(*buffer.Buffer)) error {
	if t.viewport.kind != viewportInline {
		return nil
	}
	if height < 0 {
		height = 0
	}
	area := layout.NewRect(0, 0, t.area.Width, height)
	insert := buffer.Empty(area)
	if render != nil {
		render(insert)
	}

	cells := insert.Cells
	drawnHeight := t.area.Y
	bufferHeight := height
	viewportHeight := t.area.Height
	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	screenHeight := size.Height

	for bufferHeight+viewportHeight > screenHeight {
		toDraw := minInt(bufferHeight, screenHeight)
		scrollUp := maxInt(0, drawnHeight+toDraw-screenHeight)
		if err := t.scrollUp(scrollUp); err != nil {
			return err
		}
		cells, err = t.drawLines(drawnHeight-scrollUp, toDraw, cells)
		if err != nil {
			return err
		}
		drawnHeight += toDraw - scrollUp
		bufferHeight -= toDraw
	}

	scrollUp := maxInt(0, drawnHeight+bufferHeight+viewportHeight-screenHeight)
	if err := t.scrollUp(scrollUp); err != nil {
		return err
	}
	if _, err := t.drawLines(drawnHeight-scrollUp, bufferHeight, cells); err != nil {
		return err
	}
	drawnHeight += bufferHeight - scrollUp

	t.area.Y = drawnHeight
	t.previous.Resize(t.area)
	t.current.Resize(t.area)
	return t.Clear()
}

func (t *Terminal) drawLines(yOffset, linesToDraw int, cells []buffer.Cell) ([]buffer.Cell, error) {
	width := t.area.Width
	count := width * linesToDraw
	if count > len(cells) {
		count = len(cells)
	}
	toDraw := cells[:count]
	remainder := cells[count:]
	if linesToDraw <= 0 {
		return remainder, nil
	}
	diffs := make([]buffer.CellDiff, 0, len(toDraw))
	for i, cell := range toDraw {
		diffs = append(diffs, buffer.CellDiff{
			X:    i % width,
			Y:    yOffset + i/width,
			Cell: cell,
		})
	}
	if err := t.backend.Draw(diffs); err != nil {
		return nil, err
	}
	if err := t.backend.Flush(); err != nil {
		return nil, err
	}
	return remainder, nil
}

func (t *Terminal) scrollUp(lines int) error {
	if lines <= 0 {
		return nil
	}
	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	if err := t.SetCursorPosition(layout.Position{X: 0, Y: size.Height - 1}); err != nil {
		return err
	}
	return t.backend.AppendLines(lines)
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

func (t *Terminal) resizeInlineArea(terminalArea layout.Rect) (layout.Rect, layout.Position, error) {
	offset := 0
	if t.cursorPosition != nil {
		offset = t.cursorPosition.Y - t.area.Y
		if offset < 0 {
			offset = 0
		}
	}

	originalCursor, err := t.backend.GetCursorPosition()
	if err != nil {
		return layout.Rect{}, layout.Position{}, err
	}
	nextArea, _, err := computeInlineArea(t.backend, t.viewport.height, layout.Size{Width: terminalArea.Width, Height: terminalArea.Height}, offset)
	if err != nil {
		return layout.Rect{}, layout.Position{}, err
	}
	return nextArea, originalCursor, nil
}

func computeInlineArea(backend Backend, height int, size layout.Size, offsetInPreviousViewport int) (layout.Rect, layout.Position, error) {
	pos, err := backend.GetCursorPosition()
	if err != nil {
		return layout.Rect{}, layout.Position{}, err
	}
	row := pos.Y
	maxHeight := height
	if maxHeight > size.Height {
		maxHeight = size.Height
	}
	if maxHeight < 0 {
		maxHeight = 0
	}

	linesAfterCursor := height - offsetInPreviousViewport - 1
	if linesAfterCursor < 0 {
		linesAfterCursor = 0
	}
	if err := backend.AppendLines(linesAfterCursor); err != nil {
		return layout.Rect{}, layout.Position{}, err
	}

	availableLines := size.Height - row - 1
	if availableLines < 0 {
		availableLines = 0
	}
	missingLines := linesAfterCursor - availableLines
	if missingLines > 0 {
		row -= missingLines
		if row < 0 {
			row = 0
		}
	}
	row -= offsetInPreviousViewport
	if row < 0 {
		row = 0
	}

	return layout.NewRect(0, row, size.Width, maxHeight), pos, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (t *Terminal) Clear() error {
	originalCursor, err := t.backend.GetCursorPosition()
	if err != nil {
		return err
	}
	if err := t.clearViewport(); err != nil {
		return err
	}
	return t.backend.SetCursorPosition(originalCursor)
}

func (t *Terminal) clearViewport() error {
	switch t.viewport.kind {
	case viewportFullscreen:
		if err := t.backend.ClearRegion(ClearAll); err != nil {
			return err
		}
	case viewportFixed:
		if err := t.clearFixedViewport(t.area); err != nil {
			return err
		}
	case viewportInline:
		if err := t.backend.SetCursorPosition(layout.Position{X: t.area.X, Y: t.area.Y}); err != nil {
			return err
		}
		if err := t.backend.ClearRegion(ClearAfterCursor); err != nil {
			return err
		}
	}
	t.previous.Reset()
	return nil
}

func (t *Terminal) clearFixedViewport(area layout.Rect) error {
	if area.Width == 0 || area.Height == 0 {
		return nil
	}
	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	isFullWidth := area.X == 0 && area.Width == size.Width
	endsAtBottom := area.Bottom() == size.Height
	if isFullWidth && endsAtBottom {
		if err := t.backend.SetCursorPosition(layout.Position{X: area.X, Y: area.Y}); err != nil {
			return err
		}
		return t.backend.ClearRegion(ClearAfterCursor)
	}
	if isFullWidth {
		for y := area.Y; y < area.Bottom(); y++ {
			if err := t.backend.SetCursorPosition(layout.Position{X: 0, Y: y}); err != nil {
				return err
			}
			if err := t.backend.ClearRegion(ClearCurrentLine); err != nil {
				return err
			}
		}
		return nil
	}

	clearCell := buffer.NewCell(" ")
	diffs := make([]buffer.CellDiff, 0, area.Width*area.Height)
	for y := area.Y; y < area.Bottom(); y++ {
		for x := area.X; x < area.Right(); x++ {
			diffs = append(diffs, buffer.CellDiff{X: x, Y: y, Cell: clearCell})
		}
	}
	return t.backend.Draw(diffs)
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
