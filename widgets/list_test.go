package widgets_test

import (
	"math"
	"slices"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestListState_selected(t *testing.T) {
	state := widgets.ListState{}
	if selected, ok := state.Selected(); ok || selected != 0 {
		t.Fatalf("initial Selected() = (%d, %v), want (0, false)", selected, ok)
	}

	state.Select(1)
	if selected, ok := state.Selected(); !ok || selected != 1 {
		t.Fatalf("after Select(1), Selected() = (%d, %v), want (1, true)", selected, ok)
	}

	state.ClearSelection()
	if selected, ok := state.Selected(); ok || selected != 0 {
		t.Fatalf("after ClearSelection(), Selected() = (%d, %v), want (0, false)", selected, ok)
	}
}

func renderListLines(list widgets.List, width, height int) []string {
	return renderList(list, width, height).Lines()
}

func renderList(list widgets.List, width, height int) *buffer.Buffer {
	buf := buffer.Empty(layout.NewRect(0, 0, width, height))
	list.Render(buf.Area, buf)
	return buf
}

func TestListState_select(t *testing.T) {
	state := widgets.ListState{}
	state.SetOffset(4)

	state.Select(2)
	if selected, ok := state.Selected(); !ok || selected != 2 {
		t.Fatalf("after Select(2), Selected() = (%d, %v), want (2, true)", selected, ok)
	}
	if got := state.Offset(); got != 4 {
		t.Fatalf("after Select(2), Offset() = %d, want 4", got)
	}

	state.ClearSelection()
	if _, ok := state.Selected(); ok {
		t.Fatal("after ClearSelection(), Selected() ok = true, want false")
	}
	if got := state.Offset(); got != 0 {
		t.Fatalf("after ClearSelection(), Offset() = %d, want 0", got)
	}
}

func TestListState_stateNavigation(t *testing.T) {
	state := widgets.ListState{}
	state.SelectFirst()
	assertListStateSelected(t, state, 0)

	state.SelectPrevious()
	assertListStateSelected(t, state, 0)

	state.SelectNext()
	assertListStateSelected(t, state, 1)

	state.SelectPrevious()
	assertListStateSelected(t, state, 0)

	state.SelectLast()
	assertListStateSelected(t, state, math.MaxInt)

	state.SelectNext()
	assertListStateSelected(t, state, math.MaxInt)

	state.SelectPrevious()
	assertListStateSelected(t, state, math.MaxInt-1)

	state.SelectNext()
	assertListStateSelected(t, state, math.MaxInt)

	state.ClearSelection()
	state.SelectNext()
	assertListStateSelected(t, state, 0)

	state.ClearSelection()
	state.SelectPrevious()
	assertListStateSelected(t, state, math.MaxInt)

	state = widgets.ListState{}
	state.Select(2)
	state.ScrollDownBy(4)
	assertListStateSelected(t, state, 6)

	state = widgets.ListState{}
	state.ScrollUpBy(3)
	assertListStateSelected(t, state, 0)

	state.Select(6)
	state.ScrollUpBy(4)
	assertListStateSelected(t, state, 2)

	state.ScrollUpBy(4)
	assertListStateSelected(t, state, 0)

	state.ScrollDownBy(-3)
	assertListStateSelected(t, state, 0)
}

func TestListState_withHelpers(t *testing.T) {
	state := widgets.NewListState().WithOffset(3).WithSelected(2)

	if got := state.Offset(); got != 3 {
		t.Fatalf("Offset() = %d, want 3", got)
	}
	assertListStateSelected(t, state, 2)

	state = state.WithoutSelected()
	if selected, ok := state.Selected(); ok || selected != 0 {
		t.Fatalf("WithoutSelected().Selected() = (%d, %v), want (0, false)", selected, ok)
	}
	if got := state.Offset(); got != 0 {
		t.Fatalf("WithoutSelected().Offset() = %d, want 0", got)
	}
}

func assertListStateSelected(t *testing.T, state widgets.ListState, want int) {
	t.Helper()
	selected, ok := state.Selected()
	if !ok || selected != want {
		t.Fatalf("Selected() = (%d, %v), want (%d, true)", selected, ok, want)
	}
}

func TestList_shouldShowLength(t *testing.T) {
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
	})

	if list.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", list.Len())
	}
	if list.IsEmpty() {
		t.Fatal("IsEmpty() = true, want false")
	}

	empty := widgets.NewList(nil)
	if empty.Len() != 0 {
		t.Fatalf("empty Len() = %d, want 0", empty.Len())
	}
	if !empty.IsEmpty() {
		t.Fatal("empty IsEmpty() = false, want true")
	}
}

