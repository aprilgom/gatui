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
	for y := range 5 {
		for x := range 5 {
			if y == 4 && x < 4 {
				continue
			}
			assertCellStyle(t, buf, x, y, style.NewStyle().Bg(style.Yellow))
		}
	}
	for x := range 4 {
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
	for x := range 3 {
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

func TestCanvas_shouldExposeCanvasMarkerParityAPI(t *testing.T) {
	_ = []widgets.CanvasMarker{
		widgets.CanvasMarkerDot,
		widgets.CanvasMarkerBlock,
		widgets.CanvasMarkerBar,
		widgets.CanvasMarkerBraille,
		widgets.CanvasMarkerHalfBlock,
		widgets.CanvasMarkerQuadrant,
		widgets.CanvasMarkerSextant,
		widgets.CanvasMarkerOctant,
		widgets.CanvasMarkerCustom("x"),
	}
}

func TestCanvasMapResolution_shouldString(t *testing.T) {
	tests := []struct {
		resolution widgets.MapResolution
		expected   string
	}{
		{resolution: widgets.MapResolutionLow, expected: "Low"},
		{resolution: widgets.MapResolutionHigh, expected: "High"},
		{resolution: widgets.MapResolution(99), expected: "MapResolution(99)"},
	}

	for _, tt := range tests {
		if actual := tt.resolution.String(); actual != tt.expected {
			t.Fatalf("String() = %q, want %q", actual, tt.expected)
		}
	}
}

func TestCanvasMap_shouldExposeDefaultLowResolution(t *testing.T) {
	mapShape := widgets.NewMap()

	if mapShape.Resolution != widgets.MapResolutionLow {
		t.Fatalf("NewMap().Resolution = %v, want %v", mapShape.Resolution, widgets.MapResolutionLow)
	}
	if mapShape.Color != style.Reset {
		t.Fatalf("NewMap().Color = %v, want %v", mapShape.Color, style.Reset)
	}
}

func TestCanvas_shouldDrawLowResolutionMap(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 80, 40))

	widgets.NewCanvas().
		Marker(widgets.CanvasMarkerDot).
		XBounds(-180, 180).
		YBounds(-90, 90).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewMap())
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"                                                                                ",
		"                               •                                                ",
		"               • •• •••••••• ••   ••••    •••••  ••• ••     •••                 ",
		"             •••••••••••••••       •      ••••      • •   •••••••     •••       ",
		"    • •••• ••••••••••••••• ••     ••  •     •••    ••  ••••    ••  ••••••• •••  ",
		"•••••     •••••••••••• •••• •  ••••••     •••• • ••• •••••                     •",
		"   ••  • •   •••• ••••••••  ••••   ••  • •• •  •••                        •• •••",
		"    •••• •••   •••••• •••••   •       •• ••••••                       • •••••   ",
		"•••••     •••     •  ••   ••         •••••••                          ••  •• •• ",
		"            ••    ••••  •••••          ••       •  • •                ••        ",
		"            •  •    •••••••           •• •••• ••• •• •  ••          • ••        ",
		"            •          ••             ••••••••• • ••             •••• •         ",
		"             ••       ••              • • • •• •                  •••••         ",
		"              ••   •••               •      ••••  •               • •           ",
		"               •  •   ••             •         ••  •• •           •             ",
		"    ••          • •••••••           •           •   •  •   •   •• •             ",
		"                 •••••••••          •           •• •   •  • •• •  ••            ",
		"                    ••  ••          •            •••     •   •••  ••            ",
		"                     •••  • •        •  •         •     ••  •••  •••            ",
		"                      •               •  ••                   • ••              ",
		"                   •  •     •••                • •            •••   •••         ",
		"                                •         •     •              • •    •••       ",
		"  •                                        •    • •                  • • •      ",
		"                       •       •                • •               ••• ••       •",
		"                        •      •          •    • ••              •      •   •   ",
		"                        •    •                   •               •       •      ",
		"                        •   •              •   •                    •           ",
		"                           ••               ••                   ••  ••  •   •  ",
		"                       •  •                                           •••    •• ",
		"                       •  •                                            ••   ••  ",
		"                       • •                                                      ",
		"                       •••••                                                    ",
		"                                                                                ",
		"                          ••                                                    ",
		"                         •••           •       • ••••• • •••• • • •• •• ••      ",
		"            •    • • ••••••        ••••••••• • ••      ••                  •••  ",
		"•    ••• •••• ••••   • •  • ••• • •                                        ••• •",
		"   •• •                •  ••  • ••                                         ••   ",
		"•      •                                                                      • ",
		"                                                                                ",
	})
	assertCellStyle(t, buf, 31, 1, style.NewStyle().Fg(style.Reset))
}

