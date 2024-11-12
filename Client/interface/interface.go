package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

func ShowError(err error, w fyne.Window) {
	dialog.ShowError(err, w)
}

var LoginSuccess = false

func UI() {
	go interconnect()
	a := app.New()
	a.Settings().SetTheme(myCustomTheme())

	w := a.NewWindow("Mahjong")
	w.Resize(fyne.NewSize(1024, 600))
	w.SetContent(makeGUI())

	x := a.NewWindow("Login or Register")
	x.Resize(fyne.NewSize(300, 100))

	x.SetContent(LORinterface(&x, &w))
	x.Show()

	a.Run()

}