func TestNewListFromStrings_shouldMatchExplicitListItems(t *testing.T) {
	actual := renderListLines(widgets.NewListFromStrings([]string{"one", "two"}), 6, 2)
	expected := renderListLines(widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("one"),
		widgets.ListItemFromString("two"),
	}), 6, 2)

	if !slices.Equal(actual, expected) {
		t.Fatalf("lines mismatch\nactual:   %#v\nexpected: %#v", actual, expected)
	}
}

func TestNewListFromLines_shouldMatchExplicitListItems(t *testing.T) {
	lines := []text.Line{
		text.LineFromString("one"),
		text.NewLine(text.StyledSpan("two", style.NewStyle().Fg(style.Red))),
	}
	actual := renderListLines(widgets.NewListFromLines(lines), 6, 2)
	expected := renderListLines(widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLine(lines[0]),
		widgets.ListItemFromLine(lines[1]),
	}), 6, 2)

	if !slices.Equal(actual, expected) {
		t.Fatalf("lines mismatch\nactual:   %#v\nexpected: %#v", actual, expected)
	}
}

func TestNewListFromText_shouldMatchExplicitListItems(t *testing.T) {
	items := []text.Text{
		text.FromString("one"),
		text.NewText(text.LineFromString("two"), text.LineFromString("three")),
	}
	actual := renderListLines(widgets.NewListFromText(items), 6, 3)
	expected := renderListLines(widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromText(items[0]),
		widgets.ListItemFromText(items[1]),
	}), 6, 3)

	if !slices.Equal(actual, expected) {
		t.Fatalf("lines mismatch\nactual:   %#v\nexpected: %#v", actual, expected)
	}
}

func TestNewListFromStrings_shouldNotAliasCallerSlice(t *testing.T) {
	items := []string{"one", "two"}
	list := widgets.NewListFromStrings(items)

	items[0] = "changed"

	assertLines(t, renderList(list, 8, 2), []string{
		"one     ",
		"two     ",
	})
}

func TestList_canBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 2))
	baseStyle := style.NewStyle().
		Fg(style.Black).
		Bg(style.White).
		AddModifier(style.ModifierBold | style.ModifierDim | style.ModifierItalic)
	itemStyle := style.NewStyle().
		Fg(style.Red).
		AddModifier(style.ModifierItalic)

	widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.NewLine(text.StyledSpan("ab", itemStyle))),
	}).
		Fg(style.Black).
		Bg(style.White).
		Bold().
		Dim().
		Italic().
		Cyan().
		Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"ab    ",
		"      ",
	})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red).Bg(style.White).AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).Bg(style.White).AddModifier(style.ModifierBold|style.ModifierDim|style.ModifierItalic))
	for x := 2; x < 6; x++ {
		assertCellStyle(t, buf, x, 0, baseStyle.Fg(style.Cyan))
	}
	for x := range 6 {
		assertCellStyle(t, buf, x, 1, baseStyle.Fg(style.Cyan))
	}
}

func TestListItem_style(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	itemStyle := style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierItalic)

	widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("ab").Style(itemStyle),
	}).Render(buf.Area, buf)

	assertLines(t, buf, []string{"ab  "})
	assertCellStyle(t, buf, 0, 0, itemStyle)
	assertCellStyle(t, buf, 1, 0, itemStyle)
	assertCellStyle(t, buf, 2, 0, itemStyle)
}

