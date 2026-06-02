package widgets

import (
	"math"
	"strings"
	"testing"

	"gatui/buffer"
	"gatui/layout"
	"gatui/style"
	"gatui/text"
)

func TestBlock_new(t *testing.T) {
	block := NewBlock()

	if block.borders != NoBorders {
		t.Fatalf("NewBlock().borders = %v, want %v", block.borders, NoBorders)
	}
	if block.padding != PaddingZero() {
		t.Fatalf("NewBlock().padding = %#v, want %#v", block.padding, PaddingZero())
	}
	if block.style != style.NewStyle() {
		t.Fatalf("NewBlock().style = %#v, want %#v", block.style, style.NewStyle())
	}

	bordered := BorderedBlock()
	if bordered.borders != AllBorders {
		t.Fatalf("BorderedBlock().borders = %v, want %v", bordered.borders, AllBorders)
	}
}

func TestBlock_titleStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 7, 3))
	block := BorderedBlock().
		Title(text.LineFromString("Title")).
		TitleStyle(style.NewStyle().Fg(style.Red))

	block.Render(buf.Area, buf)

	for x := 1; x <= 5; x++ {
		cell, ok := buf.CellAt(x, 0)
		if !ok {
			t.Fatalf("missing cell at (%d,0)", x)
		}
		if cell.Style != style.NewStyle().Fg(style.Red) {
			t.Fatalf("cell(%d,0).Style = %#v, want red foreground", x, cell.Style)
		}
	}
}

func TestBlock_titleBorderStyle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	block := BorderedBlock().
		Title(text.LineFromString("test")).
		BorderStyle(style.NewStyle().Fg(style.Yellow))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"┌test────┐",
		"│        │",
		"└────────┘",
	})
	for x := range 10 {
		assertBlockCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.Yellow))
		assertBlockCellStyle(t, buf, x, 2, style.NewStyle().Fg(style.Yellow))
	}
	assertBlockCellStyle(t, buf, 0, 1, style.NewStyle().Fg(style.Yellow))
	assertBlockCellStyle(t, buf, 9, 1, style.NewStyle().Fg(style.Yellow))
	for x := 1; x <= 8; x++ {
		assertBlockCellStyle(t, buf, x, 1, style.NewStyle())
	}
}

func TestBlock_renderCustomBorderSet(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))
	block := BorderedBlock().BorderSet(BorderSet{
		TopLeft:          "1",
		TopRight:         "2",
		BottomLeft:       "3",
		BottomRight:      "4",
		VerticalLeft:     "L",
		VerticalRight:    "R",
		HorizontalTop:    "T",
		HorizontalBottom: "B",
	})

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"1TTTTTTTT2",
		"L        R",
		"3BBBBBBBB4",
	})
}

func TestBlock_renderRoundedBorder(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))

	BorderedBlock().BorderSet(RoundedBorderSet).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"╭────────╮",
		"│        │",
		"╰────────╯",
	})
}

func TestBlock_renderDoubleBorder(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))

	BorderedBlock().BorderSet(DoubleBorderSet).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"╔════════╗",
		"║        ║",
		"╚════════╝",
	})
}

func TestBlock_renderSolidBorder(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))

	BorderedBlock().BorderSet(SolidBorderSet).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"┏━━━━━━━━┓",
		"┃        ┃",
		"┗━━━━━━━━┛",
	})
}

func TestMergeStrategy_merge(t *testing.T) {
	tests := []struct {
		name     string
		strategy MergeStrategy
		prev     string
		next     string
		want     string
	}{
		{name: "replace", strategy: MergeStrategyReplace, prev: "│", next: "━", want: "━"},
		{name: "exact plain cross", strategy: MergeStrategyExact, prev: "│", next: "─", want: "┼"},
		{name: "exact mixed plain thick cross", strategy: MergeStrategyExact, prev: "│", next: "━", want: "┿"},
		{name: "exact falls back to replace", strategy: MergeStrategyExact, prev: "┘", next: "╔", want: "╔"},
		{name: "fuzzy double thick", strategy: MergeStrategyFuzzy, prev: "┘", next: "╔", want: "╬"},
		{name: "fuzzy rounded plain", strategy: MergeStrategyFuzzy, prev: "┘", next: "╭", want: "┼"},
		{name: "fuzzy dashed thick", strategy: MergeStrategyFuzzy, prev: "╎", next: "╍", want: "┿"},
		{name: "non border previous wins when next is border", strategy: MergeStrategyExact, prev: "a", next: "╭", want: "a"},
		{name: "non border next wins", strategy: MergeStrategyExact, prev: "┌", next: "a", want: "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeBorderSymbols(tt.strategy, tt.prev, tt.next); got != tt.want {
				t.Fatalf("mergeBorderSymbols(%v, %q, %q) = %q, want %q", tt.strategy, tt.prev, tt.next, got, tt.want)
			}
		})
	}
}

