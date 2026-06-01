package terminal

import (
	"fmt"

	"gatui/buffer"
	"gatui/layout"
)

type ClearType int

const (
	ClearAll ClearType = iota
	ClearAfterCursor
	ClearBeforeCursor
	ClearCurrentLine
	ClearUntilNewLine
)

func (c ClearType) String() string {
	switch c {
	case ClearAll:
		return "All"
	case ClearAfterCursor:
		return "AfterCursor"
	case ClearBeforeCursor:
		return "BeforeCursor"
	case ClearCurrentLine:
		return "CurrentLine"
	case ClearUntilNewLine:
		return "UntilNewLine"
	default:
		return "Unknown"
	}
}

func ParseClearType(value string) (ClearType, error) {
	switch value {
	case "All":
		return ClearAll, nil
	case "AfterCursor":
		return ClearAfterCursor, nil
	case "BeforeCursor":
		return ClearBeforeCursor, nil
	case "CurrentLine":
		return ClearCurrentLine, nil
	case "UntilNewLine":
		return ClearUntilNewLine, nil
	default:
		return ClearType(0), fmt.Errorf("unknown clear type: %s", value)
	}
}

func (t *Terminal) Clear() error {
	originalCursor, err := t.backend.GetCursorPosition()
	if err != nil {
		return err
	}
	if err := t.clearViewport(); err != nil {
		return err
	}
	return t.backend.SetCursorPosition(originalCursor)
}

func (t *Terminal) clearViewport() error {
	switch t.viewport.kind {
	case viewportFullscreen:
		if err := t.backend.ClearRegion(ClearAll); err != nil {
			return err
		}
	case viewportFixed:
		if err := t.clearFixedViewport(t.area); err != nil {
			return err
		}
	case viewportInline:
		if err := t.backend.SetCursorPosition(layout.Position{X: t.area.X, Y: t.area.Y}); err != nil {
			return err
		}
		if err := t.backend.ClearRegion(ClearAfterCursor); err != nil {
			return err
		}
	}
	t.previous.Reset()
	return nil
}

func (t *Terminal) clearFixedViewport(area layout.Rect) error {
	if area.Width == 0 || area.Height == 0 {
		return nil
	}
	size, err := t.backend.Size()
	if err != nil {
		return err
	}
	isFullWidth := area.X == 0 && area.Width == size.Width
	endsAtBottom := area.Bottom() == size.Height
	if isFullWidth && endsAtBottom {
		if err := t.backend.SetCursorPosition(layout.Position{X: area.X, Y: area.Y}); err != nil {
			return err
		}
		return t.backend.ClearRegion(ClearAfterCursor)
	}
	if isFullWidth {
		for y := area.Y; y < area.Bottom(); y++ {
			if err := t.backend.SetCursorPosition(layout.Position{X: 0, Y: y}); err != nil {
				return err
			}
			if err := t.backend.ClearRegion(ClearCurrentLine); err != nil {
				return err
			}
		}
		return nil
	}

	clearCell := buffer.NewCell(" ")
	diffs := make([]buffer.CellDiff, 0, area.Width*area.Height)
	for y := area.Y; y < area.Bottom(); y++ {
		for x := area.X; x < area.Right(); x++ {
			diffs = append(diffs, buffer.CellDiff{X: x, Y: y, Cell: clearCell})
		}
	}
	return t.backend.Draw(diffs)
}
