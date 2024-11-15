package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var inroom bool = false

func makeBanner_bottom_bar() fyne.CanvasObject {

	for index := range myCards.Card {
		card := canvas.NewImageFromResource(resourceWordfivePng)
		card.FillMode = canvas.ImageFillContain
		return container.NewStack(container.NewPadded(card))

	}

}

func makeGUI() fyne.CanvasObject {
	received_content := canvas.NewText("", color.Black)
	received_content.TextSize = 12

	top := makeBannner_top(&received_content)
	top_bar := widget.NewLabel("Top")
	left_bar := widget.NewLabel("Left")
	right_bar := widget.NewLabel("Right")
	bottom_bar := [13]fyne.CanvasObject{
		makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(),
	}

	content := canvas.NewText("", color.Black)
	content.TextSize = 12

	/* msg, err = dealer.Recv()
	fmt.Println("msg:", msg, "err:", err)
	content.Text = string(msg.Frames[0])
	content.Refresh() */

	dividers := [5]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	//return container.NewBorder(makeBannner(), nil, left, right, content)
	objs := []fyne.CanvasObject{content, top, top_bar, left_bar, right_bar, bottom_bar[0], bottom_bar[1], bottom_bar[2], bottom_bar[3], bottom_bar[4], bottom_bar[5], bottom_bar[6], bottom_bar[7], bottom_bar[8], bottom_bar[9], bottom_bar[10], bottom_bar[11], bottom_bar[12],
		dividers[0], dividers[1], dividers[2], dividers[3], dividers[4]}
	return container.New(NewFysionLayout(top, top_bar, left_bar, right_bar, content, bottom_bar, dividers), objs...)
}