func TestBlock_renderMergedBordersInMinimalBufferDoesNotPanic(t *testing.T) {
	for _, strategy := range []MergeStrategy{MergeStrategyExact, MergeStrategyFuzzy} {
		t.Run(strategy.String(), func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

			BorderedBlock().MergeBorders(strategy).Render(buf.Area, buf)

			assertBlockLines(t, buf, []string{"┼"})
		})
	}
}

func TestBlock_renderMergedBorders(t *testing.T) {
	tests := []struct {
		name     string
		strategy MergeStrategy
		first    Block
		firstAt  layout.Rect
		second   Block
		secondAt layout.Rect
		want     []string
	}{
		{
			name:     "replace touching corners",
			strategy: MergeStrategyReplace,
			first:    BorderedBlock(),
			firstAt:  layout.NewRect(0, 0, 5, 5),
			second:   BorderedBlock().BorderType(BorderTypeThick),
			secondAt: layout.NewRect(4, 4, 5, 5),
			want: []string{
				"┌───┐    ",
				"│   │    ",
				"│   │    ",
				"│   │    ",
				"└───┏━━━┓",
				"    ┃   ┃",
				"    ┃   ┃",
				"    ┃   ┃",
				"    ┗━━━┛",
			},
		},
		{
			name:     "exact overlapping rectangles",
			strategy: MergeStrategyExact,
			first:    BorderedBlock(),
			firstAt:  layout.NewRect(0, 0, 5, 5),
			second:   BorderedBlock().BorderType(BorderTypeThick),
			secondAt: layout.NewRect(2, 2, 5, 5),
			want: []string{
				"┌───┐    ",
				"│   │    ",
				"│ ┏━┿━┓  ",
				"│ ┃ │ ┃  ",
				"└─╂─┘ ┃  ",
				"  ┃   ┃  ",
				"  ┗━━━┛  ",
			},
		},
		{
			name:     "fuzzy touching vertical edges",
			strategy: MergeStrategyFuzzy,
			first:    BorderedBlock().BorderType(BorderTypeRounded),
			firstAt:  layout.NewRect(0, 0, 5, 5),
			second:   BorderedBlock(),
			secondAt: layout.NewRect(4, 0, 5, 5),
			want: []string{
				"╭───┬───┐",
				"│   │   │",
				"│   │   │",
				"│   │   │",
				"╰───┴───┘",
			},
		},
		{
			name:     "fuzzy touching horizontal edges",
			strategy: MergeStrategyFuzzy,
			first:    BorderedBlock().BorderType(BorderTypeLightDoubleDashed),
			firstAt:  layout.NewRect(0, 0, 5, 5),
			second:   BorderedBlock().BorderType(BorderTypeHeavyDoubleDashed),
			secondAt: layout.NewRect(0, 4, 5, 5),
			want: []string{
				"┌╌╌╌┐    ",
				"╎   ╎    ",
				"╎   ╎    ",
				"╎   ╎    ",
				"┢╍╍╍┪    ",
				"╏   ╏    ",
				"╏   ╏    ",
				"╏   ╏    ",
				"┗╍╍╍┛    ",
			},
		},
		{
			name:     "exact double dashed falls back where unrepresentable",
			strategy: MergeStrategyExact,
			first:    BorderedBlock(),
			firstAt:  layout.NewRect(0, 0, 5, 5),
			second:   BorderedBlock().BorderType(BorderTypeDouble),
			secondAt: layout.NewRect(4, 0, 5, 5),
			want: []string{
				"┌───╔═══╗",
				"│   ║   ║",
				"│   ║   ║",
				"│   ║   ║",
				"└───╚═══╝",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 9, len(tt.want)))

			tt.first.MergeBorders(tt.strategy).Render(tt.firstAt, buf)
			tt.second.MergeBorders(tt.strategy).Render(tt.secondAt, buf)

			assertBlockLines(t, buf, tt.want)
		})
	}
}

