package main

import (
	"log"
	"sort"
	"strings"

	gatuitcell "github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/terminal"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
	tcell "github.com/gdamore/tcell/v2"
)

const itemHeight = 4

var infoText = []string{
	"(Esc/q) quit | (Up/k) move up | (Down/j) move down | (Left/h) move left | (Right/l) move right",
	"(Shift+Right/L) next color | (Shift+Left/H) previous color",
}

type palette struct {
	headerBG style.Color
	accent   style.Color
	cell     style.Color
}

var palettes = []palette{
	{headerBG: style.RGBColor(30, 58, 138), accent: style.RGBColor(96, 165, 250), cell: style.RGBColor(37, 99, 235)},
	{headerBG: style.RGBColor(6, 78, 59), accent: style.RGBColor(52, 211, 153), cell: style.RGBColor(5, 150, 105)},
	{headerBG: style.RGBColor(49, 46, 129), accent: style.RGBColor(129, 140, 248), cell: style.RGBColor(79, 70, 229)},
	{headerBG: style.RGBColor(127, 29, 29), accent: style.RGBColor(248, 113, 113), cell: style.RGBColor(220, 38, 38)},
}

type tableColors struct {
	bufferBG              style.Color
	headerBG              style.Color
	headerFG              style.Color
	rowFG                 style.Color
	selectedRowStyleFG    style.Color
	selectedColumnStyleFG style.Color
	selectedCellStyleFG   style.Color
	normalRowColor        style.Color
	altRowColor           style.Color
	footerBorderColor     style.Color
}

func newTableColors(p palette) tableColors {
	return tableColors{
		bufferBG:              style.RGBColor(2, 6, 23),
		headerBG:              p.headerBG,
		headerFG:              style.RGBColor(226, 232, 240),
		rowFG:                 style.RGBColor(226, 232, 240),
		selectedRowStyleFG:    p.accent,
		selectedColumnStyleFG: p.accent,
		selectedCellStyleFG:   p.cell,
		normalRowColor:        style.RGBColor(2, 6, 23),
		altRowColor:           style.RGBColor(15, 23, 42),
		footerBorderColor:     p.accent,
	}
}

type data struct {
	name    string
	address string
	email   string
}

func (d data) values() [3]string {
	return [3]string{d.name, d.address, d.email}
}

type app struct {
	state           widgets.TableState
	items           []data
	longestItemLens [3]int
	scrollState     widgets.ScrollbarState
	colors          tableColors
	colorIndex      int
}

func newApp() app {
	items := generateNames()
	return app{
		state:           widgets.NewTableState().WithSelected(0),
		items:           items,
		longestItemLens: constraintLens(items),
		scrollState:     widgets.NewScrollbarState((len(items) - 1) * itemHeight),
		colors:          newTableColors(palettes[0]),
	}
}

func (a *app) nextRow() {
	selected, ok := a.state.Selected()
	if !ok || selected >= len(a.items)-1 {
		selected = 0
	} else {
		selected++
	}
	a.state.Select(selected)
	a.scrollState = a.scrollState.Position(selected * itemHeight)
}

func (a *app) previousRow() {
	selected, ok := a.state.Selected()
	if !ok || selected == 0 {
		selected = len(a.items) - 1
	} else {
		selected--
	}
	a.state.Select(selected)
	a.scrollState = a.scrollState.Position(selected * itemHeight)
}

func (a *app) nextColumn() {
	a.state.SelectNextColumn()
}

func (a *app) previousColumn() {
	a.state.SelectPreviousColumn()
}

func (a *app) nextColor() {
	a.colorIndex = (a.colorIndex + 1) % len(palettes)
}

func (a *app) previousColor() {
	a.colorIndex = (a.colorIndex + len(palettes) - 1) % len(palettes)
}

func (a *app) setColors() {
	a.colors = newTableColors(palettes[a.colorIndex])
}

func (a *app) run(term *terminal.Terminal, screen tcell.Screen) error {
	for {
		if _, err := term.Draw(a.render); err != nil {
			return err
		}

		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if a.handleKey(ev) {
				return nil
			}
		}
	}
}

