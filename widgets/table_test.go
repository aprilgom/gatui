package widgets_test

import (
	"fmt"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestTable_shouldRenderColumnSpacingHeaderAndBlock(t *testing.T) {
	tests := []struct {
		name          string
		columnSpacing int
		expected      []string
	}{
		{
			name:          "no space between columns",
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1Head2Head3             │",
				"│                            │",
				"│Row11Row12Row13             │",
				"│Row21Row22Row23             │",
				"│Row31Row32Row33             │",
				"│Row41Row42Row43             │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name:          "one space between columns",
			columnSpacing: 1,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1 Head2 Head3           │",
				"│                            │",
				"│Row11 Row12 Row13           │",
				"│Row21 Row22 Row23           │",
				"│Row31 Row32 Row33           │",
				"│Row41 Row42 Row43           │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name:          "large spacing before pushing a column off",
			columnSpacing: 6,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1      Head2      Head3 │",
				"│                            │",
				"│Row11      Row12      Row13 │",
				"│Row21      Row22      Row23 │",
				"│Row31      Row32      Row33 │",
				"│Row41      Row42      Row43 │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name:          "large spacing pushes part of third column off",
			columnSpacing: 7,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1       Head       Head3│",
				"│                            │",
				"│Row11       Row1       Row13│",
				"│Row21       Row2       Row23│",
				"│Row31       Row3       Row33│",
				"│Row41       Row4       Row43│",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 30, 10))
			tableFixture([]layout.Constraint{
				layout.Length(5),
				layout.Length(5),
				layout.Length(5),
			}).ColumnSpacing(tt.columnSpacing).Render(buf.Area, buf)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestTable_shouldResolveFixedPercentageMixedAndRatioWidths(t *testing.T) {
	tests := []struct {
		name          string
		widths        []layout.Constraint
		columnSpacing int
		expected      []string
	}{
		{
			name: "fixed zero widths",
			widths: []layout.Constraint{
				layout.Length(0),
				layout.Length(0),
				layout.Length(0),
			},
			columnSpacing: 1,
			expected: []string{
				"┌────────────────────────────┐",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "fixed slim columns",
			widths: []layout.Constraint{
				layout.Length(1),
				layout.Length(1),
				layout.Length(1),
			},
			columnSpacing: 1,
			expected: []string{
				"┌────────────────────────────┐",
				"│H H H                       │",
				"│                            │",
				"│R R R                       │",
				"│R R R                       │",
				"│R R R                       │",
				"│R R R                       │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "percentage equal widths",
			widths: []layout.Constraint{
				layout.Percentage(50),
				layout.Percentage(50),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1         Head2         │",
				"│                            │",
				"│Row11         Row12         │",
				"│Row21         Row22         │",
				"│Row31         Row32         │",
				"│Row41         Row42         │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "percentage zero widths",
			widths: []layout.Constraint{
				layout.Percentage(0),
				layout.Percentage(0),
				layout.Percentage(0),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "percentage slim columns",
			widths: []layout.Constraint{
				layout.Percentage(11),
				layout.Percentage(11),
				layout.Percentage(11),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│HeaHeaHea                   │",
				"│                            │",
				"│RowRowRow                   │",
				"│RowRowRow                   │",
				"│RowRowRow                   │",
				"│RowRowRow                   │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "percentage thirds",
			widths: []layout.Constraint{
				layout.Percentage(33),
				layout.Percentage(33),
				layout.Percentage(33),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1    Head2    Head3     │",
				"│                            │",
				"│Row11    Row12    Row13     │",
				"│Row21    Row22    Row23     │",
				"│Row31    Row32    Row33     │",
				"│Row41    Row42    Row43     │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "mixed zero widths",
			widths: []layout.Constraint{
				layout.Percentage(0),
				layout.Length(0),
				layout.Percentage(0),
			},
			columnSpacing: 1,
			expected: []string{
				"┌────────────────────────────┐",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "mixed slim columns",
			widths: []layout.Constraint{
				layout.Percentage(11),
				layout.Length(20),
				layout.Percentage(11),
			},
			columnSpacing: 1,
			expected: []string{
				"┌────────────────────────────┐",
				"│Hea Head2                Hea│",
				"│                            │",
				"│Row Row12                Row│",
				"│Row Row22                Row│",
				"│Row Row32                Row│",
				"│Row Row42                Row│",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "mixed constraints",
			widths: []layout.Constraint{
				layout.Percentage(33),
				layout.Length(10),
				layout.Percentage(33),
			},
			columnSpacing: 1,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1     Head2      Head3  │",
				"│                            │",
				"│Row11     Row12      Row13  │",
				"│Row21     Row22      Row23  │",
				"│Row31     Row32      Row33  │",
				"│Row41     Row42      Row43  │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "mixed more than one hundred percent",
			widths: []layout.Constraint{
				layout.Percentage(60),
				layout.Length(10),
				layout.Percentage(60),
			},
			columnSpacing: 1,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1      Head2      Head3 │",
				"│                            │",
				"│Row11      Row12      Row13 │",
				"│Row21      Row22      Row23 │",
				"│Row31      Row32      Row33 │",
				"│Row41      Row42      Row43 │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "ratio zero widths",
			widths: []layout.Constraint{
				layout.Ratio(0, 1),
				layout.Ratio(0, 1),
				layout.Ratio(0, 1),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "ratio slim columns",
			widths: []layout.Constraint{
				layout.Ratio(1, 9),
				layout.Ratio(1, 9),
				layout.Ratio(1, 9),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│HeaHeaHea                   │",
				"│                            │",
				"│RowRowRow                   │",
				"│RowRowRow                   │",
				"│RowRowRow                   │",
				"│RowRowRow                   │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "ratio thirds",
			widths: []layout.Constraint{
				layout.Ratio(1, 3),
				layout.Ratio(1, 3),
				layout.Ratio(1, 3),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1    Head2     Head3    │",
				"│                            │",
				"│Row11    Row12     Row13    │",
				"│Row21    Row22     Row23    │",
				"│Row31    Row32     Row33    │",
				"│Row41    Row42     Row43    │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
		{
			name: "ratio halves",
			widths: []layout.Constraint{
				layout.Ratio(1, 2),
				layout.Ratio(1, 2),
			},
			columnSpacing: 0,
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1         Head2         │",
				"│                            │",
				"│Row11         Row12         │",
				"│Row21         Row22         │",
				"│Row31         Row32         │",
				"│Row41         Row42         │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 30, 10))
			tableFixture(tt.widths).ColumnSpacing(tt.columnSpacing).Render(buf.Area, buf)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestTable_excessAreaHighlightSymbolAndColumnSpacingAllocation(t *testing.T) {
	assertTableWithSelection(t, widgets.HighlightSpacingNever, 15, 0, nil, []string{
		"ABCDE12345     ",
		"               ",
		"               ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingNever, 15, 0, new(0), []string{
		"ABCDE12345     ",
		"               ",
		"               ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingWhenSelected, 15, 0, nil, []string{
		"ABCDE12345     ",
		"               ",
		"               ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingWhenSelected, 15, 0, new(0), []string{
		">>>ABCDE12345  ",
		"               ",
		"               ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 15, 0, nil, []string{
		"   ABCDE12345  ",
		"               ",
		"               ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 15, 0, new(0), []string{
		">>>ABCDE12345  ",
		"               ",
		"               ",
	})
}

func TestTable_insufficientAreaHighlightSymbolAllocationWithNoColumnSpacing(t *testing.T) {
	assertTableWithSelection(t, widgets.HighlightSpacingNever, 10, 0, nil, []string{
		"ABCDE12345",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingWhenSelected, 10, 0, nil, []string{
		"ABCDE12345",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 10, 0, nil, []string{
		"   ABCD123",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingNever, 10, 0, new(0), []string{
		"ABCDE12345",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingWhenSelected, 10, 0, new(0), []string{
		">>>ABCD123",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 10, 0, new(0), []string{
		">>>ABCD123",
		"          ",
		"          ",
	})
}

func TestTable_insufficientAreaHighlightSymbolAndColumnSpacingAllocation(t *testing.T) {
	assertTableWithSelection(t, widgets.HighlightSpacingNever, 10, 1, nil, []string{
		"ABCDE 1234",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingWhenSelected, 10, 1, nil, []string{
		"ABCDE 1234",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 10, 1, nil, []string{
		"   ABC 123",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 9, 1, nil, []string{
		"   ABC 12",
		"         ",
		"         ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 8, 1, nil, []string{
		"   AB 12",
		"        ",
		"        ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 7, 1, nil, []string{
		"   AB 1",
		"       ",
		"       ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingNever, 10, 1, new(0), []string{
		"ABCDE 1234",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingWhenSelected, 10, 1, new(0), []string{
		">>>ABC 123",
		"          ",
		"          ",
	})
	assertTableWithSelection(t, widgets.HighlightSpacingAlways, 10, 1, new(0), []string{
		">>>ABC 123",
		"          ",
		"          ",
	})
}

func TestTable_maxConstraint(t *testing.T) {
	assertConstraintTable(t, 20, nil, []layout.Constraint{layout.Max(4), layout.Max(4)}, []string{
		"ABCD 1234           ",
	})
	assertConstraintTable(t, 20, new(0), []layout.Constraint{layout.Max(4), layout.Max(4)}, []string{
		">>>ABCD 1234        ",
	})
	assertConstraintTable(t, 7, nil, []layout.Constraint{layout.Max(4), layout.Max(4)}, []string{
		"ABC 123",
	})
	assertConstraintTable(t, 7, new(0), []layout.Constraint{layout.Max(4), layout.Max(4)}, []string{
		">>>AB 1",
	})
}

func TestTable_minConstraint(t *testing.T) {
	assertConstraintTable(t, 20, nil, []layout.Constraint{layout.Min(4), layout.Min(4)}, []string{
		"ABCDE      12345    ",
	})
	assertConstraintTable(t, 20, new(0), []layout.Constraint{layout.Min(4), layout.Min(4)}, []string{
		">>>ABCDE    12345   ",
	})
	assertConstraintTable(t, 7, nil, []layout.Constraint{layout.Min(4), layout.Min(4)}, []string{
		"ABC 123",
	})
	assertConstraintTable(t, 7, new(0), []layout.Constraint{layout.Min(4), layout.Min(4)}, []string{
		">>>AB 1",
	})
}

func TestTable_underconstrainedFlex(t *testing.T) {
	assertConstraintTable(t, 62, nil, []layout.Constraint{layout.Min(10), layout.Min(10), layout.Min(1)}, []string{
		"ABCDE                12345                Z                   ",
	})
}

func TestTable_underconstrainedSegmentSize(t *testing.T) {
	assertConstraintTable(t, 62, nil, []layout.Constraint{layout.Min(10), layout.Min(10), layout.Min(1)}, []string{
		"ABCDE                12345                Z                   ",
	})
	assertConstraintTable(t, 23, nil, []layout.Constraint{layout.Min(10), layout.Min(10), layout.Min(1)}, []string{
		"ABCDE      12345      Z",
	})
}

func TestTable_new(t *testing.T) {
	rows := []widgets.TableRow{
		widgets.NewTableRow([]widgets.TableCell{widgets.TableCellFromString("A")}),
	}
	widths := []layout.Constraint{layout.Percentage(100)}
	table := widgets.NewTable(rows, widths)
	rows[0] = widgets.NewTableRow([]widgets.TableCell{widgets.TableCellFromString("B")})
	widths[0] = layout.Length(1)

	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{"A  "})
}

func TestTable_default(t *testing.T) {
	buf := buffer.WithLines([]string{"abc", "def"})

	widgets.NewTable([]widgets.TableRow{}, []layout.Constraint{}).Render(buf.Area, buf)

	assertLines(t, buf, []string{"abc", "def"})
	assertAllCellsStyle(t, buf, style.NewStyle())
}

func TestTable_collectRows(t *testing.T) {
	rows := make([]widgets.TableRow, 0, 4)
	for i := 0; i < 4; i++ {
		cells := make([]string, 0, 4)
		for j := 0; j < 4; j++ {
			cells = append(cells, fmt.Sprintf("%d*%d = %d", i, j, i*j))
		}
		rows = append(rows, widgets.TableRowFromStrings(cells))
	}
	widths := []layout.Constraint{
		layout.Percentage(25),
		layout.Percentage(25),
		layout.Percentage(25),
		layout.Percentage(25),
	}
	table := widgets.NewTable(rows, widths)
	rows[0] = widgets.TableRowFromStrings([]string{"mutated"})
	widths[0] = layout.Length(1)

	buf := buffer.Empty(layout.NewRect(0, 0, 40, 4))
	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"0*0 = 0   0*1 = 0   0*2 = 0   0*3 = 0   ",
		"1*0 = 0   1*1 = 1   1*2 = 2   1*3 = 3   ",
		"2*0 = 0   2*1 = 2   2*2 = 4   2*3 = 6   ",
		"3*0 = 0   3*1 = 3   3*2 = 6   3*3 = 9   ",
	})
}

func TestTable_rows(t *testing.T) {
	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A"}),
	}
	table := widgets.NewTable(nil, []layout.Constraint{layout.Length(1)}).Rows(rows)
	rows[0] = widgets.TableRowFromStrings([]string{"B"})

	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{"A  "})
}

func TestTable_stylize(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A"}),
	}, []layout.Constraint{layout.Length(1)}).
		Fg(style.Red).
		Bg(style.Blue).
		Bold().
		Dim().
		Italic().
		Cyan()

	table.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 0, style.NewStyle().
		Fg(style.Cyan).
		Bg(style.Blue).
		AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
}

func TestTable_listStateEmptyList(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 10))
	state := widgets.NewTableState()
	state.SelectFirst()
	state.SelectColumn(0)
	state.SelectCell(0, 0)
	state.SetOffset(4)

	widgets.NewTable(nil, []layout.Constraint{layout.Percentage(100)}).RenderStateful(buf.Area, buf, &state)

	if selected, ok := state.Selected(); ok {
		t.Fatalf("Selected() = %d, true, want false", selected)
	}
	if selectedColumn, ok := state.SelectedColumn(); ok {
		t.Fatalf("SelectedColumn() = %d, true, want false", selectedColumn)
	}
	if row, column, ok := state.SelectedCell(); ok {
		t.Fatalf("SelectedCell() = %d,%d,true, want false", row, column)
	}
	if got := state.Offset(); got != 0 {
		t.Fatalf("Offset() = %d, want 0", got)
	}
}

func TestTable_listStateSingleItem(t *testing.T) {
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Item 1"}),
	}, []layout.Constraint{layout.Percentage(100)}).
		RowHighlightStyle(style.NewStyle().Fg(style.Red))
	state := widgets.NewTableState()

	state.SelectFirst()
	table.RenderStateful(layout.NewRect(0, 0, 10, 10), buffer.Empty(layout.NewRect(0, 0, 10, 10)), &state)
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("Selected() after first = %d, %v; want 0, true", selected, ok)
	}
	if got := state.Offset(); got != 0 {
		t.Fatalf("Offset() after first = %d, want 0", got)
	}

	state.SelectLast()
	table.RenderStateful(layout.NewRect(0, 0, 10, 10), buffer.Empty(layout.NewRect(0, 0, 10, 10)), &state)
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("Selected() after last = %d, %v; want 0, true", selected, ok)
	}

	state.SelectPrevious()
	table.RenderStateful(layout.NewRect(0, 0, 10, 10), buffer.Empty(layout.NewRect(0, 0, 10, 10)), &state)
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("Selected() after previous = %d, %v; want 0, true", selected, ok)
	}

	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	state.SelectNext()
	table.RenderStateful(buf.Area, buf, &state)
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("Selected() after next = %d, %v; want 0, true", selected, ok)
	}
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red))
}

func TestTable_renderEmptyArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 3))
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Cell1", "Cell2"}),
	}, []layout.Constraint{layout.Length(5), layout.Length(5)})

	table.Render(layout.NewRect(0, 0, 0, 0), buf)

	assertLines(t, buf, []string{
		"               ",
		"               ",
		"               ",
	})
}

func TestTable_renderInMinimalBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Cell1", "Cell2", "Cell3"}),
		widgets.TableRowFromStrings([]string{"Cell4", "Cell5", "Cell6"}),
	}, []layout.Constraint{layout.Length(10), layout.Length(10), layout.Length(10)}).
		Header(widgets.TableRowFromStrings([]string{"Header1", "Header2", "Header3"})).
		Footer(widgets.TableRowFromStrings([]string{"Footer1", "Footer2", "Footer3"}))

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{" "})
}

