package ui

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/go-zeromq/zmq4"
)

const (
	DRAW_CARD = iota
	GAME_START_WAIT
	GAME_START
	CHECK_SELF_MING
	DISCARD_CARD
	WAITING_FOR_GET_DISCARD_CARD
	CHECK_MING
	CHOOSE_SELF_MING
	CHOOSE_MING
	WAITING_FOR_GET_OTHER_MING
	WAITING_FOR_GET_SELF_MING
	END_ROUND
	WAITING_NEXT_ROUND
)

var press_button = make(chan bool)

func dealer_recv() string {
	msg, err := dealer.Recv()
	if err != nil {
		fmt.Println("Error receiving message:", err)
		os.Exit(1)
	}
	return strings.ToUpper(string(msg.Frames[0]))
}

func ming_without_chi_confirm(mingkind string, cardorkind string, selfming bool) {
	dialog.ShowConfirm("是否要鳴牌", fmt.Sprintf("Confirm to %s %s?", mingkind, cardorkind), func(confirm bool) {
		if confirm {
			mingkind = mingkind[:1] + strings.ToLower(mingkind[1:])
			cardorkind = strings.ToLower(cardorkind)
			sendmessage := fmt.Sprintf("%s %d %s", mingkind, pos.Pos[ID], cardorkind)
			dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))

			action = WAITING_FOR_GET_OTHER_MING

			selfming = false
			press_button <- true
			myCards.SortCard()
			updateGUI()

		} else {
			dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
			if selfming {
				action = DISCARD_CARD
			} else {
				action = WAITING_FOR_GET_OTHER_MING
			}
			selfming = false
			press_button <- true

		}

	}, fyne.CurrentApp().Driver().AllWindows()[0])

}

func ming_with_chi_confirm(mingkind string, cardorkind []string, card string) {
	dialog.ShowConfirm("是否要鳴牌", fmt.Sprintf("Confirm to %s?", mingkind), func(confirm bool) {

		if confirm {
			chi_button(cardorkind, card)
			return
		} else {
			dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
			action = WAITING_FOR_GET_OTHER_MING
			press_button <- true
			return
		}

	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

func chi_button(cardorkind []string, card string) {

	cardkind := string(card[0])
	number, _ := strconv.Atoi(string(card[1]))
	var dialogWindow *dialog.CustomDialog
	var button []*widget.Button

	for _, v := range cardorkind {
		switch v {
		case "0":
			button = append(button, widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", card, card, cardkind+strconv.Itoa(number+1), cardkind+strconv.Itoa(number+2)), func() {
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(fmt.Sprintf("Chi %d 0", pos.Pos[ID]))))

				action = WAITING_FOR_GET_OTHER_MING
				press_button <- true
				dialogWindow.Hide()
				updateGUI()
			}))

		case "1":
			button = append(button, widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", card, cardkind+strconv.Itoa(number-1), card, cardkind+strconv.Itoa(number+1)), func() {
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(fmt.Sprintf("Chi %d 1", pos.Pos[ID]))))

				action = WAITING_FOR_GET_OTHER_MING
				press_button <- true
				dialogWindow.Hide()
				updateGUI()
			}))

		case "2":
			button = append(button, widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", card, cardkind+strconv.Itoa(number-2), cardkind+strconv.Itoa(number-1), card), func() {
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(fmt.Sprintf("Chi %d 2", pos.Pos[ID]))))

				action = WAITING_FOR_GET_OTHER_MING
				press_button <- true
				dialogWindow.Hide()
				updateGUI()
			}))
		}
	}
	cancelbutton := widget.NewButton("Cancel", func() {
		dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))

		action = WAITING_FOR_GET_OTHER_MING
		press_button <- true
		dialogWindow.Hide()
		updateGUI()
	})
	var conobj []fyne.CanvasObject
	for _, v := range button {
		conobj = append(conobj, v)
	}
	conobj = append(conobj, cancelbutton)
	container := container.NewHBox(conobj...)
	dialogWindow = dialog.NewCustomWithoutButtons("Chi", container, fyne.CurrentApp().Driver().AllWindows()[0])
	dialogWindow.Show()

}

