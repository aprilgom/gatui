package text

import (
	"strings"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"

	"github.com/mattn/go-runewidth"
)

type Span struct {
	Content string
	Style   style.Style
}

func NewSpan(content string) Span {
	return Span{Content: content, Style: style.NewStyle()}
}

func StyledSpan(content string, spanStyle style.Style) Span {
	return Span{Content: content, Style: spanStyle}
}

func (s Span) PatchStyle(spanStyle style.Style) Span {
	s.Style = s.Style.Patch(spanStyle)
	return s
}

func (s Span) ResetStyle() Span {
	return s.PatchStyle(style.ResetStyle())
}

func (s Span) Width() int {
	return runewidth.StringWidth(s.Content)
}

func (s Span) Render(area layout.Rect, buf *buffer.Buffer) {
	renderSpan(s, area, buf, 0)
}

func (s Span) LeftLine() Line {
	return NewLine(s).Left()
}

func (s Span) CenterLine() Line {
	return NewLine(s).Center()
}

func (s Span) RightLine() Line {
	return NewLine(s).Right()
}

func (s Span) Fg(color style.Color) Span {
	s.Style = s.Style.Fg(color)
	return s
}

func (s Span) Bg(color style.Color) Span {
	s.Style = s.Style.Bg(color)
	return s
}

func (s Span) Bold() Span {
	s.Style = s.Style.AddModifier(style.ModifierBold)
	return s
}

func (s Span) Italic() Span {
	s.Style = s.Style.AddModifier(style.ModifierItalic)
	return s
}

func (s Span) Cyan() Span {
	return s.Fg(style.Cyan)
}

func (s Span) OnCyan() Span {
	return s.Bg(style.Cyan)
}

func (s Span) LightBlue() Span {
	return s.Fg(style.LightBlue)
}

type Line struct {
	Spans     []Span
	LineStyle style.Style
	Alignment *layout.Alignment
}

func NewLine(spans ...Span) Line {
	return Line{Spans: append([]Span(nil), spans...), LineStyle: style.NewStyle()}
}

func LineFromString(content string) Line {
	return NewLine(NewSpan(content))
}

func StyledLine(content string, lineStyle style.Style) Line {
	return LineFromString(content).Style(lineStyle)
}

func (l Line) PatchStyle(lineStyle style.Style) Line {
	l.LineStyle = l.LineStyle.Patch(lineStyle)
	return l
}

func (l Line) ResetStyle() Line {
	return l.PatchStyle(style.ResetStyle())
}

func (l Line) PushSpan(span Span) Line {
	l.Spans = append(append([]Span(nil), l.Spans...), span)
	return l
}

func (l Line) AppendSpans(spans ...Span) Line {
	l.Spans = append(append([]Span(nil), l.Spans...), spans...)
	return l
}

func (l Line) Width() int {
	width := 0
	for _, span := range l.Spans {
		width += span.Width()
	}
	return width
}

func (l Line) Render(area layout.Rect, buf *buffer.Buffer) {
	l.RenderWithAlignment(area, buf, nil)
}

func (l Line) RenderWithAlignment(area layout.Rect, buf *buffer.Buffer, fallback *layout.Alignment) {
	if buf == nil {
		return
	}
	area = area.Intersection(buf.Area)
	if area.Width == 0 || area.Height == 0 {
		return
	}
	area.Height = 1
	lineWidth := l.Width()
	if lineWidth == 0 {
		return
	}

	buf.SetStyle(area, l.LineStyle)

	alignment := l.Alignment
	if alignment == nil {
		alignment = fallback
	}

	if lineWidth <= area.Width {
		x := area.X + alignedRenderOffset(lineWidth, area.Width, alignment)
		renderLineSpans(l.Spans, layout.NewRect(x, area.Y, area.Right()-x, 1), buf, 0)
		return
	}

	skipWidth := 0
	if alignment != nil {
		switch *alignment {
		case layout.Center:
			skipWidth = (lineWidth - area.Width) / 2
		case layout.Right:
			skipWidth = lineWidth - area.Width
		}
	}
	renderLineSpans(l.Spans, area, buf, skipWidth)
}

func (l Line) Style(lineStyle style.Style) Line {
	l.LineStyle = lineStyle
	return l
}

func (l Line) Fg(color style.Color) Line {
	l.LineStyle = l.LineStyle.Fg(color)
	return l
}

func (l Line) Bg(color style.Color) Line {
	l.LineStyle = l.LineStyle.Bg(color)
	return l
}

func (l Line) Bold() Line {
	l.LineStyle = l.LineStyle.AddModifier(style.ModifierBold)
	return l
}

func (l Line) Italic() Line {
	l.LineStyle = l.LineStyle.AddModifier(style.ModifierItalic)
	return l
}

func (l Line) Cyan() Line {
	return l.Fg(style.Cyan)
}

func (l Line) Left() Line {
	alignment := layout.Left
	l.Alignment = &alignment
	return l
}

func (l Line) Center() Line {
	alignment := layout.Center
	l.Alignment = &alignment
	return l
}

func (l Line) Right() Line {
	alignment := layout.Right
	l.Alignment = &alignment
	return l
}

type Text struct {
	Lines     []Line
	Style     style.Style
	Alignment *layout.Alignment
}

func NewText(lines ...Line) Text {
	return Text{Lines: append([]Line(nil), lines...), Style: style.NewStyle()}
}

func FromString(content string) Text {
	parts := strings.Split(content, "\n")
	lines := make([]Line, 0, len(parts))
	for _, part := range parts {
		lines = append(lines, LineFromString(part))
	}
	return NewText(lines...)
}

