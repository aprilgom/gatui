package widgets_test

import (
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
			selected: intPtr(0),
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
			selected: intPtr(1),
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
			selected: intPtr(3),
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
		{name: "first when selected", selected: intPtr(0), spacing: widgets.HighlightSpacingWhenSelected, expected: []string{
			"┌────────────────────────────┐",
			"│   Head1 Head2 Head3        │",
			"│                            │",
			"│>> Row11 Row12 Row13        │",
			"│   Row21 Row22 Row23        │",
			"│                            │",
			"│   Row31 Row32 Row33        │",
			"└────────────────────────────┘",
		}},
		{name: "first always", selected: intPtr(0), spacing: widgets.HighlightSpacingAlways, expected: []string{
			"┌────────────────────────────┐",
			"│   Head1 Head2 Head3        │",
			"│                            │",
			"│>> Row11 Row12 Row13        │",
			"│   Row21 Row22 Row23        │",
			"│                            │",
			"│   Row31 Row32 Row33        │",
			"└────────────────────────────┘",
		}},
		{name: "first never", selected: intPtr(0), spacing: widgets.HighlightSpacingNever, expected: []string{
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