func put_button(selfming bool, nowdiscard string) []fyne.CanvasObject {
	var button []*widget.Button

	for _, v := range mingset {
		fmt.Println("Mingset:", v)
		mingslice := strings.Split(v, " ")
		cardorkind := strings.ToLower(mingslice[1])
		switch mingslice[0] {
		case "CHI":
			button = append(button, widget.NewButton("吃", func() {
				chi_button(mingslice[1:], nowdiscard)
			}))

		case "PONG":
			button = append(button, widget.NewButton("碰", func() {
				sendmessage := fmt.Sprintf("%s %d %s", "Pong", pos.Pos[ID], cardorkind)
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))

				action = WAITING_FOR_GET_OTHER_MING

				selfming = false
				press_button <- true
				myCards.SortCard()
				updateGUI()
			}))

		case "GANG":
			button = append(button, widget.NewButton("槓", func() {
				sendmessage := fmt.Sprintf("%s %d %s", "Gang", pos.Pos[ID], cardorkind)
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))

				action = WAITING_FOR_GET_OTHER_MING

				press_button <- true
				myCards.SortCard()
				updateGUI()
			}))

		case "HU":
			button = append(button, widget.NewButton("胡", func() {
				sendmessage := fmt.Sprintf("%s %d %s", "Hu", pos.Pos[ID], cardorkind)
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))

				action = WAITING_FOR_GET_OTHER_MING

				press_button <- true
				updateGUI()
			}))

		}
	}
	cancelbutton := widget.NewButton("取消", func() {
		dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
		if selfming {
			action = DISCARD_CARD
		} else {
			action = WAITING_FOR_GET_OTHER_MING
		}
		selfming = false
		press_button <- true
		updateGUI()
	})
	var conobj []fyne.CanvasObject
	for _, v := range button {

		conobj = append(conobj, v)
	}
	conobj = append(conobj, cancelbutton)

	return conobj
}