func TestTable_renderDefault(t *testing.T) {
	buf := buffer.WithLines([]string{"abc", "def"})

	widgets.NewTable(nil, nil).Render(buf.Area, buf)

	assertLines(t, buf, []string{"abc", "def"})
	assertAllCellsStyle(t, buf, style.NewStyle())
}

func TestTable_renderInZeroSizeBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Cell1", "Cell2", "Cell3"}),
		widgets.TableRowFromStrings([]string{"Cell4", "Cell5", "Cell6"}),
	}, []layout.Constraint{layout.Length(10), layout.Length(10), layout.Length(10)}).
		Header(widgets.TableRowFromStrings([]string{"Header1", "Header2", "Header3"})).
		Footer(widgets.TableRowFromStrings([]string{"Footer1", "Footer2", "Footer3"}))

	assertNotPanics(t, func() {
		table.Render(buf.Area, buf)
	})
}

func TestTable_shouldPatchRowAndCellStyles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 1))
	row := widgets.NewTableRow([]widgets.TableCell{
		widgets.TableCellFromString("A"),
		widgets.NewTableCell(text.FromString("B")).Style(style.NewStyle().Fg(style.Cyan)),
	}).Style(style.NewStyle().Bg(style.Yellow))
	table := widgets.NewTable([]widgets.TableRow{row}, []layout.Constraint{
		layout.Length(1),
		layout.Length(1),
	}).ColumnSpacing(1)

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{"A B     "})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Bg(style.Yellow))
	assertCellStyle(t, buf, 2, 0, style.NewStyle().Fg(style.Cyan).Bg(style.Yellow))
}

