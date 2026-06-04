package main

import (
	"fmt"
	"time"

	tcellbackend "github.com/aprilgom/gatui/backend/tcell"
	"github.com/aprilgom/gatui/layout"
	"github.com/aprilgom/gatui/style"
	"github.com/aprilgom/gatui/terminal"
	"github.com/aprilgom/gatui/text"
	"github.com/aprilgom/gatui/widgets"
	tcell "github.com/gdamore/tcell/v2"
)

type App struct {
	x          float64
	y          float64
	ball       widgets.Circle
	playground layout.Rect
	vx         float64
	vy         float64
	marker     widgets.CanvasMarker
	points     []layout.Position
	isDrawing  bool
}

func NewApp() *App {
	return &App{
		ball:       widgets.NewCircle(20, 40, 10, style.Yellow),
		playground: layout.NewRect(10, 10, 200, 100),
		vx:         1,
		vy:         1,
		marker:     widgets.CanvasMarkerDot,
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	backend, err := tcellbackend.NewWithScreen(screen)
	if err != nil {
		panic(err)
	}
	defer backend.Close()
	screen.EnableMouse()

	term, err := terminal.New(backend)
	if err != nil {
		panic(err)
	}

	app := NewApp()
	if err := app.Run(term, screen); err != nil {
		panic(err)
	}
}

func (a *App) Run(term *terminal.Terminal, screen tcell.Screen) error {
	events := make(chan tcell.Event)
	go func() {
		for {
			events <- screen.PollEvent()
		}
	}()

	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for {
		if _, err := term.Draw(a.Render); err != nil {
			return err
		}

		select {
		case event := <-events:
			if a.HandleEvent(term, event) {
				return nil
			}
		case <-ticker.C:
			a.OnTick()
		}
	}
}

func (a *App) HandleEvent(term *terminal.Terminal, event tcell.Event) bool {
	switch event := event.(type) {
	case *tcell.EventKey:
		return a.HandleKeyEvent(event)
	case *tcell.EventMouse:
		a.HandleMouseEvent(event)
	case *tcell.EventResize:
		width, height := event.Size()
		_ = term.Resize(layout.NewRect(0, 0, width, height))
	}
	return false
}

func (a *App) HandleKeyEvent(event *tcell.EventKey) bool {
	if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
		return true
	}

	switch event.Key() {
	case tcell.KeyDown:
		a.y += 1
	case tcell.KeyUp:
		a.y -= 1
	case tcell.KeyRight:
		a.x += 1
	case tcell.KeyLeft:
		a.x -= 1
	case tcell.KeyEnter:
		a.CycleMarker()
	}

	switch event.Rune() {
	case 'j':
		a.y += 1
	case 'k':
		a.y -= 1
	case 'l':
		a.x += 1
	case 'h':
		a.x -= 1
	}

	return false
}

func (a *App) HandleMouseEvent(event *tcell.EventMouse) {
	x, y := event.Position()
	button := event.Buttons()
	if button&tcell.Button1 != 0 {
		a.isDrawing = true
		a.points = append(a.points, layout.NewPosition(x, y))
		return
	}
	if a.isDrawing {
		a.isDrawing = false
	}
}

func (a *App) CycleMarker() {
	switch a.marker {
	case widgets.CanvasMarkerDot:
		a.marker = widgets.CanvasMarkerBraille
	case widgets.CanvasMarkerBraille:
		a.marker = widgets.CanvasMarkerBlock
	case widgets.CanvasMarkerBlock:
		a.marker = widgets.CanvasMarkerHalfBlock
	case widgets.CanvasMarkerHalfBlock:
		a.marker = widgets.CanvasMarkerQuadrant
	case widgets.CanvasMarkerQuadrant:
		a.marker = widgets.CanvasMarkerSextant
	case widgets.CanvasMarkerSextant:
		a.marker = widgets.CanvasMarkerOctant
	case widgets.CanvasMarkerOctant:
		a.marker = widgets.CanvasMarkerCustom("x")
	case widgets.CanvasMarkerCustom("x"):
		a.marker = widgets.CanvasMarkerBar
	default:
		a.marker = widgets.CanvasMarkerDot
	}
}

