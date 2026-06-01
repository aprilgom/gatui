package style

type Modifier uint16

const (
	ModifierBold Modifier = 1 << iota
	ModifierDim
	ModifierItalic
	ModifierUnderlined
	ModifierReversed
)