func TestTable_shouldPatchCellLineAndSpanStylesInOrder(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	state := widgets.NewTableState()
	state.Select(0)
	state.SelectCell(0, 0)
	line := text.NewLine(
		text.StyledSpan("A", style.NewStyle().Fg(style.Red)),
	).Style(style.NewStyle().Fg(style.Yellow))
	row := widgets.NewTableRow([]widgets.TableCell{
		widgets.NewTableCell(text.NewText(line)).Style(style.NewStyle().Bg(style.Green)),
	}).Style(style.NewStyle().AddModifier(style.ModifierBold))
	table := widgets.NewTable([]widgets.TableRow{row}, []layout.Constraint{layout.Length(1)}).
		Style(style.NewStyle().Bg(style.Blue)).
		RowHighlightStyle(style.NewStyle().AddModifier(style.ModifierItalic)).
		CellHighlightStyle(style.NewStyle().Fg(style.Cyan).AddModifier(style.ModifierDim))

	table.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{"A   "})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red).Bg(style.Green).AddModifier(style.ModifierBold|style.ModifierItalic|style.ModifierDim))
}

func TestTable_renderWithAlignment(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 28, 4))
	rows := []widgets.TableRow{
		widgets.NewTableRow([]widgets.TableCell{
			widgets.NewTableCell(text.FromString("left")),
			widgets.NewTableCell(text.FromString("text").Center()),
			widgets.NewTableCell(text.NewText(text.LineFromString("line").Right()).Center()),
		}),
		widgets.NewTableRow([]widgets.TableCell{
			widgets.NewTableCell(text.NewText(text.LineFromString("wide line").Right())),
			widgets.NewTableCell(text.FromString("wide text").Center()),
			widgets.TableCellFromString("default"),
		}),
	}
	table := widgets.NewTable(rows, []layout.Constraint{
		layout.Length(8),
		layout.Length(8),
		layout.Length(8),
	}).ColumnSpacing(1)

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"left       text       line  ",
		"ide line wide tex default   ",
		"                            ",
		"                            ",
	})
}

func TestTable_shouldClearHiddenCellWhenWideGraphemeFits(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	table := widgets.NewTable([]widgets.TableRow{
		widgets.NewTableRow([]widgets.TableCell{
			widgets.NewTableCell(text.FromString("コ")),
		}),
	}, []layout.Constraint{layout.Length(2)})

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{"コ "})
	assertCellSymbol(t, buf, 0, 0, "コ")
	assertCellSymbol(t, buf, 1, 0, " ")
}

func TestTable_shouldRenderElementsStyledIndividually(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 25, 4))
	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Row11", "Row12", "Row13"}).
			Style(style.NewStyle().Fg(style.Green)),
		widgets.NewTableRow([]widgets.TableCell{
			widgets.TableCellFromString("Row21"),
			widgets.TableCellFromString("Row22").Style(style.NewStyle().Fg(style.Yellow)),
			widgets.NewTableCell(text.NewText(text.NewLine(
				text.NewSpan("Row"),
				text.StyledSpan("23", style.NewStyle().Fg(style.Blue)),
			))).Style(style.NewStyle().Fg(style.Red)),
		}).Style(style.NewStyle().Fg(style.LightGreen)),
	}
	state := widgets.NewTableState()
	state.Select(0)
	state.SelectColumn(1)
	table := widgets.NewTable(rows, []layout.Constraint{
		layout.Length(5),
		layout.Length(5),
		layout.Length(5),
	}).
		Block(widgets.NewBlock().Borders(widgets.LeftBorder | widgets.RightBorder)).
		HighlightSymbol(">> ").
		RowHighlightStyle(style.NewStyle().AddModifier(style.ModifierBold)).
		ColumnHighlightStyle(style.NewStyle().AddModifier(style.ModifierItalic)).
		CellHighlightStyle(style.NewStyle().AddModifier(style.ModifierDim)).
		ColumnSpacing(1)

	table.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"│>> Row11 Row12 Row13   │",
		"│   Row21 Row22 Row23   │",
		"│                       │",
		"│                       │",
	})
	assertCellStyle(t, buf, 4, 0, style.NewStyle().Fg(style.Green).AddModifier(style.ModifierBold))
	assertCellStyle(t, buf, 10, 0, style.NewStyle().Fg(style.Green).AddModifier(style.ModifierBold|style.ModifierItalic))
	assertCellStyle(t, buf, 10, 1, style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierItalic))
	assertCellStyle(t, buf, 19, 1, style.NewStyle().Fg(style.Blue))
	assertCellStyle(t, buf, 8, 1, style.NewStyle().Fg(style.LightGreen))
}

func TestTableState_ClearSelection_shouldResetOffset(t *testing.T) {
	state := widgets.NewTableState()
	state.Select(2)
	state.SetOffset(3)

	state.ClearSelection()

	if _, ok := state.Selected(); ok {
		t.Fatal("expected selection to be cleared")
	}
	if got := state.Offset(); got != 0 {
		t.Fatalf("offset = %d, want 0", got)
	}
}

func TestTableState_new(t *testing.T) {
	state := widgets.NewTableState()

	if got := state.Offset(); got != 0 {
		t.Fatalf("Offset() = %d, want 0", got)
	}
	if selected, ok := state.Selected(); ok {
		t.Fatalf("Selected() = %d, true, want false", selected)
	}
	if selectedColumn, ok := state.SelectedColumn(); ok {
		t.Fatalf("SelectedColumn() = %d, true, want false", selectedColumn)
	}
	if row, column, ok := state.SelectedCell(); ok {
		t.Fatalf("SelectedCell() = %d,%d,true, want false", row, column)
	}
}

func TestTableState_offset(t *testing.T) {
	state := widgets.NewTableState()
	state.SetOffset(3)

	if got := state.Offset(); got != 3 {
		t.Fatalf("Offset() = %d, want 3", got)
	}
}

func TestTableState_offsetMut(t *testing.T) {
	state := widgets.NewTableState()

	state.SetOffset(3)
	if got := state.Offset(); got != 3 {
		t.Fatalf("Offset() after SetOffset(3) = %d, want 3", got)
	}

	state.SetOffset(-1)
	if got := state.Offset(); got != 0 {
		t.Fatalf("Offset() after SetOffset(-1) = %d, want 0", got)
	}

	state = widgets.NewTableState().WithOffset(4)
	if got := state.Offset(); got != 4 {
		t.Fatalf("Offset() after WithOffset(4) = %d, want 4", got)
	}
}

