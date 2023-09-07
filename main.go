package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
)

// paper is a widget that can be drawn on, it is a container to detect mouse events
type paper struct {
	*fyne.Container
	// desktop.Hoverable
	lines     []*canvas.Line
	isDrawing bool
	lastPos   fyne.Position
}

func (p *paper) Tapped(where *fyne.PointEvent) {
	log.Printf("I have been tapped %+v", where)
}

// intercept mouse secondary tap
func (p *paper) TappedSecondary(_ *fyne.PointEvent) {
}

func (p *paper) MouseUp(w *desktop.MouseEvent) {
	p.isDrawing = false
	p.lastPos = fyne.Position{}
}

func (p *paper) MouseDown(w *desktop.MouseEvent) {
	p.isDrawing = true
	log.Printf("I have been MouseDown %+v", w)
}

func (p *paper) MouseIn(_ *desktop.MouseEvent) {
	log.Printf("I have been MouseIn")
}

func (p *paper) MouseOut() {
	log.Printf("I have been MouseOut")
}

func (p *paper) MouseMoved(e *desktop.MouseEvent) {
	if p.isDrawing {
		if p.lastPos != (fyne.Position{}) {
			line := canvas.NewLine(color.Black)
			line.StrokeWidth = 5
			line.Position1 = p.lastPos
			line.Position2 = e.Position
			p.Add(line)

			p.lines = append(p.lines, line)
		}
		p.lastPos = e.Position
	} else {
		p.lastPos = fyne.Position{}
	}

	log.Printf("I have been MouseMoved")
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Trupaper")
	myCanvas := myWindow.Canvas()

	// myWidget := container.NewWithoutLayout()

	ppr := &paper{
		container.NewWithoutLayout(),
		[]*canvas.Line{},
		false,
		fyne.Position{},
	}

	line := canvas.NewLine(color.Black)
	line.Position1 = fyne.NewPos(0, 0)
	line.Position2 = fyne.NewPos(100, 100)
	line.StrokeWidth = 5

	ppr.Add(line)

	myCanvas.SetContent(ppr) // R/O
	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.ShowAndRun()
}
