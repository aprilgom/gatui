package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestCanvas_shouldDrawLabels(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))

	widgets.NewCanvas().
		BackgroundColor(style.Yellow).
		XBounds(0, 5).
		YBounds(0, 5).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Print(0, 0, text.StyledSpan("test", style.NewStyle().Fg(style.Blue)))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"     ",
		"     ",
		"     ",
		"     ",
		"test ",
	})
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if y == 4 && x < 4 {
				continue
			}
			assertCellStyle(t, buf, x, y, style.NewStyle().Bg(style.Yellow))
		}
	}
	for x := 0; x < 4; x++ {
		assertCellStyle(t, buf, x, 4, style.NewStyle().Fg(style.Blue).Bg(style.Yellow))
	}
}

func TestCanvas_shouldSkipLabelsOutsideBounds(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))

	widgets.NewCanvas().
		BackgroundColor(style.Yellow).
		XBounds(0, 5).
		YBounds(0, 5).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Print(-1, 0, text.NewSpan("x"))
			ctx.Print(0, 6, text.NewSpan("y"))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
	})
}

func TestCanvas_shouldTruncateLabelsAtRightEdge(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	widgets.NewCanvas().
		XBounds(0, 5).
		YBounds(0, 1).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Print(3, 0, text.StyledSpan("test", style.NewStyle().Fg(style.Blue)))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"   te"})
	assertCellStyle(t, buf, 3, 0, style.NewStyle().Fg(style.Blue))
	assertCellStyle(t, buf, 4, 0, style.NewStyle().Fg(style.Blue))
}

func TestCanvas_shouldIgnoreZeroSizedArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	assertNotPanics(t, func() {
		widgets.NewCanvas().
			BackgroundColor(style.Yellow).
			XBounds(0, 1).
			YBounds(0, 1).
			Paint(func(ctx *widgets.CanvasContext) {
				ctx.Print(0, 0, text.NewSpan("x"))
			}).
			Render(layout.NewRect(0, 0, 0, 1), buf)
	})

	assertLines(t, buf, []string{" "})
	assertCellStyle(t, buf, 0, 0, style.NewStyle())
}

func TestCanvas_shouldOnlyApplyBackground_whenBoundsAreInvalid(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))

	assertNotPanics(t, func() {
		widgets.NewCanvas().
			BackgroundColor(style.Yellow).
			XBounds(1, 1).
			YBounds(0, 1).
			Paint(func(ctx *widgets.CanvasContext) {
				ctx.Print(1, 0, text.NewSpan("x"))
			}).
			Render(buf.Area, buf)
	})

	assertLines(t, buf, []string{"   "})
	for x := 0; x < 3; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().Bg(style.Yellow))
	}
}