func TestCanvas_shouldDrawHighResolutionMapWithBrailleMarker(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 80, 40))

	widgets.NewCanvas().
		Marker(widgets.CanvasMarkerBraille).
		XBounds(-180, 180).
		YBounds(-90, 90).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.Map{Resolution: widgets.MapResolutionHigh, Color: style.Reset})
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"                                                                                ",
		"                   ⢀⣀⣤⠄⠤⠤⣤⣀⡀⣀⣀⡄⠄⢄⣀⣄⡄⢀⡀                                          ",
		"             ⢀⣀⣤⠰⢤⣼⡯⢽⡟⣀⢶⣺⡛⠁       ⠈⢰⠃⠁    ⢖⣒⣾⡟⠂  ⠈⠛⠁        ⠺⢩⢖⡄                ",
		"            ⡬⢍⣿⣟⣿⣻⣿⣿⣿⡾⣯⡀⠈⠁⠁⢦      ⢀⡿       ⠈       ⢠⢶⠘⠋⡁⣀⢠⠤⠖⠘⠉⠁⠈⠼⡧⡄⣄⡀ ⢫⣗⠒⠆      ",
		"⣓  ⣠⠖⠓⠒⠢⠤⢄⠤⠶⠽⠽⣶⣃⣽⡮⣿⡷⣗⣤⡭⣍⢓⡄ ⠸⣷   ⢀⣀⠿⠇       ⢀⠔⠒⠲⠄⢄⢀⡀⢙⣑⡄⠴⡍⣟⠉          ⠑⠉⠉  ⠑⠐⠦⠤⣤⠤⢞",
		"⠶⢧⣗⢾⡆         ⠈⠈⠁⠈⠉⢀⣹⣶⣩⣽⣐⢮⠃ ⣇ ⢀⡔⠊ ⢰⣖⣲    ⢀⡐⠁⣰⠦ ⢲⣶⠛⠋    ⠐⠋                      ⡤",
		"  ⠉⣮⣀⣀⣴⡤⣠⡀         ⡎ ⠛⢫⠙⢫⢫  ⠈⠦⠼          ⡃⡀⢸⠼⣤⡄                        ⡀⣀⣀⡐⡶⣣⢤⠖⠉",
		"   ⢀⡽⠟⠃  ⠈⠱⡀       ⠙⠢⣀⣨⠆⠈⠁⢧⡀          ⣸⣷ ⢹⣷⣼⣸⠃                       ⢀⡐⢀ ⠁⡚⣨⠆   ",
		"          ⠘⢳⡀        ⠈⠾  ⣀⣀⣽         ⠸⢼⣇⡧⠋⠉⠁                          ⠉⣿  ⠢⠂    ",
		"           ⠈⢻           ⠜⢹⣵⠻⠇         ⠈⢻  ⢀⡀  ⢠⣠⡤ ⢀⢤                  ⢰⣯        ",
		"            ⢼          ⢀⣾⠛⠉          ⠐⡖⠒⡰⢺⣞⣵⡄⢀⣏⡭⣙⡄⢕⢫⡀             ⢀ ⢠⠖⢱⡿⠃       ",
		"            ⠸         ⠠⡎             ⠰⣅⣰⣃⣘⡣⡿⢻⡿⣁⣀  ⠸⣽             ⠐⣿⣽⣫ ⡸⡇        ",
		"             ⠳⣄       ⡰⠃             ⢀⠎⠉  ⢧⡀⣠⣛⠈⢻                  ⢻⠘⢺⡿⠚⠁        ",
		"              ⢻⣇  ⣠⠲⠖⢲⡇              ⡸     ⠉⠃⠈⠉⣿  ⢰⣆              ⢸ ⠈⠁          ",
		"              ⠈⢿⣆ ⡟  ⣘⣻             ⡸          ⢸⢇ ⠈⠯⢿⡒⠲⡀   ⢀⡀    ⣀⢾             ",
		"    ⠈⢳          ⠸⡀⢳⣠⢾⠉⢹⣦⣤⣀          ⡇           ⡿⡄  ⢰⠃ ⠑⡂ ⢠⠏⢣  ⣼⡮⠁⢈⡀            ",
		"                 ⠙⠲⢆⡿⢦⠈⠉⠁⠁          ⡇           ⠱⣇⣀⠼⠃   ⡃⢰⠃ ⠸⢶ ⠘⠄ ⢾⡁            ",
		"                    ⠙⣾⣀⡴⡶⢤⣤         ⢳            ⠻⠵⡆    ⠸⣸   ⢸⡳⡤⠃⢀⡾⣿            ",
		"                     ⠘⢻⠁  ⠈⠦⣄        ⢧⣀⣀⠤⣀        ⢐⠁    ⠈⠩⠆  ⣘⣧⠁ ⡸⡔⢿            ",
		"                      ⡸     ⢨         ⠁  ⠉⡇      ⢀⠎          ⢻⢿⠄⡴⢑⣧⡠⡄           ",
		"                      ⡇     ⠈⠋⠦⡄         ⠈⡆     ⢠⠃            ⢏⡇⢧⣼⣾⣧⣽⣿⠶⢤⡀⣤      ",
		"                      ⣇        ⠈⡇         ⢸     ⢸             ⠈⠶⣦⣄⣋⣁⡀⠸⣵⢠⣻⠋⠷⣄    ",
		"                      ⠰⡀       ⣰⠁         ⢘⠆    ⢸ ⢠⡀              ⠙⠋⢠⠦⡄⣷⠙⠃ ⠙    ",
		"⠄                      ⠣⡀      ⡃          ⢸     ⣸⢡⢾⠆               ⡞⠛⠘⢧⡏⡆   ⠸⠄ ⡤",
		"                        ⠱     ⢠⠃          ⠸⡀   ⢸⠁⢸⢨              ⡤⠚     ⠱⡀  ⢦  ⠁",
		"                        ⠅    ⡖⠉            ⡇   ⡜ ⠸⠔              ⡇       ⢳      ",
		"                        ⡇   ⢀⠃             ⢱⡀ ⢰⠃                 ⣇  ⢀⡀   ⢸      ",
		"                       ⢀⠃  ⡦⠏              ⠈⠷⠖⠃                  ⠾⠴⠊⠁⠹⣦  ⡞    ⣄ ",
		"                       ⢸  ⡤⠃                                          ⠘⢲⠖⠃    ⣽⡆",
		"                       ⢸ ⣸⠁                                            ⠈⠿   ⢀⢼⠏ ",
		"                       ⠞ ⡗                             ⣄                    ⠈⠋  ",
		"                       ⢧⡼⡁⠲⠂                                                    ",
		"                        ⠙⠉                                                      ",
		"                           ⡀                                                    ",
		"                         ⣴⠏⠁                      ⣀⡤⢤⣀⣀  ⢀⣀⣤⣀⣀⡴⣄⡤⢤⣀⠤⠤⠴⣄⣀⡀       ",
		"                 ⣀⣀    ⣠⣿⡍⣆          ⣠⣤⣤⠤⠴⠶⠖⠲⠤⠔⠛⠒⠉   ⠈⠨⣇⠖⠋              ⠈⠉⠓⠢⠤⢄  ",
		"     ⡀ ⣠⠤⠴⠒⠚⠛⠛⠒⠢⠤⠿⠙⠉⠉⠑⢋⣚⣉⠥⠚      ⢀⣀⡠⠟⠁                                      ⡴⠋  ",
		"   ⠐⠶⣛⣫⡤              ⠐⢏⣀⣤⣤ ⣴⣋⢇⢀⣮⡥                                         ⣴⠓   ",
		"⠤⠤⠤⠤⡀⣈⢣⣠⡄                 ⠉⠊⠉⠉⠉                                            ⠈⠓⠆⠤⠤",
		"                                                                                ",
	})
	assertCellStyle(t, buf, 19, 1, style.NewStyle().Fg(style.Reset))
	assertCellStyle(t, buf, 0, 4, style.NewStyle().Fg(style.Reset))
}

