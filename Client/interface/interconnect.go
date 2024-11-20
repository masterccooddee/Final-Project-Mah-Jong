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

func LORinterface(loginwindow *fyne.Window, openwindow *fyne.Window) fyne.CanvasObject {
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
						case "PONG":
							if len(mingcard) == 2 {
								mingcard[1] = strings.ToLower(mingcard[1])
								fmt.Printf("PONG %s", mingcard[1])

								dialog.ShowConfirm("Confirm", fmt.Sprintf("Confirm to pong %s?", mingcard[1]), func(confirm bool) {
									if confirm {
										mingcardamount++
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
							fmt.Printf("CHI %s", mingcard[1])
							dialog.ShowConfirm("Confirm", fmt.Sprintf("Confirm to chi %s?", mingcard[1]), func(confirm bool) {
								if confirm {
									sendmessage := fmt.Sprintf("CHI %d %s", pos.Pos[ID], mingcard[1])
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
