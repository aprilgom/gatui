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
	style   style.Style
}

func NewListItem(content text.Text) ListItem {
	return ListItem{content: content, style: style.NewStyle()}
}

func ListItemFromString(content string) ListItem {
	return NewListItem(text.FromString(content))
}

func ListItemFromLines(lines ...text.Line) ListItem {
	return NewListItem(text.NewText(lines...))
}

func ListItemFromSpan(span text.Span) ListItem {
	return NewListItem(text.TextFromSpan(span))
}

func ListItemFromLine(line text.Line) ListItem {
	return NewListItem(text.TextFromLine(line))
}

func ListItemFromSpans(spans ...text.Span) ListItem {
	return NewListItem(text.NewText(text.LineFromSpans(spans...)))
}

func ListItemFromText(content text.Text) ListItem {
	return NewListItem(content)
}

func (i ListItem) Style(itemStyle style.Style) ListItem {
	i.style = itemStyle
	return i
}

func (i ListItem) Fg(color style.Color) ListItem {
	i.style = i.style.Fg(color)
	return i
}

func (i ListItem) Bg(color style.Color) ListItem {
	i.style = i.style.Bg(color)
	return i
}

func (i ListItem) Bold() ListItem {
	i.style = i.style.AddModifier(style.ModifierBold)
	return i
}

func (i ListItem) Dim() ListItem {
	i.style = i.style.AddModifier(style.ModifierDim)
	return i
}

func (i ListItem) Italic() ListItem {
	i.style = i.style.AddModifier(style.ModifierItalic)
	return i
}

func (i ListItem) Cyan() ListItem {
	return i.Fg(style.Cyan)
}

func (i ListItem) Height() int {
	height := i.content.Height()
	if height == 0 {
		return 1
	}
	return height
}

func (i ListItem) Width() int {
	return i.content.Width()
}

func (i ListItem) height() int {
	return i.Height()
}

type ListState struct {
	offset   int
	selected *int
}

func NewListState() ListState {
	return ListState{}
}

func (s ListState) WithOffset(offset int) ListState {
	s.SetOffset(offset)
	return s
}

func (s ListState) WithSelected(index int) ListState {
	s.Select(index)
	return s
}

func (s ListState) WithoutSelected() ListState {
	s.ClearSelection()
	return s
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

func (s *ListState) SelectNext() {
	selected := 0
	if s.selected != nil {
		selected = saturatingAdd(*s.selected, 1)
	}
	s.Select(selected)
}

func (s *ListState) SelectPrevious() {
	selected := maxIntValue
	if s.selected != nil {
		selected = saturatingSub(*s.selected, 1)
	}
	s.Select(selected)
}

func (s *ListState) SelectFirst() {
	s.Select(0)
}

func (s *ListState) SelectLast() {
	s.Select(maxIntValue)
}

func (s *ListState) ScrollDownBy(amount int) {
	if amount < 0 {
		amount = 0
	}
	selected := 0
	if s.selected != nil {
		selected = *s.selected
	}
	s.Select(saturatingAdd(selected, amount))
}

func (s *ListState) ScrollUpBy(amount int) {
	if amount < 0 {
		amount = 0
	}
	selected := 0
	if s.selected != nil {
		selected = *s.selected
	}
	s.Select(saturatingSub(selected, amount))
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
	alignment             *layout.Alignment
}

func NewList(items []ListItem) List {
	return List{
		items:            append([]ListItem(nil), items...),
		style:            style.NewStyle(),
		highlightStyle:   style.NewStyle(),
		highlightSpacing: HighlightSpacingWhenSelected,
	}
}

func NewListFromStrings(items []string) List {
	listItems := make([]ListItem, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, ListItemFromString(item))
	}
	return NewList(listItems)
}

func NewListFromLines(items []text.Line) List {
	listItems := make([]ListItem, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, ListItemFromLine(item))
	}
	return NewList(listItems)
}

func NewListFromText(items []text.Text) List {
	listItems := make([]ListItem, 0, len(items))
	for _, item := range items {
		listItems = append(listItems, ListItemFromText(item))
	}
	return NewList(listItems)
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

func (l List) Fg(color style.Color) List {
	l.style = l.style.Fg(color)
	return l
}

func (l List) Bg(color style.Color) List {
	l.style = l.style.Bg(color)
	return l
}

func (l List) Bold() List {
	l.style = l.style.AddModifier(style.ModifierBold)
	return l
}

func (l List) Dim() List {
	l.style = l.style.AddModifier(style.ModifierDim)
	return l
}

func (l List) Italic() List {
	l.style = l.style.AddModifier(style.ModifierItalic)
	return l
}

func (l List) Cyan() List {
	return l.Fg(style.Cyan)
}

func (l List) Alignment(alignment layout.Alignment) List {
	l.alignment = &alignment
	return l
}

func (l List) Left() List {
	return l.Alignment(layout.Left)
}

func (l List) Center() List {
	return l.Alignment(layout.Center)
}

func (l List) Right() List {
	return l.Alignment(layout.Right)
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
		selected := max(*state.selected, 0)
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
	offset := max(min(state.offset, len(l.items)-1), 0)
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
	if last <= first {
		if indexToDisplay >= 0 && indexToDisplay < len(l.items) {
			return indexToDisplay, indexToDisplay + 1
		}
		return first, min(first+1, len(l.items))
	}
	return first, last
}

func (l List) renderItem(item ListItem, area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	baseStyle := l.style.Patch(item.style)
	buf.SetStyle(area, baseStyle)
	lines := item.content.Lines
	if len(lines) == 0 {
		lines = []text.Line{text.LineFromString("")}
	}
	for y := 0; y < area.Height && y < len(lines); y++ {
		alignment := lines[y].Alignment
		if alignment == nil {
			alignment = item.content.Alignment
		}
		if alignment == nil {
			alignment = l.alignment
		}
		line := lines[y].PatchStyle(item.content.Style)
		line.RenderWithAlignment(layout.NewRect(area.X, area.Y+y, area.Width, 1), buf, alignment)
	}
}
