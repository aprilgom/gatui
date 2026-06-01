# layout Package Guide

`layout` owns geometry and area splitting.

Key files: `types.go`, `rect_test.go`, `../api_contract_test.go`.

Keep this package dependency-free unless there is a strong reason. It should remain usable by `buffer` and `widgets` without import cycles.

## Setup

Install Go 1.23+ and `golangci-lint`, then work from the repository root.

Verification:

```sh
go test ./layout
go test ./...
```

## Done Criteria

Complete changes only after the focused layout test and full test suite pass.

## Known Caveats

Known failure risk: `Rect` behavior affects `buffer` bounds and every widget render area.

Do not add backend, secret, generated, or vendor concerns here.

## Cross-Module Dependencies

See also: [root guide](../AGENTS.md), [architecture](../ARCHITECTURE.md).

Port Ratatui geometry behavior test-first. Prefer table tests for scalar geometry and buffer snapshot tests for visual overlap behavior.
