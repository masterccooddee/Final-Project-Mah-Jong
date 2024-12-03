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
var GUI *fyne.Container
var w fyne.Window

func UI() {
	interconnect()

	a := app.New()
	a.Settings().SetTheme(myCustomTheme())

	w := a.NewWindow("Mahjong")
	w.Resize(fyne.NewSize(1024, 600))

	GUI = makeGUI()
	w.SetContent(GUI)

	w.CenterOnScreen()
	w.SetOnClosed(func() {
		conn.Write([]byte("LOGOUT"))
		dealer.Close()
		conn.Close()
	})
	defer w.Close()

	x := a.NewWindow("Login or Register")
	x.Resize(fyne.NewSize(300, 100))
	x.CenterOnScreen()

	y := a.NewWindow("Chi")
	y.Resize(fyne.NewSize(800, 100))
	y.CenterOnScreen()

	x.SetContent(LORinterface(&x, &w, &y))
	x.Show()

	a.Run()

}
