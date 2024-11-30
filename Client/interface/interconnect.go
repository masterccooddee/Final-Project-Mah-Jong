package ui

import (
	"context"
	"fmt"
	"image/color"
	"net"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

var mingselect *fyne.Container

// for room
var conn net.Conn
var msg zmq4.Msg
var err error

var receivedMessage string
var RoomID string
var in string
var gamestart bool = false
var mingcardamount int = 0

// for chi
var RightCard string
var LeftCard string
var Right2Card string
var Left2Card string

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

			err := dealer.Dial("tcp://localhost:7125")
			if err != nil {
				//fmt.Println("Error connecting dealer:", err)
				return
			}

			(*loginwindow).Close()
			(*openwindow).Show()
			fmt.Println("ID:", ID)
			behavior_handler()
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

			err := dealer.Dial("tcp://localhost:7125")
			if err != nil {
				//fmt.Println("Error connecting dealer:", err)
				return
			}

			(*loginwindow).Close()
			(*openwindow).Show()
			fmt.Println("ID:", ID)
			behavior_handler()
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

	if err != nil {
		//fmt.Println("Error dialing", err)
		return
	}

}
