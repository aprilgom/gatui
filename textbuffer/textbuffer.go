package textbuffer

import (
	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/text"
)

func SetSpan(buf *buffer.Buffer, x, y int, span text.Span, maxWidth int) (endX, endY int) {
	return buf.SetStringN(x, y, span.Content, maxWidth, span.Style)
}

func SetLine(buf *buffer.Buffer, x, y int, line text.Line, maxWidth int) (endX, endY int) {
	if buf == nil || maxWidth <= 0 {
		return x, y
	}

	endX, endY = x, y
	remainingWidth := maxWidth
	for _, span := range line.Spans {
		if remainingWidth <= 0 {
			return endX, endY
		}
		spanStyle := line.LineStyle.Patch(span.Style)
		nextX, nextY := buf.SetStringN(endX, endY, span.Content, remainingWidth, spanStyle)
		remainingWidth -= nextX - endX
		endX, endY = nextX, nextY
	}
	return endX, endY
}
