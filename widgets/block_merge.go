package widgets

import "fmt"

type MergeStrategy uint8

const (
	MergeStrategyReplace MergeStrategy = iota
	MergeStrategyExact
	MergeStrategyFuzzy
)

func (s MergeStrategy) String() string {
	switch s {
	case MergeStrategyReplace:
		return "Replace"
	case MergeStrategyExact:
		return "Exact"
	case MergeStrategyFuzzy:
		return "Fuzzy"
	default:
		return fmt.Sprintf("MergeStrategy(%d)", s)
	}
}

type borderLineStyle uint8

const (
	borderLineNothing borderLineStyle = iota
	borderLinePlain
	borderLineRounded
	borderLineDouble
	borderLineThick
	borderLineDoubleDash
	borderLineDoubleDashThick
	borderLineTripleDash
	borderLineTripleDashThick
	borderLineQuadrupleDash
	borderLineQuadrupleDashThick
)

type borderSymbol struct {
	right borderLineStyle
	up    borderLineStyle
	left  borderLineStyle
	down  borderLineStyle
}

func newBorderSymbol(right, up, left, down borderLineStyle) borderSymbol {
	return borderSymbol{right: right, up: up, left: left, down: down}
}

var borderSymbols = map[string]borderSymbol{
	"─": newBorderSymbol(borderLinePlain, borderLineNothing, borderLinePlain, borderLineNothing),
	"━": newBorderSymbol(borderLineThick, borderLineNothing, borderLineThick, borderLineNothing),
	"│": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineNothing, borderLinePlain),
	"┃": newBorderSymbol(borderLineNothing, borderLineThick, borderLineNothing, borderLineThick),
	"┄": newBorderSymbol(borderLineTripleDash, borderLineNothing, borderLineTripleDash, borderLineNothing),
	"┅": newBorderSymbol(borderLineTripleDashThick, borderLineNothing, borderLineTripleDashThick, borderLineNothing),
	"┆": newBorderSymbol(borderLineNothing, borderLineTripleDash, borderLineNothing, borderLineTripleDash),
	"┇": newBorderSymbol(borderLineNothing, borderLineTripleDashThick, borderLineNothing, borderLineTripleDashThick),
	"┈": newBorderSymbol(borderLineQuadrupleDash, borderLineNothing, borderLineQuadrupleDash, borderLineNothing),
	"┉": newBorderSymbol(borderLineQuadrupleDashThick, borderLineNothing, borderLineQuadrupleDashThick, borderLineNothing),
	"┊": newBorderSymbol(borderLineNothing, borderLineQuadrupleDash, borderLineNothing, borderLineQuadrupleDash),
	"┋": newBorderSymbol(borderLineNothing, borderLineQuadrupleDashThick, borderLineNothing, borderLineQuadrupleDashThick),
	"┌": newBorderSymbol(borderLinePlain, borderLineNothing, borderLineNothing, borderLinePlain),
	"┍": newBorderSymbol(borderLineThick, borderLineNothing, borderLineNothing, borderLinePlain),
	"┎": newBorderSymbol(borderLinePlain, borderLineNothing, borderLineNothing, borderLineThick),
	"┏": newBorderSymbol(borderLineThick, borderLineNothing, borderLineNothing, borderLineThick),
	"┐": newBorderSymbol(borderLineNothing, borderLineNothing, borderLinePlain, borderLinePlain),
	"┑": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineThick, borderLinePlain),
	"┒": newBorderSymbol(borderLineNothing, borderLineNothing, borderLinePlain, borderLineThick),
	"┓": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineThick, borderLineThick),
	"└": newBorderSymbol(borderLinePlain, borderLinePlain, borderLineNothing, borderLineNothing),
	"┕": newBorderSymbol(borderLineThick, borderLinePlain, borderLineNothing, borderLineNothing),
	"┖": newBorderSymbol(borderLinePlain, borderLineThick, borderLineNothing, borderLineNothing),
	"┗": newBorderSymbol(borderLineThick, borderLineThick, borderLineNothing, borderLineNothing),
	"┘": newBorderSymbol(borderLineNothing, borderLinePlain, borderLinePlain, borderLineNothing),
	"┙": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineThick, borderLineNothing),
	"┚": newBorderSymbol(borderLineNothing, borderLineThick, borderLinePlain, borderLineNothing),
	"┛": newBorderSymbol(borderLineNothing, borderLineThick, borderLineThick, borderLineNothing),
	"├": newBorderSymbol(borderLinePlain, borderLinePlain, borderLineNothing, borderLinePlain),
	"┝": newBorderSymbol(borderLineThick, borderLinePlain, borderLineNothing, borderLinePlain),
	"┞": newBorderSymbol(borderLinePlain, borderLineThick, borderLineNothing, borderLinePlain),
	"┟": newBorderSymbol(borderLinePlain, borderLinePlain, borderLineNothing, borderLineThick),
	"┠": newBorderSymbol(borderLinePlain, borderLineThick, borderLineNothing, borderLineThick),
	"┡": newBorderSymbol(borderLineThick, borderLineThick, borderLineNothing, borderLinePlain),
	"┢": newBorderSymbol(borderLineThick, borderLinePlain, borderLineNothing, borderLineThick),
	"┣": newBorderSymbol(borderLineThick, borderLineThick, borderLineNothing, borderLineThick),
	"┤": newBorderSymbol(borderLineNothing, borderLinePlain, borderLinePlain, borderLinePlain),
	"┥": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineThick, borderLinePlain),
	"┦": newBorderSymbol(borderLineNothing, borderLineThick, borderLinePlain, borderLinePlain),
	"┧": newBorderSymbol(borderLineNothing, borderLinePlain, borderLinePlain, borderLineThick),
	"┨": newBorderSymbol(borderLineNothing, borderLineThick, borderLinePlain, borderLineThick),
	"┩": newBorderSymbol(borderLineNothing, borderLineThick, borderLineThick, borderLinePlain),
	"┪": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineThick, borderLineThick),
	"┫": newBorderSymbol(borderLineNothing, borderLineThick, borderLineThick, borderLineThick),
	"┬": newBorderSymbol(borderLinePlain, borderLineNothing, borderLinePlain, borderLinePlain),
	"┭": newBorderSymbol(borderLinePlain, borderLineNothing, borderLineThick, borderLinePlain),
	"┮": newBorderSymbol(borderLineThick, borderLineNothing, borderLinePlain, borderLinePlain),
	"┯": newBorderSymbol(borderLineThick, borderLineNothing, borderLineThick, borderLinePlain),
	"┰": newBorderSymbol(borderLinePlain, borderLineNothing, borderLinePlain, borderLineThick),
	"┱": newBorderSymbol(borderLinePlain, borderLineNothing, borderLineThick, borderLineThick),
	"┲": newBorderSymbol(borderLineThick, borderLineNothing, borderLinePlain, borderLineThick),
	"┳": newBorderSymbol(borderLineThick, borderLineNothing, borderLineThick, borderLineThick),
	"┴": newBorderSymbol(borderLinePlain, borderLinePlain, borderLinePlain, borderLineNothing),
	"┵": newBorderSymbol(borderLinePlain, borderLinePlain, borderLineThick, borderLineNothing),
	"┶": newBorderSymbol(borderLineThick, borderLinePlain, borderLinePlain, borderLineNothing),
	"┷": newBorderSymbol(borderLineThick, borderLinePlain, borderLineThick, borderLineNothing),
	"┸": newBorderSymbol(borderLinePlain, borderLineThick, borderLinePlain, borderLineNothing),
	"┹": newBorderSymbol(borderLinePlain, borderLineThick, borderLineThick, borderLineNothing),
	"┺": newBorderSymbol(borderLineThick, borderLineThick, borderLinePlain, borderLineNothing),
	"┻": newBorderSymbol(borderLineThick, borderLineThick, borderLineThick, borderLineNothing),
	"┼": newBorderSymbol(borderLinePlain, borderLinePlain, borderLinePlain, borderLinePlain),
	"┽": newBorderSymbol(borderLinePlain, borderLinePlain, borderLineThick, borderLinePlain),
	"┾": newBorderSymbol(borderLineThick, borderLinePlain, borderLinePlain, borderLinePlain),
	"┿": newBorderSymbol(borderLineThick, borderLinePlain, borderLineThick, borderLinePlain),
	"╀": newBorderSymbol(borderLinePlain, borderLineThick, borderLinePlain, borderLinePlain),
	"╁": newBorderSymbol(borderLinePlain, borderLinePlain, borderLinePlain, borderLineThick),
	"╂": newBorderSymbol(borderLinePlain, borderLineThick, borderLinePlain, borderLineThick),
	"╃": newBorderSymbol(borderLinePlain, borderLineThick, borderLineThick, borderLinePlain),
	"╄": newBorderSymbol(borderLineThick, borderLineThick, borderLinePlain, borderLinePlain),
	"╅": newBorderSymbol(borderLinePlain, borderLinePlain, borderLineThick, borderLineThick),
	"╆": newBorderSymbol(borderLineThick, borderLinePlain, borderLinePlain, borderLineThick),
	"╇": newBorderSymbol(borderLineThick, borderLineThick, borderLineThick, borderLinePlain),
	"╈": newBorderSymbol(borderLineThick, borderLinePlain, borderLineThick, borderLineThick),
	"╉": newBorderSymbol(borderLinePlain, borderLineThick, borderLineThick, borderLineThick),
	"╊": newBorderSymbol(borderLineThick, borderLineThick, borderLinePlain, borderLineThick),
	"╋": newBorderSymbol(borderLineThick, borderLineThick, borderLineThick, borderLineThick),
	"╌": newBorderSymbol(borderLineDoubleDash, borderLineNothing, borderLineDoubleDash, borderLineNothing),
	"╍": newBorderSymbol(borderLineDoubleDashThick, borderLineNothing, borderLineDoubleDashThick, borderLineNothing),
	"╎": newBorderSymbol(borderLineNothing, borderLineDoubleDash, borderLineNothing, borderLineDoubleDash),
	"╏": newBorderSymbol(borderLineNothing, borderLineDoubleDashThick, borderLineNothing, borderLineDoubleDashThick),
	"═": newBorderSymbol(borderLineDouble, borderLineNothing, borderLineDouble, borderLineNothing),
	"║": newBorderSymbol(borderLineNothing, borderLineDouble, borderLineNothing, borderLineDouble),
	"╒": newBorderSymbol(borderLineDouble, borderLineNothing, borderLineNothing, borderLinePlain),
	"╓": newBorderSymbol(borderLinePlain, borderLineNothing, borderLineNothing, borderLineDouble),
	"╔": newBorderSymbol(borderLineDouble, borderLineNothing, borderLineNothing, borderLineDouble),
	"╕": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineDouble, borderLinePlain),
	"╖": newBorderSymbol(borderLineNothing, borderLineNothing, borderLinePlain, borderLineDouble),
	"╗": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineDouble, borderLineDouble),
	"╘": newBorderSymbol(borderLineDouble, borderLinePlain, borderLineNothing, borderLineNothing),
	"╙": newBorderSymbol(borderLinePlain, borderLineDouble, borderLineNothing, borderLineNothing),
	"╚": newBorderSymbol(borderLineDouble, borderLineDouble, borderLineNothing, borderLineNothing),
	"╛": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineDouble, borderLineNothing),
	"╜": newBorderSymbol(borderLineNothing, borderLineDouble, borderLinePlain, borderLineNothing),
	"╝": newBorderSymbol(borderLineNothing, borderLineDouble, borderLineDouble, borderLineNothing),
	"╞": newBorderSymbol(borderLineDouble, borderLinePlain, borderLineNothing, borderLinePlain),
	"╟": newBorderSymbol(borderLinePlain, borderLineDouble, borderLineNothing, borderLineDouble),
	"╠": newBorderSymbol(borderLineDouble, borderLineDouble, borderLineNothing, borderLineDouble),
	"╡": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineDouble, borderLinePlain),
	"╢": newBorderSymbol(borderLineNothing, borderLineDouble, borderLinePlain, borderLineDouble),
	"╣": newBorderSymbol(borderLineNothing, borderLineDouble, borderLineDouble, borderLineDouble),
	"╤": newBorderSymbol(borderLineDouble, borderLineNothing, borderLineDouble, borderLinePlain),
	"╥": newBorderSymbol(borderLinePlain, borderLineNothing, borderLinePlain, borderLineDouble),
	"╦": newBorderSymbol(borderLineDouble, borderLineNothing, borderLineDouble, borderLineDouble),
	"╧": newBorderSymbol(borderLineDouble, borderLinePlain, borderLineDouble, borderLineNothing),
	"╨": newBorderSymbol(borderLinePlain, borderLineDouble, borderLinePlain, borderLineNothing),
	"╩": newBorderSymbol(borderLineDouble, borderLineDouble, borderLineDouble, borderLineNothing),
	"╪": newBorderSymbol(borderLineDouble, borderLinePlain, borderLineDouble, borderLinePlain),
	"╫": newBorderSymbol(borderLinePlain, borderLineDouble, borderLinePlain, borderLineDouble),
	"╬": newBorderSymbol(borderLineDouble, borderLineDouble, borderLineDouble, borderLineDouble),
	"╭": newBorderSymbol(borderLineRounded, borderLineNothing, borderLineNothing, borderLineRounded),
	"╮": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineRounded, borderLineRounded),
	"╯": newBorderSymbol(borderLineNothing, borderLineRounded, borderLineRounded, borderLineNothing),
	"╰": newBorderSymbol(borderLineRounded, borderLineRounded, borderLineNothing, borderLineNothing),
	"╴": newBorderSymbol(borderLineNothing, borderLineNothing, borderLinePlain, borderLineNothing),
	"╵": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineNothing, borderLineNothing),
	"╶": newBorderSymbol(borderLinePlain, borderLineNothing, borderLineNothing, borderLineNothing),
	"╷": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineNothing, borderLinePlain),
	"╸": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineThick, borderLineNothing),
	"╹": newBorderSymbol(borderLineNothing, borderLineThick, borderLineNothing, borderLineNothing),
	"╺": newBorderSymbol(borderLineThick, borderLineNothing, borderLineNothing, borderLineNothing),
	"╻": newBorderSymbol(borderLineNothing, borderLineNothing, borderLineNothing, borderLineThick),
	"╼": newBorderSymbol(borderLineThick, borderLineNothing, borderLinePlain, borderLineNothing),
	"╽": newBorderSymbol(borderLineNothing, borderLinePlain, borderLineNothing, borderLineThick),
	"╾": newBorderSymbol(borderLinePlain, borderLineNothing, borderLineThick, borderLineNothing),
	"╿": newBorderSymbol(borderLineNothing, borderLineThick, borderLineNothing, borderLinePlain),
}

