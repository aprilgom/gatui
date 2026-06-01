package widgets

import (
	"fmt"
	"math"
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

type Gauge struct {
	block      *Block
	ratio      float64
	label      *text.Span
	style      style.Style
	gaugeStyle style.Style
	useUnicode bool
}

type LineGauge struct {
	block          *Block
	ratio          float64
	label          *text.Line
	style          style.Style
	filledStyle    style.Style
	unfilledStyle  style.Style
	filledSymbol   string
	unfilledSymbol string
}

type BarData struct {
	Label string
	Value uint64
}

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

type GraphType int

const (
	GraphTypeScatter GraphType = iota
	GraphTypeLine
)

type Axis struct {
	title           *text.Line
	bounds          [2]float64
	labels          []text.Line
	axisStyle       style.Style
	labelsAlignment layout.Alignment
}

type ChartPoint struct {
	X float64
	Y float64
}

type Dataset struct {
	name      string
	data      []ChartPoint
	graphType GraphType
	style     style.Style
}

type Chart struct {
	datasets []Dataset
	block    *Block
	xAxis    Axis
	yAxis    Axis
}

type Bar struct {
	label      string
	value      uint64
	textValue  string
	style      style.Style
	valueStyle style.Style
	labelStyle style.Style
}

func NewAxis() Axis {
	return Axis{
		bounds:          [2]float64{0, 0},
		axisStyle:       style.NewStyle(),
		labelsAlignment: layout.Left,
	}
}

func (a Axis) Title(title text.Line) Axis {
	a.title = &title
	return a
}

func (a Axis) TitleString(title string) Axis {
	return a.Title(text.LineFromString(title))
}

func (a Axis) Bounds(minimum, maximum float64) Axis {
	a.bounds = [2]float64{minimum, maximum}
	return a
}

func (a Axis) Labels(labels []text.Line) Axis {
	a.labels = append([]text.Line(nil), labels...)
	return a
}

func (a Axis) LabelStrings(labels []string) Axis {
	lines := make([]text.Line, 0, len(labels))
	for _, label := range labels {
		lines = append(lines, text.LineFromString(label))
	}
	return a.Labels(lines)
}

func (a Axis) Style(axisStyle style.Style) Axis {
	a.axisStyle = axisStyle
	return a
}

func (a Axis) LabelsAlignment(alignment layout.Alignment) Axis {
	a.labelsAlignment = alignment
	return a
}

func NewDataset() Dataset {
	return Dataset{graphType: GraphTypeScatter, style: style.NewStyle()}
}

func (d Dataset) Name(name string) Dataset {
	d.name = name
	return d
}

func (d Dataset) Data(data []layout.Position) Dataset {
	d.data = make([]ChartPoint, 0, len(data))
	for _, point := range data {
		d.data = append(d.data, ChartPoint{X: float64(point.X), Y: float64(point.Y)})
	}
	return d
}

func (d Dataset) DataPoints(data []ChartPoint) Dataset {
	d.data = append([]ChartPoint(nil), data...)
	return d
}

func (d Dataset) GraphType(graphType GraphType) Dataset {
	d.graphType = graphType
	return d
}

func (d Dataset) Style(datasetStyle style.Style) Dataset {
	d.style = datasetStyle
	return d
}

func NewChart(datasets []Dataset) Chart {
	return Chart{
		datasets: append([]Dataset(nil), datasets...),
		xAxis:    NewAxis(),
		yAxis:    NewAxis(),
	}
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

func (t Table) Style(tableStyle style.Style) Table {
	t.style = tableStyle
	return t
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
	bodyHeight := tableArea.Y + tableArea.Height - y - footerHeight
	if bodyHeight < 0 {
		bodyHeight = 0
	}
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
		selected := *state.selected
		if selected < 0 {
			selected = 0
		}
		if selected >= len(t.rows) {
			selected = len(t.rows) - 1
		}
		state.selected = &selected
	}
	columnCount := len(t.widths)
	if columnCount == 0 {
		state.ClearColumnSelection()
		state.ClearCellSelection()
		return
	}
	if state.selectedColumn != nil {
		selectedColumn := *state.selectedColumn
		if selectedColumn < 0 {
			selectedColumn = 0
		}
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

func (t Table) visibleBounds(state *TableState, height int) (int, int) {
	if height <= 0 || len(t.rows) == 0 {
		return 0, 0
	}
	offset := state.offset
	if offset > len(t.rows)-1 {
		offset = len(t.rows) - 1
	}
	if offset < 0 {
		offset = 0
	}
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
	for column := 0; column < startColumn; column++ {
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
	if len(cell.content.Lines) > 0 {
		line = cell.content.Lines[0]
	}
	cells := cellsFromLine(line)
	for i := 0; i < width && i < len(cells) && x+i < right; i++ {
		rendered := cells[i]
		rendered.Style = cellStyle.Patch(rendered.Style)
		buf.SetCell(x+i, y, rendered)
	}
}

func (c Chart) Block(block Block) Chart {
	c.block = &block
	return c
}

func (c Chart) XAxis(axis Axis) Chart {
	c.xAxis = axis
	return c
}

func (c Chart) YAxis(axis Axis) Chart {
	c.yAxis = axis
	return c
}

func (c Chart) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	chartArea := area
	if c.block != nil {
		c.block.Render(area, buf)
		chartArea = c.block.Inner(area)
	}
	if chartArea.Width == 0 || chartArea.Height == 0 {
		return
	}

	layout := c.layout(chartArea)
	c.renderYAxis(buf, layout)
	c.renderXAxis(buf, layout)
	c.renderYLabels(buf, layout)
	c.renderXLabels(buf, layout)
	c.renderYTitle(buf, layout)
	c.renderDatasets(buf, layout)
}

type chartAxisLayout struct {
	area       layout.Rect
	axisX      int
	graphLeft  int
	graphRight int
	axisY      int
	labelY     int
	hasXAxis   bool
	hasYAxis   bool
	yLabelW    int
}

func (c Chart) layout(area layout.Rect) chartAxisLayout {
	yLabelW := c.maxWidthLeftOfYAxis(area)
	hasXAxis := len(c.xAxis.labels) >= 2 && area.Height >= 2
	hasYAxis := len(c.yAxis.labels) > 0
	axisY := area.Y + area.Height - 1
	labelY := axisY
	if hasXAxis {
		axisY = area.Y + area.Height - 2
		labelY = area.Y + area.Height - 1
	}
	axisX := area.X + yLabelW
	graphLeft := axisX
	if hasYAxis {
		graphLeft++
	}
	if axisX > area.X+area.Width {
		axisX = area.X + area.Width
	}
	if graphLeft > area.X+area.Width {
		graphLeft = area.X + area.Width
	}
	return chartAxisLayout{
		area:       area,
		axisX:      axisX,
		graphLeft:  graphLeft,
		graphRight: area.X + area.Width,
		axisY:      axisY,
		labelY:     labelY,
		hasXAxis:   hasXAxis,
		hasYAxis:   hasYAxis,
		yLabelW:    yLabelW,
	}
}

func (c Chart) maxWidthLeftOfYAxis(area layout.Rect) int {
	maxWidth := 0
	hasYAxis := len(c.yAxis.labels) > 0
	for _, label := range c.yAxis.labels {
		maxWidth = maxInt(maxWidth, lineWidth(label))
	}
	if len(c.xAxis.labels) > 0 {
		firstWidth := lineWidth(c.xAxis.labels[0])
		switch c.xAxis.labelsAlignment {
		case layout.Left:
			if hasYAxis && firstWidth > 0 {
				firstWidth--
			}
			maxWidth = maxInt(maxWidth, firstWidth)
		case layout.Center:
			maxWidth = maxInt(maxWidth, firstWidth/2)
		case layout.Right:
		}
	}
	return minInt(maxWidth, area.Width/3)
}

func (c Chart) renderYAxis(buf *buffer.Buffer, l chartAxisLayout) {
	if !l.hasYAxis || l.axisX >= l.graphRight || l.area.Height == 0 {
		return
	}
	for y := l.area.Y; y <= l.axisY && y < l.area.Y+l.area.Height; y++ {
		buf.SetCell(l.axisX, y, buffer.Cell{Symbol: "│", Style: c.yAxis.axisStyle})
	}
}

func (c Chart) renderXAxis(buf *buffer.Buffer, l chartAxisLayout) {
	if !l.hasXAxis || l.graphLeft >= l.graphRight || l.axisY < l.area.Y || l.axisY >= l.area.Y+l.area.Height {
		return
	}
	start := l.graphLeft
	if l.hasYAxis {
		buf.SetCell(l.axisX, l.axisY, buffer.Cell{Symbol: "└", Style: c.yAxis.axisStyle.Patch(c.xAxis.axisStyle)})
	}
	for x := start; x < l.graphRight; x++ {
		buf.SetCell(x, l.axisY, buffer.Cell{Symbol: "─", Style: c.xAxis.axisStyle})
	}
}

func (c Chart) renderYLabels(buf *buffer.Buffer, l chartAxisLayout) {
	if len(c.yAxis.labels) < 2 || l.yLabelW <= 0 {
		return
	}
	top := l.area.Y
	bottom := l.axisY - 1
	if !l.hasXAxis {
		bottom = l.area.Y + l.area.Height - 1
	}
	if bottom < top {
		return
	}
	last := len(c.yAxis.labels) - 1
	for i, label := range c.yAxis.labels {
		y := bottom
		if last > 0 {
			y = bottom - (i * (bottom - top) / last)
		}
		c.renderLabel(buf, label, layout.NewRect(l.area.X, y, l.yLabelW, 1), c.yAxis.labelsAlignment, c.yAxis.axisStyle)
	}
}

func (c Chart) renderXLabels(buf *buffer.Buffer, l chartAxisLayout) {
	labels := c.xAxis.labels
	if len(labels) < 2 || l.graphLeft >= l.graphRight {
		return
	}
	graphWidth := l.graphRight - l.graphLeft
	widthBetweenTicks := graphWidth / len(labels)
	if widthBetweenTicks <= 0 {
		widthBetweenTicks = 1
	}

	firstArea := c.firstXLabelArea(l, lineWidth(labels[0]), widthBetweenTicks)
	firstAlignment := layout.Right
	switch c.xAxis.labelsAlignment {
	case layout.Center:
		firstAlignment = layout.Center
	case layout.Right:
		firstAlignment = layout.Left
	}
	c.renderLabel(buf, labels[0], firstArea, firstAlignment, c.xAxis.axisStyle)

	for i := 1; i < len(labels)-1; i++ {
		x := l.graphLeft + i*widthBetweenTicks + 1
		c.renderLabel(buf, labels[i], layout.NewRect(x, l.labelY, maxInt(0, widthBetweenTicks-1), 1), layout.Center, c.xAxis.axisStyle)
	}

	x := l.graphRight - widthBetweenTicks
	c.renderLabel(buf, labels[len(labels)-1], layout.NewRect(x, l.labelY, widthBetweenTicks, 1), layout.Right, c.xAxis.axisStyle)
}

func (c Chart) firstXLabelArea(l chartAxisLayout, labelWidth, maxWidthAfterYAxis int) layout.Rect {
	minX := l.area.X
	maxX := l.graphLeft
	switch c.xAxis.labelsAlignment {
	case layout.Center:
		maxX = l.graphLeft + minInt(maxWidthAfterYAxis, labelWidth)
	case layout.Right:
		minX = maxInt(l.area.X, l.graphLeft-1)
		maxX = l.graphLeft + maxWidthAfterYAxis
	}
	if maxX > l.graphRight {
		maxX = l.graphRight
	}
	if maxX < minX {
		maxX = minX
	}
	return layout.NewRect(minX, l.labelY, maxX-minX, 1)
}

func (c Chart) renderYTitle(buf *buffer.Buffer, l chartAxisLayout) {
	if c.yAxis.title == nil || l.graphLeft >= l.graphRight {
		return
	}
	cells := cellsFromLine(*c.yAxis.title)
	x := l.graphLeft
	for _, cell := range cells {
		if x >= l.graphRight {
			return
		}
		cell.Style = c.yAxis.axisStyle.Patch(cell.Style)
		buf.SetCell(x, l.area.Y, cell)
		x++
	}
}

func (c Chart) renderDatasets(buf *buffer.Buffer, l chartAxisLayout) {
	graphArea := c.graphArea(l)
	if graphArea.Width <= 0 || graphArea.Height <= 0 {
		return
	}
	xMin, xMax := c.xAxis.bounds[0], c.xAxis.bounds[1]
	yMin, yMax := c.yAxis.bounds[0], c.yAxis.bounds[1]
	if xMin == xMax || yMin == yMax {
		return
	}
	for _, dataset := range c.datasets {
		switch dataset.graphType {
		case GraphTypeLine:
			c.renderLineDataset(buf, graphArea, dataset, xMin, xMax, yMin, yMax)
		case GraphTypeScatter:
			for _, point := range dataset.data {
				c.plotPoint(buf, graphArea, dataset.style, point, xMin, xMax, yMin, yMax)
			}
		}
	}
}

func (c Chart) graphArea(l chartAxisLayout) layout.Rect {
	bottom := l.axisY
	if l.hasXAxis {
		bottom--
	}
	if bottom < l.area.Y {
		return layout.NewRect(l.graphLeft, l.area.Y, 0, 0)
	}
	return layout.NewRect(l.graphLeft, l.area.Y, l.graphRight-l.graphLeft, bottom-l.area.Y+1)
}

func (c Chart) renderLineDataset(buf *buffer.Buffer, area layout.Rect, dataset Dataset, xMin, xMax, yMin, yMax float64) {
	var previous *layout.Position
	for _, point := range dataset.data {
		mapped, ok := c.mapPoint(area, point, xMin, xMax, yMin, yMax)
		if !ok {
			previous = nil
			continue
		}
		if previous == nil {
			c.plotMappedPoint(buf, dataset.style, mapped)
		} else {
			c.plotLine(buf, dataset.style, *previous, mapped)
		}
		previous = &mapped
	}
}

func (c Chart) plotPoint(buf *buffer.Buffer, area layout.Rect, pointStyle style.Style, point ChartPoint, xMin, xMax, yMin, yMax float64) {
	mapped, ok := c.mapPoint(area, point, xMin, xMax, yMin, yMax)
	if !ok {
		return
	}
	c.plotMappedPoint(buf, pointStyle, mapped)
}

func (c Chart) mapPoint(area layout.Rect, point ChartPoint, xMin, xMax, yMin, yMax float64) (layout.Position, bool) {
	if point.X < xMin || point.X > xMax || point.Y < yMin || point.Y > yMax {
		return layout.Position{}, false
	}
	xRatio := (point.X - xMin) / (xMax - xMin)
	yRatio := (yMax - point.Y) / (yMax - yMin)
	x := area.X + int(math.Round(xRatio*float64(area.Width-1)))
	y := area.Y + int(math.Round(yRatio*float64(area.Height-1)))
	if x < area.X || x >= area.X+area.Width || y < area.Y || y >= area.Y+area.Height {
		return layout.Position{}, false
	}
	return layout.Position{X: x, Y: y}, true
}

func (c Chart) plotLine(buf *buffer.Buffer, lineStyle style.Style, start, end layout.Position) {
	dx := end.X - start.X
	if dx < 0 {
		dx = -dx
	}
	dy := end.Y - start.Y
	if dy < 0 {
		dy = -dy
	}
	steps := maxInt(dx, dy)
	if steps == 0 {
		c.plotMappedPoint(buf, lineStyle, start)
		return
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := int(math.Round(float64(start.X) + float64(end.X-start.X)*t))
		y := int(math.Round(float64(start.Y) + float64(end.Y-start.Y)*t))
		c.plotMappedPoint(buf, lineStyle, layout.Position{X: x, Y: y})
	}
}

func (c Chart) plotMappedPoint(buf *buffer.Buffer, pointStyle style.Style, point layout.Position) {
	if cell, ok := buf.CellAt(point.X, point.Y); ok {
		if cell.Symbol != " " {
			return
		}
		cell.Symbol = "•"
		cell.Style = cell.Style.Patch(pointStyle)
		buf.SetCell(point.X, point.Y, cell)
	}
}

func (c Chart) renderLabel(buf *buffer.Buffer, label text.Line, area layout.Rect, alignment layout.Alignment, baseStyle style.Style) {
	if area.Width <= 0 || area.Height <= 0 {
		return
	}
	cells := cellsFromLine(label)
	if len(cells) > area.Width {
		cells = cells[:area.Width]
	}
	offset := alignedOffset(len(cells), area.Width, alignment)
	for i, cell := range cells {
		cell.Style = baseStyle.Patch(cell.Style)
		buf.SetCell(area.X+offset+i, area.Y, cell)
	}
}

func NewGauge() Gauge {
	return Gauge{style: style.NewStyle(), gaugeStyle: style.NewStyle()}
}

func (g Gauge) Block(block Block) Gauge {
	g.block = &block
	return g
}

func (g Gauge) Percent(percent int) Gauge {
	if percent < 0 || percent > 100 {
		panic("gauge percent must be between 0 and 100")
	}
	g.ratio = float64(percent) / 100
	return g
}

func (g Gauge) Ratio(ratio float64) Gauge {
	if ratio < 0 || ratio > 1 {
		panic("gauge ratio must be between 0 and 1")
	}
	g.ratio = ratio
	return g
}

func (g Gauge) Label(label text.Span) Gauge {
	g.label = &label
	return g
}

func (g Gauge) LabelString(label string) Gauge {
	return g.Label(text.NewSpan(label))
}

func (g Gauge) Style(gaugeStyle style.Style) Gauge {
	g.style = gaugeStyle
	return g
}

func (g Gauge) GaugeStyle(gaugeStyle style.Style) Gauge {
	g.gaugeStyle = gaugeStyle
	return g
}

func (g Gauge) UseUnicode(useUnicode bool) Gauge {
	g.useUnicode = useUnicode
	return g
}

func (g Gauge) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, g.style)
	gaugeArea := area
	if g.block != nil {
		g.block.Render(area, buf)
		gaugeArea = g.block.Inner(area)
	}
	g.renderGauge(gaugeArea, buf)
}

func (g Gauge) renderGauge(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, g.gaugeStyle)

	label := g.effectiveLabel()
	labelRunes := []rune(label.Content)
	if len(labelRunes) > area.Width {
		labelRunes = labelRunes[:area.Width]
	}
	labelX := area.X + (area.Width-len(labelRunes))/2
	labelY := area.Y + area.Height/2

	filledWidth := float64(area.Width) * g.ratio
	end := area.X + int(math.Round(filledWidth))
	if g.useUnicode {
		end = area.X + int(math.Floor(filledWidth))
	}
	if end > area.X+area.Width {
		end = area.X + area.Width
	}

	for y := area.Y; y < area.Y+area.Height; y++ {
		for x := area.X; x < end; x++ {
			if y == labelY && x >= labelX && x < labelX+len(labelRunes) {
				buf.SetCell(x, y, buffer.Cell{Symbol: " ", Style: g.swappedGaugeStyle()})
				continue
			}
			buf.SetCell(x, y, buffer.Cell{Symbol: "█", Style: g.gaugeStyle})
		}
		if g.useUnicode && g.ratio < 1 && end < area.X+area.Width {
			buf.SetCell(end, y, buffer.Cell{Symbol: unicodeBlock(filledWidth - math.Floor(filledWidth)), Style: g.gaugeStyle})
		}
	}

	for i, r := range labelRunes {
		x := labelX + i
		cellStyle := g.gaugeStyle.Patch(label.Style)
		if x < end {
			cellStyle = g.swappedGaugeStyle().Patch(label.Style)
		}
		buf.SetCell(x, labelY, buffer.Cell{Symbol: string(r), Style: cellStyle})
	}
}

