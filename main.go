package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trupaper/ui"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Trupaper")

	drawable := ui.NewPaper()
	drawable.Resize(fyne.NewSize(500, 500))
	ttop := container.NewGridWithColumns(3,
		widget.NewButton("<<", drawable.Left),
		widget.NewSeparator(),
		widget.NewButton(">>", drawable.Right),
	)
	content := container.NewBorder(ttop, nil, nil, nil, drawable)
	myWindow.SetContent(content)

	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.ShowAndRun()
}
