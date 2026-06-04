package widgets

import (
	"testing"

	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/symbols"
)

func renderTestShadow(shadow Shadow) *buffer.Buffer {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 4))
	shadow.Render(layout.NewRect(0, 0, 2, 2), buf)
	return buf
}

func TestShadow_overlayRendersStyleWithoutChangingSymbols(t *testing.T) {
	buf := buffer.WithLines([]string{"abcd", "efgh", "ijkl", "mnop"})
	shadow := NewShadowOverlay().Style(style.NewStyle().Fg(style.Red).Bg(style.Blue))

	shadow.Render(layout.NewRect(0, 0, 2, 2), buf)

	assertShadowCell(t, buf, 2, 1, "g", style.NewStyle().Fg(style.Red).Bg(style.Blue))
	assertShadowCell(t, buf, 1, 2, "j", style.NewStyle().Fg(style.Red).Bg(style.Blue))
	assertShadowCell(t, buf, 2, 2, "k", style.NewStyle().Fg(style.Red).Bg(style.Blue))
	assertShadowCell(t, buf, 1, 1, "f", style.NewStyle())
}

func TestShadow_symbolFiltersFillOnlyVisibleShadowCells(t *testing.T) {
	tests := []struct {
		name   string
		shadow Shadow
		symbol string
	}{
		{name: "custom symbol", shadow: NewShadowSymbol("$"), symbol: "$"},
		{name: "block", shadow: NewShadowBlock(), symbol: "█"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := renderTestShadow(tt.shadow)

			assertShadowCell(t, buf, 2, 1, tt.symbol, style.NewStyle())
			assertShadowCell(t, buf, 1, 2, tt.symbol, style.NewStyle())
			assertShadowCell(t, buf, 2, 2, tt.symbol, style.NewStyle())
			assertShadowCell(t, buf, 1, 1, " ", style.NewStyle())
		})
	}
}

func TestShadow_shadeConstructorsShouldUseShadeSymbols(t *testing.T) {
	tests := []struct {
		name   string
		shadow Shadow
		want   string
	}{
		{name: "block", shadow: NewShadowBlock(), want: symbols.ShadeFull},
		{name: "light", shadow: NewShadowLightShade(), want: symbols.ShadeLight},
		{name: "medium", shadow: NewShadowMediumShade(), want: symbols.ShadeMedium},
		{name: "dark", shadow: NewShadowDarkShade(), want: symbols.ShadeDark},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := renderTestShadow(tt.shadow)

			assertShadowCell(t, buf, 2, 1, tt.want, style.NewStyle())
		})
	}
}

func TestShadow_renderIsClippedToBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 3, 2))
	shadow := NewShadowSymbol("#")

	shadow.Render(layout.NewRect(0, 0, 2, 1), buf)

	assertShadowCell(t, buf, 2, 1, "#", style.NewStyle())
}

func TestShadow_customFilterIsApplied(t *testing.T) {
	buf := renderTestShadow(NewShadowCustom(plusEffect{}))

	assertShadowCell(t, buf, 2, 1, "+", style.NewStyle())
	assertShadowCell(t, buf, 1, 2, "+", style.NewStyle())
	assertShadowCell(t, buf, 2, 2, "+", style.NewStyle())
}

func TestShadow_dimmedFilterDimsBackground(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 4))
	buf.SetStyle(buf.Area, style.NewStyle().Bg(style.Red))
	shadow := NewShadowCustom(Dimmed())

	shadow.Render(layout.NewRect(0, 0, 2, 2), buf)

	assertShadowCell(t, buf, 2, 1, " ", style.NewStyle().Bg(style.Black).AddModifier(style.ModifierDim))
	assertShadowCell(t, buf, 1, 1, " ", style.NewStyle().Bg(style.Red))
}

func TestBlock_renderShadowAfterBorders(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 4))

	BorderedBlock().Shadow(NewShadowDarkShade()).Render(layout.NewRect(0, 0, 3, 2), buf)

	assertBlockLines(t, buf, []string{
		"┌─┐  ",
		"└─┘▓ ",
		" ▓▓▓ ",
		"     ",
	})
}

type plusEffect struct{}

func (plusEffect) Apply(shadowArea layout.Rect, baseArea layout.Rect, buf *buffer.Buffer) {
	forEachShadowCell(shadowArea, baseArea, buf, func(cell *buffer.Cell) {
		cell.SetSymbol("+")
	})
}

func assertShadowCell(t *testing.T, buf *buffer.Buffer, x, y int, symbol string, cellStyle style.Style) {
	t.Helper()
	cell, ok := buf.CellAt(x, y)
	if !ok {
		t.Fatalf("missing cell at (%d,%d)", x, y)
	}
	if cell.DisplaySymbol() != symbol {
		t.Fatalf("cell(%d,%d).Symbol = %q, want %q", x, y, cell.DisplaySymbol(), symbol)
	}
	if cell.Style != cellStyle {
		t.Fatalf("cell(%d,%d).Style = %#v, want %#v", x, y, cell.Style, cellStyle)
	}
}