func StyledText(content string, textStyle style.Style) Text {
	t := FromString(content)
	t.Style = textStyle
	return t
}

func (t Text) PatchStyle(textStyle style.Style) Text {
	t.Style = t.Style.Patch(textStyle)
	return t
}

func (t Text) ResetStyle() Text {
	return t.PatchStyle(style.ResetStyle())
}

func (t Text) Align(alignment layout.Alignment) Text {
	t.Alignment = &alignment
	return t
}

func (t Text) PushLine(line Line) Text {
	t.Lines = append(append([]Line(nil), t.Lines...), line)
	return t
}

func (t Text) PushSpan(span Span) Text {
	if len(t.Lines) == 0 {
		return t.PushLine(NewLine(span))
	}

	lines := append([]Line(nil), t.Lines...)
	last := len(lines) - 1
	lines[last] = lines[last].PushSpan(span)
	t.Lines = lines
	return t
}

func (t Text) Width() int {
	width := 0
	for _, line := range t.Lines {
		lineWidth := line.Width()
		if lineWidth > width {
			width = lineWidth
		}
	}
	return width
}

func (t Text) Height() int {
	return len(t.Lines)
}

func (t Text) Render(area layout.Rect, buf *buffer.Buffer) {
	if buf == nil {
		return
	}
	area = area.Intersection(buf.Area)
	if area.Width == 0 || area.Height == 0 {
		return
	}

	buf.SetStyle(area, t.Style)
	for y, line := range t.Lines {
		if y >= area.Height {
			break
		}
		lineArea := layout.NewRect(area.X, area.Y+y, area.Width, 1)
		line.RenderWithAlignment(lineArea, buf, t.Alignment)
	}
}

func (t Text) Fg(color style.Color) Text {
	return t.PatchStyle(style.NewStyle().Fg(color))
}

func (t Text) Bg(color style.Color) Text {
	return t.PatchStyle(style.NewStyle().Bg(color))
}

func (t Text) Bold() Text {
	return t.PatchStyle(style.NewStyle().AddModifier(style.ModifierBold))
}

func (t Text) Italic() Text {
	return t.PatchStyle(style.NewStyle().AddModifier(style.ModifierItalic))
}

func (t Text) Cyan() Text {
	return t.Fg(style.Cyan)
}

func (t Text) Left() Text {
	return t.Align(layout.Left)
}

func (t Text) Center() Text {
	return t.Align(layout.Center)
}

func (t Text) Right() Text {
	return t.Align(layout.Right)
}

func renderLineSpans(spans []Span, area layout.Rect, buf *buffer.Buffer, skipWidth int) {
	x := area.X
	for _, span := range spans {
		spanWidth := span.Width()
		if skipWidth >= spanWidth {
			skipWidth -= spanWidth
			continue
		}
		if skipWidth > 0 {
			x = renderSpan(span, layout.NewRect(x, area.Y, area.Right()-x, 1), buf, skipWidth)
			skipWidth = 0
		} else {
			x = renderSpan(span, layout.NewRect(x, area.Y, area.Right()-x, 1), buf, 0)
		}
		if x >= area.Right() {
			return
		}
	}
}

func renderSpan(span Span, area layout.Rect, buf *buffer.Buffer, skipWidth int) int {
	if buf == nil {
		return area.X
	}
	area = area.Intersection(buf.Area)
	if area.Width == 0 || area.Height == 0 {
		return area.X
	}

	x := area.X
	right := area.Right()
	renderedAny := false
	for _, r := range span.Content {
		if r == '\n' {
			continue
		}
		symbol := string(r)
		width := runewidth.StringWidth(symbol)
		if width == 0 {
			if !renderedAny {
				setSpanCellSymbol(buf, x, area.Y, symbol, span.Style, false)
				renderedAny = true
			} else if x == area.X {
				setSpanCellSymbol(buf, x, area.Y, symbol, span.Style, true)
			} else {
				setSpanCellSymbol(buf, x-1, area.Y, symbol, span.Style, true)
			}
			continue
		}
		if skipWidth >= width {
			skipWidth -= width
			continue
		}
		if skipWidth > 0 {
			x += width - skipWidth
			skipWidth = 0
			continue
		}
		if x+width > right {
			break
		}

		setSpanCellSymbol(buf, x, area.Y, symbol, span.Style, renderedAny && x == area.X)
		for hidden := 1; hidden < width; hidden++ {
			buf.SetCell(x+hidden, area.Y, buffer.Cell{Symbol: " ", Style: style.NewStyle()})
		}
		x += width
		renderedAny = true
	}
	return x
}

func setSpanCellSymbol(buf *buffer.Buffer, x, y int, symbol string, spanStyle style.Style, appendSymbol bool) {
	cellStyle := style.NewStyle()
	cellSymbol := symbol
	if cell, ok := buf.CellAt(x, y); ok {
		cellStyle = cell.Style
		if appendSymbol {
			cellSymbol = cell.Symbol + symbol
		}
	}
	buf.SetCell(x, y, buffer.Cell{Symbol: cellSymbol, Style: cellStyle.Patch(spanStyle)})
}

func alignedRenderOffset(lineWidth, areaWidth int, alignment *layout.Alignment) int {
	if alignment == nil || lineWidth >= areaWidth {
		return 0
	}
	switch *alignment {
	case layout.Center:
		return (areaWidth - lineWidth) / 2
	case layout.Right:
		return areaWidth - lineWidth
	default:
		return 0
	}
}