func TestListItem_canBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))
	item := widgets.ListItemFromString("x").
		Fg(style.Black).
		Bg(style.White).
		Bold().
		Dim().
		Italic().
		Cyan()

	widgets.NewList([]widgets.ListItem{item}).Render(buf.Area, buf)

	want := style.NewStyle().
		Fg(style.Cyan).
		Bg(style.White).
		AddModifier(style.ModifierBold | style.ModifierDim | style.ModifierItalic)
	assertLines(t, buf, []string{"x "})
	assertCellStyle(t, buf, 0, 0, want)
	assertCellStyle(t, buf, 1, 0, want)
}

func TestListItem_height(t *testing.T) {
	tests := []struct {
		name string
		item widgets.ListItem
		want int
	}{
		{name: "single", item: widgets.ListItemFromString("a"), want: 1},
		{name: "multi", item: widgets.ListItemFromString("a\nb"), want: 2},
		{name: "empty text", item: widgets.ListItemFromText(text.NewText()), want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Height(); got != tt.want {
				t.Fatalf("Height() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestListItem_width(t *testing.T) {
	tests := []struct {
		name string
		item widgets.ListItem
		want int
	}{
		{name: "single", item: widgets.ListItemFromString("12345"), want: 5},
		{name: "multi", item: widgets.ListItemFromString("12345\n1234567"), want: 7},
		{name: "wide grapheme", item: widgets.ListItemFromString("aコb"), want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Width(); got != tt.want {
				t.Fatalf("Width() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestListItem_constructors(t *testing.T) {
	tests := []struct {
		name string
		item widgets.ListItem
		want []string
	}{
		{name: "string", item: widgets.ListItemFromString("ab"), want: []string{"ab  "}},
		{name: "span", item: widgets.ListItemFromSpan(text.NewSpan("ab")), want: []string{"ab  "}},
		{name: "line", item: widgets.ListItemFromLine(text.LineFromString("ab")), want: []string{"ab  "}},
		{name: "spans", item: widgets.ListItemFromSpans(text.NewSpan("a"), text.NewSpan("b")), want: []string{"ab  "}},
		{name: "text", item: widgets.ListItemFromText(text.FromString("ab")), want: []string{"ab  "}},
		{name: "lines", item: widgets.ListItemFromLines(text.LineFromString("a"), text.LineFromString("b")), want: []string{"a   ", "b   "}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 4, len(tt.want)))

			widgets.NewList([]widgets.ListItem{tt.item}).Render(buf.Area, buf)

			assertLines(t, buf, tt.want)
		})
	}
}

func TestList_renderItemStylePatchesBeforeTextStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	listStyle := style.NewStyle().Fg(style.White).Bg(style.Black)
	itemStyle := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold)
	textStyle := style.NewStyle().Bg(style.Blue).AddModifier(style.ModifierDim)
	lineStyle := style.NewStyle().Fg(style.Green).AddModifier(style.ModifierItalic)
	spanStyle := style.NewStyle().Fg(style.Yellow)
	content := text.NewText(
		text.NewLine(text.StyledSpan("ab", spanStyle)).Style(lineStyle),
	).PatchStyle(textStyle)

	widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromText(content).Style(itemStyle),
	}).Style(listStyle).Render(buf.Area, buf)

	wantText := listStyle.Patch(itemStyle).Patch(textStyle).Patch(lineStyle).Patch(spanStyle)
	wantFill := listStyle.Patch(itemStyle).Patch(textStyle).Patch(lineStyle)
	assertLines(t, buf, []string{"ab   "})
	assertCellStyle(t, buf, 0, 0, wantText)
	assertCellStyle(t, buf, 1, 0, wantText)
	assertCellStyle(t, buf, 2, 0, wantFill)
}

func TestList_renderingCanBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 3))
	blockStyle := style.NewStyle().Fg(style.Blue)
	borderStyle := style.NewStyle().Fg(style.Green)
	titleStyle := style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierBold)
	listStyle := style.NewStyle().Bg(style.White)
	itemStyle := style.NewStyle().Fg(style.Red)
	highlightStyle := style.NewStyle().Bg(style.Cyan).AddModifier(style.ModifierItalic)
	state := widgets.ListState{}
	state.Select(0)

	widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.NewLine(text.StyledSpan("ab", itemStyle))),
	}).
		Style(listStyle).
		Block(widgets.BorderedBlock().
			Style(blockStyle).
			BorderStyle(borderStyle).
			Title(text.NewLine(text.NewSpan("T"))).
			TitleStyle(titleStyle)).
		HighlightStyle(highlightStyle).
		HighlightSymbol(">").
		RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"┌T────┐",
		"│>ab  │",
		"└─────┘",
	})
	assertCellStyle(t, buf, 0, 0, borderStyle.Patch(listStyle))
	assertCellStyle(t, buf, 1, 0, titleStyle.Patch(listStyle))
	assertCellStyle(t, buf, 2, 0, borderStyle.Patch(listStyle))
	assertCellStyle(t, buf, 0, 1, borderStyle.Patch(listStyle))
	assertCellStyle(t, buf, 1, 1, listStyle.Patch(highlightStyle))
	assertCellStyle(t, buf, 2, 1, listStyle.Patch(itemStyle).Patch(highlightStyle))
	assertCellStyle(t, buf, 3, 1, listStyle.Patch(itemStyle).Patch(highlightStyle))
	assertCellStyle(t, buf, 4, 1, listStyle.Patch(highlightStyle))
	assertCellStyle(t, buf, 6, 1, borderStyle.Patch(listStyle))
}

