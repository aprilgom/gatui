package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type Tabs struct {
	titles         []text.Line
	block          *Block
	style          style.Style
	highlightStyle style.Style
	selected       *int
	divider        string
	paddingLeft    string
	paddingRight   string
}

func NewTabs(titles []text.Line) Tabs {
	selected := 0
	return Tabs{
		titles:         append([]text.Line(nil), titles...),
		style:          style.NewStyle(),
		highlightStyle: style.NewStyle().AddModifier(style.ModifierReversed),
		selected:       &selected,
		divider:        "│",
		paddingLeft:    " ",
		paddingRight:   " ",
	}
}

func TabsFromStrings(titles []string) Tabs {
	lines := make([]text.Line, 0, len(titles))
	for _, title := range titles {
		lines = append(lines, text.LineFromString(title))
	}
	return NewTabs(lines)
}

func (t Tabs) Block(block Block) Tabs {
	t.block = &block
	return t
}

func (t Tabs) Select(index int) Tabs {
	t.selected = &index
	return t
}

func (t Tabs) ClearSelection() Tabs {
	t.selected = nil
	return t
}

func (t Tabs) Style(tabStyle style.Style) Tabs {
	t.style = tabStyle
	return t
}

func (t Tabs) Fg(color style.Color) Tabs {
	t.style = t.style.Fg(color)
	return t
}

func (t Tabs) Bg(color style.Color) Tabs {
	t.style = t.style.Bg(color)
	return t
}

func (t Tabs) Bold() Tabs {
	t.style = t.style.AddModifier(style.ModifierBold)
	return t
}

func (t Tabs) Dim() Tabs {
	t.style = t.style.AddModifier(style.ModifierDim)
	return t
}

func (t Tabs) Italic() Tabs {
	t.style = t.style.AddModifier(style.ModifierItalic)
	return t
}

func (t Tabs) Cyan() Tabs {
	return t.Fg(style.Cyan)
}

func (t Tabs) HighlightStyle(highlightStyle style.Style) Tabs {
	t.highlightStyle = highlightStyle
	return t
}

func (t Tabs) Divider(divider string) Tabs {
	t.divider = divider
	return t
}

func (t Tabs) Padding(left, right string) Tabs {
	t.paddingLeft = left
	t.paddingRight = right
	return t
}

func (t Tabs) PaddingLeft(left string) Tabs {
	t.paddingLeft = left
	return t
}

func (t Tabs) PaddingRight(right string) Tabs {
	t.paddingRight = right
	return t
}

func (t Tabs) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, t.style)
	tabsArea := area
	if t.block != nil {
		t.block.Render(area, buf)
		tabsArea = t.block.Inner(area)
	}
	if tabsArea.Width == 0 || tabsArea.Height == 0 {
		return
	}
	x := tabsArea.X
	right := tabsArea.X + tabsArea.Width
	for index, title := range t.titles {
		if x >= right {
			return
		}
		x = writeStringWithin(buf, x, tabsArea.Y, right, t.paddingLeft, t.style)
		selected := t.selected != nil && *t.selected == index
		x = t.renderTitle(buf, title, x, tabsArea.Y, right, selected)
		x = writeStringWithin(buf, x, tabsArea.Y, right, t.paddingRight, t.style)
		if index < len(t.titles)-1 {
			x = writeStringWithin(buf, x, tabsArea.Y, right, t.divider, t.style)
		}
	}
}

func (t Tabs) renderTitle(buf *buffer.Buffer, title text.Line, x, y, right int, selected bool) int {
	for _, grapheme := range title.StyledGraphemes(t.style) {
		width := buffer.CellWidth(grapheme.Symbol)
		if width == 0 {
			continue
		}
		if x+width > right {
			return x
		}

		cellStyle := grapheme.Style
		if selected {
			cellStyle = cellStyle.Patch(t.highlightStyle)
		}
		buf.SetCell(x, y, buffer.Cell{Symbol: grapheme.Symbol, Style: cellStyle})
		for trailing := 1; trailing < width; trailing++ {
			buf.SetCell(x+trailing, y, buffer.Cell{Symbol: " ", Style: style.NewStyle()})
		}
		x += width
	}
	return x
}
