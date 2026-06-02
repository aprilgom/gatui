package widgets_test

import (
	"strconv"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestChart_shouldNotPanicOnSmallAreas(t *testing.T) {
	for _, size := range []layout.Size{
		{Width: 0, Height: 0},
		{Width: 0, Height: 1},
		{Width: 1, Height: 0},
		{Width: 1, Height: 1},
		{Width: 2, Height: 2},
	} {
		t.Run("", func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, size.Width, size.Height))
			chart := widgets.NewChart([]widgets.Dataset{
				widgets.NewDataset().
					Style(style.NewStyle().Fg(style.Magenta)).
					Data([]layout.Position{{X: 0, Y: 0}}),
			}).
				Block(widgets.BorderedBlock().Title(text.LineFromString("Plot"))).
				XAxis(widgets.NewAxis().
					Bounds(0, 0).
					LabelStrings([]string{"0.0", "1.0"})).
				YAxis(widgets.NewAxis().
					Bounds(0, 0).
					LabelStrings([]string{"0.0", "1.0"}))

			assertNotPanics(t, func() {
				chart.Render(buf.Area, buf)
			})
		})
	}
}

func TestChart_shouldHandleLongLabels(t *testing.T) {
	tests := []struct {
		name     string
		xLabels  []string
		yLabels  []string
		xAlign   layout.Alignment
		expected []string
	}{
		{
			name:    "x left labels",
			xLabels: []string{"AAAA", "B"},
			xAlign:  layout.Left,
			expected: []string{
				"          ",
				"          ",
				"          ",
				"   ───────",
				"AAA      B",
			},
		},
		{
			name:    "x right label wider than slot",
			xLabels: []string{"A", "BBBB"},
			xAlign:  layout.Left,
			expected: []string{
				"          ",
				"          ",
				"          ",
				" ─────────",
				"A     BBBB",
			},
		},
		{
			name:    "x first label truncated",
			xLabels: []string{"AAAAAAAAAAA", "B"},
			xAlign:  layout.Left,
			expected: []string{
				"          ",
				"          ",
				"          ",
				"   ───────",
				"AAA      B",
			},
		},
		{
			name:    "x and y labels",
			xLabels: []string{"A", "B"},
			yLabels: []string{"CCCCCCC", "D"},
			xAlign:  layout.Left,
			expected: []string{
				"D  │      ",
				"   │      ",
				"CCC│      ",
				"   └──────",
				"   A     B",
			},
		},
		{
			name:    "center x alignment with y labels",
			xLabels: []string{"AAAAAAAAAA", "B"},
			yLabels: []string{"C", "D"},
			xAlign:  layout.Center,
			expected: []string{
				"D  │      ",
				"   │      ",
				"C  │      ",
				"   └──────",
				"AAAAAAA  B",
			},
		},
		{
			name:    "right x alignment with y labels",
			xLabels: []string{"AAAAAAA", "B"},
			yLabels: []string{"C", "D"},
			xAlign:  layout.Right,
			expected: []string{
				"D│        ",
				" │        ",
				"C│        ",
				" └────────",
				" AAAAA   B",
			},
		},
		{
			name:    "right x alignment with long last label",
			xLabels: []string{"AAAAAAA", "BBBBBBB"},
			yLabels: []string{"C", "D"},
			xAlign:  layout.Right,
			expected: []string{
				"D│        ",
				" │        ",
				"C│        ",
				" └────────",
				" AAAAABBBB",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xAxis := widgets.NewAxis().Bounds(0, 1).LabelsAlignment(tt.xAlign)
			if tt.xLabels != nil {
				xAxis = xAxis.LabelStrings(tt.xLabels)
			}
			yAxis := widgets.NewAxis().Bounds(0, 1)
			if tt.yLabels != nil {
				yAxis = yAxis.LabelStrings(tt.yLabels)
			}
			assertChartLines(t, 10, 5, xAxis, yAxis, tt.expected)
		})
	}
}

