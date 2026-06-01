package text

import (
	"gatui/buffer"
	"gatui/layout"
	"gatui/style"

	"github.com/rivo/uniseg"
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
	return uniseg.StringWidth(s.Content)
}

func (s Span) StyledGraphemes(baseStyle style.Style) []StyledGrapheme {
	graphemeStyle := baseStyle.Patch(s.Style)
	graphemes := uniseg.NewGraphemes(s.Content)
	styled := make([]StyledGrapheme, 0)
	for graphemes.Next() {
		symbol := graphemes.Str()
		if containsControl(symbol) {
			continue
		}
		styled = append(styled, NewStyledGrapheme(symbol, graphemeStyle))
	}
	return styled
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