func TestBorderType_string(t *testing.T) {
	tests := []struct {
		borderType BorderType
		want       string
	}{
		{BorderTypePlain, "Plain"},
		{BorderTypeRounded, "Rounded"},
		{BorderTypeDouble, "Double"},
		{BorderTypeThick, "Thick"},
		{BorderTypeLightDoubleDashed, "LightDoubleDashed"},
		{BorderTypeHeavyDoubleDashed, "HeavyDoubleDashed"},
		{BorderTypeLightTripleDashed, "LightTripleDashed"},
		{BorderTypeHeavyTripleDashed, "HeavyTripleDashed"},
		{BorderTypeLightQuadrupleDashed, "LightQuadrupleDashed"},
		{BorderTypeHeavyQuadrupleDashed, "HeavyQuadrupleDashed"},
		{BorderTypeQuadrantInside, "QuadrantInside"},
		{BorderTypeQuadrantOutside, "QuadrantOutside"},
	}

	for _, tt := range tests {
		if got := tt.borderType.String(); got != tt.want {
			t.Fatalf("%#v.String() = %q, want %q", tt.borderType, got, tt.want)
		}
	}
}

func TestParseBorderType(t *testing.T) {
	tests := []struct {
		value string
		want  BorderType
	}{
		{"Plain", BorderTypePlain},
		{"Rounded", BorderTypeRounded},
		{"Double", BorderTypeDouble},
		{"Thick", BorderTypeThick},
		{"LightDoubleDashed", BorderTypeLightDoubleDashed},
		{"HeavyDoubleDashed", BorderTypeHeavyDoubleDashed},
		{"LightTripleDashed", BorderTypeLightTripleDashed},
		{"HeavyTripleDashed", BorderTypeHeavyTripleDashed},
		{"LightQuadrupleDashed", BorderTypeLightQuadrupleDashed},
		{"HeavyQuadrupleDashed", BorderTypeHeavyQuadrupleDashed},
		{"QuadrantInside", BorderTypeQuadrantInside},
		{"QuadrantOutside", BorderTypeQuadrantOutside},
	}

	for _, tt := range tests {
		got, err := ParseBorderType(tt.value)
		if err != nil {
			t.Fatalf("ParseBorderType(%q) returned unexpected error: %v", tt.value, err)
		}
		if got != tt.want {
			t.Fatalf("ParseBorderType(%q) = %#v, want %#v", tt.value, got, tt.want)
		}
	}

	if _, err := ParseBorderType(""); err == nil {
		t.Fatal("ParseBorderType(\"\") returned nil error, want error")
	}
}

func TestBlock_renderBorderTypeDashedAndQuadrantBorders(t *testing.T) {
	tests := []struct {
		name       string
		borderType BorderType
		want       []string
	}{
		{
			name:       "light double dashed",
			borderType: BorderTypeLightDoubleDashed,
			want: []string{
				"┌╌╌╌╌╌╌╌╌┐",
				"╎        ╎",
				"└╌╌╌╌╌╌╌╌┘",
			},
		},
		{
			name:       "heavy double dashed",
			borderType: BorderTypeHeavyDoubleDashed,
			want: []string{
				"┏╍╍╍╍╍╍╍╍┓",
				"╏        ╏",
				"┗╍╍╍╍╍╍╍╍┛",
			},
		},
		{
			name:       "light triple dashed",
			borderType: BorderTypeLightTripleDashed,
			want: []string{
				"┌┄┄┄┄┄┄┄┄┐",
				"┆        ┆",
				"└┄┄┄┄┄┄┄┄┘",
			},
		},
		{
			name:       "heavy triple dashed",
			borderType: BorderTypeHeavyTripleDashed,
			want: []string{
				"┏┅┅┅┅┅┅┅┅┓",
				"┇        ┇",
				"┗┅┅┅┅┅┅┅┅┛",
			},
		},
		{
			name:       "light quadruple dashed",
			borderType: BorderTypeLightQuadrupleDashed,
			want: []string{
				"┌┈┈┈┈┈┈┈┈┐",
				"┊        ┊",
				"└┈┈┈┈┈┈┈┈┘",
			},
		},
		{
			name:       "heavy quadruple dashed",
			borderType: BorderTypeHeavyQuadrupleDashed,
			want: []string{
				"┏┉┉┉┉┉┉┉┉┓",
				"┋        ┋",
				"┗┉┉┉┉┉┉┉┉┛",
			},
		},
		{
			name:       "quadrant inside",
			borderType: BorderTypeQuadrantInside,
			want: []string{
				"▗▄▄▄▄▄▄▄▄▖",
				"▐        ▌",
				"▝▀▀▀▀▀▀▀▀▘",
			},
		},
		{
			name:       "quadrant outside",
			borderType: BorderTypeQuadrantOutside,
			want: []string{
				"▛▀▀▀▀▀▀▀▀▜",
				"▌        ▐",
				"▙▄▄▄▄▄▄▄▄▟",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))

			BorderedBlock().BorderType(tt.borderType).Render(buf.Area, buf)

			assertBlockLines(t, buf, tt.want)
		})
	}
}

