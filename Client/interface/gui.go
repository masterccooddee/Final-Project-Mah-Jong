package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeBannner_top() fyne.CanvasObject {
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

func makeBanner_bottom_bar() fyne.CanvasObject {

	card := canvas.NewImageFromResource(resourceWordfivePng)
	card.FillMode = canvas.ImageFillContain
	return container.NewStack(container.NewPadded(card))
}

func makeGUI() fyne.CanvasObject {
	top := makeBannner_top()
	top_bar := widget.NewLabel("Top")
	left_bar := widget.NewLabel("Left")
	right_bar := widget.NewLabel("Right")
	bottom_bar := [13]fyne.CanvasObject{
		makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(),
	}

	content := canvas.NewRectangle(color.Gray{Y: 0xEE})

	dividers := [5]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	//return container.NewBorder(makeBannner(), nil, left, right, content)
	objs := []fyne.CanvasObject{content, top, top_bar, left_bar, right_bar, bottom_bar[0], bottom_bar[1], bottom_bar[2], bottom_bar[3], bottom_bar[4], bottom_bar[5], bottom_bar[6], bottom_bar[7], bottom_bar[8], bottom_bar[9], bottom_bar[10], bottom_bar[11], bottom_bar[12],
		dividers[0], dividers[1], dividers[2], dividers[3], dividers[4]}
	return container.New(NewFysionLayout(top, top_bar, left_bar, right_bar, content, bottom_bar, dividers), objs...)
}
