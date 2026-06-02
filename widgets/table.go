package widgets

import (
	"strings"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type TableCell struct {
	content    text.Text
	style      style.Style
	columnSpan int
}

type TableRow struct {
	cells        []TableCell
	height       int
	topMargin    int
	bottomMargin int
	style        style.Style
}

type Table struct {
	rows                 []TableRow
	widths               []layout.Constraint
	header               *TableRow
	footer               *TableRow
	block                *Block
	columnSpacing        int
	style                style.Style
	rowHighlightStyle    style.Style
	columnHighlightStyle style.Style
	cellHighlightStyle   style.Style
	highlightSymbol      string
	highlightSpacing     HighlightSpacing
}

func NewTableCell(content text.Text) TableCell {
	return TableCell{content: content, style: style.NewStyle(), columnSpan: 1}
}

func TableCellFromString(content string) TableCell {
	return NewTableCell(text.FromString(content))
}

func (c TableCell) Style(cellStyle style.Style) TableCell {
	c.style = cellStyle
	return c
}

func (c TableCell) Content() text.Text {
	return c.content
}

func (c TableCell) Fg(color style.Color) TableCell {
	c.style = c.style.Fg(color)
	return c
}

func (c TableCell) Bg(color style.Color) TableCell {
	c.style = c.style.Bg(color)
	return c
}

func (c TableCell) Bold() TableCell {
	c.style = c.style.AddModifier(style.ModifierBold)
	return c
}

func (c TableCell) Dim() TableCell {
	c.style = c.style.AddModifier(style.ModifierDim)
	return c
}

func (c TableCell) Italic() TableCell {
	c.style = c.style.AddModifier(style.ModifierItalic)
	return c
}

func (c TableCell) Cyan() TableCell {
	return c.Fg(style.Cyan)
}

func (c TableCell) ColumnSpan(span int) TableCell {
	c.columnSpan = span
	return c
}

func NewTableRow(cells []TableCell) TableRow {
	return TableRow{cells: append([]TableCell(nil), cells...), height: 1, style: style.NewStyle()}
}

func TableRowFromStrings(values []string) TableRow {
	cells := make([]TableCell, 0, len(values))
	for _, value := range values {
		cells = append(cells, TableCellFromString(value))
	}
	return NewTableRow(cells)
}

func (r TableRow) Height(height int) TableRow {
	r.height = maxInt(0, height)
	return r
}

func (r TableRow) TopMargin(margin int) TableRow {
	r.topMargin = maxInt(0, margin)
	return r
}

func (r TableRow) BottomMargin(margin int) TableRow {
	r.bottomMargin = maxInt(0, margin)
	return r
}

func (r TableRow) Style(rowStyle style.Style) TableRow {
	r.style = rowStyle
	return r
}

type TableState struct {
	offset         int
	selected       *int
	selectedColumn *int
	selectedCell   *tableCellSelection
}

type tableCellSelection struct {
	row    int
	column int
}

const tableStateMaxIndex = int(^uint(0) >> 1)

func NewTableState() TableState {
	return TableState{}
}

func (s TableState) WithOffset(offset int) TableState {
	s.SetOffset(offset)
	return s
}

func (s TableState) WithSelected(index int) TableState {
	s.Select(index)
	return s
}

func (s TableState) WithSelectedColumn(index int) TableState {
	s.SelectColumn(index)
	return s
}

func (s TableState) WithSelectedCell(row, column int) TableState {
	s.SelectCell(row, column)
	return s
}

func (s *TableState) Select(index int) {
	s.selected = &index
}

func (s *TableState) ClearSelection() {
	s.selected = nil
	s.offset = 0
}

func (s TableState) Selected() (int, bool) {
	if s.selected == nil {
		return 0, false
	}
	return *s.selected, true
}

func (s *TableState) SelectColumn(index int) {
	s.selectedColumn = &index
}

func (s *TableState) ClearColumnSelection() {
	s.selectedColumn = nil
}

func (s TableState) SelectedColumn() (int, bool) {
	if s.selectedColumn == nil {
		return 0, false
	}
	return *s.selectedColumn, true
}

func (s *TableState) SelectCell(row, column int) {
	s.Select(row)
	s.SelectColumn(column)
	s.selectedCell = &tableCellSelection{row: row, column: column}
}

func (s *TableState) ClearCellSelection() {
	s.ClearSelection()
	s.ClearColumnSelection()
	s.selectedCell = nil
}

func (s TableState) SelectedCell() (row int, column int, ok bool) {
	if s.selectedCell == nil {
		return 0, 0, false
	}
	return s.selectedCell.row, s.selectedCell.column, true
}

func (s TableState) Offset() int {
	return s.offset
}

func (s *TableState) SetOffset(offset int) {
	if offset < 0 {
		offset = 0
	}
	s.offset = offset
}

func (s *TableState) SelectNext() {
	s.ScrollDownBy(1)
}

func (s *TableState) SelectPrevious() {
	s.ScrollUpBy(1)
}

func (s *TableState) SelectFirst() {
	s.Select(0)
}

func (s *TableState) SelectLast() {
	s.Select(tableStateMaxIndex)
}

func (s *TableState) SelectNextColumn() {
	s.ScrollRightBy(1)
}

func (s *TableState) SelectPreviousColumn() {
	s.ScrollLeftBy(1)
}

func (s *TableState) SelectFirstColumn() {
	s.SelectColumn(0)
}

func (s *TableState) SelectLastColumn() {
	s.SelectColumn(tableStateMaxIndex)
}

func (s *TableState) ScrollDownBy(amount int) {
	if amount <= 0 {
		return
	}
	selected := 0
	if s.selected != nil {
		selected = *s.selected
	}
	s.Select(saturatingAddInt(selected, amount))
}

func (s *TableState) ScrollUpBy(amount int) {
	if amount <= 0 {
		return
	}
	selected := 0
	if s.selected != nil {
		selected = *s.selected
	}
	s.Select(saturatingSubInt(selected, amount))
}

func (s *TableState) ScrollRightBy(amount int) {
	if amount <= 0 {
		return
	}
	selected := 0
	if s.selectedColumn != nil {
		selected = *s.selectedColumn
	}
	s.SelectColumn(saturatingAddInt(selected, amount))
}

func (s *TableState) ScrollLeftBy(amount int) {
	if amount <= 0 {
		return
	}
	selected := 0
	if s.selectedColumn != nil {
		selected = *s.selectedColumn
	}
	s.SelectColumn(saturatingSubInt(selected, amount))
}

func saturatingAddInt(value, amount int) int {
	if amount > tableStateMaxIndex-value {
		return tableStateMaxIndex
	}
	return value + amount
}

func saturatingSubInt(value, amount int) int {
	if amount >= value {
		return 0
	}
	return value - amount
}

func NewTable(rows []TableRow, widths []layout.Constraint) Table {
	return Table{
		rows:                 append([]TableRow(nil), rows...),
		widths:               append([]layout.Constraint(nil), widths...),
		columnSpacing:        1,
		style:                style.NewStyle(),
		rowHighlightStyle:    style.NewStyle(),
		columnHighlightStyle: style.NewStyle(),
		cellHighlightStyle:   style.NewStyle(),
		highlightSpacing:     HighlightSpacingWhenSelected,
	}
}

func (t Table) Header(header TableRow) Table {
	t.header = &header
	return t
}

func (t Table) Footer(footer TableRow) Table {
	t.footer = &footer
	return t
}

func (t Table) Block(block Block) Table {
	t.block = &block
	return t
}

func (t Table) ColumnSpacing(spacing int) Table {
	t.columnSpacing = maxInt(0, spacing)
	return t
}

func (t Table) Rows(rows []TableRow) Table {
	t.rows = append([]TableRow(nil), rows...)
	return t
}

func (t Table) Style(tableStyle style.Style) Table {
	t.style = tableStyle
	return t
}

func (t Table) Fg(color style.Color) Table {
	t.style = t.style.Fg(color)
	return t
}

func (t Table) Bg(color style.Color) Table {
	t.style = t.style.Bg(color)
	return t
}

func (t Table) Bold() Table {
	t.style = t.style.AddModifier(style.ModifierBold)
	return t
}

func (t Table) Dim() Table {
	t.style = t.style.AddModifier(style.ModifierDim)
	return t
}

func (t Table) Italic() Table {
	t.style = t.style.AddModifier(style.ModifierItalic)
	return t
}

func (t Table) Cyan() Table {
	return t.Fg(style.Cyan)
}

func (t Table) HighlightSymbol(symbol string) Table {
	t.highlightSymbol = symbol
	return t
}

func (t Table) HighlightSpacing(spacing HighlightSpacing) Table {
	t.highlightSpacing = spacing
	return t
}

func (t Table) RowHighlightStyle(rowHighlightStyle style.Style) Table {
	t.rowHighlightStyle = rowHighlightStyle
	return t
}

func (t Table) ColumnHighlightStyle(columnHighlightStyle style.Style) Table {
	t.columnHighlightStyle = columnHighlightStyle
	return t
}

func (t Table) CellHighlightStyle(cellHighlightStyle style.Style) Table {
	t.cellHighlightStyle = cellHighlightStyle
	return t
}

func (t Table) Render(area layout.Rect, buf *buffer.Buffer) {
	state := TableState{}
	t.RenderStateful(area, buf, &state)
}

func (t Table) RenderStateful(area layout.Rect, buf *buffer.Buffer, state *TableState) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, t.style)
	tableArea := area
	if t.block != nil {
		t.block.Render(area, buf)
		tableArea = t.block.Inner(area)
	}
	if tableArea.Width == 0 || tableArea.Height == 0 {
		return
	}
	if state == nil {
		state = &TableState{}
	}
	t.clampState(state)

	selected := state.selected != nil
	addSpacing := t.highlightSpacing.shouldAdd(selected)
	symbolWidth := len([]rune(t.highlightSymbol))
	rowArea := tableArea
	if addSpacing {
		rowArea.X += symbolWidth
		if rowArea.Width >= symbolWidth {
			rowArea.Width -= symbolWidth
		} else {
			rowArea.Width = 0
		}
	}
	widths := t.resolveColumnWidths(rowArea.Width)
	y := tableArea.Y
	if t.header != nil {
		y = t.renderRow(*t.header, widths, rowArea, y, -1, state, buf)
		y += t.header.bottomMargin
	}

	footerHeight := 0
	footerY := tableArea.Y + tableArea.Height
	if t.footer != nil {
		footerHeight = t.footer.topMargin + normalizedRowHeight(*t.footer) + t.footer.bottomMargin
		footerY -= footerHeight
	}
	bodyHeight := max(tableArea.Y+tableArea.Height-y-footerHeight, 0)
	first, last := t.visibleBounds(state, bodyHeight)
	state.offset = first
	for index := first; index < last; index++ {
		if y >= tableArea.Y+tableArea.Height {
			return
		}
		row := t.rows[index]
		rowHeight := normalizedRowHeight(row)
		drawHeight := minInt(rowHeight, tableArea.Y+tableArea.Height-y)
		baseRowArea := layout.NewRect(tableArea.X, y, tableArea.Width, drawHeight)
		isSelected := state.selected != nil && *state.selected == index
		y = t.renderRow(row, widths, rowArea, y, index, state, buf)
		if isSelected {
			buf.SetStyle(baseRowArea, t.rowHighlightStyle)
			t.renderRow(row, widths, rowArea, baseRowArea.Y, index, state, buf)
		}
		if addSpacing && symbolWidth > 0 {
			symbol := strings.Repeat(" ", symbolWidth)
			symbolStyle := t.style
			if isSelected {
				symbol = t.highlightSymbol
				symbolStyle = t.style.Patch(t.rowHighlightStyle)
			}
			writeString(buf, tableArea.X, baseRowArea.Y, symbol, symbolWidth, symbolStyle)
			for line := 1; line < drawHeight; line++ {
				blankStyle := t.style
				if isSelected {
					blankStyle = t.style.Patch(t.rowHighlightStyle)
				}
				writeString(buf, tableArea.X, baseRowArea.Y+line, strings.Repeat(" ", symbolWidth), symbolWidth, blankStyle)
			}
		}
		y += row.bottomMargin
	}
	if t.footer != nil {
		t.renderRow(*t.footer, widths, rowArea, footerY, -1, state, buf)
	}
}

