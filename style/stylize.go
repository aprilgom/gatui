package style

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
