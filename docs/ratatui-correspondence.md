# Ratatui to Gatui Correspondence

This document maps Ratatui Rust API paths to Gatui Go packages and types. It is written for LLMs, coding agents, and Go users porting code from Ratatui.

## Agent Porting Rules

- Search this document for the exact Ratatui path first, for example `ratatui::widgets::Paragraph`.
- Prefer the listed Gatui package, type, or function over inventing a wrapper.
- If a Ratatui API is not listed here, inspect local Gatui code and tests before assuming it exists.
- Treat unlisted behavior as unported unless the implementation or tests show an equivalent.
- Keep input polling outside Gatui core. Applications should use `tcell` or another input library directly, then call Gatui terminal APIs for drawing and resize handling.
- Translate Rust builders, traits, and ownership patterns into the existing Go API style. Do not copy Rust method chains mechanically if the Gatui type exposes ordinary Go fields, methods, or constructors.

## Package Map

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui` | `gatui` | Root package documentation and porting orientation. |
| `ratatui::layout` | `gatui/layout` | Geometry, constraints, and area splitting. |
| `ratatui::style` | `gatui/style` | Colors, modifiers, and style values. |
| `ratatui::text` | `gatui/text` | Styled text primitives. |
| `ratatui::buffer` | `gatui/buffer` | Off-screen cell grid and diffing. |
| `ratatui::symbols` | `gatui/symbols` | Terminal glyph sets and drawing symbols. |
| `ratatui::widgets` | `gatui/widgets` | Renderable UI components. |
| `ratatui::Terminal` | `gatui/terminal` | Terminal drawing orchestration. |
| `ratatui::Frame` | `gatui/terminal.Frame` | Single-frame rendering context. |
| `ratatui::backend` | `gatui/terminal.Backend` plus backend packages | Backend interface and concrete implementations. |
| `ratatui_crossterm` | no direct Gatui package | Gatui currently provides a `tcell` backend instead of a crossterm backend. |
| `ratatui::init`, `ratatui::restore`, `ratatui::run` | applications own setup and teardown | Use `backend/tcell` and `terminal` directly; input/event loops stay outside core. |

## Layout

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui::layout::Rect` | `layout.Rect` | Rectangle in terminal cells. |
| `ratatui::layout::Position` | check `layout` primitives | Use current Gatui geometry primitives; do not assume a one-to-one type unless exported. |
| `ratatui::layout::Size` | check `layout` primitives | Use current Gatui geometry primitives. |
| `ratatui::layout::Margin` | `layout.Margin` | Inner/outer area adjustments. |
| `ratatui::layout::Offset` | `layout.Offset` | Positional offset. |
| `ratatui::layout::Constraint` | `layout.Constraint` | Length, min, max, percentage, ratio, and fill-like layout constraints as implemented by Gatui. |
| `ratatui::layout::Layout` | `layout.Layout` | Splits a `Rect` into child areas. |
| `ratatui::layout::Direction` | `layout.Direction` | Horizontal or vertical splitting. |
| `ratatui::layout::Flex` | check `layout` | Use only if exported in Gatui. |
| `ratatui::layout::Alignment` | check `layout` or widget-specific APIs | Some alignment behavior may live on widgets. |

## Style

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui::style::Style` | `style.Style` | Value describing foreground, background, underline, and modifiers. |
| `ratatui::style::Color` | `style.Color` | Terminal color value. |
| `ratatui::style::Modifier` | `style.Modifier` | Bold, italic, underline, reversed, and related text modifiers. |
| `ratatui::style::Stylize` | `style` helper methods | Gatui may expose Go-style helpers instead of a Rust trait. |
| `ratatui::style::palette` | not a guaranteed Gatui mapping | Check `style` before porting palette constants. |

## Text

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui::text::Span` | `text.Span` | Styled inline text segment. |
| `ratatui::text::Line` | `text.Line` | Sequence of spans rendered as one line. |
| `ratatui::text::Text` | `text.Text` | Multi-line styled text. |
| `ratatui::text::Masked` | `text.Masked` | Masked text representation if supported by the current package. |
| `ratatui::text::StyledGrapheme` | `text.StyledGrapheme` | Grapheme plus style where exported. |

