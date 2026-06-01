# text Package Guide

`text` owns styled text primitives: `Span`, `Line`, and `Text`.

Key files: `types.go`, `../style/types.go`, `../widgets/types.go`.

## Setup

Install Go 1.23+ and `golangci-lint`, then work from the repository root.

Allowed dependency:

- `style`

Avoid dependencies on `buffer`, `layout`, or `widgets`.

Verification:

```sh
go test ./text ./widgets ./...
```

## Done Criteria

Complete changes only after text users still compile and full tests pass.

## Known Caveats

Unicode width and grapheme behavior is incomplete. Make changes test-first and isolate any future width dependency behind small helpers.

Do not add backend, secret, generated, or vendor concerns here.

## Cross-Module Dependencies

See also: [root guide](../AGENTS.md), [architecture](../ARCHITECTURE.md).
