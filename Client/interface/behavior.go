package ui

import (
	"encoding/json"
	"fmt"
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

func behavior_handler() {
	go behavior()
}

func Behavior() {
	for {
		//fmt.Println("RoomID:", RoomID)
		if RoomID != "" {
			msg, err = dealer.Recv()
			if err != nil {
				//fmt.Println("Error receiving message:", err)
				return
			}
			receivedMessage = strings.ToUpper(string(msg.Frames[0]))
			var myPosition int
			if receivedMessage == "GAME START" {
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

				msg, _ = dealer.Recv()
				fmt.Println("Received message:", string(msg.Frames[0]))
				newcc = string(msg.Frames[0])
				updateGUI()
			} else {
				cardthrow := strings.Split(receivedMessage, " ")

				//fmt.Println("Position ", cardthrow[0], " throw ", cardthrow[1])
				msg, _ = dealer.Recv()
				fmt.Println("Received message:", string(msg.Frames[0]))
				mingcard := strings.Split(string(msg.Frames[0]), " ")
				mingcard[0] = strings.ToUpper(mingcard[0])

				switch mingcard[0] {
				case "GANG":
					if len(mingcard) == 2 {
						mingcard[1] = strings.ToLower(mingcard[1])
						fmt.Printf("GANG %s", mingcard[1])

						dialog.ShowConfirm("Confirm", fmt.Sprintf("Confirm to Gang %s?", mingcard[1]), func(confirm bool) {
							if confirm {
								sendmessage := fmt.Sprintf("Gang %d %s", pos.Pos[ID], mingcard[1])
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))
								myCards.removeCard(mingcard[1])
								myCards.removeCard(mingcard[1])
								myCards.removeCard(mingcard[1])
								myCards.SortCard()
								newcc = "Finish Gang"
								updateGUI()
							}
							dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
						}, fyne.CurrentApp().Driver().AllWindows()[0])
					}
				case "PONG":
					if len(mingcard) == 2 {
						mingcard[1] = strings.ToLower(mingcard[1])
						fmt.Printf("PONG %s", mingcard[1])

						dialog.ShowConfirm("Confirm", fmt.Sprintf("Confirm to pong %s?", mingcard[1]), func(confirm bool) {
							if confirm {
								sendmessage := fmt.Sprintf("Pong %d %s", pos.Pos[ID], mingcard[1])
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))
								myCards.removeCard(mingcard[1])
								myCards.removeCard(mingcard[1])
								myCards.SortCard()
								newcc = "Finish Pong"
								mingcardamount++
								fmt.Println("{mingcardamount}", mingcardamount)
								updateGUI()
							} else {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
							}
						}, fyne.CurrentApp().Driver().AllWindows()[0])
					} else if len(mingcard) == 4 {

					}
				case "CHI":
					var dialogWindow *dialog.CustomDialog

					dialog.ShowConfirm("Confirm", "Confirm to chi?", func(confirm bool) {
						kind := string(cardthrow[1][0])
						number, _ := strconv.Atoi(string(cardthrow[1][1]))
						sendmessage := fmt.Sprintf("Chi %d", pos.Pos[ID])
						var CheckPos [3]bool

						RightCard = kind + strconv.Itoa(number+1)
						LeftCard = kind + strconv.Itoa(number-1)
						Right2Card = kind + strconv.Itoa(number+2)
						Left2Card = kind + strconv.Itoa(number-2)

						if confirm {
							button0 := widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", cardthrow[1], cardthrow[1], RightCard, Right2Card), func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage+" 0")))
								RightCard = strings.TrimSpace(RightCard)
								Right2Card = strings.TrimSpace(Right2Card)
								myCards.removeCard(RightCard)
								myCards.removeCard(Right2Card)
								myCards.SortCard()
								newcc = "Finish Chi 0"
								mingcardamount++
								fmt.Println("{mingcardamount}", mingcardamount)
								dialogWindow.Hide()
								updateGUI()
							})
							button1 := widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", cardthrow[1], LeftCard, cardthrow[1], RightCard), func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage+" 1")))
								LeftCard = strings.TrimSpace(LeftCard)
								RightCard = strings.TrimSpace(RightCard)
								myCards.removeCard(LeftCard)
								myCards.removeCard(RightCard)
								myCards.SortCard()
								newcc = "Finish Chi 1"
								mingcardamount++
								fmt.Println("{mingcardamount}", mingcardamount)
								dialogWindow.Hide()
								updateGUI()
							})
							button2 := widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", cardthrow[1], Left2Card, LeftCard, cardthrow[1]), func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage+" 2")))
								Left2Card = strings.TrimSpace(Left2Card)
								LeftCard = strings.TrimSpace(LeftCard)
								myCards.removeCard(Left2Card)
								myCards.removeCard(LeftCard)
								myCards.SortCard()
								newcc = "Finish Chi 2"
								mingcardamount++
								fmt.Println("{mingcardamount}", mingcardamount)
								dialogWindow.Hide()
								updateGUI()
							})
							cancelbutton := widget.NewButton("Cancel", func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
								newcc = "Cancel Chi"
								dialogWindow.Hide()
								updateGUI()
							})

							for _, MingPos := range mingcard[1:] {
								pos, _ := strconv.Atoi(MingPos)
								CheckPos[pos] = true
							}

							for i, v := range CheckPos {
								if !v {
									switch i {
									case 0:
										button0.Hide()
									case 1:
										button1.Hide()
									case 2:
										button2.Hide()
									}
								}
							}
							container := container.NewHBox(button0, button1, button2, cancelbutton)
							dialogWindow = dialog.NewCustomWithoutButtons("Chi", container, fyne.CurrentApp().Driver().AllWindows()[0])
							dialogWindow.Show()
						} else {
							dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
						}
					}, fyne.CurrentApp().Driver().AllWindows()[0])

				default:
					fmt.Println("Received message:", string(msg.Frames[0]))
					newcc = strings.TrimSpace(string(msg.Frames[0]))
				}
				updateGUI()
			}
			updateGUI()
		}

	}
}

