package widgets

import (
	"math"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

type GraphType int

const (
	GraphTypeScatter GraphType = iota
	GraphTypeLine
	GraphTypeBar
	GraphTypeArea
)

type LegendPosition int

const (
	LegendPositionTopRight LegendPosition = iota
	LegendPositionTopLeft
	LegendPositionTop
	LegendPositionLeft
	LegendPositionRight
	LegendPositionBottomLeft
	LegendPositionBottom
	LegendPositionBottomRight
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
	hasName   bool
	data      []ChartPoint
	graphType GraphType
	style     style.Style
	fillToY   float64
}

type Chart struct {
	datasets           []Dataset
	block              *Block
	xAxis              Axis
	yAxis              Axis
	legendPosition     *LegendPosition
	hiddenLegendWidth  layout.Constraint
	hiddenLegendHeight layout.Constraint
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
	d.hasName = true
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

func (d Dataset) FillToY(y float64) Dataset {
	d.fillToY = y
	return d
}

func (d Dataset) Style(datasetStyle style.Style) Dataset {
	d.style = datasetStyle
	return d
}

func NewChart(datasets []Dataset) Chart {
	defaultLegendPosition := LegendPositionTopRight
	return Chart{
		datasets:           append([]Dataset(nil), datasets...),
		xAxis:              NewAxis(),
		yAxis:              NewAxis(),
		legendPosition:     &defaultLegendPosition,
		hiddenLegendWidth:  layout.Ratio(1, 4),
		hiddenLegendHeight: layout.Ratio(1, 4),
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

func (c Chart) LegendPosition(position LegendPosition) Chart {
	c.legendPosition = &position
	return c
}

func (c Chart) HideLegend() Chart {
	c.legendPosition = nil
	return c
}

func (c Chart) HiddenLegendConstraints(width, height layout.Constraint) Chart {
	c.hiddenLegendWidth = width
	c.hiddenLegendHeight = height
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
	c.renderLegend(buf, layout)
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
		case GraphTypeArea:
			c.renderAreaDataset(buf, graphArea, dataset, xMin, xMax, yMin, yMax)
		case GraphTypeBar:
			c.renderBarDataset(buf, graphArea, dataset, xMin, xMax, yMin, yMax)
		case GraphTypeLine:
			c.renderLineDataset(buf, graphArea, dataset, xMin, xMax, yMin, yMax)
		case GraphTypeScatter:
			for _, point := range dataset.data {
				c.plotPoint(buf, graphArea, dataset.style, point, xMin, xMax, yMin, yMax)
			}
		}
	}
}

func (c Chart) renderLegend(buf *buffer.Buffer, l chartAxisLayout) {
	legendArea, ok := c.legendArea(l.area)
	if !ok {
		return
	}
	BorderedBlock().Render(legendArea, buf)
	innerWidth := maxInt(0, legendArea.Width-2)
	y := legendArea.Y + 1
	for _, dataset := range c.datasets {
		if !dataset.hasName || y >= legendArea.Y+legendArea.Height-1 {
			continue
		}
		x := legendArea.X + 1
		for _, r := range dataset.name {
			if x >= legendArea.X+legendArea.Width-1 {
				break
			}
			buf.SetCell(x, y, buffer.Cell{Symbol: string(r), Style: dataset.style})
			x++
		}
		for ; x < legendArea.X+1+innerWidth; x++ {
			buf.SetCell(x, y, buffer.Cell{Symbol: " ", Style: dataset.style})
		}
		y++
	}
}

func (c Chart) legendArea(area layout.Rect) (layout.Rect, bool) {
	if c.legendPosition == nil {
		return layout.Rect{}, false
	}
	innerWidth := 0
	namedDatasets := 0
	for _, dataset := range c.datasets {
		if !dataset.hasName {
			continue
		}
		namedDatasets++
		innerWidth = maxInt(innerWidth, len([]rune(dataset.name)))
	}
	if namedDatasets == 0 {
		return layout.Rect{}, false
	}
	legendWidth := innerWidth + 2
	legendHeight := namedDatasets + 2
	if legendWidth > area.Width || legendHeight > area.Height {
		return layout.Rect{}, false
	}
	if !c.legendFits(area, legendWidth, legendHeight) {
		return layout.Rect{}, false
	}

	x := area.X
	y := area.Y
	switch *c.legendPosition {
	case LegendPositionTopLeft:
	case LegendPositionTop:
		x = area.X + (area.Width-legendWidth)/2
	case LegendPositionTopRight:
		x = area.X + area.Width - legendWidth
	case LegendPositionLeft:
		y = area.Y + (area.Height-legendHeight)/2
	case LegendPositionRight:
		x = area.X + area.Width - legendWidth
		y = area.Y + (area.Height-legendHeight)/2
	case LegendPositionBottomLeft:
		y = area.Y + area.Height - legendHeight
	case LegendPositionBottom:
		x = area.X + (area.Width-legendWidth)/2
		y = area.Y + area.Height - legendHeight
	case LegendPositionBottomRight:
		x = area.X + area.Width - legendWidth
		y = area.Y + area.Height - legendHeight
	}
	return layout.NewRect(x, y, legendWidth, legendHeight), true
}

func (c Chart) legendFits(area layout.Rect, legendWidth, legendHeight int) bool {
	maxWidth, widthAlways := legendConstraintSize(c.hiddenLegendWidth, area.Width)
	maxHeight, heightAlways := legendConstraintSize(c.hiddenLegendHeight, area.Height)
	widthFits := widthAlways || legendWidth <= maxWidth
	heightFits := heightAlways || legendHeight <= maxHeight
	return widthFits && heightFits
}

func legendConstraintSize(constraint layout.Constraint, total int) (int, bool) {
	if constraint.IsLength() {
		return constraint.Value(), false
	}
	if constraint.IsPercentage() {
		return total * constraint.Value() / 100, false
	}
	if constraint.IsRatio() {
		denominator := constraint.Denominator()
		if denominator == 0 {
			return 0, false
		}
		return total * constraint.Value() / denominator, false
	}
	return total, true
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

func (c Chart) renderBarDataset(buf *buffer.Buffer, area layout.Rect, dataset Dataset, xMin, xMax, yMin, yMax float64) {
	for _, point := range dataset.data {
		start, ok := c.mapPoint(area, ChartPoint{X: point.X, Y: 0}, xMin, xMax, yMin, yMax)
		if !ok {
			continue
		}
		end, ok := c.mapPoint(area, point, xMin, xMax, yMin, yMax)
		if !ok {
			continue
		}
		c.plotLine(buf, dataset.style, start, end)
	}
}

func (c Chart) renderAreaDataset(buf *buffer.Buffer, area layout.Rect, dataset Dataset, xMin, xMax, yMin, yMax float64) {
	fillToY := math.Min(math.Max(dataset.fillToY, yMin), yMax)
	for i := 1; i < len(dataset.data); i++ {
		start, ok := c.mapPoint(area, dataset.data[i-1], xMin, xMax, yMin, yMax)
		if !ok {
			continue
		}
		end, ok := c.mapPoint(area, dataset.data[i], xMin, xMax, yMin, yMax)
		if !ok {
			continue
		}
		fillPoint, ok := c.mapPoint(area, ChartPoint{X: dataset.data[i-1].X, Y: fillToY}, xMin, xMax, yMin, yMax)
		if !ok {
			continue
		}
		forEachChartLinePoint(start, end, func(point layout.Position) {
			fillStart := minInt(point.Y, fillPoint.Y)
			fillEnd := maxInt(point.Y, fillPoint.Y)
			for y := fillStart; y <= fillEnd; y++ {
				c.plotMappedPoint(buf, dataset.style, layout.Position{X: point.X, Y: y})
			}
		})
	}
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
	forEachChartLinePoint(start, end, func(point layout.Position) {
		c.plotMappedPoint(buf, lineStyle, point)
	})
}

func forEachChartLinePoint(start, end layout.Position, fn func(layout.Position)) {
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
		fn(start)
		return
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := int(math.Round(float64(start.X) + float64(end.X-start.X)*t))
		y := int(math.Round(float64(start.Y) + float64(end.Y-start.Y)*t))
		fn(layout.Position{X: x, Y: y})
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