func TestChart_shouldHandleXAxisLabelAlignment(t *testing.T) {
	tests := []struct {
		name      string
		alignment layout.Alignment
		expected  []string
	}{
		{
			name:      "left",
			alignment: layout.Left,
			expected: []string{
				"          ",
				"          ",
				"          ",
				"   ───────",
				"AAA   B  C",
			},
		},
		{
			name:      "center",
			alignment: layout.Center,
			expected: []string{
				"          ",
				"          ",
				"          ",
				"  ────────",
				"AAAA B   C",
			},
		},
		{
			name:      "right",
			alignment: layout.Right,
			expected: []string{
				"          ",
				"          ",
				"          ",
				"──────────",
				"AAA B    C",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xAxis := widgets.NewAxis().
				LabelStrings([]string{"AAAA", "B", "C"}).
				LabelsAlignment(tt.alignment)

			assertChartLines(t, 10, 5, xAxis, widgets.NewAxis(), tt.expected)
		})
	}
}

func TestChart_shouldHandleYAxisLabelAlignment(t *testing.T) {
	tests := []struct {
		name      string
		alignment layout.Alignment
		expected  []string
	}{
		{
			name:      "left",
			alignment: layout.Left,
			expected: []string{
				"D   │               ",
				"    │               ",
				"C   │               ",
				"    └───────────────",
				"AAAAA              B",
			},
		},
		{
			name:      "center",
			alignment: layout.Center,
			expected: []string{
				" D  │               ",
				"    │               ",
				" C  │               ",
				"    └───────────────",
				"AAAAA              B",
			},
		},
		{
			name:      "right",
			alignment: layout.Right,
			expected: []string{
				"   D│               ",
				"    │               ",
				"   C│               ",
				"    └───────────────",
				"AAAAA              B",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xAxis := widgets.NewAxis().LabelStrings([]string{"AAAAA", "B"})
			yAxis := widgets.NewAxis().
				LabelStrings([]string{"C", "D"}).
				LabelsAlignment(tt.alignment)

			assertChartLines(t, 20, 5, xAxis, yAxis, tt.expected)
		})
	}
}

func TestChart_shouldAllowZeroLengthBounds(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 100, 100))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().Data([]layout.Position{{X: 0, Y: 0}}),
	}).
		Block(widgets.BorderedBlock().Title(text.LineFromString("Plot"))).
		XAxis(widgets.NewAxis().Bounds(0, 0).LabelStrings([]string{"0.0", "1.0"})).
		YAxis(widgets.NewAxis().Bounds(0, 0).LabelStrings([]string{"0.0", "1.0"}))

	assertNotPanics(t, func() {
		chart.Render(buf.Area, buf)
	})
}

func TestChart_shouldStyleTopLine(t *testing.T) {
	titleStyle := style.NewStyle().Fg(style.Red).Bg(style.LightBlue)
	buf := buffer.Empty(layout.NewRect(0, 0, 9, 5))
	chart := widgets.NewChart([]widgets.Dataset{}).
		YAxis(widgets.NewAxis().
			Title(text.NewLine(text.StyledSpan("abc", titleStyle))).
			Bounds(0, 1).
			LabelStrings([]string{"a", "b"})).
		XAxis(widgets.NewAxis().Bounds(0, 1))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"b│abc    ",
		" │       ",
		" │       ",
		" │       ",
		"a│       ",
	})
	for x := 2; x <= 4; x++ {
		assertCellStyle(t, buf, x, 0, titleStyle)
	}
}

func TestChart_shouldStyleTopLineWithLineDataset(t *testing.T) {
	titleStyle := style.NewStyle().Fg(style.Red).Bg(style.LightBlue)
	dataStyle := style.NewStyle().Fg(style.Blue)
	buf := buffer.Empty(layout.NewRect(0, 0, 9, 5))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeLine).
			Style(dataStyle).
			DataPoints([]widgets.ChartPoint{
				{X: 0, Y: 1},
				{X: 1, Y: 1},
			}),
	}).
		YAxis(widgets.NewAxis().
			Title(text.NewLine(text.StyledSpan("abc", titleStyle))).
			Bounds(0, 1).
			LabelStrings([]string{"a", "b"})).
		XAxis(widgets.NewAxis().Bounds(0, 1))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"b│abc••••",
		" │       ",
		" │       ",
		" │       ",
		"a│       ",
	})
	for x := 2; x <= 4; x++ {
		assertCellStyle(t, buf, x, 0, titleStyle)
	}
	for x := 5; x <= 8; x++ {
		assertCellStyle(t, buf, x, 0, dataStyle)
	}
}

