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

## Current Status

Implemented slices:

- Package structure with type definitions split into focused files per package.
- Public API contract smoke test.
- Go port coverage represented by package tests across `layout`, `buffer`, `style`, `text`, `terminal`, and `widgets`.
- P0 ROI implementation pass completed; the 215 remaining missing-test actions are rebucketed across P1-P4 in the correspondence backlog.

See [AGENTS.md](AGENTS.md) for Codex workflow guidance, [ARCHITECTURE.md](ARCHITECTURE.md) for the package map, and [../correspondence/roi/README.md](../correspondence/roi/README.md) for the current ROI backlog.
