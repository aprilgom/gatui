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
	tcell "github.com/gdamore/tcell/v3"
)

type app struct {
	listState widgets.ListState
}

func newApp() app {
	return app{listState: widgets.NewListState().WithSelected(0)}
}

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

	app := newApp()
	if err := app.draw(term); err != nil {
		return fmt.Errorf("draw initial frame: %w", err)
	}

	for {
		event := <-screen.EventQ()
		switch event := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
			if err := app.draw(term); err != nil {
				return fmt.Errorf("draw resized frame: %w", err)
			}
		case *tcell.EventKey:
			if app.handleKey(event) {
				return nil
			}
			if err := app.draw(term); err != nil {
				return fmt.Errorf("draw frame: %w", err)
			}
		}
	}
}

func (a *app) handleKey(event *tcell.EventKey) bool {
	switch {
	case event.Str() == "q" || event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyCtrlC:
		return true
	case event.Str() == "j" || event.Key() == tcell.KeyDown:
		a.listState.SelectNext()
	case event.Str() == "k" || event.Key() == tcell.KeyUp:
		a.listState.SelectPrevious()
	}
	return false
}

func (a *app) draw(term *terminal.Terminal) error {
	_, err := term.Draw(func(frame *terminal.Frame) {
		a.render(frame)
	})
	return err
}

func (a *app) render(frame *terminal.Frame) {
	areas := layout.NewVerticalLayout(
		layout.Length(1),
		layout.Fill(1),
		layout.Fill(1),
	).Spacing(1).SplitN(frame.Area(), 3)

	title := text.NewLine(
		text.NewSpan("List Widget").Bold(),
		text.NewSpan(" (Press 'q' to quit and arrow keys to navigate)"),
	).Center()
	frame.RenderWidget(title, areas[0])

	a.renderList(frame, areas[1])
	renderBottomList(frame, areas[2])
}

func (a *app) renderList(frame *terminal.Frame, area layout.Rect) {
	items := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
	list := widgets.NewListFromStrings(items).
		Fg(style.White).
		HighlightStyle(style.StyleFromModifier(style.ModifierReversed)).
		HighlightSymbol("> ")

	frame.RenderStatefulWidget(list, area, &a.listState)
}

func renderBottomList(frame *terminal.Frame, area layout.Rect) {
	items := []widgets.ListItem{
		widgets.ListItemFromString("[Remy]: I'm building one now.\nIt even supports multiline text!"),
		widgets.ListItemFromString("[Gusteau]: With enough passion, yes."),
		widgets.ListItemFromString("[Remy]: But can anyone build a TUI in Rust?"),
		widgets.ListItemFromString("[Gusteau]: Anyone can cook!"),
	}
	list := widgets.NewList(items).
		Fg(style.White).
		HighlightStyle(style.NewStyle().Fg(style.Yellow).AddModifier(style.ModifierItalic)).
		HighlightSymbol("> ").
		RepeatHighlightSymbol(true)

	state := widgets.NewListState().WithSelected(0)
	frame.RenderStatefulWidget(list, area, &state)
}