func (g Gauge) effectiveLabel() text.Span {
	if g.label != nil {
		return *g.label
	}
	return text.NewSpan(strconv.Itoa(int(math.Round(g.ratio*100))) + "%")
}

func (g Gauge) swappedGaugeStyle() style.Style {
	swapped := g.gaugeStyle
	swapped.Foreground = g.gaugeStyle.Background
	swapped.Background = g.gaugeStyle.Foreground
	return swapped
}

func unicodeBlock(frac float64) string {
	switch int(math.Round(frac * 8)) {
	case 1:
		return "▏"
	case 2:
		return "▎"
	case 3:
		return "▍"
	case 4:
		return "▌"
	case 5:
		return "▋"
	case 6:
		return "▊"
	case 7:
		return "▉"
	case 8:
		return "█"
	default:
		return " "
	}
}

func NewLineGauge() LineGauge {
	return LineGauge{
		style:          style.NewStyle(),
		filledStyle:    style.NewStyle(),
		unfilledStyle:  style.NewStyle(),
		filledSymbol:   "─",
		unfilledSymbol: "─",
	}
}

func (g LineGauge) Block(block Block) LineGauge {
	g.block = &block
	return g
}

func (g LineGauge) Ratio(ratio float64) LineGauge {
	if ratio < 0 || ratio > 1 {
		panic("line gauge ratio must be between 0 and 1")
	}
	g.ratio = ratio
	return g
}

