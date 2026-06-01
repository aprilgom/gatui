package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

type Clear struct{}

func (Clear) Render(area layout.Rect, buf *buffer.Buffer) {
	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < area.X+area.Width; x++ {
			buf.SetCell(x, y, buffer.Cell{Symbol: " ", Style: style.NewStyle()})
		}
	}
}
