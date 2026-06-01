package widgets_test

import (
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

func assertChartLines(t *testing.T, width, height int, xAxis widgets.Axis, yAxis widgets.Axis, expected []string) {
	t.Helper()
	buf := buffer.Empty(layout.NewRect(0, 0, width, height))
	chart := widgets.NewChart([]widgets.Dataset{}).XAxis(xAxis).YAxis(yAxis)

	chart.Render(buf.Area, buf)

	assertLines(t, buf, expected)
}