func (g LineGauge) Label(label text.Line) LineGauge {
	g.label = &label
	return g
}

func (g LineGauge) LabelString(label string) LineGauge {
	return g.Label(text.LineFromString(label))
}

func (g LineGauge) Style(lineGaugeStyle style.Style) LineGauge {
	g.style = lineGaugeStyle
	return g
}

func (g LineGauge) FilledStyle(filledStyle style.Style) LineGauge {
	g.filledStyle = filledStyle
	return g
}

func (g LineGauge) UnfilledStyle(unfilledStyle style.Style) LineGauge {
	g.unfilledStyle = unfilledStyle
	return g
}

func (g LineGauge) FilledSymbol(symbol string) LineGauge {
	g.filledSymbol = symbol
	return g
}

func (g LineGauge) UnfilledSymbol(symbol string) LineGauge {
	g.unfilledSymbol = symbol
	return g
}

func (g LineGauge) Render(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	buf.SetStyle(area, g.style)
	gaugeArea := area
	if g.block != nil {
		g.block.Render(area, buf)
		gaugeArea = g.block.Inner(area)
	}
	g.renderLineGauge(gaugeArea, buf)
}

func (g LineGauge) renderLineGauge(area layout.Rect, buf *buffer.Buffer) {
	if area.Width == 0 || area.Height == 0 {
		return
	}
	right := area.X + area.Width
	x := g.writeLabel(area.X, area.Y, right, buf)
	start := x + 1
	if start >= right {
		return
	}
	remainingWidth := right - start
	filledWidth := int(math.Floor(float64(remainingWidth) * g.ratio))
	filledEnd := start + filledWidth
	for col := start; col < filledEnd; col++ {
		buf.SetCell(col, area.Y, buffer.Cell{Symbol: g.filledSymbol, Style: g.filledStyle})
	}
	for col := filledEnd; col < right; col++ {
		buf.SetCell(col, area.Y, buffer.Cell{Symbol: g.unfilledSymbol, Style: g.unfilledStyle})
	}
}

