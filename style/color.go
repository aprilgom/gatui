package style

import (
	"fmt"
	"strconv"
	"strings"
)

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

func (c Color) String() string {
	switch c {
	case Default:
		return "Default"
	case Reset:
		return "Reset"
	case Black:
		return "Black"
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Yellow:
		return "Yellow"
	case Blue:
		return "Blue"
	case Magenta:
		return "Magenta"
	case Cyan:
		return "Cyan"
	case White:
		return "White"
	case LightBlue:
		return "LightBlue"
	case LightGreen:
		return "LightGreen"
	case Gray:
		return "Gray"
	case DarkGray:
		return "DarkGray"
	case LightRed:
		return "LightRed"
	case LightYellow:
		return "LightYellow"
	case LightMagenta:
		return "LightMagenta"
	case LightCyan:
		return "LightCyan"
	case LightWhite:
		return "LightWhite"
	}
	if index, ok := c.indexedIndex(); ok {
		return strconv.Itoa(int(index))
	}
	if r, g, b, ok := c.rgbComponents(); ok {
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
	return fmt.Sprintf("Color(%d)", c)
}

func ParseColor(value string) (Color, error) {
	if value == "lightwhite" {
		return LightWhite, nil
	}
	normalized := strings.ToLower(value)
	normalized = strings.ReplaceAll(normalized, " ", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")
	normalized = strings.ReplaceAll(normalized, "bright", "light")
	normalized = strings.ReplaceAll(normalized, "grey", "gray")
	normalized = strings.ReplaceAll(normalized, "silver", "gray")
	normalized = strings.ReplaceAll(normalized, "lightblack", "darkgray")
	normalized = strings.ReplaceAll(normalized, "lightwhite", "white")
	normalized = strings.ReplaceAll(normalized, "lightgray", "white")

	switch normalized {
	case "default":
		return Default, nil
	case "reset":
		return Reset, nil
	case "black":
		return Black, nil
	case "red":
		return Red, nil
	case "green":
		return Green, nil
	case "yellow":
		return Yellow, nil
	case "blue":
		return Blue, nil
	case "magenta":
		return Magenta, nil
	case "cyan":
		return Cyan, nil
	case "gray":
		return Gray, nil
	case "darkgray":
		return DarkGray, nil
	case "lightred":
		return LightRed, nil
	case "lightgreen":
		return LightGreen, nil
	case "lightyellow":
		return LightYellow, nil
	case "lightblue":
		return LightBlue, nil
	case "lightmagenta":
		return LightMagenta, nil
	case "lightcyan":
		return LightCyan, nil
	case "white":
		return White, nil
	}

	if index, err := strconv.ParseUint(value, 10, 8); err == nil {
		return IndexedColor(uint8(index)), nil
	}
	if color, ok := parseHexColor(value); ok {
		return color, nil
	}
	return Default, fmt.Errorf("failed to parse color %q", value)
}

func parseHexColor(value string) (Color, bool) {
	if !strings.HasPrefix(value, "#") || len(value) != 7 {
		return Default, false
	}
	r, err := strconv.ParseUint(value[1:3], 16, 8)
	if err != nil {
		return Default, false
	}
	g, err := strconv.ParseUint(value[3:5], 16, 8)
	if err != nil {
		return Default, false
	}
	b, err := strconv.ParseUint(value[5:7], 16, 8)
	if err != nil {
		return Default, false
	}
	return RGBColor(uint8(r), uint8(g), uint8(b)), true
}

func (c Color) indexedIndex() (uint8, bool) {
	if c >= indexedColorBase && c < indexedColorBase+256 {
		return uint8(c - indexedColorBase), true
	}
	return 0, false
}

func (c Color) rgbComponents() (uint8, uint8, uint8, bool) {
	if c < rgbColorBase || c > rgbColorBase+0xFFFFFF {
		return 0, 0, 0, false
	}
	value := uint32(c - rgbColorBase)
	return uint8(value >> 16), uint8(value >> 8), uint8(value), true
}
