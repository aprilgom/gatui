package widgets

import (
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/text"
)

func TestCellDisplayWidth_shouldMeasureASCIIAsOneColumn(t *testing.T) {
	cells := cellsFromLine(text.LineFromString("abc"))

	if got, want := cellsDisplayWidth(cells), 3; got != want {
		t.Fatalf("cellsDisplayWidth = %d, want %d", got, want)
	}
}

func TestCellDisplayWidth_shouldMeasureCJKAsTwoColumns(t *testing.T) {
	cells := cellsFromLine(text.LineFromString("コ"))

	if got, want := cellsDisplayWidth(cells), 2; got != want {
		t.Fatalf("cellsDisplayWidth = %d, want %d", got, want)
	}
}

func TestCellDisplayWidth_shouldMeasureMixedTextByDisplayColumns(t *testing.T) {
	cells := cellsFromLine(text.LineFromString("aコンピ"))

	if got, want := cellsDisplayWidth(cells), 7; got != want {
		t.Fatalf("cellsDisplayWidth = %d, want %d", got, want)
	}
}

func TestCellDisplayWidth_shouldUseBufferCellWidth(t *testing.T) {
	cells := cellsFromLine(text.LineFromString("ﾞ"))

	if got, want := cellsDisplayWidth(cells), buffer.CellWidth("ﾞ"); got != want {
		t.Fatalf("cellsDisplayWidth = %d, want %d", got, want)
	}
}
