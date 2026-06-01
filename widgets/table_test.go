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