func (a *app) handleKey(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		return true
	case tcell.KeyUp:
		a.previousRow()
	case tcell.KeyDown:
		a.nextRow()
	case tcell.KeyLeft:
		if ev.Modifiers()&tcell.ModShift != 0 {
			a.previousColor()
		} else {
			a.previousColumn()
		}
	case tcell.KeyRight:
		if ev.Modifiers()&tcell.ModShift != 0 {
			a.nextColor()
		} else {
			a.nextColumn()
		}
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q':
			return true
		case 'j':
			a.nextRow()
		case 'k':
			a.previousRow()
		case 'h':
			a.previousColumn()
		case 'l':
			a.nextColumn()
		case 'H':
			a.previousColor()
		case 'L':
			a.nextColor()
		}
	}
	return false
}

func (a *app) render(frame *terminal.Frame) {
	rects := layout.NewVerticalLayout(layout.Min(5), layout.Length(4)).Split(frame.Area())
	a.setColors()
	a.renderTable(frame, rects[0])
	a.renderScrollbar(frame, rects[0])
	a.renderFooter(frame, rects[1])
}

func (a *app) renderTable(frame *terminal.Frame, area layout.Rect) {
	headerStyle := style.NewStyle().Fg(a.colors.headerFG).Bg(a.colors.headerBG)
	selectedRowStyle := style.NewStyle().AddModifier(style.ModifierReversed).Fg(a.colors.selectedRowStyleFG)
	selectedColumnStyle := style.NewStyle().Fg(a.colors.selectedColumnStyleFG)
	selectedCellStyle := style.NewStyle().AddModifier(style.ModifierReversed).Fg(a.colors.selectedCellStyleFG)

	header := widgets.TableRowFromStrings([]string{"Name", "Address", "Email"}).Style(headerStyle).Height(1)
	rows := make([]widgets.TableRow, 0, len(a.items))
	for i, item := range a.items {
		rowBG := a.colors.normalRowColor
		if i%2 == 1 {
			rowBG = a.colors.altRowColor
		}
		rows = append(rows, a.rowForItem(i, item).Style(style.NewStyle().Fg(a.colors.rowFG).Bg(rowBG)).Height(itemHeight))
	}

	table := widgets.NewTable(rows, []layout.Constraint{
		layout.Length(a.longestItemLens[0] + 1),
		layout.Min(a.longestItemLens[1] + 1),
		layout.Min(a.longestItemLens[2]),
	}).
		Header(header).
		RowHighlightStyle(selectedRowStyle).
		ColumnHighlightStyle(selectedColumnStyle).
		CellHighlightStyle(selectedCellStyle).
		HighlightSymbol(" > ").
		Bg(a.colors.bufferBG).
		HighlightSpacing(widgets.HighlightSpacingAlways)

	frame.RenderStatefulWidget(table, area, &a.state)
}

func (a *app) rowForItem(index int, item data) widgets.TableRow {
	values := item.values()
	cells := make([]widgets.TableCell, 0, len(values))
	for column, value := range values {
		if index == 3 && column == 1 {
			cells = append(cells, widgets.NewTableCell(text.FromString("\n[no address or email address is available for this person]\n")).ColumnSpan(2))
			break
		}
		cells = append(cells, widgets.NewTableCell(text.FromString("\n"+value+"\n")))
	}
	return widgets.NewTableRow(cells)
}

func (a *app) renderScrollbar(frame *terminal.Frame, area layout.Rect) {
	scrollbar := widgets.NewScrollbar(widgets.ScrollbarOrientationVerticalRight).
		ClearBeginSymbol().
		ClearEndSymbol()
	frame.RenderStatefulWidget(scrollbar, area.Inner(layout.NewMargin(1, 1)), &a.scrollState)
}