func TestBlock_borderTypeOverwritesBorderSet(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))

	BorderedBlock().BorderSet(DoubleBorderSet).BorderType(BorderTypeRounded).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"╭────────╮",
		"│        │",
		"╰────────╯",
	})
}

func TestBlock_renderPartialBorders(t *testing.T) {
	tests := []struct {
		name    string
		borders Borders
		want    []string
	}{
		{
			name:    "all",
			borders: TopBorder | LeftBorder | RightBorder | BottomBorder,
			want: []string{
				"┌────────┐",
				"│        │",
				"└────────┘",
			},
		},
		{
			name:    "top left",
			borders: TopBorder | LeftBorder,
			want: []string{
				"┌─────────",
				"│         ",
				"│         ",
			},
		},
		{
			name:    "top right",
			borders: TopBorder | RightBorder,
			want: []string{
				"─────────┐",
				"         │",
				"         │",
			},
		},
		{
			name:    "bottom left",
			borders: BottomBorder | LeftBorder,
			want: []string{
				"│         ",
				"│         ",
				"└─────────",
			},
		},
		{
			name:    "bottom right",
			borders: BottomBorder | RightBorder,
			want: []string{
				"         │",
				"         │",
				"─────────┘",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(layout.NewRect(0, 0, 10, 3))

			NewBlock().Borders(tt.borders).Render(buf.Area, buf)

			assertBlockLines(t, buf, tt.want)
		})
	}
}

func TestBlock_styleIntoWorksFromUserView(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 3))
	block := BorderedBlock().
		Style(style.NewStyle().Bg(style.Green)).
		BorderStyle(style.NewStyle().Fg(style.Cyan)).
		Title(text.LineFromString("T")).
		TitleStyle(style.NewStyle().Fg(style.Red))

	block.Render(buf.Area, buf)

	assertBlockCellStyle(t, buf, 0, 0, style.NewStyle().Fg(style.Cyan).Bg(style.Green))
	assertBlockCellStyle(t, buf, 1, 0, style.NewStyle().Fg(style.Red).Bg(style.Green))
	assertBlockCellStyle(t, buf, 1, 1, style.NewStyle())
}

func TestBlock_leftTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))

	NewBlock().Title(text.LineFromString("L12")).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L12       "})
}

func TestBlock_leftTitleTruncated(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))

	NewBlock().Title(text.LineFromString("L1234567890")).Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L123456789"})
}

func TestBlock_centerTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Center).
		Title(text.LineFromString("C12"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"   C12    "})
}

func TestBlock_centerTitleTruncated(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Center).
		Title(text.LineFromString("C1234567890"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"C123456789"})
}

func TestBlock_centerTitleTruncatesLeftTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		Title(text.LineFromString("L1234")).
		Title(text.LineFromString("C5678").Center())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L1C5678   "})
}

func TestBlock_rightTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Right).
		Title(text.LineFromString("R12"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"       R12"})
}

