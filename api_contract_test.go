package gatui_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestPublicAPISurface_shouldExposeInitialRatatuiPortTypes(t *testing.T) {
	area := layout.NewRect(0, 0, 20, 3)
	buf := buffer.Empty(area)

	span := text.NewSpan("hello").Fg(style.Green).Bold()
	line := text.NewLine(span)
	content := text.NewText(line)

	paragraph := widgets.NewParagraph(content).Wrap(widgets.Wrap{Trim: true})
	var widget widgets.Widget = paragraph
	widget.Render(area, buf)

	block := widgets.BorderedBlock().
		Title(text.NewLine(text.StyledSpan("title", style.NewStyle().Fg(style.LightBlue)))).
		Cyan()
	block.Render(area, buf)

	widgets.Clear{}.Render(area, buf)

	_ = layout.Position{X: 1, Y: 2}
	_ = layout.Size{Width: 20, Height: 3}
	_ = layout.Margin{Horizontal: 1, Vertical: 1}
	_ = layout.NewLayout(layout.Vertical).Constraints(layout.Length(1), layout.Min(0))
	_ = layout.Center
	_ = style.NewStyle().Fg(style.Red).Bg(style.Black).AddModifier(style.ModifierItalic)
	_ = style.Styled[text.Span]{Value: span, Style: style.NewStyle()}
	_ = text.FromString("hello\nworld").Cyan().Bold()
	_ = text.LineFromString("right").Right()
	_ = buffer.WithLines([]string{"hello"})
	_, _ = buf.CellAt(0, 0)
	buf.SetFg(area, style.Cyan)
	buf.SetBg(area, style.Black)
	buf.SetModifier(area, style.ModifierBold)
	_ = widgets.NewParagraph(text.FromString("body")).
		Block(block).
		Alignment(layout.Right).
		Scroll(0, 1).
		Cyan()
	_ = widgets.AllBorders
}