func (a *app) renderFooter(frame *terminal.Frame, area layout.Rect) {
	lines := make([]text.Line, 0, len(infoText))
	for _, line := range infoText {
		lines = append(lines, text.LineFromString(line))
	}
	footer := widgets.NewParagraph(text.NewText(lines...).Center()).
		Style(style.NewStyle().Fg(a.colors.rowFG).Bg(a.colors.bufferBG)).
		Block(widgets.BorderedBlock().
			BorderType(widgets.BorderTypeDouble).
			BorderStyle(style.NewStyle().Fg(a.colors.footerBorderColor)))
	frame.RenderWidget(footer, area)
}

func generateNames() []data {
	items := []data{
		{name: "Ada Lovelace", address: "12 St. James Square\nLondon, UK SW1Y 4LB", email: "ada.lovelace@example.com"},
		{name: "Alan Turing", address: "48 Science Road\nCambridge, UK CB2 1TN", email: "alan.turing@example.com"},
		{name: "Barbara Liskov", address: "77 Kendall Street\nCambridge, MA 02142", email: "barbara.liskov@example.com"},
		{name: "Donald Knuth", address: "1930 University Avenue\nStanford, CA 94305", email: "donald.knuth@example.com"},
		{name: "Edsger Dijkstra", address: "22 Algorithm Lane\nAustin, TX 78712", email: "edsger.dijkstra@example.com"},
		{name: "Frances Allen", address: "5 Compiler Court\nPoughkeepsie, NY 12601", email: "frances.allen@example.com"},
		{name: "Grace Hopper", address: "9 Navy Yard Plaza\nArlington, VA 22202", email: "grace.hopper@example.com"},
		{name: "Hedy Lamarr", address: "101 Frequency Way\nLos Angeles, CA 90028", email: "hedy.lamarr@example.com"},
		{name: "John Backus", address: "64 Fortran Drive\nSan Jose, CA 95112", email: "john.backus@example.com"},
		{name: "Katherine Johnson", address: "37 Orbital Avenue\nHampton, VA 23666", email: "katherine.johnson@example.com"},
		{name: "Ken Thompson", address: "1 Unix Place\nBerkeley, CA 94704", email: "ken.thompson@example.com"},
		{name: "Leslie Lamport", address: "42 Distributed Road\nRedmond, WA 98052", email: "leslie.lamport@example.com"},
		{name: "Margaret Hamilton", address: "11 Apollo Street\nBoston, MA 02108", email: "margaret.hamilton@example.com"},
		{name: "Mary Wilkes", address: "8 Console Circle\nBaltimore, MD 21201", email: "mary.wilkes@example.com"},
		{name: "Radia Perlman", address: "26 Spanning Tree Blvd\nSomerville, MA 02144", email: "radia.perlman@example.com"},
		{name: "Tim Berners-Lee", address: "80 Hypertext Road\nGeneva, CH 1211", email: "tim.berners-lee@example.com"},
		{name: "Ursula Burns", address: "70 Xerox Park\nRochester, NY 14644", email: "ursula.burns@example.com"},
		{name: "Vint Cerf", address: "33 Internet Avenue\nReston, VA 20190", email: "vint.cerf@example.com"},
		{name: "Whitfield Diffie", address: "19 Crypto Lane\nPalo Alto, CA 94301", email: "whitfield.diffie@example.com"},
		{name: "Yukihiro Matsumoto", address: "2 Ruby Street\nMatsue, JP 690-0000", email: "yukihiro.matsumoto@example.com"},
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].name < items[j].name
	})
	return items
}

func constraintLens(items []data) [3]int {
	var lens [3]int
	for _, item := range items {
		lens[0] = max(lens[0], buffer.CellWidth(item.name))
		for _, line := range strings.Split(item.address, "\n") {
			lens[1] = max(lens[1], buffer.CellWidth(line))
		}
		lens[2] = max(lens[2], buffer.CellWidth(item.email))
	}
	return lens
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	backend, err := gatuitcell.NewWithScreen(screen)
	if err != nil {
		log.Fatal(err)
	}
	defer backend.Close()

	term, err := terminal.New(backend)
	if err != nil {
		log.Fatal(err)
	}
	app := newApp()
	if err := app.run(term, screen); err != nil {
		log.Fatal(err)
	}
}