func TestList_renderingBlock(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 5))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("one"),
		widgets.ListItemFromString("two"),
		widgets.ListItemFromString("three"),
	}).Block(widgets.BorderedBlock().Title(text.LineFromString("L")))

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌L─────┐",
		"│one   │",
		"│two   │",
		"│three │",
		"└──────┘",
	})
}

func TestList_emptyList(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 2))
	state := widgets.ListState{}
	state.Select(3)
	state.SetOffset(5)

	widgets.NewList(nil).RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"      ",
		"      ",
	})
	if selected, ok := state.Selected(); ok {
		t.Fatalf("Selected() = %d, true, want false", selected)
	}
	if state.Offset() != 0 {
		t.Fatalf("Offset() = %d, want 0", state.Offset())
	}
}

func TestList_doesNotRenderInSmallSpace(t *testing.T) {
	tests := []struct {
		name     string
		area     layout.Rect
		block    bool
		expected []string
	}{
		{name: "one by one", area: layout.NewRect(0, 0, 1, 1), expected: []string{"i"}},
		{name: "one by two", area: layout.NewRect(0, 0, 1, 2), expected: []string{"i", "i"}},
		{name: "two by one", area: layout.NewRect(0, 0, 2, 1), expected: []string{"it"}},
		{name: "block one by one", area: layout.NewRect(0, 0, 1, 1), block: true, expected: []string{"┌"}},
		{name: "block one by two", area: layout.NewRect(0, 0, 1, 2), block: true, expected: []string{"┌", "└"}},
		{name: "block two by one", area: layout.NewRect(0, 0, 2, 1), block: true, expected: []string{"┌┐"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(tt.area)
			list := widgets.NewList([]widgets.ListItem{
				widgets.ListItemFromString("item"),
				widgets.ListItemFromString("item"),
			})
			if tt.block {
				list = list.Block(widgets.BorderedBlock())
			}

			assertNotPanics(t, func() {
				list.Render(buf.Area, buf)
			})

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestList_renderingCombinations(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 14, 6))
	state := widgets.ListState{}
	state.Select(1)
	listStyle := style.NewStyle().Bg(style.White)
	itemStyle := style.NewStyle().Fg(style.Red)
	highlightStyle := style.NewStyle().Bg(style.Blue).AddModifier(style.ModifierBold)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("top").Center()),
		widgets.ListItemFromLines(
			text.NewLine(text.StyledSpan("sel", itemStyle)).Center(),
			text.LineFromString("line").Right(),
		),
		widgets.ListItemFromString("tail"),
	}).
		Style(listStyle).
		Block(widgets.BorderedBlock()).
		HighlightStyle(highlightStyle).
		HighlightSymbol(">>").
		RepeatHighlightSymbol(true).
		HighlightSpacing(widgets.HighlightSpacingAlways).
		Right()

	list.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"┌────────────┐",
		"│     top    │",
		"│>>   sel    │",
		"│>>      line│",
		"│        tail│",
		"└────────────┘",
	})
	for x := 1; x <= 12; x++ {
		assertCellStyle(t, buf, x, 3, listStyle.Patch(highlightStyle))
	}
	for _, x := range []int{1, 2, 3, 4, 5, 9, 10, 11, 12} {
		assertCellStyle(t, buf, x, 2, listStyle.Patch(highlightStyle))
	}
	for x := 6; x <= 8; x++ {
		assertCellStyle(t, buf, x, 2, listStyle.Patch(highlightStyle).Patch(itemStyle))
	}
}

