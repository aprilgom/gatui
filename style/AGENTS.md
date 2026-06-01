# style Package Guide

`style` owns backend-neutral color, modifier, style, `Styled`, and `Stylize` concepts.

Key files: `types.go`, `../text/types.go`, `../buffer/types.go`.

## Setup

Install Go 1.23+ and `golangci-lint`, then work from the repository root.

Do not import `tcell` or any terminal backend here. Backend-specific conversion belongs under a future backend package.

Verification:

```sh
go test ./style ./text ./buffer ./widgets ./...
```

## Done Criteria

Complete changes only after style users still compile and full tests pass.

## Known Caveats

Known failure risk: style values are shared by `text`, `buffer`, and future backend conversion.

Do not add backend, secret, generated, or vendor concerns here.

## Cross-Module Dependencies

See also: [root guide](../AGENTS.md), [architecture](../ARCHITECTURE.md).

When adding modifiers or colors, keep names close to Ratatui unless Go naming requires adjustment.
