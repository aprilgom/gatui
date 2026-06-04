package terminal

import (
	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
)

func (t *Terminal) InsertBefore(height int, render func(*buffer.Buffer)) error {
	if t.viewport.kind != viewportInline {
		return nil
	}
	if height < 0 {
		height = 0
	}
	if backend, ok := t.backend.(ScrollingRegionBackend); ok {
		return t.insertBeforeScrollingRegions(backend, height, render)
	}
	return t.insertBeforeNoScrollingRegions(height, render)
}

func (t *Terminal) insertBeforeNoScrollingRegions(height int, render func(*buffer.Buffer)) error {
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

func (t *Terminal) insertBeforeScrollingRegions(backend ScrollingRegionBackend, height int, render func(*buffer.Buffer)) error {
	area := layout.NewRect(0, 0, t.area.Width, height)
	insert := buffer.Empty(area)
	if render != nil {
		render(insert)
	}
	cells := insert.Cells

	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	if t.area.Height == size.Height {
		first := true
		for len(cells) > 0 {
			if first {
				cells, err = t.drawLines(0, 1, cells)
			} else {
				cells, err = t.drawLinesOverCleared(0, 1, cells)
			}
			if err != nil {
				return err
			}
			first = false
			if err := backend.ScrollRegionUp(0, 1, 1); err != nil {
				return err
			}
		}
		topLine := append([]buffer.Cell(nil), t.previous.Cells[:t.area.Width]...)
		_, err = t.drawLinesOverCleared(0, 1, topLine)
		return err
	}

	remainingHeight := height
	viewportTop := t.area.Y
	viewportBottom := t.area.Bottom()
	screenBottom := size.Height
	if viewportBottom < screenBottom {
		toDraw := minInt(remainingHeight, screenBottom-viewportBottom)
		if err := backend.ScrollRegionDown(viewportTop, viewportBottom+toDraw, toDraw); err != nil {
			return err
		}
		cells, err = t.drawLinesOverCleared(viewportTop, toDraw, cells)
		if err != nil {
			return err
		}
		t.setViewportArea(layout.NewRect(t.area.X, viewportTop+toDraw, t.area.Width, t.area.Height))
		remainingHeight -= toDraw
	}

	viewportTop = t.area.Y
	for remainingHeight > 0 {
		toDraw := minInt(remainingHeight, viewportTop)
		if err := backend.ScrollRegionUp(0, viewportTop, toDraw); err != nil {
			return err
		}
		cells, err = t.drawLinesOverCleared(viewportTop-toDraw, toDraw, cells)
		if err != nil {
			return err
		}
		remainingHeight -= toDraw
	}
	return nil
}

func (t *Terminal) drawLines(yOffset, linesToDraw int, cells []buffer.Cell) ([]buffer.Cell, error) {
	width := t.area.Width
	count := min(width*linesToDraw, len(cells))
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

func (t *Terminal) drawLinesOverCleared(yOffset, linesToDraw int, cells []buffer.Cell) ([]buffer.Cell, error) {
	width := t.area.Width
	count := min(width*linesToDraw, len(cells))
	toDraw := cells[:count]
	remainder := cells[count:]
	if linesToDraw <= 0 {
		return remainder, nil
	}
	area := layout.NewRect(0, yOffset, width, linesToDraw)
	old := buffer.Empty(area)
	next := &buffer.Buffer{
		Area:  area,
		Cells: append([]buffer.Cell(nil), toDraw...),
	}
	if err := t.backend.Draw(old.Diff(next)); err != nil {
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

func (t *Terminal) resizeInlineArea(terminalArea layout.Rect) (layout.Rect, layout.Position, error) {
	offset := 0
	if t.cursorPosition != nil {
		offset = max(t.cursorPosition.Y-t.area.Y, 0)
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
	maxHeight := max(min(height, size.Height), 0)

	linesAfterCursor := max(height-offsetInPreviousViewport-1, 0)
	if err := backend.AppendLines(linesAfterCursor); err != nil {
		return layout.Rect{}, layout.Position{}, err
	}

	availableLines := max(size.Height-row-1, 0)
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
