package terminal

import (
	"gatui/buffer"
	"gatui/layout"
)

type Backend interface {
	Size() (layout.Size, error)
	WindowSize() (WindowSize, error)
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

type WindowSize struct {
	ColumnsRows layout.Size
	Pixels      layout.Size
}

type ScrollingRegionBackend interface {
	ScrollRegionUp(startY, endY, count int) error
	ScrollRegionDown(startY, endY, count int) error
}
