package main

import (
	"fmt"
	"os"

	"github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/terminal"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
	tcelllib "github.com/gdamore/tcell/v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	screen, err := tcelllib.NewScreen()
	if err != nil {
		return err
	}

	backend, err := tcell.NewWithScreen(screen)
	if err != nil {
		return err
	}
	defer backend.Close()

	term, err := terminal.New(backend)
	if err != nil {
		return err
	}

	for {
		if _, err := term.Draw(func(frame *terminal.Frame) {
			frame.RenderWidget(widgets.NewParagraph(text.FromString("Hello World!")), frame.Area())
		}); err != nil {
			return err
		}

		if _, ok := (<-screen.EventQ()).(*tcelllib.EventKey); ok {
			break
		}
	}

	return nil
}