var borderSymbolsBySegments = makeBorderSymbolsBySegments()

func makeBorderSymbolsBySegments() map[borderSymbol]string {
	bySegments := make(map[borderSymbol]string, len(borderSymbols))
	for symbol, segments := range borderSymbols {
		bySegments[segments] = symbol
	}
	return bySegments
}

func mergeBorderSymbols(strategy MergeStrategy, prev, next string) string {
	if strategy == MergeStrategyReplace {
		return next
	}
	prevSymbol, prevOK := borderSymbols[prev]
	nextSymbol, nextOK := borderSymbols[next]
	switch {
	case prevOK && nextOK:
		merged := prevSymbol.merge(nextSymbol)
		if strategy == MergeStrategyFuzzy {
			merged = merged.fuzzy(nextSymbol)
		}
		if symbol, ok := borderSymbolsBySegments[merged]; ok {
			return symbol
		}
		return next
	case !prevOK && nextOK:
		return prev
	default:
		return next
	}
}

func (s borderSymbol) merge(other borderSymbol) borderSymbol {
	return newBorderSymbol(
		mergeLineStyle(s.right, other.right),
		mergeLineStyle(s.up, other.up),
		mergeLineStyle(s.left, other.left),
		mergeLineStyle(s.down, other.down),
	)
}

