package widgets

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/symbols"
)

type ScrollbarOrientation int

const (
	ScrollbarOrientationVerticalRight ScrollbarOrientation = iota
	ScrollbarOrientationVerticalLeft
	ScrollbarOrientationHorizontalBottom
	ScrollbarOrientationHorizontalTop
)

func (o ScrollbarOrientation) isVertical() bool {
	return o == ScrollbarOrientationVerticalRight || o == ScrollbarOrientationVerticalLeft
}

type ScrollDirection int

const (
	ScrollDirectionForward ScrollDirection = iota
	ScrollDirectionBackward
)

type ScrollbarState struct {
	contentLength         int
	position              int
	viewportContentLength int
}

func NewScrollbarState(contentLength int) ScrollbarState {
	if contentLength < 0 {
		contentLength = 0
	}
	return ScrollbarState{contentLength: contentLength}
}

func (s ScrollbarState) Position(position int) ScrollbarState {
	s.position = position
	return s
}

func (s ScrollbarState) ContentLength(length int) ScrollbarState {
	if length < 0 {
		length = 0
	}
	s.contentLength = length
	return s
}

func (s ScrollbarState) ViewportContentLength(length int) ScrollbarState {
	if length < 0 {
		length = 0
	}
	s.viewportContentLength = length
	return s
}

func (s *ScrollbarState) Next() {
	if s == nil {
		return
	}
	s.position = minInt(s.position+1, maxInt(0, s.contentLength-1))
}

func (s *ScrollbarState) Previous() {
	if s == nil {
		return
	}
	s.position = maxInt(0, s.position-1)
}

func (s *ScrollbarState) First() {
	if s == nil {
		return
	}
	s.position = 0
}

func (s *ScrollbarState) Last() {
	if s == nil {
		return
	}
	s.position = maxInt(0, s.contentLength-1)
}

func (s *ScrollbarState) Scroll(direction ScrollDirection) {
	if s == nil {
		return
	}
	switch direction {
	case ScrollDirectionBackward:
		s.Previous()
	default:
		s.Next()
	}
}

func (s ScrollbarState) PositionValue() int {
	return s.position
}

func (s ScrollbarState) ContentLengthValue() int {
	return s.contentLength
}

type Scrollbar struct {
	orientation ScrollbarOrientation
	thumbSymbol string
	trackSymbol *string
	beginSymbol *string
	endSymbol   *string
	thumbStyle  style.Style
	trackStyle  style.Style
	beginStyle  style.Style
	endStyle    style.Style
}

func NewScrollbar(orientation ScrollbarOrientation) Scrollbar {
	symbolSet := symbols.HorizontalScrollbarSet
	if orientation.isVertical() {
		symbolSet = symbols.VerticalScrollbarSet
	}
	track := symbolSet.Track
	begin := symbolSet.Begin
	end := symbolSet.End
	return Scrollbar{
		orientation: orientation,
		thumbSymbol: symbolSet.Thumb,
		trackSymbol: &track,
		beginSymbol: &begin,
		endSymbol:   &end,
		thumbStyle:  style.NewStyle(),
		trackStyle:  style.NewStyle(),
		beginStyle:  style.NewStyle(),
		endStyle:    style.NewStyle(),
	}
}

func (s Scrollbar) BeginSymbol(symbol string) Scrollbar {
	s.beginSymbol = &symbol
	return s
}

func (s Scrollbar) ClearBeginSymbol() Scrollbar {
	s.beginSymbol = nil
	return s
}

func (s Scrollbar) EndSymbol(symbol string) Scrollbar {
	s.endSymbol = &symbol
	return s
}

func (s Scrollbar) ClearEndSymbol() Scrollbar {
	s.endSymbol = nil
	return s
}

func (s Scrollbar) TrackSymbol(symbol string) Scrollbar {
	s.trackSymbol = &symbol
	return s
}

func (s Scrollbar) ClearTrackSymbol() Scrollbar {
	s.trackSymbol = nil
	return s
}

func (s Scrollbar) ThumbSymbol(symbol string) Scrollbar {
	s.thumbSymbol = symbol
	return s
}

func (s Scrollbar) ThumbStyle(thumbStyle style.Style) Scrollbar {
	s.thumbStyle = thumbStyle
	return s
}

func (s Scrollbar) TrackStyle(trackStyle style.Style) Scrollbar {
	s.trackStyle = trackStyle
	return s
}

func (s Scrollbar) BeginStyle(beginStyle style.Style) Scrollbar {
	s.beginStyle = beginStyle
	return s
}

func (s Scrollbar) EndStyle(endStyle style.Style) Scrollbar {
	s.endStyle = endStyle
	return s
}

func (s Scrollbar) Render(area layout.Rect, buf *buffer.Buffer) {
	state := ScrollbarState{}
	s.RenderStateful(area, buf, &state)
}

