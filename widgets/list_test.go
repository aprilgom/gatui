package widgets_test

import (
	"slices"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

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
