package text

import (
	"unicode"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"

	"github.com/rivo/uniseg"
)

func renderLineSpans(spans []Span, area layout.Rect, buf *buffer.Buffer, skipWidth int) {
	x := area.X
	for _, span := range spans {
		spanWidth := span.Width()
		if skipWidth >= spanWidth {
			skipWidth -= spanWidth
			continue
		}
		if skipWidth > 0 {
			x = renderSpan(span, layout.NewRect(x, area.Y, area.Right()-x, 1), buf, skipWidth)
			skipWidth = 0
		} else {
			x = renderSpan(span, layout.NewRect(x, area.Y, area.Right()-x, 1), buf, 0)
		}
		if x >= area.Right() {
			return
		}
	}
}

func renderSpan(span Span, area layout.Rect, buf *buffer.Buffer, skipWidth int) int {
	if buf == nil {
		return area.X
	}
	area = area.Intersection(buf.Area)
	if area.Width == 0 || area.Height == 0 {
		return area.X
	}

	x := area.X
	right := area.Right()
	renderedAny := false
	graphemes := uniseg.NewGraphemes(span.Content)
	for graphemes.Next() {
		symbol := graphemes.Str()
		if symbol == "\n" {
			continue
		}
		width := buffer.CellWidth(symbol)
		if width == 0 {
			if !renderedAny {
				setSpanCellSymbol(buf, x, area.Y, symbol, span.Style, false)
				renderedAny = true
			} else if x == area.X {
				setSpanCellSymbol(buf, x, area.Y, symbol, span.Style, true)
			} else {
				setSpanCellSymbol(buf, x-1, area.Y, symbol, span.Style, true)
			}
			continue
		}
		if skipWidth >= width {
			skipWidth -= width
			continue
		}
		if skipWidth > 0 {
			x += width - skipWidth
			skipWidth = 0
			continue
		}
		if x+width > right {
			break
		}

		setSpanCellSymbol(buf, x, area.Y, symbol, span.Style, renderedAny && x == area.X)
		for hidden := 1; hidden < width; hidden++ {
			buf.SetCell(x+hidden, area.Y, buffer.Cell{Symbol: " ", Style: style.NewStyle()})
		}
		x += width
		renderedAny = true
	}
	return x
}

func setSpanCellSymbol(buf *buffer.Buffer, x, y int, symbol string, spanStyle style.Style, appendSymbol bool) {
	cellStyle := style.NewStyle()
	cellSymbol := symbol
	if cell, ok := buf.CellAt(x, y); ok {
		cellStyle = cell.Style
		if appendSymbol {
			cellSymbol = cell.Symbol + symbol
		}
	}
	buf.SetCell(x, y, buffer.Cell{Symbol: cellSymbol, Style: cellStyle.Patch(spanStyle)})
}

func containsControl(symbol string) bool {
	for _, r := range symbol {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func alignedRenderOffset(lineWidth, areaWidth int, alignment *layout.Alignment) int {
	if alignment == nil || lineWidth >= areaWidth {
		return 0
	}
	switch *alignment {
	case layout.Center:
		return (areaWidth - lineWidth) / 2
	case layout.Right:
		return areaWidth - lineWidth
	default:
		return 0
	}
}
