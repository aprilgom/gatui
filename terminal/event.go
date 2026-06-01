package terminal

import "gatui/layout"

type EventType int

const (
	EventUnknown EventType = iota
	EventKey
	EventMouse
	EventResize
)

type Event interface {
	Type() EventType
}

type KeyCode int

const (
	KeyUnknown KeyCode = iota
	KeyRune
	KeyEnter
	KeyEsc
	KeyBackspace
	KeyTab
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown
	KeyDelete
	KeyInsert
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

type KeyModifier int

const (
	ModifierNone KeyModifier = 0
	ModifierCtrl KeyModifier = 1 << iota
	ModifierAlt
	ModifierShift
)

type KeyEvent struct {
	Code      KeyCode
	Rune      rune
	Modifiers KeyModifier
}

func (KeyEvent) Type() EventType {
	return EventKey
}

type MouseButton int

const (
	MouseButtonNone MouseButton = iota
	MouseButtonLeft
	MouseButtonRight
	MouseButtonMiddle
	MouseWheelUp
	MouseWheelDown
)

type MouseEvent struct {
	Position  layout.Position
	Button    MouseButton
	Modifiers KeyModifier
}

func (MouseEvent) Type() EventType {
	return EventMouse
}

type ResizeEvent struct {
	Size layout.Size
}

func (ResizeEvent) Type() EventType {
	return EventResize
}

type UnknownEvent struct{}

func (UnknownEvent) Type() EventType {
	return EventUnknown
}
