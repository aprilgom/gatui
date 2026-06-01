package widgets

import (
	"strconv"
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

type BarData struct {
	Label string
	Value uint64
}

type Bar struct {
	label      string
	value      uint64
	textValue  string
	style      style.Style
	valueStyle style.Style
	labelStyle style.Style
}

func NewBar(value uint64) Bar {
	return Bar{value: value, textValue: uintToString(value), style: style.NewStyle(), valueStyle: style.NewStyle(), labelStyle: style.NewStyle()}
}

func BarWithLabel(label string, value uint64) Bar {
	return NewBar(value).Label(label)
}

func (b Bar) Label(label string) Bar {
	b.label = label
	return b
}

func (b Bar) TextValue(value string) Bar {
	b.textValue = value
	return b
}

func (b Bar) Style(barStyle style.Style) Bar {
	b.style = barStyle
	return b
}

func (b Bar) ValueStyle(valueStyle style.Style) Bar {
	b.valueStyle = valueStyle
	return b
}

func (b Bar) LabelStyle(labelStyle style.Style) Bar {
	b.labelStyle = labelStyle
	return b
}

type BarGroup struct {
	label string
	bars  []Bar
}

func NewBarGroup(bars []Bar) BarGroup {
	return BarGroup{bars: append([]Bar(nil), bars...)}
}

func (g BarGroup) Label(label string) BarGroup {
	g.label = label
	return g
}

func (g BarGroup) Bars(bars []Bar) BarGroup {
	g.bars = append([]Bar(nil), bars...)
	return g
}

type BarChart struct {
	groups     []BarGroup
	block      *Block
	max        uint64
	barWidth   int
	barGap     int
	groupGap   int
	barStyle   style.Style
	valueStyle style.Style
	labelStyle style.Style
}

func NewBarChart() BarChart {
	return BarChart{
		barWidth:   1,
		barGap:     1,
		groupGap:   1,
		barStyle:   style.NewStyle(),
		valueStyle: style.NewStyle(),
		labelStyle: style.NewStyle(),
	}
}

func (c BarChart) Data(group BarGroup) BarChart {
	c.groups = append(c.groups, group)
	return c
}

func (c BarChart) DataPairs(data []BarData) BarChart {
	bars := make([]Bar, 0, len(data))
	for _, item := range data {
		bars = append(bars, BarWithLabel(item.Label, item.Value))
	}
	return c.Data(NewBarGroup(bars))
}

func (c BarChart) Block(block Block) BarChart {
	c.block = &block
	return c
}

func (c BarChart) Max(max uint64) BarChart {
	c.max = max
	return c
}

func (c BarChart) BarWidth(width int) BarChart {
	c.barWidth = width
	return c
}

func (c BarChart) BarGap(gap int) BarChart {
	c.barGap = gap
	return c
}

func (c BarChart) GroupGap(gap int) BarChart {
	c.groupGap = gap
	return c
}

func (c BarChart) BarStyle(barStyle style.Style) BarChart {
	c.barStyle = barStyle
	return c
}

func (c BarChart) ValueStyle(valueStyle style.Style) BarChart {
	c.valueStyle = valueStyle
	return c
}

func (c BarChart) LabelStyle(labelStyle style.Style) BarChart {
	c.labelStyle = labelStyle
	return c
}

func (c BarChart) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	chartArea := area
	if c.block != nil {
		c.block.Render(area, buf)
		chartArea = c.block.Inner(area)
	}
	if chartArea.Width == 0 || chartArea.Height == 0 || len(c.groups) == 0 || c.barWidth <= 0 {
		return
	}
	labelRows := 1
	if c.hasGroupLabels() {
		labelRows = 2
	}
	if chartArea.Height <= labelRows {
		return
	}
	max := c.effectiveMax()
	barHeight := chartArea.Height - labelRows
	barLabelY := chartArea.Y + barHeight
	groupLabelY := barLabelY
	if labelRows == 2 {
		groupLabelY = barLabelY + 1
	}
	x := chartArea.X
	right := chartArea.X + chartArea.Width
	for groupIndex, group := range c.groups {
		if groupIndex > 0 {
			x += nonNegative(c.groupGap) + nonNegative(c.barGap)
		}
		groupStart := x
		for barIndex, bar := range group.bars {
			if barIndex > 0 {
				x += nonNegative(c.barGap)
			}
			if x >= right {
				break
			}
			width := c.barWidth
			if x+width > right {
				width = right - x
			}
			c.renderBar(buf, x, chartArea.Y, width, barHeight, max, bar)
			c.renderCentered(buf, x, barLabelY, width, bar.label, c.labelStyle.Patch(bar.labelStyle))
			x += c.barWidth
		}
		if group.label != "" && groupStart < right {
			writeStringWithin(buf, groupStart, groupLabelY, right, group.label, c.labelStyle)
		}
	}
}