func (a *App) OnTick() {
	if a.ball.X-a.ball.Radius < float64(a.playground.Left()) ||
		a.ball.X+a.ball.Radius > float64(a.playground.Right()) {
		a.vx = -a.vx
	}
	if a.ball.Y-a.ball.Radius < float64(a.playground.Top()) ||
		a.ball.Y+a.ball.Radius > float64(a.playground.Bottom()) {
		a.vy = -a.vy
	}
	a.ball.X += a.vx
	a.ball.Y += a.vy
}

func (a *App) Render(frame *terminal.Frame) {
	header := text.NewText(
		text.LineFromString("Canvas Example").Bold(),
		text.LineFromString("<q> Quit | <enter> Change Marker | <hjkl> Move"),
	).Center()

	vertical := layout.NewVerticalLayout(
		layout.Length(header.Height()),
		layout.Fill(1),
		layout.Fill(1),
	)
	areas := vertical.SplitN(frame.Area(), 3)
	frame.RenderWidget(header, areas[0])

	horizontal := layout.NewHorizontalLayout(layout.Fill(1), layout.Fill(1))
	up := horizontal.SplitN(areas[1], 2)
	down := horizontal.SplitN(areas[2], 2)

	frame.RenderWidget(a.MapCanvas(), down[0])
	frame.RenderWidget(a.DrawCanvas(up[0]), up[0])
	frame.RenderWidget(a.PongCanvas(), up[1])
	frame.RenderWidget(a.BoxesCanvas(down[1]), down[1])
}

func (a *App) MapCanvas() widgets.Widget {
	return widgets.NewCanvas().
		Block(widgets.BorderedBlock().Title(text.LineFromString("World"))).
		Marker(a.marker).
		XBounds(-180, 180).
		YBounds(-90, 90).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(widgets.Map{Resolution: widgets.MapResolutionHigh, Color: style.Green})
			ctx.Print(a.x, -a.y, text.NewSpan("You are here").Fg(style.Yellow))
		})
}

func (a *App) DrawCanvas(area layout.Rect) widgets.Widget {
	return widgets.NewCanvas().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Draw here"))).
		Marker(a.marker).
		XBounds(0, float64(area.Width)).
		YBounds(0, float64(area.Height)).
		Paint(func(ctx *widgets.CanvasContext) {
			points := make([]widgets.CanvasPoint, 0, len(a.points))
			for _, point := range a.points {
				points = append(points, widgets.CanvasPoint{
					X: float64(point.X - area.Left()),
					Y: float64(area.Bottom() - point.Y),
				})
			}
			ctx.Draw(widgets.NewPoints(points, style.White))
		})
}

func (a *App) PongCanvas() widgets.Widget {
	return widgets.NewCanvas().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Pong"))).
		Marker(a.marker).
		XBounds(10, 210).
		YBounds(10, 110).
		Paint(func(ctx *widgets.CanvasContext) {
			ctx.Draw(a.ball)
		})
}

func (a *App) BoxesCanvas(area layout.Rect) widgets.Widget {
	top := float64(area.Height)*2 - 4
	return widgets.NewCanvas().
		Block(widgets.BorderedBlock().Title(text.LineFromString("Rects"))).
		Marker(a.marker).
		XBounds(0, float64(area.Width)).
		YBounds(0, top).
		Paint(func(ctx *widgets.CanvasContext) {
			for i := 0; i <= 11; i++ {
				x := float64(i*i+3*i)/2 + 2
				size := float64(i)
				ctx.Draw(widgets.NewRectangle(x, 2, size, size, style.Red))
				ctx.Draw(widgets.NewRectangle(x, 21, size, size, style.Blue))
			}
			for i := range 100 {
				if i%10 != 0 {
					ctx.Print(float64(i)+1, 0, text.NewSpan(fmt.Sprintf("%d", i%10)))
				}
				if i%2 == 0 && i%10 != 0 {
					ctx.Print(0, float64(i), text.NewSpan(fmt.Sprintf("%d", i%10)))
				}
			}
		})
}
