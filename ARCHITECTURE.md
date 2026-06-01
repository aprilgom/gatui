# gatui Architecture

`gatui` mirrors Ratatui's separation between pure rendering logic and terminal backends. Core packages should remain testable without a terminal.

```mermaid
flowchart TD
    layout[layout<br/>Rect, Margin, Offset, Constraint]
    style[style<br/>Color, Modifier, Style]
    text[text<br/>Span, Line, Text]
    widgets[widgets<br/>Widget, Paragraph, Block, Clear]
    buffer[buffer<br/>Buffer, Cell]
    backend[future backend/tcell<br/>terminal flush and events]

    style --> text
    style --> buffer
    layout --> buffer
    layout --> widgets
    text --> widgets
    widgets --> buffer
    buffer --> backend
```

## Package Responsibilities

- `layout` owns geometry and area splitting. It must not depend on widgets or terminal backends.
- `style` owns style value types and chainable style helpers. It must stay backend-neutral.
- `text` owns styled textual content. It may depend on `style`, but not on `buffer` or `widgets`.
- `buffer` owns the off-screen cell grid. It may depend on `layout` and `style`.
- `widgets` owns renderable components. Widgets write into `buffer.Buffer` inside a `layout.Rect`.
- Future `backend/tcell` should translate `buffer.Cell` and `style.Style` into `tcell` drawing calls.

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
