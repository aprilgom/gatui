package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/text"
)

type Widget interface {
	Render(area layout.Rect, buf *buffer.Buffer)
}

type WidgetRef interface {
	RenderRef(area layout.Rect, buf *buffer.Buffer)
}

type Wrap struct {
	Trim bool
}

type Paragraph struct {
	text text.Text
	wrap *Wrap
}

func NewParagraph(content text.Text) Paragraph {
	return Paragraph{text: content}
}

func (p Paragraph) Wrap(wrap Wrap) Paragraph {
	p.wrap = &wrap
	return p
}

func (p Paragraph) Render(area layout.Rect, buf *buffer.Buffer) {
	y := area.Y
	for _, line := range p.text.Lines {
		x := area.X
		for _, span := range line.Spans {
			for _, r := range span.Content {
				if x >= area.X+area.Width || y >= area.Y+area.Height {
					return
				}
				buf.SetCell(x, y, buffer.Cell{Symbol: string(r), Style: span.Style})
				x++
			}
		}
		y++
	}
}

type Block struct {
	title text.Line
}

func NewBlock() Block {
	return Block{}
}

func (b Block) Title(title text.Line) Block {
	b.title = title
	return b
}

func (b Block) Render(area layout.Rect, buf *buffer.Buffer) {
	x := area.X
	for _, span := range b.title.Spans {
		for _, r := range span.Content {
			if x >= area.X+area.Width {
				return
			}
			buf.SetCell(x, area.Y, buffer.Cell{Symbol: string(r), Style: span.Style})
			x++
		}
	}
}

type Clear struct{}

func (Clear) Render(area layout.Rect, buf *buffer.Buffer) {
	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < area.X+area.Width; x++ {
			buf.SetCell(x, y, buffer.Cell{Symbol: " "})
		}
	}
}
