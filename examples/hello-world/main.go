package main

import (
	"fmt"

	gatuitcell "github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/terminal"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
	tcell "github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(fmt.Errorf("create tcell screen: %w", err))
	}

	backend, err := gatuitcell.NewWithScreen(screen)
	if err != nil {
		panic(fmt.Errorf("initialize tcell backend: %w", err))
	}
	defer backend.Close()

	term, err := terminal.New(backend)
	if err != nil {
		panic(fmt.Errorf("initialize terminal: %w", err))
	}

	if _, err := term.Draw(render); err != nil {
		panic(fmt.Errorf("draw initial frame: %w", err))
	}

	for {
		event := screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
			if _, err := term.Draw(render); err != nil {
				panic(fmt.Errorf("draw resized frame: %w", err))
			}
		case *tcell.EventKey:
			if event.Rune() == 'q' || event.Key() == tcell.KeyCtrlC {
				return
			}
		}
	}
}

func render(frame *terminal.Frame) {
	greeting := widgets.NewParagraph(text.FromString("Hello World! (press 'q' to quit)"))
	frame.RenderWidget(greeting, frame.Area())
}