func TestCanvas_shouldRenderCharMarkersWithOneCellResolution(t *testing.T) {
	tests := []struct {
		name     string
		marker   widgets.CanvasMarker
		expected string
	}{
		{name: "bar", marker: widgets.CanvasMarkerBar, expected: "▄"},
		{name: "custom", marker: widgets.CanvasMarkerCustom("x"), expected: "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

			widgets.NewCanvas().
				Marker(tt.marker).
				Paint(func(ctx *widgets.CanvasContext) {
					ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{{X: 0, Y: 0}}, style.Blue))
				}).
				Render(buf.Area, buf)

			assertLines(t, buf, []string{tt.expected})
			assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Blue))
		})
	}
}

func TestCanvas_shouldCombineBraillePointsInsideOneTerminalCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	widgets.NewCanvas().
		Marker(widgets.CanvasMarkerBraille).
		XBounds(0, 1).
		YBounds(0, 3).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{
				{X: 0, Y: 3},
				{X: 1, Y: 3},
				{X: 0, Y: 2},
				{X: 1, Y: 2},
			}, style.Red))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{"⠛"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red))
}

func TestCanvas_shouldRenderHalfBlockUpperLowerAndFullCells(t *testing.T) {
	tests := []struct {
		name     string
		points   []widgets.CanvasPoint
		expected string
		style    style.Style
	}{
		{
			name:     "upper",
			points:   []widgets.CanvasPoint{{X: 0, Y: 1}},
			expected: "▀",
			style:    style.NewStyle().Fg(style.Red),
		},
		{
			name:     "lower",
			points:   []widgets.CanvasPoint{{X: 0, Y: 0}},
			expected: "▄",
			style:    style.NewStyle().Fg(style.Red),
		},
		{
			name: "full",
			points: []widgets.CanvasPoint{
				{X: 0, Y: 1},
				{X: 0, Y: 0},
			},
			expected: "█",
			style:    style.NewStyle().Fg(style.Red).Bg(style.Red),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

			widgets.NewCanvas().
				Marker(widgets.CanvasMarkerHalfBlock).
				XBounds(0, 1).
				YBounds(0, 1).
				Paint(func(ctx *widgets.CanvasContext) {
					ctx.Draw(widgets.NewPoints(tt.points, style.Red))
				}).
				Render(buf.Area, buf)

			assertLines(t, buf, []string{tt.expected})
			assertCellStyle(t, buf, 0, 0, tt.style)
		})
	}
}

