package widgets_test

import (
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
)

func TestParagraph_ScrollPosition_shouldRenderLikeScroll(t *testing.T) {
	content := text.FromString("one\ntwo\nthree\nfour")
	area := layout.NewRect(0, 0, 4, 2)
	fromScroll := buffer.Empty(area)
	fromPosition := buffer.Empty(area)

	widgets.NewParagraph(content).Scroll(1, 1).Render(area, fromScroll)
	widgets.NewParagraph(content).
		ScrollPosition(widgets.ParagraphScroll{Y: 1, X: 1}).
		Render(area, fromPosition)

	assertLines(t, fromPosition, fromScroll.Lines())
}

func TestParagraph_ScrollOffset_shouldRenderLikeScrollWithWrap(t *testing.T) {
	content := text.FromString("alpha beta gamma delta")
	area := layout.NewRect(0, 0, 8, 2)
	fromScroll := buffer.Empty(area)
	fromOffset := buffer.Empty(area)

	widgets.NewParagraph(content).
		Wrap(widgets.Wrap{Trim: true}).
		Scroll(1, 0).
		Render(area, fromScroll)
	widgets.NewParagraph(content).
		Wrap(widgets.Wrap{Trim: true}).
		ScrollOffset(widgets.ParagraphScroll{Y: 1, X: 0}).
		Render(area, fromOffset)

	assertLines(t, fromOffset, fromScroll.Lines())
}
