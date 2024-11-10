package ui

import (
	"bufio"
	"context"
	"encoding/json"
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
}

var conn net.Conn
var err error
var RoomID string
var msg zmq4.Msg
var in string

/* func serverexit(conn net.Conn) {
	data := make([]byte, 4096)

	var num int
	for {

		num, _ = conn.Read(data)
		if num == 0 {
			fmt.Println("\nConnection closed")
			os.Exit(1)
		}

		in2 := string(data[:num])
		in = strings.TrimSpace(in2)
		out := strings.Split(in, " ")
		if out[0] == "True" {
			if out[1] == "Room" {
				RoomID = out[2]
			}
		}

		fmt.Printf("\rFrom Server -> %s\n", string(data[:num]))

	}

} */

func LORinterface() fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	received_content := container.NewVBox()
	Send_Login_content := container.NewVBox(widget.NewButton("Login", func() {
		fmt.Println("Logging:", input.Text)
		conn.Write([]byte("LOGIN " + input.Text))
		input.SetText("")
	}))

	Send_Register_content := container.NewVBox(widget.NewButton("Register", func() {
		fmt.Println("Registering:", input.Text)
		conn.Write([]byte("REG " + input.Text))
		input.SetText("")
	}))

	form := widget.NewForm(widget.NewFormItem("Username", input))

	content := container.NewVBox(form, Send_Login_content, Send_Register_content, received_content)

	go func() {
		for {
			data := make([]byte, 4096)
			num, err := conn.Read(data)
			if err != nil {
				fmt.Println("Error reading:", err)
				return
			}
			if num > 0 {
				message := string(data[:num])
				message = strings.TrimSpace(message)
				received_text := canvas.NewText(message, color.Black)
				msg_split := strings.Split(message, " ")

				received_content.RemoveAll()
				received_content.Add(received_text)
				received_content.Refresh()

				if msg_split[0] == "Register" || msg_split[0] == "Welcome" {
					LoginSuccess = true
					return
				}
			}
		}
	}()

	return content
}

func interconnect() {
	conn, err = net.Dial("tcp", "localhost:8080")
	defer conn.Close()

	if err != nil {
		fmt.Println("Error dialing", err)
		return
	}

	//go serverexit(conn)

	var dealer zmq4.Socket
	var ID string
	for {
		fmt.Print("Enter text: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		out := strings.Split(text, " ")
		if out[0] == "LOGIN" {
			out[1] = strings.TrimSpace(out[1])
			ID = out[1]
			dealer = zmq4.NewDealer(context.Background(), zmq4.WithID(zmq4.SocketIdentity(out[1])))
			defer dealer.Close()

			err := dealer.Dial("tcp://localhost:7125")
			if err != nil {
				fmt.Println("Error connecting dealer:", err)
				return
			}
			fmt.Println(out[1])

			go func() {

				// DEALER 接收消息
				msg, err = dealer.Recv()
				if err != nil {
					fmt.Println("Error receiving message:", err)
					return
				}
				fmt.Println("Received message:", string(msg.Frames[0]))
				msg, _ = dealer.Recv()
				var pos Position
				json.Unmarshal(msg.Frames[0], &pos)
				fmt.Println(pos.Pos)
				fmt.Println(pos.Pos[ID])

			}()
		}
		if strings.TrimSpace(out[0]) == "CHG" {
			fmt.Println("RoomID: ", RoomID)
			for {
				fmt.Print("Enter text2: ")
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(text)))
			}
		}

	}
}
