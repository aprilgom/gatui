# Gatui LLM Usage Guide

This guide is for LLMs, coding agents, and Go users building terminal UIs with
Gatui. Prefer these patterns before inventing wrappers or copying Ratatui Rust
code mechanically.

## First Rules

- Run commands from the module root: `/Users/aprilgom/gatui/gatui`.
- Start from an existing example under `examples/` when possible.
- Use `backend/tcell` for real terminal drawing.
- Keep input polling outside Gatui core. Read keyboard, mouse, and resize events
  through `tcell`, then redraw with `terminal.Terminal`.
- Call `tcellbackend.NewWithScreen(screen)` before tcell features that require
  an initialized screen, such as `screen.EnableMouse()`.
- Use `terminal/testbackend` for headless tests instead of starting a real
  terminal.
- Check [ratatui-correspondence.md](ratatui-correspondence.md) before porting a
  Ratatui path or method name.

## Minimal App Skeleton

Use this as the default shape for a small interactive app:

```go
package main

import (
	"fmt"
	"os"

	gatuitcell "github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/terminal"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
	tcell "github.com/gdamore/tcell/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	backend, err := gatuitcell.NewWithScreen(screen)
	if err != nil {
		return err
	}
	defer backend.Close()

	term, err := terminal.New(backend)
	if err != nil {
		return err
	}

	for {
		if _, err := term.Draw(render); err != nil {
			return err
		}
		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if event.Rune() == 'q' || event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyCtrlC {
				return nil
			}
		}
	}
}

func render(frame *terminal.Frame) {
	frame.RenderWidget(widgets.NewParagraph(text.FromString("Hello Gatui")), frame.Area())
}
```

For mouse input, enable mouse after backend initialization:

```go
backend, err := gatuitcell.NewWithScreen(screen)
if err != nil {
	return err
}
defer backend.Close()
screen.EnableMouse()
```

## Widget Cookbook

### Paragraph

Use `Paragraph` for plain or styled text blocks.

```go
paragraph := widgets.NewParagraph(text.FromString("Status: ready"))
frame.RenderWidget(paragraph, area)
```

### Block

Use `Block` to frame another widget or mark an area with a title.

```go
block := widgets.NewBlock().
	Title(text.NewLine(text.NewSpan("Overview"))).
	Borders(widgets.AllBorders)
frame.RenderWidget(block, area)
```

### Gauge

Use `Gauge` for progress with a filled bar and optional label.

```go
gauge := widgets.NewGauge().
	Block(widgets.NewBlock().
		Title(text.NewLine(text.NewSpan("Progress"))).
		Borders(widgets.AllBorders)).
	Percent(65).
	LabelString("65%")
frame.RenderWidget(gauge, area)
```

Use `LineGauge` when vertical space is tight. See
`examples/widget-gauge/main.go`.

### List

Use `List` with `ListState` for selectable rows.

```go
items := []widgets.ListItem{
	widgets.ListItemFromString("Inbox"),
	widgets.ListItemFromString("Archive"),
}
list := widgets.NewList(items).
	Block(widgets.NewBlock().
		Title(text.NewLine(text.NewSpan("Folders"))).
		Borders(widgets.AllBorders)).
	HighlightSymbol("> ")
state := widgets.NewListState().WithSelected(0)
frame.RenderStatefulWidget(list, area, &state)
```

Keep `ListState` in app state when selection must persist across frames.

### Table

Use `Table` with explicit widths. Use `TableState` for row, column, or cell
selection.

```go
header := widgets.TableRowFromStrings([]string{"Name", "Role"}).Bold()
rows := []widgets.TableRow{
	widgets.TableRowFromStrings([]string{"Ada", "Engineer"}),
	widgets.TableRowFromStrings([]string{"Grace", "Compiler"}),
}
table := widgets.NewTable(rows, []layout.Constraint{
	layout.Percentage(50),
	layout.Percentage(50),
}).
	Header(header).
	Block(widgets.NewBlock().
		Title(text.NewLine(text.NewSpan("People"))).
		Borders(widgets.AllBorders)).
	HighlightSymbol("> ")
state := widgets.NewTableState().WithSelected(0)
frame.RenderStatefulWidget(table, area, &state)
```

### Canvas

Use `Canvas` when drawing shapes, maps, points, or simple animations.

```go
canvas := widgets.NewCanvas().
	XBounds(-180, 180).
	YBounds(-90, 90).
	Marker(widgets.CanvasMarkerBraille).
	Paint(func(ctx *widgets.CanvasContext) {
		ctx.Draw(widgets.Map{Resolution: widgets.MapResolutionHigh, Color: style.Green})
		ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{{X: -74, Y: 40.7}}, style.Red))
	})
frame.RenderWidget(canvas, area)
```

Canvas does not own a `Block`; render a title or border in a neighboring layout
area if the screen needs framing. See `examples/widget-canvas/main.go`.

## Choosing an Example

| Need | Start from |
| --- | --- |
| Smallest smoke test | `examples/hello-world` |
| Terminal lifecycle and one draw loop | `examples/minimal` |
| Animated progress | `examples/app-gauge` |
| Interactive table | `examples/app-table` |
| Interactive list | `examples/app-list` |
| Mouse input and canvas animation | `examples/app-canvas` |
| Isolated gauge rendering | `examples/widget-gauge` |
| Isolated table rendering | `examples/widget-table` |
| Isolated list rendering | `examples/widget-list` |
| Isolated canvas rendering | `examples/widget-canvas` |

## Ratatui Translation Notes

- `ratatui::Terminal::draw` maps to `terminal.Terminal.Draw`.
- `ratatui::Frame::render_widget` maps to `frame.RenderWidget`.
- `ratatui::Frame::render_stateful_widget` maps to
  `frame.RenderStatefulWidget`.
- `ratatui_crossterm` has no direct Gatui equivalent. Use `backend/tcell`.
- Ratatui builder chains usually become Go constructor and method chains, but
  names and argument shapes differ. Inspect the target Gatui type before
  porting method names.

For a broader path-by-path mapping, see
[ratatui-correspondence.md](ratatui-correspondence.md).

## Verification

After editing Gatui code or examples, run:

```sh
go fix ./...
go test ./...
golangci-lint run ./...
```

For example-only work, also run:

```sh
go test ./examples/...
```
