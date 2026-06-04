package gatui_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/symbols"
	"gatui/terminal"
	"gatui/terminal/testbackend"
	"gatui/text"
	"gatui/textbuffer"
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
	_ = layout.NewVerticalLayout(layout.Length(1), layout.Fill(1)).
		Direction(layout.Horizontal).
		Margin(1, 2).
		UniformMargin(1).
		HorizontalMargin(2).
		VerticalMargin(3)
	_, _ = layout.NewHorizontalLayout(layout.Length(1)).SplitWithSpacers(area)
	_ = layout.Center
	_ = style.NewStyle().Fg(style.Red).Bg(style.Black).AddModifier(style.ModifierItalic)
	_ = style.Styled[text.Span]{Value: span, Style: style.NewStyle()}
	_ = text.FromString("hello\nworld").Cyan().Bold()
	_ = text.LineFromString("right").Right()
	_, _ = textbuffer.SetSpan(buf, 0, 0, span, 5)
	_, _ = textbuffer.SetLine(buf, 0, 0, line, 5)
	_ = buffer.WithLines([]string{"hello"})
	_, _ = buf.CellAt(0, 0)
	buf.SetFg(area, style.Cyan)
	buf.SetBg(area, style.Black)
	buf.SetModifier(area, style.ModifierBold)
	_ = symbols.PlainBorderSet
	_ = symbols.NineLevelBarSet
	_ = symbols.NineLevelSparklineBarSet()
	_ = symbols.CanvasMarkerBraille
	_ = symbols.HorizontalScrollbarSet
	_ = symbols.MergeBorderSymbols(symbols.MergeStrategyExact, "│", "─")
	_ = widgets.BorderSet(symbols.PlainBorderSet)
	_ = widgets.NineLevelBarSet
	_ = widgets.ThreeLevelSparklineBarSet()
	_ = widgets.CanvasMarkerBraille
	_ = widgets.NewParagraph(text.FromString("body")).
		Block(block).
		Alignment(layout.Right).
		Scroll(0, 1).
		Cyan()
	_ = widgets.AllBorders

	backend := testbackend.New(20, 3)
	_ = terminal.FullscreenViewport()
	_ = terminal.FullscreenViewport().String()
	_ = terminal.FixedViewport(area)
	_ = terminal.FixedViewport(area).String()
	_ = terminal.InlineViewport(2)
	_ = terminal.InlineViewport(2).String()
	_ = terminal.DefaultTerminalOptions()
	_, _ = terminal.NewWithOptions(backend, terminal.TerminalOptions{
		Viewport: terminal.FixedViewport(area),
	})
	term, err := terminal.New(backend)
	if err != nil {
		t.Fatalf("terminal.New returned error: %v", err)
	}
	completed, err := term.Draw(func(frame *terminal.Frame) {
		_ = frame.Size()
		frame.RenderWidget(widgets.NewParagraph(text.FromString("terminal")), frame.Area())
		frame.RenderStatefulWidget(widgets.NewList([]widgets.ListItem{
			widgets.ListItemFromString("list"),
		}), frame.Area(), &widgets.ListState{})
		frame.RenderStatefulWidget(widgets.NewTable([]widgets.TableRow{
			widgets.TableRowFromStrings([]string{"table"}),
		}, []layout.Constraint{layout.Length(5)}), frame.Area(), &widgets.TableState{})
		frame.RenderStatefulWidget(
			widgets.NewScrollbar(widgets.ScrollbarOrientationHorizontalTop),
			frame.Area(),
			&widgets.ScrollbarState{},
		)
		frame.SetCursorPosition(layout.Position{X: 1, Y: 0})
	})
	if err != nil {
		t.Fatalf("terminal draw returned error: %v", err)
	}
	_ = completed.Area
	_ = completed.Buffer
	_ = completed.Count
	_, _ = term.TryDraw(func(frame *terminal.Frame) error {
		frame.RenderWidget(widgets.Clear{}, frame.Area())
		return nil
	})
	backend.SetSize(10, 2)
	_, _ = backend.WindowSize()
	_ = backend.String()
	_, _ = term.Size()
	_ = term.Autoresize()
	_ = term.Area()
	frame := term.Frame()
	frame.Buffer().SetSymbol(0, 0, "x")
	term.CurrentBuffer().SetSymbol(1, 0, "y")
	_ = term.Flush()
	_ = term.InsertBefore(1, func(buf *buffer.Buffer) {
		buf.SetSymbol(0, 0, "i")
	})
	term.SwapBuffers()
	_ = term.HideCursor()
	_ = term.ShowCursor()
	_ = term.SetCursorPosition(layout.Position{X: 1, Y: 0})
	_, _ = term.GetCursorPosition()
	_ = term.Resize(layout.NewRect(0, 0, 10, 2))
	_ = term.Clear()
	_ = term.Backend()
	_ = terminal.ClearAll
	_ = terminal.ClearAll.String()
	_, _ = terminal.ParseClearType("All")
	_ = terminal.ClearAfterCursor
	_ = terminal.ClearBeforeCursor
	_ = terminal.ClearCurrentLine
	_ = terminal.ClearUntilNewLine
	_ = testbackend.WithLines([]string{"seed"})
	_ = backend.Draws()
	_ = backend.FlushCount()
	_ = backend.ClearCount()
	_ = backend.HideCursorCount()
	_ = backend.ShowCursorCount()
	_ = backend.CursorPositions()
	_ = backend.CursorVisible()
	_ = backend.CursorPosition()
	_ = backend.AppendLines(1)
	_ = backend.AppendLinesCalls()
	_ = testbackend.NewNoScroll(2, 1).String()
}