func TestList_noStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 2))

	widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("ab"),
	}).Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"ab  ",
		"    ",
	})
	assertAllCellsStyle(t, buf, style.NewStyle())
}

func TestList_renderInMinimalBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	assertNotPanics(t, func() {
		widgets.NewList([]widgets.ListItem{
			widgets.ListItemFromString("ab"),
		}).Render(buf.Area, buf)
	})

	assertLines(t, buf, []string{"a"})
	assertCellStyle(t, buf, 0, 0, style.NewStyle())
}

func TestList_renderInZeroSizeBuffer(t *testing.T) {
	tests := []struct {
		area     layout.Rect
		expected []string
	}{
		{area: layout.NewRect(0, 0, 0, 1), expected: []string{""}},
		{area: layout.NewRect(0, 0, 1, 0), expected: []string{}},
	}
	for _, tt := range tests {
		buf := buffer.Empty(tt.area)

		assertNotPanics(t, func() {
			widgets.NewList([]widgets.ListItem{
				widgets.ListItemFromString("ab"),
			}).Render(tt.area, buf)
		})

		assertLines(t, buf, tt.expected)
	}
}

func TestList_shouldHighlightSelectedItem(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	state := widgets.ListState{}
	state.Select(1)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
	}).HighlightStyle(style.NewStyle().Bg(style.Yellow)).HighlightSymbol(">> ")

	list.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   Item 1 ",
		">> Item 2 ",
		"   Item 3 ",
	})
	for x := range 10 {
		assertCellStyle(t, buf, x, 1, style.NewStyle().Bg(style.Yellow))
	}
}

func TestList_shouldHighlightSelectedItemWithWideSymbol(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	state := widgets.ListState{}
	state.Select(1)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
	}).HighlightStyle(style.NewStyle().Bg(style.Yellow)).HighlightSymbol("▶  ")

	list.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   Item 1 ",
		"▶  Item 2 ",
		"   Item 3 ",
	})
}

func TestList_shouldPreserveMultiSpanItemStyles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))
	line := text.NewLine(
		text.StyledSpan("ab", style.NewStyle().Fg(style.Red)),
		text.StyledSpan("cd", style.NewStyle().Fg(style.Green)),
	).Style(style.NewStyle().Bg(style.Blue))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(line),
	}).Style(style.NewStyle().AddModifier(style.ModifierBold))

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"abcd  "})
	assertCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold))
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).Bg(style.Blue).AddModifier(style.ModifierBold))
	assertCellStyle(t, buf, 2, 0, style.NewStyle().Fg(style.Green).Bg(style.Blue).AddModifier(style.ModifierBold))
	assertCellStyle(t, buf, 3, 0, style.NewStyle().Fg(style.Green).Bg(style.Blue).AddModifier(style.ModifierBold))
}

func TestList_shouldClipWideGraphemeByCellWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 1))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("aコb")),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"aコ"})
	assertCellSymbol(t, buf, 1, 0, "コ")
	assertCellSymbol(t, buf, 2, 0, " ")
}

func TestList_alignmentEvenLineEvenArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 4))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Odd").Left()),
		widgets.ListItemFromLines(text.LineFromString("Even").Center()),
		widgets.ListItemFromLines(text.LineFromString("Width").Right()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Odd   ", " Even ", " Width", "      "})
}

func TestList_alignmentEvenLineOddArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 4))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Odd").Left()),
		widgets.ListItemFromLines(text.LineFromString("Even").Center()),
		widgets.ListItemFromLines(text.LineFromString("Width").Right()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Odd   ", " Even ", " Width", "      "})
}

