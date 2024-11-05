package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeBannner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.GridIcon(), func() {}),
	)

	logo1 := canvas.NewImageFromResource(resourceLogo1Png)
	logo1.FillMode = canvas.ImageFillContain

	logo2 := canvas.NewImageFromResource(resourceLogo2Png)
	logo2.FillMode = canvas.ImageFillContain

	logo3 := canvas.NewImageFromResource(resourceLogo3Png)
	logo3.FillMode = canvas.ImageFillContain

	return container.NewStack(toolbar, container.NewPadded(logo1))
}

func makeGUI() fyne.CanvasObject {
	top := makeBannner()
	top_bar := widget.NewLabel("Top")
	left_bar := widget.NewLabel("Left")
	right_bar := widget.NewLabel("Right")
	bottom_bar := widget.NewLabel("Bottom")

	content := canvas.NewRectangle(color.Gray{Y: 0xEE})

	dividers := [5]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	//return container.NewBorder(makeBannner(), nil, left, right, content)
	objs := []fyne.CanvasObject{content, top, top_bar, left_bar, right_bar, bottom_bar, dividers[0], dividers[1], dividers[2], dividers[3], dividers[4]}
	return container.New(NewFysionLayout(top, top_bar, left_bar, right_bar, bottom_bar, content, dividers), objs...)
}