func TestBlock_rightTitleTruncated(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		TitleAlignment(layout.Right).
		Title(text.LineFromString("R1234567890"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"R123456789"})
}

func TestBlock_rightTitleTruncatesLeftTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		Title(text.LineFromString("L12345")).
		Title(text.LineFromString("R67890").Right())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"L123R67890"})
}

func TestBlock_rightTitleTruncatesCenterTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 1))
	block := NewBlock().
		Title(text.LineFromString("C12345").Center()).
		Title(text.LineFromString("R67890").Right())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{"  C1R67890"})
}

func TestBlock_titleAlignmentOverridesBlockTitleAlignment(t *testing.T) {
	tests := []struct {
		blockAlignment layout.Alignment
		lineAlignment  func(text.Line) text.Line
		want           string
	}{
		{blockAlignment: layout.Right, lineAlignment: func(line text.Line) text.Line { return line.Left() }, want: "test    "},
		{blockAlignment: layout.Left, lineAlignment: func(line text.Line) text.Line { return line.Center() }, want: "  test  "},
		{blockAlignment: layout.Center, lineAlignment: func(line text.Line) text.Line { return line.Right() }, want: "    test"},
	}
	for _, tt := range tests {
		buf := buffer.Empty(layout.NewRect(0, 0, 8, 1))
		block := NewBlock().
			TitleAlignment(tt.blockAlignment).
			Title(tt.lineAlignment(text.LineFromString("test")))

		block.Render(buf.Area, buf)

		assertBlockLines(t, buf, []string{tt.want})
	}
}

func TestBlock_titleStyleOverridesBlockTitleStyle(t *testing.T) {
	for _, alignment := range []layout.Alignment{layout.Left, layout.Center, layout.Right} {
		buf := buffer.Empty(layout.NewRect(0, 0, 4, 1))
		block := NewBlock().
			TitleAlignment(alignment).
			TitleStyle(style.NewStyle().Fg(style.Green).Bg(style.Red)).
			Title(text.LineFromString("test").Fg(style.Yellow))

		block.Render(buf.Area, buf)

		for x := range 4 {
			assertBlockCellStyle(t, buf, x, 0, style.NewStyle().Fg(style.Yellow).Bg(style.Red))
		}
	}
}

func TestBlock_titlePosition(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 4, 2))
	block := NewBlock().
		TitlePosition(TitlePositionBottom).
		Title(text.LineFromString("test"))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"    ",
		"test",
	})
}

func TestBlock_titleTopBottom(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 11, 3))
	block := BorderedBlock().
		TitleTop(text.LineFromString("A").Left()).
		TitleTop(text.LineFromString("B").Center()).
		TitleTop(text.LineFromString("C").Right()).
		TitleBottom(text.LineFromString("D").Left()).
		TitleBottom(text.LineFromString("E").Center()).
		TitleBottom(text.LineFromString("F").Right())

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"┌A───B───C┐",
		"│         │",
		"└D───E───F┘",
	})
}

func TestWidgetsBlock_rendersOnSmallAreas(t *testing.T) {
	block := BorderedBlock().
		Padding(PaddingUniform(1)).
		TitleTop(text.LineFromString("Top")).
		TitleBottom(text.LineFromString("Bottom")).
		Shadow(NewShadowBlock())

	tests := []struct {
		name string
		area layout.Rect
		want []string
	}{
		{name: "0x0", area: layout.NewRect(0, 0, 0, 0), want: nil},
		{name: "1x0", area: layout.NewRect(0, 0, 1, 0), want: nil},
		{name: "0x1", area: layout.NewRect(0, 0, 0, 1), want: []string{""}},
		{name: "1x1", area: layout.NewRect(0, 0, 1, 1), want: []string{"┌"}},
		{name: "2x1", area: layout.NewRect(0, 0, 2, 1), want: []string{"┌┐"}},
		{name: "1x2", area: layout.NewRect(0, 0, 1, 2), want: []string{"┌", "└"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := buffer.Empty(tt.area)

			block.Render(buf.Area, buf)

			assertBlockLines(t, buf, tt.want)
		})
	}

	t.Run("shadow clipped in larger buffer", func(t *testing.T) {
		buf := buffer.Empty(layout.NewRect(0, 0, 3, 3))

		block.Render(layout.NewRect(0, 0, 1, 1), buf)

		assertBlockLines(t, buf, []string{
			"┌  ",
			" █ ",
			"   ",
		})
	})
}

