package ui

import (
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var inroom bool = false
var static_name = map[string]*fyne.StaticResource{
	"1":  resource1Png,
	"2":  resource2Png,
	"3":  resource3Png,
	"4":  resource4Png,
	"5":  resource5Png,
	"6":  resource6Png,
	"7":  resource7Png,
	"w1": resourceW1Png,
	"w2": resourceW2Png,
	"w3": resourceW3Png,
	"w4": resourceW4Png,
	"w5": resourceW5Png,
	"w6": resourceW6Png,
	"w7": resourceW7Png,
	"w8": resourceW8Png,
	"w9": resourceW9Png,
	"t1": resourceT1Png,
	"t2": resourceT2Png,
	"t3": resourceT3Png,
	"t4": resourceT4Png,
	"t5": resourceT5Png,
	"t6": resourceT6Png,
	"t7": resourceT7Png,
	"t8": resourceT8Png,
	"t9": resourceT9Png,
	"l1": resourceL1Png,
	"l2": resourceL2Png,
	"l3": resourceL3Png,
	"l4": resourceL4Png,
	"l5": resourceL5Png,
	"l6": resourceL6Png,
	"l7": resourceL7Png,
	"l8": resourceL8Png,
	"l9": resourceL9Png,
}

func makeBanner_bottom_bar() []fyne.CanvasObject {
	cardslice := []fyne.CanvasObject{}
	if myCards.Card == nil {
		/* for i := 0; i < 13; i++ {
			cc := canvas.NewImageFromResource(resourceBackPng)
			cc.FillMode = canvas.ImageFillContain
			cardslice = append(cardslice, cc)
		} */
		return nil
	} else {
		for _, card := range myCards.Card {
			if _, ok := static_name[card]; ok {
				cc := canvas.NewImageFromResource(static_name[card])
				cc.FillMode = canvas.ImageFillContain
				cardslice = append(cardslice, cc)
			}
		}
	}
	return cardslice
}

func makeGUI() fyne.CanvasObject {
	received_content := canvas.NewText("", color.Black)
	received_content.TextSize = 12

	top := makeBannner_top(&received_content)
	top_bar := widget.NewLabel("Top")
	left_bar := widget.NewLabel("Left")
	right_bar := widget.NewLabel("Right")
	bottom_bar := makeBanner_bottom_bar()

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
	objs := []fyne.CanvasObject{content, top, top_bar, left_bar, right_bar}
	objs = append(objs, bottom_bar...)
	objs = append(objs, dividers[:]...)
	return container.New(NewFysionLayout(top, top_bar, left_bar, right_bar, content, bottom_bar, dividers), objs...)
}

func updateGUI() {
	//content.Text = "Hello, Fyne!"
	//content.Refresh()
	for range time.Tick(1 * time.Millisecond) {
		GUI = makeGUI()
		log.Print("GUI Updated")
	}
}
