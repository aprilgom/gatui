package text

import (
	"strings"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

type Text struct {
	Lines     []Line
	Style     style.Style
	Alignment *layout.Alignment
}

func NewText(lines ...Line) Text {
	return Text{Lines: append([]Line(nil), lines...), Style: style.NewStyle()}
}

func TextFromLines(lines []Line) Text {
	return NewText(lines...)
}

func TextFromSpan(span Span) Text {
	return NewText(NewLine(span))
}

func TextFromLine(line Line) Text {
	return NewText(line)
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

func (t Text) ToText() Text {
	return t
}

func (t Text) String() string {
	lines := make([]string, 0, len(t.Lines))
	for _, line := range t.Lines {
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n")
}

func (t Text) Extend(lines ...Line) Text {
	return t.AppendText(NewText(lines...))
}

func (t Text) AppendText(other Text) Text {
	t.Lines = append(append([]Line(nil), t.Lines...), other.Lines...)
	return t
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
