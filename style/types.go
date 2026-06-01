package style

type Color int

const (
	Default Color = iota
	Reset
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	LightBlue
	LightGreen
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

func ResetStyle() Style {
	return Style{Foreground: Reset, Background: Reset}
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

func (s Style) Patch(other Style) Style {
	if other.Foreground != Default {
		s.Foreground = other.Foreground
	}
	if other.Background != Default {
		s.Background = other.Background
	}
	s.Modifiers |= other.Modifiers
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
