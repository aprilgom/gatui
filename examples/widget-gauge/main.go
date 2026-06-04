package main

import (
	"fmt"
	"os"

	gatuitcell "github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/symbols"
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

	for {
		if _, err := term.Draw(render); err != nil {
			return fmt.Errorf("draw frame: %w", err)
		}

		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if event.Rune() == 'q' || event.Key() == tcell.KeyCtrlC {
				return nil
			}
		}
	}
}

func render(frame *terminal.Frame) {
	areas := layout.NewVerticalLayout(
		layout.Length(1),
		layout.Max(2),
		layout.Fill(1),
	).Spacing(1).SplitN(frame.Area(), 3)

	title := text.NewLine(
		text.NewSpan("Gauge Widget").Bold(),
		text.NewSpan(" (Press 'q' to quit)"),
	).Center()
	frame.RenderWidget(title, areas[0])

	renderGauge(frame, areas[1])
	renderLineGauge(frame, areas[2])
}

func renderGauge(frame *terminal.Frame, area layout.Rect) {
	gauge := widgets.NewGauge().
		Bold().
		GaugeStyle(style.NewStyle().Fg(style.Blue).Bg(style.Black)).
		LabelString("Year Progress").
		Percent(80)
	frame.RenderWidget(gauge, area)
}

func renderLineGauge(frame *terminal.Frame, area layout.Rect) {
	lineGauge := widgets.NewLineGauge().
		FilledStyle(style.NewStyle().Fg(style.White).Bg(style.Red).AddModifier(style.ModifierBold)).
		UnfilledStyle(style.NewStyle().Fg(style.Gray).Bg(style.Black)).
		LabelString("HP").
		Ratio(0.42).
		FilledSymbol(symbols.LineThickHorizontal).
		UnfilledSymbol(symbols.LineThickHorizontal)
	frame.RenderWidget(lineGauge, area)
}
