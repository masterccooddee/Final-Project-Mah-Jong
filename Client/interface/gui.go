package ui

import (
	"fmt"
	"image/color"
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
var bottom_bar [13]fyne.CanvasObject

func makeBanner_bottom_bar() [13]fyne.CanvasObject {
	// myCards.addCard()
	// myCards.Card = myCards.Card[:13]
	// myCards.SortCard()
	cardslice := [13]fyne.CanvasObject{}
	if myCards.Card == nil {
		for i := 0; i < 13; i++ {
			cc := canvas.NewImageFromResource(resource5Png)
			cc.FillMode = canvas.ImageFillContain
			cardslice[i] = cc
		}
		return cardslice
	} else {
		fmt.Println("myCards.Card:", myCards.Card)
		for i, card := range myCards.Card {
			if _, ok := static_name[card]; ok {
				cc := canvas.NewImageFromResource(static_name[card])
				cc.FillMode = canvas.ImageFillStretch
				cardslice[i] = cc
			}
		}
	}
	return cardslice
}

var top_bar *widget.Label

func makeGUI() *fyne.Container {
	received_content := canvas.NewText("", color.Black)
	received_content.TextSize = 12

	top := makeBannner_top(&received_content)
	top_bar = widget.NewLabel("Top")
	left_bar := widget.NewLabel("Left")
	right_bar := widget.NewLabel("Right")
	bottom_bar = makeBanner_bottom_bar()

	content := canvas.NewText("", color.Black)
	content.TextSize = 12

	dividers := [5]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	objs := []fyne.CanvasObject{content, top, top_bar, left_bar, right_bar, bottom_bar[0], bottom_bar[1], bottom_bar[2], bottom_bar[3], bottom_bar[4], bottom_bar[5], bottom_bar[6], bottom_bar[7], bottom_bar[8], bottom_bar[9], bottom_bar[10], bottom_bar[11], bottom_bar[12]}
	objs = append(objs, dividers[:]...)
	return container.New(NewFysionLayout(top, top_bar, left_bar, right_bar, content, bottom_bar, dividers), objs...)
}

func updateGUI() {
	for range time.Tick(1 * time.Second) {
		top_bar.SetText("Top " + time.Now().Format("15:04:05"))
		canvas.Refresh(top_bar)
		// bottom_bar = makeBanner_bottom_bar()
		// for i := 0; i < 13; i++ {
		// 	GUI.Objects[5+i] = bottom_bar[i]
		// }
		// GUI.Refresh()

		/* for i, obj := range GUI.Objects {
			if label, ok := obj.(*canvas.Image); ok {
				fmt.Println("label:", label, "i:", i)
			}
		} */
		bottom_bar = makeBanner_bottom_bar()
		for c := range bottom_bar {
			if c == 0 {
				bottom_bar[0].Move(fyne.NewPos(sideWidth, GUI.Size().Height-sideWidth))
				bottom_bar[0].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
			} else {
				bottom_bar[c].Move(fyne.NewPos(sideWidth+(GUI.Size().Width-sideWidth*2-150*GUI.Size().Width/1024)/13*((float32)(c)), GUI.Size().Height-sideWidth))
				bottom_bar[c].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
			}
		}
		for i := 0; i < 13; i++ {
			GUI.Objects[5+i] = bottom_bar[i]
		}
		//GUI.Objects[5] = bottom_bar[0]

	}
}
