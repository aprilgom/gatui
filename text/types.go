package text

import (
	"strings"

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

func (s Span) Width() int {
	return runewidth.StringWidth(s.Content)
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

func (l Line) Width() int {
	width := 0
	for _, span := range l.Spans {
		width += span.Width()
	}
	return width
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
	Lines []Line
	Style style.Style
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
