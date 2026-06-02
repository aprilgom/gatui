package widgets

import (
	"strings"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type ListItem struct {
	content text.Text
}

func NewListItem(content text.Text) ListItem {
	return ListItem{content: content}
}

func ListItemFromString(content string) ListItem {
	return NewListItem(text.FromString(content))
}

func ListItemFromLines(lines ...text.Line) ListItem {
	return NewListItem(text.NewText(lines...))
}

func (i ListItem) height() int {
	height := len(i.content.Lines)
	if height == 0 {
		return 1
	}
	return height
}

type ListState struct {
	offset   int
	selected *int
}

func (s *ListState) Select(index int) {
	s.selected = &index
}

func (s *ListState) ClearSelection() {
	s.selected = nil
	s.offset = 0
}

func (s ListState) Selected() (int, bool) {
	if s.selected == nil {
		return 0, false
	}
	return *s.selected, true
}

func (s ListState) Offset() int {
	return s.offset
}

func (s *ListState) SetOffset(offset int) {
	if offset < 0 {
		offset = 0
	}
	s.offset = offset
}

type HighlightSpacing int

const (
	HighlightSpacingWhenSelected HighlightSpacing = iota
	HighlightSpacingAlways
	HighlightSpacingNever
)

func (h HighlightSpacing) shouldAdd(selected bool) bool {
	switch h {
	case HighlightSpacingAlways:
		return true
	case HighlightSpacingNever:
		return false
	default:
		return selected
	}
}

type List struct {
	items                 []ListItem
	block                 *Block
	style                 style.Style
	highlightStyle        style.Style
	highlightSymbol       string
	repeatHighlightSymbol bool
	highlightSpacing      HighlightSpacing
}

func NewList(items []ListItem) List {
	return List{
		items:            append([]ListItem(nil), items...),
		style:            style.NewStyle(),
		highlightStyle:   style.NewStyle(),
		highlightSpacing: HighlightSpacingWhenSelected,
	}
}

func (l List) Len() int {
	return len(l.items)
}

func (l List) IsEmpty() bool {
	return len(l.items) == 0
}

func (l List) Block(block Block) List {
	l.block = &block
	return l
}

func (l List) Style(listStyle style.Style) List {
	l.style = listStyle
	return l
}

func (l List) HighlightStyle(highlightStyle style.Style) List {
	l.highlightStyle = highlightStyle
	return l
}

func (l List) HighlightSymbol(symbol string) List {
	l.highlightSymbol = symbol
	return l
}

func (l List) RepeatHighlightSymbol(repeat bool) List {
	l.repeatHighlightSymbol = repeat
	return l
}

func (l List) HighlightSpacing(spacing HighlightSpacing) List {
	l.highlightSpacing = spacing
	return l
}

func (l List) Render(area layout.Rect, buf *buffer.Buffer) {
	state := ListState{}
	l.RenderStateful(area, buf, &state)
}

func (l List) RenderStateful(area layout.Rect, buf *buffer.Buffer, state *ListState) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, l.style)
	listArea := area
	if l.block != nil {
		l.block.Render(area, buf)
		listArea = l.block.Inner(area)
	}
	if listArea.Width == 0 || listArea.Height == 0 {
		return
	}
	if state == nil {
		state = &ListState{}
	}
	if len(l.items) == 0 {
		state.ClearSelection()
		return
	}
	if state.offset >= len(l.items) {
		state.offset = len(l.items) - 1
	}
	if state.offset < 0 {
		state.offset = 0
	}
	if state.selected != nil {
		selected := *state.selected
		if selected < 0 {
			selected = 0
		}
		if selected >= len(l.items) {
			selected = len(l.items) - 1
		}
		state.selected = &selected
	}

	first, last := l.visibleBounds(state, listArea.Height)
	state.offset = first

	symbolWidth := len([]rune(l.highlightSymbol))
	addSpacing := l.highlightSpacing.shouldAdd(state.selected != nil)
	currentY := listArea.Y
	for index := first; index < last && currentY < listArea.Y+listArea.Height; index++ {
		item := l.items[index]
		itemHeight := item.height()
		rowHeight := itemHeight
		if currentY+rowHeight > listArea.Y+listArea.Height {
			rowHeight = listArea.Y + listArea.Height - currentY
		}
		rowArea := layout.NewRect(listArea.X, currentY, listArea.Width, rowHeight)
		isSelected := state.selected != nil && *state.selected == index
		itemArea := rowArea
		if addSpacing {
			itemArea.X += symbolWidth
			if itemArea.Width >= symbolWidth {
				itemArea.Width -= symbolWidth
			} else {
				itemArea.Width = 0
			}
		}
		l.renderItem(item, itemArea, buf)
		if isSelected {
			buf.SetStyle(rowArea, l.highlightStyle)
		}
		if addSpacing && symbolWidth > 0 {
			for line := 0; line < rowHeight; line++ {
				symbol := strings.Repeat(" ", symbolWidth)
				symbolStyle := l.style
				if isSelected && (line == 0 || l.repeatHighlightSymbol) {
					symbol = l.highlightSymbol
					symbolStyle = l.style.Patch(l.highlightStyle)
				} else if isSelected {
					symbolStyle = l.style.Patch(l.highlightStyle)
				}
				writeString(buf, listArea.X, currentY+line, symbol, symbolWidth, symbolStyle)
			}
		}
		currentY += itemHeight
	}
}

func (l List) RenderStatefulRef(area layout.Rect, buf *buffer.Buffer, state any) {
	if state == nil {
		l.RenderStateful(area, buf, nil)
		return
	}
	listState, ok := state.(*ListState)
	if !ok {
		panic("gatui: invalid state type for List")
	}
	l.RenderStateful(area, buf, listState)
}

func (l List) visibleBounds(state *ListState, height int) (int, int) {
	if height <= 0 || len(l.items) == 0 {
		return 0, 0
	}
	offset := state.offset
	if offset > len(l.items)-1 {
		offset = len(l.items) - 1
	}
	if offset < 0 {
		offset = 0
	}
	first := offset
	last := offset
	usedHeight := 0
	for last < len(l.items) {
		itemHeight := l.items[last].height()
		if usedHeight+itemHeight > height {
			break
		}
		usedHeight += itemHeight
		last++
	}
	if last == first {
		last = first + 1
		usedHeight = l.items[first].height()
	}
	indexToDisplay := offset
	if state.selected != nil {
		indexToDisplay = *state.selected
	}
	for indexToDisplay >= last && last < len(l.items) {
		usedHeight += l.items[last].height()
		last++
		for usedHeight > height && first < last {
			usedHeight -= l.items[first].height()
			first++
		}
	}
	for indexToDisplay < first && first > 0 {
		first--
		usedHeight += l.items[first].height()
		for usedHeight > height && last > first {
			last--
			usedHeight -= l.items[last].height()
		}
	}
	return first, last
}

func (l List) renderItem(item ListItem, area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	lines := item.content.Lines
	if len(lines) == 0 {
		lines = []text.Line{text.LineFromString("")}
	}
	for y := 0; y < area.Height && y < len(lines); y++ {
		renderLine(layout.NewRect(area.X, area.Y+y, area.Width, 1), buf, lines[y], l.style)
	}
}
