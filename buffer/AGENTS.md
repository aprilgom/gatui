# buffer Package Guide

`buffer` owns the off-screen cell grid used by widgets.

Key files: `types.go`, `../layout/types.go`, `../style/types.go`.

## Setup

Install Go 1.23+ and `golangci-lint`, then work from the repository root.

Dependencies allowed:

- `layout`
- `style`

Avoid dependencies on `widgets`, `text`, or terminal backends.

Verification:

```sh
go test ./buffer ./widgets ./...
```

## Done Criteria

Complete changes only after focused tests and `go test ./...` pass.

## Known Caveats

Known failure risk: changes here can break widget rendering snapshots because widgets write directly into `Buffer`.

Do not add backend, secret, generated, or vendor concerns here.

## Cross-Module Dependencies

See also: [root guide](../AGENTS.md), [architecture](../ARCHITECTURE.md).

When changing cell write behavior, add or update snapshot-style tests that inspect `Buffer.Lines()` or individual cells.