func TestChart_shouldAllowEmptyLineDataset(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 5))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeLine).
			DataPoints([]widgets.ChartPoint{}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"})).
		YAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"}))

	assertNotPanics(t, func() {
		chart.Render(buf.Area, buf)
	})
}

func TestChart_shouldHandlePlottingOverflows(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 5))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeLine).
			DataPoints([]widgets.ChartPoint{
				{X: -1_000_000, Y: 0},
				{X: 1, Y: 1},
				{X: 1_000_000, Y: 2},
			}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 1_000_000_000).LabelStrings([]string{"0", "1B"})).
		YAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"}))

	assertNotPanics(t, func() {
		chart.Render(buf.Area, buf)
	})
}

func TestChart_shouldPlotScatterPoint(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 9, 5))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			Style(style.NewStyle().Fg(style.Green)).
			DataPoints([]widgets.ChartPoint{{X: 0.5, Y: 0.5}}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 1)).
		YAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"}))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"1│       ",
		" │       ",
		" │   •   ",
		" │       ",
		"0│       ",
	})
	assertCellStyle(t, buf, 5, 2, style.NewStyle().Fg(style.Green))
}

func TestChart_shouldPlotLineBetweenTwoPoints(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 3))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeLine).
			DataPoints([]widgets.ChartPoint{
				{X: 0, Y: 0.5},
				{X: 1, Y: 0.5},
			}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 1)).
		YAxis(widgets.NewAxis().Bounds(0, 1))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"      ",
		"••••••",
		"      ",
	})
}

func TestChart_shouldRenderBarGraphType(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 11))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeBar).
			DataPoints([]widgets.ChartPoint{
				{X: 0, Y: 0},
				{X: 2, Y: 1},
				{X: 4, Y: 4},
				{X: 6, Y: 8},
				{X: 8, Y: 9},
				{X: 10, Y: 10},
			}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 10)).
		YAxis(widgets.NewAxis().Bounds(0, 10))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"          •",
		"        • •",
		"      • • •",
		"      • • •",
		"      • • •",
		"      • • •",
		"    • • • •",
		"    • • • •",
		"    • • • •",
		"  • • • • •",
		"• • • • • •",
	})
}

func TestChart_shouldRenderAreaGraphType(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 11))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeArea).
			FillToY(0).
			DataPoints([]widgets.ChartPoint{
				{X: 0, Y: 0},
				{X: 5, Y: 5},
				{X: 10, Y: 5},
			}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 10)).
		YAxis(widgets.NewAxis().Bounds(0, 10))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"           ",
		"           ",
		"           ",
		"           ",
		"           ",
		"     ••••••",
		"    •••••••",
		"   ••••••••",
		"  •••••••••",
		" ••••••••••",
		"•••••••••••",
	})
}

func TestChart_shouldAllowEmptyBarAndAreaDatasets(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 5))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeBar).
			DataPoints([]widgets.ChartPoint{}),
		widgets.NewDataset().
			GraphType(widgets.GraphTypeArea).
			DataPoints([]widgets.ChartPoint{}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"})).
		YAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"}))

	assertNotPanics(t, func() {
		chart.Render(buf.Area, buf)
	})
}

func TestChart_shouldClampAreaFillToYToBounds(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			GraphType(widgets.GraphTypeArea).
			FillToY(-10).
			DataPoints([]widgets.ChartPoint{
				{X: 0, Y: 4},
				{X: 4, Y: 4},
			}),
	}).
		XAxis(widgets.NewAxis().Bounds(0, 4)).
		YAxis(widgets.NewAxis().Bounds(0, 4))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"•••••",
		"•••••",
		"•••••",
		"•••••",
		"•••••",
	})
}

