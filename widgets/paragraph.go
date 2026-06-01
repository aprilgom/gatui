package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type Wrap struct {
	Trim bool
}

type Paragraph struct {
	text      text.Text
	style     style.Style
	wrap      *Wrap
	block     *Block
	alignment layout.Alignment
	scrollY   int
	scrollX   int
}

func NewParagraph(content text.Text) Paragraph {
	return Paragraph{text: content, style: style.NewStyle(), alignment: layout.Left}
}

func (p Paragraph) Wrap(wrap Wrap) Paragraph {
	p.wrap = &wrap
	return p
}

func (p Paragraph) Block(block Block) Paragraph {
	p.block = &block
	return p
}

func (p Paragraph) Alignment(alignment layout.Alignment) Paragraph {
	p.alignment = alignment
	return p
}

func (p Paragraph) Scroll(y, x int) Paragraph {
	p.scrollY = y
	p.scrollX = x
	return p
}

func (p Paragraph) Style(paragraphStyle style.Style) Paragraph {
	p.style = paragraphStyle
	return p
}

func (p Paragraph) Fg(color style.Color) Paragraph {
	p.style = p.style.Fg(color)
	return p
}

func (p Paragraph) Bg(color style.Color) Paragraph {
	p.style = p.style.Bg(color)
	return p
}

func (p Paragraph) Bold() Paragraph {
	p.style = p.style.AddModifier(style.ModifierBold)
	return p
}

func (p Paragraph) Italic() Paragraph {
	p.style = p.style.AddModifier(style.ModifierItalic)
	return p
}

func (p Paragraph) Cyan() Paragraph {
	return p.Fg(style.Cyan)
}

func (p Paragraph) LineCount(width int) int {
	if width < 1 {
		return 0
	}
	verticalSpace := 0
	if p.block != nil {
		width = maxInt(0, width-p.block.horizontalSpace())
		verticalSpace = p.block.verticalSpace()
	}
	if p.wrap == nil {
		return len(p.text.Lines) + verticalSpace
	}
	return len(p.renderLines(width)) + verticalSpace
}

func (p Paragraph) LineWidth() int {
	width := p.text.Width()
	if p.block != nil {
		width += p.block.horizontalSpace()
	}
	return width
}

func (p Paragraph) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, p.style)
	textArea := area
	if p.block != nil {
		p.block.Render(area, buf)
		textArea = p.block.Inner(area)
	}
	lines := p.renderLines(textArea.Width)
	if p.scrollY < len(lines) {
		lines = lines[p.scrollY:]
	} else {
		lines = nil
	}
	for y := 0; y < textArea.Height && y < len(lines); y++ {
		line := lines[y]
		if p.scrollX > 0 && line.alignment == layout.Left {
			line = line.skip(p.scrollX)
		}
		offset := paragraphLineOffset(line.width(), textArea.Width, line.alignment)
		x := textArea.X + offset
		for _, cell := range line.cells {
			if x >= textArea.X+textArea.Width {
				break
			}
			cell.Style = p.style.Patch(cell.Style)
			buf.SetCell(x, textArea.Y+y, cell)
			x += cellDisplayWidth(cell)
		}
	}
}

func paragraphLineOffset(lineWidth, areaWidth int, alignment layout.Alignment) int {
	if lineWidth >= areaWidth {
		return 0
	}
	switch alignment {
	case layout.Center:
		return areaWidth/2 - lineWidth/2
	case layout.Right:
		return areaWidth - lineWidth
	default:
		return 0
	}
}

func (p Paragraph) renderLines(width int) []renderLine {
	if width <= 0 {
		return nil
	}
	var lines []renderLine
	for _, line := range p.text.Lines {
		alignment := p.alignment
		if p.text.Alignment != nil {
			alignment = *p.text.Alignment
		}
		if line.Alignment != nil {
			alignment = *line.Alignment
		}
		cells := cellsFromLineWithStyle(line, p.text.Style)
		if p.wrap == nil {
			lines = append(lines, renderLine{cells: append([]buffer.Cell(nil), cells...), alignment: alignment})
			continue
		}
		for _, wrapped := range wrapCells(cells, width, p.wrap.Trim) {
			lines = append(lines, renderLine{cells: wrapped, alignment: alignment})
		}
	}
	return lines
}
