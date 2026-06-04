package main

import (
	"fmt"
	"time"

	gatuitcell "github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/buffer"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/terminal"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
	tcell "github.com/gdamore/tcell/v3"
)

var (
	gauge1Color      = style.RGBColor(153, 27, 27)
	gauge2Color      = style.RGBColor(22, 101, 52)
	gauge3Color      = style.RGBColor(30, 64, 175)
	gauge4Color      = style.RGBColor(154, 52, 18)
	customLabelColor = style.RGBColor(226, 232, 240)
)

type appState int

const (
	stateRunning appState = iota
	stateStarted
	stateQuitting
)

type app struct {
	state           appState
	progressColumns int
	progress1       int
	progress2       float64
	progress3       float64
	progress4       float64
}

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
	if err := term.HideCursor(); err != nil {
		panic(fmt.Errorf("hide cursor: %w", err))
	}

	if err := new(app).run(term, screen); err != nil {
		panic(err)
	}
}

func (a *app) run(term *terminal.Terminal, screen tcell.Screen) error {
	events := make(chan tcell.Event)
	go func() {
		for {
			events <- <-screen.EventQ()
		}
	}()

	ticker := time.NewTicker(time.Second / 20)
	defer ticker.Stop()

	if _, err := term.Draw(func(frame *terminal.Frame) {
		frame.RenderWidget(a, frame.Area())
	}); err != nil {
		return fmt.Errorf("draw initial frame: %w", err)
	}

	for a.state != stateQuitting {
		select {
		case event := <-events:
			if err := a.handleEvent(event, term, screen); err != nil {
				return err
			}
		case <-ticker.C:
			a.update(term.Area().Width)
			if _, err := term.Draw(func(frame *terminal.Frame) {
				frame.RenderWidget(a, frame.Area())
			}); err != nil {
				return fmt.Errorf("draw frame: %w", err)
			}
		}
	}
	return nil
}

func (a *app) handleEvent(event tcell.Event, term *terminal.Terminal, screen tcell.Screen) error {
	switch event := event.(type) {
	case *tcell.EventResize:
		screen.Sync()
		if _, err := term.Draw(func(frame *terminal.Frame) {
			frame.RenderWidget(a, frame.Area())
		}); err != nil {
			return fmt.Errorf("draw resized frame: %w", err)
		}
	case *tcell.EventKey:
		switch {
		case event.Str() == " " || event.Key() == tcell.KeyEnter:
			a.start()
		case event.Str() == "q" || event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyCtrlC:
			a.quit()
		}
	}
	return nil
}

func (a *app) update(terminalWidth int) {
	if a.state != stateStarted || terminalWidth <= 0 {
		return
	}

	a.progressColumns = clamp(a.progressColumns+1, 0, terminalWidth)
	a.progress1 = a.progressColumns * 100 / terminalWidth
	a.progress2 = float64(a.progressColumns) * 100 / float64(terminalWidth)

	a.progress3 = clampFloat(a.progress3+0.1, 40, 100)
	a.progress4 = clampFloat(a.progress4+0.1, 40, 100)
}

func (a *app) start() {
	a.state = stateStarted
}

func (a *app) quit() {
	a.state = stateQuitting
}

func (a *app) Render(area layout.Rect, buf *buffer.Buffer) {
	outer := layout.NewVerticalLayout(
		layout.Length(2),
		layout.Min(0),
		layout.Length(1),
	).SplitN(area, 3)
	headerArea, gaugeArea, footerArea := outer[0], outer[1], outer[2]

	gauges := layout.NewVerticalLayout(
		layout.Ratio(1, 4),
		layout.Ratio(1, 4),
		layout.Ratio(1, 4),
		layout.Ratio(1, 4),
	).SplitN(gaugeArea, 4)

	renderHeader(headerArea, buf)
	renderFooter(footerArea, buf)

	a.renderGauge1(gauges[0], buf)
	a.renderGauge2(gauges[1], buf)
	a.renderGauge3(gauges[2], buf)
	a.renderGauge4(gauges[3], buf)
}

func renderHeader(area layout.Rect, buf *buffer.Buffer) {
	widgets.NewParagraph(text.FromString("Gatui Gauge Example")).
		Alignment(layout.Center).
		Fg(customLabelColor).
		Bold().
		Render(area, buf)
}

func renderFooter(area layout.Rect, buf *buffer.Buffer) {
	widgets.NewParagraph(text.FromString("Press ENTER to start")).
		Alignment(layout.Center).
		Fg(customLabelColor).
		Bold().
		Render(area, buf)
}

func (a *app) renderGauge1(area layout.Rect, buf *buffer.Buffer) {
	widgets.NewGauge().
		Block(titleBlock("Gauge with percentage")).
		GaugeStyle(style.NewStyle().Fg(gauge1Color)).
		Percent(a.progress1).
		Render(area, buf)
}

func (a *app) renderGauge2(area layout.Rect, buf *buffer.Buffer) {
	label := text.NewSpan(fmt.Sprintf("%.1f/100", a.progress2)).
		Italic().
		Bold().
		Fg(customLabelColor)

	widgets.NewGauge().
		Block(titleBlock("Gauge with ratio and custom label")).
		GaugeStyle(style.NewStyle().Fg(gauge2Color)).
		Ratio(a.progress2/100).
		Label(label).
		Render(area, buf)
}

func (a *app) renderGauge3(area layout.Rect, buf *buffer.Buffer) {
	widgets.NewGauge().
		Block(titleBlock("Gauge with ratio (no unicode)")).
		GaugeStyle(style.NewStyle().Fg(gauge3Color)).
		Ratio(a.progress3/100).
		LabelString(fmt.Sprintf("%.1f%%", a.progress3)).
		Render(area, buf)
}

func (a *app) renderGauge4(area layout.Rect, buf *buffer.Buffer) {
	widgets.NewGauge().
		Block(titleBlock("Gauge with ratio (unicode)")).
		GaugeStyle(style.NewStyle().Fg(gauge4Color)).
		Ratio(a.progress4/100).
		LabelString(fmt.Sprintf("%.1f%%", a.progress4)).
		UseUnicode(true).
		Render(area, buf)
}

func titleBlock(title string) widgets.Block {
	return widgets.NewBlock().
		Borders(widgets.NoBorders).
		Padding(widgets.PaddingVertical(1)).
		Title(text.LineFromString(title).Center()).
		Fg(customLabelColor)
}

func clamp(value, low, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func clampFloat(value, low, high float64) float64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
