package symbols

type BorderSet struct {
	TopLeft          string
	TopRight         string
	BottomLeft       string
	BottomRight      string
	VerticalLeft     string
	VerticalRight    string
	HorizontalTop    string
	HorizontalBottom string
}

var (
	PlainBorderSet = BorderSet{
		TopLeft:          "┌",
		TopRight:         "┐",
		BottomLeft:       "└",
		BottomRight:      "┘",
		VerticalLeft:     "│",
		VerticalRight:    "│",
		HorizontalTop:    "─",
		HorizontalBottom: "─",
	}
	RoundedBorderSet = BorderSet{
		TopLeft:          "╭",
		TopRight:         "╮",
		BottomLeft:       "╰",
		BottomRight:      "╯",
		VerticalLeft:     "│",
		VerticalRight:    "│",
		HorizontalTop:    "─",
		HorizontalBottom: "─",
	}
	DoubleBorderSet = BorderSet{
		TopLeft:          "╔",
		TopRight:         "╗",
		BottomLeft:       "╚",
		BottomRight:      "╝",
		VerticalLeft:     "║",
		VerticalRight:    "║",
		HorizontalTop:    "═",
		HorizontalBottom: "═",
	}
	SolidBorderSet = BorderSet{
		TopLeft:          "┏",
		TopRight:         "┓",
		BottomLeft:       "┗",
		BottomRight:      "┛",
		VerticalLeft:     "┃",
		VerticalRight:    "┃",
		HorizontalTop:    "━",
		HorizontalBottom: "━",
	}
	LightDoubleDashedBorderSet = BorderSet{
		TopLeft:          "┌",
		TopRight:         "┐",
		BottomLeft:       "└",
		BottomRight:      "┘",
		VerticalLeft:     "╎",
		VerticalRight:    "╎",
		HorizontalTop:    "╌",
		HorizontalBottom: "╌",
	}
	HeavyDoubleDashedBorderSet = BorderSet{
		TopLeft:          "┏",
		TopRight:         "┓",
		BottomLeft:       "┗",
		BottomRight:      "┛",
		VerticalLeft:     "╏",
		VerticalRight:    "╏",
		HorizontalTop:    "╍",
		HorizontalBottom: "╍",
	}
	LightTripleDashedBorderSet = BorderSet{
		TopLeft:          "┌",
		TopRight:         "┐",
		BottomLeft:       "└",
		BottomRight:      "┘",
		VerticalLeft:     "┆",
		VerticalRight:    "┆",
		HorizontalTop:    "┄",
		HorizontalBottom: "┄",
	}
	HeavyTripleDashedBorderSet = BorderSet{
		TopLeft:          "┏",
		TopRight:         "┓",
		BottomLeft:       "┗",
		BottomRight:      "┛",
		VerticalLeft:     "┇",
		VerticalRight:    "┇",
		HorizontalTop:    "┅",
		HorizontalBottom: "┅",
	}
	LightQuadrupleDashedBorderSet = BorderSet{
		TopLeft:          "┌",
		TopRight:         "┐",
		BottomLeft:       "└",
		BottomRight:      "┘",
		VerticalLeft:     "┊",
		VerticalRight:    "┊",
		HorizontalTop:    "┈",
		HorizontalBottom: "┈",
	}
	HeavyQuadrupleDashedBorderSet = BorderSet{
		TopLeft:          "┏",
		TopRight:         "┓",
		BottomLeft:       "┗",
		BottomRight:      "┛",
		VerticalLeft:     "┋",
		VerticalRight:    "┋",
		HorizontalTop:    "┉",
		HorizontalBottom: "┉",
	}
	QuadrantInsideBorderSet = BorderSet{
		TopLeft:          "▗",
		TopRight:         "▖",
		BottomLeft:       "▝",
		BottomRight:      "▘",
		VerticalLeft:     "▐",
		VerticalRight:    "▌",
		HorizontalTop:    "▄",
		HorizontalBottom: "▀",
	}
	QuadrantOutsideBorderSet = BorderSet{
		TopLeft:          "▛",
		TopRight:         "▜",
		BottomLeft:       "▙",
		BottomRight:      "▟",
		VerticalLeft:     "▌",
		VerticalRight:    "▐",
		HorizontalTop:    "▀",
		HorizontalBottom: "▄",
	}
)
