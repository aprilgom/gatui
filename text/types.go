package text

import "gatui/style"

type Span struct {
	Content string
	Style   style.Style
}

func NewSpan(content string) Span {
	return Span{Content: content, Style: style.NewStyle()}
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

type Line struct {
	Spans []Span
}

func NewLine(spans ...Span) Line {
	return Line{Spans: append([]Span(nil), spans...)}
}

type Text struct {
	Lines []Line
}

func NewText(lines ...Line) Text {
	return Text{Lines: append([]Line(nil), lines...)}
}
