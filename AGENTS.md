# gatui Codex Guide

## Project Scope

`gatui` is a Go port of Ratatui's core terminal UI model. It is intentionally small and grows from upstream behavior tests.

- `layout`: rectangles, margins, offsets, constraints, and area splitting. See `layout/rect.go`, `layout/constraint.go`, and `layout/layout.go`.
- `buffer`: cell grid used by widgets before backend flush. See `buffer/cell.go`, `buffer/buffer.go`, and `buffer/diff.go`.
- `style`: colors, modifiers, and chainable style helpers. See `style/color.go`, `style/modifier.go`, and `style/style.go`.
- `text`: styled graphemes, spans, lines, and multi-line text. See `text/grapheme.go`, `text/span.go`, `text/line.go`, and `text/text.go`.
- `symbols`: reusable terminal symbol sets and glyph helpers. See `symbols/border.go`, `symbols/bar.go`, `symbols/sparkline.go`, `symbols/canvas.go`, and `symbols/scrollbar.go`.
- `widgets`: renderable UI components such as `Paragraph`, `Block`, `Gauge`, `Tabs`, and `Clear`. See widget-specific files under `widgets/`.

## Setup

Install Go 1.23+ and `golangci-lint`, then run:

```sh
go test ./...
golangci-lint run ./...
.githooks/pre-commit
```

Enable the shared hook with:

```sh
git config core.hooksPath .githooks
```

## Development Workflow

Port behavior test-first: read the upstream behavior, add the equivalent Go test, confirm the expected failure, implement the smallest API, then run `.githooks/pre-commit`.

Done criteria: changed behavior has tests, `go test ./...` passes, `golangci-lint run ./...` reports `0 issues`, and final reports include exact verification commands.

## Architecture Map

```mermaid
flowchart TD
    layout[layout<br/>rect/constraint/layout]
    style[style<br/>color/modifier/style]
    text[text<br/>grapheme/span/line/text/render]
    symbols[symbols<br/>border/bar/canvas/scrollbar]
    widgets[widgets<br/>component files]
    buffer[buffer<br/>cell/buffer/diff]
    tests[layout/rect_test.go]
    api[api_contract_test.go]

    style --> text
    style --> buffer
    layout --> buffer
    layout --> widgets
    symbols --> widgets
    text --> widgets
    widgets --> buffer
    tests --> layout
    tests --> buffer
    api --> layout
    api --> buffer
    api --> style
    api --> text
    api --> widgets
```

See `ARCHITECTURE.md` for the dependency graph and package ownership.

## Cross-Module Dependencies

Note: `layout` and `style` are foundational. `buffer` depends on `layout` and `style`. `text` depends on `buffer`, `layout`, and `style` for rendering styled text into buffers. `symbols` stays independent and provides shared terminal glyph sets. `widgets` depends on `buffer`, `layout`, `style`, `symbols`, and `text`. Backend code belongs under backend packages and must not leak into core packages.

## Safety Boundaries

- Treat external reference trees as read-only unless the user explicitly asks otherwise.
- Do not commit generated local audit dashboards or temporary outputs unless the user asks.
- Do not use destructive Git commands such as `git reset --hard` or `git checkout --` without explicit approval.
- Preserve user changes in a dirty worktree. If a touched file has unexpected edits, inspect and work with them.
- Do not add secrets. Keep `.env`, `.env.*`, credentials, tokens, and local logs out of Git.
- Prefer additive APIs while the port is young; avoid broad renames unless tests and callers are updated together.

## Known Gaps

- Ratatui parity is partial; current ported coverage is represented by package tests under `layout`, `buffer`, `style`, `text`, `terminal`, and `widgets`.
