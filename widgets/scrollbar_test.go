package widgets_test

import (
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/widgets"
)

func TestScrollbar_shouldRenderSimplestHorizontalNoArrows(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		position int
	}{
		{name: "position_0", expected: "#-", position: 0},
		{name: "position_1", expected: "-#", position: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, len([]rune(tt.expected)), 1))
			state := widgets.NewScrollbarState(2).Position(tt.position)

			scrollbarNoArrows().RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, []string{tt.expected})
		})
	}
}

func TestScrollbar_shouldRenderHorizontalNoArrows(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		position int
	}{
		{name: "position_0", expected: "#####-----", position: 0},
		{name: "position_1", expected: "-#####----", position: 1},
		{name: "position_2", expected: "-#####----", position: 2},
		{name: "position_3", expected: "--#####---", position: 3},
		{name: "position_4", expected: "--#####---", position: 4},
		{name: "position_5", expected: "---#####--", position: 5},
		{name: "position_6", expected: "---#####--", position: 6},
		{name: "position_7", expected: "----#####-", position: 7},
		{name: "position_8", expected: "----#####-", position: 8},
		{name: "position_9", expected: "-----#####", position: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, len([]rune(tt.expected)), 1))
			state := widgets.NewScrollbarState(10).Position(tt.position)

			scrollbarNoArrows().RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, []string{tt.expected})
		})
	}
}

func TestScrollbar_shouldRenderNothing_whenContentLengthIsZero(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	state := widgets.NewScrollbarState(0)

	scrollbarNoArrows().RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{"          "})
}

func TestScrollbar_shouldClampOutOfBoundsPosition(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	state := widgets.NewScrollbarState(10).Position(100)

	scrollbarNoArrows().RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{"-----#####"})
}

func TestScrollbar_customViewportLength(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		position int
	}{
		{name: "position_0", expected: "##--------", position: 0},
		{name: "position_1", expected: "-##-------", position: 1},
		{name: "position_2", expected: "--##------", position: 2},
		{name: "position_3", expected: "---##-----", position: 3},
		{name: "position_4", expected: "----##----", position: 4},
		{name: "position_5", expected: "-----##---", position: 5},
		{name: "position_6", expected: "-----##---", position: 6},
		{name: "position_7", expected: "------##--", position: 7},
		{name: "position_8", expected: "-------##-", position: 8},
		{name: "position_9", expected: "--------##", position: 9},
		{name: "position_out_of_bounds", expected: "--------##", position: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, len([]rune(tt.expected)), 1))
			state := widgets.NewScrollbarState(10).
				Position(tt.position).
				ViewportContentLength(2)

			scrollbarNoArrows().RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, []string{tt.expected})
		})
	}
}

func TestScrollbar_shouldRenderVerticalRightWithCustomSymbols(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		position int
	}{
		{name: "position_0", expected: "<####---->", position: 0},
		{name: "position_2", expected: "<-####--->", position: 2},
		{name: "position_4", expected: "<--####-->", position: 4},
		{name: "position_6", expected: "<---####->", position: 6},
		{name: "position_9", expected: "<----####>", position: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 5, len([]rune(tt.expected))))
			state := widgets.NewScrollbarState(10).Position(tt.position)

			widgets.NewScrollbar(widgets.ScrollbarOrientationVerticalRight).
				BeginSymbol("<").
				EndSymbol(">").
				TrackSymbol("-").
				ThumbSymbol("#").
				RenderStateful(buf.Area, buf, &state)

			expected := make([]string, 0, len([]rune(tt.expected)))
			for _, r := range tt.expected {
				expected = append(expected, "    "+string(r))
			}
			assertLines(t, buf, expected)
		})
	}
}

func TestScrollbarState_shouldNavigate(t *testing.T) {
	state := widgets.NewScrollbarState(3)

	state.Next()
	state.Next()
	state.Next()
	if got := state.PositionValue(); got != 2 {
		t.Fatalf("position after nexts = %d, want 2", got)
	}

	state.Previous()
	if got := state.PositionValue(); got != 1 {
		t.Fatalf("position after previous = %d, want 1", got)
	}

	state.First()
	if got := state.PositionValue(); got != 0 {
		t.Fatalf("position after first = %d, want 0", got)
	}

	state.Last()
	if got := state.PositionValue(); got != 2 {
		t.Fatalf("position after last = %d, want 2", got)
	}

	state.Scroll(widgets.ScrollDirectionBackward)
	if got := state.PositionValue(); got != 1 {
		t.Fatalf("position after backward scroll = %d, want 1", got)
	}

	state.Scroll(widgets.ScrollDirectionForward)
	if got := state.PositionValue(); got != 2 {
		t.Fatalf("position after forward scroll = %d, want 2", got)
	}
}

