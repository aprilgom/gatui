# gatui

`gatui` is a Go port of Ratatui's terminal UI model.

The project currently focuses on core, backend-independent packages:

- `layout`
- `style`
- `text`
- `buffer`
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

- Initial package structure.
- Public API contract smoke test.
- Go port coverage represented by `layout/rect_test.go` for `Rect.Inner`, `Rect.Outer`, `Rect.Offset`, `Rect.Intersection`, and `Rect.Clamp`.

See [AGENTS.md](AGENTS.md) for Codex workflow guidance and [ARCHITECTURE.md](ARCHITECTURE.md) for the package map.
