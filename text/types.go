package text

import (
	"strings"

	"gatui/layout"
	"gatui/style"
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
	Alignment *layout.Alignment
}

func NewLine(spans ...Span) Line {
	return Line{Spans: append([]Span(nil), spans...)}
}

func LineFromString(content string) Line {
	return NewLine(NewSpan(content))
}

func (l Line) Fg(color style.Color) Line {
	for i := range l.Spans {
		l.Spans[i] = l.Spans[i].Fg(color)
	}
	return l
}

func (l Line) Bg(color style.Color) Line {
	for i := range l.Spans {
		l.Spans[i] = l.Spans[i].Bg(color)
	}
	return l
}

func (l Line) Bold() Line {
	for i := range l.Spans {
		l.Spans[i] = l.Spans[i].Bold()
	}
	return l
}

func (l Line) Italic() Line {
	for i := range l.Spans {
		l.Spans[i] = l.Spans[i].Italic()
	}
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
}

func NewText(lines ...Line) Text {
	return Text{Lines: append([]Line(nil), lines...)}
}

func FromString(content string) Text {
	parts := strings.Split(content, "\n")
	lines := make([]Line, 0, len(parts))
	for _, part := range parts {
		lines = append(lines, LineFromString(part))
	}
	return NewText(lines...)
}

func (t Text) Fg(color style.Color) Text {
	for i := range t.Lines {
		t.Lines[i] = t.Lines[i].Fg(color)
	}
	return t
}

func (t Text) Bg(color style.Color) Text {
	for i := range t.Lines {
		t.Lines[i] = t.Lines[i].Bg(color)
	}
	return t
}

func (t Text) Bold() Text {
	for i := range t.Lines {
		t.Lines[i] = t.Lines[i].Bold()
	}
	return t
}

func (t Text) Italic() Text {
	for i := range t.Lines {
		t.Lines[i] = t.Lines[i].Italic()
	}
	return t
}

func (t Text) Cyan() Text {
	return t.Fg(style.Cyan)
}