func TestScrollbar_doNotRenderWithEmptyArea(t *testing.T) {
	tests := []struct {
		name string
		area layout.Rect
	}{
		{name: "zero_height", area: layout.NewRect(0, 0, 10, 0)},
		{name: "zero_width", area: layout.NewRect(0, 0, 0, 10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 10, 10))
			state := widgets.NewScrollbarState(10)

			widgets.NewScrollbar(widgets.ScrollbarOrientationVerticalRight).
				BeginSymbol("<").
				EndSymbol(">").
				TrackSymbol("-").
				ThumbSymbol("#").
				RenderStateful(tt.area, buf, &state)

			assertLines(t, buf, []string{
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
				"          ",
			})
		})
	}
}

func TestScrollbar_partLengthsReturnZerosWhenAreaDimensionIsZeroEvenWithoutArrows(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 5, 1))
	state := widgets.NewScrollbarState(10).
		Position(3).
		ViewportContentLength(2)

	scrollbarNoArrows().RenderStateful(layout.NewRect(0, 0, 0, 1), buf, &state)

	assertLines(t, buf, []string{"     "})
}

func TestScrollbar_partLengthsReturnZerosWhenTrackLenIsZero(t *testing.T) {
	tests := []struct {
		name        string
		orientation widgets.ScrollbarOrientation
		area        layout.Rect
		expected    []string
	}{
		{
			name:        "horizontal_width_equal_arrows",
			orientation: widgets.ScrollbarOrientationHorizontalTop,
			area:        layout.NewRect(0, 0, 2, 1),
			expected:    []string{"  "},
		},
		{
			name:        "horizontal_width_less_than_arrows",
			orientation: widgets.ScrollbarOrientationHorizontalTop,
			area:        layout.NewRect(0, 0, 1, 1),
			expected:    []string{" "},
		},
		{
			name:        "vertical_height_equal_arrows",
			orientation: widgets.ScrollbarOrientationVerticalLeft,
			area:        layout.NewRect(0, 0, 1, 2),
			expected:    []string{" ", " "},
		},
		{
			name:        "vertical_height_less_than_arrows",
			orientation: widgets.ScrollbarOrientationVerticalLeft,
			area:        layout.NewRect(0, 0, 1, 1),
			expected:    []string{" "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(tt.area)
			state := widgets.NewScrollbarState(10).
				Position(5).
				ViewportContentLength(2)

			widgets.NewScrollbar(tt.orientation).
				BeginSymbol("<").
				EndSymbol(">").
				TrackSymbol("-").
				ThumbSymbol("#").
				RenderStateful(tt.area, buf, &state)

			assertLines(t, buf, tt.expected)
		})
	}
}

func TestScrollbar_renderInMinimalBuffer(t *testing.T) {
	for _, orientation := range []widgets.ScrollbarOrientation{
		widgets.ScrollbarOrientationVerticalLeft,
		widgets.ScrollbarOrientationVerticalRight,
		widgets.ScrollbarOrientationHorizontalTop,
		widgets.ScrollbarOrientationHorizontalBottom,
	} {
		buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))
		state := widgets.NewScrollbarState(10).Position(5)

		widgets.NewScrollbar(orientation).RenderStateful(buf.Area, buf, &state)

		assertLines(t, buf, []string{" "})
	}
}

func TestScrollbar_renderInZeroSizeBuffer(t *testing.T) {
	for _, orientation := range []widgets.ScrollbarOrientation{
		widgets.ScrollbarOrientationVerticalLeft,
		widgets.ScrollbarOrientationVerticalRight,
		widgets.ScrollbarOrientationHorizontalTop,
		widgets.ScrollbarOrientationHorizontalBottom,
	} {
		buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))
		state := widgets.NewScrollbarState(10).Position(5)

		widgets.NewScrollbar(orientation).RenderStateful(buf.Area, buf, &state)
	}
}

