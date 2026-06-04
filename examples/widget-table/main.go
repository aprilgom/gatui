package main

import (
	"fmt"
	"os"

	gatuitcell "github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/terminal"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
	tcell "github.com/gdamore/tcell/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("create tcell screen: %w", err)
	}

	backend, err := gatuitcell.NewWithScreen(screen)
	if err != nil {
		return fmt.Errorf("initialize tcell backend: %w", err)
	}
	defer backend.Close()

	term, err := terminal.New(backend)
	if err != nil {
		return fmt.Errorf("initialize terminal: %w", err)
	}

	tableState := widgets.NewTableState()
	tableState.SelectFirst()
	tableState.SelectFirstColumn()

	draw := func() error {
		_, err := term.Draw(func(frame *terminal.Frame) {
			render(frame, &tableState)
		})
		return err
	}

	if err := draw(); err != nil {
		return fmt.Errorf("draw initial frame: %w", err)
	}

	for {
		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyCtrlC, tcell.KeyEsc:
				return nil
			case tcell.KeyDown:
				tableState.SelectNext()
			case tcell.KeyUp:
				tableState.SelectPrevious()
			case tcell.KeyRight:
				tableState.SelectNextColumn()
			case tcell.KeyLeft:
				tableState.SelectPreviousColumn()
			default:
				switch event.Rune() {
				case 'q':
					return nil
				case 'j':
					tableState.SelectNext()
				case 'k':
					tableState.SelectPrevious()
				case 'l':
					tableState.SelectNextColumn()
				case 'h':
					tableState.SelectPreviousColumn()
				case 'g':
					tableState.SelectFirst()
				case 'G':
					tableState.SelectLast()
				}
			}
		}

		if err := draw(); err != nil {
			return fmt.Errorf("draw frame: %w", err)
		}
	}
}

func render(frame *terminal.Frame, tableState *widgets.TableState) {
	areas := layout.NewVerticalLayout(layout.Length(1), layout.Fill(1)).
		Spacing(1).
		SplitN(frame.Area(), 2)

	title := text.NewLine(
		text.NewSpan("Table Widget").Bold(),
		text.NewSpan(" (Press 'q' to quit and arrow keys to navigate)"),
	).Center()
	frame.RenderWidget(title, areas[0])

	renderTable(frame, areas[1], tableState)
}

func renderTable(frame *terminal.Frame, area layout.Rect, tableState *widgets.TableState) {
	header := widgets.TableRowFromStrings([]string{"Ingredient", "Quantity", "Macros"}).
		Bold().
		BottomMargin(1)

	rows := []widgets.TableRow{
		widgets.TableRowFromStrings([]string{"Eggplant", "1 medium", "25 kcal, 6g carbs, 1g protein"}),
		widgets.TableRowFromStrings([]string{"Tomato", "2 large", "44 kcal, 10g carbs, 2g protein"}),
		widgets.TableRowFromStrings([]string{"Zucchini", "1 medium", "33 kcal, 7g carbs, 2g protein"}),
		widgets.TableRowFromStrings([]string{"Bell Pepper", "1 medium", "24 kcal, 6g carbs, 1g protein"}),
		widgets.TableRowFromStrings([]string{"Garlic", "2 cloves", "9 kcal, 2g carbs, 0.4g protein"}),
	}

	footer := widgets.TableRowFromStrings([]string{
		"Ratatouille Recipe",
		"",
		"135 kcal, 31g carbs, 6.4g protein",
	}).Italic()

	widths := []layout.Constraint{
		layout.Percentage(30),
		layout.Percentage(20),
		layout.Percentage(50),
	}

	table := widgets.NewTable(rows, widths).
		Header(header).
		Footer(footer).
		ColumnSpacing(1).
		Fg(style.White).
		RowHighlightStyle(style.NewStyle().Bg(style.Black).AddModifier(style.ModifierBold)).
		ColumnHighlightStyle(style.NewStyle().Fg(style.Gray)).
		CellHighlightStyle(style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierReversed)).
		HighlightSymbol("🍴 ")

	frame.RenderStatefulWidget(table, area, tableState)
}