func TestTableState_selected(t *testing.T) {
	state := widgets.NewTableState()

	if selected, ok := state.Selected(); ok {
		t.Fatalf("Selected() before Select = %d, true, want false", selected)
	}

	state.Select(2)
	if selected, ok := state.Selected(); !ok || selected != 2 {
		t.Fatalf("Selected() after Select(2) = %d, %v; want 2, true", selected, ok)
	}
}

func TestTableState_selectedCell(t *testing.T) {
	state := widgets.NewTableState()

	if row, column, ok := state.SelectedCell(); ok {
		t.Fatalf("SelectedCell() before SelectCell = %d,%d,true, want false", row, column)
	}

	state.SelectCell(2, 4)
	if row, column, ok := state.SelectedCell(); !ok || row != 2 || column != 4 {
		t.Fatalf("SelectedCell() after SelectCell(2,4) = %d,%d,%v; want 2,4,true", row, column, ok)
	}
	if selected, ok := state.Selected(); !ok || selected != 2 {
		t.Fatalf("Selected() after SelectCell(2,4) = %d, %v; want 2, true", selected, ok)
	}
	if selectedColumn, ok := state.SelectedColumn(); !ok || selectedColumn != 4 {
		t.Fatalf("SelectedColumn() after SelectCell(2,4) = %d, %v; want 4, true", selectedColumn, ok)
	}
}

func TestTableState_selectedColumn(t *testing.T) {
	state := widgets.NewTableState()

	if selectedColumn, ok := state.SelectedColumn(); ok {
		t.Fatalf("SelectedColumn() before SelectColumn = %d, true, want false", selectedColumn)
	}

	state.SelectColumn(3)
	if selectedColumn, ok := state.SelectedColumn(); !ok || selectedColumn != 3 {
		t.Fatalf("SelectedColumn() after SelectColumn(3) = %d, %v; want 3, true", selectedColumn, ok)
	}
}

func TestTableState_selectedColumnMut(t *testing.T) {
	state := widgets.NewTableState()

	state.SelectColumn(3)
	if selectedColumn, ok := state.SelectedColumn(); !ok || selectedColumn != 3 {
		t.Fatalf("SelectedColumn() after SelectColumn(3) = %d, %v; want 3, true", selectedColumn, ok)
	}

	state.ClearColumnSelection()
	if selectedColumn, ok := state.SelectedColumn(); ok {
		t.Fatalf("SelectedColumn() after ClearColumnSelection() = %d, true, want false", selectedColumn)
	}

	state = widgets.NewTableState().WithSelectedColumn(4)
	if selectedColumn, ok := state.SelectedColumn(); !ok || selectedColumn != 4 {
		t.Fatalf("SelectedColumn() after WithSelectedColumn(4) = %d, %v; want 4, true", selectedColumn, ok)
	}
}

func TestTableState_selectedMut(t *testing.T) {
	state := widgets.NewTableState()

	state.Select(2)
	if selected, ok := state.Selected(); !ok || selected != 2 {
		t.Fatalf("Selected() after Select(2) = %d, %v; want 2, true", selected, ok)
	}

	state.ClearSelection()
	if selected, ok := state.Selected(); ok {
		t.Fatalf("Selected() after ClearSelection() = %d, true, want false", selected)
	}

	state = widgets.NewTableState().WithSelected(3)
	if selected, ok := state.Selected(); !ok || selected != 3 {
		t.Fatalf("Selected() after WithSelected(3) = %d, %v; want 3, true", selected, ok)
	}
}

func TestTableState_shouldSupportFluentSetters(t *testing.T) {
	state := widgets.NewTableState().
		WithOffset(4).
		WithSelected(2).
		WithSelectedColumn(1).
		WithSelectedCell(3, 5)

	if got := state.Offset(); got != 4 {
		t.Fatalf("offset = %d, want 4", got)
	}
	if selected, ok := state.Selected(); !ok || selected != 3 {
		t.Fatalf("selected = %d, %v; want 3, true", selected, ok)
	}
	if selectedColumn, ok := state.SelectedColumn(); !ok || selectedColumn != 5 {
		t.Fatalf("selected column = %d, %v; want 5, true", selectedColumn, ok)
	}
	if row, column, ok := state.SelectedCell(); !ok || row != 3 || column != 5 {
		t.Fatalf("selected cell = %d,%d,%v; want 3,5,true", row, column, ok)
	}
}

func TestTableState_SelectCell_shouldSynchronizeRowAndColumnSelection(t *testing.T) {
	state := widgets.NewTableState()

	state.SelectCell(2, 4)

	if selected, ok := state.Selected(); !ok || selected != 2 {
		t.Fatalf("selected = %d, %v; want 2, true", selected, ok)
	}
	if selectedColumn, ok := state.SelectedColumn(); !ok || selectedColumn != 4 {
		t.Fatalf("selected column = %d, %v; want 4, true", selectedColumn, ok)
	}
	if row, column, ok := state.SelectedCell(); !ok || row != 2 || column != 4 {
		t.Fatalf("selected cell = %d,%d,%v; want 2,4,true", row, column, ok)
	}
}

func TestTableState_ClearCellSelection_shouldClearSelectionsAndResetOffset(t *testing.T) {
	state := widgets.NewTableState().
		WithOffset(3).
		WithSelectedCell(1, 2)

	state.ClearCellSelection()

	if _, ok := state.Selected(); ok {
		t.Fatal("expected row selection to be cleared")
	}
	if _, ok := state.SelectedColumn(); ok {
		t.Fatal("expected column selection to be cleared")
	}
	if _, _, ok := state.SelectedCell(); ok {
		t.Fatal("expected cell selection to be cleared")
	}
	if got := state.Offset(); got != 0 {
		t.Fatalf("offset = %d, want 0", got)
	}
}

func TestTableState_shouldNavigateRows(t *testing.T) {
	state := widgets.NewTableState()

	state.SelectFirst()
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("selected after first = %d, %v; want 0, true", selected, ok)
	}

	state.SelectPrevious()
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("selected after previous at start = %d, %v; want 0, true", selected, ok)
	}

	state.SelectNext()
	if selected, ok := state.Selected(); !ok || selected != 1 {
		t.Fatalf("selected after next = %d, %v; want 1, true", selected, ok)
	}

	state.SelectPrevious()
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("selected after previous = %d, %v; want 0, true", selected, ok)
	}

	state.SelectLast()
	if selected, ok := state.Selected(); !ok || selected <= 1 {
		t.Fatalf("selected after last = %d, %v; want sentinel > 1, true", selected, ok)
	}
}

func TestTableState_shouldScrollRows(t *testing.T) {
	state := widgets.NewTableState()

	state.ScrollDownBy(3)
	if selected, ok := state.Selected(); !ok || selected != 3 {
		t.Fatalf("selected after scroll down from empty = %d, %v; want 3, true", selected, ok)
	}

	state.ScrollUpBy(10)
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("selected after scroll up past start = %d, %v; want 0, true", selected, ok)
	}

	state.ScrollDownBy(-1)
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("selected after negative scroll down = %d, %v; want 0, true", selected, ok)
	}
}

func TestTableState_shouldNavigateColumns(t *testing.T) {
	state := widgets.NewTableState()

	state.SelectFirstColumn()
	state.SelectPreviousColumn()
	if selected, ok := state.SelectedColumn(); !ok || selected != 0 {
		t.Fatalf("selected column after previous at start = %d, %v; want 0, true", selected, ok)
	}

	state.SelectNextColumn()
	if selected, ok := state.SelectedColumn(); !ok || selected != 1 {
		t.Fatalf("selected column after next = %d, %v; want 1, true", selected, ok)
	}

	state.SelectPreviousColumn()
	if selected, ok := state.SelectedColumn(); !ok || selected != 0 {
		t.Fatalf("selected column after previous = %d, %v; want 0, true", selected, ok)
	}

	state.SelectLastColumn()
	if selected, ok := state.SelectedColumn(); !ok || selected <= 1 {
		t.Fatalf("selected column after last = %d, %v; want sentinel > 1, true", selected, ok)
	}
}

