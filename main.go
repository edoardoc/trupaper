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
	widget.BaseWidget
	linesContainer *fyne.Container

	lines     []*canvas.Line
	isDrawing bool
	lastPos   fyne.Position
}

func newPaper() *paper {
	p := &paper{}
	p.ExtendBaseWidget(p)
	p.linesContainer = container.NewWithoutLayout()
	return p
}

func (p *paper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.linesContainer)
}

func (p *paper) MouseUp(w *desktop.MouseEvent) {
	p.isDrawing = false
	p.lastPos = fyne.Position{}
	// this is the moment where I update the remote db

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
			p.lines = append(p.lines, line)

			p.linesContainer.Add(line)
			p.linesContainer.Refresh()
		}
		p.lastPos = e.Position
	} else {
		p.lastPos = fyne.Position{}
	}
	log.Printf("Event MouseMoved")
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Trupaper")

	drawable := newPaper()
	drawable.Resize(fyne.NewSize(500, 500))
	top := canvas.NewText("top bar", color.Black)
	content := container.NewBorder(top, nil, nil, nil, drawable)
	myWindow.SetContent(content)

	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.ShowAndRun()
}