func TestChart_shouldRenderOverlappingLineDatasetsWithDatasetMarkers(t *testing.T) {
	tests := []struct {
		name              string
		marker            widgets.CanvasMarker
		expected          []string
		redMarkerCells    []layout.Position
		blueBlockCells    []layout.Position
		centerOverlapCell layout.Position
	}{
		{
			name:   "dot",
			marker: widgets.CanvasMarkerDot,
			expected: []string{
				"•   █",
				" • █ ",
				"  •  ",
				" █ • ",
				"█   •",
			},
			redMarkerCells: []layout.Position{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 3, Y: 3},
				{X: 4, Y: 4},
			},
			blueBlockCells: []layout.Position{
				{X: 4, Y: 0},
				{X: 3, Y: 1},
				{X: 1, Y: 3},
				{X: 0, Y: 4},
			},
			centerOverlapCell: layout.Position{X: 2, Y: 2},
		},
		{
			name:   "braille",
			marker: widgets.CanvasMarkerBraille,
			expected: []string{
				"⢣   █",
				" ⢣ █ ",
				"  ⢣  ",
				" █ ⢣ ",
				"█   ⢣",
			},
			redMarkerCells: []layout.Position{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 3, Y: 3},
				{X: 4, Y: 4},
			},
			blueBlockCells: []layout.Position{
				{X: 4, Y: 0},
				{X: 3, Y: 1},
				{X: 1, Y: 3},
				{X: 0, Y: 4},
			},
			centerOverlapCell: layout.Position{X: 2, Y: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, 5))
			chart := widgets.NewChart([]widgets.Dataset{
				widgets.NewDataset().
					GraphType(widgets.GraphTypeLine).
					Marker(widgets.CanvasMarkerBlock).
					Style(style.NewStyle().Fg(style.Blue)).
					DataPoints([]widgets.ChartPoint{
						{X: 0, Y: 0},
						{X: 5, Y: 5},
					}),
				widgets.NewDataset().
					GraphType(widgets.GraphTypeLine).
					Marker(tt.marker).
					Style(style.NewStyle().Fg(style.Red)).
					DataPoints([]widgets.ChartPoint{
						{X: 0, Y: 5},
						{X: 5, Y: 0},
					}),
			}).
				XAxis(widgets.NewAxis().Bounds(0, 5)).
				YAxis(widgets.NewAxis().Bounds(0, 5))

			chart.Render(buf.Area, buf)

			assertLines(t, buf, tt.expected)
			for _, point := range tt.redMarkerCells {
				assertCellStyle(t, buf, point.X, point.Y, style.NewStyle().Fg(style.Red))
			}
			for _, point := range tt.blueBlockCells {
				assertCellStyle(t, buf, point.X, point.Y, style.NewStyle().Fg(style.Blue).Bg(style.Blue))
			}
			assertCellStyle(t, buf, tt.centerOverlapCell.X, tt.centerOverlapCell.Y, style.NewStyle().Fg(style.Red).Bg(style.Blue))
		})
	}
}

func TestChart_datasetsWithoutNameShouldNotContributeToLegendHeight(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 50, 25))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().Name("data1"),
		widgets.NewDataset(),
		widgets.NewDataset().Name(""),
	})

	chart.Render(buf.Area, buf)

	assertLines(t, firstLines(buf, 4), []string{
		"                                           ┌─────┐",
		"                                           │data1│",
		"                                           │     │",
		"                                           └─────┘",
	})
}

func TestChart_shouldNotRenderLegend_whenNoDatasetsAreNamed(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 4))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset(),
		widgets.NewDataset(),
		widgets.NewDataset(),
	})

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"            ",
		"            ",
		"            ",
		"            ",
	})
}