func TestTableState_shouldScrollColumns(t *testing.T) {
	state := widgets.NewTableState()

	state.ScrollRightBy(3)
	if selected, ok := state.SelectedColumn(); !ok || selected != 3 {
		t.Fatalf("selected column after scroll right from empty = %d, %v; want 3, true", selected, ok)
	}

	state.ScrollLeftBy(10)
	if selected, ok := state.SelectedColumn(); !ok || selected != 0 {
		t.Fatalf("selected column after scroll left past start = %d, %v; want 0, true", selected, ok)
	}

	state.ScrollRightBy(-1)
	if selected, ok := state.SelectedColumn(); !ok || selected != 0 {
		t.Fatalf("selected column after negative scroll right = %d, %v; want 0, true", selected, ok)
	}
}

func TestTable_shouldRenderMultilineRowsWithSelection(t *testing.T) {
	tests := []struct {
		name     string
		selected *int
		expected []string
	}{
		{
			name: "none",
			expected: []string{
				"┌────────────────────────────┐",
				"│Head1 Head2 Head3           │",
				"│                            │",
				"│Row11 Row12 Row13           │",
				"│Row21 Row22 Row23           │",
				"│                            │",
				"│Row31 Row32 Row33           │",
				"└────────────────────────────┘",
			},
		},
		{
			name:     "first",
			selected: new(0),
			expected: []string{
				"┌────────────────────────────┐",
				"│   Head1 Head2 Head3        │",
				"│                            │",
				"│>> Row11 Row12 Row13        │",
				"│   Row21 Row22 Row23        │",
				"│                            │",
				"│   Row31 Row32 Row33        │",
				"└────────────────────────────┘",
			},
		},
		{
			name:     "second",
			selected: new(1),
			expected: []string{
				"┌────────────────────────────┐",
				"│   Head1 Head2 Head3        │",
				"│                            │",
				"│   Row11 Row12 Row13        │",
				"│>> Row21 Row22 Row23        │",
				"│                            │",
				"│   Row31 Row32 Row33        │",
				"└────────────────────────────┘",
			},
		},
		{
			name:     "fourth",
			selected: new(3),
			expected: []string{
				"┌────────────────────────────┐",
				"│   Head1 Head2 Head3        │",
				"│                            │",
				"│   Row31 Row32 Row33        │",
				"│>> Row41 Row42 Row43        │",
				"│                            │",
				"│                            │",
				"└────────────────────────────┘",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 30, 8))
			state := widgets.NewTableState()
			if tt.selected != nil {
				state.Select(*tt.selected)
			}

			multilineTable().RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestTable_shouldRespectHighlightSpacing(t *testing.T) {
	tests := []struct {
		name     string
		selected *int
		spacing  widgets.HighlightSpacing
		expected []string
	}{
		{name: "none when selected", spacing: widgets.HighlightSpacingWhenSelected, expected: []string{
			"┌────────────────────────────┐",
			"│Head1 Head2 Head3           │",
			"│                            │",
			"│Row11 Row12 Row13           │",
			"│Row21 Row22 Row23           │",
			"│                            │",
			"│Row31 Row32 Row33           │",
			"└────────────────────────────┘",
		}},
		{name: "none always", spacing: widgets.HighlightSpacingAlways, expected: []string{
			"┌────────────────────────────┐",
			"│   Head1 Head2 Head3        │",
			"│                            │",
			"│   Row11 Row12 Row13        │",
			"│   Row21 Row22 Row23        │",
			"│                            │",
			"│   Row31 Row32 Row33        │",
			"└────────────────────────────┘",
		}},
		{name: "none never", spacing: widgets.HighlightSpacingNever, expected: []string{
			"┌────────────────────────────┐",
			"│Head1 Head2 Head3           │",
			"│                            │",
			"│Row11 Row12 Row13           │",
			"│Row21 Row22 Row23           │",
			"│                            │",
			"│Row31 Row32 Row33           │",
			"└────────────────────────────┘",
		}},
		{name: "first when selected", selected: new(0), spacing: widgets.HighlightSpacingWhenSelected, expected: []string{
			"┌────────────────────────────┐",
			"│   Head1 Head2 Head3        │",
			"│                            │",
			"│>> Row11 Row12 Row13        │",
			"│   Row21 Row22 Row23        │",
			"│                            │",
			"│   Row31 Row32 Row33        │",
			"└────────────────────────────┘",
		}},
		{name: "first always", selected: new(0), spacing: widgets.HighlightSpacingAlways, expected: []string{
			"┌────────────────────────────┐",
			"│   Head1 Head2 Head3        │",
			"│                            │",
			"│>> Row11 Row12 Row13        │",
			"│   Row21 Row22 Row23        │",
			"│                            │",
			"│   Row31 Row32 Row33        │",
			"└────────────────────────────┘",
		}},
		{name: "first never", selected: new(0), spacing: widgets.HighlightSpacingNever, expected: []string{
			"┌────────────────────────────┐",
			"│Head1 Head2 Head3           │",
			"│                            │",
			"│Row11 Row12 Row13           │",
			"│Row21 Row22 Row23           │",
			"│                            │",
			"│Row31 Row32 Row33           │",
			"└────────────────────────────┘",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 30, 8))
			state := widgets.NewTableState()
			if tt.selected != nil {
				state.Select(*tt.selected)
			}

			multilineTable().HighlightSpacing(tt.spacing).RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestTable_shouldClampOffsetWhenRowsAreRemoved(t *testing.T) {
	state := widgets.NewTableState()
	state.Select(5)
	state.SetOffset(5)
	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Row1"}),
		widgets.TableRowFromStrings([]string{"Row2"}),
		widgets.TableRowFromStrings([]string{"Row3"}),
		widgets.TableRowFromStrings([]string{"Row4"}),
		widgets.TableRowFromStrings([]string{"Row5"}),
		widgets.TableRowFromStrings([]string{"Row6"}),
	}
	widgets.NewTable(rows, []layout.Constraint{layout.Length(4)}).RenderStateful(layout.NewRect(0, 0, 6, 2), buffer.Empty(layout.NewRect(0, 0, 6, 2)), &state)

	widgets.NewTable(rows[:1], []layout.Constraint{layout.Length(4)}).RenderStateful(layout.NewRect(0, 0, 6, 2), buffer.Empty(layout.NewRect(0, 0, 6, 2)), &state)

	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("selected = %d, %v; want 0, true", selected, ok)
	}
	if got := state.Offset(); got != 0 {
		t.Fatalf("offset = %d, want 0", got)
	}
}

func TestTable_shouldNotPanicWithSelectedFirstRowAndPercentageColumns(t *testing.T) {
	state := widgets.NewTableState()
	state.Select(0)
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Row11", "Row12", "Row13"}),
	}, []layout.Constraint{
		layout.Percentage(50),
		layout.Percentage(50),
		layout.Percentage(50),
	}).HighlightSymbol(">> ")

	assertNotPanics(t, func() {
		table.RenderStateful(layout.NewRect(0, 0, 3, 1), buffer.Empty(layout.NewRect(0, 0, 3, 1)), &state)
	})
}

func TestTable_RenderStateful_shouldHandleNilState(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 1))

	widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Row11"}),
	}, []layout.Constraint{layout.Length(5)}).RenderStateful(buf.Area, buf, nil)

	assertLines(t, buf, []string{"Row11       "})
}

func TestTable_RenderStateful_shouldClampSelectedIndex(t *testing.T) {
	tests := []struct {
		name     string
		selected int
		want     int
	}{
		{name: "negative", selected: -10, want: 0},
		{name: "past end", selected: 99, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := widgets.NewTableState()
			state.Select(tt.selected)

			widgets.NewTable([]widgets.TableRow{
				widgets.TableRowFromStrings([]string{"Row1"}),
				widgets.TableRowFromStrings([]string{"Row2"}),
			}, []layout.Constraint{layout.Length(4)}).RenderStateful(layout.NewRect(0, 0, 6, 2), buffer.Empty(layout.NewRect(0, 0, 6, 2)), &state)

			if selected, ok := state.Selected(); !ok || selected != tt.want {
				t.Fatalf("selected = %d, %v; want %d, true", selected, ok, tt.want)
			}
		})
	}
}

