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
	mainContainer           *fyne.Container // this holds everything
	committedLinesContainer *fyne.Container // this holds the drawing that is in the db
	linesContainer          *fyne.Container // this holds whatever we are drawing at the moment (if its empty we are not drawing aka the mouse is not down)
	lines                   []*canvas.Line

	isDrawing bool
	lastPos   fyne.Position
}

func newPaper() *paper {
	p := &paper{}
	p.ExtendBaseWidget(p)
	p.committedLinesContainer = container.NewWithoutLayout()
	p.linesContainer = container.NewWithoutLayout() // linesContainer is empty upon start
	// here I should load the committedLines from the remote db

	p.mainContainer = container.NewWithoutLayout()
	p.mainContainer.Add(p.committedLinesContainer)
	p.mainContainer.Add(p.linesContainer)
	return p
}

func (p *paper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.mainContainer)
}

func (p *paper) MouseUp(w *desktop.MouseEvent) {
	p.isDrawing = false
	p.lastPos = fyne.Position{}
	p.commitCurrentLines()

}

func (p *paper) commitCurrentLines() {
	// this is the moment where I update the remote db,

	//  then add all the linesContainer objects to the committedLinesContainer

	// I need to empty the linesContainer
	p.linesContainer.RemoveAll()
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