func TestChart_shouldPatchDatasetStyleIntoLegendName(t *testing.T) {
	longStyle := style.NewStyle().Fg(style.Red)
	shortStyle := style.NewStyle().Fg(style.Green)
	buf := buffer.Empty(layout.NewRect(0, 0, 20, 5))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().Name("Very long name").Style(longStyle),
		widgets.NewDataset().Name("Short name").Style(shortStyle),
	}).HiddenLegendConstraints(layout.Length(100), layout.Length(100))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"    ┌──────────────┐",
		"    │Very long name│",
		"    │Short name    │",
		"    └──────────────┘",
		"                    ",
	})
	for x := 5; x <= 18; x++ {
		assertCellStyle(t, buf, x, 1, longStyle)
	}
	for x := 5; x <= 14; x++ {
		assertCellStyle(t, buf, x, 2, shortStyle)
	}
}

func TestAxis_canBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 4))
	axis := widgets.NewAxis().
		Bounds(0, 1).
		LabelStrings([]string{"0", "1"}).
		TitleString("Y").
		Fg(style.Black).
		Bg(style.White).
		Bold().
		Dim().
		Italic().
		Cyan()
	want := style.NewStyle().
		Fg(style.Cyan).
		Bg(style.White).
		AddModifier(style.ModifierBold | style.ModifierDim | style.ModifierItalic)
	chart := widgets.NewChart([]widgets.Dataset{}).YAxis(axis)

	chart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 1, 0, want)
	assertCellStyle(t, buf, 0, 3, want)
	assertCellStyle(t, buf, 2, 0, want)
}

func TestDataset_canBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 6))
	dataset := widgets.NewDataset().
		Name("D").
		DataPoints([]widgets.ChartPoint{{X: 0.5, Y: 0.5}}).
		Fg(style.Black).
		Bg(style.White).
		Bold().
		Dim().
		Italic().
		Cyan()
	want := style.NewStyle().
		Fg(style.Cyan).
		Bg(style.White).
		AddModifier(style.ModifierBold | style.ModifierDim | style.ModifierItalic)
	chart := widgets.NewChart([]widgets.Dataset{dataset}).
		HiddenLegendConstraints(layout.Length(100), layout.Length(100)).
		XAxis(widgets.NewAxis().Bounds(0, 1)).
		YAxis(widgets.NewAxis().Bounds(0, 1))

	chart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 6, 3, style.NewStyle().Fg(style.Cyan))
	assertCellStyle(t, buf, 10, 1, want)
}

func TestChart_canBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 2))
	chart := widgets.NewChart([]widgets.Dataset{}).
		Fg(style.Black).
		Bg(style.White).
		Bold().
		Dim().
		Italic().
		Cyan()
	want := style.NewStyle().
		Fg(style.Cyan).
		Bg(style.White).
		AddModifier(style.ModifierBold | style.ModifierDim | style.ModifierItalic)

	chart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 0, want)
	assertCellStyle(t, buf, 3, 1, want)
}

func TestChart_shouldApplyChartStyleToPlotBackground(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 5))
	chart := widgets.NewChart([]widgets.Dataset{}).
		Style(style.NewStyle().Bg(style.Blue)).
		XAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"})).
		YAxis(widgets.NewAxis().Bounds(0, 1).LabelStrings([]string{"0", "1"}))

	chart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 5, 1, style.NewStyle().Bg(style.Blue))
}

func TestChart_shouldPatchAxisStyleOverChartStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 5))
	chart := widgets.NewChart([]widgets.Dataset{}).
		Style(style.NewStyle().Bg(style.Blue)).
		XAxis(widgets.NewAxis().
			Bounds(0, 1).
			LabelStrings([]string{"0", "1"}).
			TitleString("X").
			Fg(style.Red)).
		YAxis(widgets.NewAxis().
			Bounds(0, 1).
			LabelStrings([]string{"0", "1"}).
			Fg(style.Green))

	chart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 1, 2, style.NewStyle().Fg(style.Green).Bg(style.Blue))
	assertCellStyle(t, buf, 2, 3, style.NewStyle().Fg(style.Red).Bg(style.Blue))
	assertCellStyle(t, buf, 0, 2, style.NewStyle().Fg(style.Green).Bg(style.Blue))
	assertCellStyle(t, buf, 9, 4, style.NewStyle().Fg(style.Red).Bg(style.Blue))
}

