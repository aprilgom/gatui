package buffer

import "github.com/rivo/uniseg"

// CellWidth returns the display width of value in terminal cells.
func CellWidth(value string) int {
	if value == "" {
		return 0
	}
	if len(value) == 1 && value[0] >= 0x20 && value[0] < 0x7f {
		return 1
	}

	width := uniseg.StringWidth(value)
	for _, r := range value {
		if isHalfwidthVoicingMark(r) {
			width++
		}
	}
	return width
}