func (s Scrollbar) RenderStateful(area layout.Rect, buf *buffer.Buffer, state *ScrollbarState) {
	if state == nil || state.contentLength == 0 || area.Width == 0 || area.Height == 0 {
		return
	}
	scrollbarArea := s.scrollbarArea(area)
	trackLength := s.trackLength(scrollbarArea)
	if trackLength == 0 {
		return
	}

	trackStartLen, thumbLen, trackEndLen := s.partLengths(scrollbarArea, *state)
	cells := s.barCells(trackStartLen, thumbLen, trackEndLen)
	for i, cell := range cells {
		x, y := scrollbarArea.X, scrollbarArea.Y
		if s.orientation.isVertical() {
			y += i
		} else {
			x += i
		}
		if cell != nil {
			buf.SetCell(x, y, *cell)
		}
	}
}

func (s Scrollbar) RenderStatefulRef(area layout.Rect, buf *buffer.Buffer, state any) {
	if state == nil {
		s.RenderStateful(area, buf, nil)
		return
	}
	scrollbarState, ok := state.(*ScrollbarState)
	if !ok {
		panic("gatui: invalid state type for Scrollbar")
	}
	s.RenderStateful(area, buf, scrollbarState)
}

func (s Scrollbar) scrollbarArea(area layout.Rect) layout.Rect {
	switch s.orientation {
	case ScrollbarOrientationVerticalLeft:
		return layout.NewRect(area.X, area.Y, 1, area.Height)
	case ScrollbarOrientationVerticalRight:
		return layout.NewRect(area.X+area.Width-1, area.Y, 1, area.Height)
	case ScrollbarOrientationHorizontalTop:
		return layout.NewRect(area.X, area.Y, area.Width, 1)
	default:
		return layout.NewRect(area.X, area.Y+area.Height-1, area.Width, 1)
	}
}

func (s Scrollbar) trackLength(area layout.Rect) int {
	length := area.Width
	if s.orientation.isVertical() {
		length = area.Height
	}
	length -= optionalSymbolWidth(s.beginSymbol)
	length -= optionalSymbolWidth(s.endSymbol)
	return maxInt(0, length)
}

func (s Scrollbar) partLengths(area layout.Rect, state ScrollbarState) (int, int, int) {
	trackLength := s.trackLength(area)
	if trackLength == 0 {
		return 0, 0, 0
	}

	viewportLength := state.viewportContentLength
	if viewportLength == 0 {
		if s.orientation.isVertical() {
			viewportLength = area.Height
		} else {
			viewportLength = area.Width
		}
	}

	maxPosition := maxInt(0, state.contentLength-1)
	startPosition := clampInt(state.position, 0, maxPosition)
	maxViewportPosition := maxPosition + viewportLength
	if maxViewportPosition == 0 {
		return 0, trackLength, 0
	}

	thumbLen := clampInt(roundingDivide(viewportLength*trackLength, maxViewportPosition), 1, trackLength)
	thumbStart := clampInt(roundingDivide(startPosition*trackLength, maxViewportPosition), 0, trackLength-1)
	trackEnd := maxInt(0, trackLength-thumbStart-thumbLen)
	return thumbStart, thumbLen, trackEnd
}

func (s Scrollbar) barCells(trackStartLen, thumbLen, trackEndLen int) []*buffer.Cell {
	cells := make([]*buffer.Cell, 0, optionalSymbolWidth(s.beginSymbol)+trackStartLen+thumbLen+trackEndLen+optionalSymbolWidth(s.endSymbol))
	cells = appendOptionalCell(cells, s.beginSymbol, s.beginStyle)
	cells = appendRepeatedOptionalCells(cells, s.trackSymbol, s.trackStyle, trackStartLen)
	cells = appendRepeatedOptionalCells(cells, &s.thumbSymbol, s.thumbStyle, thumbLen)
	cells = appendRepeatedOptionalCells(cells, s.trackSymbol, s.trackStyle, trackEndLen)
	cells = appendOptionalCell(cells, s.endSymbol, s.endStyle)
	return cells
}

func appendOptionalCell(cells []*buffer.Cell, symbol *string, cellStyle style.Style) []*buffer.Cell {
	if symbol == nil {
		return cells
	}
	return append(cells, &buffer.Cell{Symbol: *symbol, Style: cellStyle})
}

func appendRepeatedOptionalCells(cells []*buffer.Cell, symbol *string, cellStyle style.Style, count int) []*buffer.Cell {
	for range count {
		if symbol == nil {
			cells = append(cells, nil)
			continue
		}
		cells = append(cells, &buffer.Cell{Symbol: *symbol, Style: cellStyle})
	}
	return cells
}

func optionalSymbolWidth(symbol *string) int {
	if symbol == nil {
		return 0
	}
	return len([]rune(*symbol))
}

func roundingDivide(numerator, denominator int) int {
	if denominator == 0 {
		return 0
	}
	return (numerator + denominator/2) / denominator
}

func clampInt(value, minimum, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