func TestCanvas_shouldRenderPatternMarkersFromPseudoPixels(t *testing.T) {
	tests := []struct {
		name     string
		marker   widgets.CanvasMarker
		yMax     float64
		points   []widgets.CanvasPoint
		expected string
	}{
		{
			name:   "quadrant",
			marker: widgets.CanvasMarkerQuadrant,
			yMax:   1,
			points: []widgets.CanvasPoint{
				{X: 0, Y: 1},
				{X: 1, Y: 0},
			},
			expected: "▚",
		},
		{
			name:   "sextant",
			marker: widgets.CanvasMarkerSextant,
			yMax:   2,
			points: []widgets.CanvasPoint{
				{X: 0, Y: 2},
				{X: 1, Y: 1},
				{X: 0, Y: 0},
			},
			expected: "🬗",
		},
		{
			name:   "octant",
			marker: widgets.CanvasMarkerOctant,
			yMax:   3,
			points: []widgets.CanvasPoint{
				{X: 0, Y: 3},
				{X: 0, Y: 2},
				{X: 1, Y: 1},
				{X: 1, Y: 0},
			},
			expected: "▚",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

			widgets.NewCanvas().
				Marker(tt.marker).
				XBounds(0, 1).
				YBounds(0, tt.yMax).
				Paint(func(ctx *widgets.CanvasContext) {
					ctx.Draw(widgets.NewPoints(tt.points, style.Green))
				}).
				Render(buf.Area, buf)

			assertLines(t, buf, []string{tt.expected})
			assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Green))
		})
	}
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

