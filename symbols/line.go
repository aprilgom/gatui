package symbols

const (
	LineVertical                   = "│"
	LineDoubleVertical             = "║"
	LineThickVertical              = "┃"
	LineLightDoubleDashVertical    = "╎"
	LineHeavyDoubleDashVertical    = "╏"
	LineLightTripleDashVertical    = "┆"
	LineHeavyTripleDashVertical    = "┇"
	LineLightQuadrupleDashVertical = "┊"
	LineHeavyQuadrupleDashVertical = "┋"

	LineHorizontal                   = "─"
	LineDoubleHorizontal             = "═"
	LineThickHorizontal              = "━"
	LineLightDoubleDashHorizontal    = "╌"
	LineHeavyDoubleDashHorizontal    = "╍"
	LineLightTripleDashHorizontal    = "┄"
	LineHeavyTripleDashHorizontal    = "┅"
	LineLightQuadrupleDashHorizontal = "┈"
	LineHeavyQuadrupleDashHorizontal = "┉"

	LineTopRight        = "┐"
	LineRoundedTopRight = "╮"
	LineDoubleTopRight  = "╗"
	LineThickTopRight   = "┓"

	LineTopLeft        = "┌"
	LineRoundedTopLeft = "╭"
	LineDoubleTopLeft  = "╔"
	LineThickTopLeft   = "┏"

	LineBottomRight        = "┘"
	LineRoundedBottomRight = "╯"
	LineDoubleBottomRight  = "╝"
	LineThickBottomRight   = "┛"

	LineBottomLeft        = "└"
	LineRoundedBottomLeft = "╰"
	LineDoubleBottomLeft  = "╚"
	LineThickBottomLeft   = "┗"

	LineVerticalLeft       = "┤"
	LineDoubleVerticalLeft = "╣"
	LineThickVerticalLeft  = "┫"

	LineVerticalRight       = "├"
	LineDoubleVerticalRight = "╠"
	LineThickVerticalRight  = "┣"

	LineHorizontalDown       = "┬"
	LineDoubleHorizontalDown = "╦"
	LineThickHorizontalDown  = "┳"

	LineHorizontalUp       = "┴"
	LineDoubleHorizontalUp = "╩"
	LineThickHorizontalUp  = "┻"

	LineCross       = "┼"
	LineDoubleCross = "╬"
	LineThickCross  = "╋"
)

type LineSet struct {
	Vertical       string
	Horizontal     string
	TopRight       string
	TopLeft        string
	BottomRight    string
	BottomLeft     string
	VerticalLeft   string
	VerticalRight  string
	HorizontalDown string
	HorizontalUp   string
	Cross          string
}

