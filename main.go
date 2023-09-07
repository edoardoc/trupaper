package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type paper struct {
	widget.BaseWidget
	// desktop.MouseEvent
	// this should be able to intercept mouse events
	desktop.Hoverable
}

func newPaper() *paper {
	ppr := &paper{}
	return ppr
}

// intercept mouse tap
func (t *paper) Tapped(where *fyne.PointEvent) {
	log.Printf("I have been tapped %+v", where)
}

// intercept mouse secondary tap
func (t *paper) TappedSecondary(_ *fyne.PointEvent) {
}

func (t *paper) MouseUp(w *desktop.MouseEvent) {
	log.Printf("I have been MouseUp %+v", w)
}
func (t *paper) MouseDown(w *desktop.MouseEvent) {
	log.Printf("I have been MouseDown %+v", w)
}

// intercept mouse hover
func (t *paper) MouseIn(_ *desktop.MouseEvent) {
	log.Printf("I have been MouseIn")
}

func (t *paper) MouseOut() {
	log.Printf("I have been MouseOut")
}

func (t *paper) MouseMoved(_ *desktop.MouseEvent) {
	log.Printf("I have been MouseMoved")
}

func main() {
	a := app.New()
	w := a.NewWindow("Trupaper")
	w.SetContent(newPaper())
	w.Resize(fyne.NewSize(500, 500))
	w.ShowAndRun()
}
