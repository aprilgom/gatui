package buffer_test

import (
	"reflect"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
)

func TestWithLines_shouldCreateBlankPaddedBuffer(t *testing.T) {
	buf := buffer.WithLines([]string{"ab", "c"})

	if got, want := buf.Area, layout.NewRect(0, 0, 2, 2); got != want {
		t.Fatalf("area = %#v, want %#v", got, want)
	}
	if got, want := buf.Lines(), []string{"ab", "c "}; !reflect.DeepEqual(got, want) {
		t.Fatalf("lines = %#v, want %#v", got, want)
	}
}

func TestSetStyleHelpers_shouldPatchExistingCells(t *testing.T) {
	buf := buffer.WithLines([]string{"ab"})
	buf.SetFg(layout.NewRect(0, 0, 1, 1), style.Red)
	buf.SetModifier(layout.NewRect(0, 0, 1, 1), style.ModifierBold)

	cell, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatal("expected cell")
	}
	want := style.NewStyle().Fg(style.Red).AddModifier(style.ModifierBold)
	if cell.Style != want {
		t.Fatalf("style = %#v, want %#v", cell.Style, want)
	}
}