func TestTable_RenderStateful_shouldClampSelectedColumnAndCell(t *testing.T) {
	state := widgets.NewTableState()
	state.SelectColumn(99)
	state.SelectCell(99, 99)

	widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A", "B"}),
		widgets.TableRowFromStrings([]string{"C", "D"}),
	}, []layout.Constraint{layout.Length(1), layout.Length(1)}).
		RenderStateful(layout.NewRect(0, 0, 4, 2), buffer.Empty(layout.NewRect(0, 0, 4, 2)), &state)

	if selected, ok := state.SelectedColumn(); !ok || selected != 1 {
		t.Fatalf("selected column = %d, %v; want 1, true", selected, ok)
	}
	if row, column, ok := state.SelectedCell(); !ok || row != 1 || column != 1 {
		t.Fatalf("selected cell = %d,%d,%v; want 1,1,true", row, column, ok)
	}
}

func TestTable_RenderStateful_shouldClampSelectLastToFinalRow(t *testing.T) {
	state := widgets.NewTableState()
	state.SelectLast()

	widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A"}),
		widgets.TableRowFromStrings([]string{"B"}),
		widgets.TableRowFromStrings([]string{"C"}),
	}, []layout.Constraint{layout.Length(1)}).
		RenderStateful(layout.NewRect(0, 0, 1, 3), buffer.Empty(layout.NewRect(0, 0, 1, 3)), &state)

	if selected, ok := state.Selected(); !ok || selected != 2 {
		t.Fatalf("selected = %d, %v; want 2, true", selected, ok)
	}
}

func TestTable_RenderStateful_shouldClampSelectLastColumnToFinalColumn(t *testing.T) {
	state := widgets.NewTableState()
	state.SelectLastColumn()

	widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A", "B", "C"}),
	}, []layout.Constraint{layout.Length(1), layout.Length(1), layout.Length(1)}).
		RenderStateful(layout.NewRect(0, 0, 5, 1), buffer.Empty(layout.NewRect(0, 0, 5, 1)), &state)

	if selectedColumn, ok := state.SelectedColumn(); !ok || selectedColumn != 2 {
		t.Fatalf("selected column = %d, %v; want 2, true", selectedColumn, ok)
	}
}

func TestTable_RenderStateful_WithSelectedCell_shouldHighlightSameCellAsSelectCell(t *testing.T) {
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A", "B"}),
	}, []layout.Constraint{layout.Length(1), layout.Length(1)}).
		CellHighlightStyle(style.NewStyle().Fg(style.Red))

	selectCellBuf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	selectCellState := widgets.NewTableState()
	selectCellState.SelectCell(0, 1)
	table.RenderStateful(selectCellBuf.Area, selectCellBuf, &selectCellState)

	withSelectedCellBuf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	withSelectedCellState := widgets.NewTableState().WithSelectedCell(0, 1)
	table.RenderStateful(withSelectedCellBuf.Area, withSelectedCellBuf, &withSelectedCellState)

	cell, ok := selectCellBuf.CellAt(2, 0)
	if !ok {
		t.Fatal("expected selected cell")
	}
	assertCellStyle(t, withSelectedCellBuf, 2, 0, cell.Style)
}

func TestTable_ClearColumnAndCellSelection_shouldRemoveHighlights(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	state := widgets.NewTableState()
	state.SelectColumn(0)
	state.SelectCell(0, 1)
	state.ClearColumnSelection()
	state.ClearCellSelection()
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A", "B"}),
	}, []layout.Constraint{layout.Length(1), layout.Length(1)}).
		ColumnHighlightStyle(style.NewStyle().Fg(style.Red)).
		CellHighlightStyle(style.NewStyle().Fg(style.Blue))

	table.RenderStateful(buf.Area, buf, &state)

	assertCellStyle(t, buf, 0, 0, style.NewStyle())
	assertCellStyle(t, buf, 2, 0, style.NewStyle())
}

func TestTable_shouldPatchCellHighlightOverRowAndColumnHighlight(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	state := widgets.NewTableState()
	state.Select(0)
	state.SelectColumn(1)
	state.SelectCell(0, 1)
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"A", "B"}),
	}, []layout.Constraint{layout.Length(1), layout.Length(1)}).
		RowHighlightStyle(style.NewStyle().Fg(style.Green).AddModifier(style.ModifierBold)).
		ColumnHighlightStyle(style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierItalic)).
		CellHighlightStyle(style.NewStyle().Fg(style.Red).AddModifier(style.ModifierDim))

	table.RenderStateful(buf.Area, buf, &state)

	assertCellStyle(t, buf, 2, 0, style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold|style.ModifierItalic|style.ModifierDim))
}

func TestTableCell_ColumnSpan_shouldRenderAcrossPhysicalColumns(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		widths   []layout.Constraint
		rows     []widgets.TableRow
		expected []string
	}{
		{
			name:  "zero span skips cell",
			width: 15,
			widths: []layout.Constraint{
				layout.Length(5),
				layout.Length(5),
			},
			rows: []widgets.TableRow{
				widgets.NewTableRow([]widgets.TableCell{
					widgets.TableCellFromString("Cell1").ColumnSpan(0),
					widgets.TableCellFromString("Cell2"),
				}),
				widgets.TableRowFromStrings([]string{"Cell3", "Cell4"}),
			},
			expected: []string{
				"Cell2          ",
				"Cell3 Cell4    ",
			},
		},
		{
			name:  "two column span includes spacing",
			width: 15,
			widths: []layout.Constraint{
				layout.Length(5),
				layout.Length(5),
			},
			rows: []widgets.TableRow{
				widgets.NewTableRow([]widgets.TableCell{
					widgets.TableCellFromString("Cell1").ColumnSpan(2),
					widgets.TableCellFromString("Cell2"),
				}),
				widgets.TableRowFromStrings([]string{"Cell3", "Cell4"}),
			},
			expected: []string{
				"Cell1          ",
				"Cell3 Cell4    ",
			},
		},
		{
			name:  "first cell spans first two of three columns",
			width: 17,
			widths: []layout.Constraint{
				layout.Length(5),
				layout.Length(5),
				layout.Length(5),
			},
			rows: []widgets.TableRow{
				widgets.NewTableRow([]widgets.TableCell{
					widgets.TableCellFromString("Cell1").ColumnSpan(2),
					widgets.TableCellFromString("Cell2"),
				}),
				widgets.TableRowFromStrings([]string{"Cell3", "Cell4", "Cell5"}),
			},
			expected: []string{
				"Cell1       Cell2",
				"Cell3 Cell4 Cell5",
			},
		},
		{
			name:  "middle cell spans remaining columns",
			width: 17,
			widths: []layout.Constraint{
				layout.Length(5),
				layout.Length(5),
				layout.Length(5),
			},
			rows: []widgets.TableRow{
				widgets.NewTableRow([]widgets.TableCell{
					widgets.TableCellFromString("Cell1"),
					widgets.TableCellFromString("Cell2").ColumnSpan(2),
					widgets.TableCellFromString("Cell3"),
				}),
				widgets.TableRowFromStrings([]string{"Cell4", "Cell5", "Cell6"}),
			},
			expected: []string{
				"Cell1 Cell2      ",
				"Cell4 Cell5 Cell6",
			},
		},
		{
			name:  "long text truncates to spanned width",
			width: 15,
			widths: []layout.Constraint{
				layout.Length(5),
				layout.Length(5),
				layout.Length(5),
			},
			rows: []widgets.TableRow{
				widgets.NewTableRow([]widgets.TableCell{
					widgets.TableCellFromString("11111111111111111111").ColumnSpan(2),
					widgets.TableCellFromString("22222222222222222222"),
				}),
				widgets.NewTableRow([]widgets.TableCell{
					widgets.TableCellFromString("33333333333333333333"),
					widgets.TableCellFromString("44444444444444444444").ColumnSpan(2),
					widgets.TableCellFromString("55555555555555555555"),
				}),
			},
			expected: []string{
				"111111111 22222",
				"3333 4444444444",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, tt.width, 2))

			widgets.NewTable(tt.rows, tt.widths).Render(buf.Area, buf)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestTableCell_content(t *testing.T) {
	content := text.FromString("cell\ncontent").Cyan()
	cell := widgets.NewTableCell(content)

	got := cell.Content()
	if got.String() != content.String() {
		t.Fatalf("Content().String() = %q, want %q", got.String(), content.String())
	}
	if got.Style != content.Style {
		t.Fatalf("Content().Style = %#v, want %#v", got.Style, content.Style)
	}
	if got.Alignment != nil || content.Alignment != nil {
		t.Fatalf("Content().Alignment = %#v, want %#v", got.Alignment, content.Alignment)
	}
	if len(got.Lines) != len(content.Lines) {
		t.Fatalf("len(Content().Lines) = %d, want %d", len(got.Lines), len(content.Lines))
	}
}

func TestTableCell_new(t *testing.T) {
	cell := widgets.NewTableCell(text.FromString("simple string"))
	row := widgets.NewTableRow([]widgets.TableCell{cell})
	table := widgets.NewTable([]widgets.TableRow{row}, []layout.Constraint{layout.Length(13)})
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 1))

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{"simple string  "})
}

