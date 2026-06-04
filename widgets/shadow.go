package widgets

import (
	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/symbols"
)

type CellEffect interface {
	Apply(shadowArea layout.Rect, baseArea layout.Rect, buf *buffer.Buffer)
}

type Shadow struct {
	effect CellEffect
	style  style.Style
	offset layout.Offset
}

func NewShadowOverlay() Shadow {
	return Shadow{
		style:  style.NewStyle(),
		offset: layout.NewOffset(1, 1),
	}
}

func NewShadowSymbol(symbol string) Shadow {
	return NewShadowCustom(symbolEffect{symbol: symbol})
}

func NewShadowBlock() Shadow {
	return NewShadowSymbol(symbols.ShadeFull)
}

func NewShadowLightShade() Shadow {
	return NewShadowSymbol(symbols.ShadeLight)
}

func NewShadowMediumShade() Shadow {
	return NewShadowSymbol(symbols.ShadeMedium)
}

func NewShadowDarkShade() Shadow {
	return NewShadowSymbol(symbols.ShadeDark)
}

func NewShadowCustom(effect CellEffect) Shadow {
	shadow := NewShadowOverlay()
	shadow.effect = effect
	return shadow
}

func (s Shadow) Style(shadowStyle style.Style) Shadow {
	s.style = shadowStyle
	return s
}

func (s Shadow) Offset(offset layout.Offset) Shadow {
	s.offset = offset
	return s
}

func (s Shadow) Render(area layout.Rect, buf *buffer.Buffer) {
	if buf == nil || area.IsEmpty() {
		return
	}
	shadowArea := area.Offset(s.offset).Intersection(buf.Area)
	forEachShadowCell(shadowArea, area, buf, func(cell *buffer.Cell) {
		cell.SetStyle(s.style)
	})
	if s.effect != nil {
		s.effect.Apply(shadowArea, area, buf)
	}
}

type symbolEffect struct {
	symbol string
}

func (e symbolEffect) Apply(shadowArea layout.Rect, baseArea layout.Rect, buf *buffer.Buffer) {
	forEachShadowCell(shadowArea, baseArea, buf, func(cell *buffer.Cell) {
		cell.SetSymbol(e.symbol)
	})
}

type dimmedEffect struct{}

func Dimmed() CellEffect {
	return dimmedEffect{}
}

func (dimmedEffect) Apply(shadowArea layout.Rect, baseArea layout.Rect, buf *buffer.Buffer) {
	forEachShadowCell(shadowArea, baseArea, buf, func(cell *buffer.Cell) {
		cell.SetStyle(style.NewStyle().Bg(style.Black).AddModifier(style.ModifierDim))
	})
}

func forEachShadowCell(shadowArea layout.Rect, baseArea layout.Rect, buf *buffer.Buffer, apply func(*buffer.Cell)) {
	if buf == nil || shadowArea.IsEmpty() {
		return
	}
	for y := shadowArea.Y; y < shadowArea.Bottom(); y++ {
		for x := shadowArea.X; x < shadowArea.Right(); x++ {
			if baseArea.Contains(layout.NewPosition(x, y)) {
				continue
			}
			cell, ok := buf.CellRef(x, y)
			if !ok {
				continue
			}
			apply(cell)
		}
	}
}
