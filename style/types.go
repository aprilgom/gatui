package style

type Color int

const (
	Default Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type Modifier uint16

const (
	ModifierBold Modifier = 1 << iota
	ModifierDim
	ModifierItalic
	ModifierUnderlined
	ModifierReversed
)

type Style struct {
	Foreground Color
	Background Color
	Modifiers  Modifier
}

func NewStyle() Style {
	return Style{Foreground: Default, Background: Default}
}

func (s Style) Fg(color Color) Style {
	s.Foreground = color
	return s
}

func (s Style) Bg(color Color) Style {
	s.Background = color
	return s
}

func (s Style) AddModifier(modifier Modifier) Style {
	s.Modifiers |= modifier
	return s
}

type Styled[T any] struct {
	Value T
	Style Style
}

type Stylize[T any] interface {
	Fg(Color) T
	Bg(Color) T
	Bold() T
	Italic() T
}
