package ui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-zeromq/zmq4"
)

func makeRoomInterface(received_content **canvas.Text) fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	(*received_content) = canvas.NewText("", color.Black)
	(*received_content).TextSize = 12
	//roommake command
	roommake := container.NewVBox(widget.NewButton("Make New Room", func() {
		//fmt.Println("Make New Room")
		conn.Write([]byte("ROOM MAKE"))
		input.SetText("")
		recv := rrecv()
		//fmt.Println(recv)
		msg := strings.Split(recv, " ")
		if msg[0] == "True" {
			RoomID = msg[2]
			inroom = true
			//fmt.Println("Make RoomID: " + RoomID + " and join")
			(*received_content).Text = "Make RoomID: " + RoomID
			(*received_content).Color = color.Black
			(*received_content).Refresh()
			//(*w).Close()
		} else {
			//fmt.Println("You are already in a room")
			(*received_content).Text = "You are already in a room"
			(*received_content).Color = color.RGBA{255, 0, 0, 255}
			(*received_content).Refresh()
		}
	}))

	//roomjoin command
	roomjoin := container.NewVBox(widget.NewButton("Join Room", func() {
		RoomID = input.Text
		if RoomID == "" {
			(*received_content).Text = "Please enter a room ID"
			(*received_content).Color = color.RGBA{255, 0, 0, 255}
			(*received_content).Refresh()
			return
		} else {
			conn.Write([]byte("ROOM JOIN " + RoomID))
		}

		input.SetText("")
		recv := rrecv()
		//fmt.Println(recv)
		msg := strings.Split(recv, " ")
		if msg[0] == "True" {
			inroom = true
			msg[2] = strings.TrimSpace(msg[2])
			RoomID = msg[2]
			//fmt.Println("Join RoomID: " + RoomID)
			(*received_content).Text = "Successfully Join RoomID: " + RoomID
			(*received_content).Color = color.Black
			(*received_content).Refresh()
			//(*w).Close()
		} else if msg[0] == "False" {
			if msg[1] == "command" {
				//fmt.Println("False Command")
				(*received_content).Text = "False Command"
				(*received_content).Color = color.RGBA{255, 0, 0, 255}
				(*received_content).Refresh()
			} else if msg[3] == "full" {
				//fmt.Println("Room is full")
				(*received_content).Text = "Room is full"
				(*received_content).Color = color.RGBA{255, 0, 0, 255}
				(*received_content).Refresh()
			} else {
				//fmt.Println("Room not exist")
				(*received_content).Text = "Room not exist"
				(*received_content).Color = color.RGBA{255, 0, 0, 255}
				(*received_content).Refresh()
			}
		} else {
			//fmt.Println("You are already in a room")
			(*received_content).Text = "You are already in a room"
			(*received_content).Color = color.RGBA{255, 0, 0, 255}
			(*received_content).Refresh()
		}
	}))

	//roomFind command
	roomfind := container.NewVBox(widget.NewButton("Find Room", func() {
		//fmt.Println("Find Room")
		conn.Write([]byte("ROOM FIND"))
		recv := rrecv()
		//fmt.Println(recv)
		msg := strings.Split(recv, " ")
		if msg[0] == "True" {
			//fmt.Println("msg[2] " + msg[2])
			msg[2] = strings.TrimSpace(msg[2])
			RoomID = msg[2]
			(*received_content).Text = "Find RoomID: " + msg[2]
			(*received_content).Color = color.Black
			(*received_content).Refresh()
			//(*w).Close()
		} else {
			//fmt.Println("You are already in a room")
			(*received_content).Text = "You are already in a room"
			(*received_content).Color = color.RGBA{255, 0, 0, 255}
			(*received_content).Refresh()
		}
	}))

	form := widget.NewForm(widget.NewFormItem("Room ID", input))
	content := container.NewVBox(form, roommake, roomjoin, roomfind, (*received_content))

	return content
}

func makeRoomChatInterface() fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	received_content := canvas.NewText("", color.Black)
	received_content.TextSize = 12

	//roomchat command
	roomchatbutton := container.NewVBox(widget.NewButton("Send Message", func() {
		if RoomID == "" {
			received_content.Text = "You are not in a room"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
		} else {
			//roomchat command
			dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(input.Text)))
		}
	}))

	form := widget.NewForm(widget.NewFormItem("Message", input))

	content := container.NewVBox(form, roomchatbutton, received_content)

	return content
}

func makeBannner_top(received_content **canvas.Text) fyne.CanvasObject {
	toolbar := widget.NewToolbar()
	toolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.GridIcon(), func() {

			leaveRoomAction := fyne.NewMenuItem("Leave Room", func() {
				// 傳送 LEAVE 命令給伺服器
				conn.Write([]byte("ROOM LEAVE"))
				recv := rrecv()
				//fmt.Println(recv)
				msg := strings.Split(recv, " ")[0]
				if msg == "True" {
					inroom = false
					(*received_content).Text = "Leave Room " + RoomID
					(*received_content).Color = color.Black
					RoomID = ""
					(*received_content).Refresh()
					//fmt.Println("Leave Room")
				} else {
					//fmt.Println("You are not in a room")
					(*received_content).Text = "You are not in a room"
					(*received_content).Color = color.RGBA{255, 0, 0, 255}
					(*received_content).Refresh()
				}
			})

			roomAction := fyne.NewMenuItem("Room", func() {
				w := fyne.CurrentApp().NewWindow("Room")
				w.Resize(fyne.NewSize(300, 200))
				w.SetContent(makeRoomInterface(received_content))
				w.CenterOnScreen()
				w.Show()
			})

			roomchatAction := fyne.NewMenuItem("Room Chat", func() {
				w := fyne.CurrentApp().NewWindow("Room Chat")
				w.Resize(fyne.NewSize(300, 300))
				w.SetContent(makeRoomChatInterface())
				w.CenterOnScreen()
				w.Show()
			})

			menu := fyne.NewMenu("", leaveRoomAction, roomAction, roomchatAction)
			popUpMenu := widget.NewPopUpMenu(menu, fyne.CurrentApp().Driver().CanvasForObject(toolbar))
			popUpMenu.ShowAtPosition(fyne.CurrentApp().Driver().AbsolutePositionForObject(toolbar).Add(fyne.NewPos(0, toolbar.Size().Height)))
		}),
	)

	logo1 := canvas.NewImageFromResource(resource7Png)
	logo1.FillMode = canvas.ImageFillContain

	return container.NewStack(toolbar, container.NewPadded(logo1))
}