var (
	NormalLineSet = LineSet{
		Vertical:       LineVertical,
		Horizontal:     LineHorizontal,
		TopRight:       LineTopRight,
		TopLeft:        LineTopLeft,
		BottomRight:    LineBottomRight,
		BottomLeft:     LineBottomLeft,
		VerticalLeft:   LineVerticalLeft,
		VerticalRight:  LineVerticalRight,
		HorizontalDown: LineHorizontalDown,
		HorizontalUp:   LineHorizontalUp,
		Cross:          LineCross,
	}
	RoundedLineSet = LineSet{
		Vertical:       LineVertical,
		Horizontal:     LineHorizontal,
		TopRight:       LineRoundedTopRight,
		TopLeft:        LineRoundedTopLeft,
		BottomRight:    LineRoundedBottomRight,
		BottomLeft:     LineRoundedBottomLeft,
		VerticalLeft:   LineVerticalLeft,
		VerticalRight:  LineVerticalRight,
		HorizontalDown: LineHorizontalDown,
		HorizontalUp:   LineHorizontalUp,
		Cross:          LineCross,
	}
	DoubleLineSet = LineSet{
		Vertical:       LineDoubleVertical,
		Horizontal:     LineDoubleHorizontal,
		TopRight:       LineDoubleTopRight,
		TopLeft:        LineDoubleTopLeft,
		BottomRight:    LineDoubleBottomRight,
		BottomLeft:     LineDoubleBottomLeft,
		VerticalLeft:   LineDoubleVerticalLeft,
		VerticalRight:  LineDoubleVerticalRight,
		HorizontalDown: LineDoubleHorizontalDown,
		HorizontalUp:   LineDoubleHorizontalUp,
		Cross:          LineDoubleCross,
	}
	ThickLineSet = LineSet{
		Vertical:       LineThickVertical,
		Horizontal:     LineThickHorizontal,
		TopRight:       LineThickTopRight,
		TopLeft:        LineThickTopLeft,
		BottomRight:    LineThickBottomRight,
		BottomLeft:     LineThickBottomLeft,
		VerticalLeft:   LineThickVerticalLeft,
		VerticalRight:  LineThickVerticalRight,
		HorizontalDown: LineThickHorizontalDown,
		HorizontalUp:   LineThickHorizontalUp,
		Cross:          LineThickCross,
	}
	LightDoubleDashedLineSet = LineSet{
		Vertical:       LineLightDoubleDashVertical,
		Horizontal:     LineLightDoubleDashHorizontal,
		TopRight:       LineTopRight,
		TopLeft:        LineTopLeft,
		BottomRight:    LineBottomRight,
		BottomLeft:     LineBottomLeft,
		VerticalLeft:   LineVerticalLeft,
		VerticalRight:  LineVerticalRight,
		HorizontalDown: LineHorizontalDown,
		HorizontalUp:   LineHorizontalUp,
		Cross:          LineCross,
	}
	HeavyDoubleDashedLineSet = LineSet{
		Vertical:       LineHeavyDoubleDashVertical,
		Horizontal:     LineHeavyDoubleDashHorizontal,
		TopRight:       LineThickTopRight,
		TopLeft:        LineThickTopLeft,
		BottomRight:    LineThickBottomRight,
		BottomLeft:     LineThickBottomLeft,
		VerticalLeft:   LineThickVerticalLeft,
		VerticalRight:  LineThickVerticalRight,
		HorizontalDown: LineThickHorizontalDown,
		HorizontalUp:   LineThickHorizontalUp,
		Cross:          LineThickCross,
	}
	LightTripleDashedLineSet = LineSet{
		Vertical:       LineLightTripleDashVertical,
		Horizontal:     LineLightTripleDashHorizontal,
		TopRight:       LineTopRight,
		TopLeft:        LineTopLeft,
		BottomRight:    LineBottomRight,
		BottomLeft:     LineBottomLeft,
		VerticalLeft:   LineVerticalLeft,
		VerticalRight:  LineVerticalRight,
		HorizontalDown: LineHorizontalDown,
		HorizontalUp:   LineHorizontalUp,
		Cross:          LineCross,
	}
	HeavyTripleDashedLineSet = LineSet{
		Vertical:       LineHeavyTripleDashVertical,
		Horizontal:     LineHeavyTripleDashHorizontal,
		TopRight:       LineThickTopRight,
		TopLeft:        LineThickTopLeft,
		BottomRight:    LineThickBottomRight,
		BottomLeft:     LineThickBottomLeft,
		VerticalLeft:   LineThickVerticalLeft,
		VerticalRight:  LineThickVerticalRight,
		HorizontalDown: LineThickHorizontalDown,
		HorizontalUp:   LineThickHorizontalUp,
		Cross:          LineThickCross,
	}
	LightQuadrupleDashedLineSet = LineSet{
		Vertical:       LineLightQuadrupleDashVertical,
		Horizontal:     LineLightQuadrupleDashHorizontal,
		TopRight:       LineTopRight,
		TopLeft:        LineTopLeft,
		BottomRight:    LineBottomRight,
		BottomLeft:     LineBottomLeft,
		VerticalLeft:   LineVerticalLeft,
		VerticalRight:  LineVerticalRight,
		HorizontalDown: LineHorizontalDown,
		HorizontalUp:   LineHorizontalUp,
		Cross:          LineCross,
	}
	HeavyQuadrupleDashedLineSet = LineSet{
		Vertical:       LineHeavyQuadrupleDashVertical,
		Horizontal:     LineHeavyQuadrupleDashHorizontal,
		TopRight:       LineThickTopRight,
		TopLeft:        LineThickTopLeft,
		BottomRight:    LineThickBottomRight,
		BottomLeft:     LineThickBottomLeft,
		VerticalLeft:   LineThickVerticalLeft,
		VerticalRight:  LineThickVerticalRight,
		HorizontalDown: LineThickHorizontalDown,
		HorizontalUp:   LineThickHorizontalUp,
		Cross:          LineThickCross,
	}
)
