package ui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var inroom bool = false

func makeBannner_top() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.GridIcon(), func() {}),
	)

	logo1 := canvas.NewImageFromResource(resourceLogo1Png)
	logo1.FillMode = canvas.ImageFillContain

	logo2 := canvas.NewImageFromResource(resourceLogo2Png)
	logo2.FillMode = canvas.ImageFillContain

	logo3 := canvas.NewImageFromResource(resourceLogo3Png)
	logo3.FillMode = canvas.ImageFillContain

	return container.NewStack(toolbar, container.NewPadded(logo1))
}

func makeBanner_bottom_bar() fyne.CanvasObject {

	card := canvas.NewImageFromResource(resourceWordfivePng)
	card.FillMode = canvas.ImageFillContain
	return container.NewStack(container.NewPadded(card))
}

func makeRoomInterface() fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	received_content := canvas.NewText("", color.Black)
	received_content.TextSize = 12

	roommake := container.NewVBox(widget.NewButton("Make New Room", func() {
		fmt.Println("Make New Room")
		conn.Write([]byte("ROOM MAKE"))
		recv := rrecv()
		fmt.Println(recv)
		msg := strings.Split(recv, " ")
		if msg[0] == "True" {
			RoomID = msg[2]
			inroom = true
			fmt.Println("Make RoomID: " + RoomID + " and join")
			received_content.Text = "Make RoomID: " + RoomID
			received_content.Refresh()
		} else {
			fmt.Println("You are already in a room")
			received_content.Text = "You are already in a room"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
		}
	}))
	roomjoin := container.NewVBox(widget.NewButton("Join Room", func() {
		fmt.Println("Join Room")
		if RoomID == "" {
			received_content.Text = "Please enter a room ID"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()

		} else {
			conn.Write([]byte("ROOM JOIN " + RoomID))
		}

		recv := rrecv()
		fmt.Println(recv)
		msg := strings.Split(recv, " ")
		if msg[0] == "True" {
			inroom = true
			RoomID = msg[2]
			fmt.Println("Join RoomID: " + RoomID)
			received_content.Text = "Succesfully Join RoomID: " + RoomID
			received_content.Refresh()
		} else if msg[0] == "False" {
			if msg[1] == "command" {
				fmt.Println("False Command")
				received_content.Text = "False Command"
				received_content.Color = color.RGBA{255, 0, 0, 255}
				received_content.Refresh()
			} else if msg[3] == "full" {
				fmt.Println("Room is full")
				received_content.Text = "Room is full"
				received_content.Color = color.RGBA{255, 0, 0, 255}
				received_content.Refresh()
			} else {
				fmt.Println("Room not exist")
				received_content.Text = "Room not exist"
				received_content.Color = color.RGBA{255, 0, 0, 255}
				received_content.Refresh()
			}
		} else {
			fmt.Println("You are already in a room")
			received_content.Text = "You are already in a room"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
		}
	}))

	/* roomleave := container.NewVBox(widget.NewButton("Leave Room", func() {
		fmt.Println("Leave Room")
		conn.Write([]byte("ROOM LEAVE " + RoomID))
		recv := rrecv()
		fmt.Println(recv)
		msg := strings.Split(recv, " ")
		if msg[0] == "True" {
			inroom = false
			RoomID = "-1"
			fmt.Println("Leave Room")
			received_content.Text = "Leave Room"
			received_content.Refresh()
		} else {
			fmt.Println("You are not in a room")
			received_content.Text = "You are not in a room"
			received_content.Color = color.RGBA{255, 0, 0, 255}
			received_content.Refresh()
		}
	})) */
	form := widget.NewForm(widget.NewFormItem("Room ID", input))
	content := container.NewVBox(form, roommake, roomjoin, received_content)

	return content

}

func makeGUI() fyne.CanvasObject {
	top := makeBannner_top()
	top_bar := widget.NewLabel("Top")
	left_bar := widget.NewLabel("Left")
	right_bar := widget.NewLabel("Right")
	bottom_bar := [13]fyne.CanvasObject{
		makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(), makeBanner_bottom_bar(),
	}

	content := makeRoomInterface()

	dividers := [5]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator(),
	}
	//return container.NewBorder(makeBannner(), nil, left, right, content)
	objs := []fyne.CanvasObject{content, top, top_bar, left_bar, right_bar, bottom_bar[0], bottom_bar[1], bottom_bar[2], bottom_bar[3], bottom_bar[4], bottom_bar[5], bottom_bar[6], bottom_bar[7], bottom_bar[8], bottom_bar[9], bottom_bar[10], bottom_bar[11], bottom_bar[12],
		dividers[0], dividers[1], dividers[2], dividers[3], dividers[4]}
	return container.New(NewFysionLayout(top, top_bar, left_bar, right_bar, content, bottom_bar, dividers), objs...)
}
