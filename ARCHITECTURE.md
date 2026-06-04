# gatui Architecture

`gatui` mirrors Ratatui's separation between pure rendering logic and terminal backends. Core packages should remain testable without a terminal.

```mermaid
flowchart TD
    layout[layout<br/>Rect, Margin, Offset, Constraint]
    style[style<br/>Color, Modifier, Style]
    text[text<br/>StyledGrapheme, Span, Line, Text]
    symbols[symbols<br/>BorderSet, BarSet, CanvasMarker]
    widgets[widgets<br/>Widget, Paragraph, Block, Clear]
    buffer[buffer<br/>Buffer, Cell]
    backend[backend/tcell<br/>terminal drawing and flush]

    style --> text
    style --> buffer
    layout --> buffer
    layout --> text
    layout --> widgets
    symbols --> widgets
    buffer --> text
    text --> widgets
    widgets --> buffer
    buffer --> backend
```

## Package Responsibilities

- `layout` owns geometry and area splitting. It must not depend on widgets or terminal backends.
- `style` owns style value types and chainable style helpers. It must stay backend-neutral.
- `text` owns styled textual content and rendering text primitives into buffers. It may depend on `buffer`, `layout`, and `style`, but not on `widgets` or terminal backends.
- `buffer` owns the off-screen cell grid. It may depend on `layout` and `style`.
- `symbols` owns reusable terminal symbol sets and glyph helpers for borders, bars, sparklines, canvas markers, and scrollbars. It must stay independent of widgets.
- `widgets` owns renderable components. Widgets write into `buffer.Buffer` inside a `layout.Rect`.
- `backend/tcell` translates `buffer.Cell` and `style.Style` into `tcell` drawing calls.
- Input and event polling are not core responsibilities. Applications should handle keyboard, mouse, and resize events directly with `tcell` or an equivalent input library, then call `terminal.Resize` when the terminal area changes.

## Porting Strategy

Port upstream behavior in small tested slices:

1. Start with pure core behavior tests from the upstream reference.
2. Add buffer snapshot helpers only when tests need them.
3. Port widget tests after the underlying `layout`, `style`, `text`, and `buffer` behavior is covered.
4. Keep tcell out of core packages so tests remain headless.

## Verification Gates

Every behavior change should pass:

```sh
go test ./...
golangci-lint run ./...
```

The pre-commit gate runs:

```sh
go fix ./...
go test ./...
golangci-lint run ./...
```