func TestCanvas_shouldDrawFilledLinesWithDotMarker(t *testing.T) {
	tests := []struct {
		name     string
		line     widgets.FilledLine
		expected []string
	}{
		{
			name: "off grid",
			line: widgets.NewFilledLine(-1, 0, -1, 10, 0, style.Red),
			expected: []string{
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
			},
		},
		{
			name: "horizontal fill to bottom",
			line: widgets.NewFilledLine(0, 0, 10, 0, 0, style.Red),
			expected: []string{
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"••••••••••",
			},
		},
		{
			name: "horizontal fill to top",
			line: widgets.NewFilledLine(0, 0, 10, 0, 10, style.Red),
			expected: []string{
				"••••••••••",
				"••••••••••",
				"••••••••••",
				"••••••••••",
				"••••••••••",
				"••••••••••",
				"••••••••••",
				"••••••••••",
				"••••••••••",
				"••••••••••",
			},
		},
		{
			name: "diagonal fill to bottom",
			line: widgets.NewFilledLine(0, 0, 10, 10, 0, style.Red),
			expected: []string{
				"         •",
				"        ••",
				"       •••",
				"      ••••",
				"     •••••",
				"    ••••••",
				"   •••••••",
				"  ••••••••",
				" •••••••••",
				"••••••••••",
			},
		},
		{
			name: "split fill",
			line: widgets.NewFilledLine(0, 0, 10, 10, 5, style.Red),
			expected: []string{
				"         •",
				"        ••",
				"       •••",
				"      ••••",
				"     •••••",
				"••••••••••",
				"••••      ",
				"•••       ",
				"••        ",
				"•         ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 10, 10))

			assertNotPanics(t, func() {
				widgets.NewCanvas().
					XBounds(0, 10).
					YBounds(0, 10).
					Marker(widgets.CanvasMarkerDot).
					Paint(func(ctx *widgets.CanvasContext) {
						ctx.Draw(tt.line)
					}).
					Render(buf.Area, buf)
			})

			assertLines(t, buf, tt.expected)
			for y, line := range tt.expected {
				for x, symbol := range []rune(line) {
					if symbol == '•' {
						assertCellStyle(t, buf, x, y, style.NewStyle().Fg(style.Red))
					}
				}
			}
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

func TestCanvas_shouldDrawCircleWithBrailleMarker(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 5))

	widgets.NewCanvas().
		Marker(widgets.CanvasMarkerBraille).
		XBounds(-10, 10).
		YBounds(-10, 10).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewCircle(5, 2, 5, style.Default))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"      ⣀⣀⣀ ",
		"     ⡞⠁ ⠈⢣",
		"     ⢇⡀ ⢀⡼",
		"      ⠉⠉⠉ ",
		"          ",
	})
}

func TestCanvas_shouldSkipOffGridCircleWithoutPanicking(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))

	assertNotPanics(t, func() {
		widgets.NewCanvas().
			XBounds(0, 2).
			YBounds(0, 2).
			Paint(func(ctx *widgets.CanvasContext) {
				ctx.Draw(widgets.NewCircle(10, 10, 1, style.Red))
			}).
			Render(buf.Area, buf)
	})

	assertLines(t, buf, []string{
		"   ",
		"   ",
		"   ",
	})
}

func TestCanvas_shouldDrawSmallCircleWithDotMarker(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))

	widgets.NewCanvas().
		Marker(widgets.CanvasMarkerDot).
		XBounds(0, 4).
		YBounds(0, 4).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.NewCircle(2, 2, 1, style.Green))
		}).
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"     ",
		" ••• ",
		" • • ",
		" ••• ",
		"     ",
	})
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
