package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/widgets"
)

func TestList_StatefulRef_shouldRenderWithState(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	state := widgets.ListState{}
	state.Select(1)
	var widget widgets.StatefulWidget = widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item 1"),
		widgets.ListItemFromString("Item 2"),
		widgets.ListItemFromString("Item 3"),
	}).HighlightSymbol(">> ")

	widget.RenderStatefulRef(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   Item 1 ",
		">> Item 2 ",
		"   Item 3 ",
	})
}

func TestList_StatefulRef_shouldPanicOnWrongStateType(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	var widget widgets.StatefulWidget = widgets.NewList([]widgets.ListItem{
		widgets.ListItemFromString("Item"),
	})

	assertPanicMessage(t, "gatui: invalid state type for List", func() {
		widget.RenderStatefulRef(buf.Area, buf, &widgets.TableState{})
	})
}

func TestTable_StatefulRef_shouldRenderWithState(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 2))
	state := widgets.NewTableState().WithSelected(1)
	var widget widgets.StatefulWidget = widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"one"}),
		widgets.TableRowFromStrings([]string{"two"}),
	}, []layout.Constraint{layout.Length(3)}).HighlightSymbol(">> ")

	widget.RenderStatefulRef(buf.Area, buf, &state)

	assertLines(t, buf, []string{
		"   one ",
		">> two ",
	})
}

func TestTable_StatefulRef_shouldPanicOnWrongStateType(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 1))
	var widget widgets.StatefulWidget = widgets.NewTable([]widgets.TableRow{
		widgets.TableRowFromStrings([]string{"one"}),
	}, []layout.Constraint{layout.Length(3)})

	assertPanicMessage(t, "gatui: invalid state type for Table", func() {
		widget.RenderStatefulRef(buf.Area, buf, &widgets.ListState{})
	})
}

func TestScrollbar_StatefulRef_shouldRenderWithState(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))
	state := widgets.NewScrollbarState(2).Position(1)
	var widget widgets.StatefulWidget = widgets.NewScrollbar(widgets.ScrollbarOrientationHorizontalTop).ClearBeginSymbol().ClearEndSymbol()

	widget.RenderStatefulRef(buf.Area, buf, &state)

	assertLines(t, buf, []string{"═█"})
}

func TestScrollbar_StatefulRef_shouldPanicOnWrongStateType(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 2, 1))
	var widget widgets.StatefulWidget = widgets.NewScrollbar(widgets.ScrollbarOrientationHorizontalTop)

	assertPanicMessage(t, "gatui: invalid state type for Scrollbar", func() {
		widget.RenderStatefulRef(buf.Area, buf, &widgets.ListState{})
	})
}

func assertPanicMessage(t *testing.T, want string, fn func()) {
	t.Helper()
	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("panic = nil, want %q", want)
		}
		if got != want {
			t.Fatalf("panic = %#v, want %q", got, want)
		}
	}()
	fn()
}