func mergeLineStyle(prev, next borderLineStyle) borderLineStyle {
	if next == borderLineNothing {
		return prev
	}
	return next
}

func (s borderSymbol) fuzzy(other borderSymbol) borderSymbol {
	if !s.isStraight() {
		s = s.replace(borderLineDoubleDash, borderLinePlain)
		s = s.replace(borderLineTripleDash, borderLinePlain)
		s = s.replace(borderLineQuadrupleDash, borderLinePlain)
		s = s.replace(borderLineDoubleDashThick, borderLineThick)
		s = s.replace(borderLineTripleDashThick, borderLineThick)
		s = s.replace(borderLineQuadrupleDashThick, borderLineThick)
	}
	if !s.isCorner() {
		s = s.replace(borderLineRounded, borderLinePlain)
	}
	if s.contains(borderLineDouble) && s.contains(borderLineThick) {
		if other.contains(borderLineDouble) {
			s = s.replace(borderLineThick, borderLineDouble)
		} else {
			s = s.replace(borderLineDouble, borderLineThick)
		}
	}
	if _, ok := borderSymbolsBySegments[s]; !ok {
		if other.contains(borderLineDouble) {
			s = s.replace(borderLinePlain, borderLineDouble)
		} else {
			s = s.replace(borderLineDouble, borderLinePlain)
		}
	}
	return s
}

func (s borderSymbol) isStraight() bool {
	return s.up == s.down && s.left == s.right &&
		(s.up == borderLineNothing || s.left == borderLineNothing)
}

func (s borderSymbol) isCorner() bool {
	switch {
	case s.down == borderLineNothing && s.left == borderLineNothing:
		return s.up == s.right
	case s.up == borderLineNothing && s.left == borderLineNothing:
		return s.right == s.down
	case s.up == borderLineNothing && s.right == borderLineNothing:
		return s.down == s.left
	case s.right == borderLineNothing && s.down == borderLineNothing:
		return s.up == s.left
	default:
		return false
	}
}

func (s borderSymbol) contains(style borderLineStyle) bool {
	return s.up == style || s.right == style || s.down == style || s.left == style
}

func (s borderSymbol) replace(from, to borderLineStyle) borderSymbol {
	if s.up == from {
		s.up = to
	}
	if s.right == from {
		s.right = to
	}
	if s.down == from {
		s.down = to
	}
	if s.left == from {
		s.left = to
	}
	return s
}