func changecolor() {

	switch pos_history[0] {
	case myPosition:
		grid.Objects[7].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	case (myPosition + 1) % 4:
		grid.Objects[5].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	case (myPosition + 2) % 4:
		grid.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	case (myPosition + 3) % 4:
		grid.Objects[3].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	if len(pos_history) == 2 {
		switch pos_history[1] {
		case myPosition:
			grid.Objects[7].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		case (myPosition + 1) % 4:
			grid.Objects[5].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		case (myPosition + 2) % 4:
			grid.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		case (myPosition + 3) % 4:
			grid.Objects[3].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		}
		pos_history = pos_history[1:]
	} else {
		switch pos_history[0] {
		case myPosition:
			grid.Objects[7].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		case (myPosition + 1) % 4:
			grid.Objects[5].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		case (myPosition + 2) % 4:
			grid.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		case (myPosition + 3) % 4:
			grid.Objects[3].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
		}
	}

	GUI.Refresh()
}

func behavior_handler() {
	go behavior()
}

var action int
var discardcard [4][]string
var throwcard = make(chan string)
var nowdiscard string
var selfming bool
var putnewcard bool
var mingset []string
var point [4]string
var now_pos int
var pos_history []int
var myPosition int

func behavior() {
	action = GAME_START_WAIT

	var getcard string

	for {
		fmt.Println("Action: ", action)
		switch action {
		case GAME_START_WAIT:
			msg := dealer_recv()
			mingbuttonlist.Hide()
			if msg == "GAME START" {
				action = GAME_START
			} else {
				action = GAME_START_WAIT
			}

		case GAME_START:
			mingbuttonlist.Hide()
			msg, _ = dealer.Recv()
			json.Unmarshal(msg.Frames[0], &pos)
			myCards = pos.Ma
			fmt.Println("My Cards:", myCards.Card)
			myCards.SortCard()
			fmt.Println("My Cards after sorting:", myCards.Card)
			myPosition = pos.Pos[ID]
			gamestart = true
			fmt.Println("My Position:", myPosition)
			// 逆時針標記其他玩家的位置
			playerPositions = make(map[int]string)
			now_pos = 0
			pos_history = append(pos_history, now_pos)
			changecolor()
			var relativePos [4]int

			// 逆時針標記其他玩家的位置
			relativePos[0] = myPosition           // 自己的位置
			relativePos[1] = (myPosition + 1) % 4 // 右邊的玩家
			relativePos[2] = (myPosition + 2) % 4 // 對面的玩家
			relativePos[3] = (myPosition + 3) % 4 // 左邊的玩家

			var AbsPos [4]string

			for i := range relativePos {
				j := (myPosition + i) % 4
				switch j {
				case 0:
					AbsPos[i] = "東"
				case 1:
					AbsPos[i] = "南"
				case 2:
					AbsPos[i] = "西"
				case 3:
					AbsPos[i] = "北"
				}
			}

			for playerID, position := range pos.Pos {
				playerPositions[position] = playerID
			}

			// 打印玩家位置
			grid.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerPositions[relativePos[2]] + AbsPos[2] // 對面的玩家
			grid.Objects[3].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerPositions[relativePos[3]] + AbsPos[3] // 左邊的玩家
			grid.Objects[5].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerPositions[relativePos[1]] + AbsPos[1] // 右邊的玩家
			grid.Objects[7].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerPositions[relativePos[0]] + AbsPos[0] // 自己的位置
			grid.Objects[4].(*fyne.Container).Objects[0].(*canvas.Text).Text = "東一局 0本場"

			updateGUI()

			action = DRAW_CARD

		case DRAW_CARD:
			msg := dealer_recv()
			if strings.Contains(msg, "DRAW") {
				action = WAITING_FOR_GET_SELF_MING
			} else {
				getcard = strings.ToLower(msg)
				newcc = getcard
				putnewcard = true
				//myCards.Card = append(myCards.Card, getcard)
				updateGUI()
				action = CHECK_SELF_MING
				GUI.Refresh()

			}

		case DISCARD_CARD:
			grid.Objects[7].(*fyne.Container).Objects[0].(*canvas.Text).Color = color.RGBA{R: 242, G: 76, B: 76, A: 255}
			GUI.Refresh()
			discardcard[myPosition] = append(discardcard[myPosition], <-throwcard)

		case CHECK_SELF_MING:
			msg := dealer_recv()
			if strings.Contains(msg, "NO") {
				action = DISCARD_CARD
			} else {
				selfming = true
				mingset = nil
				mingset = append(mingset, strings.Split(msg, ",")...)
				action = CHOOSE_SELF_MING
			}

		case CHOOSE_SELF_MING:

			mingbuttonlist.Objects = put_button(selfming, nowdiscard)
			for _, obj := range mingbuttonlist.Objects {
				if center, ok := obj.(*fyne.Container); ok {
					if text, ok := center.Objects[0].(*canvas.Text); ok {
						text.TextSize = 20 // 設置你想要的字體大小
					}
				}
			}
			mingbuttonlist.Show()
			GUI.Refresh()
			<-press_button
			//顯示按鈕給按
			//CANCEL action -> DISCARD_CARD
			//else -> WAITING_FOR_GET_OTHER_MING

		case WAITING_FOR_GET_DISCARD_CARD:
			msg := strings.ToLower(dealer_recv())
			msgslice := strings.Split(msg, " ")
			discardpos, _ := strconv.Atoi(msgslice[0])
			fmt.Println("Discard Position:", discardpos, "Discard Card:", msgslice[1])

			discardcard[discardpos] = append(discardcard[discardpos], msgslice[1])
			nowdiscard = msgslice[1]
			action = CHECK_MING

		case CHECK_MING:
			msg := dealer_recv()
			if strings.Contains(msg, "NO") {
				action = WAITING_FOR_GET_OTHER_MING
			} else {
				fmt.Println("CHECK_MING:", msg)
				mingset = nil
				mingset = append(mingset, strings.Split(msg, ",")...)

				action = CHOOSE_MING
			}

		case WAITING_FOR_GET_SELF_MING:
			mingbuttonlist.Hide()
			GUI.Refresh()
			msg, _ := dealer.Recv()
			fmt.Println("WAITING_FOR_GET_SELF_MING:", msg)
			if strings.Contains(string(msg.Frames[0]), "NO") {
				action = WAITING_FOR_GET_DISCARD_CARD
			} else {
				//顯示在桌上
				if strings.Contains(string(msg.Frames[0]), "Hu") {
					action = END_ROUND
					continue
				}
				action = WAITING_FOR_GET_DISCARD_CARD
			}

		case CHOOSE_MING:
			//顯示按鈕給按
			mingbuttonlist.Objects = put_button(selfming, nowdiscard)
			for _, obj := range mingbuttonlist.Objects {
				if center, ok := obj.(*fyne.Container); ok {
					if text, ok := center.Objects[0].(*canvas.Text); ok {
						text.TextSize = 20 // 設置你想要的字體大小
					}
				}
			}
			mingbuttonlist.Show()
			GUI.Refresh()
			<-press_button

		case WAITING_FOR_GET_OTHER_MING:
			//做相應處理
			mingbuttonlist.Hide()
			msg, _ := dealer.Recv()
			msgslice := strings.Split(string(msg.Frames[0]), " ")
			cardorkind := msgslice[1]
			fmt.Println("WAITING_FOR_GET_OTHER_MING:", msg)
			if strings.Contains(string(msg.Frames[0]), "NO") { //沒人鳴牌
				now_pos = (now_pos + 1) % 4
				pos_history = append(pos_history, now_pos)
				changecolor()
				action = DRAW_CARD
			} else {
				if msgslice[2] == ID { //如果是自己鳴，接下來要丟一張
					now_pos = myPosition
					pos_history = append(pos_history, now_pos)
					changecolor()
					if msgslice[0] == "Pong" {
						fmt.Println("PPong", cardorkind)
						myCards.removeCard(cardorkind)
						myCards.removeCard(cardorkind)
						myCards.SortCard()
						updateGUI()

					} else if msgslice[0] == "Gang" {
						if !selfming {
							myCards.removeCard(cardorkind)
							myCards.removeCard(cardorkind)
							myCards.removeCard(cardorkind)
						} else {
							myCards.removeCard(cardorkind)
						}

						myCards.SortCard()
						updateGUI()
						action = DRAW_CARD
						continue

					} else if msgslice[0] == "Chi" {
						cardkind := string(nowdiscard[0])
						number, _ := strconv.Atoi(string(nowdiscard[1]))
						switch cardorkind {
						case "0":
							myCards.removeCard(cardkind + strconv.Itoa(number+1))
							myCards.removeCard(cardkind + strconv.Itoa(number+2))
							myCards.SortCard()
							updateGUI()

						case "1":
							myCards.removeCard(cardkind + strconv.Itoa(number+1))
							myCards.removeCard(cardkind + strconv.Itoa(number-1))
							myCards.SortCard()
							updateGUI()

						case "2":
							myCards.removeCard(cardkind + strconv.Itoa(number-2))
							myCards.removeCard(cardkind + strconv.Itoa(number-1))
							myCards.SortCard()
							updateGUI()
						}
					} else if msgslice[0] == "Hu" {
						action = END_ROUND
						continue
					}
					action = DISCARD_CARD
				} else {
					now_pos = pos.Pos[msgslice[2]]
					pos_history = append(pos_history, now_pos)
					changecolor()
					if msgslice[0] == "Hu" {
						action = END_ROUND
					} else {
						action = WAITING_FOR_GET_DISCARD_CARD
					}
				}

			}

		case END_ROUND:
			selfming = false
			mingset = nil
			putnewcard = false
			pos_history = nil
			msg := dealer_recv()
			pointslice := strings.Split(msg, " ")
			copy(point[:], pointslice)
			action = WAITING_NEXT_ROUND

		case WAITING_NEXT_ROUND:
			msg := dealer_recv()
			if msg == "NEXT ROUND" {
				action = GAME_START
			} else {
				action = WAITING_NEXT_ROUND
			}

		}

	}

}
