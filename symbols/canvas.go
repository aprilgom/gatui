package symbols

import "strings"

type CanvasMarker string

const (
	CanvasMarkerDot       CanvasMarker = "dot"
	CanvasMarkerBlock     CanvasMarker = "block"
	CanvasMarkerBar       CanvasMarker = "bar"
	CanvasMarkerBraille   CanvasMarker = "braille"
	CanvasMarkerHalfBlock CanvasMarker = "half_block"
	CanvasMarkerQuadrant  CanvasMarker = "quadrant"
	CanvasMarkerSextant   CanvasMarker = "sextant"
	CanvasMarkerOctant    CanvasMarker = "octant"
)

func CanvasMarkerCustom(symbol string) CanvasMarker {
	runes := []rune(symbol)
	if len(runes) == 0 {
		return CanvasMarker("custom: ")
	}
	return CanvasMarker("custom:" + string(runes[0]))
}

func (m CanvasMarker) Kind() CanvasMarker {
	if strings.HasPrefix(string(m), "custom:") {
		return "custom"
	}
	return m
}

func (m CanvasMarker) CustomSymbol() string {
	symbol, ok := strings.CutPrefix(string(m), "custom:")
	if !ok || symbol == "" {
		return " "
	}
	return symbol
}

func (m CanvasMarker) CellResolution() (int, int) {
	switch m.Kind() {
	case CanvasMarkerBraille, CanvasMarkerOctant:
		return 2, 4
	case CanvasMarkerHalfBlock:
		return 1, 2
	case CanvasMarkerQuadrant:
		return 2, 2
	case CanvasMarkerSextant:
		return 2, 3
	default:
		return 1, 1
	}
}

func BrailleSymbol(pattern uint8) rune {
	masks := [8]rune{0x01, 0x08, 0x02, 0x10, 0x04, 0x20, 0x40, 0x80}
	code := rune(0x2800)
	for i, mask := range masks {
		if pattern&(1<<uint(i)) != 0 {
			code += mask
		}
	}
	return code
}

func QuadrantSymbol(pattern uint8) rune {
	return quadrantSymbols[pattern&0x0f]
}

var quadrantSymbols = [16]rune{' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛', '▗', '▚', '▐', '▜', '▄', '▙', '▟', '█'}

func SextantSymbol(pattern uint8) string {
	symbols := [64]string{
		" ", "🬀", "🬁", "🬂", "🬃", "🬄", "🬅", "🬆",
		"🬇", "🬈", "🬉", "🬊", "🬋", "🬌", "🬍", "🬎",
		"🬏", "🬐", "🬑", "🬒", "🬓", "▌", "🬔", "🬕",
		"🬖", "🬗", "🬘", "🬙", "🬚", "🬛", "🬜", "🬝",
		"🬞", "🬟", "🬠", "🬡", "🬢", "🬣", "🬤", "🬥",
		"🬦", "🬧", "▐", "🬨", "🬩", "🬪", "🬫", "🬬",
		"🬭", "🬮", "🬯", "🬰", "🬱", "🬲", "🬳", "🬴",
		"🬵", "🬶", "🬷", "🬸", "🬹", "🬺", "🬻", "█",
	}
	return symbols[pattern&0x3f]
}

func OctantSymbol(pattern uint8) string {
	symbols := [256]string{
		" ", "𜺨", "𜺫", "🮂", "𜴀", "▘", "𜴁", "𜴂",
		"𜴃", "𜴄", "▝", "𜴅", "𜴆", "𜴇", "𜴈", "▀",
		"𜴉", "𜴊", "𜴋", "𜴌", "🯦", "𜴍", "𜴎", "𜴏",
		"𜴐", "𜴑", "𜴒", "𜴓", "𜴔", "𜴕", "𜴖", "𜴗",
		"𜴘", "𜴙", "𜴚", "𜴛", "𜴜", "𜴝", "𜴞", "𜴟",
		"🯧", "𜴠", "𜴡", "𜴢", "𜴣", "𜴤", "𜴥", "𜴦",
		"𜴧", "𜴨", "𜴩", "𜴪", "𜴫", "𜴬", "𜴭", "𜴮",
		"𜴯", "𜴰", "𜴱", "𜴲", "𜴳", "𜴴", "𜴵", "🮅",
		"𜺣", "𜴶", "𜴷", "𜴸", "𜴹", "𜴺", "𜴻", "𜴼",
		"𜴽", "𜴾", "𜴿", "𜵀", "𜵁", "𜵂", "𜵃", "𜵄",
		"▖", "𜵅", "𜵆", "𜵇", "𜵈", "▌", "𜵉", "𜵊",
		"𜵋", "𜵌", "▞", "𜵍", "𜵎", "𜵏", "𜵐", "▛",
		"𜵑", "𜵒", "𜵓", "𜵔", "𜵕", "𜵖", "𜵗", "𜵘",
		"𜵙", "𜵚", "𜵛", "𜵜", "𜵝", "𜵞", "𜵟", "𜵠",
		"𜵡", "𜵢", "𜵣", "𜵤", "𜵥", "𜵦", "𜵧", "𜵨",
		"𜵩", "𜵪", "𜵫", "𜵬", "𜵭", "𜵮", "𜵯", "𜵰",
		"𜺠", "𜵱", "𜵲", "𜵳", "𜵴", "𜵵", "𜵶", "𜵷",
		"𜵸", "𜵹", "𜵺", "𜵻", "𜵼", "𜵽", "𜵾", "𜵿",
		"𜶀", "𜶁", "𜶂", "𜶃", "𜶄", "𜶅", "𜶆", "𜶇",
		"𜶈", "𜶉", "𜶊", "𜶋", "𜶌", "𜶍", "𜶎", "𜶏",
		"▗", "𜶐", "𜶑", "𜶒", "𜶓", "▚", "𜶔", "𜶕",
		"𜶖", "𜶗", "▐", "𜶘", "𜶙", "𜶚", "𜶛", "▜",
		"𜶜", "𜶝", "𜶞", "𜶟", "𜶠", "𜶡", "𜶢", "𜶣",
		"𜶤", "𜶥", "𜶦", "𜶧", "𜶨", "𜶩", "𜶪", "𜶫",
		"▂", "𜶬", "𜶭", "𜶮", "𜶯", "𜶰", "𜶱", "𜶲",
		"𜶳", "𜶴", "𜶵", "𜶶", "𜶷", "𜶸", "𜶹", "𜶺",
		"𜶻", "𜶼", "𜶽", "𜶾", "𜶿", "𜷀", "𜷁", "𜷂",
		"𜷃", "𜷄", "𜷅", "𜷆", "𜷇", "𜷈", "𜷉", "𜷊",
		"𜷋", "𜷌", "𜷍", "𜷎", "𜷏", "𜷐", "𜷑", "𜷒",
		"𜷓", "𜷔", "𜷕", "𜷖", "𜷗", "𜷘", "𜷙", "𜷚",
		"▄", "𜷛", "𜷜", "𜷝", "𜷞", "▙", "𜷟", "𜷠",
		"𜷡", "𜷢", "▟", "𜷣", "▆", "𜷤", "𜷥", "█",
	}
	return symbols[pattern]
}
