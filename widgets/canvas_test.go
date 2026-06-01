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

func TestCanvas_shouldDrawPointsInsideBoundsWithDotMarker(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))

	widgets.NewCanvas().
		XBounds(0, 4).
		YBounds(0, 4).
		Marker(widgets.CanvasMarkerDot).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{
				{X: 0, Y: 0},
				{X: 2, Y: 2},
				{X: 4, Y: 4},
			}, style.Red))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"    •",
		"     ",
		"  •  ",
		"     ",
		"•    ",
	})
	assertCellStyle(t, buf, 0, 4, style.NewStyle().Fg(style.Red))
	assertCellStyle(t, buf, 2, 2, style.NewStyle().Fg(style.Red))
	assertCellStyle(t, buf, 4, 0, style.NewStyle().Fg(style.Red))
}

func TestCanvas_shouldSkipPointsOutsideBounds(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))

	widgets.NewCanvas().
		XBounds(0, 2).
		YBounds(0, 2).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{
				{X: -1, Y: 1},
				{X: 1, Y: 3},
			}, style.Red))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"   ",
		"   ",
		"   ",
	})
}

func TestCanvas_shouldPreserveBackgroundForDotMarker(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	widgets.NewCanvas().
		BackgroundColor(style.Yellow).
		Marker(widgets.CanvasMarkerDot).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{{X: 0, Y: 0}}, style.Blue))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"•"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Blue).Bg(style.Yellow))
}

func TestCanvas_shouldApplyForegroundAndBackgroundForBlockMarker(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	widgets.NewCanvas().
		BackgroundColor(style.Yellow).
		Marker(widgets.CanvasMarkerBlock).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{{X: 0, Y: 0}}, style.Blue))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"█"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Blue).Bg(style.Blue))
}

func TestCanvas_shouldDrawHorizontalVerticalDiagonalAndClippedLines(t *testing.T) {
	tests := []struct {
		name     string
		line     widgets.CanvasLine
		expected []string
	}{
		{
			name: "horizontal",
			line: widgets.NewCanvasLine(0, 2, 4, 2, style.Red),
			expected: []string{
				"     ",
				"     ",
				"•••••",
				"     ",
				"     ",
			},
		},
		{
			name: "vertical",
			line: widgets.NewCanvasLine(2, 0, 2, 4, style.Red),
			expected: []string{
				"  •  ",
				"  •  ",
				"  •  ",
				"  •  ",
				"  •  ",
			},
		},
		{
			name: "diagonal",
			line: widgets.NewCanvasLine(0, 0, 4, 4, style.Red),
			expected: []string{
				"    •",
				"   • ",
				"  •  ",
				" •   ",
				"•    ",
			},
		},
		{
			name: "clipped",
			line: widgets.NewCanvasLine(-2, 2, 2, 2, style.Red),
			expected: []string{
				"     ",
				"     ",
				"•••  ",
				"     ",
				"     ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))

			widgets.NewCanvas().
				XBounds(0, 4).
				YBounds(0, 4).
				Paint(func(ctx *widgets.CanvasContext) {
					ctx.Draw(tt.line)
				}).
				Render(buf.Area, buf)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestCanvas_shouldSkipOffGridLinesWithoutPanicking(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))

	assertNotPanics(t, func() {
		widgets.NewCanvas().
			XBounds(0, 2).
			YBounds(0, 2).
			Paint(func(ctx *widgets.CanvasContext) {
				ctx.Draw(widgets.NewCanvasLine(-3, -3, -1, -1, style.Red))
			}).
			Render(buf.Area, buf)
	})

	assertLines(t, buf, []string{
		"   ",
		"   ",
		"   ",
	})
}

func TestCanvas_shouldDrawRectangleEdgesForDotAndBlockMarkers(t *testing.T) {
	tests := []struct {
		name     string
		marker   widgets.CanvasMarker
		expected []string
	}{
		{
			name:   "dot",
			marker: widgets.CanvasMarkerDot,
			expected: []string{
				"•••••",
				"•   •",
				"•   •",
				"•   •",
				"•••••",
			},
		},
		{
			name:   "block",
			marker: widgets.CanvasMarkerBlock,
			expected: []string{
				"█████",
				"█   █",
				"█   █",
				"█   █",
				"█████",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))

			widgets.NewCanvas().
				XBounds(0, 4).
				YBounds(0, 4).
				Marker(tt.marker).
				Paint(func(ctx *widgets.CanvasContext) {
					ctx.Draw(widgets.NewRectangle(0, 0, 4, 4, style.Green))
				}).
				Render(buf.Area, buf)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestCanvas_shouldRenderLabelsAfterShapes(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))

	widgets.NewCanvas().
		XBounds(0, 4).
		YBounds(0, 1).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewCanvasLine(0, 0, 4, 0, style.Red))
			ctx.Print(2, 0, text.StyledSpan("xy", style.NewStyle().Fg(style.Blue)))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"••xy•"})
	assertCellStyle(t, buf, 2, 0, style.NewStyle().Fg(style.Blue))
	assertCellStyle(t, buf, 3, 0, style.NewStyle().Fg(style.Blue))
}
