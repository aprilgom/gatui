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
	Gray
	DarkGray
	LightRed
	LightYellow
	LightMagenta
	LightCyan
	LightWhite
)

const (
	indexedColorBase Color = 1 << 8
	rgbColorBase     Color = 1 << 16
)

func ColorFromANSIColor(color Color) (Color, bool) {
	if color >= Black && color <= LightWhite {
		return color, true
	}
	return Default, false
}

func IndexedColor(index uint8) Color {
	color, _ := ColorFromIndexedColor(index)
	return color
}

func ColorFromIndexedColor(index uint8) (Color, bool) {
	return indexedColorBase + Color(index), true
}

func RGBColor(r, g, b uint8) Color {
	color, _ := ColorFromRGBColor(r, g, b)
	return color
}

func ColorFromRGBColor(r, g, b uint8) (Color, bool) {
	return rgbColorBase + Color(uint32(r)<<16|uint32(g)<<8|uint32(b)), true
}

func ColorFromUint32(value uint32) (Color, bool) {
	if value > 0xFFFFFF {
		return Default, false
	}
	return ColorFromRGBColor(uint8(value>>16), uint8(value>>8), uint8(value))
}
