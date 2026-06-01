package style

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
