package widgets

import (
	"strings"
	"unicode"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
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
	text      text.Text
	wrap      *Wrap
	block     *Block
	alignment layout.Alignment
	scrollY   int
	scrollX   int
}

func NewParagraph(content text.Text) Paragraph {
	return Paragraph{text: content, alignment: layout.Left}
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

func (p Paragraph) Fg(color style.Color) Paragraph {
	p.text = p.text.Fg(color)
	return p
}

func (p Paragraph) Bg(color style.Color) Paragraph {
	p.text = p.text.Bg(color)
	return p
}

func (p Paragraph) Bold() Paragraph {
	p.text = p.text.Bold()
	return p
}

func (p Paragraph) Italic() Paragraph {
	p.text = p.text.Italic()
	return p
}

func (p Paragraph) Cyan() Paragraph {
	return p.Fg(style.Cyan)
}

func (p Paragraph) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
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
		if p.scrollX > 0 && p.alignment == layout.Left {
			line = line.skip(p.scrollX)
		}
		offset := alignedOffset(line.width(), textArea.Width, p.alignment)
		x := textArea.X + offset
		for _, cell := range line.cells {
			if x >= textArea.X+textArea.Width {
				break
			}
			buf.SetCell(x, textArea.Y+y, cell)
			x++
		}
	}
}

func (p Paragraph) renderLines(width int) []renderLine {
	if width <= 0 {
		return nil
	}
	var lines []renderLine
	for _, line := range p.text.Lines {
		alignment := p.alignment
		if line.Alignment != nil {
			alignment = *line.Alignment
		}
		cells := cellsFromLine(line)
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

type Block struct {
	title   text.Line
	borders Borders
	style   style.Style
}

func NewBlock() Block {
	return Block{style: style.NewStyle()}
}

func BorderedBlock() Block {
	return NewBlock().Borders(AllBorders)
}

func (b Block) Title(title text.Line) Block {
	b.title = title
	return b
}

func (b Block) Borders(borders Borders) Block {
	b.borders = borders
	return b
}

func (b Block) Inner(area layout.Rect) layout.Rect {
	if b.borders == NoBorders {
		return area
	}
	return area.Inner(layout.NewMargin(1, 1))
}

func (b Block) Fg(color style.Color) Block {
	b.style = b.style.Fg(color)
	return b
}

func (b Block) Bg(color style.Color) Block {
	b.style = b.style.Bg(color)
	return b
}

func (b Block) Bold() Block {
	b.style = b.style.AddModifier(style.ModifierBold)
	return b
}

func (b Block) Italic() Block {
	b.style = b.style.AddModifier(style.ModifierItalic)
	return b
}

func (b Block) Cyan() Block {
	return b.Fg(style.Cyan)
}

func (b Block) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	if b.borders != NoBorders {
		b.renderBorders(area, buf)
	}
	titleX := area.X
	if b.borders != NoBorders {
		titleX++
	}
	x := titleX
	for _, span := range b.title.Spans {
		for _, r := range span.Content {
			if x >= area.X+area.Width {
				return
			}
			buf.SetCell(x, area.Y, buffer.Cell{Symbol: string(r), Style: b.style.Patch(span.Style)})
			x++
		}
	}
}

func (b Block) renderBorders(area layout.Rect, buf *buffer.Buffer) {
	right := area.X + area.Width - 1
	bottom := area.Y + area.Height - 1
	for x := area.X; x <= right; x++ {
		buf.SetCell(x, area.Y, buffer.Cell{Symbol: "─", Style: b.style})
		if bottom != area.Y {
			buf.SetCell(x, bottom, buffer.Cell{Symbol: "─", Style: b.style})
		}
	}
	for y := area.Y; y <= bottom; y++ {
		buf.SetCell(area.X, y, buffer.Cell{Symbol: "│", Style: b.style})
		if right != area.X {
			buf.SetCell(right, y, buffer.Cell{Symbol: "│", Style: b.style})
		}
	}
	buf.SetCell(area.X, area.Y, buffer.Cell{Symbol: "┌", Style: b.style})
	if right != area.X {
		buf.SetCell(right, area.Y, buffer.Cell{Symbol: "┐", Style: b.style})
	}
	if bottom != area.Y {
		buf.SetCell(area.X, bottom, buffer.Cell{Symbol: "└", Style: b.style})
		if right != area.X {
			buf.SetCell(right, bottom, buffer.Cell{Symbol: "┘", Style: b.style})
		}
	}
}

type Clear struct{}

func (Clear) Render(area layout.Rect, buf *buffer.Buffer) {
	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < area.X+area.Width; x++ {
			buf.SetCell(x, y, buffer.Cell{Symbol: " ", Style: style.NewStyle()})
		}
	}
}

type Borders uint8

const (
	NoBorders  Borders = 0
	AllBorders Borders = 1
)

type renderLine struct {
	cells     []buffer.Cell
	alignment layout.Alignment
}

func (l renderLine) width() int {
	return len(l.cells)
}

func (l renderLine) skip(count int) renderLine {
	if count >= len(l.cells) {
		return renderLine{alignment: l.alignment}
	}
	l.cells = l.cells[count:]
	return l
}

func cellsFromLine(line text.Line) []buffer.Cell {
	var cells []buffer.Cell
	for _, span := range line.Spans {
		for _, r := range span.Content {
			cells = append(cells, buffer.Cell{Symbol: string(r), Style: span.Style})
		}
	}
	return cells
}

func wrapCells(cells []buffer.Cell, width int, trim bool) [][]buffer.Cell {
	var lines [][]buffer.Cell
	for len(cells) > 0 {
		if trim {
			cells = trimLeftCells(cells)
		}
		if len(cells) <= width {
			lines = append(lines, trimRightCells(append([]buffer.Cell(nil), cells...), trim))
			break
		}
		breakAt := width
		for i := width; i >= 0; i-- {
			if i < len(cells) && isSpaceCell(cells[i]) {
				breakAt = i
				break
			}
		}
		if breakAt == 0 {
			breakAt = width
		}
		line := append([]buffer.Cell(nil), cells[:breakAt]...)
		lines = append(lines, trimRightCells(line, trim))
		cells = cells[breakAt:]
	}
	if len(lines) == 0 {
		lines = append(lines, nil)
	}
	return lines
}

func trimLeftCells(cells []buffer.Cell) []buffer.Cell {
	for len(cells) > 0 && isSpaceCell(cells[0]) {
		cells = cells[1:]
	}
	return cells
}

func trimRightCells(cells []buffer.Cell, trim bool) []buffer.Cell {
	if !trim {
		return cells
	}
	for len(cells) > 0 && isSpaceCell(cells[len(cells)-1]) {
		cells = cells[:len(cells)-1]
	}
	return cells
}

func isSpaceCell(cell buffer.Cell) bool {
	return strings.TrimFunc(cell.Symbol, unicode.IsSpace) == ""
}

func alignedOffset(lineWidth, areaWidth int, alignment layout.Alignment) int {
	if lineWidth >= areaWidth {
		return 0
	}
	switch alignment {
	case layout.Center:
		return (areaWidth - lineWidth) / 2
	case layout.Right:
		return areaWidth - lineWidth
	default:
		return 0
	}
}