func TestTableCell_stylize(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	cell := widgets.TableCellFromString("A").
		Fg(style.Red).
		Bg(style.Blue).
		Bold().
		Dim().
		Italic().
		Cyan()
	table := widgets.NewTable([]widgets.TableRow{
		widgets.NewTableRow([]widgets.TableCell{cell}),
	}, []layout.Constraint{layout.Length(1)})

	table.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 0, style.NewStyle().
		Fg(style.Cyan).
		Bg(style.Blue).
		AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
}

func TestTableRow_bottomMargin(t *testing.T) {
	row := widgets.TableRowFromStrings([]string{"row"}).BottomMargin(2)
	if got := row.BottomMarginValue(); got != 2 {
		t.Fatalf("BottomMarginValue() = %d, want 2", got)
	}

	clamped := widgets.TableRowFromStrings([]string{"row"}).BottomMargin(-1)
	if got := clamped.BottomMarginValue(); got != 0 {
		t.Fatalf("BottomMarginValue() after negative margin = %d, want 0", got)
	}

	buf := buffer.Empty(layout.NewRect(0, 0, 8, 4))
	table := widgets.NewTable([]widgets.TableRow{
		row,
		widgets.TableRowFromStrings([]string{"next"}),
	}, []layout.Constraint{layout.Length(4)})

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"row     ",
		"        ",
		"        ",
		"next    ",
	})
}

func TestTableRow_cells(t *testing.T) {
	row := widgets.NewTableRow([]widgets.TableCell{
		widgets.TableCellFromString("one"),
		widgets.TableCellFromString("two"),
	})

	got := row.Cells()
	if len(got) != 2 {
		t.Fatalf("len(Cells()) = %d, want 2", len(got))
	}
	if got[0].Content().String() != "one" || got[1].Content().String() != "two" {
		t.Fatalf("Cells() content = %q, %q; want one, two", got[0].Content().String(), got[1].Content().String())
	}

	got[0] = widgets.TableCellFromString("mutated")
	again := row.Cells()
	if again[0].Content().String() != "one" {
		t.Fatalf("Cells() returned internal slice; first cell = %q, want one", again[0].Content().String())
	}
}

func TestTableRow_height(t *testing.T) {
	row := widgets.TableRowFromStrings([]string{"row"}).Height(2)
	if got := row.HeightValue(); got != 2 {
		t.Fatalf("HeightValue() = %d, want 2", got)
	}

	clamped := widgets.TableRowFromStrings([]string{"row"}).Height(-1)
	if got := clamped.HeightValue(); got != 0 {
		t.Fatalf("HeightValue() after negative height = %d, want 0", got)
	}

	buf := buffer.Empty(layout.NewRect(0, 0, 8, 3))
	table := widgets.NewTable([]widgets.TableRow{row}, []layout.Constraint{layout.Length(4)})

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"row     ",
		"        ",
		"        ",
	})
}

func TestTableRow_new(t *testing.T) {
	cells := []widgets.TableCell{widgets.TableCellFromString("original")}
	row := widgets.NewTableRow(cells)
	cells[0] = widgets.TableCellFromString("mutated")

	if got := row.HeightValue(); got != 1 {
		t.Fatalf("HeightValue() = %d, want default 1", got)
	}
	if got := row.TopMarginValue(); got != 0 {
		t.Fatalf("TopMarginValue() = %d, want default 0", got)
	}
	if got := row.BottomMarginValue(); got != 0 {
		t.Fatalf("BottomMarginValue() = %d, want default 0", got)
	}
	if got := row.Cells()[0].Content().String(); got != "original" {
		t.Fatalf("NewTableRow copied cell content = %q, want original", got)
	}
}

func TestTableRow_stylize(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	row := widgets.NewTableRow([]widgets.TableCell{
		widgets.TableCellFromString("A"),
	}).Fg(style.Red).
		Bg(style.Blue).
		Bold().
		Dim().
		Italic().
		Cyan()
	table := widgets.NewTable([]widgets.TableRow{row}, []layout.Constraint{layout.Length(1)})

	table.Render(buf.Area, buf)

	assertCellStyle(t, buf, 0, 0, style.NewStyle().
		Fg(style.Cyan).
		Bg(style.Blue).
		AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
}

func TestTableRow_topMargin(t *testing.T) {
	row := widgets.TableRowFromStrings([]string{"row"}).TopMargin(2)
	if got := row.TopMarginValue(); got != 2 {
		t.Fatalf("TopMarginValue() = %d, want 2", got)
	}

	clamped := widgets.TableRowFromStrings([]string{"row"}).TopMargin(-1)
	if got := clamped.TopMarginValue(); got != 0 {
		t.Fatalf("TopMarginValue() after negative margin = %d, want 0", got)
	}

	buf := buffer.Empty(layout.NewRect(0, 0, 8, 3))
	table := widgets.NewTable([]widgets.TableRow{row}, []layout.Constraint{layout.Length(4)})

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"        ",
		"        ",
		"row     ",
	})
}

func TestTableCell_ColumnSpan_shouldRespectHighlightSymbolSpacing(t *testing.T) {
	tests := []struct {
		name     string
		spacing  widgets.HighlightSpacing
		selected bool
		expected []string
	}{
		{name: "always without selection", spacing: widgets.HighlightSpacingAlways, expected: []string{"   ABCDEFG 1234"}},
		{name: "always with selection", spacing: widgets.HighlightSpacingAlways, selected: true, expected: []string{">>>ABCDEFG 1234"}},
		{name: "when selected without selection", spacing: widgets.HighlightSpacingWhenSelected, expected: []string{"ABCDEFGHI 12345"}},
		{name: "when selected with selection", spacing: widgets.HighlightSpacingWhenSelected, selected: true, expected: []string{">>>ABCDEFG 1234"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 15, 1))
			state := widgets.NewTableState()
			if tt.selected {
				state.Select(0)
			}
			table := widgets.NewTable([]widgets.TableRow{
				widgets.NewTableRow([]widgets.TableCell{
					widgets.TableCellFromString("ABCDEFGHIJK").ColumnSpan(2),
					widgets.TableCellFromString("12345678901"),
					widgets.TableCellFromString("XYZXYZXYZXY"),
				}),
			}, []layout.Constraint{
				layout.Length(5),
				layout.Length(5),
				layout.Length(5),
			}).
				HighlightSpacing(tt.spacing).
				HighlightSymbol(">>>").
				ColumnSpacing(1)

			table.RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestTableCell_ColumnSpan_shouldApplyColumnAndCellHighlightsToSpannedCell(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 1))
	state := widgets.NewTableState()
	state.Select(0)
	state.SelectColumn(1)
	state.SelectCell(0, 1)
	table := widgets.NewTable([]widgets.TableRow{
		widgets.NewTableRow([]widgets.TableCell{
			widgets.TableCellFromString("A").ColumnSpan(2),
		}).Style(style.NewStyle().Fg(style.Cyan)),
	}, []layout.Constraint{
		layout.Length(5),
		layout.Length(5),
	}).
		RowHighlightStyle(style.NewStyle().Fg(style.Green).AddModifier(style.ModifierBold)).
		ColumnHighlightStyle(style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierItalic)).
		CellHighlightStyle(style.NewStyle().Fg(style.Red).AddModifier(style.ModifierDim))

	table.RenderStateful(buf.Area, buf, &state)

	expected := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold | style.ModifierItalic | style.ModifierDim)
	for x := range 11 {
		assertCellStyle(t, buf, x, 0, expected)
	}
}