var action int
var discardcard [4][]string
var throwcard = make(chan string)
var nowdiscard string
var selfming bool
var putnewcard bool

func behavior() {
	action = GAME_START_WAIT
	var myPosition int
	var getcard string
	var mingchoose string

	for {
		fmt.Println("Action: ", action)
		switch action {
		case GAME_START_WAIT:
			msg := dealer_recv()
			if msg == "GAME START" {
				action = GAME_START
			} else {
				action = GAME_START_WAIT
			}
		case GAME_START:
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
			discardcard[myPosition] = append(discardcard[myPosition], <-throwcard)

		case CHECK_SELF_MING:
			msg := dealer_recv()
			if strings.Contains(msg, "NO") {
				action = DISCARD_CARD
			} else {
				mingchoose = msg
				selfming = true
				action = CHOOSE_SELF_MING
			}

		case CHOOSE_SELF_MING:
			var choose [5]bool // 0: PONG, 1: CHI, 2: GANG, 3: HU, 4: CANCEL
			if strings.Contains(mingchoose, "GANG") {
				choose[2] = true
				ming_without_chi_confirm("GANG", getcard, selfming)

			} else if strings.Contains(mingchoose, "HU") {
				choose[3] = true
				ming_without_chi_confirm("HU", getcard, selfming)
				action = END_ROUND
			}

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
				ming := strings.Split(msg, ",")
				for _, v := range ming {
					mingslice := strings.Split(v, " ")
					mingkind := mingslice[0]
					cardorkind := mingslice[1]
					switch mingkind {
					case "PONG":
						ming_without_chi_confirm("PONG", cardorkind, selfming)

					case "GANG":
						ming_without_chi_confirm("GANG", cardorkind, selfming)

					case "CHI":
						ming_with_chi_confirm("CHI", mingslice[1:], nowdiscard)

					}
				}
				<-press_button
				//action = CHOOSE_MING
			}

		case WAITING_FOR_GET_SELF_MING:
			msg := dealer_recv()
			if msg == "NO" {
				action = WAITING_FOR_GET_DISCARD_CARD
			} else {
				//顯示在桌上
				fmt.Println("WAITING_FOR_GET_SELF_MING:", msg)
				action = WAITING_FOR_GET_DISCARD_CARD
			}

		case CHOOSE_MING:
			//顯示按鈕給按

			action = WAITING_FOR_GET_OTHER_MING

		case WAITING_FOR_GET_OTHER_MING:
			//做相應處理
			msg := dealer_recv()
			msgslice := strings.Split(msg, " ")
			cardorkind := strings.ToLower(msgslice[1])
			fmt.Println("WAITING_FOR_GET_OTHER_MING:", msg)
			if strings.Contains(msg, "NO") { //沒人鳴牌
				action = DRAW_CARD
			} else {
				if msgslice[2] == strings.ToUpper(ID) { //如果是自己鳴，接下來要丟一張
					if strings.Contains(msg, "PONG") {
						fmt.Println("PPong", cardorkind)
						myCards.removeCard(cardorkind)
						myCards.removeCard(cardorkind)
						myCards.SortCard()
						updateGUI()

						mingcardamount += 2
					} else if strings.Contains(msg, "GANG") {
						myCards.removeCard(cardorkind)
						myCards.removeCard(cardorkind)
						myCards.removeCard(cardorkind)
						myCards.SortCard()
						updateGUI()

						mingcardamount += 3
					} else if strings.Contains(msg, "CHI") {
						cardkind := string(nowdiscard[0])
						number, _ := strconv.Atoi(string(nowdiscard[1]))
						switch cardorkind {
						case "0":
							myCards.removeCard(cardkind + strconv.Itoa(number+1))
							myCards.removeCard(cardkind + strconv.Itoa(number+2))
							myCards.SortCard()
							updateGUI()
							mingcardamount += 2
							fmt.Println("{mingcardamount}", mingcardamount)

						case "1":
							myCards.removeCard(cardkind + strconv.Itoa(number+1))
							myCards.removeCard(cardkind + strconv.Itoa(number-1))
							myCards.SortCard()
							updateGUI()
							mingcardamount += 2
							fmt.Println("{mingcardamount}", mingcardamount)

						case "2":
							myCards.removeCard(cardkind + strconv.Itoa(number-2))
							myCards.removeCard(cardkind + strconv.Itoa(number-1))
							myCards.SortCard()
							updateGUI()
							mingcardamount += 2
							fmt.Println("{mingcardamount}", mingcardamount)
						}
					}
					action = DISCARD_CARD
				} else {
					action = WAITING_FOR_GET_DISCARD_CARD
				}

			}

		}

	}

}
