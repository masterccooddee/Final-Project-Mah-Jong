package ui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/go-zeromq/zmq4"
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
var bottom_bar [14]fyne.CanvasObject

type TappableCard struct {
	widget.Icon
}

func NewTappableCard(res fyne.Resource) *TappableCard {
	icon := &TappableCard{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	return icon
}

var tapped bool = false
var tpapped_time int
var tap_item []*TappableCard
var win_size []fyne.Size

func (i *TappableCard) MouseIn(_ *desktop.MouseEvent) {
	fmt.Println("MouseIn")
	i.Move(fyne.NewPos(i.Position().X, i.Position().Y-10))
}

//var nowthrowpos int = 0

func (i *TappableCard) Tapped(_ *fyne.PointEvent) {

	// tap_item = append(tap_item, i)
	// win_size = append(win_size, GUI.Size())
	// if tpapped_time > 0 {

	// 	if tap_item[0].Position().X != tap_item[1].Position().X {
	// 		tpapped_time = 0
	// 		tap_item[0].Move(fyne.NewPos(tap_item[0].Position().X, tap_item[0].Position().Y+30))
	// 		tap_item = tap_item[1:]

	// 	}
	// }
	// fmt.Println("Tapped")
	// if tpapped_time < 2 {
	// 	tpapped_time++
	// }

	// if tpapped_time == 1 {
	// 	i.Move(fyne.NewPos(i.Position().X, i.Position().Y-30))
	// }
	//if tpapped_time == 2 {
	sendname := strings.TrimSuffix(i.Resource.Name(), ".png")

	if action == DISCARD_CARD {
		myCards.removeCard(sendname)
		myCards.SortCard()
		dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendname)))
		fmt.Println("Send:", sendname)
		tpapped_time = 0
		throwcard <- sendname
		action = WAITING_FOR_GET_OTHER_MING

		updateGUI()
		GUI.Refresh()
		updateGUI()

	} else {
		fmt.Println("Not your turn")
	}

	//}

}

func makeFromSlice(sl []string) []string {
	result := make([]string, len(sl))
	copy(result, sl)
	return result
}

func makeBanner_bottom_bar() [14]fyne.CanvasObject {
	fmt.Println("myCards.Card:", myCards.Card)
	cardslice := [14]fyne.CanvasObject{}
	if myCards.Card == nil {
		for i := 0; i < 14; i++ {
			cc := canvas.NewRectangle(color.White)
			cc.Hide()
			//cc := canvas.NewImageFromResource(resourceBackPng)
			//cc.FillMode = canvas.ImageFillStretch
			cardslice[i] = cc
			mingcardamount = 0
		}
		return cardslice
	} else {
		var displaycards [14]string

		copy(displaycards[:], myCards.Card)
		if len(myCards.Card) < 14 {
			for i := len(myCards.Card); i < 14; i++ {
				displaycards[i] = "x"
			}
		}
		//fmt.Println("Display_card1:", displaycards)
		if putnewcard {
			displaycards[13] = newcc
			myCards.Card = append(myCards.Card, newcc)
			putnewcard = false
		}
		//fmt.Println("Display_card2:", displaycards)
		for i, card := range displaycards {
			if card == "x" {
				cc := canvas.NewRectangle(color.White)
				cc.Hide()
				cardslice[i] = cc
			} else {
				cc := NewTappableCard(static_name[card])
				cardslice[i] = cc
			}
		}

		return cardslice
	}

}

// var new_card fyne.CanvasObject
var top_bar *widget.Label
var grid *fyne.Container
var mingbuttonlist *fyne.Container
var gangbutton *widget.Button
var pongbutton *widget.Button
var chibutton *widget.Button

func makeGUI() *fyne.Container {
	received_content := canvas.NewText("", color.Black)
	received_content.TextSize = 12

	top := makeBannner_top(&received_content)
	top_bar = widget.NewLabel("Top")
	left_bar := widget.NewLabel("Left")
	right_bar := widget.NewLabel("Right")
	bottom_bar = makeBanner_bottom_bar()

	grid = container.NewGridWithColumns(3,
		container.NewCenter(canvas.NewText("", color.Black)),
		container.NewCenter(canvas.NewText("", color.Black)), //1 : Front
		container.NewCenter(canvas.NewText("", color.Black)),
		container.NewCenter(canvas.NewText("", color.Black)), //3 : Left
		container.NewCenter(canvas.NewText("", color.Black)),
		container.NewCenter(canvas.NewText("", color.Black)), //5 : Right
		container.NewCenter(canvas.NewText("", color.Black)),
		container.NewCenter(canvas.NewText("", color.Black)), //7 : Myself
		container.NewCenter(canvas.NewText("", color.Black)),
	)
	gangbutton = widget.NewButton("Gang", func() {})
	pongbutton = widget.NewButton("Pong", func() {})
	chibutton = widget.NewButton("Chi", func() {})

	mingbuttonlist = container.NewHBox(gangbutton, pongbutton, chibutton)

	// 修改 grid 中所有 canvas.Text 的 TextSize
	for _, obj := range grid.Objects {
		if center, ok := obj.(*fyne.Container); ok {
			if text, ok := center.Objects[0].(*canvas.Text); ok {
				text.TextSize = 15 // 設置你想要的字體大小
			}
		}
	}

	dividers := [5]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	objs := []fyne.CanvasObject{mingbuttonlist, grid, top, top_bar, left_bar, right_bar, bottom_bar[0], bottom_bar[1], bottom_bar[2], bottom_bar[3], bottom_bar[4], bottom_bar[5], bottom_bar[6], bottom_bar[7], bottom_bar[8], bottom_bar[9], bottom_bar[10], bottom_bar[11], bottom_bar[12], bottom_bar[13]}
	objs = append(objs, dividers[:]...)
	return container.New(NewFysionLayout(top, top_bar, left_bar, right_bar, grid, mingbuttonlist, bottom_bar, dividers), objs...)
}

func updateGUI() {

	top_bar.SetText("Top " + time.Now().Format("15:04:05"))
	canvas.Refresh(top_bar)

	bottom_bar = makeBanner_bottom_bar()
	for c := 0; c < 14; c++ {
		//fmt.Println("c:", c)
		if c == 0 {
			bottom_bar[0].Move(fyne.NewPos(sideWidth, GUI.Size().Height-sideWidth))
			bottom_bar[0].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
		} else if c == 13 {
			bottom_bar[c].Move(fyne.NewPos(sideWidth+(GUI.Size().Width-sideWidth*2-150*GUI.Size().Width/1024)/13*(float32)(14), GUI.Size().Height-sideWidth))
			bottom_bar[c].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
		} else {
			bottom_bar[c].Move(fyne.NewPos(sideWidth+(GUI.Size().Width-sideWidth*2-150*GUI.Size().Width/1024)/13*((float32)(c)), GUI.Size().Height-sideWidth))
			bottom_bar[c].Resize(fyne.NewSize((GUI.Size().Width-sideWidth*2)/13, sideWidth))
		}
	}
	for i := 0; i < 14; i++ {
		GUI.Objects[5+i] = bottom_bar[i]
	}
	GUI.Refresh()

}