func TestWidgetsBlock_titlesOverlap(t *testing.T) {
	block := NewBlock().
		TitleTop(text.LineFromString("left").Left()).
		TitleTop(text.LineFromString("CENTER").Center()).
		TitleTop(text.LineFromString("right").Right()).
		TitleBottom(text.LineFromString("bottom-left").Left()).
		TitleBottom(text.LineFromString("BOTTOM").Center()).
		TitleBottom(text.LineFromString("bottom-right").Right())
	buf := buffer.Empty(layout.NewRect(0, 0, 10, 2))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"leCENright",
		"ttom-right",
	})
	assertBlockCellStyle(t, buf, 2, 0, style.NewStyle())
	assertBlockCellStyle(t, buf, 6, 0, style.NewStyle())
}

func TestBlock_hasTitleAtPosition(t *testing.T) {
	block := NewBlock()
	if block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("empty block has top title")
	}
	if block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("empty block has bottom title")
	}

	block = NewBlock().TitleTop(text.LineFromString("test"))
	if !block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("TitleTop block missing top title")
	}
	if block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("TitleTop block has bottom title")
	}

	block = NewBlock().TitleBottom(text.LineFromString("test"))
	if block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("TitleBottom block has top title")
	}
	if !block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("TitleBottom block missing bottom title")
	}

	block = NewBlock().
		TitleTop(text.LineFromString("test")).
		TitleBottom(text.LineFromString("test"))
	if !block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("mixed block missing top title")
	}
	if !block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("mixed block missing bottom title")
	}

	block = NewBlock().
		Title(text.LineFromString("top")).
		TitlePosition(TitlePositionBottom).
		Title(text.LineFromString("bottom"))
	if !block.hasTitleAtPosition(TitlePositionTop) {
		t.Fatalf("default-position block missing top title")
	}
	if !block.hasTitleAtPosition(TitlePositionBottom) {
		t.Fatalf("default-position block missing bottom title")
	}
}

func TestBlock_titlesAreaHandlesEmptyAreaWithoutPanicking(t *testing.T) {
	block := NewBlock()

	got := block.titlesArea(layout.NewRect(0, 0, 0, 0), TitlePositionBottom)

	if got != layout.NewRect(0, 0, 0, 1) {
		t.Fatalf("titlesArea = %#v, want %#v", got, layout.NewRect(0, 0, 0, 1))
	}
}

func TestBlock_titlesAreaSaturatesWhenLeftBorderOffsetOverflows(t *testing.T) {
	block := NewBlock().Borders(LeftBorder)

	got := block.titlesArea(layout.NewRect(layout.MaxCoordinate, 0, 1, 1), TitlePositionTop)

	if got != layout.NewRect(layout.MaxCoordinate, 0, 1, 1) {
		t.Fatalf("titlesArea = %#v, want %#v", got, layout.NewRect(layout.MaxCoordinate, 0, 1, 1))
	}
}

func TestBlock_renderRightAlignedEmptyTitle(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 15, 3))
	block := NewBlock().
		TitleAlignment(layout.Right).
		Title(text.LineFromString(""))

	block.Render(buf.Area, buf)

	assertBlockLines(t, buf, []string{
		"               ",
		"               ",
		"               ",
	})
}

func TestBlock_renderCenterTitlesHandlesTitleWidthIncrementOverflow(t *testing.T) {
	block := NewBlock().Title(text.LineFromString(strings.Repeat("a", layout.MaxCoordinate)).Center())
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	block.renderCenterTitles(TitlePositionTop, layout.NewRect(0, 0, 1, 1), buf)

	assertBlockLines(t, buf, []string{" "})
}

func TestBlock_renderCenterTitlesHandlesTotalWidthOverflow(t *testing.T) {
	block := NewBlock().
		Title(text.LineFromString(strings.Repeat("a", 40_000)).Center()).
		Title(text.LineFromString(strings.Repeat("b", 30_000)).Center())
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	block.renderCenterTitles(TitlePositionTop, layout.NewRect(0, 0, 1, 1), buf)

	assertBlockLines(t, buf, []string{" "})
}

