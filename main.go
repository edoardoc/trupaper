package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// paper is a widget that can be drawn on, it is a container to detect mouse events
type paper struct {
	widget.Icon

	lines     []*canvas.Line
	isDrawing bool
	lastPos   fyne.Position
}

func newPaper() *paper {
	icon := &paper{}
	// icon.ExtendBaseWidget(icon)
	return icon
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
			mainContainer.Add(line)

			p.lines = append(p.lines, line)
		}
		p.lastPos = e.Position
	} else {
		p.lastPos = fyne.Position{}
	}

	log.Printf("I have been MouseMoved")
}

var mainContainer *fyne.Container

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Trupaper")
	myCanvas := myWindow.Canvas()
	mainContainer = container.NewWithoutLayout()

	drawable := newPaper()
	drawable.Resize(fyne.NewSize(500, 500))
	mainContainer.Add(drawable)

	myCanvas.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.ShowAndRun()
}