func (t Table) RenderStatefulRef(area layout.Rect, buf *buffer.Buffer, state any) {
	if state == nil {
		t.RenderStateful(area, buf, nil)
		return
	}
	tableState, ok := state.(*TableState)
	if !ok {
		panic("gatui: invalid state type for Table")
	}
	t.RenderStateful(area, buf, tableState)
}

func (t Table) clampState(state *TableState) {
	if len(t.rows) == 0 {
		state.ClearSelection()
		state.ClearColumnSelection()
		state.ClearCellSelection()
		return
	}
	if state.offset >= len(t.rows) {
		state.offset = len(t.rows) - 1
	}
	if state.offset < 0 {
		state.offset = 0
	}
	if state.selected != nil {
		selected := max(*state.selected, 0)
		if selected >= len(t.rows) {
			selected = len(t.rows) - 1
		}
		state.selected = &selected
	}
	columnCount := t.columnCount()
	if columnCount == 0 {
		state.ClearColumnSelection()
		state.ClearCellSelection()
		return
	}
	if state.selectedColumn != nil {
		selectedColumn := max(*state.selectedColumn, 0)
		if selectedColumn >= columnCount {
			selectedColumn = columnCount - 1
		}
		state.selectedColumn = &selectedColumn
	}
	if state.selectedCell != nil {
		selectedCell := *state.selectedCell
		if selectedCell.row < 0 {
			selectedCell.row = 0
		}
		if selectedCell.row >= len(t.rows) {
			selectedCell.row = len(t.rows) - 1
		}
		if selectedCell.column < 0 {
			selectedCell.column = 0
		}
		if selectedCell.column >= columnCount {
			selectedCell.column = columnCount - 1
		}
		state.selectedCell = &selectedCell
	}
}

