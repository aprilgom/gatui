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

const (
	BorderQuadrantTopLeft                       = "▘"
	BorderQuadrantTopRight                      = "▝"
	BorderQuadrantBottomLeft                    = "▖"
	BorderQuadrantBottomRight                   = "▗"
	BorderQuadrantTopHalf                       = HalfBlockUpper
	BorderQuadrantBottomHalf                    = HalfBlockLower
	BorderQuadrantLeftHalf                      = BlockHalf
	BorderQuadrantRightHalf                     = "▐"
	BorderQuadrantTopLeftBottomLeftBottomRight  = "▙"
	BorderQuadrantTopLeftTopRightBottomLeft     = "▛"
	BorderQuadrantTopLeftTopRightBottomRight    = "▜"
	BorderQuadrantTopRightBottomLeftBottomRight = "▟"
	BorderQuadrantTopLeftBottomRight            = "▚"
	BorderQuadrantTopRightBottomLeft            = "▞"
	BorderQuadrantBlock                         = BlockFull
	BorderOneEighthTop                          = "▔"
	BorderOneEighthBottom                       = "▁"
	BorderOneEighthLeft                         = BlockOneEighth
	BorderOneEighthRight                        = "▕"
)

var (
	PlainBorderSet = BorderSet{
		TopLeft:          LineTopLeft,
		TopRight:         LineTopRight,
		BottomLeft:       LineBottomLeft,
		BottomRight:      LineBottomRight,
		VerticalLeft:     LineVertical,
		VerticalRight:    LineVertical,
		HorizontalTop:    LineHorizontal,
		HorizontalBottom: LineHorizontal,
	}
	RoundedBorderSet = BorderSet{
		TopLeft:          LineRoundedTopLeft,
		TopRight:         LineRoundedTopRight,
		BottomLeft:       LineRoundedBottomLeft,
		BottomRight:      LineRoundedBottomRight,
		VerticalLeft:     LineVertical,
		VerticalRight:    LineVertical,
		HorizontalTop:    LineHorizontal,
		HorizontalBottom: LineHorizontal,
	}
	DoubleBorderSet = BorderSet{
		TopLeft:          LineDoubleTopLeft,
		TopRight:         LineDoubleTopRight,
		BottomLeft:       LineDoubleBottomLeft,
		BottomRight:      LineDoubleBottomRight,
		VerticalLeft:     LineDoubleVertical,
		VerticalRight:    LineDoubleVertical,
		HorizontalTop:    LineDoubleHorizontal,
		HorizontalBottom: LineDoubleHorizontal,
	}
	SolidBorderSet = BorderSet{
		TopLeft:          LineThickTopLeft,
		TopRight:         LineThickTopRight,
		BottomLeft:       LineThickBottomLeft,
		BottomRight:      LineThickBottomRight,
		VerticalLeft:     LineThickVertical,
		VerticalRight:    LineThickVertical,
		HorizontalTop:    LineThickHorizontal,
		HorizontalBottom: LineThickHorizontal,
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
		TopLeft:          BorderQuadrantBottomRight,
		TopRight:         BorderQuadrantBottomLeft,
		BottomLeft:       BorderQuadrantTopRight,
		BottomRight:      BorderQuadrantTopLeft,
		VerticalLeft:     BorderQuadrantRightHalf,
		VerticalRight:    BorderQuadrantLeftHalf,
		HorizontalTop:    BorderQuadrantBottomHalf,
		HorizontalBottom: BorderQuadrantTopHalf,
	}
	QuadrantOutsideBorderSet = BorderSet{
		TopLeft:          BorderQuadrantTopLeftTopRightBottomLeft,
		TopRight:         BorderQuadrantTopLeftTopRightBottomRight,
		BottomLeft:       BorderQuadrantTopLeftBottomLeftBottomRight,
		BottomRight:      BorderQuadrantTopRightBottomLeftBottomRight,
		VerticalLeft:     BorderQuadrantLeftHalf,
		VerticalRight:    BorderQuadrantRightHalf,
		HorizontalTop:    BorderQuadrantTopHalf,
		HorizontalBottom: BorderQuadrantBottomHalf,
	}
	OneEighthWideBorderSet = BorderSet{
		TopLeft:          BorderOneEighthBottom,
		TopRight:         BorderOneEighthBottom,
		BottomLeft:       BorderOneEighthTop,
		BottomRight:      BorderOneEighthTop,
		VerticalLeft:     BorderOneEighthLeft,
		VerticalRight:    BorderOneEighthRight,
		HorizontalTop:    BorderOneEighthBottom,
		HorizontalBottom: BorderOneEighthTop,
	}
	OneEighthTallBorderSet = BorderSet{
		TopLeft:          BorderOneEighthRight,
		TopRight:         BorderOneEighthLeft,
		BottomLeft:       BorderOneEighthRight,
		BottomRight:      BorderOneEighthLeft,
		VerticalLeft:     BorderOneEighthRight,
		VerticalRight:    BorderOneEighthLeft,
		HorizontalTop:    BorderOneEighthTop,
		HorizontalBottom: BorderOneEighthBottom,
	}
	ProportionalWideBorderSet = BorderSet{
		TopLeft:          BorderQuadrantBottomHalf,
		TopRight:         BorderQuadrantBottomHalf,
		BottomLeft:       BorderQuadrantTopHalf,
		BottomRight:      BorderQuadrantTopHalf,
		VerticalLeft:     BorderQuadrantBlock,
		VerticalRight:    BorderQuadrantBlock,
		HorizontalTop:    BorderQuadrantBottomHalf,
		HorizontalBottom: BorderQuadrantTopHalf,
	}
	ProportionalTallBorderSet = BorderSet{
		TopLeft:          BorderQuadrantBlock,
		TopRight:         BorderQuadrantBlock,
		BottomLeft:       BorderQuadrantBlock,
		BottomRight:      BorderQuadrantBlock,
		VerticalLeft:     BorderQuadrantBlock,
		VerticalRight:    BorderQuadrantBlock,
		HorizontalTop:    BorderQuadrantTopHalf,
		HorizontalBottom: BorderQuadrantBottomHalf,
	}
	FullBorderSet = BorderSet{
		TopLeft:          BlockFull,
		TopRight:         BlockFull,
		BottomLeft:       BlockFull,
		BottomRight:      BlockFull,
		VerticalLeft:     BlockFull,
		VerticalRight:    BlockFull,
		HorizontalTop:    BlockFull,
		HorizontalBottom: BlockFull,
	}
	EmptyBorderSet = BorderSet{
		TopLeft:          BlockEmpty,
		TopRight:         BlockEmpty,
		BottomLeft:       BlockEmpty,
		BottomRight:      BlockEmpty,
		VerticalLeft:     BlockEmpty,
		VerticalRight:    BlockEmpty,
		HorizontalTop:    BlockEmpty,
		HorizontalBottom: BlockEmpty,
	}
)
