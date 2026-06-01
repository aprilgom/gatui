# text Package Guide

`text` owns styled text primitives: `Span`, `Line`, and `Text`.

Key files: `grapheme.go`, `span.go`, `line.go`, `text.go`, `render.go`, `masked.go`, `../style/style.go`, `../widgets/paragraph.go`.

## Setup

Install Go 1.23+ and `golangci-lint`, then work from the repository root.

Allowed dependencies:

- `buffer`
- `layout`
- `style`

Avoid dependencies on `widgets` or terminal backends.

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