func (t Table) columnCount() int {
	columnCount := len(t.widths)
	for _, row := range t.rows {
		rowColumns := 0
		for _, cell := range row.cells {
			if cell.columnSpan <= 0 {
				continue
			}
			rowColumns += cell.columnSpan
		}
		columnCount = maxInt(columnCount, rowColumns)
	}
	if t.header != nil {
		columnCount = maxInt(columnCount, rowColumnCount(*t.header))
	}
	if t.footer != nil {
		columnCount = maxInt(columnCount, rowColumnCount(*t.footer))
	}
	return columnCount
}

func rowColumnCount(row TableRow) int {
	columnCount := 0
	for _, cell := range row.cells {
		if cell.columnSpan <= 0 {
			continue
		}
		columnCount += cell.columnSpan
	}
	return columnCount
}

func (t Table) visibleBounds(state *TableState, height int) (int, int) {
	if height <= 0 || len(t.rows) == 0 {
		return 0, 0
	}
	offset := max(min(state.offset, len(t.rows)-1), 0)
	first := offset
	last := offset
	usedHeight := 0
	for last < len(t.rows) {
		rowHeight := normalizedRowHeight(t.rows[last])
		if usedHeight+rowHeight > height {
			break
		}
		usedHeight += rowHeight
		last++
	}
	if last == first {
		last = first + 1
		usedHeight = normalizedRowHeight(t.rows[first])
	}
	indexToDisplay := offset
	if state.selected != nil {
		indexToDisplay = *state.selected
	}
	for indexToDisplay >= last && last < len(t.rows) {
		usedHeight += normalizedRowHeight(t.rows[last])
		last++
		for usedHeight > height && first < last {
			usedHeight -= normalizedRowHeight(t.rows[first])
			first++
		}
	}
	for indexToDisplay < first && first > 0 {
		first--
		usedHeight += normalizedRowHeight(t.rows[first])
		for usedHeight > height && last > first {
			last--
			usedHeight -= normalizedRowHeight(t.rows[last])
		}
	}
	return first, last
}

