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
		layout.Fill(1),
	).Spacing(1).SplitN(frame.Area(), 2)

	title := text.NewLine(
		text.NewSpan("Canvas Widget").Bold(),
		text.NewSpan(" (Press 'q' to quit)"),
	).Center()
	frame.RenderWidget(title, areas[0])

	renderCanvas(frame, areas[1])
}

func renderCanvas(frame *terminal.Frame, area layout.Rect) {
	canvas := widgets.NewCanvas().
		XBounds(-180, 180).
		YBounds(-90, 90).
		Marker(widgets.CanvasMarkerBraille).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.Map{Resolution: widgets.MapResolutionHigh, Color: style.White})
			ctx.Layer()
			ctx.Draw(widgets.NewCanvasLine(0, 10, 10, 10, style.Blue))
			ctx.Draw(widgets.NewRectangle(10, 20, 10, 10, style.Green))
			ctx.Draw(widgets.NewPoints([]widgets.CanvasPoint{
				{X: 2.3522, Y: 48.8566},    // Paris
				{X: -122.3321, Y: 47.6062}, // Seattle
				{X: -79.3837, Y: 43.6511},  // Toronto
				{X: 32.8597, Y: 39.9334},   // Ankara
			}, style.Red))
		})

	frame.RenderWidget(canvas, area)
}
