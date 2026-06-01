package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestTabs_shouldNotPanic_whenAreaIsNarrow(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))
	tabs := widgets.TabsFromStrings([]string{"Tab1", "Tab2"})

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" "})
}

func TestTabs_shouldTruncateLastItem_whenAreaIsNarrow(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	tabs := widgets.TabsFromStrings([]string{"Tab1", "Tab2"})

	tabs.Render(layout.NewRect(0, 0, 9, 1), buf)

	assertLines(t, buf, []string{" Tab1 │ T "})
	for x := 1; x <= 4; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().AddModifier(style.ModifierReversed))
	}
	assertCellStyle(t, buf, 0, 0, style.NewStyle())
	assertCellStyle(t, buf, 5, 0, style.NewStyle())
	assertCellStyle(t, buf, 6, 0, style.NewStyle())
}

func TestTabs_shouldHighlightSelectedTab(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 13, 1))
	tabs := widgets.TabsFromStrings([]string{"Tab1", "Tab2"}).Select(1)

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" Tab1 │ Tab2 "})
	for x := 8; x <= 11; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().AddModifier(style.ModifierReversed))
	}
	assertCellStyle(t, buf, 1, 0, style.NewStyle())
}

func TestTabs_shouldPatchHighlightStyleOverSelectedTitleSpans(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	tabs := widgets.NewTabs([]text.Line{
		text.NewLine(text.StyledSpan("A", style.NewStyle().Fg(style.Cyan))),
	}).HighlightStyle(style.NewStyle().Fg(style.Yellow))

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" A   "})
	assertCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Yellow))
}

func TestTabs_shouldRenderCustomDividerAndPadding(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 1))
	tabs := widgets.TabsFromStrings([]string{"A", "B"}).Divider("/").Padding("[", "]")

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{"[A]/[B]    "})
}

func TestTabs_shouldRenderInsideBlock(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	tabs := widgets.TabsFromStrings([]string{"A", "B"}).Block(widgets.BorderedBlock())

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{
		"┌────────┐",
		"│ A │ B  │",
		"└────────┘",
	})
}

func TestTabs_shouldNotHighlight_whenSelectionIsCleared(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))
	tabs := widgets.NewTabs([]text.Line{text.LineFromString("Tab1")}).
		ClearSelection()

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" Tab1 "})
	for x := 0; x < 6; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle())
	}
}