func normalizedRowHeight(row TableRow) int {
	if row.height <= 0 {
		return 1
	}
	return row.height
}

func (t Table) resolveColumnWidths(width int) []int {
	widths := make([]int, len(t.widths))
	if width <= 0 || len(widths) == 0 {
		return widths
	}
	spacingTotal := t.columnSpacing * maxInt(0, len(widths)-1)
	available := maxInt(0, width-spacingTotal)
	hasLength := false
	hasPercentage := false
	percentageSum := 0
	for _, constraint := range t.widths {
		hasLength = hasLength || constraint.IsLength()
		hasPercentage = hasPercentage || constraint.IsPercentage()
		if constraint.IsPercentage() {
			percentageSum += constraint.Value()
		}
	}
	if hasLength && hasPercentage && percentageSum > 100 {
		return t.resolveOverspecifiedMixedWidths(available)
	}
	for i, constraint := range t.widths {
		switch {
		case constraint.IsLength():
			widths[i] = maxInt(0, constraint.Value())
		case constraint.IsPercentage():
			percentageBase := available
			if hasLength && hasPercentage {
				percentageBase = width
			}
			widths[i] = maxInt(0, percentageBase*constraint.Value()/100)
		case constraint.IsRatio():
			if constraint.Denominator() > 0 {
				widths[i] = maxInt(0, available*constraint.Value()/constraint.Denominator())
			}
		default:
			widths[i] = available
		}
	}
	t.distributeRatioRemainder(widths, available)
	t.fitWidths(widths, available, hasLength && !hasPercentage)
	return widths
}