func (c BarChart) renderBar(buf *buffer.Buffer, x, y, width, height int, max uint64, bar Bar) {
	if width <= 0 || height <= 0 || max == 0 || bar.value == 0 {
		return
	}
	eighths := int((bar.value * uint64(height) * 8) / max)
	if eighths > height*8 {
		eighths = height * 8
	}
	barStyle := c.barStyle.Patch(bar.style)
	buf.SetStyle(layout.NewRect(x, y, width, height), barStyle)
	for rowFromBottom := 0; rowFromBottom < height; rowFromBottom++ {
		rowEighths := eighths - rowFromBottom*8
		if rowEighths <= 0 {
			continue
		}
		symbol := "█"
		if rowEighths < 8 {
			symbol = partialBarSymbol(rowEighths)
		}
		for dx := 0; dx < width; dx++ {
			buf.SetCell(x+dx, y+height-1-rowFromBottom, buffer.Cell{Symbol: symbol, Style: barStyle})
		}
	}
	value := bar.textValue
	if value == "" {
		value = uintToString(bar.value)
	}
	valueStyle := barStyle.Patch(c.valueStyle).Patch(bar.valueStyle)
	c.renderCentered(buf, x, y+height-1, width, value, valueStyle)
}

func (c BarChart) renderCentered(buf *buffer.Buffer, x, y, width int, value string, cellStyle style.Style) {
	runes := []rune(value)
	if len(runes) > width {
		runes = runes[:width]
	}
	offset := (width - len(runes)) / 2
	for i, r := range runes {
		buf.SetCell(x+offset+i, y, buffer.Cell{Symbol: string(r), Style: cellStyle})
	}
}

func (c BarChart) hasGroupLabels() bool {
	for _, group := range c.groups {
		if group.label != "" {
			return true
		}
	}
	return false
}

func (c BarChart) effectiveMax() uint64 {
	if c.max > 0 {
		return c.max
	}
	var max uint64
	for _, group := range c.groups {
		for _, bar := range group.bars {
			if bar.value > max {
				max = bar.value
			}
		}
	}
	return max
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
	for _, span := range title.Spans {
		cellStyle := t.style.Patch(span.Style)
		if selected {
			cellStyle = cellStyle.Patch(t.highlightStyle)
		}
		for _, r := range span.Content {
			if x >= right {
				return x
			}
			buf.SetCell(x, y, buffer.Cell{Symbol: string(r), Style: cellStyle})
			x++
		}
	}
	return x
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
		cells := cellsFromLine(lines[y])
		for x := 0; x < area.Width && x < len(cells); x++ {
			cell := cells[x]
			cell.Style = l.style.Patch(cell.Style)
			buf.SetCell(area.X+x, area.Y+y, cell)
		}
	}
}

func writeString(buf *buffer.Buffer, x, y int, value string, width int, cellStyle style.Style) {
	runes := []rune(value)
	for i := 0; i < width; i++ {
		symbol := " "
		if i < len(runes) {
			symbol = string(runes[i])
		}
		buf.SetCell(x+i, y, buffer.Cell{Symbol: symbol, Style: cellStyle})
	}
}

func writeStringWithin(buf *buffer.Buffer, x, y, right int, value string, cellStyle style.Style) int {
	for _, r := range value {
		if x >= right {
			return x
		}
		buf.SetCell(x, y, buffer.Cell{Symbol: string(r), Style: cellStyle})
		x++
	}
	return x
}

func partialBarSymbol(eighths int) string {
	switch eighths {
	case 1:
		return "▁"
	case 2:
		return "▂"
	case 3:
		return "▃"
	case 4:
		return "▄"
	case 5:
		return "▅"
	case 6:
		return "▆"
	case 7:
		return "▇"
	default:
		return " "
	}
}

func uintToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func nonNegative(value int) int {
	if value < 0 {
		return 0
	}
	return value
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
