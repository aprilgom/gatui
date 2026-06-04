package text

import (
	"strings"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

type Line struct {
	Spans     []Span
	LineStyle style.Style
	Alignment *layout.Alignment
}

func NewLine(spans ...Span) Line {
	return Line{Spans: append([]Span(nil), spans...), LineStyle: style.NewStyle()}
}

func LineFromSpans(spans ...Span) Line {
	return NewLine(spans...)
}

func LineFromSpan(span Span) Line {
	return NewLine(span)
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

func (l Line) String() string {
	var builder strings.Builder
	for _, span := range l.Spans {
		builder.WriteString(span.Content)
	}
	return builder.String()
}

func (l Line) Extend(spans ...Span) Line {
	return l.AppendSpans(spans...)
}

func (l Line) PushSpan(span Span) Line {
	l.Spans = append(append([]Span(nil), l.Spans...), span)
	return l
}

func (l Line) AppendSpans(spans ...Span) Line {
	l.Spans = append(append([]Span(nil), l.Spans...), spans...)
	return l
}

func (l Line) AddLine(other Line) Line {
	return l.AppendSpans(other.Spans...)
}

func (l Line) Width() int {
	width := 0
	for _, span := range l.Spans {
		width += span.Width()
	}
	return width
}

func (l Line) StyledGraphemes(baseStyle style.Style) []StyledGrapheme {
	lineStyle := baseStyle.Patch(l.LineStyle)
	styled := make([]StyledGrapheme, 0)
	for _, span := range l.Spans {
		styled = append(styled, span.StyledGraphemes(lineStyle)...)
	}
	return styled
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
	buf.SetStyle(area, l.LineStyle)

	lineWidth := l.Width()
	if lineWidth == 0 {
		return
	}

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