func TestTable_shouldRenderFooter(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 3))
	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Cell1", "Cell2"}),
		widgets.TableRowFromStrings([]string{"Cell3", "Cell4"}),
	}
	table := widgets.NewTable(rows, []layout.Constraint{
		layout.Length(5),
		layout.Length(5),
	}).Footer(widgets.TableRowFromStrings([]string{"Foot1", "Foot2"}))

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"Cell1 Cell2    ",
		"Cell3 Cell4    ",
		"Foot1 Foot2    ",
	})
}

func TestTable_renderWithHeaderMargin(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 5))
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"row"}),
	}, []layout.Constraint{layout.Length(4)}).
		Style(style.NewStyle().Bg(style.Blue)).
		Header(widgets.TableRowFromStrings([]string{"head"}).TopMargin(1).BottomMargin(1))

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"            ",
		"head        ",
		"            ",
		"row         ",
		"            ",
	})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Bg(style.Blue))
	assertCellStyle(t, buf, 0, 2, style.NewStyle().Bg(style.Blue))
}

func TestTable_renderWithOverflowDoesNotPanic(t *testing.T) {
	rowWithMoreCells := widgets.TableRowFromStrings([]string{"a", "b", "c"})
	rowWithLargeMargin := widgets.TableRowFromStrings([]string{"margin"}).TopMargin(5).BottomMargin(5)
	rowWithLargeSpan := widgets.NewTableRow([]widgets.TableCell{
		widgets.TableCellFromString("span").ColumnSpan(10),
	})
	table := widgets.NewTable([]widgets.TableRow{
		rowWithMoreCells,
		widgets.NewTableRow(nil),
		rowWithLargeMargin,
		rowWithLargeSpan,
	}, []layout.Constraint{layout.Length(1)}).
		Header(widgets.TableRowFromStrings([]string{"h1", "h2"}).TopMargin(3).BottomMargin(3)).
		ColumnSpacing(4)

	assertNotPanics(t, func() {
		table.Render(layout.NewRect(0, 0, 2, 2), buffer.Empty(layout.NewRect(0, 0, 2, 2)))
		table.Render(layout.NewRect(0, 0, 0, 2), buffer.Empty(layout.NewRect(0, 0, 2, 2)))
		table.Render(layout.NewRect(0, 0, 2, 0), buffer.Empty(layout.NewRect(0, 0, 2, 2)))
	})
}

func TestTable_renderWithRowMargin(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 6))
	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"one"}).TopMargin(1).BottomMargin(1),
		widgets.TableRowFromStrings([]string{"two"}),
	}
	table := widgets.NewTable(rows, []layout.Constraint{layout.Length(4)}).
		Style(style.NewStyle().Bg(style.Blue))

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"            ",
		"one         ",
		"            ",
		"two         ",
		"            ",
		"            ",
	})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Bg(style.Blue))
	assertCellStyle(t, buf, 0, 2, style.NewStyle().Bg(style.Blue))
}

func TestTable_renderWithSelectedColumnAndIncorrectWidthCountDoesNotPanic(t *testing.T) {
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"a", "b", "c"}),
		widgets.TableRowFromStrings([]string{"d"}),
	}, []layout.Constraint{layout.Length(1)}).
		ColumnHighlightStyle(style.NewStyle().Bg(style.Red)).
		CellHighlightStyle(style.NewStyle().Bg(style.Blue))
	state := widgets.NewTableState()
	state.Select(0)
	state.SelectColumn(2)
	state.SelectCell(0, 2)

	buf := buffer.Empty(layout.NewRect(0, 0, 4, 2))
	assertNotPanics(t, func() {
		table.RenderStateful(buf.Area, buf, &state)
	})
	assertLines(t, buf, []string{
		"a   ",
		"d   ",
	})
	assertCellStyle(t, buf, 0, 0, style.NewStyle())
}

func TestTable_shouldRenderFooterWithTopMargin(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 3))
	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Cell1", "Cell2"}),
	}
	table := widgets.NewTable(rows, []layout.Constraint{
		layout.Length(5),
		layout.Length(5),
	}).Footer(widgets.TableRowFromStrings([]string{"Foot1", "Foot2"}).TopMargin(1))

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"Cell1 Cell2    ",
		"               ",
		"Foot1 Foot2    ",
	})
}

func TestTable_shouldRenderFooterAtBottomWhenBodyIsShort(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 5))
	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Cell1", "Cell2"}),
	}
	table := widgets.NewTable(rows, []layout.Constraint{
		layout.Length(5),
		layout.Length(5),
	}).Footer(widgets.TableRowFromStrings([]string{"Foot1", "Foot2"}))

	table.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"Cell1 Cell2    ",
		"               ",
		"               ",
		"               ",
		"Foot1 Foot2    ",
	})
}

func TestTable_shouldRenderHeaderAndFooterOnEmptyTable(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 30, 6))
	table := widgets.NewTable(nil, []layout.Constraint{
		layout.Length(6),
		layout.Length(6),
		layout.Length(6),
	}).
		Header(widgets.TableRowFromStrings([]string{"Head1", "Head2", "Head3"}).BottomMargin(1)).
		Footer(widgets.TableRowFromStrings([]string{"Foot1", "Foot2", "Foot3"}).TopMargin(1)).
		Block(widgets.BorderedBlock()).
		ColumnSpacing(1)

	table.RenderStateful(buf.Area, buf, &widgets.TableState{})

	assertLines(t, buf, []string{
		"┌────────────────────────────┐",
		"│Head1  Head2  Head3         │",
		"│                            │",
		"│                            │",
		"│Foot1  Foot2  Foot3         │",
		"└────────────────────────────┘",
	})
}

func assertTableWithSelection(t *testing.T, highlightSpacing widgets.HighlightSpacing, width, columnSpacing int, selected *int, expected []string) {
	t.Helper()
	buf := buffer.Empty(layout.NewRect(0, 0, width, 3))
	state := widgets.NewTableState()
	if selected != nil {
		state.Select(*selected)
	}
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"ABCDE", "12345"}),
	}, []layout.Constraint{
		layout.Length(5),
		layout.Length(5),
	}).
		HighlightSpacing(highlightSpacing).
		HighlightSymbol(">>>").
		ColumnSpacing(columnSpacing)

	table.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, expected)
}

func assertConstraintTable(t *testing.T, width int, selected *int, widths []layout.Constraint, expected []string) {
	t.Helper()
	buf := buffer.Empty(layout.NewRect(0, 0, width, len(expected)))
	state := widgets.NewTableState()
	if selected != nil {
		state.Select(*selected)
	}
	table := widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"ABCDE", "12345", "Z"}),
	}, widths).
		HighlightSpacing(widgets.HighlightSpacingWhenSelected).
		HighlightSymbol(">>>").
		ColumnSpacing(1)

	table.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, expected)
}

func tableFixture(widths []layout.Constraint) widgets.Table {
	return widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Row11", "Row12", "Row13"}),
		widgets.TableRowFromStrings([]string{"Row21", "Row22", "Row23"}),
		widgets.TableRowFromStrings([]string{"Row31", "Row32", "Row33"}),
		widgets.TableRowFromStrings([]string{"Row41", "Row42", "Row43"}),
	}, widths).
		Header(widgets.TableRowFromStrings([]string{"Head1", "Head2", "Head3"}).BottomMargin(1)).
		Block(widgets.BorderedBlock())
}

func multilineTable() widgets.Table {
	return widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Row11", "Row12", "Row13"}),
		widgets.TableRowFromStrings([]string{"Row21", "Row22", "Row23"}).Height(2),
		widgets.TableRowFromStrings([]string{"Row31", "Row32", "Row33"}),
		widgets.TableRowFromStrings([]string{"Row41", "Row42", "Row43"}).Height(2),
	}, []layout.Constraint{
		layout.Length(5),
		layout.Length(5),
		layout.Length(5),
	}).
		Header(widgets.TableRowFromStrings([]string{"Head1", "Head2", "Head3"}).BottomMargin(1)).
		Block(widgets.BorderedBlock()).
		HighlightSymbol(">> ").
		ColumnSpacing(1)
}
