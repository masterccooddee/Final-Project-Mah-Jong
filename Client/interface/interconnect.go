package ui

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"net"
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

type Position struct {
	Pos map[string]int
	Ma  mao
}

var myCards mao
var newcc string
var playerPositions map[int]string

// for roomchat
var pos Position
var dealer zmq4.Socket
var ID string

// for room
var conn net.Conn
var msg zmq4.Msg
var err error

var receivedMessage string
var RoomID string
var in string
var gamestart bool = false
var mingcardamount int

func rrecv() string {
	data := make([]byte, 4096)

	var num int
	num, _ = conn.Read(data)
	if num == 0 {
		//fmt.Println("\nConnection closed")
		os.Exit(1)
	}

	in2 := string(data[:num])
	in = strings.TrimSpace(in2)
	return in

}

func LORinterface(loginwindow *fyne.Window, openwindow *fyne.Window, chiwindow *fyne.Window) fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	received_content := canvas.NewText("", color.Black)
	received_content.TextSize = 12

	Send_Login_content := container.NewVBox(widget.NewButton("Login", func() {
		//fmt.Println("Logging:", input.Text)
		if input.Text == "" {
			//fmt.Println("Please enter a username")
			received_content.Text = "Please enter a username"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
			return
		} else {
			conn.Write([]byte("LOGIN " + input.Text))
		}

		ID = strings.TrimSpace(input.Text)
		input.SetText("")
		recv := rrecv()
		//fmt.Println(recv)
		recv = strings.Split(recv, " ")[0]
		if recv == "Welcome" {
			LoginSuccess = true
			dealer = zmq4.NewDealer(context.Background(), zmq4.WithID(zmq4.SocketIdentity(ID)))
			defer dealer.Close()

			err := dealer.Dial("tcp://localhost:7125")
			if err != nil {
				//fmt.Println("Error connecting dealer:", err)
				return
			}

			(*loginwindow).Close()
			(*openwindow).Show()

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

						rightPlayer := (myPosition + 1) % 4
						frontPlayer := (myPosition + 2) % 4
						leftPlayer := (myPosition + 3) % 4

						for playerID, position := range pos.Pos {
							playerPositions[position] = playerID
						}

						// 打印玩家位置
						grid.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerPositions[frontPlayer] // 對面的玩家
						grid.Objects[3].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerPositions[leftPlayer]  // 左邊的玩家
						grid.Objects[5].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerPositions[rightPlayer] // 右邊的玩家
						grid.Objects[7].(*fyne.Container).Objects[0].(*canvas.Text).Text = ID                           // 自己的位置

						msg, _ = dealer.Recv()
						fmt.Println("Received message:", string(msg.Frames[0]))
						newcc = string(msg.Frames[0])
					} else {
						cardthrow := strings.Split(receivedMessage, " ")

						fmt.Println("Position ", cardthrow[0], " throw ", cardthrow[1])
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
										return
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
										return
									}
									dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
								}, fyne.CurrentApp().Driver().AllWindows()[0])
							}
						case "CHI":
							kind := string(cardthrow[1][0])
							number, _ := strconv.Atoi(string(cardthrow[1][1]))
							sendmessage := fmt.Sprintf("Chi %d", pos.Pos[ID])
							var CheckPos [3]bool

							RightCard := kind + strconv.Itoa(number+1)
							LeftCard := kind + strconv.Itoa(number-1)
							Right2Card := kind + strconv.Itoa(number+2)
							Left2Card := kind + strconv.Itoa(number-2)

							button0 := widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", cardthrow[1], cardthrow[1], RightCard, Right2Card), func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage+" 0")))
								myCards.removeCard(RightCard)
								myCards.removeCard(Right2Card)
								myCards.SortCard()
								(*chiwindow).Close()
							})
							button1 := widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", cardthrow[1], LeftCard, cardthrow[1], RightCard), func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage+" 1")))
								myCards.removeCard(LeftCard)
								myCards.removeCard(RightCard)
								myCards.SortCard()
								(*chiwindow).Close()
							})
							button2 := widget.NewButton(fmt.Sprintf("Chi %s (%s %s %s)", cardthrow[1], Left2Card, LeftCard, cardthrow[1]), func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage+" 2")))
								myCards.removeCard(Left2Card)
								myCards.removeCard(LeftCard)
								myCards.SortCard()
								(*chiwindow).Close()
							})
							cancelbutton := widget.NewButton("Cancel", func() {
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("Cancel")))
								(*chiwindow).Close()
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
							(*chiwindow).SetContent(container)
							(*chiwindow).Show()
						default:
							fmt.Println("Received message:", string(msg.Frames[0]))
							newcc = strings.TrimSpace(string(msg.Frames[0]))
						}
					}
				}
			}
		} else {
			//fmt.Println("Login failed")
			received_content.Text = "Login failed"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
		}
	}))

	Send_Register_content := container.NewVBox(widget.NewButton("Register", func() {
		//fmt.Println("Registering:", input.Text)
		if input.Text == "" {
			//fmt.Println("Please enter a username")
			received_content.Text = "Please enter a username"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
			return
		} else {
			conn.Write([]byte("REG " + input.Text))
		}
		ID = strings.TrimSpace(input.Text)
		input.SetText("")
		recv := rrecv()
		//fmt.Println(recv)
		recv = strings.Split(recv, " ")[0]
		if recv == "Register" {
			LoginSuccess = true

			dealer = zmq4.NewDealer(context.Background(), zmq4.WithID(zmq4.SocketIdentity(ID)))
			defer dealer.Close()

			err := dealer.Dial("tcp://localhost:7125")
			if err != nil {
				//fmt.Println("Error connecting dealer:", err)
				return
			}

			(*loginwindow).Close()
			(*openwindow).Show()

			for {
				//fmt.Println("RoomID:", RoomID)
				if RoomID != "" {
					msg, err = dealer.Recv()
					if err != nil {
						//fmt.Println("Error receiving message:", err)
						break
					}
					receivedMessage = strings.ToUpper(string(msg.Frames[0]))
					mingcard := strings.Split(receivedMessage, " ")
					switch mingcard[0] {
					case "GAME":
						msg, _ = dealer.Recv()
						json.Unmarshal(msg.Frames[0], &pos)
						myCards = pos.Ma
						myCards.SortCard()
						myPosition := pos.Pos[ID]
						fmt.Println("My Position:", myPosition)
						// 逆時針標記其他玩家的位置
						playerPositions = make(map[int]string)
						for playerID, position := range pos.Pos {
							relativePosition := (position - myPosition + 4) % 4 // 逆時針計算相對位置
							playerPositions[relativePosition] = playerID
						}

						// 打印玩家位置
						for i := 0; i < 4; i++ {
							if playerID, ok := playerPositions[i]; ok {
								if i == 0 {
									fmt.Println("Myself:", playerID)
									grid.Objects[7].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerID
								} else {
									fmt.Printf("Player %d: %s\n", i, playerID)
									grid.Objects[7-i*2].(*fyne.Container).Objects[0].(*canvas.Text).Text = playerID
								}
							}
						}

					case "PONG":
						mingcard[1] = strings.ToLower(mingcard[1])
						fmt.Printf("PONG %s", mingcard[1])
						dialog.ShowConfirm("Confirm", fmt.Sprintf("Confirm to pong %s?", mingcard[1]), func(confirm bool) {
							if confirm {
								sendmessage := fmt.Sprintf("PONG %s %s", strconv.Itoa(pos.Pos[ID]), mingcard[1])
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))
								myCards.removeCard(mingcard[1])
								myCards.removeCard(mingcard[1])
								myCards.SortCard()
								return
							}
							dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("CANCEL")))
						}, fyne.CurrentApp().Driver().AllWindows()[0])
					case "CHI":
						fmt.Printf("CHI %s", mingcard[1])
						dialog.ShowConfirm("Confirm", fmt.Sprintf("Confirm to chi %s?", mingcard[1]), func(confirm bool) {
							if confirm {
								sendmessage := fmt.Sprintf("CHI %s %s", strconv.Itoa(pos.Pos[ID]), mingcard[1])
								dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(sendmessage)))
								return
							}
							dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte("CANCEL")))
						}, fyne.CurrentApp().Driver().AllWindows()[0])
					default:
						fmt.Println("Received message:", string(msg.Frames[0]))
						newcc = string(msg.Frames[0])
					}
				}
			}
		} else {
			//fmt.Println("ID already exist")
			received_content.Text = "ID already exist"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
		}
	}))

	form := widget.NewForm(widget.NewFormItem("Username", input))

	content := container.NewVBox(form, Send_Login_content, Send_Register_content, received_content)

	return content
}

func interconnect() {
	conn, err = net.Dial("tcp", ":8080")
	defer conn.Close()

	if err != nil {
		//fmt.Println("Error dialing", err)
		return
	}

	//go serverexit(conn)

	for {
		//fmt.Print("Enter text: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		out := strings.Split(text, " ")
		if out[0] == "LOGIN" {
			out[1] = strings.TrimSpace(out[1])
			ID = out[1]
			//fmt.Println("LOGIN")

		}
		if strings.TrimSpace(out[0]) == "CHG" {
			//fmt.Println("RoomID: ", RoomID)
			for {
				//fmt.Print("Enter text2: ")
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(text)))
			}
		}

	}
}
