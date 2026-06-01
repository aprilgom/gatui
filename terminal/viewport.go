package terminal

import (
	"fmt"

	"gatui/layout"
)

type viewportKind int

const (
	viewportFullscreen viewportKind = iota
	viewportFixed
	viewportInline
)

type Viewport struct {
	kind   viewportKind
	area   layout.Rect
	height int
}

func FullscreenViewport() Viewport {
	return Viewport{kind: viewportFullscreen}
}

func FixedViewport(area layout.Rect) Viewport {
	return Viewport{kind: viewportFixed, area: area}
}

func InlineViewport(height int) Viewport {
	return Viewport{kind: viewportInline, height: height}
}

func (v Viewport) String() string {
	switch v.kind {
	case viewportFullscreen:
		return "Fullscreen"
	case viewportInline:
		return fmt.Sprintf("Inline(%d)", v.height)
	case viewportFixed:
		return fmt.Sprintf("Fixed(%s)", formatViewportRect(v.area))
	default:
		return "Unknown"
	}
}

func formatViewportRect(area layout.Rect) string {
	return fmt.Sprintf("Rect { x: %d, y: %d, width: %d, height: %d }", area.X, area.Y, area.Width, area.Height)
}