func TestList_alignmentOddLineEvenArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 4))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Odd").Left()),
		widgets.ListItemFromLines(text.LineFromString("Even").Center()),
		widgets.ListItemFromLines(text.LineFromString("Width").Right()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Odd     ", "  Even  ", "   Width", "        "})
}

func TestList_alignmentOddLineOddArea(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 4))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Odd").Left()),
		widgets.ListItemFromLines(text.LineFromString("Even").Center()),
		widgets.ListItemFromLines(text.LineFromString("Width").Right()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Odd    ", " Even  ", "  Width", "       "})
}

func TestList_alignmentLineEqualToWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Exact").Left()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Exact", "     "})
}

func TestList_alignmentLineGreaterThanWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 2))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Large line").Left()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Large", "     "})
}

func TestList_alignmentLineLessThanWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 2))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Small").Center()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"  Small   ", "          "})
}

func TestList_alignmentZeroAreaWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Text").Left()),
	})

	list.Render(layout.NewRect(0, 0, 4, 0), buf)

	assertLines(t, buf, []string{"    "})
}

func TestList_alignmentZeroLineWidth(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 2))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("This line has zero width").Center()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"", ""})
}

func TestList_withAlignment(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 4))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Left").Left()),
		widgets.ListItemFromLines(text.LineFromString("Center").Center()),
		widgets.ListItemFromLines(text.LineFromString("Right").Right()),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Left      ", "  Center  ", "     Right", "          "})
}

func TestList_alignmentPrecedence(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 4))
	textAligned := text.NewText(text.LineFromString("text")).Center()
	lineAligned := text.NewText(text.LineFromString("line").Right()).Center()
	listAligned := text.NewText(text.LineFromString("list"))
	defaultAligned := text.NewText(text.LineFromString("left"))
	list := widgets.NewList([]widgets.ListItem{
		widgets.NewListItem(lineAligned),
		widgets.NewListItem(textAligned),
		widgets.NewListItem(listAligned),
		widgets.NewListItem(defaultAligned),
	}).Right()

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"    line", "  text  ", "    list", "    left"})
}

func TestList_alignmentDefaultsLeft(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 8, 1))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("left"),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"left    "})
}

func TestList_alignmentUsesItemAreaAfterHighlightSymbol(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	state := widgets.ListState{}
	state.Select(0)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("abc").Center()),
	}).HighlightSymbol(">> ")

	list.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{">>   abc  "})
}

func TestList_shouldTruncateItems(t *testing.T) {
	tests := []struct {
		name     string
		selected bool
		expected []string
	}{
		{
			name:     "selected",
			selected: true,
			expected: []string{
				"┌──────┐  ",
				"│>> A v│  ",
				"└──────┘  ",
			},
		},
		{
			name: "not selected",
			expected: []string{
				"┌──────┐  ",
				"│A very│  ",
				"└──────┘  ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
			state := widgets.ListState{}
			if tt.selected {
				state.Select(0)
			}
			list := widgets.NewList([]widgets.ListItem{
				widgets.ListItemFromString("A very long line"),
				widgets.ListItemFromString("A very long line"),
			}).Block(widgets.BorderedBlock()).HighlightSymbol(">> ")

			list.RenderStateful(layout.NewRect(0, 0, 8, 3), buf, &state)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestList_shouldClampOffsetIfItemsAreRemoved(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 4))
	state := widgets.ListState{}
	state.Select(5)
	widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 0"),
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
		widgets.ListItemFromString("Item 4"),
		widgets.ListItemFromString("Item 5"),
	}).HighlightSymbol(">> ").RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   Item 2 ",
		"   Item 3 ",
		"   Item 4 ",
		">> Item 5 ",
	})
	if state.Offset() != 2 {
		t.Fatalf("Offset() = %d, want 2", state.Offset())
	}

	buf = buffer.Empty(layout.NewRect(0, 0, 10, 4))
	state.Select(1)
	widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 3"),
	}).HighlightSymbol(">> ").RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		">> Item 3 ",
		"          ",
		"          ",
		"          ",
	})
	if selected, ok := state.Selected(); !ok || selected != 0 {
		t.Fatalf("Selected() = %d, %v, want 0, true", selected, ok)
	}
}