## Buffer

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui::buffer::Buffer` | `buffer.Buffer` | Off-screen grid of terminal cells. |
| `ratatui::buffer::Cell` | `buffer.Cell` | A styled terminal cell. |
| `ratatui::buffer::Cell::set_symbol` | `buffer.Cell` methods | Use the existing Go methods; names follow Go conventions. |
| `ratatui::buffer::Buffer::diff` | `buffer` diff APIs | Check `buffer/diff.go` for the exact exported shape. |

## Symbols

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui::symbols` | `symbols` | Shared glyph sets. |
| `ratatui::symbols::border` | `symbols` border sets | Border glyph definitions for widgets. |
| `ratatui::symbols::line` | `symbols` line sets | Line drawing glyphs. |
| `ratatui::symbols::bar` | `symbols` bar sets | Bar chart and gauge glyphs. |
| `ratatui::symbols::scrollbar` | `symbols` scrollbar sets | Scrollbar glyphs. |
| `ratatui::symbols::Marker` | `symbols.CanvasMarker` or canvas symbols | Canvas marker naming differs in Gatui. |

## Widgets

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui::widgets::Widget` | `widgets.Widget` | Interface for renderable widgets. |
| `ratatui::widgets::StatefulWidget` | `widgets.StatefulWidget` or state-specific APIs | Check widget package interfaces before porting. |
| `ratatui::widgets::Block` | `widgets.Block` | Borders, titles, padding, and visual framing. |
| `ratatui::widgets::Paragraph` | `widgets.Paragraph` | Text rendering widget. |
| `ratatui::widgets::Clear` | `widgets.Clear` | Clears an area before drawing. |
| `ratatui::widgets::Gauge` | `widgets.Gauge` | Progress gauge. |
| `ratatui::widgets::LineGauge` | `widgets.LineGauge` | Compact line gauge. |
| `ratatui::widgets::BarChart` | `widgets.BarChart` | Bar chart widget. |
| `ratatui::widgets::Sparkline` | `widgets.Sparkline` | Sparkline widget. |
| `ratatui::widgets::List` | `widgets.List` | List widget. |
| `ratatui::widgets::ListState` | `widgets.ListState` | List selection/scroll state where exported. |
| `ratatui::widgets::Table` | `widgets.Table` | Table widget. |
| `ratatui::widgets::TableState` | `widgets.TableState` | Table selection/scroll state where exported. |
| `ratatui::widgets::Tabs` | `widgets.Tabs` | Tabs widget. |
| `ratatui::widgets::Scrollbar` | `widgets.Scrollbar` | Scrollbar widget. |
| `ratatui::widgets::ScrollbarState` | `widgets.ScrollbarState` | Scrollbar state where exported. |
| `ratatui::widgets::Canvas` | `widgets.Canvas` | Canvas widget and drawing world. |
| `ratatui::widgets::Calendar` | `widgets.Calendar` | Calendar widget. |
| `ratatui::widgets::Fill` | `widgets.Fill` | Fills an area with a symbol/style. |
| `ratatui::widgets::Shadow` | `widgets.Shadow` | Shadow rendering helper/widget. |

## Terminal and Backends

| Ratatui Rust path | Gatui Go path | Notes |
| --- | --- | --- |
| `ratatui::Terminal` | `terminal.Terminal` | Owns draw orchestration over a backend. |
| `ratatui::Frame` | `terminal.Frame` | Passed to draw functions for one render pass. |
| `ratatui::CompletedFrame` | check `terminal` | Only use if exported in Gatui. |
| `ratatui::Viewport` | `terminal.Viewport` | Terminal viewport mode where exported. |
| `ratatui::backend::Backend` | `terminal.Backend` | Backend interface implemented by concrete backends. |
| `ratatui::backend::TestBackend` | `terminal/testbackend` | Headless backend for tests. |
| `ratatui_crossterm::CrosstermBackend` | `backend/tcell` | Gatui uses tcell for terminal drawing. |

## Gatui-Only Helper Packages

| Gatui Go path | Purpose | Notes |
| --- | --- | --- |
| `gatui/textbuffer` | Helper functions for rendering `text` values into `buffer.Buffer`. | This is a Gatui helper package, not a direct Ratatui module. |

## Common Translation Patterns

### Rendering

Ratatui Rust code commonly renders inside `Terminal::draw` with a closure receiving `&mut Frame`. Gatui code should use the current `terminal` package draw API and render widgets into the `terminal.Frame` or underlying `buffer.Buffer` using existing widget interfaces.

### Builders

Ratatui often uses Rust builder chains:

```rust
Paragraph::new("hello").block(Block::bordered())
```

Do not assume Gatui has the same chain. Check the Go type and use the exported constructor, fields, or methods that already exist.

### Traits and Interfaces

Ratatui traits such as `Widget`, `StatefulWidget`, and backend traits correspond to Go interfaces where Gatui exposes them. Use the interface names from `widgets` and `terminal`; avoid introducing adapter interfaces unless a missing API is being intentionally ported.

### Input

Ratatui examples may show crossterm event polling. Gatui core does not own input polling. Go applications should read input through `tcell` or another library and then ask Gatui to redraw.