func (g LineGauge) writeLabel(x, y, right int, buf *buffer.Buffer) int {
	for _, cell := range cellsFromLine(g.effectiveLineGaugeLabel()) {
		if x >= right {
			return x
		}
		buf.SetCell(x, y, cell)
		x++
	}
	return x
}

func (g LineGauge) effectiveLineGaugeLabel() text.Line {
	if g.label != nil {
		return *g.label
	}
	return text.LineFromString(fmt.Sprintf("%3.0f%%", g.ratio*100))
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
	left := 0
	right := 0
	top := 0
	bottom := 0
	if b.borders.Has(LeftBorder) {
		left = 1
	}
	if b.borders.Has(RightBorder) {
		right = 1
	}
	if b.borders.Has(TopBorder) {
		top = 1
	}
	if b.borders.Has(BottomBorder) {
		bottom = 1
	}
	inner := area
	inner.X += left
	inner.Y += top
	inner.Width = maxInt(0, inner.Width-left-right)
	inner.Height = maxInt(0, inner.Height-top-bottom)
	return inner
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
	if b.borders.Has(LeftBorder) {
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
	if b.borders.Has(TopBorder) {
		for x := area.X; x <= right; x++ {
			buf.SetCell(x, area.Y, buffer.Cell{Symbol: "─", Style: b.style})
		}
	}
	if b.borders.Has(BottomBorder) && bottom != area.Y {
		for x := area.X; x <= right; x++ {
			buf.SetCell(x, bottom, buffer.Cell{Symbol: "─", Style: b.style})
		}
	}
	if b.borders.Has(LeftBorder) {
		for y := area.Y; y <= bottom; y++ {
			buf.SetCell(area.X, y, buffer.Cell{Symbol: "│", Style: b.style})
		}
	}
	if b.borders.Has(RightBorder) && right != area.X {
		for y := area.Y; y <= bottom; y++ {
			buf.SetCell(right, y, buffer.Cell{Symbol: "│", Style: b.style})
		}
	}
	if b.borders.Has(TopBorder) && b.borders.Has(LeftBorder) {
		buf.SetCell(area.X, area.Y, buffer.Cell{Symbol: "┌", Style: b.style})
	}
	if b.borders.Has(TopBorder) && b.borders.Has(RightBorder) && right != area.X {
		buf.SetCell(right, area.Y, buffer.Cell{Symbol: "┐", Style: b.style})
	}
	if b.borders.Has(BottomBorder) && b.borders.Has(LeftBorder) && bottom != area.Y {
		buf.SetCell(area.X, bottom, buffer.Cell{Symbol: "└", Style: b.style})
	}
	if b.borders.Has(BottomBorder) && b.borders.Has(RightBorder) && right != area.X && bottom != area.Y {
		buf.SetCell(right, bottom, buffer.Cell{Symbol: "┘", Style: b.style})
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

func lineWidth(line text.Line) int {
	return len(cellsFromLine(line))
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
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
	NoBorders    Borders = 0
	TopBorder    Borders = 1 << 0
	RightBorder  Borders = 1 << 1
	BottomBorder Borders = 1 << 2
	LeftBorder   Borders = 1 << 3
	AllBorders   Borders = TopBorder | RightBorder | BottomBorder | LeftBorder
)

func (b Borders) Has(border Borders) bool {
	return b&border != 0
}

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