func TestList_shouldDisplayMultilineItems(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 6))
	state := widgets.ListState{}
	state.Select(1)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Item 1"), text.LineFromString("Item 1a")),
		widgets.ListItemFromLines(text.LineFromString("Item 2"), text.LineFromString("Item 2b")),
		widgets.ListItemFromLines(text.LineFromString("Item 3"), text.LineFromString("Item 3c")),
	}).HighlightStyle(style.NewStyle().Bg(style.Yellow)).HighlightSymbol(">> ")

	list.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   Item 1 ",
		"   Item 1a",
		">> Item 2 ",
		"   Item 2b",
		"   Item 3 ",
		"   Item 3c",
	})
	for x := range 10 {
		assertCellStyle(t, buf, x, 2, style.NewStyle().Bg(style.Yellow))
		assertCellStyle(t, buf, x, 3, style.NewStyle().Bg(style.Yellow))
	}
}

func TestList_paddingFlicker(t *testing.T) {
	state := widgets.ListState{}
	state.SetOffset(2)
	state.Select(4)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 0"),
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
		widgets.ListItemFromString("Item 4"),
		widgets.ListItemFromString("Item 5"),
		widgets.ListItemFromString("Item 6"),
		widgets.ListItemFromString("Item 7"),
	}).HighlightSymbol(">> ").HighlightSpacing(widgets.HighlightSpacingAlways)

	first := buffer.Empty(layout.NewRect(0, 0, 10, 5))
	list.RenderStateful(first.Area, first, &state)
	offsetAfterFirstRender := state.Offset()
	linesAfterFirstRender := first.Lines()

	second := buffer.Empty(layout.NewRect(0, 0, 10, 5))
	list.RenderStateful(second.Area, second, &state)

	if state.Offset() != offsetAfterFirstRender {
		t.Fatalf("Offset() after second render = %d, want %d", state.Offset(), offsetAfterFirstRender)
	}
	if actual := second.Lines(); !slices.Equal(actual, linesAfterFirstRender) {
		t.Fatalf("lines mismatch after second render\nactual:   %#v\nexpected: %#v", actual, linesAfterFirstRender)
	}
	assertLines(t, first, []string{
		"   Item 2 ",
		"   Item 3 ",
		">> Item 4 ",
		"   Item 5 ",
		"   Item 6 ",
	})
}

func TestList_paddingInconsistentItemSizes(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	state := widgets.ListState{}
	state.Select(3)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 0"),
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
		widgets.ListItemFromLines(
			text.LineFromString("Item 4"),
			text.LineFromString("Test"),
			text.LineFromString("Test"),
		),
		widgets.ListItemFromString("Item 5"),
	}).HighlightSymbol(">> ").HighlightSpacing(widgets.HighlightSpacingAlways)

	list.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   Item 1 ",
		"   Item 2 ",
		">> Item 3 ",
	})
	if state.Offset() != 1 {
		t.Fatalf("Offset() = %d, want 1", state.Offset())
	}
}

func TestList_paddingOffsetPushbackBreak(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 4))
	state := widgets.ListState{}
	state.SetOffset(1)
	state.Select(2)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(
			text.LineFromString("Item 0"),
			text.LineFromString("Test"),
			text.LineFromString("Test"),
		),
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
	}).HighlightSymbol(">> ").HighlightSpacing(widgets.HighlightSpacingAlways)

	assertNotPanics(t, func() {
		list.RenderStateful(buf.Area, buf, &state)
	})

	assertLines(t, buf, []string{
		"   Item 1 ",
		">> Item 2 ",
		"   Item 3 ",
		"          ",
	})
	if state.Offset() != 1 {
		t.Fatalf("Offset() = %d, want 1", state.Offset())
	}

	buf = buffer.Empty(layout.NewRect(0, 0, 10, 3))
	state = widgets.ListState{}
	state.Select(1)
	list = widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 0"),
		widgets.ListItemFromLines(
			text.LineFromString("Item 1"),
			text.LineFromString("Test"),
			text.LineFromString("Test"),
			text.LineFromString("More"),
		),
		widgets.ListItemFromString("Item 2"),
	}).HighlightSymbol(">> ").HighlightSpacing(widgets.HighlightSpacingAlways)

	assertNotPanics(t, func() {
		list.RenderStateful(buf.Area, buf, &state)
	})

	assertLines(t, buf, []string{
		">> Item 1 ",
		"   Test   ",
		"   Test   ",
	})
	if state.Offset() != 1 {
		t.Fatalf("Offset() with oversized selected item = %d, want 1", state.Offset())
	}
}

