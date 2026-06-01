package text

import (
	"unicode"

	"gatui/style"
)

const (
	nbsp = "\u00A0"
	zwsp = "\u200B"
)

type StyledGrapheme struct {
	Symbol string
	Style  style.Style
}

func NewStyledGrapheme(symbol string, graphemeStyle style.Style) StyledGrapheme {
	return StyledGrapheme{Symbol: symbol, Style: graphemeStyle}
}

func (g StyledGrapheme) IsWhitespace() bool {
	if g.Symbol == zwsp {
		return true
	}
	if g.Symbol == nbsp {
		return false
	}
	for _, r := range g.Symbol {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func (g StyledGrapheme) Fg(color style.Color) StyledGrapheme {
	g.Style = g.Style.Fg(color)
	return g
}

func (g StyledGrapheme) Bg(color style.Color) StyledGrapheme {
	g.Style = g.Style.Bg(color)
	return g
}

func (g StyledGrapheme) Bold() StyledGrapheme {
	g.Style = g.Style.AddModifier(style.ModifierBold)
	return g
}

func (g StyledGrapheme) Italic() StyledGrapheme {
	g.Style = g.Style.AddModifier(style.ModifierItalic)
	return g
}

func (g StyledGrapheme) Cyan() StyledGrapheme {
	return g.Fg(style.Cyan)
}