func TestScrollbar_renderWithoutTrackSymbols(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		position int
	}{
		{name: "position_0", expected: "█████     ", position: 0},
		{name: "position_1", expected: " █████    ", position: 1},
		{name: "position_2", expected: " █████    ", position: 2},
		{name: "position_3", expected: "  █████   ", position: 3},
		{name: "position_4", expected: "  █████   ", position: 4},
		{name: "position_5", expected: "   █████  ", position: 5},
		{name: "position_6", expected: "   █████  ", position: 6},
		{name: "position_7", expected: "    █████ ", position: 7},
		{name: "position_8", expected: "    █████ ", position: 8},
		{name: "position_9", expected: "     █████", position: 9},
		{name: "position_out_of_bounds", expected: "     █████", position: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, len([]rune(tt.expected)), 1))
			state := widgets.NewScrollbarState(10).Position(tt.position)

			widgets.NewScrollbar(widgets.ScrollbarOrientationHorizontalBottom).
				ClearTrackSymbol().
				ClearBeginSymbol().
				ClearEndSymbol().
				RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, []string{tt.expected})
		})
	}
}

func TestScrollbar_renderWithoutTrackSymbolsOverContent(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		position int
	}{
		{name: "position_0", expected: "█████-----", position: 0},
		{name: "position_1", expected: "-█████----", position: 1},
		{name: "position_2", expected: "-█████----", position: 2},
		{name: "position_3", expected: "--█████---", position: 3},
		{name: "position_4", expected: "--█████---", position: 4},
		{name: "position_5", expected: "---█████--", position: 5},
		{name: "position_6", expected: "---█████--", position: 6},
		{name: "position_7", expected: "----█████-", position: 7},
		{name: "position_8", expected: "----█████-", position: 8},
		{name: "position_9", expected: "-----█████", position: 9},
		{name: "position_out_of_bounds", expected: "-----█████", position: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.WithLines([]string{"----------"})
			state := widgets.NewScrollbarState(10).Position(tt.position)

			widgets.NewScrollbar(widgets.ScrollbarOrientationHorizontalBottom).
				ClearTrackSymbol().
				ClearBeginSymbol().
				ClearEndSymbol().
				RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, []string{tt.expected})
		})
	}
}

func TestScrollbar_thumbVisibleOnVerySmallTrack(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		position int
	}{
		{name: "position_0", expected: "#----", position: 0},
		{name: "position_10", expected: "#----", position: 10},
		{name: "position_20", expected: "-#---", position: 20},
		{name: "position_30", expected: "-#---", position: 30},
		{name: "position_40", expected: "--#--", position: 40},
		{name: "position_50", expected: "--#--", position: 50},
		{name: "position_60", expected: "---#-", position: 60},
		{name: "position_70", expected: "---#-", position: 70},
		{name: "position_80", expected: "----#", position: 80},
		{name: "position_90", expected: "----#", position: 90},
		{name: "position_out_of_bounds", expected: "----#", position: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, len([]rune(tt.expected)), 1))
			state := widgets.NewScrollbarState(100).
				Position(tt.position).
				ViewportContentLength(2)

			scrollbarNoArrows().RenderStateful(buf.Area, buf, &state)

			assertLines(t, buf, []string{tt.expected})
		})
	}
}

func TestScrollbar_shouldApplyStyles(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
	state := widgets.NewScrollbarState(4)
	beginStyle := style.NewStyle().Fg(style.Red)
	thumbStyle := style.NewStyle().Fg(style.Green)
	trackStyle := style.NewStyle().Fg(style.Blue)
	endStyle := style.NewStyle().Fg(style.Yellow)

	widgets.NewScrollbar(widgets.ScrollbarOrientationHorizontalTop).
		BeginSymbol("<").
		EndSymbol(">").
		TrackSymbol("-").
		ThumbSymbol("#").
		BeginStyle(beginStyle).
		ThumbStyle(thumbStyle).
		TrackStyle(trackStyle).
		EndStyle(endStyle).
		RenderStateful(buf.Area, buf, &state)

	assertLines(t, buf, []string{"<#->"})
	assertCellStyle(t, buf, 0, 0, beginStyle)
	assertCellStyle(t, buf, 1, 0, thumbStyle)
	assertCellStyle(t, buf, 2, 0, trackStyle)
	assertCellStyle(t, buf, 3, 0, endStyle)
}

func scrollbarNoArrows() widgets.Scrollbar {
	return widgets.NewScrollbar(widgets.ScrollbarOrientationHorizontalTop).
		ClearBeginSymbol().
		ClearEndSymbol().
		TrackSymbol("-").
		ThumbSymbol("#")
}