func TestList_shouldRepeatHighlightSymbol(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 6))
	state := widgets.ListState{}
	state.Select(1)
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromLines(text.LineFromString("Item 1"), text.LineFromString("Item 1a")),
		widgets.ListItemFromLines(text.LineFromString("Item 2"), text.LineFromString("Item 2b")),
		widgets.ListItemFromLines(text.LineFromString("Item 3"), text.LineFromString("Item 3c")),
	}).HighlightStyle(style.NewStyle().Bg(style.Yellow)).HighlightSymbol(">> ").RepeatHighlightSymbol(true)

	list.RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   Item 1 ",
		"   Item 1a",
		">> Item 2 ",
		">> Item 2b",
		"   Item 3 ",
		"   Item 3c",
	})
}

func TestList_shouldNotIgnoreEmptyStringItems(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 4))
	list := widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString(""),
		widgets.ListItemFromString(""),
		widgets.ListItemFromString("Item 4"),
	})

	list.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Item 1", "      ", "      ", "Item 4"})
}

func TestList_highlightSpacing(t *testing.T) {
	tests := []struct {
		name     string
		selected *int
		spacing  widgets.HighlightSpacing
		expected []string
	}{
		{
			name:    "none when selected",
			spacing: widgets.HighlightSpacingWhenSelected,
			expected: []string{
				"┌─────────────┐",
				"│Item 1       │",
				"│Item 1a      │",
				"│Item 2       │",
				"│Item 2b      │",
				"│Item 3       │",
				"│Item 3c      │",
				"└─────────────┘",
			},
		},
		{
			name:    "none always",
			spacing: widgets.HighlightSpacingAlways,
			expected: []string{
				"┌─────────────┐",
				"│   Item 1    │",
				"│   Item 1a   │",
				"│   Item 2    │",
				"│   Item 2b   │",
				"│   Item 3    │",
				"│   Item 3c   │",
				"└─────────────┘",
			},
		},
		{
			name:     "first never",
			selected: new(0),
			spacing:  widgets.HighlightSpacingNever,
			expected: []string{
				"┌─────────────┐",
				"│Item 1       │",
				"│Item 1a      │",
				"│Item 2       │",
				"│Item 2b      │",
				"│Item 3       │",
				"│Item 3c      │",
				"└─────────────┘",
			},
		},
		{
			name:     "first when selected",
			selected: new(0),
			spacing:  widgets.HighlightSpacingWhenSelected,
			expected: []string{
				"┌─────────────┐",
				"│>> Item 1    │",
				"│   Item 1a   │",
				"│   Item 2    │",
				"│   Item 2b   │",
				"│   Item 3    │",
				"│   Item 3c   │",
				"└─────────────┘",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 15, 8))
			state := widgets.ListState{}
			if tt.selected != nil {
				state.Select(*tt.selected)
			}
			list := widgets.NewList([]widgets.ListItem{
				widgets.ListItemFromLines(text.LineFromString("Item 1"), text.LineFromString("Item 1a")),
				widgets.ListItemFromLines(text.LineFromString("Item 2"), text.LineFromString("Item 2b")),
				widgets.ListItemFromLines(text.LineFromString("Item 3"), text.LineFromString("Item 3c")),
			}).Block(widgets.BorderedBlock()).HighlightSymbol(">> ").HighlightSpacing(tt.spacing)

			list.RenderStateful(buf.Area, buf, &state)

			if actual := buf.Lines(); !slices.Equal(actual, tt.expected) {
				t.Fatalf("lines mismatch\nactual:   %#v\nexpected: %#v", actual, tt.expected)
			}
		})
	}
}
