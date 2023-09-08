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
	lines     []*canvas.Line
	isDrawing bool
	lastPos   fyne.Position
	desktop.Hoverable
}

func newPaper() *paper {
	icon := &paper{
		Container: &fyne.Container{
			Objects: []fyne.CanvasObject{},
			Layout:  nil,
			Hidden:  false,
			// Size:     fyne.Size{},
			// Position: fyne.Position{},
		},
		lines:     []*canvas.Line{},
		isDrawing: false,
		lastPos:   fyne.Position{},
	}
	// icon.ExtendBaseWidget(icon)
	return icon
}

func (p *paper) MouseUp(w *desktop.MouseEvent) {
	p.isDrawing = false
	p.lastPos = fyne.Position{}
}

func (p *paper) MouseDown(w *desktop.MouseEvent) {
	p.isDrawing = true
	log.Printf("Event MouseDown %+v", w)
}

func (p *paper) MouseIn(_ *desktop.MouseEvent) {
	log.Printf("Event MouseIn")
}

func (p *paper) MouseOut() {
	log.Printf("Event MouseOut")
}

func (p *paper) MouseMoved(e *desktop.MouseEvent) {
	if p.isDrawing {
		if p.lastPos != (fyne.Position{}) {
			line := canvas.NewLine(color.Black)
			line.StrokeWidth = 5
			line.Position1 = p.lastPos
			line.Position2 = e.Position
			p.Container.Add(line)

			p.lines = append(p.lines, line)
		}
		p.lastPos = e.Position
	} else {
		p.lastPos = fyne.Position{}
	}

	log.Printf("Event MouseMoved")
}

var mainContainer *fyne.Container

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Trupaper")
	mainContainer = container.NewWithoutLayout()

	drawable := newPaper()
	drawable.Resize(fyne.NewSize(500, 500))
	mainContainer.Add(drawable)

	// top := canvas.NewText("top bar", color.White)
	// left := canvas.NewText("left", color.White)
	// middle := canvas.NewText("content", color.White)
	// content := container.NewBorder(top, nil, left, nil, middle)
	// myWindow.SetContent(content)

	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.ShowAndRun()
}
