package widgets

import "gatui/symbols"

type MergeStrategy = symbols.MergeStrategy

const (
	MergeStrategyReplace = symbols.MergeStrategyReplace
	MergeStrategyExact   = symbols.MergeStrategyExact
	MergeStrategyFuzzy   = symbols.MergeStrategyFuzzy
)

func mergeBorderSymbols(strategy MergeStrategy, prev, next string) string {
	return symbols.MergeBorderSymbols(strategy, prev, next)
}
