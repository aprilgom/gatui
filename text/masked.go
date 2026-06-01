package text

import "strings"

type Masked struct {
	value string
	mask  rune
}

func NewMasked(value string, mask rune) Masked {
	return Masked{value: value, mask: mask}
}

func (m Masked) MaskChar() rune {
	return m.mask
}

func (m Masked) Value() string {
	return strings.Repeat(string(m.mask), len([]rune(m.value)))
}

func (m Masked) String() string {
	return m.Value()
}

func (m Masked) Text() Text {
	return FromString(m.Value())
}