func (t Table) resolveOverspecifiedMixedWidths(available int) []int {
	widths := make([]int, len(t.widths))
	fixedTotal := 0
	for i, constraint := range t.widths {
		if constraint.IsLength() {
			widths[i] = maxInt(0, constraint.Value())
			fixedTotal += widths[i]
		}
	}
	remaining := maxInt(0, available-fixedTotal)
	percentageLeft := 0
	for _, constraint := range t.widths {
		if constraint.IsPercentage() {
			percentageLeft++
		}
	}
	for i, constraint := range t.widths {
		if !constraint.IsPercentage() {
			continue
		}
		percentageLeft--
		percentageWidth := remaining
		if percentageLeft > 0 {
			percentageWidth = (remaining*constraint.Value() + 99) / 100
		}
		if percentageWidth > remaining {
			percentageWidth = remaining
		}
		widths[i] = percentageWidth
		remaining -= percentageWidth
	}
	t.fitWidths(widths, available, false)
	return widths
}

func (t Table) distributeRatioRemainder(widths []int, available int) {
	ratioIndexes := make([]int, 0, len(widths))
	total := 0
	ratioSum := 0.0
	for i, width := range widths {
		total += width
		if t.widths[i].IsRatio() && t.widths[i].Value() > 0 && t.widths[i].Denominator() > 0 {
			ratioIndexes = append(ratioIndexes, i)
			ratioSum += float64(t.widths[i].Value()) / float64(t.widths[i].Denominator())
		}
	}
	remainder := available - total
	for remainder > 0 && len(ratioIndexes) > 0 && ratioSum >= 1 {
		start := len(ratioIndexes) / 2
		for offset := range ratioIndexes {
			if remainder == 0 {
				return
			}
			index := ratioIndexes[(start+offset)%len(ratioIndexes)]
			widths[index]++
			remainder--
		}
	}
}

func (t Table) fitWidths(widths []int, available int, preferMiddle bool) {
	total := 0
	for _, width := range widths {
		total += width
	}
	for total > available {
		if !preferMiddle {
			for i := len(widths) - 1; i >= 0 && total > available; i-- {
				for widths[i] > 0 && total > available {
					widths[i]--
					total--
				}
			}
			continue
		}
		for _, i := range shrinkOrder(len(widths), preferMiddle) {
			if total <= available {
				break
			}
			if widths[i] > 0 {
				widths[i]--
				total--
			}
		}
	}
}

