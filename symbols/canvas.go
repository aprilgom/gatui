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

const (
	CanvasDotSymbol   = "•"
	CanvasBlockSymbol = BlockFull
	CanvasBarSymbol   = HalfBlockLower
)

func (m CanvasMarker) String() string {
	switch m.Kind() {
	case CanvasMarkerDot:
		return "Dot"
	case CanvasMarkerBlock:
		return "Block"
	case CanvasMarkerBar:
		return "Bar"
	case CanvasMarkerBraille:
		return "Braille"
	case CanvasMarkerHalfBlock:
		return "HalfBlock"
	case CanvasMarkerQuadrant:
		return "Quadrant"
	case CanvasMarkerSextant:
		return "Sextant"
	case CanvasMarkerOctant:
		return "Octant"
	case "custom":
		return "Custom"
	default:
		return string(m)
	}
}

func ParseCanvasMarker(s string) (CanvasMarker, bool) {
	switch s {
	case "Dot":
		return CanvasMarkerDot, true
	case "Block":
		return CanvasMarkerBlock, true
	case "Bar":
		return CanvasMarkerBar, true
	case "Braille":
		return CanvasMarkerBraille, true
	case "HalfBlock":
		return CanvasMarkerHalfBlock, true
	case "Quadrant":
		return CanvasMarkerQuadrant, true
	case "Sextant":
		return CanvasMarkerSextant, true
	case "Octant":
		return CanvasMarkerOctant, true
	default:
		return "", false
	}
}

func IsCanvasDatasetSymbol(symbol string) bool {
	if symbol == CanvasDotSymbol || symbol == CanvasBlockSymbol || symbol == CanvasBarSymbol {
		return true
	}
	runes := []rune(symbol)
	if len(runes) != 1 {
		return false
	}
	return runes[0] >= 0x2800 && runes[0] <= 0x28ff
}

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
	return Braille[pattern]
}

var Braille = makeBraille()

func makeBraille() [256]rune {
	masks := [8]rune{0x01, 0x08, 0x02, 0x10, 0x04, 0x20, 0x40, 0x80}
	var braille [256]rune
	for pattern := range braille {
		code := rune(0x2800)
		for i, mask := range masks {
			if pattern&(1<<uint(i)) != 0 {
				code += mask
			}
		}
		braille[pattern] = code
	}
	return braille
}

func QuadrantSymbol(pattern uint8) rune {
	return Quadrants[pattern&0x0f]
}

var Quadrants = [16]rune{' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛', '▗', '▚', '▐', '▜', '▄', '▙', '▟', '█'}

func SextantSymbol(pattern uint8) string {
	return Sextants[pattern&0x3f]
}

var Sextants = [64]string{
	" ", "🬀", "🬁", "🬂", "🬃", "🬄", "🬅", "🬆",
	"🬇", "🬈", "🬉", "🬊", "🬋", "🬌", "🬍", "🬎",
	"🬏", "🬐", "🬑", "🬒", "🬓", "▌", "🬔", "🬕",
	"🬖", "🬗", "🬘", "🬙", "🬚", "🬛", "🬜", "🬝",
	"🬞", "🬟", "🬠", "🬡", "🬢", "🬣", "🬤", "🬥",
	"🬦", "🬧", "▐", "🬨", "🬩", "🬪", "🬫", "🬬",
	"🬭", "🬮", "🬯", "🬰", "🬱", "🬲", "🬳", "🬴",
	"🬵", "🬶", "🬷", "🬸", "🬹", "🬺", "🬻", "█",
}

func OctantSymbol(pattern uint8) string {
	return Octants[pattern]
}

var Octants = [256]string{
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