func TestChart_shouldPatchDatasetStyleOverChartStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 6))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().
			Name("D").
			DataPoints([]widgets.ChartPoint{{X: 0.5, Y: 0.5}}).
			Fg(style.Red),
	}).
		Style(style.NewStyle().Bg(style.Blue)).
		HiddenLegendConstraints(layout.Length(100), layout.Length(100)).
		XAxis(widgets.NewAxis().Bounds(0, 1)).
		YAxis(widgets.NewAxis().Bounds(0, 1))

	chart.Render(buf.Area, buf)

	assertCellStyle(t, buf, 6, 3, style.NewStyle().Fg(style.Red).Bg(style.Blue))
	assertCellStyle(t, buf, 10, 1, style.NewStyle().Fg(style.Red).Bg(style.Blue))
}

func TestChart_shouldRenderStyledAxisLabelsWithCellWidthClipping(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 5))
	axisStyle := style.NewStyle().Bg(style.Blue)
	label := text.NewLine(
		text.StyledSpan("A", style.NewStyle().Fg(style.Red)),
		text.StyledSpan("コ", style.NewStyle().Fg(style.Green)),
	).Style(style.NewStyle().AddModifier(style.ModifierBold))
	chart := widgets.NewChart([]widgets.Dataset{}).
		XAxis(widgets.NewAxis().Bounds(0, 1).Labels([]text.Line{label, text.LineFromString("Z")}).Style(axisStyle))

	chart.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"      ",
		"      ",
		"      ",
		"  ────",
		"A    Z",
	})
	assertCellStyle(t, buf, 0, 4, style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold))
	assertCellSymbol(t, buf, 1, 4, " ")
}

func TestChart_shouldRenderTopLeftLegend(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 30, 20))
	chart := widgets.NewChart([]widgets.Dataset{
		widgets.NewDataset().Name("Ds1"),
	}).LegendPosition(widgets.LegendPositionTopLeft)

	chart.Render(buf.Area, buf)

	lines := firstLines(buf, 3)
	assertLines(t, lines, []string{
		"┌───┐                         ",
		"│Ds1│                         ",
		"└───┘                         ",
	})
}

func TestChart_shouldHideLegend_whenHiddenLegendConstraintsAreExceeded(t *testing.T) {
	datasets := make([]widgets.Dataset, 0, 10)
	for i := range 10 {
		datasets = append(datasets, widgets.NewDataset().Name("Dataset #"+strconv.Itoa(i)))
	}

	shown := buffer.Empty(layout.NewRect(0, 0, 100, 100))
	widgets.NewChart(datasets).Render(shown.Area, shown)
	assertLines(t, firstLines(shown, 1), []string{
		"                                                                                        ┌──────────┐",
	})

	hidden := buffer.Empty(layout.NewRect(0, 0, 100, 100))
	widgets.NewChart(datasets).
		HiddenLegendConstraints(layout.Ratio(1, 10), layout.Ratio(1, 4)).
		Render(hidden.Area, hidden)
	assertLines(t, firstLines(hidden, 1), []string{
		"                                                                                                    ",
	})
}

func assertChartLines(t *testing.T, width, height int, xAxis widgets.Axis, yAxis widgets.Axis, expected []string) {
	t.Helper()
	buf := buffer.Empty(layout.NewRect(0, 0, width, height))
	chart := widgets.NewChart([]widgets.Dataset{}).XAxis(xAxis).YAxis(yAxis)

	chart.Render(buf.Area, buf)

	assertLines(t, buf, expected)
}

func firstLines(buf *buffer.Buffer, n int) *buffer.Buffer {
	lines := buf.Lines()
	return buffer.WithLines(lines[:n])
}
