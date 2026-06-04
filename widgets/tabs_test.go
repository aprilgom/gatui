package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
	"gatui/widgets"
)

func TestTabs_canBeStylized(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))
	tabs := widgets.TabsFromStrings([]string{"Tab"}).
		Fg(style.Red).
		Bg(style.White).
		Bold().
		Dim().
		Italic().
		Cyan()

	tabs.Render(buf.Area, buf)

	want := style.NewStyle().
		Fg(style.Cyan).
		Bg(style.White).
		AddModifier(style.ModifierBold | style.ModifierDim | style.ModifierItalic)
	for x := range 6 {
		if x >= 1 && x <= 3 {
			assertCellStyle(t, buf, x, 0, want.AddModifier(style.ModifierReversed))
			continue
		}
		assertCellStyle(t, buf, x, 0, want)
	}
}

func TestTabs_new(t *testing.T) {
	titles := []text.Line{text.LineFromString("One")}
	tabs := widgets.NewTabs(titles)
	titles[0] = text.LineFromString("Changed")
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" One  "})
	assertCellStyle(t, buf, 1, 0, style.NewStyle().AddModifier(style.ModifierReversed))
}

func TestTabs_newFromVecOfStr(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 13, 1))
	tabs := widgets.TabsFromStrings([]string{"Tab1", "Tab2"})

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" Tab1 │ Tab2 "})
	for x := 1; x <= 4; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().AddModifier(style.ModifierReversed))
	}
}

func TestTabs_collectStrings_shouldMatchRatatui(t *testing.T) {
	titles := []string{"Tab0", "Tab1", "Tab2", "Tab3", "Tab4"}
	tabs := widgets.TabsFromStrings(titles)
	titles[0] = "Changed"
	buf := buffer.Empty(layout.NewRect(0, 0, 34, 1))

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" Tab0 │ Tab1 │ Tab2 │ Tab3 │ Tab4 "})
	for x := 1; x <= 4; x++ {
		assertCellStyle(t, buf, x, 0, style.NewStyle().AddModifier(style.ModifierReversed))
	}
	assertCellStyle(t, buf, 8, 0, style.NewStyle())
}

func TestDefaultTabs_shouldMatchRatatui(t *testing.T) {
	baseStyle := style.NewStyle().Fg(style.Blue)
	buf := buffer.Empty(layout.NewRect(0, 0, 6, 1))
	buf.SetString(0, 0, "seeded", baseStyle)

	widgets.DefaultTabs().Render(buf.Area, buf)

	assertLines(t, buf, []string{"seeded"})
	for x := range 6 {
		assertCellStyle(t, buf, x, 0, baseStyle)
	}

	widgets.DefaultTabs().Select(0).Render(buf.Area, buf)

	assertLines(t, buf, []string{"seeded"})
	for x := range 6 {
		assertCellStyle(t, buf, x, 0, baseStyle)
	}
}

func TestTabs_renderNew(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 19, 1))
	tabs := widgets.NewTabs([]text.Line{
		text.LineFromString("Tab0"),
		text.LineFromString("Tab1"),
		text.LineFromString("Tab2"),
	})

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" Tab0 │ Tab1 │ Tab2"})
}

func TestTabs_renderNoPadding(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 12, 1))
	tabs := widgets.TabsFromStrings([]string{"Tab0", "Tab1", "Tab2"}).Padding("", "")

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{"Tab0│Tab1│Ta"})
}

func TestTabs_selectBeforeTitles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 13, 1))
	tabs := widgets.TabsFromStrings([]string{"Tab1", "Tab2"}).Select(4)

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" Tab1 │ Tab2 "})
	for x := range 13 {
		assertCellStyle(t, buf, x, 0, style.NewStyle())
	}
}

func TestTabs_unicodeWidthCJK(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 13, 1))
	tabs := widgets.TabsFromStrings([]string{"コン", "Tab"})

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{" コン │ Tab  "})
	assertCellSymbol(t, buf, 1, 0, "コ")
	assertCellSymbol(t, buf, 2, 0, " ")
	assertCellSymbol(t, buf, 3, 0, "ン")
	assertCellSymbol(t, buf, 4, 0, " ")
	assertCellSymbol(t, buf, 6, 0, "│")
	assertCellStyle(t, buf, 1, 0, style.NewStyle().AddModifier(style.ModifierReversed))
	assertCellStyle(t, buf, 3, 0, style.NewStyle().AddModifier(style.ModifierReversed))
}

func TestTabs_unicodeWidthCJKCustomPaddingAndDivider(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 14, 1))
	tabs := widgets.TabsFromStrings([]string{"コ", "ン"}).Padding("[]", "[]").Divider("界")

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{"[]コ[]界[]ン[]"})
	assertCellSymbol(t, buf, 2, 0, "コ")
	assertCellSymbol(t, buf, 3, 0, " ")
	assertCellSymbol(t, buf, 6, 0, "界")
	assertCellSymbol(t, buf, 7, 0, " ")
	assertCellSymbol(t, buf, 10, 0, "ン")
	assertCellSymbol(t, buf, 11, 0, " ")
}

func TestTabs_unicodeWidthEmptyTitles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	tabs := widgets.TabsFromStrings([]string{"", ""})

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{"  │  "})
}

func TestTabs_unicodeWidthNoPadding(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 1))
	tabs := widgets.TabsFromStrings([]string{"コ", "ン", "A"}).Padding("", "")

	tabs.Render(buf.Area, buf)

	assertLines(t, buf, []string{"コ│ン│A"})
	assertCellSymbol(t, buf, 0, 0, "コ")
	assertCellSymbol(t, buf, 1, 0, " ")
	assertCellSymbol(t, buf, 2, 0, "│")
	assertCellSymbol(t, buf, 3, 0, "ン")
	assertCellSymbol(t, buf, 4, 0, " ")
	assertCellSymbol(t, buf, 5, 0, "│")
}

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
	for x := range 6 {
		assertCellStyle(t, buf, x, 0, style.NewStyle())
	}
}