func shrinkOrder(length int, preferMiddle bool) []int {
	if length <= 0 {
		return nil
	}
	if !preferMiddle {
		order := make([]int, 0, length)
		for i := length - 1; i >= 0; i-- {
			order = append(order, i)
		}
		return order
	}
	order := make([]int, 0, length)
	middle := length / 2
	order = append(order, middle)
	for offset := 1; len(order) < length; offset++ {
		left := middle - offset
		right := middle + offset
		if left >= 0 {
			order = append(order, left)
		}
		if right < length {
			order = append(order, right)
		}
	}
	return order
}

func (t Table) renderRow(row TableRow, widths []int, area layout.Rect, y int, rowIndex int, state *TableState, buf *buffer.Buffer) int {
	y += row.topMargin
	rowHeight := row.height
	if rowHeight <= 0 {
		rowHeight = 1
	}
	for line := 0; line < rowHeight && y+line < area.Y+area.Height; line++ {
		physicalColumn := 0
		for _, cell := range row.cells {
			if physicalColumn >= len(widths) {
				break
			}
			span := cell.columnSpan
			if span <= 0 {
				continue
			}
			span = minInt(span, len(widths)-physicalColumn)
			x, width := t.cellSpanArea(widths, physicalColumn, span, area)
			if width > 0 {
				cellArea := layout.NewRect(x, y+line, width, 1)
				cellStyle := row.style.
					Patch(cell.style).
					Patch(t.cellHighlightForSpan(rowIndex, physicalColumn, span, state))
				buf.SetStyle(cellArea, cellStyle)
			}
			if line == 0 && width > 0 {
				t.renderCell(cell, row.style, t.cellHighlightForSpan(rowIndex, physicalColumn, span, state), x, y+line, area.X+area.Width, width, buf)
			}
			physicalColumn += span
		}
	}
	return y + rowHeight
}

func (t Table) cellSpanArea(widths []int, startColumn, span int, area layout.Rect) (int, int) {
	x := area.X
	for column := range startColumn {
		if column > 0 {
			x += t.columnSpacing
		}
		x += widths[column]
	}
	if startColumn > 0 {
		x += t.columnSpacing
	}
	width := 0
	for offset := 0; offset < span && startColumn+offset < len(widths); offset++ {
		if offset > 0 {
			width += t.columnSpacing
		}
		width += widths[startColumn+offset]
	}
	return x, minInt(width, area.X+area.Width-x)
}

func (t Table) cellHighlightForSpan(rowIndex, startColumn, span int, state *TableState) style.Style {
	cellStyle := style.NewStyle()
	if state == nil {
		return cellStyle
	}
	if rowIndex >= 0 && state.selected != nil && *state.selected == rowIndex {
		cellStyle = cellStyle.Patch(t.rowHighlightStyle)
	}
	endColumn := startColumn + span
	if state.selectedColumn != nil && *state.selectedColumn >= startColumn && *state.selectedColumn < endColumn {
		cellStyle = cellStyle.Patch(t.columnHighlightStyle)
	}
	if rowIndex >= 0 && state.selectedCell != nil &&
		state.selectedCell.row == rowIndex &&
		state.selectedCell.column >= startColumn &&
		state.selectedCell.column < endColumn {
		cellStyle = cellStyle.Patch(t.cellHighlightStyle)
	}
	return cellStyle
}

func (t Table) renderCell(cell TableCell, rowStyle, highlightStyle style.Style, x, y, right, width int, buf *buffer.Buffer) {
	if x >= right || width <= 0 {
		return
	}
	cellStyle := t.style.Patch(rowStyle).Patch(cell.style).Patch(highlightStyle)
	line := text.LineFromString("")
	var alignment *layout.Alignment
	if len(cell.content.Lines) > 0 {
		line = cell.content.Lines[0]
		alignment = cell.content.Alignment
	}
	renderLineAligned(layout.NewRect(x, y, minInt(width, right-x), 1), buf, line, cellStyle, alignment)
}