func TestBlock_renderCenteredTitlesWithTruncationHandlesTitleAdvanceOverflow(t *testing.T) {
	block := NewBlock().
		Title(text.LineFromString(strings.Repeat("a", layout.MaxCoordinate)).Center()).
		Title(text.LineFromString("b").Center())
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	block.renderCenterTitles(TitlePositionTop, layout.NewRect(0, 0, layout.MaxCoordinate, 1), buf)

	assertBlockLines(t, buf, []string{"a"})
}

func TestBlock_renderCenteredTitlesWithoutTruncationHandlesMaximumX(t *testing.T) {
	block := NewBlock()
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	block.renderCenterTitles(TitlePositionTop, layout.NewRect(layout.MaxCoordinate-1, 0, 1, 1), buf)

	assertBlockLines(t, buf, []string{" "})
}

func TestBlock_renderCenteredTitlesWithoutTruncationHandlesTitleAdvanceOverflow(t *testing.T) {
	block := NewBlock().Title(text.LineFromString(strings.Repeat("a", layout.MaxCoordinate)).Center())
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	block.renderCenterTitles(TitlePositionTop, layout.NewRect(0, 0, layout.MaxCoordinate, 1), buf)

	assertBlockLines(t, buf, []string{"a"})
}

func TestBlock_renderLeftTitlesHandlesTitleAdvanceOverflow(t *testing.T) {
	block := NewBlock().Title(text.LineFromString(strings.Repeat("a", layout.MaxCoordinate)))
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	block.renderLeftTitles(TitlePositionTop, layout.NewRect(0, 0, 1, 1), buf)

	assertBlockLines(t, buf, []string{"a"})
}

func TestBlock_renderInMinimalBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 1, 1))

	BorderedBlock().Title(text.LineFromString("Title")).Render(buf.Area, buf)

	cell, ok := buf.CellAt(0, 0)
	if !ok {
		t.Fatalf("missing cell at (0,0)")
	}
	if cell.Symbol != "┌" {
		t.Fatalf("cell(0,0).Symbol = %q, want %q", cell.Symbol, "┌")
	}
}

func TestBlock_renderInZeroSizeBuffer(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))

	BorderedBlock().Title(text.LineFromString("Title")).Render(buf.Area, buf)
}

func TestBlock_renderCornersHandlesEmptyAreaWithoutPanicking(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))

	BorderedBlock().renderBorders(buf.Area, buf)
}

func TestBlock_renderSidesHandlesEmptyAreaWithoutPanicking(t *testing.T) {
	buf := buffer.Empty(layout.NewRect(0, 0, 0, 0))

	NewBlock().Borders(LeftBorder|RightBorder).renderBorders(buf.Area, buf)
}

func TestBlock_innerSaturatesWhenPaddingSumOverflows(t *testing.T) {
	block := BorderedBlock().Padding(NewPadding(math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt))

	inner := block.Inner(layout.NewRect(1, 2, 3, 4))

	if inner.Width != 0 || inner.Height != 0 {
		t.Fatalf("inner size = %dx%d, want 0x0", inner.Width, inner.Height)
	}
}

func TestBlock_verticalSpaceSaturatesWhenSpaceOverflows(t *testing.T) {
	block := NewBlock().
		Borders(TopBorder | BottomBorder).
		Padding(NewPadding(0, 0, math.MaxInt, math.MaxInt))

	if got := block.verticalSpace(); got != math.MaxInt {
		t.Fatalf("verticalSpace = %d, want %d", got, math.MaxInt)
	}
}

func assertBlockCellStyle(t *testing.T, buf *buffer.Buffer, x, y int, want style.Style) {
	t.Helper()
	cell, ok := buf.CellAt(x, y)
	if !ok {
		t.Fatalf("missing cell at (%d,%d)", x, y)
	}
	if cell.Style != want {
		t.Fatalf("cell(%d,%d).Style = %#v, want %#v", x, y, cell.Style, want)
	}
}

func assertBlockLines(t *testing.T, buf *buffer.Buffer, want []string) {
	t.Helper()
	if got := buf.Lines(); !equalStringSlices(got, want) {
		t.Fatalf("buffer lines = %#v, want %#v", got, want)
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
