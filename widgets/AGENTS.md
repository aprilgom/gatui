# widgets Package Guide

`widgets` owns renderable UI components.

Key files: `types.go`, `../buffer/types.go`, `../layout/types.go`.

## Setup

Install Go 1.23+ and `golangci-lint`, then work from the repository root.

Allowed dependencies:

- `buffer`
- `layout`
- `text`
- `style` when needed

Avoid terminal backend dependencies. Widgets render into `*buffer.Buffer`; backends flush buffers later.

Verification:

```sh
go test ./widgets ./buffer ./layout ./...
```

## Done Criteria

Complete changes only after focused widget tests and full tests pass.

## Known Caveats

Known failure risk: widgets are integration points for `layout`, `text`, `style`, and `buffer`.

Do not add backend, secret, generated, or vendor concerns here.

## Cross-Module Dependencies

See also: [root guide](../AGENTS.md), [architecture](../ARCHITECTURE.md).

When porting Ratatui widget tests, ensure lower-level `layout`, `text`, `style`, and `buffer` behavior is already covered or add focused tests first.
