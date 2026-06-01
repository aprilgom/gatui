# gatui Codex Guide

## Project Scope

`gatui` is a Go port of Ratatui's core terminal UI model. It is intentionally small and grows from upstream behavior tests.

- `layout`: rectangles, margins, offsets, constraints, and area splitting. See `layout/types.go`.
- `buffer`: cell grid used by widgets before backend flush. See `buffer/types.go`.
- `style`: colors, modifiers, and chainable style helpers. See `style/types.go`.
- `text`: spans, lines, and multi-line text. See `text/types.go`.
- `widgets`: renderable UI components such as `Paragraph`, `Block`, and `Clear`. See `widgets/types.go`.

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
    layout[layout/types.go]
    style[style/types.go]
    text[text/types.go]
    widgets[widgets/types.go]
    buffer[buffer/types.go]
    tests[layout/rect_test.go]
    api[api_contract_test.go]

    style --> text
    style --> buffer
    layout --> buffer
    layout --> widgets
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

Note: `layout` and `style` are foundational. `text` depends on `style` only. `buffer` depends on `layout` and `style`. `widgets` depends on `buffer`, `layout`, and `text`. Backend code should be added in a future backend package and must not leak into core packages.

## Safety Boundaries

- Treat external reference trees as read-only unless the user explicitly asks otherwise.
- Do not commit generated local audit dashboards or temporary outputs unless the user asks.
- Do not use destructive Git commands such as `git reset --hard` or `git checkout --` without explicit approval.
- Preserve user changes in a dirty worktree. If a touched file has unexpected edits, inspect and work with them.
- Do not add secrets. Keep `.env`, `.env.*`, credentials, tokens, and local logs out of Git.
- Prefer additive APIs while the port is young; avoid broad renames unless tests and callers are updated together.

## Known Gaps

- No terminal backend is implemented yet. `tcell` is planned as a backend, not a dependency of core packages.
- Unicode width/grapheme behavior is still minimal.
- Ratatui parity is partial; current ported coverage is represented by `layout/rect_test.go`.
