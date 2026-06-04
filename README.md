# gatui

`gatui` is a Go port of Ratatui's terminal UI model.

The project currently focuses on core, backend-independent packages:

- `layout`
- `style`
- `text`
- `buffer`
- `symbols`
- `widgets`

The `tcell` backend handles terminal drawing, while input and event polling stay outside core. Applications should read keyboard, mouse, and resize events directly with `tcell` or an equivalent input library.

## Quick Start

```sh
go test ./...
golangci-lint run ./...
```

Enable the shared pre-commit hook:

```sh
git config core.hooksPath .githooks
```

Then run:

```sh
.githooks/pre-commit
```

## Examples

Runnable examples live under `examples/`. Run them from the module root:

```sh
go run ./examples/hello-world
go run ./examples/minimal
go run ./examples/app-gauge
go run ./examples/app-table
go run ./examples/app-list
go run ./examples/app-canvas
go run ./examples/widget-gauge
go run ./examples/widget-table
go run ./examples/widget-list
go run ./examples/widget-canvas
```

The app examples demonstrate small interactive terminal programs. The widget
examples focus on individual widget rendering. Most examples exit with `q`,
`Esc`, or `Ctrl+C`; list and table examples also support arrow keys or `hjkl`.

For LLMs, coding agents, and humans building a new app, start with
[docs/llm-usage.md](docs/llm-usage.md). It includes the minimal app skeleton,
widget cookbook snippets, common pitfalls, and Ratatui-to-Gatui translation
rules.

## Current Status

Implemented slices:

- Package structure with type definitions split into focused files per package.
- Public API contract smoke test.
- Go port coverage represented by package tests across `layout`, `buffer`, `style`, `text`, `terminal`, and `widgets`.
- Runnable Ratatui-inspired examples under `examples/`.
- P0 ROI implementation pass completed; the 215 remaining missing-test actions are rebucketed across P1-P4 in the correspondence backlog.

See [AGENTS.md](AGENTS.md) for Codex workflow guidance, [ARCHITECTURE.md](ARCHITECTURE.md) for the package map, and [../correspondence/roi/README.md](../correspondence/roi/README.md) for the current ROI backlog.
