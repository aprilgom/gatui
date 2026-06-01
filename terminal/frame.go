package terminal

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/widgets"
)

type Frame struct {
	area           layout.Rect
	buffer         *buffer.Buffer
	count          int
	cursorPosition *layout.Position
}

func (f *Frame) Area() layout.Rect {
	return f.area
}

func (f *Frame) Size() layout.Size {
	return layout.Size{Width: f.area.Width, Height: f.area.Height}
}

func (f *Frame) Buffer() *buffer.Buffer {
	return f.buffer
}

func (f *Frame) Count() int {
	return f.count
}

func (f *Frame) RenderWidget(widget widgets.Widget, area layout.Rect) {
	if widget == nil {
		return
	}
	widget.Render(area, f.buffer)
}

func (f *Frame) RenderStatefulWidget(widget widgets.StatefulWidget, area layout.Rect, state any) {
	if widget == nil {
		return
	}
	widget.RenderStatefulRef(area, f.buffer, state)
}

func (f *Frame) SetCursorPosition(pos layout.Position) {
	f.cursorPosition = &pos
}
